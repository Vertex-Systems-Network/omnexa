# Omnexa AI Project Context

Status: **P04.04 closure / P04.05 activation candidate; protected main remains authoritative until governed promotion, merge and read-back**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative pre-transaction checkpoint

Fresh protected main is `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`.

At that checkpoint:

- P00-P03 are complete and their exit gates remain satisfied;
- P04 is canonically ACTIVE at 3 / 10;
- P04.04 remains the sole canonical active package but its implementation and completion evidence are accepted;
- P04.05 contract/handoff preparation is accepted but locked;
- P04.06-P04.10, business features, strategic X runtime and AI/model/agent runtime remain locked.

## Accepted P04.04 evidence

- final Supervisor integration/promotion head: `ef09b878577d25a4a1186cb8fe84205b08a24851`;
- promotion PR: #193;
- promotion Governance run/job: `33810095507 / 100829646792` — PASS;
- implementation merge/read-back: `66c072b5caf42ceecb88d30cd1a1ee4e910322e6`;
- completion evidence: `docs/roadmap/evidence/P04.04_COMPLETION_2026-09-04.md`;
- evidence carrier #194 merge/read-back: `4445c21f1e6b03e84859d31ce7b32169b9c4cccc`.

The accepted P04.04 wave is provider-neutral, preserves duplicate-publication truth, and grants no P04.05+ or business/AI authority.

## Accepted P04.05 preparation

- source preparation PR #195;
- unchanged preparation promotion PR #196;
- exact preparation head `211fea2077d7a1bf94be48f32f047b27273a4515`;
- merge/read-back `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`;
- contract `docs/roadmap/work-packages/P04.05.md`;
- handoff `docs/ai/handoffs/P04.05.md`.

Preparation alone grants no runtime or migration authority.

## Candidate result

Only after this closure transaction passes exact-head source Governance, is promoted unchanged, passes promotion Governance, merges through protected main and is read back:

- P04 becomes 4 / 10 done;
- P04.01-P04.04 are done with retained evidence;
- P04.05 is the sole active package;
- P04.06-P04.10 remain planned/locked;
- `kernel_code_authorized=true` only for P04.05;
- `business_feature_code_authorized=false`.

## P04.05 bounded implementation law after activation

P04.05 may later implement only the provider-neutral consumer inbox/deduplication and local application-idempotency primitive under `kernel.events`:

- bind canonical EventID to stable consumer/owner/tenant/route processing scope;
- keep EventID from becoming a global cross-consumer lock;
- commit protected local mutation and inbox completion in the same local PostgreSQL transaction;
- never persist completion before the protected mutation commits;
- committed redelivery returns an explicit already-applied outcome without rerunning the mutation;
- concurrent same-scope attempts cannot both commit;
- conflicting identity/content reuse and tenant/owner/consumer rebinding fail closed;
- checkpoint progress and inbox completion remain separate facts;
- external side effects are not made exactly once;
- no retry/backoff/DLQ, schema registry, job runtime, provider selection, business feature or AI runtime is introduced.

Activation adds no migration and no runtime branch. A later post-activation read-back must precede any fresh wave, task/branch/lease creation or schema mutation.

## Exact next action

Govern and promote only the P04.04 closure / P04.05 activation transaction. After protected-main read-back, reconcile post-activation continuity and then create a separate P04.05 implementation wave. Do not auto-advance P04.06.
