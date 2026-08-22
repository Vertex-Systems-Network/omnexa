# Omnexa Module & Submodule Execution Blueprint

Status: **Mandatory planning baseline**  
Purpose: pre-plan module/submodule execution deeply enough that implementation agents execute an approved decomposition instead of repeatedly redesigning the work.

## 1. Why this exists

Omnexa contains large product domains that naturally decompose into submodules such as page builder, template builder, theme system, catalog, checkout, reporting, workflows and many others. Re-planning these boundaries during implementation creates architecture drift, duplicate ownership and inconsistent quality.

This blueprint therefore defines the mandatory planning shape for every future module and submodule before executable work begins.

It is **planning only**. It does not activate a future roadmap phase, bypass `STATE.json`, or authorize implementation outside the current active package.

## 2. Runner-deferred delivery workflow

For an already-authorized work package, use this sequence:

1. read canonical state, architecture, ownership and the active package spec;
2. complete the implementation decomposition and all required subtasks on the feature branch;
3. implement source, tests, docs, migrations/contracts and deterministic verification scripts;
4. perform all available static/unit/self-review preparation before opening the final executable PR;
5. open the PR only when the branch is implementation-complete enough for canonical verification;
6. run the required GitHub-hosted governance lane at the end as the authoritative integration gate;
7. fix any defects found by the hosted lane without weakening tests or policy;
8. merge only after the required hosted check is green;
9. create immutable completion evidence and perform the separate state/ledger transition.

This permits runner work to be deferred until the end of implementation while preserving the rule that no package becomes `done`, and no protected merge occurs, without canonical hosted evidence.

## 3. Mandatory hierarchy

Every business/platform module must be planned using the following hierarchy:

```text
Phase
  -> Module / domain
      -> Submodule / capability family
          -> Work package
              -> Task
                  -> Test/evidence item
```

A submodule is not an independent architecture by default. It remains inside its owning module unless an approved ADR explicitly creates another product boundary.

## 4. Required submodule plan

Before coding a submodule, its plan must contain all of the following:

- stable identifier and human-readable name;
- owning phase and owning module/domain;
- purpose and user/business outcome;
- authoritative write owner;
- dependencies: platform, required module, optional module and forbidden dependencies;
- data model/schema ownership;
- capabilities provided and consumed;
- API/event/workflow/UI contribution contracts;
- permissions and authorization boundary;
- tenant/organization scope where applicable;
- settings and feature flags;
- lifecycle behavior: install/enable/disable/upgrade/archive/purge where applicable;
- migration/fresh-install/upgrade strategy;
- import/export implications;
- files/storage/cache/search/reporting implications;
- security/data-classification/secrets considerations;
- accessibility/localization/RTL requirements for UI surfaces;
- failure/degradation/idempotency/retry behavior;
- observability/audit requirements;
- explicit in-scope and out-of-scope boundaries;
- implementation work packages and ordered subtasks;
- acceptance criteria and gate mapping;
- Definition-of-Done evidence requirements.

If any of these are unknown, they are a planning defect to resolve before implementation rather than a reason to invent semantics while coding.

## 5. Standard ordered implementation subtasks

Unless a submodule-specific plan justifies a different order, AI agents should execute these subtasks in sequence:

### S01 — Ownership and contracts

- confirm owner and write boundary;
- define capability identifiers and request/response semantics;
- define API/event/workflow/UI contribution contracts that are actually required;
- define permissions and policy checks;
- define dependency degradation behavior.

### S02 — Data and migrations

- define owned entities/value objects;
- define persistence schema and indexes;
- create versioned migrations under owner scope;
- implement fresh-install and supported-upgrade path;
- define immutable snapshots where audit/business continuity requires them.

### S03 — Domain/application behavior

- implement invariants;
- implement commands/use cases;
- implement validation and structured failures;
- implement idempotency/retry semantics where execution may repeat;
- keep provider-specific behavior behind adapters.

### S04 — Provider/integration adapters

- implement repository/provider adapters;
- set bounded timeout/cancellation behavior;
- classify unavailable/degraded states;
- protect secrets and sensitive content;
- add contract tests against synthetic providers where required.

### S05 — Permissions, tenancy and security

- server-side permission enforcement;
- tenant/organization scoping;
- cross-tenant negative tests;
- security/data-classification checks;
- privileged operation audit requirements;
- SSRF/path traversal/upload/download/secret handling tests where applicable.

### S06 — UI and UX contributions

Only when the active package authorizes UI:

- navigation/route/slot declaration;
- loading/empty/error/success/disabled states;
- responsive layout and design-token use;
- keyboard/focus semantics;
- semantic HTML and accessible labels/names;
- WCAG 2.2 AA target;
- W3C validation;
- WAVE evaluation;
- manual keyboard/focus/screen-reader/zoom/reflow evidence;
- localization/RTL and reduced-motion behavior where applicable.

