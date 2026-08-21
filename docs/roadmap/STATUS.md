# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active — exit verification**
- Current work package: **P00.10 — Foundation architecture freeze review**
- Architecture baseline: **FROZEN — Foundation v1**
- Executable CI entry gate: **SATISFIED**
- P01 entry: **BLOCKED BY ISSUE #3 ONLY**
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

P00.01–P00.09 are accepted and frozen as **Omnexa Foundation Architecture v1**. Material changes to frozen foundation semantics require change control and a superseding accepted ADR.

Normative freeze/entry records:

- `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/contracts/governance/foundation-freeze.schema.json`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`
- `scripts/validate_freeze_review.py`

## P01 implementation-entry gate

### EG-03 / Issue #14 — SATISFIED

The executable CI blocker is resolved. PR #20 moved the canonical governance workflow to the local Windows self-hosted selector:

```text
[self-hosted, Windows, X64]
```

Verified evidence:

- runner: `LOCAL-WIN-01`;
- workflow run: `32522919774`;
- governance job: `96898839560`;
- Git `2.55.0.windows.5`;
- Python `3.13.7`;
- governance validator PASS;
- development specification validator PASS;
- threat/operations validator PASS;
- foundation-freeze validator PASS;
- PR #20 merged to `main` as `c2ab2cd679c295a8dec84b1879acb9a9e02ad67d`;
- Issue #14 closed/completed.

GitHub-hosted runner quota is no longer required for the canonical governance lane. ADR-0006 is retained as historical evidence but is not the active evidence path while the self-hosted lane remains operational.

### EG-02 / Issue #3 — REMAINING P01 BLOCKER

`main` branch/ruleset protection must still be applied and verified before executable P01 merges. Required properties include:

- PR-based integration;
- required `governance` check;
- force-push blocked;
- branch deletion blocked;
- conversation-resolution policy;
- controlled break-glass bypass only.

Until Issue #3 is verified, P00.10 remains in exit verification and kernel code stays locked.

### Issue #4 — external-distribution blocker

Licensing/IP/trademark strategy does **not** block private internal P01 engineering after Issue #3 is cleared, but remains a hard gate before public/external distribution, self-hosted customer delivery or public launch.

## Frozen foundation summary

- governance/change control/AI execution;
- canonical vocabulary/domain ownership/dependency rules;
- UUIDv7, exact money, time/locale/error primitives;
- stable HTTP/OpenAPI and event contracts;
- security/data classification/tenant isolation/authorization/audit;
- G0–G8 quality, testing, CI and release semantics;
- governed monorepo/local development/toolchain/config model;
- threat model, operational criticality, SLO/RPO/RTO/error budgets, SEV0–SEV3 and reliability readiness.

## Exact next transition

P00.10 may become `done` only when Issue #3 has verified branch-protection evidence. The narrow transition must:

1. mark P00.10 done;
2. mark P00 done;
3. retire ADR-0006 from active use;
4. activate P01;
5. set `kernel_code_authorized = true`;
6. keep `business_feature_code_authorized = false`;
7. record branch-protection evidence;
8. define the first P01 kernel work package.

Until then, **do not begin canonical kernel/product implementation**.
