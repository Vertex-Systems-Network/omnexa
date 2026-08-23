# Omnexa P01 Implementation Entry Gate

Status: **SATISFIED — P01 execution active under strict sequential package control**  
Owner phase: **P01 — Omnexa Kernel**

Foundation Architecture v1 remains frozen. This document preserves the controls that authorize bounded executable kernel work; it does not authorize business-domain implementation.

## Entry controls

### EG-01 — Foundation architecture freeze accepted
State: **SATISFIED**

Evidence: `FOUNDATION_FREEZE_REVIEW.md`, `FOUNDATION_FREEZE.json`, ADR-0010.

### EG-02 — `main` integration protection
State: **SATISFIED**  
Tracker: **Issue #3 — closed/completed**

Verified repository behavior includes:

- `protected: true` for `main`;
- PR-only integration with strict required `governance`;
- intentionally failing governance probe run `32540836431` blocked merge;
- direct-update probe commit `44ca19e80c5fccccebfd8d4f96dde6dc5af14bc2` was rejected;
- force updates remain rejected;
- conversation-resolution probe run `32541439589` blocked merge until the review thread was resolved;
- green integration is permitted only when all required controls are satisfied.

Current single-maintainer policy uses zero required approvals and does not require Code Owner review until an independent reviewer exists. Strict up-to-date branch enforcement remains mandatory; stale green runs do not authorize merge.

### EG-03 — executable verification lane
State: **SATISFIED**  
Tracker: **Issue #14 — closed/completed**

Canonical governance CI is GitHub-hosted only on `ubuntu-24.04`. The job is named `governance` and fails closed unless the runner environment is GitHub-hosted Linux/X64. Local/self-hosted governance runners are prohibited.

### EG-04 — canonical executable verification command
State: **SATISFIED BY COMPLETED P01 PACKAGES**

Completed package regression verifiers remain mandatory in the same GitHub-hosted governance job. Repository Go quality remains enforced by `bash scripts/verify_go_quality.sh` with pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

Canonical completion evidence exists for P01.01 through P01.09 under `docs/roadmap/evidence/`. Latest completed-package implementation evidence is P01.09 PR #59, final strict-up-to-date head `61e9c1115d05300ac9aedf5a555138c6a5a5be1e`, run `32605309150`, job `97109396616`, merge `0bcafbfc52324acba1df9d8eff84a264dda0f233`.

### EG-05 — implementation locks transition atomically
State: **SATISFIED**

Current bounded state after the governed P01.09 closure transition:

- P00/P00.10: `done`;
- P01: `active`;
- P01.01-P01.09: `done`;
- P01.10: `active`;
- P01.11-P01.12: `planned`;
- `kernel_code_authorized=true` only for P01.10;
- `business_feature_code_authorized=false`.

## Active P01 package sequence

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution.

P01.01 through P01.09 are done. P01.09 canonical evidence is `docs/roadmap/evidence/P01.09_COMPLETION_2026-08-22.md`. **P01.10 — Feature flag & configuration registry** is the sole active package. P01.11-P01.12 remain planned.

Advancing again requires P01.10 to reach `done` with required GitHub-hosted G0/G1/G2/G3/G5/G6/G7 evidence, P01.01-P01.09 regressions and repository Go quality all PASS, followed by a governed state reconciliation.

## Active P01.10 bounds

P01.10 may implement only the runtime flag/configuration registry defined in `docs/roadmap/work-packages/P01.10.md`: typed definitions, stable identifiers/descriptions/owner metadata, explicit deterministic defaults/fallbacks, runtime evaluation, future-scope-aware evaluation inputs without P02 identity, version/change metadata hooks, bounded cache-safe refresh/invalidation, explicitly declared operational kill switches and a deterministic test provider.

It may not implement product experimentation/analytics, tenant admin UI, pricing/entitlement/licensing, authorization based solely on flags, business-module flags before their owner modules exist, P01.11 audit transport, P01.12 developer CLI, P02+ identity/tenancy/business behavior or AI/model/agent/planner functionality.

Flags do not grant authority. Security controls fail closed and cannot be disabled through undeclared generic flags. Sensitive configuration remains governed by classification/secrets policy.

## Future browser UI planning requirement

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` remains a planning requirement for future authorized browser UI packages. It does not alter P01 execution authorization.

## External distribution / Issue #4

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. The repository remains public and the current `LICENSE` remains GPLv3. This external distribution gate does not block currently authorized P01 kernel engineering.

## Current decision

```text
Foundation Architecture v1: FROZEN
P00: DONE
Repository visibility: PUBLIC
EG-02 / Issue #3: SATISFIED / CLOSED
EG-03 / Issue #14: SATISFIED / GITHUB-HOSTED ONLY
P01: ACTIVE
P01.01-P01.09: DONE
P01.10: ACTIVE — Feature flag & configuration registry
P01.11-P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```
