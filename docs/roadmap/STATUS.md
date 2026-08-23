# Omnexa Program Status

Last reconciled: **2026-08-24**

## Current position

- Program: **Kernel Program**
- Phase: **P02 — Identity, Tenancy & Organization**
- Phase state: **active**
- Current work package: **P02.06 — Relationship/context-aware authorization**
- P02 progress: **5 / 10 done**
- P02.01-P02.05: **DONE**
- P02.06: **ACTIVE**
- P02.07-P02.10: **PLANNED**
- P01: **DONE — 12 / 12**
- P01 exit gate: **SATISFIED**
- P02 entry gate: **SATISFIED**
- P02 exit gate: **NOT SATISFIED — PHASE ACTIVE**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR P02.06**
- Business-feature implementation: **NOT AUTHORIZED**

## P02.05 completion

P02.05 implementation completed through PR #77. Final exact head `2df8d2a8bef0cea60256a832986d6f8495c80378` passed canonical GitHub-hosted run/job `32660848145 / 97246683239` and squash-merged as `7b6a59e83c9bd696e6e008385b4413d529254171`.

The accepted lane ran on runner `GitHub Actions 1000015357`, GitHub-hosted Ubuntu 24.04.4 LTS / X64, image `ubuntu-24.04 / 20260816.277.1`, Go 1.26.7 and PostgreSQL 18.6. Repository Go quality passed with gofmt over 100 Go files, golangci-lint v2.12.2 reporting 0 issues and govulncheck v1.7.0 reporting no vulnerabilities. P01.01-P01.12 regressions, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.04 regressions and P02.05 G0-G8 all passed.

P02.05 evidence proves stable capability-oriented permission identifiers, deterministic Role permission composition, tenant/organization scoped direct assignments, deny-by-default evaluation, privileged mutation checks, anti-escalation behavior, assignment revocation, role-name non-bypass, classification-safe required audit records and owner-bounded PostgreSQL persistence/migration evidence. It deliberately does not include P02.06 relationship/object/context policy.

Diagnostic run `32656689041 / 97236397635` remains FAIL for 29 corrected `govet` shadow findings and one corrected `gosec` G304 test-helper finding. Diagnostic run `32660398632 / 97245574862` remains FAIL because the initial P02.05 verifier incorrectly judged transitive prerequisite imports against a direct-owner allowlist after all runtime/integration evidence had passed. Both failures were corrected without suppressions, gate weakening or acceptance changes and are not completion evidence. Immutable completion evidence: `docs/roadmap/evidence/P02.05_COMPLETION_2026-08-24.md`.

## Retained P02.01-P02.04 completion

P02.01 remains complete through PR #69, exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical run/job `32635243643 / 97183883007`, and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 remains complete through PR #71, exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical run/job `32637760875 / 97189971101`, and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`. Its trusted tenant-boundary invariants and regression verifier remain mandatory.

P02.03 remains complete through PR #73, exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`, canonical run/job `32640790333 / 97197453122`, and merge `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`. Its organization hierarchy/scope invariants and regression verifier remain mandatory.

P02.04 remains complete through PR #75, exact head `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f`, canonical run/job `32653747461 / 97229198036`, and merge `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`. Its authentication/session lifecycle and current-context reauthorization invariants remain mandatory.

## Active P02.06 boundary

P02.06 is owned by `kernel.authorization`. It may implement only the relationship/context-aware authorization layer defined by `docs/roadmap/work-packages/P02.06.md`: relationship/context policy evaluation layered on the accepted P02.05 RBAC foundation; tenant, organization and object-scope relationships; capability-bound deny-by-default decisions; contextual conditions that cannot grant beyond trusted principal/scope relationships; field/export distinction hooks; disclosure-safe deny behavior and material authorization audit hooks; and same-scope allow plus wrong-tenant/wrong-org/wrong-object negative tests.

P02.06 retains all P02.05 RBAC invariants. Client IDs, tenant IDs and object IDs remain references rather than authority. Tenant membership alone is insufficient for every child scope. Role names and internal/background call origin never bypass policy. Contextual rules cannot widen authority beyond current trusted principal, tenant and organization relationships.

P02.06 does **not** authorize P02.07 MFA/passkeys, P02.08 service-account/API credential scope, P02.09 tenant settings, P02.10 phase-exit product behavior, P03 module permission registration/capability registry, business-domain object policies not owned by an active domain, support impersonation product surface, business features/UI, deployment/Kubernetes authority or AI/model/agent runtime.

## P01 completion retained

P01.01-P01.12 remain complete with final P01.12 implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, and merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. All P01 executable regressions remain mandatory during P02.

## Protected integration / CI

`main` remains protected with PR-only integration, strict required `governance`, blocked direct/force updates, required conversation resolution and strict up-to-date enforcement. Canonical required CI remains GitHub-hosted `ubuntu-24.04` only.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block internal P02 engineering.

## Exact next work

After this P02.05 closure/P02.06 activation transition passes, merges and protected `main` plus canonical `STATE.json` confirm P02.06 as the sole active package, **STOP this execution session**. P02.06 implementation starts only in a later governed execution session from the then-current protected `main`. Do not auto-advance to P02.07.
