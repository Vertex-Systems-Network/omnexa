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
P02.01-P02.10: DONE
Current work package: NONE
P03: PLANNED — NOT ACTIVATED
kernel_code_authorized: false
business_feature_code_authorized: false
```

There is no active implementation package at the terminal P02 checkpoint. P02 completion does not authorize P03. The next possible governed work is P03 specification/readiness preparation and a separate explicit activation transition; implementation remains locked until that transition is accepted on protected `main`.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.10.md` before material work. These files are subordinate snapshots/indexes and never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or canonical GitHub evidence.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. P01/P02 entry/exit gates and transition checklist;
4. `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.10.md` and P02.10 completion evidence;
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

P01.01-P01.12 remain `done`, P01 exit remains **SATISFIED**, and all P01 regression verifiers remain mandatory. Final P01 evidence remains implementation PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, run/job `32629072886 / 97168916985`, merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

## Completed P02 retained

P02.01-P02.10 are `done` and `docs/governance/P02_EXIT_GATE.md` is **SATISFIED**. The terminal implementation evidence is:

- implementation PR: #88
- final exact head: `975e4925060a035780ca13b68c5437634ed0f4ea`
- canonical run/job: `32904678957 / 97986011269` — PASS
- implementation merge: `88799aa41da8ce8c22540146d157d488565e2ce9`
- runner: `GitHub Actions 1000023269`, GitHub-hosted Ubuntu 24.04.4 LTS / X64
- image: `ubuntu-24.04 / 20260816.277.1`, Go 1.26.7, PostgreSQL 18.6, Valkey 9.1.1, S3 mock 5.1.0
- repository Go quality, P01.01-P01.12, `omnexa db migrate`, `omnexa verify all`, P02.01-P02.09 and P02.10 G0-G8: PASS
- completion evidence: `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`

Diagnostic run `32903969206 / 97983773781` remains explicit FAIL evidence for the corrected undefined helper reference and is not acceptance evidence.

P02 retained invariants remain binding: User is not business Person; trusted tenant context comes from authoritative tenant/membership state; no global tenant fallback; organization hierarchy is tenant-contained and Organization is not business Party Organization; authentication/strong authentication/credential possession prove identity rather than authority; secrets remain non-disclosable; sessions are revocable and current context is reauthorized; authorization is deny-by-default with exact trusted scope and no role-name/internal/background bypass; contextual policy narrows only; service accounts are distinct non-human principals with revocable/rotatable exact-scope credentials; settings cannot create authority and have no global/user override; audit is classification-safe, separate from ordinary logs and contains no credentials/authentication material; required-audit protected mutations cannot silently claim success when audit delivery fails.

## Terminal P02 lock

At this checkpoint:

- `kernel_code_authorized=false`;
- `business_feature_code_authorized=false`;
- there is no active package;
- P03 remains planned;
- P03 module runtime, business features, P04 workflows/events, generic audit UI/export, support impersonation product behavior and AI/model/agent runtime are not authorized.

A later preparation/readiness session may define governed P03 package specifications and entry controls. That planning work must not be confused with implementation authority.

## Completed kernel capability rules retained

Protected audit remains separate from ordinary logs. P01.11 audit is immutable/tamper-evident, classification-aware and append-oriented; required-audit failure cannot silently claim success and audit write does not imply read/export authority. P01.10 configuration flags cannot grant authority. P01.09 jobs remain non-authoritative. P01.08 diagnostics remain operational evidence rather than authority. Cache/storage/observability remain infrastructure primitives without tenancy/authorization authority. The developer CLI remains convenience tooling only.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

Canonical completion evidence comes from GitHub-hosted `ubuntu-24.04`. All completed P01/P02 regressions remain mandatory until a future governed change explicitly replaces them.

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

## Exact next transition

This closure records P02 as DONE 10 / 10 and P02 exit as SATISFIED, with no active implementation package. After the closure PR merges, verify protected `main` and canonical `STATE.json`, confirm P03 remains planned and both implementation locks remain false, then **STOP**.

P03 specification/readiness preparation and explicit activation belong to a later governed execution session. Never implement P03 before that separate activation is accepted.
