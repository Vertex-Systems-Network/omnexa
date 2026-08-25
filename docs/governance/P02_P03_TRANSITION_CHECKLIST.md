# P02 → P03 Transition Checklist

Status: **READY — ACTIVATION PENDING**  
Owner transition: `P02 done -> P03.01 active`

This checklist prepares a governance/state-only activation transition. It contains no P03 runtime/schema implementation.

## A. Preconditions

- [x] P00 remains done and Foundation Architecture v1 remains frozen.
- [x] P01.01-P01.12 are done and regressions remain mandatory.
- [x] P02.01-P02.10 are done.
- [x] `docs/governance/P02_EXIT_GATE.md` is SATISFIED.
- [x] protected `main` and required `governance` controls remain the integration authority.
- [x] canonical CI is GitHub-hosted `ubuntu-24.04` only.
- [x] P03 purpose/exit criteria are defined by `docs/roadmap/MASTER_PLAN.md`.
- [x] module ownership/dependency rules remain defined by `DOMAIN_OWNERSHIP.md`, `DEPENDENCY_MATRIX.md` and `MODULE_STANDARD.md`.
- [x] P03 package sequence contains exactly 11 strict sequential work packages.
- [x] P03 entry and exit gates are defined.
- [x] P03 AI-native compatibility mapping is documented without strategic-program activation.
- [x] P03 readiness/package validators are wired into canonical governance.
- [x] no P03 runtime/schema implementation exists in this readiness preparation.
- [x] Issue #4 remains an external distribution gate rather than implicit P03 authorization.

## B. Activation state to be applied separately

The later activation PR must atomically establish:

- [ ] `current_phase=P03`.
- [ ] `current_work_package=P03.01`.
- [ ] P00/P01/P02 remain `done`.
- [ ] P03 -> `active`.
- [ ] P03.01 -> `active`.
- [ ] P03.02-P03.11 remain `planned`.
- [ ] `kernel_code_authorized=true` only for P03.01.
- [ ] `business_feature_code_authorized=false`.
- [ ] P03 entry gate -> `SATISFIED`.
- [ ] P03 package manifest -> `active / implementation_authorized=true`.
- [ ] canonical `p03_preparation` metadata is added/reconciled for the active package.
- [ ] no P03 runtime/schema code is included in the activation PR.

## C. Activation verification

Before activation merge:

- [ ] exact-head `governance` PASS on GitHub-hosted `ubuntu-24.04`.
- [ ] repository Go quality PASS.
- [ ] P01.01-P01.12 regression verifiers PASS.
- [ ] P02.01-P02.10 regression verifiers PASS.
- [ ] P03 preparation validator PASS in active mode.
- [ ] P03 package-sequence validator PASS in active mode.
- [ ] branch is current with protected `main`.
- [ ] PR is mergeable with required conversations resolved.
- [ ] diff contains no P03 runtime/schema implementation.

## D. Post-merge rule

After activation merges, verify protected `main` and canonical `STATE.json`, identify P03.01 as the sole authorized implementation scope, then **STOP that execution session**. P03.01 implementation starts only in a later governed execution session.

## Current readiness checkpoint

```text
P02: DONE — 10 / 10
P02 exit: SATISFIED
P03 specifications: 11 / 11 prepared
P03: PLANNED / NOT ACTIVE
P03.01-P03.11: PLANNED
kernel_code_authorized: false
business_feature_code_authorized: false
```
