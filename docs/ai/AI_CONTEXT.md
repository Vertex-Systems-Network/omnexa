# Omnexa AI Project Context

Status: **P04.04 ACTIVE — post-activation continuity; subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative checkpoint

Fresh protected main is `50edeae03ad52e435d142b2f22c803f08a5c7f1a`.

At that exact protected-main checkpoint:

- Foundation Architecture v1 is FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED / historical.
- P04 is ACTIVE — 3 / 10 canonically done.
- P04.01-P04.03 are DONE with accepted evidence.
- P04.04 is the sole ACTIVE package.
- P04.05-P04.10 are planned/locked.
- `kernel_code_authorized=true` only for P04.04.
- `business_feature_code_authorized=false`.
- strategic X-program and AI/model/agent runtime remain unauthorized.

## Accepted activation chain

P04.03 closure / P04.04 activation is accepted evidence:

- source PR #171;
- exact source/promotion head `2cb7a92217a81b921e40a274007e1c50dc4a8fcf`;
- source Governance `33546204586 / 99984199589` — PASS;
- source SELF REVIEW; zero unresolved threads;
- unchanged promotion PR #172;
- promotion Governance `33547080086 / 99987132610` — PASS;
- promotion SELF REVIEW; zero unresolved threads;
- guarded merge/read-back `50edeae03ad52e435d142b2f22c803f08a5c7f1a`.

No runtime code or migration was part of activation.

## P04.04 accepted implementation boundary

P04.04 owns producer-side transactional outbox reliability under `kernel.events`.

Implementation must:

- preserve the canonical P04.01 envelope and stable EventID;
- atomically enqueue the canonical event in the same local PostgreSQL transaction as the authorized owner mutation;
- reuse P01 `database.InTransaction`, not create a second transaction framework;
- make committed pending state durable/recoverable after restart;
- relay through the accepted P04.02 publication abstraction;
- keep publication failures pending;
- permit duplicate relay after publish-success/crash-before-mark and never claim end-to-end exactly once;
- advance published state only after successful publish;
- remain race-safe under concurrent relay attempts;
- preserve owner/tenant/event identity isolation;
- fail safely on malformed persistence/database failures without payload/credential/topology leakage;
- introduce no global ordering guarantee.

## Persistence/migration preflight

Durable local PostgreSQL outbox persistence is required for production P04.04 semantics. Activation granted no blanket schema authority.

Before schema mutation on the fresh implementation branch:

1. inspect existing `kernel.events` migrations and retained P01 migration conventions;
2. choose the exact next immutable owner-scoped migration version/path;
3. record the exact tables/columns/indexes/constraints/data budget;
4. preserve fresh-install/upgrade ledger behavior and owner isolation;
5. define rollback/forward-recovery that cannot destructively erase an event that may already have been published;
6. exclude P04.05 inbox and P04.06 retry/DLQ semantics.

## Locked scope

Do not implement concrete Kafka/NATS/RabbitMQ/Redis Streams providers, consumer inbox/dedup, retry/backoff/DLQ/quarantine, schema-registry runtime, background-job runtime ownership, business handlers/features, strategic X runtime or AI/model/agent runtime.

## Exact next action

Govern and land this post-activation continuity correction first. After protected-main read-back, create a **fresh separate P04.04 implementation branch** from that exact SHA, record the exact runtime/migration path budget, and implement only `docs/roadmap/work-packages/P04.04.md`.

Do not reuse activation/promotion/continuity branches. Do not auto-advance P04.05.
