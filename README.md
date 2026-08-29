# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current transition candidate:** **P03 — Module Runtime is ACTIVE at 9 / 11 with P03.01-P03.09 complete and P03.10 — Module Health Reporting as the sole active package after this closure merges.** `kernel_code_authorized=true` only for P03.10 after that merge; `business_feature_code_authorized=false`.

Protected main currently contains the completed P03.09 implementation at `ea402964c45a630fd6723e0e4a6754555a6a4994`. `docs/roadmap/STATE.json` remains the canonical machine-readable execution cursor; this closure carrier does not grant P03.10 implementation authority before merge and protected-main readback.

## Project progress

```text
P00  Product Constitution & Architecture Freeze  [██████████] 10/10  DONE
P01  Omnexa Kernel                               [██████████] 12/12  DONE
P02  Identity, Tenancy & Organization            [██████████] 10/10  DONE
      └─ Exit: SATISFIED
P03  Module Runtime                              [█████████░]  9/11  ACTIVE after closure
      ├─ P03.01 — Module Manifest Schema: DONE
      ├─ P03.02 — Registry & Deterministic Discovery: DONE
      ├─ P03.03 — Dependency Graph Resolver: DONE
      ├─ P03.04 — Module Lifecycle State Machine: DONE
      ├─ P03.05 — Module Settings & Feature Flags: DONE
      ├─ P03.06 — Capability Registry: DONE
      ├─ P03.07 — Permission Registration: DONE
      ├─ P03.08 — UI Contribution Registry Contract: DONE
      ├─ P03.09 — Migration Ownership Registry: DONE
      ├─ Current after closure: P03.10 — Module Health Reporting: ACTIVE
      └─ P03.11: PLANNED / LOCKED
P04+ Future phases                               [░░░░░░░░░░]        PLANNED / LOCKED
```

## Mandatory contributor / AI start here

Read `AGENTS.md` first. Then read `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, `docs/roadmap/work-packages/P03.10.md`, completed `docs/ai/handoffs/P03.09.md` and activation-candidate `docs/ai/handoffs/P03.10.md` before material P03.10 work. Until this closure merges and protected main is re-read, P03.10 runtime work remains locked.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/governance/P03_ENTRY_GATE.md`
- `docs/governance/P03_EXIT_GATE.md`
- `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P03.09.md`
- `docs/roadmap/work-packages/P03.10.md`
- `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`
- `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`
- `docs/roadmap/evidence/P03.03_COMPLETION_2026-08-28.md`
- `docs/roadmap/evidence/P03.04_COMPLETION_2026-08-28.md`
- `docs/roadmap/evidence/P03.05_COMPLETION_2026-08-28.md`
- `docs/roadmap/evidence/P03.06_COMPLETION_2026-08-28.md`
- `docs/roadmap/evidence/P03.07_COMPLETION_2026-08-28.md`
- `docs/roadmap/evidence/P03.08_COMPLETION_2026-08-29.md`
- `docs/roadmap/evidence/P03.09_COMPLETION_2026-08-29.md`
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

P03.03 — Dependency Graph Resolver is complete through PR #98, exact head `4dcaca22911fbb81b1d25af316fef146c4a71ff3`, run/job `33112808869 / 98659824107`, merge `774fab8b0350ffb2776517e3f1361f76bc2c68f9`.

P03.04 — Module Lifecycle State Machine is complete through PR #100, exact head `cddb42d4466e7f97a7547c4cf5ea0812c768ff0b`, run/job `33125377739 / 98702150001`, merge `13701e7647c1e084dfe4288d4b27b3ddd75e72c2`, evidence `docs/roadmap/evidence/P03.04_COMPLETION_2026-08-28.md`, retained verifier `scripts/verify_p03_04.sh`.

P03.05 — Module Settings & Feature Flags is complete through PR #103, exact head `c52b48be1a82eb27670f03bdd4e1be4df6eb9f54`, run/job `33132237120 / 98724184966`, merge `0c6b075c272aeac5a6e5f9d4210b1c5a30a040ce`, evidence `docs/roadmap/evidence/P03.05_COMPLETION_2026-08-28.md`, retained verifier `scripts/verify_p03_05.sh`.

P03.06 — Capability Registry is complete through PR #107, exact head `c895f44a1383d1c1d9c5fd23c95d7864810353c3`, run/job `33181421854 / 98883286556`, merge `13dbe8a393c20cabeb8aac60d073a6c66775efd3`, evidence `docs/roadmap/evidence/P03.06_COMPLETION_2026-08-28.md`, retained verifier `scripts/verify_p03_06.sh`.

