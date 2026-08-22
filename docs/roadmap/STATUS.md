# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.06 — Object & file storage abstraction**
- P01 progress: **5 / 12 done**
- P01.01: **DONE**
- P01.02: **DONE**
- P01.03: **DONE**
- P01.04: **DONE**
- P01.05: **DONE**
- P01.06: **ACTIVE**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**
- P01.07–P01.12: **PLANNED / NOT ACTIVE**

## Completed P01 packages

### P01.01 — Go workspace/build skeleton

Canonical evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`.

PR #40 merged as `7257977264d788663083fa215462b1828f1e5afb` after GitHub-hosted run `32562869345` / job `97007065640` passed the pinned Go, format/static, unit, dependency-boundary and build/smoke gates.

### P01.02 — Configuration & environment system

Canonical evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

PR #42 merged as `c857bb9e7df1e347226653eeaded024d6ecd0271` after corrected GitHub-hosted run `32563880800` / job `97009520624` passed P01.01 regression, format/static, unit/race, precedence/startup, secret-safe negative, dependency-boundary and build/package gates. Initial run `32563817196` remains failure evidence and was fixed rather than bypassed.

### P01.03 — Structured error & result conventions

Canonical evidence: `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`.

PR #44 merged as `bdeda5fad09a2369b2a6852e5c62550db50047ea` after GitHub-hosted run `32565935613` / job `97014452248` passed P01.01/P01.02 regressions, Go `1.26.7`, format/static, unit/race, wrapping, public/private redaction, validation bounds/order, dependency-boundary and build/package. Earlier diagnostic failures remain preserved and are not counted as PASS.

### P01.04 — PostgreSQL connection & migration foundation

Canonical evidence: `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`.

PR #46 was verified by GitHub-hosted run `32567842071` / job `97019012280` and squash-merged as `6068202415dd124d3e74a196b6e0bbca5d75c4cd`.

Verified baseline: Go `1.26.7`, pgx `v5.10.0`, PostgreSQL `18.6 (Debian 18.6-1.pgdg13+2)`, P01.01-P01.03 regressions PASS, P01.04 G1/G2/G3/G4/G5/G7 PASS and G0 governance/state validators PASS. Synthetic evidence covered ping/version, bounded pool exhaustion, unavailable connection, transaction commit/rollback, fresh/idempotent/upgrade migrations, failed-migration rollback without ledger advance, migration drift and advisory-lock contention.

### P01.05 — Cache abstraction

Canonical evidence: `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`.

PR #48 final source head `c39ba2f0bee8da7c185608f8ca30d6090dc4e824` passed canonical GitHub-hosted run `32571147128` / job `97026673348` and squash-merged to protected `main` as `725cbbd87e9456e1be02306ce3788e43ab139bd5`.

Verified baseline:

- Go `1.26.7`;
- Valkey `9.1.1`;
- valkey-go `v1.0.75`;
- P01.01-P01.04 regressions PASS;
- P01.05 G1/G2/G3/G5/G6/G7 PASS;
- deterministic namespaced/versioned keys, bounded TTL/value size, JSON serialization, miss-vs-provider-error behavior, get/set/delete, timeout/cancellation, FLUSHDB non-authority and provider restart/reconnect PASS.

P01.05 also established the permanent repository-wide Go quality gate. On the final run, `golangci-lint v2.12.2` reported **0 issues**, `govulncheck v1.7.0` reported **no reachable vulnerabilities**, and `gofmt` passed. During development the new gate found reachable `GO-2026-5970` through `golang.org/x/text v0.36.0`; the canonical module graph was corrected to `golang.org/x/text v0.39.0` and `golang.org/x/sync v0.21.0` rather than suppressing the finding.

Diagnostic/failure runs `32569193459`, `32570210066`, `32570578340` and intentionally failing metadata diagnostic `32571033650` remain preserved as non-PASS history.

P01.05 did not introduce business cache models, sessions/authentication, tenancy, module runtime, event/job fabric, object storage, telemetry, health or later capability implementation.

## Active P01 package — P01.06

P01.06 is the sole active package. It owns only the S3-compatible object/file storage abstraction foundation:

- S3-compatible provider interface/adapter;
- bucket/container configuration boundary;
- deterministic object key namespace/version conventions;
- streaming upload/download with bounded memory behavior;
- untrusted content length/type/file metadata handling;
- integrity/checksum hooks;
- put/get/head/delete and missing-object behavior;
- bounded provider timeout/cancellation/error mapping through completed P01.03 failures;
- P01.02 configuration integration;
- synthetic S3-compatible contract/integration evidence;
- P01.01-P01.05 regressions and repository Go quality preservation.

Explicitly prohibited in P01.06: media library/CMS/file UI, public URL/CDN/image processing, tenant document/business models, sessions/authentication, tenancy implementation, module runtime, event/job fabric, malware-scanning implementation beyond hook boundaries, retention/legal-hold semantics, logging/OpenTelemetry, health, scheduler, feature registry, audit transport or any later P01/P02+ capability.

Object keys never imply authorization. Path traversal must not escape the provider namespace. Storage credentials are RESTRICTED. Signed URLs, secrets and sensitive object content must not appear in public failures/logs/artifacts. Large-object flows must be streamed and bounded.

P01.07 may activate only after P01.06 reaches `done` with required G0/G1/G2/G3/G5/G6/G7 GitHub-hosted evidence.

## Repository Go quality

`docs/quality/GO_CODE_QUALITY.md` defines the permanent Go quality gate. Canonical `governance` executes `bash scripts/verify_go_quality.sh` before package regressions. Tool versions are pinned and the gate fails closed for formatting, configured lint/static findings or reachable Go vulnerabilities.

## Future web UI quality/accessibility plan

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` records the future browser-UI implementation requirement for humans and AI systems.

When an authorized UI package begins, browser surfaces target **WCAG 2.2 AA**, standards-based semantic HTML/CSS, W3C validation and WAVE evaluation plus manual keyboard/focus/screen-reader/zoom-reflow evidence. WAVE is an evaluation input, not an accessibility certification, and a missing required WAVE API/license is `BLOCKED`, not PASS.

This planning requirement does **not** authorize P12/P13/P17 or any other business/UI implementation during P01.

## Protected integration / CI

`main.protected=true` remains verified. PR-only integration, required `governance`, failed-check rejection, direct/force update rejection and conversation resolution remain the canonical integration controls.

Canonical governance CI remains GitHub-hosted only on `ubuntu-24.04`; local/self-hosted runners are prohibited. Completed P01 verifiers remain regression gates in the same required `governance` job.

## Dependabot

Repository-managed `.github/dependabot.yml` is present for weekly GitHub Actions dependency updates. Dependency updates remain governed through normal protected PRs.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block authorized P01 kernel engineering.

## Exact next work

Implement **P01.06 only**, establish a fail-closed P01.06 verifier and governed S3-compatible synthetic provider, obtain GitHub-hosted object storage contract/streaming/integrity/outage/security/build evidence, reconcile P01.06 to `done`, then activate only P01.07. Keep `business_feature_code_authorized=false` and P02+ implementation locked.
