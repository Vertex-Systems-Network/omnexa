# Omnexa Program Status

Last reconciled: **2026-08-31 transition candidate**

## Current candidate position

- Program: **Kernel Program**
- Phase: **P04 — Data, Jobs & Event Fabric**
- Phase state: **ACTIVE**
- P04 progress after this transition merges: **1 / 10 done**
- P04.01: **DONE — accepted implementation evidence recorded**
- Current work package after this transition merges: **P04.02 — Publish/Subscribe Abstraction & Ownership Boundaries**
- P04.02: **SOLE ACTIVE CANDIDATE**
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
- Kernel implementation authority after transition merge + main read-back: **P04.02 ONLY**
- Business-feature implementation: **NOT AUTHORIZED**

This branch is a governance/state/continuity transition candidate. None of the P04.02 authority described here is effective merely because the branch contains candidate state. Authority begins only after the complete exact-final-head transition passes Omnexa Governance, is reviewed, is promoted/merged through the governed path, and protected `main` plus canonical state are re-read.

## Accepted P04.01 evidence chain

P04.01 implementation is complete and already accepted on `main`:

- P04 activation / P04.01 authorization merge: `dc269932481a1e742926ac247f28f94bdf43937d`;
- implementation exact head: `3c4ee6e79b76a6f042f3ee77b518943f3d7064e4`;
- implementation promotion PR: `#147`;
- promotion-specific Omnexa Governance: `#548` / run `33331934481`;
- implementation merge: `6b1a01c009e9a08d45613a00e9b63c0c272bf020`;
- completion evidence: `docs/roadmap/evidence/P04.01_COMPLETION_2026-08-31.md`;
- P04.02 contract/handoff preparation promotion PR: `#149`;
- preparation promotion Governance: `#551`;
- preparation merge / exact base for this transition: `16dfca22cd0430b8b2727e55628cc6fba381b5ef`.

Historical P01.01-P01.12, P02.01-P02.10 and P03.01-P03.11 completion evidence remains immutable and all retained regression verifiers remain mandatory.

## P04.02 candidate authority boundary

P04.02 remains provider-neutral and contract-first. Only after this transition merges and protected `main` is re-read may a **new separate implementation branch** implement the accepted `docs/roadmap/work-packages/P04.02.md` contract.

Candidate P04.02 scope is limited to:

- publish/subscribe abstraction and explicit ownership boundaries;
- event publisher/subscriber interfaces and deterministic in-process/reference behavior where the accepted spec permits it;
- explicit tenant/correlation propagation using the completed P04.01 envelope;
- fail-closed module/owner identity and authorization boundaries;
- focused positive/adversarial tests and package verification.

Explicitly unauthorized by this transition:

- Kafka, NATS, RabbitMQ, Redis Streams or another broker/provider selection;
- durable stream/consumer/checkpoint runtime;
- database migrations;
- outbox/inbox persistence;
- retry/backoff/DLQ/quarantine runtime;
- schema-registry runtime;
- background-job execution changes;
- business-domain event handlers;
- business-feature code;
- strategic X-program runtime;
- AI/model/agent runtime.

## State/continuity reconciliation rule

The complete transition must make `docs/roadmap/STATE.json`, `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`, `AGENTS.md`, this status document and `docs/ai/` continuity snapshots agree on one cursor:

```text
P04 ACTIVE — 1 / 10 done
P04.01 DONE
P04.02 SOLE ACTIVE PACKAGE
P04.03-P04.10 PLANNED / LOCKED
kernel_code_authorized=true — P04.02 only after transition merge + read-back
business_feature_code_authorized=false
```

Any contradictory candidate state is a merge blocker.

## Protected integration / CI

`main` remains governed through PR-only integration. Canonical acceptance evidence is GitHub-hosted `ubuntu-24.04` only. P04.01 implementation Governance #548 and P04.02-preparation Governance #551 are prerequisites and do **not** substitute for this transition carrier's own exact-final-head Omnexa Governance.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions and grants no kernel/business implementation authority.

## Exact next work

1. Reconcile canonical `STATE.json`, package sequence, `AGENTS.md`, this status document and AI continuity to the P04.01-done / P04.02-active candidate.
2. Confirm this carrier contains governance/state/continuity changes only and no P04.02 runtime implementation.
3. Require exact-final-head Omnexa Governance PASS and clean review/thread state.
4. Record SELF REVIEW honestly if no independent reviewer is available.
5. Use a fresh unchanged promotion carrier if required by the repository process.
6. Merge with expected-head guard and re-read protected `main` plus canonical state.
7. Only then create a new isolated P04.02 implementation branch from the exact post-transition `main` SHA.
