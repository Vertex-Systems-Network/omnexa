# Omnexa Master Plan

Status: **Canonical Roadmap v1**  
Execution mode: gated, evidence-based, architecture-first

## 1. Roadmap rules

1. No implementation phase may be skipped because a later feature appears easier or more visible.
2. A phase may start only when its mandatory predecessors satisfy their exit gates.
3. Foundation phases are deliberately sequential. Parallel business-domain streams open only after module, tenancy, event and workflow contracts are proven.
4. Every phase is split into work packages with explicit acceptance criteria before coding begins.
5. `docs/roadmap/STATE.json` is the machine-readable state; this document defines intent, dependencies and gates.
6. A phase is `done` only when required evidence exists under `DEFINITION_OF_DONE.md`.
7. Scope additions are routed through `CHANGE_CONTROL.md` rather than silently absorbed.

## 2. Product progression

```text
P00 Governance / Constitution
      ↓
P01 Kernel
      ↓
P02 Identity & Tenancy
      ↓
P03 Module Runtime
      ↓
P04 Data & Event Fabric
      ↓
P05 Workflow OS
      ↓
P06 Business Foundation
      ↓
┌───────────────────────────────────────────────┐
│ Governed domain expansion                    │
│ P07 CRM  P08 ERP  P09 Commerce ...           │
└───────────────────────────────────────────────┘
      ↓
Integration + Low-code + Data
      ↓
AI Platform + Agents
      ↓
Developer Ecosystem + Marketplace
      ↓
Global / Enterprise / Scale
      ↓
Industry + Autonomous Business OS
```

---

# FOUNDATION PROGRAM

## P00 — Product Constitution & Architecture Freeze

**Purpose:** establish immutable-enough rules so future contributors and AI systems cannot redefine Omnexa accidentally.

### Work packages

- P00.01 Repository governance baseline
- P00.02 Product/domain glossary and naming standard
- P00.03 ID, money, time, locale and error conventions
- P00.04 API contract standard
- P00.05 Event contract standard
- P00.06 Security and data-classification baseline
- P00.07 Testing/CI/release standard
- P00.08 Local developer and repository structure specification
- P00.09 Initial threat model and operational SLO targets
- P00.10 Foundation architecture freeze review

### Exit gate

- governance documents reconciled;
- baseline ADRs accepted;
- module/API/event conventions machine-testable or specified precisely enough to implement validators;
- technology baseline recorded;
- repository structure approved;
- initial security/threat model documented;
- no unresolved architecture contradiction.

**No business module implementation before P00 exit.**

---

## P01 — Omnexa Kernel

**Purpose:** create the minimum platform runtime every later module will use.

### Work packages

- P01.01 Go workspace/build skeleton
- P01.02 configuration and environment system
- P01.03 structured error/result conventions
- P01.04 PostgreSQL connection/migration foundation
- P01.05 cache abstraction
- P01.06 object/file storage abstraction
- P01.07 structured logging and OpenTelemetry baseline
- P01.08 health/readiness/diagnostics
- P01.09 job/scheduler primitives
- P01.10 feature-flag/configuration registry
- P01.11 audit transport foundation
- P01.12 developer CLI baseline

### Exit gate

Fresh install from zero must start, migrate, health-check, emit telemetry and run tests reproducibly without hidden manual steps.

---

## P02 — Identity, Tenancy & Organization

**Purpose:** establish security and organizational context before domain data exists.

### Work packages

- user/service identity;
- tenant lifecycle;
- organization/company/legal-entity/branch/team/location hierarchy;
- authentication/session model;
- RBAC foundation;
- relationship/context-aware authorization foundation;
- MFA/passkey-ready flows;
- service accounts and API credentials;
- tenant-scoped settings;
- identity/permission audit trails.

### Exit gate

Tests prove cross-tenant isolation, object/scope permissions, role differences, service-account scoping and session invalidation behavior.

---

## P03 — Module Runtime

**Purpose:** prove Omnexa is actually modular before adding business breadth.

### Work packages

- module manifest schema;
- registry/discovery;
- dependency graph resolver;
- install/enable/disable/suspend/archive/detach/purge state machine;
- module settings and feature flags;
- capability registry;
- permission registration;
- UI contribution registry contract;
- migration ownership registry;
- health reporting;
- signed package hooks for later marketplace use.

### Exit gate

Reference test modules demonstrate:

1. required dependency enforcement;
2. optional dependency degradation;
3. safe disable/re-enable;
4. upgrade/migration path;
5. forbidden cross-module dependency detection;
6. health/state accuracy;
7. no unrelated module corruption after lifecycle operations.

---

## P04 — Data, Jobs & Event Fabric

**Purpose:** provide reliable asynchronous and decoupled communication.

### Work packages

