# Submodule Execution Profiles

Status: **Mandatory planning accelerator**

The owning dossier defines **what** a submodule owns, its flows and options. This document defines the default **how-to-decompose** task shape so a future AI does not redesign implementation sequencing from scratch.

All profiles inherit `DOSSIER_STANDARD.md` S01-S10. A profile adds mandatory task focus; it never changes `STATE.json` authorization.

## EP01 — Governed Record / CRUD Domain

Use for contacts, leads, workforce records, projects, suppliers and similar owned records.

- S01: record identity, lifecycle statuses, capability CRUD/search, permission names, duplicate/merge/archive rules.
- S02: owned tables, uniqueness/indexes, references, fresh/upgrade migrations, immutable historical snapshots.
- S03: create/update/archive/restore/merge validation, optimistic concurrency where needed.
- S04: optional reference/search/import adapters only; no foreign private writes.
- S05: row/object scope, tenant/org isolation, field visibility, bulk-operation limits.
- S06: list/detail/create/edit, bulk actions, loading/empty/error/conflict states, accessibility/localization.
- S07: disable/re-enable, import retries, archive/purge dependency checks.
- S08: CRUD, duplicate, concurrency, permission/tenant, migration, import/export, UI/e2e evidence.

## EP02 — Transaction / State Machine

Use for opportunities, orders, returns, payment intents, approvals, tickets, production orders.

- S01: allowed states/transitions, command capabilities, transition permissions, events.
- S02: aggregate/state/history/outbox ownership, transition uniqueness/idempotency keys.
- S03: transition guards, invariants, duplicate/retry semantics, compensation/reversal.
- S04: downstream capability/provider orchestration behind adapters.
- S05: transition-specific permissions, tenant/resource scope, privileged approvals.
- S06: timeline/status/actions with disabled-state reasons and conflict recovery.
- S07: restart/replay, partial downstream failure, timeout/retry, terminal-state recovery.
- S08: transition matrix, illegal transition, duplicate command, fault/restart, permission and integration tests.

## EP03 — Financial / Immutable Ledger

Use for GL, journals, inventory movement ledgers and other append/accounting invariants.

- S01: posting/movement contract, balancing/conservation invariant, source ownership, reversal policy.
- S02: append-only ledger/journal schema, period/effective-time indexes, immutable audit references.
- S03: deterministic validation/posting, idempotent source keys, reversal/correction rather than mutation.
- S04: tax/FX/bank/inventory/domain inputs as governed references/capabilities.
- S05: separation-of-duty, posting/approval permissions, closed-period/tenant controls.
- S06: journal/movement explorer, reconciliation/exception surfaces; never client-side authoritative totals.
- S07: crash between source/posting, duplicate delivery, period close/reopen policy, reconciliation rebuild.
- S08: invariant/property tests, duplicate posting, reversal, currency/rounding, migration and audit evidence.

## EP04 — Registry / Runtime Infrastructure

Use for module, capability, prompt, model, block, field, semantic-model and provider registries.

- S01: manifest/schema/version/compatibility and lookup contract.
- S02: registry persistence/cache/index where owned; uniqueness/version migration.
- S03: register/validate/activate/deprecate/resolve operations.
- S04: plugin/provider discovery adapters with explicit trust boundary.
- S05: admin/write permissions, malicious declaration validation, secret/network declarations.
- S06: registry/admin explorer when authorized.
- S07: duplicate registration, incompatible upgrade, missing optional provider, disable/re-enable.
- S08: schema compatibility, dependency, lifecycle and negative declaration tests.

## EP05 — Event / Async / Workflow Runtime

Use for event fabric, outbox/inbox, scheduler-like durable primitives and workflow runtime components.

- S01: envelope/identity, delivery semantics, trigger/action contract, ordering assumptions.
- S02: durable queues/streams/checkpoints/outbox/inbox/state schema.
- S03: publish/consume/transition logic, idempotency and deterministic retry decisions.
- S04: transport adapters separated from runtime semantics.
- S05: tenant/actor/correlation propagation, payload classification, poison input limits.
- S06: operator/diagnostic/timeline surfaces only when authorized.
- S07: crash/restart/replay, duplicate delivery, poison message, backoff/quarantine, compensation.
- S08: fault injection, replay, idempotency, restart and load/concurrency evidence.

## EP06 — Builder / Versioned Definition Runtime

Use for Page Builder, Template Builder, Form Builder, Custom Object Builder, Dashboard Builder and similar authoring products.

- S01: definition schema, registry dependencies, author/publish permissions, runtime/authoring boundary.
- S02: versioned definition/revision storage and compatible migration of saved definitions.
- S03: server validation, version diff/restore, runtime renderer/executor, publish/activate rules.
- S04: contributed blocks/actions/data providers through registries; missing dependency degradation.
- S05: definition/action permissions, unsafe expression/import validation, tenant isolation.
- S06: editor shell, keyboard-equivalent drag/drop, focus, responsive preview, history, autosave/conflict, W3C/WAVE/manual accessibility.
- S07: module removal, definition upgrade, import failure rollback, editor recovery/offline/degraded behavior.
- S08: schema/runtime/editor parity, migration, e2e, accessibility, large-document performance evidence.

