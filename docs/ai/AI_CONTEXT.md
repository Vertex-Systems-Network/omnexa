# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected main `66f8b4cc630f6cd865e440a62478df365e042a31` contains the completed P03.07 implementation. This separate P03.07 closure / P03.08 activation carrier records P03.07 completion and makes P03.08 authoritative only if its exact final head passes canonical governance and merges.

Transition candidate state:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 7 / 11 done after this closure merges.
- P03.01-P03.07: DONE with canonical completion evidence.
- current work package after closure merge: P03.08 — UI Contribution Registry Contract.
- P03.09-P03.11: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P03.08 after closure merge.
- `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` and live protected-main state remain authoritative. No P03.08 implementation may start from this closure branch merely because continuity docs describe the intended transition.

## P03.07 canonical completion evidence

- implementation issue: #110 — completed
- implementation PR: #111 — merged
- final exact implementation head: `28e36b3ac3183f28ec500f1e70b1fefe02c0c325`
- implementation merge / protected-main closure base: `66f8b4cc630f6cd865e440a62478df365e042a31`
- canonical run/job: `33195104185 / 98930123416` — PASS
- runner: GitHub-hosted `ubuntu-24.04`, Linux/X64
- Go: repository-pinned `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.07_COMPLETION_2026-08-28.md`
- retained verifier: `scripts/verify_p03_07.sh`

Governance #474 / `33192567020 / 98921494281` and Governance #476 / `33194438411 / 98927852853` remain diagnostic FAIL evidence only. The accepted exact-head run passed repository Go quality, all retained P01/P02/P03.01-P03.06 regressions and the dedicated P03.07 verifier.

## Retained P03 prerequisites

P03.01 remains DONE through PR #92 / head `87da3302605c852ae5bf43d473aaa01a9e1aaa74` / run-job `33009396644 / 98311433013` / merge `4229e2a28442bf475afed143bab359a770d48053`.

P03.02 remains DONE through PR #94 / head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05` / run-job `33022405704 / 98355747775` / merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`.

P03.03 remains DONE through PR #98 / head `4dcaca22911fbb81b1d25af316fef146c4a71ff3` / run-job `33112808869 / 98659824107` / merge `774fab8b0350ffb2776517e3f1361f76bc2c68f9`.

P03.04 remains DONE through PR #100 / head `cddb42d4466e7f97a7547c4cf5ea0812c768ff0b` / run-job `33125377739 / 98702150001` / merge `13701e7647c1e084dfe4288d4b27b3ddd75e72c2`.

P03.05 remains DONE through issue #102 / PR #103 / head `c52b48be1a82eb27670f03bdd4e1be4df6eb9f54` / run-job `33132237120 / 98724184966` / merge `0c6b075c272aeac5a6e5f9d4210b1c5a30a040ce`.

P03.06 remains DONE through issue #106 / PR #107 / head `c895f44a1383d1c1d9c5fd23c95d7864810353c3` / run-job `33181421854 / 98883286556` / merge `13dbe8a393c20cabeb8aac60d073a6c66775efd3`.

P01/P02 and completed P03 regressions remain mandatory during later P03.08 implementation. Historical completion evidence is never rewritten by later transitions.

## Retained architecture/security baselines

ADR-0012 remains accepted. P03.03 preserves bounded manifest-version dependency semantics. P03.04 preserves fail-closed lifecycle, non-destructive disable/re-enable, dependency/reverse-dependency protections, guarded purge, deterministic recovery and unrelated-module isolation. P03.05 preserves `kernel.configuration` as authoritative for settings/feature flags and keeps flags non-authorizing. P03.06 preserves validated capability metadata as non-invoking and non-authorizing. P03.07 preserves `kernel.authorization` as the sole deny-by-default permission enforcement/policy authority.

P03.07 registration creates no role/principal grant, role-name bypass or tenant-scope authority. Unknown/unavailable module permissions deny. Optional capability association is descriptive only. Historical policy/role/audit references are preserved through lifecycle changes.

## P03.08 implementation contract

Owner: `kernel.modules`.

P03.08 implementation starts only after this P03.07 closure / P03.08 activation carrier passes exact-final-head governance, merges, and protected `main` plus canonical state are re-read. A new separate implementation branch must be created from that exact post-merge SHA.

Authorized P03.08 scope is limited to the canonical package specification and includes:

- stable contribution ID/module owner/slot metadata;
- declarative contribution categories bounded to navigation/pages/widgets/settings/builder-slot metadata;
- permission requirement metadata without grants;
- module/lifecycle availability metadata;
- feature-flag condition metadata without authorization authority;
- optional-dependency condition and selective fallback/degradation metadata;
- deterministic collision/unknown-slot validation;
- versioned contribution metadata suitable for future Experience/Portal/Product Federation consumers without implementing those runtimes;
- focused positive/adversarial tests and a dedicated P03.08 verifier.

UI contribution metadata is declarative only. It never replaces backend authorization and cannot expose secret values, raw tenant/org authority, unrestricted executable code, private handlers or cross-module database-write shortcuts.

## Explicitly unauthorized

- P03.08 runtime implementation on this closure carrier;
- rendering framework/component implementation;
- CMS/Experience Builder runtime;
- portal runtime;
- product federation/unified work-surface runtime;
- frontend-only authorization;
- P03.09 migration ownership/execution;
- P03.10 module health runtime;
- P03.11 trust hooks/phase-exit runtime;
- generic remote RPC/service mesh;
- workflow orchestration;
- package installation/download/marketplace runtime;
- P04+ event/data orchestration;
- business modules/features;
- strategic X-program runtime;
- AI/model/agent runtime;
- weakening retained regressions, branch protection, security or governance gates.

## Exact next action

1. Finish atomic P03.07 closure / P03.08 activation reconciliation under GitHub issue #112.
2. Preserve every failed/cancelled/stale run as diagnostic evidence; do not relabel it.
3. Keep the carrier inside the approved governance/evidence/continuity boundary and include no P03.08 runtime code.
4. Require the exact final closure head to pass canonical GitHub-hosted governance with all retained regressions including `scripts/verify_p03_07.sh`.
5. Merge only if current with protected `main` and repository review/conversation gates permit it.
6. Re-read protected `main`, `STATE.json`, `STATUS.md`, package sequence and P03.08 handoff after merge and record the exact new SHA.
7. Identify a separate P03.08 implementation branch from that exact SHA as the next authorized action and stop the closure transaction without implementing P03.08.
8. Keep P03.09+ locked; P03.08 completion/state transition requires a later separate closure PR.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.08.md`, `docs/roadmap/evidence/P03.07_COMPLETION_2026-08-28.md`, accepted `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, completed `docs/ai/handoffs/P03.07.md` and activation-candidate `docs/ai/handoffs/P03.08.md`.
