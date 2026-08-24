# Omnexa AI-Native Business OS Architecture

Status: **Proposed strategic architecture under ADR-0011**

## 1. Target system category

Omnexa targets an **AI-native composable Business Operating System**, not a monolithic ERP with an AI chatbot added later.

The architecture should eventually combine:

```text
Humans / Devices / External Systems / AI Agents
                    ↓
            Omnexa Work / Portals / APIs
                    ↓
        Identity + Tenant + Policy + Approval
                    ↓
      Governed AI / Workflow / Capability Plane
                    ↓
 ┌──────────────────┼───────────────────┐
 │                  │                   │
Business Graph   Process Graph      System Graph
 │                  │                   │
 └──────────────────┼───────────────────┘
                    ↓
          Versioned Domain Capabilities
                    ↓
 CRM | Finance | Commerce | Payments | Inventory | HR | ...
                    ↓
              Universal Kernel
                    ↓
 Data / Event / Storage / Cache / Telemetry / Deployment
```

Horizontal controls span every layer:

- security/tenancy;
- data governance/classification;
- GRC/SoD;
- audit;
- quality/evidence;
- reliability;
- performance/capacity;
- package/supply-chain trust;
- cost/FinOps;
- AI/model governance;
- outcome/control loops.

## 2. Three-graph model

### 2.1 Business Graph

Answers: **what business entities exist and how are they related?**

Examples:

- Person works for Organization;
- Customer placed Order;
- Order contains Product;
- Order produced Invoice;
- Invoice settled by Payment;
- Product stored at Warehouse;
- Supplier supplies Product;
- Employee belongs to Team and Cost Center;
- Asset located at Site and maintained under Contract.

Properties:

- semantic and business-oriented;
- typed/temporal relationships;
- source/provenance recorded;
- permission-aware traversal;
- tenant/org/legal-entity scope preserved;
- derived from authoritative domain state/events/contracts;
- not a generic cross-domain write authority.

### 2.2 Process Graph

Answers: **what actually happened in a business process, in what order, with which waits/loops/approvals/retries?**

Sources may include:

- P05 workflow history;
- P04 events;
- domain activities;
- human approvals;
- tasks/cases;
- integration/provider transitions.

Uses:

- process mining;
- conformance checking;
- rework/loop detection;
- bottleneck/wait analysis;
- automation opportunities;
- before/after process comparisons;
- business outcome linkage.

### 2.3 System Graph

Answers: **what software, runtime, data, permission and infrastructure path produced the behavior?**

Examples:

`API -> commerce.create-order -> Inventory capability -> Payment workflow -> Provider adapter -> Event -> Finance projection`

System Graph must capture evidence/provenance and keep declared/static/runtime/test/modelled information separate.

## 3. Canonical evidence model

Any intelligent recommendation or graph path must distinguish:

- **declared** — contract/manifest/architecture says the relationship exists;
- **static** — source/build analysis infers it;
- **observed** — captured runtime execution;
- **tested** — a controlled test exercised it;
- **production-observed** — production telemetry observed it;
- **modelled** — scenario/simulation predicts it;
- **ai-inferred** — AI hypothesizes it from available evidence.

Promotion rules:

- `ai-inferred` never becomes observed automatically;
- `modelled` never becomes actual business outcome automatically;
- one runtime trace does not prove all branches/concurrency behavior;
- test success does not prove production outcome;
- production telemetry does not grant business write authority.

## 4. AI-native control planes

### 4.1 AI Gateway / Model Plane

Owns:

- provider abstraction;
- model routing/fallback;
- quotas/rate/cost policy;
- model identity/version;
- request classification/redaction;
- approved context assembly;
- model/tool trace correlation.

### 4.2 Tool / Capability Plane

AI tools are typed adapters over versioned Omnexa capabilities.

Every tool declares:

- owner;
- input/output schema;
- allowed principal types;
- tenant/org scope;
- data classification;
- permission/policy;
- risk tier;
- approval requirement;
- idempotency/concurrency behavior;
- audit/event effect;
- retry/reconciliation behavior;
- rollback/compensation where meaningful.

No tool should be a disguised generic SQL/shell/filesystem/network superpower.

### 4.3 Agent Plane

Agents may be:

- assistant/read-only;
- recommendation/draft;
- bounded executor;
- workflow participant;
- supervisor/coordinator;
- domain specialist;
- multi-agent team member.

Agents inherit authority from authenticated user/service/agent identity plus current policy. Agent identity does not imply blanket domain authority.

### 4.4 AI Control Tower

The future Control Tower is the governance/observability plane for all AI assets rather than a model chat UI.

It tracks lifecycle, owner, risk, lineage, data, permissions, evals, cost/value, runtime behavior, incidents and kill/revocation state.

## 5. AI architecture roles

The following are **engineering roles/capability profiles**, not unrestricted personalities.

### AI Architect

