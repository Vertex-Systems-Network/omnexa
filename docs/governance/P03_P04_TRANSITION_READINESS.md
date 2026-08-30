# Omnexa P03 → P04 Transition Readiness

Status: **READINESS IN PROGRESS — NOT AN ACTIVATION**

Owner transition: **P03 Module Runtime → P04 Data, Jobs & Event Fabric**
Working lane: `agent/20260830-omnexa-abd-255`

## Authority boundary

This document is a readiness/preparation artifact only. It does **not**:

- change `docs/roadmap/STATE.json`;
- activate P04;
- authorize kernel or business-feature implementation;
- authorize schema, API, runtime, migration or product behavior changes;
- reuse P03 implementation carriers as P04 authority.

Until a separate canonical activation transaction is accepted:

```text
P03 exit = SATISFIED
P04 = PLANNED / NOT ACTIVATED
current_work_package = NONE
kernel_code_authorized = false
business_feature_code_authorized = false
```

## Canonical P04 purpose

P04 establishes the reliable asynchronous and decoupled communication fabric required by later workflow and business-domain phases. Its terminal proof must show that replay, duplicate delivery, consumer failure, restart and poison-event scenarios cannot double-apply protected business mutations.

## Draft bounded package sequence

The following sequence refines the canonical P04 scope for readiness review. Package identifiers are **draft planning identifiers** until the activation transaction canonically accepts them.

| Draft package | Scope | Primary dependency | Implementation authority now |
| --- | --- | --- | --- |
| `P04.01` | Event envelope, identity, version, tenant and correlation/causation contract | P03 exit | LOCKED |
| `P04.02` | Publish/subscribe abstraction and ownership boundaries | P04.01 | LOCKED |
| `P04.03` | Durable stream/consumer baseline and checkpoint model | P04.01-P04.02 | LOCKED |
| `P04.04` | Transactional outbox reliability primitive | P04.01-P04.03 | LOCKED |
| `P04.05` | Consumer inbox/deduplication + idempotency primitive | P04.01-P04.04 | LOCKED |
| `P04.06` | Retry/backoff, terminal failure and dead-letter/quarantine policy | P04.03-P04.05 | LOCKED |
| `P04.07` | Event schema registry, compatibility and validation | P04.01-P04.06 | LOCKED |
| `P04.08` | Background job ownership, tenant context and event/job correlation | P04.01-P04.07 | LOCKED |
| `P04.09` | Reliability observability, diagnostics and operator recovery contracts | P04.03-P04.08 | LOCKED |
| `P04.10` | Aggregate acceptance: replay/duplicate/failure/restart/poison-event proof | P04.01-P04.09 | LOCKED |

## Readiness decisions required before activation

### R1 — Authority and ownership

- One canonical event envelope authority must be named.
- Event publication ownership must remain with the domain/module that owns the state transition.
- P04 may transport events but must not create cross-domain write authority.
- Migration execution authority must remain consistent with the already accepted P01/P03 migration ownership contracts.

### R2 — Delivery semantics

The activation package must explicitly choose and document the platform guarantee. The design must assume duplicates can occur and make protected mutations idempotent; an undocumented “exactly once” assumption is not acceptable.

Required concepts before implementation:

- stable event ID;
- producer identity;
- tenant/context identity;
- occurred/published timestamps;
- schema/event version;
- correlation and causation IDs;
- deduplication key/window/retention semantics;
- consumer checkpoint semantics;
- retry classification;
- poison-event quarantine/recovery.

### R3 — Transaction boundary

Outbox/inbox behavior must define:

- database transaction boundary;
- crash windows;
- relay retry behavior;
- duplicate relay behavior;
- consumer transaction boundary;
- redelivery after crash;
- reconciliation/repair path;
- retention and cleanup without deleting required audit evidence.

### R4 — Tenant and security boundary

Every package must preserve:

- tenant isolation established by P02;
- deny-by-default authorization where actions are triggered from events/jobs;
- classification-safe payloads and diagnostics;
- no secret values in events, logs, DLQ payload views or package metadata;
- bounded payload size and untrusted-input validation;
- no event-driven bypass around module lifecycle or permission controls.

### R5 — Failure and replay model

Acceptance must include deterministic proof for at least:

1. duplicate publish;
2. duplicate delivery;
3. consumer crash before commit;
4. consumer crash after commit but before acknowledgement/checkpoint;
5. producer restart with pending outbox rows;
6. consumer restart with pending inbox state;
7. poison event and quarantine;
8. retry exhaustion;
9. schema incompatibility;
10. cross-tenant event attempt;
11. replay of historical events;
12. restart/replay without double-applying a protected mutation.

### R6 — Technology decision gate

Readiness must separate **required semantics** from a specific broker/vendor. A broker/stream implementation may be selected only after the abstraction, persistence, operational and local-development requirements are explicit. P04 must not leak vendor-specific lifecycle logic into business modules.

### R7 — Evidence and CI

Before implementation activation, the governed P04 package sequence should define per-package:

- exact owned paths/path budget;
- tests and retained regressions;
- migration impact;
- threat/security impact;
- failure-mode evidence;
- exact-head CI requirements;
- rollback/forward-recovery expectations;
- documentation/state updates allowed at completion.

## Activation prerequisites ledger

| Prerequisite | Readiness state |
| --- | --- |
| P03 exit gate satisfied | PASS |
| P04 canonical purpose/scope exists | PASS |
| P04 implementation currently authorized | **NO** |
| Separate P04 activation transaction accepted | PENDING |
| Canonical P04 package sequence accepted | PENDING |
| First package exact scope/path budget accepted | PENDING |
| Event/delivery semantics decision accepted | PENDING |
| Security/tenant failure model accepted | PENDING |
| CI/evidence plan accepted | PENDING |

## Required activation transaction

The later activation transaction must be a separate governed change that, at minimum:

1. records this readiness review or its accepted successor;
2. creates/accepts the canonical P04 package sequence;
3. sets the current phase/work package intentionally;
4. grants only the minimum implementation authority required for the first P04 package;
5. records exact branch/PR/base evidence;
6. keeps later P04 packages locked until their predecessor evidence passes.

Until that transaction is accepted, this branch remains documentation/readiness-only.
