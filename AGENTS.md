# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## 1. Mission

Omnexa is a **Composable Enterprise Business Operating System**, not a conventional ERP or a collection of unrelated applications. ERP, CRM, commerce, POS, website/CMS, portals, payments, workflow, integrations, analytics, low-code and AI are governed domains running on one platform foundation.

The repository must evolve as one coherent platform with independently installable modules and explicit ownership boundaries.

## 2. Mandatory read order

Before changing code, schema, infrastructure, APIs, events, tests or documentation, read:

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
19. `docs/roadmap/MASTER_PLAN.md`
20. `docs/roadmap/STATUS.md`
21. `docs/roadmap/STATE.json`
22. `docs/governance/AI_EXECUTION_POLICY.md`
23. `docs/governance/CHANGE_CONTROL.md`
24. `docs/governance/DEFINITION_OF_DONE.md`
25. relevant ADRs under `docs/adr/`

If canonical documents conflict, stop implementation and resolve the conflict through change control.

## 3. Canonical execution state

`docs/roadmap/STATE.json` is the machine-readable execution source of truth.

Only the work package marked `active` is authorized, except a narrowly scoped dependency-blocking defect or repository-maintenance correction.

During P00:

- kernel implementation is forbidden;
- business-feature implementation is forbidden;
- only architecture, governance, specifications and narrow maintenance are permitted.

An AI system must never claim progress beyond repository evidence.

## 4. Core architecture invariants

Unless superseded through an accepted ADR and plan reconciliation:

1. Kernel before business modules.
2. Every authoritative write model/capability has one owner.
3. A module owns its write model/schema; other modules do not write its private tables.
4. Cross-domain communication uses governed APIs/capabilities, events, workflows or approved read projections.
5. Dependency direction follows `DEPENDENCY_MATRIX.md`.
6. Tenant and organization boundaries are explicit.
7. Every protected action is authorization-aware.
8. Business/security-significant mutations are auditable.
9. Public/external contracts are versioned.
10. Optional-module failure/removal must not corrupt unrelated domains.
11. Retriable integrations/jobs/events are idempotent where required.
12. AI acts through governed capabilities and never unrestricted database authority.
13. Omnexa begins as a strict modular monolith; service extraction requires evidence.
14. Infrastructure complexity must be justified by requirements/evidence.

## 5. Foundation primitive invariants — P00.03

- canonical new entity/request/event/workflow/job/audit identifiers use UUIDv7;
- PostgreSQL uses native `uuid` for canonical identifiers;
- tenant-owned persisted state uses `tenant_id` where applicable;
- monetary values use exact decimal semantics with explicit currency; binary floating point is forbidden for money;
- JSON monetary/rate decimals are strings; high-precision PostgreSQL baseline is `NUMERIC(38,18)`;
- absolute instants use UTC/`timestamptz`; business dates remain dates; civil-time recurrence carries an IANA timezone;
- locales use BCP 47; locale/country/currency/timezone are separate concepts; RTL is first-class;
- public errors use stable machine codes plus safe structured problem details and never expose secrets, stack traces or SQL internals.

## 6. Stable HTTP API invariants — P00.04

Stable public/partner/cross-domain HTTP contracts follow `API_STANDARD.md` and ADR-0003.

- route major versions use `/api/v{major}/{domain}/{resources}`;
- OpenAPI 3.2.0 is the canonical stable HTTP contract description baseline;
- JSON fields use lowercase `snake_case`;
- errors use `application/problem+json` plus stable Omnexa machine codes/request IDs;
- scalable lists default to opaque cursor pagination;
- filters/sorts/includes are allowlisted, bounded and authorization-aware;
- protected retriable mutations define `Idempotency-Key` behavior;
- lost-update-sensitive mutations use explicit optimistic concurrency such as `ETag` / `If-Match` where applicable;
- business lifecycle actions are explicit capability/action operations, not arbitrary status patches;
- client-provided tenant/organization IDs are context only and never authorization authority;
- generated routes/SDKs/docs derive from the governed contract rather than replacing it.

