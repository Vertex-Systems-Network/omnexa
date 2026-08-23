# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## Current canonical state

`docs/roadmap/STATE.json` is the machine-readable execution source of truth.

```text
Foundation Architecture v1: FROZEN
P00: DONE — 10 / 10
Repository visibility: PUBLIC
Issue #3: SATISFIED / CLOSED
Issue #14: SATISFIED
Canonical CI: GITHUB-HOSTED ONLY / ubuntu-24.04
Local/self-hosted governance runners: PROHIBITED
P01: DONE — 12 / 12
P01 exit gate: SATISFIED
P02: ACTIVE — 1 / 10 done
P02.01: DONE — Principal & user identity foundation
Current work package: P02.02 — Tenant lifecycle & trusted tenant context
P02.03-P02.10: PLANNED
kernel_code_authorized: true — P02.02 only
business_feature_code_authorized: false
```

Implementation authority is bounded to the single active package P02.02. No later P02 package, P03+, business feature, deployment administration or AI/model/agent runtime is authorized.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.02.md` before material work. These files are subordinate snapshots/indexes and never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or canonical GitHub evidence.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. P01/P02 entry/exit gates and transition checklist;
4. `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json` and the active `docs/roadmap/work-packages/P02.02.md`;
5. Product Constitution, architecture, glossary, naming, ownership and dependency matrix;
6. identifier/money/time/locale/error/API/event standards;
7. security/data-classification/threat model;
8. testing/CI/release/quality standards including `docs/quality/GO_CODE_QUALITY.md`;
9. repository/local-development/toolchain/configuration/developer-command standards;
10. SLO/incident/reliability standards;
11. AI Execution Policy, Change Control and Definition of Done;
12. relevant accepted ADRs, especially ADR-0010.

If canonical documents conflict, reconcile through change control before implementation.

## Frozen architecture laws

1. Kernel before business modules.
2. One authoritative owner per write model/capability.
3. Cross-module direct DB writes and private implementation imports are forbidden.
4. Cross-domain integration uses governed APIs/capabilities, events, workflows or approved read projections.
5. Tenant/org boundaries, authorization, audit, observability and versioned contracts are mandatory.
6. Optional-module failure/removal cannot corrupt unrelated domains.
7. Retriable work is idempotent where required.
8. AI acts only through governed capabilities; no unrestricted DB authority.
9. Strict modular monolith first; service extraction requires evidence plus ADR.
10. Infrastructure complexity must be earned.

Frozen primitives include UUIDv7 IDs, exact-decimal money with explicit currency, UTC/timestamptz instants with IANA civil-time semantics, BCP 47 locale/RTL support, stable safe structured errors, versioned HTTP/OpenAPI contracts, CloudEvents-compatible event envelopes, at-least-once/idempotent event handling, four data classes (`PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`) and deny-by-default authorization/tenant isolation.

## Protected integration and CI

Issue #3 is satisfied. `main` is protected with PR-only integration, strict required `governance`, blocked direct/force updates, failed-check merge rejection, required conversation resolution and strict up-to-date enforcement.

Canonical governance CI is **GitHub-hosted only** on `ubuntu-24.04`, Linux/X64. Do not reintroduce `self-hosted`, local evidence fanout or local-runner fallback.

The permanent repository Go quality gate runs through `bash scripts/verify_go_quality.sh` with pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`. Do not weaken checks merely to obtain green CI, use `@latest`, or silently auto-fix source in CI.

Strict protection requires implementation/closure/activation PRs to be current with protected `main` before merge. Stale green runs are not merge permission.

## Completed P01 prerequisite retained

P01.01-P01.12 remain `done`, P01 exit remains **SATISFIED**, and all P01 regression verifiers remain mandatory during P02. Final P01 evidence remains implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

## Completed P02.01 evidence retained

P02.01 is complete under `kernel.identity`.

- implementation PR: #69
- final exact head: `76919a9588f70aeea7e00f5214b82dcbf34cbee7`
- canonical run/job: `32635243643 / 97183883007` — PASS
- implementation merge: `44882e91e49d0364d841b511edbfd0619d05de1f`
- GitHub-hosted Ubuntu 24.04.4 LTS / X64, image `ubuntu-24.04 / 20260816.277.1`, Go 1.26.7
- repository Go quality, P01.01-P01.12 regressions and P02.01 G0-G8: PASS
- completion evidence: `docs/roadmap/evidence/P02.01_COMPLETION_2026-08-23.md`

