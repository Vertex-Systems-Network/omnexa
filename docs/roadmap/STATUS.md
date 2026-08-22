# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.05 — Cache abstraction**
- P01 progress: **4 / 12 done**
- P01.01: **DONE**
- P01.02: **DONE**
- P01.03: **DONE**
- P01.04: **DONE**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**
- P01.06–P01.12: **PLANNED / NOT ACTIVE**

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

### P01.04 — PostgreSQL connection & migration foundation

Canonical evidence: `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`.

PR #46 was verified by GitHub-hosted run `32567842071` / job `97019012280` and squash-merged to protected `main` as `6068202415dd124d3e74a196b6e0bbca5d75c4cd`.

Verified baseline:

- Go `1.26.7`;
- pgx `v5.10.0`;
- PostgreSQL `18.6 (Debian 18.6-1.pgdg13+2)`;
- P01.01-P01.03 regressions PASS;
- P01.04 G1/G2/G3/G4/G5/G7 PASS;
- G0 governance/state validators PASS.

Synthetic PostgreSQL evidence covered server ping/version, bounded pool exhaustion, unavailable connection, transaction commit/rollback, fresh migration, idempotent rerun, deterministic upgrade, failed-migration rollback without ledger advance, rewritten-migration drift detection and advisory-lock contention.

Diagnostic run `32567607242` / job `97018421698` remains preserved as FAIL history after exposing a formatting defect and canonical module metadata. It is not counted as PASS.

P01.04 did not introduce an ORM, tenant/organization schema, module-runtime schema, event outbox/inbox, business schema/data, cache/storage, telemetry, health or later capability implementation.

## Active P01 package — P01.05

P01.05 is the sole active package. It owns only the Redis-compatible cache abstraction foundation:

- cache interface/provider adapter;
- deterministic namespaced/versioned keys;
- explicit TTL/expiry semantics;
- typed serialization boundary;
- get/set/delete and only justified bounded atomic primitives;
- P01.02 configuration integration and P01.03 safe failure mapping;
- cache-miss versus provider-failure distinction;
- bounded timeout/cancellation behavior;
- synthetic Redis-compatible integration evidence;
- provider flush/restart semantics proving cache is never authoritative;
- P01.01-P01.04 regression preservation.

Explicitly prohibited in P01.05: business-domain keys/models, sessions/authentication, tenant/organization behavior, module runtime, event/job fabric, object storage, telemetry, health, later P01 capability implementation or use of cache as system of record.

P01.06 may activate only after P01.05 reaches `done` with required G0/G1/G2/G3/G5/G6/G7 GitHub-hosted evidence.

## Future web UI quality/accessibility plan

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` now records the future browser-UI implementation requirement for humans and AI systems.

When an authorized UI package begins, browser surfaces target **WCAG 2.2 AA**, standards-based semantic HTML/CSS, W3C validation and WAVE evaluation plus manual keyboard/focus/screen-reader/zoom-reflow evidence. WAVE is an evaluation input, not an accessibility certification, and a missing required WAVE API/license is `BLOCKED`, not PASS.

This planning requirement does **not** authorize P12/P13/P17 or any other business/UI implementation during P01.

## Protected integration / CI

`main.protected=true` remains verified. PR-only integration, required `governance`, failed-check rejection, direct/force update rejection and conversation resolution remain enforced.

Canonical governance CI remains GitHub-hosted only on `ubuntu-24.04`; local/self-hosted runners are prohibited. Completed P01 verifiers remain regression gates in the same required `governance` job.

## Dependabot

Repository-managed `.github/dependabot.yml` is present for weekly GitHub Actions dependency updates. Dependency updates remain governed through normal protected PRs.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block authorized P01 kernel engineering.

## Exact next work

Implement **P01.05 only**, obtain GitHub-hosted cache contract/integration/outage/resilience/security/build evidence, reconcile P01.05 to `done`, then activate only P01.06. Keep `business_feature_code_authorized=false` and P02+ implementation locked.
