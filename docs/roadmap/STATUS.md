# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.02 — Configuration & environment system**
- P01 progress: **1 / 12 done**
- P01.01: **DONE — hosted executable evidence recorded**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**
- P01.03–P01.12: **PLANNED / NOT ACTIVE**

## P01.01 completion

P01.01 merged through PR #40 as `7257977264d788663083fa215462b1828f1e5afb` after GitHub-hosted governance run `32562869345`, job `97007065640`, completed SUCCESS.

Completion evidence is canonicalized at `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md` and records:

- Go `1.26.7` exact pin: PASS;
- G0 governance: PASS;
- G1 format/static/workspace/dependency-boundary: PASS;
- G2 unit/smoke: PASS;
- G7 build/package: PASS;
- GitHub-hosted Ubuntu 24.04 / X64 execution: PASS;
- no P01.02+, P02/P03 or business-domain behavior introduced.

## Active P01 package — P01.02

P01.02 is the sole active package. It owns the typed configuration/environment system only:

- deterministic configuration loading and validation;
- explicit precedence across defaults, config file and environment variables;
- governed environment identity;
- required/optional setting semantics;
- secret-safe redaction;
- deterministic isolated test overrides;
- provenance diagnostics that never reveal secret values;
- fail-closed startup on invalid required configuration.

Explicitly prohibited in P01.02: structured application error model beyond narrow config errors (P01.03), PostgreSQL/migrations, cache/storage clients, telemetry, health endpoints, jobs, feature flags, audit transport, tenancy/organizations, module runtime and business configuration.

P01.03 may activate only after P01.02 reaches `done` with required G0/G1/G2/G5/G7 hosted evidence.

## Protected integration / CI

`main.protected=true` remains verified. PR-only integration, required `governance`, failed-check rejection, direct/force update rejection and conversation resolution are enforced. Current single-maintainer review policy keeps required approvals at zero until an independent reviewer exists.

Canonical governance CI remains GitHub-hosted only on `ubuntu-24.04`; local/self-hosted runners are prohibited.

## Dependabot

Repository-managed `.github/dependabot.yml` is merged. GitHub Actions dependencies are checked weekly. The P01.01 run emitted a non-blocking Node runtime deprecation warning for `actions/checkout@v4`; Dependabot remains the governed dependency-update path.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block the active P01 kernel package but must be resolved before decisions that depend on commercial/public distribution policy.

## Exact next work

Implement **P01.02 only**, obtain the required hosted configuration/security/build evidence, reconcile P01.02 to `done`, then activate only P01.03. Keep `business_feature_code_authorized=false` and P02+ implementation locked.