- event envelope/version rules;
- publish/subscribe abstraction;
- durable stream/consumer baseline;
- outbox/inbox reliability pattern;
- idempotency primitives;
- retry/backoff policy;
- dead-letter/quarantine path;
- correlation/causation propagation;
- event schema registry/validation;
- background job ownership/tenant context.

### Exit gate

Replay, duplicate delivery, consumer failure, restart and poison-event scenarios pass without double-applying protected business mutations.

---

## P05 — Omnexa Flow / Workflow OS

**Purpose:** run long-lived cross-domain business processes safely.

### Work packages

- workflow definition/version model;
- trigger/action registry;
- conditions/branches;
- timers/waits;
- retries/timeouts;
- human approvals;
- parallel paths;
- compensation/saga semantics;
- workflow state persistence;
- audit/timeline;
- workflow testing/simulation API;
- visual-designer contract for later UI.

### Exit gate

A reference workflow survives process restart and demonstrates failure, retry, compensation, wait, approval and idempotent resume.

---

## P06 — Universal Business Foundation

**Purpose:** establish reusable business primitives without turning them into a giant shared-domain database.

### Work packages

- Party/Person/Organization references;
- addresses/contact points;
- locations;
- currency/money semantics;
- tax-context primitives;
- product/service reference primitives;
- files/documents/notes/activity references;
- common numbering/reference services where approved;
- import/export primitives;
- search/index integration foundation.

### Exit gate

Reference domains can consume common primitives without direct write coupling or duplicating platform identity concepts.

---

# CORE BUSINESS PROGRAM

## P07 — CRM, Sales & Customer 360

Scope:

- leads;
- contacts/organizations;
- opportunities/pipelines;
- activities;
- quotes/proposals foundation;
- customer timeline/360 projection;
- sales territories/ownership;
- scoring/sequences later within the phase only after core records stabilize.

Exit: tenant-safe CRM lifecycle and published capability/event contracts proven.

---

## P08 — Finance & ERP Core

Scope:

- chart of accounts;
- general ledger;
- journals;
- accounts receivable/payable;
- invoices/credit notes;
- expenses;
- cash/bank abstractions;
- tax integration points;
- budgeting/fixed-assets progression;
- procurement foundation where accounting ownership is clear.

Exit: double-entry invariants, period controls, immutable audit history, multi-currency rules and reconciliation tests proven.

---

## P09 — Commerce OS

Scope:

- products/variants/catalogs;
- price lists;
- carts/checkouts;
- orders;
- discounts/promotions;
- returns/refunds orchestration;
- subscriptions/bundles progression;
- B2C/B2B/wholesale channel model;
- marketplace/multi-vendor primitives only behind approved module boundaries.

Exit: channel-independent order lifecycle works through published contracts without provider lock-in.

---

## P10 — Payment Fabric

Scope:

- gateway/provider abstraction;
- authorize/capture/void/refund;
- token/provider references;
- webhooks;
- recurring/mandate abstraction;
- settlement/reconciliation;
- payout/dispute expansion;
- PCI-scope-minimizing integration patterns.

Exit: at least two provider adapters demonstrate that commerce/finance do not contain provider-specific lifecycle logic.

---

## P11 — POS & Edge

Scope:

- POS application/runtime;
- local durable transaction queue;
- offline authorization policy;
- sync/reconciliation protocol;
- receipt/device abstractions;
- barcode/cash drawer/printer/scale/payment-terminal adapter model;
- shift/session management.

Exit: offline sale/reconnect/replay/conflict scenarios are deterministic and auditable.

---

## P12 — Experience Builder & CMS

Scope:

- site/page/section/component model;
- responsive design tokens/layout;
- data binding;
- CMS content types/fields/relations;
- drafts/versions/publishing/scheduling;
- localization/SEO;
- module-contributed blocks;
- forms/actions;
- theme/package contract.

Exit: module removal degrades contributed blocks safely; content is API-addressable and versioned.

---

## P13 — Portal Platform

Scope:

- portal runtime;
- customer/vendor/employee/partner portal profiles;
- capability-driven navigation;
- scoped authentication/authorization;
- self-service documents, requests and status surfaces.

Exit: one runtime hosts multiple portal personas without duplicating identity/security stacks.

---

## P14 — HR, Projects & Service Operations

Scope families:

- employees/workforce records;
- leave/attendance/time foundation;
- projects/tasks/timesheets;
- service tickets/cases;
- scheduling/resource allocation;
- payroll integration boundaries rather than premature country-specific payroll hard-coding.

Exit: domain boundaries between HR, project costing, service and finance are explicit and contract-tested.

---

## P15 — Supply Chain, Warehouse & Manufacturing

Scope:

- inventory ownership and stock movements;
- warehouses/bins;
- purchasing/suppliers;
- transfers;
- fulfillment/logistics abstractions;
- BOM/routings;
- production orders;
- quality/maintenance progression;
- traceability/lot/serial where required.

Exit: inventory invariants and financial/commerce integrations operate without cross-domain direct writes.

