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
P01.03: DONE — Structured error & result conventions
P01.04: DONE — PostgreSQL connection & migration foundation
P01.05: DONE — Cache abstraction
P01.06: DONE — Object & file storage abstraction
P01.07: DONE — Structured logging & OpenTelemetry baseline
P01.08: ACTIVE — Health, readiness & diagnostics
P01.09-P01.12: PLANNED
kernel_code_authorized: true
business_feature_code_authorized: false
```

Kernel authorization is bounded to the sole active package. It is not permission to implement P01.09+, P02+, module runtime or business features.

## Persistent AI continuity

A new AI session must use `docs/ai/` as the durable continuity/handoff index after verifying canonical state. Read `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current `docs/ai/handoffs/<WORK_PACKAGE>.md` before material work. These files are subordinate snapshots/indexes: they never override this contract, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, accepted ADRs or other canonical governance. If continuity content disagrees with authoritative repository/GitHub evidence, treat the continuity value as **STALE**, re-verify the authoritative sources and do not silently overwrite canonical state.

## Mandatory read order

Before material work read:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json` and `docs/roadmap/STATUS.md`;
3. `docs/governance/FOUNDATION_FREEZE.json` and `docs/governance/P01_ENTRY_GATE.md`;
4. `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`;
5. the active package specification (`P01.08.md` currently);
6. Product Constitution, system/module architecture, glossary, naming, ownership and dependency matrix;
7. identifier/money/time/locale/error/API/event standards;
8. security/data-classification/threat model;
9. testing/CI/release/quality standards, including `docs/quality/GO_CODE_QUALITY.md`;
10. repository/local-development/toolchain/configuration/developer-command standards;
11. `docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` whenever browser UI is in an authorized package;
12. SLO/incident/reliability standards;
13. AI Execution Policy, Change Control and Definition of Done;
14. relevant ADRs, especially ADR-0010.

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

Issue #3 / EG-02 is satisfied. `main` is protected with PR-only integration, strict required `governance`, blocked direct/force updates, failed-check merge rejection and required conversation resolution.

Canonical governance CI is **GitHub-hosted only**:

```yaml
runs-on: ubuntu-24.04
```

The job must fail closed unless `RUNNER_ENVIRONMENT=github-hosted`, `RUNNER_OS=Linux` and `RUNNER_ARCH=X64`. Do not reintroduce `self-hosted`, `LOCAL-WIN-*`, local evidence fanout or local-runner fallback.

The permanent repository Go quality gate runs through `bash scripts/verify_go_quality.sh` with pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`. Do not remove it from required governance, weaken configured checks merely to obtain green CI, use `@latest`, or silently auto-fix source in CI.

Strict protection requires an implementation/closure PR to be current with protected `main` before merge. Stale green runs are not merge permission; synchronize the branch and obtain a fresh full green lane.

