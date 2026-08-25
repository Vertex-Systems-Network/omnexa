package authorization

import (
	"context"
	"time"
)

// These methods extend the existing P02.05 in-memory test double so prior tests
// continue exercising exactly their historical human-user surface. P02.08 service
// account behavior is covered separately by focused tests/integration evidence.
func (repository *fakeRepository) createServiceAccountAssignment(context.Context, ServiceAccountAssignment) error {
	return nil
}

func (repository *fakeRepository) getServiceAccountAssignment(context.Context, Scope, AssignmentID) (ServiceAccountAssignment, error) {
	return ServiceAccountAssignment{}, assignmentNotFoundFailure()
}

func (repository *fakeRepository) revokeServiceAccountAssignment(context.Context, Scope, AssignmentID, time.Time) (ServiceAccountAssignment, error) {
	return ServiceAccountAssignment{}, assignmentNotFoundFailure()
}

func (repository *fakeRepository) hasServiceAccountPermission(context.Context, ServiceAccountSubject, PermissionID) (bool, error) {
	return false, nil
}
