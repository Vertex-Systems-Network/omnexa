# Omnexa Roadmap Status

Last reconciled: 2026-09-06 — **P04.05 CLOSURE / P04.06 ACTIVATION CANDIDATE**

## Authoritative candidate base

Protected main at transition start:

`3f547180eb5e839439834eb2ce7977324803df18`

Canonical state is changed only by the separately governed activation transaction. Until the exact activation source passes Governance, is promoted unchanged through fresh Governance, merges with expected-head protection and protected main is re-read, the candidate result is not accepted main authority.

Candidate result:

- Foundation Architecture v1 remains FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED.
- P04 remains ACTIVE.
- P04 candidate progress becomes **5 / 10 done**.
- P04.01-P04.05 are DONE with retained accepted evidence.
- P04.06 becomes the **sole ACTIVE package** only after accepted merge/read-back.
- P04.07-P04.10 remain planned/locked.
- `kernel_code_authorized=true` is bounded to P04.06 only after accepted activation + required post-activation continuity.
- `business_feature_code_authorized=false`.
- strategic X-program and AI/model/agent product runtime remain unauthorized.
- canonical CI remains GitHub-hosted `ubuntu-24.04` only.

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This status file is subordinate and cannot activate a package by itself.

## Accepted P04.05 implementation/completion

P04.05 is accepted through the bounded T01-T04 implementation chain. Terminal implementation evidence:

- final Supervisor source/promotion head `dd713fe3217a0d092ab3ff31115ac031ae8c0303`;
- source PR `#210` / promotion PR `#211`;
- source Governance `#669 / 33927932705` — PASS;
- promotion Governance `#670 / 33985111334`, job `101357077040` — PASS;
- implementation merge/read-back `0c66a3371dbf2fa942a95b7d0475b06235392474`;
- completion evidence `docs/roadmap/evidence/P04.05_COMPLETION_2026-09-05.md`;
- completion evidence PR `#213` / Governance `#672` — PASS;
- completion evidence merge/read-back `e44ece77ddf7b821c03997266ca0c68c07162910`.

Accepted P04.05 semantics remain local application-idempotency only. P04.03 checkpoint progress and P04.05 inbox completion are separate facts; external/non-transactional side effects are not made end-to-end exactly once.

## Accepted P04.06 preparation

P04.06 was separately prepared before this activation transaction:

- source preparation PR `#215`;
- unchanged promotion PR `#216`;
- exact preparation head `7babb9c39185636b3af5184d5a7bd31cedbc37a0`;
- source Governance `#674 / 33987003924` — PASS;
- promotion Governance `#675 / 33987472967` — PASS;
- preparation merge/read-back `3f547180eb5e839439834eb2ce7977324803df18`;
- contract `docs/roadmap/work-packages/P04.06.md`;
- handoff `docs/ai/handoffs/P04.06.md`.

Preparation did not mutate canonical package state, reserve schema, select a provider or create a runtime wave.

## P04.05 closure / P04.06 activation candidate

Coordination issue: `#217`.

The activation transaction is governance/state/continuity only. It reconciles:

- P04.05 from ACTIVE to DONE and attaches its accepted completion evidence;
- P04.06 from PLANNED to sole ACTIVE and attaches its separately prepared specification;
- P04 done count from 4 to 5;
- bounded kernel authority from P04.05 to P04.06 only;
- old P04.05 worker/supervisor leases to historical/released state;
- active runtime orchestration to zero worker slots/tasks/migration reservations.

The activation carrier itself contains no retry worker, scheduler, quarantine persistence, migration, provider, broker-native DLQ, P04.07+ runtime, business feature or AI/model product runtime.

## Activated P04.06 boundary

After accepted activation **and a later separately governed post-activation continuity carrier**, a fresh implementation wave may implement only the prepared provider-neutral P04.06 scope:

- stable structured failure disposition;
- finite retry attempts and capped deterministic backoff;
- authoritative UTC retry eligibility;
- one-at-a-time retry claim/CAS/lease behavior;
- P04.03 checkpoint, P04.05 inbox and P04.06 retry/quarantine as distinct facts;
- stale retry suppression when P04.05 says already applied;
- owner/consumer/route/tenant fail-closed identity binding;
- terminal/exhausted logical quarantine;
- quarantine persistence before checkpoint may advance past poison delivery;
- restart recovery after quarantine commit/checkpoint gap without handler reinvocation;
- bounded safe diagnostic evidence with no raw payload/secrets/provider/database diagnostics.

P01.09 in-memory job retry remains precedent only and is not durable event-fabric retry persistence/timing authority.

## Runtime and migration remain blocked during activation

At activation candidate time:

- P04.06 worker slots: `0`;
- P04.06 runtime tasks: `0`;
- P04.06 runtime branches: `0`;
- P04.06 migration reservations: `0`;
- provider/broker choice: `none`.

Accepted `kernel.events` migration v2 remains historical P04.05 inbox ownership and is not a live lease.

Before any later P04.06 schema mutation, fresh protected main after post-activation continuity must be re-read and the exact next immutable `kernel.events` migration path/version/data budget must be recorded and governed. No version 3 is reserved here.

## Still unauthorized

- runtime mutation on this activation source/promotion;
- P04.07 schema registry/runtime;
- P04.08 background-job ownership changes;
- P04.09 broad operator recovery UX;
- P04.10 controlled replay/poison aggregate runtime;
- provider-specific retry headers or broker-native DLQ integration;
- global ordering or end-to-end exactly-once claims;
- production business handlers/features;
- strategic X runtime;
- AI/model/agent product runtime.

## Exact next work

1. Complete this activation source's canonical state reconciliation.
2. Require exact-head Omnexa Governance and honest SELF REVIEW with zero unresolved threads.
3. Verify protected-main freshness and exact changed-path budget.
4. Create an unchanged promotion of the exact reviewed head.
5. Require fresh promotion Governance, review/thread cleanliness and expected-head freshness.
6. Merge/read back protected main.
7. Confirm P04.01-P04.05 DONE, P04.06 sole ACTIVE, P04.07+ locked, and zero runtime leases.
8. Govern/promote/read back a **separate P04.06 post-activation continuity carrier**.
9. Only then create a fresh P04.06 implementation wave and conduct fresh migration preflight.
10. Do not auto-advance P04.07.
