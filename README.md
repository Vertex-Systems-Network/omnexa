# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P02 — Identity, Tenancy & Organization is DONE at 10 / 10 and its exit gate is SATISFIED. There is no active implementation package. P03 remains PLANNED / NOT ACTIVATED.** `kernel_code_authorized=false`; `business_feature_code_authorized=false`.

## Project progress

```text
P00  Product Constitution & Architecture Freeze  [██████████] 10/10  DONE
P01  Omnexa Kernel                               [██████████] 12/12  DONE
P02  Identity, Tenancy & Organization            [██████████] 10/10  DONE
      └─ Exit:    SATISFIED
      └─ Current: NONE — terminal P02 checkpoint
P03+ Future phases                               [░░░░░░░░░░]        PLANNED / LOCKED
```

The bars report only comparable package completion **inside each governed phase**. Omnexa does not publish a synthetic overall roadmap percentage across P00-P27 because later phases contain unequal scope and would make a single percentage misleading. The authoritative execution cursor remains `docs/roadmap/STATE.json`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_ENTRY_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P02.10.md`
- `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.02_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.03_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.05_COMPLETION_2026-08-24.md`
- `docs/roadmap/evidence/P02.06_COMPLETION_2026-08-24.md`
- `docs/roadmap/evidence/P02.07_COMPLETION_2026-08-24.md`
- `docs/roadmap/evidence/P02.08_COMPLETION_2026-08-25.md`
- `docs/roadmap/evidence/P02.09_COMPLETION_2026-08-26.md`
- `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`
- `docs/quality/GO_CODE_QUALITY.md`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`

For strategic architecture/roadmap work, also inspect the proposed AI-native expansion under `docs/adr/ADR-0011-ai-native-business-os-strategic-expansion.md`, `docs/roadmap/STRATEGIC_PROGRAMS.json`, `docs/roadmap/STRATEGIC_CROSS_CUTTING_PROGRAMS.md`, `docs/roadmap/STRATEGIC_ACCEPTANCE_GATES.md`, `docs/roadmap/AI_NATIVE_STRATEGIC_AUDIT_2026.md`, `docs/architecture/AI_NATIVE_BUSINESS_OS_ARCHITECTURE.md` and `docs/architecture/PRODUCT_FEDERATION_AND_APP_MESH.md`. Proposed strategic files do **not** activate future implementation by themselves.

## Core laws

- Kernel before business modules.
- One authoritative owner per write model/capability.
- Cross-module direct DB writes/private implementation imports are forbidden.
- Cross-domain communication uses governed APIs/capabilities/events/workflows/read projections.
- Tenant scope, authorization, audit, observability and contract versioning are mandatory.
- Optional modules fail/degrade independently.
- AI acts only through governed authorized capabilities; no unrestricted raw DB/object-store/business-state authority.
- Strict modular monolith first; service extraction requires evidence and ADR.
- Architecture/roadmap changes require change control and reconciliation.

## Protected GitHub integration and executable CI

Issue #3 is closed and `main` is protected with PR-only integration, strict required `governance`, blocked direct/force updates, failed-check merge rejection, required conversation resolution and up-to-date branch enforcement.

Canonical required CI uses GitHub-hosted `ubuntu-24.04` only and fails closed unless the runner is GitHub-hosted Linux/X64. Local/self-hosted governance runners are prohibited. Permanent repository-wide Go quality runs through `bash scripts/verify_go_quality.sh` using pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

## P01 completion retained

P01.01-P01.12 are complete with canonical executable evidence. Final P01.12 evidence: PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. P01 regressions remain mandatory.

## P02 completion

P02.01-P02.10 are complete with retained package-specific evidence and mandatory regression verifiers. The terminal package, P02.10, completed through implementation PR #88, final exact head `975e4925060a035780ca13b68c5437634ed0f4ea`, canonical GitHub-hosted run/job `32904678957 / 97986011269`, and protected-main implementation merge `88799aa41da8ce8c22540146d157d488565e2ce9`.

P02.10 canonical evidence passed repository Go quality, P01.01-P01.12, real `omnexa db migrate`, real `omnexa verify all`, P02.01-P02.09, aggregate PostgreSQL exit proof and P02.10 G0-G8. It proves classification-safe audit integration for accepted P02 security lifecycle hooks, required-audit mutation behavior, same-tenant success/cross-tenant denial, migration replay/idempotency, no-secret audit behavior and preservation of identity/tenancy/organization/authz/service-account/session/settings boundaries.

Diagnostic run `32903969206 / 97983773781` remains historical **FAIL** evidence for a corrected undefined helper reference and is not acceptance evidence. Immutable completion evidence is `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`.

`docs/governance/P02_EXIT_GATE.md` is **SATISFIED**. P02 completion does not implicitly activate P03.

## Current implementation lock

- `kernel_code_authorized=false`.
- `business_feature_code_authorized=false`.
- P03 remains planned and inactive.
- The next governed work is P03 specification/readiness preparation and, only after separate acceptance, an explicit activation transition.
- All proposed `X..` strategic programs remain planning-only until accepted, dependency-ready and activated through canonical work-package/state governance.

## Strategic AI-native direction

The proposed strategic layer does not replace P00-P27. It adds cross-cutting intelligence/control systems so Omnexa can evolve beyond feature aggregation into:

**deterministic domain truth + universal workflow + governed data/Business Graph + Process Graph + System Graph + AI Control Tower + model governance + simulation + continuous controls + measurable outcomes + constrained autonomy.**

Standalone first-party products should join through `XPF-200 Product Federation & App Mesh` as native, embedded, federated or edge products according to architecture evidence rather than being force-merged into one codebase/database.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current checkpoint: **P00 done 10 / 10; P01 done 12 / 12; P02 done 10 / 10 with exit SATISFIED; no active package; P03+ and business implementation locked.**

`docs/roadmap/STRATEGIC_PROGRAMS.json` records proposed cross-cutting X-programs. It does not create a second execution cursor or override `STATE.json`.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
