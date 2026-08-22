# Foundation Module Dossiers — P02 to P06

Status: **Mandatory future planning baseline**

These plans predefine submodule boundaries, flows and option surfaces. They do not authorize P02-P06 implementation while `STATE.json` remains in P01.

## P02 — Identity, Tenancy & Organization

Architecture: one platform identity/security context shared by all later modules. Business modules may consume identity/organization references but may not create parallel auth/session/role stacks.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P02.A | Identity Accounts | invite/register/provision -> verify -> activate -> suspend/archive | registration policy, identity status rules, profile-edit policy |
| P02.B | Authentication & Sessions | credential/passkey/OIDC input -> authenticate -> policy -> session -> revoke/expire | session TTL, device/session limits, re-auth windows, provider enablement |
| P02.C | Tenant Lifecycle | create -> initialize -> activate -> suspend -> archive/purge governance | tenant naming, lifecycle policy, isolation mode, default locale/timezone |
| P02.D | Organization Hierarchy | legal entity/company -> branch/team/location relationships -> membership | hierarchy types, allowed nesting, default org context |
| P02.E | RBAC | permission registry -> role composition -> assignment -> evaluation | custom roles, role templates, deny/allow composition policy |
| P02.F | Context-aware Authorization | actor + tenant/org + resource + action -> policy decision | scope inheritance, ownership rules, break-glass policy hooks |
| P02.G | MFA & Passkeys | enroll -> verify -> recover -> revoke | required populations, recovery factors, step-up triggers |
| P02.H | Service Accounts & API Credentials | create service identity -> scope -> issue/rotate/revoke credential | credential lifetime, IP/network policy hooks, allowed capabilities |
| P02.I | Tenant Settings | resolve platform defaults -> tenant override -> org/user-safe projections | override scopes, validation, inheritance |
| P02.J | Identity Security Audit | auth/session/role/credential facts -> immutable audit stream | retention/export hooks, alert classifications |

Authoritative writes: P02 owns users/service identities, sessions, tenants, org hierarchy, memberships, role assignments and identity-security audit facts. Business domains store stable references, not copies of mutable authorization truth.

Permissions: identity.admin, tenant.manage, organization.manage, role.manage, credential.manage, session.revoke, security.audit.read plus narrower self-service capabilities.

Events: identity.activated/suspended, session.revoked, tenant.activated/suspended, organization.changed, membership.changed, role.assignment.changed, credential.rotated/revoked.

Security/test gates: cross-tenant denial, object/scope permission matrix, session invalidation, privilege escalation negatives, credential redaction, MFA recovery abuse cases, idempotent invitation/provisioning.

Delivery order: A -> C/D -> E/F -> B/G -> H -> I -> J, with contract tests after each ownership boundary.

## P03 — Module Runtime

Architecture: modules are versioned product boundaries registered into one modular-monolith runtime. Registry data is platform-owned; each module retains private schema and migrations.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P03.A | Manifest Schema | package -> parse -> validate -> normalized manifest | contract version, runtime type, declared scopes |
| P03.B | Registry & Discovery | discover/install candidate -> register -> query capabilities | discovery roots, allowed publishers/source classes |
| P03.C | Dependency Resolver | requested module -> required/optional graph -> cycle/conflict check -> plan | version constraints, optional degradation policy |
| P03.D | Lifecycle State Machine | install -> enable -> disable -> upgrade -> archive/detach/purge | lifecycle policy, destructive-operation confirmation |
| P03.E | Module Settings & Flags | schema -> validated values -> scoped resolution | setting scope, secret classification, feature defaults |
| P03.F | Capability Registry | manifest declarations -> register -> resolve provider -> invoke boundary | capability version policy, provider priority only where allowed |
| P03.G | Permission Registry | module permissions -> validation -> platform policy registration | naming/version rules, default role templates |
| P03.H | UI Contribution Registry | module slot declaration -> permission/availability filter -> render contribution | allowed slots, fallback/degradation behavior |
| P03.I | Migration Ownership Registry | module -> owned schema/migrations -> ordered execution | owner namespace, supported upgrade baselines |
| P03.J | Module Health | lifecycle/dependency/provider state -> health projection | required vs optional dependency severity |
| P03.K | Package Trust Hooks | package metadata -> signature/trust hook -> install policy | trust source, signature requirement hooks for P22 |

Forbidden: direct imports of another module private package, writes to another module tables, module-specific auth stacks, lifecycle deletion of referenced business evidence without explicit purge policy.

Events: module.installed/enabled/disabled/upgraded/archived/purged, capability.available/unavailable, module.health.changed.

Delivery order: A -> B/C -> F/G/I -> D -> E/H/J -> K -> reference-module lifecycle matrix.

## P04 — Data, Jobs & Event Fabric

