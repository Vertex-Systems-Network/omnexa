# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

After this closure is accepted and merged, the intended canonical checkpoint is:

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 7 / 10 done.
- P02.01-P02.07: DONE.
- sole active package: `P02.08 — Service accounts & API credentials`.
- owner: `kernel.identity`.
- `kernel_code_authorized=true` only for P02.08 after closure merge and post-merge verification.
- `business_feature_code_authorized=false`.
- P02.09-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

Until this closure merges, protected `main` remains authoritative and records the accepted P02.07 implementation while P02.07 is still canonically active. This continuity file creates no implementation authority.

## P02.07 completion evidence

- implementation PR: #81
- final exact implementation head: `51ccaa12c3534f74fba6eab9d4698ee483ef4ffd`
- canonical run/job: `32669167972 / 97267175953` — PASS
- implementation merge: `5642f5da1eb24e70b67e5ec757d9f4584c4e3f5c`
- runner: `GitHub Actions 1000016379`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- PostgreSQL: 18.6
- repository Go quality, P01.01-P01.12, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.06 and applicable P02.07 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.07_COMPLETION_2026-08-24.md`

Retained diagnostic failure remains visible rather than relabeled: `32668735841 / 97266149110` on exact head `0ab9873367586ab0191a91f658da275f449de796` failed the P02.07 verifier because a redundant regex matched explanatory `authorization.` text in a comment. The final correction removed only that comment-sensitive token while independent authorization-import and executable authority-symbol guards remained. No runtime, schema, test, acceptance, security or gate behavior changed.

## Retained P02.01-P02.07 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

P02.04 retains the `kernel.identity` authentication/session foundation: approved adaptive password hashing, disclosure-safe authentication, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/revocation/replay rejection, password/account invalidation, safe device/session inventory, tenant/organization context re-authorization and classification-safe security lifecycle hooks. Authentication remains distinct from authorization.

P02.05 retains the `kernel.authorization` direct RBAC foundation: stable capability-oriented permission identifiers, deterministic Role composition, trusted tenant/organization scoped assignments, deny-by-default direct permission checks, privileged server-side mutations with anti-escalation, assignment revocation, role-name non-bypass and classification-safe required audit records.

P02.06 retains the `kernel.authorization` contextual layer on top of P02.05: direct RBAC is always evaluated first; trusted relationship evidence must match exact principal/object/tenant/organization scope; contextual constraints may only narrow authority; internal/background call origin never bypasses authorization; sensitive-field/export capability checks remain distinct from ordinary read; material denials/privileged decisions use safe required audit; relationship/constraint dependency failure fails closed.

P02.07 retains the `kernel.identity` strong-authentication foundation: passkey factor lifecycle is deterministic; passkey verification uses an injected verifier rather than custom protocol/private-key crypto; one-time challenges are exact User/session bound, expiring and replay-safe; recovery codes are one-time with digest-only persistence; step-up proof is session-bound and non-authorizing; factor removal follows explicit session invalidation policy; restricted factor/challenge/recovery material does not enter ordinary telemetry/audit; and migration v3 is covered by fresh/idempotent/P02.04-upgrade/immutable-ledger evidence.

## Candidate active P02.08 boundary

Only after this closure merges and protected `main` + canonical `STATE.json` confirm P02.08 ACTIVE, P02.08 may implement only:

- Service Account as a distinct non-human principal under `kernel.identity`;
- tenant/organization binding and capability/permission scope composition through accepted P02 authorization foundations;
- API credential issue/identify/verify/rotate/revoke/expire lifecycle;
- one-time secret presentation where applicable and non-reversible verifier storage;
- classification-safe credential inventory/audit metadata while raw secrets remain `RESTRICTED`;
- last-used/rotation metadata only where classification-safe and useful;
- deterministic allowed-scope plus wrong-tenant/wrong-scope/revoked/expired negative evidence;
- applicable owner-bounded persistence, migrations and focused identity/security evidence.

P02.08 must represent Service Accounts as distinct non-human principals rather than fake human Users. Raw API credentials remain `RESTRICTED`, cannot be logged/traced/audited/stored reversibly, must be least-privilege, rotatable, revocable and tenant/org bound, and possession never bypasses current authorization. No generic platform superkey/master token is authorized.

P02.08 must not implement OAuth developer applications, external connector/provider integration, device/POS identities, AI agent execution identity, P02.09 tenant settings, P02.10 exit-product behavior, business API scopes/features/UI, deployment authority or other future scope.

## Exact next authorized action

After this P02.07 closure/P02.08 activation PR merges, verify protected `main` and canonical `STATE.json`, identify P02.08 as the sole active package, then **STOP the execution session**.

A later governed execution session must start from the then-current protected-main SHA, re-read canonical state, open PRs and the frozen P02.08 acceptance criteria, and only then may implement P02.08. Do not auto-advance to P02.09.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.08.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.08.md`.
