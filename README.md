# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state after this closure transition merges:** **P03 — Module Runtime is ACTIVE at 2 / 11 with P03.01-P03.02 complete and P03.03 — Dependency Graph Resolver as the sole active package.** `kernel_code_authorized=true` only for P03.03 after activation; `business_feature_code_authorized=false`. Until then, protected `main` remains authoritative at P03.02.

## Project progress

```text
P00  Product Constitution & Architecture Freeze  [██████████] 10/10  DONE
P01  Omnexa Kernel                               [██████████] 12/12  DONE
P02  Identity, Tenancy & Organization            [██████████] 10/10  DONE
      └─ Exit: SATISFIED
P03  Module Runtime                              [██░░░░░░░░]  2/11  ACTIVE AFTER CLOSURE MERGE
      ├─ P03.01 — Module Manifest Schema: DONE
      ├─ P03.02 — Registry & Deterministic Discovery: DONE
      └─ Current after closure: P03.03 — Dependency Graph Resolver
      └─ P03.04-P03.11: PLANNED / LOCKED
P04+ Future phases                               [░░░░░░░░░░]        PLANNED / LOCKED
```

The bars report only comparable package completion **inside each governed phase**. Omnexa does not publish a synthetic overall roadmap percentage across P00-P27 because later phases contain unequal scope and would make a single percentage misleading. The authoritative execution cursor remains `docs/roadmap/STATE.json`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and, after canonical activation is verified, `docs/ai/handoffs/P03.03.md`.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/governance/P03_ENTRY_GATE.md`
- `docs/governance/P03_EXIT_GATE.md`
- `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P03.01.md`
- `docs/roadmap/work-packages/P03.02.md`
- `docs/roadmap/work-packages/P03.03.md`
- `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`
- `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`
- `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`
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

## Completed prerequisites

P01.01-P01.12 are complete with canonical executable evidence. Final P01.12 evidence: PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

P02.01-P02.10 are complete. Terminal P02.10 evidence: PR #88, exact head `975e4925060a035780ca13b68c5437634ed0f4ea`, canonical run/job `32904678957 / 97986011269`, merge `88799aa41da8ce8c22540146d157d488565e2ce9`. Diagnostic run `32903969206 / 97983773781` remains historical **FAIL** evidence and is not acceptance evidence.

P03.01 — Module Manifest Schema is complete. Canonical completion identity: implementation PR #92, exact head `87da3302605c852ae5bf43d473aaa01a9e1aaa74`, run/job `33009396644 / 98311433013`, implementation merge `4229e2a28442bf475afed143bab359a770d48053`, evidence `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`.

P03.02 — Registry & Deterministic Discovery is complete with implementation PR #94, final exact head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05`, canonical run/job `33022405704 / 98355747775`, implementation merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`, and evidence `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`. Its dedicated `scripts/verify_p03_02.sh` remains a mandatory regression.

An earlier P03.02 implementation candidate exposed an `errorlint` wrapped-error inspection defect and was corrected with `errors.As` without gate, acceptance or scope weakening. Only the final exact successful head/run above is acceptance evidence.

All completed P01/P02/P03.01/P03.02 regression verifiers remain mandatory during P03.03.

## P03.03 activation

The P03.02 closure transition advances exactly one package and contains no P03.03 implementation. After the closure PR passes its own exact-head canonical governance and merges to protected `main`, the active package is **P03.03 — Dependency Graph Resolver**, owned by `kernel.modules`.

P03.03 may implement only the dependency-resolution contract in `docs/roadmap/work-packages/P03.03.md`: version-aware required/optional dependency resolution, platform dependency validation, deterministic topological order, cycle/incompatible-version rejection, undeclared/forbidden/private dependency detection hooks and selective optional-dependency degradation metadata.

Required or incompatible dependencies and circular required graphs fail closed. Missing optional dependencies must not globally fail unrelated modules. Resolver output creates no permission/capability/tenant/database/private-schema authority, and identical validated registry/manifests must resolve deterministically.

Not authorized by P03.03: lifecycle mutations/persistence, package installation/download, P03.04-P03.11 implementation, P04 event orchestration, full System Graph runtime, trust/advisory scanning, direct cross-module private imports/writes, business modules, or AI/model/agent runtime.

The P03 AI-native compatibility mapping for XQ-100, XSG-100, XTRUST-100, XPF-200 and XPERF-100 remains planning-only; it does not activate those strategic runtimes.

## Current implementation lock

- `kernel_code_authorized=true` **only for P03.03 after the closure transition merges**.
- Until that merge, protected `main` remains authoritative at P03.02.
- `business_feature_code_authorized=false`.
- P03.04-P03.11 remain planned/locked.
- P04-P27 remain planned/locked.
- Strategic X programs remain non-authorizing until separately accepted/dependency-ready/activated.

## Strategic AI-native direction

The proposed strategic layer does not replace P00-P27. It adds cross-cutting intelligence/control systems so Omnexa can evolve beyond feature aggregation into:

**deterministic domain truth + universal workflow + governed data/Business Graph + Process Graph + System Graph + AI Control Tower + model governance + simulation + continuous controls + measurable outcomes + constrained autonomy.**

Standalone first-party products should join through `XPF-200 Product Federation & App Mesh` as native, embedded, federated or edge products according to architecture evidence rather than being force-merged into one codebase/database.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Closure-candidate checkpoint: **P00 done; P01 done 12 / 12; P02 done 10 / 10 with exit SATISFIED; P03 active 2 / 11 after closure; P03.01-P03.02 done; P03.03 activation target.** Protected `main` becomes authoritative for P03.03 only after this closure transition passes governance and merges.

`docs/roadmap/STRATEGIC_PROGRAMS.json` records proposed cross-cutting X-programs. It does not create a second execution cursor or override `STATE.json`.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