Initial run `32635051321 / 97183427697` remains diagnostic FAIL for corrected `govet` shadow findings and is not completion evidence.

P02.01 invariants remain binding: User is not business Person; non-human principals are not fake users; identity attributes do not grant authority; authentication secrets/sessions were not introduced; owner-bounded `kernel.identity` persistence remains authoritative.

## Active P02.02 boundary

Owner: `kernel.tenancy`.

P02.02 may implement only the Tenant lifecycle and trusted tenant-context foundation described by `docs/roadmap/work-packages/P02.02.md`, including:

- authoritative Tenant identity/lifecycle/state;
- explicit tenant-scoped persistence/query semantics where applicable;
- trusted tenant context derived from governed identity/execution relationships rather than client payload claims;
- only the minimum relationship/membership primitive needed for later scoped authorization;
- fail-closed same-tenant/cross-tenant boundaries;
- applicable fresh/upgrade migration evidence;
- focused positive/negative tests and a dedicated verifier.

P02.02 invariants:

- a client-provided `tenant_id` is never authorization authority;
- there is no global-tenant fallback;
- tenant isolation applies independently of data classification;
- no hidden support-superuser or cross-tenant bypass;
- `kernel.tenancy` owns this write boundary;
- P02.01 User identity remains separate from tenant membership/authority.

Explicitly forbidden in P02.02:

- P02.03 organization hierarchy;
- P02.04 authentication/session implementation;
- P02.05-P02.06 RBAC or relationship/context policy beyond the minimum trusted tenant relationship primitive authorized by P02.02;
- P02.07 MFA/passkeys;
- P02.08 service-account/API credential lifecycle;
- P02.09 tenant settings;
- P02.10 phase-exit audit product behavior;
- P03 module runtime;
- business modules/features;
- deployment/Kubernetes authority;
- AI/model/agent runtime.

## Completed kernel capability rules retained

Protected audit remains separate from ordinary logs. P01.11 audit is immutable/tamper-evident, classification-aware and append-oriented; required-audit failure cannot silently claim success and audit write does not imply read/export authority. P01.10 configuration flags cannot grant authority. P01.09 jobs remain non-authoritative. P01.08 diagnostics remain operational evidence rather than authority. Cache/storage/observability remain infrastructure primitives without tenancy/authorization authority. The developer CLI remains convenience tooling only.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

P02.02 implementation must retain repository Go quality, P01.01-P01.12 regressions and P02.01 regression evidence; add risk-appropriate same-tenant/cross-tenant tests; and add migration evidence if persistence changes. Canonical completion evidence must come from GitHub-hosted `ubuntu-24.04`.

## Repository/local-development rules

Canonical roots: `apps/`, `kernel/`, `modules/`, `platform/`, `shared/`, `infrastructure/`, `scripts/`, `docs/`, `generated/`.

- folder != microservice;
- module private code/schema/migrations stay with owner;
- generated output is derivative, not source of truth;
- repository toolchains/dependencies are pinned;
- secrets are separate from committed config;
- production sensitive data is prohibited locally by default;
- Linux is canonical backend/CI environment;
- supported workflows must not depend on hidden manual SQL/file/UI steps.

## Required work protocol

For every material change:

1. verify phase/package and locks in `STATE.json`;
2. inspect authorized specification/governance scope and frozen standards;
3. preserve ownership/dependency boundaries;
4. implement only explicitly authorized scope;
5. add positive/negative evidence appropriate to risk;
6. run canonical GitHub-hosted `governance`;
7. inspect diff/status before merge;
8. merge only when required checks are green and branch is current with protected `main`;
9. reconcile state/status/continuity only after completion evidence exists;
10. use ADR/change control before changing frozen architecture.

After a closure activates a new package or phase, identify the next authorized action and **STOP**; do not implement that newly activated scope in the same execution session.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement unactivated future-phase scope; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3. It does not block internal P02 engineering.

## Exact next transition

This closure marks P02.01 done and activates P02.02 as the sole next implementation scope. After the closure PR merges, verify protected `main` and canonical `STATE.json`, identify P02.02 implementation as the next authorized action, then **STOP**. P02.02 implementation starts only in a later governed execution session from the then-current protected `main`.
