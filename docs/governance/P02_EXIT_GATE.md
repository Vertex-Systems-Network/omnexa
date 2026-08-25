# Omnexa P02 Exit Gate

Status: **NOT SATISFIED — PHASE ACTIVE**  
Owner phase: **P02 — Identity, Tenancy & Organization**

This document defines the evidence required before P02 may be reconciled `done`. It is a gate definition, not implementation authority. Current progress is 9 / 10 done: P02.01-P02.09 are complete and P02.10 is active.

## Mandatory exit evidence

P02 exit requires all P02.01-P02.10 packages to be `done` and exact canonical GitHub-hosted `ubuntu-24.04` evidence proving:

- cross-tenant access cannot escape the trusted tenant boundary;
- organization/object/scope permissions enforce same-scope allow and wrong-scope deny behavior;
- role/permission differences are deterministic and no role name creates a bypass;
- service accounts are non-human principals with bounded, rotatable/revocable credentials and scope;
- session invalidation behavior is explicit: interactive sessions expire, rotate/revoke and invalidate after required account/security changes;
- authentication secrets/tokens/recovery/API credential material are not leaked to ordinary logs, traces, errors, audit payloads or test fixtures;
- identity, tenant, organization, authorization and tenant-settings mutations use their authoritative owners;
- privileged identity/permission/settings actions emit classification-safe attributable audit evidence through `kernel.audit` where applicable;
- fresh and supported-upgrade database migrations pass whenever P02 persistence changes;
- repository Go quality and all completed P01/P02 regressions pass;
- applicable G0-G8 gates are recorded using only PASS/FAIL/BLOCKED/NOT RUN/N/A.

## Current retained evidence

P02.01 completion is recorded in `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md` from implementation PR #69, final exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical run/job `32635243643 / 97183883007` PASS and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 completion is recorded in `docs/roadmap/evidence/P02.02_COMPLETION_2026-08-23.md` from implementation PR #71, final exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical run/job `32637760875 / 97189971101` PASS and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`.

P02.03 completion is recorded in `docs/roadmap/evidence/P02.03_COMPLETION_2026-08-23.md` from implementation PR #73, final exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`, canonical run/job `32640790333 / 97197453122` PASS and merge `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`.

P02.04 completion is recorded in `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md` from implementation PR #75, final exact head `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f`, canonical run/job `32653747461 / 97229198036` PASS and merge `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`.

P02.05 completion is recorded in `docs/roadmap/evidence/P02.05_COMPLETION_2026-08-24.md` from implementation PR #77, final exact head `2df8d2a8bef0cea60256a832986d6f8495c80378`, canonical run/job `32660848145 / 97246683239` PASS and merge `7b6a59e83c9bd696e6e008385b4413d529254171`.

P02.06 completion is recorded in `docs/roadmap/evidence/P02.06_COMPLETION_2026-08-24.md` from implementation PR #79, final exact head `dbbd105fd5f2543ca7dd5df93375eaf1057928fc`, canonical run/job `32664834112 / 97256520050` PASS and merge `083c2866f0cd0773b85201750c2196bfd2fcc167`.

P02.07 completion is recorded in `docs/roadmap/evidence/P02.07_COMPLETION_2026-08-24.md` from implementation PR #81, final exact head `51ccaa12c3534f74fba6eab9d4698ee483ef4ffd`, canonical run/job `32669167972 / 97267175953` PASS and merge `5642f5da1eb24e70b67e5ec757d9f4584c4e3f5c`.

P02.08 completion is recorded in `docs/roadmap/evidence/P02.08_COMPLETION_2026-08-25.md` from implementation PR #84, final exact head `43bdcf525ce5e0cfdb9dc0707fbafee7cd552543`, canonical run/job `32885950897 / 97926598423` PASS and merge `32eb7187eb229327585551e4e28b0d596de78bd9`.

P02.09 completion is recorded in `docs/roadmap/evidence/P02.09_COMPLETION_2026-08-26.md` from implementation PR #86, final exact head `0618904a18f82231469dd173aeb3d9d51edb73ed`, canonical run/job `32895186252 / 97956097639` PASS and merge `8ef86d2644b5ed455b3610192b8379d94204692f`.

P02.04 proves the interactive authentication/session portion of the phase exit: explicit access/refresh/session expiry, rotation/revocation, replay denial, password-change and account-lifecycle invalidation, secret non-disclosure, current tenant/organization context reauthorization and fresh/upgrade persistence evidence.

P02.05 proves the direct RBAC portion of the phase exit: stable capability permission identifiers, deterministic Role composition, deny-by-default direct decisions, tenant/organization exact-scope assignments, anti-escalation on privileged mutation, assignment revocation, role-name non-bypass, classification-safe required audit records and owner-bounded PostgreSQL persistence.

P02.06 proves the relationship/context-aware authorization portion of the phase exit: accepted direct RBAC remains mandatory; trusted relationship evidence is exact principal/object/tenant/organization scoped; contextual constraints cannot widen authority; wrong tenant/org/object, missing permission and internal/background bypass attempts deny deterministically; field/export capability boundaries may be stricter than ordinary read; material denials and privileged decisions are audited safely; and dependency failures fail closed. P02.06 introduced no new persistence, so its G4 is N/A for new migration while retained P02.05 migration evidence passed.

P02.07 proves the human strong-authentication portion of the phase exit: passkey factor lifecycle is deterministic; challenges are exact User/session bound, expiring and replay-safe; recovery codes are one-time and persisted only as digests; step-up proof is session-bound and non-authorizing; factor removal follows explicit session invalidation policy; restricted factor/challenge/recovery material is excluded from ordinary telemetry/audit; and fresh/idempotent/P02.04-upgrade migration evidence passed.

P02.08 proves the service-account/API-credential portion of the phase exit: Service Accounts are distinct non-human principals; credential proof is exact tenant/organization bound; raw secrets are one-time issuance material with SHA-256 verifier-only persistence; revoked, expired and superseded credentials deny; rotation invalidates the prior credential transactionally; current principal/tenant/assignment state remains authoritative; direct RBAC exact-scope permission composition passes; and fresh/idempotent/P02.07+P02.05 supported-upgrade migration evidence passed.

P02.09 proves the tenant-settings portion of the phase exit: trusted setting scope derives only from accepted tenant/organization context; exact organization overrides fall back only to the enclosing tenant and then the registered definition default; no global/user override exists; protected reads and all writes pass through current authorization; generic RESTRICTED/secret values are rejected; security-significant changes are audited without values; cross-tenant/wrong-org access denies; and fresh/idempotent/P02.08 supported-upgrade migration evidence passed. P02.10 still must provide the aggregate identity/permission audit-trail coverage and final P02 exit proof.

## Exit-denial conditions

P02 exit is blocked by any known required-gate failure, unresolved cross-tenant or privilege-escalation defect, hidden administrator bypass, unverified migration path, secret leakage, ambiguous identity/organization/authorization/configuration ownership, incomplete P02 package sequence, or premature P03/business/AI scope.

## Transition rule

P03 remains planned until P02 exit is satisfied through a separate governed closure. Completion of P02.10 does not implicitly activate P03.
