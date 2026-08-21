# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** **Omnexa Foundation Architecture v1 is FROZEN.** P00 remains active only in **exit verification**. Current package: **P00.10 — Foundation architecture freeze review**.

> **Implementation lock:** P01 kernel implementation is **BLOCKED** until Issue #3 (`main` protection) and Issue #14 (working executable CI lane) are verified. Kernel and business-feature code remain unauthorized.

> **Temporary CI:** ADR-0006 permits P00 documentation/specification-only manual evidence while GitHub Actions is unavailable. Hosted CI is `BLOCKED`/`NOT RUN`, never `PASS`; the exception expires at P00 exit and cannot authorize P01 executable work.

## Mandatory contributor / AI start here

Read `AGENTS.md`, then the governance, architecture, security, quality, development, operations and roadmap documents referenced there. `docs/roadmap/STATE.json` is the machine-readable execution source of truth.

Freeze/entry-gate sources:

- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`

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

## Frozen foundation v1

P00.01–P00.09 freeze:

- governance, AI execution, change control and roadmap discipline;
- glossary, naming, domain ownership and dependency direction;
- UUIDv7, exact-money, time/locale/error primitives;
- stable HTTP/OpenAPI and event contracts;
- security, data classification, tenant isolation, authorization and audit;
- G0–G8 testing/CI/release semantics and exact evidence vocabulary;
- governed monorepo/local-development/toolchain/config model;
- threat model T01–T24, TIER_0–TIER_3 operational criticality, recovery classes A–D, SLO/RPO/RTO/error budgets, SEV0–SEV3 and reliability readiness.

Technology baseline remains Go + TypeScript/React with selective Rust/Python; PostgreSQL, Redis-compatible cache, S3-compatible storage, NATS/JetStream-class messaging and OpenTelemetry.

## P01 entry gate

P01 is **not authorized** yet.

### Issue #3 — P01 entry blocker

Before executable P01 merges: protect `main` with PR-based integration, blocked force-push/deletion, controlled bypass, conversation resolution and required verification checks when the CI lane exists.

### Issue #14 — P01 entry blocker

Before executable P01 merges: provide an approved CI/self-hosted/provider lane that actually runs the repository-owned verification commands in a clean reproducible environment. GitHub Actions is not architecturally mandatory; a compliant equivalent is acceptable.

### Issue #4 — external distribution blocker

Licensing/IP/trademark resolution is required before public/external distribution or self-hosted customer delivery, but does not block private internal P01 engineering after the P01 entry gate is cleared.

## Exact next transition

A narrow governance transition may close P00 only after Issue #3 and #14 are verified. It must expire ADR-0006, mark P00.10/P00 done, activate P01, set `kernel_code_authorized = true`, keep business features locked, record evidence and define the first P01 kernel work package.

Do **not** combine that transition with unrelated kernel feature code.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current canonical state is **Foundation v1 frozen; P00.10 exit verification; P01 blocked**.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**