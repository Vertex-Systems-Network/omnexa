# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is ACTIVE. P01.01-P01.05 are DONE; P01.06 — Object & file storage abstraction is the sole active package.** `kernel_code_authorized=true`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Relevant sources:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.06.md`
- `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`
- `docs/quality/GO_CODE_QUALITY.md`
- `docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`

## Core laws

- Kernel before business modules.
- One authoritative owner per write model/capability.
- Cross-module direct DB writes/private implementation imports are forbidden.
- Cross-domain communication uses governed APIs/capabilities/events/workflows/read projections.
- Tenant scope, authorization, audit, observability and contract versioning are mandatory.
- Optional modules fail/degrade independently.
- AI acts only through governed authorized capabilities; no unrestricted raw DB authority.
- Strict modular monolith first; service extraction requires evidence and ADR.
- Architecture/roadmap changes require change control and reconciliation.

## Foundation v1

P00 froze governance/AI/change control, terminology/ownership/dependencies, UUIDv7/exact-money/time/locale/error primitives, HTTP/OpenAPI and event contracts, security/data classification/tenant isolation, G0–G8 quality/release semantics, repository/local-development rules and threat/SLO/recovery/incident baselines.

Technology baseline remains Go + TypeScript/React with selective Rust/Python; PostgreSQL, Redis-compatible cache, S3-compatible storage, NATS/JetStream-class messaging and OpenTelemetry.

## Protected GitHub integration

Issue #3 is **closed/completed** and `main` is protected. Verified behavior includes PR-only integration, required `governance`, blocked direct/force updates, failed-check merge rejection and required conversation resolution. Current single-maintainer policy uses zero required approvals; increase independent review/Code Owner requirements when a second reviewer exists.

## Executable CI — GitHub-hosted only

Canonical required job:

```yaml
runs-on: ubuntu-24.04
```

The `governance` job fails unless `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. Local/self-hosted runners are prohibited.

Permanent repository-wide Go quality verification runs before package regressions through `bash scripts/verify_go_quality.sh`, using pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

Completed package evidence:

- P01.01: PR #40, run `32562869345`, job `97007065640`, merge `7257977264d788663083fa215462b1828f1e5afb`.
- P01.02: PR #42, run `32563880800`, job `97009520624`, merge `c857bb9e7df1e347226653eeaded024d6ecd0271`.
- P01.03: PR #44, run `32565935613`, job `97014452248`, merge `bdeda5fad09a2369b2a6852e5c62550db50047ea`.
- P01.04: PR #46, run `32567842071`, job `97019012280`, merge `6068202415dd124d3e74a196b6e0bbca5d75c4cd`.
- P01.05: PR #48, run `32571147128`, job `97026673348`, merge `725cbbd87e9456e1be02306ce3788e43ab139bd5`.

Repository-managed Dependabot configuration is present for weekly GitHub Actions dependency updates.

## P01 package status

### P01.01 — done

Pinned Go workspace/build skeleton and deterministic hosted verification. Evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`.

### P01.02 — done

Typed configuration/environment system with deterministic precedence, strict unknown-key handling, secret-safe redaction/provenance, race-tested isolated loader state and fail-closed startup validation. Evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

### P01.03 — done

Transport-neutral structured failure primitives with stable codes/categories, safe public projection, private causes, Go wrapping semantics, explicit retryability, bounded validation details and correlation metadata hooks. Evidence: `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`.

### P01.04 — done

PostgreSQL connection/migration foundation with pgx `v5.10.0`, PostgreSQL `18.6`, bounded pool/config behavior, safe provider failures, transaction helper, advisory-locked owner-scoped migration ledger, deterministic fresh/upgrade/failure tests and GitHub-hosted G0/G1/G2/G3/G4/G5/G7 evidence. Evidence: `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`.

### P01.05 — done

Redis-compatible cache foundation with Valkey `9.1.1` and valkey-go `v1.0.75`: deterministic namespaced/versioned keys, bounded TTL/value semantics, typed JSON serialization, explicit miss/error behavior, get/set/delete, safe provider failures, bounded cancellation/timeouts, FLUSHDB non-authority and provider restart/reconnect evidence. Evidence: `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`.

P01.05 also established the permanent Go code-quality gate. The canonical completion run reported golangci-lint `0 issues` and no reachable govulncheck vulnerabilities. A reachable `GO-2026-5970` finding discovered during development was fixed by moving the canonical dependency graph to `golang.org/x/text v0.39.0` rather than suppressing it.

### P01.06 — active

P01.06 implements only the governed S3-compatible object/file storage foundation: provider adapter, bucket/container configuration, explicitly namespaced/versioned object keys, streamed bounded upload/download, untrusted metadata handling, integrity/checksum hooks, put/get/head/delete/missing-object behavior, safe provider timeout/cancellation/error mapping and synthetic provider verification.

It must **not** pull forward media library/CMS/file UI, public URL/CDN/image processing, tenant/business document models, sessions/authentication, tenancy, module runtime, event/job fabric, retention/legal hold, telemetry/logging, health, scheduler, feature registry, audit transport or later capability semantics.

Object keys never imply authorization. Path traversal must not escape the provider namespace. Storage credentials are RESTRICTED; signed URLs, secrets and sensitive content must not appear in diagnostics. Large-object flows must be streamed and bounded.

P01.07–P01.12 remain planned and strict sequential activation applies.

## Go code quality

`docs/quality/GO_CODE_QUALITY.md` defines the permanent repository-wide Go quality policy. The required governance job validates formatting, curated static analysis and reachable dependency vulnerabilities with pinned tool versions before completed-package regression verifiers.

## Future browser UI quality/accessibility

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is the AI/human execution plan for future authorized browser UI work. Production UI targets **WCAG 2.2 AA**, semantic standards-based HTML/CSS, W3C validation, WAVE evaluation and manual keyboard/focus/screen-reader/zoom-reflow checks.

WAVE is treated as an evaluation tool rather than a certification. Required WAVE automation that lacks an approved key/license is `BLOCKED`, not PASS. This plan does not authorize P12/P13/P17 or any business/UI implementation during P01.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate. It does not block the currently authorized P01 kernel engineering scope.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current canonical state: **P00 done; P01 active; P01.01-P01.05 done; P01.06 active; P01 progress 5 / 12 done; kernel implementation authorized only for P01.06; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
