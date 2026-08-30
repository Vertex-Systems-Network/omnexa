# Omnexa Program Status

Last reconciled: **2026-08-30**

## Current position

- Program: **Kernel Program**
- Phase: **P04 — Data, Jobs & Event Fabric**
- Phase state: **ACTIVE**
- Current work package: **P04.01 — Event Envelope & Identity Contract**
- P04 progress: **0 / 10 done**
- P04.01: **ACTIVE**
- P04.02-P04.10: **PLANNED / LOCKED**
- P04 entry gate: **SATISFIED CANDIDATE — PR #144 exact-head Governance + merge/read-back still required**
- P03: **DONE — 11 / 11**; exit **SATISFIED**
- P02: **DONE — 10 / 10**; exit **SATISFIED**
- P01: **DONE — 12 / 12**; exit **SATISFIED**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation authority: **AUTHORIZED FOR P04.01 ONLY AFTER ACTIVATION PR MERGE + PROTECTED-MAIN READ-BACK**
- Business-feature implementation: **NOT AUTHORIZED**

`docs/roadmap/STATE.json` is the canonical machine-readable candidate cursor. This activation branch itself is governance/state/continuity only and must not contain P04.01 runtime implementation.

## P04 activation chain

Accepted preparation/evidence chain:

- P03 terminal exit: SATISFIED with P03.01-P03.11 complete;
- P04 readiness promotion merge: `384c82b6d5b8bd89b70eb94640ea3be7109bc537`;
- bounded P04 planning/contract promotion merge: `73561ad68e453c2741dec82d0b6fb01a7dcea1be`;
- historical-validator compatibility prerequisite merge: `f97ee9a509f14247771626dc89b685e22ecd9457`;
- activation carrier: GitHub PR #144 from exact base `f97ee9a509f14247771626dc89b685e22ecd9457`.

PR #144 proposes only the canonical transition to P04/P04.01. Its candidate state becomes effective only after exact-final-head Omnexa Governance PASS, review/conversation preflight, merge, and protected-main state read-back.

## P04.01 authority boundary

P04.01 is contract-first and provider-neutral. After activation merge, its later separate implementation branch may implement only:

- stable event ID/type/version/producer identity;
- explicit tenant context and fail-closed validation;
- occurrence time, correlation and causation identity;
- bounded/classification-safe payload contract;
- deterministic safe validation errors;
- duplicate/replay-aware semantics with no global ordering assumption;
- focused positive/negative tests and package verification.

Explicitly unauthorized:

- P04.02-P04.10 implementation;
- Kafka/NATS/RabbitMQ/Redis Streams or another broker selection;
- outbox/inbox schema or database migration;
- durable consumers/checkpoints;
- retry/DLQ/quarantine runtime;
- business-domain event handlers;
- background-job execution changes;
- business-feature code;
- AI/model/agent authority.

## Historical prerequisite retention

P01.01-P01.12, P02.01-P02.10 and P03.01-P03.11 remain completed historical prerequisites with immutable completion evidence and mandatory regression verifiers. P04 activation does not rewrite their acceptance history.

The retained P01/P02/P03/foundation validators are extended only enough to recognize valid active-P04 state while continuing to require completed predecessor phase state, satisfied exits, locked predecessor manifests and `business_feature_code_authorized=false`.

## Protected integration / CI

`main` remains protected and PR-only. Canonical governance evidence is GitHub-hosted `ubuntu-24.04` only. Earlier P03/P04-readiness passes are prerequisites, not merge permission for PR #144.

The exact final head of PR #144 must pass the complete Omnexa Governance lane, including all retained P01/P02/P03 regressions and P04 activation validation.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions and grants no implementation authority.

## Exact next work

1. Finish PR #144 state/continuity reconciliation with no runtime source changes.
2. Require exact-final-head Omnexa Governance PASS and clean review/thread/comment state.
3. Record SELF REVIEW honestly if independent review is unavailable.
4. Promote/merge the unchanged accepted activation head through the governed PR path.
5. Re-read protected `main`, `STATE.json`, `STATUS.md`, P04 sequence and P04 entry gate.
6. Only then create a **new separate P04.01 implementation branch** from the exact post-activation main SHA.
7. Do not auto-advance to P04.02.
