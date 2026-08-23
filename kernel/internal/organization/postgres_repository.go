package organization

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository persists only kernel.organization-owned hierarchy and scoped
// relationship state. P02.02 tenant context is consumed through a narrow resolver.
type PostgresRepository struct {
	pool           *pgxpool.Pool
	tenantContexts TenantContextResolver
}

// NewPostgresRepository creates the P02.03 owner-bounded PostgreSQL repository.
func NewPostgresRepository(pool *pgxpool.Pool, tenantContexts TenantContextResolver) (*PostgresRepository, error) {
	if pool == nil || tenantContexts == nil {
		return nil, repositoryInvalidFailure()
	}
	return &PostgresRepository{pool: pool, tenantContexts: tenantContexts}, nil
}

// CreateNode persists one hierarchy node inside an explicit trusted Tenant scope.
func (repository *PostgresRepository) CreateNode(ctx context.Context, scope tenancy.Scope, node Node) error {
	if repository == nil || repository.pool == nil || repository.tenantContexts == nil {
		return repositoryInvalidFailure()
	}
	if !scope.Valid() || scope.TenantID() != node.tenantID {
		return scopeDeniedFailure()
	}
	if validationErr := node.validate(); validationErr != nil {
		return validationErr
	}
	if node.parentID != nil {
		parent, parentErr := repository.GetNode(ctx, scope, *node.parentID)
		if parentErr != nil {
			return hierarchyParentInvalidFailure()
		}
		if !parentKindAllowed(node.kind, parent.kind) {
			return hierarchyTransitionFailure()
		}
	}

	_, executionErr := repository.pool.Exec(
		ctx,
		`INSERT INTO omnexa_organization.nodes
		 (id, tenant_id, node_kind, parent_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		string(node.id),
		string(node.tenantID),
		string(node.kind),
		optionalNodeIDString(node.parentID),
		node.createdAt,
		node.updatedAt,
	)
	if executionErr != nil {
		return nodePersistenceFailure(executionErr)
	}
	return nil
}

func optionalNodeIDString(id *NodeID) any {
	if id == nil {
		return nil
	}
	return string(*id)
}

// GetNode retrieves one node only inside the supplied trusted Tenant scope.
func (repository *PostgresRepository) GetNode(ctx context.Context, scope tenancy.Scope, id NodeID) (Node, error) {
	if repository == nil || repository.pool == nil || repository.tenantContexts == nil {
		return Node{}, repositoryInvalidFailure()
	}
	if !scope.Valid() || !id.Valid() {
		return Node{}, scopeDeniedFailure()
	}
	return scanNode(repository.pool.QueryRow(
		ctx,
		`SELECT id, tenant_id, node_kind, parent_id::text, created_at, updated_at
		 FROM omnexa_organization.nodes
		 WHERE id = $1 AND tenant_id = $2`,
		string(id),
		string(scope.TenantID()),
	))
}

type rowScanner interface {
	Scan(...any) error
}

func scanNode(row rowScanner) (Node, error) {
	var storedID string
	var tenantID string
	var kind string
	var parent pgtype.Text
	var createdAt time.Time
	var updatedAt time.Time
	scanErr := row.Scan(&storedID, &tenantID, &kind, &parent, &createdAt, &updatedAt)
	if errors.Is(scanErr, pgx.ErrNoRows) {
		return Node{}, nodeNotFoundFailure()
	}
	if scanErr != nil {
		return Node{}, repositoryFailure(scanErr)
	}
	var parentID *NodeID
	if parent.Valid {
		value := NodeID(parent.String)
		parentID = &value
	}
	return rehydrateNode(
		NodeID(storedID),
		tenancy.TenantID(tenantID),
		NodeKind(kind),
		parentID,
		createdAt,
		updatedAt,
	)
}

// MoveNode changes one non-root node parent. The table lock serializes hierarchy
// mutations so two concurrent moves cannot jointly introduce a cycle.
func (repository *PostgresRepository) MoveNode(
	ctx context.Context,
	scope tenancy.Scope,
	nodeID NodeID,
	newParentID NodeID,
	changedAt time.Time,
) (Node, error) {
	if repository == nil || repository.pool == nil || repository.tenantContexts == nil {
		return Node{}, repositoryInvalidFailure()
	}
	if !scope.Valid() || !nodeID.Valid() || !newParentID.Valid() || changedAt.IsZero() {
		return Node{}, hierarchyTransitionFailure()
	}
	if nodeID == newParentID {
		return Node{}, hierarchyCycleFailure()
	}

	tx, beginErr := repository.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if beginErr != nil {
		return Node{}, repositoryFailure(beginErr)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, lockErr := tx.Exec(ctx, `LOCK TABLE omnexa_organization.nodes IN SHARE ROW EXCLUSIVE MODE`); lockErr != nil {
		return Node{}, repositoryFailure(lockErr)
	}

	node, nodeErr := scanNode(tx.QueryRow(
		ctx,
		`SELECT id, tenant_id, node_kind, parent_id::text, created_at, updated_at
		 FROM omnexa_organization.nodes
		 WHERE id = $1 AND tenant_id = $2
		 FOR UPDATE`,
		string(nodeID),
		string(scope.TenantID()),
	))
	if nodeErr != nil {
		return Node{}, nodeErr
	}
	if node.kind == NodeKindOrganization || changedAt.UTC().Before(node.updatedAt) {
		return Node{}, hierarchyTransitionFailure()
	}

	parent, parentErr := scanNode(tx.QueryRow(
		ctx,
		`SELECT id, tenant_id, node_kind, parent_id::text, created_at, updated_at
		 FROM omnexa_organization.nodes
		 WHERE id = $1 AND tenant_id = $2
		 FOR UPDATE`,
		string(newParentID),
		string(scope.TenantID()),
	))
	if parentErr != nil {
		return Node{}, hierarchyParentInvalidFailure()
	}
	if !parentKindAllowed(node.kind, parent.kind) {
		return Node{}, hierarchyTransitionFailure()
	}

	var createsCycle bool
	cycleErr := tx.QueryRow(
		ctx,
		`WITH RECURSIVE descendants(id) AS (
		     SELECT id
		     FROM omnexa_organization.nodes
		     WHERE tenant_id = $1 AND parent_id = $2
		     UNION
		     SELECT child.id
		     FROM omnexa_organization.nodes AS child
		     JOIN descendants AS ancestor ON child.parent_id = ancestor.id
		     WHERE child.tenant_id = $1
		 )
		 SELECT EXISTS (SELECT 1 FROM descendants WHERE id = $3)`,
		string(scope.TenantID()),
		string(nodeID),
		string(newParentID),
	).Scan(&createsCycle)
	if cycleErr != nil {
		return Node{}, repositoryFailure(cycleErr)
	}
	if createsCycle {
		return Node{}, hierarchyCycleFailure()
	}

	instant := changedAt.UTC()
	updated, updateErr := scanNode(tx.QueryRow(
		ctx,
		`UPDATE omnexa_organization.nodes
		 SET parent_id = $3, updated_at = $4
		 WHERE id = $1 AND tenant_id = $2 AND updated_at <= $4
		 RETURNING id, tenant_id, node_kind, parent_id::text, created_at, updated_at`,
		string(nodeID),
		string(scope.TenantID()),
		string(newParentID),
		instant,
	))
	if updateErr != nil {
		return Node{}, updateErr
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return Node{}, repositoryFailure(commitErr)
	}
	return updated, nil
}

// Ancestors returns nearest-parent-first ancestry and fails closed if persisted
// data is missing, cross-tenant or cyclic.
func (repository *PostgresRepository) Ancestors(ctx context.Context, scope tenancy.Scope, id NodeID) ([]Node, error) {
	if repository == nil || repository.pool == nil || repository.tenantContexts == nil {
		return nil, repositoryInvalidFailure()
	}
	current, currentErr := repository.GetNode(ctx, scope, id)
	if currentErr != nil {
		return nil, currentErr
	}
	seen := map[NodeID]struct{}{current.id: {}}
	ancestors := make([]Node, 0, 4)
	for current.parentID != nil {
		parentID := *current.parentID
		if _, exists := seen[parentID]; exists {
			return nil, hierarchyCycleFailure()
		}
		parent, parentErr := repository.GetNode(ctx, scope, parentID)
		if parentErr != nil {
			return nil, hierarchyParentInvalidFailure()
		}
		if !parentKindAllowed(current.kind, parent.kind) {
			return nil, hierarchyTransitionFailure()
		}
		seen[parentID] = struct{}{}
		ancestors = append(ancestors, parent)
		current = parent
	}
	return ancestors, nil
}

// CreateMembership persists one scoped User relationship. A live P02.02
// Tenant/User relationship for the target principal is a structural prerequisite,
// not an authorization grant.
func (repository *PostgresRepository) CreateMembership(
	ctx context.Context,
	scope tenancy.Scope,
	membership Membership,
) error {
	if repository == nil || repository.pool == nil || repository.tenantContexts == nil {
		return repositoryInvalidFailure()
	}
	if !scope.Valid() || membership.tenantID != scope.TenantID() {
		return scopeDeniedFailure()
	}
	if validationErr := membership.validate(); validationErr != nil {
		return validationErr
	}
	node, nodeErr := repository.GetNode(ctx, scope, membership.scopeID)
	if nodeErr != nil || node.tenantID != membership.tenantID {
		return scopeDeniedFailure()
	}
	if _, tenantContextErr := repository.tenantContexts.ResolveContext(ctx, membership.principalID, membership.tenantID); tenantContextErr != nil {
		return scopeDeniedFailure()
	}

	_, executionErr := repository.pool.Exec(
		ctx,
		`INSERT INTO omnexa_organization.scoped_memberships
		 (id, tenant_id, scope_id, principal_id, relationship_state, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		string(membership.id),
		string(membership.tenantID),
		string(membership.scopeID),
		string(membership.principalID),
		string(membership.state),
		membership.createdAt,
		membership.updatedAt,
	)
	if executionErr != nil {
		return membershipPersistenceFailure(executionErr)
	}
	return nil
}

// RevokeMembership terminally revokes one relationship inside a trusted Tenant.
func (repository *PostgresRepository) RevokeMembership(
	ctx context.Context,
	scope tenancy.Scope,
	membershipID MembershipID,
	changedAt time.Time,
) (Membership, error) {
	if repository == nil || repository.pool == nil || repository.tenantContexts == nil {
		return Membership{}, repositoryInvalidFailure()
	}
	if !scope.Valid() || !membershipID.Valid() || changedAt.IsZero() {
		return Membership{}, membershipConflictFailure()
	}

	var storedID string
	var tenantID string
	var scopeID string
	var principalID string
	var state string
	var createdAt time.Time
	var updatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_organization.scoped_memberships
		 SET relationship_state = 'revoked', updated_at = $3
		 WHERE id = $1
		   AND tenant_id = $2
		   AND relationship_state = 'active'
		   AND updated_at <= $3
		 RETURNING id, tenant_id, scope_id, principal_id, relationship_state, created_at, updated_at`,
		string(membershipID),
		string(scope.TenantID()),
		changedAt.UTC(),
	).Scan(&storedID, &tenantID, &scopeID, &principalID, &state, &createdAt, &updatedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return Membership{}, membershipConflictFailure()
	}
	if queryErr != nil {
		return Membership{}, membershipPersistenceFailure(queryErr)
	}
	return rehydrateMembership(
		MembershipID(storedID),
		tenancy.TenantID(tenantID),
		NodeID(scopeID),
		identity.UserID(principalID),
		MembershipState(state),
		createdAt,
		updatedAt,
	)
}

// ResolveContext derives a non-authorizing organization relationship context from
// a fresh P02.02 trusted Tenant context and active scoped membership.
func (repository *PostgresRepository) ResolveContext(
	ctx context.Context,
	trusted tenancy.TrustedContext,
	requestedScopeID NodeID,
) (ScopedContext, error) {
	if repository == nil || repository.pool == nil || repository.tenantContexts == nil {
		return ScopedContext{}, repositoryInvalidFailure()
	}
	if !trusted.Valid() || !requestedScopeID.Valid() {
		return ScopedContext{}, contextUntrustedFailure()
	}
	refreshed, refreshErr := repository.tenantContexts.ResolveContext(ctx, trusted.PrincipalID(), trusted.TenantID())
	if refreshErr != nil || refreshed.MembershipID() != trusted.MembershipID() {
		return ScopedContext{}, contextUntrustedFailure()
	}
	tenantScope, tenantScopeErr := refreshed.ScopeFor(refreshed.TenantID())
	if tenantScopeErr != nil {
		return ScopedContext{}, contextUntrustedFailure()
	}
	if _, nodeErr := repository.GetNode(ctx, tenantScope, requestedScopeID); nodeErr != nil {
		return ScopedContext{}, contextUntrustedFailure()
	}

	var membershipID string
	var tenantID string
	var scopeID string
	var principalID string
	var state string
	var createdAt time.Time
	var updatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`SELECT id, tenant_id, scope_id, principal_id, relationship_state, created_at, updated_at
		 FROM omnexa_organization.scoped_memberships
		 WHERE tenant_id = $1
		   AND scope_id = $2
		   AND principal_id = $3
		   AND relationship_state = 'active'`,
		string(refreshed.TenantID()),
		string(requestedScopeID),
		string(refreshed.PrincipalID()),
	).Scan(&membershipID, &tenantID, &scopeID, &principalID, &state, &createdAt, &updatedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return ScopedContext{}, contextUntrustedFailure()
	}
	if queryErr != nil {
		return ScopedContext{}, repositoryFailure(queryErr)
	}
	membership, membershipErr := rehydrateMembership(
		MembershipID(membershipID),
		tenancy.TenantID(tenantID),
		NodeID(scopeID),
		identity.UserID(principalID),
		MembershipState(state),
		createdAt,
		updatedAt,
	)
	if membershipErr != nil {
		return ScopedContext{}, contextUntrustedFailure()
	}
	return newScopedContext(refreshed, membership)
}
