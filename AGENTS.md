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
5. `docs/architecture/GLOSSARY.md`
6. `docs/architecture/NAMING_STANDARD.md`
7. `docs/architecture/DOMAIN_OWNERSHIP.md`
8. `docs/architecture/DEPENDENCY_MATRIX.md`
9. `docs/roadmap/MASTER_PLAN.md`
10. `docs/roadmap/STATUS.md`
11. `docs/roadmap/STATE.json`
12. `docs/governance/AI_EXECUTION_POLICY.md`
13. `docs/governance/CHANGE_CONTROL.md`
14. `docs/governance/DEFINITION_OF_DONE.md`
15. Relevant ADRs under `docs/adr/`

If these documents conflict, stop and resolve the conflict through the change-control process before implementation.

## 3. Canonical execution state

`docs/roadmap/STATE.json` is the machine-readable canonical execution state.

Only work that is explicitly marked `active` may be implemented unless a dependency-blocking defect requires a narrowly scoped fix.

An AI agent MUST NOT:

- start a future phase early;
- silently add a new product domain;
- invent synonyms that conflict with the canonical glossary/naming standard;
- create a second authoritative owner for an existing concept;
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
2. **One authoritative owner.** Every authoritative write model/capability has one owning domain recorded by the ownership model.
3. **Domain ownership.** A module owns its write model and schema. Other modules may use contracts, events or approved read models only.
4. **No hidden coupling.** No module may depend on another module's internal tables, classes, private endpoints or undocumented behavior.
5. **Dependency direction is governed.** `DEPENDENCY_MATRIX.md` defines allowed cross-domain mechanisms; an `X` path requires redesign or an approved ADR, not a shortcut.
6. **Tenant-aware by default.** Tenant and organization boundaries must be explicit at data, authorization and service layers.
7. **Authorization-aware by default.** Every protected action goes through policy/capability checks.
8. **Auditable by default.** Security-sensitive and business-significant mutations emit attributable audit records.
9. **Versioned contracts.** Public APIs, events, module manifests and externally consumed schemas are versioned.
10. **Failure isolation.** Disabling or removing an optional module must not corrupt unrelated modules.
11. **Idempotent integration.** Events, webhooks, retries and background jobs must be safe to replay where required.
12. **AI uses governed capabilities.** AI agents never receive unrestricted database write access.
13. **Modular-monolith first.** Services are extracted only when measurable scale, isolation or team ownership justifies it.
14. **No speculative infrastructure.** Complexity must be earned by requirements and evidence.

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
3. Identify canonical terminology and the authoritative owning domain before creating a new entity/capability.
4. State affected modules/contracts/data ownership and allowed dependency direction.
5. Implement the smallest complete change satisfying the active acceptance criteria.
6. Add or update tests at the appropriate layer.
7. Run required quality gates from `DEFINITION_OF_DONE.md`.
8. Record evidence: tests, builds, migration checks, contract checks and relevant CI run IDs.
9. Update `STATUS.md` and `STATE.json` only when evidence proves the transition.
10. If architecture changed, add/update an ADR and reconcile all dependent documents in the same change.

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
- canonical domain ownership for an established concept;
- canonical terminology when semantics/ownership change materially;
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
- authoritative domain owner(s);
- dependency direction/mechanism;
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
- domain ownership is ambiguous or duplicated;
- an implementation uses a forbidden dependency path;
- status claims exceed evidence;
- documentation and machine-readable state disagree;
- runtime depends on hidden manual steps.

## 11. Repository and legal guardrails

- Follow `CONTRIBUTING.md` and `SECURITY.md`.
- Do not intentionally push implementation directly to `main`; use governed PRs.
- Hosted branch/ruleset targets are defined in `docs/governance/REPOSITORY_HARDENING.md`.
- The existing `LICENSE` file must not be replaced or treated as the final business licensing strategy without explicit owner authorization and the licensing/IP decision process.
- AI systems must never make trademark/legal ownership decisions autonomously.

## 12. Scope-drift rule

If a requested feature is valuable but outside the active work package, record it as planned backlog or propose a plan amendment. Do not absorb it silently into the current implementation.

The target is not maximum feature count. The target is a platform that can grow to very high feature count **without architectural decay**.
