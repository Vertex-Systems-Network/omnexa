# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is ACTIVE. P01.01 — Go workspace/build skeleton is the sole active package.** `kernel_code_authorized=true`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Relevant entry/sequence sources:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.01.md`
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

The `governance` job fails unless `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. Local/self-hosted runners are prohibited. PR #37 run `32541439589` is current positive hosted evidence.

Repository-managed Dependabot configuration is present for weekly GitHub Actions dependency updates.

## P01.01 — active scope

P01.01 establishes only the reproducible Go workspace/build skeleton: repository-owned Go toolchain declaration, workspace/module layout, minimal kernel process, build metadata and deterministic format/vet/test/build/smoke commands.

It must **not** implement configuration, persistence/migrations, cache/storage, telemetry, health endpoints, jobs, feature flags, audit transport, full developer CLI, identity/tenancy, module runtime or any business-domain behavior. Those belong to later packages/phases.

P01.02–P01.12 remain planned and strict sequential activation applies.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate. It does not block the currently authorized P01 kernel engineering scope.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00–P27. Current canonical state: **P00 done; P01 active; P01.01 active; P01 progress 0 / 12 done; kernel implementation authorized only for P01.01; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
