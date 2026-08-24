# Omnexa Strategic Cross-Cutting Programs

Status: **Proposed under ADR-0011**  
Execution impact: **none until accepted and explicitly activated**

## 1. Purpose

Omnexa already has a strong P00-P27 domain roadmap. This document adds cross-cutting platform systems that should mature across multiple existing phases without renumbering them, duplicating domain ownership or opening future implementation early.

These programs exist because Omnexa's target is not “another ERP.” The target is an AI-native Business Operating System that combines:

1. **System of Record** — deterministic authoritative business state;
2. **System of Workflow** — durable cross-domain execution;
3. **System of Context** — semantic Business Graph and governed data fabric;
4. **System of Intelligence** — AI, process, performance and decision intelligence;
5. **System of Governance** — security, GRC, AI control, evidence and policy;
6. **System of Simulation/Optimization** — scenario modelling before consequential autonomy;
7. **System of Work** — role/intent-driven unified operating surface.

## 2. Non-negotiable architecture rules

- Domain write authority remains with the owning domain.
- Business Graph, Process Graph, System Graph, analytics, AI and search are not alternate business write models.
- AI acts through the same versioned capabilities and policy boundaries as humans/workflows/integrations.
- Cross-tenant and cross-organization context remains explicit.
- Derived stores are rebuildable unless an accepted ADR makes them authoritative.
- High-risk financial/security/AI actions require evidence and independent review proportional to risk.
- Specialized graph/vector/stream/warehouse/runtime-isolation infrastructure is introduced only after measured need.
- Current P02.08 remains the sole active implementation package until canonical state explicitly changes.

## 3. Foundation wave — before broad business-domain parallelism

### XQ-100 — AI-Native Engineering Governance & Quality OS

Mature existing governance into a closed-loop engineering system:

`Signal/Research -> Problem -> VOC/CTQ -> Plan -> Architecture -> Data -> Security -> UX -> Implementation -> QA -> Performance -> Reliability -> Release -> Outcome -> Improve/Control`

Method selection:

- new/materially redesigned high-impact systems: proportional DMADV;
- defects/incidents/regressions: proportional DMAIC;
- high/critical complex failure paths: FMEA plus security threat modelling.

AI-development orchestration adds exact run identity, active-scope manifest, allowed paths/tools, policy digest, scope leases, concurrent-writer conflict protection, test-oracle integrity, evidence attestation, independent exact-head review, bounded retry/cost, dependency intake, and governance self-modification protection.

### XDG-100 — Enterprise Data Governance, Catalog, Lineage & Master Data

Adds the governance layer between domain-owned data and enterprise-wide consumption.

Required concepts:

- authoritative owner and source;
- data contracts and classifications;
- tenant/org/legal-entity scope;
- lineage and transformations;
- derived stores and rebuild semantics;
- retention/export/delete propagation;
- residency/localization rules;
- master/reference data;
- golden-record/entity-resolution policy;
- data-quality metrics;
- AI/search/analytics exposure policy;
- backup/recovery treatment.

Master Data Management must not become a giant shared write database. A “golden record” is a governed resolution/product of authoritative domain sources, not permission for a generic MDM service to overwrite private domain state.

### XBG-100 — Unified Business Graph & Semantic Layer

Turns the Product Constitution's strategic moat into an explicit platform program.

The Business Graph relates domain identities such as:

`Party <-> Customer <-> Contract <-> Order <-> Product <-> Location <-> Shipment <-> Invoice <-> Payment <-> Account`

It must support typed and temporal relationships, provenance, tenant/org scoping, semantic concepts, policy-aware traversal and versioned graph query contracts.

Rules:

- domain IDs remain canonical;
- graph facts identify source/provenance;
- graph projections can be rebuilt;
- graph traversal never bypasses object/field authorization;
- AI receives only authorized graph sub-context;
- graph storage technology remains implementation-neutral until measured.

### XSG-100 — System Graph & Flow Intelligence Foundation

