# Omnexa P04 Activation Gap Ledger

Status: **READINESS EVIDENCE — IMPLEMENTATION REMAINS LOCKED**

Working lane: `agent/20260830-omnexa-abd-255`
Companion: `docs/governance/P03_P04_TRANSITION_READINESS.md`
Canonical authority: protected `main` + `docs/roadmap/STATE.json`

## Current canonical boundary

Fresh protected-main state records:

- `current_phase = P03`
- `current_work_package = null`
- P03 exit gate = satisfied
- `kernel_code_authorized = false`
- `business_feature_code_authorized = false`
- P04 readiness/preparation is allowed
- P04 implementation requires a later, separate governed activation transaction

Therefore this branch cannot authorize event/runtime/schema/migration/business implementation by itself.

## Gap ledger

| ID | Required before P04 implementation | Current state | Fail-closed effect |
| --- | --- | --- | --- |
| `P04-G01` | Separate canonical P04 activation transaction | MISSING | No P04 implementation |
| `P04-G02` | Canonically accepted P04 package-sequence artifact | MISSING | Draft package IDs are planning-only |
| `P04-G03` | Exact first-package owned paths/path budget | MISSING | No source-tree mutation authorized |
| `P04-G04` | Delivery guarantee and duplicate-delivery semantics decision | MISSING | No broker/stream implementation selection |
| `P04-G05` | Event envelope authority and versioning contract | MISSING | No event schema becomes canonical |
| `P04-G06` | Outbox/inbox transaction and crash-window contract | MISSING | No reliability persistence migration |
| `P04-G07` | Tenant/security propagation and fail-closed consumer rules | MISSING | No background/event-triggered privileged action |
| `P04-G08` | Retry, poison-event, DLQ/quarantine and operator recovery contract | MISSING | No terminal-failure automation |
| `P04-G09` | Exact-head CI/evidence requirements per package | MISSING | No package acceptance claim |
| `P04-G10` | Replay/duplicate/restart aggregate acceptance plan | DRAFT ONLY | No P04 exit claim |

## Readiness conclusions

### 1. The current draft sequence is useful but not canonical

`P03_P04_TRANSITION_READINESS.md` proposes `P04.01` through `P04.10`. Those identifiers remain planning labels until a separate activation change writes the accepted sequence into the canonical roadmap/governance structure.

### 2. Technology selection is intentionally premature

No Kafka/Redpanda/NATS/RabbitMQ/cloud-broker/database-queue choice is authorized from readiness alone. Required delivery, persistence, local-development, operational and recovery semantics must be accepted before vendor selection.

### 3. "Exactly once" must not be an unstated platform assumption

The design must tolerate duplicate publication/delivery and prove protected mutations remain idempotent across crash/replay windows. A provider marketing guarantee is not sufficient acceptance evidence.

### 4. Event transport cannot create write authority

The module/domain that owns a state transition remains the authority for that transition. P04 transport, consumers and jobs may carry work but must not bypass module lifecycle, tenant isolation or authorization rules established by prior phases.

## Minimum activation transaction contents

A later activation carrier must atomically and explicitly record:

1. the exact protected-main base being activated;
2. the canonical P04 package sequence;
3. the first active package only;
4. exact owned paths/path budget for that package;
5. the event/delivery semantics decisions needed by that package;
6. tenant/security/failure constraints;
7. migration/API impact, if any;
8. exact CI and retained-regression evidence requirements;
9. rollback or forward-recovery expectations;
10. `docs/roadmap/STATE.json` authorization changes.

If the protected-main base changes materially before activation, the activation carrier must re-read and reconcile instead of assuming this readiness branch is still authoritative.

## Explicit non-authorization

This ledger does not authorize:

- kernel code;
- business feature code;
- event runtime code;
- broker integration;
- schema or migration changes;
- API behavior changes;
- background-job execution changes;
- P04 activation or completion claims.
