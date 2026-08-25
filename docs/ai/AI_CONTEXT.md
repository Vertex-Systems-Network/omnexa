# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

After this closure is accepted and merged, the intended canonical checkpoint is:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: DONE — 10 / 10; exit gate SATISFIED.
- P02.01-P02.10: DONE.
- current work package: NONE.
- P03: PLANNED — NOT ACTIVATED.
- `kernel_code_authorized=false`.
- `business_feature_code_authorized=false`.
- P03+, business modules and AI/model/agent runtime remain unauthorized until a separate governed preparation/readiness and activation transition.

Until this closure merges, protected `main` remains authoritative. This continuity file creates no implementation authority.

## P02.10 completion evidence

- implementation PR: #88
- final exact implementation head: `975e4925060a035780ca13b68c5437634ed0f4ea`
- canonical run/job: `32904678957 / 97986011269` — PASS
- implementation merge: `88799aa41da8ce8c22540146d157d488565e2ce9`
- runner: `GitHub Actions 1000023269`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- PostgreSQL: 18.6
- Valkey: 9.1.1
- S3 mock: 5.1.0
- repository Go quality, P01.01-P01.12, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.09 and P02.10 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`

Retained diagnostic failure remains visible rather than relabeled: `32903969206 / 97983773781` failed repository Go quality because the organization decorator referenced an undefined `invalidScopeFailure`. The defect was corrected with the existing owner-consistent `scopeDeniedFailure`; the final exact head passed the complete canonical lane.

## Retained P02 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

P02.04 retains the `kernel.identity` authentication/session foundation: approved adaptive password hashing, disclosure-safe authentication, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/revocation/replay rejection, password/account invalidation, safe device/session inventory, tenant/organization context re-authorization and classification-safe security lifecycle hooks. Authentication remains distinct from authorization.

P02.05 retains the `kernel.authorization` direct RBAC foundation: stable capability-oriented permission identifiers, deterministic Role composition, trusted tenant/organization scoped assignments, deny-by-default direct permission checks, privileged server-side mutations with anti-escalation, assignment revocation, role-name non-bypass and classification-safe required audit records.

P02.06 retains the contextual authorization layer: direct RBAC first; trusted relationship evidence exact to principal/object/tenant/organization scope; contextual constraints narrow only; internal/background origin never bypasses; field/export capability may be stricter than read; material denials/privileged decisions use safe audit; dependency failure fails closed.

P02.07 retains strong-authentication semantics: injected passkey verification, exact principal/session-bound expiring/replay-safe challenges, one-time digest-only recovery codes, session-bound non-authorizing step-up, explicit factor-removal invalidation and no restricted factor material in ordinary telemetry/audit.

P02.08 retains distinct non-human Service Account/API credentials: exact tenant/organization binding, one-time high-entropy issuance, verifier-only persistence, verify/rotate/revoke/expire lifecycle, supersession/revocation/expiry denial, direct RBAC composition and raw-secret non-disclosure. Credential possession never creates authority.

P02.09 retains tenant-scoped settings: trusted tenant/organization scope, deterministic organization→tenant→registered-definition-default precedence, no global/user override, current authorization on protected reads/mutations, owner-bounded persistence, generic RESTRICTED/secret values rejected, value-free security audit and cross-tenant/wrong-org fail-closed behavior. Settings never create authority.

P02.10 retains the accepted aggregate audit/exit behavior: secret-free identity/session/strong-auth/service-account lifecycle hooks bridge to `kernel.audit`; material tenancy/organization mutations use required-audit owner-preserving decorators; required-audit failure cannot silently claim success; same-tenant success and cross-tenant denial are proven; complete P02 migrations replay idempotently; and all P01/P02 regressions remain mandatory.

## Terminal P02 boundary

P02 completion creates no P03 authority. At this checkpoint:

- there is no active implementation package;
- kernel implementation is locked;
- business-feature implementation is locked;
- P03 is only planned;
- no generic audit UI/export, support impersonation product surface, module runtime, business domain, workflow/event, or AI/model/agent implementation is authorized by P02 completion.

## Exact next authorized action

After this P02 terminal closure merges, verify protected `main` and canonical `STATE.json`, confirm P02 is DONE 10 / 10 with exit SATISFIED, current work package NONE, P03 PLANNED and both implementation locks false, then **STOP the execution session**.

A later governed session may prepare P03 specifications/readiness and an explicit activation transition. It must not implement P03 before that activation is accepted.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.10.md`, `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.10.md`.