Creates a machine-readable topology for Omnexa itself.

Node examples:

- module;
- capability/API;
- event;
- workflow;
- job;
- database/schema;
- cache;
- object store;
- external provider;
- permission/policy;
- secret reference;
- AI tool/agent/model;
- deployment/service;
- test/evidence artifact.

Edge examples:

- calls;
- reads;
- writes;
- emits;
- consumes;
- authorizes;
- requires capability;
- crosses trust boundary;
- sends to/receives from;
- retries;
- compensates;
- depends on;
- tested by;
- owned by.

Evidence classes remain separate: `declared`, `static`, `observed`, `tested`, `production-observed`, `modelled`, `ai-inferred`.

The GUI can later project architecture/runtime/data/security/error/state/transaction/retry/network/DB/cache/package/deployment/performance/test/cost/incident lenses, but the drawing is not the source of truth.

### XTRUST-100 — Module, Connector, Agent & Package Trust Runtime

P03 establishes contracts; P16/P21/P22/P24 mature enforcement.

Target controls:

`Package -> Quarantine -> Publisher/Signature -> Provenance/SBOM -> Manifest -> Capability/Data/Network/Secret Review -> Static/Security Tests -> Runtime Trust Tier -> Sandbox/Bounded Execution -> Health -> Enable -> Monitor -> Revoke/Kill`

Third-party code does not gain generic DB/filesystem/network/secrets because it is “installed.” Brokered access and declared scopes are the preferred model.

### XPERF-100 — Performance, Capacity & Runtime Attribution Foundation

Before many business domains run in parallel, Omnexa needs attributable performance evidence:

- request/capability/module latency;
- DB/query/cache timing;
- event/job/workflow delay;
- provider/network latency;
- memory/CPU/connections;
- tenant/noisy-neighbor attribution;
- queue age/backlog;
- performance budgets and regression baselines;
- load/capacity profiles.

This complements, not replaces, current SLO/reliability standards.

## 4. Business-control wave

### XPS-100 — Payment Security & Provider Certification Gate

P10 payment provider abstraction is necessary but generic module acceptance is not enough for financial-provider activation.

Default architecture:

- provider-hosted redirect/iframe/hosted-fields/tokenized flows;
- raw PAN/CVV/track/PIN excluded from standard Omnexa runtime/logs/cache/queues/backups/AI;
- purpose-specific payment capabilities;
- Core payment domain validates amount/currency/order/state;
- payment Secret Broker and Network Broker;
- signed/fresh/replay-safe webhooks;
- explicit async/3DS/SCA states;
- duplicate/concurrent mutation protection;
- ambiguous timeout -> reconcile before retry;
- protected payment-entry surfaces;
- provider sandbox certification;
- credential rotation, provider kill switch and reconciliation after compromise.

### XGRC-100 / XGRC-200 — SoD, Maker-Checker, GRC & Continuous Controls

Basic financial governance should arrive with Finance/ERP, not wait until enterprise hardening.

XGRC-100 at P08:

- segregation-of-duties conflicts;
- sensitive access;
- maker-checker;
- approval authority;
- financial role constraints;
- privilege escalation tests;
- period-control linkage.

XGRC-200 at P24:

- control library;
- risk registers;
- access reviews/role certification;
- continuous transaction/configuration monitoring;
- control testing;
- exceptions/waivers;
- compliance evidence automation;
- legal hold/retention linkage;
- risk/fraud indicators.

### XEPM-200 — Enterprise Planning & Performance Management

Omnexa should support planning as a first-class enterprise capability, not only budgeting fields inside Finance.

Includes:

- budgets/versions;
- rolling forecasts;
- driver-based plans;
- cash/workforce/sales/supply-chain/capex planning;
- profitability/cost allocation;
- consolidation/close support;
- variance analysis;
- scenario versions and what-if planning.

## 5. Intelligence wave

### XPROC-200 — Process Intelligence / Process Mining

P05 executes designed workflows. XPROC-200 explains what actually happened across workflows/events/manual actions.

