# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active — exit verification**
- Current work package: **P00.10 — Foundation architecture freeze review**
- Architecture baseline: **FROZEN — Foundation v1**
- Executable CI entry gate: **SATISFIED — ANY AVAILABLE WINDOWS/X64 SELF-HOSTED RUNNER**
- Latest routing proof: **LOCAL-WIN-02 / run 32535324900 / job 96935023669 / PASS**
- P01 entry: **BLOCKED BY ISSUE #3 / GITHUB PLAN LIMITATION**
- P01.01 preparation: **SPECIFIED / PLANNED / BLOCKED**
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

## P01 implementation-entry gate

### EG-03 / Issue #14 — SATISFIED

The executable CI blocker remains resolved. The canonical workflow now exposes one required job named `governance` on:

```yaml
runs-on: [self-hosted, Windows, X64]
```

GitHub may schedule that job on any available Windows/X64 self-hosted runner. Runner-name pinning, discovery fanout and target evidence aggregation have been removed; the validators themselves remain fail-closed.

Verified current routing evidence:

- workflow run `32535324900`: SUCCESS;
- job `96935023669`: SUCCESS;
- actual runner `LOCAL-WIN-02` / Windows X64;
- machine `ABDUL-HANAN`;
- Git `2.55.0.windows.5`;
- Python `3.13.7`;
- branch-protection PowerShell parse PASS;
- governance validator PASS;
- development specification validator PASS;
- threat/operations validator PASS;
- foundation-freeze validator PASS;
- P01-preparation validator PASS.

The earlier LOCAL-WIN-4 evidence from PR #23/run `32528329184` remains historical provenance only. It is no longer a scheduling requirement.

GitHub-hosted standard runners may be introduced later if capacity/policy permits without weakening P00.07 quality semantics. The current self-hosted pool is already operational and avoids dependence on hosted-minute entitlement.

### EG-02 / Issue #3 — BLOCKED BY CURRENT GITHUB PLAN

The owner executed the merged branch-protection administration tooling from an authenticated admin workstation. GitHub returned:

```text
gh: Upgrade to GitHub Pro or make this repository public to enable this feature. (HTTP 403)
```

Verified context:

- repository is organization-owned and `private`;
- authenticated user has repository `admin` permission;
- branch-protection scripts parse successfully;
- the API attempt reached GitHub and was rejected by product-plan entitlement;
- `main` remains `protected: false`.

Do not retry the same API operation until one of these changes:

1. upgrade to a plan that supports private branch protection/rulesets;
2. intentionally change repository visibility to public through a separate owner/legal/security decision; or
3. approve a superseding governance ADR that deliberately replaces EG-02 with a compensating control.

Changing this repository to public merely to clear EG-02 is not an automatic workaround because Issue #4 remains an external distribution/IP/licensing gate.

Until EG-02 is satisfied or deliberately superseded, P00.10 remains in exit verification and kernel code stays locked.

## P01.01 implementation-readiness preparation

P01 itself remains blocked, but its first executable package is pre-specified so no architecture decisions need to be invented after the gate clears.

Prepared artifacts:

- `docs/roadmap/work-packages/P01.01.md` — controlling specification for **Go workspace/build skeleton**;
- `docs/governance/P00_P01_TRANSITION_CHECKLIST.md` — atomic P00→P01 state/lock/evidence transition;
- `scripts/validate_p01_preparation.py` — fail-closed validation that P01 is prepared but still locked.

P01.01 remains `planned`. Its scope is deliberately limited to the Go workspace/build/process skeleton, initial build/version metadata and applicable G0/G1/G2/G7 evidence. Configuration, DB, cache, storage, telemetry, health, jobs, flags, audit, identity/tenancy, module runtime and business domains remain later packages/phases.

The preparation validator rejects known P01.01 executable paths such as `go.mod`, `go.work` and `kernel/cmd/omnexa/main.go` while `kernel_code_authorized=false` and verifies the runner-name-agnostic governance workflow shape.

## Issue #4 — external-distribution blocker

Licensing/IP/trademark strategy does **not** block private internal P01 engineering after the P01 entry gate is cleared, but remains a hard gate before public/external distribution, self-hosted customer delivery or public launch.

The owner/legal decision is structured in:

- `docs/governance/LICENSING_DECISION.md` — canonical gate;
- `docs/governance/LICENSING_DECISION_BRIEF.md` — explicit owner choices for distribution model, core licensing direction, contributor IP, marketplace/extension boundary, dependency policy and trademark/name clearance.

No repository `LICENSE` change is authorized by preparation alone.

## Exact next transition

P00.10 may become `done` only when Issue #3 has verified branch-protection evidence **or** an explicitly accepted superseding governance ADR replaces EG-02 with a compensating control. The narrow transition must follow `P00_P01_TRANSITION_CHECKLIST.md` and:

1. mark P00.10 done;
2. mark P00 done;
3. retire ADR-0006 from active use;
4. activate P01 and P01.01;
5. set `kernel_code_authorized = true`;
6. keep `business_feature_code_authorized = false`;
7. record the applicable protection/compensating-control evidence;
8. preserve P01.02–P01.12 as planned.

Until then, **do not begin canonical kernel/product implementation**.
