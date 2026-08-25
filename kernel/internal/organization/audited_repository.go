package organization

import (
	"context"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

const organizationAuditActorReference = "kernel.organization"

// AuditedRepository decorates the kernel.organization owner repository with
// required, classification-safe audit delivery for material hierarchy/membership
// mutations. Reads and context resolution remain non-authorizing delegated calls.
type AuditedRepository struct {
	delegate Repository
	writer   *audit.Writer
}

// NewAuditedRepository returns the P02.10 organization mutation boundary.
func NewAuditedRepository(delegate Repository, writer *audit.Writer) (*AuditedRepository, error) {
	if delegate == nil || writer == nil {
		return nil, repositoryInvalidFailure()
	}
	return &AuditedRepository{delegate: delegate, writer: writer}, nil
}

func (repository *AuditedRepository) CreateNode(ctx context.Context, scope tenancy.Scope, node Node) error {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return repositoryInvalidFailure()
	}
	if err := repository.delegate.CreateNode(ctx, scope, node); err != nil {
		return err
	}
	return repository.auditMutation(ctx, scope, "organization.node.create", "organization.node", string(node.ID()))
}

func (repository *AuditedRepository) GetNode(ctx context.Context, scope tenancy.Scope, nodeID NodeID) (Node, error) {
	if repository == nil || repository.delegate == nil {
		return Node{}, repositoryInvalidFailure()
	}
	return repository.delegate.GetNode(ctx, scope, nodeID)
}

func (repository *AuditedRepository) MoveNode(ctx context.Context, scope tenancy.Scope, nodeID NodeID, parentID NodeID, at time.Time) (Node, error) {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return Node{}, repositoryInvalidFailure()
	}
	moved, err := repository.delegate.MoveNode(ctx, scope, nodeID, parentID, at)
	if err != nil {
		return Node{}, err
	}
	if err = repository.auditMutation(ctx, scope, "organization.node.move", "organization.node", string(moved.ID())); err != nil {
		return Node{}, err
	}
	return moved, nil
}

func (repository *AuditedRepository) Ancestors(ctx context.Context, scope tenancy.Scope, nodeID NodeID) ([]Node, error) {
	if repository == nil || repository.delegate == nil {
		return nil, repositoryInvalidFailure()
	}
	return repository.delegate.Ancestors(ctx, scope, nodeID)
}

func (repository *AuditedRepository) CreateMembership(ctx context.Context, scope tenancy.Scope, membership Membership) error {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return repositoryInvalidFailure()
	}
	if err := repository.delegate.CreateMembership(ctx, scope, membership); err != nil {
		return err
	}
	return repository.auditMutation(ctx, scope, "organization.membership.create", "organization.membership", string(membership.ID()))
}

func (repository *AuditedRepository) RevokeMembership(ctx context.Context, scope tenancy.Scope, membershipID MembershipID, at time.Time) (Membership, error) {
	if repository == nil || repository.delegate == nil || repository.writer == nil {
		return Membership{}, repositoryInvalidFailure()
	}
	revoked, err := repository.delegate.RevokeMembership(ctx, scope, membershipID, at)
	if err != nil {
		return Membership{}, err
	}
	if err = repository.auditMutation(ctx, scope, "organization.membership.revoke", "organization.membership", string(revoked.ID())); err != nil {
		return Membership{}, err
	}
	return revoked, nil
}

func (repository *AuditedRepository) ResolveContext(ctx context.Context, trusted tenancy.TrustedContext, nodeID NodeID) (ScopedContext, error) {
	if repository == nil || repository.delegate == nil {
		return ScopedContext{}, repositoryInvalidFailure()
	}
	return repository.delegate.ResolveContext(ctx, trusted, nodeID)
}

func (repository *AuditedRepository) auditMutation(ctx context.Context, scope tenancy.Scope, action, targetKind, targetReference string) error {
	if !scope.Valid() {
		return invalidScopeFailure()
	}
	_, err := repository.writer.Write(ctx, audit.RequirementRequired, audit.RecordInput{
		Classification: audit.ClassificationConfidential,
		Actor:          audit.Actor{Kind: "system", Reference: organizationAuditActorReference},
		Action:         action,
		Target:         audit.Target{Kind: targetKind, Reference: targetReference},
		Scope:          audit.Scope{TenantID: string(scope.TenantID())},
		Outcome:        audit.OutcomeSucceeded,
		CorrelationID:  action + ":" + targetReference,
		Reason:         "kernel.organization protected mutation",
		Privileged:     true,
	})
	return err
}

var _ Repository = (*AuditedRepository)(nil)
