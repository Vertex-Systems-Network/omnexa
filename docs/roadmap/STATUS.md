# Omnexa Program Status

Last reconciled: **2026-08-24**

## Current position

- Program: **Kernel Program**
- Phase: **P02 — Identity, Tenancy & Organization**
- Phase state: **active**
- Current work package: **P02.08 — Service accounts & API credentials**
- P02 progress: **7 / 10 done**
- P02.01-P02.07: **DONE**
- P02.08: **ACTIVE**
- P02.09-P02.10: **PLANNED**
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
- Kernel implementation: **AUTHORIZED ONLY FOR P02.08**
- Business-feature implementation: **NOT AUTHORIZED**

## P02.07 completion

P02.07 implementation completed through PR #81. Final exact head `51ccaa12c3534f74fba6eab9d4698ee483ef4ffd` passed canonical GitHub-hosted run/job `32669167972 / 97267175953` and squash-merged as `5642f5da1eb24e70b67e5ec757d9f4584c4e3f5c`.

The accepted lane ran on runner `GitHub Actions 1000016379`, GitHub-hosted Ubuntu 24.04.4 LTS / X64, image `ubuntu-24.04 / 20260816.277.1`, Go 1.26.7 and PostgreSQL 18.6. Repository Go quality passed over 111 Go files, golangci-lint v2.12.2 reported 0 issues and govulncheck v1.7.0 reported no vulnerabilities. P01.01-P01.12 regressions, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.06 regressions and applicable P02.07 G0-G8 all passed.

P02.07 evidence proves deterministic passkey factor enrollment/verification/removal; approved verification through an injected `PasskeyVerifier` instead of custom WebAuthn/private-key cryptography; exact User/session-bound challenge expiry/replay enforcement; recovery-code digest-only persistence and replay denial; session-bound non-authorizing step-up proof; explicit factor-removal session invalidation policy; restricted-material redaction from ordinary telemetry/audit; and fresh/idempotent/P02.04-upgrade/immutable-ledger migration evidence.

Diagnostic run `32668735841 / 97266149110` on head `0ab9873367586ab0191a91f658da275f449de796` remains FAIL because a redundant verifier regex matched the explanatory word `authorization.` in a comment that stated strong authentication cannot replace authorization. The final fix removed only that comment-sensitive token; independent authorization-import and executable authority-symbol guards remained. No runtime, schema, test, acceptance, security or gate behavior changed. Immutable completion evidence: `docs/roadmap/evidence/P02.07_COMPLETION_2026-08-24.md`.

## Retained P02.01-P02.06 completion

P02.01 remains complete through PR #69, exact head `76919a9588f70aeea7e00f5214b82dcbf34cbee7`, canonical run/job `32635243643 / 97183883007`, and merge `44882e91e49d0364d841b511edbfd0619d05de1f`.

P02.02 remains complete through PR #71, exact head `a63bd45523ed35c4b11d11c8abc0cb42ce9e11d7`, canonical run/job `32637760875 / 97189971101`, and merge `2ed0d9a5855f84ac8b7265c23ff6b8b7799b779d`. Its trusted tenant-boundary invariants and regression verifier remain mandatory.

P02.03 remains complete through PR #73, exact head `20bcafb9d2ccb5829e44f5b69130a4cd5b9e816c`, canonical run/job `32640790333 / 97197453122`, and merge `03b3d42a67d98638129b7f9d2b2f49467ae1fcec`. Its organization hierarchy/scope invariants and regression verifier remain mandatory.

P02.04 remains complete through PR #75, exact head `83a1d9e9f47e05f2e6fa7e50874dd7bfce51437f`, canonical run/job `32653747461 / 97229198036`, and merge `769423a94ec03a9f2d7b9e667b9d4527fb0660bf`. Its authentication/session lifecycle and current-context reauthorization invariants remain mandatory.

P02.05 remains complete through PR #77, exact head `2df8d2a8bef0cea60256a832986d6f8495c80378`, canonical run/job `32660848145 / 97246683239`, and merge `7b6a59e83c9bd696e6e008385b4413d529254171`. Its direct RBAC, anti-escalation, exact-scope and migration invariants remain mandatory.

P02.06 remains complete through PR #79, exact head `dbbd105fd5f2543ca7dd5df93375eaf1057928fc`, canonical run/job `32664834112 / 97256520050`, and merge `083c2866f0cd0773b85201750c2196bfd2fcc167`. Its contextual authorization, field/export separation, internal/background non-bypass and fail-closed dependency invariants remain mandatory.

## Active P02.08 boundary

P02.08 is owned by `kernel.identity`. It may implement only the Service Account and API credential lifecycle defined by `docs/roadmap/work-packages/P02.08.md`: a distinct non-human Service Account principal; tenant/organization binding and scope composition through accepted P02 authorization foundations; API credential issue/identify/verify/rotate/revoke/expire lifecycle; one-time secret presentation where applicable; non-reversible verifier storage; classification-safe inventory/audit metadata; optional classification-safe last-used/rotation metadata; and deterministic allowed-scope plus wrong-tenant/wrong-scope/revoked/expired evidence.

Non-human principals are never fake human Users. Raw API credentials are `RESTRICTED`, never logged, traced, audited or stored reversibly. Credentials remain least-privilege, rotatable, revocable and tenant/organization bound. Credential possession never bypasses current authorization state, and no generic platform superkey/master-token mechanism is authorized.

P02.08 does **not** authorize OAuth developer applications, external connector/provider integration, device/POS identities, AI agent execution identity, P02.09 tenant settings, P02.10 exit behavior, business API scopes/features/UI, deployment/Kubernetes authority or other future scope.

## P01 completion retained

P01.01-P01.12 remain complete with final P01.12 implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, and merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. All P01 executable regressions remain mandatory during P02.

## Protected integration / CI

`main` remains protected with PR-only integration, strict required `governance`, blocked direct/force updates, required conversation resolution and strict up-to-date enforcement. Canonical required CI remains GitHub-hosted `ubuntu-24.04` only.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block internal P02 engineering.

## Exact next work

After this P02.07 closure/P02.08 activation transition passes, merges and protected `main` plus canonical `STATE.json` confirm P02.08 as the sole active package, **STOP this execution session**. P02.08 implementation starts only in a later governed execution session from the then-current protected `main`. Do not auto-advance to P02.09.
