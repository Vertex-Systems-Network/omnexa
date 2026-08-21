# Omnexa P01 Implementation Entry Gate

Status: **BLOCKED — main protection not yet applied**  
Owner phase: **P00.10 -> P01 transition**

P01 kernel implementation is not authorized merely because the P00 architecture documents exist. This gate is the operational handoff between architecture freeze and executable kernel work.

Prepared handoff artifacts:

- `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P01.01.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `scripts/validate_p01_preparation.py`
- `scripts/validate_p01_package_specs.py`

These artifacts reduce transition ambiguity but do not satisfy EG-02 or authorize executable code.

## Required before P01 executable merge

### EG-01 — Foundation architecture freeze accepted
State: **SATISFIED**

Evidence:
- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- ADR-0010.

### EG-02 — `main` integration protection
State: **ACTIONABLE_UNPROTECTED**
Tracker: **Issue #3**

Required policy remains:
- PR-based integration enforced for `main`;
- force push blocked;
- branch deletion blocked;
- conversation resolution policy configured;
- bypass restricted to explicit break-glass actors;
- required `governance` check attached to the protection/ruleset;
- protection verified through GitHub, not assumed from documentation.

#### Current 2026-08-22 evidence

Repository visibility is now intentionally **public**.

Verified context:
- repository visibility: `public`;
- repository owner: `Vertex-Systems-Network` organization;
- authenticated user permission: repository `admin`;
- branch-protection tooling syntax/structure: PASS;
- live branch metadata: `protected: false`.

The earlier private-repository attempt returned HTTP 403 with an upgrade/public-repository message. That historical plan-entitlement blocker is no longer the active cause because repository visibility has changed to public. EG-02 is now technically actionable, but it remains **unsatisfied** until live GitHub configuration reports the required protection and negative tests are recorded.

Do not authorize P01 merely because the plan blocker disappeared. `protected: false` is still a hard stop.

### EG-03 — executable verification lane
State: **SATISFIED**
Tracker: **Issue #14 — closed/completed**

The canonical governance workflow is now **GitHub-hosted only**:

```yaml
runs-on: ubuntu-24.04
```

The required job remains named `governance`. It fails closed unless the runtime reports `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. Local/self-hosted runners are prohibited from the canonical governance workflow.

Current hosted evidence:
- workflow run `32537207455`: **SUCCESS**;
- governance job `96940269306`: **SUCCESS**;
- runner: `GitHub Actions 1000006777`;
- runner environment: `github-hosted`;
- operating system: Ubuntu 24.04.4 LTS / X64;
- runner image: `ubuntu-24.04`, image version `20260816.277.1`;
- Git `2.55.0`;
- Python `3.12.3`;
- PowerShell `7.6.5`;
- branch-protection PowerShell tooling parse: PASS;
- `scripts/validate_governance.py`: PASS;
- `scripts/validate_development_spec.py`: PASS;
- `scripts/validate_operations_spec.py`: PASS;
- `scripts/validate_freeze_review.py`: PASS on the initial migration head;
- `scripts/validate_p01_preparation.py`: PASS;
- `scripts/validate_p01_package_specs.py`: PASS.

Earlier LOCAL-WIN runner evidence remains historical provenance only. It is not an active routing option and must not be reintroduced into the canonical workflow.

### EG-04 — canonical local verification command exists
State: **PENDING P01 BOOTSTRAP / MUST BE FIRST-CLASS**

P01 may create the executable `omnexa verify`/equivalent command family defined by P00.07/P00.08, but the first executable PR must not merge unless its applicable verification runs through the satisfied GitHub-hosted executable-CI lane.

The first package controlling this bootstrap is `docs/roadmap/work-packages/P01.01.md`.

### EG-05 — implementation locks transition atomically
State: **BLOCKED BY EG-02 ONLY**

The transition PR must follow `docs/governance/P00_P01_TRANSITION_CHECKLIST.md` and:

- verify Issue #3 / live hosted `main` protection, or reference a superseding accepted governance ADR that deliberately replaces EG-02;
- mark P00.10 done;
- mark P00 phase done;
- retire the temporary P00 CI exception from active use;
- activate P01 and P01.01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- keep P01.02–P01.12 planned;
- record branch-protection/compensating-control and GitHub-hosted CI evidence.

Do not combine the transition with kernel implementation. The first kernel-code PR is separate and follows the merged transition.

## Prepared P01 package sequence

P01.01–P01.12 are specification-complete under the machine-readable sequence at `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`.

P01.01 remains the next package and is **Go workspace/build skeleton**. All packages remain `planned`; strict sequential one-active-package activation applies. No executable P01 implementation is authorized yet.

## External distribution / public visibility

Issue #4 licensing/IP/trademark remains open. The repository is now public while the current repository `LICENSE` remains GPLv3 and trademark clearance is unresolved. That visibility change must be reconciled explicitly with the owner/legal decision path; it does not by itself satisfy or replace the P01 technical entry gates.

The owner/legal decision worksheet is `docs/governance/LICENSING_DECISION_BRIEF.md`.

## Current decision

```text
P00 architecture: FROZEN
Executable CI lane: SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04
Hosted proof: run 32537207455 / job 96940269306
Local/self-hosted governance runner: PROHIBITED
Repository visibility: PUBLIC
EG-02: ACTIONABLE_UNPROTECTED
Live main protection: false
P00 exit: BLOCKED ON ISSUE #3 / EG-02 ONLY
P01.01-P01.12 specs: PREPARED / PLANNED
P01 implementation: NOT AUTHORIZED
Kernel code: LOCKED
Business feature code: LOCKED
```
