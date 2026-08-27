# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Protected main `77ca52b4041013d1785b00aac6655aad7f3fe91f` establishes:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 2 / 11 done.
- P03.01-P03.02: DONE with canonical completion evidence.
- current work package: P03.03 — Dependency Graph Resolver.
- P03.04-P03.11: PLANNED / LOCKED.
- `kernel_code_authorized=true` only for P03.03.
- `business_feature_code_authorized=false`.

The prior P03.02 closure/P03.03 activation candidate has already merged through PR #95. Any continuity wording that says protected main is still at P03.02 is stale and superseded.

## P03.01 canonical completion evidence

- implementation PR: #92
- exact implementation head: `87da3302605c852ae5bf43d473aaa01a9e1aaa74`
- implementation merge: `4229e2a28442bf475afed143bab359a770d48053`
- canonical run/job: `33009396644 / 98311433013` — PASS
- completion evidence: `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`
- retained verifier: `scripts/verify_p03_01.sh`

## P03.02 canonical completion evidence

- implementation PR: #94
- exact implementation head: `0c46db41b0d724a08ea1a78545b3c2debdd8cd05`
- implementation merge: `2e38969dbbbcfcf4765a114f449dc3fa960061d7`
- canonical run/job: `33022405704 / 98355747775` — PASS
- runner: GitHub-hosted `ubuntu-24.04`
- Go: repository-pinned `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`
- retained verifier: `scripts/verify_p03_02.sh`

The accepted governance job passed every retained P01/P02/P03.01 regression plus the dedicated P03.02 verifier. An earlier candidate exposed a wrapped-error `errorlint` defect; it was corrected using `errors.As` without gate or scope weakening. Only the final exact successful head/run above is acceptance evidence.

P01 and P02 remain completed prerequisites. Terminal P02.10 evidence remains PR #88, exact head `975e4925060a035780ca13b68c5437634ed0f4ea`, run/job `32904678957 / 97986011269` PASS, merge `88799aa41da8ce8c22540146d157d488565e2ce9`, evidence `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`.

All P01/P02/P03.01/P03.02 regressions remain mandatory during P03.03.

## ADR-0012 accepting/reconciliation checkpoint

The active P03.03 package exposed a Class C prerequisite: schema v1 cannot express per-module dependency version constraints and the public P03.02 registry does not carry dependency declarations.

ADR-0012 accepts the forward architecture for this prerequisite. Its accepted status becomes authoritative only after the accepting/reconciliation PR passes exact-head canonical governance and merges to protected `main`.

The accepted baseline is:

- explicit bounded top-level schema-version dispatch;
- preserved strict schema-v1 parser/validator behavior;
- separate strict schema-v2 parser/validator;
- schema-v2 required/optional module dependencies as exact `{id, constraint}` records;
- 1–16 comparator / 256-byte bounded strict SemVer grammar with `=`, `>`, `>=`, `<`, `<=` and logical AND only;
- no parser fallback, string overloading or implicit compatibility inference;
- one discovered version per module ID under retained P03.02 public registry semantics;
- discovery atomically binds each registry identity to its exact normalized validated manifest snapshot;
- resolver dependency declarations come only from that bound snapshot rather than a second raw-manifest set;
- required edges alone drive deterministic install/enable topological order and release-blocking cycles;
- optional absence/incompatibility yields selective degradation and optional edges do not participate in the required global cycle gate;
- no multi-version/SAT solving, automatic dependency selection, external compatibility matrix or remote package acquisition;
- resolver output creates no permission/capability/tenant/database/private-access authority.

Historical P03.01/P03.02 evidence remains immutable. ADR reconciliation is architecture prerequisite evidence only; it is not P03.03 implementation completion evidence.

## Active P03.03 implementation contract

Owner: `kernel.modules`.

After the ADR-0012 reconciliation merges and protected `main` is re-read, P03.03 implementation must start on a **new separate branch from the exact new main SHA**.

Authorized implementation is limited to:

- bounded v1/v2 schema dispatch with no fallback;
- schema-v2 structured dependency validation;
- strict SemVer parsing/comparison and bounded constraints;
- package-private registry-to-normalized-manifest binding;
- version-aware required and optional dependency resolution;
- platform dependency validation;
- deterministic required graph topological ordering;
- self-dependency, required-cycle and incompatible-version rejection;
- undeclared/forbidden/private dependency detection hooks;
- schema-v1 required/optional migration/degradation behavior;
- selective optional-dependency degradation metadata;
- stable safe deterministic diagnostics;
- dedicated `scripts/verify_p03_03.sh` and canonical GitHub-hosted governance wiring.

## Explicitly unauthorized

- implementation code on the ADR accepting/reconciliation branch;
- P03.04 lifecycle runtime/persistence;
- P03.05-P03.11 later registries/trust/exit runtime;
- package installation/download or remote marketplace runtime;
- multi-version/SAT solver or automatic dependency selection;
- P04 event dependency orchestration;
- full System Graph runtime or trust/advisory scanning;
- direct cross-module private imports/writes;
- business modules/features;
- AI/model/agent runtime.

The XQ-100/XSG-100/XTRUST-100/XPF-200/XPERF-100 alignment remains planning-only and non-authorizing.

## Exact next action

1. Complete the ADR-0012 accepting/reconciliation PR without implementation code.
2. Require exact-final-head canonical GitHub-hosted governance plus all retained P01/P02/P03.01/P03.02 regressions.
3. Merge only if current with protected `main` and repository gates permit it.
4. Re-read protected `main`, `STATE.json`, ADR-0012 and P03.03 docs after merge and record the exact new SHA.
5. Create a new separate P03.03 implementation branch from that SHA and implement only the accepted scope.
6. Keep P03.04+ locked; P03.03 completion/state transition remains a separate later closure PR.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.03.md`, `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`, `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`, accepted `docs/adr/ADR-0012-versioned-module-dependency-requirements.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P03.03.md`.
