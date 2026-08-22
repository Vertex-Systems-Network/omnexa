# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.03 — Structured error & result conventions**
- P01 progress: **2 / 12 done**
- P01.01: **DONE**
- P01.02: **DONE**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**
- P01.04–P01.12: **PLANNED / NOT ACTIVE**

## Completed P01 packages

### P01.01 — Go workspace/build skeleton

Canonical evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`.

PR #40 merged as `7257977264d788663083fa215462b1828f1e5afb` after GitHub-hosted run `32562869345` / job `97007065640` passed the pinned Go, format/static, unit, dependency-boundary and build/smoke gates.

### P01.02 — Configuration & environment system

Canonical evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

PR #42 merged as `c857bb9e7df1e347226653eeaded024d6ecd0271` after corrected GitHub-hosted run `32563880800` / job `97009520624` passed:

- P01.01 regression verification;
- Go `1.26.7` exact pin;
- format/static/dependency boundary;
- unit + race tests;
- deterministic precedence/startup smoke;
- secret-safe negative tests;
- build/package.

Initial run `32563817196` failed on a test typing defect and remains preserved as failure evidence; it was fixed rather than bypassed.

## Active P01 package — P01.03

P01.03 is the sole active package. It owns transport-neutral structured error/result primitives only:

- stable machine error codes;
- safe public message/detail separated from private causes;
- wrapping/unwrapping with standard Go error semantics;
- explicit retryability/category metadata;
- deterministic bounded validation-field errors;
- correlation metadata boundary without telemetry emission;
- negative redaction/security tests;
- P01.01/P01.02 regression preservation.

Explicitly prohibited in P01.03: HTTP adapter behavior beyond shared hooks, PostgreSQL/provider mapping, cache/storage behavior, logging/OpenTelemetry emission, health endpoints, jobs, tenancy/authz/module runtime and business-domain error catalogs.

P01.04 may activate only after P01.03 reaches `done` with required G0/G1/G2/G5/G7 hosted evidence.

## Protected integration / CI

`main.protected=true` remains verified. PR-only integration, required `governance`, failed-check rejection, direct/force update rejection and conversation resolution remain enforced.

Canonical governance CI remains GitHub-hosted only on `ubuntu-24.04`; local/self-hosted runners are prohibited. Completed P01 verifiers remain regression gates in the same required `governance` job.

## Dependabot

Repository-managed `.github/dependabot.yml` is present for weekly GitHub Actions dependency updates. Dependency updates remain governed through normal protected PRs.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block authorized P01 kernel engineering.

## Exact next work

Implement **P01.03 only**, obtain hosted structured-error/security/build evidence, reconcile P01.03 to `done`, then activate only P01.04. Keep `business_feature_code_authorized=false` and P02+ implementation locked.
