# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.04 — PostgreSQL connection & migration foundation**
- P01 progress: **3 / 12 done**
- P01.01: **DONE**
- P01.02: **DONE**
- P01.03: **DONE**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**
- P01.05–P01.12: **PLANNED / NOT ACTIVE**

## Completed P01 packages

### P01.01 — Go workspace/build skeleton

Canonical evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`.

PR #40 merged as `7257977264d788663083fa215462b1828f1e5afb` after GitHub-hosted run `32562869345` / job `97007065640` passed the pinned Go, format/static, unit, dependency-boundary and build/smoke gates.

### P01.02 — Configuration & environment system

Canonical evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

PR #42 merged as `c857bb9e7df1e347226653eeaded024d6ecd0271` after corrected GitHub-hosted run `32563880800` / job `97009520624` passed P01.01 regression, format/static, unit/race, precedence/startup, secret-safe negative, dependency-boundary and build/package gates.

Initial run `32563817196` remains preserved as failure evidence and was fixed rather than bypassed.

### P01.03 — Structured error & result conventions

Canonical evidence: `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`.

PR #44 merged as `bdeda5fad09a2369b2a6852e5c62550db50047ea` after GitHub-hosted run `32565935613` / job `97014452248` passed:

- P01.01 and P01.02 regression verification;
- Go `1.26.7` exact pin;
- format/static and dependency-boundary checks;
- unit + race tests;
- `errors.Is` / `errors.As` wrapping behavior;
- private-cause/public-redaction negative tests;
- deterministic bounded validation detail;
- build/package.

Earlier formatting/diagnostic failures remain recorded in the canonical completion evidence and are not counted as PASS.

## Active P01 package — P01.04

P01.04 is the sole active package. It owns only the PostgreSQL connection and migration foundation:

- pool/connection construction from P01.02 configuration;
- safe provider failure mapping through P01.03;
- bounded connection behavior;
- transaction helper boundary without domain writes;
- migration runner and version ledger;
- deterministic fresh and upgrade migrations;
- migration coordination/ownership rules;
- synthetic-data PostgreSQL integration tests;
- P01.01-P01.03 regression preservation.

Explicitly prohibited in P01.04: tenant/organization tables, module-runtime schema, event outbox/inbox, business schemas/data, cross-module SQL, cache/storage, telemetry, health endpoints, production HA/backups or later domain repositories.

P01.05 may activate only after P01.04 reaches `done` with required G0/G1/G2/G3/G4/G5/G7 GitHub-hosted evidence.

## Protected integration / CI

`main.protected=true` remains verified. PR-only integration, required `governance`, failed-check rejection, direct/force update rejection and conversation resolution remain enforced.

Canonical governance CI remains GitHub-hosted only on `ubuntu-24.04`; local/self-hosted runners are prohibited. Completed P01 verifiers remain regression gates in the same required `governance` job.

## Dependabot

Repository-managed `.github/dependabot.yml` is present for weekly GitHub Actions dependency updates. Dependency updates remain governed through normal protected PRs.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block authorized P01 kernel engineering.

## Exact next work

Implement **P01.04 only**, obtain hosted PostgreSQL connection/migration/security/build evidence, reconcile P01.04 to `done`, then activate only P01.05. Keep `business_feature_code_authorized=false` and P02+ implementation locked.
