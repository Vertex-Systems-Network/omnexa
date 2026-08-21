# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** **Omnexa Foundation Architecture v1 is FROZEN.** P00 remains active only in **exit verification**. Current package: **P00.10 — Foundation architecture freeze review**.

> **Implementation lock:** the executable CI gate is now **SATISFIED** on the local Windows self-hosted runner. P01 kernel implementation remains **BLOCKED only by Issue #3 (`main` protection)**. Kernel and business-feature code remain unauthorized until that gate is cleared.

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

P00.01–P00.09 freeze governance/AI/change control, terminology/ownership/dependencies, UUIDv7/exact-money/time/locale/error primitives, HTTP/OpenAPI and event contracts, security/data classification/tenant isolation, G0–G8 quality/release semantics, monorepo/local-development rules, and the threat/SLO/recovery/incident baseline.

Technology baseline remains Go + TypeScript/React with selective Rust/Python; PostgreSQL, Redis-compatible cache, S3-compatible storage, NATS/JetStream-class messaging and OpenTelemetry.

## Executable CI — verified

The canonical governance workflow now uses:

```text
runs-on: [self-hosted, Windows, X64]
```

Verified PR #20 evidence:

- runner `LOCAL-WIN-01` / Windows X64;
- Git `2.55.0.windows.5`;
- Python `3.13.7`;
- workflow run `32522919774` SUCCESS;
- governance validator PASS;
- development-spec validator PASS;
- operations validator PASS;
- foundation-freeze validator PASS;
- PR #20 merged as `c2ab2cd679c295a8dec84b1879acb9a9e02ad67d`;
- Issue #14 closed/completed.

The old hosted-runner quota problem no longer blocks the governance lane. ADR-0006 remains historical evidence, not the active path while self-hosted CI is operational.

## P01 entry gate

P01 is **not authorized yet**.

### Issue #3 — remaining P01 entry blocker

Before executable P01 merges, protect `main` with PR-based integration, required `governance` check, blocked force-push/deletion, controlled bypass and conversation resolution.

### Issue #14 — satisfied

Executable CI has been restored and verified through the local Windows self-hosted runner.

### Issue #4 — external distribution blocker

Licensing/IP/trademark resolution is required before public/external distribution or self-hosted customer delivery, but does not block private internal P01 engineering after Issue #3 is cleared.

## Exact next transition

A narrow governance transition may close P00 after Issue #3 is verified. It must retire ADR-0006 from active use, mark P00.10/P00 done, activate P01, set `kernel_code_authorized = true`, keep business features locked, record branch-protection evidence and define the first P01 kernel work package.

Do **not** combine that transition with unrelated kernel feature code.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current canonical state is **Foundation v1 frozen; executable CI satisfied; P00.10 exit verification; P01 blocked only by Issue #3**.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
