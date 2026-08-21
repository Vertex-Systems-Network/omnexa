# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## Mission and execution state

Omnexa is a **Composable Enterprise Business Operating System**. `docs/roadmap/STATE.json` is the machine-readable execution source of truth. Only the package marked `active` is authorized except narrowly scoped dependency-blocking or repository-maintenance fixes.

During P00, kernel and business-feature implementation remain forbidden. Architecture, governance, specifications and narrow maintenance only.

## Mandatory read order

Before material work, read:

1. `AGENTS.md`
2. product constitution + system/module architecture
3. glossary, naming, domain ownership and dependency matrix
4. identifier/money/time/locale/error standards
5. API and Event standards
6. Security Standard, Data Classification, Security Control Matrix and Threat Model
7. Testing, CI, Release and Quality Gate standards
8. repository/local-development/toolchain/configuration/developer-command standards
9. SLO, Incident and Reliability standards
10. roadmap `MASTER_PLAN.md`, `STATUS.md`, `STATE.json`
11. AI Execution Policy, Change Control, Definition of Done
12. active CI exception if present
13. relevant accepted ADRs.

If canonical documents conflict, stop implementation and resolve through change control.

## Temporary hosted-CI exception

When `github_actions_ci.state = temporary_p00_exception`, ADR-0006 applies. Hosted Actions is `BLOCKED`/`NOT RUN`, never `PASS`. Manual evidence may advance **P00 documentation/specification-only** work. It cannot authorize P01+ executable merges or waive build, migration, test, security-scan or release gates. The exception expires at P00 exit/before P01 implementation.

## Architecture invariants

1. Kernel before business modules.
2. One authoritative owner per write model/capability.
3. Cross-module direct DB writes and private implementation imports are forbidden.
4. Cross-domain integration uses governed capabilities/APIs, events, workflows or approved read projections.
5. Tenant/org boundaries, authorization, audit, observability and versioned contracts are mandatory.
6. Optional-module failure/removal cannot corrupt unrelated domains.
7. Retriable jobs/events/integrations are idempotent where required.
8. AI acts through governed capabilities only; no unrestricted DB authority.
9. Strict modular monolith first; service extraction requires evidence + ADR.
10. Infrastructure complexity must be earned.

## Foundation primitives — P00.03

- UUIDv7 canonical IDs; PostgreSQL native `uuid`;
- tenant-owned state uses `tenant_id` where applicable;
- exact-decimal money + explicit currency; no binary floating point;
- UTC/`timestamptz` instants + IANA civil-time semantics;
- BCP 47 locale with independent country/currency/timezone and first-class RTL;
- stable safe structured errors with no secrets/SQL/stack traces.

## HTTP APIs — P00.04

- `/api/v{major}/{domain}/{resources}`;
- OpenAPI 3.2.0 canonical contracts;
- lowercase `snake_case` JSON;
- Problem Details + stable machine codes;
- cursor pagination;
- explicit idempotency/concurrency;
- explicit business actions;
- client tenant/org IDs never become authorization authority.

## Events — P00.05

- `<domain>.<subject>.<past_tense_fact>.v<major>`;
- producer-owned CloudEvents-compatible envelope with UUIDv7 identity;
- trusted producer-derived tenant context;
- at-least-once delivery + idempotent consumers;
- outbox/inbox or equivalent durable guarantees;
- no global ordering assumption;
- bounded retry/DLQ;
- replay preserves identity and cannot duplicate protected side effects.

## Security — P00.06

Data classes: `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`.

- tenant-owned business data defaults to at least `CONFIDENTIAL` unless explicitly published;
- secrets/private keys/auth equivalents are `RESTRICTED`;
- production sensitive data does not flow to ordinary dev/test;
- authn never substitutes for authz;
- authz = RBAC + relationships + contextual policy + bounded capabilities;
- tenant isolation spans OLTP/cache/files/search/analytics/events/backups/AI/vector data;
- privileged support/export/purge/replay/financial/high-impact AI actions require explicit policy/audit;
- webhooks/integrations/devices/modules/AI are independent trust boundaries;
- prompt/tool injection cannot create authority.

## Quality — P00.07

Gate classes: `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package, `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never treat blocked/unrun/N/A as PASS.

Affected risks require positive and negative evidence: cross-tenant denial, authz deny, duplicate/idempotency/replay, fresh+upgrade migrations, optional-module lifecycle/degradation and security boundaries. Flaky tests are defects. Releases prefer immutable build-once/promote artifacts with source SHA and gate evidence.

## Repository and local development — P00.08

Canonical monorepo roots: `apps/`, `kernel/`, `modules/`, `platform/`, `shared/`, `infrastructure/`, `scripts/`, `docs/`, `generated/`.

