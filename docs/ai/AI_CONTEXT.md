# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected main `0c6b075c272aeac5a6e5f9d4210b1c5a30a040ce` contains the completed P03.05 implementation. This separate closure/state-transition carrier will record P03.05 completion and make P03.06 authoritative only if its exact final head passes canonical governance and merges.

Transition candidate state:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 5 / 11 done after this closure merges.
- P03.01-P03.05: DONE with canonical completion evidence.
- current work package after closure merge: P03.06 — Capability Registry.
- P03.07-P03.11: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P03.06 after closure merge.
- `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` and live protected-main state remain authoritative. No P03.06 implementation may start from this branch merely because continuity docs describe the intended transition.

## P03.05 canonical completion evidence

- implementation issue: #102 — completed
- implementation PR: #103 — merged
- activation/base main: `f38daddf2a4b71290cadef02cbad8c1afaaca15f`
- final exact implementation head: `c52b48be1a82eb27670f03bdd4e1be4df6eb9f54`
- implementation merge / protected-main closure base: `0c6b075c272aeac5a6e5f9d4210b1c5a30a040ce`
- canonical run/job: `33132237120 / 98724184966` — PASS
- runner: GitHub-hosted `ubuntu-24.04`
- Go: repository-pinned `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.05_COMPLETION_2026-08-28.md`
- retained verifier: `scripts/verify_p03_05.sh`

The accepted exact-head run passed repository Go quality, all retained P01/P02/P03.01/P03.02/P03.03/P03.04 regressions and the dedicated P03.05 verifier. Earlier failed/stale P03.05 candidates remain diagnostic history only; only the exact final head/run above is acceptance authority.

## Retained P03 prerequisites

P03.01 remains DONE through PR #92 / head `87da3302605c852ae5bf43d473aaa01a9e1aaa74` / run-job `33009396644 / 98311433013` / merge `4229e2a28442bf475afed143bab359a770d48053`.

P03.02 remains DONE through PR #94 / head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05` / run-job `33022405704 / 98355747775` / merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`.

P03.03 remains DONE through PR #98 / head `4dcaca22911fbb81b1d25af316fef146c4a71ff3` / run-job `33112808869 / 98659824107` / merge `774fab8b0350ffb2776517e3f1361f76bc2c68f9`.

P03.04 remains DONE through PR #100 / head `cddb42d4466e7f97a7547c4cf5ea0812c768ff0b` / run-job `33125377739 / 98702150001` / merge `13701e7647c1e084dfe4288d4b27b3ddd75e72c2`.

P01/P02 and completed P03 regressions remain mandatory during later P03.06 implementation. Historical completion evidence is not rewritten by later transitions.

## Accepted ADR-0012 baseline

ADR-0012 remains accepted on protected `main`. P03.03 implemented its bounded v1/v2 manifest dispatch, strict schema-v2 dependency requirements, strict bounded SemVer comparison, registry-bound validated manifest snapshots, required dependency fail-closed semantics, deterministic required graph ordering/cycle rejection and selective optional degradation.

Resolver metadata creates no permission, capability, tenant, database, lifecycle or private-package authority.

## P03.04 lifecycle baseline retained

P03.04 established the explicit fail-closed module lifecycle boundary. Retained invariants include deterministic legal transitions, dependency/reverse-dependency guards, non-destructive disable/re-enable, authorization/audit/dependency guarded purge, idempotent replay, CAS/concurrency protection, deterministic recovery and unrelated-module isolation.

Lifecycle state creates no permission, capability, tenant, database or configuration authority.

## P03.05 configuration baseline retained

P03.05 preserves manifest schema v1/v2 and binds manifest-declared settings/feature flags to the existing `kernel.configuration` authority.

Retained invariants include:

- exact fully-qualified configuration keys rather than an encoded mini-language;
- validated discovery snapshot provenance rather than a second raw-manifest parse authority;
- existing typed `kernel.configuration.Definition` ownership/type/default semantics;
- explicit global/scoped registration scope;
- existing P02.09 `SettingPolicy` validation for scoped registrations;
- trusted tenant/organization context only;
- lifecycle-aware non-destructive configuration history;
- collision/ownership/class/scope validation that fails closed;
- no permission-granting feature flag behavior;
- no duplicate configuration subsystem.

Settings/feature flags create no permission, capability, tenant, private-package or database authority.

## P03.06 implementation contract

Owner: `kernel.modules`.

P03.06 implementation starts only after this P03.05 closure / P03.06 activation carrier passes exact-final-head governance, merges, and protected `main` plus canonical state are re-read. A new separate implementation branch must be created from that exact post-merge SHA.

Authorized P03.06 scope is limited to the canonical package specification and includes:

- stable capability ID and bounded major-version metadata;
- provider module/owner and consumer declarations;
- lifecycle-derived availability state;
- authorization and tenant/org scope requirement metadata references;
- contract metadata references without private implementation exposure;
- collision/version compatibility validation;
- lifecycle-aware registration/withdrawal availability;
- deterministic lookup suitable for later attribution/graph consumers;
- dedicated positive/adversarial tests and P03.06 verifier.

Capability registration is metadata/availability only. It does not grant invocation authority. The owning capability boundary remains responsible for validation, authorization, trusted tenant/org enforcement and audit.

## Explicitly unauthorized

- P03.06 runtime implementation on this closure carrier;
- business capability implementations;
- P03.07 permission registration;
- P03.08 UI contribution runtime;
- P03.09 migration ownership/execution;
- P03.10 module health runtime;
- P03.11 trust hooks/phase-exit runtime;
- generic remote RPC/service mesh;
- workflow orchestration;
- product federation runtime;
- package installation/download/marketplace runtime;
- P04+ event/data orchestration;
- business modules/features;
- full strategic X-program runtime;
- AI/model/agent runtime;
- weakening retained regressions, branch protection, security or governance gates.

## Exact next action

1. Finish atomic P03.05 closure / P03.06 activation reconciliation under GitHub issue #104.
2. Preserve every failed/cancelled/stale run as diagnostic evidence; do not relabel it.
3. Keep the carrier inside the approved 12-file governance/evidence/continuity boundary and include no P03.06 runtime code.
4. Require the exact final closure head to pass canonical GitHub-hosted governance with all retained regressions including `scripts/verify_p03_05.sh`.
5. Merge only if current with protected `main` and repository review/conversation gates permit it.
6. Re-read protected `main`, `STATE.json`, `STATUS.md`, package sequence and P03.06 handoff after merge and record the exact new SHA.
7. Create a separate P03.06 implementation branch from that exact SHA.
8. Keep P03.07+ locked; P03.06 completion/state transition requires a later separate closure PR.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.06.md`, `docs/roadmap/evidence/P03.05_COMPLETION_2026-08-28.md`, accepted `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, completed `docs/ai/handoffs/P03.05.md` and active-after-closure `docs/ai/handoffs/P03.06.md`.