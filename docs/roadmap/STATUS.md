# Omnexa Program Status

Last reconciled: **2026-08-31 post-activation read-back**

## Current authoritative position

- Program: **Kernel Program**
- Phase: **P04 — Data, Jobs & Event Fabric**
- Phase state: **ACTIVE**
- P04 progress: **1 / 10 done**
- P04.01: **DONE — accepted implementation evidence recorded**
- Current work package: **P04.02 — Publish/Subscribe Abstraction & Ownership Boundaries**
- P04.02: **SOLE ACTIVE PACKAGE**
- P04.03-P04.10: **PLANNED / LOCKED**
- P04 entry gate: **SATISFIED**
- P03: **DONE — 11 / 11**; exit **SATISFIED / HISTORICAL**
- P02: **DONE — 10 / 10**; exit **SATISFIED**
- P01: **DONE — 12 / 12**; exit **SATISFIED**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Canonical CI: **GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation authority: **P04.02 ONLY**
- Business-feature implementation: **NOT AUTHORIZED**

Protected `main` read-back after the P04.01 closure / P04.02 activation transaction is `226f5236d80c806a7b07d65e7870981de320c05c`. Canonical `docs/roadmap/STATE.json` on that SHA confirms `current_phase=P04`, `current_work_package=P04.02`, `kernel_code_authorized=true` for the bounded P04.02 contract, and `business_feature_code_authorized=false`.

## Accepted P04.01 / P04.02 transition evidence chain

P04.01 implementation is complete and accepted:

- P04 activation / P04.01 authorization merge: `dc269932481a1e742926ac247f28f94bdf43937d`;
- implementation exact head: `3c4ee6e79b76a6f042f3ee77b518943f3d7064e4`;
- implementation promotion PR: `#147`;
- promotion-specific Omnexa Governance: `#548` / run `33331934481`;
- implementation merge: `6b1a01c009e9a08d45613a00e9b63c0c272bf020`;
- completion evidence: `docs/roadmap/evidence/P04.01_COMPLETION_2026-08-31.md`;
- P04.02 contract/handoff preparation promotion PR: `#149`;
- preparation promotion Governance: `#551`;
- preparation merge: `16dfca22cd0430b8b2727e55628cc6fba381b5ef`;
- P04.01 closure / P04.02 activation source PR: `#150`;
- source exact head: `8f00090a2055a2283d61b30a66b0772ddf354648`;
- source Governance: `#554` / run `33334373676`;
- activation promotion PR: `#151`;
- promotion Governance: `#555` / run `33335001265`;
- accepted transition merge / authoritative main read-back: `226f5236d80c806a7b07d65e7870981de320c05c`.

Historical P01.01-P01.12, P02.01-P02.10 and P03.01-P03.11 completion evidence remains immutable and all retained regression verifiers remain mandatory.

## P04.02 authoritative implementation boundary

P04.02 is provider-neutral and contract-first. A **new separate implementation branch from exact accepted main** may implement only the accepted `docs/roadmap/work-packages/P04.02.md` contract.

Authorized P04.02 scope is limited to:

- publish/subscribe abstraction and explicit ownership boundaries;
- event publisher/subscriber interfaces and deterministic in-process/reference behavior where the accepted spec permits it;
- explicit tenant/correlation propagation using the completed P04.01 envelope;
- fail-closed module/owner identity and authorization boundaries;
- focused positive/adversarial tests and package verification.

Explicitly unauthorized:

- Kafka, NATS, RabbitMQ, Redis Streams or another broker/provider selection;
- durable stream/consumer/checkpoint runtime;
- database migrations;
- outbox/inbox persistence;
- retry/backoff/DLQ/quarantine runtime;
- schema-registry runtime;
- background-job execution changes;
- business-domain event handlers;
- P04.03+ implementation;
- business-feature code;
- strategic X-program runtime;
- AI/model/agent runtime.

## State/continuity rule

`docs/roadmap/STATE.json`, `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`, `AGENTS.md`, this status document and `docs/ai/` continuity snapshots must agree on:

```text
P04 ACTIVE — 1 / 10 done
P04.01 DONE
P04.02 SOLE ACTIVE PACKAGE
P04.03-P04.10 PLANNED / LOCKED
kernel_code_authorized=true — P04.02 only
business_feature_code_authorized=false
```

Any future contradictory cursor/authority state is a stop-the-line reconciliation defect.

## Protected integration / CI

`main` remains governed through PR-only integration. Canonical acceptance evidence is GitHub-hosted `ubuntu-24.04` only. Source Governance #554 and promotion Governance #555 are the accepted P04.01→P04.02 transition evidence; future P04.02 implementation must obtain its own exact-head Governance and later separate completion transition.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions and grants no kernel/business implementation authority.

## Exact next work

1. Complete this bounded post-merge continuity correction so subordinate snapshots stop describing an already-merged transaction as a candidate.
2. Require exact-head Omnexa Governance and governed promotion/merge for the continuity-only correction.
3. Re-read protected `main` and canonical `STATE.json`; confirm P04.02 remains sole active and all locks remain intact.
4. Create a new isolated P04.02 implementation branch from that exact corrected `main` SHA.
5. Implement only the accepted P04.02 contract; do not auto-advance to P04.03.
