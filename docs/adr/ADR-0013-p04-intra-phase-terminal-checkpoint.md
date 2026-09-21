# ADR-0013 — P04 Intra-Phase Terminal Checkpoint

Status: **Accepted**

> This ADR becomes authoritative only after the change-control carrier that contains it passes exact-head Governance, unchanged promotion when required, guarded protected merge and protected-main readback.

## Context

P04 uses a strict sequential one-active-package execution model. The accepted P04.07 closure plan intentionally forbids automatic P04.08 activation in the same transaction.

Fresh P04.07 T05 readiness evidence is accepted on protected main `3f34301f0cd37aa43392c26e906a0a719b6db2a0`:

- dedicated `P04.07 Readiness` run `35655774421` / #2 — PASS;
- job `106519447117` — PASS;
- `bash scripts/verify_p04_07.sh` — PASS;
- verifier gates G0-G11 — PASS;
- source evidence PR #299 Governance #807 — PASS;
- unchanged promotion PR #300 Governance #808 — PASS.

Immediately before T06 closure, a governance contradiction was discovered. `scripts/validate_p04_activation.py` accepts only pre-P04 planning or a P04 state with exactly one active work package. It therefore rejects the required safe state after P04.07 closure where P04 remains active as a phase, no work package is active, and P04.08 remains planned/locked.

Issue #301 governs this change.

## Problem

Without a terminal checkpoint mode, the repository has only unsafe choices:

1. falsely keep P04.07 active after closure;
2. auto-activate P04.08, violating the accepted closure plan; or
3. weaken/bypass the canonical validator.

All three are prohibited.

## Decision

P04 supports a fail-closed **intra-phase terminal checkpoint** between sequential work packages.

For a non-final completed strict prefix:

- `current_phase == "P04"`;
- `current_work_package == null`;
- the P04 phase remains `active`;
- `phases[].P04.active_work_package == null`;
- a non-empty strict prefix of P04 packages is `done`;
- every completed package retains its accepted spec, completion evidence and PASS tracking;
- no package is `active`;
- all remaining packages are `planned` and retain `spec: null` until separately prepared/accepted;
- `implementation_lock.kernel_code_authorized == false`;
- `implementation_lock.business_feature_code_authorized == false`;
- `P04_PACKAGE_SEQUENCE.implementation_authorized == false`;
- the preparation cursor may identify only the next planned package and cannot grant implementation authority;
- P04 entry remains SATISFIED;
- package order and dependency rules remain unchanged;
- no package auto-advances.

The immediate T06 target is P04.01-P04.07 DONE with P04.08-P04.10 PLANNED/LOCKED.

## Implementation

After this ADR/change-control carrier is accepted, one atomic T06 carrier may modify:

- `scripts/validate_p04_activation.py` to add the terminal-checkpoint mode;
- `docs/roadmap/evidence/P04.07_COMPLETION_2026-09-21.md`;
- `docs/roadmap/STATE.json`;
- `docs/roadmap/STATUS.md`;
- `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`;
- `docs/ai/handoffs/P04.07.md`;
- `docs/ai/AI_STATE.yaml`;
- `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`;
- `docs/ai/P04.07_COMPLETION_PLAN.json`;
- `docs/ai/compact/**`;
- `README.md`.

The atomic carrier must run canonical Governance against the actual proposed closure state. Existing planning and one-active-package modes must remain fail-closed and unchanged in meaning.

## Non-decisions

This ADR does **not**:

- activate or prepare P04.08;
- authorize P04.08 runtime or implementation;
- authorize migration 4;
- select a provider/hosted schema registry;
- authorize remote references or external schema I/O;
- change `scripts/verify_p04_07.sh`;
- weaken, skip or remove Governance;
- grant business-feature or AI/model/agent product-runtime authority.

## Alternatives considered

### Auto-activate P04.08 in T06

Rejected. The accepted closure plan explicitly requires a separate later activation transaction.

### Keep P04.07 active after completion

Rejected. That would make canonical state contradict accepted completion evidence.

### Bypass or remove the validator

Rejected. Red governance must not be converted to green by weakening a control.

### Separate validator implementation and closure carriers

Safe but unnecessarily serial. The accepted ADR permits one atomic post-decision T06 carrier so the updated validator is exercised by canonical Governance against the exact closure state it is intended to validate.

## Compatibility impact

No product runtime or public contract changes. This changes only governance-state validation semantics to represent an implementation-locked gap between sequential P04 packages.

## Migration impact

None. Migration 4 remains unreserved and unauthorized.

## Security / tenancy impact

The checkpoint is more restrictive than active-package mode: kernel implementation authority is false, business authority remains false, and future packages remain locked.

## Operational impact

None in product runtime.

## Rollback / forward-fix

Before protected merge, reject/revert the candidate. After protected merge, a later governed change may supersede this ADR; do not bypass the validator.

## Documents/work packages affected

- `scripts/validate_p04_activation.py`
- P04.07 T06 closure plan/evidence/state
- P04 package sequence and progress mirrors
- AI/compact continuity state

Parent coordination: #289  
Change-control issue: #301
