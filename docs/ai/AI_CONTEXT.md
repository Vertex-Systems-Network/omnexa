# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

After this closure is accepted and merged, the intended canonical checkpoint is:

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 6 / 10 done.
- P02.01-P02.06: DONE.
- sole active package: `P02.07 — MFA/passkey-ready flows`.
- owner: `kernel.identity`.
- `kernel_code_authorized=true` only for P02.07 after closure merge and post-merge verification.
- `business_feature_code_authorized=false`.
- P02.08-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

Until this closure merges, protected `main` remains authoritative and records the accepted P02.06 implementation while P02.06 is still canonically active. This continuity file creates no implementation authority.

## P02.06 completion evidence

- implementation PR: #79
- final exact implementation head: `dbbd105fd5f2543ca7dd5df93375eaf1057928fc`
- canonical run/job: `32664834112 / 97256520050` — PASS
- implementation merge: `083c2866f0cd0773b85201750c2196bfd2fcc167`
- runner: `GitHub Actions 1000015801`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- PostgreSQL: 18.6
- repository Go quality, P01.01-P01.12, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.05 and applicable P02.06 G0-G8: PASS
- P02.06 G4: N/A for new persistence; retained P02.05 migration regression PASS
- evidence: `docs/roadmap/evidence/P02.06_COMPLETION_2026-08-24.md`

Retained diagnostic failure remains visible rather than relabeled: `32664671013 / 97256120056` on exact head `68a295957e30453afdfa8e4303ad8924b08dc530` failed repository Go quality for one formatter alignment-space difference in `contextual_errors.go`. The correction was formatter-only; lint remained 0 issues, govulncheck found no vulnerabilities, and no behavior, acceptance criterion or gate was changed.

## Retained P02.01-P02.06 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

P02.04 retains the `kernel.identity` authentication/session foundation: approved adaptive password hashing, disclosure-safe authentication, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/revocation/replay rejection, password/account invalidation, safe device/session inventory, tenant/organization context re-authorization and classification-safe security lifecycle hooks. Authentication remains distinct from authorization.

P02.05 retains the `kernel.authorization` direct RBAC foundation: stable capability-oriented permission identifiers, deterministic Role composition, trusted tenant/organization scoped assignments, deny-by-default direct permission checks, privileged server-side mutations with anti-escalation, assignment revocation, role-name non-bypass and classification-safe required audit records.

P02.06 retains the `kernel.authorization` contextual layer on top of P02.05: direct RBAC is always evaluated first; trusted relationship evidence must match exact principal/object/tenant/organization scope; contextual constraints may only narrow authority; internal/background call origin never bypasses authorization; sensitive-field/export capability checks remain distinct from ordinary read; material denials/privileged decisions use safe required audit; relationship/constraint dependency failure fails closed. P02.06 did not take ownership of business-object persistence.

## Candidate active P02.07 boundary

Only after this closure merges and protected `main` + canonical `STATE.json` confirm P02.07 ACTIVE, P02.07 may implement only:

- MFA factor enrollment, verification and removal lifecycle;
- passkey/WebAuthn-ready credential and challenge contracts consistent with approved platform security primitives;
- approved additional factor semantics where implemented;
- recovery-code lifecycle with one-way/secure handling;
- strong-auth/step-up policy hooks for privileged operations without replacing authorization;
- replay, expiry and principal/session binding validation for authentication challenges;
- synthetic fixtures and no-secret telemetry/audit behavior;
- applicable owner-bounded persistence, migrations and deterministic identity/security evidence.

P02.07 must treat factor secrets, recovery codes and authentication-equivalent material as `RESTRICTED`. No secret value may enter ordinary logs, traces, errors or audit payloads. Expired, replayed or wrong-principal/session challenges fail closed. Factor removal or security-policy change may trigger session invalidation according to policy. Strong authentication never replaces authorization.

P02.07 must not implement P02.08 service accounts/API credentials, P02.09 tenant settings, P02.10 exit-product behavior, P24 enterprise SSO/SAML/SCIM, business portal UI/features, custom cryptographic algorithms/private-key management outside approved platform primitives, deployment authority or AI/model/agent runtime.

## Exact next authorized action

After this P02.06 closure/P02.07 activation PR merges, verify protected `main` and canonical `STATE.json`, identify P02.07 as the sole active package, then **STOP the execution session**.

A later governed execution session must start from the then-current protected-main SHA, re-read canonical state, open PRs and the frozen P02.07 acceptance criteria, and only then may implement P02.07. Do not auto-advance to P02.08.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.07.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.07.md`.
