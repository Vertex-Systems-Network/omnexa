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
P01.01-P01.10: DONE
P01.11: ACTIVE — Audit transport foundation
P01.12: PLANNED
P01 progress: 10 / 12 done
kernel_code_authorized: true
business_feature_code_authorized: false
```

Kernel authorization is bounded to the sole active package. It is not permission to implement P01.12, P02+, module runtime or business features.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current `docs/ai/handoffs/<WORK_PACKAGE>.md` before material work. These files are subordinate snapshots/indexes and never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or other canonical governance. If continuity disagrees with authoritative repository/GitHub evidence, mark it **STALE** and follow authoritative sources.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/FOUNDATION_FREEZE.json` and `docs/governance/P01_ENTRY_GATE.md`;
4. `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`;
5. the active package specification (`P01.11.md` currently);
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

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` enforces a completed prefix, exactly one active package and a planned suffix.

P01.01 through P01.10 are complete with canonical evidence under `docs/roadmap/evidence/`. Latest completed evidence is `docs/roadmap/evidence/P01.10_COMPLETION_2026-08-23.md`, based on implementation PR #61, exact head `4c9914e4641d0d6e94a895d0fcd16c3a6bf4d962`, final canonical run/job `32609018028 / 97118796940`, and implementation merge `9d11b9250eb74ca2ade531ee58e8f905468cf103`.

Current active package: **P01.11 — Audit Transport Foundation**.

P01.12 becomes active only after P01.11 reaches `done` with required evidence and a separate governed closure transition. More than one active P01 package is forbidden.

## P01.11 audit transport rules

Authorized P01.11 scope is limited to:

- stable classification-aware audit record envelope;
- actor/action/target/scope/outcome/correlation/reason/approval metadata without implementing P02 identities;
- append-oriented sink interface;
- explicit required-audit failure semantics;
- classification/redaction enforcement boundary;
- immutable UUIDv7/timestamp conventions;
- impersonation/privileged-action metadata representation;
- deterministic local/test sink;
- bounded transport-health observability without protected-payload logging;
- P01.01-P01.10 regression verification;
- permanent repository Go quality verification;
- GitHub-hosted G0/G1/G2/G3/G5/G6/G7 evidence.

P01.11 must not implement P02 identity/tenant/role catalogs, business-domain audit event definitions, compliance export/reporting UI, long-term legal retention/hold, replacement of domain events with audit records, P01.12 CLI, durable messaging/outbox/inbox pull-forward, business modules or AI/model/agent behavior.

Protected audit is separate from ordinary logs. Required-audit sink failure cannot silently claim success. Records must reject secrets/auth equivalents by default. Audit write capability does not imply audit read/export authority. Actor/scope identifiers are descriptive metadata and do not grant authorization, tenancy or identity authority.

## Completed P01.10 configuration-registry rules retained

Runtime configuration remains distinct from P01.02 static process configuration. Definitions are typed, stable and owner-declared. Defaults/fallbacks are deterministic. Cache/refresh/invalidation remains bounded and non-authoritative. Kill switches are explicitly declared disable-only mitigation controls and cannot grant authority or weaken authorization/data isolation. Future tenant/org/user references are opaque UUIDv7 metadata only. Generic runtime configuration is not a secrets store. P01.10 regression verification remains mandatory.

## Completed lower-package invariants retained

P01.09 jobs remain deterministic, bounded and non-authoritative with bounded retry/idempotency and graceful shutdown semantics. P01.08 health/readiness remains safe operational evidence rather than authority. P01.07 observability remains diagnostic infrastructure rather than audit/correctness authority. P01.06 storage and P01.05 cache remain infrastructure primitives without tenancy/authorization authority. Their regression verifiers remain mandatory.

## Business-feature lock

`business_feature_code_authorized=false` remains mandatory for all P01 work. Do not implement CRM, ERP, commerce, payment, POS, CMS, portal, HR/projects, supply chain, integrations, builders, BI, AI-agent business behavior or any other business domain.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

For P01.11, required executable evidence is G0/G1/G2/G3/G5/G6/G7 plus repository Go quality and completed P01.01-P01.10 regression verification. P02 tenancy/catalog behavior, durable event replay and business module lifecycle remain outside P01.11.

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

The AI must not automatically advance to the next work package. After a package closure activates the next package, identify the next authorized action and **STOP**; do not implement that next package in the same execution session.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement later P01/P02/P03/business scope early; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3. It does not block bounded P01 kernel engineering.

## Exact next transition

After the P01.10 closure PR merges, begin P01.11 only in a **new governed execution session**. Implement the bounded `kernel.audit` transport foundation, add its fail-closed package verifier, preserve P01.01-P01.10 regressions and repository Go quality, obtain exact-head GitHub-hosted G0/G1/G2/G3/G5/G6/G7 evidence, then use a separate closure transition before P01.12.
