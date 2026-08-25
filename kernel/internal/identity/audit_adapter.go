package identity

import (
	"context"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
)

const (
	identityAuditActorSystem = "kernel.identity"
	identityAuditActorUser   = "user"
	identityAuditActorSystemKind = "system"
	identityAuditActorServiceAccount = "service_account"
)

// AuditAdapter bridges the secret-free P02 identity lifecycle hooks into the
// append-only kernel.audit transport without granting audit read/export authority.
// Existing hook interfaces are best-effort by design; protected mutations that
// require acknowledged audit delivery continue to use audit.Writer directly.
type AuditAdapter struct {
	writer *audit.Writer
}

// NewAuditAdapter creates the P02.10 identity audit bridge.
func NewAuditAdapter(writer *audit.Writer) (*AuditAdapter, error) {
	if writer == nil {
		return nil, repositoryInvalidFailure()
	}
	return &AuditAdapter{writer: writer}, nil
}

// RecordSecurityEvent implements SecurityAuditHook. The event vocabulary contains
// no password, token, passkey response, recovery code, digest, or authorization decision.
func (adapter *AuditAdapter) RecordSecurityEvent(event SecurityAuditEvent) {
	if adapter == nil || adapter.writer == nil || event.OccurredAt.IsZero() || event.Action == "" {
		return
	}
	actor := audit.Actor{Kind: identityAuditActorSystemKind, Reference: identityAuditActorSystem}
	target := audit.Target{Kind: "identity.security", Reference: string(event.Action)}
	if event.PrincipalID.Valid() {
		actor = audit.Actor{Kind: identityAuditActorUser, Reference: string(event.PrincipalID)}
		target = audit.Target{Kind: "identity.user", Reference: string(event.PrincipalID)}
	}
	if event.SessionID.Valid() {
		target = audit.Target{Kind: "identity.session", Reference: string(event.SessionID)}
	}
	adapter.writeBestEffort(event.OccurredAt, audit.RecordInput{
		Classification: audit.ClassificationConfidential,
		Actor:          actor,
		Action:         string(event.Action),
		Target:         target,
		Scope:          audit.Scope{Platform: true},
		Outcome:        identityAuditOutcome(event.Succeeded),
		CorrelationID:  identityAuditCorrelation(string(event.Action), target.Reference, event.OccurredAt),
	})
}

// RecordServiceAccountEvent implements ServiceAccountAuditHook. Raw API credential
// material is structurally absent from ServiceAccountAuditEvent and therefore cannot
// enter audit records through this bridge.
func (adapter *AuditAdapter) RecordServiceAccountEvent(event ServiceAccountAuditEvent) {
	if adapter == nil || adapter.writer == nil || event.OccurredAt.IsZero() || event.Action == "" {
		return
	}
	actor := audit.Actor{Kind: identityAuditActorSystemKind, Reference: identityAuditActorSystem}
	target := audit.Target{Kind: "identity.service_account", Reference: string(event.Action)}
	if event.ServiceAccountID.Valid() {
		target.Reference = string(event.ServiceAccountID)
		if event.Action == ServiceAccountAuditCredentialVerified {
			actor = audit.Actor{Kind: identityAuditActorServiceAccount, Reference: string(event.ServiceAccountID)}
		}
	}
	if event.CredentialID.Valid() {
		target = audit.Target{Kind: "identity.api_credential", Reference: string(event.CredentialID)}
	}
	adapter.writeBestEffort(event.OccurredAt, audit.RecordInput{
		Classification: audit.ClassificationConfidential,
		Actor:          actor,
		Action:         string(event.Action),
		Target:         target,
		Scope:          audit.Scope{Platform: true},
		Outcome:        identityAuditOutcome(event.Succeeded),
		CorrelationID:  identityAuditCorrelation(string(event.Action), target.Reference, event.OccurredAt),
	})
}

func (adapter *AuditAdapter) writeBestEffort(_ time.Time, input audit.RecordInput) {
	_, _ = adapter.writer.Write(context.Background(), audit.RequirementBestEffort, input)
}

func identityAuditOutcome(succeeded bool) audit.Outcome {
	if succeeded {
		return audit.OutcomeSucceeded
	}
	return audit.OutcomeFailed
}

func identityAuditCorrelation(action, target string, occurredAt time.Time) string {
	return action + ":" + target + ":" + occurredAt.UTC().Format("20060102T150405.000000000Z")
}

var _ SecurityAuditHook = (*AuditAdapter)(nil)
var _ ServiceAccountAuditHook = (*AuditAdapter)(nil)
