# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is ACTIVE. P01.01-P01.09 are DONE; P01.10 — Feature flag & configuration registry is the sole active package.** `kernel_code_authorized=true`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.10.md`
- `docs/roadmap/evidence/P01.09_COMPLETION_2026-08-22.md`
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

## Protected GitHub integration and executable CI

Issue #3 is closed and `main` is protected with PR-only integration, strict required `governance`, blocked direct/force updates, failed-check merge rejection, required conversation resolution and up-to-date branch enforcement.

Canonical required CI uses GitHub-hosted `ubuntu-24.04` only and fails closed unless the runner is GitHub-hosted Linux/X64. Local/self-hosted governance runners are prohibited. Permanent repository-wide Go quality runs before package regressions through `bash scripts/verify_go_quality.sh`, using pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

## P01 package status

P01.01 through P01.08 are completed with canonical evidence under `docs/roadmap/evidence/`.

### P01.09 — done

P01.09 established deterministic process-local job registration/execution, UUIDv7 execution IDs, bounded worker/queue concurrency, bounded retry/backoff with explicit idempotency protection, repeatable completion handles, graceful queued/synchronous drain/cancel semantics, UTC-normalized one-shot/fixed-interval schedules and safe P01.07 observability propagation. Final implementation evidence: PR #59, run `32605309150`, job `97109396616`, merge `0bcafbfc52324acba1df9d8eff84a264dda0f233`. Canonical evidence: `docs/roadmap/evidence/P01.09_COMPLETION_2026-08-22.md`.

### P01.10 — active

P01.10 owns only the governed runtime feature flag/configuration registry: typed definitions, stable identifiers and owner metadata, explicit deterministic defaults/fallbacks, runtime evaluation, future-scope-aware evaluation inputs without P02 identity, version/change metadata hooks, bounded refresh/invalidation, explicitly declared operational kill switches and a deterministic test provider.

P01.10 must not implement product experimentation/analytics, tenant admin UI, pricing/entitlement/licensing, authorization based solely on flags, business-module flags before their owners exist, P01.11/P01.12 behavior, P02+ implementation or AI/model/agent/planner functionality. Flags never grant authority or bypass authorization/data isolation; sensitive configuration remains governed by classification/secrets policy.

P01.11-P01.12 remain planned and strict sequential activation applies.

## Future browser UI quality/accessibility

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is the AI/human execution plan for future authorized browser UI work. It does not authorize business/UI implementation during P01.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate. It does not block currently authorized P01 kernel engineering.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current canonical state after the governed P01.09 closure transition: **P00 done; P01 active; P01.01-P01.09 done; P01.10 active; P01 progress 9 / 12 done; kernel implementation authorized only for P01.10; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