P03.07 — Permission Registration is complete through PR #111, exact head `28e36b3ac3183f28ec500f1e70b1fefe02c0c325`, run/job `33195104185 / 98930123416`, merge `66f8b4cc630f6cd865e440a62478df365e042a31`, evidence `docs/roadmap/evidence/P03.07_COMPLETION_2026-08-28.md`, retained verifier `scripts/verify_p03_07.sh`.

P03.08 — UI Contribution Registry Contract is complete through PR #116, exact head `65dc38c6d60d1535c97a5dda59fb49490df59ec6`, run/job `33216021914 / 98999758150`, merge `55ec376146c4c43f24b079050a35f58eec13c479`, evidence `docs/roadmap/evidence/P03.08_COMPLETION_2026-08-29.md`, retained verifier `scripts/verify_p03_08.sh`.

P03.09 — Migration Ownership Registry is complete through promotion PR #122:

- implementation issue #120 — completed
- draft carrier #121 — closed unmerged after exact-head Governance #494 / `33222631307` SUCCESS
- final exact head `8c4da1c1c9e11dfe2f1fa4b81b730140a9f24d56`
- promotion-specific canonical run/job `33223035182 / 99020954655` (#495) — PASS
- implementation merge `ea402964c45a630fd6723e0e4a6754555a6a4994`
- evidence `docs/roadmap/evidence/P03.09_COMPLETION_2026-08-29.md`
- retained verifier `scripts/verify_p03_09.sh`

Governance #493 / `33222404123` remains diagnostic FAIL evidence for formatter alignment only. Historical failed/cancelled/stale candidates remain their original evidence states. All completed P01/P02/P03.01-P03.09 regressions remain mandatory during later P03.10 implementation.

## P03.09 retained boundary

Migration ownership registration remains execution-free metadata/planning around the retained P01 migration foundation. P01 remains the sole SQL execution, checksum-ledger, owner-lock and transactional retry authority. Migration identity is owner/module/version/declaration bound; cross-owner targets and owner/version order conflicts fail closed; destructive/backfill intent is explicit and requires bounded strategy/recovery metadata; raw SQL, arbitrary filesystem path, execution callback, secret and raw tenant/org authority are not accepted registry surfaces.

## P03.10 boundary

P03.10 — Module Health Reporting becomes the sole active package only after this closure passes exact-final-head governance and merges to protected `main`. Its implementation must begin on a new branch from the exact post-merge main SHA.

P03.10 may implement only classification-safe module health/state diagnostics: stable module identity/version/state, required/optional dependency health, migration compatibility/pending/failure status, capability/permission/UI registration availability summaries, bounded healthy/degraded/unavailable/failed reason categories, deterministic fixtures, adversarial tests and a dedicated P03.10 verifier around the retained P01 health/readiness/observability foundation.

Health is diagnostic evidence only. It cannot grant authorization, mutate lifecycle state, expose secrets/raw tenant authority/internal stack traces, replace P01 health/readiness, or make migration inconsistency appear healthy. Required dependency failure must fail closed; optional dependency absence degrades selectively; one module failure must not corrupt unrelated healthy-module reporting where isolation is feasible.

Not authorized: System Graph storage/collector runtime, performance-budget analytics, business KPI monitoring, secrets/debug dumps or unrestricted topology disclosure, lifecycle command authority, authorization grants, P03.11, generic RPC/service mesh, workflow orchestration, P04+, business modules, package acquisition/marketplace runtime, strategic X-program runtime, AI/model/agent runtime, or weakening governance/security/regressions.

## Current implementation lock

- `kernel_code_authorized=true` only for P03.10 **after** this closure merges and protected main is re-read.
- `business_feature_code_authorized=false`.
- P03.11 remains planned/locked.
- P04-P27 remain planned/locked.
- Strategic X programs remain non-authorizing until separately governed.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Exact next action

Finish the P03.09 closure / P03.10 activation reconciliation under issue #123. Require its exact final head to pass canonical GitHub-hosted governance, merge only if current and permitted, then re-read protected `main`, canonical state/status/package sequence and P03.10 handoff. Record the exact post-closure main SHA and stop the closure transaction; P03.10 implementation belongs to a new separate branch and must not be implemented on this carrier.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**