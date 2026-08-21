# Omnexa Product Constitution

Status: **Baseline v1**  
Authority: repository governance  
Change class: **architectural**

## 1. Product category

Omnexa is a **Composable Enterprise Business Operating System**.

It is not defined as a monolithic ERP. ERP is one domain family inside a broader platform that unifies business identity, operations, finance, commerce, customer experience, content, payments, automation, data, integrations and governed AI execution.

## 2. Product objective

Build a platform on which an organization can progressively run most digital business functions without forcing those functions into one tightly coupled application.

The platform must support:

- independent installation, enablement, suspension, upgrade and removal of modules;
- shared identity, tenancy, policy, workflow, event, audit, configuration and observability foundations;
- API-first and event-driven integration between domains;
- SaaS, dedicated, private-cloud, self-hosted and later edge/on-prem deployment models;
- multi-organization, multi-company, multi-brand, multi-branch and multi-location structures;
- multi-language, RTL, multi-currency, timezone, tax and country-pack extensibility;
- low-code and third-party extensions without bypassing platform governance;
- AI that can understand and act through authorized platform capabilities.

## 3. Strategic moat

Omnexa must not compete by accumulating disconnected features faster than established suites.

The intended moat is the combination of:

1. **Universal Kernel** — one governed foundation for identity, tenancy, authorization, configuration, files, events, workflows, audit and observability.
2. **Extreme Modularity** — domains own their data and communicate through stable contracts.
3. **Universal Workflow** — cross-domain business processes can be composed, resumed, compensated and audited.
4. **Unified Business Graph** — people, organizations, products, locations, transactions and relationships have coherent platform identities.
5. **Governed AI Execution** — AI can reason across authorized business context and execute through auditable capabilities, not raw database access.

## 4. Product laws

The following laws apply across every phase.

### 4.1 Kernel before modules
Shared infrastructure is implemented once in the platform kernel. Domains consume it rather than reproduce it.

### 4.2 Domain data ownership
Each domain owns its write model. Cross-domain direct database writes are forbidden.

### 4.3 Contract-first interoperability
Cross-domain behavior uses documented APIs, events, workflow actions or approved read models.

### 4.4 Tenant isolation
Every tenant-owned entity, query path, background job, event and cache key must have an explicit tenant isolation strategy.

### 4.5 Policy before action
Protected actions are authorized through platform policy/capability checks. UI visibility is not authorization.

### 4.6 Auditability
Security-sensitive and business-significant state transitions are attributable, timestamped and queryable.

### 4.7 Versioning
Externally consumed contracts are versioned. Breaking changes require migration and compatibility strategy.

### 4.8 Failure isolation
Optional modules must be removable or disableable without corrupting unrelated domains.

### 4.9 Evidence-based scaling
Microservices, sharding, multi-region complexity and specialized infrastructure are introduced only after measurable need.

### 4.10 AI least authority
AI receives only the tools and scopes necessary for a task. High-risk actions support policy controls and human approval.

## 5. Organizational model

The platform hierarchy must support at least:

```text
Platform
  └── Tenant
      └── Organization
          ├── Legal Entity
          ├── Company
          ├── Business Unit
          ├── Brand
          ├── Branch
          ├── Store
          ├── Warehouse
          ├── Department
          ├── Team
          └── Location
```

A person may hold different relationships and roles across different scopes.

The authorization design therefore combines role-based access with relationship/context-aware policy evaluation.

## 6. Universal business primitives

The kernel/foundation must standardize identifiers and semantics for common concepts without forcing all domain data into one giant schema.

Foundational concepts include:

- Identity
- Party
- Person
- Organization
- Address
- Contact Point
- Location
- Money
- Currency
- Tax Context
- Product/Service reference
- File/Document
- Note/Activity reference
- Time/Audit metadata

Domain concepts such as Invoice, Opportunity, Purchase Order or Shipment remain owned by their respective domains.

## 7. Domain families

Omnexa is planned around the following major families:

- Platform Kernel
- Identity & Organization
- Workflow & Automation
- CRM / Sales / Customer 360
- Finance / Accounting / ERP
- Procurement / Inventory / Warehouse
- Commerce / Marketplace / Subscription
- Payments / Billing / Settlement
- POS / Edge
- Website / CMS / Experience Builder
- Portals
- HR / Workforce
- Projects / Service Operations
- Manufacturing / MRP / Maintenance
- Integration Fabric
- Low-code App Builder
- Analytics / BI / Data Platform
- AI Platform / Agents
- Developer Platform / SDK
- Marketplace / Extension Exchange
- Globalization / Country Packs
- Enterprise Governance / Security / Compliance
- Industry Packs

This catalog is a planning boundary, not permission to implement all domains simultaneously.

## 8. Approved baseline architecture

Initial execution strategy:

- **modular monolith first** for the control plane and early business domains;
- independently owned module boundaries inside one deployable platform;
- event fabric for decoupled domain notifications and async work;
- durable workflow abstraction for long-running cross-domain processes;
- PostgreSQL as primary transactional store;
- Redis-compatible ephemeral/cache layer;
- S3-compatible object storage;
- OpenTelemetry-compatible observability;
- service extraction only when operational evidence justifies it.

## 9. Language policy

Approved baseline:

- Go for kernel/backend platform services and domain services;
- TypeScript for web/admin/builder surfaces and primary extension SDK;
- Rust for edge, local agent, sandboxing, native device integration or high-assurance components when justified;
- Python for AI/data workloads only where its ecosystem materially improves the implementation.

Polyglot use is controlled. A new language requires an ADR explaining runtime ownership, supply-chain implications, observability, deployment and maintenance cost.

## 10. Deployment policy

The architecture must not depend on a single commercial cloud service.

Target deployment progression:

1. local developer environment;
2. single-region Omnexa Cloud;
3. dedicated tenant deployment;
4. private/self-hosted deployment;
5. regional data planes;
6. multi-region control plane/data-plane topology;
7. enterprise edge or air-gapped patterns where justified.

## 11. Security posture

Security is a platform concern, not a late hardening phase.

Foundation requirements include:

- tenant isolation;
- encryption in transit and at rest;
- secret management;
- MFA/passkey-ready identity;
- OIDC/OAuth2 compatibility;
- enterprise SSO/SAML/SCIM expansion path;
- granular authorization;
- attributable audit logs;
- secure-by-default APIs;
- rate limiting and abuse controls;
- signed/versioned modules and extension permissions;
- dependency/SBOM and supply-chain controls;
- backup, recovery and retention policy hooks.

## 12. AI posture

AI is an execution client of the platform, not a privileged bypass.

AI agents:

- consume authorized business context;
- call versioned capability/tool interfaces;
- inherit tenant and user/service-account scope;
- produce auditable actions;
- must support human approval for sensitive operations;
- must never receive unrestricted production database mutation authority.

## 13. Non-goals during foundation phases

Until explicitly opened in the roadmap, the project will not:

- build every business module in parallel;
- prematurely split domains into independent microservices;
- optimize for hypothetical planet-scale load before correctness and modularity are proven;
- allow module-specific authentication or tenant models;
- add duplicate payment, media, notification, workflow or permission subsystems inside domains;
- accept undocumented breaking contracts for delivery speed.

## 14. Success definition

Omnexa succeeds when new domains can be added without architectural decay, existing optional modules can be disabled without collateral corruption, business processes can safely span domains, and the platform can evolve from one deployable system into distributed services without changing its public capability contracts unnecessarily.
