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

Canonical completion evidence exists for P01.01 through P01.10 under `docs/roadmap/evidence/`. Latest completed-package implementation evidence is P01.10 PR #61, final exact head `4c9914e4641d0d6e94a895d0fcd16c3a6bf4d962`, run `32609018028`, job `97118796940`, implementation merge `9d11b9250eb74ca2ade531ee58e8f905468cf103`.

### EG-05 — implementation locks transition atomically
State: **SATISFIED**

Current bounded state proposed by the governed P01.10 closure transition:

- P00/P00.10: `done`;
- P01: `active`;
- P01.01-P01.10: `done`;
- P01.11: `active`;
- P01.12: `planned`;
- `kernel_code_authorized=true` only for P01.11;
- `business_feature_code_authorized=false`.

## Active P01 package sequence

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution.

P01.01 through P01.10 are done. P01.10 canonical evidence is `docs/roadmap/evidence/P01.10_COMPLETION_2026-08-23.md`. **P01.11 — Audit Transport Foundation** is the sole next active package in this closure transition. P01.12 remains planned.

Advancing again requires P01.11 to reach `done` with required GitHub-hosted G0/G1/G2/G3/G5/G6/G7 evidence, P01.01-P01.10 regressions and repository Go quality all PASS, followed by a separate governed state reconciliation.

## Active P01.11 bounds

P01.11 may implement only the audit transport foundation defined in `docs/roadmap/work-packages/P01.11.md`: stable classification-aware audit envelope; actor/action/target/scope/outcome/correlation/reason/approval metadata without P02 identity; append-oriented sink interface; explicit required-audit failure semantics; classification/redaction enforcement; UUIDv7/timestamp immutability conventions; impersonation/privileged-action metadata; deterministic local/test sink; and bounded protected-payload-safe transport health observability.

It may not implement P02 identity/tenant/role catalogs, business-domain audit events, compliance export/reporting UI, long-term legal hold/retention, domain-event replacement, P01.12 CLI, durable messaging/outbox/inbox pull-forward, later business modules or AI/model/agent behavior.

Audit write access does not imply audit read/export authority. Audit metadata is descriptive and does not create authorization, tenancy or identity authority during P01.11. Protected audit remains separate from ordinary logs, and required-audit sink failure must not silently claim success.

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
P01.01-P01.10: DONE
P01.11: ACTIVE — Audit transport foundation
P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```
