# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected main `774fab8b0350ffb2776517e3f1361f76bc2c68f9` contains the completed P03.03 implementation. Closure PR #99 is the separate state-transition carrier that will make P03.04 authoritative only if its exact final head passes canonical governance and merges.

Transition candidate state:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 3 / 11 done after this closure merges.
- P03.01-P03.03: DONE with canonical completion evidence.
- current work package after closure merge: P03.04 — Module Lifecycle State Machine.
- P03.05-P03.11: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P03.04 after closure merge.
- `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` and live protected-main state remain authoritative. No P03.04 implementation may start from this branch merely because continuity docs describe the intended transition.

## P03.03 canonical completion evidence

- implementation PR: #98
- final exact implementation head: `4dcaca22911fbb81b1d25af316fef146c4a71ff3`
- implementation merge: `774fab8b0350ffb2776517e3f1361f76bc2c68f9`
- canonical run/job: `33112808869 / 98659824107` — PASS
- runner: GitHub-hosted `ubuntu-24.04`
- Go: repository-pinned `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.03_COMPLETION_2026-08-28.md`
- retained verifier: `scripts/verify_p03_03.sh`

The accepted run passed repository Go quality, all retained P01/P02/P03.01/P03.02 regressions and the dedicated P03.03 verifier. Earlier #412/#414 failures and #413 cancellation remain diagnostic history only.

## Retained P03 prerequisites

P03.01 remains DONE through PR #92 / head `87da3302605c852ae5bf43d473aaa01a9e1aaa74` / run-job `33009396644 / 98311433013` / merge `4229e2a28442bf475afed143bab359a770d48053`.

P03.02 remains DONE through PR #94 / head `0c46db41b0d724a08ea1a78545b3c2debdd8cd05` / run-job `33022405704 / 98355747775` / merge `2e38969dbbbcfcf4765a114f449dc3fa960061d7`.

P01/P02 and completed P03 regressions remain mandatory during P03.04. Historical completion evidence is not rewritten by later transitions.

## Accepted ADR-0012 baseline

ADR-0012 remains accepted on protected `main`. P03.03 implemented its bounded v1/v2 manifest dispatch, strict schema-v2 dependency requirements, strict bounded SemVer comparison, registry-bound validated manifest snapshots, required dependency fail-closed semantics, deterministic required graph ordering/cycle rejection and selective optional degradation.

Resolver metadata creates no permission, capability, tenant, database, lifecycle or private-package authority.

## P03.04 implementation contract

Owner: `kernel.modules`.

P03.04 implementation starts only after closure PR #99 merges and protected `main` is re-read. A new separate implementation branch must be created from that exact post-merge SHA.

Authorized P03.04 scope is limited to the canonical package specification and includes:

- explicit lifecycle states and transition validation;
- dependency/reverse-dependency lifecycle preconditions;
- safe non-destructive disable/re-enable semantics;
- suspend/archive/detach/purge distinctions;
- explicit destructive purge boundary with audit requirements;
- idempotency, retry and concurrency protection;
- deterministic partial-failure recovery semantics;
- upgrade lifecycle coordination without pulling later migration execution forward.

## Explicitly unauthorized

- P03.04 runtime implementation on closure PR #99;
- P03.05-P03.11 later runtime;
- P04+ event/data orchestration;
- business modules/features;
- package installation/download/marketplace runtime;
- full strategic X-program runtime;
- AI/model/agent runtime;
- weakening retained regressions, branch protection, security or governance gates.

## Exact next action

1. Finish atomic P03.03 closure / P03.04 activation reconciliation in PR #99.
2. Preserve #419/#420 and any later failed runs as diagnostic evidence; do not relabel them.
3. Require the exact final closure head to pass canonical GitHub-hosted governance with all retained regressions.
4. Merge only if current with protected `main` and repository review/conversation gates permit it.
5. Re-read protected `main` and `STATE.json` after merge and record the exact new SHA.
6. Create a separate P03.04 implementation branch from that exact SHA.
7. Keep P03.05+ locked; P03.04 completion/state transition requires a later separate closure PR.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.04.md`, `docs/roadmap/evidence/P03.03_COMPLETION_2026-08-28.md`, accepted `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md`, `docs/ai/handoffs/P03.03.md` and `docs/ai/handoffs/P03.04.md`.
