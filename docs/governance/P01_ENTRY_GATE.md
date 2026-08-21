# Omnexa P01 Implementation Entry Gate

Status: **BLOCKED — GitHub plan-limited hosted protection**  
Owner phase: **P00.10 -> P01 transition**

P01 kernel implementation is not authorized merely because the P00 architecture documents exist. This gate is the operational handoff between architecture freeze and executable kernel work.

## Required before P01 executable merge

### EG-01 — Foundation architecture freeze accepted
State: **SATISFIED**

Evidence:
- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- ADR-0010.

### EG-02 — `main` integration protection
State: **BLOCKED_BY_PLAN**
Tracker: **Issue #3**

Required policy remains:
- PR-based integration enforced for `main`;
- force push blocked;
- branch deletion blocked;
- conversation resolution policy configured;
- bypass restricted to explicit break-glass actors;
- required `governance` check attached to the protection/ruleset;
- protection verified through GitHub, not assumed from documentation.

#### Decisive 2026-08-22 evidence

The owner executed the merged administration tooling against the private organization repository and GitHub returned:

```text
gh: Upgrade to GitHub Pro or make this repository public to enable this feature. (HTTP 403)
```

Verified context:
- repository visibility: `private`;
- repository owner: `Vertex-Systems-Network` organization;
- authenticated user permission: repository `admin`;
- branch-protection tooling syntax/structure: PASS on `LOCAL-WIN-4`;
- live branch metadata after the failed API write: `protected: false`.

GitHub's plan model makes private branch protection/rulesets unavailable on the current plan. This is therefore a hosted product-plan blocker, not an authentication, authorization, script or CI defect.

Do **not** retry the same API call until one of these conditions changes:

1. the organization/repository is upgraded to a GitHub plan that supports private branch protection/rulesets;
2. repository visibility is intentionally changed to public through a separate owner/legal/security decision; or
3. a superseding governance ADR explicitly accepts a compensating control and changes EG-02.

Option 3 is an architecture/governance change and must not be inferred from a generic `continue` instruction.

### EG-03 — executable verification lane
State: **SATISFIED**
Tracker: **Issue #14 — closed/completed**

Current canonical evidence:
- PR #23 migrated `Omnexa Governance` to require successful validation evidence produced specifically by runner `LOCAL-WIN-4`;
- runner `LOCAL-WIN-4` executed on Windows X64 / machine `ABDUL-HANAN` from `C:\actions-runner-4\_work`;
- Git `2.55.0.windows.5` and Python `3.13.7` were available;
- workflow run `32528329184` executed all four validators successfully on `LOCAL-WIN-4` and the final job named `governance` completed **SUCCESS**;
- `scripts/validate_governance.py` PASS;
- `scripts/validate_development_spec.py` PASS;
- `scripts/validate_operations_spec.py` PASS;
- `scripts/validate_freeze_review.py` PASS;
- PR #23 merged to `main` as `1a14362e2ed52a20d66cec6f28b93a2ee457f9a9`.

Historical evidence remains PR #20 / run `32522919774`, which originally restored the self-hosted CI lane. GitHub-hosted runner quota is no longer required for the canonical governance lane.

`LOCAL-WIN-4` currently has no unique schedulable Actions label. The canonical workflow therefore fails closed: it fans out only across local Windows/X64 self-hosted runners, runs validators only where `RUNNER_NAME == LOCAL-WIN-4`, uploads target pass evidence only from that runner, and exposes a final `governance` job that fails if no LOCAL-WIN-4 evidence exists.

### EG-04 — canonical local verification command exists
State: **PENDING P01 BOOTSTRAP / MUST BE FIRST-CLASS**

P01 may create the executable `omnexa verify`/equivalent command family defined by P00.07/P00.08, but the first executable PR must not merge unless its applicable verification runs through the satisfied executable-CI lane.

### EG-05 — implementation locks transition atomically
State: **BLOCKED BY EG-02 ONLY**

The transition PR must:

- verify Issue #3 / hosted `main` protection, or reference a superseding accepted governance ADR that deliberately replaces EG-02;
- mark P00.10 done;
- mark P00 phase done;
- retire the temporary P00 CI exception from active use;
- activate P01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- define the first P01 work package;
- record branch-protection/compensating-control and CI evidence.

Do not combine the transition with unrelated kernel feature code.

## External distribution gate

Issue #4 licensing/IP/trademark is **not** an internal P01 engineering blocker, but remains a hard gate before external/public distribution, self-hosted customer delivery or public launch.

Changing this private repository to public merely to obtain GitHub Free branch protection is not an implicit workaround; it requires an explicit owner/legal/security decision and must consider Issue #4 first.

## Current decision

```text
P00 architecture: FROZEN
Executable CI lane: SATISFIED ON LOCAL-WIN-4
EG-02: BLOCKED_BY_PLAN
P00 exit: BLOCKED ON ISSUE #3 / EG-02 ONLY
P01 implementation: NOT AUTHORIZED
Kernel code: LOCKED
Business feature code: LOCKED
```
