# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.01 — Go workspace/build skeleton**
- P01 progress: **0 / 12 done**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- P00.10: **DONE**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**
- P01.02–P01.12: **PLANNED / NOT ACTIVE**

## P00 exit result

P00.01–P00.10 are complete. P00.01–P00.09 remain frozen as Omnexa Foundation Architecture v1; P00.10 verified and exited the foundation phase. Material reinterpretation of frozen architecture requires change control and a superseding accepted ADR.

ADR-0006's temporary P00 CI exception is expired and historical-only. It cannot authorize any current bypass.

## P01 entry controls

### EG-02 / Issue #3 — SATISFIED

Verified controls include:

- live `main.protected=true`;
- PR-only integration;
- required `governance` status check;
- direct fast-forward update rejected;
- force update rejected;
- failed-governance PR #34 rejected;
- unresolved conversation on CODEOWNERS-path PR #37 rejected until resolution;
- resolved, green PR #37 merged successfully;
- green Dependabot PR #35 merged successfully;
- force pushes and branch deletion remain blocked by configured ruleset.

Current single-maintainer review policy uses zero required approvals and no required Code Owner review to avoid self-review deadlock. Tighten this when an independent reviewer exists.

### EG-03 / Issue #14 — SATISFIED

Canonical governance CI runs only on:

```yaml
runs-on: ubuntu-24.04
```

The required job is `governance` and must prove `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. PR #37 run `32541439589` is current positive evidence. No self-hosted/local fallback is permitted.

## Dependabot

Repository-managed `.github/dependabot.yml` is merged. GitHub Actions dependencies are checked weekly and minor/patch updates are grouped. Go/npm ecosystems will be added only when their manifests exist in governed implementation scope.

## Active P01 package

P01.01 is the sole active package. Scope is limited to the Go workspace/build skeleton: toolchain declaration, workspace/module layout, minimal kernel process, build metadata and deterministic format/vet/test/build/smoke verification.

Explicitly prohibited in P01.01: configuration system, persistence/migrations, cache/storage, telemetry, health endpoints, jobs, feature flags, audit transport, full CLI, identity/tenancy, module runtime and business modules.

P01.02 may activate only after P01.01 reaches `done` with required `G0/G1/G2/G7` hosted evidence.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This gate does not block the authorized P01 kernel engineering scope, but it must be resolved before claims or decisions that depend on commercial/public distribution policy.

## Exact next work

Implement **P01.01 only** in a separate executable PR after this governance transition is merged and verified on `main`. Keep `business_feature_code_authorized=false` and P01.02–P01.12 planned.
