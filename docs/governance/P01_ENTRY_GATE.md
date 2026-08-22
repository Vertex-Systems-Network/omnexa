# Omnexa P01 Implementation Entry Gate

Status: **SATISFIED — P01 execution active under strict sequential package control**  
Owner phase: **P01 — Omnexa Kernel**

Foundation Architecture v1 remains frozen. This document preserves the evidence that cleared executable kernel entry; it does not authorize business-domain implementation.

## Entry controls

### EG-01 — Foundation architecture freeze accepted
State: **SATISFIED**

Evidence: `FOUNDATION_FREEZE_REVIEW.md`, `FOUNDATION_FREEZE.json`, ADR-0010.

### EG-02 — `main` integration protection
State: **SATISFIED**  
Tracker: **Issue #3 — closed/completed**

Verified 2026-08-22:

- repository visibility is `public`;
- live GitHub branch metadata reports `main protected: true`;
- required check is `governance`;
- PR #34 intentionally failed `governance` in run `32540836431`; GitHub rejected merge with `Required status check "governance" is failing`;
- controlled direct-update probe commit `44ca19e80c5fccccebfd8d4f96dde6dc5af14bc2` was rejected;
- force-update probe was rejected;
- PR #37 run `32541439589` proved conversation resolution enforcement: an unresolved review thread blocked merge until explicit resolution;
- green PR #35 merged normally;
- force pushes and deletion of `main` remain blocked by the configured ruleset.

Current single-maintainer review policy intentionally requires `0` approvals and does not require Code Owner review. Tighten this when an independent reviewer exists.

### EG-03 — executable verification lane
State: **SATISFIED**  
Tracker: **Issue #14 — closed/completed**

Canonical governance CI is GitHub-hosted only:

```yaml
runs-on: ubuntu-24.04
```

The job is named `governance` and fails closed unless `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. Local/self-hosted runners are prohibited.

### EG-04 — canonical executable verification command
State: **SATISFIED BY COMPLETED P01 PACKAGES**

P01.01 established the pinned Go workspace/build verifier. P01.02 added its configuration/security verifier. Completed-package regression verifiers remain mandatory in the same GitHub-hosted governance job.

Canonical evidence:

- `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`;
- `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

### EG-05 — implementation locks transition atomically
State: **SATISFIED**

Current bounded state:

- P00/P00.10: `done`;
- P01: `active`;
- P01.01: `done`;
- P01.02: `done`;
- P01.03: `active`;
- `kernel_code_authorized=true`;
- `business_feature_code_authorized=false`;
- P01.04–P01.12 remain `planned`;
- ADR-0006 remains historical-only.

## Active P01 package sequence

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution.

- P01.01 — **Go workspace / build skeleton**: `done`.
- P01.02 — **Configuration & environment system**: `done`; PR #42, run `32563880800`, job `97009520624`, merge `c857bb9e7df1e347226653eeaded024d6ecd0271`.
- P01.03 — **Structured error & result conventions**: sole `active` package.
- P01.04–P01.12: `planned`.

Advancing again requires P01.03 to reach `done` with required hosted evidence and a governed state reconciliation.

## External distribution / Issue #4

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. The repository remains public and the current `LICENSE` remains GPLv3. This external distribution gate does not block the currently authorized P01 kernel engineering scope.

## Current decision

```text
Foundation Architecture v1: FROZEN
P00: DONE
Repository visibility: PUBLIC
EG-02 / Issue #3: SATISFIED / CLOSED
Live main protection: protected:true
EG-03 / Issue #14: SATISFIED / GITHUB-HOSTED ONLY
P01: ACTIVE
P01.01: DONE
P01.02: DONE
P01.03: ACTIVE — Structured error & result conventions
P01.04-P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```
