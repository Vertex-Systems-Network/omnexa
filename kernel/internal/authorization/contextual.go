package authorization

import (
	"context"
	"regexp"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
)

const (
	maxObjectKindRunes     = 64
	contextAuditAction     = "authorization.context.evaluate"
	contextAuditTargetKind = "authorization.object"
	contextAuditReason     = "contextual authorization decision"
)

var objectKindPattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]{0,63}$`)

// ObjectKind is a stable opaque object class reference. It does not declare a
// business-domain policy or grant authority over objects of that class.
type ObjectKind string

// Valid reports whether kind is safe bounded reference metadata.
func (kind ObjectKind) Valid() bool {
	return len(kind) <= maxObjectKindRunes && objectKindPattern.MatchString(string(kind))
}

// ObjectID is the canonical UUIDv7 reference of one governed object.
type ObjectID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id ObjectID) Valid() bool {
	parsed, parseErr := uuid.Parse(string(id))
	return parseErr == nil && parsed.Version() == 7
}

// ObjectReference is an opaque object reference. Identifiers are inputs to
// relationship resolution and never authorization evidence by themselves.
type ObjectReference struct {
	kind ObjectKind
	id   ObjectID
}

// NewObjectReference creates a validated object reference with no authority.
func NewObjectReference(kind ObjectKind, id ObjectID) (ObjectReference, error) {
	object := ObjectReference{kind: kind, id: id}
	if !object.Valid() {
		return ObjectReference{}, invalidContextFailure()
	}
	return object, nil
}

// Kind returns the opaque object class.
func (object ObjectReference) Kind() ObjectKind { return object.kind }

// ID returns the canonical object UUIDv7 reference.
func (object ObjectReference) ID() ObjectID { return object.id }

// Valid reports whether the object reference is canonical.
func (object ObjectReference) Valid() bool { return object.kind.Valid() && object.id.Valid() }

// Equal reports exact object-reference equality.
func (object ObjectReference) Equal(other ObjectReference) bool {
	return object.Valid() && other.Valid() && object.kind == other.kind && object.id == other.id
}

// AccessKind distinguishes normal resource reads from stricter sensitive-field
// and export boundaries. The distinction chooses a permission; it never grants one.
type AccessKind string

const (
	AccessRead           AccessKind = "read"
	AccessSensitiveField AccessKind = "sensitive_field"
	AccessExport         AccessKind = "export"
)

func (access AccessKind) valid() bool {
	return access == AccessRead || access == AccessSensitiveField || access == AccessExport
}

func (access AccessKind) privileged() bool {
	return access == AccessSensitiveField || access == AccessExport
}

// CapabilityBoundary maps one governed object action to three distinct stable
// permissions. P03 may later register module permissions; this type only consumes
// PermissionID values and does not register a capability vocabulary.
type CapabilityBoundary struct {
	readPermission      PermissionID
	sensitivePermission PermissionID
	exportPermission    PermissionID
}

// NewCapabilityBoundary requires distinct read/sensitive/export permissions so
// ordinary read authority can never imply sensitive-field or export authority.
func NewCapabilityBoundary(readPermission, sensitivePermission, exportPermission PermissionID) (CapabilityBoundary, error) {
	boundary := CapabilityBoundary{
		readPermission:      readPermission,
		sensitivePermission: sensitivePermission,
		exportPermission:    exportPermission,
	}
	if !boundary.Valid() {
		return CapabilityBoundary{}, invalidContextFailure()
	}
	return boundary, nil
}

// Valid reports whether all permission hooks are valid and pairwise distinct.
func (boundary CapabilityBoundary) Valid() bool {
	return boundary.readPermission.Valid() && boundary.sensitivePermission.Valid() && boundary.exportPermission.Valid() &&
		boundary.readPermission != boundary.sensitivePermission &&
		boundary.readPermission != boundary.exportPermission &&
		boundary.sensitivePermission != boundary.exportPermission
}

func (boundary CapabilityBoundary) requiredPermission(access AccessKind) (PermissionID, error) {
	if !boundary.Valid() || !access.valid() {
		return "", invalidContextFailure()
	}
	switch access {
	case AccessRead:
		return boundary.readPermission, nil
	case AccessSensitiveField:
		return boundary.sensitivePermission, nil
	case AccessExport:
		return boundary.exportPermission, nil
	default:
		return "", invalidContextFailure()
	}
}

// CallerOrigin is descriptive execution origin only. No value is privileged and
// internal/background execution receives no authorization shortcut.
type CallerOrigin string

const (
	CallerInteractive CallerOrigin = "interactive"
	CallerInternal    CallerOrigin = "internal"
	CallerBackground  CallerOrigin = "background"
)

func (origin CallerOrigin) valid() bool {
	return origin == CallerInteractive || origin == CallerInternal || origin == CallerBackground
}

// ContextMetadata is bounded safe metadata for contextual authorization audit.
type ContextMetadata struct {
	CorrelationID string
}

func (metadata ContextMetadata) validate() error {
	if !validReference(metadata.CorrelationID, maxCorrelationRunes) {
		return invalidContextFailure()
	}
	return nil
}

// ContextRequest is one exact subject/object/capability authorization request.
type ContextRequest struct {
	subject  Subject
	object   ObjectReference
	boundary CapabilityBoundary
	access   AccessKind
	origin   CallerOrigin
	metadata ContextMetadata
}

// NewContextRequest creates one validated contextual authorization request.
func NewContextRequest(
	subject Subject,
	object ObjectReference,
	boundary CapabilityBoundary,
	access AccessKind,
	origin CallerOrigin,
	metadata ContextMetadata,
) (ContextRequest, error) {
	request := ContextRequest{
		subject:  subject,
		object:   object,
		boundary: boundary,
		access:   access,
		origin:   origin,
		metadata: metadata,
	}
	if request.validate() != nil {
		return ContextRequest{}, invalidContextFailure()
	}
	return request, nil
}

func (request ContextRequest) validate() error {
	if !request.subject.Valid() || !request.object.Valid() || !request.boundary.Valid() || !request.access.valid() || !request.origin.valid() {
		return invalidContextFailure()
	}
	return request.metadata.validate()
}

// Subject returns the trusted P02.05 authorization subject.
func (request ContextRequest) Subject() Subject { return request.subject }

// Object returns the opaque object reference being governed.
func (request ContextRequest) Object() ObjectReference { return request.object }

// Access returns the requested field/export access mode.
func (request ContextRequest) Access() AccessKind { return request.access }

// Origin returns descriptive caller origin. It never implies authority.
func (request ContextRequest) Origin() CallerOrigin { return request.origin }

// RelationshipQuery asks a trusted resolver for current principal/object
// relationship evidence. Requested tenant/org IDs are intentionally absent:
// authoritative scope comes back from the resolver rather than from client input.
type RelationshipQuery struct {
	principalID identity.UserID
	object      ObjectReference
}

// PrincipalID returns the principal whose current relationship is requested.
func (query RelationshipQuery) PrincipalID() identity.UserID { return query.principalID }

// Object returns the exact object reference to resolve.
func (query RelationshipQuery) Object() ObjectReference { return query.object }

func (query RelationshipQuery) valid() bool {
	return query.principalID.Valid() && query.object.Valid()
}

// RelationshipEvidence is current resolver-owned evidence binding a principal
// and object to one exact tenant or organization scope.
type RelationshipEvidence struct {
	principalID    identity.UserID
	tenantID       tenancy.TenantID
	organizationID organization.NodeID
	object         ObjectReference
}

// NewTenantRelationshipEvidence creates current tenant-scoped resolver evidence.
func NewTenantRelationshipEvidence(
	principalID identity.UserID,
	tenantID tenancy.TenantID,
	object ObjectReference,
) (RelationshipEvidence, error) {
	evidence := RelationshipEvidence{principalID: principalID, tenantID: tenantID, object: object}
	if !evidence.valid() {
		return RelationshipEvidence{}, invalidContextFailure()
	}
	return evidence, nil
}

// NewOrganizationRelationshipEvidence creates current organization-scoped resolver evidence.
func NewOrganizationRelationshipEvidence(
	principalID identity.UserID,
	tenantID tenancy.TenantID,
	organizationID organization.NodeID,
	object ObjectReference,
) (RelationshipEvidence, error) {
	evidence := RelationshipEvidence{
		principalID:    principalID,
		tenantID:       tenantID,
		organizationID: organizationID,
		object:         object,
	}
	if !evidence.valid() || !organizationID.Valid() {
		return RelationshipEvidence{}, invalidContextFailure()
	}
	return evidence, nil
}

func (evidence RelationshipEvidence) valid() bool {
	return evidence.principalID.Valid() && evidence.tenantID.Valid() && evidence.object.Valid() &&
		(evidence.organizationID == "" || evidence.organizationID.Valid())
}

func (evidence RelationshipEvidence) matches(query RelationshipQuery, scope Scope) bool {
	if !evidence.valid() || !query.valid() || !scope.Valid() {
		return false
	}
	if evidence.principalID != query.principalID || !evidence.object.Equal(query.object) || evidence.tenantID != scope.TenantID() {
		return false
	}
	switch scope.Kind() {
	case ScopeTenant:
		return evidence.organizationID == ""
	case ScopeOrganization:
		return evidence.organizationID == scope.OrganizationID()
	default:
		return false
	}
}

// RelationshipResolution is the resolver result. Related=false is a normal
// deterministic deny result, not an error or disclosure of why no relationship exists.
type RelationshipResolution struct {
	related  bool
	evidence RelationshipEvidence
}

// NoRelationship returns a disclosure-safe negative relationship result.
func NoRelationship() RelationshipResolution { return RelationshipResolution{} }

// RelationshipFound returns validated current relationship evidence.
func RelationshipFound(evidence RelationshipEvidence) (RelationshipResolution, error) {
	if !evidence.valid() {
		return RelationshipResolution{}, invalidContextFailure()
	}
	return RelationshipResolution{related: true, evidence: evidence}, nil
}

func (resolution RelationshipResolution) matches(query RelationshipQuery, scope Scope) bool {
	return resolution.related && resolution.evidence.matches(query, scope)
}

// RelationshipResolver is a trusted server-side capability implemented by the
// owner of current object relationship facts. IDs in the query are references,
// and only a resolver result can satisfy this part of policy evaluation.
type RelationshipResolver interface {
	Resolve(context.Context, RelationshipQuery) (RelationshipResolution, error)
}

// ContextFacts are exact already-RBAC/relationship-bounded facts supplied to a
// server-side constraint evaluator. The evaluator may narrow authority only.
type ContextFacts struct {
	subject Subject
	object  ObjectReference
	access  AccessKind
	origin  CallerOrigin
}

// Subject returns the trusted subject already accepted by direct RBAC.
func (facts ContextFacts) Subject() Subject { return facts.subject }

// Object returns the exact relationship-bounded object.
func (facts ContextFacts) Object() ObjectReference { return facts.object }

// Access returns the requested access mode.
func (facts ContextFacts) Access() AccessKind { return facts.access }

// Origin returns descriptive caller origin with no bypass semantics.
func (facts ContextFacts) Origin() CallerOrigin { return facts.origin }

func (facts ContextFacts) valid() bool {
	return facts.subject.Valid() && facts.object.Valid() && facts.access.valid() && facts.origin.valid()
}

// ContextConstraintEvaluator applies server-side contextual narrowing after
// direct RBAC and trusted relationship evidence have both succeeded. Returning
// true cannot grant authority that those earlier checks denied.
type ContextConstraintEvaluator interface {
	Evaluate(context.Context, ContextFacts) (bool, error)
}

// ContextService is the P02.06 relationship/context-aware layer around the
// accepted P02.05 direct-RBAC Service. It owns no parallel role or permission path.
type ContextService struct {
	rbac          *Service
	relationships RelationshipResolver
	constraints   ContextConstraintEvaluator
}

// NewContextService creates the P02.06 policy layer. Every dependency is
// required; missing relationship or constraint infrastructure fails closed.
func NewContextService(
	rbac *Service,
	relationships RelationshipResolver,
	constraints ContextConstraintEvaluator,
) (*ContextService, error) {
	if rbac == nil || rbac.repository == nil || rbac.audit == nil || relationships == nil || constraints == nil {
		return nil, serviceInvalidFailure()
	}
	return &ContextService{rbac: rbac, relationships: relationships, constraints: constraints}, nil
}

// Check evaluates RBAC, trusted object relationship and contextual narrowing in
// that order. No role name or caller origin can skip an earlier requirement.
func (service *ContextService) Check(ctx context.Context, request ContextRequest) (Decision, error) {
	if service == nil || service.rbac == nil || service.relationships == nil || service.constraints == nil {
		return DecisionDeny, serviceInvalidFailure()
	}
	if requestErr := request.validate(); requestErr != nil {
		return DecisionDeny, requestErr
	}
	requiredPermission, permissionErr := request.boundary.requiredPermission(request.access)
	if permissionErr != nil {
		return DecisionDeny, permissionErr
	}
	directDecision, directErr := service.rbac.Check(ctx, request.subject, requiredPermission)
	if directErr != nil {
		return DecisionDeny, directErr
	}
	if !directDecision.Allowed() {
		return service.auditedDeny(ctx, request)
	}

	query := RelationshipQuery{principalID: request.subject.PrincipalID(), object: request.object}
	resolution, relationshipErr := service.relationships.Resolve(ctx, query)
	if relationshipErr != nil {
		return DecisionDeny, relationshipResolutionFailure(relationshipErr)
	}
	if !resolution.matches(query, request.subject.Scope()) {
		return service.auditedDeny(ctx, request)
	}

	facts := ContextFacts{subject: request.subject, object: request.object, access: request.access, origin: request.origin}
	if !facts.valid() {
		return DecisionDeny, invalidContextFailure()
	}
	satisfied, constraintErr := service.constraints.Evaluate(ctx, facts)
	if constraintErr != nil {
		return DecisionDeny, constraintEvaluationFailure(constraintErr)
	}
	if !satisfied {
		return service.auditedDeny(ctx, request)
	}

	if request.access.privileged() {
		if auditErr := service.auditDecision(ctx, request, audit.OutcomeSucceeded); auditErr != nil {
			return DecisionDeny, auditErr
		}
	}
	return DecisionAllow, nil
}

// Require returns one disclosure-safe authorization denial unless all P02.06
// requirements allow the governed action.
func (service *ContextService) Require(ctx context.Context, request ContextRequest) error {
	decision, checkErr := service.Check(ctx, request)
	if checkErr != nil {
		return checkErr
	}
	if !decision.Allowed() {
		return deniedFailure()
	}
	return nil
}

func (service *ContextService) auditedDeny(ctx context.Context, request ContextRequest) (Decision, error) {
	if auditErr := service.auditDecision(ctx, request, audit.OutcomeDenied); auditErr != nil {
		return DecisionDeny, auditErr
	}
	return DecisionDeny, nil
}

func (service *ContextService) auditDecision(ctx context.Context, request ContextRequest, outcome audit.Outcome) error {
	scope := audit.Scope{TenantID: string(request.subject.Scope().TenantID())}
	if request.subject.Scope().Kind() == ScopeOrganization {
		scope.OrganizationID = string(request.subject.Scope().OrganizationID())
	}
	_, writeErr := service.rbac.audit.Write(ctx, audit.RequirementRequired, audit.RecordInput{
		Classification: audit.ClassificationInternal,
		Actor: audit.Actor{
			Kind:      actorKindUser,
			Reference: string(request.subject.PrincipalID()),
		},
		Action: contextAuditAction,
		Target: audit.Target{
			Kind:      contextAuditTargetKind,
			Reference: string(request.object.ID()),
		},
		Scope:         scope,
		Outcome:       outcome,
		CorrelationID: request.metadata.CorrelationID,
		Reason:        contextAuditReason,
		Privileged:    request.access.privileged(),
	})
	return writeErr
}
