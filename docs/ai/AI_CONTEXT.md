# Omnexa AI Project Context

Status: **P04.03 closure / P04.04 activation candidate — subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative checkpoint before this transaction

Fresh protected main is `962a62c7c111079ca6f2047fa748deea97c84534`.

At that exact protected-main checkpoint:

- Foundation Architecture v1 is FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED / historical prerequisite.
- P04 is ACTIVE — 2 / 10 canonically done.
- P04.01 and P04.02 are DONE.
- P04.03 is still the canonical active package, but its implementation/completion evidence is accepted.
- P04.04 contract/handoff preparation is accepted but planned/locked.
- P04.05-P04.10 are planned/locked.
- `business_feature_code_authorized=false`.
- strategic X-program runtime remains unauthorized.
- AI/model/agent runtime remains unauthorized.

This carrier proposes the state cursor change only. It contains no P04.04 runtime code or migration.

## Accepted P04.03 chain

P04.03 implementation acceptance is immutable predecessor evidence:

- source implementation PR #165;
- final exact source/promotion head `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`;
- source Governance #581 / run `33377793927` / job `99443112098` — PASS;
- unchanged promotion PR #166;
- promotion Governance #582 / run `33405463251` / job `99531835998` — PASS;
- accepted implementation merge/read-back `b94189873bef11f4870935205398f1ef44f160bf`;
- completion evidence `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`;
- completion-evidence carrier PR #167;
- completion-evidence merge/read-back `ed9c9b067c2725e9ddef4c3a2b03c4aa0b29dbcd`.

P04.03 retained laws include owner/consumer/tenant/scope-bound checkpoint progress, contiguous monotonic advancement, fail-closed stale/regressive/gapped/conflicting updates, no progress on failed/cancelled handling, restart from accepted checkpoint, explicit duplicate replay around crash/write-failure windows, race-safe checkpoint advancement, no global ordering guarantee and no exactly-once business-mutation guarantee.

P04.03 selected no concrete broker, no production checkpoint provider and no database migration.

## Accepted P04.04 preparation chain

The P04.04 contract/handoff was prepared only after P04.03 implementation/completion evidence existed:

- source preparation PR #168;
- exact preparation head `658fbb4074f76c7960680c7e04bfd7cbc07a59a9`;
- source Governance #587 / run `33409859631` / job `99546416425` — PASS;
- unchanged preparation promotion PR #169;
- promotion Governance #588 / run `33410873382` / job `99549818529` — PASS;
- preparation merge/read-back `962a62c7c111079ca6f2047fa748deea97c84534`;
- contract `docs/roadmap/work-packages/P04.04.md`;
- handoff `docs/ai/handoffs/P04.04.md`.

Preparation alone does not authorize P04.04 implementation.

## Candidate activation result

Only after the exact source transaction passes Governance, is promoted unchanged through a fresh promotion-specific carrier, that promotion passes fresh Governance, merges through protected main and is read back, the candidate becomes canonical:

- P04 ACTIVE — 3 / 10 done;
- P04.01-P04.03 DONE with accepted evidence;
- P04.04 sole ACTIVE;
- P04.05-P04.10 planned/locked;
- `kernel_code_authorized=true` for P04.04 only;
- `business_feature_code_authorized=false`;
- strategic X-program runtime unauthorized;
- AI/model/agent runtime unauthorized.

Do not treat this candidate snapshot as authority before protected-main read-back.

## P04.04 implementation boundary after activation

Owner: `kernel.events`.

A later fresh implementation branch may implement only the accepted P04.04 contract:

- provider-neutral transactional outbox reliability;
- authoritative owner mutation and canonical P04.01 event envelope committed/rolled back together in the same local PostgreSQL transaction;
- reuse of P01 `database.InTransaction`;
- durable committed-pending outbox state recoverable after restart;
- relay through the accepted P04.02 publication boundary;
- publish failure leaves work pending/recoverable;
- publish-success/crash-before-published-mark can produce duplicate publication of the same canonical event;
- published state means producer-side publication progress only, not consumer/business exactly-once completion;
- concurrent relay attempts cannot corrupt or cross-bind owner/event/tenant state;
- no global ordering guarantee from database ordering or relay scheduling;
- focused positive/adversarial tests plus a dedicated P04.04 verifier/evidence path.

P04.04 must not select Kafka/NATS/RabbitMQ/Redis Streams or another broker, claim end-to-end exactly-once semantics, implement P04.05 inbox/deduplication, P04.06 retry/DLQ, P04.07 schema registry, P04.08 background-job execution changes, business handlers/features, P04.05+, strategic X-program runtime or AI/model/agent runtime.

## Persistence / migration decision

The P04.04 contract requires durable local PostgreSQL outbox persistence for the production reliability claim, so persistence is expected.

The **activation carrier itself adds no migration and grants no blanket schema authority**.

Before the first schema mutation on the later P04.04 implementation branch:

1. re-read fresh protected main;
2. inspect the current `kernel.events` migration/version namespace;
3. record the exact migration path/version/data budget;
4. reuse the retained P01 migration runner and immutable migration identity;
5. prove fresh-install and supported-upgrade behavior;
6. prove rollback/forward-recovery;
7. prove tenant/owner isolation and safe failure handling.

No cross-owner database authority follows from P04.04 activation.

## Sequential activation validator

`scripts/validate_p04_activation.py` remains generic for strict sequential activation and now additionally requires the accepted P04.04 specification to exist and retain its core outbox/no-exactly-once/no-global-ordering/no-automatic-migration laws.

For every completed P04 predecessor, the validator requires:

- `done` state in the package sequence;
- accepted specification;
- completion-evidence path;
- matching `governance_tracking.p04_NN_completion.state=PASS`;
- exact implementation head and merge;
- canonical workflow run and job;
- GitHub-hosted `ubuntu-24.04` evidence.

P04.03 completion tracking is therefore a fail-closed prerequisite for candidate active P04.04.

## Active-package execution rule after accepted transition

Before material P04.04 runtime work:

1. re-read protected `main`, `AGENTS.md`, canonical `STATE.json`, `STATUS.md`, P04 sequence, P04.04 spec and handoff;
2. confirm P04.04 is the sole active package and P04.05-P04.10 remain locked;
3. confirm `kernel_code_authorized=true` only for P04.04 and `business_feature_code_authorized=false`;
4. create a new isolated implementation branch from the exact verified protected-main SHA;
5. record the exact owner/path/version/data migration budget before any schema mutation;
6. preserve P04.01/P04.02/P04.03 identity, ownership, duplicate and crash/restart semantics;
7. add focused atomicity, crash-window, restart, duplicate-relay, concurrency, tenant/owner isolation and malformed-state evidence;
8. require exact-final-head GitHub-hosted Governance and the governed promotion path;
9. do not auto-advance to P04.05.

## Historical prerequisites retained

P01.01-P01.12, P02.01-P02.10, P03.01-P03.11 and P04.01-P04.03 retain immutable completion evidence and mandatory regressions. Historical diagnostic failures remain diagnostic and are never rewritten as PASS.

## Exact next action

Finish the P04.03 closure / P04.04 activation source carrier, require exact-head Governance, review the exact diff, promote the unchanged source state through a fresh promotion-specific carrier, require fresh promotion Governance, merge with expected-head protection and re-read protected main. Confirm P04.04 is sole active at 3 / 10 and then **STOP**; only a later new branch from that exact main may begin P04.04 runtime implementation.
