# Intelligence Module Dossiers — P19 to P20

Status: **Mandatory future planning baseline**

## P19 — Omnexa Intelligence Platform

Architecture: AI is a governed platform consumer. Models never receive unrestricted database access and never become authoritative business-state owners. All retrieval and actions pass through identity, tenant/context, policy and approved capabilities.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P19.A | Model Provider Registry | configure adapter/model -> capability/cost/health validation -> activate | provider/model allowlist, regions, modalities |
| P19.B | AI Gateway | request -> policy/budget/rate check -> provider -> normalized response | model routing, timeout, fallback policy |
| P19.C | Prompt Registry | prompt definition -> variables/version/evaluation -> activate | version, locale, model requirements |
| P19.D | Tool/Capability Registry | approved capability -> tool schema -> permission binding -> invocation | allowed actors, confirmation level, timeout |
| P19.E | Retrieval/Knowledge Base | source authorization -> ingest/chunk/index -> query -> cited context | sources, chunking, freshness, retention |
| P19.F | Embedding/Vector Boundary | governed text/data -> embedding provider -> index reference | model/dimensions, reindex policy |
| P19.G | Conversation/Run Context | request -> bounded context/history -> model/tool loop -> result | token/context limits, retention, privacy |
| P19.H | Model Cost & Rate Policy | usage estimate/actual -> tenant/provider budget -> allow/throttle/block | budgets, quotas, model tiers |
| P19.I | Evaluation Framework | fixture/task -> run -> score/safety/quality metrics -> compare | datasets, thresholds, regression policy |
| P19.J | Traceability | model/prompt/tools/context refs -> safe trace/audit record | redaction, retention, sampling |
| P19.K | AI Permission & Data Boundary | actor/tenant/data class/tool -> policy -> allowed context/actions | data classes, field masking, tool grants |
| P19.L | Human Approval Hooks | proposed privileged action -> approval task -> approved capability execution | risk levels, approver policy, expiry |
| P19.M | Safety/Content Policy Hooks | request/response/tool proposal -> policy evaluators -> outcome | policy sets, escalation, fail-closed classes |

Required flow: authenticated request -> tenant/context -> policy -> prompt/model routing -> retrieval through authorized sources -> optional tool proposal -> approval when required -> governed capability -> trace/evaluation. Provider fallback may change model/provider, never permissions or data scope.

Security: prompt/tool injection defenses, retrieval authorization at query time, secret stripping, model-provider data-use declarations, no production-sensitive evaluation fixtures by default, bounded tool recursion/cost.

## P20 — Governed AI Agents

Architecture: an agent is a policy-bound orchestrator over P19 + P05 + domain capabilities. It is not a superuser and cannot invent permissions.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P20.A | Agent Definition Registry | role/goal/tools/policies -> validate/version -> activate | allowed tools, model class, autonomy level |
| P20.B | Agent Runtime | objective -> plan/propose -> tool calls -> observe -> complete/escalate | max steps, time/cost budget, retry bounds |
| P20.C | Memory Boundary | approved fact/context -> scoped memory -> retrieve/update | scope, TTL/retention, sensitivity |
| P20.D | Approval & Escalation | risk/action -> require human/policy -> approve/reject/timeout | thresholds, approver groups |
| P20.E | Sales Agent Pack | CRM context -> recommendations/actions through CRM capabilities | allowed CRM tools, outreach approval |
| P20.F | Finance Agent Pack | finance read context -> analysis/proposals -> approved finance actions | posting/payment approval always explicit by policy |
| P20.G | Procurement Agent Pack | demand/supplier context -> recommendation/workflow | spend limits, supplier constraints |
| P20.H | Support Agent Pack | case/customer context -> response/action proposal | response autonomy, refund/escalation limits |
| P20.I | Executive Analysis Agent | governed metrics -> analysis/scenarios -> recommendation | datasets, freshness, no unrestricted writes |
| P20.J | Agent Evaluation & Regression | benchmark scenarios -> run -> safety/task/cost score -> release gate | score thresholds, model-change gates |
| P20.K | Agent Audit/Replay | run inputs/model/tools/approvals -> immutable evidence -> replay/simulation | retention, redaction, export policy |
| P20.L | Multi-agent Coordination | parent objective -> bounded delegation -> results -> reconcile | allowed agent graph, recursion/delegation limits |

Agent action flow:

```text
objective
 -> actor/tenant/policy context
 -> agent definition
 -> evidence/retrieval
 -> proposed plan
 -> capability/tool policy
 -> human approval when required
 -> governed capability/workflow
 -> domain validation
 -> audit/event
 -> measurement/evaluation
```

Forbidden: direct SQL writes, hidden service-account escalation, tool schema that exposes arbitrary HTTP/DB access, silent approval bypass, agent-to-agent delegation outside declared graph.

## Common intelligence options

All AI/agent options must declare who may change them, tenant scope, model/provider dependency, cost impact, data-classification impact, evaluation requirements and rollback behavior. Model/prompt changes that materially change behavior require evaluation evidence before activation.