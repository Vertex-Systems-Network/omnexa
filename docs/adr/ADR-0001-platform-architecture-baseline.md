# ADR-0001 — Platform Architecture Baseline

- Status: **Accepted**
- Date: 2026-08-21
- Decision class: Architecture baseline

## Context

Omnexa is intended to cover a very broad business surface: ERP, CRM, commerce, payments, POS, websites/CMS, portals, workflow, integrations, low-code, analytics and AI. Starting by implementing these as unrelated applications or as a large network of premature microservices would create incompatible identity, data, authorization and lifecycle models.

The platform also needs to support future module installation/removal without unrelated failures and should remain evolvable from a single deployable product into distributed services where justified.

## Problem

We need a foundation that:

- enforces coherent platform-wide capabilities;
- preserves domain ownership;
- minimizes operational complexity during early development;
- allows later service extraction;
- supports strong tenancy, authorization and audit boundaries;
- gives AI contributors a stable architecture to follow.

## Decision

Omnexa adopts the following baseline:

1. Product category is **Composable Enterprise Business Operating System**.
2. Architecture starts as a **strict modular monolith with service-ready boundaries**.
3. A shared **platform kernel** owns tenancy, identity foundation, authorization, module lifecycle, configuration, files, events/jobs/workflow infrastructure, audit and observability primitives.
4. Business modules own their transactional write models and communicate through versioned capabilities, events, workflows or approved read projections.
5. Direct cross-module database writes and reliance on another module's internal implementation are prohibited.
6. Primary backend language is **Go**; primary web/admin/builder language is **TypeScript**. Rust and Python are specialized runtimes requiring justified scope.
7. Primary transactional store is **PostgreSQL**; cache uses a Redis-compatible layer; files use S3-compatible object storage; messaging follows a NATS/JetStream-class abstraction; observability follows OpenTelemetry-compatible semantics.
8. Service extraction is permitted only when measurable scaling, isolation, ownership, regulatory or runtime requirements justify it.
9. AI acts through authenticated, authorized, auditable platform capabilities and does not receive unrestricted database mutation access.
10. Architecture and phase changes follow `docs/governance/CHANGE_CONTROL.md`.

## Alternatives considered

### A. Traditional monolithic ERP
Rejected because shared tables/internal calls across every domain would make independent module lifecycle and future extraction difficult.

### B. Microservices from day one
Rejected because it would add distributed transactions, networking, deployment, schema ownership, observability and local-development complexity before domain boundaries are proven.

### C. One programming language for every component
Rejected because UI/builder, core services, edge/native and AI workloads have materially different ecosystem requirements. Polyglot use is allowed but governed.

### D. Separate standalone products for CRM, ERP, commerce and website builder
Rejected as the primary architecture because it would duplicate identity, policy, tenant, event, workflow and business-context foundations and weaken the intended unified operating model.

## Consequences

### Positive

- one coherent foundation;
- strong domain ownership;
- reduced premature infrastructure complexity;
- safer module enable/disable behavior;
- explicit path to service extraction;
- consistent security/audit/tenant controls;
- predictable rules for human and AI contributors.

### Negative / cost

- strict boundaries require more contract discipline than direct table access;
- some cross-domain features require projections/events/workflows rather than convenient joins;
- architecture tests and module tooling must be built early;
- developers must learn kernel/domain ownership rules before shipping features.

## Compatibility impact

This is the first architecture baseline and therefore establishes, rather than breaks, compatibility rules.

Future changes to the decisions above require a new ADR and reconciliation of all affected governance/roadmap artifacts.

## Migration impact

No production application schema exists at this decision point. P00 is intentionally freezing the architecture before implementation begins.

## Security and tenancy impact

Security/tenancy are elevated to kernel-level concerns and may not be reimplemented independently by business modules.

## Operational impact

The initial deployment remains operationally simple compared with a broad microservice system, while telemetry, events and module boundaries prepare for future extraction.

## Supersession

This ADR remains authoritative until explicitly superseded by another accepted ADR.
