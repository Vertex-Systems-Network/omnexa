# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

After this closure is accepted and merged, the intended canonical checkpoint is:

- Foundation Architecture v1: FROZEN.
- P00: DONE — 10 / 10.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 8 / 10 done.
- P02.01-P02.08: DONE.
- sole active package: `P02.09 — Tenant-Scoped Settings`.
- owner: `kernel.configuration`.
- `kernel_code_authorized=true` only for P02.09 after closure merge and post-merge verification.
- `business_feature_code_authorized=false`.
- P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

Until this closure merges, protected `main` remains authoritative. This continuity file creates no implementation authority.

## P02.08 completion evidence

- implementation PR: #84
- final exact implementation head: `43bdcf525ce5e0cfdb9dc0707fbafee7cd552543`
- canonical run/job: `32885950897 / 97926598423` — PASS
- implementation merge: `32eb7187eb229327585551e4e28b0d596de78bd9`
- runner: `GitHub Actions 1000022204`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260823.283.1`
- Go: 1.26.7
- PostgreSQL: 18.6
- repository Go quality, P01.01-P01.12, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.07 and applicable P02.08 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.08_COMPLETION_2026-08-25.md`

Retained diagnostic failures remain visible rather than relabeled: `32882746486 / 97915911717`, `32884311341 / 97921359088`, and `32884939579 / 97923224921`. Run `32885758158` was superseded/cancelled. The final canonical lane proves historical P02.04-P02.07 regressions remain PASS after the narrow verifier compatibility corrections.

## Retained P02.01-P02.08 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

P02.04 retains the `kernel.identity` authentication/session foundation: approved adaptive password hashing, disclosure-safe authentication, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/revocation/replay rejection, password/account invalidation, safe device/session inventory, tenant/organization context re-authorization and classification-safe security lifecycle hooks. Authentication remains distinct from authorization.

P02.05 retains the `kernel.authorization` direct RBAC foundation: stable capability-oriented permission identifiers, deterministic Role composition, trusted tenant/organization scoped assignments, deny-by-default direct permission checks, privileged server-side mutations with anti-escalation, assignment revocation, role-name non-bypass and classification-safe required audit records.

P02.06 retains the `kernel.authorization` contextual layer on top of P02.05: direct RBAC is always evaluated first; trusted relationship evidence must match exact principal/object/tenant/organization scope; contextual constraints may only narrow authority; internal/background call origin never bypasses authorization; sensitive-field/export capability checks remain distinct from ordinary read; material denials/privileged decisions use safe required audit; relationship/constraint dependency failure fails closed.

P02.07 retains the `kernel.identity` strong-authentication foundation: passkey factor lifecycle is deterministic; passkey verification uses an injected verifier rather than custom protocol/private-key crypto; one-time challenges are exact User/session bound, expiring and replay-safe; recovery codes are one-time with digest-only persistence; step-up proof is session-bound and non-authorizing; factor removal follows explicit session invalidation policy; restricted factor/challenge/recovery material does not enter ordinary telemetry/audit.

P02.08 retains the distinct non-human Service Account/API credential foundation: exact tenant/organization binding; one-time high-entropy credential issuance; SHA-256 verifier-only persistence; verify/rotate/revoke/expire lifecycle; supersession/revocation/expiry denial; direct RBAC composition through current authorization state; and raw-secret non-disclosure. Service Accounts are never fake human Users and credential possession never creates authority.

## Candidate active P02.09 boundary

Only after this closure merges and protected `main` + canonical `STATE.json` confirm P02.09 ACTIVE, P02.09 may implement only:

- tenant-scoped and approved organization-scoped setting resolution using `kernel.configuration`;
- trusted scope derived from P02 identity/tenancy context rather than arbitrary payload identifiers;
- authorization around protected setting reads/writes;
- classification-aware values and no-secret output behavior;
- deterministic precedence only for explicitly supported scopes;
- change audit hooks for security-significant settings;
- same-tenant allow plus cross-tenant/wrong-scope negative tests;
- applicable owner-bounded persistence, migrations and focused configuration/security evidence.

Settings/feature flags cannot create authority. Tenant/org scope is trusted context rather than a client assertion. `kernel.configuration` remains authoritative owner; no cross-tenant fallback or global-write shortcut is authorized.

P02.09 must not implement business-module settings, P03 module runtime, a secrets-management product surface, feature/config values that independently grant authority, deployment/environment orchestration, P02.10 phase-exit implementation or other future business/AI scope.

## Exact next authorized action

After this P02.08 closure/P02.09 activation PR merges, verify protected `main` and canonical `STATE.json`, identify P02.09 as the sole active package, then **STOP the execution session**.

A later governed execution session must start from the then-current protected-main SHA, re-read canonical state, open PRs and the frozen P02.09 acceptance criteria, and only then may implement P02.09. Do not auto-advance to P02.10.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.09.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.09.md`.
