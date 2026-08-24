# Omnexa Strategic Program Acceptance Gates

Status: **Proposed under ADR-0011**  
Rule: these gates define future closure evidence; they do not activate implementation.

## XQ-100 — AI-Native Engineering Governance & Quality OS

Exit requires:

- machine-readable development run/scope model;
- exact base/source and active-work-package binding;
- scope-delta stop/re-plan enforcement;
- test-oracle weakening detection;
- governance/workflow self-modification guard;
- machine evidence producer/attestation semantics;
- independent exact-head review requirement for high/critical changes;
- multi-agent conflict/lease handling;
- bounded attempt/cost circuit breaker;
- representative adversarial fixtures for prompt injection, stale context, fake evidence, scope race and control weakening;
- Research/VOC/CTQ + proportional DMADV/DMAIC/FMEA workflow validated on representative work.

## XDG-100 — Data Governance / Catalog / Lineage / Master Data

Exit requires:

- data catalog and authoritative owner model;
- lineage through at least one cross-domain + analytics/AI path;
- classification/tenant/org/field authorization preserved;
- retention/delete/export propagation proven;
- derived-store rebuild/reconciliation proven;
- master/reference/golden-record policy with no hidden cross-domain write authority;
- data-quality metrics and exception ownership;
- residency/AI exposure policy exercised.

## XBG-100 — Unified Business Graph

Exit requires:

- typed/temporal/provenance-aware graph contract;
- representative entities from at least three independently owned domains;
- tenant/org and field authorization negative tests;
- source-domain identity preserved;
- graph rebuild/projection consistency test;
- no direct business mutation through graph storage/query layer;
- versioned query/semantic compatibility policy;
- AI context retrieval respects source permissions.

## XSG-100 — System Graph & Flow Intelligence

Exit requires:

- typed node/edge/evidence schema;
- module/capability/event/job/data/network/permission/security relationships;
- distinct `declared/static/observed/tested/production-observed/modelled/ai-inferred` evidence;
- representative Core/module/federated-product runtime paths;
- unexpected/undeclared edge drift detection;
- sensitive topology default-deny/redaction/audit;
- bounded collector overhead;
- graph history/diff and source/package/deployment identity;
- change-impact/blast-radius analysis clearly labelled by evidence level.

## XTRUST-100 — Module/Connector/Agent/Package Trust

Exit requires:

- publisher/signature/provenance/SBOM validation;
- declared capability/data/network/secret profile;
- dependency/advisory/license policy;
- runtime trust tiers;
- sandbox/brokered network/secret/file access for untrusted executable profiles where supported;
- resource/time quotas;
- malicious-package fixtures;
- emergency revoke/kill and safe recovery;
- first-party packages subject to the same public-boundary rules.

## XPF-200 — Product Federation & App Mesh

Exit requires:

- native/embedded/federated/edge attachment contracts;
- representative standalone first-party product integration without private DB coupling;
- SSO/principal/tenant/org mapping;
- versioned capabilities/events/data contracts;
- Business/System Graph contributions;
- workflow and AI-tool federation;
- entitlement vs authorization separation;
- health/SLO/version identity;
- disconnect/revoke/deprecation with historical evidence preserved;
- cross-product tenant isolation and token audience/scope negative tests.

## XPERF-100 — Performance / Capacity Foundation

Exit requires:

- request/module/capability/query/cache/event/workflow/provider attribution;
- performance profiles/budgets for representative critical paths;
- tenant/noisy-neighbor dimensions;
- queue/event lag and saturation evidence;
- load/capacity baseline;
- source/version correlation;
- measured telemetry overhead and redaction.

## XPERF-200 — Performance Intelligence

Exit requires:

- history and version comparisons;
- root bottleneck attribution across application/data/event/provider layers;
- capacity forecasting with uncertainty/limits;
- regression alerts;
- AI recommendations grounded in measured evidence;
- scale-fabric recommendations cannot claim necessity without measured trigger.

## XGRC-100 — Financial SoD / Maker-Checker

Exit requires:

- SoD conflict model;
- sensitive-access matrix;
- maker-checker/approval-authority rules;
- allow/deny/privilege-escalation tests;
- financial period/control integration;
- immutable review evidence;
- tenant/legal-entity scope.

## XGRC-200 — Enterprise GRC

Exit requires:

- risk/control registry;
- access reviews/role certification;
- continuous transaction/configuration monitoring;
- control testing and evidence;
- policy exception/waiver lifecycle;
- legal hold/retention linkage;
- compliance evidence automation with explicit source/provenance;
- no generic “compliant” claim without applicable scoped evidence.

## XPS-100 — Payment Security

Exit requires:

- standard profile excludes raw PAN/CVV/track/PIN from Omnexa runtime/log/cache/queue/backup/AI;
- provider-hosted/tokenized integration path;
- payment-specific capabilities/state model;
- Secret/Network Brokers;
- signed/fresh/replay-safe tenant/provider-bound webhooks;
- duplicate/concurrent authorize/capture/refund safety;
- ambiguous-timeout reconciliation;
- 3DS/SCA/asynchronous state tests;
- payment-surface script/data-leak controls;
- real provider sandbox tests;
- credential rotation/kill/reconciliation procedure;
- threat model + FMEA + independent security review.

## XEPM-200 — Enterprise Planning / EPM

Exit requires:

