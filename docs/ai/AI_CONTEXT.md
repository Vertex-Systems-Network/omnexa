# Omnexa AI Project Context

Status: **P04.06 POST-ACTIVATION CONTINUITY — subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative protected-main base

Fresh protected main at this continuity start is:

`4c9f60843f2612bc4c9a10b4efca7b6a20826be3`

That protected-main state contains the accepted P04.05 closure / P04.06 activation transaction:

- P04.05 implementation acceptance through its governed T01-T04 chain;
- P04.05 completion evidence merged/read back as `e44ece77ddf7b821c03997266ca0c68c07162910`;
- separately governed P04.06 preparation source/promotion at exact head `7babb9c39185636b3af5184d5a7bd31cedbc37a0`;
- activation source PR `#218` at exact head `b718ad7316dba6fca0cafccb514df6da653abe13`;
- source Governance `#677 / 33990309561` — PASS;
- unchanged activation promotion PR `#219` at the same exact head;
- promotion Governance `#678` — PASS;
- accepted protected-main merge/read-back `4c9f60843f2612bc4c9a10b4efca7b6a20826be3`.

Activation is accepted. This separate continuity carrier exists only to reconcile stale subordinate pre-activation wording before runtime implementation begins.

## Accepted canonical result

Current protected-main truth is:

- P04 ACTIVE — `5 / 10 done`;
- P04.01-P04.05 DONE with retained evidence;
- P04.06 sole ACTIVE package;
- P04.07-P04.10 PLANNED / LOCKED;
- bounded kernel authority only for P04.06;
- business-feature authority remains false;
- strategic X and AI/model/agent product runtime remain unauthorized;
- P04.06 runtime worker slots/tasks/branches/migration reservations remain zero.

P04.06 being ACTIVE does not authorize runtime mutation on this continuity carrier or its unchanged promotion. Runtime work starts only after this continuity sequence is accepted and protected main is re-read.

## Accepted P04.05 completion evidence

The final P04.05 Supervisor verifier/promotion head is `dd713fe3217a0d092ab3ff31115ac031ae8c0303`.

- source Governance `#669 / 33927932705` — PASS;
- promotion Governance `#670 / 33985111334`, job `101357077040` — PASS;
- implementation merge/read-back `0c66a3371dbf2fa942a95b7d0475b06235392474`;
- completion evidence `docs/roadmap/evidence/P04.05_COMPLETION_2026-09-05.md`;
- completion evidence merge/read-back `e44ece77ddf7b821c03997266ca0c68c07162910`.

The accepted P04.05 guarantee is local duplicate-safe application inside the retained PostgreSQL boundary. Checkpoint and inbox completion remain separate facts; no external/end-to-end exactly-once claim exists.

## Accepted P04.06 preparation

- source preparation PR `#215`;
- unchanged promotion PR `#216`;
- exact source/promotion head `7babb9c39185636b3af5184d5a7bd31cedbc37a0`;
- source Governance `#674 / 33987003924` — PASS;
- promotion Governance `#675 / 33987472967` — PASS;
- preparation merge/read-back `3f547180eb5e839439834eb2ce7977324803df18`;
- contract `docs/roadmap/work-packages/P04.06.md`;
- handoff `docs/ai/handoffs/P04.06.md`.

The prepared contract selects no provider, grants no migration and creates no runtime wave.

## Accepted P04.05 closure / P04.06 activation

The explicit sequential transition was accepted through source `#218` and unchanged promotion `#219` at exact head `b718ad7316dba6fca0cafccb514df6da653abe13`.

The activation only changed governance/state/continuity:

- P04.05 moved from ACTIVE to DONE;
- P04.06 moved from PLANNED to sole ACTIVE;
- P04 progress moved from 4 / 10 to 5 / 10 done;
- bounded kernel authority moved from P04.05 to P04.06;
- completed P04.05 leases became historical/released;
- P04.06 runtime worker slots/tasks/branches/migration reservations remained zero.

No retry runtime, migration/schema reservation, provider, broker-native DLQ, P04.07+, business feature or AI/model runtime was introduced.

