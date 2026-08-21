# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## 1. Mission

Omnexa is a **Composable Enterprise Business Operating System**, not a conventional ERP or a collection of unrelated applications. ERP, CRM, commerce, POS, website/CMS, portals, payments, workflow, integrations, analytics, low-code and AI are governed domains running on one platform foundation.

The repository must evolve as one coherent platform with independently installable modules and explicit ownership boundaries.

## 2. Mandatory read order

Before any material change, read:

1. `AGENTS.md`
2. `docs/governance/PRODUCT_CONSTITUTION.md`
3. `docs/architecture/SYSTEM_ARCHITECTURE.md`
4. `docs/architecture/MODULE_STANDARD.md`
5. `docs/architecture/GLOSSARY.md`
6. `docs/architecture/NAMING_STANDARD.md`
7. `docs/architecture/DOMAIN_OWNERSHIP.md`
8. `docs/architecture/DEPENDENCY_MATRIX.md`
9. `docs/architecture/IDENTIFIER_STANDARD.md`
10. `docs/architecture/MONEY_STANDARD.md`
11. `docs/architecture/TIME_STANDARD.md`
12. `docs/architecture/LOCALE_STANDARD.md`
13. `docs/architecture/ERROR_STANDARD.md`
14. `docs/architecture/API_STANDARD.md`
15. `docs/architecture/EVENT_STANDARD.md`
16. `docs/security/SECURITY_STANDARD.md`
17. `docs/security/DATA_CLASSIFICATION.md`
18. `docs/security/SECURITY_CONTROL_MATRIX.md`
19. `docs/quality/TESTING_STANDARD.md`
20. `docs/quality/CI_STANDARD.md`
21. `docs/quality/RELEASE_STANDARD.md`
22. `docs/quality/QUALITY_GATE_MATRIX.md`
23. `docs/roadmap/MASTER_PLAN.md`
24. `docs/roadmap/STATUS.md`
25. `docs/roadmap/STATE.json`
26. `docs/governance/AI_EXECUTION_POLICY.md`
27. `docs/governance/CHANGE_CONTROL.md`
28. `docs/governance/DEFINITION_OF_DONE.md`
29. if active, `docs/governance/CI_EVIDENCE_EXCEPTION_2026-08-22.md`
30. relevant accepted ADRs.

If canonical documents conflict, stop implementation and resolve the conflict through change control.

## 3. Canonical execution state

`docs/roadmap/STATE.json` is the machine-readable execution source of truth. Only the work package marked `active` is authorized except a narrowly scoped dependency-blocking defect or repository-maintenance correction.

During P00:

- kernel implementation is forbidden;
- business-feature implementation is forbidden;
- only architecture, governance, specifications and narrow maintenance are permitted.

Never claim progress beyond repository evidence.

### Temporary hosted-CI exception

When `STATE.json` declares `github_actions_ci.state = temporary_p00_exception`, ADR-0006 governs evidence handling. Hosted Actions is `BLOCKED`/`NOT RUN`, never fabricated as `PASS`. Manual evidence may advance **P00 documentation/specification-only** work when the exception criteria are satisfied. The exception cannot authorize P01+ executable merges or waive build, migration, test, security-scan or release gates.

## 4. Core architecture invariants

Unless superseded by accepted ADR + plan reconciliation:

1. Kernel before business modules.
2. Every authoritative write model/capability has one owner.
3. Modules own their private write models/schemas; cross-module direct DB writes are forbidden.
4. Cross-domain integration uses governed capabilities/APIs, events, workflows or approved read projections.
5. Dependency direction follows `DEPENDENCY_MATRIX.md`.
6. Tenant and organization boundaries are explicit.
7. Every protected action is authorization-aware.
8. Business/security-significant mutations are auditable.
9. Public/external contracts are versioned.
10. Optional-module failure/removal must not corrupt unrelated domains.
11. Retriable jobs/events/integrations are idempotent where required.
12. AI acts only through governed capabilities and never unrestricted database authority.
13. Start as a strict modular monolith; extract services only when evidence justifies it.
14. Infrastructure complexity must be earned by requirements/evidence.

