# P00 → P01 Transition Checklist

Status: **EXECUTING — governance transition; effective after verified merge**  
Owner transition: `P00.10 -> P01.01`

This checklist records the atomic handoff from the completed Foundation Program to the first executable kernel package. The transition PR is governance/state-only and contains no kernel implementation.

## A. Preconditions

- [x] P00.01–P00.09 remain done/frozen.
- [x] P00.10 freeze-review artifacts are present and internally consistent.
- [x] `FOUNDATION_FREEZE.json` reports `architecture_status=FROZEN`.
- [x] Executable governance CI is GitHub-hosted `ubuntu-24.04` only.
- [x] Required job is `governance` and proves `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64.
- [x] No canonical governance workflow contains self-hosted/LOCAL-WIN routing.
- [x] Issue #14 is closed/satisfied.
- [x] Issue #3 is closed/satisfied.
- [x] Live GitHub API reports `main.protected=true`.
- [x] Failed `governance` blocks merge (PR #34 / run `32540836431`).
- [x] Direct update to `main` is rejected (probe commit `44ca19e80c5fccccebfd8d4f96dde6dc5af14bc2`).
- [x] Force update is rejected.
- [x] Conversation resolution is enforced (PR #37 / run `32541439589`).
- [x] Valid green PR integration succeeds (PR #35 merge `843c615170058ab900ba69516dbed80a47f26973`).
- [x] Deletion of `main` remains blocked by configured ruleset; destructive deletion was intentionally not attempted.
- [x] Current review policy is compatible with the single-maintainer model: zero required approvals; CODEOWNERS documents ownership.
- [x] No architecture contradiction has been introduced since ADR-0010.
- [x] Repository is free of unrelated project code.

## B. Transition PR scope

Required state on this branch:

- [x] `P00.10 -> done` / P00 exit `DONE`.
- [x] P00 phase -> `done`.
- [x] P01 phase -> `active`.
- [x] `current_phase -> P01`.
- [x] `current_work_package -> P01.01`.
- [x] `kernel_code_authorized -> true`.
- [x] `business_feature_code_authorized -> false`.
- [x] P01.01 -> `active`.
- [x] P01.02–P01.12 remain `planned`.
- [x] ADR-0006 -> expired/historical-only.
- [x] Issue #3 closing evidence recorded without deleting history.
- [x] `P01_ENTRY_GATE.md` -> SATISFIED.
- [x] `FOUNDATION_FREEZE.json` -> `p00_exit_status=DONE`.
- [x] STATUS, STATE, README and AGENTS reconciled.
- [x] Transition ledger entry prepared.
- [x] No executable Go/kernel implementation included.

## C. Verification required on this transition PR

Before merge all must be PASS on GitHub-hosted `ubuntu-24.04`:

- [ ] `scripts/validate_governance.py`;
- [ ] `scripts/validate_development_spec.py`;
- [ ] `scripts/validate_operations_spec.py`;
- [ ] `scripts/validate_freeze_review.py`;
- [ ] `scripts/validate_p01_preparation.py`;
- [ ] `scripts/validate_p01_package_specs.py`;
- [ ] required `governance` job.

These boxes become evidenced by the PR run rather than by editing them optimistically. `BLOCKED`, `NOT RUN` and `N/A` are not PASS.

## D. P01.01 activation invariants

- kernel code is authorized only for active P01.01;
- business-domain implementation remains prohibited;
- no P02 identity/tenant implementation;
- no P03 module runtime;
- no persistence/infrastructure owned by P01.04+;
- every executable PR uses GitHub-hosted canonical verification;
- `P01.01.md` controls implementation;
- strict sequential P01.01→P01.12 activation remains in force.

## E. Post-merge verification

After transition merge:

- [ ] fetch canonical `main` SHA;
- [ ] verify `STATE.json` reports P01 active/P01.01 active;
- [ ] verify `kernel_code_authorized=true`;
- [ ] verify `business_feature_code_authorized=false`;
- [ ] verify exactly one active P01 package;
- [ ] verify canonical hosted governance remains green;
- [ ] append immutable transition merge SHA/run evidence through reconciliation if not knowable inside this PR.

## F. Abort conditions

Abort if live `main` loses protection, a local/self-hosted runner is used, final governance fails, business code becomes authorized, multiple P01 packages become active, P01.01 conflicts with frozen P00 architecture, or this transition PR contains executable kernel code.

## Intended result after merge

```text
Foundation Architecture v1: FROZEN
P00: DONE — 10 / 10
P00.10: DONE
EG-02 / Issue #3: SATISFIED / CLOSED
EG-03 / Issue #14: SATISFIED
Canonical CI: GITHUB-HOSTED ONLY / ubuntu-24.04
P01: ACTIVE
P01.01: ACTIVE
P01.02-P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```
