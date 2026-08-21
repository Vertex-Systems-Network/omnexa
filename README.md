# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed, modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are planned as domain families running on one shared platform kernel.

> **Current execution lock:** the repository is in **P00 — Product Constitution & Architecture Freeze**. Kernel and business-feature implementation must not begin until the P00 exit gate is complete. Current package: **P00.05 — Event contract standard**.

## Mandatory contributor / AI start here

Any human or AI system changing this repository must read the following in order:

1. [`AGENTS.md`](AGENTS.md)
2. [`docs/governance/PRODUCT_CONSTITUTION.md`](docs/governance/PRODUCT_CONSTITUTION.md)
3. [`docs/architecture/SYSTEM_ARCHITECTURE.md`](docs/architecture/SYSTEM_ARCHITECTURE.md)
4. [`docs/architecture/MODULE_STANDARD.md`](docs/architecture/MODULE_STANDARD.md)
5. [`docs/architecture/GLOSSARY.md`](docs/architecture/GLOSSARY.md)
6. [`docs/architecture/NAMING_STANDARD.md`](docs/architecture/NAMING_STANDARD.md)
7. [`docs/architecture/DOMAIN_OWNERSHIP.md`](docs/architecture/DOMAIN_OWNERSHIP.md)
8. [`docs/architecture/DEPENDENCY_MATRIX.md`](docs/architecture/DEPENDENCY_MATRIX.md)
9. [`docs/architecture/IDENTIFIER_STANDARD.md`](docs/architecture/IDENTIFIER_STANDARD.md)
10. [`docs/architecture/MONEY_STANDARD.md`](docs/architecture/MONEY_STANDARD.md)
11. [`docs/architecture/TIME_STANDARD.md`](docs/architecture/TIME_STANDARD.md)
12. [`docs/architecture/LOCALE_STANDARD.md`](docs/architecture/LOCALE_STANDARD.md)
13. [`docs/architecture/ERROR_STANDARD.md`](docs/architecture/ERROR_STANDARD.md)
14. [`docs/architecture/API_STANDARD.md`](docs/architecture/API_STANDARD.md)
15. [`docs/roadmap/MASTER_PLAN.md`](docs/roadmap/MASTER_PLAN.md)
16. [`docs/roadmap/STATUS.md`](docs/roadmap/STATUS.md)
17. [`docs/roadmap/STATE.json`](docs/roadmap/STATE.json)
18. [`docs/governance/AI_EXECUTION_POLICY.md`](docs/governance/AI_EXECUTION_POLICY.md)
19. [`docs/governance/CHANGE_CONTROL.md`](docs/governance/CHANGE_CONTROL.md)
20. [`docs/governance/DEFINITION_OF_DONE.md`](docs/governance/DEFINITION_OF_DONE.md)
21. Relevant ADRs under [`docs/adr/`](docs/adr/)

`STATE.json` is the machine-readable canonical execution state. Work outside the active package is not implicitly authorized.

## Core architecture laws

- Kernel before business modules.
- Every authoritative write model has one owner.
- Modules own their write models and schemas.
- Cross-module direct database writes are forbidden.
- Cross-domain integration uses versioned capabilities, events, workflows or approved read projections.
- Tenant scope, authorization, audit and observability are platform requirements, not optional module features.
- Optional modules must fail/degrade independently.
- Public contracts are versioned.
- AI acts only through governed, authorized, auditable capabilities.
- Omnexa begins as a strict modular monolith and extracts services only when evidence justifies it.
- Architecture or roadmap changes require formal change control and ADR reconciliation.

## Frozen foundation primitives (P00.03)

- **Identifiers:** UUIDv7 by default; PostgreSQL native `uuid`; canonical tenant scope `tenant_id`.
- **Money:** exact decimal + explicit currency; JSON decimal strings; PostgreSQL baseline `NUMERIC(38,18)`; no binary floating-point money.
- **Time:** UTC instants + `timestamptz`; date-only business dates; IANA timezone semantics for civil time/recurrence.
- **Locale:** BCP 47 language tags; locale/country/currency/timezone are independent; RTL is first-class.
- **Errors:** stable machine error codes, safe structured problem details, request/trace correlation, no public stack traces/SQL/secrets.

These decisions are recorded by [`ADR-0002`](docs/adr/ADR-0002-foundation-data-conventions.md).

## Frozen HTTP API baseline (P00.04)

