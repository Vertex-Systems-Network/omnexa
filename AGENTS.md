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
P02: ACTIVE — 0 / 10 done
Current work package: P02.01 — Principal & user identity foundation
P02.02-P02.10: PLANNED
kernel_code_authorized: true — P02.01 only
business_feature_code_authorized: false
```

Implementation authority is bounded to the single active package P02.01. No later P02 package, P03+, business feature, deployment administration or AI/model/agent runtime is authorized.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.01.md` before material work. These files are subordinate snapshots/indexes and never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or canonical GitHub evidence.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/P01_EXIT_GATE.md`, `docs/governance/P02_ENTRY_GATE.md`, `docs/governance/P02_EXIT_GATE.md` and `docs/governance/P01_P02_TRANSITION_CHECKLIST.md`;
4. `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json` and the active `docs/roadmap/work-packages/P02.01.md`;
5. Product Constitution, system/module architecture, glossary, naming, ownership and dependency matrix;
6. identifier/money/time/locale/error/API/event standards;
7. security/data-classification/threat model;
8. testing/CI/release/quality standards including `docs/quality/GO_CODE_QUALITY.md`;
9. repository/local-development/toolchain/configuration/developer-command standards;
10. SLO/incident/reliability standards;
11. AI Execution Policy, Change Control and Definition of Done;
12. relevant accepted ADRs, especially ADR-0010.

If canonical documents conflict, resolve through change control before implementation.

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

Canonical governance CI is **GitHub-hosted only** on `ubuntu-24.04`. The job must fail closed unless `RUNNER_ENVIRONMENT=github-hosted`, `RUNNER_OS=Linux` and `RUNNER_ARCH=X64`. Do not reintroduce `self-hosted`, local evidence fanout or local-runner fallback.

The permanent repository Go quality gate runs through `bash scripts/verify_go_quality.sh` with pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`. Do not remove it, weaken checks merely to obtain green CI, use `@latest`, or silently auto-fix source in CI.

Strict protection requires implementation/closure/activation PRs to be current with protected `main` before merge. Stale green runs are not merge permission.

## Completed P01 prerequisite retained

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` records P01.01-P01.12 as `done` with zero active packages and `implementation_authorized=false`. `docs/governance/P01_EXIT_GATE.md` is **SATISFIED**. All P01 regression verifiers remain mandatory during P02.

Final P01 evidence is retained under `docs/roadmap/evidence/P01.12_COMPLETION_2026-08-23.md`, based on implementation PR #65, final exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, and merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

## P02 activation and sequencing

P02 readiness preparation merged through PR #67 as `c6301ca4a5eec5dd62bcb75481d900e40ad968bd` after final exact-head canonical run `32632920772 / 97178312240` PASS.

`docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution. At this checkpoint:

- P02.01 is `active`;
- P02.02-P02.10 are `planned`;
- P02 progress is `0 / 10 done`;
- `kernel_code_authorized=true` only for P02.01;
- `business_feature_code_authorized=false`.

The AI must not automatically advance to another package. Implementation and closure/state transition remain separate governed PRs.

## Active P02.01 boundary

Owner: `kernel.identity`.

P02.01 may implement only the canonical human principal/User identity foundation described by `docs/roadmap/work-packages/P02.01.md`, including stable UUIDv7 identity/principal IDs, deterministic User lifecycle semantics, classification-safe identity attributes, owner-bounded persistence if required, safe errors, focused positive/negative tests and applicable migration evidence.

P02.01 invariants:

- User is authentication identity, not business Person;
- non-human principals are not fake human users;
- identity attributes do not imply authorization;
- authentication secrets/session credentials are not introduced here;
- tenant claims are not authority;
- only `kernel.identity` owns this write boundary;
- no direct writes into later domains.

Explicitly forbidden in P02.01:

- P02.02 tenant lifecycle/membership authority;
- P02.03 organization hierarchy;
- P02.04 authentication/session lifecycle;
- P02.05-P02.06 RBAC or relationship/context policy;
- P02.07 MFA/passkeys;
- P02.08 service-account/API credential lifecycle;
- P02.09 tenant-scoped settings;
- P02.10 identity/permission audit-trail product behavior beyond reuse of the existing P01 audit transport where required;
- P03 module runtime;
- business modules/features;
- deployment/Kubernetes authority;
- AI/model/agent runtime.

## Completed P01 capability rules retained

Protected audit remains separate from ordinary logs. P01.11 audit is immutable/tamper-evident, classification-aware and append-oriented; required-audit failure cannot silently claim success and audit write does not imply read/export authority.

P01.10 runtime configuration remains distinct from P01.02 static process configuration; kill switches cannot grant authority. P01.09 jobs remain deterministic, bounded and non-authoritative. P01.08 health/readiness remains safe operational evidence rather than authority. P01.07 observability remains diagnostic infrastructure rather than audit/correctness authority. P01.06 storage and P01.05 cache remain infrastructure primitives without tenancy/authorization authority.

The completed `kernel.developer` CLI remains convenience tooling only; it does not grant production super-admin, identity/tenant administration, module runtime, business or deployment authority.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

P02.01 implementation must retain repository Go quality and P01.01-P01.12 regressions, add risk-appropriate positive/negative tests, and add fresh/supported-upgrade migration evidence if persistence changes. Canonical completion evidence must come from GitHub-hosted `ubuntu-24.04`.

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

## Future browser UI quality rule

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is a mandatory planning input whenever a future package authorizes browser UI. It does not authorize UI work by itself.

## Required work protocol

For every material change:

1. verify phase/package and locks in `STATE.json`;
2. inspect the authorized specification/governance scope and frozen standards;
3. preserve ownership/dependency boundaries;
4. implement only explicitly authorized scope;
5. add positive/negative evidence appropriate to risk;
6. run canonical GitHub-hosted `governance`;
7. inspect diff/status before merge;
8. merge only when required checks are green and the branch is current with protected `main`;
9. reconcile state/status/continuity only after completion evidence exists;
10. use ADR/change control before changing frozen architecture.

After a closure activates a new package or phase, identify the next authorized action and **STOP**; do not implement that newly activated scope in the same execution session.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement unactivated future-phase scope; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3. It does not block internal P02 engineering.

## Exact next transition

This activation transition makes P02.01 the sole active implementation scope. After the activation PR merges, verify protected `main` and canonical `STATE.json`, identify the next authorized action as P02.01 implementation, then **STOP**. P02.01 implementation starts only in a later governed execution session from the then-current protected `main`.
