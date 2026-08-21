# Omnexa P01 Implementation Entry Gate

Status: **BLOCKED — one remaining entry blocker**  
Owner phase: **P00.10 -> P01 transition**

P01 kernel implementation is not authorized merely because the P00 architecture documents exist. This gate is the operational handoff between architecture freeze and executable kernel work.

## Required before P01 executable merge

### EG-01 — Foundation architecture freeze accepted
State: **SATISFIED**

Evidence:
- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- ADR-0010.

### EG-02 — `main` integration protection
State: **BLOCKED**
Tracker: **Issue #3**

Required evidence:
- PR-based integration enforced for `main`;
- force push blocked;
- branch deletion blocked;
- conversation resolution policy configured;
- bypass restricted to explicit break-glass actors;
- required `governance` check attached to the protection/ruleset;
- protection verified through GitHub, not assumed from documentation.

### EG-03 — executable verification lane
State: **SATISFIED**
Tracker: **Issue #14 — closed/completed**

Verified evidence:
- PR #20 moved `Omnexa Governance` to `runs-on: [self-hosted, Windows, X64]`;
- attached runner executed as `LOCAL-WIN-01` on Windows X64;
- Git and Python were available on the clean runner workspace;
- workflow run `32522919774`, governance job `96898839560`, completed **SUCCESS**;
- `scripts/validate_governance.py` PASS;
- `scripts/validate_development_spec.py` PASS;
- `scripts/validate_operations_spec.py` PASS;
- `scripts/validate_freeze_review.py` PASS;
- PR #20 merged to `main` as `c2ab2cd679c295a8dec84b1879acb9a9e02ad67d`.

GitHub-hosted runner quota is no longer required for the canonical governance lane. Future executable gates may expand on the same self-hosted lane or another approved lane, provided P00.07 semantics remain unchanged.

### EG-04 — canonical local verification command exists
State: **PENDING P01 BOOTSTRAP / MUST BE FIRST-CLASS**

P01 may create the executable `omnexa verify`/equivalent command family defined by P00.07/P00.08, but the first executable PR must not merge unless its applicable verification runs through the satisfied executable-CI lane.

### EG-05 — implementation locks transition atomically
State: **BLOCKED BY EG-02 ONLY**

The transition PR must:

- verify Issue #3 / hosted `main` protection;
- mark P00.10 done;
- mark P00 phase done;
- retire the temporary P00 CI exception from active use;
- activate P01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- define the first P01 work package;
- record branch-protection and CI evidence.

Do not combine the transition with unrelated kernel feature code.

## External distribution gate

Issue #4 licensing/IP/trademark is **not** an internal P01 engineering blocker, but remains a hard gate before external/public distribution, self-hosted customer delivery or public launch.

## Current decision

```text
P00 architecture: FROZEN
Executable CI lane: SATISFIED
P00 exit: BLOCKED ON EG-02 / ISSUE #3 ONLY
P01 implementation: NOT AUTHORIZED
Kernel code: LOCKED
Business feature code: LOCKED
```
