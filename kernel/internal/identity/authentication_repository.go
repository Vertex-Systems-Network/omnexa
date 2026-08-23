package identity

import (
	"context"
	"crypto/sha256"
	"time"
)

type authenticationSnapshot struct {
	principalID       UserID
	state             LifecycleState
	userUpdatedAt     time.Time
	passwordHash      string
	credentialVersion uint64
}

type sessionRecord struct {
	session           Session
	credentialVersion uint64
}

// AuthenticationRepository is the P02.04 owner-bounded persistence contract.
// Raw passwords and raw access/refresh secrets are intentionally absent.
type AuthenticationRepository interface {
	EnrollPassword(context.Context, UserID, string, time.Time) (uint64, error)
	AuthenticationSnapshot(context.Context, UserID) (authenticationSnapshot, error)
	ChangePassword(context.Context, UserID, uint64, string, time.Time) (uint64, error)
	CreateSession(
		context.Context,
		sessionRecord,
		[sha256.Size]byte,
		time.Time,
		[sha256.Size]byte,
		time.Time,
	) error
	AccessSession(context.Context, [sha256.Size]byte, time.Time) (sessionRecord, error)
	RefreshSession(context.Context, [sha256.Size]byte, time.Time) (sessionRecord, error)
	RotateRefresh(
		context.Context,
		SessionID,
		[sha256.Size]byte,
		[sha256.Size]byte,
		time.Time,
		[sha256.Size]byte,
		time.Time,
		time.Time,
	) (sessionRecord, error)
	RevokeSession(context.Context, UserID, SessionID, time.Time) (sessionRecord, error)
	ListSessions(context.Context, UserID) ([]Session, error)
}
