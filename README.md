# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is ACTIVE. P01.01-P01.08 are DONE; P01.09 — Job & scheduler primitives is the sole active package.** `kernel_code_authorized=true`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.09.md`
- `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.06_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.07_COMPLETION_2026-08-22.md`
- `docs/roadmap/evidence/P01.08_COMPLETION_2026-08-22.md`
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
- AI acts only through governed authorized capabilities; no unrestricted raw DB/object-store/business-state authority.
- Strict modular monolith first; service extraction requires evidence and ADR.
- Architecture/roadmap changes require change control and reconciliation.

## Protected GitHub integration

Issue #3 is closed and `main` is protected. Verified behavior includes PR-only integration, strict required `governance`, blocked direct/force updates, failed-check merge rejection, required conversation resolution and up-to-date branch enforcement.

P01.06 PR #50 reconfirmed strict up-to-date enforcement: prior green runs did not permit merge while the branch lagged `main`. The branch was synchronized, fresh full governance run `32588244996` passed, and the PR merged normally as `f7867d9e1c570e3abbed90740970acf7b5a30bd7`.

## Executable CI — GitHub-hosted only

Canonical required job uses:

```yaml
runs-on: ubuntu-24.04
```

The `governance` job fails unless the runner is GitHub-hosted Linux/X64. Local/self-hosted governance runners are prohibited.

Permanent repository-wide Go quality verification runs before package regressions through `bash scripts/verify_go_quality.sh`, using pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

## P01 package status

### P01.01 — done
Pinned Go workspace/build skeleton. Evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`.

### P01.02 — done
Typed configuration/environment system. Evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

### P01.03 — done
Structured failure/error conventions. Evidence: `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`.

### P01.04 — done
PostgreSQL connection/migration foundation with pgx `v5.10.0` and PostgreSQL `18.6`. Evidence: `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`.

### P01.05 — done
Redis-compatible cache foundation with Valkey `9.1.1` / valkey-go `v1.0.75`, plus the permanent Go code-quality gate. Evidence: `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`.

### P01.06 — done
Governed S3-compatible object/file storage foundation with deterministic namespaced/versioned keys, bounded streaming, untrusted metadata validation, SHA-256 integrity handling, provider-safe failures, timeout/cancellation/unavailability classification, concurrent synthetic integration and provider restart evidence. Final implementation evidence: PR #50, run `32588244996`, job `97067784835`, merge `f7867d9e1c570e3abbed90740970acf7b5a30bd7`. Canonical evidence: `docs/roadmap/evidence/P01.06_COMPLETION_2026-08-22.md`.

### P01.07 — done
Governed structured logging and OpenTelemetry baseline using standard-library `log/slog`, safe classification/sensitive-key redaction, bounded correlation + W3C trace-context propagation, isolated trace/metric SDK lifecycle, vendor-neutral exporter injection, bounded flush/shutdown, disabled mode and deterministic capture tests. OpenTelemetry Go is pinned to `v1.45.0`. Final implementation evidence: PR #54, run `32595413156`, job `97085413083`, merge `245fd67d60c02a5ed546f0cd4dc934345b5b4d42`. Canonical evidence: `docs/roadmap/evidence/P01.07_COMPLETION_2026-08-22.md`.

### P01.08 — done
Governed portable health/readiness/diagnostics foundation with distinct process liveness and dependency readiness, required/optional/security-critical classifications, bounded timeout/cancellation/panic handling, startup/ready/stopping/failed lifecycle semantics, safe machine-readable diagnostics, build identity and P01.07 observability integration. Final implementation evidence: PR #56, run `32601741049`, job `97100949202`, merge `a2a93454bb283464bfa144bb5a38539041e40069`. Canonical evidence: `docs/roadmap/evidence/P01.08_COMPLETION_2026-08-22.md`.

### P01.09 — active
P01.09 owns only minimal kernel-local background job and scheduling primitives: deterministic job registration/execution, bounded worker concurrency, cancellation/deadline propagation, bounded retry/backoff, explicit idempotency hooks, one-shot/recurring maintenance schedules, observability propagation and a deterministic in-memory/test harness.

P01.09 must not pull forward NATS/JetStream durable streams/consumers, transactional outbox/inbox, distributed workflow timers, business jobs, tenant-context implementation, P01.10+ behavior, P02+ implementation or AI/model/agent/planner functionality.

P01.10-P01.12 remain planned and strict sequential activation applies.

## Future browser UI quality/accessibility

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is the AI/human execution plan for future authorized browser UI work. It does not authorize P12/P13/P17 or business/UI implementation during P01.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate. It does not block currently authorized P01 kernel engineering.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current canonical state after the P01.08 closure transition: **P00 done; P01 active; P01.01-P01.08 done; P01.09 active; P01 progress 8 / 12 done; kernel implementation authorized only for P01.09; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
