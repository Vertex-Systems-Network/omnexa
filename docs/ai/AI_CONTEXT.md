# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 1 / 10 done.
- P02.01: DONE — Principal & user identity foundation.
- sole active package after this closure merges: `P02.02 — Tenant lifecycle & trusted tenant context`.
- owner: `kernel.tenancy`.
- `kernel_code_authorized=true` only for P02.02.
- `business_feature_code_authorized=false`.
- P02.03-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

## P02.01 completion evidence

- implementation PR: #69
- final exact implementation head: `76919a9588f70aeea7e00f5214b82dcbf34cbee7`
- canonical run/job: `32635243643 / 97183883007` — PASS
- implementation merge: `44882e91e49d0364d841b511edbfd0619d05de1f`
- runner: `GitHub Actions 1000013607`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- repository Go quality, P01.01-P01.12 regressions, real PostgreSQL integration and P02.01 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md`

Initial canonical run `32635051321 / 97183427697` remains diagnostic FAIL for corrected `govet` shadow findings. No gate was weakened.

## Retained P02.01 contract

P02.01 established the owner-bounded `kernel.identity` human User foundation using UUIDv7, deterministic lifecycle semantics, classification-safe PII handling, immutable owner migration and safe structured failures. User remains distinct from business Person; identity attributes do not create authority; no tenant/authentication/session/RBAC/service-account authority was pulled forward.

## Active P02.02 boundary

P02.02 may implement only:

- authoritative Tenant identity/lifecycle/state under `kernel.tenancy`;
- explicit tenant-scoped persistence/query semantics;
- trusted tenant context derived from governed identity/execution relationships rather than request payload claims;
- minimum tenant relationship primitive required for later scoped authorization;
- same-tenant allow and cross-tenant deny evidence;
- applicable fresh/upgrade migration evidence;
- focused deterministic positive/negative tests and a dedicated verifier.

P02.02 must not implement organization hierarchy, authentication/session lifecycle, RBAC/policy, hidden support-superuser/global-tenant bypasses, tenant business data, P03 module runtime, business modules, deployment authority or AI/model/agent runtime.

## Exact next authorized action

After the P02.01 closure/P02.02 activation PR merges, verify protected `main` and `STATE.json`, identify P02.02 as the sole authorized implementation scope, then **STOP the execution session**.

In the next governed session, create a fresh implementation branch from the exact protected-main SHA and implement only P02.02. Do not auto-advance to P02.03.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.02.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.02.md`.