---

# PLATFORM EXPANSION PROGRAM

## P16 — Omnexa Connect / Integration Fabric

Deliver connector SDK, OAuth credential handling, webhooks, sync jobs, mapping/transformation, rate-limit handling and connector health.

Initial connector classes: email/calendar/storage, commerce channels, messaging, shipping, accounting, payments, AI providers and generic REST/webhook.

---

## P17 — Low-code App Builder

Deliver governed custom objects, forms, list/kanban/calendar views, relationships, permissions, workflows, dashboards and portal exposure.

Custom apps must use the same tenant/policy/audit/event foundations as first-party modules.

---

## P18 — Data, Reporting & BI

Deliver semantic metrics, dashboards, scheduled reports, cross-domain read projections, export, analytical pipeline/warehouse boundary and governed data access.

Heavy analytics must not degrade OLTP correctness/availability.

---

# INTELLIGENCE PROGRAM

## P19 — Omnexa Intelligence Platform

Deliver:

- model-provider abstraction;
- AI gateway;
- prompt/tool registry;
- retrieval/knowledge-base infrastructure;
- model/cost/rate policy;
- evaluation and traceability;
- tenant data boundaries;
- AI permissions;
- human-approval hooks.

Exit: AI can retrieve and propose actions across authorized context without unrestricted write access.

---

## P20 — Governed AI Agents

Deliver role/domain agents such as sales, finance, procurement, support and executive analysis through approved capabilities.

Exit: end-to-end agent actions are permissioned, approval-aware, auditable, replay-safe and measurable.

---

# ECOSYSTEM PROGRAM

## P21 — Developer Platform

Deliver public API surface, SDKs, CLI, sandbox tenant, module generator, contract/event explorer, local harness, developer docs and compatibility tooling.

---

## P22 — Omnexa Exchange / Marketplace

Deliver signed packages, publisher identity, declared scopes, automated validation, install/upgrade/revoke flow, ratings/release metadata and governance for modules/connectors/themes/workflows/AI tools/country packs.

---

# GLOBAL ENTERPRISE PROGRAM

## P23 — Globalization & Country Packs

Deliver multi-language/RTL, currency, locale/timezone, address formats, tax/fiscal adapters, country-pack contract, localized payments/documents and data-residency hooks.

Country-specific logic belongs in governed packs/adapters where practical rather than contaminating global core domains.

---

## P24 — Enterprise Governance, Security & Compliance

Expand SSO/SAML/SCIM, device/session policy, privileged access, policy administration, retention/legal hold hooks, audit export, security event center, key management integration, compliance evidence automation and enterprise admin controls.

Security starts earlier; this phase deepens enterprise-grade controls.

---

## P25 — Scale Fabric

Introduce only based on measured need:

- horizontal service scaling;
- read replicas;
- partitioning/sharding strategy;
- regional deployment/data planes;
- failover;
- multi-region control plane;
- disaster recovery automation;
- tenant placement/mobility;
- service extraction for justified domains.

Exit gates use measured SLO, RTO/RPO, load and fault-injection evidence.

---

# INDUSTRY / AUTONOMY PROGRAM

## P26 — Industry Packs

Compose existing platform modules into governed vertical solutions such as retail, restaurant, hospitality, healthcare, real estate, construction, logistics, manufacturing, education and professional services.

Industry packs must avoid forks of core modules.

---

## P27 — Autonomous Business OS

Target state:

```text
Business objective
 -> governed intelligence analysis
 -> evidence/recommendations
 -> policy/human approval
 -> multi-domain workflow execution
 -> measurement
 -> feedback/optimization
```

Examples may include profitability improvement, inventory optimization, purchasing recommendations, churn reduction or marketing reallocation.

The platform must preserve human control, policy limits, auditability and deterministic financial/business records even as autonomy increases.

---

# 3. Parallelism policy

Before P06 completion: foundation phases are primarily sequential.

After P06: P07-P15 may have controlled parallel teams only when:

- shared contracts are frozen enough for the planned work;
- each domain has an owner;
- no team writes another domain's schema;
- integration uses published capability/event/workflow contracts;
- CI can test module combinations.

P16-P18 depend on stable enough business capability contracts. P19-P20 depend on governed data and authorization foundations. P21-P22 depend on stable extension contracts. P23-P25 deepen global/enterprise maturity and may feed requirements back through change control.

# 4. Release philosophy

Do not wait for P27 to release useful products. Release coherent capability sets behind maturity labels:

- internal/dev preview;
- alpha;
- beta;
- production-ready;
- enterprise-ready.

A domain can mature independently while the overall roadmap continues, provided platform compatibility gates remain satisfied.

# 5. Quality principle

Progress percentage is based on accepted work packages, not lines of code, number of screens or subjective effort.

If tests or architecture evidence invalidate a completed package, its state must be reopened rather than preserving an inaccurate percentage.
