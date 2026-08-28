package authorization

import (
	"context"
	"errors"
	"testing"
)

func TestModulePermissionRequiresLiveAvailabilityAndExistingGrant(t *testing.T) {
	ctx := context.Background()
	repository := newFakeRepository()
	writer, _ := testAuditWriter(t)
	subject := testTenantSubject(
		"01890f3e-7b9a-7cc0-98c4-dc0c0c0739b1",
		"01890f3e-7b9a-7cc0-98c4-dc0c0c0739c1",
	)
	permission := PermissionID("inventory.stock.read")
	seedGrant(t, repository, subject, []PermissionID{permission})

	legacyService, err := NewService(repository, writer)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	decision, err := legacyService.Check(ctx, subject, permission)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("module permission without P03.07 provider = %q, %v; want deny, nil", decision, err)
	}

	unavailableService, err := NewServiceWithModulePermissionAvailability(
		repository,
		writer,
		staticModulePermissionAvailability{permission: permission, available: false},
	)
	if err != nil {
		t.Fatalf("NewServiceWithModulePermissionAvailability(unavailable) error = %v", err)
	}
	decision, err = unavailableService.Check(ctx, subject, permission)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("unavailable module permission = %q, %v; want deny, nil", decision, err)
	}

	availableService, err := NewServiceWithModulePermissionAvailability(
		repository,
		writer,
		staticModulePermissionAvailability{permission: permission, available: true},
	)
	if err != nil {
		t.Fatalf("NewServiceWithModulePermissionAvailability(available) error = %v", err)
	}
	decision, err = availableService.Check(ctx, subject, permission)
	if err != nil || decision != DecisionAllow {
		t.Fatalf("available permission with existing exact grant = %q, %v; want allow, nil", decision, err)
	}

	unassigned := testTenantSubject(
		"01890f3e-7b9a-7cc0-98c4-dc0c0c0739b2",
		"01890f3e-7b9a-7cc0-98c4-dc0c0c0739c1",
	)
	decision, err = availableService.Check(ctx, unassigned, permission)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("availability without role grant = %q, %v; want deny, nil", decision, err)
	}
}

func TestModulePermissionAvailabilityFailureFailsClosed(t *testing.T) {
	repository := newFakeRepository()
	writer, _ := testAuditWriter(t)
	subject := testTenantSubject(
		"01890f3e-7b9a-7cc0-98c4-dc0c0c0739b3",
		"01890f3e-7b9a-7cc0-98c4-dc0c0c0739c1",
	)
	permission := PermissionID("inventory.stock.read")
	seedGrant(t, repository, subject, []PermissionID{permission})
	availabilityErr := errors.New("synthetic lifecycle dependency failure")
	service, err := NewServiceWithModulePermissionAvailability(
		repository,
		writer,
		staticModulePermissionAvailability{permission: permission, err: availabilityErr},
	)
	if err != nil {
		t.Fatalf("NewServiceWithModulePermissionAvailability() error = %v", err)
	}
	decision, checkErr := service.Check(context.Background(), subject, permission)
	if !errors.Is(checkErr, availabilityErr) || decision != DecisionDeny {
		t.Fatalf("availability failure = %q, %v; want deny and original dependency error", decision, checkErr)
	}
}

func TestModulePermissionDefinitionRejectsKernelReservedPermission(t *testing.T) {
	definition := ModulePermissionDefinition{
		Permission: PermissionRoleRead,
		ModuleID:   "omnexa.inventory",
		Owner:      "inventory.team",
		Available:  true,
	}
	if !validModulePermissionDefinition(definition) {
		t.Fatal("definition shape should be structurally valid before reserved-kernel check")
	}
	if !kernelPermission(definition.Permission) {
		t.Fatal("kernel permission was not recognized as reserved")
	}
}

type staticModulePermissionAvailability struct {
	permission PermissionID
	available  bool
	err        error
}

func (availability staticModulePermissionAvailability) PermissionAvailable(
	_ context.Context,
	permission PermissionID,
) (bool, error) {
	if availability.err != nil {
		return false, availability.err
	}
	if permission != availability.permission {
		return false, nil
	}
	return availability.available, nil
}
