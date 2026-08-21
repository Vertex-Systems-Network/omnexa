# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed, modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are planned as domain families running on one shared platform kernel.

> **Current execution lock:** the repository is in **P00 — Product Constitution & Architecture Freeze**. Kernel and business-feature implementation must not begin until the P00 exit gate is complete. Current package: **P00.09 — Initial threat model and operational SLO targets**.

> **Temporary CI note:** GitHub Actions allowance is exhausted/disabled. ADR-0006 permits temporary P00 documentation/specification manual evidence only. Hosted CI remains `BLOCKED`/`NOT RUN`, never `PASS`, and the exception expires before P01 implementation.

## Mandatory contributor / AI start here

Read in this order before material changes:

1. `AGENTS.md`
2. Product Constitution and system/module architecture
3. glossary, naming, domain ownership and dependency matrix
4. identifier, money, time, locale and error standards
5. API and Event standards
6. Security Standard, Data Classification and Security Control Matrix
7. Testing, CI, Release and Quality Gate standards
8. Repository Structure, Local Development, Toolchain, Configuration and Developer Command standards
9. `docs/roadmap/MASTER_PLAN.md`, `STATUS.md`, `STATE.json`
10. AI Execution Policy, Change Control and Definition of Done
11. active temporary CI exception if present
12. relevant ADRs

`STATE.json` is the machine-readable execution source of truth.

## Core laws

- Kernel before business modules.
- One authoritative owner per write model/capability.
- Cross-module direct DB writes are forbidden.
- Cross-domain communication uses governed APIs/capabilities/events/workflows/read projections.
- Tenant scope, authorization, audit, observability and contract versioning are mandatory.
- Optional modules fail/degrade independently.
- AI acts only through governed authorized capabilities.
- Strict modular monolith first; extract services only when evidence justifies it.
- Architecture/roadmap changes require change control and ADR reconciliation.

## Frozen P00 foundation

### P00.03 — Primitives

UUIDv7 IDs, exact-decimal money, UTC/`timestamptz` plus IANA civil-time semantics, BCP 47 locale/RTL and stable safe error contracts.

### P00.04 — HTTP APIs

`/api/v{major}/{domain}/{resources}`, OpenAPI 3.2.0, `snake_case`, Problem Details, cursor pagination, explicit idempotency/concurrency and authorization-derived tenant context.

### P00.05 — Events

Producer-owned versioned facts, CloudEvents-compatible envelope, UUIDv7 event identity, at-least-once delivery, idempotent consumers, outbox/inbox reliability, bounded retry/DLQ and replay safety.

### P00.06 — Security

`PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`; RBAC + relationships + contextual policy + capabilities; tenant isolation; secrets/KMS; audit; privileged operations; integration/module/AI trust boundaries.

### P00.07 — Quality and release

Repository-owned CI-provider-independent verification. Gate classes G0–G8 cover governance through supply-chain/release. Evidence vocabulary is `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Releases prefer immutable build-once/promote artifacts and explicit migration/compatibility/rollback evidence.

### P00.08 — Repository and local development

Canonical governed monorepo categories:

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

Folders express ownership, not automatic microservices. Module private code/schema/migrations stay with the owner; cross-module private imports/tables are forbidden.

Default local infrastructure is containerized PostgreSQL, Redis-compatible cache, NATS/JetStream and S3-compatible object storage. Kubernetes is **not** required for the default local developer loop.

Toolchains/dependencies are repository-pinned. Configuration is explicit and separate from secrets. Development fixtures are synthetic/deterministic; production sensitive data is prohibited by default locally.

Future local command semantics include:

```text
omnexa dev bootstrap|up|down|status|reset
omnexa db status|migrate|fresh|seed
omnexa verify governance|format|lint|static|unit|contracts|integration|migrations|security|module-lifecycle|build|release|all
omnexa module create|validate|test|package
```

Linux is the canonical backend execution environment; macOS is supported where tooling permits; Windows backend development prefers WSL2. Native Windows is a separate certification target where POS/edge requirements need it.

Normative P00.08 documents are under `docs/development/`, with `docs/contracts/development/workspace.schema.json` and ADR-0008.

## Technology baseline

- Go — kernel/backend/domain services
- TypeScript + React — admin/web/builder/SDK
- Rust — justified edge/native/security-sensitive work
- Python — justified AI/data workloads
- PostgreSQL — primary OLTP
- Redis-compatible — cache/ephemeral coordination
- S3-compatible — files/media
- NATS/JetStream-class — events/messaging
- OpenTelemetry — observability

## Governance status

- Issue #3: hosted `main` branch protection still needs admin configuration.
- Issue #14: GitHub Actions quota/runner unavailable; ADR-0006 is P00-only.
- Issue #4: final licensing/IP/trademark strategy remains unresolved before external distribution/public launch.

None of these authorize early implementation.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00 through P27. Current status is **8/10 P00 packages done; P00.09 active**. P01+ remain planned.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**