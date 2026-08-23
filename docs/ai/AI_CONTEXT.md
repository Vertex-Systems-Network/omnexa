# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

After this closure is accepted and merged, the intended canonical checkpoint is:

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 4 / 10 done.
- P02.01-P02.04: DONE.
- sole active package: `P02.05 — RBAC foundation`.
- owner: `kernel.authorization`.
- `kernel_code_authorized=true` only for P02.05 after closure merge and post-merge verification.
- `business_feature_code_authorized=false`.
- P02.06-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

Until this closure merges, protected `main` remains authoritative and still records P02.04 as active. This continuity file creates no implementation authority.

## P02.04 completion evidence

- implementation PR: #75
- final exact implementation head: `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f`
- canonical run/job: `32653747461 / 97229198036` — PASS
- implementation merge: `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`
- runner: `GitHub Actions 1000014748`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- PostgreSQL: 18.6
- repository Go quality, P01.01-P01.12, `omnexa verify all`, P02.01-P02.03 and P02.04 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`

Diagnostic failure remains retained rather than relabeled: `32653382151 / 97228307755` on exact head `5f8be7f6eb0ad9e807a9dc5d0d29121604822716` failed repository Go quality for four corrected `gosec` G101 findings and one corrected staticcheck simplification. It is not completion evidence.

## Retained P02.01-P02.04 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

P02.04 retains the `kernel.identity` authentication/session foundation: approved adaptive password hashing, disclosure-safe authentication, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/revocation/replay rejection, password/account invalidation, safe device/session inventory, tenant/organization context re-authorization and classification-safe security lifecycle hooks. Authentication remains distinct from authorization.

## Candidate active P02.05 boundary

Only after this closure merges and protected `main` + canonical `STATE.json` confirm P02.05 ACTIVE, P02.05 may implement only:

- stable permission identifiers and governed capability-oriented permission semantics;
- Role as deterministic permission composition, never a bypass identity;
- authorized tenant/organization-scoped role definitions and assignments;
- assignment/revocation lifecycle and deterministic direct-RBAC evaluation;
- privileged role/permission mutation checks and classification-safe audit hooks;
- allow, deny and privilege-escalation negative tests;
- applicable persistence/migrations and focused deterministic authorization evidence.

P02.05 must deny by default. Role names such as `admin`, `owner` or `superuser` never imply platform bypass. Permission scope may not cross tenant/organization boundaries by identifier substitution.

P02.05 must not implement relationship/context policy reserved for P02.06, MFA/passkeys P02.07, service-account/API credentials P02.08, tenant settings P02.09, P02.10 exit behavior, P03 module permission registration, business permissions/features, hidden super-admin behavior, deployment authority or AI/model/agent runtime.

## Exact next authorized action

After this P02.04 closure/P02.05 activation PR merges, verify protected `main` and canonical `STATE.json`, identify P02.05 as the sole active package, then **STOP the execution session**.

A later governed execution session must start from the then-current protected-main SHA, re-read canonical state, open PRs and the frozen P02.05 acceptance criteria, and only then may implement P02.05. Do not auto-advance to P02.06.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.05.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.05.md`.
