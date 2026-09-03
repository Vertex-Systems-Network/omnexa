package events

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var errOutboxTestQueryUnsupported = errors.New("outbox test transaction query unsupported")

type outboxTestTx struct{}

func (*outboxTestTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (*outboxTestTx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errOutboxTestQueryUnsupported
}

func (*outboxTestTx) QueryRow(context.Context, string, ...any) pgx.Row {
	return nil
}

type outboxTestStore struct {
	enqueueResult EnqueueResult
	enqueueErr    error
	pending       []OutboxRecord
	pendingErr    error
	markResult    MarkResult
	markErr       error

	enqueueCalls int
	pendingCalls int
	markCalls    int
	lastTx       OutboxTransaction
	lastEnqueued OutboxRecord
	lastScope    OutboxScope
	lastEventID  EventID
	lastRevision uint64
}

func (store *outboxTestStore) Enqueue(_ context.Context, tx OutboxTransaction, record OutboxRecord) (EnqueueResult, error) {
	store.enqueueCalls++
	store.lastTx = tx
	store.lastEnqueued = record
	return store.enqueueResult, store.enqueueErr
}

func (store *outboxTestStore) Pending(_ context.Context, scope OutboxScope, _ uint32) ([]OutboxRecord, error) {
	store.pendingCalls++
	store.lastScope = scope
	return append([]OutboxRecord(nil), store.pending...), store.pendingErr
}

func (store *outboxTestStore) MarkPublished(_ context.Context, scope OutboxScope, eventID EventID, revision uint64) (MarkResult, error) {
	store.markCalls++
	store.lastScope = scope
	store.lastEventID = eventID
	store.lastRevision = revision
	return store.markResult, store.markErr
}

func TestNewOutboxRecordPreservesCanonicalEnvelopeAndCopiesPayload(t *testing.T) {
	envelope := testEnvelope(t)
	scope := testOutboxScope(envelope)
	record, err := NewOutboxRecord(scope, envelope)
	if err != nil {
		t.Fatalf("NewOutboxRecord() error = %v", err)
	}
	if record.State != OutboxPending || record.Revision != 1 {
		t.Fatalf("initial outbox state = (%v, %d), want pending revision 1", record.State, record.Revision)
	}
	if record.Envelope.ID != envelope.ID || record.Envelope.Source != envelope.Source || record.Envelope.TenantID != envelope.TenantID || record.Envelope.CorrelationID != envelope.CorrelationID || record.Envelope.CausationID != envelope.CausationID || record.Envelope.Classification != envelope.Classification {
		t.Fatal("outbox record changed canonical envelope identity or context")
	}
	before := append([]byte(nil), record.Envelope.Data...)
	envelope.Data[0] = 'X'
	if string(record.Envelope.Data) != string(before) {
		t.Fatal("outbox record retained caller payload alias")
	}
}

func TestNewOutboxRecordRejectsOwnerOrTenantRebinding(t *testing.T) {
	envelope := testEnvelope(t)
	scope := testOutboxScope(envelope)
	scope.Owner = Producer("urn:omnexa:module:other.owner")
	_, err := NewOutboxRecord(scope, envelope)
	assertFailureCode(t, err, codeOutboxConflict)

	scope = testOutboxScope(envelope)
	scope.TenantID = ""
	_, err = NewOutboxRecord(scope, envelope)
	assertFailureCode(t, err, codeOutboxConflict)
}

func TestEnqueueOutboxUsesCallerTransactionWithoutOpeningAnother(t *testing.T) {
	envelope := testEnvelope(t)
	scope := testOutboxScope(envelope)
	tx := &outboxTestTx{}
	store := &outboxTestStore{enqueueResult: OutboxEnqueued}

	result, err := EnqueueOutbox(context.Background(), tx, store, scope, envelope)
	if err != nil {
		t.Fatalf("EnqueueOutbox() error = %v", err)
	}
	if result != OutboxEnqueued || store.enqueueCalls != 1 || store.lastTx != tx {
		t.Fatalf("enqueue did not use the exact caller transaction: result=%v calls=%d exactTx=%v", result, store.enqueueCalls, store.lastTx == tx)
	}
	if store.lastEnqueued.Envelope.ID != envelope.ID || store.lastEnqueued.State != OutboxPending {
		t.Fatal("enqueue store received mutated identity or non-pending state")
	}
}

func TestEnqueueOutboxSeparatesIdempotentRetryFromConflictingIdentity(t *testing.T) {
	envelope := testEnvelope(t)
	scope := testOutboxScope(envelope)
	tx := &outboxTestTx{}

	unchanged := &outboxTestStore{enqueueResult: OutboxUnchanged}
	result, err := EnqueueOutbox(context.Background(), tx, unchanged, scope, envelope)
	if err != nil || result != OutboxUnchanged {
		t.Fatalf("same-event retry = (%v, %v), want unchanged without error", result, err)
	}

	conflict := &outboxTestStore{enqueueResult: OutboxEnqueueConflict}
	result, err = EnqueueOutbox(context.Background(), tx, conflict, scope, envelope)
	if result != OutboxEnqueueConflict {
		t.Fatalf("conflicting EventID result = %v, want conflict", result)
	}
	assertFailureCode(t, err, codeOutboxConflict)
}

func TestEnqueueOutboxSanitizesPersistenceFailures(t *testing.T) {
	const privateCause = "postgres host and restricted payload must remain private"
	envelope := testEnvelope(t)
	store := &outboxTestStore{enqueueErr: errors.New(privateCause)}
	_, err := EnqueueOutbox(context.Background(), &outboxTestTx{}, store, testOutboxScope(envelope), envelope)
	assertFailureCode(t, err, codeOutboxStoreFailed)
	if err != nil && containsPrivateText(err.Error(), privateCause) {
		t.Fatalf("enqueue failure leaked provider details: %v", err)
	}
}

func TestRelayPendingPublishesCanonicalEnvelopeThenMarksExactRevision(t *testing.T) {
	envelope := testEnvelope(t)
	record, err := NewOutboxRecord(testOutboxScope(envelope), envelope)
	if err != nil {
		t.Fatalf("NewOutboxRecord() error = %v", err)
	}
	store := &outboxTestStore{pending: []OutboxRecord{record}, markResult: OutboxMarkedPublished}
	var received Envelope
	publisher, err := NewPublisher(func(_ context.Context, candidate Envelope) error {
		received = candidate
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	relay, err := NewOutboxRelay(store, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() error = %v", err)
	}

	report, relayErr := relay.RelayPending(context.Background(), record.Scope, 8)
	if relayErr != nil {
		t.Fatalf("RelayPending() error = %v", relayErr)
	}
	if report.PendingObserved != 1 || report.AcceptedPublications != 1 || report.PublishedMarks != 1 || report.DuplicateMarkObserved != 0 {
		t.Fatalf("RelayPending() report = %#v", report)
	}
	if received.ID != envelope.ID || received.Source != envelope.Source || received.TenantID != envelope.TenantID || string(received.Data) != string(record.Envelope.Data) {
		t.Fatal("relay mutated canonical envelope")
	}
	if store.markCalls != 1 || store.lastEventID != envelope.ID || store.lastRevision != record.Revision || store.lastScope != record.Scope {
		t.Fatal("relay did not mark the exact published record revision")
	}
}

func TestRelayPublishFailureLeavesRecordUnmarked(t *testing.T) {
	envelope := testEnvelope(t)
	record, err := NewOutboxRecord(testOutboxScope(envelope), envelope)
	if err != nil {
		t.Fatalf("NewOutboxRecord() error = %v", err)
	}
	store := &outboxTestStore{pending: []OutboxRecord{record}, markResult: OutboxMarkedPublished}
	publisher, err := NewPublisher(func(context.Context, Envelope) error { return errors.New("private broker failure") })
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	relay, err := NewOutboxRelay(store, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() error = %v", err)
	}

	report, relayErr := relay.RelayPending(context.Background(), record.Scope, 1)
	assertFailureCode(t, relayErr, codeOutboxPublishFailed)
	if report.AcceptedPublications != 0 || store.markCalls != 0 {
		t.Fatalf("failed publication advanced outbox state: report=%#v marks=%d", report, store.markCalls)
	}
}

func TestRelayMarkFailureKeepsDuplicatePublicationWindowExplicit(t *testing.T) {
	envelope := testEnvelope(t)
	record, err := NewOutboxRecord(testOutboxScope(envelope), envelope)
	if err != nil {
		t.Fatalf("NewOutboxRecord() error = %v", err)
	}
	store := &outboxTestStore{
		pending:    []OutboxRecord{record},
		markResult: MarkResultUnknown,
		markErr:    errors.New("private mark failure"),
	}
	publishes := 0
	publisher, err := NewPublisher(func(context.Context, Envelope) error {
		publishes++
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	relay, err := NewOutboxRelay(store, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() error = %v", err)
	}

	first, firstErr := relay.RelayPending(context.Background(), record.Scope, 1)
	assertFailureCode(t, firstErr, codeOutboxMarkFailed)
	second, secondErr := relay.RelayPending(context.Background(), record.Scope, 1)
	assertFailureCode(t, secondErr, codeOutboxMarkFailed)
	if first.AcceptedPublications != 1 || second.AcceptedPublications != 1 || publishes != 2 {
		t.Fatalf("publish+mark crash window was hidden: first=%#v second=%#v publishes=%d", first, second, publishes)
	}
}

func TestRelayRejectsCrossScopeOrMalformedPendingStateBeforePublish(t *testing.T) {
	envelope := testEnvelope(t)
	record, err := NewOutboxRecord(testOutboxScope(envelope), envelope)
	if err != nil {
		t.Fatalf("NewOutboxRecord() error = %v", err)
	}
	requested := record.Scope
	record.Scope.Owner = Producer("urn:omnexa:module:other.owner")
	store := &outboxTestStore{pending: []OutboxRecord{record}}
	called := false
	publisher, err := NewPublisher(func(context.Context, Envelope) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	relay, err := NewOutboxRelay(store, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() error = %v", err)
	}

	_, relayErr := relay.RelayPending(context.Background(), requested, 1)
	assertFailureCode(t, relayErr, codeOutboxStateMalformed)
	if called || store.markCalls != 0 {
		t.Fatal("cross-scope stored state reached publication or mark boundary")
	}
}

func TestRelayTreatsConcurrentAlreadyPublishedMarkAsDuplicateNotExactlyOnce(t *testing.T) {
	envelope := testEnvelope(t)
	record, err := NewOutboxRecord(testOutboxScope(envelope), envelope)
	if err != nil {
		t.Fatalf("NewOutboxRecord() error = %v", err)
	}
	store := &outboxTestStore{pending: []OutboxRecord{record}, markResult: OutboxAlreadyPublished}
	publisher, err := NewPublisher(func(context.Context, Envelope) error { return nil })
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	relay, err := NewOutboxRelay(store, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() error = %v", err)
	}

	report, relayErr := relay.RelayPending(context.Background(), record.Scope, 1)
	if relayErr != nil {
		t.Fatalf("RelayPending() error = %v", relayErr)
	}
	if report.AcceptedPublications != 1 || report.PublishedMarks != 0 || report.DuplicateMarkObserved != 1 {
		t.Fatalf("concurrent duplicate mark was misreported: %#v", report)
	}
}

func TestOutboxCancellationFailsClosedBeforePersistence(t *testing.T) {
	envelope := testEnvelope(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	store := &outboxTestStore{enqueueResult: OutboxEnqueued}
	_, err := EnqueueOutbox(ctx, &outboxTestTx{}, store, testOutboxScope(envelope), envelope)
	assertFailureCode(t, err, codeOutboxInterrupted)
	if store.enqueueCalls != 0 {
		t.Fatal("cancelled enqueue reached persistence")
	}
}

func testOutboxScope(envelope Envelope) OutboxScope {
	return OutboxScope{Owner: envelope.Source, TenantID: envelope.TenantID}
}

func containsPrivateText(value, private string) bool {
	return value == private || len(private) > 0 && len(value) >= len(private) && stringContains(value, private)
}

func stringContains(value, fragment string) bool {
	for index := 0; index+len(fragment) <= len(value); index++ {
		if value[index:index+len(fragment)] == fragment {
			return true
		}
	}
	return false
}
