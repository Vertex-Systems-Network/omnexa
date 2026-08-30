# P04 Entry Gate — Activation Candidate

Status: **SATISFIED**

Phase: `P04 — Data, Jobs & Event Fabric`
Activation base: `f97ee9a509f14247771626dc89b685e22ecd9457`
Activation carrier: GitHub PR `#144`
Readiness source: `docs/governance/P03_P04_TRANSITION_READINESS.md`
Accepted readiness merge: `384c82b6d5b8bd89b70eb94640ea3be7109bc537`
Bounded planning promotion: `73561ad68e453c2741dec82d0b6fb01a7dcea1be`
Historical-validator prerequisite merge: `f97ee9a509f14247771626dc89b685e22ecd9457`
Gap ledger: `docs/governance/P04_ACTIVATION_GAP_LEDGER.md`
Transaction template: `docs/governance/P04_ACTIVATION_TRANSACTION_TEMPLATE.md`
Package sequence: `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`
First package: `docs/roadmap/work-packages/P04.01.md`

## Preconditions re-verified on activation base

- P03 exit is `SATISFIED`.
- activation base was terminal `current_phase=P03` with `current_work_package=null`;
- kernel implementation authority was false before this activation transaction;
- business-feature implementation authority was false and remains false;
- accepted P04 readiness/planning artifacts are present on protected `main`;
- retained P01/P02/P03/foundation validators have a reviewed fail-closed historical mode compatible with active P04;
- P04 activation validator exists in canonical Governance.

If protected `main` changes materially before activation merge, this gate loses merge authority until the carrier is brought current and exact-head Governance is rerun.

## Accepted entry semantics

This activation authorizes **P04.01 only**.

The following semantics are binding:

- duplicates are possible;
- consumers must be duplicate-safe where mutation is later introduced;
- no global ordering assumption;
- globally stable event identity;
- explicit fail-closed tenant context;
- first-class correlation and causation identity;
- versioned bounded event envelope;
- secrets prohibited from payload/diagnostic surfaces;
- event transport never grants cross-module write authority;
- no broker/vendor selection in P04.01.

## P04.01 implementation authority boundary

After this state-only activation carrier merges and protected-main state is re-read, a **new separate implementation branch** may implement only the event-envelope/identity contract and focused tests defined by `docs/roadmap/work-packages/P04.01.md`.

Explicitly unauthorized at entry:

- `P04.02-P04.10` implementation;
- broker/stream vendor integration;
- outbox/inbox database schema or migrations;
- durable consumer/checkpoint runtime;
- retry/DLQ/quarantine runtime;
- business-domain handlers;
- background-job execution changes;
- business-feature code.

## Activation evidence requirements

This gate is semantically SATISFIED by the candidate transaction, but merge authority remains fail-closed until all are true on the exact final PR head:

1. canonical `docs/roadmap/STATE.json` records `current_phase=P04`, `current_work_package=P04.01`, bounded kernel authority true and business-feature authority false;
2. package sequence records P04 active, P04.01 active and P04.02-P04.10 planned/locked;
3. changed-file review proves governance/state/continuity scope only;
4. exact-head Omnexa Governance passes;
5. review threads/conversations are resolved;
6. review provenance is independent review when available or explicitly `SELF REVIEW` otherwise;
7. activation base `f97ee9a509f14247771626dc89b685e22ecd9457`, PR #144 and accepted readiness/prerequisite merges remain recorded;
8. protected `main` is re-read after merge before any P04.01 implementation branch is created.

## Rollback

Before P04.01 implementation begins, an erroneous state-only activation may be reverted by a governed rollback restoring the terminal P03 state and locking P04.

After P04.01 implementation begins, do not casually roll back activation underneath accepted source changes; prefer forward correction or a coordinated reverse transition with explicit evidence.

## Current decision

`SATISFIED CANDIDATE — MERGE AUTHORITY REQUIRES EXACT-HEAD GOVERNANCE`.

This file records the accepted entry contract; it is not permission to implement from the activation branch itself.
