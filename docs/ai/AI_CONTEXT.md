# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

After this closure is accepted and merged, the intended canonical checkpoint is:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 9 / 10 done.
- P02.01-P02.09: DONE.
- sole active package: `P02.10 — Identity / Permission Audit Trails & P02 Exit Proof`.
- owner: `kernel.audit` with P02 identity/tenancy/organization/authorization/configuration producers.
- `kernel_code_authorized=true` only for P02.10 after closure merge and post-merge verification.
- `business_feature_code_authorized=false`.
- P03+, business modules and AI/model/agent runtime remain unauthorized.

Until this closure merges, protected `main` remains authoritative. This continuity file creates no implementation authority.

## P02.09 completion evidence

- implementation PR: #86
- final exact implementation head: `0618904a18f82231469dd173aeb3d9d51edb73ed`
- canonical run/job: `32895186252 / 97956097639` — PASS
- implementation merge: `8ef86d2644b5ed455b3610192b8379d94204692f`
- runner: `GitHub Actions 1000022922`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- PostgreSQL: 18.6
- repository Go quality, P01.01-P01.12, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.08 and applicable P02.09 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.09_COMPLETION_2026-08-26.md`

Retained diagnostic failure remains visible rather than relabeled: `32894368734 / 97953417878` failed repository Go quality on gosec G304 for a variable-path integration migration helper. The fixed exact-path allow-list changed no runtime/schema/acceptance semantics, and the final canonical lane passed all historical regressions.

## Retained P02.01-P02.09 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

P02.04 retains the `kernel.identity` authentication/session foundation: approved adaptive password hashing, disclosure-safe authentication, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/revocation/replay rejection, password/account invalidation, safe device/session inventory, tenant/organization context re-authorization and classification-safe security lifecycle hooks. Authentication remains distinct from authorization.

P02.05 retains the `kernel.authorization` direct RBAC foundation: stable capability-oriented permission identifiers, deterministic Role composition, trusted tenant/organization scoped assignments, deny-by-default direct permission checks, privileged server-side mutations with anti-escalation, assignment revocation, role-name non-bypass and classification-safe required audit records.

P02.06 retains the `kernel.authorization` contextual layer on top of P02.05: direct RBAC is always evaluated first; trusted relationship evidence must match exact principal/object/tenant/organization scope; contextual constraints may only narrow authority; internal/background call origin never bypasses authorization; sensitive-field/export capability checks remain distinct from ordinary read; material denials/privileged decisions use safe required audit; relationship/constraint dependency failure fails closed.

P02.07 retains the `kernel.identity` strong-authentication foundation: passkey factor lifecycle is deterministic; passkey verification uses an injected verifier rather than custom protocol/private-key crypto; one-time challenges are exact User/session bound, expiring and replay-safe; recovery codes are one-time with digest-only persistence; step-up proof is session-bound and non-authorizing; factor removal follows explicit session invalidation policy; restricted factor/challenge/recovery material does not enter ordinary telemetry/audit.

P02.08 retains the distinct non-human Service Account/API credential foundation: exact tenant/organization binding; one-time high-entropy credential issuance; SHA-256 verifier-only persistence; verify/rotate/revoke/expire lifecycle; supersession/revocation/expiry denial; direct RBAC composition through current authorization state; and raw-secret non-disclosure. Service Accounts are never fake human Users and credential possession never creates authority.

P02.09 retains the tenant-scoped settings foundation: trusted tenant/organization setting scope; deterministic exact organization -> tenant -> registered-definition-default precedence; no global/user override path; protected-read and mutation authorization through current RBAC; owner-bounded configuration persistence; classification-aware values with generic RESTRICTED/secret values rejected; value-free required audit changes; cache invalidation after mutation; cross-tenant/wrong-org fail-closed behavior; and supported-upgrade migration evidence. Settings never create authority.

## Candidate active P02.10 boundary

Only after this closure merges and protected `main` + canonical `STATE.json` confirm P02.10 ACTIVE, P02.10 may implement only:

- attributable classification-safe audit records for material identity/session/tenant/org/role/policy/service-account/settings security operations;
- explicit required-audit fail-closed behavior for privileged operations where correctness requires audit;
- aggregate P02 verification composed from completed package verifiers rather than duplicating them;
- final cross-tenant isolation, object/scope permission, role-difference, service-account and session-invalidation evidence;
- fresh + supported-upgrade migration proof for the P02 schema baseline;
- repository Go quality, P01.01-P01.12 and P02.01-P02.09 regression preservation;
- applicable G0-G8 evidence on canonical GitHub-hosted `ubuntu-24.04`.

Audit remains separate from ordinary logs and stores no credentials/authentication factors/secrets. Audit write authority never implies audit read/export authority. Required-audit protected mutations must fail closed, and cross-tenant/privilege-escalation failures are release blockers.

P02.10 must not implement generic audit read/export/admin UI, support impersonation product behavior, P03 module runtime/permission registration, business domains/business-feature audit catalogs, P04 events/workflows, AI/model/agent behavior or automatic P03 activation.

## Exact next authorized action

After this P02.09 closure/P02.10 activation PR merges, verify protected `main` and canonical `STATE.json`, identify P02.10 as the sole active package, then **STOP the execution session**.

A later governed execution session must start from the then-current protected-main SHA, re-read canonical state, open PRs and the frozen P02.10 acceptance criteria, and only then may implement P02.10. Do not auto-advance to P03.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.10.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.10.md`.
