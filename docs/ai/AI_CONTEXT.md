# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

The activation candidate establishes:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P03: ACTIVE — 0 / 11 done.
- current work package: P03.01 — Module Manifest Schema.
- P03.02-P03.11: PLANNED.
- `kernel_code_authorized=true` only for P03.01.
- `business_feature_code_authorized=false`.

Until the activation PR merges, protected `main` remains authoritative. This continuity file creates no implementation authority by itself.

## P03 readiness evidence

- readiness PR: #90
- exact readiness head: `411c61f161a589b813ea28458cd43ce97be49f01`
- canonical run/job: `32910704583 / 98004146744` — PASS
- readiness merge: `712064870ce77d01902c2450c9eed240cda9c44e`
- runner: GitHub-hosted `ubuntu-24.04`
- artifacts: P03.01-P03.11 specs, P03 entry/exit gates, transition checklist, P03 validators, AI-native compatibility matrix

## Retained P02 evidence

P02 remains a completed prerequisite. Terminal P02.10 evidence remains PR #88, exact head `975e4925060a035780ca13b68c5437634ed0f4ea`, run/job `32904678957 / 97986011269` PASS, merge `88799aa41da8ce8c22540146d157d488565e2ce9`, evidence `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`.

Diagnostic run `32903969206 / 97983773781` remains FAIL history and is not acceptance evidence. All P01/P02 regressions remain mandatory.

## Active P03.01 contract

Owner: `kernel.modules`.

P03.01 is limited to the canonical machine-readable module manifest schema and deterministic validation contract:

- stable module ID/name/version/contract version/owner metadata;
- required platform and required/optional/platform dependency declarations;
- declared contribution/security metadata;
- minimal forward-compatible publisher/provenance/SBOM/scope metadata hooks without trust enforcement;
- bounded parsing and deterministic safe validation errors;
- positive/negative schema fixtures and explicit P03.01 verification.

P03.01 must treat manifests/package metadata as untrusted data and must not execute package code while parsing or validating. Declarations never grant authorization and manifests never contain secret values.

## Explicitly unauthorized

- P03.02 registry/discovery;
- P03.03 dependency resolver;
- P03.04 lifecycle runtime/persistence;
- P03.05-P03.11 later registries/trust/exit runtime;
- marketplace, System Graph, product federation or performance-intelligence runtime;
- P04 events/workflows;
- business modules/features;
- AI/model/agent runtime.

The XQ-100/XSG-100/XTRUST-100/XPF-200/XPERF-100 alignment remains planning-only and non-authorizing.

## Exact next action

After the activation PR merges, verify protected `main` and `STATE.json` confirm P03/P03.01 active, `kernel_code_authorized=true` only for P03.01 and `business_feature_code_authorized=false`, then **STOP**.

A later governed session must branch from that exact main SHA and implement only P03.01. Implementation and closure/state transition remain separate.

## Authority and references

Mandatory sources: `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P03.01.md`, `docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P03.01.md`.
