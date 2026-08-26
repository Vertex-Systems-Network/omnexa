# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

The P03.02 closure/P03.03 activation candidate establishes after its own successful merge:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 2 / 11 done.
- P03.01-P03.02: DONE with canonical completion evidence.
- current work package after closure merge: P03.03 — Dependency Graph Resolver.
- P03.04-P03.11: PLANNED.
- `kernel_code_authorized=true` only for P03.03 after the closure PR merges.
- `business_feature_code_authorized=false`.

Until the closure PR itself passes canonical governance and merges, protected `main` remains authoritative at P03.02. This continuity file creates no implementation authority by itself.

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

## Retained prerequisite evidence

P03.01 remains complete: PR #92, exact head `87da3302605c852ae5bf43d473aaa01a9e1aaa74`, run/job `33009396644 / 98311433013` PASS, merge `4229e2a28442bf475afed143bab359a770d48053`, evidence `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`.

P01 and P02 remain completed prerequisites. Terminal P02.10 evidence remains PR #88, exact head `975e4925060a035780ca13b68c5437634ed0f4ea`, run/job `32904678957 / 97986011269` PASS, merge `88799aa41da8ce8c22540146d157d488565e2ce9`, evidence `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`.

All P01/P02/P03.01/P03.02 regressions remain mandatory during P03.03.

## Active P03.03 contract after closure merge

Owner: `kernel.modules`.

P03.03 is limited to deterministic dependency graph validation/resolution over the validated P03.01 manifest and P03.02 registry/discovery contracts:

- version-aware required and optional dependency resolution;
- platform dependency validation;
- deterministic topological order for valid required graphs;
- cycle and incompatible-version detection;
- undeclared/forbidden/private dependency detection hooks;
- selective optional-dependency degradation metadata for later lifecycle/capability availability decisions.

P03.03 must fail closed on missing/incompatible required dependencies and cycles. Missing optional dependencies must not produce unrelated global failure. Resolver output cannot grant permissions, capabilities, tenant authority, private-schema access or database authority. Resolution must remain deterministic for identical registry/manifests.

## Explicitly unauthorized

- P03.04 lifecycle runtime/persistence;
- P03.05-P03.11 later registries/trust/exit runtime;
- package installation/download;
- P04 event dependency orchestration;
- full System Graph runtime or trust/advisory scanning;
- direct cross-module private imports/writes;
- business modules/features;
- AI/model/agent runtime.

The XQ-100/XSG-100/XTRUST-100/XPF-200/XPERF-100 alignment remains planning-only and non-authorizing.

## Exact next action

After the P03.02 closure/P03.03 activation PR passes its own exact-head canonical governance and merges, verify protected `main` and `STATE.json` confirm P03.01-P03.02 `done`, P03.03 sole `active`, `kernel_code_authorized=true` only for P03.03 and `business_feature_code_authorized=false`, then **STOP this closure session**.

A later governed session must branch from that exact protected-main SHA and implement only P03.03. Implementation and closure/state transition remain separate; never auto-advance to P03.04.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.03.md`, `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`, `docs/roadmap/evidence/P03.02_COMPLETION_2026-08-27.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P03.03.md`.
