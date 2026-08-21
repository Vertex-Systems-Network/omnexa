# ADR-0006 — Temporary P00 CI Evidence Exception

Status: **Accepted — temporary**  
Date: **2026-08-22**  
Scope: **P00 architecture/governance/specification work only**

## Context

Omnexa's governance model normally requires successful GitHub Actions evidence before merging foundation work. During P00.06, the repository exhausted/disabled its GitHub Actions allowance. Multiple PR #13 workflow runs, including a rerun and three independent diagnostic jobs, failed before any runner step executed. Issue #14 records the infrastructure evidence.

The project owner explicitly authorized skipping GitHub Actions while the quota condition persists and instructed the project to continue the remaining P00 planning/specification work.

## Decision

Adopt the narrow operational exception defined in `docs/governance/CI_EVIDENCE_EXCEPTION_2026-08-22.md`.

While the exception is active:

- hosted GitHub Actions evidence is recorded as `BLOCKED` or `NOT RUN`, never fabricated as `PASS`;
- P00 documentation/specification packages may progress using repository diff inspection, evidence-file verification, state/dependency reconciliation and ADR review;
- no runtime/kernel/business implementation is authorized;
- no executable build, migration, test, security scan or release gate for P01+ may be waived by this ADR;
- PR-based integration remains mandatory even though `main` branch protection is not yet technically enforced;
- when hosted Actions return, the current governance validator is rerun and any discovered violation reopens affected work.

## Why this is acceptable

The active P00 work is architecture/governance/specification-only and the implementation lock explicitly prohibits runtime software. Therefore build, migration and runtime test gates are already `N/A` for these changes. The unavailable evidence is the hosted execution of governance checks, not executable product correctness.

Manual evidence still must prove scope, artifact presence, state consistency and change-control reconciliation.

## Rejected alternatives

### Stop all P00 work until quota reset

Rejected because it would block non-runtime architecture work solely due to hosted runner availability and the owner explicitly approved a temporary exception.

### Treat failed Actions runs as successful

Rejected. Infrastructure-blocked execution is not a pass.

### Permanently remove CI as a requirement

Rejected. P00.07 will define the durable testing/CI/release model and P01+ executable changes require real automated evidence.

## Compatibility and risk

Risk is limited by the P00 implementation lock and the docs/specification-only scope. The main residual risk is a governance-state/schema inconsistency that would normally be caught by automation. The mitigation is explicit manual verification plus mandatory revalidation when Actions returns.

## Expiry

This ADR's operational exception expires before the first P01 implementation merge, at P00 exit review, when Actions becomes available, or when revoked by the owner—whichever occurs first.

It remains in history afterward as evidence of why certain P00 merges show `BLOCKED` hosted-CI evidence.