# Omnexa System Architecture

Status: **Architecture Baseline v1**

## 1. Architectural shape

Omnexa starts as a **strict modular monolith with service-ready boundaries**.

This is deliberate. The system should gain the organizational and contract advantages of service boundaries before paying the operational cost of a broad microservice topology.

```text
Clients
  ├── Admin Web
  ├── Customer/Vendor/Employee Portals
  ├── Website/Storefront
  ├── POS / Edge
  ├── Mobile
  └── External API Consumers
          │
          ▼
      API / Edge Layer
          │
          ▼
     Omnexa Platform
  ┌───────┼──────────────────────────────────────────────┐
  │ Kernel│ Domain Modules                               │
  │       │ CRM Finance Commerce Inventory HR ...       │
  └───────┼──────────────────────────────────────────────┘
          │
  ┌───────┼────────┬────────────┬───────────────┐
  ▼       ▼        ▼            ▼               ▼
Postgres Cache  Event Fabric  Object Store  Observability
```

## 2. Platform kernel

The kernel owns capabilities that must remain consistent across every domain:

- tenant and organization context;
- identity and session foundation;
- authorization/policy evaluation;
- module registry and lifecycle;
- configuration, settings and feature flags;
- secrets abstraction;
- IDs, timestamps, locale, currency and base value objects;
- files/media abstraction;
- event publishing/subscription infrastructure;
- job/scheduler infrastructure;
- workflow runtime abstraction;
- notifications foundation;
- audit and security-event foundation;
- API conventions and error format;
- health, telemetry and diagnostics;
- licensing/entitlement hooks;
- migration/version registry.

A business module must not create its own replacement for these concerns.

## 3. Domain module boundary

A domain module contains its own:

- commands/use cases;
- domain model;
- persistence adapters/schema ownership;
- API handlers/routes;
- events published and consumed;
- workflow actions/triggers;
- permissions/capabilities;
- settings and feature flags;
- UI contributions;
- jobs;
- migrations;
- tests;
- module manifest and lifecycle hooks.

Internal module code is private unless exported as a declared capability.

## 4. Dependency direction

Preferred dependency direction:

```text
UI / API Adapters
      ↓
Application Use Cases
      ↓
Domain Model
      ↓
Ports / Contracts
      ↓
Infrastructure Adapters
```

Cross-module dependency direction:

```text
Consumer Module
      ↓
Published Capability Contract / Event Contract
      ↓
Owning Module
```

Forbidden:

```text
Consumer Module
      ↓
Owning Module internal table / ORM model / private package
```

## 5. Communication patterns

Use the lightest pattern that preserves ownership.

### 5.1 Synchronous capability call
Use when the caller needs an immediate authoritative result.

Examples:

- get customer eligibility;
- authorize a payment operation;
- reserve inventory;
- calculate tax through the tax capability.

### 5.2 Domain event
Use when an owning domain announces a fact and consumers may react independently.

Examples:

- `commerce.order.created.v1`
- `payments.capture.succeeded.v1`
- `inventory.stock.changed.v1`

Events are immutable facts. Consumers may not reinterpret an event as permission to mutate the publisher's state directly.

### 5.3 Workflow orchestration
Use for long-running multi-step business processes needing waits, retries, compensation, human approval or durable state.

Example:

```text
Order accepted
 -> fraud/customer checks
 -> reserve inventory
 -> authorize/capture payment
 -> create finance document
 -> prepare fulfillment
 -> notify customer
```

### 5.4 Read model/projection
Use when a cross-domain screen needs efficient consolidated reads. Projections are disposable/rebuildable views, not an alternate write authority.

## 6. Event standard

All durable domain events require:

- globally unique event ID;
- event type and explicit version;
- occurred-at timestamp;
- producer module/version;
- tenant and relevant organization scope;
- actor/service identity when applicable;
- correlation ID;
- causation ID;
- idempotency/replay strategy;
- documented payload schema.

Event naming convention:

```text
<domain>.<entity-or-capability>.<past-tense-fact>.v<major>
```

Example: `commerce.order.created.v1`.

## 7. Request and identity context

Every protected request establishes a context containing at minimum:

- actor identity;
- tenant ID;
- organization/company scope when applicable;
- authenticated session/service identity;
- locale/timezone;
- correlation/trace ID;
- authorization context.

Background jobs and event consumers must reconstruct an explicit service/tenant context rather than assuming a global tenant.

## 8. Authorization architecture

Authorization is evaluated server-side at capability boundaries.

The model supports:

- roles for coarse responsibilities;
- relationships for object/scope membership;
- contextual conditions for policy constraints;
- service accounts for integrations/automation;
- explicit permission scopes for extensions and AI tools.

A front-end hidden button does not constitute access control.

## 9. Data architecture

### 9.1 Transactional store
PostgreSQL is the primary OLTP baseline.

Module schema ownership should be logically explicit, for example:

