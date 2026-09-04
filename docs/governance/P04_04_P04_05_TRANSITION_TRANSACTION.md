# P04.04 Closure / P04.05 Activation Transaction

Status: **CANDIDATE — NOT AUTHORITATIVE UNTIL EXACT-HEAD SOURCE GOVERNANCE, UNCHANGED PROMOTION GOVERNANCE, PROTECTED MERGE AND READ-BACK**

Transaction base: protected `main@fa53b01cd92c8e0dd59026abff06f5f95f642d2d`.

This transaction changes governance/state/continuity only. It contains no P04.05 runtime source, migration, retry/DLQ policy, schema registry, background-job change, provider selection, business feature, strategic X runtime or AI/model/agent runtime.

## Preconditions

At the exact base:

- P04 is active at 3 / 10;
- P04.01-P04.03 are done;
- P04.04 remains canonical active but has accepted implementation/completion evidence;
- P04.05 contract/handoff preparation is accepted but locked;
- P04.06-P04.10 remain locked;
- `business_feature_code_authorized=false`.

## Accepted predecessor — P04.04

- final Supervisor integration/promotion head `ef09b878577d25a4a1186cb8fe84205b08a24851`;
- promotion PR #193;
- promotion Governance run/job `33810095507 / 100829646792` — PASS;
- implementation merge/read-back `66c072b5caf42ceecb88d30cd1a1ee4e910322e6`;
- completion evidence `docs/roadmap/evidence/P04.04_COMPLETION_2026-09-04.md`;
- evidence PR #194 merge/read-back `4445c21f1e6b03e84859d31ce7b32169b9c4cccc`.

Accepted P04.04 behavior remains provider-neutral and explicitly permits duplicate publication around publish-success/mark-failure. Producer published state does not prove downstream exactly-once mutation.

## Accepted preparation — P04.05

- source PR #195;
- unchanged promotion PR #196;
- exact preparation head `211fea2077d7a1bf94be48f32f047b27273a4515`;
- merge/read-back `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`;
- accepted contract `docs/roadmap/work-packages/P04.05.md`;
- handoff `docs/ai/handoffs/P04.05.md`.

Preparation alone grants no runtime or migration authority.

## Candidate result

After this exact source state is governed, promoted unchanged, governed again, merged and read back:

- P04 is active at 4 / 10;
- P04.01-P04.04 are done with accepted evidence;
- P04.05 is the sole active package;
- P04.06-P04.10 remain planned/locked;
- kernel code is authorized only within the accepted P04.05 contract;
- business-feature, strategic X and AI/model/agent runtime remain unauthorized.

## P04.05 bounded authority after acceptance

P04.05 may later implement only a provider-neutral local application-idempotency primitive:

- canonical EventID plus stable consumer/owner/tenant/route processing identity;
- no global cross-consumer EventID lock;
- inbox completion and protected local mutation in the same local PostgreSQL transaction;
- no completion before mutation commit;
- explicit already-applied outcome on committed duplicate redelivery;
- restart/checkpoint-gap and concurrent-delivery safety;
- fail-closed conflict and scope rebinding behavior;
- separation of checkpoint progress from inbox completion;
- no external-side-effect or end-to-end exactly-once claim.

Activation itself creates no runtime branch, task lease or migration. Those require a later fresh-main post-activation continuity/read-back and a separately governed P04.05 wave.

## Changed-path boundary

This source carrier is limited to:

- `AGENTS.md`;
- `README.md`;
- `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`;
- `docs/ai/AI_CONTEXT.md`;
- `docs/ai/AI_STATE.yaml`;
- `docs/governance/P04_04_P04_05_TRANSITION_TRANSACTION.md`;
- `docs/roadmap/STATE.json`;
- `docs/roadmap/STATUS.md`;
- `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`;
- `scripts/validate_p04_activation.py`.

Any kernel runtime, migration, workflow dependency/toolchain or future-package behavior is a stop/re-plan condition.

## Acceptance protocol

1. source exact head passes canonical GitHub-hosted `governance` on `ubuntu-24.04`;
2. exact diff and review provenance are recorded honestly, with zero unresolved threads;
3. the exact unchanged source head is promoted through a fresh promotion branch/PR;
4. promotion receives fresh Governance while current with protected main;
5. merge uses expected-head protection;
6. protected main, canonical state, package sequence and continuity are read back;
7. stop before P04.05 runtime work and perform separate post-activation continuity/wave planning.

Historical failed/cancelled/stale runs remain diagnostic evidence and are never relabeled PASS.
