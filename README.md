# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is DONE at 12 / 12. P01 exit is SATISFIED. No work package is active. P02 remains PLANNED / NOT ACTIVE.** `kernel_code_authorized=false`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.12.md`
- `docs/roadmap/evidence/P01.12_COMPLETION_2026-08-23.md`
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

## P01 completion

P01.01-P01.12 are complete with canonical executable evidence. The final P01.12 implementation delivered the bounded `kernel.developer` CLI baseline and P01 fresh-install exit proof:

- deterministic `help` and `version`;
- safe structured `health` diagnostics;
- guarded non-production `db migrate`;
- fail-closed `verify <target>` orchestration using existing governed validators/verifiers;
- exact subprocess command allowlisting without shell-string expansion;
- runtime configuration isolation from verification subprocesses;
- real canonical `verify all` execution;
- focused positive/negative/race checks;
- P01 G0-G8 exit evidence.

Final P01.12 evidence: implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. See `docs/roadmap/evidence/P01.12_COMPLETION_2026-08-23.md` and `docs/governance/P01_EXIT_GATE.md`.

The P01 exit gate proves clean configuration/startup, fresh PostgreSQL migration, cache/storage provider contracts, safe telemetry, readiness/diagnostics, jobs/configuration/audit primitives, developer verification, security/static checks, module checksum verification and build/package behavior reproducibly without hidden manual steps.

## Current implementation lock

No product-development phase is active. P02 is planned but not activated. Until a separate governed P02 readiness/activation transition merges:

- kernel/P02 implementation is not authorized;
- business-feature implementation is not authorized;
- P03+ implementation is not authorized;
- AI/model/agent runtime implementation is not authorized.

Governance/specification/readiness preparation may proceed only within the allowance recorded in `STATE.json`.

## Future browser UI quality/accessibility

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is the AI/human execution plan for future authorized browser UI work. It does not itself authorize UI implementation.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current checkpoint: **P00 done; P01 done 12 / 12; P01 exit satisfied; P02 planned/not active; both implementation locks false.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
