# Omnexa P02 Exit Gate

Status: **NOT SATISFIED — PHASE ACTIVE**  
Owner phase: **P02 — Identity, Tenancy & Organization**

This document defines the evidence required before P02 may be reconciled `done`. It is a gate definition, not implementation authority. Current progress is 2 / 10 done: P02.01-P02.02 are complete and P02.03 is active.

## Mandatory exit evidence

P02 exit requires all P02.01-P02.10 packages to be `done` and exact canonical GitHub-hosted `ubuntu-24.04` evidence proving:

- cross-tenant access cannot escape the trusted tenant boundary;
- organization/object/scope permissions enforce same-scope allow and wrong-scope deny behavior;
- role/permission differences are deterministic and no role name creates a bypass;
- service accounts are non-human principals with bounded, rotatable/revocable credentials and scope;
- session invalidation behavior is explicit: interactive sessions expire, rotate/revoke and invalidate after required account/security changes;
- authentication secrets/tokens/recovery material are not leaked to ordinary logs, traces, errors, audit payloads or test fixtures;
- identity, tenant, organization and authorization mutations use their authoritative owners;
- privileged identity/permission actions emit classification-safe attributable audit evidence through `kernel.audit`;
- fresh and supported-upgrade database migrations pass whenever P02 persistence changes;
- repository Go quality and all completed P01/P02 regressions pass;
- applicable G0-G8 gates are recorded using only PASS/FAIL/BLOCKED/NOT RUN/N/A.

## Current retained evidence

P02.01 completion is recorded in `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md` from implementation PR #69, final exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical run/job `32635243643 / 97183883007` PASS and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 completion is recorded in `docs/roadmap/evidence/P02.02_COMPLETION_2026-08-23.md` from implementation PR #71, final exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical run/job `32637760875 / 97189971101` PASS and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`.

P02.02 provides canonical same-tenant allow/cross-tenant forged-selector deny/no-global-fallback evidence and establishes the trusted tenant boundary. This satisfies the tenant-isolation foundation portion only. It does not satisfy organization/object/scope authorization, role/permission differences, service-account lifecycle or session invalidation requirements that belong to later P02 packages.

## Exit-denial conditions

P02 exit is blocked by any known required-gate failure, unresolved cross-tenant or privilege-escalation defect, hidden administrator bypass, unverified migration path, secret leakage, ambiguous identity/organization ownership, incomplete P02 package sequence, or premature P03/business/AI scope.

## Transition rule

P03 remains planned until P02 exit is satisfied through a separate governed closure. Completion of P02.10 does not implicitly activate P03.
