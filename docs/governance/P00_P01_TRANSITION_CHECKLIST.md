# P00 → P01 Transition Checklist

Status: **Prepared / not yet executable**  
Owner phase: `P00.10`  
Purpose: make the final foundation-to-kernel transition atomic, auditable and non-ambiguous.

This checklist does not itself authorize P01. It becomes executable only after the remaining P01 entry blocker is cleared through one of the governed paths in `P01_ENTRY_GATE.md`.

## A. Preconditions

All items must be true before the transition PR is opened as ready-to-merge:

- [ ] P00.01–P00.09 remain `done` and frozen.
- [ ] P00.10 review artifacts are present and internally consistent.
- [ ] `docs/governance/FOUNDATION_FREEZE.json` reports `architecture_status=FROZEN`.
- [ ] Executable governance CI is operational and produces evidence specifically from `LOCAL-WIN-4`.
- [ ] Final `governance` aggregator check passes.
- [ ] Issue #14 remains resolved/completed.
- [ ] EG-02 / Issue #3 is no longer blocked.
- [ ] If GitHub hosted protection is used, live API reports `main.protected=true` and required `governance` check is enforced.
- [ ] If EG-02 is superseded instead, an owner-approved ADR explicitly defines compensating controls, residual risk, expiry/review conditions and rollback path.
- [ ] No architecture contradiction has been introduced since ADR-0010.
- [ ] Repository is still free of unrelated project code.

## B. Transition PR scope

The transition PR is governance/state-only. It must not include kernel implementation.

Required changes:

- [ ] `P00.10.state: active -> done`.
- [ ] P00 phase state -> `done`.
- [ ] P00 `done_work_packages: 10`.
- [ ] P01 phase -> `active`.
- [ ] `current_phase -> P01`.
- [ ] `current_work_package -> P01.01`.
- [ ] `kernel_code_authorized -> true`.
- [ ] `business_feature_code_authorized -> false`.
- [ ] P01.01 -> `active`.
- [ ] P01.02–P01.12 remain `planned`.
- [ ] ADR-0006 explicitly becomes historical-only/inactive evidence.
- [ ] Issue #3 closing/supersession evidence recorded exactly.
- [ ] `P01_ENTRY_GATE.md` reflects `SATISFIED` rather than deleting history.
- [ ] `FOUNDATION_FREEZE.json` records P00 exit status `DONE`.
- [ ] `STATUS.md`, `STATE.json`, README and AGENTS are reconciled.
- [ ] Append-only execution-ledger transition entry added.

## C. Verification required on the transition PR

The transition PR must produce:

- [ ] `scripts/validate_governance.py` PASS.
- [ ] `scripts/validate_development_spec.py` PASS.
- [ ] `scripts/validate_operations_spec.py` PASS.
- [ ] `scripts/validate_freeze_review.py` PASS or its governed post-P00 successor.
- [ ] `scripts/validate_p01_preparation.py` PASS.
- [ ] LOCAL-WIN-4 target evidence PASS.
- [ ] final `governance` check PASS.

No `BLOCKED`, `NOT RUN` or manual-only exception may be relabeled PASS.

## D. P01.01 activation invariants

When P01 activates:

- kernel code is authorized **only** for the active P01 work package;
- business-domain implementation remains prohibited;
- no P02 identity/tenant implementation may be pulled forward;
- no P03 module runtime may be prebuilt;
- no persistence/infrastructure from P01.04+ may be added to P01.01;
- every executable PR must use the canonical quality lane;
- `docs/roadmap/work-packages/P01.01.md` is the controlling package specification.

## E. Post-merge verification

After the transition merge:

- [ ] fetch canonical `main` SHA;
- [ ] verify `STATE.json` from `main` reports P01 active/P01.01 active;
- [ ] verify `kernel_code_authorized=true`;
- [ ] verify `business_feature_code_authorized=false`;
- [ ] verify no second P01 package is active;
- [ ] verify latest governance run on canonical state passes;
- [ ] append immutable merge SHA/run evidence through a reconciliation PR if the transition PR cannot know its own merge SHA.

## F. Abort conditions

Do not complete the transition if any of these occur:

- hosted protection is assumed rather than observed;
- compensating-control ADR is proposed but not explicitly owner-approved;
- LOCAL-WIN-4 evidence is absent;
- final `governance` check fails/blocks;
- state would authorize business code;
- more than one P01 package becomes active;
- P01.01 specification conflicts with a frozen P00 architecture rule;
- the transition PR includes unrelated executable feature code.

## Current result

```text
P00 architecture: FROZEN
P00.10: EXIT VERIFICATION
EG-02 / Issue #3: BLOCKED_BY_PLAN
Executable CI: SATISFIED ON LOCAL-WIN-4
P01: BLOCKED
P01.01 specification: PREPARED / PLANNED
Kernel code: LOCKED
Business feature code: LOCKED
```
