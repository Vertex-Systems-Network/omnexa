# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

Terminal P03 closure is tracked by GitHub #135 / Linear ABD-208 from protected-main base `b3b9b61f963df6a05ea45cbd3c562e12974d92d0`.

Target terminal state after closure merge:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit SATISFIED.
- P02: DONE — 10 / 10; exit SATISFIED.
- P03: DONE — 11 / 11; exit SATISFIED.
- P03.01-P03.11: DONE with canonical completion evidence.
- current work package: NONE.
- `kernel_code_authorized=false`.
- `business_feature_code_authorized=false`.
- P04+: PLANNED / LOCKED.

P03 completion never auto-activates P04. P04 readiness/preparation and activation are a later separate governed flow.

## P03.11 canonical implementation evidence

- implementation issue #132 — completed;
- draft implementation carrier #133 — closed unmerged;
- exact implementation head `a083a8a86ec3a51309fa479ee49c79e1b6ec9f10`;
- draft Governance #511 / `33258092323 / 99115191521` — PASS;
- promotion implementation PR #134 — merged;
- promotion Governance #512 / `33258456851 / 99116152701` — PASS;
- implementation merge / terminal closure base `b3b9b61f963df6a05ea45cbd3c562e12974d92d0`;
- completion evidence `docs/roadmap/evidence/P03.11_COMPLETION_2026-08-29.md`;
- retained verifier `scripts/verify_p03_11.sh`.

Promotion #512 passed canonical state/spec validators, repository Go quality, P01.01-P01.12, P02.01-P02.10 and P03.01-P03.11.

## P03.11 delivered boundary

P03.11 adds no trust/certification authority. It preserves already-validated v1/v2 publisher, provenance, SBOM, data classification and security declarations in immutable registry-bound snapshots and projects deterministic typed/versioned metadata-only profiles. Symbolic secret names may be exposed; secret locators/values may not. Untrusted package code is not executed for profile discovery.

The P03 exit aggregate maps EX-01..EX-07 to retained executable tests and adds explicit unrelated-module lifecycle isolation across install, failure/recovery, enable, suspend/resume, disable, archive/restore, detach and purge.

## Retained P03 prerequisites

P03.01-P03.10 retain their existing immutable completion files and verifier chain. P01.01-P01.12 and P02.01-P02.10 remain historical prerequisites. Diagnostic failures remain diagnostic and are never rewritten as PASS.

## Retained architecture/security baselines

ADR-0012 dependency-version semantics remain accepted. `kernel.configuration`, `kernel.authorization`, P01 migration authority and P01 health/readiness authority remain authoritative at their existing boundaries. Capability/permission/UI/health/trust metadata remains non-granting. Tenant isolation, deny-by-default authorization, attributable audit and data classification remain mandatory.

## Explicitly unauthorized at terminal P03

- P04 events/jobs runtime;
- business modules/features;
- publisher onboarding/signature trust roots or package certification;
- dependency advisory/license enforcement;
- marketplace/package acquisition/distribution runtime;
- Product Federation/System Graph/Performance Intelligence runtime;
- generic service-mesh/workflow expansion;
- strategic X-program runtime;
- AI/model/agent runtime;
- weakening governance/security/regression gates.

## Exact next action

Complete terminal P03 closure branch `chore/p03-11-closure-p03-exit` under GitHub #135 / Linear ABD-208, obtain a fresh exact-final-head GitHub-hosted Governance PASS, merge only after current-with-main/review/conversation preflight with an expected-head guard, then re-read protected main and stop at P03 DONE / P04 PLANNED.
