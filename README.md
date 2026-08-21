# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is being designed as a governed, modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are planned as domain families running on one shared platform kernel.

> **Current execution lock:** the repository is in **P00 — Product Constitution & Architecture Freeze**. Kernel and business-feature implementation must not begin until the P00 exit gate is complete. Current package: **P00.08 — Local developer and repository structure specification**.

> **Temporary CI note:** GitHub Actions allowance is currently exhausted/disabled. ADR-0006 authorizes a temporary P00 documentation/specification-only manual evidence path. Hosted CI is `BLOCKED`/`NOT RUN`, never `PASS`, and this exception expires before any P01 implementation merge.

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
15. [`docs/architecture/EVENT_STANDARD.md`](docs/architecture/EVENT_STANDARD.md)
16. [`docs/security/SECURITY_STANDARD.md`](docs/security/SECURITY_STANDARD.md)
17. [`docs/security/DATA_CLASSIFICATION.md`](docs/security/DATA_CLASSIFICATION.md)
18. [`docs/security/SECURITY_CONTROL_MATRIX.md`](docs/security/SECURITY_CONTROL_MATRIX.md)
19. [`docs/quality/TESTING_STANDARD.md`](docs/quality/TESTING_STANDARD.md)
20. [`docs/quality/CI_STANDARD.md`](docs/quality/CI_STANDARD.md)
21. [`docs/quality/RELEASE_STANDARD.md`](docs/quality/RELEASE_STANDARD.md)
22. [`docs/quality/QUALITY_GATE_MATRIX.md`](docs/quality/QUALITY_GATE_MATRIX.md)
23. [`docs/roadmap/MASTER_PLAN.md`](docs/roadmap/MASTER_PLAN.md)
24. [`docs/roadmap/STATUS.md`](docs/roadmap/STATUS.md)
25. [`docs/roadmap/STATE.json`](docs/roadmap/STATE.json)
26. [`docs/governance/AI_EXECUTION_POLICY.md`](docs/governance/AI_EXECUTION_POLICY.md)
27. [`docs/governance/CHANGE_CONTROL.md`](docs/governance/CHANGE_CONTROL.md)
28. [`docs/governance/DEFINITION_OF_DONE.md`](docs/governance/DEFINITION_OF_DONE.md)
29. If active, [`docs/governance/CI_EVIDENCE_EXCEPTION_2026-08-22.md`](docs/governance/CI_EVIDENCE_EXCEPTION_2026-08-22.md)
30. Relevant ADRs under [`docs/adr/`](docs/adr/)

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

## Frozen foundation contracts

### P00.03 — Primitives

UUIDv7 identifiers, exact-decimal money, UTC/`timestamptz` + IANA civil-time semantics, BCP 47 locale/RTL rules and stable structured error contracts are frozen by ADR-0002.

### P00.04 — HTTP APIs

Stable APIs use `/api/v{major}/{domain}/{resources}`, OpenAPI 3.2.0, `snake_case`, Problem Details errors, cursor pagination, explicit idempotency/concurrency semantics and authorization-derived tenant context. See ADR-0003.

### P00.05 — Events

Events are producer-owned immutable past-tense facts using `<domain>.<subject>.<past_tense_fact>.v<major>`, CloudEvents-compatible envelopes, UUIDv7 identities, at-least-once assumptions, idempotent consumers, outbox/inbox reliability, bounded retry/DLQ and replay-safe semantics. See ADR-0004.

### P00.06 — Security and classification

Data classes are `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`. Authorization combines RBAC, relationships, contextual policy and governed capabilities. Tenant isolation, secrets/KMS, audit, privileged operations, integrations/webhooks/SSRF, modules/supply chain and AI execution are governed platform invariants. See ADR-0005.

### P00.07 — Testing, CI and release

Quality semantics are **repository-owned and CI-provider independent**. Local development and any approved CI provider must execute the same canonical gates.

Gate classes:

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

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`; blocked/unrun work is never silently green. Testing is risk-based and requires negative evidence for affected tenant, authorization, replay/idempotency, lifecycle and security boundaries. Releases use semantic-versioning semantics, immutable source/artifact identity and a build-once/promote model where possible.

Normative documents:

- [`TESTING_STANDARD.md`](docs/quality/TESTING_STANDARD.md)
- [`CI_STANDARD.md`](docs/quality/CI_STANDARD.md)
- [`RELEASE_STANDARD.md`](docs/quality/RELEASE_STANDARD.md)
- [`QUALITY_GATE_MATRIX.md`](docs/quality/QUALITY_GATE_MATRIX.md)
- [`quality-gates.schema.json`](docs/contracts/quality/quality-gates.schema.json)
- [`ADR-0007`](docs/adr/ADR-0007-testing-ci-release-baseline.md)

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

Repository-level controls include `CONTRIBUTING.md`, `SECURITY.md`, CODEOWNERS, issue/ADR templates, the governance validator/workflow definition, repository hardening specification and licensing/IP decision gate.

Current hosted blockers:

- Issue #3: `main` branch/ruleset protection is still not applied through the available connector.
- Issue #14: GitHub Actions execution is temporarily unavailable due exhausted/disabled allowance; ADR-0006 applies only to P00 documentation/specification work.
- Issue #4: final licensing/IP/trademark strategy remains a decision gate before external distribution/public launch.

None of these authorize early implementation.

## Roadmap

The canonical [`MASTER_PLAN.md`](docs/roadmap/MASTER_PLAN.md) covers P00 through P27:

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

See [`STATUS.md`](docs/roadmap/STATUS.md) and [`STATE.json`](docs/roadmap/STATE.json) for current execution state.

## Work-package discipline

Material work must satisfy [`DEFINITION_OF_DONE.md`](docs/governance/DEFINITION_OF_DONE.md) and use the governed work-package/change-control process. Execution history is append-only in [`EXECUTION_LEDGER.md`](docs/roadmap/EXECUTION_LEDGER.md).

## Accepted architecture decisions

ADR-0001 through ADR-0007 define the current foundation, with ADR-0006 being a temporary P00 operational exception rather than a permanent architecture weakening.

Do not implement an architectural change first and document it afterward.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**