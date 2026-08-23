package tenancy

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository persists only kernel.tenancy-owned Tenant and membership state.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates the P02.02 owner-bounded PostgreSQL repository.
func NewPostgresRepository(pool *pgxpool.Pool) (*PostgresRepository, error) {
	if pool == nil {
		return nil, repositoryInvalidFailure()
	}
	return &PostgresRepository{pool: pool}, nil
}

// CreateTenant persists one authoritative Tenant record.
func (repository *PostgresRepository) CreateTenant(ctx context.Context, tenant Tenant) error {
	if repository == nil || repository.pool == nil {
		return repositoryInvalidFailure()
	}
	if validationErr := tenant.validate(); validationErr != nil {
		return validationErr
	}
	_, executionErr := repository.pool.Exec(
		ctx,
		`INSERT INTO omnexa_tenancy.tenants
		 (id, lifecycle_state, created_at, updated_at)
		 VALUES ($1, $2, $3, $4)`,
		string(tenant.id), string(tenant.state), tenant.createdAt, tenant.updatedAt,
	)
	if executionErr != nil {
		return tenantPersistenceFailure(executionErr)
	}
	return nil
}

// GetTenant retrieves one authoritative Tenant by UUIDv7 identifier.
func (repository *PostgresRepository) GetTenant(ctx context.Context, id TenantID) (Tenant, error) {
	if repository == nil || repository.pool == nil {
		return Tenant{}, repositoryInvalidFailure()
	}
	if !id.Valid() {
		return Tenant{}, invalidTenantFailure()
	}

	var storedID string
	var state string
	var createdAt time.Time
	var updatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`SELECT id, lifecycle_state, created_at, updated_at
		 FROM omnexa_tenancy.tenants
		 WHERE id = $1`,
		string(id),
	).Scan(&storedID, &state, &createdAt, &updatedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return Tenant{}, tenantNotFoundFailure()
	}
	if queryErr != nil {
		return Tenant{}, repositoryFailure(queryErr)
	}
	return rehydrateTenant(TenantID(storedID), TenantState(state), createdAt, updatedAt)
}

// TransitionTenant atomically applies one expected Tenant lifecycle change.
func (repository *PostgresRepository) TransitionTenant(
	ctx context.Context,
	id TenantID,
	from TenantState,
	to TenantState,
	changedAt time.Time,
) (Tenant, error) {
	if repository == nil || repository.pool == nil {
		return Tenant{}, repositoryInvalidFailure()
	}
	if !id.Valid() || !from.Valid() || !to.Valid() || !tenantTransitionAllowed(from, to) || changedAt.IsZero() {
		return Tenant{}, tenantTransitionFailure()
	}
	instant := changedAt.UTC()

	var storedID string
	var state string
	var createdAt time.Time
	var updatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_tenancy.tenants
		 SET lifecycle_state = $3, updated_at = $4
		 WHERE id = $1
		   AND lifecycle_state = $2
		   AND updated_at <= $4
		 RETURNING id, lifecycle_state, created_at, updated_at`,
		string(id), string(from), string(to), instant,
	).Scan(&storedID, &state, &createdAt, &updatedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		existing, existingErr := repository.GetTenant(ctx, id)
		if existingErr != nil {
			return Tenant{}, existingErr
		}
		if existing.state != from || instant.Before(existing.updatedAt) {
			return Tenant{}, tenantConflictFailure()
		}
		return Tenant{}, invalidStoredTenantFailure()
	}
	if queryErr != nil {
		return Tenant{}, tenantPersistenceFailure(queryErr)
	}
	return rehydrateTenant(TenantID(storedID), TenantState(state), createdAt, updatedAt)
}

// CreateMembership persists one kernel.tenancy-owned User/Tenant relationship.
// The identity principal is referenced but kernel.identity state is never written.
func (repository *PostgresRepository) CreateMembership(ctx context.Context, membership Membership) error {
	if repository == nil || repository.pool == nil {
		return repositoryInvalidFailure()
	}
	if validationErr := membership.validate(); validationErr != nil {
		return validationErr
	}
	_, executionErr := repository.pool.Exec(
		ctx,
		`INSERT INTO omnexa_tenancy.tenant_memberships
		 (id, tenant_id, principal_id, relationship_state, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		string(membership.id), string(membership.tenantID), string(membership.principalID),
		string(membership.state), membership.createdAt, membership.updatedAt,
	)
	if executionErr != nil {
		return membershipPersistenceFailure(executionErr)
	}
	return nil
}

