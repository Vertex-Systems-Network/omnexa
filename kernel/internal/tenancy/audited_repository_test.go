package tenancy

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
)

func TestAuditedRepositoryRecordsTenantMutation(t *testing.T) {
	delegate := &tenancyAuditFakeRepository{}
	sink, err := audit.NewMemorySink(8)
	if err != nil {
		t.Fatalf("audit.NewMemorySink() error = %v", err)
	}
	writer, err := audit.NewWriter(sink, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	repository, err := NewAuditedRepository(delegate, writer)
	if err != nil {
		t.Fatalf("NewAuditedRepository() error = %v", err)
	}
	tenant, err := NewTenant()
	if err != nil {
		t.Fatalf("NewTenant() error = %v", err)
	}
	if err = repository.CreateTenant(context.Background(), tenant); err != nil {
		t.Fatalf("CreateTenant() error = %v", err)
	}
	if delegate.createTenantCalls != 1 {
		t.Fatalf("delegate create calls = %d, want 1", delegate.createTenantCalls)
	}
	records := sink.Snapshot()
	if len(records) != 1 {
		t.Fatalf("audit records = %d, want 1", len(records))
	}
	record := records[0]
	if record.Action() != "tenancy.tenant.create" || record.Actor().Reference != tenancyAuditActorReference {
		t.Fatalf("audit action/actor = %q/%#v", record.Action(), record.Actor())
	}
	if record.Scope().TenantID != string(tenant.ID()) || record.Target().Reference != string(tenant.ID()) {
		t.Fatalf("audit scope/target = %#v/%#v", record.Scope(), record.Target())
	}
	if !record.Privileged() || len(record.Fields()) != 0 {
		t.Fatalf("audit privileged/fields = %v/%#v", record.Privileged(), record.Fields())
	}
}

func TestAuditedRepositoryRequiredAuditFailureCannotClaimSuccess(t *testing.T) {
	delegate := &tenancyAuditFakeRepository{}
	writer, err := audit.NewWriter(tenancyAuditFailingSink{}, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	repository, err := NewAuditedRepository(delegate, writer)
	if err != nil {
		t.Fatalf("NewAuditedRepository() error = %v", err)
	}
	tenant, err := NewTenant()
	if err != nil {
		t.Fatalf("NewTenant() error = %v", err)
	}
	if err = repository.CreateTenant(context.Background(), tenant); err == nil {
		t.Fatal("CreateTenant() unexpectedly claimed success after required audit failure")
	}
	if delegate.createTenantCalls != 1 {
		t.Fatalf("delegate create calls = %d, want 1", delegate.createTenantCalls)
	}
	health := writer.Health()
	if health.Submitted != 1 || health.Failed != 1 {
		t.Fatalf("writer health = %#v", health)
	}
}

type tenancyAuditFailingSink struct{}

func (tenancyAuditFailingSink) Append(context.Context, audit.Record) error {
	return errors.New("synthetic required audit failure")
}

type tenancyAuditFakeRepository struct {
	createTenantCalls int
}

func (repository *tenancyAuditFakeRepository) CreateTenant(_ context.Context, _ Tenant) error {
	repository.createTenantCalls++
	return nil
}
func (*tenancyAuditFakeRepository) GetTenant(context.Context, TenantID) (Tenant, error) {
	return Tenant{}, errors.New("not implemented")
}
func (*tenancyAuditFakeRepository) TransitionTenant(context.Context, TenantID, TenantState, TenantState, time.Time) (Tenant, error) {
	return Tenant{}, errors.New("not implemented")
}
func (*tenancyAuditFakeRepository) CreateMembership(context.Context, Membership) error {
	return errors.New("not implemented")
}
func (*tenancyAuditFakeRepository) RevokeMembership(context.Context, TenantID, MembershipID, time.Time) (Membership, error) {
	return Membership{}, errors.New("not implemented")
}
func (*tenancyAuditFakeRepository) ResolveContext(context.Context, identity.UserID, TenantID) (TrustedContext, error) {
	return TrustedContext{}, errors.New("not implemented")
}
