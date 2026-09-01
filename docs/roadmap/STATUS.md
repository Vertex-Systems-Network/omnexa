# Omnexa Roadmap Status

Last reconciled: 2026-09-01 — **P04.04 ACTIVE after accepted P04.03 closure / P04.04 activation merge/read-back**

## Authoritative checkpoint

Protected main: `50edeae03ad52e435d142b2f22c803f08a5c7f1a`.

Canonical state:

- Foundation Architecture v1 is FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED.
- P04 is ACTIVE — **3 / 10 done** (canonical).
- P04.01-P04.03 are DONE with retained accepted evidence.
- P04.04 is the **sole ACTIVE package**.
- P04.05-P04.10 remain planned/locked.
- `kernel_code_authorized=true` for P04.04 only.
- `business_feature_code_authorized=false`.
- canonical CI is GitHub-hosted `ubuntu-24.04` only.

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This continuity file grants no authority beyond that state.

## Accepted P04.03 implementation/completion evidence

- source implementation PR #165;
- final implementation head `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`;
- source Governance `33377793927 / 99443112098` — PASS;
- unchanged promotion PR #166;
- promotion Governance `33405463251 / 99531835998` — PASS;
- implementation merge/read-back `b94189873bef11f4870935205398f1ef44f160bf`;
- completion evidence `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`;
- completion evidence merge/read-back `ed9c9b067c2725e9ddef4c3a2b03c4aa0b29dbcd`.

## Accepted P04.04 preparation evidence

- preparation source PR #168;
- preparation head `658fbb4074f76c7960680c7e04bfd7cbc07a59a9`;
- source Governance `33409859631 / 99546416425` — PASS;
- unchanged preparation promotion PR #169;
- promotion Governance `33410873382 / 99549818529` — PASS;
- preparation merge/read-back `962a62c7c111079ca6f2047fa748deea97c84534`;
- contract `docs/roadmap/work-packages/P04.04.md`;
- handoff `docs/ai/handoffs/P04.04.md`.

## Accepted P04.03 closure / P04.04 activation evidence

- source PR #171;
- exact source/promotion head `2cb7a92217a81b921e40a274007e1c50dc4a8fcf`;
- source Governance run/job `33546204586 / 99984199589` — PASS;
- source unresolved review threads: 0; SELF REVIEW recorded;
- unchanged promotion PR #172;
- promotion Governance run/job `33547080086 / 99987132610` — PASS;
- promotion unresolved review threads: 0; SELF REVIEW recorded;
- expected-head guarded merge/read-back `50edeae03ad52e435d142b2f22c803f08a5c7f1a`.

The accepted transaction contained no P04.04 runtime source, migration, concrete broker/provider, business handler, P04.05+ behavior or AI/model/agent runtime.

## P04.04 runtime boundary

Owner: `kernel.events`.

Authorized runtime scope after this continuity correction lands:

- provider-neutral transactional outbox reliability;
- canonical P04.01 event and authorized local owner mutation commit/rollback in the same local PostgreSQL transaction;
- reuse P01 `database.InTransaction`;
- durable committed-pending outbox persistence and restart recovery;
- relay through the accepted P04.02 publication abstraction;
- publication failure remains pending/recoverable;
- publish success followed by mark failure/crash may duplicate the same canonical event;
- published state means producer-side publication progress only;
- concurrent relay attempts must not corrupt/rebind owner/event/tenant state;
- no global ordering guarantee.

Still unauthorized:

- concrete broker/provider selection;
- end-to-end exactly-once claims;
- P04.05 inbox/deduplication;
- P04.06 retry/backoff/DLQ/quarantine;
- P04.07 schema-registry runtime;
- P04.08 background-job runtime changes;
- business features/handlers;
- strategic X-program runtime;
- AI/model/agent runtime.

## Persistence / migration preflight

Production P04.04 semantics require durable local PostgreSQL persistence, but activation itself added no schema change. Before the first runtime schema mutation, a fresh implementation branch must record the exact:

- owner: `kernel.events`;
- migration directory/path and next immutable version;
- tables/indexes/constraints/data budget;
- fresh-install and upgrade behavior;
- rollback/forward-recovery strategy;
- owner and tenant isolation proof.

No second migration runner/ledger or cross-owner SQL is authorized.

## Exact next work

1. Finish this **post-activation continuity-only** source/promotion Governance path without changing canonical `STATE.json`, package sequence, validator or runtime source.
2. Re-read protected main after continuity merge and confirm P04.04 remains sole active at 3 / 10.
3. Create a **new isolated P04.04 implementation branch** from that exact protected-main SHA.
4. Re-read the accepted P04.04 contract and repository migration/event foundations; record exact runtime + migration path budget before source mutation.
5. Implement only P04.04 with focused atomicity, rollback, restart, duplicate-relay, concurrency, tenant/owner isolation and payload-safe failure tests.
6. Do not auto-advance to P04.05.