Capabilities:

- process discovery;
- conformance vs designed process;
- path/variant analysis;
- waits/rework/loops;
- approval bottlenecks;
- SLA/SLO/business-cycle metrics;
- automation opportunity analysis;
- cross-domain process graph;
- process changes compared over time.

### XPERF-200 — Performance Intelligence Center

Builds on XPERF-100 and System Graph:

- root bottleneck attribution;
- version/package comparisons;
- query/cache/provider regressions;
- capacity forecasting;
- historical trends/alerts;
- tenant/module/provider cost-performance relationships;
- evidence-backed AI recommendations;
- scale-fabric triggers based on measured need.

### XAIC-200 — AI Control Tower

Central inventory and governance for all enterprise AI, including third-party/shadow assets where discoverable.

Asset types:

- models/providers;
- agents/agent teams;
- prompts/instructions;
- tools/capabilities;
- MCP servers/external agent endpoints;
- datasets/evaluation sets;
- knowledge sources;
- embeddings/vector collections;
- AI applications.

Each asset tracks owner, tenant/org scope, lineage, risk, data access, permissions, model/tool versions, evals, runtime traces, incidents, cost, value/ROI, approvals, lifecycle and kill/revoke state.

### XAIS-200 — Agent Studio, Multi-Agent Orchestration & Interoperability

Low-code/pro-code controlled agent builder with:

- typed tool selection;
- Business Graph/context selection;
- deterministic workflow steps plus non-deterministic agent steps;
- supervisor/specialist/agent-team patterns;
- human approvals;
- sandbox/evals;
- provider-neutral models;
- MCP and external-agent interoperability;
- version/publish/rollback;
- template/marketplace path later.

### XMLM-200 — AI Model Lifecycle / Evaluation / MLOps Governance

“AI-MLM” in Omnexa is defined as **AI Model Lifecycle Management** rather than a business domain.

Required capabilities:

- provider/model registry and routing;
- model cards and risk class;
- prompt/tool/model version lineage;
- datasets and eval-set governance;
- offline/online evaluations;
- regression and drift monitoring;
- latency/cost/quality trade-offs;
- fine-tuning/training approval where supported;
- red-team/adversarial suites;
- deprecation, fallback and rollback;
- no training on tenant-sensitive data without explicit policy/contract.

### XDOC-200 — Document Intelligence, Enterprise Knowledge & Semantic Search

Adds permission-aware knowledge ingestion and unstructured-to-structured intelligence:

- OCR/intelligent document processing;
- classification/extraction;
- document provenance/citations;
- enterprise semantic search;
- Business Graph linking;
- retention/legal-hold propagation;
- tenant/field authorization;
- RAG-ready retrieval;
- no unrestricted vector-store bypass around source permissions.

## 6. Delivery/operations wave

### XALM-200 — Configuration-as-Code & Environment ALM

Version, diff, validate and promote non-code platform configuration across environments:

- workflows;
- low-code apps/objects/forms;
- permissions/policies;
- dashboards/reports;
- agent definitions/prompts/tools;
- feature flags/settings;
- country packs;
- templates.

Flow:

`DEV -> diff/validate -> TEST -> UAT -> approval -> PROD`

Secrets remain references, never exported values.

### XMIG-200 — Migration & Data Onboarding Center

Enterprise adoption requires first-class migration from external systems.

Pipeline:

`Discover -> Profile -> Map -> Transform -> Validate -> Dry Run -> Reconcile -> Import -> Exception Queue -> Final Reconciliation`

Include idempotent resume, duplicate handling, mapping templates, incremental cutover, business/financial balancing and migration audit.

### XFIN-200 — FinOps & Resource Cost Intelligence

Attribute platform cost by tenant/module/capability/workflow/provider/agent where data supports it.

This enables capacity governance, pricing/unit economics and AI cost/value decisions without turning cost telemetry into user billing authority unless a separate billing domain explicitly owns that responsibility.

### XOUT-200 — Business Outcome & Value Intelligence

