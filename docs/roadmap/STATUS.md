# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active**
- Current work package: **P00.07 — Testing, CI and release standard**
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
| P00.06 | Security and data-classification baseline | done | security/data-classification/control matrix, classification schema + ADR-0005; hosted CI evidence temporarily BLOCKED under ADR-0006 |
| P00.07 | Testing, CI and release standard | active | Current canonical work package |
| P00.08 | Local developer and repository structure specification | planned | Requires P00.07 |
| P00.09 | Initial threat model and operational SLO targets | planned | Requires P00.06 + P00.08 |
| P00.10 | Foundation architecture freeze review | planned | Final P00 exit gate |

P00 package progress: **6 / 10 done**.

## Frozen foundation contracts

### P00.03 — Primitive semantics

- canonical entity/request/event/workflow/job/audit IDs: **UUIDv7**;
- PostgreSQL canonical identifiers use native `uuid`; tenant-owned state uses `tenant_id` where applicable;
- money uses exact decimal semantics with explicit ISO currency and no binary floating-point representation;
- absolute instants use UTC/`timestamptz`; civil-time recurrence retains IANA timezone semantics;
- locale uses BCP 47; locale/country/currency/timezone are independent concepts; RTL is first-class;
- public errors use stable machine codes and safe structured problem details.

### P00.04 — HTTP API baseline

- stable routes: `/api/v{major}/{domain}/{resources}`;
- OpenAPI **3.2.0** is the canonical stable HTTP contract description baseline;
- JSON contract fields use lowercase `snake_case`;
- errors use `application/problem+json` with Omnexa machine codes/request correlation;
- scalable collections use opaque cursor pagination;
- protected retriable mutations define `Idempotency-Key` semantics;
- lost-update-sensitive mutations use explicit optimistic concurrency where applicable;
- tenant/organization identifiers from clients are never authorization authority.

### P00.05 — Event baseline

- event type: `<domain>.<subject>.<past_tense_fact>.v<major>`;
- canonical envelope is CloudEvents-compatible structured JSON with UUIDv7 event identity;
- producer owns event meaning/schema and derives trusted tenant context;
- delivery baseline is **at least once** and consumers are idempotent;
- business-significant publication/consumption uses outbox + inbox/deduplication or equivalent guarantees;
- no global ordering assumption; retry is bounded; poison events go to governed dead-letter/quarantine;
- replay preserves original identity/payload and may not duplicate protected business side effects.

### P00.06 — Security and data classification baseline

Canonical confidentiality classes:

```text
PUBLIC
INTERNAL
CONFIDENTIAL
RESTRICTED
```

Security invariants now include:

- explicit trust boundaries and zero implicit trust;
- distinct human, service, workload, device, integration, support and AI principals;
- authentication independent from authorization;
- authorization = RBAC + relationships + contextual policy + governed capabilities;
- tenant and organization isolation across OLTP, cache, files, search, analytics, events, backups and AI/vector data;
- secrets/key material are `RESTRICTED`, excluded from source control/logging/AI and managed by approved secret/KMS mechanisms;
- production sensitive data does not flow to lower environments by default;
- audit is a protected security/business record separate from debug logs;
- privileged operations/support impersonation/export/purge/event replay/high-impact AI actions require explicit capability, policy and audit;
- external webhooks/integrations are independent trust/disclosure boundaries;
- extension/module permissions, provenance and supply-chain controls are required;
- retention, deletion, search, export, analytics and AI behavior is classification-aware;
- `RESTRICTED` data is prohibited from generic logs/search/analytics and is prohibited as AI input by default.

Normative evidence:

- `docs/security/SECURITY_STANDARD.md`
- `docs/security/DATA_CLASSIFICATION.md`
- `docs/security/SECURITY_CONTROL_MATRIX.md`
- `docs/contracts/security/data-classification.schema.json`
- `docs/adr/ADR-0005-security-data-classification-baseline.md`

## Temporary GitHub Actions exception

GitHub Actions allowance is currently exhausted/disabled. The project owner explicitly authorized continuing P00 documentation/specification work without hosted Actions until the quota condition is resolved.

The temporary policy is defined by:

- `docs/governance/CI_EVIDENCE_EXCEPTION_2026-08-22.md`
- `docs/adr/ADR-0006-temporary-p00-ci-evidence-exception.md`
- issue #14 for the runner/quota blocker evidence.

Hosted Actions evidence is therefore **BLOCKED / NOT RUN**, never recorded as PASS. P00 packages may use manual repository/state/evidence verification only while they remain documentation/specification work. This exception expires before any P01 implementation merge, at P00 exit, or sooner if Actions returns.

## Governance hardening status

File-level governance is active through CODEOWNERS, contributor/security policies, issue/ADR templates, `scripts/validate_governance.py` and the `Omnexa Governance` workflow definition.

Hosted/business decisions still tracked:

1. **Issue #3 — main branch/ruleset protection:** GitHub still reports `main` as unprotected. The connected GitHub toolset does not expose a branch-protection/ruleset write mutation, so hosted admin configuration remains required and must be verified against `docs/governance/REPOSITORY_HARDENING.md`.
2. **Issue #14 — GitHub Actions quota/runner:** hosted execution is temporarily unavailable and handled only through ADR-0006 for P00 specification work.
3. **Issue #4 — licensing/IP/trademark:** existing GPLv3 is not automatically approved as the final commercial strategy; explicit owner/legal decision is required before external distribution/public launch.

None of these items authorize early kernel/business implementation.

## Phase states

P00 is the only active phase. P01 through P27 remain `planned` until their governed entry gates are satisfied.

## Execution lock

Until P00 is complete, contributors and AI systems may work only on architecture/governance/specification activities belonging to P00 plus narrow repository-maintenance fixes.

Do not begin kernel implementation, database models, CRM, ERP, commerce, POS, website builder, payments or AI product code before the P00 exit gate.

## How status changes

A package moves to `done` only when acceptance evidence satisfies `docs/governance/DEFINITION_OF_DONE.md` and `docs/roadmap/STATE.json` is reconciled in the same governed change. While ADR-0006 is active, hosted Actions may be `BLOCKED`/`NOT RUN` for P00 documentation/specification changes only; manual evidence must be recorded and the exception may not be used for P01+ executable work.