# Omnexa AI Project Context

Status: **P04.05 CLOSURE / P04.06 ACTIVATION CANDIDATE — subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative candidate base

Fresh protected main at this activation candidate start is:

`3f547180eb5e839439834eb2ce7977324803df18`

That protected-main state contains:

- P04.05 implementation acceptance through its governed T01-T04 chain;
- P04.05 completion evidence merged/read back as `e44ece77ddf7b821c03997266ca0c68c07162910`;
- the separately governed P04.06 preparation source/promotion at exact head `7babb9c39185636b3af5184d5a7bd31cedbc37a0`;
- P04.06 preparation source Governance `#674 / 33987003924` — PASS;
- P04.06 preparation promotion Governance `#675 / 33987472967` — PASS;
- P04.06 preparation merge/read-back `3f547180eb5e839439834eb2ce7977324803df18`.

Preparation alone did not activate P04.06. This separate transaction is the explicit sequential package transition.

## Candidate canonical result

Only after this exact activation source passes fresh Governance, is reviewed, promoted byte-for-byte unchanged through fresh promotion Governance, merged with expected-head protection and protected main is re-read may the candidate result be treated as accepted:

- P04 ACTIVE — `5 / 10 done`;
- P04.01-P04.05 DONE with retained evidence;
- P04.06 sole ACTIVE package;
- P04.07-P04.10 planned/locked;
- bounded kernel authority only for P04.06;
- business-feature authority remains false;
- strategic X and AI/model/agent product runtime remain unauthorized.

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

## P04.06 activation boundary

Owner: `kernel.events`.

After activation **and a later separate post-activation continuity source/promotion/read-back**, a fresh implementation wave may implement only the accepted provider-neutral P04.06 contract:

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

## Zero-runtime activation law

This activation carrier intentionally has:

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

No migration version 3 is reserved by activation.

## Locked scope

Do not implement on this activation source/promotion:

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

1. Finish canonical state/continuity reconciliation on this source branch.
2. Require fresh exact-head Omnexa Governance.
3. Record honest SELF REVIEW and require zero unresolved threads.
4. Verify protected main freshness.
5. Promote the exact unchanged source head through a separate promotion PR.
6. Require fresh promotion-specific Governance and zero threads.
7. Merge with expected-head protection and re-read protected main.
8. Confirm P04.06 sole ACTIVE at `5 / 10` with zero runtime leases.
9. Govern/promote/read back a separate P04.06 post-activation continuity carrier.
10. Only then create a fresh isolated P04.06 runtime wave and perform migration preflight.

Do not auto-advance P04.07.
