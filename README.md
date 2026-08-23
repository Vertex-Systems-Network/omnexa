# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P01 — Omnexa Kernel is ACTIVE. P01.01-P01.10 are DONE; P01.11 — Audit Transport Foundation is the sole active package.** `kernel_code_authorized=true`; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_ENTRY_GATE.md`
- `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P01.11.md`
- `docs/roadmap/evidence/P01.10_COMPLETION_2026-08-23.md`
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

### P01.10 — done

P01.10 implemented the governed runtime feature flag/configuration registry: typed definitions, stable owner/version metadata, deterministic defaults/provider/fallbacks, future UUIDv7 scope references as opaque metadata only, bounded non-authoritative cache/refresh/invalidation, value-free change metadata, explicit fail-closed disable-only kill switches and a deterministic test provider. Flags do not grant authority or establish tenancy, and the registry is not a secrets store.

Final implementation evidence: PR #61, exact head `4c9914e4641d0d6e94a895d0fcd16c3a6bf4d962`, run `32609018028`, job `97118796940`, merge `9d11b9250eb74ca2ade531ee58e8f905468cf103`. Canonical evidence: `docs/roadmap/evidence/P01.10_COMPLETION_2026-08-23.md`.

### P01.11 — active

P01.11 owns only the governed `kernel.audit` transport foundation: a stable classification-aware audit envelope, actor/action/target/scope/outcome/correlation/reason/approval metadata without P02 identity, append-oriented sink semantics, explicit required-audit failure behavior, classification/redaction enforcement, UUIDv7/timestamp immutability, impersonation/privileged metadata representation, deterministic local/test sink and protected-payload-safe transport health observability.

P01.11 must not implement P02 identity/tenant/role catalogs, business audit event definitions, compliance/reporting UI, legal retention/hold systems, P01.12 CLI, durable messaging/outbox/inbox pull-forward, later business modules or AI/model/agent behavior. Audit write capability does not imply audit read/export authority.

P01.12 remains planned and strict sequential activation applies.

## Future browser UI quality/accessibility

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is the AI/human execution plan for future authorized browser UI work. It does not authorize business/UI implementation during P01.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate. It does not block currently authorized P01 kernel engineering.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current closure state: **P00 done; P01 active; P01.01-P01.10 done; P01.11 active; P01 progress 10 / 12 done; kernel implementation authorized only for P01.11; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
