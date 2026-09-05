# P04.05 Closure / P04.06 Activation Transaction

Status: **ACCEPTED / PROTECTED-MAIN READ-BACK COMPLETE — POST-ACTIVATION CONTINUITY REQUIRED BEFORE RUNTIME**

- Original coordination issue: `#217`
- Transition source base: protected `main@3f547180eb5e839439834eb2ce7977324803df18`
- Source PR: `#218`
- Unchanged promotion PR: `#219`
- Exact source/promotion head: `b718ad7316dba6fca0cafccb514df6da653abe13`
- Source Governance: `#677 / 33990309561` — PASS
- Promotion Governance: `#678` — PASS
- Accepted protected-main merge/read-back: `4c9f60843f2612bc4c9a10b4efca7b6a20826be3`
- Accepted package result: `P04.01-P04.05 DONE`, `P04.06 sole ACTIVE`, `P04.07-P04.10 PLANNED / LOCKED`
- Runtime/schema/provider mutation in the transaction: `none`
- Business-feature authority: `false`
- AI/model/agent product-runtime authority: `false`
- Post-activation continuity issue: `#220`

The transition is accepted canonical history. It changed governance/state/continuity only. It did not implement retry/backoff runtime, create a retry worker, reserve a migration, create a quarantine table, select a broker/provider, add a provider-native DLQ, implement P04.07+, add a business feature, or authorize AI/model/agent product runtime.

The separate issue `#220` continuity carrier now reconciles subordinate stale pre-activation wording. No P04.06 runtime or schema mutation may occur on that continuity source or its unchanged promotion.

## Preconditions satisfied

Protected main at the transition base contained both required predecessor facts:

1. accepted P04.05 implementation and completion evidence;
2. separately prepared and governed P04.06 contract/handoff.

The package did not auto-advance merely because those files existed. Source `#218` and unchanged promotion `#219` were the explicit sequential closure/activation decision.

## Accepted P04.05 implementation/completion evidence

The accepted P04.05 implementation chain culminates in:

- T04 Supervisor source PR `#210`;
- unchanged T04 promotion PR `#211`;
- exact T04 source/promotion head `dd713fe3217a0d092ab3ff31115ac031ae8c0303`;
- source Governance `#669` / run `33927932705` — PASS;
- promotion Governance `#670` / run `33985111334` / job `101357077040` — PASS;
- expected-head implementation merge/read-back `0c66a3371dbf2fa942a95b7d0475b06235392474`;
- completion evidence `docs/roadmap/evidence/P04.05_COMPLETION_2026-09-05.md`;
- completion evidence PR `#213`;
- evidence Governance `#672` — PASS;
- completion evidence merge/read-back `e44ece77ddf7b821c03997266ca0c68c07162910`.

P04.05 accepted only the local PostgreSQL consumer inbox/idempotency boundary. Checkpoint progress and inbox completion remain separate facts and external/non-transactional side effects are not made end-to-end exactly once.

## Accepted P04.06 preparation evidence

P04.06 preparation was accepted before activation:

- source preparation PR `#215`;
- unchanged preparation promotion PR `#216`;
- exact source/promotion head `7babb9c39185636b3af5184d5a7bd31cedbc37a0`;
- source Governance `#674` / run `33987003924` — PASS;
- promotion Governance `#675` / run `33987472967` — PASS;
- preparation merge/read-back `3f547180eb5e839439834eb2ce7977324803df18`;
- accepted contract `docs/roadmap/work-packages/P04.06.md`;
- handoff `docs/ai/handoffs/P04.06.md`.

Preparation selected no provider and reserved no schema/migration.

## Accepted canonical result

Protected `main@4c9f60843f2612bc4c9a10b4efca7b6a20826be3` establishes:

- P04 remains ACTIVE;
- P04 progress is `5 / 10 done`;
- P04.01-P04.05 are DONE with retained evidence;
- P04.06 is the sole ACTIVE package;
- P04.07-P04.10 remain PLANNED / LOCKED;
- `kernel_code_authorized=true` is bounded to P04.06 only;
- P04.06 runtime implementation remains gated by required post-activation continuity acceptance/read-back;
- `business_feature_code_authorized=false`;
- strategic X runtime remains unauthorized;
- AI/model/agent product runtime remains unauthorized.

