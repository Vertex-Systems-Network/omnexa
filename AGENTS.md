# Omnexa Repository Execution Contract

This file is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## 1. Mission

Omnexa is a **Composable Enterprise Business Operating System**, not a conventional ERP and not a collection of unrelated applications. ERP, CRM, commerce, POS, website/CMS, portals, payments, workflows, integrations, analytics, low-code and AI are domains running on one governed platform kernel.

The repository must evolve as one coherent platform with independently installable domain modules.

## 2. Mandatory read order before any change

Before proposing or modifying code, schema, infrastructure, APIs, tests or documentation, read these files in order:

1. `AGENTS.md`
2. `docs/governance/PRODUCT_CONSTITUTION.md`
3. `docs/architecture/SYSTEM_ARCHITECTURE.md`
4. `docs/architecture/MODULE_STANDARD.md`
5. `docs/roadmap/MASTER_PLAN.md`
6. `docs/roadmap/STATUS.md`
7. `docs/roadmap/STATE.json`
8. `docs/governance/CHANGE_CONTROL.md`
9. `docs/governance/DEFINITION_OF_DONE.md`
10. Relevant ADRs under `docs/adr/`

If these documents conflict, stop and resolve the conflict through the change-control process before implementation.

## 3. Canonical execution state

`docs/roadmap/STATE.json` is the machine-readable canonical execution state.

Only work that is explicitly marked `active` may be implemented unless a dependency-blocking defect requires a narrowly scoped fix.

An AI agent MUST NOT:

- start a future phase early;
- silently add a new product domain;
- replace an approved technology or architectural pattern without an ADR;
- couple modules through direct cross-module database writes;
- create a second implementation of an existing platform capability;
- bypass tenancy, authorization, audit, observability or versioning requirements;
- claim a task complete without acceptance evidence;
- mark status files manually inconsistent with tests or CI evidence;
- mix unrelated project code into this repository.

## 4. Architecture invariants

These rules are non-negotiable unless superseded by an approved ADR and corresponding plan revision.

1. **Kernel before business modules.** Shared platform capabilities belong in the kernel, not copied into domains.
2. **Domain ownership.** A module owns its write model and schema. Other modules may use contracts, events or approved read models only.
3. **No hidden coupling.** No module may depend on another module's internal tables, classes, private endpoints or undocumented behavior.
4. **Tenant-aware by default.** Tenant and organization boundaries must be explicit at data, authorization and service layers.
5. **Authorization-aware by default.** Every protected action goes through policy/capability checks.
6. **Auditable by default.** Security-sensitive and business-significant mutations emit attributable audit records.
7. **Versioned contracts.** Public APIs, events, module manifests and externally consumed schemas are versioned.
8. **Failure isolation.** Disabling or removing an optional module must not corrupt unrelated modules.
9. **Idempotent integration.** Events, webhooks, retries and background jobs must be safe to replay where required.
10. **AI uses governed capabilities.** AI agents never receive unrestricted database write access.
11. **Modular-monolith first.** Services are extracted only when measurable scale, isolation or team ownership justifies it.
12. **No speculative infrastructure.** Complexity must be earned by requirements and evidence.

## 5. Technology baseline

Until changed by ADR:

- Core backend/platform services: **Go**
- Web/admin/builder SDK surfaces: **TypeScript + React**
- Native/edge/security-sensitive components: **Rust only where justified**
- AI/data-science workers: **Python only where ecosystem value justifies it**
- Primary OLTP: **PostgreSQL**
- Cache/ephemeral coordination: **Redis-compatible layer**
- Object storage: **S3-compatible**
- Event/messaging baseline: **NATS/JetStream-class event fabric**
- Observability standard: **OpenTelemetry**

Technology choice does not authorize premature implementation. The active phase and its acceptance gates still control work.

## 6. Required work protocol

For every implementation task:

1. Identify the exact phase, work package and acceptance criteria.
2. Inspect current repository state before editing.
3. State affected modules/contracts/data ownership.
4. Implement the smallest complete change satisfying the active acceptance criteria.
5. Add or update tests at the appropriate layer.
6. Run required quality gates from `DEFINITION_OF_DONE.md`.
7. Record evidence: tests, builds, migration checks, contract checks and relevant CI run IDs.
8. Update `STATUS.md` and `STATE.json` only when evidence proves the transition.
9. If architecture changed, add/update an ADR and reconcile all dependent documents in the same change.

## 7. Phase and task state model

Allowed states:

- `planned`
- `ready`
- `active`
- `blocked`
- `verification`
- `done`
- `superseded`

Transitions must be evidence-backed. `done` means all acceptance gates passed, not merely that code exists.

Only one foundation phase should normally be `active` at a time until the kernel/module contracts are stable. Parallel domain execution is allowed only when the master plan explicitly opens parallel workstreams.

## 8. Change-control trigger

An ADR and master-plan reconciliation are mandatory for changes to any of the following:

- platform category or product boundary;
- language/runtime baseline;
- tenancy hierarchy;
- authorization model;
- module lifecycle or dependency model;
- cross-module communication model;
- event/versioning semantics;
- primary data ownership model;
- deployment topology baseline;
- security boundary;
- phase ordering or removal of a mandatory gate.

Do not implement the architectural change first and document it later.

## 9. Pull request discipline

Every PR must identify:

- phase/work package;
- scope and non-scope;
- architecture impact;
- data/migration impact;
- security/tenancy impact;
- contracts/events changed;
- test and CI evidence;
- rollback/compatibility considerations;
- status files updated.

Unrelated changes belong in separate PRs.

## 10. Definition of safe completion

A change is incomplete if any of these are true:

- build or required tests fail;
- migrations are non-repeatable or fresh-install path is broken;
- tenant isolation is untested where relevant;
- permission checks are missing;
- module uninstall/disable compatibility is broken where relevant;
- public contracts changed without versioning;
- status claims exceed evidence;
- documentation and machine-readable state disagree;
- runtime depends on hidden manual steps.

## 11. Scope-drift rule

If a requested feature is valuable but outside the active work package, record it as planned backlog or propose a plan amendment. Do not absorb it silently into the current implementation.

The target is not maximum feature count. The target is a platform that can grow to very high feature count **without architectural decay**.
