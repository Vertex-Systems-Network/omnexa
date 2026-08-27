# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected main `13701e7647c1e084dfe4288d4b27b3ddd75e72c2` contains the completed P03.04 implementation. This separate closure/state-transition carrier will make P03.05 authoritative only if its exact final head passes canonical governance and merges.

Transition candidate state:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 4 / 11 done after this closure merges.
- P03.01-P03.04: DONE with canonical completion evidence.
- current work package after closure merge: P03.05 — Module Settings & Feature Flags.
- P03.06-P03.11: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P03.05 after closure merge.
- `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` and live protected-main state remain authoritative. No P03.05 implementation may start from this branch merely because continuity docs describe the intended transition.

## P03.04 canonical completion evidence

- implementation PR: #100
- final exact implementation head: `cddb42d4466e7f97a7547c4cf5ea0812c768ff0b`
- implementation merge: `13701e7647c1e084dfe4288d4b27b3ddd75e72c2`
- canonical run/job: `33125377739 / 98702150001` — PASS
- runner: GitHub-hosted `ubuntu-24.04`
- Go: repository-pinned `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.04_COMPLETION_2026-08-28.md`
- retained verifier: `scripts/verify_p03_04.sh`

The accepted run passed repository Go quality, all retained P01/P02/P03.01/P03.02/P03.03 regressions and the dedicated P03.04 verifier. Earlier failed/stale P03.04 candidates remain diagnostic history only.

## Retained P03 prerequisites

P03.01 remains DONE through PR #92 / head `87da3302605c852ae5bf43d473aaa01a9e1aaa74` / run-job `33009396644 / 98311433013` / merge `4229e2a28442bf475afed143bab359a770d48053`.

P03.02 remains DONE through PR #94 / head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05` / run-job `33022405704 / 98355747775` / merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`.

P03.03 remains DONE through PR #98 / head `4dcaca22911fbb81b1d25af316fef146c4a71ff3` / run-job `33112808869 / 98659824107` / merge `774fab8b0350ffb2776517e3f1361f76bc2c68f9`.

P01/P02 and completed P03 regressions remain mandatory during P03.05. Historical completion evidence is not rewritten by later transitions.

## Accepted ADR-0012 baseline

ADR-0012 remains accepted on protected `main`. P03.03 implemented its bounded v1/v2 manifest dispatch, strict schema-v2 dependency requirements, strict bounded SemVer comparison, registry-bound validated manifest snapshots, required dependency fail-closed semantics, deterministic required graph ordering/cycle rejection and selective optional degradation.

Resolver metadata creates no permission, capability, tenant, database, lifecycle or private-package authority.

## P03.04 lifecycle baseline retained

P03.04 established the explicit fail-closed module lifecycle boundary. Retained invariants include deterministic legal transitions, dependency/reverse-dependency guards, non-destructive disable/re-enable, authorization/audit/dependency guarded purge, idempotent replay, CAS/concurrency protection, deterministic recovery and unrelated-module isolation.

Lifecycle state creates no permission, capability, tenant, database or configuration authority.

## P03.05 implementation contract

Owner: `kernel.modules` with `kernel.configuration` integration.

P03.05 implementation starts only after this closure merges and protected `main` is re-read. A new separate implementation branch must be created from that exact post-merge SHA.

Authorized P03.05 scope is limited to the canonical package specification and includes:

- registration of manifest-declared setting and feature-flag definitions;
- stable owner/module keys, types, validation and defaults;
- lifecycle-aware availability/read semantics;
- tenant/org scoping only through trusted existing configuration context;
- safe unregister/disable behavior that preserves required historical configuration;
- collision/ownership validation;
- negative tests proving feature flags do not grant permissions or bypass `kernel.authorization`;
- proof that no duplicate configuration subsystem is introduced.

`kernel.configuration` remains the authoritative setting/flag state owner. P03.05 may integrate with it but may not replace or duplicate it.

## Explicitly unauthorized

- P03.05 runtime implementation on this closure carrier;
- business-feature configuration domains;
- entitlement/licensing product logic;
- permission-granting flags or hidden authorization paths;
- P03.06-P03.11 later runtime;
- P04+ event/data orchestration;
- business modules/features;
- package installation/download/marketplace runtime;
- full strategic X-program runtime;
- AI/model/agent runtime;
- weakening retained regressions, branch protection, security or governance gates.

## Exact next action

1. Finish atomic P03.04 closure / P03.05 activation reconciliation on this governance carrier.
2. Preserve every failed/cancelled/stale run as diagnostic evidence; do not relabel it.
3. Require the exact final closure head to pass canonical GitHub-hosted governance with all retained regressions including `scripts/verify_p03_04.sh`.
4. Merge only if current with protected `main` and repository review/conversation gates permit it.
5. Re-read protected `main` and `STATE.json` after merge and record the exact new SHA.
6. Create a separate P03.05 implementation branch from that exact SHA.
7. Keep P03.06+ locked; P03.05 completion/state transition requires a later separate closure PR.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.05.md`, `docs/roadmap/evidence/P03.04_COMPLETION_2026-08-28.md`, accepted `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, completed `docs/ai/handoffs/P03.04.md` and active-after-closure `docs/ai/handoffs/P03.05.md`.