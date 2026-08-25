# Omnexa P03 Exit Gate

Status: **NOT SATISFIED — PHASE NOT ACTIVE**  
Owner phase: **P03 — Module Runtime**

This document defines the evidence required before P03 may be reconciled `done`. It is a gate definition, not implementation authority.

## Mandatory exit evidence

P03 exit requires P03.01-P03.11 to be `done` and exact canonical GitHub-hosted `ubuntu-24.04` evidence proving all applicable G0-G8 gates plus the reference-module behaviors below.

### EX-01 — Required dependency enforcement

A reference module with a missing/incompatible required dependency cannot install or enable successfully. Resolution and failure are deterministic, explicit and do not partially mutate unrelated module state.

### EX-02 — Optional dependency degradation

A reference module with an absent optional dependency remains healthy for unaffected capabilities and selectively degrades only the dependent feature. Optional absence cannot crash unrelated modules or global boot.

### EX-03 — Safe disable and re-enable

Disable is non-destructive, preserves owned data/history required for re-enable/audit/reference continuity, and makes contributed capabilities/permissions/UI/settings availability accurately reflect lifecycle state. Re-enable restores valid behavior without manual repair.

### EX-04 — Upgrade and migration path

Fresh install and supported upgrade paths execute only owner-approved migrations, preserve tenant boundaries and data integrity, and handle retry/failure/rollback-or-forward-fix semantics explicitly.

### EX-05 — Forbidden dependency detection

Canonical tests reject private implementation imports, cross-module direct writes, circular required dependencies, undeclared dependencies and other `DEPENDENCY_MATRIX.md` forbidden paths.

### EX-06 — Health and state accuracy

Module lifecycle, dependency, migration and capability availability are reported accurately through classification-safe health diagnostics. Required dependency failure is fail-closed; optional degradation is distinguishable from total failure.

### EX-07 — Unrelated module isolation

Install, enable, disable, suspend, archive, detach, purge and failed lifecycle operations on a reference module do not corrupt unrelated module state, data, registrations or health.

## Contract coverage required before exit

Evidence must also prove:

- manifest schema/version validation and stable module identity;
- deterministic registry/discovery with duplicate/conflict rejection;
- required/optional/platform/forbidden dependency semantics;
- explicit lifecycle transition validity, retry/idempotency and destructive-operation authorization;
- module settings/feature flags integrate through `kernel.configuration` without creating authority;
- capability registration is owner/version/scope aware and does not bypass owning boundaries;
- permission registration integrates with `kernel.authorization` deny-by-default semantics;
- UI contributions are declarative and never authorize backend operations;
- migration ownership is machine-enforced and cross-owner schema mutation is rejected unless explicitly approved;
- signed-package/provenance/SBOM/scope hooks are represented for later `XTRUST-100` without falsely claiming trust enforcement;
- P03 AI-native compatibility requirements remain forward-compatible without pulling strategic-program runtime forward.

## Security and regression evidence

Before P03 exit:

- P01.01-P01.12 regression verifiers pass;
- P02.01-P02.10 regression verifiers pass;
- repository Go quality passes;
- applicable tenant-isolation, authorization, audit, migration, lifecycle, concurrency and failure-isolation tests pass;
- manifests/package metadata are treated as untrusted and cannot cause implicit code execution during parse/discovery;
- secrets, credentials and RESTRICTED material are not exposed through manifests, health output or ordinary logs;
- exact-head completion evidence records PR, source SHA, workflow run/job, canonical environment and result;
- diagnostic failures remain explicit and are never relabelled PASS.

## Phase transition rule

P03 completion does not activate P04 automatically. After P03 terminal closure is accepted on protected `main`, `kernel_code_authorized=false`, `business_feature_code_authorized=false`, current work package is `NONE`, and P04 remains `planned` until a separate governed preparation/readiness and activation flow.

## Current decision

```text
P03 exit: NOT SATISFIED
P03: PLANNED / NOT ACTIVE
P03.01-P03.11: PLANNED
P04: PLANNED
implementation authority: NONE
```
