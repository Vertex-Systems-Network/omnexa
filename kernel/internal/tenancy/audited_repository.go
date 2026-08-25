package tenancy

import (
	"context"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
)

const tenancyAuditActorReference = "kernel.tenancy"

// AuditedRepository decorates the kernel.tenancy owner repository with required,
// classification-safe audit delivery for material lifecycle mutations. Reads and
// trusted-context resolution remain delegated without creating audit read authority.
type AuditedRepository struct {
	delegate Repository
	writer   *audit.Writer
}

// NewAuditedRepository returns the P02.10 tenancy mutation boundary.
func NewAuditedRepository(delegate Repository, writer *audit.Writer) (*AuditedRepository, error) {
	if delegate == nil || writer == nil {
		return nil, repositoryInvalidFailure()
	}
	return &AuditedRepository{delegate: delegate, writer: writer}, nil
}

func (repository *AuditedRepository) CreateTenant(ctx context.Context, tenant Tenant) error {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return repositoryInvalidFailure()
	}
	if err := repository.delegate.CreateTenant(ctx, tenant); err != nil {
		return err
	}
	return repository.auditMutation(ctx, "tenancy.tenant.create", "tenancy.tenant", string(tenant.ID()), tenant.ID())
}

func (repository *AuditedRepository) GetTenant(ctx context.Context, tenantID TenantID) (Tenant, error) {
	if repository == nil || repository.delegate == nil {
		return Tenant{}, repositoryInvalidFailure()
	}
	return repository.delegate.GetTenant(ctx, tenantID)
}

func (repository *AuditedRepository) TransitionTenant(ctx context.Context, tenantID TenantID, from TenantState, to TenantState, at time.Time) (Tenant, error) {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return Tenant{}, repositoryInvalidFailure()
	}
	updated, err := repository.delegate.TransitionTenant(ctx, tenantID, from, to, at)
	if err != nil {
		return Tenant{}, err
	}
	if err = repository.auditMutation(ctx, "tenancy.tenant.transition", "tenancy.tenant", string(updated.ID()), updated.ID()); err != nil {
		return Tenant{}, err
	}
	return updated, nil
}

func (repository *AuditedRepository) CreateMembership(ctx context.Context, membership Membership) error {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return repositoryInvalidFailure()
	}
	if err := repository.delegate.CreateMembership(ctx, membership); err != nil {
		return err
	}
	return repository.auditMutation(ctx, "tenancy.membership.create", "tenancy.membership", string(membership.ID()), membership.TenantID())
}

func (repository *AuditedRepository) RevokeMembership(ctx context.Context, tenantID TenantID, membershipID MembershipID, at time.Time) (Membership, error) {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return Membership{}, repositoryInvalidFailure()
	}
	revoked, err := repository.delegate.RevokeMembership(ctx, tenantID, membershipID, at)
	if err != nil {
		return Membership{}, err
	}
	if err = repository.auditMutation(ctx, "tenancy.membership.revoke", "tenancy.membership", string(revoked.ID()), revoked.TenantID()); err != nil {
		return Membership{}, err
	}
	return revoked, nil
}

func (repository *AuditedRepository) ResolveContext(ctx context.Context, principalID identity.UserID, tenantID TenantID) (TrustedContext, error) {
	if repository == nil || repository.delegate == nil {
		return TrustedContext{}, repositoryInvalidFailure()
	}
	return repository.delegate.ResolveContext(ctx, principalID, tenantID)
}

func (repository *AuditedRepository) auditMutation(ctx context.Context, action, targetKind, targetReference string, tenantID TenantID) error {
	_, err := repository.writer.Write(ctx, audit.RequirementRequired, audit.RecordInput{
		Classification: audit.ClassificationConfidential,
		Actor:          audit.Actor{Kind: "system", Reference: tenancyAuditActorReference},
		Action:         action,
		Target:         audit.Target{Kind: targetKind, Reference: targetReference},
		Scope:          audit.Scope{TenantID: string(tenantID)},
		Outcome:        audit.OutcomeSucceeded,
		CorrelationID:  action + ":" + targetReference,
		Reason:         "kernel.tenancy protected mutation",
		Privileged:     true,
	})
	return err
}

var _ Repository = (*AuditedRepository)(nil)
