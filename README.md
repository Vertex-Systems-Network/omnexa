# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** **Omnexa Foundation Architecture v1 is FROZEN.** P00 remains active only in **exit verification**. Current package: **P00.10 — Foundation architecture freeze review**.

> **Implementation lock:** the executable CI gate is **SATISFIED on any available Windows/X64 self-hosted runner**. P01 kernel implementation remains **BLOCKED by Issue #3 / GitHub plan-limited `main` protection**. P01.01 is prepared/specification-only; kernel and business-feature code remain unauthorized.

## Mandatory contributor / AI start here

Read `AGENTS.md`, then the governance, architecture, security, quality, development, operations and roadmap documents referenced there. `docs/roadmap/STATE.json` is the machine-readable execution source of truth.

Freeze/entry-gate sources:

- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`
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

## Executable CI — runner-name agnostic

The canonical governance workflow exposes one required job named `governance` on:

```yaml
runs-on: [self-hosted, Windows, X64]
```

GitHub schedules that job on whichever eligible Windows/X64 self-hosted runner is available. There is no runner-name pinning, discovery fanout or target-runner artifact aggregation. The full validator set runs directly in the required job and fails closed on any error.

Current routing proof:

- workflow run `32535324900`: SUCCESS;
- job `96935023669`: SUCCESS;
- actual runner `LOCAL-WIN-02` / Windows X64;
- machine `ABDUL-HANAN`;
- Git `2.55.0.windows.5`;
- Python `3.13.7`;
- PowerShell branch-protection tooling parse PASS;
- governance validator PASS;
- development-spec validator PASS;
- operations validator PASS;
- foundation-freeze validator PASS;
- P01-preparation validator PASS.

The earlier LOCAL-WIN-4 evidence from PR #23/run `32528329184` remains historical provenance, not a scheduling requirement.

While P01 is blocked, CI runs `scripts/validate_p01_preparation.py`, which keeps P01.01 prepared but fails if known executable kernel paths appear before the authorization transition.

## P01 entry gate

P01 is **not authorized yet**.

### Issue #3 — remaining P01 entry blocker: `BLOCKED_BY_PLAN`

The branch-protection tooling is merged and validated, but GitHub returned HTTP 403 for private branch protection on the current organization plan. `main` remains unprotected.

The gate clears only when hosted protection becomes available and is verified, or an explicitly owner-approved superseding governance ADR replaces EG-02 with compensating controls. Making the repository public is not an automatic workaround.

### Issue #14 — satisfied

Executable CI is restored and runner-name agnostic inside the approved Windows/X64 self-hosted pool.

## P01.01 — prepared, not active

The first kernel work package is fully specified at `docs/roadmap/work-packages/P01.01.md`:

**Go workspace / build skeleton**

Prepared scope includes repository-owned Go toolchain/workspace structure, deterministic build/test entrypoints, minimal process/build metadata and G0/G1/G2/G7 evidence. It explicitly excludes configuration, PostgreSQL, cache, object storage, telemetry, health, jobs, feature flags, audit, identity/tenancy, module runtime and business domains.

No `go.mod`, `go.work` or kernel entrypoint is authorized until the P00→P01 transition is merged.

## Issue #4 — external distribution blocker

Licensing/IP/trademark resolution is required before public/external distribution, self-hosted customer delivery, public launch or external contributor intake, but does not block private internal P01 engineering after the P01 entry gate is cleared.

Prepared owner/legal worksheet: `docs/governance/LICENSING_DECISION_BRIEF.md`. It does not change the current `LICENSE` or establish trademark clearance.

## Exact next transition

A narrow governance-only transition may close P00 only after EG-02 is satisfied or deliberately superseded. It must follow `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`, retire ADR-0006 from active use, mark P00.10/P00 done, activate P01/P01.01, set `kernel_code_authorized = true`, keep business features locked, and preserve P01.02–P01.12 as planned.

The first kernel implementation PR comes **after** that state transition; it is not combined with it.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current canonical state is **Foundation v1 frozen; runner-name-agnostic executable CI satisfied; P00.10 exit verification; P01.01 prepared/planned; P01 blocked by plan-limited Issue #3**.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
