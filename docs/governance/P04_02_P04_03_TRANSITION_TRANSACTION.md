# P04.02 Closure / P04.03 Activation Transaction

Status: **CANDIDATE — NOT AUTHORITY UNTIL GOVERNED PROMOTION, MERGE & PROTECTED-MAIN READ-BACK**

Transaction base: `2d454a87e03f404f081b6a87f216d0cfa8c7608d`
Source branch: `agent/20260831-omnexa-p04-02-close-p04-03-activate`

This transaction is governance/state/continuity plus narrowly required activation-validator hardening only. It contains no P04.03 durable consumer/checkpoint runtime, migration, provider selection, business handler, background-job change or AI/model/agent runtime.

## Preconditions proved from protected main

At the transaction base:

- P04 is ACTIVE at 1 / 10 done;
- P04.01 is DONE with retained accepted evidence;
- P04.02 is the sole ACTIVE package;
- P04.02 implementation is accepted and completion evidence is present;
- P04.03 contract/handoff preparation has landed but remains prepared/planned/locked;
- P04.04-P04.10 remain planned/locked;
- `kernel_code_authorized=true` for P04.02 only;
- `business_feature_code_authorized=false`;
- canonical CI remains GitHub-hosted `ubuntu-24.04` only;
- protected main remains PR-only with required `governance` and conversation resolution.

Any stale-base or materially contradictory protected-main result invalidates this transaction and requires re-plan rather than silent rebase.

## Accepted predecessor evidence — P04.02

P04.02 implementation/completion evidence is already accepted and is retained, not manufactured by this carrier:

- source implementation PR: `#155`;
- promotion implementation PR: `#156`;
- final exact implementation head: `5ec1de746eebc8734f86ec3aa2f311daae0dc18a`;
- source Governance: `#563` / run `33350757489`;
- promotion Governance: `#564` / run `33351220746` / job `99364850230`;
- accepted implementation merge/read-back: `1b378a2f44c6e3cba87b936e39e05f1a18da94cc`;
- completion evidence: `docs/roadmap/evidence/P04.02_COMPLETION_2026-08-31.md`;
- completion-evidence carrier PR: `#157`;
- completion-evidence merge/read-back: `8334bed24e793e200540bf953f7249f11242badf`.

Retained P04.02 invariants include provider-neutral publish/subscribe, stable ownership/consumer identity, validated P04.01 envelopes, fail-closed tenant mismatch behavior, duplicate-delivery/no-global-order semantics and no broker/provider-specific dependency.

## Accepted preparation evidence — P04.03

P04.03 preparation also predates this activation transaction:

- preparation source PR: `#158`;
- exact preparation head: `a20b37cbe0401626cd29e51a08efc399ee2d72e2`;
- source Governance: `#568` / run `33364934121`;
- unchanged preparation promotion PR: `#159`;
- promotion Governance: `#569` / run `33368684906` / job `99414689966`;
- preparation merge/read-back: `2d454a87e03f404f081b6a87f216d0cfa8c7608d`;
- accepted candidate contract path: `docs/roadmap/work-packages/P04.03.md`;
- prepared handoff: `docs/ai/handoffs/P04.03.md`.

Preparation alone grants no runtime authority. This transaction explicitly accepts the prepared P04.03 contract only if the complete candidate passes the governed source/promotion path and protected-main read-back.

## Atomic candidate result

Only after this exact transaction becomes authoritative, canonical state will establish:

- `current_phase=P04`;
- P04 progress `2 / 10 done`;
- P04.01 `done` with retained evidence;
- P04.02 `done` with retained evidence;
- P04.03 the **sole active** package;
- P04.04-P04.10 `planned / locked`;
- `kernel_code_authorized=true` for P04.03 only;
- `business_feature_code_authorized=false`;
- strategic X-program runtime unauthorized;
- AI/model/agent runtime unauthorized.

No P04.03 runtime source is included in this transaction. Implementation must start later on a new branch from exact post-transition protected main.

## Validator hardening required for safe advancement

The pre-transaction `scripts/validate_p04_activation.py` had a bounded governance weakness: advancement beyond P04.01 explicitly checked `p04_01_completion`, while later predecessor tracking was not generically required.

This transaction hardens strict sequential activation so every completed predecessor must prove:

