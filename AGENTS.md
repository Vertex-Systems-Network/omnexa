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
P01.01: DONE — Go workspace/build skeleton
P01.02: DONE — Configuration & environment system
P01.03: ACTIVE — Structured error & result conventions
P01.04-P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```

Kernel authorization is bounded to the sole active package. It is not permission to implement P01.04+, P02+, module runtime or business features.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/FOUNDATION_FREEZE.json` and `docs/governance/P01_ENTRY_GATE.md`;
4. `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`;
5. the active package specification (`P01.03.md` currently);
6. Product Constitution, system/module architecture, glossary, naming, ownership and dependency matrix;
7. identifier/money/time/locale/error/API/event standards;
8. security/data-classification/threat model;
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
7. Retriable work is idempotent where required.
8. AI acts only through governed capabilities; no unrestricted DB authority.
9. Strict modular monolith first; service extraction requires evidence + ADR.
10. Infrastructure complexity must be earned.

Frozen primitives include UUIDv7 IDs, exact-decimal money with explicit currency, UTC/timestamptz instants with IANA civil-time semantics, BCP 47 locale/RTL support, stable safe structured errors, versioned HTTP/OpenAPI contracts, CloudEvents-compatible event envelopes, at-least-once/idempotent event handling, four data classes (`PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`) and deny-by-default authorization/tenant isolation.

## Protected integration and CI

Issue #3 / EG-02 is satisfied. `main` is protected with PR-only integration, required `governance`, blocked direct/force updates, failed-check merge rejection and required conversation resolution.

Canonical governance CI is **GitHub-hosted only**:

```yaml
runs-on: ubuntu-24.04
```

The job must fail closed unless `RUNNER_ENVIRONMENT=github-hosted`, `RUNNER_OS=Linux` and `RUNNER_ARCH=X64`. Do not reintroduce `self-hosted`, `LOCAL-WIN-*`, local evidence fanout or local-runner fallback.

## P01 execution rule

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` enforces a completed prefix, exactly one active package and a planned suffix.

Completed:

- P01.01 — evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`;
- P01.02 — evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`.

Current active package: **P01.03 — Structured error & result conventions**.

Allowed executable scope is limited to:

- stable machine error codes;
- safe public message/detail separated from private causes;
- wrapping/unwrapping preserving standard Go error semantics;
- explicit retryability/category metadata;
- deterministic bounded validation-field errors;
- correlation metadata hooks without telemetry emission;
- negative redaction/security tests;
- completed-package regression verification;
- hosted CI evidence for these boundaries.

P01.03 must not implement HTTP transport adapters, PostgreSQL/provider-specific mapping, cache/storage behavior, logging/OpenTelemetry emission, health endpoints, jobs, feature flags, audit transport, identity/tenancy, module runtime or business-domain error catalogs.

P01.04 becomes active only after P01.03 reaches `done` with required evidence. More than one active P01 package is forbidden.

## Business-feature lock

`business_feature_code_authorized=false` remains mandatory for all P01 work. Do not implement CRM, ERP, commerce, payment, POS, CMS, portal, HR/projects, supply chain, integrations, builders, BI, AI-agent business behavior or any other business domain.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

For P01.03, required executable evidence is G0/G1/G2/G5/G7 plus completed P01.01/P01.02 regression verification. Persistence, tenancy, event replay and module lifecycle remain N/A until their owning capabilities exist.

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

1. verify active phase/package and locks in `STATE.json`;
2. inspect active package spec and frozen standards;
3. preserve ownership/dependency boundaries;
4. implement only authorized scope;
5. add positive/negative evidence appropriate to risk;
6. run canonical GitHub-hosted `governance`;
7. inspect diff/status before merge;
8. merge only when required checks are green;
9. reconcile state/status/ledger only after completion evidence exists;
10. use ADR/change control before changing frozen architecture.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement later P01/P02/P03/business scope early; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3. It does not block bounded P01 kernel engineering.

## Exact next transition

Implement P01.03 in a separate executable PR. Obtain G0/G1/G2/G5/G7 evidence, move P01.03 `active -> verification -> done`, reconcile canonical state, then activate only P01.04.
