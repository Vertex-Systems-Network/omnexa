# ADR-0011 — AI-Native Business Operating System Strategic Expansion

Status: **proposed**

## Context

Omnexa Foundation Architecture v1 correctly establishes a strict modular-monolith-first platform with domain-owned write models, governed capabilities/events/workflows, tenant isolation, policy-before-action, durable audit, OpenTelemetry-compatible observability and governed AI execution.

The P00-P27 roadmap already spans ERP, CRM, finance, commerce, payments, POS, CMS, portals, HR/projects/service, supply chain/manufacturing, integrations, low-code, BI/data, AI/agents, developer ecosystem, marketplace, globalization, enterprise governance, scale, industry packs and an Autonomous Business OS.

A 2026 strategic audit against the intended product category — a composable enterprise Business Operating System above conventional ERP/business-suite scope — identified cross-cutting systems that are either implicit, under-specified or absent as explicit governed programs. The largest gap is that the Product Constitution names **Unified Business Graph** as a strategic moat while the roadmap has no dedicated implementation program for it.

The audit also found that Omnexa's existing AI execution/continuity policy is strong procedurally, but future AI-assisted engineering needs stronger machine-level orchestration: exact-run identity, scope leases, self-modification protections, test-oracle integrity, evidence attestation, independent exact-head review, multi-agent coordination and bounded retry/cost controls.

External market direction in 2026 reinforces the need for these foundations: major enterprise platforms are converging on trusted business-data fabrics/knowledge graphs, process-aware agents, central AI control towers, agent studios, multi-agent orchestration, continuous controls, model governance and outcome/value measurement. Omnexa should not copy vendor implementations; it should provide provider-neutral governed equivalents integrated with its own kernel, Business Graph, Process Graph and System Graph.

## Problem

If Omnexa continues with only the current phase catalog, several risks emerge:

1. the Unified Business Graph remains a slogan rather than a governed platform capability;
2. data lineage/master-data/semantic context can fragment across domains and analytics;
3. workflow execution exists without a first-class process-mining/intelligence layer;
4. observability exists without a canonical system/runtime relationship graph and explainable flow intelligence;
5. AI/agent breadth can grow without a central AI asset/control/value plane;
6. AI development can remain dependent on procedural discipline instead of machine-enforced orchestration;
7. enterprise finance can mature before segregation-of-duties/continuous-control foundations are deep enough;
8. performance/capacity/cost decisions can remain distributed across tests and operations instead of becoming explainable platform intelligence;
9. low-code/agents/workflows/configuration can become difficult to promote safely across environments;
10. autonomous decisioning can arrive before simulation, outcome measurement and autonomy-risk controls are mature.

## Decision

Introduce a **Strategic Cross-Cutting Program Layer** using stable `X..` program IDs without renumbering P00-P27 and without authorizing implementation ahead of existing phase gates.

The new strategic program registry is `docs/roadmap/STRATEGIC_PROGRAMS.json`; narrative architecture and dependency rules are defined in `docs/roadmap/STRATEGIC_CROSS_CUTTING_PROGRAMS.md` and `docs/architecture/AI_NATIVE_BUSINESS_OS_ARCHITECTURE.md`.

Core strategic programs include:

- AI-native engineering governance / Quality OS;
- enterprise data governance, catalog, lineage, master/reference data and semantic layer;
- Unified Business Graph;
- System Graph and Flow Intelligence;
- Process Intelligence / process mining;
- module/package trust and runtime security;
- payment-security certification layer;
- performance/capacity intelligence;
- GRC / SoD / continuous controls;
- enterprise planning/EPM;
- AI Control Tower;
- Agent Studio and external-agent interoperability;
- model lifecycle/evaluation/MLOps governance;
- configuration-as-code/environment ALM;
- migration/data-onboarding center;
- FinOps/cost + business-outcome intelligence;
- decision simulation/business digital twin;
- unified Work/Command Center;
- document intelligence and enterprise knowledge;
- service management/CMDB/business-service graph;
- autonomous-business governor.

## Three-graph architecture

