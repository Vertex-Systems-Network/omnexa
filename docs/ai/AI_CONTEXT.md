# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

After this closure is accepted and merged, the intended canonical checkpoint is:

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 5 / 10 done.
- P02.01-P02.05: DONE.
- sole active package: `P02.06 — Relationship/context-aware authorization`.
- owner: `kernel.authorization`.
- `kernel_code_authorized=true` only for P02.06 after closure merge and post-merge verification.
- `business_feature_code_authorized=false`.
- P02.07-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

Until this closure merges, protected `main` remains authoritative and records the accepted P02.05 implementation while P02.05 is still canonically active. This continuity file creates no implementation authority.

## P02.05 completion evidence

- implementation PR: #77
- final exact implementation head: `2df8d2a8bef0cea60256a832986d6f8495c80378`
- canonical run/job: `32660848145 / 97246683239` — PASS
- implementation merge: `7b6a59e83c9bd696e6e008385b4413d529254171`
- runner: `GitHub Actions 1000015357`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- runner image: `ubuntu-24.04 / 20260816.277.1`
- Go: 1.26.7
- PostgreSQL: 18.6
- repository Go quality, P01.01-P01.12, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.04 and P02.05 G0-G8: PASS
- evidence: `docs/roadmap/evidence/P02.05_COMPLETION_2026-08-24.md`

Retained diagnostic failures remain visible rather than relabeled: `32656689041 / 97236397635` on exact head `26653b09f8ae678db74d17486b3d0bb988ce3b8d` failed repository Go quality for 29 corrected `govet` shadow findings plus one corrected `gosec` G304 test-helper finding. `32660398632 / 97245574862` on exact head `42df3d449aec49f1c1d1226985837f0501178fd4` passed repository quality, all P01/P02 prior regressions and P02.05 runtime/integration evidence but failed the initial verifier's transitive-vs-direct dependency check. Both were corrected without suppression, gate weakening or acceptance changes and are not completion evidence.

## Retained P02.01-P02.05 contract

P02.01 retains the owner-bounded `kernel.identity` human User foundation. User remains distinct from business Person and identity attributes do not create authority.

P02.02 retains the owner-bounded `kernel.tenancy` foundation: UUIDv7 Tenant lifecycle, minimal User↔Tenant membership with terminal revocation, trusted tenant context derived from active persisted Tenant/membership state, explicit tenant-safe scope equality, owner-bounded PostgreSQL persistence and no global-tenant fallback.

P02.03 retains the `kernel.organization` tenant-contained Organization/Legal Entity/Business Unit/Branch/Team/Location hierarchy, scoped memberships, deterministic cycle/cross-tenant rejection and non-authorizing organization scope context. Organization remains distinct from business Party Organization.

P02.04 retains the `kernel.identity` authentication/session foundation: approved adaptive password hashing, disclosure-safe authentication, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/revocation/replay rejection, password/account invalidation, safe device/session inventory, tenant/organization context re-authorization and classification-safe security lifecycle hooks. Authentication remains distinct from authorization.

P02.05 retains the `kernel.authorization` direct RBAC foundation: stable capability-oriented permission identifiers, deterministic Role composition, trusted tenant/organization scoped assignments, deny-by-default direct permission checks, privileged server-side mutations with anti-escalation, assignment revocation, role-name non-bypass and classification-safe required audit records. Direct RBAC does not itself implement relationship/object/context policy.

## Candidate active P02.06 boundary

Only after this closure merges and protected `main` + canonical `STATE.json` confirm P02.06 ACTIVE, P02.06 may implement only:

- relationship/context policy evaluation layered on accepted P02.05 RBAC;
- trusted tenant, organization and object-scope relationships;
- capability-bound deny-by-default authorization decisions;
- contextual conditions that cannot grant authority outside the principal's valid relationships;
- field/export distinction hooks where broader read authority must not imply sensitive-field/export authority;
- disclosure-safe deny behavior and classification-safe material authorization audit hooks;
- same-scope allow plus wrong-tenant/wrong-org/wrong-object and missing-permission negative tests;
- applicable owner-bounded persistence/migrations and focused deterministic authorization evidence.

P02.06 must preserve P02.05 deny-by-default RBAC. Client IDs, tenant IDs and object IDs are references, never authority. Tenant membership alone is insufficient for all child scopes. Role names and internal/background call origin never bypass policy. Contextual rules cannot widen beyond trusted principal, tenant, organization and object relationships.

P02.06 must not implement P02.07 MFA/passkeys, P02.08 service-account/API credentials, P02.09 tenant settings, P02.10 exit-product behavior, P03 module permission registration/capability registry, business-domain object policies not owned by an active domain, support impersonation product surface, business features/UI, deployment authority or AI/model/agent runtime.

## Exact next authorized action

After this P02.05 closure/P02.06 activation PR merges, verify protected `main` and canonical `STATE.json`, identify P02.06 as the sole active package, then **STOP the execution session**.

A later governed execution session must start from the then-current protected-main SHA, re-read canonical state, open PRs and the frozen P02.06 acceptance criteria, and only then may implement P02.06. Do not auto-advance to P02.07.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P02 entry/exit gates, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.06.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.06.md`.