- plan/version/scenario model separate from actuals;
- exact financial arithmetic/currency rules;
- budget/forecast/driver planning;
- representative cross-functional plan;
- profitability/cost-allocation controls;
- consolidation/close boundary defined;
- scenario/what-if results clearly modelled;
- authorization/approval/version/audit;
- reconciliation to Finance actuals.

## XPROC-200 — Process Intelligence

Exit requires:

- event/activity provenance model;
- actual process discovery from representative workflow/manual/domain activity;
- conformance and variant analysis;
- bottleneck/wait/rework/loop metrics;
- tenant/business-scope isolation;
- process model vs actual distinction;
- automation recommendation evidence;
- historical comparison.

## XAIC-200 — AI Control Tower

Exit requires:

- inventory for models, agents, prompts, tools, MCP endpoints, datasets and knowledge sources;
- ownership/tenant/org/risk/lineage;
- permissions/data access;
- runtime trace/eval/cost/value status;
- shadow/unmanaged asset discovery path where feasible;
- policy/compliance gates;
- incident/case workflow;
- kill/revoke/deprecate lifecycle;
- third-party assets supported without vendor-specific core coupling.

## XAIS-200 — Agent Studio / Multi-Agent / Interop

Exit requires:

- low-code/pro-code agent definition contract;
- typed tool and context scopes;
- deterministic-workflow + agent hybrid;
- supervisor/specialist/agent-team pattern;
- human approval and autonomy limits;
- sandbox/eval/test environment;
- external interoperability/MCP with identity/policy;
- version/publish/rollback;
- delegation cannot widen authority.

## XMLM-200 — Model Lifecycle / Evaluation / MLOps

Exit requires:

- provider/model registry and model cards;
- model/prompt/tool/dataset/eval lineage;
- offline and representative online evals;
- task success/groundedness/safety/latency/cost dimensions;
- drift/regression monitoring;
- routing/fallback/deprecation/rollback;
- training/fine-tune governance if supported;
- red-team suites;
- tenant-sensitive training/use policy.

## XALM-200 — Config-as-Code / Environment ALM

Exit requires:

- export/import/version format;
- semantic diff;
- dependency validation;
- secret references only;
- DEV/TEST/UAT/PROD promotion;
- approval;
- rollback/forward-fix;
- drift detection;
- representative workflow/low-code/report/agent asset promotion.

## XMIG-200 — Migration / Onboarding Center

Exit requires:

- source profiling;
- mapping/transform validation;
- dry run;
- data-quality exception queue;
- idempotent resume;
- duplicate handling;
- financial/business reconciliation;
- incremental cutover where applicable;
- migration audit;
- rollback/compensation/forward-fix policy.

## XFIN-200 — FinOps / Resource Cost Intelligence

Exit requires:

- reconciled cost-source identities;
- attribution to tenant/module/capability/workflow/provider/AI where supported;
- allocation rule versioning;
- budgets/anomaly detection;
- tenant-safe reporting;
- no implication that cost allocation is billing authority;
- AI cost/value linkage with explicit assumptions.

## XOUT-200 — Business Outcome / Value Intelligence

Exit requires:

- goal/CTQ registry;
- pre-defined metric/source/baseline;
- representative feature/process/agent outcome linkage;
- observation window;
- financial/non-financial value measures;
- uncertainty/limitations;
- no fabricated post-release outcome;
- control action when outcome regresses.

## XDOC-200 — Document Intelligence / Enterprise Knowledge

Exit requires:

- ingestion/OCR/extraction pipeline;
- source/provenance/citations;
- classification/retention/legal-hold handling;
- permission-aware search/retrieval;
- tenant-isolated indexes/vector stores;
- Business Graph linking;
- deletion propagation;
- adversarial document/prompt-injection tests for AI retrieval.

## XESM-200 — Enterprise Service Management / CMDB

Exit requires:

- configuration item/business service model;
- service catalog + incident/problem/change/request lifecycle;
- ownership/SLA/SLO linkage;
- System Graph/service dependency correlation;
- operational impact analysis;
- tenant/access controls;
- automation/agent integration through governed capabilities.

## XDEC-300 — Decision / Simulation / Digital Twin

Exit requires:

- versioned scenario inputs/assumptions;
- constraint/policy model;
- uncertainty/limits;
- modelled-vs-live separation;
- representative finance/supply/workforce/capacity scenarios;
- reproducibility where possible;
- approval before consequential action;
- outcome comparison after real execution.

## XWORK-200 — Unified Business Work / Command Center

Exit requires:

- role/intention-driven work queue;
- unified approvals/exceptions/tasks/KPIs/alerts;
- cross-domain context through public read/capability contracts;
- natural-language search/action with policy/approval;
- explainable “why shown/why allowed/why blocked”;
- accessibility/i18n/RTL;
- no UI-based authorization assumption;
- measured task-efficiency/usability outcome.

## XAUTO-300 — Autonomous Business Governor

Exit requires:

- capability-specific autonomy levels;
- risk/action/financial budgets;
- objective/constraint registry;
- approval thresholds;
- simulation-before-action rules where consequential;
- multi-agent conflict resolution;
- kill/containment/revocation;
- reversal/compensation/reconciliation;
- continuous quality/risk/outcome monitoring;
- automatic downgrade/stop when evidence violates policy.

## Universal gate

No strategic program may be reported `done` from architecture/design/code alone. Applicable source, integration, security, data, performance, reliability, target/provider and outcome evidence must be recorded under the canonical evidence vocabulary.
