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

Verified repository behavior includes `main protected:true`, PR-only integration, strict required `governance`, rejected direct/force updates, failed-check merge rejection and required conversation resolution. Current single-maintainer policy uses zero required approvals and does not require Code Owner review until an independent reviewer exists.

P01.06 integration also reconfirmed the strict up-to-date requirement: PR #50 could not merge while its branch lagged protected `main`, despite previously green workflow runs. The branch was synchronized with current `main`, fresh run `32588244996` passed, and the PR then merged normally. Protection was not bypassed or weakened.

### EG-03 — executable verification lane
State: **SATISFIED**  
Tracker: **Issue #14 — closed/completed**

Canonical governance CI is GitHub-hosted only on `ubuntu-24.04`. The job is named `governance` and fails closed unless the runner environment is GitHub-hosted Linux/X64. Local/self-hosted governance runners are prohibited.

### EG-04 — canonical executable verification command
State: **SATISFIED BY COMPLETED P01 PACKAGES**

Completed package regression verifiers remain mandatory in the same GitHub-hosted governance job. Repository Go quality remains enforced by `bash scripts/verify_go_quality.sh` with pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

Canonical completed-package evidence:

- `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`;
- `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`;
- `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`;
- `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`;
- `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`;
- `docs/roadmap/evidence/P01.06_COMPLETION_2026-08-22.md`.

Latest completed-package implementation evidence is P01.06 PR #50, final strict-up-to-date run `32588244996`, job `97067784835`, merge `f7867d9e1c570e3abbed90740970acf7b5a30bd7`.

### EG-05 — implementation locks transition atomically
State: **SATISFIED**

Current bounded state after the governed P01.06 closure transition:

- P00/P00.10: `done`;
- P01: `active`;
- P01.01-P01.06: `done`;
- P01.07: `active`;
- P01.08-P01.12: `planned`;
- `kernel_code_authorized=true` only for P01.07;
- `business_feature_code_authorized=false`.

## Active P01 package sequence

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution.

- P01.01 — **Go workspace/build skeleton**: done.
- P01.02 — **Configuration & environment system**: done.
- P01.03 — **Structured error & result conventions**: done.
- P01.04 — **PostgreSQL connection & migration foundation**: done.
- P01.05 — **Cache abstraction**: done.
- P01.06 — **Object & file storage abstraction**: done; PR #50 / run `32588244996` / job `97067784835` / merge `f7867d9e1c570e3abbed90740970acf7b5a30bd7`.
- P01.07 — **Structured logging & OpenTelemetry baseline**: sole active package.
- P01.08-P01.12: planned.

Advancing again requires P01.07 to reach `done` with its required GitHub-hosted evidence, P01.01-P01.06 regressions and repository Go quality all PASS, followed by a governed state reconciliation.

## Active P01.07 bounds

P01.07 may implement only the structured logging/OpenTelemetry-compatible observability baseline defined in `docs/roadmap/work-packages/P01.07.md`: stable structured fields, levels/defaults, correlation/trace context helpers, OpenTelemetry resource/provider lifecycle, vendor-neutral exporter configuration, bounded shutdown/flush, redaction/filtering hooks and deterministic test capture.

It may not implement product analytics/business metrics, dashboards/alerts/SLO automation, domain-specific telemetry, audit semantics, health/readiness, scheduler/feature-registry/audit transport, identity/tenancy/module/event/workflow/business behavior or AI/model/agent/planner functionality.

Secrets, tokens, private keys and raw `RESTRICTED` content are prohibited from telemetry by default. Correlation IDs are not authorization credentials. Observability outage must not become a correctness/security bypass.

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
P01.01-P01.06: DONE
P01.07: ACTIVE — Structured logging & OpenTelemetry baseline
P01.08-P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```
