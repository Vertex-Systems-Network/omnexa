# ADR-0006 — Temporary P00 CI Evidence Exception

Status: **Expired — historical evidence only**  
Date: **2026-08-22**  
Original scope: **P00 architecture/governance/specification work only**

## Context

During P00.06, GitHub Actions allowance was exhausted/disabled and multiple runs failed before runner steps executed. The project owner authorized a narrow temporary exception for P00 documentation/specification work. Issue #14 recorded the infrastructure evidence.

## Historical decision

The exception allowed P00 documentation/specification packages to progress using explicit manual evidence while hosted execution was recorded as `BLOCKED`/`NOT RUN`, never as PASS. It never authorized runtime/kernel/business implementation and never waived executable build, migration, test, security or release gates for P01+.

## Expiry

This exception expired on **2026-08-22** before executable P01 implementation because GitHub-hosted CI was restored and P00 exited.

Current consequences:

- ADR-0006 is historical provenance only;
- it cannot authorize a present or future CI bypass;
- canonical governance execution is GitHub-hosted only on `ubuntu-24.04`;
- P01 executable changes require real automated evidence;
- any future CI exception requires new explicit change control/ADR rather than reuse of this record.

## Why history is retained

The ADR explains why some P00 merges correctly carry `BLOCKED`/manual evidence instead of hosted PASS. Historical evidence must not be rewritten to make earlier progress look cleaner.

## Rejected alternatives retained from the original decision

- Treat failed Actions runs as successful — rejected.
- Permanently remove CI as a requirement — rejected.
- Extend the exception into P01 executable implementation — prohibited.

No architecture, security, tenancy, API, event or product semantics were changed by this temporary operational exception.
