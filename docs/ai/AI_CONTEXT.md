# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 3 / 10 done.
- P02.01-P02.03: DONE.
- sole active package after this closure merges: `P02.04 — Authentication & session lifecycle`.
- owner: `kernel.identity`.
- `kernel_code_authorized=true` only for P02.04.
- `business_feature_code_authorized=false`.
- P02.05-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

## P02.03 completion evidence

- implementation PR: #73
- final exact implementation head: `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`
- canonical run/job: `32640790333 / 97197453122` — PASS
- implementation merge: `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`
- runner: `GitHub Actions 1000014421`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- repository Go quality, P01.01-P01.12, `omnexa verify all`, P02.01-P02.02 and P02.03 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.03_COMPLETION_2026-08-23.md`

Diagnostic failures remain retained rather than relabeled: `32640199607 / 97196005995` failed nine corrected `govet` shadow findings; `32640419476 / 97196545810` failed the P02.03 dependency guard because its narrow allowlist omitted already-governed transitive `kernel.database`/`kernel.config` prerequisites.

## Retained P02.01-P02.03 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

## Active P02.04 boundary

P02.04 may implement only:

- authentication mechanism boundaries and approved adaptive password hashing where passwords are supported;
- session/access/refresh credential expiry, rotation and revocation;
- short-lived access credentials relative to refresh/session credentials;
- device/session inventory semantics where supported;
- required invalidation after material account/security changes;
- tenant/organization context re-authorization rather than stale-client authority;
- disclosure-safe authentication failures and synthetic security fixtures;
- classification-safe audit hooks without secret payloads;
- applicable migrations and focused deterministic security/lifecycle evidence.

P02.04 must not implement RBAC/policy decisions, MFA/passkeys, service-account/API credentials, SAML/SCIM/enterprise SSO, business login portals/UI, P03 runtime, business features, deployment authority or AI/model/agent runtime.

## Exact next authorized action

After this P02.03 closure/P02.04 activation PR merges, verify protected `main` and canonical `STATE.json`, identify P02.04 as the sole authorized implementation scope, then **STOP the execution session**.

In the next governed session, create a fresh implementation branch from the exact protected-main SHA and implement only P02.04 against its already-frozen acceptance criteria. Do not auto-advance to P02.05.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.04.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.04.md`.
