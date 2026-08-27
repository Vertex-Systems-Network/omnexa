# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current transition candidate:** **P03 — Module Runtime is ACTIVE at 3 / 11 with P03.01-P03.03 complete and P03.04 — Module Lifecycle State Machine as the sole active package after closure PR #99 merges.** `kernel_code_authorized=true` only for P03.04 after that merge; `business_feature_code_authorized=false`.

Protected main currently contains the completed P03.03 implementation at `774fab8b0350ffb2776517e3f1361f76bc2c68f9`. `docs/roadmap/STATE.json` remains the canonical machine-readable execution cursor; PR #99 does not grant P03.04 implementation authority before merge.

## Project progress

```text
P00  Product Constitution & Architecture Freeze  [██████████] 10/10  DONE
P01  Omnexa Kernel                               [██████████] 12/12  DONE
P02  Identity, Tenancy & Organization            [██████████] 10/10  DONE
      └─ Exit: SATISFIED
P03  Module Runtime                              [███░░░░░░░]  3/11  ACTIVE after closure
      ├─ P03.01 — Module Manifest Schema: DONE
      ├─ P03.02 — Registry & Deterministic Discovery: DONE
      ├─ P03.03 — Dependency Graph Resolver: DONE
      ├─ Current after closure: P03.04 — Module Lifecycle State Machine: ACTIVE
      └─ P03.05-P03.11: PLANNED / LOCKED
P04+ Future phases                               [░░░░░░░░░░]        PLANNED / LOCKED
```

## Mandatory contributor / AI start here

Read `AGENTS.md` first. Then read `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, `docs/roadmap/work-packages/P03.04.md` and `docs/ai/handoffs/P03.04.md` before material P03.04 work. If closure PR #99 is not merged, P03.04 runtime work is still locked.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/governance/P03_ENTRY_GATE.md`
- `docs/governance/P03_EXIT_GATE.md`
- `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P03.03.md`
- `docs/roadmap/work-packages/P03.04.md`
- `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`
- `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`
- `docs/roadmap/evidence/P03.03_COMPLETION_2026-08-28.md`
- `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`
- `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`
- `docs/quality/GO_CODE_QUALITY.md`

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

Canonical required CI uses GitHub-hosted `ubuntu-24.04` only. Local/self-hosted governance runners are prohibited. Repository-wide Go quality runs through `bash scripts/verify_go_quality.sh` using pinned tooling.

## Completed prerequisites and P03 evidence

P01.01-P01.12 and P02.01-P02.10 are complete with canonical executable evidence.

P03.01 — Module Manifest Schema is complete through PR #92, exact head `87da3302605c852ae5bf43d473aaa01a9e1aaa74`, run/job `33009396644 / 98311433013`, merge `4229e2a28442bf475afed143bab359a770d48053`.

P03.02 — Registry & Deterministic Discovery is complete through PR #94, exact head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05`, run/job `33022405704 / 98355747775`, merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`.

P03.03 — Dependency Graph Resolver is complete through PR #98:

- final exact head `4dcaca22911fbb81b1d25af316fef146c4a71ff3`
- canonical run/job `33112808869 / 98659824107` — PASS
- implementation merge `774fab8b0350ffb2776517e3f1361f76bc2c68f9`
- evidence `docs/roadmap/evidence/P03.03_COMPLETION_2026-08-28.md`
- retained verifier `scripts/verify_p03_03.sh`

Diagnostic #412/#414 failures and #413 cancellation remain historical non-acceptance evidence. All completed P01/P02/P03.01/P03.02/P03.03 regressions remain mandatory during P03.04.

## P03.04 boundary

P03.04 — Module Lifecycle State Machine becomes the sole active package only after closure PR #99 passes exact-final-head governance and merges to protected `main`. Its implementation must begin on a new branch from the exact post-merge main SHA.

P03.04 may implement only the canonical lifecycle contract: explicit state/transition validation, dependency and reverse-dependency lifecycle preconditions, safe disable/re-enable, suspend/archive/detach/purge distinctions, audited destructive purge, idempotency/retry/concurrency protection, deterministic partial-failure recovery and bounded upgrade lifecycle coordination.

Not authorized: P03.05-P03.11, P04+, business modules, package acquisition/marketplace runtime, strategic X-program runtime, AI/model/agent runtime, or weakening governance/security/regressions.

## Current implementation lock

- `kernel_code_authorized=true` only for P03.04 **after** PR #99 merges.
- `business_feature_code_authorized=false`.
- P03.05-P03.11 remain planned/locked.
- P04-P27 remain planned/locked.
- Strategic X programs remain non-authorizing until separately governed.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Exact next action

Finish the P03.03 closure / P03.04 activation reconciliation in draft PR #99. Require its exact final head to pass canonical GitHub-hosted governance, merge only if current and permitted, then re-read protected `main` and start P03.04 on a new separate implementation branch from that exact SHA. Do not implement P03.04 on the closure carrier and do not auto-advance to P03.05.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
