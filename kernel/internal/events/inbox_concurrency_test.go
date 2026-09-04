package events

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/jackc/pgx/v5"
)

type p0405ConcurrentAttempt struct {
	result InboxApplyResult
	err    error
}

func TestPostgresInboxSameScopeConcurrentAttemptsCannotBothMutate(t *testing.T) {
	ctx, pool := setupP0405ReliabilityDatabase(t)
	store := NewPostgresInboxStore()
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.concurrent")
	var mutationCalls atomic.Int32
	firstMutationEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondTransactionEntered := make(chan struct{})
	attempts := make(chan p0405ConcurrentAttempt, 2)

	go func() {
		result := InboxApplyResultUnknown
		err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
			var applyErr error
			result, applyErr = ApplyInbox(ctx, tx, store, binding, envelope, func(mutationCtx context.Context, mutationTx OutboxTransaction) error {
				mutationCalls.Add(1)
				if _, execErr := mutationTx.Exec(
					mutationCtx,
					`INSERT INTO omnexa_events.p04_05_reliability_owner_state (event_id, consumer_id, note) VALUES ($1::uuid, $2, 'first-concurrent-attempt')`,
					string(envelope.ID),
					binding.ConsumerID,
				); execErr != nil {
					return execErr
				}
				close(firstMutationEntered)
				select {
				case <-releaseFirst:
					return nil
				case <-mutationCtx.Done():
					return mutationCtx.Err()
				}
			})
			return applyErr
		})
		attempts <- p0405ConcurrentAttempt{result: result, err: err}
	}()

	select {
	case <-firstMutationEntered:
	case <-ctx.Done():
		t.Fatalf("first concurrent attempt did not reach mutation: %v", ctx.Err())
	}

	go func() {
		result := InboxApplyResultUnknown
		err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
			close(secondTransactionEntered)
			var applyErr error
			result, applyErr = ApplyInbox(ctx, tx, store, binding, envelope, func(mutationCtx context.Context, mutationTx OutboxTransaction) error {
				mutationCalls.Add(1)
				_, execErr := mutationTx.Exec(
					mutationCtx,
					`INSERT INTO omnexa_events.p04_05_reliability_owner_state (event_id, consumer_id, note) VALUES ($1::uuid, $2, 'second-concurrent-attempt')`,
					string(envelope.ID),
					binding.ConsumerID,
				)
				return execErr
			})
			return applyErr
		})
		attempts <- p0405ConcurrentAttempt{result: result, err: err}
	}()

	select {
	case <-secondTransactionEntered:
	case <-ctx.Done():
		t.Fatalf("second concurrent transaction did not begin: %v", ctx.Err())
	}
	close(releaseFirst)

	first := <-attempts
	second := <-attempts
	if first.err != nil || second.err != nil {
		t.Fatalf("concurrent attempts returned errors: first=%v/%v second=%v/%v", first.result, first.err, second.result, second.err)
	}
	applied := 0
	duplicates := 0
	for _, attempt := range []p0405ConcurrentAttempt{first, second} {
		switch attempt.result {
		case InboxApplied:
			applied++
		case InboxAlreadyApplied:
			duplicates++
		default:
			t.Fatalf("unexpected concurrent inbox result: %v", attempt.result)
		}
	}
	if applied != 1 || duplicates != 1 || mutationCalls.Load() != 1 {
		t.Fatalf("same-scope concurrency = applied:%d duplicate:%d mutations:%d; want 1/1/1", applied, duplicates, mutationCalls.Load())
	}
	assertP0405ReliabilityCount(
		t,
		ctx,
		pool,
		`SELECT count(*) FROM omnexa_events.p04_05_reliability_owner_state WHERE event_id=$1::uuid AND consumer_id=$2`,
		[]any{string(envelope.ID), binding.ConsumerID},
		1,
	)
	assertP0405ReliabilityCount(
		t,
		ctx,
		pool,
		`SELECT count(*) FROM omnexa_events.consumer_inbox WHERE event_id=$1::uuid AND consumer_id=$2 AND processing_state='completed'`,
		[]any{string(envelope.ID), binding.ConsumerID},
		1,
	)
}

func TestPostgresInboxSameEventDifferentConsumersRemainIndependentConcurrently(t *testing.T) {
	ctx, pool := setupP0405ReliabilityDatabase(t)
	store := NewPostgresInboxStore()
	envelope := testEnvelope(t)
	bindings := []DurableBinding{
		testInboxBinding(envelope, "billing.consumer-a"),
		testInboxBinding(envelope, "analytics.consumer-b"),
	}
	var mutationCalls atomic.Int32
	start := make(chan struct{})
	var ready sync.WaitGroup
	ready.Add(len(bindings))
	attempts := make(chan p0405ConcurrentAttempt, len(bindings))

	for _, binding := range bindings {
		go func() {
			ready.Done()
			<-start
			result := InboxApplyResultUnknown
			err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
				var applyErr error
				result, applyErr = ApplyInbox(ctx, tx, store, binding, envelope, func(mutationCtx context.Context, mutationTx OutboxTransaction) error {
					mutationCalls.Add(1)
					_, execErr := mutationTx.Exec(
						mutationCtx,
						`INSERT INTO omnexa_events.p04_05_reliability_owner_state (event_id, consumer_id, note) VALUES ($1::uuid, $2, 'independent-consumer')`,
						string(envelope.ID),
						binding.ConsumerID,
					)
					return execErr
				})
				return applyErr
			})
			attempts <- p0405ConcurrentAttempt{result: result, err: err}
		}()
	}
	ready.Wait()
	close(start)

	for range bindings {
		attempt := <-attempts
		if attempt.err != nil || attempt.result != InboxApplied {
			t.Fatalf("independent consumer application = %v/%v, want applied", attempt.result, attempt.err)
		}
	}
	if mutationCalls.Load() != 2 {
		t.Fatalf("different consumer mutation calls = %d, want 2", mutationCalls.Load())
	}
	assertP0405ReliabilityCount(
		t,
		ctx,
		pool,
		`SELECT count(*) FROM omnexa_events.consumer_inbox WHERE event_id=$1::uuid AND processing_state='completed'`,
		[]any{string(envelope.ID)},
		2,
	)
}

func TestInboxReliabilityRejectsRouteAndTenantRebindingBeforeClaim(t *testing.T) {
	envelope := testEnvelope(t)
	store := &inboxTestStore{claimResult: InboxClaimed}
	mutation := func(context.Context, OutboxTransaction) error {
		return errors.New("mutation must not execute for rebinding")
	}

	routeConflict := testInboxBinding(envelope, "billing.rebinding")
	routeConflict.EventType = EventType("billing.other.changed.v1")
	result, err := ApplyInbox(context.Background(), &outboxTestTx{}, store, routeConflict, envelope, mutation)
	if result != InboxApplyResultUnknown {
		t.Fatalf("route rebinding result = %v, want unknown", result)
	}
	assertFailureCode(t, err, codeInboxConflict)
	if store.claimCalls != 0 {
		t.Fatal("route rebinding crossed the inbox persistence seam")
	}

	tenantConflict := testInboxBinding(envelope, "billing.rebinding")
	tenantConflict.Scope.TenantID = "01990f6e-1f30-7000-8000-000000000099"
	result, err = ApplyInbox(context.Background(), &outboxTestTx{}, store, tenantConflict, envelope, mutation)
	if result != InboxApplyResultUnknown || err == nil {
		t.Fatalf("tenant rebinding result/error = %v/%v", result, err)
	}
	if store.claimCalls != 0 {
		t.Fatal("tenant rebinding crossed the inbox persistence seam")
	}
}