## 7. Event invariants — P00.05

Published events follow `EVENT_STANDARD.md` and ADR-0004.

- event type: `<domain>.<subject>.<past_tense_fact>.v<major>`;
- the authoritative producer owns the event meaning/schema/publication condition;
- canonical envelope is CloudEvents-compatible structured JSON with UUIDv7 identity;
- tenant-owned events carry trusted producer-derived `tenantid`;
- correlation, causation and tracing remain separate explicit context;
- delivery baseline is **at least once**;
- consumers are idempotent under duplicate delivery;
- business-significant publication uses transactional outbox or equivalent consistency guarantees;
- business-significant consumers use inbox/deduplication or equivalent durable idempotency controls;
- no global ordering guarantee exists; subject ordering is explicit where required;
- retries are bounded; poison/permanent failures use governed dead-letter/quarantine handling;
- replay preserves original identity/payload and may not duplicate protected business side effects;
- broker routing is infrastructure detail, not canonical business identity;
- event sourcing is not assumed platform-wide.

## 8. Security invariants — P00.06

Security follows `docs/security/SECURITY_STANDARD.md`, `DATA_CLASSIFICATION.md`, `SECURITY_CONTROL_MATRIX.md` and ADR-0005.

### Data classes

```text
PUBLIC
INTERNAL
CONFIDENTIAL
RESTRICTED
```

Rules:

- tenant-owned business records default to at least `CONFIDENTIAL` unless the owning domain explicitly publishes a `PUBLIC` projection;
- encryption does not lower classification;
- exact copies/derived representations inherit sensitivity unless approved irreversible transformation demonstrably lowers risk;
- secrets, authentication equivalents, private keys and high-risk credential material are `RESTRICTED`;
- `RESTRICTED` values do not enter ordinary logs, traces, analytics, search indexes, support screenshots or generic exports;
- `RESTRICTED` data is prohibited as AI/model input by default;
- production `CONFIDENTIAL`/`RESTRICTED` data does not flow to dev/test by default;
- classification applies to queues, DLQs, files, analytics, embeddings/vector stores and AI prompts/outputs as data stores.

### Trust and identity

- same process/network/module/VPC is not implicit trust;
- human users, service accounts, workloads, devices, integrations, support operators and AI agents are distinct principal types;
- authentication proves identity but never substitutes for authorization;
- authorization combines RBAC + relationships + contextual policy + bounded capabilities;
- roles such as `owner`, `admin` or `superuser` are not hidden unrestricted bypasses.

### Tenant isolation

- tenant scope is derived from trusted identity/policy/execution context;
- client payload/header tenant IDs are never authority;
- OLTP, cache, files, search, analytics, broker/events, backups and AI/vector data preserve tenant isolation;
- cross-tenant operations require explicit privileged capability and audit;
- affected implementation paths require negative cross-tenant tests.

### Secrets and cryptography

- secrets are never committed or emitted to logs/traces/errors/CI output/AI context;
- production secrets use approved secret-management/KMS mechanisms;
- modules do not invent custom cryptography;
- protected traffic uses authenticated encryption in transit;
- production stores/backups use approved encryption at rest;
- passwords are never reversibly stored and use an approved adaptive password-hashing mechanism when password auth exists.

### High-risk operations

Role/permission changes, tenant transfer, auth/MFA policy changes, secret/key actions, bulk export, destructive purge, payment/refund/payout actions, financial close/reversal, support impersonation, side-effecting event replay, module trust changes and high-impact AI actions require explicit capability/policy/audit and may require stronger authentication, reason capture, approval or dual control according to risk.

### Integrations, modules and AI

- webhooks/integrations are independent trust and disclosure boundaries;
- inbound webhooks require authenticity verification and replay/deduplication controls where supported;
- configurable outbound destinations must address SSRF/egress risks;
- modules/extensions declare permissions and external access requirements;
- future marketplace packages require verifiable identity/integrity/provenance controls;
- AI retrieval is tenant/object authorized before context assembly;
- prompt/tool injection cannot create authority;
- AI actions use the same protected business capabilities as humans/workflows/integrations.

