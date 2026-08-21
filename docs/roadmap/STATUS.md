# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active — exit verification**
- Current work package: **P00.10 — Foundation architecture freeze review**
- Architecture baseline: **FROZEN — Foundation v1**
- P01 entry: **BLOCKED**
- Kernel implementation: **NOT AUTHORIZED**
- Business-feature implementation: **NOT AUTHORIZED**
- P00 progress: **9 / 10 done; P00.10 verification active**

## P00 packages

| ID | Work package | State |
|---|---|---|
| P00.01 | Repository governance baseline | done |
| P00.02 | Product/domain glossary and naming | done |
| P00.03 | ID/money/time/locale/error conventions | done |
| P00.04 | API contract standard | done |
| P00.05 | Event contract standard | done |
| P00.06 | Security/data classification | done |
| P00.07 | Testing/CI/release standard | done |
| P00.08 | Local developer/repository structure | done |
| P00.09 | Threat model and operational SLO targets | done |
| P00.10 | Foundation architecture freeze review | **active / verification** |

## Freeze result

P00.01–P00.09 are accepted and frozen as **Omnexa Foundation Architecture v1**. The normative review and machine manifest are:

- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/contracts/governance/foundation-freeze.schema.json`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`
- `scripts/validate_freeze_review.py`

Material changes to frozen foundation semantics require change control and a superseding accepted ADR.

## P01 implementation-entry gate

`docs/governance/P01_ENTRY_GATE.md` is authoritative.

### Issue #3 — BLOCKS P01

`main` branch/ruleset protection must be applied and verified before executable P01 merges. Required properties include PR-based integration, blocked force-push/deletion, controlled bypass, conversation resolution and required verification checks once an executable CI lane exists.

### Issue #14 — BLOCKS P01

A working executable verification lane is required before P01. It may be GitHub Actions, another approved provider or a self-hosted lane, but it must execute the repository-owned P00.07 verification semantics in a clean reproducible environment.

The P00 temporary Actions exception is **not valid for P01 executable work**.

### Issue #4 — external-distribution blocker

Licensing/IP/trademark strategy does **not** block private internal P01 engineering after the P01 entry gate is cleared, but it remains a hard gate before public/external distribution, self-hosted customer delivery or public launch.

## Frozen foundation summary

- governance/change control/AI execution;
- canonical vocabulary/domain ownership/dependency rules;
- UUIDv7, exact money, time/locale/error primitives;
- stable HTTP/OpenAPI and event contracts;
- security/data classification/tenant isolation/authorization/audit;
- G0–G8 quality, testing, CI and release semantics;
- governed monorepo/local development/toolchain/config model;
- threat model, operational criticality, SLO/RPO/RTO/error budgets, SEV0–SEV3 and reliability readiness.

## Temporary GitHub Actions exception

ADR-0006 remains active **only because P00 has not exited yet**. Hosted CI is `BLOCKED / NOT RUN`, never PASS. It expires when P00 exits or sooner if executable CI becomes available.

## Exact next transition

P00.10 may become `done` only when Issue #3 and Issue #14 have verified exit evidence. The narrow transition must:

1. mark P00.10 done;
2. mark P00 done;
3. expire ADR-0006;
4. activate P01;
5. set `kernel_code_authorized = true`;
6. keep `business_feature_code_authorized = false`;
7. record branch-protection and executable-CI evidence;
8. define the first P01 kernel work package.

Until then, **do not begin canonical kernel/product implementation**.