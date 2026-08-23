# Omnexa Program Status

Last reconciled: **2026-08-23**

## Current position

- Program: **Kernel Program**
- Phase: **P02 — Identity, Tenancy & Organization**
- Phase state: **active**
- Current work package: **P02.05 — RBAC foundation**
- P02 progress: **4 / 10 done**
- P02.01-P02.04: **DONE**
- P02.05: **ACTIVE**
- P02.06-P02.10: **PLANNED**
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
- Kernel implementation: **AUTHORIZED ONLY FOR P02.05**
- Business-feature implementation: **NOT AUTHORIZED**

## P02.04 completion

P02.04 implementation completed through PR #75. Final exact head `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f` passed canonical GitHub-hosted run/job `32653747461 / 97229198036` and squash-merged as `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`.

The accepted lane ran on runner `GitHub Actions 1000014748`, GitHub-hosted Ubuntu 24.04.4 LTS / X64, image `ubuntu-24.04 / 20260816.277.1`, Go 1.26.7 and PostgreSQL 18.6. Repository Go quality passed with gofmt over 93 Go files, golangci-lint v2.12.2 reporting 0 issues and govulncheck v1.7.0 reporting no vulnerabilities. P01.01-P01.12 regressions, `omnexa verify all`, P02.01-P02.03 regressions and P02.04 G0-G8 all passed.

P02.04 evidence includes adaptive one-way password hashing, disclosure-safe invalid-credential behavior, opaque access/refresh credentials with digest-only persistence, explicit expiry, deterministic rotation/replay rejection/revocation, password/account-lifecycle invalidation, safe session inventory, current tenant/organization context reauthorization, classification-safe lifecycle audit hooks and real PostgreSQL fresh/idempotent/P02.01-upgrade migration evidence. Authentication remains identity proof only and does not grant authorization.

Diagnostic run `32653382151 / 97228307755` remains FAIL for four corrected `gosec` G101 heuristic findings and one corrected staticcheck simplification. It is not completion evidence. Immutable completion evidence: `docs/roadmap/evidence/P02.04_COMPLETION_2026-08-23.md`.

## Retained P02.01-P02.03 completion

P02.01 remains complete through PR #69, exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical run/job `32635243643 / 97183883007`, and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 remains complete through PR #71, exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical run/job `32637760875 / 97189971101`, and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`. Its trusted tenant-boundary invariants and regression verifier remain mandatory.

P02.03 remains complete through PR #73, exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`, canonical run/job `32640790333 / 97197453122`, and merge `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`. Its organization hierarchy/scope invariants and regression verifier remain mandatory.

## Active P02.05 boundary

P02.05 is owned by `kernel.authorization`. It may implement only the RBAC foundation defined by `docs/roadmap/work-packages/P02.05.md`: stable permission identifiers and capability-oriented permission semantics; Role as a permission composition rather than a bypass identity; tenant/organization-scoped role definitions and assignments; assignment/revocation lifecycle and deterministic direct-RBAC evaluation; privileged mutation checks and classification-safe audit hooks; and allow/deny/privilege-escalation negative tests.

P02.05 remains deny-by-default. `admin`, `owner`, `superuser` or any role name never means platform bypass. Assignment itself is privileged/auditable, and tenant/org identifier substitution cannot widen authority.

P02.05 does **not** authorize P02.06 relationship/context policy, P02.07 MFA/passkeys, P02.08 service-account/API credential scope, P03 module permission registration, business-module permissions, hidden super-admin behavior, business features, deployment/Kubernetes authority or AI/model/agent runtime.

## P01 completion retained

P01.01-P01.12 remain complete with final P01.12 implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, and merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. All P01 executable regressions remain mandatory during P02.

## Protected integration / CI

`main` remains protected with PR-only integration, strict required `governance`, blocked direct/force updates, required conversation resolution and strict up-to-date enforcement. Canonical required CI remains GitHub-hosted `ubuntu-24.04` only.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block internal P02 engineering.

## Exact next work

After this P02.04 closure/P02.05 activation transition merges, verify protected `main` and canonical `STATE.json`, identify P02.05 as the sole authorized implementation scope, then **STOP this execution session**. P02.05 implementation starts only in a later governed execution session from the then-current protected `main`. Do not auto-advance to P02.06.
