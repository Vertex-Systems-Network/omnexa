# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** **Omnexa Foundation Architecture v1 is FROZEN.** P00 remains active only in **exit verification**. Current package: **P00.10 — Foundation architecture freeze review**.

> **Implementation lock:** executable CI is **SATISFIED on GitHub-hosted `ubuntu-24.04` only**. The repository is public, so the previous private-plan branch-protection blocker is gone, but live `main` still reports `protected:false`; therefore P01 remains blocked and kernel/business-feature code remain unauthorized.

## Mandatory contributor / AI start here

Read `AGENTS.md`, then the governance, architecture, security, quality, development, operations and roadmap documents referenced there. `docs/roadmap/STATE.json` is the machine-readable execution source of truth.

Freeze/entry-gate sources:

- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.01.md`
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

## Executable CI — GitHub-hosted only

The canonical required job remains `governance` and runs only on:

```yaml
runs-on: ubuntu-24.04
```

The workflow fails unless `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. It contains no `self-hosted` selector, no `LOCAL-WIN-*` routing and no local-runner fallback.

Current proof:

- workflow run `32537207455`: SUCCESS;
- job `96940269306`: SUCCESS;
- runner `GitHub Actions 1000006777`;
- `RUNNER_ENVIRONMENT=github-hosted`;
- Ubuntu 24.04.4 LTS / X64;
- image `ubuntu-24.04` version `20260816.277.1`;
- Git `2.55.0`;
- Python `3.12.3`;
- PowerShell `7.6.5`;
- governance/development/operations/freeze/P01-preparation/P01-package-spec validators PASS.

Local/self-hosted runner evidence remains historical provenance only and is not an active CI option.

## P01 entry gate

P01 is **not authorized yet**.

### Issue #3 — remaining blocker: `ACTIONABLE_UNPROTECTED`

The repository is now public, so the former private-repository plan limitation is no longer the active blocker. However GitHub's live branch API still reports `main.protected=false`.

The gate clears only after the required PR-only protection/ruleset is applied and verified with strict `governance`, conversation resolution, force-push/deletion blocking and administrator enforcement, or after an explicitly owner-approved superseding governance ADR.

### Issue #14 — satisfied

Executable CI is restored and canonical governance now runs only on GitHub-hosted infrastructure.

## P01 package preparation

All P01.01–P01.12 specifications are prepared under strict sequential one-active-package policy. P01.01 remains the only next package:

**Go workspace / build skeleton**

No `go.mod`, `go.work` or kernel entrypoint is authorized until the P00→P01 transition is merged.

## Public visibility / Issue #4

The repository is now public while the current repository `LICENSE` remains GPLv3 and trademark clearance remains unresolved. Issue #4 requires explicit owner/legal reconciliation for product launch, commercial licensing, contribution policy and trademark claims. The prepared worksheet remains `docs/governance/LICENSING_DECISION_BRIEF.md`.

## Exact next transition

A narrow governance-only transition may close P00 only after EG-02 is satisfied or deliberately superseded. It must follow `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`, retire ADR-0006 from active use, mark P00.10/P00 done, activate P01/P01.01, set `kernel_code_authorized = true`, keep business features locked, and preserve P01.02–P01.12 as planned.

The first kernel implementation PR comes **after** that state transition; it is not combined with it.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current canonical state is **Foundation v1 frozen; GitHub-hosted-only executable CI satisfied; repository public; P00.10 exit verification; P01.01–P01.12 prepared/planned; P01 blocked only by unverified `main` protection**.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