Purpose: preserve global architecture and evolve it through evidence-backed decisions.

Responsibilities:

- classify new capability as Kernel/domain/platform/package/external;
- enforce ownership/dependency direction;
- evaluate alternatives/trade-offs;
- author/review ADRs;
- model trust/data/runtime boundaries;
- identify compatibility/migration/reliability impact;
- analyze change blast radius using System Graph when available.

Cannot:

- silently rewrite frozen architecture;
- approve its own high-risk architecture change as final authority;
- use AI preference as architectural evidence.

### AI Design Engineer

Purpose: produce professional, consistent, accessible and workflow-efficient UX across Admin, portals, builder, POS/mobile/edge and future unified Work.

Responsibilities:

- design-system/token/component governance;
- responsive/accessibility/i18n/RTL;
- workflow/task ergonomics;
- dense enterprise-data presentation;
- error/empty/loading/offline states;
- permission-aware UI;
- AI interaction/approval/explanation surfaces;
- generated UI validation against design contracts.

Design must never imply authorization. Hidden/disabled controls do not replace server policy.

### AI Performance Engineer

Purpose: optimize latency, throughput, resource efficiency and capacity without changing semantics or weakening safeguards.

Responsibilities:

- establish baselines/budgets;
- attribute cost to module/capability/query/cache/event/provider/tenant;
- detect regressions;
- analyze concurrency/backpressure;
- recommend measured cache/index/query/runtime improvements;
- evaluate scale triggers;
- correlate performance with cost and reliability.

Cannot trade correctness/security/audit for benchmark score without explicit architecture decision.

### AI Systems Engineer

Purpose: understand and improve whole-system composition.

Responsibilities:

- dependency and topology analysis;
- System Graph generation/query;
- distributed state/event/workflow interactions;
- deployment and failure-domain analysis;
- module lifecycle/removal impact;
- integration and edge/POS synchronization analysis;
- resilience and blast-radius modelling.

### AI Analyzer

Purpose: evidence-backed diagnosis across code, runtime, security, process, finance and business outcomes.

Analysis modes:

- static/code analysis;
- security source-to-sink analysis;
- runtime root-cause/cascade analysis;
- process conformance/bottleneck analysis;
- financial invariant/anomaly analysis;
- data-quality/lineage analysis;
- release regression analysis;
- agent/model behavior analysis;
- business-outcome analysis.

Must label hypothesis/confidence and never invent causality.

### AI Developer

Purpose: implement approved, bounded work through repository/public architecture contracts.

Required controls:

- exact active work package;
- run/base SHA;
- allowed paths;
- dependency/tool permissions;
- architecture/security/data/performance applicability;
- tests/migrations/rollback;
- no self-certification;
- scope-delta stop/re-plan.

### AI Domain Expert

Purpose: encode and review domain knowledge (finance, procurement, inventory, HR, tax, manufacturing, service, CRM, legal/compliance etc.) without bypassing product/domain ownership.

Expert output can include:

- rules/constraints;
- recommendation;
- policy interpretation;
- test scenarios;
- workflow design;
- risk analysis;
- knowledge-pack content.

Domain experts cannot mutate authoritative state directly; actions still invoke owning capabilities.

### AI-MLM — AI Model Lifecycle Manager

In Omnexa, `AI-MLM` means **AI Model Lifecycle Management**.

Responsibilities:

- model/provider inventory;
- model version/routing policy;
- evaluation datasets;
- model/prompt/tool lineage;
- quality/groundedness/task-success evaluation;
- latency/cost monitoring;
- drift/change detection;
- fine-tune/training governance;
- red-team testing;
- deprecation/fallback/rollback.

### AI Engineer

Purpose: convert architecture into production-grade systems and operating controls.

Responsibilities:

- platform/infrastructure engineering;
- event/job/workflow reliability;
- observability;
- security engineering;
- CI/CD/release engineering;
- data/AI runtime integration;
- capacity/DR;
- automated verification.

### AI Constructor

Purpose: turn an approved specification into coherent generated structure without inventing product authority.

Potential outputs:

- module skeletons;
- manifests;
- API/event schemas;
- permissions;
- migrations;
- CRUD/application ports;
- workflow definitions;
- UI scaffolds;
- tests/fixtures;
- docs/SDK samples;
- graph/evidence declarations.

Constructor output is generated implementation input, not accepted architecture or verified behavior.

## 6. Additional mandatory independent roles

For serious AI-native development, the requested roles are not sufficient alone. Add:

### AI Security / Red-Team Reviewer

Threat models, abuse cases, prompt/tool injection, tenant escape, IDOR, SSRF, supply chain, secrets, payment, unsafe agent delegation and evidence forgery.

### AI Data / Knowledge Engineer

Data contracts, lineage, MDM, semantic Business Graph, knowledge ingestion, retrieval authorization, derived stores and data quality.

### AI QA / Verifier