```text
kernel_*
identity_*
crm_*
finance_*
commerce_*
inventory_*
```

Exact physical schema/database partitioning may evolve, but ownership must remain clear.

### 9.2 Cache/ephemeral state
Redis-compatible infrastructure may be used for cache, short-lived coordination and rate-limiting primitives. The primary source of business truth must not live only in cache.

### 9.3 Object storage
Large files/media use S3-compatible object storage behind the platform file capability.

### 9.4 Analytics
Operational reporting may begin with carefully designed OLTP queries/read projections. A dedicated analytics/warehouse path is introduced before large analytical workloads are allowed to degrade transactional workloads.

## 10. Transaction boundary

Transactions should normally remain inside one module's owned write model.

Cross-domain atomicity should not be simulated with hidden distributed database transactions. Use workflows, sagas/compensating actions, idempotency and explicit state transitions.

## 11. Module lifecycle

Supported lifecycle:

```text
available -> installing -> enabled -> active
                         -> suspended
                         -> disabled
                         -> archived/exported
                         -> detached
                         -> purged (explicit destructive operation)
```

Disable/remove behavior must respect historical references. A module's absence must not erase required business evidence owned by other domains.

## 12. Compatibility model

Public contracts require semantic compatibility discipline:

- additive fields should be backward compatible where possible;
- breaking API/event changes require a new major contract version;
- module upgrades require migration compatibility notes;
- consumers must not depend on undocumented fields or ordering;
- deprecation must be observable and time-bounded before removal.

## 13. Technology topology

Baseline implementation:

```text
Backend: Go
Web/Admin/Builder: TypeScript + React
Edge/native where justified: Rust
AI/data workers where justified: Python
OLTP: PostgreSQL
Cache: Redis-compatible
Messaging: NATS/JetStream-class fabric
Files: S3-compatible object store
Observability: OpenTelemetry
```

Infrastructure products may be substituted only if the architectural capability and operational requirements remain satisfied and change control is followed.

## 14. Service extraction criteria

A module/service may be extracted from the modular monolith when one or more are demonstrated:

- materially different scaling profile;
- failure isolation requirement;
- independent deployment cadence owned by a team;
- regulatory/data-residency boundary;
- CPU/memory/runtime characteristics incompatible with the main process;
- availability target materially different from the platform baseline.

Extraction must preserve capability/event contracts wherever practical.

## 15. POS and edge architecture

POS/edge cannot assume permanent connectivity.

Target pattern:

```text
Omnexa Cloud
    ↕ secure sync
Omnexa Edge Runtime
    ├── local durable state
    ├── transaction/outbox queue
    ├── conflict/reconciliation policy
    └── device adapters
        ├── printer
        ├── barcode scanner
        ├── cash drawer
        ├── scale
        └── payment terminal
```

Offline permissions, settlement, conflict handling and reconciliation must be explicit rather than eventual side effects.

## 16. Payment architecture

Payment providers integrate behind an Omnexa payment orchestration contract.

Common capability vocabulary includes:

- authorize;
- capture;
- void;
- refund/partial refund;
- tokenize/provider reference;
- recurring/mandate reference;
- payout;
- settlement;
- dispute;
- reconciliation;
- webhook verification.

Business modules must not hard-code a provider-specific payment lifecycle into their core model.

## 17. Experience builder architecture

The website/CMS builder is schema/component based, not hard-wired to one business module.

Conceptual hierarchy:

```text
Site -> Page -> Section -> Component -> Element
```

Modules contribute versioned UI/data capabilities to the builder. If a module is unavailable, builder behavior must degrade safely and visibly.

## 18. AI architecture

AI sits above governed platform capabilities:

```text
Model/Agent
   ↓
AI Gateway / Tool Registry
   ↓
Identity + Tenant Context
   ↓
Policy / Approval
   ↓
Versioned Omnexa Capability
   ↓
Owning Domain
   ↓
Audit + Event
```

AI cannot bypass domain validation or become a second source of business truth.

## 19. Observability

Every major request/workflow should be traceable by correlation/trace ID across:

- gateway;
- application use case;
- database operations;
- events/jobs;
- workflow transitions;
- external integrations.

Metrics, logs and traces should follow OpenTelemetry-compatible semantics.

## 20. Resilience principles

- explicit timeouts;
- bounded retries with backoff;
- idempotency keys for replayable mutations;
- dead-letter/quarantine for poison events;
- outbox/inbox patterns where delivery guarantees require them;
- health/readiness separation;
- graceful degradation for optional dependencies;
- backup and restore tests, not backup configuration alone.

## 21. Architecture verification

The architecture is considered preserved when automated tests can increasingly prove:

- module boundary rules;
- dependency rules;
- tenant isolation;
- permission enforcement;
- event schema/version validation;
- migration safety;
- module enable/disable behavior;
- fresh-install reproducibility;
- contract compatibility.
