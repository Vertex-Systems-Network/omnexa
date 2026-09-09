# Omnexa Roadmap Status

Last reconciled: 2026-09-10 — **P04.07 POST-ACTIVATION CONTINUITY**

## Authoritative protected-main base

Protected main at continuity start:

`5af9383c3e973c4055eea66d48462b7c9a2a5858`

That commit is the accepted P04.06-closure / P04.07-activation merge/read-back from source PR `#265` and unchanged promotion PR `#266` at exact source/promotion head `02b58b4245da5eef7a3ab1090698cd10a4832d90`, with source Governance `#747 / 34411353055` and promotion Governance `#748 / 34412007672` both passing before guarded merge.

Canonical state is now:

- Foundation Architecture v1 remains FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED.
- P04 remains ACTIVE — **6 / 10 done**.
- P04.01-P04.06 are DONE with retained accepted evidence.
- P04.07 is the **sole ACTIVE package**.
- P04.08-P04.10 remain PLANNED / LOCKED.
- `kernel_code_authorized=true` is bounded to P04.07 only.
- P04.07 runtime is still blocked pending accepted Issue #267 continuity and a later separate implementation-plan/worker-plan.
- P04.07 worker slots/tasks/runtime branches remain `0`.
- `business_feature_code_authorized=false`.
- accepted `kernel.events` migration 3 remains immutable P04.06 history.
- migration 4 is not reserved or authorized.
- provider/vendor schema-registry selection remains unauthorized.
- strategic X-program and AI/model/agent product runtime remain unauthorized.
- canonical CI remains GitHub-hosted `ubuntu-24.04` only.

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This status file is a subordinate human-readable mirror and cannot widen authority by itself.

## Accepted P04.06 implementation/completion

P04.06 is DONE with bounded accepted runtime/verifier evidence through the completed multi-wave chain.

Terminal evidence:

- source PR `#257` / unchanged promotion `#258`;
- exact source/promotion head `b0d047c93e8c53e9da62a96e36c4a1a8a7ce634f`;
- source Governance `#738 / 34400644888` — PASS;
- promotion Governance `#739 / 34401644140`, job `102634706533` — PASS;
- implementation merge/read-back `bdb96bb7f0dabf5b103acd78335599cc59e92b69`;
- completion evidence `docs/roadmap/evidence/P04.06_COMPLETION_2026-09-09.md`;
- Wave 2D ledger closure source `#259` / unchanged promotion `#260`;
- zero-lease ledger merge/read-back `6050bc2d970b5cf83118115557e952e3e91e0f96`.

Accepted P04.06 scope includes structured deterministic failure disposition, bounded retry/backoff/eligibility, durable scheduled/quarantined/resolved retry state, immutable `kernel.events` migration 3, PostgreSQL CAS/claim/due-discovery behavior, quarantine-before-checkpoint and crash-gap recovery, P04.05 already-applied composition and real-PostgreSQL restart persistence.

P04.03 checkpoint, P04.05 inbox completion and P04.06 retry/quarantine remain separate facts. Nothing in P04.06 grants authorization, business success, global ordering or end-to-end exactly-once semantics.

## Accepted P04.07 preparation

P04.07 contract preparation was accepted before activation:

- source `#262` / unchanged promotion `#263`;
- exact preparation head `94789264824a34e3f2608283cb6b1c2f158e0b81`;
- source Governance `#744 / 34407394506` — PASS;
- promotion Governance `#745 / 34408137084` — PASS;
- preparation merge/read-back `23f3dca090338b5debbf1af02b6db49b2f03f30b`;
- contract `docs/roadmap/work-packages/P04.07.md`;
- handoff `docs/ai/handoffs/P04.07.md`.

Preparation defined provider-neutral payload-schema identity/versioning, immutable historical fingerprints, deterministic compatibility, bounded local-only validation, owner/producer/tenant isolation, no remote references/dynamic execution/schema-derived authorization, and strict separation from checkpoint/inbox/retry facts.

Preparation reserved no migration and selected no provider/vendor registry.

## Accepted P04.06 closure / P04.07 activation

Coordination issue `#264` is completed.

Accepted transition evidence:

- source PR `#265`;
- unchanged promotion PR `#266`;
- exact unchanged source/promotion head `02b58b4245da5eef7a3ab1090698cd10a4832d90`;
- source Governance `#747 / 34411353055` — PASS;
- promotion Governance `#748 / 34412007672` — PASS;
- SELF/Supervisor review provenance recorded without claiming independent approval;
- zero unresolved review threads before protected integration;
- protected-main freshness verified at `23f3dca090338b5debbf1af02b6db49b2f03f30b`;
- expected-head guarded merge/read-back `5af9383c3e973c4055eea66d48462b7c9a2a5858`.

