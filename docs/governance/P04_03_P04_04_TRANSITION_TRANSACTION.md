# P04.03 Closure / P04.04 Activation Transaction

Status: **CANDIDATE — NOT AUTHORITY UNTIL GOVERNED PROMOTION, MERGE & PROTECTED-MAIN READ-BACK**

Transaction base: `962a62c7c111079ca6f2047fa748deea97c84534`  
Source branch: `agent/20260831-omnexa-p04-03-close-p04-04-activate`

This transaction is governance/state/continuity plus narrowly required activation-validator hardening only. It contains no P04.04 runtime source, migration, broker/provider selection, business handler, background-job change or AI/model/agent runtime.

## Preconditions proved from protected main

At the transaction base:

- P04 is ACTIVE at 2 / 10 canonically done;
- P04.01 and P04.02 are DONE with retained accepted evidence;
- P04.03 is the sole canonical ACTIVE package;
- P04.03 implementation is accepted and completion evidence is present;
- P04.04 contract/handoff preparation has landed but remains prepared/planned/locked;
- P04.05-P04.10 remain planned/locked;
- `kernel_code_authorized=true` for P04.03 only on protected main;
- `business_feature_code_authorized=false`;
- canonical CI remains GitHub-hosted `ubuntu-24.04` only;
- protected main remains PR-only with required `governance` and conversation resolution.

Any stale-base or materially contradictory protected-main result invalidates this transaction and requires re-plan rather than silent rebase.

## Accepted predecessor evidence — P04.03

P04.03 implementation/completion evidence is already accepted and is retained, not manufactured by this carrier:

- source implementation PR: `#165`;
- promotion implementation PR: `#166`;
- final exact implementation head: `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`;
- source Governance: `#581` / run `33377793927` / job `99443112098`;
- promotion Governance: `#582` / run `33405463251` / job `99531835998`;
- accepted implementation merge/read-back: `b94189873bef11f4870935205398f1ef44f160bf`;
- completion evidence: `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`;
- completion-evidence carrier PR: `#167`;
- completion-evidence merge/read-back: `ed9c9b067c2725e9ddef4c3a2b03c4aa0b29dbcd`.

Retained P04.03 invariants include explicit owner/consumer/route/tenant/scope-bound progress, contiguous monotonic checkpoint advancement, deterministic rejection of stale/regressive/gapped/conflicting advancement, no progress after failed/cancelled handling, restart from the last accepted checkpoint, duplicate replay around crash/write-failure boundaries, race-safe progress update, no global ordering guarantee and no end-to-end exactly-once business-mutation guarantee.

P04.03 selected no broker/provider, no production checkpoint persistence and no database migration.

## Accepted preparation evidence — P04.04

P04.04 preparation also predates this activation transaction:

- preparation source PR: `#168`;
- exact preparation head: `658fbb4074f76c7960680c7e04bfd7cbc07a59a9`;
- source Governance: `#587` / run `33409859631` / job `99546416425`;
- unchanged preparation promotion PR: `#169`;
- promotion Governance: `#588` / run `33410873382` / job `99549818529`;
- preparation merge/read-back: `962a62c7c111079ca6f2047fa748deea97c84534`;
- accepted candidate contract path: `docs/roadmap/work-packages/P04.04.md`;
- prepared handoff: `docs/ai/handoffs/P04.04.md`.

Preparation alone grants no runtime authority. This transaction explicitly accepts the prepared P04.04 contract only if the complete candidate passes the governed source/promotion path and protected-main read-back.

## Atomic candidate result

Only after this exact transaction becomes authoritative, canonical state will establish:

- `current_phase=P04`;
- P04 progress `3 / 10 done`;
- P04.01 `done` with retained evidence;
- P04.02 `done` with retained evidence;
- P04.03 `done` with accepted evidence;
- P04.04 the **sole active** package;
- P04.05-P04.10 `planned / locked`;
- `kernel_code_authorized=true` for P04.04 only;
- `business_feature_code_authorized=false`;
- strategic X-program runtime unauthorized;
- AI/model/agent runtime unauthorized.

No P04.04 runtime source is included in this transaction. Implementation must start later on a new branch from exact post-transition protected main.

## Validator hardening required for safe advancement

The generic sequential validator already requires each completed predecessor to retain accepted spec/evidence plus matching PASS tracking, exact head/merge/run/job and GitHub-hosted `ubuntu-24.04` evidence.

This transaction extends that same fail-closed contract to P04.04 by:

1. requiring `docs/roadmap/work-packages/P04.04.md` as a known activation artifact;
2. loading the P04.04 contract during validation;
3. requiring core P04.04 markers for:
   - the same local PostgreSQL transaction boundary;
   - explicit duplicate-publication possibility;
   - no global event ordering guarantee;
   - no downstream exactly-once protected-mutation claim;
   - the preparation-carrier no-migration boundary;
   - P04.05 remaining later scope;
