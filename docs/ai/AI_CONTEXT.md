# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

P04.01 implementation is complete and accepted. The separate P04.01 closure / P04.02 activation transaction passed source Governance #554, promotion Governance #555, merged through PR #151, and was read back on protected `main` at `226f5236d80c806a7b07d65e7870981de320c05c` under Linear `ABD-267`.

Current authoritative state after that read-back:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit SATISFIED.
- P02: DONE — 10 / 10; exit SATISFIED.
- P03: DONE — 11 / 11; exit SATISFIED / historical prerequisite.
- P04: ACTIVE — 1 / 10 done.
- P04.01: DONE with accepted completion evidence.
- P04.02: sole ACTIVE package.
- P04.03-P04.10: PLANNED / LOCKED.
- `kernel_code_authorized=true` for P04.02 only.
- `business_feature_code_authorized=false`.

This post-activation reconciliation lane is continuity/governance only. It must not contain P04.02 runtime implementation.

## Accepted P04.01 / P04.02 chain

- P03 terminal exit remains SATISFIED.
- P04 readiness promotion merged as `384c82b6d5b8bd89b70eb94640ea3be7109bc537`.
- Bounded P04 planning / P04.01 contract promotion merged as `73561ad68e453c2741dec82d0b6fb01a7dcea1be`.
- P04/P04.01 activation merged as `dc269932481a1e742926ac247f28f94bdf43937d`.
- P04.01 implementation exact head `3c4ee6e79b76a6f042f3ee77b518943f3d7064e4` passed promotion-specific Governance #548 and merged through PR #147 as `6b1a01c009e9a08d45613a00e9b63c0c272bf020`.
- P04.01 completion evidence is `docs/roadmap/evidence/P04.01_COMPLETION_2026-08-31.md`.
- Prepared/locked P04.02 contract and handoff passed promotion-specific Governance #551 and merged through PR #149 as `16dfca22cd0430b8b2727e55628cc6fba381b5ef`.
- P04.01 closure / P04.02 activation source PR #150 exact head `8f00090a2055a2283d61b30a66b0772ddf354648` passed Governance #554 / run `33334373676`.
- Unchanged promotion PR #151 passed promotion-specific Governance #555 / run `33335001265`.
- Expected-head merge and protected-main read-back are `226f5236d80c806a7b07d65e7870981de320c05c`.

Historical acceptance is immutable; the post-activation continuity correction only removes stale candidate wording from subordinate snapshots.

## P04.02 authoritative boundary

Owner: `kernel.events`.

After this continuity-only correction lands and protected main is re-read, a **new separate branch** may implement only the provider-neutral publish/subscribe and ownership-boundary contract in `docs/roadmap/work-packages/P04.02.md`:

- publisher/subscriber abstraction without vendor lock-in;
- explicit owner/module boundaries;
- reuse of the completed P04.01 envelope and tenant/correlation semantics;
- deterministic bounded reference behavior where the accepted spec permits it;
- fail-closed ownership/authority validation;
- focused positive/adversarial tests and verifier evidence.

P04.02 must not select Kafka/NATS/RabbitMQ/Redis Streams or another broker, add durable streams/checkpoints, database migration, outbox/inbox persistence, retry/DLQ/quarantine runtime, schema-registry runtime, background-job execution changes, business-domain handlers, P04.03+, business features, strategic X-program runtime or AI/model/agent authority.

## Historical prerequisites retained

P01.01-P01.12, P02.01-P02.10 and P03.01-P03.11 retain immutable completion evidence and mandatory regression verifiers. Historical diagnostic failures remain diagnostic and are never rewritten as PASS.

## Review and merge rule for this reconciliation

This continuity correction is not P04.02 implementation. Before it may land:

1. `AGENTS.md`, `STATUS.md`, `AI_STATE.yaml`, this context and the transition receipt must agree with canonical `STATE.json`;
2. exact-final-head Omnexa Governance must PASS;
3. changed-file review must prove continuity/governance-only scope;
4. review threads/conversations must be clean;
5. independent review is recorded when available, otherwise explicit `SELF REVIEW`;
6. unchanged promotion/merge completes through the governed path;
7. protected `main` plus canonical state are re-read.

Only after that correction may a new P04.02 implementation branch be created from exact corrected `main`.

## Exact next action

Complete this post-activation continuity reconciliation without P04.02 runtime code, require exact-final-head Omnexa Governance, review/promote/merge, re-read protected main, and then start P04.02 implementation on a fresh branch. Do not auto-advance to P04.03.
