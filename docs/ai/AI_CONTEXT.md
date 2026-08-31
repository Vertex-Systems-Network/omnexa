# Omnexa AI Project Context

Status: **post-P04.03-activation continuity snapshot — subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

The P04.02 closure / P04.03 activation transaction is accepted. Protected-main activation read-back is `d74375cd0a2952ab8622117089e5eb43043e6e78`.

Current authoritative state at that checkpoint:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit SATISFIED.
- P02: DONE — 10 / 10; exit SATISFIED.
- P03: DONE — 11 / 11; exit SATISFIED / historical prerequisite.
- P04: ACTIVE — **2 / 10 done**.
- P04.01: DONE with retained accepted evidence.
- P04.02: DONE with accepted implementation/completion evidence.
- P04.03: sole ACTIVE package.
- P04.04-P04.10: PLANNED / LOCKED.
- `kernel_code_authorized=true` for P04.03 only.
- `business_feature_code_authorized=false`.
- strategic X-program runtime remains unauthorized.
- AI/model/agent runtime remains unauthorized.

Any continuity-only reconciliation carrier after activation must contain no P04.03 runtime implementation and grants no additional authority.

## Accepted P04.02 chain

P04.02 implementation acceptance remains immutable predecessor evidence:

- P04.01 closure / P04.02 activation merge: `226f5236d80c806a7b07d65e7870981de320c05c`.
- source implementation PR #155 exact final head `5ec1de746eebc8734f86ec3aa2f311daae0dc18a` passed Governance #563 / run `33350757489`.
- unchanged promotion PR #156 passed promotion-specific Governance #564 / run `33351220746` / job `99364850230`.
- accepted implementation merge/read-back: `1b378a2f44c6e3cba87b936e39e05f1a18da94cc`.
- completion evidence: `docs/roadmap/evidence/P04.02_COMPLETION_2026-08-31.md`.
- completion-evidence carrier PR #157 merged/read back as `8334bed24e793e200540bf953f7249f11242badf`.

P04.02 retained laws include provider-neutral publish/subscribe, stable owner/module and consumer identity, validated P04.01 envelopes, fail-closed tenant mismatch behavior, duplicate-delivery/no-global-order semantics and no concrete broker/provider dependency.

## Accepted P04.03 preparation and activation chain

The P04.03 contract/handoff was separately prepared first:

- source preparation PR #158;
- exact preparation head `a20b37cbe0401626cd29e51a08efc399ee2d72e2`;
- source Governance #568 / run `33364934121` — PASS;
- unchanged preparation promotion PR #159;
- promotion Governance #569 / run `33368684906` / job `99414689966` — PASS;
- preparation merge/read-back `2d454a87e03f404f081b6a87f216d0cfa8c7608d`;
- contract `docs/roadmap/work-packages/P04.03.md`;
- handoff `docs/ai/handoffs/P04.03.md`.

The separate activation transaction then completed through the governed source/promotion path:

- source activation PR #160;
- exact reviewed source/promotion head `452858a3ab9bfa827697105bd5168cf660bd62ba`;
- source Governance #571 / run `33370681216` / job `99420856278` — PASS;
- unchanged promotion PR #161;
- promotion Governance #572 / run `33371203708` / job `99422521576` — PASS;
- zero unresolved review threads before merge;
- explicit source and promotion SELF REVIEW, with no fabricated independent review;
- expected-head guarded merge/read-back `d74375cd0a2952ab8622117089e5eb43043e6e78`.

P04.03 activation changes the package cursor only. It does not retroactively change P04.02 evidence and does not authorize P04.04+.

## P04.03 authoritative implementation boundary

Owner: `kernel.events`.

After verifying current protected main and canonical state, a **new separate implementation branch** may implement only the accepted P04.03 contract:

- stable durable consumer identity compatible with P04.02 owner/consumer identity;
- explicit durable stream/shard/partition-equivalent checkpoint scope;
- monotonic checkpoint advancement within one scope;
- deterministic stale/regressive/conflicting checkpoint rejection;
- restart/resume from the last accepted checkpoint;
- failed/cancelled unacknowledged work cannot be skipped by checkpoint advancement;
- fail-closed owner/tenant/scope collision handling;
- explicit at-least-once-compatible duplicate delivery around crash/restart;
- preservation of P04.01 envelope metadata and P04.02 ownership laws;
- focused crash/restart, monotonicity, duplicate, cancellation, malformed-state and tenant/scope isolation tests;
- dedicated package verifier and exact-head acceptance evidence.

P04.03 must not select Kafka/NATS/RabbitMQ/Redis Streams or another broker, claim exactly-once delivery/business mutation, assume activation grants blanket migration authority, implement transactional outbox, inbox/deduplication persistence, retry/DLQ, schema-registry runtime, background-job execution changes, business-domain handlers, P04.04+, business features, strategic X-program runtime or AI/model/agent authority.

If persistence is genuinely required by the accepted P04.03 implementation, the exact persistence/migration path must be justified inside the P04.03 budget with rollback, compatibility and tenant-isolation evidence. Activation itself is not blanket migration permission.

## Sequential activation validator hardening retained

`scripts/validate_p04_activation.py` generically requires every completed P04 predecessor to retain:

- `done` package-sequence state;
- accepted spec and completion-evidence file;
- matching `governance_tracking.p04_NN_completion.state=PASS`;
- matching completion-evidence reference;
- exact implementation head and merge;
- canonical workflow run and job;
- GitHub-hosted evidence on `ubuntu-24.04`;
- consistent state/sequence mirrors and prepared-spec count.

P04.02 accepted completion tracking is therefore a fail-closed prerequisite for active P04.03. Later transitions must preserve this generic rule rather than reintroduce predecessor-specific exceptions.

## Historical prerequisites retained

P01.01-P01.12, P02.01-P02.10, P03.01-P03.11, P04.01 and P04.02 retain immutable completion evidence and mandatory regressions. Historical diagnostic failures remain diagnostic and are never rewritten as PASS.

## Active-package execution rule

Before material P04.03 work:

1. re-read protected `main`, `AGENTS.md`, canonical `STATE.json`, `STATUS.md`, P04 sequence, P04.03 spec and active handoff;
2. confirm P04.03 is the sole active package and P04.04-P04.10 remain locked;
3. confirm `kernel_code_authorized=true` only for P04.03 and `business_feature_code_authorized=false`;
4. create a new isolated implementation branch from the exact verified protected-main SHA;
5. record a bounded owned-path/path-budget decision before runtime mutation;
6. preserve P04.01/P04.02 identity/ownership semantics rather than inventing a competing bus or authorization model;
7. require focused positive/adversarial evidence plus exact-final-head GitHub-hosted Governance;
8. use the repository's governed promotion path for accepted implementation;
9. do not auto-advance to P04.04 after merge.

## Exact next action

Ensure this post-activation continuity snapshot and the active P04.03 handoff are present on protected main through the governed path. Then create a **fresh separate P04.03 implementation branch** from exact verified protected main and implement only `docs/roadmap/work-packages/P04.03.md`. Do not reuse activation/promotion/continuity branches and do not auto-activate P04.04.
