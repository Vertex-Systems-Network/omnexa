# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## Current canonical state

`docs/roadmap/STATE.json` is the machine-readable execution source of truth.

```text
Foundation Architecture v1: FROZEN
P00: DONE — 10 / 10
Repository visibility: PUBLIC
EG-02 / Issue #3: SATISFIED / CLOSED
EG-03 / Issue #14: SATISFIED
Canonical CI: GITHUB-HOSTED ONLY / ubuntu-24.04
Local/self-hosted governance runners: PROHIBITED
P01: ACTIVE
P01.01-P01.11: DONE
P01.12: ACTIVE — Developer CLI baseline / P01 exit proof
P01 progress: 11 / 12 done
kernel_code_authorized: true
business_feature_code_authorized: false
```

Kernel authorization is bounded to the sole active package. It is not permission to implement P02+, module runtime, business features, deployment administration or AI/model/agent runtime.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current `docs/ai/handoffs/<WORK_PACKAGE>.md` before material work. These files are subordinate snapshots/indexes and never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or other canonical governance. If continuity disagrees with authoritative repository/GitHub evidence, mark it **STALE** and follow authoritative sources.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/FOUNDATION_FREEZE.json` and `docs/governance/P01_ENTRY_GATE.md`;
4. `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`;
5. the active package specification (`P01.12.md` currently);
6. Product Constitution, system/module architecture, glossary, naming, ownership and dependency matrix;
7. identifier/money/time/locale/error/API/event standards;
8. security/data-classification/threat model;
9. testing/CI/release/quality standards including `docs/quality/GO_CODE_QUALITY.md`;
10. repository/local-development/toolchain/configuration/developer-command standards;
11. SLO/incident/reliability standards;
12. AI Execution Policy, Change Control and Definition of Done;
13. relevant accepted ADRs, especially ADR-0010.

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

Issue #3 / EG-02 is satisfied. `main` is protected with PR-only integration, strict required `governance`, blocked direct/force updates, failed-check merge rejection, required conversation resolution and strict up-to-date enforcement.

Canonical governance CI is **GitHub-hosted only** on `ubuntu-24.04`. The job must fail closed unless `RUNNER_ENVIRONMENT=github-hosted`, `RUNNER_OS=Linux` and `RUNNER_ARCH=X64`. Do not reintroduce `self-hosted`, local evidence fanout or local-runner fallback.

The permanent repository Go quality gate runs through `bash scripts/verify_go_quality.sh` with pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`. Do not remove it, weaken checks merely to obtain green CI, use `@latest`, or silently auto-fix source in CI.

Strict protection requires implementation/closure PRs to be current with protected `main` before merge. Stale green runs are not merge permission.

## P01 execution rule

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` enforces a completed prefix and exactly one active package.

P01.01 through P01.11 are complete with canonical evidence under `docs/roadmap/evidence/`. Latest completed evidence is `docs/roadmap/evidence/P01.11_COMPLETION_2026-08-23.md`, based on implementation PR #63, final exact head `1c1ab1f8d5120fb6b1e5908fdb93cffef9275940`, final canonical run/job `32610902537 / 97123708250`, and implementation merge `10c94a638b89d47da05f5481fb2db298a2da6942`.

Current active package: **P01.12 — Developer CLI Baseline**. It is the final P01 package and owns the P01 fresh-install exit proof.

P02 may not activate until P01.12 reaches `done` with required evidence and a separate governed P01-exit reconciliation. More than one active P01 package is forbidden.

## P01.12 developer CLI / P01 exit rules

Authorized P01.12 scope is limited to:

- stable repository-owned `omnexa`/developer command surface for version, verify, build/test helpers and approved local diagnostics;
- deterministic canonical `verify` orchestration mapped to applicable quality gates;
- explicit fail-closed exit codes and structured-safe output;
- command help/version metadata;
- safe composition of P01.02 configuration, P01.04 migration and P01.08 diagnostics where explicitly authorized;
- no-secret / no-RESTRICTED-output behavior;
- deterministic invocation from a clean checkout and canonical CI;
- P01.01-P01.11 regression preservation;
- full P01 fresh-install exit proof: configuration, build/start, fresh migration, cache/storage contracts, safe telemetry, readiness/diagnostics, jobs/configuration/audit primitives, canonical developer verification and required security/supply-chain/build gates.

P01.12 must not implement production super-admin authority, P02 tenant/user/role administration, P03 module install/runtime administration, P04+ domain/event/workflow commands, deployment/Kubernetes orchestration, hidden SQL/file mutation, business modules or AI/model/agent behavior.

CLI convenience never creates authority. Destructive operations require explicit environment/resource semantics. Privileged future operations must authenticate, authorize and audit through governed capabilities rather than hidden CLI bypasses.

## Completed P01.11 audit rules retained

Protected audit remains separate from ordinary logs. P01.11 records are immutable/tamper-evident, classification-aware and append-oriented. Required-audit failure cannot silently claim success. Generic audit records reject secret/auth/key/payment-sensitive fields. Audit write capability does not imply audit read/export authority, and actor/scope metadata does not grant authentication, authorization, tenancy or identity authority. P01.11 regression verification remains mandatory.

## Completed lower-package invariants retained

P01.10 runtime configuration remains distinct from P01.02 static process configuration; kill switches cannot grant authority. P01.09 jobs remain deterministic, bounded and non-authoritative. P01.08 health/readiness remains safe operational evidence rather than authority. P01.07 observability remains diagnostic infrastructure rather than audit/correctness authority. P01.06 storage and P01.05 cache remain infrastructure primitives without tenancy/authorization authority. Their regression verifiers remain mandatory.

## Business-feature lock

`business_feature_code_authorized=false` remains mandatory for all P01 work. Do not implement CRM, ERP, commerce, payment, POS, CMS, portal, HR/projects, supply chain, integrations, builders, BI, AI-agent business behavior or any other business domain.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

For P01.12 and P01 exit, record applicable G0-G8 evidence accurately, preserve repository Go quality plus P01.01-P01.11 regressions, and prove the fresh-install exit path reproducibly on canonical GitHub-hosted `ubuntu-24.04` without hidden manual steps.

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

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is a mandatory planning input whenever a future package authorizes browser UI. It does not authorize UI work during P01.

## Required work protocol

For every material change:

1. verify active phase/package and locks in `STATE.json`;
2. inspect active package spec and frozen standards;
3. preserve ownership/dependency boundaries;
4. implement only authorized scope;
5. add positive/negative evidence appropriate to risk;
6. run canonical GitHub-hosted `governance`;
7. inspect diff/status before merge;
8. merge only when required checks are green and the branch is current with protected `main`;
9. reconcile state/status/continuity only after completion evidence exists;
10. use ADR/change control before changing frozen architecture.

The AI must not automatically advance to the next work package or phase. After a package closure activates the next package, identify the next authorized action and **STOP**; do not implement that newly activated package in the same execution session.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement P02/P03/business scope early; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3. It does not block bounded P01 kernel engineering.

## Exact next transition

This closure transition records P01.11 `done` and activates only P01.12. After this closure merges, **STOP**. A new governed execution session may then implement only the bounded P01.12 developer CLI baseline and P01 exit proof. P02+, business features and AI/model/agent runtime remain locked until their own governed transitions.
