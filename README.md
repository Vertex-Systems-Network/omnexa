# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P02 — Identity, Tenancy & Organization is ACTIVE at 8 / 10 done. P02.01-P02.08 are DONE and P02.09 — Tenant-Scoped Settings is the sole active work package.** `kernel_code_authorized=true` only for P02.09 after this closure merges and protected-main state is verified; `business_feature_code_authorized=false`.

## Project progress

```text
P00  Product Constitution & Architecture Freeze  [██████████] 10/10  DONE
P01  Omnexa Kernel                               [██████████] 12/12  DONE
P02  Identity, Tenancy & Organization            [████████░░]  8/10  ACTIVE
      └─ Current: P02.09 — Tenant-Scoped Settings
      └─ Next:    P02.10 — Identity/permission audit trails & P02 exit proof
P03+ Future phases                               [░░░░░░░░░░]        PLANNED / LOCKED
```

The bars report only comparable package completion **inside each governed phase**. Omnexa does not publish a synthetic overall roadmap percentage across P00-P27 because later phases contain unequal scope and would make a single percentage misleading. The authoritative execution cursor remains `docs/roadmap/STATE.json`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_ENTRY_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P02.09.md`
- `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.02_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.03_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.05_COMPLETION_2026-08-24.md`
- `docs/roadmap/evidence/P02.06_COMPLETION_2026-08-24.md`
- `docs/roadmap/evidence/P02.07_COMPLETION_2026-08-24.md`
- `docs/roadmap/evidence/P02.08_COMPLETION_2026-08-25.md`
- `docs/quality/GO_CODE_QUALITY.md`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`

For strategic architecture/roadmap work, also inspect the proposed AI-native expansion under `docs/adr/ADR-0011-ai-native-business-os-strategic-expansion.md`, `docs/roadmap/STRATEGIC_PROGRAMS.json`, `docs/roadmap/STRATEGIC_CROSS_CUTTING_PROGRAMS.md`, `docs/roadmap/STRATEGIC_ACCEPTANCE_GATES.md`, `docs/roadmap/AI_NATIVE_STRATEGIC_AUDIT_2026.md`, `docs/architecture/AI_NATIVE_BUSINESS_OS_ARCHITECTURE.md` and `docs/architecture/PRODUCT_FEDERATION_AND_APP_MESH.md`. Proposed strategic files do **not** activate future implementation by themselves.

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

## P02.01-P02.08 completion

P02.01 completed through implementation PR #69, exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical GitHub-hosted run/job `32635243643 / 97183883007`, and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 completed through implementation PR #71, exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical GitHub-hosted run/job `32637760875 / 97189971101`, and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`.

P02.03 completed through implementation PR #73, exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`, canonical GitHub-hosted run/job `32640790333 / 97197453122`, and merge `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`.

P02.04 completed through implementation PR #75, exact head `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f`, canonical GitHub-hosted run/job `32653747461 / 97229198036`, and merge `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`, evidence `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`.

P02.05 completed through implementation PR #77, exact head `2df8d2a8bef0cea60256a832986d6f8495c80378`, canonical GitHub-hosted run/job `32660848145 / 97246683239`, and merge `7b6a59e83c9bd696e6e008385b4413d529254171`, evidence `docs/roadmap/evidence/P02.05_COMPLETION_2026-08-24.md`.

P02.06 completed through implementation PR #79, exact head `dbbd105fd5f2543ca7dd5df93375eaf1057928fc`, canonical GitHub-hosted run/job `32664834112 / 97256520050`, and merge `083c2866f0cd0773b85201750c2196bfd2fcc167`, evidence `docs/roadmap/evidence/P02.06_COMPLETION_2026-08-24.md`.

P02.07 completed through implementation PR #81, exact head `51ccaa12c3534f74fba6eab9d4698ee483ef4ffd`, canonical GitHub-hosted run/job `32669167972 / 97267175953`, and merge `5642f5da1eb24e70b67e5ec757d9f4584c4e3f5c`, evidence `docs/roadmap/evidence/P02.07_COMPLETION_2026-08-24.md`.

P02.08 completed through implementation PR #84, exact head `43bdcf525ce5e0cfdb9dc0707fbafee7cd552543`, canonical GitHub-hosted run/job `32885950897 / 97926598423`, and merge `32eb7187eb229327585551e4e28b0d596de78bd9`.

P02.08 canonical evidence passed repository Go quality, P01.01-P01.12 regressions, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.07 regressions and applicable P02.08 G0-G8. It proves distinct non-human Service Account lifecycle; exact tenant/organization credential binding; SHA-256 verifier-only credential persistence; current-state verification; transactional rotation; revocation/expiry/supersession denial; direct RBAC composition; raw-secret redaction; and fresh/idempotent/P02.07+P02.05 supported-upgrade migration evidence. Immutable evidence is `docs/roadmap/evidence/P02.08_COMPLETION_2026-08-25.md`.

Diagnostic implementation runs remain recorded as FAIL rather than being hidden: `32882746486 / 97915911717`, `32884311341 / 97921359088`, and `32884939579 / 97923224921`. Run `32885758158` was superseded/cancelled and is not acceptance evidence. The final canonical lane validates the narrow historical-verifier compatibility fixes without weakening retained P02.04-P02.07 regressions.

## Active P02.09 scope

P02.09 owner is `kernel.configuration`. Authorized scope is limited to tenant-scoped and approved organization-scoped setting resolution using the completed P01 configuration registry; trusted scope derived from P02 identity/tenancy context rather than arbitrary payload identifiers; authorization around protected setting reads/writes; classification-aware values and no-secret output behavior; deterministic precedence only for explicitly supported scopes; change audit hooks for security-significant settings; and same-tenant allow plus cross-tenant/wrong-scope negative evidence.

Settings and feature flags cannot create authority by themselves. Tenant/org scope is trusted context, not a client assertion. Values inherit normal data-classification and logging restrictions. `kernel.configuration` remains authoritative owner, with no cross-tenant fallback or global-write shortcut.

Business-module settings, P03 module runtime, a secrets-management product surface, feature/config values that independently grant authority, deployment/environment orchestration, P02.10 implementation and other future scope remain unauthorized.

## Current implementation lock

- `kernel_code_authorized=true` only for P02.09 after this closure merges and protected-main state is verified.
- `business_feature_code_authorized=false`.
- P02.10 remains planned.
- P03+ remains planned.
- All proposed `X..` strategic programs remain planning-only until accepted, dependency-ready and activated through canonical work-package/state governance.

## Strategic AI-native direction

The proposed strategic layer does not replace P00-P27. It adds cross-cutting intelligence/control systems so Omnexa can evolve beyond feature aggregation into:

**deterministic domain truth + universal workflow + governed data/Business Graph + Process Graph + System Graph + AI Control Tower + model governance + simulation + continuous controls + measurable outcomes + constrained autonomy.**

Standalone first-party products should join through `XPF-200 Product Federation & App Mesh` as native, embedded, federated or edge products according to architecture evidence rather than being force-merged into one codebase/database.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current checkpoint: **P00 done 10 / 10; P01 done 12 / 12; P02 active 8 / 10; P02.01-P02.08 done; P02.09 sole active package; P02.10 planned; business features locked.**

`docs/roadmap/STRATEGIC_PROGRAMS.json` records proposed cross-cutting X-programs. It does not create a second execution cursor or override `STATE.json`.

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