### S07 — Lifecycle and resilience

- install/enable/disable/re-enable;
- upgrade and compatibility;
- provider outage/restart/degradation;
- retry/replay/idempotency;
- archive/detach/purge semantics if owned by the module;
- prove unrelated modules remain unaffected.

### S08 — Tests and quality

- domain/unit tests;
- application/use-case tests;
- persistence/provider integration tests;
- capability/API/event contract tests;
- permission and tenancy tests;
- migration/fresh-upgrade tests;
- lifecycle tests;
- accessibility/browser tests for UI;
- lint/static/security/build/package gates.

### S09 — Documentation and AI handoff

- update module/submodule plan state;
- document configuration and developer commands;
- record exact contracts and ownership;
- update AI read order if a new controlling document exists;
- record remaining known follow-ups explicitly rather than leaving implicit TODOs.

### S10 — Final hosted verification and closure

- final diff/scope audit;
- PR review/thread audit;
- GitHub-hosted canonical governance run;
- fix failures without bypassing checks;
- merge exact verified head;
- immutable completion evidence;
- state/status/ledger reconciliation;
- activate only the next authorized package.

## 6. Submodule package sizing rule

A submodule must be split into multiple work packages when any of the following is true:

- it owns more than one independent write model;
- it exposes multiple independently versioned public capability groups;
- it contains both a runtime engine and a visual/admin builder;
- it has separately deployable/provider-specific adapters;
- migration risk is high enough to require isolated evidence;
- security/permission surface is materially different across parts;
- UI breadth is large enough that accessibility/e2e evidence would become unreviewable in one PR.

Do not split merely because files are numerous. Split by capability, ownership, risk and independently verifiable behavior.

## 7. Builder-family architecture rule

Builder products must separate **runtime contracts** from **authoring UI**.

Typical layers:

```text
Schema/model
  -> validation/versioning
  -> runtime renderer/executor
  -> reusable registry (blocks/components/templates/actions)
  -> authoring APIs
  -> visual builder UI
  -> preview/simulation
  -> publish/deploy lifecycle
```

The visual editor must not become the source of truth for undocumented runtime behavior. Saved definitions are versioned contracts validated server-side.

## 8. P12 Experience Builder & CMS — preplanned submodules

P12 must be decomposed before coding into at least the following submodules. Exact work-package IDs are assigned when P12 is activated, but the capability boundaries below are pre-approved planning targets.

### P12.A — Site & channel foundation

Subtasks:

1. site/channel identity and lifecycle;
2. domain/locale/environment references without embedding provider-specific DNS logic;
3. site settings and design-token binding;
4. navigation/menu contract;
5. preview versus published channel state;
6. capability/permission model;
7. lifecycle and version compatibility tests.

### P12.B — Page Builder

Subtasks:

1. page identity, slug/path and hierarchy model;
2. page draft/version/published state machine;
3. section/container/grid/layout schema;
4. block/component instance tree;
5. responsive breakpoint/layout constraints;
6. design-token/style binding;
7. data-binding references through governed capabilities;
8. reusable sections/symbols where approved;
9. undo/redo operation model for authoring UI;
10. preview rendering contract;
11. server-side schema validation;
12. permissions for view/edit/publish;
13. localization/RTL behavior;
14. SEO metadata integration;
15. accessibility semantics and authoring warnings;
16. version diff/restore;
17. publish/schedule integration;
18. page import/export package format;
19. module-contributed block degradation;
20. unit/contract/e2e/accessibility/lifecycle evidence.

### P12.C — Template Builder

Subtasks:

1. template identity/type/version model;
2. page/layout/section/email/document template categories only when their owning phases authorize them;
3. template slots and required regions;
4. variable/data-binding schema;
5. template inheritance/composition rules;
6. compatibility with registered components/blocks;
7. preview/test-data mechanism;
8. validation before activation;
9. template clone/version/restore;
10. assignment rules by site/content/context;
11. dependency checks during module disable;
12. import/export packaging;
13. permissions and publishing lifecycle;
14. contract/e2e/accessibility evidence.

### P12.D — Theme & design system

Subtasks:

1. theme package manifest;
2. design tokens: typography, spacing, radius, color, elevation, motion;
3. responsive tokens/breakpoints;
4. global styles with safe scoped overrides;
5. component style variants;
6. light/dark and brand variants where product requirements allow;
7. theme activation/rollback;
8. theme dependency/version compatibility;
9. accessibility contrast/token validation;
10. import/export/marketplace-ready package metadata.

### P12.E — Component & block registry

Subtasks:

