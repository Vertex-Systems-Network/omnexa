# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active — exit verification**
- Current work package: **P00.10 — Foundation architecture freeze review**
- Architecture baseline: **FROZEN — Foundation v1**
- Executable CI entry gate: **SATISFIED ON LOCAL-WIN-4**
- P01 entry: **BLOCKED BY ISSUE #3 / GITHUB PLAN LIMITATION**
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

The executable CI blocker remains resolved and the canonical lane is certified on the specifically requested organization runner `LOCAL-WIN-4`.

Verified current evidence:

- PR #23 migrated the canonical workflow to fail closed unless LOCAL-WIN-4 produces validation evidence;
- runner: `LOCAL-WIN-4` / Windows X64;
- machine/work root: `ABDUL-HANAN` / `C:\actions-runner-4\_work`;
- Git `2.55.0.windows.5`;
- Python `3.13.7`;
- workflow run `32528329184`;
- LOCAL-WIN-4 target job `96915072868`: SUCCESS;
- governance validator PASS;
- development specification validator PASS;
- threat/operations validator PASS;
- foundation-freeze validator PASS;
- final required job named `governance`: SUCCESS;
- PR #23 merged to `main` as `1a14362e2ed52a20d66cec6f28b93a2ee457f9a9`;
- Issue #14 remains closed/completed.

Because GitHub Actions does not schedule by runner name and LOCAL-WIN-4 currently has no unique Actions label, the canonical workflow fans out only across local Windows/X64 self-hosted runners, executes protected validators only when `RUNNER_NAME == LOCAL-WIN-4`, uploads target pass evidence only from that runner, and fails the final `governance` job when target evidence is absent.

GitHub-hosted runner quota is not required. ADR-0006 is historical evidence only while this executable lane remains operational.

### EG-02 / Issue #3 — BLOCKED BY CURRENT GITHUB PLAN

The owner executed the merged branch-protection administration tooling from an authenticated admin workstation. GitHub returned:

```text
gh: Upgrade to GitHub Pro or make this repository public to enable this feature. (HTTP 403)
```

Verified context:

- repository is organization-owned and `private`;
- authenticated user has repository `admin` permission;
- branch-protection scripts parse successfully on `LOCAL-WIN-4`;
- the API attempt reached GitHub and was rejected by product-plan entitlement;
- `main` remains `protected: false`.

GitHub's current feature model provides protected branches/rulesets for public repositories on GitHub Free, while private repositories require GitHub Pro/Team/Enterprise-class support. The remaining blocker is therefore the hosted GitHub plan, not CI, permissions or script correctness.

Do not retry the same API operation until one of these changes:

1. upgrade to a plan that supports private branch protection/rulesets;
2. intentionally change repository visibility to public through a separate owner/legal/security decision; or
3. approve a superseding governance ADR that deliberately replaces EG-02 with a compensating control.

Changing this repository to public merely to clear EG-02 is not an automatic workaround because Issue #4 remains an external distribution/IP/licensing gate.

Until EG-02 is satisfied or deliberately superseded, P00.10 remains in exit verification and kernel code stays locked.

### Issue #4 — external-distribution blocker

Licensing/IP/trademark strategy does **not** block private internal P01 engineering after the P01 entry gate is cleared, but remains a hard gate before public/external distribution, self-hosted customer delivery or public launch.

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

P00.10 may become `done` only when Issue #3 has verified branch-protection evidence **or** an explicitly accepted superseding governance ADR replaces EG-02 with a compensating control. The narrow transition must:

1. mark P00.10 done;
2. mark P00 done;
3. retire ADR-0006 from active use;
4. activate P01;
5. set `kernel_code_authorized = true`;
6. keep `business_feature_code_authorized = false`;
7. record the applicable protection/compensating-control evidence;
8. define the first P01 kernel work package.

Until then, **do not begin canonical kernel/product implementation**.
