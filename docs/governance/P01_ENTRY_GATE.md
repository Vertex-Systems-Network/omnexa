# Omnexa P01 Implementation Entry Gate

Status: **BLOCKED — GitHub plan-limited hosted protection**  
Owner phase: **P00.10 -> P01 transition**

P01 kernel implementation is not authorized merely because the P00 architecture documents exist. This gate is the operational handoff between architecture freeze and executable kernel work.

Prepared handoff artifacts:

- `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P01.01.md`
- `scripts/validate_p01_preparation.py`

These artifacts reduce transition ambiguity but do not satisfy EG-02 or authorize executable code.

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
- branch-protection tooling syntax/structure: PASS;
- live branch metadata after the failed API write: `protected: false`.

GitHub's plan model makes private branch protection/rulesets unavailable on the current plan. This is therefore a hosted product-plan blocker, not an authentication, authorization, script or CI defect.

Do **not** retry the same API call until one of these conditions changes:

1. the organization/repository is upgraded to a GitHub plan that supports private branch protection/rulesets;
2. repository visibility is intentionally changed to public through a separate owner/legal/security decision; or
3. a superseding governance ADR explicitly accepts a compensating control and changes EG-02.

### EG-03 — executable verification lane
State: **SATISFIED**
Tracker: **Issue #14 — closed/completed**

The canonical workflow is now **runner-name agnostic** inside the approved local pool. The required job is named `governance` and runs on:

```yaml
runs-on: [self-hosted, Windows, X64]
```

Any available Windows/X64 self-hosted runner may execute the job. The workflow no longer pins `LOCAL-WIN-4`, no longer uses discovery fanout, and no longer waits for target-runner artifacts. The same fail-closed validators remain mandatory.

Current routing-policy evidence:
- workflow run `32535324900`: **SUCCESS**;
- governance job `96935023669`: **SUCCESS**;
- actual runner: `LOCAL-WIN-02`;
- runner OS/arch: Windows / X64;
- machine: `ABDUL-HANAN`;
- Git `2.55.0.windows.5`;
- Python `3.13.7`;
- PowerShell tooling parse: PASS;
- `scripts/validate_governance.py`: PASS;
- `scripts/validate_development_spec.py`: PASS;
- `scripts/validate_operations_spec.py`: PASS;
- `scripts/validate_freeze_review.py`: PASS;
- `scripts/validate_p01_preparation.py`: PASS.

Historical `LOCAL-WIN-4` evidence remains valid provenance: PR #23, run `32528329184`, target job `96915072868`, merge `1a14362e2ed52a20d66cec6f28b93a2ee457f9a9`. It is no longer a routing requirement.

GitHub-hosted standard runners may be introduced later if capacity/policy makes them appropriate, provided P00.07 gate semantics remain unchanged. Current canonical routing prefers the already operational self-hosted Windows/X64 pool because it is available and does not depend on hosted-minute entitlement.

While P01 remains blocked, `scripts/validate_p01_preparation.py` is required and rejects known executable P01.01 paths while `kernel_code_authorized=false`.

### EG-04 — canonical local verification command exists
State: **PENDING P01 BOOTSTRAP / MUST BE FIRST-CLASS**

P01 may create the executable `omnexa verify`/equivalent command family defined by P00.07/P00.08, but the first executable PR must not merge unless its applicable verification runs through the satisfied executable-CI lane.

The first package controlling this bootstrap is `docs/roadmap/work-packages/P01.01.md`.

### EG-05 — implementation locks transition atomically
State: **BLOCKED BY EG-02 ONLY**

The transition PR must follow `docs/governance/P00_P01_TRANSITION_CHECKLIST.md` and:

- verify Issue #3 / hosted `main` protection, or reference a superseding accepted governance ADR that deliberately replaces EG-02;
- mark P00.10 done;
- mark P00 phase done;
- retire the temporary P00 CI exception from active use;
- activate P01 and P01.01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- keep P01.02–P01.12 planned;
- record branch-protection/compensating-control and CI evidence.

Do not combine the transition with kernel implementation. The first kernel-code PR is separate and follows the merged transition.

## Prepared P01.01 scope

P01.01 is prepared as **Go workspace/build skeleton** and remains `planned`.

It may not be implemented yet. Its specification deliberately excludes configuration, DB, cache, storage, telemetry, health, jobs, feature flags, audit, identity/tenancy, module runtime and business-domain behavior. Those remain later packages/phases.

## External distribution gate

Issue #4 licensing/IP/trademark is **not** an internal P01 engineering blocker, but remains a hard gate before external/public distribution, self-hosted customer delivery, public launch or external contributor intake.

The owner/legal decision worksheet is `docs/governance/LICENSING_DECISION_BRIEF.md`. It does not change `LICENSE` or establish trademark clearance.

Changing this private repository to public merely to obtain GitHub Free branch protection is not an implicit workaround; it requires an explicit owner/legal/security decision and must consider Issue #4 first.

## Current decision

```text
P00 architecture: FROZEN
Executable CI lane: SATISFIED — ANY AVAILABLE WINDOWS/X64 SELF-HOSTED RUNNER
Latest routing proof: LOCAL-WIN-02 / run 32535324900
EG-02: BLOCKED_BY_PLAN
P00 exit: BLOCKED ON ISSUE #3 / EG-02 ONLY
P01.01: PREPARED / PLANNED
P01 implementation: NOT AUTHORIZED
Kernel code: LOCKED
Business feature code: LOCKED
```
