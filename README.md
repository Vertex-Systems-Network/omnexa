# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is ACTIVE. P01.01-P01.03 are DONE; P01.04 — PostgreSQL connection & migration foundation is the sole active package.** `kernel_code_authorized=true`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Relevant sources:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.04.md`
- `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`
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

Completed package evidence:

- P01.01: PR #40, run `32562869345`, job `97007065640`, merge `7257977264d788663083fa215462b1828f1e5afb`.
- P01.02: PR #42, run `32563880800`, job `97009520624`, merge `c857bb9e7df1e347226653eeaded024d6ecd0271`.
- P01.03: PR #44, run `32565935613`, job `97014452248`, merge `bdeda5fad09a2369b2a6852e5c62550db50047ea`.

Repository-managed Dependabot configuration is present for weekly GitHub Actions dependency updates.

## P01 package status

### P01.01 — done

Pinned Go workspace/build skeleton and deterministic hosted verification. Evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`.

### P01.02 — done

Typed configuration/environment system with deterministic precedence, strict unknown-key handling, secret-safe redaction/provenance, race-tested isolated loader state and fail-closed startup validation. Evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

### P01.03 — done

Transport-neutral structured failure primitives with stable codes/categories, safe public projection, private causes, Go wrapping semantics, explicit retryability, bounded validation details and correlation metadata hooks. Evidence: `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`.

### P01.04 — active

P01.04 implements only the PostgreSQL connection and migration foundation: governed pool construction from P01.02 configuration, safe P01.03 provider failure mapping, transaction helper boundary, migration runner/version ledger, deterministic fresh/upgrade execution and synthetic PostgreSQL integration tests.

It must **not** pull forward tenant/organization tables, module-runtime schema, event outbox/inbox, business schemas/data, cache/storage, telemetry, health endpoints, production HA/backups or business-domain repositories.

P01.05–P01.12 remain planned and strict sequential activation applies.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate. It does not block the currently authorized P01 kernel engineering scope.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current canonical state: **P00 done; P01 active; P01.01-P01.03 done; P01.04 active; P01 progress 3 / 12 done; kernel implementation authorized only for P01.04; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