1. block manifest/schema;
2. property editor schema;
3. allowed nesting/parent-child rules;
4. slots and composition;
5. capability/data requirements;
6. permission/feature/entitlement conditions;
7. server-side validation;
8. runtime renderer registry;
9. editor representation;
10. graceful fallback when contributing module is absent;
11. deprecation/version migration;
12. performance and accessibility requirements.

### P12.F — CMS content model builder

Subtasks:

1. content-type identity/versioning;
2. field-type registry;
3. validation rules;
4. relationships/references;
5. taxonomy/category primitives where owned;
6. draft/version/publish lifecycle;
7. localization;
8. permissions;
9. API-addressable content contract;
10. migration strategy for schema changes;
11. import/export;
12. search/index hooks;
13. audit/version history;
14. content editor generation.

### P12.G — Publishing, revisions & scheduling

Subtasks:

1. immutable revision identity;
2. draft/review/approved/published/archived states as authorized;
3. optimistic concurrency/conflict handling;
4. scheduled publish/unpublish;
5. atomic publication set where required;
6. rollback/restore;
7. preview tokens without authorization bypass;
8. cache invalidation hooks through governed contracts;
9. event/audit emission;
10. failure/retry/idempotency evidence.

### P12.H — Forms & actions

Subtasks:

1. form schema/field registry;
2. validation and accessible error contract;
3. spam/abuse hooks;
4. submission storage ownership decision;
5. workflow/action integration;
6. file-upload integration through governed storage capability;
7. consent/privacy metadata;
8. permission and anti-CSRF/session requirements from owning platform contracts;
9. notification/action adapter hooks;
10. submission lifecycle tests.

### P12.I — SEO, localization & discoverability

Subtasks:

1. title/meta/canonical schema;
2. robots/indexing controls;
3. sitemap generation contract;
4. structured-data extension points;
5. hreflang/locale routing;
6. slug translation and redirect rules;
7. social preview metadata;
8. accessibility-language metadata;
9. validation and duplicate-conflict tests.

### P12.J — Media integration layer

This is a CMS integration over kernel/global file capabilities, **not** ownership of object storage.

Subtasks:

1. media reference model;
2. metadata/alt-text/caption model;
3. selection/browse UX;
4. upload workflow through authorized storage/media capability;
5. permitted file type/size policies;
6. image rendition/thumbnail contract only when explicitly authorized;
7. orphan/reference lifecycle;
8. permissions;
9. localization/accessibility metadata;
10. provider independence.

### P12.K — Builder authoring shell

Subtasks:

1. editor shell layout;
2. canvas/iframe isolation decision;
3. drag/drop and keyboard-equivalent operations;
4. navigator/layers tree;
5. property inspector;
6. responsive preview modes;
7. command palette/history;
8. autosave and conflict indication;
9. error boundary/recovery;
10. loading/empty/offline/degraded states;
11. W3C/WAVE/manual accessibility gates;
12. performance budgets for large documents.

### P12.L — Package/import/export compatibility

Subtasks:

1. portable site/page/template/theme package manifest;
2. schema versions;
3. dependency declarations;
4. asset/reference remapping;
5. collision strategy;
6. security validation of imported packages;
7. dry-run validation;
8. deterministic import report;
9. rollback on failed import;
10. future marketplace signing hooks without implementing marketplace early.

## 9. Other roadmap module pre-decomposition

The following capability families must be treated as distinct submodule planning units when their phases become active.

### P07 CRM

- Leads;
- Contacts & organizations;
- Opportunities & pipelines;
- Activities/tasks/notes integration;
- Sales territories & ownership;
- Quotes/proposals integration boundary;
- Customer 360/read projection;
- Scoring;
- Sequences/automation.

### P08 Finance & ERP

- Chart of accounts;
- General ledger;
- Journal engine;
- Accounts receivable;
- Accounts payable;
- Invoicing & credit notes;
- Expenses;
- Cash/bank/reconciliation;
- Tax integration boundary;
- Budgeting;
- Fixed assets;
- Procurement/accounting integration.

### P09 Commerce

- Catalog/products;
- Variants/options;
- Collections/categories;
- Price lists;
- Promotions/discounts;
- Cart;
- Checkout;
- Orders;
- Returns/refunds orchestration;
- Subscriptions;
- Bundles;
- B2B/wholesale;
- Marketplace/multi-vendor boundary.

### P10 Payments

- Provider registry;
- Payment intents/authorization;
- Capture/void;
- Refunds;
- Token/reference vault boundary;
- Webhooks;
- Recurring mandates;
- Settlement;
- Reconciliation;
- Disputes/payouts.

### P11 POS

