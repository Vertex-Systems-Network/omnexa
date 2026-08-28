# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected main `55ec376146c4c43f24b079050a35f58eec13c479` contains the completed P03.08 implementation. This separate P03.08 closure / P03.09 activation carrier records P03.08 completion and makes P03.09 authoritative only if its exact final head passes canonical governance and merges.

Transition candidate state:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 8 / 11 done after this closure merges.
- P03.01-P03.08: DONE with canonical completion evidence.
- current work package after closure merge: P03.09 — Migration Ownership Registry.
- P03.10-P03.11: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P03.09 after closure merge.
- `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` and live protected-main state remain authoritative. No P03.09 implementation may start from this closure branch merely because continuity docs describe the intended transition.

## P03.08 canonical completion evidence

- implementation issue: #115 — completed
- implementation PR: #116 — merged
- final exact implementation head: `65dc38c6d60d1535c97a5dda59fb49490df59ec6`
- implementation merge / protected-main closure base: `55ec376146c4c43f24b079050a35f58eec13c479`
- canonical run/job: `33216021914 / 98999758150` — PASS
- runner: GitHub-hosted `ubuntu-24.04`, Linux/X64
- Go: repository-pinned `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.08_COMPLETION_2026-08-29.md`
- retained verifier: `scripts/verify_p03_08.sh`

The accepted exact-head run passed repository Go quality, all retained P01/P02/P03.01-P03.07 regressions and the dedicated P03.08 verifier. Historical failed/cancelled/stale candidates remain diagnostic evidence only and are never rewritten as acceptance evidence.

## Retained P03 prerequisites

P03.01 remains DONE through PR #92 / head `87da3302605c852ae5bf43d473aaa01a9e1aaa74` / run-job `33009396644 / 98311433013` / merge `4229e2a28442bf475afed143bab359a770d48053`.

P03.02 remains DONE through PR #94 / head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05` / run-job `33022405704 / 98355747775` / merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`.

P03.03 remains DONE through PR #98 / head `4dcaca22911fbb81b1d25af316fef146c4a71ff3` / run-job `33112808869 / 98659824107` / merge `774fab8b0350ffb2776517e3f1361f76bc2c68f9`.

P03.04 remains DONE through PR #100 / head `cddb42d4466e7f97a7547c4cf5ea0812c768ff0b` / run-job `33125377739 / 98702150001` / merge `13701e7647c1e084dfe4288d4b27b3ddd75e72c2`.

P03.05 remains DONE through issue #102 / PR #103 / head `c52b48be1a82eb27670f03bdd4e1be4df6eb9f54` / run-job `33132237120 / 98724184966` / merge `0c6b075c272aeac5a6e5f9d4210b1c5a30a040ce`.

P03.06 remains DONE through issue #106 / PR #107 / head `c895f44a1383d1c1d9c5fd23c95d7864810353c3` / run-job `33181421854 / 98883286556` / merge `13dbe8a393c20cabeb8aac60d073a6c66775efd3`.

P03.07 remains DONE through issue #110 / PR #111 / head `28e36b3ac3183f28ec500f1e70b1fefe02c0c325` / run-job `33195104185 / 98930123416` / merge `66f8b4cc630f6cd865e440a62478df365e042a31`.

P01/P02 and completed P03 regressions remain mandatory during later P03.09 implementation. Historical completion evidence is never rewritten by later transitions.

## Retained architecture/security baselines

ADR-0012 remains accepted. P03.03 preserves bounded manifest-version dependency semantics. P03.04 preserves fail-closed lifecycle, non-destructive disable/re-enable, dependency/reverse-dependency protections, guarded purge, deterministic recovery and unrelated-module isolation. P03.05 preserves `kernel.configuration` as authoritative for settings/feature flags and keeps flags non-authorizing. P03.06 preserves validated capability metadata as non-invoking and non-authorizing. P03.07 preserves `kernel.authorization` as the sole deny-by-default permission enforcement/policy authority. P03.08 preserves UI-contribution metadata as declarative, secret-free, non-executable and non-authorizing.

UI visibility does not grant permission. Disabled/unavailable modules cannot present contributions as operational. Optional-dependency absence degrades only affected contribution paths. Contribution metadata cannot introduce raw tenant/org authority, private execution shortcuts or cross-module database-write authority.

## P03.09 implementation contract

Owner: `kernel.modules`, integrated with the retained P01 database/migration foundation.

P03.09 implementation starts only after this P03.08 closure / P03.09 activation carrier passes exact-final-head governance, merges, and protected `main` plus canonical state are re-read. A new separate implementation branch must be created from that exact post-merge SHA.

Authorized P03.09 scope is limited to the canonical package specification and includes:

- deterministic migration identity/order/module/version/owner registry metadata;
- ownership checks against module/domain boundaries;
- fresh-install and supported-upgrade planning/validation;
- destructive/backfill classification metadata and lifecycle coordination;
- duplicate/order/conflict detection;
- representative migration fixtures required by later P03 exit proof;
- focused positive/adversarial tests;
- a dedicated `scripts/verify_p03_09.sh` verifier;
- canonical governance wiring after retained P03.08 verification.

P03.09 extends the existing P01 migration foundation rather than replacing it. A module may mutate only owned schema unless an explicitly approved platform migration says otherwise. Direct cross-module table mutation is forbidden. Migration inputs/paths cannot become arbitrary user-controlled file or SQL execution. Destructive/backfill work must be explicit. Retries/concurrency must not double-apply a migration. Tenant boundaries and data integrity remain mandatory.

## Explicitly unauthorized

- P03.09 runtime implementation on this closure carrier;
- replacement of the P01 migration engine;
- cross-domain schema ownership changes without accepted ADR/governance;
- arbitrary marketplace/plugin SQL execution;
- environment-specific manual SQL as accepted completion evidence;
- P04 event/outbox migration work before P04 activation;
- P03.10 module health runtime;
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

1. Finish atomic P03.08 closure / P03.09 activation reconciliation under GitHub issue #117.
2. Preserve historical failed/cancelled/stale runs as their actual evidence state.
3. Keep the carrier inside the governance/evidence/continuity boundary and include no P03.09 runtime code.
4. Require the exact final closure head to pass canonical GitHub-hosted governance with all retained regressions including `scripts/verify_p03_08.sh`.
5. Merge only if current with protected `main` and repository review/conversation gates permit it.
6. Re-read protected `main`, `STATE.json`, `STATUS.md`, package sequence and P03.09 handoff after merge and record the exact new SHA.
7. Identify a separate P03.09 implementation branch from that exact SHA as the next authorized action and stop the closure transaction without implementing P03.09.
8. Keep P03.10+ locked; P03.09 completion/state transition requires a later separate closure PR.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.09.md`, `docs/roadmap/evidence/P03.08_COMPLETION_2026-08-29.md`, accepted `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, completed `docs/ai/handoffs/P03.08.md` and activation-candidate `docs/ai/handoffs/P03.09.md`.
