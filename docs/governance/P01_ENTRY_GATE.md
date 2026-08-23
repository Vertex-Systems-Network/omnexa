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

Canonical completion evidence exists for P01.01 through P01.11 under `docs/roadmap/evidence/`. Latest completed-package implementation evidence is P01.11 PR #63, final exact head `1c1ab1f8d5120fb6b1e5908fdb93cffef9275940`, run `32610902537`, job `97123708250`, implementation merge `10c94a638b89d47da05f5481fb2db298a2da6942`.

### EG-05 — implementation locks transition atomically
State: **SATISFIED**

Current bounded state proposed by the governed P01.11 closure transition:

- P00/P00.10: `done`;
- P01: `active`;
- P01.01-P01.11: `done`;
- P01.12: `active`;
- P02+: `planned`;
- `kernel_code_authorized=true` only for P01.12;
- `business_feature_code_authorized=false`.

## Active P01 package sequence

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution.

P01.01 through P01.11 are done. P01.11 canonical evidence is `docs/roadmap/evidence/P01.11_COMPLETION_2026-08-23.md`. **P01.12 — Developer CLI Baseline** is the sole active package in this closure transition and owns the final P01 exit proof.

Advancing beyond P01 requires P01.12 to reach `done` with applicable GitHub-hosted G0-G8 evidence, P01.01-P01.11 regressions, repository Go quality and fresh-install P01 exit proof all accurately recorded, followed by a separate governed P01-exit reconciliation. P02 may not activate early.

## Active P01.12 bounds

P01.12 may implement only the developer CLI/P01 exit boundary defined in `docs/roadmap/work-packages/P01.12.md`: stable repository-owned version/help/verify/build-test/approved-diagnostic command surfaces; deterministic fail-closed verification orchestration; explicit exit codes and structured-safe output; safe composition of P01 configuration/migration/diagnostic capabilities; no-secret/no-RESTRICTED output; deterministic clean-checkout/CI behavior; completed-package regression preservation; and the governed fresh-install P01 exit proof.

It may not implement production super-admin authority, P02 tenant/user/role administration, P03 module install/runtime administration, P04+ domain/event/workflow commands, deployment/Kubernetes orchestration, hidden SQL/file mutation, later business modules or AI/model/agent behavior.

CLI convenience does not create authority. Privileged future operations must authenticate/authorize/audit through governed capabilities. Destructive operations require explicit target/environment semantics rather than hidden inference.

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
P01.01-P01.11: DONE
P01.12: ACTIVE — Developer CLI baseline / P01 exit proof
P01 progress: 11 / 12 done
kernel_code_authorized: true
business_feature_code_authorized: false
```