## 5. Foundation primitive invariants — P00.03

- canonical IDs use UUIDv7; PostgreSQL uses native `uuid`;
- tenant-owned persisted state uses `tenant_id` where applicable;
- money uses exact decimal + explicit currency; binary floating point is forbidden for money;
- JSON monetary/rate decimals are strings; high-precision baseline is `NUMERIC(38,18)`;
- absolute instants use UTC/`timestamptz`; civil time carries an IANA timezone;
- locales use BCP 47; locale/country/currency/timezone are separate; RTL is first-class;
- public errors use stable machine codes + safe structured problem details and expose no secrets/SQL/stack traces.

## 6. HTTP API invariants — P00.04

- stable routes: `/api/v{major}/{domain}/{resources}`;
- OpenAPI 3.2.0 is canonical;
- JSON fields use lowercase `snake_case`;
- errors use `application/problem+json` + stable Omnexa codes/request IDs;
- scalable lists use opaque cursor pagination;
- retriable mutations define `Idempotency-Key` behavior;
- lost-update-sensitive writes use explicit optimistic concurrency where applicable;
- business lifecycle actions are explicit capability/action operations;
- client tenant/org identifiers never become authorization authority;
- generated SDK/routes/docs derive from governed contracts.

## 7. Event invariants — P00.05

- event type: `<domain>.<subject>.<past_tense_fact>.v<major>`;
- authoritative producer owns meaning/schema/publication condition;
- envelope is CloudEvents-compatible with UUIDv7 identity;
- tenant-owned events carry trusted producer-derived tenant context;
- delivery baseline is at least once; consumers are idempotent;
- business-significant publication/consumption uses outbox/inbox or equivalent durable guarantees;
- no global ordering guarantee;
- retries are bounded; poison failures use governed dead-letter/quarantine handling;
- replay preserves original identity/payload and cannot duplicate protected business side effects.

## 8. Security invariants — P00.06

Data classes: `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`.

- tenant-owned business records default to at least `CONFIDENTIAL` unless explicitly published;
- encryption does not lower classification;
- secrets/auth equivalents/private keys are `RESTRICTED`;
- `RESTRICTED` data is excluded from ordinary logs/search/analytics/generic exports and AI model input by default;
- production sensitive data does not flow to dev/test by default;
- authentication does not substitute for authorization;
- authorization combines RBAC + relationships + contextual policy + bounded capabilities;
- tenant isolation spans OLTP/cache/files/search/analytics/events/backups/AI/vector data;
- cross-tenant operations require explicit privileged capability and audit;
- modules do not invent custom cryptography;
- privileged support/export/purge/replay/financial/high-impact AI actions require explicit policy/audit and stronger controls according to risk;
- webhooks/integrations/devices/modules/AI are independent trust boundaries;
- prompt/tool injection cannot create authority.

## 9. Quality invariants — P00.07

Quality semantics are repository-owned and **CI-provider independent**. Local development and CI must execute the same underlying verification rules.

### Gate classes

```text
G0 Governance
G1 Static
G2 Unit / Component
G3 Contract / Integration
G4 Data / Migration
G5 Security / Tenancy
G6 Lifecycle / Resilience
G7 Build / Package
G8 Supply Chain / Release
```