Omnexa will distinguish three related but non-authoritative-by-themselves graph planes:

1. **Business Graph** — semantic business entities and relationships: people, organizations, products, accounts, orders, invoices, payments, assets, locations, obligations and other domain references.
2. **Process Graph** — what business process actually happened over time: workflow/activity/event sequences, waits, rework, approvals, retries, exceptions, conformance and bottlenecks.
3. **System Graph** — software/runtime topology and evidence: modules, capabilities, APIs, events, jobs, data stores, caches, external providers, permissions, trust boundaries, errors, performance, tests and deployment relationships.

None of these graphs becomes a hidden write authority. Owning domains remain authoritative for business state. Graphs are derived/declared/observed semantic and operational evidence planes.

## Evidence rule

Where graph/flow intelligence is implemented, evidence classes remain distinct:

- `declared`
- `static`
- `observed`
- `tested`
- `production-observed`
- `modelled`
- `ai-inferred`

A diagram, AI explanation or modelled scenario is not runtime/target evidence. Missing evidence remains missing rather than being promoted to PASS.

## AI-native engineering role model

Future AI-assisted engineering uses governed roles, not an unrestricted general agent. Required role capabilities include:

- AI Architect;
- AI Design Engineer;
- AI Performance Engineer;
- AI Systems Engineer;
- AI Analyzer;
- AI Developer;
- AI Domain Expert;
- AI Model Lifecycle Manager (AI-MLM);
- AI Engineer;
- AI Constructor;
- plus independent AI Security/Red-Team, Data/Knowledge, QA/Verifier, Release/Operations and Product/Outcome review roles where risk requires them.

Role separation is about authority/evidence, not model branding. One model may execute multiple low-risk roles, but high/critical work requires independent review/evidence paths.

## Compatibility impact

No current runtime or public contract is changed by this planning ADR. P00-P27 semantic IDs remain stable. Existing completed P00/P01/P02 work is not invalidated.

Future implementations must be additive and use existing kernel/module/capability/event/workflow ownership rules.

## Migration impact

None for this planning change.

Future graph/data/AI-control systems must define additive migrations, backfills, rebuild semantics for derived stores and clear uninstall/deprecation behavior.

## Security and tenancy impact

The expansion strengthens security requirements:

- Graph and AI-control surfaces are sensitive reconnaissance/authority surfaces and default-deny access;
- tenant/org scope must propagate into all graph, AI, analytics, process and model evidence;
- AI and graphs never gain direct cross-domain write authority;
- third-party modules/agents/models/connectors require declared permissions/data/network/secrets and supply-chain provenance;
- payment/security/financial evidence remains separately certified;
- sensitive telemetry and model context are redacted by classification.

## Operational impact

Adds future telemetry, lineage, graph, model-governance and process-intelligence workloads. Infrastructure choice must remain evidence-based; specialized graph/vector/stream/warehouse infrastructure is not mandated until measured need justifies it.

## Rollback / forward-fix

Because this is a planning expansion, rollback means removing the proposed strategic addendum before acceptance. After acceptance, individual strategic programs may be deferred/resequenced through normal change control, but the three-graph separation, domain-write-authority rule and AI least-authority principles require a superseding ADR to change.

## Documents/work packages affected

This proposal adds or updates:

- `docs/roadmap/MASTER_PLAN.md`
- `docs/roadmap/STRATEGIC_CROSS_CUTTING_PROGRAMS.md`
- `docs/roadmap/STRATEGIC_PROGRAMS.json`
- `docs/roadmap/AI_NATIVE_STRATEGIC_AUDIT_2026.md`
- `docs/architecture/AI_NATIVE_BUSINESS_OS_ARCHITECTURE.md`
- `docs/governance/AI_NATIVE_ENGINEERING_ROLES.md`
- `docs/governance/AI_EXECUTION_POLICY.md`
- `docs/quality/QUALITY_GATE_MATRIX.md`
- `docs/governance/DEFINITION_OF_DONE.md`

Current execution remains `P02.08` only. This ADR does not authorize P03+ or any `X..` program implementation.
