# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## Mission and execution state

Omnexa is a **Composable Enterprise Business Operating System**. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Only the package marked `active` is authorized except narrowly scoped dependency-blocking or repository-maintenance fixes.

During P00, kernel and business-feature implementation remain forbidden. Architecture, governance, specifications and narrow maintenance only.

## Mandatory read order

Before material work, read:

1. `AGENTS.md`
2. `docs/governance/PRODUCT_CONSTITUTION.md`
3. system/module architecture, glossary, naming, domain ownership and dependency matrix
4. identifier/money/time/locale/error standards
5. API and Event standards
6. Security Standard, Data Classification and Security Control Matrix
7. Testing, CI, Release and Quality Gate standards
8. `docs/development/REPOSITORY_STRUCTURE.md`
9. `docs/development/LOCAL_DEVELOPMENT.md`
10. `docs/development/TOOLCHAIN_STANDARD.md`
11. `docs/development/CONFIGURATION_STANDARD.md`
12. `docs/development/DEVELOPER_COMMANDS.md`
13. roadmap `MASTER_PLAN.md`, `STATUS.md`, `STATE.json`
14. AI Execution Policy, Change Control, Definition of Done
15. active CI exception if present
16. relevant accepted ADRs.

If canonical documents conflict, stop implementation and resolve through change control.

## Temporary hosted-CI exception

When `github_actions_ci.state = temporary_p00_exception`, ADR-0006 applies. Hosted Actions is `BLOCKED`/`NOT RUN`, never `PASS`. Manual evidence may advance **P00 documentation/specification-only** work. It cannot authorize P01+ executable merges or waive build, migration, test, security-scan or release gates.

## Architecture invariants

1. Kernel before business modules.
2. One authoritative owner per write model/capability.
3. Cross-module direct DB writes and private implementation imports are forbidden.
4. Cross-domain integration uses governed capabilities/APIs, events, workflows or approved read projections.
5. Tenant/org boundaries, authorization, audit, observability and versioned contracts are mandatory.
6. Optional-module failure/removal cannot corrupt unrelated domains.
7. Retriable jobs/events/integrations are idempotent where required.
8. AI acts through governed capabilities only; no unrestricted DB authority.
9. Strict modular monolith first; service extraction requires evidence + ADR.
10. Infrastructure complexity must be earned.

## Foundation primitives — P00.03

- UUIDv7 canonical IDs; PostgreSQL native `uuid`;
- tenant-owned state uses `tenant_id` where applicable;
- exact-decimal money + explicit currency; no binary floating point;
- UTC/`timestamptz` instants + IANA civil-time semantics;
- BCP 47 locale with independent country/currency/timezone and first-class RTL;
- stable safe structured errors with no secrets/SQL/stack traces.

## HTTP APIs — P00.04

- `/api/v{major}/{domain}/{resources}`;
- OpenAPI 3.2.0 canonical contracts;
- lowercase `snake_case` JSON;
- Problem Details + stable machine codes;
- cursor pagination;
- explicit idempotency/concurrency;
- explicit business actions;
- client tenant/org IDs never become authorization authority.

## Events — P00.05

- `<domain>.<subject>.<past_tense_fact>.v<major>`;
- producer-owned CloudEvents-compatible envelope with UUIDv7 identity;
- trusted producer-derived tenant context;
- at-least-once delivery + idempotent consumers;
- outbox/inbox or equivalent durable guarantees;
- no global ordering assumption;
- bounded retry/DLQ;
- replay preserves identity and cannot duplicate protected side effects.

## Security — P00.06

Data classes: `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`.

- tenant-owned business data defaults to at least `CONFIDENTIAL` unless explicitly published;
- secrets/private keys/auth equivalents are `RESTRICTED`;
- production sensitive data does not flow to ordinary dev/test;
- authn never substitutes for authz;
- authz = RBAC + relationships + contextual policy + bounded capabilities;
- tenant isolation spans OLTP/cache/files/search/analytics/events/backups/AI/vector data;
- privileged support/export/purge/replay/financial/high-impact AI actions require explicit policy/audit;
- webhooks/integrations/devices/modules/AI are independent trust boundaries;
- prompt/tool injection cannot create authority.

