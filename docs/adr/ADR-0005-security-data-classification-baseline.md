# ADR-0005 — Security & Data Classification Baseline

Status: **Accepted**  
Date: **2026-08-21**  
Work package: **P00.06**

## Context

Omnexa is designed to host multiple tenants, business entities, financial/commerce operations, employee/customer data, integrations, extensible modules, POS/edge devices and governed AI agents. A generic “secure by best practice” statement is insufficient because future modules and AI implementers could interpret security boundaries differently.

The platform therefore needs explicit invariants for authentication, authorization, tenant isolation, data classification, secrets, encryption, audit, external integrations, support access, modules, CI/release and AI tool execution before runtime implementation begins.

## Decision

Omnexa adopts the following security architecture baseline:

1. Security is a platform property; no module may opt out or replace a shared control with a weaker private alternative.
2. Trust is explicit across client, service, module, tenant, organization, integration, device, support, CI/release and AI boundaries.
3. Authentication proves principal identity; authorization is independently enforced by owning capabilities.
4. Authorization combines RBAC, relationship-based authorization, contextual policy and capability boundaries.
5. Tenant isolation is mandatory across primary data, cache, search, analytics, files, events, queues, backups and AI/vector data.
6. Client-supplied tenant/organization identifiers are context, never authorization authority.
7. Omnexa uses four confidentiality classes: `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`, plus handling tags for specific risk classes.
8. Tenant-owned business records default to at least `CONFIDENTIAL` unless an owner explicitly publishes a `PUBLIC` projection.
9. Secrets/credentials/key material are `RESTRICTED`, excluded from source control/logging/AI, and managed through approved secret/key mechanisms.
10. Protected traffic uses authenticated encryption in transit; production data stores/backups use approved encryption at rest. Encryption never lowers classification or replaces authorization.
11. Passwords are never reversibly stored; password authentication uses approved adaptive hashing.
12. Audit is a protected business/security record separate from debug logging.
13. Privileged operations, support impersonation, bulk export, replay, destructive purge and high-impact AI actions require explicit capability/policy/audit and may require stronger approval/authentication by risk.
14. External integrations/webhooks are independent trust/disclosure boundaries with scoped principals, authenticity checks, data minimization and egress/SSRF controls.
15. Modules/extensions declare permissions and external access; future marketplace packages require integrity/provenance controls.
16. Production `CONFIDENTIAL`/`RESTRICTED` data does not flow to lower environments by default.
17. AI agents are untrusted planners that can act only through least-privilege governed capabilities; prompt injection/content cannot grant new authority.
18. Data retention/export/deletion/search/analytics/AI behavior is classification-aware and policy-driven.
19. Security exceptions are explicit, owned, compensating-control-backed and time/review bounded.

Normative details live in:

- `docs/security/SECURITY_STANDARD.md`
- `docs/security/DATA_CLASSIFICATION.md`
- `docs/security/SECURITY_CONTROL_MATRIX.md`
- `docs/contracts/security/data-classification.schema.json`

## Consequences

### Positive

- security behavior becomes deterministic across future modules and AI implementations;
- tenant isolation is treated as a system-wide invariant rather than a query convention;
- sensitive data cannot silently leak into logs, search, analytics, events or AI pipelines;
- extension/integration/AI boundaries have explicit authority models;
- future testing can be organized around named controls and negative evidence.

### Costs

- capabilities/modules must classify data and declare security behavior;
- support/admin shortcuts are intentionally constrained;
- lower-environment debugging cannot rely on uncontrolled production copies;
- some high-risk actions require extra approval/audit/operational tooling;
- platform kernel must implement reusable identity, authorization, audit, secrets/crypto and isolation mechanisms before business domains scale.

## Rejected alternatives

### Authentication-only security

Rejected because a logged-in principal still requires tenant/object/capability authorization.

### Role-only authorization

Rejected because multi-tenant/business-scope relationships require object and contextual policy beyond global role labels.

### One “sensitive” classification

Rejected because public/internal/business-confidential/secrets require materially different handling.

### Network-is-trusted internal services

Rejected because same-network/process/module origin is not sufficient proof of authority.

### Hidden super-admin bypass

Rejected because it defeats auditability, tenant isolation and governed support access.

### Separate relaxed AI permissions

Rejected because AI execution must remain under normal capability/authorization controls.

## Compatibility

This ADR strengthens and depends on:

- ADR-0001 platform architecture;
- ADR-0002 identity/data primitive conventions;
- ADR-0003 HTTP API semantics;
- ADR-0004 event reliability/tenant-context semantics.

P00.07 may define executable security scanning/testing gates. P00.09 may add threat-model priorities and SLO/incident expectations. Neither may silently weaken these security boundaries without a superseding ADR.

## Supersession

Material changes to tenant isolation, authentication/session semantics, authorization model, classification levels, secrets/key handling, audit integrity, support impersonation, extension trust or AI execution authority require a superseding ADR plus governance/state reconciliation before implementation.