Architecture: reliable asynchronous communication with explicit ownership, at-least-once delivery and idempotent consumers. Event facts never grant cross-domain write authority.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P04.A | Event Envelope | domain fact -> versioned CloudEvents-compatible envelope | schema version, correlation/causation fields |
| P04.B | Publish/Subscribe Boundary | producer -> transport abstraction -> subscribers | transport adapter, delivery timeout |
| P04.C | Durable Stream/Consumer | append -> consumer group -> ack/checkpoint | retention, partition/key policy, concurrency bounds |
| P04.D | Transactional Outbox | domain transaction -> outbox row -> dispatcher -> transport | batch size, retry/backoff, dispatch lease |
| P04.E | Inbox/Deduplication | received event -> identity check -> inbox -> handler -> outcome | dedupe window/retention, handler timeout |
| P04.F | Idempotency Primitives | operation key -> claim -> execute -> record/replay result | TTL/retention by operation class |
| P04.G | Retry/Backoff | retryable failure -> policy -> scheduled retry | attempts, exponential/jitter bounds |
| P04.H | Dead-letter/Quarantine | terminal failure -> quarantine -> inspect/retry/discard governance | thresholds, retention, operator permissions |
| P04.I | Correlation/Causation | request/workflow -> propagated IDs across jobs/events | generation/propagation rules |
| P04.J | Event Schema Registry | schema publish -> compatibility check -> version registry | compatibility mode, deprecation window |
| P04.K | Background Work Context | job dispatch -> tenant/actor/correlation context -> worker validation | context fields, lease/heartbeat bounds |

Required flows: successful publish; process crash between domain commit and dispatch; duplicate delivery; poison event; consumer restart; replay; retry exhaustion.

Delivery order: A/J -> B/C -> D/E/F -> G/H/I/K -> replay/fault-injection evidence.

## P05 — Omnexa Flow / Workflow OS

Architecture: versioned workflow definitions separated from durable workflow instances. Workflow actions invoke governed capabilities; they do not write foreign domain tables.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P05.A | Definition & Versioning | author definition -> validate -> version -> activate | version activation, draft retention |
| P05.B | Trigger Registry | event/manual/schedule trigger -> normalize -> start instance | trigger filters, duplicate-start key |
| P05.C | Action Registry | action declaration -> capability binding -> input/output validation | action timeout, retry class |
| P05.D | Conditions & Branches | state/context -> expression -> deterministic path | expression limits, null/error behavior |
| P05.E | Timers & Waits | schedule wait -> durable timer -> resume | timezone semantics, max duration |
| P05.F | Retry & Timeout | step failure -> retry policy -> terminal outcome | attempts, backoff, timeout |
| P05.G | Human Approvals | create task -> notify surface -> approve/reject/expire -> resume | approver rules, escalation/expiry |
| P05.H | Parallel Paths | fork -> independent branches -> join policy | all/any/quorum join behavior |
| P05.I | Compensation/Saga | committed steps -> downstream failure -> compensation plan | compensation order, retry/terminal policy |
| P05.J | State Persistence | transition -> durable checkpoint -> resume after restart | checkpoint cadence, history retention |
| P05.K | Audit/Timeline | transition/action/actor -> append history | redaction/export hooks |
| P05.L | Simulation/Test API | definition + synthetic context -> dry-run/simulate -> report | side-effect policy, fixture set |
| P05.M | Visual Designer Contract | schema/registry -> authoring API contract for later UI | node palette, validation feedback |

Required flows: start -> branch -> wait -> approval -> action -> retry -> compensation -> completion; process restart at each durable boundary; duplicate trigger/start.

Delivery order: A/J -> B/C/D -> E/F -> G/H/I -> K/L -> M contract only.

## P06 — Universal Business Foundation

Architecture: reusable references/value models, not a giant shared business database. Each later domain remains authoritative for its own transaction model.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P06.A | Party References | create/reference person or organization identity -> consume stable reference | party type extensibility, status rules |
| P06.B | Addresses & Contact Points | validate/store contact point -> link to owning subject | address/contact type registry, verification hooks |
| P06.C | Locations | define location reference -> hierarchy/geocode hook -> consume | location type, timezone/default locale |
| P06.D | Money & Currency | amount + ISO currency -> exact operations/conversion boundary | rounding policy by owning domain, FX provider hook |
| P06.E | Tax Context | subject/location/product context -> tax adapter input | jurisdiction/type classifications, provider hooks |
| P06.F | Product/Service References | stable reference to sellable/procurable subject without owning commerce catalog | reference class, lifecycle snapshot rules |
| P06.G | Files/Documents/Notes/Activities References | domain record -> governed reference -> metadata projection | allowed link types, visibility/classification |
| P06.H | Numbering & References | request namespace -> allocate deterministic/business reference | prefix/sequence reset scope, collision policy |
| P06.I | Import/Export Primitives | source -> map/validate -> staged import -> commit/report | batch size, error policy, dry-run |
| P06.J | Search/Index Integration | authoritative change -> projection/index hook -> query reference | index class, freshness target, rebuild policy |

Data rule: shared primitives may own their own reference records, but domain-specific invoice/customer/order/employee/product state stays in its domain.

Delivery order: A/B/C -> D/E -> F/G/H -> I/J, with at least two reference domains proving no direct-write coupling.

## Common foundation option governance

All configurable options in P02-P06 must declare scope (`platform`, `tenant`, `organization`, `module`, `user` where allowed), default, validation, sensitivity, change permission, audit requirement and whether a change is immediate or requires restart/migration. Secrets must never be represented as ordinary feature flags.

## Common evidence

Every activated submodule requires ownership/dependency validation, unit tests, persistence/provider integration where applicable, authorization/tenant negatives once P02 exists, migration fresh/upgrade evidence, error/redaction checks, lifecycle/resilience tests, code-quality/security/build gates and exact PR/CI evidence.