## 9. Technology baseline

Until changed by ADR:

- Go — kernel/backend and primary domain services;
- TypeScript + React — admin/web/builder/primary extension SDK surfaces;
- Rust — edge/native/security-sensitive components only where justified;
- Python — AI/data workloads only where ecosystem value justifies it;
- PostgreSQL — primary OLTP;
- Redis-compatible — cache/ephemeral coordination;
- S3-compatible object storage — files/media;
- NATS/JetStream-class fabric — event/messaging baseline;
- OpenTelemetry — observability semantics.

Technology choice never authorizes premature phase work.

## 10. Required work protocol

For every material change:

1. identify phase/work package and acceptance criteria;
2. inspect repository state before editing;
3. identify canonical terminology and authoritative owner;
4. apply primitive, API, event and security/classification standards at affected boundaries;
5. state ownership/dependency direction and security/tenant implications;
6. implement the smallest complete authorized change;
7. add/update positive and negative tests appropriate to risk;
8. run required quality gates;
9. record build/test/migration/contract/security/CI evidence;
10. update STATUS/STATE only when evidence supports the transition;
11. add/update ADR and dependent docs before implementing an architectural change.

## 11. Forbidden AI/contributor behavior

Do not:

- start future phases early;
- add an unplanned domain silently;
- create duplicate authoritative ownership;
- invent alternative primitive/API/event/security semantics;
- write directly across module-private schemas;
- bypass tenancy, authorization, audit, classification or contract versioning;
- grant AI a private write path or broader authority;
- log/store secrets for debugging convenience;
- copy sensitive production data to lower environments without explicit security approval;
- weaken a shared security control inside a module;
- introduce hidden super-admin/support bypasses;
- claim `done` without acceptance evidence;
- mix unrelated project code into this repository.

## 12. Change-control triggers

An ADR plus plan/document reconciliation is required before materially changing:

- product/platform boundary;
- domain ownership or canonical terminology semantics;
- identifier/money/time/locale/error primitives;
- stable API semantics;
- event envelope/naming/versioning/delivery/replay semantics;
- trust/security boundary, tenant isolation, authentication/session or authorization model;
- confidentiality classification model, secrets/key handling, audit integrity or AI execution authority;
- language/runtime baseline;
- module lifecycle/dependency model;
- deployment topology baseline;
- phase ordering or mandatory gates.

Never implement the architecture change first and document it later.

## 13. Pull request requirements

Every material PR states:

- phase/work package;
- scope/non-scope;
- architecture impact;
- authoritative owner(s) and dependency mechanism;
- primitive/API/event impacts;
- security/tenancy/data-classification impact;
- migration/data impact;
- tests and CI evidence;
- compatibility/rollback considerations;
- status/state reconciliation.

Unrelated work belongs in separate PRs.

## 14. Safe completion

A change is incomplete when any relevant required gate fails, including:

- build/tests/static/contract checks;
- fresh-install or migration repeatability;
- tenant/object authorization negative tests;
- public contract/version compatibility;
- duplicate-delivery/replay safety;
- security/classification requirements;
- secret disclosure controls;
- ownership/dependency rules;
- documentation/state consistency.

P00.07 defines the canonical testing/CI/release gate taxonomy that future phases must execute.

## 15. Repository and legal guardrails

- use governed feature branches/PRs rather than intentional direct `main` changes;
- follow `CONTRIBUTING.md` and `SECURITY.md`;
- hosted protection target is defined in `docs/governance/REPOSITORY_HARDENING.md` and tracked by issue #3;
- do not replace the existing license or treat it as the approved final commercial strategy without explicit owner/legal decision tracked by issue #4;
- AI systems do not make trademark/legal ownership decisions autonomously.

## 16. Scope drift

Useful work outside the active package is recorded for later or proposed through plan change. It is not silently absorbed into the current package.

The objective is not maximum feature count. It is a platform that can grow to very high feature count **without architectural decay**.
