# Omnexa P02 Exit Gate

Status: **NOT SATISFIED — PHASE NOT ACTIVE**  
Owner phase: **P02 — Identity, Tenancy & Organization**

This document defines the evidence required before P02 may be reconciled `done`. It is a gate definition, not implementation authority.

## Mandatory exit evidence

P02 exit requires all P02.01-P02.10 packages to be `done` and exact canonical GitHub-hosted `ubuntu-24.04` evidence proving:

- cross-tenant access cannot escape the trusted tenant boundary;
- organization/object/scope permissions enforce same-scope allow and wrong-scope deny behavior;
- role/permission differences are deterministic and no role name creates a bypass;
- service accounts are non-human principals with bounded, rotatable/revocable credentials and scope;
- interactive sessions expire, rotate/revoke and invalidate after required account/security changes;
- authentication secrets/tokens/recovery material are not leaked to ordinary logs, traces, errors, audit payloads or test fixtures;
- identity, tenant, organization and authorization mutations use their authoritative owners;
- privileged identity/permission actions emit classification-safe attributable audit evidence through `kernel.audit`;
- fresh and supported-upgrade database migrations pass whenever P02 persistence changes;
- repository Go quality and all completed P01/P02 regressions pass;
- applicable G0-G8 gates are recorded using only PASS/FAIL/BLOCKED/NOT RUN/N/A.

## Exit-denial conditions

P02 exit is blocked by any known required-gate failure, unresolved cross-tenant or privilege-escalation defect, hidden administrator bypass, unverified migration path, secret leakage, ambiguous identity/organization ownership, or premature P03/business/AI scope.

## Transition rule

P03 remains planned until P02 exit is satisfied through a separate governed closure. Completion of P02.10 does not implicitly activate P03.
