package identity

import (
	"context"
	"time"
)

// Repository is the owner-bounded persistence contract for human User identity.
// It exposes no tenant, authentication/session, role or business-Person authority.
type Repository interface {
	Create(context.Context, User) error
	Get(context.Context, UserID) (User, error)
	Transition(context.Context, UserID, LifecycleState, LifecycleState, time.Time) (User, error)
}
