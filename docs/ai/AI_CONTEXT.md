# Omnexa AI Project Context

Status: **P04.02 closure / P04.03 activation candidate — subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected `main@2d454a87e03f404f081b6a87f216d0cfa8c7608d` remains authoritative until this complete closure/activation transaction passes exact-final-head Omnexa Governance, receives substantive exact-head review, is promoted unchanged through a fresh promotion carrier, passes promotion-specific Governance, merges through the protected path, and protected main plus canonical state are re-read.

Authoritative protected-main state before this transaction:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit SATISFIED.
- P02: DONE — 10 / 10; exit SATISFIED.
- P03: DONE — 11 / 11; exit SATISFIED / historical prerequisite.
- P04: ACTIVE — 1 / 10 done.
- P04.01: DONE.
- P04.02: sole ACTIVE package with accepted implementation/completion evidence.
- P04.03: prepared/planned/locked.
- P04.04-P04.10: planned/locked.
- `kernel_code_authorized=true` for P04.02 only.
- `business_feature_code_authorized=false`.

Candidate atomic result after this transaction becomes authoritative:

- P04 remains ACTIVE — **2 / 10 done**.
- P04.01 remains DONE.
- P04.02 becomes DONE with retained accepted evidence.
- P04.03 becomes the sole ACTIVE package.
- P04.04-P04.10 remain PLANNED / LOCKED.
- `kernel_code_authorized=true` for P04.03 only.
- `business_feature_code_authorized=false`.
- strategic X-program and AI/model/agent runtime remain unauthorized.

This branch is governance/state/continuity plus activation-validator hardening only. It must not contain P04.03 runtime implementation.

## Accepted P04.02 chain

P04.02 implementation acceptance predates this transaction and remains immutable predecessor evidence:

- P04.01 closure / P04.02 activation merge: `226f5236d80c806a7b07d65e7870981de320c05c`.
- source implementation PR #155 exact final head `5ec1de746eebc8734f86ec3aa2f311daae0dc18a` passed Governance #563 / run `33350757489`.
- unchanged promotion PR #156 passed promotion-specific Governance #564 / run `33351220746` / job `99364850230`.
- accepted implementation merge/read-back: `1b378a2f44c6e3cba87b936e39e05f1a18da94cc`.
- completion evidence: `docs/roadmap/evidence/P04.02_COMPLETION_2026-08-31.md`.
- completion-evidence carrier PR #157 merged/read back as `8334bed24e793e200540bf953f7249f11242badf`.

P04.02 accepted behavior remains provider-neutral publish/subscribe plus ownership boundaries only. Its completion does not grant P04.03 authority by itself.

## Accepted P04.03 preparation chain

The P04.03 contract/handoff was separately prepared without activation authority:

- source preparation PR #158;
- exact preparation head `a20b37cbe0401626cd29e51a08efc399ee2d72e2`;
- source Governance #568 / run `33364934121` — PASS;
- unchanged promotion PR #159;
- promotion Governance #569 / run `33368684906` / job `99414689966` — PASS;
- preparation merge/read-back `2d454a87e03f404f081b6a87f216d0cfa8c7608d`;
- contract `docs/roadmap/work-packages/P04.03.md`;
- handoff `docs/ai/handoffs/P04.03.md`.

Preparation is not implementation authority. This separate closure/activation transaction is required before a P04.03 runtime branch may exist.

## P04.03 candidate implementation boundary

Owner: `kernel.events`.

Only after this transaction merges and protected main is re-read may a **new separate branch** implement the accepted P04.03 contract:

- stable durable consumer identity compatible with P04.02 owner/consumer identity;
- explicit durable stream/partition-equivalent checkpoint scope;
- monotonic checkpoint advancement within one scope;
- deterministic stale/regressive/conflicting checkpoint rejection;
- restart/resume from the last accepted checkpoint;
- failed/cancelled unacknowledged work cannot be skipped by checkpoint advancement;
- fail-closed tenant/owner/scope collisions;
- at-least-once-compatible duplicate delivery around crash/restart;
- preservation of P04.01 envelope metadata and P04.02 ownership identity;
- focused crash/restart, monotonicity, duplicate, cancellation, malformed-state and tenant/scope isolation tests.

P04.03 must not select Kafka/NATS/RabbitMQ/Redis Streams or another broker, claim exactly-once delivery/business mutation, implement transactional outbox, inbox/deduplication persistence, retry/DLQ, schema-registry runtime, background-job execution changes, business-domain handlers, P04.04+, business features, strategic X-program runtime or AI/model/agent authority.

## Sequential activation validator hardening

This transaction also corrects a governance weakness in `scripts/validate_p04_activation.py`. The old validator had a P04.01-specific predecessor tracking exception. The candidate validator now generically requires every completed P04 predecessor to retain:

- `done` package-sequence state;
- accepted spec and completion-evidence file;
- matching `governance_tracking.p04_NN_completion.state=PASS`;
- matching completion-evidence reference;
- exact implementation head and merge;
- canonical workflow run and job;
- GitHub-hosted evidence on `ubuntu-24.04`;
- consistent state/sequence mirrors and prepared-spec count.

Therefore P04.02 accepted completion tracking is now a fail-closed prerequisite to P04.03 activation, and later P04 package transitions do not need new predecessor-specific exceptions.

## Historical prerequisites retained

P01.01-P01.12, P02.01-P02.10, P03.01-P03.11 and P04.01 retain immutable completion evidence and mandatory regressions. Historical diagnostic failures remain diagnostic and are never rewritten as PASS.

## Review and merge rule for this transaction

Before this activation candidate may land:

1. `AGENTS.md`, canonical `STATE.json`, `STATUS.md`, P04 package sequence, `AI_STATE.yaml`, this context and the transition receipt must agree on the candidate cursor and locks;
2. exact-final-head Omnexa Governance must PASS;
3. changed-file review must prove governance/state/continuity/validator-only scope with no P04.03 runtime;
4. review threads/conversations must be clean;
5. independent review is recorded when available, otherwise explicit `SELF REVIEW` only;
6. the exact reviewed head is promoted unchanged through a fresh promotion-specific carrier;
7. the promotion carrier receives its own fresh Governance PASS;
8. expected-head protected merge succeeds while current with main;
9. protected `main`, canonical state, sequence and continuity are re-read.

Only after that read-back may a new P04.03 implementation branch be created.

## Exact next action

Finish this closure/activation candidate, run exact-final-head Governance, perform substantive self-review, promote the unchanged reviewed head, require promotion-specific Governance, guarded-merge, re-read protected main, and only then start P04.03 implementation on a fresh branch. Do not implement P04.03 on this transaction branch and do not auto-advance to P04.04.
