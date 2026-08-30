# Omnexa P04 Activation Transaction Template

Status: **PREPARED TEMPLATE — NOT EXECUTED / NOT AUTHORITY**

Companions:

- `docs/governance/P03_P04_TRANSITION_READINESS.md`
- `docs/governance/P04_ACTIVATION_GAP_LEDGER.md`

This file defines the deterministic transaction a later, separate governed carrier must execute. It does not itself mutate canonical state or authorize implementation.

## Preconditions — all mandatory

At activation time the carrier must re-read fresh protected `main` and prove:

1. `docs/roadmap/STATE.json.current_phase == "P03"` unless a separately accepted successor state explicitly supersedes this transition;
2. `current_work_package == null`;
3. the P03 exit gate is still satisfied;
4. no new blocker invalidates the P03 retained evidence;
5. the protected-main head equals the exact activation base recorded by the carrier;
6. the accepted P04 sequence and first-package contract exist on that exact base;
7. required governance/CI is green on the activation carrier itself.

Any stale-base or contradictory-state result aborts the transaction. Do not auto-rebase an activation decision across a material governance change.

## Canonical package sequence proposed for acceptance

The readiness sequence to be explicitly accepted is:

1. `P04.01` — event envelope, identity, version, tenant, correlation and causation contract;
2. `P04.02` — publish/subscribe abstraction and ownership boundaries;
3. `P04.03` — durable stream/consumer baseline and checkpoint model;
4. `P04.04` — transactional outbox reliability primitive;
5. `P04.05` — consumer inbox/deduplication and idempotency primitive;
6. `P04.06` — retry/backoff, terminal failure and dead-letter/quarantine policy;
7. `P04.07` — event schema registry, compatibility and validation;
8. `P04.08` — background-job ownership, tenant context and event/job correlation;
9. `P04.09` — reliability observability, diagnostics and operator recovery;
10. `P04.10` — aggregate replay/duplicate/failure/restart/poison-event acceptance.

A later activation review may refine the sequence, but it must do so explicitly and atomically rather than silently diverging from this readiness record.

## Minimum first-package authority

The activation transaction should authorize **P04.01 only**.

P04.01 implementation scope must remain contract-first. Before source paths are opened, the activation carrier must record exact owned paths and a path budget for:

- event envelope/domain contract definitions;
- focused contract/schema tests;
- public exports only where required by the accepted architecture;
- package-specific documentation/checkpoint evidence.

P04.01 must not silently include:

- broker/vendor integration;
- outbox/inbox persistence;
- database migration;
- durable consumer runtime;
- retry/DLQ implementation;
- business-domain event handlers;
- background-job execution changes.

Those belong to later packages unless the activation review explicitly changes the decomposition with evidence.

## Required semantic decisions recorded by activation

The activation carrier must either accept these values or explicitly replace them:

- delivery model: duplicates are possible and consumers must be idempotent;
- event identity: globally stable event ID;
- ordering: no global ordering assumption unless explicitly proven by a later package;
- tenant context: explicit, validated and fail-closed;
- correlation/causation: first-class envelope fields;
- schema evolution: versioned and compatibility-checked;
- secrets: prohibited from event payloads and diagnostic surfaces;
- authority: event transport never creates cross-module write authority;
- retry/replay: may re-deliver, therefore protected mutations must remain duplicate-safe.

No provider claim of "exactly once" may replace application-level idempotency evidence.

## State mutation blueprint

The later activation carrier must update canonical state in one governed transaction. The exact repository schema must be re-read at execution time; do not paste this pseudo-patch blindly.

Conceptually the transaction must record:

```text
current_phase = P04
current_work_package = P04.01
P04 activated = true
P04.01 state = ACTIVE
kernel implementation authority = minimum required for P04.01
business feature code authority = false
P04.02..P04.10 = LOCKED
activation_base = <exact protected-main SHA>
activation_pr = <activation PR>
readiness_evidence = <accepted readiness PR/head>
```

If canonical `STATE.json` uses different field names or a different authority model at execution time, the carrier must conform to that fresh schema rather than introduce parallel state keys.

## Required activation evidence

Before merge, the activation PR must record:

- exact base SHA and exact head SHA;
- changed-file list proving governance/state-only scope unless the canonical activation convention requires another owned artifact;
- clean review threads/comments;
- required governance/CI success on the exact head;
- proof that P03 evidence was not rewritten;
- explicit non-authorization of P04.02+;
- rollback procedure for an erroneous state-only activation before P04.01 implementation begins.

## Post-activation rule

Only after the activation transaction merges may a **new** P04.01 implementation carrier be created from fresh protected `main`.

Do not reuse the readiness branch or activation branch as the P04.01 implementation branch. This keeps readiness, authorization and implementation evidence independently reviewable.
