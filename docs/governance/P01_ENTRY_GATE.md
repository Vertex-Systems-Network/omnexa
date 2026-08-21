# Omnexa P01 Implementation Entry Gate

Status: **BLOCKED**  
Owner phase: **P00.10 -> P01 transition**

P01 kernel implementation is not authorized merely because the P00 architecture documents exist. This gate is the operational handoff between architecture freeze and executable kernel work.

## Required before P01 executable merge

### EG-01 — Foundation architecture freeze accepted
State: **SATISFIED**

Evidence:
- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- ADR-0010 once merged.

### EG-02 — `main` integration protection
State: **BLOCKED**
Tracker: **Issue #3**

Required evidence:
- PR-based integration enforced for `main`;
- force push blocked;
- branch deletion blocked;
- conversation resolution policy configured;
- bypass restricted to explicit break-glass actors;
- protection verified through GitHub, not assumed from documentation.

### EG-03 — executable verification lane
State: **BLOCKED**
Tracker: **Issue #14**

Required evidence:
- approved CI provider or self-hosted lane can allocate a runner and execute repository commands;
- clean checkout can run governance/static/unit/build gates applicable to kernel bootstrap;
- CI environment is reproducible from repository-pinned toolchains;
- secrets/permissions are least-privileged;
- blocked/unrun jobs are never green;
- required branch check(s) can be attached to the `main` protection policy.

GitHub Actions is a preferred available integration but is not architecturally mandatory; an approved equivalent lane satisfying P00.07 is acceptable.

### EG-04 — canonical local verification command exists
State: **PENDING P01 BOOTSTRAP / MUST BE FIRST-CLASS**

P01 may create the executable `omnexa verify`/equivalent command family defined by P00.07/P00.08, but the first executable PR must not merge unless its applicable verification can run through EG-03.

### EG-05 — implementation locks transition atomically
State: **BLOCKED BY EG-02/EG-03**

The transition PR must:

- mark P00.10 done;
- mark P00 phase done;
- expire the temporary P00 CI exception;
- activate P01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- define the first P01 work package;
- record protection/CI evidence.

Do not combine the transition with unrelated kernel feature code.

## External distribution gate

Issue #4 licensing/IP/trademark is **not** an internal P01 engineering blocker, but remains a hard gate before external/public distribution, self-hosted customer delivery or public launch.

## Current decision

```text
P00 architecture: frozen
P00 exit: blocked on EG-02 + EG-03
P01 implementation: NOT AUTHORIZED
Kernel code: LOCKED
Business feature code: LOCKED
```
