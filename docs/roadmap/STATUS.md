# Omnexa Program Status

Last reconciled: **2026-08-21**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active**
- Current work package: **P00.03 — ID, money, time, locale and error conventions**
- Business-feature implementation: **NOT AUTHORIZED YET**
- Kernel implementation: **NOT AUTHORIZED YET**

## P00 work packages

| ID | Work package | State | Evidence / note |
|---|---|---|---|
| P00.01 | Repository governance baseline | done | `AGENTS.md`, Product Constitution, Architecture, Module Standard, Change Control, DoD, Master Plan, State and baseline ADR |
| P00.02 | Product/domain glossary and naming standard | done | `GLOSSARY.md`, `NAMING_STANDARD.md`, `DOMAIN_OWNERSHIP.md`, `DEPENDENCY_MATRIX.md`, contribution/security/hardening controls and governance CI baseline |
| P00.03 | ID, money, time, locale and error conventions | active | Current canonical work package |
| P00.04 | API contract standard | planned | Depends on P00.02/P00.03 |
| P00.05 | Event contract standard | planned | Depends on P00.02/P00.03 |
| P00.06 | Security and data-classification baseline | planned | Depends on P00.02/P00.03 |
| P00.07 | Testing/CI/release standard | planned | Full executable quality/release gates; initial governance-only workflow already exists |
| P00.08 | Local developer and repository structure specification | planned | Implementation skeleton specification only until earlier conventions settle |
| P00.09 | Initial threat model and operational SLO targets | planned | Must precede foundation freeze |
| P00.10 | Foundation architecture freeze review | planned | Final P00 exit gate |

P00 package progress: **2 / 10 done**.

## Governance hardening status

File-level governance is now specified/implemented through:

- `.github/CODEOWNERS`;
- contribution and security policies;
- architecture-change and bug templates;
- ADR template;
- repository hardening specification;
- dependency-free governance state validator;
- GitHub Actions governance workflow.

Two hosted/business decisions remain explicitly tracked and must not be forgotten:

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
