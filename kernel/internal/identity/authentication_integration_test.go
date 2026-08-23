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

func TestPostgresAuthenticationSessionLifecycleIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_04_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_04_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		t.Fatalf("pool.Ping() error = %v", err)
	}
	foundation, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(kernel.foundation) error = %v", err)
	}
	if err = foundation.Run(ctx); err != nil {
		t.Fatalf("foundation migrator Run() error = %v", err)
	}
	defer resetP0204Identity(t, pool)

	v1SQL, err := os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	if err != nil {
		t.Fatalf("read identity v1 migration error = %v", err)
	}
	v2SQL, err := os.ReadFile("../../migrations/kernel.identity/2_create_authentication_sessions.sql")
	if err != nil {
		t.Fatalf("read identity v2 migration error = %v", err)
	}
	v1 := database.Migration{Version: 1, Name: "create_identity_foundation", SQL: string(v1SQL)}
	v2 := database.Migration{Version: 2, Name: "create_authentication_sessions", SQL: string(v2SQL)}

	resetP0204Identity(t, pool)
	freshMigrator, err := database.NewMigrator(pool, "kernel.identity", []database.Migration{v1, v2}, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(fresh) error = %v", err)
	}
	if err = freshMigrator.Run(ctx); err != nil {
		t.Fatalf("fresh P02.04 migration error = %v", err)
	}
	if err = freshMigrator.Run(ctx); err != nil {
		t.Fatalf("idempotent P02.04 migration error = %v", err)
	}
	var freshLedgerCount int
	if err = pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity'").Scan(&freshLedgerCount); err != nil {
		t.Fatalf("fresh ledger query error = %v", err)
	}
	if freshLedgerCount != 2 {
		t.Fatalf("fresh identity ledger count = %d, want 2", freshLedgerCount)
	}

	resetP0204Identity(t, pool)
	v1Migrator, err := database.NewMigrator(pool, "kernel.identity", []database.Migration{v1}, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(v1) error = %v", err)
	}
	if err = v1Migrator.Run(ctx); err != nil {
		t.Fatalf("P02.01 baseline migration error = %v", err)
	}
	identityRepository, err := NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("NewPostgresRepository() error = %v", err)
	}
	createdAt := time.Date(2026, time.August, 23, 14, 0, 0, 0, time.UTC)
	user, err := newUserAt(fixedUserID, "auth-owner@example.com", createdAt)
	if err != nil {
		t.Fatalf("newUserAt() error = %v", err)
	}
	if err = identityRepository.Create(ctx, user); err != nil {
		t.Fatalf("Create(user) error = %v", err)
	}
	active, err := identityRepository.Transition(ctx, user.ID(), LifecycleProvisioned, LifecycleActive, createdAt.Add(time.Minute))
	if err != nil {
		t.Fatalf("Transition(active) error = %v", err)
	}

	upgradeMigrator, err := database.NewMigrator(pool, "kernel.identity", []database.Migration{v1, v2}, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(upgrade) error = %v", err)
	}
	if err = upgradeMigrator.Run(ctx); err != nil {
		t.Fatalf("P02.01 -> P02.04 upgrade migration error = %v", err)
	}
	preserved, err := identityRepository.Get(ctx, active.ID())
	if err != nil || preserved.State() != LifecycleActive {
		t.Fatalf("upgrade did not preserve P02.01 user: state=%q error=%v", preserved.State(), err)
	}

	authRepository, err := NewPostgresAuthenticationRepository(pool)
	if err != nil {
		t.Fatalf("NewPostgresAuthenticationRepository() error = %v", err)
	}
	reauthorizer := &fakeContextReauthorizer{}
	audit := &captureSecurityAudit{}
	service, err := NewAuthenticationService(
		authRepository,
		NewPBKDF2PasswordHasher(),
		DefaultSessionPolicy(),
		reauthorizer,
		audit,
	)
	if err != nil {
		t.Fatalf("NewAuthenticationService() error = %v", err)
	}

	password := "Synthetic-Auth-Password-2026!"
	credentialAt := createdAt.Add(2 * time.Minute)
	if err = service.EnrollPassword(ctx, active.ID(), password, credentialAt); err != nil {
		t.Fatalf("EnrollPassword() error = %v", err)
	}
	wrongAt := credentialAt.Add(time.Second)
	_, wrongErr := service.AuthenticatePassword(ctx, active.ID(), "wrong-synthetic-password", wrongAt)
	missingID := UserID("01890f3e-7b9a-7cc0-98c4-dc0c0c073993")
	_, missingErr := service.AuthenticatePassword(ctx, missingID, "wrong-synthetic-password", wrongAt)
	if !failure.IsCode(wrongErr, codeAuthenticationFailed) || !failure.IsCode(missingErr, codeAuthenticationFailed) || wrongErr.Error() != missingErr.Error() {
		t.Fatalf("disclosure-safe auth mismatch wrong=%v missing=%v", wrongErr, missingErr)
	}
	authentication, err := service.AuthenticatePassword(ctx, active.ID(), password, credentialAt.Add(2*time.Second))
	if err != nil {
		t.Fatalf("AuthenticatePassword(valid) error = %v", err)
	}

	sessionContext, err := NewSessionContext(fixedTenantContextID, fixedOrganizationContextID)
	if err != nil {
		t.Fatalf("NewSessionContext() error = %v", err)
	}
	issuedAt := credentialAt.Add(3 * time.Second)
	issued, err := service.IssueSession(ctx, authentication, sessionContext, "synthetic-browser", issuedAt)
	if err != nil {
		t.Fatalf("IssueSession() error = %v", err)
	}
	if issued.AccessExpiresAt().Sub(issuedAt) != 15*time.Minute || issued.RefreshExpiresAt().Sub(issuedAt) != 7*24*time.Hour || issued.Session().ExpiresAt().Sub(issuedAt) != 30*24*time.Hour {
		t.Fatalf("issued lifetimes do not match explicit policy")
	}
	if _, err = service.ValidateAccess(ctx, issued.AccessSecret(), issuedAt.Add(time.Second)); err != nil {
		t.Fatalf("ValidateAccess(valid) error = %v", err)
	}

	rotatedAt := issuedAt.Add(2 * time.Minute)
	rotated, err := service.RotateRefresh(ctx, issued.RefreshSecret(), rotatedAt)
	if err != nil {
		t.Fatalf("RotateRefresh(valid) error = %v", err)
	}
	if _, err = service.ValidateAccess(ctx, issued.AccessSecret(), rotatedAt.Add(time.Second)); !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("old access after refresh rotation error = %v", err)
	}
	if _, err = service.RotateRefresh(ctx, issued.RefreshSecret(), rotatedAt.Add(time.Second)); !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("replayed refresh error = %v", err)
	}
	if _, err = service.ValidateAccess(ctx, rotated.AccessSecret(), rotatedAt.Add(time.Second)); err != nil {
		t.Fatalf("rotated access error = %v", err)
	}

	reauthorizer.reject = true
	if _, err = service.ValidateAccess(ctx, rotated.AccessSecret(), rotatedAt.Add(2*time.Second)); !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("stale context access error = %v", err)
	}
	reauthorizer.reject = false

	changedAt := rotatedAt.Add(3 * time.Minute)
	newPassword := "Synthetic-Changed-Password-2026!"
	if err = service.ChangePassword(ctx, authentication, newPassword, changedAt); err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}
	if _, err = service.ValidateAccess(ctx, rotated.AccessSecret(), changedAt.Add(time.Second)); !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("password-change session invalidation error = %v", err)
	}
	if _, err = service.IssueSession(ctx, authentication, sessionContext, "stale-proof", changedAt.Add(time.Second)); !failure.IsCode(err, codeAuthenticationInvalid) {
		t.Fatalf("stale authentication proof error = %v", err)
	}

	freshAuthentication, err := service.AuthenticatePassword(ctx, active.ID(), newPassword, changedAt.Add(2*time.Second))
	if err != nil {
		t.Fatalf("AuthenticatePassword(changed) error = %v", err)
	}
	platformContext, err := NewSessionContext("", "")
	if err != nil {
		t.Fatalf("NewSessionContext(platform) error = %v", err)
	}
	lifecycleSession, err := service.IssueSession(ctx, freshAuthentication, platformContext, "lifecycle-browser", changedAt.Add(3*time.Second))
	if err != nil {
		t.Fatalf("IssueSession(lifecycle) error = %v", err)
	}
	suspendedAt := changedAt.Add(4 * time.Second)
	if _, err = identityRepository.Transition(ctx, active.ID(), LifecycleActive, LifecycleSuspended, suspendedAt); err != nil {
		t.Fatalf("Transition(suspended) error = %v", err)
	}
	if _, err = service.ValidateAccess(ctx, lifecycleSession.AccessSecret(), suspendedAt.Add(time.Second)); !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("suspended user access error = %v", err)
	}
	reactivatedAt := suspendedAt.Add(2 * time.Second)
	if _, err = identityRepository.Transition(ctx, active.ID(), LifecycleSuspended, LifecycleActive, reactivatedAt); err != nil {
		t.Fatalf("Transition(reactivated) error = %v", err)
	}
	if _, err = service.ValidateAccess(ctx, lifecycleSession.AccessSecret(), reactivatedAt.Add(time.Second)); !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("pre-suspension session resurrected after reactivation: %v", err)
	}

	reauthenticated, err := service.AuthenticatePassword(ctx, active.ID(), newPassword, reactivatedAt.Add(2*time.Second))
	if err != nil {
		t.Fatalf("AuthenticatePassword(reactivated) error = %v", err)
	}
	revocable, err := service.IssueSession(ctx, reauthenticated, platformContext, "revocable-browser", reactivatedAt.Add(3*time.Second))
	if err != nil {
		t.Fatalf("IssueSession(revocable) error = %v", err)
	}
	revokedAt := reactivatedAt.Add(4 * time.Second)
	revoked, err := service.RevokeSession(ctx, reauthenticated, revocable.Session().ID(), revokedAt)
	if err != nil || revoked.RevokedAt().IsZero() {
		t.Fatalf("RevokeSession() revoked/error = %v/%v", revoked.RevokedAt(), err)
	}
	if _, err = service.RevokeSession(ctx, reauthenticated, revocable.Session().ID(), revokedAt.Add(time.Second)); err != nil {
		t.Fatalf("RevokeSession(idempotent) error = %v", err)
	}
	if _, err = service.ValidateAccess(ctx, revocable.AccessSecret(), revokedAt.Add(time.Second)); !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("revoked access error = %v", err)
	}
	sessions, err := service.ListSessions(ctx, reauthenticated, revokedAt.Add(2*time.Second))
	if err != nil || len(sessions) < 3 {
		t.Fatalf("ListSessions() count/error = %d/%v", len(sessions), err)
	}

	var storedPasswordHash string
	if err = pool.QueryRow(ctx, "SELECT password_hash FROM omnexa_identity.password_credentials WHERE principal_id = $1", string(active.ID())).Scan(&storedPasswordHash); err != nil {
		t.Fatalf("password hash query error = %v", err)
	}
	if storedPasswordHash == newPassword || !strings.HasPrefix(storedPasswordHash, "$pbkdf2-sha256$i=600000$") {
		t.Fatalf("password persistence is not one-way/governed")
	}
	var invalidSecretDigests int
	if err = pool.QueryRow(
		ctx,
		`SELECT
		   (SELECT count(*) FROM omnexa_identity.access_credentials WHERE octet_length(secret_digest) <> 32) +
		   (SELECT count(*) FROM omnexa_identity.refresh_credentials WHERE octet_length(secret_digest) <> 32)`,
	).Scan(&invalidSecretDigests); err != nil {
		t.Fatalf("secret digest query error = %v", err)
	}
	if invalidSecretDigests != 0 {
		t.Fatalf("persisted non-digest session secret rows = %d", invalidSecretDigests)
	}
	var forbiddenColumns int
	if err = pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_identity'
		   AND column_name IN ('role_id','permission_id','policy_id','mfa_secret','passkey','api_key','service_account_id')`,
	).Scan(&forbiddenColumns); err != nil {
		t.Fatalf("forbidden-column query error = %v", err)
	}
	if forbiddenColumns != 0 {
		t.Fatalf("P02.04 pulled %d future authz/MFA/service-account columns forward", forbiddenColumns)
	}
	auditJSON, err := json.Marshal(audit.events)
	if err != nil {
		t.Fatalf("json.Marshal(audit events) error = %v", err)
	}
	auditPayload := string(auditJSON)
	for _, forbidden := range []string{password, newPassword, issued.AccessSecret().Reveal(), issued.RefreshSecret().Reveal(), rotated.AccessSecret().Reveal(), rotated.RefreshSecret().Reveal()} {
		if strings.Contains(auditPayload, forbidden) {
			t.Fatalf("audit hook leaked restricted authentication material")
		}
	}

	drifted := []database.Migration{v1, {
		Version: 2,
		Name:    "create_authentication_sessions",
		SQL:     string(v2SQL) + "\n-- rewritten migration must fail closed\n",
	}}
	driftMigrator, err := database.NewMigrator(pool, "kernel.identity", drifted, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(drifted) error = %v", err)
	}
	if err = driftMigrator.Run(ctx); err == nil {
		t.Fatalf("rewritten applied P02.04 migration unexpectedly succeeded")
	}

	accessDigest := sha256.Sum256([]byte(revocable.AccessSecret().Reveal()))
	var storedDigestCount int
	if err = pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_identity.access_credentials WHERE secret_digest = $1",
		accessDigest[:],
	).Scan(&storedDigestCount); err != nil {
		t.Fatalf("secret digest lookup error = %v", err)
	}
	if storedDigestCount != 1 {
		t.Fatalf("opaque access secret was not represented by exactly one digest")
	}
}

func resetP0204Identity(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE"); err != nil {
		t.Fatalf("drop identity schema error = %v", err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity'"); err != nil {
		t.Fatalf("reset identity migration ledger error = %v", err)
	}
}
