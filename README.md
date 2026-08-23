# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is ACTIVE. P01.01-P01.11 are DONE; P01.12 — Developer CLI Baseline is the sole active package and owns the P01 exit proof.** `kernel_code_authorized=true`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.12.md`
- `docs/roadmap/evidence/P01.11_COMPLETION_2026-08-23.md`
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

### P01.11 — done

P01.11 implemented the governed `kernel.audit` transport foundation: immutable classification-aware UUIDv7/UTC audit records, tamper-evident integrity, descriptive actor/action/target/scope/outcome/correlation/reason/approval and privileged/impersonation metadata without P02 authority, append-only sink capability, explicit required-audit failure/best-effort degradation semantics, bounded deterministic memory sink, prohibited-secret handling and protected-payload-safe transport health.

Final implementation evidence: PR #63, exact head `1c1ab1f8d5120fb6b1e5908fdb93cffef9275940`, run `32610902537`, job `97123708250`, merge `10c94a638b89d47da05f5481fb2db298a2da6942`. Canonical evidence: `docs/roadmap/evidence/P01.11_COMPLETION_2026-08-23.md`.

Audit write capability does not imply read/export authority, and actor/scope metadata does not grant identity, tenancy or authorization authority.

### P01.12 — active

P01.12 owns only the governed developer CLI baseline and P01 exit proof: deterministic help/version/verify behavior, fail-closed verification orchestration, explicit structured-safe output/exit codes, safe composition of existing P01 configuration/migration/diagnostic capabilities, no-secret/no-RESTRICTED output, clean-checkout/CI reproducibility and the complete fresh-install P01 exit path.

P01.12 must not implement production super-admin authority, P02 tenant/user/role administration, P03 module runtime administration, P04+ domain/event/workflow commands, deployment/Kubernetes orchestration, hidden SQL/file mutation, business modules or AI/model/agent behavior.

P01 progress is **11 / 12 done**. P02 remains planned until P01.12 completes and P01 exit is reconciled through a separate governed transition.

## Future browser UI quality/accessibility

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is the AI/human execution plan for future authorized browser UI work. It does not authorize business/UI implementation during P01.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate. It does not block currently authorized P01 kernel engineering.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current closure state: **P00 done; P01 active; P01.01-P01.11 done; P01.12 active; P01 progress 11 / 12 done; kernel implementation authorized only for P01.12; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
