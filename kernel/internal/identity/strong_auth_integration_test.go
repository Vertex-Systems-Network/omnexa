package identity

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresStrongAuthenticationIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_07_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_07_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	pool, poolErr := pgxpool.New(ctx, databaseURL)
	if poolErr != nil {
		t.Fatalf("pgxpool.New() error = %v", poolErr)
	}
	defer pool.Close()
	if pingErr := pool.Ping(ctx); pingErr != nil {
		t.Fatalf("pool.Ping() error = %v", pingErr)
	}
	foundation, foundationErr := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if foundationErr != nil {
		t.Fatalf("NewMigrator(kernel.foundation) error = %v", foundationErr)
	}
	if migrationErr := foundation.Run(ctx); migrationErr != nil {
		t.Fatalf("foundation migrator Run() error = %v", migrationErr)
	}
	defer resetP0204Identity(t, pool)

	v1SQL, readV1Err := os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	if readV1Err != nil {
		t.Fatalf("read identity v1 migration error = %v", readV1Err)
	}
	v2SQL, readV2Err := os.ReadFile("../../migrations/kernel.identity/2_create_authentication_sessions.sql")
	if readV2Err != nil {
		t.Fatalf("read identity v2 migration error = %v", readV2Err)
	}
	v3SQL, readV3Err := os.ReadFile("../../migrations/kernel.identity/3_create_strong_authentication.sql")
	if readV3Err != nil {
		t.Fatalf("read identity v3 migration error = %v", readV3Err)
	}
	v1 := database.Migration{Version: 1, Name: "create_identity_foundation", SQL: string(v1SQL)}
	v2 := database.Migration{Version: 2, Name: "create_authentication_sessions", SQL: string(v2SQL)}
	v3 := database.Migration{Version: 3, Name: "create_strong_authentication", SQL: string(v3SQL)}

	resetP0204Identity(t, pool)
	freshMigrator, freshErr := database.NewMigrator(pool, "kernel.identity", []database.Migration{v1, v2, v3}, 5*time.Second)
	if freshErr != nil {
		t.Fatalf("NewMigrator(fresh P02.07) error = %v", freshErr)
	}
	if migrationErr := freshMigrator.Run(ctx); migrationErr != nil {
		t.Fatalf("fresh P02.07 migration error = %v", migrationErr)
	}
	if migrationErr := freshMigrator.Run(ctx); migrationErr != nil {
		t.Fatalf("idempotent P02.07 migration error = %v", migrationErr)
	}
	var freshLedgerCount int
	if queryErr := pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity'").Scan(&freshLedgerCount); queryErr != nil {
		t.Fatalf("fresh identity ledger query error = %v", queryErr)
	}
	if freshLedgerCount != 3 {
		t.Fatalf("fresh identity ledger count = %d, want 3", freshLedgerCount)
	}

	resetP0204Identity(t, pool)
	baselineMigrator, baselineErr := database.NewMigrator(pool, "kernel.identity", []database.Migration{v1, v2}, 5*time.Second)
	if baselineErr != nil {
		t.Fatalf("NewMigrator(P02.04 baseline) error = %v", baselineErr)
	}
	if migrationErr := baselineMigrator.Run(ctx); migrationErr != nil {
		t.Fatalf("P02.04 baseline migration error = %v", migrationErr)
	}

	identityRepository, identityRepositoryErr := NewPostgresRepository(pool)
	if identityRepositoryErr != nil {
		t.Fatalf("NewPostgresRepository() error = %v", identityRepositoryErr)
	}
	authRepository, authRepositoryErr := NewPostgresAuthenticationRepository(pool)
	if authRepositoryErr != nil {
		t.Fatalf("NewPostgresAuthenticationRepository() error = %v", authRepositoryErr)
	}
	audit := &captureSecurityAudit{}
	authService, authServiceErr := NewAuthenticationService(
		authRepository,
		NewPBKDF2PasswordHasher(),
		DefaultSessionPolicy(),
		nil,
		audit,
	)
	if authServiceErr != nil {
		t.Fatalf("NewAuthenticationService() error = %v", authServiceErr)
	}
	createdAt := time.Date(2026, time.August, 24, 1, 0, 0, 0, time.UTC)
	user, userErr := newUserAt(fixedUserID, "p0207-owner@example.com", createdAt)
	if userErr != nil {
		t.Fatalf("newUserAt() error = %v", userErr)
	}
	if createErr := identityRepository.Create(ctx, user); createErr != nil {
		t.Fatalf("Create(user) error = %v", createErr)
	}
	active, transitionErr := identityRepository.Transition(ctx, user.ID(), LifecycleProvisioned, LifecycleActive, createdAt.Add(time.Minute))
	if transitionErr != nil {
		t.Fatalf("Transition(active) error = %v", transitionErr)
	}
	password := strings.Join([]string{"Synthetic", "P02", "07", "Password", "2026!"}, "-")
	passwordAt := createdAt.Add(2 * time.Minute)
	if enrollErr := authService.EnrollPassword(ctx, active.ID(), password, passwordAt); enrollErr != nil {
		t.Fatalf("EnrollPassword() error = %v", enrollErr)
	}
	passwordProof, authenticateErr := authService.AuthenticatePassword(ctx, active.ID(), password, passwordAt.Add(time.Second))
	if authenticateErr != nil {
		t.Fatalf("AuthenticatePassword() error = %v", authenticateErr)
	}
	platformContext, contextErr := NewSessionContext("", "")
	if contextErr != nil {
		t.Fatalf("NewSessionContext(platform) error = %v", contextErr)
	}
	issuedAt := passwordAt.Add(2 * time.Second)
	issued, issueErr := authService.IssueSession(ctx, passwordProof, platformContext, "p0207-browser", issuedAt)
	if issueErr != nil {
		t.Fatalf("IssueSession() error = %v", issueErr)
	}

	upgradeMigrator, upgradeErr := database.NewMigrator(pool, "kernel.identity", []database.Migration{v1, v2, v3}, 5*time.Second)
	if upgradeErr != nil {
		t.Fatalf("NewMigrator(upgrade P02.07) error = %v", upgradeErr)
	}
	if migrationErr := upgradeMigrator.Run(ctx); migrationErr != nil {
		t.Fatalf("P02.04 -> P02.07 upgrade migration error = %v", migrationErr)
	}
	validated, validateErr := authService.ValidateAccess(ctx, issued.AccessSecret(), issuedAt.Add(3*time.Second))
	if validateErr != nil || validated.Session().ID() != issued.Session().ID() {
		t.Fatalf("P02.04 session not preserved after P02.07 upgrade: %q/%v", validated.Session().ID(), validateErr)
	}

	strongRepository, strongRepositoryErr := NewPostgresStrongAuthenticationRepository(pool)
	if strongRepositoryErr != nil {
		t.Fatalf("NewPostgresStrongAuthenticationRepository() error = %v", strongRepositoryErr)
	}
	strongService, strongServiceErr := newStrongAuthenticationService(
		strongRepository,
		&fakePasskeyVerifier{},
		DefaultStrongAuthenticationPolicy(),
		audit,
		&incrementingReader{},
	)
	if strongServiceErr != nil {
		t.Fatalf("newStrongAuthenticationService() error = %v", strongServiceErr)
	}

	enrollmentAt := issuedAt.Add(4 * time.Second)
	enrollment, beginErr := strongService.BeginPasskeyEnrollment(ctx, validated, "primary-passkey", enrollmentAt)
	if beginErr != nil {
		t.Fatalf("BeginPasskeyEnrollment() error = %v", beginErr)
	}
	challengeDigest := sha256.Sum256([]byte(enrollment.Challenge().Reveal()))
	var challengeDigestCount int
	if queryErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_identity.authentication_challenges WHERE id = $1 AND secret_digest = $2",
		string(enrollment.ChallengeID()), challengeDigest[:],
	).Scan(&challengeDigestCount); queryErr != nil {
		t.Fatalf("challenge digest lookup error = %v", queryErr)
	}
	if challengeDigestCount != 1 {
		t.Fatalf("raw challenge was not represented by exactly one digest")
	}
	response, responseErr := NewPasskeyResponse([]byte("synthetic-registration-response"))
	if responseErr != nil {
		t.Fatalf("NewPasskeyResponse(registration) error = %v", responseErr)
	}
	wrongPrincipal := strongAuthenticatedSession(fixedOtherStrongUserID, issued.Session().ID(), enrollmentAt)
	if _, wrongPrincipalErr := strongService.CompletePasskeyEnrollment(
		ctx, wrongPrincipal, enrollment.ChallengeID(), enrollment.Challenge(), response, enrollmentAt.Add(time.Second),
	); !failure.IsCode(wrongPrincipalErr, codeStrongChallengeInvalid) {
		t.Fatalf("wrong-principal enrollment error = %v", wrongPrincipalErr)
	}
	wrongSession := strongAuthenticatedSession(active.ID(), fixedOtherStrongSessionID, enrollmentAt)
	if _, wrongSessionErr := strongService.CompletePasskeyEnrollment(
		ctx, wrongSession, enrollment.ChallengeID(), enrollment.Challenge(), response, enrollmentAt.Add(time.Second),
	); !failure.IsCode(wrongSessionErr, codeStrongChallengeInvalid) {
		t.Fatalf("wrong-session enrollment error = %v", wrongSessionErr)
	}
	strongProof, completeErr := strongService.CompletePasskeyEnrollment(
		ctx, validated, enrollment.ChallengeID(), enrollment.Challenge(), response, enrollmentAt.Add(2*time.Second),
	)
	if completeErr != nil || !strongProof.valid() {
		t.Fatalf("CompletePasskeyEnrollment() proof/error = %+v/%v", strongProof, completeErr)
	}
	if _, replayErr := strongService.CompletePasskeyEnrollment(
		ctx, validated, enrollment.ChallengeID(), enrollment.Challenge(), response, enrollmentAt.Add(3*time.Second),
	); !failure.IsCode(replayErr, codeStrongChallengeInvalid) {
		t.Fatalf("replayed enrollment challenge error = %v", replayErr)
	}

	expiredEnrollment, expiredBeginErr := strongService.BeginPasskeyEnrollment(ctx, validated, "expiring-passkey", enrollmentAt.Add(4*time.Second))
	if expiredBeginErr != nil {
		t.Fatalf("BeginPasskeyEnrollment(expiring) error = %v", expiredBeginErr)
	}
	if _, expiredErr := strongService.CompletePasskeyEnrollment(
		ctx,
		validated,
		expiredEnrollment.ChallengeID(),
		expiredEnrollment.Challenge(),
		response,
		expiredEnrollment.ExpiresAt().Add(time.Second),
	); !failure.IsCode(expiredErr, codeStrongChallengeInvalid) {
		t.Fatalf("expired enrollment challenge error = %v", expiredErr)
	}

	assertionAt := enrollmentAt.Add(10 * time.Second)
	assertion, assertionErr := strongService.BeginPasskeyAssertion(ctx, validated, enrollment.Factor().ID(), assertionAt)
	if assertionErr != nil {
		t.Fatalf("BeginPasskeyAssertion() error = %v", assertionErr)
	}
	assertionResponse, assertionResponseErr := NewPasskeyResponse([]byte("synthetic-assertion-response"))
	if assertionResponseErr != nil {
		t.Fatalf("NewPasskeyResponse(assertion) error = %v", assertionResponseErr)
	}
	if _, wrongBindingErr := strongService.CompletePasskeyAssertion(
		ctx,
		wrongSession,
		enrollment.Factor().ID(),
		assertion.ChallengeID(),
		assertion.Challenge(),
		assertionResponse,
		assertionAt.Add(time.Second),
	); !failure.IsCode(wrongBindingErr, codeStrongChallengeInvalid) {
		t.Fatalf("wrong-session assertion error = %v", wrongBindingErr)
	}
	asserted, assertionCompleteErr := strongService.CompletePasskeyAssertion(
		ctx,
		validated,
		enrollment.Factor().ID(),
		assertion.ChallengeID(),
		assertion.Challenge(),
		assertionResponse,
		assertionAt.Add(2*time.Second),
	)
	if assertionCompleteErr != nil || strongService.RequireStepUp(validated, asserted, assertionAt.Add(3*time.Second)) != nil {
		t.Fatalf("CompletePasskeyAssertion()/RequireStepUp() error = %v", assertionCompleteErr)
	}
	if _, replayErr := strongService.CompletePasskeyAssertion(
		ctx,
		validated,
		enrollment.Factor().ID(),
		assertion.ChallengeID(),
		assertion.Challenge(),
		assertionResponse,
		assertionAt.Add(4*time.Second),
	); !failure.IsCode(replayErr, codeStrongChallengeInvalid) {
		t.Fatalf("replayed assertion challenge error = %v", replayErr)
	}

	recoveryAt := assertionAt.Add(5 * time.Second)
	bundle, recoveryIssueErr := strongService.IssueRecoveryCodes(ctx, validated, asserted, recoveryAt)
	if recoveryIssueErr != nil || len(bundle.Codes()) != 8 {
		t.Fatalf("IssueRecoveryCodes() count/error = %d/%v", len(bundle.Codes()), recoveryIssueErr)
	}
	codes := bundle.Codes()
	firstRecoveryDigest := sha256.Sum256([]byte(codes[0].Reveal()))
	var recoveryDigestCount int
	if queryErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_identity.recovery_codes WHERE secret_digest = $1",
		firstRecoveryDigest[:],
	).Scan(&recoveryDigestCount); queryErr != nil {
		t.Fatalf("recovery digest lookup error = %v", queryErr)
	}
	if recoveryDigestCount != 1 {
		t.Fatalf("raw recovery code was not represented by exactly one digest")
	}
	recoveryProof, recoveryErr := strongService.UseRecoveryCode(ctx, validated, codes[0], recoveryAt.Add(time.Second))
	if recoveryErr != nil || strongService.RequireStepUp(validated, recoveryProof, recoveryAt.Add(2*time.Second)) != nil {
		t.Fatalf("UseRecoveryCode()/RequireStepUp() error = %v", recoveryErr)
	}
	if _, replayRecoveryErr := strongService.UseRecoveryCode(ctx, validated, codes[0], recoveryAt.Add(3*time.Second)); !failure.IsCode(replayRecoveryErr, codeRecoveryCodeInvalid) {
		t.Fatalf("replayed recovery code error = %v", replayRecoveryErr)
	}

	removedAt := recoveryAt.Add(4 * time.Second)
	removed, removeErr := strongService.RemoveFactor(ctx, validated, recoveryProof, enrollment.Factor().ID(), removedAt)
	if removeErr != nil || removed.State() != StrongFactorRevoked {
		t.Fatalf("RemoveFactor() state/error = %q/%v", removed.State(), removeErr)
	}
	if _, accessErr := authService.ValidateAccess(ctx, issued.AccessSecret(), removedAt.Add(time.Second)); !failure.IsCode(accessErr, codeSessionInvalid) {
		t.Fatalf("factor removal did not invalidate interactive session: %v", accessErr)
	}

	var forbiddenColumns int
	if queryErr := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_identity'
		   AND column_name IN (
		       'challenge_secret','raw_challenge','recovery_code','raw_recovery_code',
		       'private_key','authenticator_secret','api_key','api_secret','service_account_id',
		       'saml_assertion','sso_provider_id','tenant_setting_id'
		   )`,
	).Scan(&forbiddenColumns); queryErr != nil {
		t.Fatalf("forbidden-column query error = %v", queryErr)
	}
	if forbiddenColumns != 0 {
		t.Fatalf("P02.07 persisted %d forbidden raw-secret/P02.08+/P24/settings columns", forbiddenColumns)
	}
	var forbiddenTables int
	if queryErr := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.tables
		 WHERE table_schema = 'omnexa_identity'
		   AND table_name IN ('service_accounts','api_credentials','sso_providers','saml_credentials','tenant_settings')`,
	).Scan(&forbiddenTables); queryErr != nil {
		t.Fatalf("forbidden-table query error = %v", queryErr)
	}
	if forbiddenTables != 0 {
		t.Fatalf("P02.07 pulled %d P02.08+/P24/settings tables forward", forbiddenTables)
	}

	auditJSON, marshalAuditErr := json.Marshal(audit.events)
	if marshalAuditErr != nil {
		t.Fatalf("json.Marshal(audit events) error = %v", marshalAuditErr)
	}
	auditPayload := string(auditJSON)
	for _, restricted := range []string{
		enrollment.Challenge().Reveal(),
		expiredEnrollment.Challenge().Reveal(),
		assertion.Challenge().Reveal(),
		codes[0].Reveal(),
		string(response.Reveal()),
		string(assertionResponse.Reveal()),
	} {
		if strings.Contains(auditPayload, restricted) {
			t.Fatalf("security audit leaked RESTRICTED P02.07 authentication material")
		}
	}

	drifted := []database.Migration{v1, v2, {
		Version: 3,
		Name:    "create_strong_authentication",
		SQL:     string(v3SQL) + "\n-- rewritten applied P02.07 migration must fail closed\n",
	}}
	driftMigrator, driftErr := database.NewMigrator(pool, "kernel.identity", drifted, 5*time.Second)
	if driftErr != nil {
		t.Fatalf("NewMigrator(drifted P02.07) error = %v", driftErr)
	}
	if migrationErr := driftMigrator.Run(ctx); migrationErr == nil {
		t.Fatalf("rewritten applied P02.07 migration unexpectedly succeeded")
	}
}
