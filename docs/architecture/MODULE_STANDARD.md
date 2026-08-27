# Omnexa Module Standard

Status: **Mandatory baseline v1 with ADR-0012 dependency-contract evolution**

## 1. Goal

Every business capability in Omnexa must fit a repeatable module contract so the platform can grow without hidden coupling.

A module is not merely a folder. It is a versioned product boundary with declared ownership, dependencies, permissions, events, migrations, UI contributions and lifecycle behavior.

## 2. Required module metadata

Every installable module must expose a machine-readable manifest containing the canonical module metadata.

Schema v1 remains supported under its accepted P03.01 contract. Under accepted ADR-0012, schema v2 evolves only required/optional **module dependency** declarations into structured version requirements while preserving the rest of the module semantics.

Representative schema-v2 shape:

```yaml
schema_version: 2
id: omnexa.<domain>
name: Human Readable Name
version: 1.0.0
contract_version: 1
status: stable|beta|experimental
owner: team-or-domain
runtime: go|typescript|rust|python
required_platform_version: ">=1.0.0"
dependencies:
  - id: omnexa.example
    constraint: ">=1.2.0 <2.0.0"
optional_dependencies: []
platform_dependencies: []
capabilities_provided: []
capabilities_consumed: []
permissions: []
events_published: []
events_consumed: []
workflow_triggers: []
workflow_actions: []
ui_slots: []
settings: []
feature_flags: []
data_classification: []
migrations: []
lifecycle_hooks: []
health_checks: []
```

Schema evolution rules:

- schema v1 dependency strings retain their historical meaning and are never silently reinterpreted as version constraints;
- schema v2 `dependencies` and `optional_dependencies` contain exact `{id, constraint}` records;
- schema-version dispatch is explicit and version-specific; parser fallback between schema versions is forbidden;
- unknown fields remain rejected according to the selected schema;
- module dependency constraints use the bounded grammar accepted by ADR-0012;
- platform dependencies remain platform-capability identifiers and are not converted to module dependency records;
- `required_platform_version` remains the separate platform compatibility declaration.

The exact manifest syntax may evolve only through governed schema versions; the mandatory semantics must remain represented.

## 3. Module ownership

Each module owns:

- its domain invariants;
- commands/use cases;
- transactional write model;
- internal schema/tables;
- migrations for owned schema;
- authoritative API/capability contracts;
- events announcing its own facts;
- permission definitions for its capabilities;
- tests for its invariants;
- lifecycle compatibility behavior.

A module does not own kernel concerns such as tenant identity, global files, global audit transport, core policy runtime or event transport.

## 4. Dependency classes

### 4.1 Required dependency

The module cannot function at all without the dependency. Keep these rare.

For schema v2, every required module dependency must declare an explicit version constraint. Missing or incompatible required dependencies fail closed before install/enable eligibility.

### 4.2 Optional dependency

The module exposes additional capability when another module is present, but degrades safely when it is absent or incompatible.

For schema v2, every optional module dependency declares an explicit version constraint. Optional absence/incompatibility produces selective degradation metadata and does not globally invalidate unrelated modules.

Optional edges do not create required install/enable ordering authority and do not participate in the required-dependency global cycle gate.

### 4.3 Platform dependency

Kernel contracts required by modules. These do not create domain coupling. Platform dependency identifiers remain distinct from versioned module dependency records.

### 4.4 Forbidden dependency

Any dependency on another module's private package, internal table, undocumented endpoint, migration detail or implementation-specific field.

### 4.5 Dependency identity/version rules

- a module may not depend on itself in required or optional classes;
- duplicate dependency IDs within a class are invalid;
- the same module ID cannot appear in conflicting dependency classes;
- required dependency edges alone determine the release-blocking topological graph;
- one discovered module ID maps to at most one discovered version under the current registry contract;
- P03.03 does not perform multi-version/SAT/package-manager solving;
- resolver dependency declarations must come from the same normalized validated manifest snapshot that discovery used to create the registry identity; independently reparsed raw manifests are not resolver authority.

## 5. Capability contract

A provided capability must document:

- stable capability name;
- major version;
- request shape;
- response shape;
- authorization requirements;
- tenant/organization scope;
- validation and business errors;
- idempotency behavior where applicable;
- audit/event side effects;
- compatibility policy.

Example capability identifier:

```text
inventory.reserve-stock.v1
```

## 6. Event contract

A published event must declare:

- event name and major version;
- producer domain;
- trigger condition;
- payload schema;
- tenant/organization fields;
- actor/correlation/causation metadata;
- ordering assumptions, if any;
- replay/idempotency guidance;
- sensitive-field classification.

Consumers must treat events as facts, not remote procedure calls.

## 7. Permission contract

Permissions should use stable names:

```text
crm.contact.read
crm.contact.create
crm.contact.update
crm.contact.delete
finance.invoice.approve
payments.refund.execute
```

Permissions must be checked server-side at the owning capability boundary.

## 8. UI contribution contract

Modules may contribute navigation, pages, widgets, builder blocks or settings through declared slots rather than patching unrelated UI directly.

Each UI contribution must define:

- slot/location;
- permission requirement;
- module availability condition;
- feature flag/entitlement condition;
- fallback behavior when dependencies are absent.

## 9. Data ownership

A module may write only:

- its own transactional data;
- kernel-owned metadata through kernel capabilities;
- explicit integration/read projections designed for that purpose.

Cross-module foreign keys should be considered carefully. Prefer stable platform IDs/references when hard database coupling would make independent lifecycle impossible.

Historical records requiring external context should store the immutable snapshot necessary for audit/business continuity when appropriate.

## 10. Module lifecycle hooks

Modules may implement:

- `pre_install`
- `install`
- `post_install`
- `pre_enable`
- `enable`
- `disable`
- `upgrade`
- `rollback` where supported
- `archive`
- `export`
- `detach`
- `purge`
- `health_check`

Lifecycle handlers must be retry-aware and idempotent where execution can be repeated.

## 11. Disable versus purge

Disable is non-destructive. It makes active features unavailable while preserving data required for future re-enable, reporting, audit or references.

Purge is destructive and must be explicit, authorized, audited and dependency-checked.

A module must never interpret normal uninstall/disable as permission to silently erase business evidence referenced elsewhere.

## 12. Migration rules

Module migrations must:

- operate only on owned schema unless an approved platform migration says otherwise;
- be ordered and versioned;
- support fresh install from zero;
- support upgrade from every supported baseline;
- avoid environment-specific manual edits;
- include backfill strategy for new required fields;
- document destructive operations;
- be tested against representative data;
- preserve tenant boundaries.

Manifest schema migration is not a data-ownership grant. Dependency-bearing schema-v1 manifests must explicitly migrate to schema v2 before they can receive full required-dependency resolver eligibility under ADR-0012.

## 13. Failure behavior

A module must classify dependencies as:

- required and unavailable/incompatible -> fail closed with explicit safe compatibility/health error;
- optional and unavailable/incompatible -> degrade the affected feature selectively;
- external provider unavailable -> preserve local state and retry/compensate according to policy.

A module failure should not crash unrelated domains where isolation is feasible.

## 14. Testing contract

Every module is expected to provide, as applicable:

- domain unit tests;
- application/use-case tests;
- persistence integration tests;
- permission tests;
- tenant-isolation tests;
- contract tests for provided capabilities;
- event schema tests;
- migration/fresh-install tests;
- lifecycle enable/disable tests;
- optional-dependency degradation tests;
- dependency version-compatibility tests;
- required-cycle and deterministic-order tests;
- idempotency/retry tests for async behavior.

Manifest/resolver implementation must retain schema-v1 regressions and add schema-v2 parser, strict SemVer, constraint-bound, registry-binding and degradation fixtures required by ADR-0012.

## 15. Versioning

Module versioning uses semantic intent:

- patch: compatible bug fix/internal improvement;
- minor: backward-compatible capability addition;
- major: breaking module behavior or public contract transition.

Public API/event/capability contract versions are independent where needed. A module major bump does not justify silently breaking all contracts.

For schema-v2 dependency resolution, module versions and constraint operands use one strict SemVer 2.0.0 comparison path. Build metadata does not affect precedence. Constraint syntax is the bounded comparator grammar defined by ADR-0012 rather than package-manager shorthand.

## 16. Security declaration

Modules must declare:

- permissions required;
- sensitive data categories handled;
- secrets needed;
- external network destinations if applicable;
- webhooks/endpoints exposed;
- file types handled;
- privileged operations;
- AI tools/capabilities exposed.

Marketplace/third-party modules will eventually require signed packages and explicit installation consent for these scopes.

Dependency metadata never grants authorization, capability access, private package access or database authority.

## 17. Module acceptance gate

A new module is not accepted merely because it loads.

Minimum acceptance:

1. manifest valid for its declared schema version;
2. dependency graph valid and version-compatible under the accepted resolver contract;
3. owned schema/migrations pass from fresh install;
4. enable/disable behavior passes;
5. permission and tenant tests pass;
6. public contracts are documented/versioned;
7. events validate against schema;
8. no forbidden cross-module imports/writes;
9. health check reports accurately;
10. documentation and roadmap ownership are reconciled.

For P03.03, schema-v1 dependency-bearing compatibility behavior and schema-v2 dependency requirements are governed by ADR-0012. No implicit compatibility inference is permitted.