// RevokeMembership terminally revokes one relationship inside an explicit Tenant
// boundary. Wrong-tenant, missing, already-revoked and stale requests all fail closed.
func (repository *PostgresRepository) RevokeMembership(
	ctx context.Context,
	tenantID TenantID,
	membershipID MembershipID,
	changedAt time.Time,
) (Membership, error) {
	if repository == nil || repository.pool == nil {
		return Membership{}, repositoryInvalidFailure()
	}
	if !tenantID.Valid() || !membershipID.Valid() || changedAt.IsZero() {
		return Membership{}, membershipTransitionFailure()
	}
	instant := changedAt.UTC()

	var storedID string
	var storedTenantID string
	var principalID string
	var state string
	var createdAt time.Time
	var updatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_tenancy.tenant_memberships
		 SET relationship_state = 'revoked', updated_at = $3
		 WHERE id = $1
		   AND tenant_id = $2
		   AND relationship_state = 'active'
		   AND updated_at <= $3
		 RETURNING id, tenant_id, principal_id, relationship_state, created_at, updated_at`,
		string(membershipID), string(tenantID), instant,
	).Scan(&storedID, &storedTenantID, &principalID, &state, &createdAt, &updatedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return Membership{}, membershipConflictFailure()
	}
	if queryErr != nil {
		return Membership{}, membershipPersistenceFailure(queryErr)
	}
	return rehydrateMembership(
		MembershipID(storedID),
		TenantID(storedTenantID),
		identity.UserID(principalID),
		MembershipState(state),
		createdAt,
		updatedAt,
	)
}

// ResolveContext derives an execution-scoped TrustedContext from authoritative
// membership and Tenant state. requestedTenantID is only a selector: changing it
// cannot create authority because an active persisted relationship must match.
func (repository *PostgresRepository) ResolveContext(
	ctx context.Context,
	principalID identity.UserID,
	requestedTenantID TenantID,
) (TrustedContext, error) {
	if repository == nil || repository.pool == nil {
		return TrustedContext{}, repositoryInvalidFailure()
	}
	if !principalID.Valid() || !requestedTenantID.Valid() {
		return TrustedContext{}, contextUntrustedFailure()
	}

	var tenantID string
	var tenantState string
	var tenantCreatedAt time.Time
	var tenantUpdatedAt time.Time
	var membershipID string
	var membershipTenantID string
	var membershipPrincipalID string
	var membershipState string
	var membershipCreatedAt time.Time
	var membershipUpdatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`SELECT t.id, t.lifecycle_state, t.created_at, t.updated_at,
		        m.id, m.tenant_id, m.principal_id, m.relationship_state, m.created_at, m.updated_at
		 FROM omnexa_tenancy.tenant_memberships AS m
		 JOIN omnexa_tenancy.tenants AS t ON t.id = m.tenant_id
		 WHERE m.principal_id = $1
		   AND m.tenant_id = $2
		   AND m.relationship_state = 'active'
		   AND t.lifecycle_state = 'active'`,
		string(principalID), string(requestedTenantID),
	).Scan(
		&tenantID, &tenantState, &tenantCreatedAt, &tenantUpdatedAt,
		&membershipID, &membershipTenantID, &membershipPrincipalID, &membershipState,
		&membershipCreatedAt, &membershipUpdatedAt,
	)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return TrustedContext{}, contextUntrustedFailure()
	}
	if queryErr != nil {
		return TrustedContext{}, repositoryFailure(queryErr)
	}

	tenant, tenantErr := rehydrateTenant(TenantID(tenantID), TenantState(tenantState), tenantCreatedAt, tenantUpdatedAt)
	if tenantErr != nil {
		return TrustedContext{}, contextUntrustedFailure()
	}
	membership, membershipErr := rehydrateMembership(
		MembershipID(membershipID),
		TenantID(membershipTenantID),
		identity.UserID(membershipPrincipalID),
		MembershipState(membershipState),
		membershipCreatedAt,
		membershipUpdatedAt,
	)
	if membershipErr != nil {
		return TrustedContext{}, contextUntrustedFailure()
	}
	return newTrustedContext(tenant, membership)
}
