# P04.05 Closure / P04.06 Activation Transaction

Status: **CANDIDATE — EXACT-HEAD GOVERNANCE / UNCHANGED PROMOTION / PROTECTED READ-BACK REQUIRED**

- Coordination issue: `#217`
- Transaction base: protected `main@3f547180eb5e839439834eb2ce7977324803df18`
- Candidate package result: `P04.01-P04.05 DONE`, `P04.06 sole ACTIVE`, `P04.07-P04.10 PLANNED / LOCKED`
- Runtime/schema/provider mutation in this transaction: `none`
- Business-feature authority: `false`
- AI/model/agent product-runtime authority: `false`

This transaction changes governance/state/continuity only. It does not implement retry/backoff runtime, create a retry worker, reserve a migration, create a quarantine table, select a broker/provider, add a provider-native DLQ, implement P04.07+, add a business feature, or authorize AI/model/agent product runtime.

## Preconditions

Protected main already contains both required predecessor facts:

1. accepted P04.05 implementation and completion evidence;
2. separately prepared and governed P04.06 contract/handoff.

No package is allowed to auto-advance merely because those files exist. This transaction is the explicit sequential closure/activation decision.

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

P04.05 accepted only the local PostgreSQL consumer inbox/idempotency boundary. It retained the explicit limits that checkpoint progress and inbox completion are separate facts and external/non-transactional side effects are not made end-to-end exactly once.

## Accepted P04.06 preparation evidence

P04.06 preparation is accepted but not runtime-active at the transaction base:

- source preparation PR `#215`;
- unchanged preparation promotion PR `#216`;
- exact source/promotion head `7babb9c39185636b3af5184d5a7bd31cedbc37a0`;
- source Governance `#674` / run `33987003924` — PASS;
- promotion Governance `#675` / run `33987472967` — PASS;
- preparation merge/read-back `3f547180eb5e839439834eb2ce7977324803df18`;
- accepted contract `docs/roadmap/work-packages/P04.06.md`;
- handoff `docs/ai/handoffs/P04.06.md`.

Preparation selected no provider and reserved no schema/migration.

## Candidate canonical result

Only after this exact source head passes fresh Governance, is promoted byte-for-byte unchanged through a second fresh Governance run, merges with expected-head protection and protected main is re-read, canonical state may be interpreted as:

- P04 remains ACTIVE;
- P04 progress becomes `5 / 10 done`;
- P04.01-P04.05 are DONE with retained evidence;
- P04.06 is the sole ACTIVE package;
- P04.07-P04.10 remain PLANNED / LOCKED with no accepted runtime spec;
- `kernel_code_authorized=true` is bounded to P04.06 only;
- `business_feature_code_authorized=false`;
- strategic X runtime remains unauthorized;
- AI/model/agent product runtime remains unauthorized.

## P04.05 wave release

The completed P04.05 wave is historical after this transition:

- `P04.05-T01` inbox core — released;
- `P04.05-T02` PostgreSQL inbox persistence / migration v2 — released;
- `P04.05-T03` reliability/concurrency proof — released;
- `P04.05-T04` Supervisor verifier — released.

The accepted `kernel.events` migration version 2 remains immutable historical ownership evidence. It is **not** a live migration reservation for P04.06.

This transaction intentionally creates:

- zero P04.06 worker slots;
- zero P04.06 runtime tasks;
- zero P04.06 runtime branches;
- zero P04.06 migration reservations.

The repository writer cap remains three worker writers maximum when a later governed implementation wave exists.

## Activated P04.06 boundary

Activation authorizes only the prepared provider-neutral P04.06 package after the required post-activation continuity gate. The prepared contract includes:

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

Activation itself does not implement those semantics.

## Persistence / migration gate remains closed

No P04.06 schema is reserved here. Before a later runtime wave can mutate schema it must, from fresh protected main after post-activation continuity acceptance:

1. re-read `kernel.events` migration history;
2. identify the exact next immutable migration version/path/name;
3. record the exact table/index/constraint/data budget;
4. define attempt/quarantine state-transition and CAS law;
5. prove tenant/owner/consumer isolation and bounded due-scan/index behavior;
6. define forward recovery that cannot resurrect quarantined or already-completed protected mutations;
7. confirm P04.07+, provider/business and AI runtime are absent.

No `kernel.events` version 3 is reserved by this transaction.

## Required integration path

1. Source activation PR on the exact candidate head.
2. Fresh exact-head Omnexa Governance.
3. Substantive SELF REVIEW with honest non-independent provenance; zero unresolved threads.
4. Verify protected main still equals the expected base.
5. Create a separate promotion branch at the exact unchanged source head.
6. Fresh promotion-specific Omnexa Governance.
7. Promotion SELF REVIEW, zero unresolved threads and exact unchanged-head verification.
8. Expected-head guarded merge.
9. Re-read protected main, canonical STATE, package sequence and active-plan authority.
10. Confirm P04.06 is sole ACTIVE at `5 / 10`, with zero runtime leases.
11. Govern a separate P04.06 post-activation continuity carrier.
12. Only after that continuity source + unchanged promotion are accepted/read back may a fresh P04.06 implementation wave or migration preflight/reservation be created.

Do not auto-advance P04.07.
