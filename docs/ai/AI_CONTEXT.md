# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

The P03.01 closure/P03.02 activation candidate establishes:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 1 / 11 done.
- P03.01: DONE with canonical completion evidence.
- current work package after closure merge: P03.02 — Registry & Deterministic Discovery.
- P03.03-P03.11: PLANNED.
- `kernel_code_authorized=true` only for P03.02 after the closure PR merges.
- `business_feature_code_authorized=false`.

Until the closure PR itself passes canonical governance and merges, protected `main` remains authoritative at the prior cursor. This continuity file creates no implementation authority by itself.

## P03.01 canonical completion evidence

- implementation PR: #92
- exact implementation head: `87da3302605c852ae5bf43d473aaa01a9e1aaa74`
- implementation merge: `4229e2a28442bf475afed143bab359a770d48053`
- canonical run/job: `33009396644 / 98311433013` — PASS
- runner: GitHub-hosted Ubuntu 24.04.4 LTS / X64
- Go: 1.26.7
- completion evidence: `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`
- retained verifier: `scripts/verify_p03_01.sh`

Earlier run `33009067521 / 98310300730` remains diagnostic FAIL history only. It detected pinned-Go formatter drift; lint and vulnerability checks were clean, and the exact Go 1.26.7 formatter output was applied before the accepted exact-head run.

## Retained P01/P02 evidence

P01 and P02 remain completed prerequisites. Terminal P02.10 evidence remains PR #88, exact head `975e4925060a035780ca13b68c5437634ed0f4ea`, run/job `32904678957 / 97986011269` PASS, merge `88799aa41da8ce8c22540146d157d488565e2ce9`, evidence `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`.

All P01/P02/P03.01 regressions remain mandatory during P03.02.

## Active P03.02 contract after closure merge

Owner: `kernel.modules`.

P03.02 is limited to deterministic registry/discovery over validated P03.01 manifests:

- registry records for validated module identity/version/owner/source metadata;
- deterministic discovery from explicit approved repository/runtime sources;
- duplicate/conflicting module ID/version rejection;
- stable lookup/list ordering and stable source/version identity;
- structural separation of discovered/available metadata from installed/enabled lifecycle state;
- safe operator diagnostics suitable for later health integration.

P03.02 must not execute module/package code or lifecycle hooks. It must not perform arbitrary filesystem/network scanning, resolve dependency graphs, grant permissions/capability authority, or infer installed/enabled lifecycle state from discovery alone.

## Explicitly unauthorized

- P03.03 dependency resolver;
- P03.04 lifecycle runtime/persistence;
- P03.05-P03.11 later registries/trust/exit runtime;
- remote marketplace/catalog download, package signature/trust enforcement or sandbox runtime;
- System Graph, product federation or performance-intelligence runtime;
- P04 events/workflows;
- business modules/features;
- AI/model/agent runtime.

The XQ-100/XSG-100/XTRUST-100/XPF-200/XPERF-100 alignment remains planning-only and non-authorizing.

## Exact next action

After the P03.01 closure/P03.02 activation PR passes its own exact-head canonical governance and merges, verify protected `main` and `STATE.json` confirm P03.01 `done`, P03.02 sole `active`, `kernel_code_authorized=true` only for P03.02 and `business_feature_code_authorized=false`, then **STOP this closure session**.

A later governed session must branch from that exact protected-main SHA and implement only P03.02. Implementation and closure/state transition remain separate; never auto-advance to P03.03.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.02.md`, `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P03.02.md`.