Allowed evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`.

Rules:

- `BLOCKED`, `NOT RUN` and `N/A` are never silently equivalent to PASS;
- every change maps risk to the applicable G0-G8 gates;
- tenant-owned paths require positive same-scope plus negative cross-scope tests where affected;
- authorization changes require allow + deny evidence;
- retriable/event/async behavior requires duplicate/idempotency/replay evidence;
- persistence changes require fresh-install + supported-upgrade migration evidence;
- optional modules require install/enable/disable/re-enable/upgrade/degradation evidence;
- flaky tests are defects and quarantine requires owner, issue and expiry;
- aggregate coverage never substitutes for high-risk invariant coverage;
- production releases prefer build-once/promote immutable artifacts;
- release identity must trace to source SHA, artifact digest and applicable gate evidence;
- P00 temporary CI exceptions cannot authorize P01+ executable merges/releases.

Future tooling must expose stable local verification commands equivalent to the canonical `verify:*` command family defined in `CI_STANDARD.md`.

## 10. Technology baseline

- Go — kernel/backend and primary domain services;
- TypeScript + React — admin/web/builder/SDK surfaces;
- Rust — edge/native/security-sensitive components where justified;
- Python — AI/data workloads where ecosystem value justifies it;
- PostgreSQL — primary OLTP;
- Redis-compatible — cache/ephemeral coordination;
- S3-compatible — files/media;
- NATS/JetStream-class — event/messaging baseline;
- OpenTelemetry — observability semantics.

Technology choice never authorizes premature phase work.

## 11. Required work protocol

For every material change:

1. identify phase/work package and acceptance criteria;
2. inspect current repository state;
3. identify canonical terminology and authoritative owner;
4. apply primitive/API/event/security/quality standards at affected boundaries;
5. map affected risk to G0-G8 quality gates;
6. state ownership/dependency/security/tenant implications;
7. implement the smallest complete authorized change;
8. add positive + required negative tests;
9. execute required canonical verification commands, except only where an active P00 exception records hosted execution `BLOCKED/NOT RUN`;
10. record evidence accurately;
11. reconcile STATUS/STATE only when evidence supports transition;
12. add/update ADR + dependent docs before architectural change implementation.

## 12. Forbidden behavior

Do not:

- start future phases early;
- silently add domains or duplicate ownership;
- invent conflicting primitive/API/event/security/quality semantics;
- cross-write module-private schemas;
- bypass tenancy, authorization, audit, classification or versioning;
- grant AI broader/private write authority;
- log/store secrets for convenience;
- use production sensitive data in lower environments without approved exception;
- introduce hidden super-admin/support bypasses;
- claim `done` without evidence;
- call blocked/unrun CI PASS;
- use ADR-0006 as precedent for executable P01+ work;
- weaken required G0-G8 gates to make a PR green;
- mix unrelated project code into the repository.

## 13. Change-control triggers

ADR + plan reconciliation is required for material changes to product boundary, domain ownership, primitives, stable API/event semantics, security/tenant/authorization model, classification/secrets/audit/AI authority, quality gate taxonomy/evidence semantics, release artifact identity/promotion, language/runtime baseline, module lifecycle/dependencies, deployment topology or phase ordering/gates.

Never implement the architecture change first and document it later.

## 14. Pull request requirements

Every material PR states phase/work package, scope/non-scope, architecture impact, owner/dependency mechanism, primitive/API/event impacts, security/tenancy/classification impact, quality gate mapping, data/migration impact, exact evidence state, compatibility/rollback considerations and status/state reconciliation.

Unrelated work belongs in separate PRs.

## 15. Safe completion

A change is incomplete if an applicable required gate fails. Infrastructure-blocked execution remains `BLOCKED`; it is not a pass. P00.07 is the canonical quality baseline for future implementation.

## 16. Repository and legal guardrails

- use governed feature branches/PRs rather than intentional direct `main` changes;
- follow `CONTRIBUTING.md` and `SECURITY.md`;
- branch protection target is tracked by issue #3;
- GitHub Actions quota/runner blocker is tracked by issue #14;
- final licensing/IP/trademark strategy is tracked by issue #4;
- AI systems do not make legal/trademark ownership decisions autonomously.

## 17. Scope drift

Useful work outside the active package is recorded for later or proposed through plan change; it is not silently absorbed.

The objective is a platform that can grow to very high feature count **without architectural decay**.