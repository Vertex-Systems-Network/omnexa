# Omnexa Roadmap Status

Last reconciled: 2026-09-04 — **P04.04 closure / P04.05 activation candidate; protected main remains authoritative until governed promotion, merge and read-back**

## Authoritative pre-transaction checkpoint

Protected main: `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`.

- Foundation Architecture v1 is FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED.
- P04 is canonically ACTIVE — 3 / 10 done.
- P04.04 is still the canonical active package, with accepted implementation/completion evidence.
- P04.05 contract/handoff preparation is accepted but locked.
- P04.06-P04.10 and business/AI runtime remain locked.
- canonical CI is GitHub-hosted `ubuntu-24.04` only.

## Accepted P04.04 implementation/completion

P04.04 is accepted through its bounded multi-agent chain:

- T01 core: source #184 / promotion #186;
- T02 PostgreSQL persistence/migration: source #185 / promotion #187;
- T03 reliability/concurrency: source #188 / promotion #189;
- T02 regression-fixture repair: source #191 / promotion #192;
- T04 Supervisor verifier: source #190 / promotion #193;
- final exact implementation head `ef09b878577d25a4a1186cb8fe84205b08a24851`;
- promotion Governance `33810095507 / 100829646792` — PASS;
- implementation merge/read-back `66c072b5caf42ceecb88d30cd1a1ee4e910322e6`;
- completion evidence `docs/roadmap/evidence/P04.04_COMPLETION_2026-09-04.md`;
- evidence carrier #194 merge/read-back `4445c21f1e6b03e84859d31ce7b32169b9c4cccc`.

The accepted implementation preserves provider-neutral, at-least-once-compatible producer semantics. It does not select a broker, grant business authority, or claim end-to-end exactly once.

## Accepted P04.05 preparation

- preparation source #195;
- unchanged preparation promotion #196;
- exact preparation head `211fea2077d7a1bf94be48f32f047b27273a4515`;
- merge/read-back `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`;
- contract `docs/roadmap/work-packages/P04.05.md`;
- handoff `docs/ai/handoffs/P04.05.md`.

Preparation adds no runtime or migration and grants no authority before the separate activation transaction completes.

## Candidate atomic result

Only after exact-head source Governance, unchanged promotion Governance, protected merge and protected-main read-back:

- P04 ACTIVE — 4 / 10 done;
- P04.01-P04.04 DONE with retained evidence;
- P04.05 sole ACTIVE package;
- P04.06-P04.10 planned/locked;
- `kernel_code_authorized=true` only for P04.05;
- `business_feature_code_authorized=false`;
- strategic X and AI/model/agent runtime unauthorized.

## P04.05 runtime boundary after activation

Owner: `kernel.events`.

Authorized later runtime scope is limited to:

- stable consumer-processing identity bound to canonical EventID, consumer/owner/tenant/route scope;
- same-transaction local PostgreSQL protected mutation and inbox completion;
- deterministic already-applied handling after commit;
- restart and checkpoint-gap duplicate safety;
- concurrency, conflict and tenant/owner/consumer isolation;
- payload-safe failure behavior;
- focused verifier and retained regressions.

Still unauthorized:

- P04.06 retry/backoff/terminal failure/DLQ/quarantine;
- P04.07 schema registry;
- P04.08 background-job runtime;
- broad P04.09 operator recovery;
- concrete broker/provider selection;
- external-side-effect or end-to-end exactly-once claims;
- business features, strategic X runtime, or AI/model/agent runtime.

## Exact next work

1. Pass exact-head source Governance for this state/continuity-only transaction.
2. Review the exact diff and promote the unchanged source head.
3. Pass promotion Governance, merge with expected-head protection and re-read protected main.
4. Confirm P04.05 is sole active at 4 / 10 and P04.06+ remains locked.
5. Reconcile post-activation continuity in a separate carrier.
6. Only then create a fresh P04.05 wave/branches/leases and record the exact runtime/migration budget.
7. Do not auto-advance to P04.06.
