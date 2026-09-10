# P04.06 Closure / P04.07 Activation Transaction

Status: **ACCEPTED — PROTECTED-MAIN READ-BACK COMPLETE**

## Acceptance receipt

- Coordination issue: `#264` — COMPLETED
- Transaction base: protected `main@23f3dca090338b5debbf1af02b6db49b2f03f30b`
- Source PR: `#265`
- Unchanged promotion PR: `#266`
- Exact source/promotion head: `02b58b4245da5eef7a3ab1090698cd10a4832d90`
- Source Governance: `#747 / 34411353055` — PASS
- Promotion Governance: `#748 / 34412007672` — PASS
- Review provenance: honest SELF/Supervisor review; independent approval was not claimed
- Unresolved review threads before integration: `0`
- Protected-main freshness immediately before merge: `23f3dca090338b5debbf1af02b6db49b2f03f30b`
- Expected-head guarded merge/read-back: `5af9383c3e973c4055eea66d48462b7c9a2a5858`
- Accepted package result: `P04.01-P04.06 DONE`, `P04.07 sole ACTIVE`, `P04.08-P04.10 PLANNED / LOCKED`
- P04 progress: `6 / 10 done`
- Runtime/schema/provider mutation in the transaction: `none`
- P04.07 worker slots/tasks/runtime branches: `0`
- Migration 4 reservation: `none`
- Provider/vendor schema-registry selection: `none`
- Business-feature authority: `false`
- AI/model/agent product-runtime authority: `false`

This transaction changed governance/state/continuity only. It did not implement a schema registry, compatibility engine, payload validator, registry persistence, migration 4, remote-reference resolver, provider integration, background job, business feature or AI/model/agent product runtime.

## Accepted prerequisites

Protected main at the transition base contained both required predecessor facts:

1. accepted P04.06 implementation/verifier evidence plus released Wave 2D leases;
2. separately prepared and governed P04.07 contract/handoff.

No package auto-advanced merely because those files existed. Issue #264 and the #265/#266 protected integration were the explicit sequential closure/activation decision.

## Accepted P04.06 implementation/completion evidence

P04.06 implementation culminated in the governed T07 Supervisor verifier/evidence chain:

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

P04.07 preparation was separately governed before activation:

- source preparation PR `#262`;
- unchanged preparation promotion PR `#263`;
- exact preparation source/promotion head `94789264824a34e3f2608283cb6b1c2f158e0b81`;
- source Governance `#744 / 34407394506` — PASS;
- promotion Governance `#745 / 34408137084` — PASS;
- preparation merge/read-back `23f3dca090338b5debbf1af02b6db49b2f03f30b`;
- accepted contract `docs/roadmap/work-packages/P04.07.md`;
- handoff `docs/ai/handoffs/P04.07.md`.

Preparation selected no provider/vendor registry, no schema technology and reserved no migration.

## Accepted canonical result

Protected `main@5af9383c3e973c4055eea66d48462b7c9a2a5858` canonically records:

- P04 remains ACTIVE;
- P04 progress is `6 / 10 done`;
- P04.01-P04.06 are DONE with retained evidence;
- P04.07 is the sole ACTIVE package;
- P04.08-P04.10 remain PLANNED / LOCKED;
- `kernel_code_authorized=true` is bounded to P04.07 only;
- `business_feature_code_authorized=false`;
- strategic X runtime remains unauthorized;
- AI/model/agent product runtime remains unauthorized.

P04.07 being ACTIVE is **not** equivalent to a runtime worker lease. Runtime remains blocked until post-activation continuity is accepted and a later separate implementation-plan/worker-plan is itself governed and accepted.

## P04.06 closure law

P04.06 is historical DONE with `docs/roadmap/evidence/P04.06_COMPLETION_2026-09-09.md` retained.

- its accepted runtime/verifier files remain regression authority;
- `kernel.events` migration 3 remains immutable historical retry/quarantine ownership evidence;
- no P04.06 worker/Supervisor lease remains active;
- P04.06 completion does not grant P04.07+ runtime beyond explicit later governance.

## Activated P04.07 contract boundary

P04.07 may eventually implement only the provider-neutral contract already accepted in `docs/roadmap/work-packages/P04.07.md`:

- explicit payload schema version separate from P04.01 envelope version;
- stable owner/event-type/version schema identity;
- immutable accepted historical schema versions and deterministic fingerprints;
- deterministic compatibility policy with `backward`, `forward`, `full` and `exact` modes;
- explicitly frozen predecessor comparison-set semantics;
- unknown/unregistered/unsupported schema identity/version failing before protected handler mutation;
- bounded deterministic payload validation;
- no remote HTTP(S), DNS, arbitrary filesystem or registry-to-registry schema/reference fetching;
- no dynamic code, eval, arbitrary plugin, shell or schema callback execution;
- no schema-derived authorization/capabilities/tenant membership;
- owner/producer/tenant isolation with no tenant-specific V1 contract forks;
- strict separation from P04.03 checkpoint, P04.05 inbox and P04.06 retry/quarantine facts.

Validation success only proves payload conformance to an accepted schema contract. It grants no authorization and proves neither exactly-once processing nor global ordering.

## Zero-runtime activation law

The accepted transition created:

- zero P04.07 worker slots;
- zero P04.07 runtime tasks;
- zero P04.07 runtime branches;
- zero P04.07 migration reservations;
- zero provider/vendor schema-registry selections.

The active multi-agent plan remains a zero-lease activation plan. No runtime worker may start merely because P04.07 is canonical ACTIVE.

## Persistence / migration gate remains closed

Accepted `kernel.events` migration 3 remains immutable predecessor evidence. The activation did **not** reserve, assume or authorize migration 4.

A later separately governed P04.07 implementation-plan/worker-plan must re-read fresh protected main and decide whether registry state is:

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

## Post-activation continuity — Issue #267

Issue `#267` is the separate continuity-only reconciliation required after activation read-back.

Its exact path budget is:

1. `AGENTS.md`
2. `README.md`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/handoffs/P04.07.md`
7. `docs/governance/P04_06_P04_07_TRANSITION_TRANSACTION.md`

It must not modify `docs/roadmap/STATE.json`, the P04 package sequence, activation validators, `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`, runtime Go source, migrations, workflows or modules.

Throughout Issue #267 source/promotion/read-back, P04.07 runtime leases/tasks/branches and migration reservations remain zero.

## Hard non-scope until later implementation-plan acceptance

Do not implement or authorize from this transaction or Issue #267 continuity:

- P04.07 registry/compatibility/validator runtime;
- migration 4/schema mutation;
- JSON Schema/Avro/Protobuf or another schema technology as implementation authority;
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

## Next governed sequence

1. Govern Issue #267 continuity source from exact accepted protected main `5af9383c3e973c4055eea66d48462b7c9a2a5858`.
2. Review exact seven-file scope and require zero unresolved threads.
3. Promote the exact unchanged continuity head through fresh Governance.
4. Merge with expected-head protection and re-read protected main.
5. Confirm P04.07 remains sole ACTIVE at `6 / 10` and runtime/migration/provider leases remain zero.
6. Create a **separate** P04.07 implementation-plan/worker-plan from that fresh main.
7. Freeze exact implementation paths/leases, schema representation/canonicalization/fingerprint rules, compatibility predecessor set, validator limits, local-reference policy and persistence choice.
8. Only after that separate plan is governed, promoted and read back may the exact bounded P04.07 runtime lease begin.

Do not auto-advance P04.08.