## Quality — P00.07

Quality is repository-owned and CI-provider independent.

Gate classes: `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package, `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never treat blocked/unrun/N/A as PASS.

Affected risks require positive and negative evidence: cross-tenant denial, authz deny, duplicate/idempotency/replay, fresh+upgrade migrations, optional-module lifecycle/degradation, security boundaries. Flaky tests are defects. Releases prefer immutable build-once/promote artifacts with source SHA and gate evidence.

## Repository and local development — P00.08

Canonical monorepo ownership roots are:

```text
apps/
kernel/
modules/
platform/
shared/
infrastructure/
scripts/
docs/
generated/
```

Rules:

- folder != microservice; deployment extraction requires evidence/ADR;
- module private implementation/schema/migrations stay with the owning module;
- `shared/` is for truly platform-wide contracts/primitives, not hidden domain logic;
- generated output is reproducible derivative output, never source of truth;
- default local infrastructure = containerized PostgreSQL + Redis-compatible + NATS/JetStream + S3-compatible storage; add more only when active phases require them;
- Kubernetes is not required for the default developer loop;
- local and CI use the same semantic verification commands;
- toolchains/dependencies are repository-pinned; unversioned global tools and conflicting JS lockfiles are forbidden;
- config precedence is explicit; secrets are separate from committed config;
- production sensitive data is prohibited by default locally; use synthetic deterministic multi-tenant fixtures;
- Linux is canonical backend execution; macOS supported where upstream tooling permits; Windows backend development prefers WSL2; native Windows is a separate certification target where required;
- supported workflows must not depend on hidden manual SQL/file/UI steps.

Reserved semantic command families include `omnexa dev`, `omnexa db`, `omnexa verify` and `omnexa module`, as defined in `DEVELOPER_COMMANDS.md`.

## Technology baseline

Go backend/core; TypeScript+React web/admin/builder/SDK; Rust justified edge/native/security; Python justified AI/data; PostgreSQL OLTP; Redis-compatible cache; S3-compatible storage; NATS/JetStream-class messaging; OpenTelemetry observability.

## Required work protocol

For every material change:

1. identify active phase/package and acceptance criteria;
2. inspect repository state;
3. identify canonical terminology and authoritative owner;
4. apply primitive/API/event/security/quality/development standards;
5. map risk to G0-G8 gates;
6. preserve repository/module boundaries;
7. implement the smallest complete authorized change;
8. add required positive + negative evidence;
9. execute canonical local/CI verification unless an active P00 exception records hosted execution BLOCKED/NOT RUN;
10. record evidence accurately and reconcile STATUS/STATE;
11. document architectural changes through ADR/change control before implementation.

## Forbidden behavior

Do not start future phases; silently add domains; duplicate ownership; invent conflicting contracts/security/quality/toolchain conventions; cross-write/import module-private internals; bypass tenancy/authz/audit/classification; grant AI private write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; call blocked CI PASS; rely on hidden global/latest tools; require Kubernetes for ordinary local dev without an ADR; or mix unrelated project code.

## Change-control triggers

ADR + plan reconciliation is required for material changes to product boundary, domain ownership, primitive/API/event/security semantics, quality/evidence/release semantics, monorepo ownership model, local execution/toolchain/config model, module lifecycle/dependencies, deployment topology or phase gates.

## Pull request evidence

Every material PR states phase/package, scope/non-scope, owners/dependencies, primitive/API/event/security impacts, quality gate mapping, repository/toolchain/config impacts, data/migrations, exact evidence states, compatibility/rollback and state reconciliation.

## Hosted/legal guardrails

- issue #3: `main` protection still needs hosted configuration;
- issue #14: GitHub Actions quota/runner currently blocked;
- issue #4: final licensing/IP/trademark strategy remains unresolved before external distribution/public launch.

None authorize early implementation.

## Scope drift

Useful work outside the active package is recorded for later or proposed through change control. The objective is a platform that can grow without architectural decay.