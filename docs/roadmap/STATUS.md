# Omnexa Roadmap Status

Last reconciled: 2026-09-01 — **P04.03 closure / P04.04 activation candidate; protected main remains authoritative until governed promotion, merge and read-back**

## Authoritative pre-transaction checkpoint

Protected main base: `6b8cf73bbb8d6c9bc977b69d9a110f01e706b221`.

At that exact base:

- Foundation Architecture v1 is FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED.
- P04 is ACTIVE — **2 / 10 canonically done**.
- P04.01 and P04.02 are DONE.
- P04.03 remains the canonical active package on protected main, but its implementation/completion evidence is already accepted.
- P04.04 contract/handoff preparation is accepted on protected main but remains planned/locked.
- P04.05-P04.10 remain planned/locked.
- canonical CI is GitHub-hosted `ubuntu-24.04` only.
- `business_feature_code_authorized=false`.

This carrier changes state/continuity only. It does not contain P04.04 runtime source, migration, provider selection, business logic or AI/model/agent runtime.

## Accepted P04.03 implementation/completion evidence

P04.03 — Durable Stream/Consumer Baseline & Checkpoint Model is accepted predecessor evidence:

- source implementation PR #165;
- exact final implementation head `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`;
- source Governance #581 / run `33377793927` / job `99443112098` — PASS;
- unchanged promotion PR #166;
- promotion Governance #582 / run `33405463251` / job `99531835998` — PASS;
- accepted implementation merge/read-back `b94189873bef11f4870935205398f1ef44f160bf`;
- completion evidence `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`;
- completion-evidence carrier PR #167;
- completion-evidence merge/read-back `ed9c9b067c2725e9ddef4c3a2b03c4aa0b29dbcd`.

Retained P04.03 laws:

- durable progress is owner/consumer/route/tenant/scope bound;
- checkpoint advancement is contiguous and monotonic;
- stale/regressive/gapped/conflicting advancement fails deterministically;
- failed/cancelled handling cannot advance progress;
- restart resumes from the last accepted checkpoint;
- handler success followed by checkpoint-write failure may replay the same event;
- concurrent same-scope work cannot corrupt accepted checkpoint progress;
- checkpoint state is consumption progress only and never authorization;
- no global ordering or exactly-once business-mutation guarantee exists;
- no production checkpoint provider/migration was introduced.

## Accepted P04.04 preparation evidence

P04.04 — Transactional Outbox Reliability Primitive was separately prepared and remains locked until this activation transaction completes:

- preparation source PR #168;
- exact preparation head `658fbb4074f76c7960680c7e04bfd7cbc07a59a9`;
- source Governance #587 / run `33409859631` / job `99546416425` — PASS;
- unchanged preparation promotion PR #169;
- promotion Governance #588 / run `33410873382` / job `99549818529` — PASS;
- preparation merge/read-back `962a62c7c111079ca6f2047fa748deea97c84534`;
- contract `docs/roadmap/work-packages/P04.04.md`;
- handoff `docs/ai/handoffs/P04.04.md`.

Preparation alone grants no runtime authority.

## Candidate atomic result

Only after this exact source state passes Governance, is promoted unchanged through a fresh promotion-specific carrier, that promotion passes fresh Governance, merges through protected main and is read back, the canonical result becomes:

- P04 ACTIVE — **3 / 10 done**;
- P04.01 DONE;
- P04.02 DONE;
- P04.03 DONE with accepted completion evidence;
- P04.04 the **sole ACTIVE package**;
- P04.05-P04.10 planned/locked;
- `kernel_code_authorized=true` for P04.04 only;
- `business_feature_code_authorized=false`;
- strategic X-program runtime unauthorized;
- AI/model/agent runtime unauthorized.

## P04.04 accepted runtime boundary after activation

Owner: `kernel.events`.

A later separate P04.04 implementation branch may implement only the accepted contract:

- provider-neutral transactional outbox reliability;
- authoritative owner mutation and the canonical P04.01 envelope commit/rollback in the same local PostgreSQL transaction;
- reuse of P01 `database.InTransaction`;
- durable committed-pending outbox state recoverable after restart;
- relay through the accepted P04.02 publication abstraction;
- publication failure remains pending/recoverable;
- publish-success/crash-before-published-mark may duplicate the canonical event;
- published state is producer-side publication progress only;
- concurrent relay attempts cannot corrupt or rebind owner/event/tenant state;
- no global ordering guarantee from row order, sequence, timestamps or relay scheduling.

The P04.04 contract requires durable local PostgreSQL outbox persistence for production semantics, so persistence is expected. **This activation carrier adds no migration and grants no blanket schema authority.** Before any runtime schema mutation, the separate implementation branch must re-read fresh protected main and record the exact `kernel.events` migration path/version/data budget under the retained P01 migration rules, including fresh-install, upgrade, rollback/forward-recovery and tenant/owner isolation evidence.

P04.04 does not authorize:

- Kafka, NATS, RabbitMQ, Redis Streams or another concrete broker/provider choice;
- end-to-end exactly-once delivery/business-mutation claims;
- P04.05 inbox/deduplication persistence;
- P04.06 retry/backoff/DLQ/quarantine policy;
- P04.07 schema-registry runtime;
- P04.08 background-job execution changes;
- P04.05+ implementation;
- business-domain handlers/features;
- strategic X-program runtime;
- AI/model/agent runtime.

## Governance and integration requirements

This transaction must follow the same accepted source/promotion pattern as prior high-value P04 transitions:

1. exact final source head passes canonical GitHub-hosted Governance;
2. changed paths stay within the recorded transition budget;
3. exact-head review is recorded honestly;
4. unresolved review threads are zero;
5. source is not itself merge authority;
6. the reviewed source head is promoted unchanged through a fresh promotion-specific carrier;
7. promotion receives fresh Governance while current with protected main;
8. merge uses expected-head protection;
9. protected main is re-read after merge before any new runtime branch is created.

Historical failed/cancelled/stale runs remain diagnostic evidence and are never relabeled PASS.

## Changed-path boundary

This closure/activation transaction is limited to:

- `AGENTS.md`;
- `README.md`;
- `docs/roadmap/STATE.json`;
- `docs/roadmap/STATUS.md`;
- `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`;
- `docs/ai/AI_STATE.yaml`;
- `docs/ai/AI_CONTEXT.md`;
- `docs/governance/P04_03_P04_04_TRANSITION_TRANSACTION.md`;
- `scripts/validate_p04_activation.py`.

No `kernel/`, `modules/`, migration, broker/provider, job-runtime or business source path belongs in this carrier.

## Exact next action

Finish this source transaction and obtain exact-head Governance. Then review the exact diff and promote the unchanged state through a fresh promotion-specific carrier. After promotion Governance, expected-head merge and protected-main read-back, confirm P04.04 is sole active at 3 / 10 and **STOP**. The next authorized action is creation of a separate P04.04 implementation branch from that exact protected-main SHA; runtime implementation must not occur on this activation carrier.
