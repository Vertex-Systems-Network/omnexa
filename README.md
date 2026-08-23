# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Omnexa Foundation Architecture v1 is **FROZEN** and P00 is **DONE**.

> **Current execution state:** **P02 — Identity, Tenancy & Organization is ACTIVE at 6 / 10 done. P02.01-P02.06 are DONE and P02.07 — MFA/passkey-ready flows is the sole active work package.** `kernel_code_authorized=true` only for P02.07; `business_feature_code_authorized=false`.

## Mandatory contributor / AI start here

Read `AGENTS.md` first. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Durable AI continuation starts with `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant handoff after canonical state is verified.

Key references:

- `docs/governance/FOUNDATION_FREEZE.json`
- `docs/governance/P01_EXIT_GATE.md`
- `docs/governance/P02_ENTRY_GATE.md`
- `docs/governance/P02_EXIT_GATE.md`
- `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`
- `docs/roadmap/work-packages/P02.07.md`
- `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.02_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.03_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`
- `docs/roadmap/evidence/P02.05_COMPLETION_2026-08-24.md`
- `docs/roadmap/evidence/P02.06_COMPLETION_2026-08-24.md`
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

## P02.01-P02.06 completion

P02.01 completed through implementation PR #69, exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical GitHub-hosted run/job `32635243643 / 97183883007`, and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 completed through implementation PR #71, exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical GitHub-hosted run/job `32637760875 / 97189971101`, and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`.

P02.03 completed through implementation PR #73, exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`, canonical GitHub-hosted run/job `32640790333 / 97197453122`, and merge `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`.

P02.04 completed through implementation PR #75, exact head `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f`, canonical GitHub-hosted run/job `32653747461 / 97229198036`, and merge `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`.

P02.05 completed through implementation PR #77, exact head `2df8d2a8bef0cea60256a832986d6f8495c80378`, canonical GitHub-hosted run/job `32660848145 / 97246683239`, and merge `7b6a59e83c9bd696e6e008385b4413d529254171`.

P02.06 completed through implementation PR #79, exact head `dbbd105fd5f2543ca7dd5df93375eaf1057928fc`, canonical GitHub-hosted run/job `32664834112 / 97256520050`, and merge `083c2866f0cd0773b85201750c2196bfd2fcc167`.

P02.06 canonical evidence passed repository Go quality, P01.01-P01.12 regressions, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.05 regressions and applicable P02.06 G0-G8. It proves P02.05 RBAC remains mandatory before contextual checks; trusted relationship evidence must match exact principal/object/tenant/organization scope; contextual constraints can narrow but never widen authority; internal/background caller origin cannot bypass policy; sensitive-field/export permissions remain distinct from ordinary read; material denials and privileged decisions use safe required audit; and resolver/evaluator failures fail closed. P02.06 introduced no new persistence, so G4 is N/A for new migration while retained P02.05 migration evidence passed. Immutable evidence is `docs/roadmap/evidence/P02.06_COMPLETION_2026-08-24.md`.

Diagnostic run `32664671013 / 97256120056` remains FAIL for a corrected one-space `gofmt` alignment issue in `contextual_errors.go`. It is not completion evidence and no behavior, acceptance criterion or gate was changed.

## Active P02.07 scope

P02.07 owner is `kernel.identity`. Authorized scope is limited to MFA factor enrollment/verification/removal lifecycle; passkey/WebAuthn-ready credential/challenge contracts using approved platform primitives; approved additional factor semantics where implemented; one-way/secure recovery-code lifecycle; strong-auth/step-up hooks for privileged operations that never replace authorization; replay/expiry/principal-session binding checks; synthetic fixtures; no-secret telemetry/audit behavior; and applicable owner-bounded persistence/migration evidence.

Factor secrets, recovery codes and authentication-equivalent material are `RESTRICTED`. No factor material may appear in logs, traces, errors or audit payloads. Expired, replayed or wrong-principal/session challenges fail closed. Factor removal/security-policy change may invalidate sessions according to policy. Strong authentication never replaces or bypasses authorization.

P02.08 service accounts/API credentials, P02.09 tenant settings, P02.10 phase-exit behavior, P24 enterprise SSO/SAML/SCIM, business portal UI/features, custom cryptographic algorithms/private-key management outside approved platform primitives, deployment authority and AI/model/agent runtime remain unauthorized.

## Current implementation lock

- `kernel_code_authorized=true` only for P02.07.
- `business_feature_code_authorized=false`.
- P02.08-P02.10 remain planned.
- P03+ remains planned.

## Public visibility / Issue #4

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains an external distribution/public-launch licensing, IP and trademark decision gate.

## Roadmap

`docs/roadmap/MASTER_PLAN.md` governs P00-P27. Current checkpoint: **P00 done; P01 done 12 / 12; P02 active 6 / 10; P02.01-P02.06 done; P02.07 sole active package; business features locked.**

## Product principle

**Universal Kernel + Extreme Modularity + Universal Workflow + Unified Business Graph + Governed AI Execution.**
