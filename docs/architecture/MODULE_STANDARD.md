# Omnexa Module Standard

Status: **Mandatory baseline v1**

## 1. Goal

Every business capability in Omnexa must fit a repeatable module contract so the platform can grow without hidden coupling.

A module is not merely a folder. It is a versioned product boundary with declared ownership, dependencies, permissions, events, migrations, UI contributions and lifecycle behavior.

## 2. Module and submodule decomposition

Large modules are decomposed into governed submodules/capability families using:

- `docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md` — ordered S01-S10 execution template;
- `docs/roadmap/modules/SUBMODULE_CATALOG.json` — machine-readable P02-P27 submodule/family IDs;
- the owning dossier under `docs/roadmap/modules/` — architecture, flows, options and delivery order;
- `docs/governance/AI_MODULE_EXECUTION_PROTOCOL.md` — continuation/no-replanning/handoff rules.

Required hierarchy:

```text
Phase -> Module -> Submodule -> Work package -> Task -> Evidence
```

A submodule is normally an internal capability boundary of its owning module, not a separate product/domain owner. It may own its own focused write model, contracts, tests and lifecycle behavior only where those responsibilities remain consistent with the module's authoritative ownership.

Before implementation begins, each submodule plan must identify ownership, dependencies, data/schema, provided/consumed capabilities, permissions, tenant/security implications, UI/event/workflow contributions, migrations, lifecycle, ordered subtasks, acceptance criteria and evidence requirements.

The owning dossier must additionally define primary happy/failure/retry/degradation flows and the configurable option surface. If a requested option or flow is absent, resolve the planning gap before inventing runtime semantics.

When a preplanned decomposition exists, implementation agents continue from the next incomplete task instead of restarting architecture planning. Replanning is reserved for a genuine architecture conflict, changed requirement, missing owner/dependency, security/regulatory constraint or approved change-control event.

Builder-style modules must separate runtime schema/contracts from visual authoring UI. The visual builder is an editor for versioned server-validated definitions; it is not an undocumented source of runtime semantics.

Preplanning is not implementation authorization. `docs/roadmap/STATE.json` remains the sole executable phase/work-package lock.

## 3. Required module metadata

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

## 4. Module ownership

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

Every new persistence change must answer **who owns this write?** before schema/code is added.

## 5. Dependency classes

### 5.1 Required dependency
The module cannot function at all without the dependency. Keep these rare.

### 5.2 Optional dependency
The module exposes additional capability when another module is present, but degrades safely when it is absent.

### 5.3 Platform dependency
Kernel/platform contracts required by modules. These do not create domain coupling.

### 5.4 Forbidden dependency
Any dependency on another module's private package, internal table, undocumented endpoint, migration detail or implementation-specific field.

## 6. Capability contract

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

## 7. Event contract

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

## 8. Permission contract

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

## 9. Option, setting and policy contract

Every configurable option must declare:

- stable key/name;
- type and allowed values/bounds;
- default;
- scope (`platform`, `tenant`, `organization`, `module`, `user`, or explicit domain scope);
- read/change permission;
- sensitivity/data classification;
- audit requirement;
- effective timing (immediate, next request, restart, migration or scheduled effective date);
- compatibility/rollback behavior;
- classification as setting, feature flag or policy.

Security, tenancy, accounting and data-ownership invariants are not ordinary user-editable settings. Secrets are not feature flags.

## 10. Flow contract

Every non-trivial submodule must document at least:

1. happy path;
2. validation failure;
3. authorization/tenant denial where applicable;
4. required/optional dependency unavailable behavior;
5. timeout/cancellation;
6. duplicate/retry/idempotency behavior where applicable;
7. disable/re-enable/upgrade/restart behavior;
8. rollback/compensation when partial state is possible.

Builder/editor families additionally document author -> validate -> preview -> version -> publish/activate -> restore/degrade behavior.

## 11. UI contribution contract

Modules may contribute navigation, pages, widgets, builder blocks or settings through declared slots rather than patching unrelated UI directly.

Each UI contribution must define:

- slot/location;
- permission requirement;
- module availability condition;
- feature flag/entitlement condition;
- fallback behavior when dependencies are absent.

Builder submodules such as page builder, template builder, form builder, dashboard builder or custom-object builder must also provide keyboard-equivalent authoring operations and satisfy the repository browser accessibility/standards plan when UI implementation is authorized.

## 12. Data ownership

A module may write only:

- its own transactional data;
- kernel-owned metadata through kernel capabilities;
- explicit integration/read projections designed for that purpose.

Cross-module foreign keys should be considered carefully. Prefer stable platform IDs/references when hard database coupling would make independent lifecycle impossible.

Historical records requiring external context should store the immutable snapshot necessary for audit/business continuity when appropriate.

## 13. Module lifecycle hooks

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

## 14. Disable versus purge

Disable is non-destructive. It makes active features unavailable while preserving data required for future re-enable, reporting, audit or references.

Purge is destructive and must be explicit, authorized, audited and dependency-checked.

A module must never interpret normal uninstall/disable as permission to silently erase business evidence referenced elsewhere.

## 15. Migration rules

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

## 16. Failure behavior

A module must classify dependencies as:

- required and unavailable -> fail closed with explicit health error;
- optional and unavailable -> degrade feature selectively;
- external provider unavailable -> preserve local state and retry/compensate according to policy.

A module failure should not crash unrelated domains where isolation is feasible.

## 17. Testing contract

Every module/submodule is expected to provide, as applicable:

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
- idempotency/retry tests for async behavior;
- option/default/scope validation tests;
- documented flow-path tests including failure/cancellation paths;
- browser accessibility/standards evidence for authorized UI surfaces.

## 18. Versioning

Module versioning uses semantic intent:

- patch: compatible bug fix/internal improvement;
- minor: backward-compatible capability addition;
- major: breaking module behavior or public contract transition.

Public API/event/capability contract versions are independent where needed. A module major bump does not justify silently breaking all contracts.

## 19. Security declaration

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

## 20. Module acceptance gate

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
10. module/submodule ID exists in the canonical activated plan/catalog as applicable;
11. architecture, flows and options are reconciled with the owning dossier;
12. decomposition is reconciled with the S01-S10 execution blueprint;
13. documentation and roadmap ownership are reconciled.