- POS runtime;
- Device abstraction;
- Local transaction queue;
- Offline policy;
- Sync/reconciliation;
- Shift/session;
- Receipt;
- Barcode/scanner;
- cash drawer/printer/scale adapters;
- terminal/payment adapter integration.

### P13 Portals

- Portal runtime;
- Portal profile builder/configuration;
- Customer portal;
- Vendor portal;
- Employee portal;
- Partner portal;
- Capability-driven navigation;
- Scoped document/request/status surfaces.

### P14 HR, Projects & Service

- Employee/workforce records;
- Leave;
- Attendance/time;
- Projects;
- Tasks;
- Timesheets;
- Resource scheduling;
- Service tickets/cases;
- Service scheduling;
- Payroll integration boundary.

### P15 Supply Chain & Manufacturing

- Inventory ledger/movements;
- Warehouses/bins;
- Purchasing;
- Suppliers;
- Transfers;
- Fulfillment;
- Logistics/shipping boundary;
- BOM;
- Routings;
- Production orders;
- Quality;
- Maintenance;
- lot/serial traceability.

### P16 Integration Fabric

- Connector SDK;
- Credential/OAuth management;
- Webhook ingress/egress;
- Sync engine;
- Mapping/transformation;
- Rate-limit/retry engine;
- Connector health;
- Provider-specific connectors as separate adapters.

### P17 Low-code App Builder

- Custom object builder;
- Field/schema builder;
- Form builder;
- List view builder;
- Kanban builder;
- Calendar view builder;
- Relationship builder;
- Permission builder;
- Workflow integration;
- Dashboard builder;
- Portal exposure;
- App package/version lifecycle.

### P18 Data, Reporting & BI

- Semantic metric registry;
- Dataset/read-model contracts;
- Report builder;
- Dashboard builder;
- Visualization registry;
- Scheduled reports;
- Export;
- Cross-domain read projections;
- analytical pipeline/warehouse boundary;
- governed data-access policy.

### P19 Intelligence Platform

- Model/provider registry;
- AI gateway;
- Prompt registry;
- Tool/capability registry;
- Knowledge sources;
- Retrieval/indexing;
- Model/cost/rate policy;
- Evaluation;
- traces/audit;
- AI authorization/tenant context;
- human approval hooks.

### P20 Governed Agents

- Agent definition registry;
- role/domain agent packages;
- planning/reasoning policy boundary;
- capability/tool execution;
- approval workflow;
- memory/context policy;
- evaluation/quality;
- audit/replay;
- cost/rate controls.

### P21 Developer Platform

- Public API surface;
- SDK generation;
- CLI;
- sandbox tenant;
- module generator;
- contract explorer;
- event explorer;
- local test harness;
- developer documentation;
- compatibility tooling.

### P22 Marketplace

- Publisher identity;
- package signing;
- package manifest validation;
- scopes/consent;
- listing/catalog;
- install;
- upgrade;
- revoke;
- ratings/release metadata;
- module/connector/theme/workflow/AI-tool/country-pack package types.

### P23 Globalization

- Language/translation;
- RTL;
- currency/formatting;
- timezone;
- address formats;
- tax/fiscal adapters;
- country-pack runtime;
- localized payment/document adapters;
- data-residency hooks.

### P24 Enterprise Governance

- SSO/SAML;
- SCIM;
- device/session policy;
- privileged access;
- policy administration;
- retention/legal hold hooks;
- audit export;
- security event center;
- key-management integration;
- compliance evidence automation.

### P25 Scale Fabric

- horizontal scaling;
- read replicas;
- partitioning;
- sharding policy;
- regional data planes;
- failover;
- multi-region control plane;
- DR automation;
- tenant placement/mobility;
- justified service extraction.

### P26 Industry Packs

Each vertical pack is a composition plan, not a fork. Pre-plan separate packs for retail, restaurant, hospitality, healthcare, real estate, construction, logistics, manufacturing, education and professional services using existing module capabilities.

### P27 Autonomous Business OS

- Objective definition;
- governed analysis;
- recommendation generation;
- policy/approval;
- multi-domain workflow execution;
- measurement;
- feedback/optimization;
- safety/audit/control plane.

## 10. Planning completeness rule

A future AI agent must not respond to an activated module with a fresh generic "first we need to plan" cycle when this blueprint and the owning phase plan already define the decomposition. The agent should:

1. load the canonical module/submodule plan;
2. select the next incomplete work package/subtask;
3. verify dependencies and active authorization;
4. execute it;
5. record evidence and remaining blockers.

Replanning is required only when implementation uncovers a genuine architecture conflict, missing owner, changed requirement or change-control event.

## 11. No premature implementation

This blueprint pre-plans future scope but does not unlock it. `docs/roadmap/STATE.json` remains authoritative for what may be implemented now.