Tie product/process/agent changes to measurable goals rather than code volume.

Examples:

- order cycle time;
- DSO;
- invoice error rate;
- inventory turns;
- support resolution time;
- conversion;
- approval latency;
- workflow automation rate;
- AI task success;
- AI ROI;
- failure/rework reduction.

Outcomes must be observed, not fabricated by AI.

### XESM-200 — Enterprise Service Management, CMDB & Business Service Graph

P14 service cases are not equivalent to a ServiceNow-class enterprise-service platform.

Later capability may include:

- service catalog;
- incident/problem/change/request;
- CMDB/configuration items;
- business services;
- service dependency graph;
- SLA/SLO linkage;
- asset/service ownership;
- automation/agent integration;
- operational impact/blast-radius analysis.

## 7. Autonomy wave

### XDEC-300 — Decision Intelligence, Scenario Simulation & Business Digital Twin

Autonomous recommendations should have a simulation layer before high-consequence execution.

Inputs may include Business Graph, Process Graph, plans/forecasts, historical metrics, capacity, constraints and current state.

Outputs remain clearly classified as modelled/predicted until verified.

Examples:

- supplier price increase;
- warehouse outage;
- workforce shortage;
- promotion/pricing change;
- capacity change;
- cash-flow scenario;
- procurement strategy;
- inventory rebalancing.

### XWORK-200 — Unified Business Work & Command Center

Role/intention-driven workspace above module navigation:

- work queue;
- approvals;
- exceptions;
- alerts/KPIs;
- customer/vendor/employee context;
- cases/tasks;
- workflows;
- AI recommendations;
- natural-language search/action;
- explainable reason for every surfaced item.

### XAUTO-300 — Autonomous Business Governor

P27 autonomy receives an explicit control plane rather than “agents with more permissions.”

Required concepts:

- autonomy levels per capability;
- policy/risk/action budgets;
- approval thresholds;
- objective/constraint registry;
- simulation-before-action rules;
- multi-agent conflict resolution;
- kill/containment;
- reversible/compensating execution requirements;
- continuous outcome/control monitoring;
- automatic reduction of autonomy when quality/risk evidence degrades.

## 8. Future domain-gap candidates

These are not automatically approved phases; each needs research/ownership/dependency analysis before promotion:

- marketing automation / CDP / journey orchestration;
- contract lifecycle management and e-signature;
- deeper treasury/cash/risk management;
- strategic sourcing, supplier-risk and business-network capabilities;
- enterprise asset management / field service depth;
- subscription billing/revenue management/revenue recognition;
- fraud/financial-crime/risk engines for regulated industries;
- native collaboration/communications/knowledge productivity where owning the capability creates more strategic value than integrating existing providers.

## 9. Dependency principles

Do not implement all X-programs together. They are dependency overlays.

Recommended order:

`XQ-100 -> XDG-100 -> XBG-100`  
`P03/P04 -> XSG-100 + XTRUST-100 + XPERF-100`  
`P08 -> XGRC-100`  
`P10 -> XPS-100`  
`P18 -> XPROC-200 + XPERF-200 + XEPM-200 + XFIN-200 + XOUT-200`  
`P19/P20 -> XAIC-200 + XMLM-200 + XAIS-200`  
`P17/P21 -> XALM-200 + XMIG-200`  
`P24 -> XGRC-200 + mature XTRUST/XAIC`  
`P25 -> scale based on System/Performance evidence`  
`P27 -> XDEC-300 + XWORK-200 + XAUTO-300`

## 10. Acceptance philosophy

A strategic system is not complete because its UI exists.

Each program eventually needs, as applicable:

- architecture/data/security threat model;
- tenant/org negative tests;
- source + target evidence separation;
- migration/rebuild/rollback behavior;
- performance/reliability budgets;
- system-graph evidence contribution;
- independent review for high/critical paths;
- operational runbooks;
- business outcome/CTQ evidence;
- deprecation and recovery strategy.