## P01 execution rule

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` enforces a completed prefix, exactly one active package and a planned suffix.

Completed:

- P01.01 — evidence: `docs/roadmap/evidence/P01.01_COMPLETION_2026-08-22.md`;
- P01.02 — evidence: `docs/roadmap/evidence/P01.02_COMPLETION_2026-08-22.md`;
- P01.03 — evidence: `docs/roadmap/evidence/P01.03_COMPLETION_2026-08-22.md`;
- P01.04 — evidence: `docs/roadmap/evidence/P01.04_COMPLETION_2026-08-22.md`;
- P01.05 — evidence: `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`;
- P01.06 — evidence: `docs/roadmap/evidence/P01.06_COMPLETION_2026-08-22.md`;
- P01.07 — evidence: `docs/roadmap/evidence/P01.07_COMPLETION_2026-08-22.md`.

Current active package: **P01.08 — Health, readiness & diagnostics**.

Allowed executable scope is limited to:

- semantically distinct liveness/readiness primitives;
- dependency check registry with criticality classification;
- bounded timeout/cancellation behavior;
- startup/readiness transition model;
- safe diagnostic summary with stable machine states;
- build/version identity integration;
- integration with the completed P01.07 observability boundary;
- deterministic healthy/degraded/unready test harness;
- completed P01.01-P01.07 regression verification;
- permanent repository Go quality verification;
- GitHub-hosted G0/G1/G2/G3/G5/G6/G7 evidence.

P01.08 must not implement a public business status page, P03 tenant/module health aggregation, SLO alerting/orchestration automation, Kubernetes-specific architecture, privileged sensitive diagnostics, scheduler primitives, feature registry, audit transport, developer CLI behavior, identity/tenancy/module/event/workflow/business behavior or AI/model/agent/planner functionality.

P01.09 becomes active only after P01.08 reaches `done` with required evidence. More than one active P01 package is forbidden.

## Business-feature lock

`business_feature_code_authorized=false` remains mandatory for all P01 work. Do not implement CRM, ERP, commerce, payment, POS, CMS, portal, HR/projects, supply chain, integrations, builders, BI, AI-agent business behavior or any other business domain.

## Quality and release rules

Gate classes remain `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package and `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never relabel blocked/unrun/N/A as PASS. Flaky tests are defects.

For P01.08, required executable evidence is G0/G1/G2/G3/G5/G6/G7 plus repository Go quality and completed P01.01-P01.07 regression verification. Tenancy, event replay and module lifecycle remain N/A until their owning capabilities exist.

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

## P01.08 health/readiness rules

- Liveness and readiness are distinct operational semantics; do not collapse them into one generic health boolean.
- Dependency checks must be timeout-bounded and classify required versus optional/conditionally relevant dependencies explicitly.
- Optional dependency degradation must not automatically kill the process; required security-critical capability failure must fail closed where applicable.
- Diagnostic output must not expose connection strings, host secrets, credentials, SQL, object keys or sensitive payloads.
- Diagnostic states are operational signals, not authorization, tenancy or business-state authority.
- Use the completed P01.07 observability boundary; do not create a parallel logging/tracing system.
- Do not implement public status pages, module/tenant health aggregation, SLO automation or Kubernetes-specific architecture by anticipation.
- P01.08 must not pull forward scheduler/feature-registry/audit/CLI or later-phase behavior.

## Completed P01.07 observability rules retained

Telemetry remains diagnostic infrastructure and not a correctness dependency or business source of truth. Secrets, auth tokens, private keys, raw `RESTRICTED` payloads and ordinary sensitive business content remain prohibited from logs/telemetry by default. Correlation/trace IDs are identifiers, not authorization credentials. Exporter failure and shutdown/flush remain bounded. Structured logging, traces and metrics preserve vendor-neutral boundaries. Audit semantics remain owned by P01.11. P01.07 regression verification stays mandatory.

## Completed P01.06 storage rules retained

Storage remains an infrastructure primitive below capability/domain boundaries. Object keys do not imply authorization or tenancy. Path/file-name input cannot escape the governed provider namespace. Provider credentials remain `RESTRICTED`; signed URLs, secrets and sensitive object content are not emitted by diagnostics. Content type, file name, length and caller metadata remain untrusted. Large object flows stream with bounded memory. Integrity mismatch fails safely. Missing object, provider failure, timeout and caller cancellation remain distinguishable. P01.06 regression verification stays mandatory.

## Completed P01.05 cache rules retained

Cache remains non-authoritative. Flush, restart, eviction or cache loss cannot corrupt protected state. Cache provider credentials remain sensitive, cache miss remains distinct from provider failure, and keys remain namespaced/versioned without invented tenancy. P01.05 regression verification stays mandatory.

## Future browser UI quality rule

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` is a mandatory planning input whenever a future package authorizes browser UI.

Future production browser UI must target WCAG 2.2 AA and standards-based semantic HTML/CSS, with W3C validation, WAVE evaluation and manual keyboard/focus/screen-reader/zoom-reflow checks appropriate to the affected surface.

AI systems must not use ARIA as a mechanical substitute for native semantics, disable validators/accessibility checks merely to obtain green CI, claim WAVE/WCAG/W3C compliance solely from automated scans, silently ignore WAVE Errors/Contrast Errors/Alerts, or expose WAVE secrets.

A required WAVE automation dependency without an approved key/license is `BLOCKED`, not PASS. This future UI rule is planning only during P01.

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
9. reconcile state/status/ledger only after completion evidence exists;
10. use ADR/change control before changing frozen architecture.

## Forbidden behavior

Do not use local/self-hosted runners for canonical governance; silently add domains; duplicate ownership; invent conflicting contracts/security/quality semantics; bypass tenancy/authz/audit/classification; grant AI unrestricted write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; claim untested evidence; implement later P01/P02/P03/business scope early; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Issue #4

Issue #4 remains the external distribution/public-launch licensing/IP/trademark gate. Repository visibility is public and current `LICENSE` remains GPLv3. It does not block bounded P01 kernel engineering.

## Exact next transition

Implement P01.08 in a separate executable PR. Add the package-specific fail-closed verification required by `docs/roadmap/work-packages/P01.08.md`, obtain G0/G1/G2/G3/G5/G6/G7 GitHub-hosted evidence plus repository Go quality and P01.01-P01.07 regressions, move P01.08 `active -> verification -> done`, reconcile canonical state, then activate only P01.09.