## EP07 — External Provider / Connector Adapter

Use for payment providers, shipping, OAuth connectors, country fiscal adapters, model providers.

- S01: provider-neutral capability and normalized states/errors.
- S02: only provider references/config state owned locally; secrets classified separately.
- S03: normalized request/result mapping and idempotency boundary.
- S04: auth/credential lifecycle, timeouts, retry/rate limits, webhook/cursor handling.
- S05: SSRF/allowlists, secret redaction, signature/replay checks, data-classification constraints.
- S06: provider configuration/health UI when authorized, with secret-safe states.
- S07: provider outage/restart, expired credential, rate limit, duplicate webhook, fallback/degrade.
- S08: synthetic provider contract tests plus at least two adapters when provider neutrality is an exit gate.

## EP08 — Projection / Reporting / Analytics

Use for Customer 360, semantic metrics, reports, dashboards and cross-domain read projections.

- S01: source contracts, semantic/query contract, freshness and authorization rules.
- S02: projection/index/warehouse schema and rebuild checkpoints.
- S03: projection/update/query/aggregation logic; no source-domain writes.
- S04: event/capability ingestion and analytical-store adapters.
- S05: row/column/data-class policy, export controls, sensitive aggregation rules.
- S06: report/dashboard UX, stale/fresh indicators, accessible charts/tables/exports.
- S07: rebuild, lag, source outage, duplicate event, expensive-query cancellation.
- S08: aggregation correctness, access leakage negatives, rebuild/freshness and performance gates.

## EP09 — Identity / Security / Policy

Use for auth, RBAC, context policy, MFA, privileged access, SSO/SCIM and DLP.

- S01: actor/resource/action policy contract and deny-by-default rules.
- S02: identities/assignments/policies/security evidence schema with restricted secrets separated.
- S03: authenticate/evaluate/assign/revoke/step-up flows and recovery semantics.
- S04: IdP/KMS/device provider adapters behind bounded interfaces.
- S05: privilege escalation, cross-tenant, replay, recovery abuse, credential/session revocation tests.
- S06: security/admin surfaces with confirmation, least privilege and accessibility.
- S07: key/provider outage, policy rollback, session invalidation, emergency access expiry.
- S08: exhaustive policy matrix, negative security, audit and recovery evidence.

## EP10 — Edge / Offline Client

Use for POS/device/offline runtime.

- S01: server/edge ownership, command envelope, offline-allowed operation classes.
- S02: encrypted local queue/cache schema with bounded retention; never shadow authoritative ledgers.
- S03: local command capture, optimistic UX and deterministic sync protocol.
- S04: device adapters and connectivity/provider abstraction.
- S05: device identity, local data protection, offline risk limits, tamper/replay defenses.
- S06: offline/online/conflict/queue/device states and accessible operator recovery UX.
- S07: power loss, reconnect, duplicate sync, stale data, queue saturation, device failure.
- S08: fault/restart/replay and physical-device contract simulation evidence.

## EP11 — AI / Agent Runtime

Use for AI gateway, retrieval, tool runtime and governed agents.

- S01: model/tool/agent contract, permissions, approval levels, evaluation requirements.
- S02: run/trace/memory/reference schema without copying unauthorized source data.
- S03: context assembly, model invocation, tool proposal/execution, bounded loop semantics.
- S04: model/vector/provider adapters with cost/rate/data-use policy.
- S05: prompt/tool injection, data-scope leakage, secret filtering, approval bypass negatives.
- S06: explainability/approval/run-inspection UI when authorized.
- S07: model outage/fallback, tool timeout, max-step/cost stop, replay/simulation.
- S08: task quality, safety, permission, cost and model-change regression evaluation.

## EP12 — Composition Pack / Operations / Scale

Use for industry packs, marketplace packages, globalization packs and measured scale/DR operations.

- S01: manifest/dependencies/compatibility/objective and operator authority.
- S02: configuration/evidence/placement metadata; avoid duplicating domain data.
- S03: orchestration/composition/plan logic over governed capabilities.
- S04: infrastructure/provider/package adapters.
- S05: declared scopes, residency, signing/trust, privileged operations.
- S06: install/configuration/operations/control surfaces when authorized.
- S07: upgrade/rollback/failover/recovery/disable combinations.
- S08: compatibility matrix, disaster/fault rehearsal, combination and rollback evidence.

## Profile composition rule

A submodule may use one primary profile plus one secondary profile when genuinely mixed, e.g. Checkout = EP02 + EP07, Dashboard Builder = EP06 + EP08. Do not invent a new profile merely for naming preference. If none fits materially, update this canonical profile set through planning change control.