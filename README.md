# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P02 — Identity, Tenancy & Organization is ACTIVE at 4 / 10 done. P02.01-P02.04 are DONE and P02.05 — RBAC foundation is the sole active work package.** `kernel_code_authorized=true` only for P02.05; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_ENTRY_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P02.05.md`
- `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.02_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.03_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`
- `docs/quality/GO_CODE_QUALITY.md`
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

P01.01-P01.12 are complete with canonical executable evidence. Final P01.12 evidence: PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. P01 regressions remain mandatory during P02.

## P02.01-P02.04 completion

P02.01 completed through implementation PR #69, exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical GitHub-hosted run/job `32635243643 / 97183883007`, and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 completed through implementation PR #71, exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical GitHub-hosted run/job `32637760875 / 97189971101`, and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`.

P02.03 completed through implementation PR #73, exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`, canonical GitHub-hosted run/job `32640790333 / 97197453122`, and merge `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`.

P02.04 completed through implementation PR #75, exact head `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f`, canonical GitHub-hosted run/job `32653747461 / 97229198036`, and merge `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`.

P02.04 canonical evidence passed repository Go quality, P01.01-P01.12 regressions, `omnexa verify all`, P02.01-P02.03 regressions and P02.04 G0-G8, including real PostgreSQL fresh/idempotent/P02.01-upgrade migration evidence; governed adaptive password hashing; disclosure-safe authentication; opaque digest-only session credential persistence; explicit expiry; rotation/replay rejection/revocation; password/account lifecycle invalidation; safe session inventory; and current tenant/organization context reauthorization. Evidence is retained in `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`.

Diagnostic run `32653382151 / 97228307755` remains FAIL for four corrected `gosec` G101 findings and one corrected staticcheck simplification. It is not completion evidence.

## Active P02.05 scope

P02.05 owner is `kernel.authorization`. Authorized scope is limited to stable capability-oriented permission identifiers; Role as deterministic permission composition rather than bypass identity; tenant/organization-scoped role definitions and assignments; assignment/revocation lifecycle; direct deny-by-default RBAC evaluation; privileged role/permission mutation checks and classification-safe audit hooks; and allow/deny/privilege-escalation negative tests.

`admin`, `owner`, `superuser` or any role name never means platform bypass. Role assignment is itself privileged and auditable. Tenant/organization identifier substitution cannot widen permission authority.

P02.06 relationship/context policy, P02.07 MFA/passkeys, P02.08 service accounts/API credential scope, P03 module permission registration, business-module permissions, hidden super-admin logic, business behavior, deployment authority and AI/model/agent runtime remain unauthorized.

## Current implementation lock

- `kernel_code_authorized=true` only for P02.05.
- `business_feature_code_authorized=false`.
- P02.06-P02.10 remain planned.
- P03+ remains planned.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current checkpoint: **P00 done; P01 done 12 / 12; P02 active 4 / 10; P02.01-P02.04 done; P02.05 sole active package; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
