# Omnexa Roadmap Status

Last reconciled: 2026-10-06 — **P04.08 activation candidate; protected main read-back still at pre-activation checkpoint**

## Canonical protected-main status

docs/roadmap/STATE.json on protected main ec726440213209af2018fea9ff552504f9a5496f remains authoritative until an activation carrier is accepted and read back:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit SATISFIED.
- P02: DONE — 10 / 10; exit SATISFIED.
- P03: DONE — 11 / 11; exit SATISFIED.
- P04: ACTIVE — 7 / 10 done (the seven accepted packages remain complete).
- Current protected-main package at base ec726440: NONE; the candidate STATE.json proposes current_work_package P04.08.
- Current work package: P04.08 (activation candidate; effective only after protected-main acceptance/read-back).
- P04.08-P04.10 remain planned/locked on the starting protected main.
- Kernel and business-feature authority on protected main: false.
- Migration 4, provider/vendor, queue/broker, business feature and AI product runtime: not authorized.

## Current activation candidate

A Supervisor-only governance/evidence transaction proposes P04.08 as the sole active package, with bounded kernel.jobs authority after acceptance. The current activation lease has zero worker slots and does not lease runtime paths. The proposed state is not effective until exact-head Governance, unchanged promotion, clean required review threads, protected merge and fresh main read-back.

- Fresh-main base: ec726440213209af2018fea9ff552504f9a5496f.
- P04.08 preparation source PR #316 / Governance 37312529538 PASS; unchanged promotion PR #317 / Governance 37315910306 PASS; merged main read-back ec726440213209af2018fea9ff552504f9a5496f.
- Historical PR #316 carried a CHANGES_REQUESTED review and unresolved threads at promotion time; those are disclosed for continuity and do not constitute approval of this activation candidate.
- Activation evidence record: docs/roadmap/evidence/P04.08_ACTIVATION_2026-10-06.md.
- Current activation PR/run/merge: pending.

Migration 4 remains unreserved and unauthorized. Provider/queue/broker selection, P04.09/P04.10, business features and AI/model/agent product runtime remain out of scope.

## P04.07 T05 readiness, ADR-0013 and T06 closure state

Fresh completion-readiness evidence is accepted:

- readiness trigger head `2129e0834dc5436a6c65759a0ca3643b01a38ada`;
- `P04.07 Readiness` run `35655774421` / #2 — PASS;
- job `106519447117` — PASS;
- `bash scripts/verify_p04_07.sh` — G0-G11 PASS;
- source evidence PR #299 / Governance #807 — PASS;
- promotion PR #300 / Governance #808 — PASS;
- accepted evidence main `3f34301f0cd37aa43392c26e906a0a719b6db2a0`.

The terminal-checkpoint decision was accepted through Issue #301 / ADR-0013 / #302/#303. Failed T06 source #304 / Governance #813 then exposed the omitted canonical-validator surface. Issue #305 amended that implementation boundary through source #306 / Governance #814, unchanged promotion #307 / Governance #815, and protected-main readback `3a314e10bcc06fd5d90222087a35e12c99e7fdff`.

Rebuilt T06 R2 source #308 / Governance #817 proved the canonical-state reconciliation at step 8, then failed step 11 in `scripts/validate_freeze_review.py`. The same stale active-P04 assumption was found in `scripts/validate_p03_preparation.py` and `scripts/validate_p03_package_specs.py`. Issue #309 consolidated those three downstream surfaces and was accepted through source #310 / Governance #818, unchanged promotion #311 / Governance #819, and protected-main readback `0ea2c5a2b34c00b32aaecf00f3ae0d217704b562`.

This fresh-main R3 T06 candidate closes P04.07 into an implementation-locked intra-P04 checkpoint while reconciling all five accepted validator surfaces:

- P04.01-P04.07 DONE;
- no active P04 work package;
- P04.08-P04.10 PLANNED / LOCKED;
- kernel/business implementation authority false;
- no migration 4/provider/remote/P04.08/business/AI runtime expansion.

The closure is canonical only when this exact R3 state passes source Governance, unchanged-promotion Governance, guarded merge and protected-main readback. Failed #304 and #308 remain immutable historical FAIL evidence and grant no merge authority.

## P04.07 accepted chain

