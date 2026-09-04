# P04.04 Closure / P04.05 Activation Transaction

Status: **ACCEPTED — MERGED AND READ BACK FROM PROTECTED MAIN**

Transaction base: `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`  
Exact source/promotion head: `6907253d375125a7ff096fb434c3433dbc17b331`  
Accepted merge/read-back: `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a`

This transaction changed governance/state/continuity only. It contained no P04.05 runtime source, migration, retry/DLQ policy, schema registry, background-job change, provider selection, business feature, strategic X runtime or AI/model/agent runtime.

## Accepted transition evidence

- source PR: `#197`;
- exact source head: `6907253d375125a7ff096fb434c3433dbc17b331`;
- source Governance run/job: `33819597433 / 100859169990` — PASS;
- source review: SELF REVIEW; independent approval not claimed;
- source unresolved review threads: 0;
- unchanged promotion PR: `#198`;
- promotion exact head: `6907253d375125a7ff096fb434c3433dbc17b331`;
- promotion Governance run/job: `33820312475 / 100861375770` — PASS;
- promotion review: SELF REVIEW; independent approval not claimed;
- promotion unresolved review threads: 0;
- expected-head guarded merge/read-back: `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a`.

A later separately governed maintenance PR #199 updated only two GitHub Actions checkout references and produced protected main `35b438a77838518400758e8b877170918cc3278f`; it did not widen or alter this accepted activation result.

## Preconditions retained

At the transaction base:

- P04 was ACTIVE at 3 / 10;
- P04.01-P04.03 were DONE;
- P04.04 was the sole active package with accepted implementation/completion evidence;
- P04.05 contract/handoff preparation was accepted but locked;
- P04.06-P04.10 were locked;
- business-feature, strategic X and AI/model/agent runtime were unauthorized.

## Accepted predecessor evidence — P04.04

- final implementation head `ef09b878577d25a4a1186cb8fe84205b08a24851`;
- promotion Governance `33810095507 / 100829646792`;
- implementation merge/read-back `66c072b5caf42ceecb88d30cd1a1ee4e910322e6`;
- completion evidence `docs/roadmap/evidence/P04.04_COMPLETION_2026-09-04.md`;
- completion evidence merge/read-back `4445c21f1e6b03e84859d31ce7b32169b9c4cccc`.

## Accepted preparation evidence — P04.05

- source preparation PR #195;
- unchanged promotion PR #196;
- exact preparation head `211fea2077d7a1bf94be48f32f047b27273a4515`;
- preparation merge/read-back `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`;
- accepted contract `docs/roadmap/work-packages/P04.05.md`;
- handoff `docs/ai/handoffs/P04.05.md`.

## Canonical protected-main result

Protected-main read-back `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a` established:

- P04 ACTIVE — `4 / 10 done`;
- P04.01-P04.04 DONE with retained evidence;
- P04.05 sole ACTIVE package;
- P04.06-P04.10 planned/locked;
- `kernel_code_authorized=true` only for P04.05;
- `business_feature_code_authorized=false`;
- strategic X runtime unauthorized;
- AI/model/agent runtime unauthorized.

Canonical `docs/roadmap/STATE.json` already records that state and is intentionally unchanged by this later post-activation continuity reconciliation.

## Accepted P04.05 implementation boundary

P04.05 may implement only the accepted provider-neutral local application-idempotency primitive under `kernel.events`:

- canonical EventID plus stable consumer/owner/tenant/route processing identity;
- no global cross-consumer EventID lock;
- protected local mutation and inbox completion in the same retained PostgreSQL transaction;
- no completion before mutation commit;
- explicit already-applied outcome on committed duplicate redelivery;
- restart/checkpoint-gap and concurrent-delivery safety;
- fail-closed conflicting identity/content reuse and scope rebinding;
- checkpoint progress separate from inbox completion;
- no external-side-effect or end-to-end exactly-once claim.

Activation granted no blanket migration authority. Durable production inbox persistence requires the later fresh runtime wave to record the exact `kernel.events` migration path/version/data budget before schema mutation.

## Explicitly excluded

- P04.06 retry/backoff/DLQ/quarantine;
- P04.07 schema registry runtime;
- P04.08 background-job runtime changes;
- broad P04.09 operator recovery;
- concrete broker/provider selection;
- business handlers/features;
- distributed transaction or global ordering guarantees;
- external/end-to-end exactly-once claims;
- strategic X runtime;
- AI/model/agent runtime.

## Post-activation continuity rule

Runtime implementation must not occur on activation, promotion or continuity carriers. The accepted activation transaction explicitly required a later fresh-main continuity/read-back before creating P04.05 runtime workers/branches.

This continuity correction is bounded to subordinate continuity surfaces only. It does not modify canonical `STATE.json`, P04 package sequence, validators, workflow policy or runtime source.

## Exact next action

1. Govern the separate P04.05 post-activation continuity source from fresh protected main.
2. Require exact-head source Governance, exact diff/review inspection and zero unresolved threads.
3. Promote the **unchanged continuity source state** through a fresh promotion carrier and require fresh promotion Governance.
4. Merge only with expected-head protection while current with protected main, then re-read canonical state and all continuity surfaces.
5. Confirm P04 remains `4 / 10`, P04.05 remains sole ACTIVE, P04.06-P04.10 remain locked and no authority expanded.
6. Only then create a **fresh P04.05 implementation wave** from that exact protected-main SHA with explicit non-overlapping path leases.
7. Record the exact `kernel.events` runtime/migration path/version/data budget before schema mutation.
8. Implement only P04.05 and do not auto-advance P04.06.