## P04.05 wave release

The completed P04.05 wave is historical after the accepted transition:

- `P04.05-T01` inbox core — released;
- `P04.05-T02` PostgreSQL inbox persistence / migration v2 — released;
- `P04.05-T03` reliability/concurrency proof — released;
- `P04.05-T04` Supervisor verifier — released.

The accepted `kernel.events` migration version 2 remains immutable historical ownership evidence. It is **not** a live migration reservation for P04.06.

The accepted transition created:

- zero P04.06 worker slots;
- zero P04.06 runtime tasks;
- zero P04.06 runtime branches;
- zero P04.06 migration reservations.

The repository writer cap remains three worker writers maximum when a later governed implementation wave exists.

## Active P04.06 boundary

P04.06 is active, but implementation is permitted only after the required post-activation continuity gate. The prepared provider-neutral contract includes:

- stable failure disposition rather than raw error-string policy;
- finite deterministic retry attempts and capped backoff;
- authoritative UTC retry eligibility;
- one-at-a-time claim/CAS/lease semantics;
- strict separation of P04.03 checkpoint, P04.05 inbox and P04.06 retry/quarantine evidence;
- stale retry suppression after P04.05 already-applied completion;
- fail-closed owner/consumer/route/tenant rebinding;
- terminal/exhausted quarantine that must become durable before checkpoint may advance past a poison event;
- restart recovery when quarantine commits but checkpoint advancement fails;
- bounded safe quarantine metadata with no raw payload/secrets/provider/database diagnostics;
- logical local quarantine only, with no broker-native DLQ/provider selection.

The accepted activation itself did not implement those semantics, and the post-activation continuity carrier must not implement them either.

## Persistence / migration gate remains closed

No P04.06 schema is reserved by activation or continuity. Before a later runtime wave can mutate schema it must, from fresh protected main after post-activation continuity acceptance:

1. re-read `kernel.events` migration history;
2. identify the exact next immutable migration version/path/name;
3. record the exact table/index/constraint/data budget;
4. define attempt/quarantine state-transition and CAS law;
5. prove tenant/owner/consumer isolation and bounded due-scan/index behavior;
6. define forward recovery that cannot resurrect quarantined or already-completed protected mutations;
7. confirm P04.07+, provider/business and AI runtime are absent.

No `kernel.events` version 3 is reserved.

## Post-activation continuity transaction

Issue `#220` is the sole current coordination carrier for subordinate post-activation wording reconciliation.

Exact allowed changed-path budget:

1. `AGENTS.md`
2. `README.md`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/handoffs/P04.06.md`
7. `docs/governance/P04_05_P04_06_TRANSITION_TRANSACTION.md`

Hard non-scope:

- `docs/roadmap/STATE.json`;
- package sequence;
- activation validator;
- active worker plan;
- kernel Go runtime;
- migration/schema reservation;
- retry worker/scheduler;
- provider/broker-native DLQ;
- P04.07+ runtime;
- business features;
- strategic X runtime;
- AI/model/agent product runtime.

## Required continuity integration path

1. Create the continuity source from exact accepted protected `main@4c9f60843f2612bc4c9a10b4efca7b6a20826be3`.
2. Reconcile only the exact seven allowed subordinate files.
3. Require fresh exact-head Omnexa Governance.
4. Record substantive SELF REVIEW with honest non-independent provenance and zero unresolved threads.
5. Verify protected main still equals the expected base before promotion.
6. Create a separate promotion branch at the exact unchanged continuity source head.
7. Require fresh promotion-specific Omnexa Governance.
8. Record promotion SELF REVIEW, zero unresolved threads and exact unchanged-head verification.
9. Merge with expected-head protection.
10. Re-read protected main, canonical `STATE.json`, package sequence and active-plan authority.
11. Confirm P04.01-P04.05 remain DONE, P04.06 remains sole ACTIVE at `5 / 10`, P04.07-P04.10 remain locked and P04.06 runtime leases/reservations remain zero.
12. Only then create a fresh P04.06 implementation wave and conduct fresh migration preflight before any schema reservation.

Do not auto-advance P04.07.
