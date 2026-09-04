# Omnexa AI Project Context

Status: **P04.05 ACTIVE — post-activation continuity; subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative checkpoint

P04.05 activation merged/read back as `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a`.

Fresh protected main at this continuity carrier start is `35b438a77838518400758e8b877170918cc3278f`. The intervening #199 merge only updates two GitHub Actions checkout references and does not change P04.05 state or runtime authority.

At this checkpoint:

- P00-P03 are complete and their exit gates remain satisfied;
- P04 is ACTIVE at 4 / 10;
- P04.01-P04.04 are DONE with retained evidence;
- P04.05 is the sole ACTIVE package;
- P04.06-P04.10 remain planned/locked;
- `kernel_code_authorized=true` only for P04.05;
- `business_feature_code_authorized=false`;
- strategic X-program and AI/model/agent runtime remain unauthorized.

## Accepted P04.04 evidence

- final Supervisor integration/promotion head `ef09b878577d25a4a1186cb8fe84205b08a24851`;
- promotion Governance `33810095507 / 100829646792` — PASS;
- implementation merge/read-back `66c072b5caf42ceecb88d30cd1a1ee4e910322e6`;
- completion evidence `docs/roadmap/evidence/P04.04_COMPLETION_2026-09-04.md`;
- evidence carrier #194 merge/read-back `4445c21f1e6b03e84859d31ce7b32169b9c4cccc`.

The accepted implementation is provider-neutral, preserves duplicate-publication truth and grants no P04.05+ or business/AI authority beyond the later explicit activation transaction.

## Accepted P04.05 preparation

- source preparation PR #195;
- unchanged preparation promotion PR #196;
- exact preparation head `211fea2077d7a1bf94be48f32f047b27273a4515`;
- merge/read-back `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`;
- contract `docs/roadmap/work-packages/P04.05.md`;
- handoff `docs/ai/handoffs/P04.05.md`.

## Accepted P04.05 activation chain

P04.04 closure / P04.05 activation is accepted evidence:

- source PR #197;
- exact source/promotion head `6907253d375125a7ff096fb434c3433dbc17b331`;
- source Governance `33819597433 / 100859169990` — PASS;
- source SELF REVIEW only; independent approval not claimed; zero unresolved threads;
- unchanged promotion PR #198;
- promotion Governance `33820312475 / 100861375770` — PASS;
- promotion SELF REVIEW only; independent approval not claimed; zero unresolved threads;
- guarded merge/read-back `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a`.

No P04.05 runtime source or migration was part of activation.

## P04.05 implementation law

Owner: `kernel.events`.

After this separate post-activation continuity correction is governed, promoted unchanged, merged and read back, a fresh implementation wave may implement only the provider-neutral consumer inbox/deduplication and local application-idempotency primitive:

- bind canonical EventID to stable consumer/owner/tenant/route processing scope;
- do not create a global cross-consumer EventID lock;
- commit protected local mutation and inbox completion in the same retained local PostgreSQL transaction;
- never persist completion before the protected mutation commits;
- committed redelivery returns explicit already-applied without rerunning the protected mutation;
- concurrent same-scope attempts cannot both commit;
- conflicting canonical content or tenant/owner/consumer rebinding fails closed;
- checkpoint progress and inbox completion remain separate facts;
- external/non-transactional side effects are not made exactly once;
- no retry/backoff/DLQ, schema-registry, job-runtime, provider selection, business feature or AI runtime is introduced.

## Persistence / migration gate

Durable local PostgreSQL inbox persistence is required for production semantics, but activation and continuity grant no blanket migration authority. Before the first schema mutation on the later fresh runtime branch:

1. re-read exact protected main and current `kernel.events` migration history;
2. choose the exact next immutable owner-scoped migration path/version;
3. record exact tables/columns/indexes/constraints/data budget;
4. preserve retained P01 fresh-install/upgrade ledger behavior and P03.09 ownership;
5. define forward-recovery/rollback that cannot reopen an already-applied duplicate mutation;
6. prove tenant/owner/consumer isolation;
7. exclude P04.06+ and business-domain schema semantics.

## Locked scope

Do not implement concrete broker/provider selection, retry/backoff/DLQ/quarantine, schema-registry runtime, background-job runtime changes, broad operator recovery, business handlers/features, strategic X runtime, AI/model/agent runtime, global ordering or external/end-to-end exactly-once claims.

## Exact next action

Govern and land this post-activation continuity correction first, using exact-head source Governance, exact diff/review/thread verification, unchanged promotion Governance, expected-head protected merge and protected-main read-back. Then create a **fresh separate P04.05 implementation wave** from that exact main SHA, record exact runtime/migration budgets and implement only `docs/roadmap/work-packages/P04.05.md`.

Do not reuse activation/promotion/continuity branches and do not auto-advance to P04.06.
