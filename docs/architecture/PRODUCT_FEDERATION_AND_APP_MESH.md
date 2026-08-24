# Omnexa Product Federation & App Mesh

Status: **Proposed strategic architecture under ADR-0011**

## 1. Purpose

Omnexa is intended to become the umbrella Business Operating System for products/systems developed across the organization, including standalone products such as a CMS/site-builder platform like Nexora and future workforce, industry, commerce, operations or specialist systems.

“Be part of Omnexa” must **not** mean copying every repository into one process/database or forcing every product to abandon its own architecture.

The target is a governed **Product Federation / App Mesh**.

## 2. Product attachment modes

A product may participate as one of four modes:

1. **Native module** — lives inside Omnexa module runtime and follows Omnexa module ownership/lifecycle directly.
2. **Embedded application** — separate application/runtime presented inside Omnexa Work/portal/UI shell through a governed app contract.
3. **Federated connected product** — standalone service/product with its own deployment/data plane, connected through versioned capabilities/events/data contracts/SSO.
4. **Edge/local product** — POS/device/local/air-gapped runtime synchronized through approved edge/federation contracts.

Mode is an architecture decision based on ownership, lifecycle, scale, regulatory/data-residency and runtime constraints—not on convenience.

## 3. Federation contract

Every federated product declares:

```text
product_id
product_version
publisher/owner
attachment_mode
supported Omnexa versions
tenant/org mapping
identity/SSO contract
capabilities provided/consumed
permissions/scopes
events published/consumed
workflow triggers/actions
data contracts + classifications
Business Graph contributions
System Graph contributions
AI tools/agent exposure
network destinations
secret requirements
storage/data-residency ownership
entitlements/licensing
health/readiness
SLOs
upgrade/deprecation
incident/support owner
kill/disconnect semantics
```

## 4. Identity and tenancy

Preferred user experience is shared identity and explicit tenant/org federation, not duplicate accounts.

Rules:

- Omnexa identity may authenticate/federate into a product, but the product still enforces its own approved capability boundaries;
- product-local identities, where unavoidable, require mapping and lifecycle policy;
- tenant/org mapping is explicit and cannot default globally;
- cross-product tokens are audience/scoped/expiring and revocable;
- a product must not trust browser-supplied tenant IDs as authority.

## 5. Data ownership

Federation never means one shared write database.

Example:

```text
Omnexa CRM owns Customer
Nexora-like product owns Site/Page/Theme
Omnexa Finance owns Invoice
Omnexa Payments owns Payment
```

Cross-product context uses:

- APIs/capabilities;
- events;
- workflow actions;
- approved projections;
- Business Graph references.

If a standalone product has its own authoritative domain, Omnexa references/coordinates it rather than duplicating its write state.

## 6. Business Graph federation

Federated products can contribute graph facts with source/provenance.

Example:

```text
Customer --owns--> Brand
Brand --operates--> Website
Website --implemented-by--> Nexora Product
Website --generates--> Lead
Lead --owned-by--> Omnexa CRM
```

Graph federation is read/semantic context. Mutation still flows through owning capabilities.

## 7. Workflow federation

P05 workflows can orchestrate federated products through declared actions/events.

Example:

```text
New Brand Approved
 -> create website workspace in Nexora-like product
 -> provision domain
 -> create CRM campaign
 -> create Finance cost center
 -> await approvals
 -> publish
```

Each step preserves idempotency, compensation/reconciliation and product authority.

## 8. AI federation

Federated product capabilities may become Omnexa AI tools only through the same governed tool registry.

AI must not receive raw administrative APIs simply because a product is first-party.

Tool declaration includes:

- product owner;
- capability version;
- tenant/org scope;
- risk;
- input/output schema;
- data classification;
- approval requirement;
- idempotency/recovery;
- audit/event behavior.

## 9. Experience federation

Omnexa Work may provide unified navigation/search/tasks/approvals/context while preserving the underlying product boundary.

UI contribution options:

- native route/component;
- embedded signed app surface;
- deep-link with federated identity;
- contextual cards/actions;
- unified search/entity results.

Embedding never grants backend authority.

## 10. Entitlement and commercial federation

A future entitlement/catalog layer may express which tenant can use which products/modules/plans.

Entitlement checks may gate availability but are not a substitute for authorization.

Commercial data such as subscriptions/licenses may be integrated with Finance/Commerce/Billing owners rather than duplicated inside the federation layer.

## 11. Operational federation

Every product supplies health/SLO/version/deployment identity suitable for platform-level visibility.

System Graph can map:

`Omnexa workflow -> federated product capability -> external/provider dependency`

Incident/blast-radius views can show affected tenants/products while preserving access controls.

## 12. Product trust classes

Suggested classes:

- first-party native;
- first-party federated;
- verified partner;
- third-party marketplace;
- external unmanaged integration.

First-party status does not bypass security, versioning, data, network or audit controls.

## 13. Lifecycle

Federated product lifecycle must support:

- register;
- verify compatibility;
- bind tenant/org;
- authorize scopes;
- activate;
- health monitor;
- upgrade;
- suspend;
- revoke/disconnect;
- export/transition;
- deprecate.

Disconnect must not delete historical business/audit evidence owned elsewhere.

## 14. Why this is necessary

Without federation, an “everything in Omnexa” strategy creates one of two failures:

1. a giant coupled monolith where every product shares implementation/data internals; or
2. many disconnected products with duplicate identity, permissions, workflows, AI and customer context.

Product Federation gives Omnexa a third option: **one governed operating system with multiple independently evolvable product boundaries**.
