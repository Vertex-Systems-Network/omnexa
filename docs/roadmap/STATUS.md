# Omnexa Program Status

Last reconciled: **2026-08-23**

## Current position

- Program: **Kernel Program**
- Phase: **P02 — Identity, Tenancy & Organization**
- Phase state: **active**
- Current work package: **P02.04 — Authentication & session lifecycle**
- P02 progress: **3 / 10 done**
- P02.01-P02.03: **DONE**
- P02.04: **ACTIVE**
- P02.05-P02.10: **PLANNED**
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
- Kernel implementation: **AUTHORIZED ONLY FOR P02.04**
- Business-feature implementation: **NOT AUTHORIZED**

## P02.03 completion

P02.03 implementation completed through PR #73. Final exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c` passed canonical GitHub-hosted run/job `32640790333 / 97197453122` and squash-merged as `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`.

The final lane ran on GitHub-hosted Ubuntu 24.04.4 LTS / X64, image `ubuntu-24.04 / 20260816.277.1`, runner `GitHub Actions 1000014421`, Go 1.26.7. Repository Go quality passed with gofmt over 86 files, golangci-lint v2.12.2 reporting 0 issues and govulncheck v1.7.0 reporting no vulnerabilities. P01.01-P01.12 regressions, `omnexa verify all`, P02.01-P02.02 regressions and P02.03 G0-G8 all passed, including real PostgreSQL fresh/idempotent/P02.02-upgrade migration evidence, same-tenant hierarchy/membership positives, cross-tenant parent/membership negatives, cycle rejection and non-authorizing scope context.

Diagnostic run `32640199607 / 97196005995` remains FAIL for corrected `govet` variable-shadow findings. Diagnostic run `32640419476 / 97196545810` remains FAIL because the P02.03 verifier initially omitted already-governed transitive `kernel.database`/`kernel.config` prerequisites from its narrow dependency allowlist. Neither diagnostic failure is completion evidence. Immutable completion evidence: `docs/roadmap/evidence/P02.03_COMPLETION_2026-08-23.md`.

## Retained P02.01-P02.02 completion

P02.01 remains complete through PR #69, exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical run/job `32635243643 / 97183883007`, and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 remains complete through PR #71, exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical run/job `32637760875 / 97189971101`, and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`. Its trusted tenant-boundary invariants and regression verifier remain mandatory.

## Active P02.04 boundary

P02.04 is owned by `kernel.identity`. It may implement only the authentication and interactive-session lifecycle defined by `docs/roadmap/work-packages/P02.04.md`: authentication mechanism boundaries; approved adaptive password hashing where passwords are supported; session/access/refresh credential expiry, rotation and revocation; short-lived access credentials; device/session inventory semantics where supported; required invalidation after material account/security changes; tenant/organization context re-authorization; disclosure-safe authentication failures; and classification-safe audit hooks without secret payloads.

Authentication proves identity; it does not grant business authority. Passwords may never be plaintext or reversibly stored. Refresh/session secrets remain `RESTRICTED` and must never be logged. Bearer possession cannot bypass current authorization/policy state.

P02.04 does **not** authorize RBAC/policy decisions, MFA/passkeys, service accounts/API credentials, SAML/SCIM/enterprise SSO, business login portals/UI, business features, P03+, deployment/Kubernetes authority or AI/model/agent runtime.

## P01 completion retained

P01.01-P01.12 remain complete with final P01.12 implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, and merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. All P01 executable regressions remain mandatory during P02.

## Protected integration / CI

`main` remains protected with PR-only integration, strict required `governance`, blocked direct/force updates, required conversation resolution and strict up-to-date enforcement. Canonical required CI remains GitHub-hosted `ubuntu-24.04` only.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block internal P02 engineering.

## Exact next work

After this P02.03 closure/P02.04 activation transition merges, **STOP this execution session**. In the next governed execution session, create a fresh implementation branch from the exact protected `main` SHA and implement only P02.04 against its existing acceptance criteria. Do not auto-advance to P02.05.
