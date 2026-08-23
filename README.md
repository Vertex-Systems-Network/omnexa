# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P02 — Identity, Tenancy & Organization is ACTIVE at 0 / 10 done. P02.01 — Principal & user identity foundation is the sole active work package.** `kernel_code_authorized=true` only for P02.01; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_ENTRY_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/governance/P01_P02_TRANSITION_CHECKLIST.md`
- `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P02.01.md`
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

Canonical required CI uses GitHub-hosted `ubuntu-24.04` only and fails closed unless the runner is GitHub-hosted Linux/X64. Local/self-hosted governance runners are prohibited. Permanent repository-wide Go quality runs through `bash scripts/verify_go_quality.sh` using pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

## P01 completion retained

P01.01-P01.12 are complete with canonical executable evidence. Final P01.12 evidence: implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. `docs/governance/P01_EXIT_GATE.md` remains **SATISFIED**, and all P01 regressions remain mandatory during P02.

## P02 readiness and active scope

P02 readiness preparation merged through PR #67 as `c6301ca4a5eec5dd62bcb75481d900e40ad968bd` after final exact-head canonical run `32632920772 / 97178312240` passed.

P02 uses a strict sequential 10-package sequence. Only P02.01 is active. Its owner is `kernel.identity`, and its authorized scope is the canonical human principal/User identity foundation with deterministic lifecycle semantics, safe identity attributes and owner-bounded persistence where required.

P02.01 explicitly does **not** authorize tenant lifecycle/membership authority, organization hierarchy, authentication/session behavior, RBAC/policy, MFA/passkeys, service-account/API credential lifecycle, P03 module runtime, business modules or AI/model/agent runtime. User remains distinct from business Person.

## Current implementation lock

- `kernel_code_authorized=true` only for P02.01.
- `business_feature_code_authorized=false`.
- P02.02-P02.10 remain planned.
- P03+ remains planned.
- AI/model/agent runtime implementation remains unauthorized.

## Future browser UI quality/accessibility

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is the AI/human execution plan for future authorized browser UI work. It does not itself authorize UI implementation.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current checkpoint: **P00 done; P01 done 12 / 12; P01 exit satisfied; P02 active 0 / 10; P02.01 sole active package; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
