# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 2 / 10 done.
- P02.01-P02.02: DONE.
- sole active package after this closure merges: `P02.03 — Organization hierarchy & scoped memberships`.
- owner: `kernel.organization`.
- `kernel_code_authorized=true` only for P02.03.
- `business_feature_code_authorized=false`.
- P02.04-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

## P02.02 completion evidence

- implementation PR: #71
- final exact implementation head: `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`
- canonical run/job: `32637760875 / 97189971101` — PASS
- implementation merge: `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`
- runner: `GitHub Actions 1000014047`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- repository Go quality, P01.01-P01.12, P02.01 regression and P02.02 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.02_COMPLETION_2026-08-23.md`

The exact implementation head passed its first canonical GitHub-hosted run; there is no P02.02 diagnostic FAIL to relabel.

## Retained P02.01-P02.02 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 established the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback. Client/request `tenant_id` remains only a selector, never authority.

## Active P02.03 boundary

P02.03 may implement only:

- Organization, Legal Entity, Business Unit, Branch, Team and Location hierarchy semantics under `kernel.organization`;
- tenant-bound parent/child validation and scoped organization memberships;
- deterministic hierarchy traversal/validation with cycle and cross-tenant rejection;
- organization/sub-scope context primitives for later policy evaluation without granting authority;
- classification-safe persistence and applicable fresh/upgrade migration evidence;
- focused deterministic positive/negative tests and a dedicated verifier.

P02.03 must not implement business Party/Person/Customer/Supplier models, authentication/session lifecycle, RBAC/policy enforcement beyond relationship primitives, MFA/passkeys, service-account/API credential lifecycle, tenant settings, HR/warehouse/CRM behavior, P03 module runtime, business UI/features, deployment authority or AI/model/agent runtime.

## Exact next authorized action

After the P02.02 closure/P02.03 activation PR merges, verify protected `main` and `STATE.json`, identify P02.03 as the sole authorized implementation scope, then **STOP the execution session**.

In the next governed session, create a fresh implementation branch from the exact protected-main SHA and implement only P02.03. Do not auto-advance to P02.04.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.03.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.03.md`.
