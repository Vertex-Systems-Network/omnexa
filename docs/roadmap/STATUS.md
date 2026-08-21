# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active**
- Current work package: **P00.08 — Local developer and repository structure specification**
- Business-feature implementation: **NOT AUTHORIZED YET**
- Kernel implementation: **NOT AUTHORIZED YET**

## P00 work packages

| ID | Work package | State | Evidence / note |
|---|---|---|---|
| P00.01 | Repository governance baseline | done | Product Constitution, architecture/module standards, AI policy, roadmap/state, change control and ADR-0001 |
| P00.02 | Product/domain glossary and naming standard | done | glossary, naming, ownership/dependency model, repository governance and baseline CI |
| P00.03 | ID, money, time, locale and error conventions | done | primitive standards + ADR-0002 |
| P00.04 | API contract standard | done | `API_STANDARD.md`, OpenAPI template + ADR-0003 |
| P00.05 | Event contract standard | done | `EVENT_STANDARD.md`, event envelope schema + ADR-0004 |
| P00.06 | Security and data-classification baseline | done | security/data-classification/control matrix, classification schema + ADR-0005; manual P00 evidence under ADR-0006 because hosted Actions is BLOCKED |
| P00.07 | Testing, CI and release standard | done | testing/CI/release standards, G0-G8 gate matrix, quality evidence schema + ADR-0007; hosted execution BLOCKED under ADR-0006 |
| P00.08 | Local developer and repository structure specification | active | Current canonical work package |
| P00.09 | Initial threat model and operational SLO targets | planned | Requires P00.06 + P00.08 |
| P00.10 | Foundation architecture freeze review | planned | Final P00 exit gate |

P00 package progress: **7 / 10 done**.

## Frozen quality baseline — P00.07

Omnexa quality semantics are **repository-owned and CI-provider independent**. GitHub Actions, another CI service and local development must execute the same canonical verification semantics.

Canonical test layers include static/structural, unit, component, contract, integration, migration, security/negative, module lifecycle, end-to-end, performance/resilience, compatibility and disaster-recovery/rehearsal testing as applicable.

Important invariants:

- affected tenant-owned paths require same-tenant success plus cross-tenant denial evidence;
- authorization changes require allow/deny/privilege-escalation evidence;
- retriable/event/async paths require duplicate/idempotency/replay evidence;
- persistence changes require fresh-install plus supported-upgrade migration evidence;
- optional modules require disable/re-enable/degradation isolation evidence;
- money/time/localization/security boundaries require domain-specific negative/boundary tests;
- flaky tests are defects and may be quarantined only through named, expiring governance;
- aggregate coverage percentage never substitutes for critical invariant coverage.

Quality gate classes are:

```text
G0 Governance
G1 Static
G2 Unit / Component
G3 Contract / Integration
G4 Data / Migration
G5 Security / Tenancy
G6 Lifecycle / Resilience
G7 Build / Package
G8 Supply Chain / Release
```

Allowed evidence states are exactly:

```text
PASS
FAIL
BLOCKED
NOT RUN
N/A
```

`BLOCKED`, `NOT RUN` and `N/A` are never silently converted into PASS.

Release architecture uses semantic-versioning semantics, immutable source/artifact identity and a **build once, promote** model where artifact type permits. Stable releases eventually require applicable test, migration, security, packaging, SBOM/provenance/signature and operational-readiness evidence.

Normative evidence:

- `docs/quality/TESTING_STANDARD.md`
- `docs/quality/CI_STANDARD.md`
- `docs/quality/RELEASE_STANDARD.md`
- `docs/quality/QUALITY_GATE_MATRIX.md`
- `docs/contracts/quality/quality-gates.schema.json`
- `docs/adr/ADR-0007-testing-ci-release-baseline.md`

## Previously frozen baselines

P00.03 freezes identifier/money/time/locale/error semantics. P00.04 freezes HTTP/OpenAPI semantics. P00.05 freezes event/reliability/replay semantics. P00.06 freezes security and data-classification semantics. These remain mandatory and P00.07 defines how future implementation proves them.

## Temporary GitHub Actions exception

GitHub Actions allowance is currently exhausted/disabled. The project owner explicitly authorized continuing P00 documentation/specification work without hosted Actions until the quota condition is resolved.

The temporary policy is defined by:

- `docs/governance/CI_EVIDENCE_EXCEPTION_2026-08-22.md`
- `docs/adr/ADR-0006-temporary-p00-ci-evidence-exception.md`
- issue #14 for runner/quota blocker evidence.

Hosted Actions evidence is **BLOCKED / NOT RUN**, never recorded as PASS. P00.07 was manually reviewed for scope, mandatory evidence, state/dependency ordering and quality-contract reconciliation. The exception expires before any P01 implementation merge, at P00 exit, or sooner if Actions returns.

## Governance hardening status

1. **Issue #3 — main branch/ruleset protection:** `main` remains unprotected; hosted admin configuration is still required.
2. **Issue #14 — GitHub Actions quota/runner:** hosted execution is temporarily unavailable and handled only through ADR-0006 for P00 specification work.
3. **Issue #4 — licensing/IP/trademark:** existing GPLv3 is not automatically approved as the final commercial strategy; explicit owner/legal decision is required before external distribution/public launch.

None of these authorize early kernel/business implementation.

## Phase states

P00 is the only active phase. P01 through P27 remain `planned` until their governed entry gates are satisfied.

## Execution lock

Until P00 is complete, contributors and AI systems may work only on architecture/governance/specification activities belonging to P00 plus narrow repository-maintenance fixes.

Do not begin kernel implementation, database models, CRM, ERP, commerce, POS, website builder, payments or AI product code before the P00 exit gate.

## How status changes

A package moves to `done` only when its required evidence is recorded and `STATE.json` is reconciled. While ADR-0006 is active, hosted Actions may be `BLOCKED`/`NOT RUN` for P00 documentation/specification changes only; this exception cannot be used for P01+ executable work.