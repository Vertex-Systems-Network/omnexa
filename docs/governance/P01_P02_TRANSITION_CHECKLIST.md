# P01 → P02 Transition Checklist

Status: **ACTIVATION CANDIDATE — GOVERNANCE VERIFICATION REQUIRED**  
Owner transition: `P01 done -> P02.01 active`

This checklist governs a state-only activation transition. It contains no P02 runtime/schema implementation.

## A. Preconditions

- [x] P00 remains done and Foundation Architecture v1 remains frozen.
- [x] P01.01-P01.12 are done.
- [x] `docs/governance/P01_EXIT_GATE.md` is SATISFIED.
- [x] protected `main` and required `governance` controls remain the integration authority.
- [x] canonical CI is GitHub-hosted `ubuntu-24.04` only.
- [x] P02 purpose/exit criteria are defined by `docs/roadmap/MASTER_PLAN.md`.
- [x] P02 ownership maps to `kernel.identity`, `kernel.tenancy`, `kernel.organization`, `kernel.authorization`, `kernel.configuration` and `kernel.audit`.
- [x] P02 package sequence contains exactly 10 strict sequential work packages.
- [x] P02 entry and exit gates are defined.
- [x] P02 readiness/package validators are wired into canonical governance.
- [x] Preparation PR #67 merged after final exact-head canonical run `32632920772 / 97178312240` PASS.
- [x] No P02 implementation exists in the readiness preparation.
- [x] Issue #4 remains an external distribution gate rather than an internal P02 blocker.

## B. Activation state in this transition

This activation candidate atomically establishes:

- [x] `current_phase=P02`.
- [x] `current_work_package=P02.01`.
- [x] P00/P01 remain `done`.
- [x] P02 -> `active`.
- [x] P02.01 -> `active`.
- [x] P02.02-P02.10 remain `planned`.
- [x] `kernel_code_authorized=true` only for P02.01.
- [x] `business_feature_code_authorized=false`.
- [x] P02 entry gate -> `SATISFIED`.
- [x] P02 package manifest -> `active / implementation_authorized=true`.
- [x] no P02 runtime/schema code is included in the activation transition.

## C. Required pre-merge verification

Before activation merge, all of the following must be observed on the exact PR head:

- [ ] `governance` PASS on GitHub-hosted `ubuntu-24.04`.
- [ ] repository Go quality PASS.
- [ ] P01.01-P01.12 regression verifiers PASS.
- [ ] P02 preparation validator PASS in active mode.
- [ ] P02 package-sequence validator PASS in active mode.
- [ ] branch is current with protected `main`.
- [ ] PR is mergeable with required conversations resolved.
- [ ] diff contains no P02 runtime/schema implementation.

The PR/merge record is the evidence for section C; these boxes are intentionally not used as a substitute for live GitHub evidence.

## D. Post-merge rule

After activation merges, verify protected `main` and canonical `STATE.json`, identify P02.01 as the sole authorized implementation scope, then **STOP that execution session**. P02.01 implementation starts only in a later governed execution session.

## Activation target

```text
P01: DONE — 12 / 12
P01 exit: SATISFIED
P02 specifications: 10 / 10 prepared
P02: ACTIVE
P02.01: ACTIVE
P02.02-P02.10: PLANNED
kernel_code_authorized: true — P02.01 only
business_feature_code_authorized: false
```
