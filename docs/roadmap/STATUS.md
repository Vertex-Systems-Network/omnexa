# Omnexa Program Status

Last reconciled: **2026-08-21**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active**
- Current work package: **P00.05 — Event contract standard**
- Business-feature implementation: **NOT AUTHORIZED YET**
- Kernel implementation: **NOT AUTHORIZED YET**

## P00 work packages

| ID | Work package | State | Evidence / note |
|---|---|---|---|
| P00.01 | Repository governance baseline | done | `AGENTS.md`, Product Constitution, Architecture, Module Standard, Change Control, DoD, Master Plan, State and baseline ADR |
| P00.02 | Product/domain glossary and naming standard | done | `GLOSSARY.md`, `NAMING_STANDARD.md`, `DOMAIN_OWNERSHIP.md`, `DEPENDENCY_MATRIX.md`, contribution/security/hardening controls and governance CI baseline |
| P00.03 | ID, money, time, locale and error conventions | done | `IDENTIFIER_STANDARD.md`, `MONEY_STANDARD.md`, `TIME_STANDARD.md`, `LOCALE_STANDARD.md`, `ERROR_STANDARD.md`, ADR-0002 plus governance validation |
| P00.04 | API contract standard | done | `API_STANDARD.md`, OpenAPI 3.2 foundation template, ADR-0003 and governance validation |
| P00.05 | Event contract standard | active | Current canonical work package: event envelope/versioning/ownership/reliability/replay rules |
| P00.06 | Security and data-classification baseline | ready | Dependency gate satisfied; remains non-active while P00.05 executes |
| P00.07 | Testing/CI/release standard | planned | Requires P00.04/P00.05/P00.06 |
| P00.08 | Local developer and repository structure specification | planned | Requires P00.03/P00.07 |
| P00.09 | Initial threat model and operational SLO targets | planned | Must precede foundation freeze |
| P00.10 | Foundation architecture freeze review | planned | Final P00 exit gate |

P00 package progress: **4 / 10 done**.

## P00.03 frozen primitives

- Canonical new entity/request/event/workflow/job/audit identifiers: **UUIDv7**; PostgreSQL native `uuid`; canonical tenant scope field: `tenant_id`.
- Money: exact decimal + explicit currency; JSON decimal strings; PostgreSQL baseline `NUMERIC(38,18)`; no binary floating point; explicit rounding/conversion policies.
- Time: UTC instants + `timestamptz`; IANA timezones; business dates remain dates; recurring civil-time schedules retain timezone semantics.
- Locale: BCP 47 language tags; country codes separate; locale never implies currency/timezone; RTL is a first-class UI requirement.
- Errors: stable machine codes + safe structured problem representation + request/trace correlation; no public stack traces/SQL/secrets.

## P00.04 frozen HTTP API baseline

- Stable route major versioning: `/api/v{major}/{domain}/{resources}`.
- OpenAPI **3.2.0** is the canonical HTTP contract description baseline.
- JSON contract fields use lowercase `snake_case`; success bodies use `data` with controlled metadata/page envelopes.
- Business lifecycle commands use explicit action operations rather than arbitrary `status` mutations.
- Errors use `application/problem+json` with canonical Omnexa machine codes/request IDs.
- Protected retriable mutations define `Idempotency-Key` semantics; lost-update-sensitive resources may use `ETag` + `If-Match`.
- Cursor pagination (`page_size`, `page_cursor`) is the scalable default; filters/sorts are allowlisted.
- Tenant selection is untrusted input and never authorization authority.
- Compatibility/deprecation rules and an OpenAPI machine template are now frozen.

These conventions are platform contracts. P00.05/P00.06/P00.07 may add event/security/CI details but must not contradict them without change control/ADR.

## Governance hardening status

File-level governance is specified/implemented through:

- `.github/CODEOWNERS`;
- contribution and security policies;
- architecture-change and bug templates;
- ADR template;
- repository hardening specification;
- dependency-free governance state validator;
- GitHub Actions governance workflow.

Two hosted/business decisions remain explicitly tracked:

1. **Issue #3 — main branch/ruleset protection:** GitHub-hosted branch protection must be enabled and verified against `docs/governance/REPOSITORY_HARDENING.md`.
2. **Issue #4 — licensing/IP/trademark:** existing GPLv3 must not be treated as an automatically approved final commercial strategy; owner/legal decision is required before external distribution/public launch.

Neither issue authorizes early kernel/business implementation.

## Phase states

| Phase | State |
|---|---|
| P00 Product Constitution & Architecture Freeze | active |
| P01 Omnexa Kernel | planned |
| P02 Identity, Tenancy & Organization | planned |
| P03 Module Runtime | planned |
| P04 Data, Jobs & Event Fabric | planned |
| P05 Omnexa Flow / Workflow OS | planned |
| P06 Universal Business Foundation | planned |
| P07 CRM, Sales & Customer 360 | planned |
| P08 Finance & ERP Core | planned |
| P09 Commerce OS | planned |
| P10 Payment Fabric | planned |
| P11 POS & Edge | planned |
| P12 Experience Builder & CMS | planned |
| P13 Portal Platform | planned |
| P14 HR, Projects & Service Operations | planned |
| P15 Supply Chain, Warehouse & Manufacturing | planned |
| P16 Omnexa Connect / Integration Fabric | planned |
| P17 Low-code App Builder | planned |
| P18 Data, Reporting & BI | planned |
| P19 Omnexa Intelligence Platform | planned |
| P20 Governed AI Agents | planned |
| P21 Developer Platform | planned |
| P22 Omnexa Exchange / Marketplace | planned |
| P23 Globalization & Country Packs | planned |
| P24 Enterprise Governance, Security & Compliance | planned |
| P25 Scale Fabric | planned |
| P26 Industry Packs | planned |
| P27 Autonomous Business OS | planned |

## Execution lock

Until P00 is complete, contributors and AI systems may work only on architecture/governance/specification activities belonging to P00, plus narrowly scoped repository-maintenance fixes.

Do not begin kernel implementation, database models, CRM, ERP, commerce, POS, website builder, payments or AI product code before the P00 exit gate.

## How status changes

A package moves to `done` only when its acceptance evidence satisfies `docs/governance/DEFINITION_OF_DONE.md` and `docs/roadmap/STATE.json` is reconciled in the same change.