The transition changed governance/state/continuity only:

- P04.06 moved from ACTIVE to DONE with retained completion evidence;
- P04.07 moved from PLANNED to sole ACTIVE;
- P04 done count moved from 5 to 6;
- bounded kernel package authority moved from P04.06 to P04.07;
- active runtime worker/task/branch count remained zero;
- migration 3 remained immutable historical P04.06 evidence;
- migration 4 remained unreserved;
- provider/vendor schema-registry selection remained unauthorized;
- P04.08-P04.10 remained locked.

## P04.07 active boundary

Owner: `kernel.events`.

P04.07 is canonically ACTIVE, but runtime implementation does not begin from package activation alone. The accepted provider-neutral contract in `docs/roadmap/work-packages/P04.07.md` requires:

- payload schema version separate from P04.01 envelope version;
- stable owner/event-type/version identity and immutable historical fingerprints;
- deterministic `backward`, `forward`, `full` and `exact` compatibility;
- an explicitly frozen predecessor comparison set;
- bounded deterministic local-only validation before protected mutation;
- fail-closed unknown/unregistered/unsupported schema identity/version;
- no remote HTTP(S)/DNS/registry/filesystem schema fetching;
- no dynamic code/plugin/shell/schema-callback execution;
- no schema-derived authorization;
- owner/producer/tenant isolation and no tenant-specific V1 schema forks;
- separation from checkpoint, inbox and retry/quarantine state.

Validation success proves only schema conformance. It does not grant authorization or exactly-once/global-ordering semantics.

## Issue #267 post-activation continuity

Issue `#267` is the current continuity-only carrier.

It may reconcile exactly:

1. `AGENTS.md`
2. `README.md`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/handoffs/P04.07.md`
7. `docs/governance/P04_06_P04_07_TRANSITION_TRANSACTION.md`

It must not modify canonical `STATE.json`, the P04 package sequence, activation validators, `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`, runtime Go source, migrations, workflows or modules.

Throughout this continuity source and its unchanged promotion:

- P04.07 worker slots: `0`;
- P04.07 runtime tasks: `0`;
- P04.07 runtime branches: `0`;
- P04.07 migration reservations: `0`;
- schema-registry/provider choice: `none`.

## Separate implementation-plan gate

After Issue #267 continuity is accepted/read back, a later separately governed P04.07 implementation-plan/worker-plan must freeze the first bounded implementation slice before any runtime change.

That plan must explicitly decide:

- exact implementation paths and non-overlapping worker leases;
- supported bounded schema representation/dialect or subset;
- canonicalization/fingerprint algorithm and encoding;
- predecessor-set compatibility semantics;
- validator size/depth/time/memory/collection limits and stable safe errors;
- local-reference policy and graph limits if supported;
- registry persistence form;
- only if durable persistence is required, the exact next `kernel.events` migration owner/version/path/data budget after fresh migration preflight.

No worker may infer migration 4 or a vendor registry from package activation or continuity.

## Still unauthorized

- P04.07 runtime mutation on Issue #267 source/promotion;
- a P04.07 runtime worker/task/branch before separate implementation-plan acceptance;
- migration 4 reservation/schema mutation before separately governed persistence decision;
- provider/vendor registry selection or remote schema fetching;
- dynamic code/plugin/shell execution from schema content;
- P04.08 background-job ownership/runtime;
- P04.09 broad operator recovery UX;
- P04.10 replay/release/poison aggregate runtime;
- global ordering or end-to-end exactly-once claims;
- production business handlers/features;
- strategic X runtime;
- AI/model/agent product runtime.

## Exact next work

1. Complete the exact seven-file Issue `#267` continuity source from protected `main@5af9383c3e973c4055eea66d48462b7c9a2a5858`.
2. Verify exact changed-path scope and absence of runtime/state/worker-plan/migration drift.
3. Require exact-head Omnexa Governance.
4. Record honest SELF/Supervisor review and require zero unresolved threads.
5. Re-verify protected-main freshness.
6. Create an unchanged promotion at the exact reviewed source head.
7. Require fresh promotion Governance, promotion SELF/Supervisor review and zero unresolved threads.
8. Merge with expected-head protection and re-read protected main.
9. Confirm P04.01-P04.06 DONE, P04.07 sole ACTIVE at **6 / 10 done**, P04.08+ locked and zero runtime leases/reservations.
10. Only then create a separate P04.07 implementation-plan/worker-plan carrier.

Do not start P04.07 runtime or auto-advance P04.08 from this continuity carrier.
