# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are planned as domain families running on one shared platform foundation.

> **Current execution lock:** repository phase is **P00 — Product Constitution & Architecture Freeze**. Current package: **P00.10 — Foundation architecture freeze review**. Kernel and business-feature implementation remain unauthorized until the freeze review and P01 entry prerequisites are satisfied.

> **Temporary CI note:** GitHub Actions allowance is exhausted/disabled. ADR-0006 permits temporary P00 documentation/specification-only manual evidence. Hosted CI remains `BLOCKED`/`NOT RUN`, never `PASS`, and the exception expires at P00 exit/before P01 implementation.

## Mandatory contributor / AI start here

Read `AGENTS.md`, then the governance, architecture, security, quality, development, operations and roadmap documents referenced there. `docs/roadmap/STATE.json` is the machine-readable execution source of truth.

## Core laws

- Kernel before business modules.
- One authoritative owner per write model/capability.
- Cross-module direct DB writes/private imports are forbidden.
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
Producer-owned versioned facts, CloudEvents-compatible envelopes, UUIDv7 event identity, at-least-once delivery, idempotent consumers, outbox/inbox reliability, bounded retry/DLQ and replay safety.

### P00.06 — Security
`PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`; RBAC + relationships + contextual policy + capabilities; tenant isolation; secrets/KMS; audit; privileged operations; integration/module/AI trust boundaries.

### P00.07 — Quality and release
Repository-owned CI-provider-independent verification. G0–G8 cover governance through supply-chain/release. Evidence vocabulary: `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Releases prefer immutable build-once/promote artifacts with explicit migration/compatibility/rollback evidence.

### P00.08 — Repository and local development
Governed monorepo roots: `apps/`, `kernel/`, `modules/`, `platform/`, `shared/`, `infrastructure/`, `scripts/`, `docs/`, `generated/`. Local baseline uses containerized PostgreSQL, Redis-compatible cache, NATS/JetStream and S3-compatible object storage. Kubernetes is not required for ordinary local development. Toolchains/dependencies are repository-pinned; secrets are separate; production sensitive data is prohibited by default locally.

### P00.09 — Threat model and operational reliability
Foundation threat model covers tenant escape, authn/authz abuse, privilege escalation, injection/SSRF/webhooks, replay/idempotency, financial integrity, module/supply-chain compromise, CI/release credentials, POS/edge compromise, backup/search/vector leakage, AI prompt/tool abuse, insider misuse, noisy-neighbor/DDoS/provider outage, migration corruption, region failure, audit tampering, secrets exposure and misconfiguration.

Operational tiers: `TIER_0` integrity-critical, `TIER_1` core transactions, `TIER_2` supporting interactive, `TIER_3` optional/background. Initial mature-production availability objectives are 99.99%, 99.95%, 99.9% and 99.5% respectively.

Recovery targets:

```text
A: RPO <= 5m,  RTO <= 30m
B: RPO <= 15m, RTO <= 2h
C: RPO <= 24h, RTO <= 8h
D: rebuild-based derived state
```

These are architecture targets until recovery rehearsal proves them. Zero-tolerance conditions include cross-tenant disclosure, unauthorized privileged mutation, duplicate protected financial side effects, material financial/ledger integrity violation and lost acknowledged durable work. Incident model is SEV0–SEV3.

Normative P00.09 documents: `docs/security/THREAT_MODEL.md`, `docs/operations/`, `docs/contracts/operations/operational-targets.schema.json`, ADR-0009 and `scripts/validate_operations_spec.py`.

## Technology baseline

Go — core/backend; TypeScript + React — web/admin/builder/SDK; Rust — justified edge/native/security; Python — justified AI/data; PostgreSQL — OLTP; Redis-compatible — cache; S3-compatible — storage; NATS/JetStream-class — messaging; OpenTelemetry — observability.

## Governance status

- Issue #3: hosted `main` branch protection still requires admin configuration.
- Issue #14: GitHub Actions quota/runner unavailable; ADR-0006 is P00-only.
- Issue #4: final licensing/IP/trademark strategy remains unresolved before external distribution/public launch.

P00.10 must classify remaining blockers and P01 entry prerequisites explicitly; none authorizes early implementation.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current status: **9/10 P00 packages done; P00.10 active**. P01+ remain planned.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**