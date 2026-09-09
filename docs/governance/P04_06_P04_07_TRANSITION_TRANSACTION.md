# P04.06 Closure / P04.07 Activation Transaction

Status: **CANDIDATE — EXACT-HEAD GOVERNANCE / UNCHANGED PROMOTION / PROTECTED READ-BACK REQUIRED**

- Coordination issue: `#264`
- Transaction base: protected `main@23f3dca090338b5debbf1af02b6db49b2f03f30b`
- Candidate package result: `P04.01-P04.06 DONE`, `P04.07 sole ACTIVE`, `P04.08-P04.10 PLANNED / LOCKED`
- Runtime/schema/provider mutation in this transaction: `none`
- P04.07 worker slots/tasks/branches: `0`
- Migration 4 reservation: `none`
- Provider/vendor schema-registry selection: `none`
- Business-feature authority: `false`
- AI/model/agent product-runtime authority: `false`

This transaction changes governance/state/continuity only. It does not implement a schema registry, compatibility engine, payload validator, registry persistence, migration 4, remote-reference resolver, provider integration, background job, business feature or AI/model/agent product runtime.

## Preconditions

Protected main already contains both required predecessor facts:

1. accepted P04.06 implementation/verifier evidence plus released Wave 2D leases;
2. separately prepared and governed P04.07 contract/handoff.

No package may auto-advance merely because these files exist. This transaction is the explicit sequential closure/activation decision.

## Accepted P04.06 implementation/completion evidence

P04.06 implementation culminates in the governed T07 Supervisor verifier/evidence chain:

- T07 source PR `#257`;
- unchanged T07 promotion PR `#258`;
- exact T07 source/promotion head `b0d047c93e8c53e9da62a96e36c4a1a8a7ce634f`;
- source Governance `#738 / 34400644888` — PASS;
- promotion Governance `#739 / 34401644140`, job `102634706533` — PASS;
- expected-head implementation merge/read-back `bdb96bb7f0dabf5b103acd78335599cc59e92b69`;
- completion evidence `docs/roadmap/evidence/P04.06_COMPLETION_2026-09-09.md`;
- Wave 2D ledger closure source `#259`;
- unchanged ledger closure promotion `#260`;
- zero-slot/zero-task protected-main reconciliation `6050bc2d970b5cf83118115557e952e3e91e0f96`.

Accepted P04.06 includes:

- structured deterministic failure disposition with unknown failure fail-closed behavior;
- finite attempt budgets and capped deterministic backoff;
- authoritative UTC retry eligibility;
- durable scheduled/quarantined/resolved state;
- immutable `kernel.events` migration 3;
- PostgreSQL exact-identity load/create, revision CAS, due discovery and bounded claims;
- authoritative claim gating before retry execution;
- interruption without false attempt consumption;
- quarantine before checkpoint advancement;
- quarantine/checkpoint crash-gap recovery without handler replay;
- owner/consumer/route/stream/partition/tenant fail-closed isolation;
- P04.05 already-applied inbox composition into retry resolution;
- real-PostgreSQL restart preservation.

P04.03 checkpoint, P04.05 inbox completion and P04.06 retry/quarantine remain distinct facts. No P04.06 state grants authorization, business success, global ordering or end-to-end exactly-once semantics.

## Accepted P04.07 preparation evidence

P04.07 preparation was separately governed before this activation transaction:

- source preparation PR `#262`;
- unchanged preparation promotion PR `#263`;
- exact source/promotion head `94789264824a34e3f2608283cb6b1c2f158e0b81`;
- source Governance `#744 / 34407394506` — PASS;
- promotion Governance `#745 / 34408137084` — PASS;
- preparation merge/read-back `23f3dca090338b5debbf1af02b6db49b2f03f30b`;
- accepted contract `docs/roadmap/work-packages/P04.07.md`;
- handoff `docs/ai/handoffs/P04.07.md`.

Preparation selected no provider/vendor registry and reserved no migration.

## Candidate canonical result

Only after this exact source head passes fresh Governance, receives honest exact-head SELF/Supervisor review, is promoted byte-for-byte unchanged through a second fresh Governance run, merges with expected-head protection and protected main is re-read may canonical state be interpreted as:

- P04 remains ACTIVE;
- P04 progress becomes `6 / 10 done`;
- P04.01-P04.06 are DONE with retained evidence;
- P04.07 is the sole ACTIVE package;
- P04.08-P04.10 remain PLANNED / LOCKED;
- `kernel_code_authorized=true` is bounded to P04.07 only;
- `business_feature_code_authorized=false`;
- strategic X runtime remains unauthorized;
- AI/model/agent product runtime remains unauthorized.

P04.07 being ACTIVE is **not** equivalent to a runtime worker lease. Runtime remains blocked until a separate post-activation continuity/worker plan is itself governed, promoted and read back.

## P04.06 closure law

After accepted transition read-back:

- P04.06 becomes historical DONE with `docs/roadmap/evidence/P04.06_COMPLETION_2026-09-09.md` retained;
- its accepted runtime/verifier files remain regression authority;
- `kernel.events` migration 3 remains immutable historical retry/quarantine ownership evidence;
- no P04.06 worker/Supervisor lease remains active;
- P04.06 completion does not grant P04.07+ runtime beyond the explicit activation/continuity rules.