- Stable routes use `/api/v{major}/{domain}/{resources}`.
- Canonical stable HTTP contracts use OpenAPI 3.2.0.
- JSON fields use lowercase `snake_case`.
- Errors use `application/problem+json` plus stable Omnexa codes and request correlation.
- Scalable lists use opaque cursor pagination.
- Protected retriable mutations define `Idempotency-Key` semantics.
- Lost-update-sensitive mutations use explicit optimistic concurrency such as `ETag` / `If-Match` where applicable.
- Filters, sorts and includes are allowlisted, bounded and authorization-aware.
- Business actions are explicit capability/action operations rather than arbitrary status patches.
- Client-supplied tenant or organization identifiers never become authorization authority.

The contract is defined by [`API_STANDARD.md`](docs/architecture/API_STANDARD.md), [`openapi-template.yaml`](docs/contracts/http/openapi-template.yaml), and [`ADR-0003`](docs/adr/ADR-0003-http-api-contract-baseline.md).

## Technology baseline

Until superseded by an accepted ADR:

- **Go** — platform kernel/backend and primary domain services
- **TypeScript + React** — admin, web, builder and primary extension SDK surfaces
- **Rust** — edge/native/security-sensitive components where justified
- **Python** — AI/data workloads where ecosystem value justifies it
- **PostgreSQL** — primary transactional store
- **Redis-compatible layer** — cache/ephemeral coordination
- **S3-compatible object storage** — files/media
- **NATS/JetStream-class fabric** — event/messaging baseline
- **OpenTelemetry** — observability semantics

## Governance hardening

Repository-level controls include:

- [`CONTRIBUTING.md`](CONTRIBUTING.md)
- [`SECURITY.md`](SECURITY.md)
- [`.github/CODEOWNERS`](.github/CODEOWNERS)
- architecture/bug issue templates
- governance CI via [`.github/workflows/governance.yml`](.github/workflows/governance.yml)
- dependency-free state validator at [`scripts/validate_governance.py`](scripts/validate_governance.py)
- hosted repository ruleset target in [`docs/governance/REPOSITORY_HARDENING.md`](docs/governance/REPOSITORY_HARDENING.md)
- licensing/IP decision gate in [`docs/governance/LICENSING_DECISION.md`](docs/governance/LICENSING_DECISION.md)

Hosted branch protection and licensing/IP/trademark decisions remain explicitly tracked; they do not authorize early implementation.

## Roadmap

The canonical plan is [`docs/roadmap/MASTER_PLAN.md`](docs/roadmap/MASTER_PLAN.md), covering **P00 through P27**:

```text
P00 Architecture/Governance
 -> P01 Kernel
 -> P02 Identity/Tenancy
 -> P03 Module Runtime
 -> P04 Event/Data Fabric
 -> P05 Workflow OS
 -> P06 Business Foundation
 -> P07-P15 Core Business Domains
 -> P16-P18 Integration/Low-code/Data
 -> P19-P20 Intelligence/Agents
 -> P21-P22 Developer Ecosystem/Marketplace
 -> P23-P25 Global/Enterprise/Scale
 -> P26 Industry Packs
 -> P27 Autonomous Business OS
```

See [`docs/roadmap/STATUS.md`](docs/roadmap/STATUS.md) for the current human-readable status and [`docs/roadmap/STATE.json`](docs/roadmap/STATE.json) for canonical machine state.

## Work-package discipline

Material work must use [`docs/governance/WORK_PACKAGE_TEMPLATE.md`](docs/governance/WORK_PACKAGE_TEMPLATE.md) and satisfy [`docs/governance/DEFINITION_OF_DONE.md`](docs/governance/DEFINITION_OF_DONE.md).

Execution history is recorded in the append-only [`docs/roadmap/EXECUTION_LEDGER.md`](docs/roadmap/EXECUTION_LEDGER.md).

## Architecture decisions

Architecture decisions live under [`docs/adr/`](docs/adr/). Current accepted baselines include:

- [`ADR-0001 — Platform Architecture Baseline`](docs/adr/ADR-0001-platform-architecture-baseline.md)
- [`ADR-0002 — Foundation Data & Contract Conventions`](docs/adr/ADR-0002-foundation-data-conventions.md)
- [`ADR-0003 — HTTP API Contract Baseline`](docs/adr/ADR-0003-http-api-contract-baseline.md)

Do not implement an architectural change first and document it afterward.

## Product principle

Omnexa is not trying to win by collecting the largest number of disconnected features. The intended long-term advantage is:

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
