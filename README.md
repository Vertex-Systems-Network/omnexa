# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed, modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are planned as domain families running on one shared platform foundation.

> **Execution lock:** Omnexa is still in **P00 — Product Constitution & Architecture Freeze**. Kernel and business-feature implementation are not authorized until the P00 exit gate is complete. Current package: **P00.07 — Testing, CI and release standard**.

## Mandatory contributor / AI start here

Read in this order before making material changes:

1. [`AGENTS.md`](AGENTS.md)
2. [`docs/governance/PRODUCT_CONSTITUTION.md`](docs/governance/PRODUCT_CONSTITUTION.md)
3. architecture standards under [`docs/architecture/`](docs/architecture/)
4. security standards under [`docs/security/`](docs/security/)
5. [`docs/roadmap/MASTER_PLAN.md`](docs/roadmap/MASTER_PLAN.md)
6. [`docs/roadmap/STATUS.md`](docs/roadmap/STATUS.md)
7. [`docs/roadmap/STATE.json`](docs/roadmap/STATE.json)
8. governance policies under [`docs/governance/`](docs/governance/)
9. relevant ADRs under [`docs/adr/`](docs/adr/)

`STATE.json` is the canonical machine-readable execution state. Work outside the active package is not implicitly authorized.

## Core architecture laws

- Kernel before business modules.
- Every authoritative write model/capability has one owner.
- Modules own private write state; cross-module direct database writes are forbidden.
- Cross-domain integration uses versioned capabilities/APIs, events, workflows or approved read projections.
- Tenant scope, authorization, audit, classification and observability are platform properties.
- Optional modules must fail/degrade independently.
- Public/external contracts are versioned.
- AI acts only through governed, authorized, auditable capabilities.
- Omnexa starts as a strict modular monolith and extracts services only when evidence justifies it.
- Architecture changes require formal change control and ADR reconciliation.

## Frozen P00 baselines

### P00.03 — Foundation primitives

- UUIDv7 canonical IDs; PostgreSQL native `uuid`; canonical tenant scope `tenant_id`.
- exact-decimal money with explicit currency; no binary floating-point money.
- UTC/`timestamptz` instants; date-only business dates; IANA timezone semantics for civil recurrence.
- BCP 47 locale model; locale/country/currency/timezone are independent; RTL is first-class.
- stable machine error codes with disclosure-safe structured problem details.

See ADR-0002 and the primitive standards in `docs/architecture/`.

### P00.04 — Stable HTTP API

- stable route major versioning: `/api/v{major}/{domain}/{resources}`;
- OpenAPI 3.2.0 canonical contract baseline;
- lowercase `snake_case` JSON fields;
- `application/problem+json` errors;
- cursor pagination, bounded filters/sorts/includes;
- `Idempotency-Key` for protected retriable mutations;
- explicit optimistic concurrency where lost updates matter;
- client-provided tenant/organization IDs never become authorization authority.

See [`API_STANDARD.md`](docs/architecture/API_STANDARD.md) and ADR-0003.

### P00.05 — Events

- producer-owned `<domain>.<subject>.<past_tense_fact>.v<major>` events;
- CloudEvents-compatible structured envelope with UUIDv7 identity;
- producer-derived tenant context and explicit correlation/causation/tracing;
- at-least-once delivery assumption with idempotent consumers;
- transactional outbox + inbox/deduplication or equivalent guarantees for business-significant flows;
- explicit subject ordering only, bounded retry, dead-letter/quarantine, replay-safe side effects.

See [`EVENT_STANDARD.md`](docs/architecture/EVENT_STANDARD.md) and ADR-0004.

### P00.06 — Security & data classification

Canonical confidentiality classes:

```text
PUBLIC
INTERNAL
CONFIDENTIAL
RESTRICTED
```

Security baseline includes:

- zero implicit trust across client, service, module, tenant, device, integration, support, CI/release and AI boundaries;
- distinct human/service/workload/device/integration/support/AI principals;
- authentication separate from authorization;
- authorization = RBAC + relationship + contextual policy + capability boundaries;
- tenant isolation across OLTP, cache, files, search, analytics, events, backups and AI/vector data;
- classification-aware logging, search, analytics, export, retention, deletion and AI behavior;
- secrets/key material managed as `RESTRICTED` and excluded from ordinary logs/source/AI;
- governed privileged actions, support impersonation, event replay and high-impact AI execution;
- external integration/webhook authenticity and SSRF/egress controls;
- extension/module permission declaration and future package provenance/integrity controls;
- sensitive production data excluded from lower environments by default.

Normative documents:

- [`SECURITY_STANDARD.md`](docs/security/SECURITY_STANDARD.md)
- [`DATA_CLASSIFICATION.md`](docs/security/DATA_CLASSIFICATION.md)
- [`SECURITY_CONTROL_MATRIX.md`](docs/security/SECURITY_CONTROL_MATRIX.md)
- [`data-classification.schema.json`](docs/contracts/security/data-classification.schema.json)
- [`ADR-0005`](docs/adr/ADR-0005-security-data-classification-baseline.md)

## Technology baseline

Until superseded by an accepted ADR:

- **Go** — kernel/backend and primary domain services
- **TypeScript + React** — web/admin/builder/primary extension SDK surfaces
- **Rust** — edge/native/security-sensitive components where justified
- **Python** — AI/data workloads where ecosystem value justifies it
- **PostgreSQL** — primary transactional store
- **Redis-compatible layer** — cache/ephemeral coordination
- **S3-compatible object storage** — files/media
- **NATS/JetStream-class fabric** — event/messaging baseline
- **OpenTelemetry** — observability semantics

## Governance hardening

Repository controls include CODEOWNERS, contribution/security policies, issue/ADR templates, governance CI and the dependency-free state validator.

Two explicit non-code decisions remain tracked:

- **Issue #3:** hosted `main` branch/ruleset protection. GitHub still reports `main` as unprotected and the connected GitHub toolset does not expose a protection/ruleset write mutation.
- **Issue #4:** final licensing/IP/trademark strategy. Existing GPLv3 is not treated as implicit approval for the eventual commercial distribution model.

Neither item authorizes early kernel/business implementation.

## Roadmap

Canonical roadmap: [`docs/roadmap/MASTER_PLAN.md`](docs/roadmap/MASTER_PLAN.md), P00 through P27.

```text
P00 Architecture/Governance
 -> P01 Kernel
 -> P02 Identity/Tenancy
 -> P03 Module Runtime
 -> P04 Data/Event Fabric
 -> P05 Workflow OS
 -> P06 Business Foundation
 -> P07-P15 Business Domains
 -> P16-P18 Integration/Low-code/Data
 -> P19-P20 Intelligence/Agents
 -> P21-P22 Developer Ecosystem/Marketplace
 -> P23-P25 Global/Enterprise/Scale
 -> P26 Industry Packs
 -> P27 Autonomous Business OS
```

Current human status: [`STATUS.md`](docs/roadmap/STATUS.md). Canonical machine state: [`STATE.json`](docs/roadmap/STATE.json).

## Work-package discipline

Material work must satisfy [`DEFINITION_OF_DONE.md`](docs/governance/DEFINITION_OF_DONE.md). Execution history is append-only in [`EXECUTION_LEDGER.md`](docs/roadmap/EXECUTION_LEDGER.md).

Current accepted architecture baselines are ADR-0001 through ADR-0005. Do not implement an architectural change first and document it afterward.

## Product principle

Omnexa's intended long-term advantage is:

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