4. adding `governance_tracking.p04_03_completion=PASS` with exact accepted P04.03 implementation evidence;
5. preserving strict predecessor count, phase mirror and prepared-spec count rules.

This is fail-closed governance hardening, not a gate weakening and not P04.04 runtime implementation.

## P04.04 authority boundary after accepted transition

A fresh post-transition implementation branch may implement only the accepted P04.04 transactional-outbox boundary:

- stable owner-bound outbox identity linked to the canonical P04.01 event identity;
- authoritative owner mutation and canonical event committed/rolled back together in the same local PostgreSQL transaction;
- reuse of P01 `database.InTransaction`;
- committed pending outbox state remains recoverable after restart;
- relay consumes committed pending state only;
- relay publishes through the accepted P04.02 publication abstraction;
- publication failure leaves the record pending/recoverable;
- crash after successful publication but before the published mark can cause duplicate relay of the same canonical event;
- published state records producer-side publication progress only;
- concurrent relay attempts cannot corrupt or cross-bind owner/event/tenant state;
- no global ordering guarantee is inferred from row order, sequence, timestamps or relay scheduling;
- focused atomicity/crash/restart/duplicate/concurrency/tenant-isolation/malformed-state evidence.

## Persistence / migration decision

The accepted P04.04 contract requires durable local PostgreSQL outbox state for the production reliability claim, so **production persistence is expected**.

This activation transaction itself:

- adds no migration;
- changes no database schema;
- grants no blanket migration authority;
- grants no cross-owner database authority.

Before the first schema mutation on the later separate P04.04 implementation branch, fresh protected main must be re-read and an exact owner/path/version/data budget must be recorded. That runtime preflight must prove:

- the canonical `kernel.events` migration owner and exact path under `kernel/migrations/<owner>/<version>_<name>.sql`;
- immutable migration identity and compatibility/change class;
- fresh-install and supported-upgrade behavior through the retained P01 migration runner;
- rollback/forward-recovery;
- owner/tenant isolation;
- bounded persisted fields/classification and secret-free diagnostics.

If the exact migration namespace/path cannot be proved on that later fresh branch, schema mutation remains blocked until reconciled through governance/change control.

## Explicitly unauthorized

This transaction and later P04.04 authority do **not** authorize:

- Kafka, NATS, RabbitMQ, Redis Streams or another concrete broker/provider;
- an end-to-end exactly-once delivery or business-mutation claim;
- consumer inbox/deduplication persistence (`P04.05`);
- retry/backoff/DLQ/quarantine runtime (`P04.06`);
- schema-registry runtime (`P04.07`);
- background-job ownership/execution changes (`P04.08`);
- business-domain handlers or cross-module private-table writes;
- P04.05+ implementation;
- business-feature authority;
- strategic X-program runtime;
- AI/model/agent runtime.

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

Before P04.04 implementation begins, an erroneous state-only activation can be reverted through a separately governed transaction that restores P04.03 as active without rewriting accepted P04.01-P04.03 implementation evidence or the prepared P04.04 contract history.

Once a later P04.04 implementation has landed, rollback must preserve any accepted outbox persistence compatibility, canonical event identity, tenant/isolation and investigation/evidence obligations and must not pretend historical evidence never existed.

## Source acceptance requirements

Before source promotion:

1. exact final source head is recorded;
2. changed-file list stays within the declared path budget;
3. canonical Omnexa Governance passes on that exact head;
4. repository Go quality and all retained P01/P02/P03/P04.01-P04.03 regressions pass;
5. substantive exact-head SELF REVIEW or genuine independent review is recorded honestly;
6. review threads/conversations are zero unresolved;
7. protected main is still the exact compatible base or the transaction is re-planned.

Source evidence is not merge authority.

## Promotion / merge requirements

The reviewed exact source head must be promoted unchanged through a **fresh promotion-specific carrier**. That promotion must obtain its own fresh exact-head Omnexa Governance, clean review threads, honest promotion review, current protected-main read-back and expected-head guarded merge.

Only protected-main read-back after the promotion merge makes P04.04 active.

## Post-merge rule

After accepted merge:

1. re-read protected main SHA;
2. re-read `STATE.json`, `STATUS.md`, P04 sequence, this receipt and AI continuity;
3. confirm P04 is ACTIVE 3 / 10, P04.01-P04.03 DONE, P04.04 sole ACTIVE, P04.05-P04.10 locked;
4. confirm `kernel_code_authorized=true` for P04.04 only and `business_feature_code_authorized=false`;
5. identify creation of a **new isolated P04.04 implementation branch** from that exact main SHA as the next authorized action;
6. **STOP** — do not implement P04.04 on this activation branch or promotion branch;
7. do not auto-activate P04.05 after any later P04.04 implementation merge.
