# P04.03 Closure / P04.04 Activation Transaction

Status: **ACCEPTED — MERGED AND READ BACK FROM PROTECTED MAIN**

Transaction base: `6b8cf73bbb8d6c9bc977b69d9a110f01e706b221`  
Exact source/promotion head: `2cb7a92217a81b921e40a274007e1c50dc4a8fcf`  
Accepted merge/read-back: `50edeae03ad52e435d142b2f22c803f08a5c7f1a`

This transaction changed governance/state/continuity plus narrowly required activation-validator hardening only. It contained no P04.04 runtime source, migration, broker/provider selection, business handler, background-job change or AI/model/agent runtime.

## Accepted transition evidence

- source PR: `#171`;
- source branch: `agent/20260901-p04-03-close-p04-04-activate-v2`;
- exact source head: `2cb7a92217a81b921e40a274007e1c50dc4a8fcf`;
- source Governance run/job: `33546204586 / 99984199589` — PASS;
- source review: SELF REVIEW; independent approval not claimed;
- source unresolved review threads: 0;
- unchanged promotion PR: `#172`;
- promotion branch: `promotion/20260901-p04-03-close-p04-04-activate`;
- promotion exact head: `2cb7a92217a81b921e40a274007e1c50dc4a8fcf`;
- promotion Governance run/job: `33547080086 / 99987132610` — PASS;
- promotion review: SELF REVIEW; independent approval not claimed;
- promotion unresolved review threads: 0;
- expected-head guarded merge/read-back: `50edeae03ad52e435d142b2f22c803f08a5c7f1a`.

## Preconditions retained

At the transaction base:

- P04 was ACTIVE at 2 / 10;
- P04.01-P04.02 were DONE;
- P04.03 was the sole active package with accepted implementation/completion evidence;
- P04.04 contract/handoff preparation was accepted but locked;
- P04.05-P04.10 were locked;
- business-feature and AI/model/agent runtime were unauthorized.

## Accepted predecessor evidence — P04.03

- source implementation PR #165;
- promotion implementation PR #166;
- final implementation head `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`;
- source Governance `33377793927 / 99443112098`;
- promotion Governance `33405463251 / 99531835998`;
- implementation merge/read-back `b94189873bef11f4870935205398f1ef44f160bf`;
- completion evidence `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`;
- completion evidence merge/read-back `ed9c9b067c2725e9ddef4c3a2b03c4aa0b29dbcd`.

## Accepted preparation evidence — P04.04

- source preparation PR #168;
- preparation head `658fbb4074f76c7960680c7e04bfd7cbc07a59a9`;
- source Governance `33409859631 / 99546416425`;
- unchanged promotion PR #169;
- promotion Governance `33410873382 / 99549818529`;
- preparation merge/read-back `962a62c7c111079ca6f2047fa748deea97c84534`;
- accepted contract `docs/roadmap/work-packages/P04.04.md`;
- handoff `docs/ai/handoffs/P04.04.md`.

## Canonical protected-main result

Protected-main read-back `50edeae03ad52e435d142b2f22c803f08a5c7f1a` establishes:

- P04 ACTIVE — `3 / 10 done`;
- P04.01-P04.03 DONE with retained evidence;
- P04.04 sole ACTIVE package;
- P04.05-P04.10 planned/locked;
- `kernel_code_authorized=true` only for P04.04;
- `business_feature_code_authorized=false`;
- strategic X runtime unauthorized;
- AI/model/agent runtime unauthorized.

Canonical `docs/roadmap/STATE.json` already records that state and is intentionally unchanged by the later post-activation continuity reconciliation.

## Accepted P04.04 implementation boundary

P04.04 may implement only the accepted provider-neutral producer-side transactional outbox reliability contract:

- canonical P04.01 event + authorized owner mutation in the same local PostgreSQL transaction;
- P01 transaction foundation reuse;
- durable committed-pending recovery;
- relay through P04.02;
- publication failure leaves pending state;
- publish-success/crash-before-mark may duplicate;
- producer published state is not consumer/business exactly once;
- concurrency, owner, tenant and canonical identity must fail closed;
- no global ordering.

Activation granted no blanket migration authority. Durable production persistence requires a fresh runtime branch to record the exact `kernel.events` migration path/version/data budget before schema mutation.

## Explicitly excluded

- P04.05+ runtime;
- concrete broker/provider;
- consumer inbox/dedup;
- retry/backoff/DLQ/quarantine;
- schema registry;
- background-job runtime;
- business handlers/features;
- distributed transaction;
- global ordering;
- end-to-end exactly-once claims;
- strategic X runtime;
- AI/model/agent runtime.

## Exact next action

1. Reconcile subordinate post-activation continuity only; do not change canonical `STATE.json`, package sequence, validator or runtime source.
2. Require exact-head source Governance, review, zero unresolved threads, unchanged promotion Governance, expected-head merge and protected-main read-back for that continuity correction.
3. After continuity read-back, create a **fresh P04.04 implementation branch** from the exact protected-main SHA.
4. Record exact runtime/migration path/version/data budget before schema mutation.
5. Implement only P04.04 and do not auto-advance P04.05.
