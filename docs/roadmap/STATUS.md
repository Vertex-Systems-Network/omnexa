# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active — exit verification**
- Current work package: **P00.10 — Foundation architecture freeze review**
- Architecture baseline: **FROZEN — Foundation v1**
- Repository visibility: **PUBLIC**
- Executable CI entry gate: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Hosted proof: **run 32537207455 / job 96940269306 / PASS**
- Local/self-hosted governance runners: **PROHIBITED**
- P01 entry: **BLOCKED BY ISSUE #3 — LIVE MAIN STILL UNPROTECTED**
- P01.01–P01.12 preparation: **12 / 12 SPECIFIED / PLANNED / BLOCKED**
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
- `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`
- `docs/contracts/governance/foundation-freeze.schema.json`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`
- `scripts/validate_freeze_review.py`
- `scripts/validate_p01_preparation.py`
- `scripts/validate_p01_package_specs.py`

## P01 implementation-entry gate

### EG-03 / Issue #14 — SATISFIED

Canonical governance CI now runs only on GitHub-hosted infrastructure:

```yaml
runs-on: ubuntu-24.04
```

The workflow fails closed unless `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. It has no self-hosted/local fallback.

Verified evidence:

- workflow run `32537207455`: SUCCESS;
- job `96940269306`: SUCCESS;
- runner `GitHub Actions 1000006777`;
- environment `github-hosted`;
- Ubuntu 24.04.4 LTS / X64;
- runner image `ubuntu-24.04`, version `20260816.277.1`;
- Git `2.55.0`;
- Python `3.12.3`;
- PowerShell `7.6.5`;
- branch-protection PowerShell syntax parse PASS;
- governance validator PASS;
- development specification validator PASS;
- threat/operations validator PASS;
- foundation-freeze validator PASS on the initial hosted migration head;
- P01-preparation validator PASS;
- P01 package-specification validator PASS.

Earlier local/self-hosted runs remain historical provenance only. They are not an active routing option.

### EG-02 / Issue #3 — ACTIONABLE_UNPROTECTED

The repository is now public. The earlier private-repository HTTP 403/plan limitation is therefore historical, not the current blocker.

Verified current context:

- repository is organization-owned and `public`;
- authenticated user has repository `admin` permission;
- branch-protection scripts parse successfully on GitHub-hosted Ubuntu;
- live branch API still reports `main protected: false`.

The next operation is to apply and verify the hosted protection policy. Until live evidence reports `protected:true` with strict required `governance`, PR-only integration, conversation resolution, force-push/deletion blocking and administrator enforcement, P00.10 stays in exit verification and kernel code stays locked.

A manual GitHub-hosted administration workflow is prepared at `.github/workflows/main-protection-admin.yml`; it requires a short-lived repository Administration secret `OMNEXA_GITHUB_ADMIN_TOKEN`. The ordinary `GITHUB_TOKEN` is intentionally insufficient.

## P01 implementation-readiness preparation

All twelve P01 work-package specifications are prepared at `docs/roadmap/work-packages/P01.01.md` through `P01.12.md`, with machine sequence `P01_PACKAGE_SEQUENCE.json` enforcing strict sequential one-active-package activation.

P01.01 remains the next package and remains `planned`. The preparation validators reject premature executable paths while `kernel_code_authorized=false`.

## Issue #4 — public visibility / licensing-IP-trademark

The repository is now public while the current repository `LICENSE` remains GPLv3 and trademark clearance remains unresolved. Issue #4 requires explicit owner/legal reconciliation for product launch, commercial licensing, contribution policy and trademark claims.

No repository `LICENSE` change or trademark claim is authorized by this CI/governance migration.

## Exact next transition

P00.10 may become `done` only when Issue #3 has verified live branch-protection evidence **or** an explicitly accepted superseding governance ADR replaces EG-02 with compensating controls. The narrow transition must follow `P00_P01_TRANSITION_CHECKLIST.md` and:

1. mark P00.10 done;
2. mark P00 done;
3. retire ADR-0006 from active use;
4. activate P01 and P01.01;
5. set `kernel_code_authorized = true`;
6. keep `business_feature_code_authorized = false`;
7. record the applicable protection/compensating-control and GitHub-hosted CI evidence;
8. preserve P01.02–P01.12 as planned.

Until then, **do not begin canonical kernel/product implementation**.