- folder != microservice; extraction requires evidence/ADR;
- module private implementation/schema/migrations stay with owner;
- generated output is reproducible derivative output, never source of truth;
- default local infra = PostgreSQL + Redis-compatible + NATS/JetStream + S3-compatible storage in containers;
- Kubernetes is not required for ordinary local development;
- local and CI use the same semantic verification rules;
- toolchains/dependencies are repository-pinned;
- secrets are separate from committed config;
- production sensitive data is prohibited by default locally;
- use synthetic deterministic multi-tenant fixtures;
- Linux is canonical backend environment; Windows backend development prefers WSL2; native Windows is a separate certification target where required;
- supported workflows must not depend on hidden manual SQL/file/UI steps.

## Threat and operational reliability — P00.09

The foundation threat model is a minimum for every future phase. Each new material trust boundary/provider/privileged capability requires a delta threat model.

Required threat classes include tenant escape, object/relationship authz bypass, session takeover, privilege escalation, injection, SSRF, webhook spoofing/replay, event/job replay, financial duplication/integrity loss, module/supply-chain compromise, CI/release credential theft, POS/edge compromise, backup/export/search/vector leakage, AI prompt/tool abuse, insider/support misuse, noisy-neighbor/DDoS/provider outage, migration corruption, region failure, audit tampering, secret exposure and misconfiguration.

Operational criticality: `TIER_0` integrity-critical, `TIER_1` core transaction, `TIER_2` interactive supporting, `TIER_3` optional/background. Initial mature-production availability objectives: 99.99%, 99.95%, 99.9%, 99.5% respectively.

Recovery classes:

- A: target RPO <= 5m, RTO <= 30m;
- B: target RPO <= 15m, RTO <= 2h;
- C: target RPO <= 24h, RTO <= 8h;
- D: rebuild-based derived state.

These are targets until recovery rehearsal proves them.

Zero-tolerance conditions include cross-tenant disclosure, unauthorized privileged mutation, duplicate protected financial side effects, material financial/ledger integrity violation and lost acknowledged durable work. Error budgets never excuse these conditions.

Incident model is `SEV0`–`SEV3`. Security/privacy/integrity can outrank availability severity. Material production capabilities eventually require owner, SLI/SLO, dependencies, degradation behavior, observability, runbook, recovery treatment, alert ownership and threat-model delta.

Fail-closed rules: never bypass authorization because a dependency is down; never treat unverified payment as success; never drop protected audit; never silently lose durable work.

## Technology baseline

Go backend/core; TypeScript+React web/admin/builder/SDK; Rust justified edge/native/security; Python justified AI/data; PostgreSQL OLTP; Redis-compatible cache; S3-compatible storage; NATS/JetStream-class messaging; OpenTelemetry observability.

## Required work protocol

For every material change:

1. identify active phase/package and acceptance criteria;
2. inspect repository state;
3. identify canonical terminology and authoritative owner;
4. apply primitive/API/event/security/quality/development/operations standards;
5. map risk to G0-G8 gates and operational criticality/recovery class where relevant;
6. preserve repository/module boundaries;
7. implement the smallest complete authorized change;
8. add required positive + negative evidence and threat-model delta where needed;
9. execute canonical verification unless an active P00 exception records hosted execution `BLOCKED/NOT RUN`;
10. record evidence accurately and reconcile STATUS/STATE;
11. document architectural changes through ADR/change control before implementation.

## Forbidden behavior

Do not start future phases; silently add domains; duplicate ownership; invent conflicting contracts/security/quality/toolchain/SLO semantics; cross-write/import module-private internals; bypass tenancy/authz/audit/classification; grant AI private write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; call blocked CI PASS; claim untested RPO/RTO as achieved; spend availability error budget on security/integrity violations; require Kubernetes for ordinary local dev without ADR; or mix unrelated project code.

## Change-control triggers

ADR + plan reconciliation is required for material changes to product boundary, domain ownership, primitive/API/event/security semantics, quality/evidence/release semantics, monorepo/toolchain/config model, threat model baseline, criticality/SLO/RPO/RTO/error-budget/incident semantics, module lifecycle/dependencies, deployment topology or phase gates.

## Pull request evidence

Every material PR states phase/package, scope/non-scope, owners/dependencies, primitive/API/event/security impacts, quality gate mapping, operational criticality/recovery impacts where relevant, repository/toolchain/config impacts, data/migrations, exact evidence states, compatibility/rollback and state reconciliation.

## Hosted/legal guardrails

- issue #3: `main` protection still needs hosted configuration;
- issue #14: GitHub Actions quota/runner currently blocked;
- issue #4: final licensing/IP/trademark strategy remains unresolved before external distribution/public launch.

P00.10 must classify remaining blockers before P01. None authorize early implementation.

## Scope drift

Useful work outside the active package is recorded for later or proposed through change control. The objective is a platform that can grow without architectural decay.