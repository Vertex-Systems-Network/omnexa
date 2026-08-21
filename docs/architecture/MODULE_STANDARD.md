# Omnexa Module Standard

Status: **Mandatory baseline v1**

## 1. Goal

Every business capability in Omnexa must fit a repeatable module contract so the platform can grow without hidden coupling.

A module is not merely a folder. It is a versioned product boundary with declared ownership, dependencies, permissions, events, migrations, UI contributions and lifecycle behavior.

## 2. Required module metadata

Every installable module must expose a machine-readable manifest containing at least:

```yaml
id: omnexa.<domain>
name: Human Readable Name
version: 1.0.0
contract_version: 1
status: stable|beta|experimental
owner: team-or-domain
runtime: go|typescript|rust|python
required_platform_version: ">=x.y.z"
dependencies: []
optional_dependencies: []
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

The exact manifest syntax may evolve, but these semantics must remain represented.

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

### 4.2 Optional dependency
The module exposes additional capability when another module is present, but degrades safely when it is absent.

### 4.3 Platform dependency
Kernel contracts required by every module. These do not create domain coupling.

### 4.4 Forbidden dependency
Any dependency on another module's private package, internal table, undocumented endpoint, migration detail or implementation-specific field.

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

## 13. Failure behavior

A module must classify dependencies as:

- required and unavailable -> fail closed with explicit health error;
- optional and unavailable -> degrade feature selectively;
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
- idempotency/retry tests for async behavior.

## 15. Versioning

Module versioning uses semantic intent:

- patch: compatible bug fix/internal improvement;
- minor: backward-compatible capability addition;
- major: breaking module behavior or public contract transition.

Public API/event/capability contract versions are independent where needed. A module major bump does not justify silently breaking all contracts.

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

## 17. Module acceptance gate

A new module is not accepted merely because it loads.

Minimum acceptance:

1. manifest valid;
2. dependency graph valid;
3. owned schema/migrations pass from fresh install;
4. enable/disable behavior passes;
5. permission and tenant tests pass;
6. public contracts are documented/versioned;
7. events validate against schema;
8. no forbidden cross-module imports/writes;
9. health check reports accurately;
10. documentation and roadmap ownership are reconciled.