Owns independent test/eval review and prevents a code-authoring agent from weakening the oracle to obtain green evidence.

### AI Release / Operations Engineer

Binds reviewed/tested source to built/promoted artifact, verifies deployment/rollback/health/SLO evidence and handles incident/recovery controls.

### AI Product / Outcome Analyst

Connects Research/VOC/CTQs to measurable business outcomes; prevents “shipped” from being confused with “successful.”

## 7. AI-development orchestration

Future autonomous/semi-autonomous development should use a run manifest concept:

```text
run_id
base_sha
active_phase/work_package
policy_digest
plan_digest
allowed_paths
allowed_commands/tools
network_policy
secret_policy
dependency_policy
target_mutation_policy
required_reviews
risk_tier
evidence_authorities
attempt_budget
cost_budget
```

### Scope delta gate

If implementation discovers a material new dependency, migration, permission, secret, external network destination, public contract, destructive action, trust boundary or future-domain feature:

`STOP -> classify -> update plan/ADR/registry if required -> re-authorize -> continue`

### Self-modification guard

An AI authoring a product change must not silently weaken:

- governance;
- CI/workflows;
- security standards;
- quality gates;
- required tests;
- branch/ruleset enforcement;
- evidence definitions;
- active scope locks.

Material control-plane changes require separate review context and exact-head evidence.

### Test-oracle integrity

A failing test cannot be “fixed” by removing/weaking the assertion unless the specification itself is reviewed and changed.

Track:

- deleted tests;
- skipped tests;
- broadened ignores;
- reduced thresholds;
- removed negative cases;
- changed fixtures that hide the failure.

### Evidence authority

AI-authored prose is not machine evidence. Evidence should identify producer/tool/source SHA/environment/result/artifact digest and target/provider identity where applicable.

## 8. AI feature execution lifecycle

A product AI action should follow:

```text
Intent
  ↓
Context Assembly
  ↓
Risk / Policy
  ↓
Plan
  ↓
Tool Selection
  ↓
Approval if required
  ↓
Capability Invocation
  ↓
Owning Domain Validation
  ↓
Mutation / Result
  ↓
Audit + Event
  ↓
Postcondition / Reconciliation
  ↓
Outcome / Feedback
```

High-consequence autonomous action adds:

```text
Simulation / Policy Envelope
  ↓
Constraint Check
  ↓
Approval / Autonomy Budget
```

## 9. AI security invariants

- retrieved text, documents, emails, webpages, logs, issues and module content are untrusted data and can contain prompt injection;
- system/repository policy has higher authority than retrieved content;
- model output is untrusted until validated;
- agent-to-agent messages are not authority by themselves;
- external MCP/tool metadata is untrusted supply-chain input;
- tool arguments require typed validation and policy checks;
- secrets and RESTRICTED data are excluded from model context unless an explicit approved design says otherwise;
- tenant boundaries apply to embeddings/vector indexes/caches/evals/traces;
- AI cannot use system internals discovered through System Graph as permission to act;
- high-risk actions support re-authentication/human approval/dual control as required;
- incident response can revoke agent/model/tool/provider access without corrupting business history.

## 10. AI performance/cost architecture

AI quality is multi-dimensional:

- correctness/task success;
- groundedness/provenance;
- safety/policy compliance;
- latency;
- cost;
- determinism/replayability where required;
- human review burden;
- business outcome.

Model selection/routing should optimize against explicit task/risk budgets, not one global “best model.”

## 11. Design architecture for AI-native enterprise work

AI must reduce navigation burden but not hide system state.

Key UI patterns:

- explainable recommendations;
- preview-before-apply;
- diff/impact preview;
- approval reason/risk;
- source/citation/evidence visibility;
- reversible draft states;
- batch action review;
- exception inbox;
- role-aware command center;
- confidence/uncertainty labels;
- modelled vs observed distinction;
- timeline/audit reconstruction.

## 12. Autonomous Business OS maturity levels

Suggested autonomy levels:

- **L0 Observe** — search/summarize/explain only;
- **L1 Recommend** — proposes actions, no mutation;
- **L2 Draft** — creates draft records/plans for approval;
- **L3 Execute bounded** — performs low/medium-risk actions within policy;
- **L4 Orchestrate** — coordinates cross-domain workflow/agents with defined budgets/approvals;
- **L5 Optimize constrained** — continuously improves objectives inside hard policy/financial/risk envelopes.

No domain should jump to higher autonomy merely because the model is capable. Autonomy is capability-specific and evidence-earned.

## 13. What makes Omnexa materially beyond conventional ERP

The differentiator is not the number of modules. It is the integration of:

`deterministic domain truth + universal workflow + semantic business context + process intelligence + system intelligence + governed agents + simulation + continuous controls + measurable outcomes`.

That allows Omnexa to evolve from recording work to understanding, coordinating, simulating and safely optimizing work without surrendering deterministic business authority to AI.
