# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

P04 activation is tracked by GitHub PR #144 / Linear ABD-267 from exact protected-main base `f97ee9a509f14247771626dc89b685e22ecd9457`.

Candidate state after accepted activation merge:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit SATISFIED.
- P02: DONE — 10 / 10; exit SATISFIED.
- P03: DONE — 11 / 11; exit SATISFIED / historical prerequisite.
- P04: ACTIVE — 0 / 10 done.
- P04.01: sole ACTIVE package.
- P04.02-P04.10: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P04.01 after activation merge/read-back.
- `business_feature_code_authorized=false`.

The activation branch is governance/state/continuity only. It must not contain P04.01 runtime source implementation.

## Accepted activation prerequisite chain

- P03 terminal exit remains SATISFIED.
- P04 readiness promotion merged as `384c82b6d5b8bd89b70eb94640ea3be7109bc537`.
- Bounded P04 planning / P04.01 contract promotion merged as `73561ad68e453c2741dec82d0b6fb01a7dcea1be`.
- Historical P03 validator compatibility prerequisite merged as `f97ee9a509f14247771626dc89b685e22ecd9457` after source Governance #528 and promotion Governance #529 PASS.
- PR #144 is the separate canonical activation transaction.

## P04.01 boundary

Owner: `kernel.events`.

P04.01 may later implement only the provider-neutral event envelope and validation contract:

- globally stable event identity/type/version/producer;
- explicit tenant context;
- occurred time, correlation and causation identity;
- bounded/classification-safe payload contract;
- deterministic safe validation errors;
- duplicate/replay-aware semantics;
- no global ordering assumption;
- focused positive/negative tests and verifier evidence.

P04.01 must not select Kafka/NATS/RabbitMQ/Redis Streams or another broker, add outbox/inbox persistence, add a database migration, implement durable consumers/checkpoints, retry/DLQ runtime, business-domain handlers, background-job execution changes, P04.02+, business features or AI/model/agent authority.

## Historical prerequisites retained

P01.01-P01.12, P02.01-P02.10 and P03.01-P03.11 retain immutable completion evidence and mandatory regression verifiers. Historical diagnostic failures remain diagnostic and are never rewritten as PASS.

Retained foundation/P01/P02/P03 validators are extended only to recognize a valid active-P04 checkpoint while still requiring completed predecessor state, satisfied exits, locked predecessor manifests, protected-main evidence, GitHub-hosted canonical CI and `business_feature_code_authorized=false`.

## Review and merge rule

PR #144 is not accepted merely because canonical files say `active`. Candidate activation becomes authoritative only after:

1. exact-final-head Omnexa Governance PASS;
2. changed-file review proves governance/state/continuity-only scope;
3. review threads/conversations are clean;
4. independent review is recorded when available, otherwise explicit `SELF REVIEW`;
5. governed promotion/merge completes;
6. protected `main` plus STATE/STATUS/P04 entry/sequence are re-read.

Only then may a new P04.01 implementation branch be created from exact post-activation `main`.

## Exact next action

Finish PR #144 continuity reconciliation, run the full exact-head governance lane, review the candidate, promote/merge without changing accepted bytes, re-read protected main, then start P04.01 implementation on a fresh branch. Do not auto-advance to P04.02.
