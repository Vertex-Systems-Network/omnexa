# Omnexa Roadmap Status

Last reconciled: 2026-09-23 — **P04.07 T06 R3 CLOSURE CANDIDATE / ADR-0013 TERMINAL CHECKPOINT**

## Canonical status

`docs/roadmap/STATE.json` remains the machine-readable source of truth.

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit SATISFIED.
- P02: DONE — 10 / 10; exit SATISFIED.
- P03: DONE — 11 / 11; exit SATISFIED.
- P04: ACTIVE — **7 / 10 done**.
- Current work package: **NONE** — implementation-locked intra-P04 terminal checkpoint after P04.07.
- P04.01-P04.07: DONE with retained accepted evidence.
- P04.07: DONE in this closure state; canonical only from protected-main acceptance of this exact tree.
- P04.08-P04.10: PLANNED / LOCKED.
- `kernel_code_authorized=false` at the terminal checkpoint.
- `business_feature_code_authorized=false`.
- migration 4: not reserved or authorized.
- provider/vendor schema registry: not selected or authorized.
- strategic X runtime: unauthorized.
- AI/model/agent product runtime: unauthorized.

P04.07 completion is based on its accepted first runtime slice plus fresh T05 G0-G11 readiness evidence and ADR-0013 terminal-checkpoint governance. P04.08 is not activated by this closure.

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

No package auto-advances from this terminal checkpoint.

A future P04.08 preparation/activation requires a fresh separately governed transaction from current protected main. Until then, P04.08-P04.10 remain PLANNED/LOCKED, kernel implementation authority is false, migration 4 remains unreserved, and provider/business/AI runtime scope remains unauthorized.
