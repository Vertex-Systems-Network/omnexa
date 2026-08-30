# P04 Entry Gate — Activation Candidate

Status: **CANDIDATE — NOT SATISFIED UNTIL CANONICAL STATE TRANSITION MERGES**

Phase: `P04 — Data, Jobs & Event Fabric`
Activation base: `384c82b6d5b8bd89b70eb94640ea3be7109bc537`
Readiness source: `docs/governance/P03_P04_TRANSITION_READINESS.md`
Gap ledger: `docs/governance/P04_ACTIVATION_GAP_LEDGER.md`
Transaction template: `docs/governance/P04_ACTIVATION_TRANSACTION_TEMPLATE.md`
Package sequence candidate: `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`
First package: `docs/roadmap/work-packages/P04.01.md`

## Preconditions re-verified on activation base

- P03 exit is `SATISFIED`.
- `current_phase` is `P03`.
- `current_work_package` is `null`.
- kernel implementation authority is currently `false`.
- business-feature implementation authority is `false`.
- accepted P04 readiness documents are present on protected `main` through merge `384c82b6d5b8bd89b70eb94640ea3be7109bc537`.

If protected `main` changes materially before activation merge, this gate must be re-read and the carrier brought current before acceptance.

## Accepted entry semantics candidate

The activation transaction is allowed to authorize only `P04.01` and only after canonical `STATE.json` records the transition.

The following semantics are accepted for the first package contract:

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

After this entry gate and canonical state transition are accepted, P04.01 may implement only the event-envelope/identity contract and focused tests defined by `docs/roadmap/work-packages/P04.01.md`.

Explicitly unauthorized at entry:

- `P04.02-P04.10` implementation;
- broker/stream vendor integration;
- outbox/inbox database schema or migrations;
- durable consumer/checkpoint runtime;
- retry/DLQ/quarantine runtime;
- business-domain handlers;
- background-job execution changes;
- business-feature code.

## Evidence required before gate can become SATISFIED

1. canonical `docs/roadmap/STATE.json` transition is included in this same governed activation carrier;
2. state validates under repository governance with `current_phase=P04`, `current_work_package=P04.01`, minimum P04.01 kernel authority true and business-feature authority false;
3. P04 package sequence is accepted with P04.01 active and P04.02-P04.10 locked;
4. changed-file review proves activation/governance scope only;
5. exact-head Omnexa Governance passes;
6. review threads are resolved;
7. review provenance is recorded as independent review when available or explicitly `SELF REVIEW` otherwise;
8. activation base/final head/PR and accepted readiness merge are recorded in the canonical evidence fields.

## Rollback

Before P04.01 implementation begins, an erroneous state-only activation may be reverted by a governed rollback restoring P03 terminal state and locking all P04 packages.

After P04.01 implementation begins, do not casually roll back activation underneath accepted source changes; prefer forward correction or a coordinated reverse transition with explicit evidence.

## Current decision

`NOT YET SATISFIED` — package sequence and P04.01 contract are prepared, but canonical `STATE.json` has not yet been changed on this carrier.