## Activated P04.07 contract boundary

After accepted transition and later separate post-activation continuity/worker-plan acceptance, P04.07 may implement only the provider-neutral contract already prepared in `docs/roadmap/work-packages/P04.07.md`:

- explicit payload schema version separate from P04.01 envelope version;
- stable owner/event-type/version schema identity;
- immutable accepted historical schema versions and deterministic fingerprints;
- deterministic compatibility policy with `backward`, `forward`, `full` and `exact` modes;
- declared predecessor comparison-set semantics;
- unknown/unregistered/unsupported schema identity/version failing before protected handler mutation;
- bounded deterministic payload validation;
- no remote HTTP(S), DNS, arbitrary filesystem or registry-to-registry schema/reference fetching;
- no dynamic code, eval, arbitrary plugin, shell or schema callback execution;
- no schema-derived authorization/capabilities/tenant membership;
- owner/producer/tenant isolation with no tenant-specific V1 contract forks;
- strict separation from P04.03 checkpoint, P04.05 inbox and P04.06 retry/quarantine facts.

Validation success only proves payload conformance to an accepted schema contract. It grants no authorization and proves neither exactly-once processing nor global ordering.

## Zero-runtime activation law

This transaction intentionally creates:

- zero P04.07 worker slots;
- zero P04.07 runtime tasks;
- zero P04.07 runtime branches;
- zero P04.07 migration reservations;
- zero provider/vendor schema-registry selections.

The active multi-agent plan remains a zero-lease activation plan. No runtime worker may start merely because P04.07 becomes canonical ACTIVE.

## Persistence / migration gate remains closed

Accepted `kernel.events` migration 3 remains immutable predecessor evidence. This transaction does **not** reserve, assume or authorize migration 4.

After transition read-back, a later separate governed P04.07 post-activation continuity/worker plan must re-read fresh protected main and explicitly decide whether registry state is:

- generated/static repository-local registration;
- durable PostgreSQL-backed state; or
- another provider-neutral local representation consistent with the accepted P04.07 contract.

If durable persistence is required, before any schema mutation that later plan must record and govern:

1. exact current `kernel.events` migration ledger;
2. exact next immutable migration owner/version/path/name;
3. exact schema/table/index/constraint/data budget;
4. immutable event owner/type/version fingerprint uniqueness law;
5. bounded canonical schema storage and limits;
6. compatibility-policy/version persistence semantics;
7. owner/producer/tenant isolation;
8. local-reference graph bounds if references are supported;
9. fresh-install and supported-upgrade/ledger evidence;
10. rollback/forward-recovery behavior that cannot remap accepted historical schema versions.

No vendor-hosted schema registry is authorized by activation.

## Exact transition path budget

Issue `#264` authorizes changes only to:

1. `docs/roadmap/STATE.json`
2. `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`
7. `docs/governance/P04_06_P04_07_TRANSITION_TRANSACTION.md`

The source/promotion must not modify:

- `AGENTS.md`;
- `README.md`;
- accepted P04.06 evidence;
- accepted P04.07 spec/handoff;
- runtime Go source;
- migrations;
- workflows;
- modules.

Any later subordinate wording reconciliation belongs to the separately governed post-activation continuity carrier.

## Hard non-scope

Do not implement or authorize on this transaction:

- P04.07 registry/compatibility/validator runtime;
- migration 4/schema mutation;
- JSON Schema/Avro/Protobuf or another technology selection as implementation authority;
- Confluent/Kafka/RabbitMQ/NATS or another provider/vendor registry selection;
- remote schema/reference fetching;
- dynamic schema execution;
- P04.08 background-job ownership/execution;
- P04.09 broad operator recovery/requeue/replay UX;
- P04.10 replay/release/poison-event aggregate runtime;
- production business handlers/features;
- cross-module private-table writes;
- global ordering or exactly-once transport/external-side-effect claims;
- strategic X runtime;
- AI/model/agent product runtime.

## Required integration path

1. Source transition branch from exact protected `main@23f3dca090338b5debbf1af02b6db49b2f03f30b`.
2. Reconcile only the seven authorized transition files.
3. Require fresh exact-head Omnexa Governance.
4. Record substantive SELF/Supervisor review with honest non-independent provenance and zero unresolved threads.
5. Verify protected main still equals the expected base before promotion.
6. Create a separate promotion branch at the exact unchanged source head.
7. Require fresh promotion-context Omnexa Governance.
8. Record promotion SELF/Supervisor review, zero unresolved threads and exact unchanged-head verification.
9. Merge with expected-head protection.
10. Re-read protected main, canonical `STATE.json`, package sequence and active-plan authority.
11. Confirm P04.01-P04.06 DONE, P04.07 sole ACTIVE at `6 / 10`, P04.08-P04.10 locked and zero runtime leases/reservations.
12. Govern a separate P04.07 post-activation continuity/worker-plan carrier.
13. Only after that later carrier is accepted/read back may an exact bounded P04.07 runtime lease or migration preflight begin.

Do not auto-advance P04.08.