1. package state is `done` in the P04 sequence;
2. accepted package spec remains present;
3. completion-evidence path remains present and exists;
4. matching `governance_tracking.p04_NN_completion.state == PASS` exists;
5. tracking `completion_evidence` matches the package sequence evidence;
6. canonical `final_exact_head`, `implementation_merge`, `workflow_run` and `job` are retained;
7. evidence remains `github-hosted` on `ubuntu-24.04`;
8. `STATE.json` phase mirror and P04 prepared-spec count remain consistent with the strict predecessor count.

The validator also recognizes the prepared P04.03 specification and enforces its core no-exactly-once/no-global-ordering/non-authorizing checkpoint laws.

This is fail-closed governance hardening, not a gate weakening and not P04.03 runtime implementation.

## P04.03 authority boundary after accepted transition

A fresh post-transition implementation branch may implement only the accepted P04.03 durable-consumer/checkpoint boundary:

- stable owner-bound durable consumer identity;
- explicit stream/shard/partition-equivalent checkpoint scope;
- monotonic checkpoint advancement within one scope;
- deterministic stale/regressive/conflicting checkpoint rejection;
- restart/resume from the last accepted checkpoint;
- no advancement past failed/cancelled unacknowledged work;
- fail-closed tenant/owner/scope collision handling;
- at-least-once-compatible duplicate delivery around crash/restart;
- preservation of P04.01 metadata and P04.02 ownership laws;
- focused positive/adversarial evidence.

## Explicitly unauthorized

This transaction and later P04.03 authority do **not** authorize:

- Kafka, NATS, RabbitMQ, Redis Streams or another concrete broker/provider;
- an end-to-end exactly-once delivery or business-mutation claim;
- transactional outbox (`P04.04`);
- consumer inbox/deduplication persistence (`P04.05`);
- retry/backoff/DLQ/quarantine runtime (`P04.06`);
- schema-registry runtime (`P04.07`);
- background-job execution changes (`P04.08`);
- business-domain handlers or cross-module private-table writes;
- P04.04+ implementation;
- business-feature authority;
- strategic X-program runtime;
- AI/model/agent runtime.

If a later P04.03 implementation requires persistence schema, it must separately prove that exact persistence path/migration is inside the accepted P04.03 budget; this activation transaction itself adds no migration and grants no automatic migration authority.

## Changed-path budget

Expected source-transaction paths are bounded to:

- `AGENTS.md`;
- `docs/roadmap/STATE.json`;
- `docs/roadmap/STATUS.md`;
- `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`;
- `docs/ai/AI_STATE.yaml`;
- `docs/ai/AI_CONTEXT.md`;
- this transaction receipt;
- `scripts/validate_p04_activation.py`.

No `kernel/`, `modules/`, migration, provider, job-runtime or business source path belongs in this carrier.

## Rollback boundary

Before P04.03 implementation begins, an erroneous state-only activation can be reverted through a separately governed transaction that restores P04.02 as active without rewriting accepted P04.01/P04.02 implementation evidence or the prepared P04.03 contract history.

Once a later P04.03 implementation has landed, rollback must preserve any accepted durable-state compatibility/tenant/isolation obligations and must not pretend historical evidence never existed.

## Source acceptance requirements

Before source promotion:

1. exact final source head is recorded;
2. changed-file list stays within the declared path budget;
3. canonical Omnexa Governance passes on that exact head;
4. repository Go quality and all retained P01/P02/P03/P04 regressions pass;
5. substantive exact-head SELF REVIEW or genuine independent review is recorded honestly;
6. review threads/conversations are zero unresolved;
7. protected main is still the exact compatible base or the transaction is re-planned.

Source evidence is not merge authority.

## Promotion / merge requirements

The reviewed exact source head must be promoted unchanged through a **fresh promotion-specific carrier**. That promotion must obtain its own fresh exact-head Omnexa Governance, clean review threads, honest promotion review, current protected-main read-back and expected-head guarded merge.

Only protected-main read-back after the promotion merge makes P04.03 active.

## Post-merge rule

After accepted merge:

1. re-read protected main SHA;
2. re-read `STATE.json`, `STATUS.md`, P04 sequence, this receipt and AI continuity;
3. confirm P04 is ACTIVE 2 / 10, P04.01-P04.02 DONE, P04.03 sole ACTIVE, P04.04-P04.10 locked;
4. confirm `kernel_code_authorized=true` for P04.03 only and `business_feature_code_authorized=false`;
5. create a **new isolated P04.03 implementation branch** from that exact main SHA;
6. do not reuse this activation branch or promotion branch for runtime implementation;
7. do not auto-activate P04.04 after any later P04.03 implementation merge.