Activation:

- source `#265`;
- unchanged promotion `#266`;
- protected-main merge/read-back `5af9383c3e973c4055eea66d48462b7c9a2a5858`.

Post-activation continuity:

- Issue `#267`;
- source `#268`;
- promotion `#269`;
- merge/read-back `3df6c0ec0d1034134b7417fe34e813a31b3ab821`.

Wave-1 plan:

- coordination Issue `#270`;
- source `#271`;
- promotion `#272`;
- merge/read-back `07ab591f37ccf16df74cbadd3cc641c195ebbc43`.

Accepted first-slice implementation:

- T01 schema registry: source `#273`, promotion `#274`, merge/read-back `5d61af89743625bd4e39d515a11cd711d1fe01e4`;
- T02 compatibility: source `#275`, promotion `#278`, merge/read-back `d06fa216ace09e806ef4dc1c5005cd26d424391e`;
- T03 payload validation: source `#276`, promotion `#279`, merge/read-back `e39a9dab525e410dad11f7059e41a1d052b42fd1`;
- T04 Supervisor verifier/evidence: source `#280`, promotion `#281`, exact head `52b3e5c9449f6f31c47cae9234347fbd0b6b8770`, source Governance `34540920140`, promotion Governance `34541808091`, merge/read-back `5c5153ee15c70646d28926743e2d19ab941013d2`.

The accepted first-slice verifier is `scripts/verify_p04_07.sh`. Historical evidence is `docs/roadmap/evidence/P04.07_FIRST_SLICE_2026-09-11.md`.

## Live lease state

Issue #267 and Issue #270 are historical accepted governance carriers. Their old branch/task/lease wording must not be interpreted as live authority.

Current live Wave-1 state:

- worker slots: `0`;
- worker tasks: `0`;
- worker branches: `0`;
- Supervisor T04 lease: `0`;
- migration reservations: `0`.

`docs/ai/ACTIVE_MULTI_AGENT_PLAN.json` records Wave 1 as completed historical leases. A new agent does not inherit those leases merely by reading the file or reusing a branch name.

## P04.07 retained boundary

Owner: `kernel.events`.

The accepted first slice remains bounded to provider-neutral local schema-registry, compatibility and payload-validation behavior. Retained laws include:

- payload schema version separate from P04.01 envelope version;
- stable owner/event-type/version identity;
- immutable accepted historical fingerprints;
- deterministic bounded canonicalization and compatibility;
- deterministic bounded payload validation before protected mutation;
- fail-closed unknown/unregistered/unsupported schema identity/version;
- no remote HTTP(S), DNS, hosted registry or arbitrary filesystem schema/reference fetching;
- no dynamic code, eval, arbitrary plugin, shell or schema-callback execution;
- no schema-derived authorization, capability or tenant-membership authority;
- no tenant-specific V1 schema forks;
- no migration 4;
- no provider/vendor schema registry selection;
- no P04.08+ runtime;
- no business-feature or AI/model/agent product-runtime authority;
- no global ordering or end-to-end exactly-once claim.

Validation success proves schema conformance only.

## Security continuity — Issue #285

Issue `#285` corrects a stale AI instruction/confused-deputy risk: repository-authoritative mirrors still presented pre-runtime Issue #267/#270 state as current after T01-T04 had already been accepted.

The reconciliation updates only governance/continuity surfaces. It does not change runtime code, migrations, provider choices, package state, architecture, business authority or AI product-runtime authority.

Required security outcome:

- historical plans and branches are explicitly historical;
- accepted T01-T04 facts are recorded;
- no historical write lease is reusable;
- contradictory live-looking instructions are removed;
- fail-closed authority behavior is preserved;
- canonical Governance must pass on the exact proposed head before merge.

## Next action

Complete the current P04.08 activation candidate through the source PR, exact-head Governance, resolved review threads, unchanged promotion, promotion Governance, protected merge and fresh main read-back. Only after acceptance, close/release this Supervisor-only governance lease and create a separate fresh-main P04.08 implementation lease. Until protected-main read-back, the canonical cursor remains P04 7/10 with no active package and all implementation authority locked. Migration 4 remains unreserved; no provider, queue or broker is selected; business and AI product runtime remain unauthorized.
