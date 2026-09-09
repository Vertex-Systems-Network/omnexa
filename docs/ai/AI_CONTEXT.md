# Omnexa AI Project Context

Status: **P04.06 CLOSURE / P04.07 ACTIVATION CANDIDATE — subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Authoritative candidate base

Fresh protected main at this transition start is:

`23f3dca090338b5debbf1af02b6db49b2f03f30b`

That protected-main state contains both prerequisites for a separate P04.06 closure / P04.07 activation decision:

1. accepted P04.06 implementation/verifier evidence with all Wave 2D leases released;
2. separately governed P04.07 preparation accepted through source `#262`, unchanged promotion `#263`, and protected-main read-back `23f3dca090338b5debbf1af02b6db49b2f03f30b`.

Preparation alone did not activate P04.07. Issue `#264` is the explicit sequential package-transition carrier.

## Accepted P04.06 implementation/completion evidence

P04.06 remains canonical ACTIVE at the transition base, but its bounded implementation evidence is complete:

- final T07 source PR `#257` / unchanged promotion `#258`;
- exact source/promotion head `b0d047c93e8c53e9da62a96e36c4a1a8a7ce634f`;
- source Governance `#738 / 34400644888` — PASS;
- promotion Governance `#739 / 34401644140`, job `102634706533` — PASS;
- protected merge/read-back `bdb96bb7f0dabf5b103acd78335599cc59e92b69`;
- completion evidence `docs/roadmap/evidence/P04.06_COMPLETION_2026-09-09.md`;
- Wave 2D ledger closure source `#259` / promotion `#260`;
- zero-lease protected-main reconciliation `6050bc2d970b5cf83118115557e952e3e91e0f96`.

Accepted P04.06 runtime includes deterministic failure disposition/backoff, durable scheduled/quarantined/resolved retry state, immutable `kernel.events` migration 3, PostgreSQL CAS/claim/due-discovery semantics, checkpoint/quarantine crash-gap recovery, P04.05 already-applied composition and real-PostgreSQL restart persistence.

Checkpoint, inbox completion and retry/quarantine remain separate facts. No P04.06 fact grants authorization, business success, global ordering or end-to-end exactly-once semantics.

## Accepted P04.07 preparation

P04.07 preparation is accepted but still not active at the transition base:

- source preparation PR `#262`;
- unchanged promotion PR `#263`;
- exact source/promotion head `94789264824a34e3f2608283cb6b1c2f158e0b81`;
- source Governance `#744 / 34407394506` — PASS;
- promotion Governance `#745 / 34408137084` — PASS;
- preparation protected-main merge/read-back `23f3dca090338b5debbf1af02b6db49b2f03f30b`;
- contract `docs/roadmap/work-packages/P04.07.md`;
- handoff `docs/ai/handoffs/P04.07.md`.

Prepared P04.07 boundaries include:

- explicit payload-schema version separate from the P04.01 envelope version;
- immutable accepted historical schema versions and deterministic fingerprints;
- explicit deterministic `backward`, `forward`, `full` and `exact` compatibility modes;
- unknown/unregistered/unsupported schema identity/version failing before protected mutation;
- bounded deterministic local-only payload validation;
- no remote schema/reference fetching;
- no dynamic code/plugin/shell execution through schema content;
- no schema-derived authorization;
- owner/producer/tenant isolation and no tenant-specific schema forks in V1;
- strict separation from P04.03 checkpoint, P04.05 inbox and P04.06 retry/quarantine facts;
- provider/vendor registry technology left unselected;
- persistence and migration held behind a later explicit decision gate.

## Candidate canonical result

Only after Issue `#264` source exact-head Governance, honest SELF/Supervisor review, unchanged promotion Governance, zero unresolved threads, protected-main freshness, expected-head guarded merge and protected-main read-back may the repository treat the following as canonical:

- P04 ACTIVE — `6 / 10 done`;
- P04.01-P04.06 DONE with retained evidence;
- P04.07 sole ACTIVE package;
- P04.08-P04.10 PLANNED / LOCKED;
- `kernel_code_authorized=true` bounded to P04.07 only;
- `business_feature_code_authorized=false`;
- zero P04.07 runtime worker slots/tasks/branches;
- migration 4 not reserved or authorized;
- provider/vendor registry selection remains unauthorized;
- strategic X and AI/model/agent product runtime remain unauthorized.

## Zero-runtime activation law

This transition is governance/state/continuity only. It must not implement P04.07 runtime.

At candidate transition and immediately after accepted read-back:

- P04.07 runtime worker slots: `0`;
- P04.07 runtime tasks: `0`;
- P04.07 runtime branches: `0`;
- P04.07 migration reservations: `0`;
- provider/vendor schema registry: `none`.

Accepted `kernel.events` migration 3 remains immutable P04.06 history. Migration 4 is not assumed, reserved or authorized by activation.

## Post-activation implementation gate

P04.07 becoming ACTIVE does not itself create runtime authority for a branch or migration.

After accepted transition read-back, a separate governed post-activation continuity/worker-plan carrier must:

1. re-read fresh protected main, `STATE.json`, package sequence, P04.07 contract/handoff and current `kernel.events` migration history;
2. reconcile subordinate human/AI continuity wording where explicitly authorized;
3. decide the bounded implementation slice and exact non-overlapping worker leases;
4. explicitly decide whether registry persistence is repository-static/generated, PostgreSQL-backed or another local provider-neutral representation;
5. if durable persistence is required, separately record the exact next migration owner/version/path/data budget before schema mutation;
6. require exact-head Governance, review, unchanged promotion and protected read-back;
7. only then permit the exact P04.07 runtime lease.

No implementation worker may infer migration 4 or a provider registry from this transition.

## Transition path budget

Issue `#264` authorizes changes only to:

1. `docs/roadmap/STATE.json`
2. `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`
7. `docs/governance/P04_06_P04_07_TRANSITION_TRANSACTION.md`

The transition must not modify runtime Go source, migrations, workflows, modules, README, `AGENTS.md`, accepted P04.06 evidence, or the accepted P04.07 contract/handoff.

## Locked scope

Do not implement on this source/promotion:

- P04.07 registry/compatibility/validator runtime;
- migration 4 or any schema mutation;
- provider/vendor registry integration;
- remote schema/reference fetching;
- dynamic schema execution;
- P04.08 background-job ownership/runtime;
- P04.09 broad operator recovery;
- P04.10 replay/release/aggregate poison-event runtime;
- business handlers/features;
- strategic X runtime;
- AI/model/agent product runtime;
- global ordering or end-to-end exactly-once claims.

## Exact next action

1. Complete the seven-file Issue `#264` transition source from exact protected `main@23f3dca090338b5debbf1af02b6db49b2f03f30b`.
2. Require fresh exact-head Omnexa Governance.
3. Record honest SELF/Supervisor review and verify zero unresolved threads.
4. Re-verify protected-main freshness.
5. Promote the exact unchanged source head through a separate promotion PR.
6. Require fresh promotion-specific Governance and promotion SELF/Supervisor review.
7. Merge with expected-head protection and re-read protected main.
8. Confirm P04.01-P04.06 DONE, P04.07 sole ACTIVE at `6 / 10`, P04.08+ locked and zero runtime leases/reservations.
9. Only then govern the separate P04.07 post-activation continuity/worker plan.

Do not start P04.07 runtime automatically.
