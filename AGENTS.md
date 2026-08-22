# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## Current canonical state

Omnexa is a **Composable Enterprise Business Operating System**. `docs/roadmap/STATE.json` is the machine-readable execution source of truth.

```text
Foundation Architecture v1: FROZEN
P00: DONE — 10 / 10
Repository visibility: PUBLIC
EG-02 / Issue #3: SATISFIED / CLOSED
EG-03 / Issue #14: SATISFIED
Canonical CI: GITHUB-HOSTED ONLY / ubuntu-24.04
Local/self-hosted governance runners: PROHIBITED
P01: ACTIVE
P01.01: DONE — Go workspace/build skeleton
P01.02: ACTIVE — Configuration & environment system
P01.03-P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```

Kernel implementation permission is bounded to the active work package. `kernel_code_authorized=true` is **not** permission to implement P01.03+, P02+, module runtime or business features.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/FOUNDATION_FREEZE.json` and `docs/governance/P01_ENTRY_GATE.md`;
4. `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`;
5. the active work-package specification (`P01.02.md` currently);
6. Product Constitution, system/module architecture, glossary, naming, ownership and dependency matrix;
7. primitive/API/event standards;
8. security/data classification/threat model;
9. testing/CI/release/quality standards;
10. repository/local-development/toolchain/configuration/developer-command standards;
11. SLO/incident/reliability standards;
12. AI Execution Policy, Change Control and Definition of Done;
13. relevant ADRs, especially ADR-0010.

If canonical documents conflict, resolve through change control before implementation.

## Frozen architecture laws

1. Kernel before business modules.
2. One authoritative owner per write model/capability.
3. Cross-module direct DB writes and private implementation imports are forbidden.
4. Cross-domain integration uses governed APIs/capabilities, events, workflows or approved read projections.
5. Tenant/org boundaries, authorization, audit, observability and versioned contracts are mandatory.
6. Optional-module failure/removal cannot corrupt unrelated domains.
7. Retriable jobs/events/integrations are idempotent where required.
8. AI acts through governed capabilities only; no unrestricted DB authority.
9. Strict modular monolith first; service extraction requires evidence + ADR.
10. Infrastructure complexity must be earned.

Frozen primitives include UUIDv7 IDs, exact-decimal money with explicit currency, UTC/timestamptz instants with IANA civil-time semantics, BCP 47 locale/RTL support, stable safe structured errors, versioned HTTP/OpenAPI contracts, CloudEvents-compatible event envelopes, at-least-once/idempotent event handling, four data classes (`PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`) and deny-by-default authorization/tenant isolation.

## Protected integration and CI

Issue #3 / EG-02 is satisfied. Verified behavior includes `main.protected=true`, PR-only integration, required `governance`, blocked direct/force updates, failed-check merge rejection and required conversation resolution. Branch deletion remains blocked by the configured ruleset; destructive deletion of the default branch was intentionally not performed.

Current single-maintainer review policy uses zero required approvals and no required Code Owner review. Tighten to independent approval + Code Owner review when a second reviewer exists.

Canonical governance CI is **GitHub-hosted only**:

```yaml
runs-on: ubuntu-24.04
```

The job must fail closed unless `RUNNER_ENVIRONMENT=github-hosted`, `RUNNER_OS=Linux` and `RUNNER_ARCH=X64`. Do not use or reintroduce `self-hosted`, `LOCAL-WIN-*`, local evidence fanout or local-runner fallback.

## P01 execution rule

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` enforces strict sequential one-active-package execution.

Completed package: **P01.01 — Go workspace/build skeleton**. Canonical completion evidence is `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`.

Current active package: **P01.02 — Configuration & environment system**.

Allowed executable work is limited to:

- typed configuration loading and validation;
- explicit precedence across defaults, governed config files and environment variables;
- governed environment identity;
- required/optional setting semantics;
- secret-safe redaction;
- deterministic isolated test overrides;
- configuration provenance that reveals source category/key but never secret value;
- fail-closed startup/config construction on invalid required settings;
- tests and hosted CI evidence needed for those boundaries.

P01.02 must not implement P01.03 structured application error conventions beyond narrow package-local configuration errors; PostgreSQL/migrations (P01.04); cache/storage; telemetry; health endpoints; jobs; feature flags; audit transport; full developer CLI; identity/tenancy (P02); module runtime (P03); or business configuration/features.

P01.03 becomes active only after P01.02 reaches `done` with required evidence. More than one active P01 package is forbidden.

## Business-feature lock

`business_feature_code_authorized=false` remains mandatory. Do not implement CRM, ERP, commerce, payment, POS, CMS, portal, HR/projects, supply chain, integrations, builders, BI, AI-agent business behavior or any other domain feature during P01.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects. Releases prefer immutable build-once/promote artifacts with source SHA and gate evidence.

For P01.02 required executable evidence includes G0/G1/G2/G5/G7 with explicit negative tests for invalid required configuration, invalid environment identity and secret redaction. Database migration, event replay and module lifecycle remain N/A until their owning capabilities exist.

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

## Threat/reliability rules

Every material new trust boundary/provider/privileged capability requires a threat-model delta. Fail closed: never bypass authorization because a dependency is down; never treat unverified payment as success; never drop protected audit; never silently lose durable work.

## Technology baseline

Go backend/core; TypeScript+React web/admin/builder/SDK; Rust only when justified; Python only for justified AI/data work; PostgreSQL OLTP; Redis-compatible cache; S3-compatible storage; NATS/JetStream-class messaging; OpenTelemetry observability.

The active P01 package does **not** authorize adding later infrastructure merely because it is part of the technology baseline.

## Required work protocol

For every material change:

1. verify active phase/package and locks in `STATE.json`;
2. inspect the active package spec and relevant frozen standards;
3. preserve ownership/dependency boundaries;
4. implement only authorized scope;
5. add positive/negative evidence appropriate to affected risk;
6. run the canonical GitHub-hosted `governance` lane;
7. inspect diff/status before merge;
8. merge only when required checks are green and protection rules allow it;
9. reconcile state/status/ledger only after completion evidence exists;
10. use ADR/change control before changing frozen architecture.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; cross-write/import module-private internals; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement later P01/P02/P03/business scope early; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## ADR-0006 and Issue #4

ADR-0006 is expired/historical-only and cannot authorize a current CI bypass.

Issue #4 remains an external distribution/public-launch licensing/IP/trademark gate. The repository is public and current `LICENSE` remains GPLv3. Issue #4 does not block the currently authorized P01 kernel engineering scope.

## Exact next transition

Implement P01.02 in a separate executable PR. Obtain required G0/G1/G2/G5/G7 evidence, move P01.02 `active -> verification -> done`, reconcile canonical state, then activate only P01.03.

## Scope drift

Useful work outside the active package is recorded for later or proposed through change control. The objective is a platform that can grow without architectural decay.
