# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## Current canonical state

`docs/roadmap/STATE.json` is the machine-readable execution source of truth.

```text
Foundation Architecture v1: FROZEN
P00: DONE — 10 / 10
Repository visibility: PUBLIC
Issue #3: SATISFIED / CLOSED
Canonical CI: GITHUB-HOSTED ONLY / ubuntu-24.04
Local/self-hosted governance runners: PROHIBITED
P01: DONE — 12 / 12
P01 exit gate: SATISFIED
P02: DONE — 10 / 10
P02 exit gate: SATISFIED
P03: ACTIVE — 1 / 11 done
P03.01: DONE
Current work package after this closure transition merges: P03.02 — Registry & Deterministic Discovery
P03.03-P03.11: PLANNED
kernel_code_authorized: true — P03.02 only after closure merge
business_feature_code_authorized: false
```

Implementation authority becomes bounded to P03.02 only after this closure/state-transition PR passes canonical governance and merges to protected `main`. Until then, protected `main` remains authoritative at the prior cursor. P03.03+, P04+, business features, strategic X-program runtime, deployment administration and AI/model/agent runtime remain unauthorized.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and, after activation is verified on protected `main`, `docs/ai/handoffs/P03.02.md` before material work. These files are subordinate snapshots/indexes and never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or canonical GitHub evidence.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/P02_EXIT_GATE.md`, `docs/governance/P03_ENTRY_GATE.md`, `docs/governance/P03_EXIT_GATE.md` and `docs/governance/P02_P03_TRANSITION_CHECKLIST.md`;
4. `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, completed `docs/roadmap/work-packages/P03.01.md`, its completion evidence, and active `docs/roadmap/work-packages/P03.02.md`;
5. Product Constitution, architecture, glossary, naming, ownership and dependency matrix;
6. identifier/money/time/locale/error/API/event standards;
7. security/data-classification/threat model;
8. testing/CI/release/quality standards including `docs/quality/GO_CODE_QUALITY.md`;
9. repository/local-development/toolchain/configuration/developer-command standards;
10. SLO/incident/reliability standards;
11. AI Execution Policy, Change Control and Definition of Done;
12. relevant accepted ADRs, especially ADR-0010, plus proposed strategic overlays only as non-authorizing compatibility context.

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

## Completed P01/P02/P03.01 prerequisites retained

P01.01-P01.12 remain `done`; P01 exit remains **SATISFIED**. Final P01 evidence remains PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

P02.01-P02.10 remain `done`; P02 exit remains **SATISFIED**. Terminal P02 evidence remains:

- implementation PR: #88
- final exact head: `975e4925060a035780ca13b68c5437634ed0f4ea`
- canonical run/job: `32904678957 / 97986011269` — PASS
- implementation merge: `88799aa41da8ce8c22540146d157d488565e2ce9`
- completion evidence: `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`

P03.01 — Module Manifest Schema is complete with canonical evidence:

- implementation PR: #92
- final exact head: `87da3302605c852ae5bf43d473aaa01a9e1aaa74`
- canonical run/job: `33009396644 / 98311433013` — PASS
- implementation merge: `4229e2a28442bf475afed143bab359a770d48053`
- completion evidence: `docs/roadmap/evidence/P03.01_COMPLETION_2026-08-26.md`
- retained verifier: `scripts/verify_p03_01.sh`

Diagnostic runs remain explicit FAIL history and are never acceptance evidence. All completed P01/P02/P03.01 regressions remain mandatory during P03.02.

## P03 sequencing

`docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution. This closure candidate advances exactly one boundary:

- P03.01 is `done`;
- P03.02 becomes `active` only after the closure PR merges;
- P03.03-P03.11 remain `planned`;
- P03 progress becomes `1 / 11 done`;
- `kernel_code_authorized=true` only for P03.02 after closure merge;
- `business_feature_code_authorized=false`.

The AI must not automatically advance to another package. Implementation and closure/state transition remain separate governed PRs.

## Active P03.02 boundary after closure merge

Owner: `kernel.modules`.

P03.02 may implement only Registry & Deterministic Discovery described by `docs/roadmap/work-packages/P03.02.md`:

- registry records derived from already validated P03.01 manifest metadata;
- deterministic discovery from explicit approved repository/runtime sources;
- stable lookup/list ordering and stable source/version identity;
- duplicate/conflicting module identity/version rejection;
- structural distinction between discovered/available metadata and installed/enabled lifecycle state;
- operator-safe discovery diagnostics.

P03.02 invariants:

- discovery consumes validated metadata and cannot execute package/module code or lifecycle hooks;
- arbitrary filesystem or network scanning is forbidden;
- duplicate identity must fail closed rather than resolve through nondeterministic enumeration or last-write-wins behavior;
- registry metadata cannot grant permissions/capability authority;
- sensitive source/path detail must respect classification/redaction;
- one authoritative owner remains `kernel.modules`;
- existing P02 tenancy/authorization/audit rules and P03.01 manifest validation remain binding.

Explicitly forbidden in P03.02:

- P03.03 dependency graph resolution;
- P03.04 lifecycle runtime/persistence;
- P03.05-P03.11 settings/capability/permission/UI/migration/health/trust runtime;
- remote marketplace/catalog download, signature/trust enforcement or sandbox runtime;
- System Graph, product federation or performance-intelligence runtime;
- P04 events/workflows;
- business modules/features;
- AI/model/agent runtime.

The P03 AI-native compatibility matrix for `XQ-100`, `XSG-100`, `XTRUST-100`, `XPF-200` and `XPERF-100` remains planning-only and does not authorize strategic runtime.

## Completed kernel capability rules retained

Protected audit remains separate from ordinary logs. P01.11 audit is immutable/tamper-evident, classification-aware and append-oriented; required-audit failure cannot silently claim success and audit write does not imply read/export authority. P01.10 configuration flags cannot grant authority. P01.09 jobs remain non-authoritative. P01.08 diagnostics remain operational evidence rather than authority. Cache/storage/observability remain infrastructure primitives without tenancy/authorization authority. The developer CLI remains convenience tooling only.

P02 retained invariants remain binding: User is not business Person; trusted tenant context comes from authoritative tenant/membership state; no global tenant fallback; organization hierarchy is tenant-contained; authentication and credential possession prove identity rather than authority; sessions are revocable and current context is reauthorized; authorization is deny-by-default with exact trusted scope; contextual policy narrows only; service accounts are distinct non-human principals; settings cannot create authority; audit is classification-safe and secret-free; required-audit protected mutations cannot silently claim success when audit delivery fails.

P03.01 retained invariants remain binding: manifests are untrusted declarative metadata; parsing/validation executes no package code; secret values are prohibited; declared permissions/capabilities do not create authorization; stable schema/version/owner/dependency semantics remain the only valid manifest contract.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

Canonical completion evidence comes from GitHub-hosted `ubuntu-24.04`. All completed P01/P02/P03.01 regressions remain mandatory until a future governed change explicitly replaces them.

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

After a closure activates a new package or phase, identify the next authorized action and **STOP**; do not implement newly activated scope in the same execution session.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement unactivated future-phase scope; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3.

## Exact next action

After this P03.01 closure/P03.02 activation PR passes its exact-head canonical governance and merges, verify protected `main` and canonical `STATE.json`, confirm P03.01 is done and P03.02 is the sole active implementation-authorized package, then **STOP this closure execution session**.

P03.02 implementation starts only in a later governed session from the then-current exact protected-main SHA. Never implement P03.03+ from this transition.
