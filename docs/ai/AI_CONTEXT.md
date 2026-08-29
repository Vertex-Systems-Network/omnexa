# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected main `ea402964c45a630fd6723e0e4a6754555a6a4994` contains the completed P03.09 implementation. This separate P03.09 closure / P03.10 activation carrier records P03.09 completion and makes P03.10 authoritative only if its exact final head passes canonical governance and merges.

Transition candidate state:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 9 / 11 done after this closure merges.
- P03.01-P03.09: DONE with canonical completion evidence.
- current work package after closure merge: P03.10 — Module Health Reporting.
- P03.11: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P03.10 after closure merge.
- `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` and live protected-main state remain authoritative. No P03.10 implementation may start from this closure branch merely because continuity docs describe the intended transition.

## P03.09 canonical completion evidence

- implementation issue: #120 — completed
- draft implementation carrier: #121 — closed unmerged
- promotion implementation PR: #122 — merged
- final exact implementation head: `8c4da1c1c9e11dfe2f1fa4b81b730140a9f24d56`
- implementation merge / protected-main closure base: `ea402964c45a630fd6723e0e4a6754555a6a4994`
- promotion-specific canonical run/job: `33223035182 / 99020954655` (#495) — PASS
- runner: GitHub-hosted `ubuntu-24.04`, Linux/X64
- Go: repository-pinned `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.09_COMPLETION_2026-08-29.md`
- retained verifier: `scripts/verify_p03_09.sh`

Governance #493 / `33222404123` remains diagnostic FAIL evidence for formatter alignment only. Governance #494 / `33222631307` is successful exact-head draft-carrier evidence; #495 is promotion-specific acceptance authority. Historical evidence states are never rewritten.

## Retained P03 prerequisites

P03.01 remains DONE through PR #92 / head `87da3302605c852ae5bf43d473aaa01a9e1aaa74` / run-job `33009396644 / 98311433013` / merge `4229e2a28442bf475afed143bab359a770d48053`.

P03.02 remains DONE through PR #94 / head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05` / run-job `33022405704 / 98355747775` / merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`.

P03.03 remains DONE through PR #98 / head `4dcaca22911fbb81b1d25af316fef146c4a71ff3` / run-job `33112808869 / 98659824107` / merge `774fab8b0350ffb2776517e3f1361f76bc2c68f9`.

P03.04 remains DONE through PR #100 / head `cddb42d4466e7f97a7547c4cf5ea0812c768ff0b` / run-job `33125377739 / 98702150001` / merge `13701e7647c1e084dfe4288d4b27b3ddd75e72c2`.

P03.05 remains DONE through issue #102 / PR #103 / head `c52b48be1a82eb27670f03bdd4e1be4df6eb9f54` / run-job `33132237120 / 98724184966` / merge `0c6b075c272aeac5a6e5f9d4210b1c5a30a040ce`.

P03.06 remains DONE through issue #106 / PR #107 / head `c895f44a1383d1c1d9c5fd23c95d7864810353c3` / run-job `33181421854 / 98883286556` / merge `13dbe8a393c20cabeb8aac60d073a6c66775efd3`.

P03.07 remains DONE through issue #110 / PR #111 / head `28e36b3ac3183f28ec500f1e70b1fefe02c0c325` / run-job `33195104185 / 98930123416` / merge `66f8b4cc630f6cd865e440a62478df365e042a31`.

P03.08 remains DONE through issue #115 / PR #116 / head `65dc38c6d60d1535c97a5dda59fb49490df59ec6` / run-job `33216021914 / 98999758150` / merge `55ec376146c4c43f24b079050a35f58eec13c479`.

P01/P02 and completed P03 regressions remain mandatory during later P03.10 implementation. Historical completion evidence is never rewritten by later transitions.

## Retained architecture/security baselines

ADR-0012 remains accepted. P03.03 preserves bounded manifest-version dependency semantics. P03.04 preserves fail-closed lifecycle, non-destructive disable/re-enable, dependency/reverse-dependency protections, guarded purge, deterministic recovery and unrelated-module isolation. P03.05 preserves `kernel.configuration` as authoritative for settings/feature flags and keeps flags non-authorizing. P03.06 preserves validated capability metadata as non-invoking and non-authorizing. P03.07 preserves `kernel.authorization` as the sole deny-by-default permission enforcement/policy authority. P03.08 preserves UI-contribution metadata as declarative, secret-free, non-executable and non-authorizing. P03.09 preserves P01 as sole migration execution/checksum/locking authority while adding fail-closed owner/version planning metadata.

P03.09 migration metadata cannot grant cross-owner database authority, embed raw SQL or arbitrary file execution, hide destructive/backfill intent, carry secrets/raw tenant authority or bypass the retained P01 migration ledger/locking path.

## P03.10 implementation contract

Owner: `kernel.modules`, integrated with the retained P01 health/readiness/observability foundation.

P03.10 implementation starts only after this P03.09 closure / P03.10 activation carrier passes exact-final-head governance, merges, and protected `main` plus canonical state are re-read. A new separate implementation branch must be created from that exact post-merge SHA.

Authorized P03.10 scope is limited to the canonical package specification and includes:

- stable module ID/version/state health identity;
- required versus optional dependency health/degradation;
- migration compatibility/pending/failure status needed for safe operation;
- capability/permission/UI registration availability summaries where applicable;
- explicit healthy/degraded/unavailable/failed diagnostics with bounded reason categories;
- deterministic reference-module health fixtures;
- focused positive/adversarial tests;
- a dedicated `scripts/verify_p03_10.sh` verifier;
- canonical governance wiring after retained P03.09 verification.

P03.10 reuses P01 health/readiness/observability rather than replacing it. Health remains diagnostic evidence, never authorization or lifecycle command authority. Required dependency failure must fail closed while optional absence degrades selectively. Migration inconsistency cannot be presented as healthy. Health output must remain classification-safe and secret-free, and one module failure must not corrupt unrelated healthy-module reporting where isolation is feasible.

## Explicitly unauthorized

- P03.10 runtime implementation on this closure carrier;
- System Graph storage/collector runtime;
- performance-budget analytics;
- business KPI monitoring;
- secrets/debug dumps or unrestricted topology disclosure;
- lifecycle command/control authority through health reporting;
- authorization or permission grants through health state;
- P03.11 trust hooks/phase-exit runtime;
- P04+ runtime;
- business modules/features;
- package acquisition/marketplace runtime;
- generic remote RPC/service mesh;
- workflow orchestration;
- strategic X-program runtime;
- AI/model/agent runtime;
- weakening retained regressions, branch protection, security or governance gates.

## Exact next action

1. Finish atomic P03.09 closure / P03.10 activation reconciliation under GitHub issue #123 / Linear ABD-187.
2. Preserve historical failed/cancelled/stale runs as their actual evidence state.
3. Keep the carrier inside the governance/evidence/continuity boundary and include no P03.10 runtime code.
4. Require the exact final closure head to pass canonical GitHub-hosted governance with all retained regressions including `scripts/verify_p03_09.sh`.
5. Merge only if current with protected `main` and repository review/conversation gates permit it.
6. Re-read protected `main`, `STATE.json`, `STATUS.md`, package sequence and P03.10 handoff after merge and record the exact new SHA.
7. Identify a separate P03.10 implementation branch from that exact SHA as the next authorized action and stop the closure transaction without implementing P03.10.
8. Keep P03.11 locked; P03.10 completion/state transition requires a later separate closure PR.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.10.md`, `docs/roadmap/evidence/P03.09_COMPLETION_2026-08-29.md`, accepted `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, completed `docs/ai/handoffs/P03.09.md` and activation-candidate `docs/ai/handoffs/P03.10.md`.
