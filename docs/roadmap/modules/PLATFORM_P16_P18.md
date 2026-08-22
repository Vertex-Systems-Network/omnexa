# Platform Expansion Module Dossiers — P16 to P18

Status: **Mandatory future planning baseline**

## P16 — Omnexa Connect / Integration Fabric

Architecture: connectors are adapters around a common sync/webhook/credential runtime. Business modules own business state; connectors translate provider state into governed capabilities/events and never become alternate domain owners.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P16.A | Connector SDK | manifest + schemas -> generated/runtime adapter contract -> validation | SDK/runtime version, capability declarations |
| P16.B | Credential & OAuth Management | authorize -> store restricted token/reference -> refresh/rotate/revoke | grant type, scopes, refresh policy |
| P16.C | Webhook Ingress | receive -> authenticate -> dedupe -> normalize -> dispatch | signature/replay windows, body limits |
| P16.D | Webhook Egress | domain event -> transform -> sign -> deliver -> retry/quarantine | timeout, retry, signing profile |
| P16.E | Sync Engine | schedule/manual trigger -> cursor/page -> map -> apply governed capability -> checkpoint | direction, batch size, conflict policy |
| P16.F | Mapping/Transformation | source schema -> mapping/version -> validation -> target contract | field mapping, transforms, defaults |
| P16.G | Rate-limit & Retry Runtime | provider response -> quota/backoff state -> reschedule | concurrency, token bucket, retry limits |
| P16.H | Connector Health | credential/API/sync state -> health projection -> operator action | failure thresholds, alert severity |
| P16.I | Generic REST/Webhook Connector | configured endpoint/schema -> authenticated request/event flow | auth mode, allowed hosts, request limits |
| P16.J+ | Provider Adapters | provider-specific auth/API translation behind SDK | provider endpoints/features only |

Security: SSRF allowlists, secret redaction, OAuth state/PKCE, webhook replay defense, encrypted restricted credentials, bounded payloads. Tests: cursor restart, duplicate webhook, provider rate-limit, expired credentials, schema drift, partial batch failure, optional connector disable.

## P17 — Low-code App Builder

Architecture: custom apps use governed platform identity, data, permissions, workflow, audit and events. Definitions are versioned server-validated contracts; the visual UI is an authoring client, not the runtime source of truth.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P17.A | Custom App Package | create app -> manifest/version -> activate/archive/export | app metadata, dependencies, environments |
| P17.B | Custom Object Builder | define object -> fields/constraints -> validate -> migrate/version | naming, retention hooks, record limits |
| P17.C | Field/Schema Registry | choose field type -> constraints/defaults -> migration compatibility | types, required/index/unique, computed rules |
| P17.D | Relationship Builder | define relation -> ownership/cardinality -> validate lifecycle | one/many, delete behavior, lookup policy |
| P17.E | Form Builder | schema -> sections/fields/validation -> action binding -> publish | layouts, conditional visibility, submit actions |
| P17.F | List/Table View Builder | columns/filter/sort/actions -> saved view | page size, column rules, bulk actions |
| P17.G | Kanban Builder | grouping/status mapping -> cards -> drag/keyboard transition | group field, transition permissions |
| P17.H | Calendar View Builder | date mapping -> event display/actions | timezone, duration fields, range |
| P17.I | Permission Builder | object/action/field scope -> policy declaration -> role assignment | CRUD/action/field visibility rules |
| P17.J | Workflow Integration | record event/action -> P05 trigger/action registry | allowed triggers/actions, idempotency |
| P17.K | Dashboard Builder | governed metrics/widgets -> layout -> publish | widget types, refresh/freshness |
| P17.L | Portal Exposure | app capability/view -> portal profile slot | persona, permissions, route visibility |
| P17.M | Runtime CRUD/API | validated definition -> generated runtime capability/API -> policy -> persistence | pagination, concurrency, bulk limits |
| P17.N | Authoring Shell | object/form/view/workflow designers -> preview/version/publish | autosave, history, validation modes |

Critical rule: no arbitrary code execution from user definitions. Expressions/actions use bounded registries. Schema changes require migration plan, dry-run and rollback/compatibility evidence. UI must provide keyboard alternatives for drag/drop, WCAG 2.2 AA target and W3C/WAVE/manual evidence.

## P18 — Data, Reporting & BI

Architecture: analytics uses governed read projections/semantic models; heavy analytical workloads must not mutate or overload OLTP write models.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P18.A | Semantic Model Registry | source contract -> dimensions/measures -> validate/version | naming, aggregation, classification |
| P18.B | Metrics Engine | semantic query -> policy -> compute/read projection -> result | freshness, aggregation window, cache policy |
| P18.C | Report Builder | dataset/model -> fields/filters/grouping -> preview -> save/publish | pagination, export limits, formatting |
| P18.D | Dashboard Runtime | widgets -> metric/report queries -> layout -> refresh | refresh cadence, stale indicator |
| P18.E | Dashboard Builder | choose widgets -> configure queries -> layout -> publish | grid, responsive behavior, sharing |
| P18.F | Scheduled Reports | report + recipients/destination -> schedule -> render/export/deliver | timezone, frequency, format, retry |
| P18.G | Export Service | authorized query -> bounded export job -> storage reference | row/file limits, formats, expiry |
| P18.H | Cross-domain Read Projections | events/capabilities -> projection -> freshness/health | projection lag target, rebuild policy |
| P18.I | Analytical Pipeline/Warehouse Boundary | governed extract -> transform/load -> analytical store | batch/stream cadence, retention |
| P18.J | Data Access Governance | actor + semantic object -> row/column/classification policy | masking, export permission, audit |
| P18.K | Query/Performance Guardrails | query plan -> cost/limit enforcement -> execute/cancel | timeout, rows, concurrency, complexity |

Tests: row/column authorization, export leakage negatives, semantic aggregation correctness, stale/fresh indicators, scheduled retry/idempotency, projection rebuild, high-cost query rejection.

## Common platform flow rules

Every visual builder separates definition/version storage, server validation, runtime execution and authoring UI. Every provider adapter is replaceable. Every scheduled/retriable operation uses P04/P05 primitives where available rather than private retry engines.