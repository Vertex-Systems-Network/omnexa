# Omnexa Program Status

Last reconciled: **2026-08-23**

## Current position

- Program: **Kernel Program**
- Phase: **P02 — Identity, Tenancy & Organization**
- Phase state: **active**
- Current work package: **P02.01 — Principal & user identity foundation**
- P02 progress: **0 / 10 done**
- P02.01: **ACTIVE**
- P02.02-P02.10: **PLANNED**
- P01: **DONE — 12 / 12**
- P01 exit gate: **SATISFIED**
- P02 entry gate: **SATISFIED**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR P02.01**
- Business-feature implementation: **NOT AUTHORIZED**

## P02 readiness and activation

P02 readiness preparation completed through PR #67. Final exact-head canonical run `32632920772 / 97178312240` passed on the GitHub-hosted `ubuntu-24.04` lane and the readiness PR merged as `c6301ca4a5eec5dd62bcb75481d900e40ad968bd`.

The prepared phase contains exactly 10 strict sequential packages. This activation transition moves only P02/P02.01 to active and leaves P02.02-P02.10 planned. No P02 runtime/schema implementation is part of the activation transition.

P02.01 is owned by `kernel.identity` and is limited to the canonical human principal/User identity foundation. User remains distinct from business Person. Tenant lifecycle, organization hierarchy, authentication/session behavior, RBAC/policy, MFA/passkeys, service-account credential lifecycle, P03 module runtime, business domains and AI/model/agent behavior remain outside P02.01.

## P01 completion retained

P01.01-P01.12 remain complete with final P01.12 implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, and merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

`docs/governance/P01_EXIT_GATE.md` remains **SATISFIED**. P01 executable regressions continue to run in canonical governance during P02.

## Repository Go quality

`docs/quality/GO_CODE_QUALITY.md` remains the permanent Go quality gate. Canonical `governance` runs on GitHub-hosted `ubuntu-24.04` using pinned Go 1.26.7, golangci-lint v2.12.2 and govulncheck v1.7.0.

## Protected integration / CI

`main` remains protected with PR-only integration, strict required `governance`, blocked direct/force updates, required conversation resolution and strict up-to-date enforcement. Local/self-hosted governance runners remain prohibited.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block internal P02 engineering.

## Exact next work

After this P02.01 activation transition merges, **STOP this execution session**. In a later governed session, create a fresh implementation branch from protected `main` and implement only P02.01 with focused positive/negative tests, applicable migration evidence, a dedicated verifier and canonical GitHub-hosted CI. Do not auto-advance to P02.02.