## P04.06 active boundary

Owner: `kernel.events`.

After this separate post-activation continuity source/promotion/read-back is accepted, a fresh implementation wave may implement only the accepted provider-neutral P04.06 contract:

- structured deterministic failure disposition rather than raw error-string matching;
- finite attempt budget and capped deterministic backoff;
- authoritative UTC retry eligibility;
- one-at-a-time retry claim/CAS/lease semantics;
- checkpoint, inbox and retry/quarantine evidence as distinct facts;
- P04.05 already-applied completion suppressing stale retry mutation;
- fail-closed owner/consumer/route/tenant rebinding;
- terminal/retry-exhausted logical quarantine;
- quarantine persistence before checkpoint may advance past a poison delivery;
- crash-gap recovery where committed quarantine suppresses handler reinvocation while checkpoint catches up;
- bounded safe quarantine evidence without raw payload/secrets/provider/database diagnostics.

P01.09's in-memory bounded job retry is precedent only. It is not durable event-fabric retry persistence/timing authority.

## Zero-runtime continuity law

This continuity carrier intentionally has:

- zero P04.06 runtime worker slots;
- zero P04.06 runtime tasks;
- zero P04.06 runtime branches;
- zero P04.06 migration reservations;
- no retry worker/scheduler source;
- no quarantine persistence source;
- no workflow/provider dependency change.

Completed P04.05 T01-T04 leases are historical/released. Accepted `kernel.events` migration v2 remains immutable inbox history, not a live P04.06 lease.

## Persistence / migration gate

Before any later P04.06 schema mutation, after post-activation continuity is accepted, fresh protected main must be re-read and the runtime wave must record:

1. exact current `kernel.events` migration history;
2. exact next immutable owner-scoped migration version/path/name;
3. exact retry/quarantine tables/indexes/constraints/data budget;
4. processing identity + state-transition/CAS law;
5. due-scan and eligibility index bounds;
6. owner/consumer/tenant isolation;
7. rollback/forward recovery that cannot resurrect quarantined or P04.05-completed mutations;
8. explicit exclusion of P04.07+, provider/business schema and AI/model runtime.

No migration version 3 is reserved by activation or this continuity carrier.

## Exact continuity budget

Issue `#220` authorizes reconciliation only in:

1. `AGENTS.md`
2. `README.md`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/handoffs/P04.06.md`
7. `docs/governance/P04_05_P04_06_TRANSITION_TRANSACTION.md`

Do not mutate `docs/roadmap/STATE.json`, package sequence, activation validator, active worker plan, kernel Go runtime, schema/migrations, retry runtime, provider/broker integration, P04.07+, business features, strategic X runtime or AI/model/agent runtime.

## Locked scope on this carrier

Do not implement:

- P04.06 runtime or schema;
- retry worker/scheduler;
- provider/broker-native retry or DLQ integration;
- broad operator requeue/replay/release UX;
- P04.07 schema-registry runtime;
- P04.08 background-job changes;
- P04.09/P04.10 recovery/replay runtime;
- business handlers/features;
- strategic X runtime;
- AI/model/agent product runtime;
- global ordering or end-to-end exactly-once claims.

## Exact next action

1. Keep this continuity source rooted at accepted protected `main@4c9f60843f2612bc4c9a10b4efca7b6a20826be3`.
2. Verify the exact seven-file changed-path budget and absence of runtime/schema changes.
3. Require fresh exact-head Omnexa Governance.
4. Record honest SELF REVIEW and require zero unresolved threads.
5. Verify protected-main freshness.
6. Promote the exact unchanged continuity source head through a separate promotion PR.
7. Require fresh promotion-specific Governance and promotion SELF REVIEW with zero unresolved threads.
8. Merge with expected-head protection and re-read protected main.
9. Confirm P04.06 sole ACTIVE at `5 / 10` with zero runtime leases.
10. Only then create a fresh isolated P04.06 runtime wave and perform migration preflight.

Do not auto-advance P04.07.
