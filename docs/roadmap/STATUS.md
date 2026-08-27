# Omnexa Program Status

Last reconciled: **2026-08-27**

## Current position

- Program: **Kernel Program**
- Phase: **P03 — Module Runtime**
- Phase state: **ACTIVE**
- Current work package: **P03.03 — Dependency Graph Resolver**
- P03 progress: **2 / 11 done**
- P03.01: **DONE**
- P03.02: **DONE**
- P03.03: **ACTIVE**
- P03.04-P03.11: **PLANNED / LOCKED**
- P03 entry gate: **SATISFIED**
- P03 exit gate: **NOT SATISFIED**
- P02: **DONE — 10 / 10**
- P02 exit gate: **SATISFIED**
- P01: **DONE — 12 / 12**
- P01 exit gate: **SATISFIED**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation authority: **P03.03 ONLY**
- Business-feature implementation: **NOT AUTHORIZED**

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This status file is explanatory and must not grant authority beyond that state.

## Protected-main activation identity

P03.02 closure / P03.03 activation merged through PR #95 as protected-main commit:

- activation merge: `77ca52b4041013d1785b00aac6655aad7f3fe91f`
- P03.01-P03.02: `done`
- P03.03: sole `active` package
- `kernel_code_authorized=true` bounded to P03.03
- `business_feature_code_authorized=false`

The earlier closure-candidate wording is superseded by this merged protected-main state.

## P03.01 retained completion boundary

P03.01 — Module Manifest Schema remains complete with canonical evidence:

- implementation PR: #92
- exact implementation head: `87da3302605c852ae5bf43d473aaa01a9e1aaa74`
- canonical run/job: `33009396644 / 98311433013` — PASS
- implementation merge: `4229e2a28442bf475afed143bab359a770d48053`
- completion evidence: `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`
- retained verifier: `scripts/verify_p03_01.sh`

P03.01 remains `done`. Accepted ADR-0012 adds forward schema-v2 compatibility work under active P03.03 without rewriting historical P03.01 evidence.

## P03.02 retained completion boundary

P03.02 — Registry & Deterministic Discovery remains complete with canonical evidence:

- implementation PR: #94
- exact implementation head: `0c46db41b0d724a08ea1a78545b3c2debdd8cd05`
- canonical run/job: `33022405704 / 98355747775` — PASS
- implementation merge: `2e38969dbbbcfcf4765a114f449dc3fa960061d7`
- completion evidence: `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`
- retained verifier: `scripts/verify_p03_02.sh`

The accepted P03.02 public boundary remains deterministic, explicit-source, non-executing and fail-closed for duplicate/conflicting identities. Registry metadata remains distinct from installed/enabled lifecycle state and creates no authorization.

Accepted ADR-0012 permits only additive package-private normalized validated-manifest binding needed by P03.03; public P03.02 lookup/list/cardinality and duplicate/version-conflict behavior remains unchanged.

## Accepted ADR-0012 architecture prerequisite

The ADR-0012 accepting/reconciliation PR establishes the governing version-aware dependency contract for P03.03 and becomes authoritative only when that PR merges to protected `main`.

Accepted architecture:

- manifest schema v1 remains parseable under retained P03.01 semantics;
- explicit bounded schema-version dispatch selects separate strict v1/v2 decoders;
- schema v2 uses exact `{id, constraint}` required/optional module dependency records;
- strict bounded SemVer comparator grammar applies;
- self-dependencies, duplicates and cross-class conflicts fail deterministically;
- discovery atomically binds each registry identity to its exact normalized validated manifest snapshot;
- resolver policy may consume only that bound snapshot, not an independently reparsed raw-manifest set;
- one discovered module ID maps to at most one discovered version;
- required edges alone drive install/enable ordering and release-blocking cycle detection;
- optional absence/incompatibility produces selective degradation and does not globally invalidate unrelated modules;
- resolver output remains metadata/eligibility only and grants no authority;
- multi-version/SAT solving, external compatibility matrix and package acquisition remain out of scope.

## Active P03.03 boundary

P03.03 is limited to deterministic dependency graph validation/resolution over the retained P03.01/P03.02 contracts plus accepted ADR-0012 forward evolution.

Authorized P03.03 implementation surfaces after the ADR reconciliation merge include:

- bounded v1/v2 manifest version dispatch;
- strict schema-v2 dependency records and constraints;
- package-private registry-bound normalized manifest snapshots;
- version-aware required/optional dependency resolution;
- platform dependency validation;
- deterministic required-edge topological ordering;
- required-cycle/incompatible-version/self-dependency rejection;
- undeclared/forbidden/private dependency detection hooks;
- explicit schema-v1 migration/degradation behavior;
- selective optional-dependency degradation metadata;
- dedicated `scripts/verify_p03_03.sh` and canonical governance wiring.

Resolver output cannot grant permission/capability/tenant/database authority. P03.04 lifecycle runtime, P03.05-P03.11 later runtime, package install/download, multi-version solving, P04+, business features, full System Graph/trust-advisory runtime and AI/model/agent runtime remain unauthorized.

## Retained prerequisite chain

P01.01-P01.12 and P02.01-P02.10 remain complete with canonical package-specific evidence. Terminal P02.10 evidence remains:

- implementation PR: #88
- exact head: `975e4925060a035780ca13b68c5437634ed0f4ea`
- canonical run/job: `32904678957 / 97986011269` — PASS
- implementation merge: `88799aa41da8ce8c22540146d157d488565e2ce9`
- canonical runner: GitHub-hosted `ubuntu-24.04`
- completion evidence: `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`

Diagnostic run `32903969206 / 97983773781` remains explicit **FAIL** history and is not acceptance evidence.

All completed P01/P02/P03.01/P03.02 regression verifiers remain mandatory during P03.03.

## Protected integration / CI

`main` remains the protected PR-only integration authority. Canonical governance evidence is GitHub-hosted `ubuntu-24.04` only; local/self-hosted governance evidence is prohibited.

The ADR-0012 accepting/reconciliation PR must be current with protected `main` and pass its own exact-head canonical governance, repository Go quality, retained P01.01-P01.12 and P02.01-P02.10 regressions, completed P03.01/P03.02 regression verification and P03 package/state validators before merge.

ADR acceptance/reconciliation evidence is **not** P03.03 implementation completion evidence.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. It does not broaden P03 implementation authority.

## Exact next work

1. Complete and merge the ADR-0012 accepting/reconciliation PR only after exact-head governance passes and current repository merge gates permit it.
2. Re-read protected `main` and `STATE.json` after that merge and record the exact new SHA.
3. Create a separate P03.03 implementation branch from that exact main SHA.
4. Implement only the accepted P03.03 contract and obtain exact-head canonical evidence.
5. Keep P03.04+ locked; package completion/state transition belongs to a later separate closure PR.
