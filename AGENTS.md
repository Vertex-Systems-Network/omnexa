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
P01: DONE
P01.01-P01.12: DONE
P01 progress: 12 / 12 done
P01 exit gate: SATISFIED
Current work package: NONE
P02: PLANNED / NOT ACTIVE
kernel_code_authorized: false
business_feature_code_authorized: false
```

No executable product-development scope is active at this checkpoint. P01 completion does not implicitly activate P02, P03+, business features, deployment administration or AI/model/agent runtime. New implementation authority requires a separate governed activation reflected in `STATE.json`.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant current handoff before material work. These files are subordinate snapshots/indexes and never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or canonical GitHub evidence. If continuity disagrees with authoritative repository/GitHub evidence, mark it **STALE** and follow authoritative sources.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/FOUNDATION_FREEZE.json`, `docs/governance/P01_ENTRY_GATE.md` and `docs/governance/P01_EXIT_GATE.md`;
4. `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` and relevant next-phase planning/specification artifacts;
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

Issue #3 / EG-02 is satisfied. `main` is protected with PR-only integration, strict required `governance`, blocked direct/force updates, failed-check merge rejection, required conversation resolution and strict up-to-date enforcement.

Canonical governance CI is **GitHub-hosted only** on `ubuntu-24.04`. The job must fail closed unless `RUNNER_ENVIRONMENT=github-hosted`, `RUNNER_OS=Linux` and `RUNNER_ARCH=X64`. Do not reintroduce `self-hosted`, local evidence fanout or local-runner fallback.

The permanent repository Go quality gate runs through `bash scripts/verify_go_quality.sh` with pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`. Do not remove it, weaken checks merely to obtain green CI, use `@latest`, or silently auto-fix source in CI.

Strict protection requires implementation/closure PRs to be current with protected `main` before merge. Stale green runs are not merge permission.

## P01 completion rule

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` records P01.01-P01.12 as `done`. While P01 was executing it enforced a completed prefix and exactly one active package; at the terminal 12/12 checkpoint it must contain zero active packages and `implementation_authorized=false`.

P01 completion evidence is retained under `docs/roadmap/evidence/`. Final package evidence is `docs/roadmap/evidence/P01.12_COMPLETION_2026-08-23.md`, based on implementation PR #65, final exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, final canonical run/job `32629072886 / 97168916985`, and implementation merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

`docs/governance/P01_EXIT_GATE.md` is **SATISFIED**. The canonical GitHub-hosted lane proved fresh configuration/startup, PostgreSQL migration, cache/storage provider contracts, safe logs/telemetry, readiness/diagnostics, jobs/configuration/audit primitives, the repository-owned developer verification path, security/static checks, module checksum verification and build/package behavior without hidden manual steps.

P02 remains `planned`. P01 completion is not permission to start P02 implementation.

## Completed P01.12 developer CLI rules retained

The completed `kernel.developer` boundary provides:

- deterministic repository-owned `omnexa` help/version behavior;
- safe structured `health` diagnostics;
- guarded non-production `db migrate` using the existing migration foundation;
- deterministic fail-closed `verify <target>` orchestration mapped to governed quality gates;
- exact executable-plus-argument allowlisting with no shell-string expansion;
- runtime `OMNEXA_*` isolation from verification subprocesses;
- canonical `verify all` composition preserving repository Go quality and P01.01-P01.11 regressions;
- explicit N/A for P01 module lifecycle rather than inventing P03 behavior;
- no-secret/no-RESTRICTED output behavior.

CLI convenience never creates authority. Production super-admin commands, P02 identity/tenant administration, P03 module runtime administration, P04+ domain/event/workflow commands, deployment/Kubernetes orchestration, hidden SQL/file mutation, business modules and AI/model/agent behavior remain outside this completed boundary.

## Completed lower-package invariants retained

Protected audit remains separate from ordinary logs. P01.11 audit is immutable/tamper-evident, classification-aware and append-oriented; required-audit failure cannot silently claim success and audit write does not imply read/export authority.

P01.10 runtime configuration remains distinct from P01.02 static process configuration; kill switches cannot grant authority. P01.09 jobs remain deterministic, bounded and non-authoritative. P01.08 health/readiness remains safe operational evidence rather than authority. P01.07 observability remains diagnostic infrastructure rather than audit/correctness authority. P01.06 storage and P01.05 cache remain infrastructure primitives without tenancy/authorization authority. All P01 regression verifiers remain mandatory.

## Implementation lock after P01

`kernel_code_authorized=false` and `business_feature_code_authorized=false` are mandatory at this checkpoint. Do not implement Identity/Tenancy, module runtime, CRM, ERP, commerce, payment, POS, CMS, portal, HR/projects, supply chain, integrations, builders, BI, AI-agent behavior or any other future-phase scope until a separate governed activation explicitly changes the lock.

Governance/specification/readiness work for the next phase may proceed only within the non-implementation allowance recorded in `STATE.json`.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

P01 completed only after applicable G0-G8 evidence, repository Go quality, P01.01-P01.12 regressions and the P01 fresh-install exit path were proven on canonical GitHub-hosted `ubuntu-24.04`.

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

The AI must not automatically advance to the next work package or phase. After a closure activates a new package or phase, identify the next authorized action and **STOP**; do not implement that newly activated scope in the same execution session.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement unactivated future-phase scope; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3. It does not invalidate completed internal kernel engineering, but it must be resolved before any distribution/public-launch action that depends on that decision.

## Exact next transition

This closure reconciles P01.12 `done` and P01 `done` at 12/12 while keeping P02 `planned` and both implementation locks false. After this closure merges, **STOP**. A new governed execution session may prepare P02 work-package specifications/readiness and an explicit P02 activation transition. No P02 implementation may begin until that transition has merged and `STATE.json` authorizes it.
