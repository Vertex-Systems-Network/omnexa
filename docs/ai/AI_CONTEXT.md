# Omnexa AI Project Context

Status: **P04.07 POST-ACTIVATION CONTINUITY — subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative protected-main base

Fresh protected main at this continuity start is:

`5af9383c3e973c4055eea66d48462b7c9a2a5858`

That protected-main state contains the accepted P04.06 closure / P04.07 activation transaction:

- source PR `#265`;
- unchanged promotion PR `#266`;
- exact source/promotion head `02b58b4245da5eef7a3ab1090698cd10a4832d90`;
- source Governance `#747 / 34411353055` — PASS;
- promotion Governance `#748 / 34412007672` — PASS;
- honest SELF/Supervisor review provenance with no fabricated independent approval;
- zero unresolved review threads before protected integration;
- expected-head guarded merge/read-back `5af9383c3e973c4055eea66d48462b7c9a2a5858`.

The activation changed governance/state/continuity only. P04.07 runtime, migration 4 and provider/vendor registry selection were not introduced.

## Current canonical result

Protected-main truth is now:

- P04 ACTIVE — `6 / 10 done`;
- P04.01-P04.06 DONE with retained accepted evidence;
- P04.07 sole ACTIVE package;
- P04.08-P04.10 PLANNED / LOCKED;
- `kernel_code_authorized=true` is bounded to P04.07 only;
- `business_feature_code_authorized=false`;
- P04.07 runtime worker slots/tasks/branches: `0`;
- P04.07 migration reservations: `0`;
- migration 4: not reserved / not authorized;
- provider/vendor schema registry: not selected / not authorized;
- strategic X and AI/model/agent product runtime: unauthorized.

`docs/roadmap/STATE.json` and the P04 package sequence are canonical. This continuity carrier does not mutate either one.

## Accepted P04.06 completion

P04.06 is DONE with completion evidence `docs/roadmap/evidence/P04.06_COMPLETION_2026-09-09.md`.

Terminal accepted implementation facts include:

- final T07 source `#257` / unchanged promotion `#258`;
- exact head `b0d047c93e8c53e9da62a96e36c4a1a8a7ce634f`;
- source Governance `34400644888` and promotion Governance `34401644140` — PASS;
- implementation merge/read-back `bdb96bb7f0dabf5b103acd78335599cc59e92b69`;
- Wave 2D zero-lease reconciliation through `#259/#260` and protected merge `6050bc2d970b5cf83118115557e952e3e91e0f96`;
- immutable accepted `kernel.events` migration 3.

Checkpoint, inbox completion, retry/quarantine and schema validation remain distinct facts. No one of them grants authorization, global ordering or end-to-end exactly-once semantics.

## Accepted P04.07 preparation and activation

P04.07 preparation was accepted before activation:

- source `#262` / unchanged promotion `#263`;
- exact head `94789264824a34e3f2608283cb6b1c2f158e0b81`;
- source Governance `34407394506` — PASS;
- promotion Governance `34408137084` — PASS;
- preparation merge/read-back `23f3dca090338b5debbf1af02b6db49b2f03f30b`;
- contract `docs/roadmap/work-packages/P04.07.md`;
- handoff `docs/ai/handoffs/P04.07.md`.

The later #265/#266 transaction then canonically activated P04.07 without creating runtime authority.

## P04.07 active contract boundary

Owner: `kernel.events`.

The accepted P04.07 contract requires:

- explicit payload-schema version separate from the P04.01 envelope version;
- stable owner/event-type/version schema identity;
- immutable accepted historical schema versions and deterministic fingerprints;
- deterministic compatibility policy with `backward`, `forward`, `full` and `exact` modes;
- an explicitly frozen predecessor comparison set before implementation acceptance;
- unknown/unregistered/unsupported schema identity/version failing before protected mutation;
- bounded deterministic local-only payload validation;
- no remote HTTP(S), DNS, registry-to-registry or arbitrary-filesystem schema/reference fetching;
- no dynamic code/eval/plugin/shell/schema-callback execution;
- no schema-derived authorization/capabilities/tenant membership;
- owner/producer/tenant isolation with no tenant-specific V1 schema forks;
- strict separation from P04.03 checkpoint, P04.05 inbox and P04.06 retry/quarantine facts.

Validation success proves only payload conformance to an accepted schema contract. It does not grant authorization or prove exactly-once/global-ordering semantics.

## Issue #267 continuity-only budget

Issue `#267` may reconcile stale post-activation wording only in:

1. `AGENTS.md`
2. `README.md`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/handoffs/P04.07.md`
7. `docs/governance/P04_06_P04_07_TRANSITION_TRANSACTION.md`

It must not mutate `docs/roadmap/STATE.json`, the P04 package sequence, activation validators, `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`, runtime Go source, migrations, workflows or modules.

## Zero-runtime continuity law

Throughout Issue #267 source/promotion/read-back:

- P04.07 worker slots: `0`;
- P04.07 runtime tasks: `0`;
- P04.07 runtime branches: `0`;
- P04.07 migration reservations: `0`;
- provider/vendor schema-registry selection: `none`.

Accepted migration 3 remains immutable P04.06 history. Migration 4 is not inferred from package activation.

## Separate implementation-plan gate

After Issue #267 continuity is accepted and protected main is re-read, a **new separately governed P04.07 implementation-plan/worker-plan** must decide the first runtime slice before any runtime mutation.

That later plan must explicitly freeze at least:

1. exact implementation paths and non-overlapping worker leases;
2. supported bounded schema representation/dialect or subset;
3. deterministic canonicalization and fingerprint algorithm/encoding;
4. predecessor-set compatibility semantics;
5. validator size/depth/time/memory/collection limits and safe error contract;
6. local-reference policy and graph bounds if references are supported;
7. registry persistence form;
8. exact migration owner/version/path/data budget only if durable persistence is actually required after fresh-main migration preflight;
9. retained P04.01-P04.06 regression requirements;
10. explicit exclusion of P04.08+, business features, provider/vendor registry selection and AI/model/agent product runtime.

No runtime worker may infer migration 4, a vendor registry, or a schema technology from activation or this continuity carrier.

## Locked scope

Do not implement on this continuity source/promotion:

- P04.07 registry/compatibility/validator runtime;
- migration 4 or any schema mutation;
- provider/vendor registry integration;
- remote schema/reference fetching;
- dynamic schema execution;
- P04.08 background-job ownership/runtime;
- P04.09 broad operator recovery;
- P04.10 replay/release/aggregate poison-event runtime;
- production business handlers/features;
- strategic X runtime;
- AI/model/agent product runtime;
- global ordering or end-to-end exactly-once claims.

## Exact next action

1. Complete Issue `#267` continuity source from exact protected `main@5af9383c3e973c4055eea66d48462b7c9a2a5858`.
2. Verify exactly the seven authorized continuity paths changed and no runtime/state/worker-plan drift exists.
3. Require fresh exact-head Omnexa Governance.
4. Record honest SELF/Supervisor review and verify zero unresolved threads.
5. Re-verify protected-main freshness.
6. Promote the exact unchanged continuity head through a separate promotion PR.
7. Require fresh promotion Governance and promotion SELF/Supervisor review.
8. Merge with expected-head protection and re-read protected main.
9. Confirm P04.07 remains sole ACTIVE at `6 / 10` with zero runtime leases/reservations.
10. Only then create the separate P04.07 implementation-plan/worker-plan carrier.

Do not start P04.07 runtime automatically and do not auto-advance P04.08.
