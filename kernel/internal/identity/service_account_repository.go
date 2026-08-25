package identity

import (
	"context"
	"crypto/sha256"
	"time"
)

// ServiceAccountRepository is the owner-bounded P02.08 persistence contract.
// Raw API credentials are intentionally absent; only non-reversible digests cross
// this boundary after one-time secret issuance.
type ServiceAccountRepository interface {
	CreateServiceAccount(context.Context, ServiceAccount) error
	GetServiceAccount(context.Context, ServiceAccountID) (ServiceAccount, error)
	TransitionServiceAccount(context.Context, ServiceAccountID, LifecycleState, LifecycleState, time.Time) (ServiceAccount, error)
	CreateAPICredential(context.Context, APICredential, [sha256.Size]byte) error
	AuthenticateAPICredential(context.Context, APICredentialID, [sha256.Size]byte, time.Time) (ServiceAccount, APICredential, error)
	RotateAPICredential(
		context.Context,
		ServiceAccountID,
		APICredentialID,
		APICredential,
		[sha256.Size]byte,
		time.Time,
	) (APICredential, error)
	RevokeAPICredential(context.Context, ServiceAccountID, APICredentialID, time.Time) (APICredential, error)
}
