# Omnexa Repository Execution Contract

This is the highest-priority repository instruction for human contributors and AI coding systems. It applies to the entire repository.

## Current canonical state

Omnexa is a **Composable Enterprise Business Operating System**.

`docs/roadmap/STATE.json` is the machine-readable execution source of truth.

Current program state:

```text
Foundation Architecture v1: FROZEN
P00: ACTIVE — exit verification
P00.10: ACTIVE — review_state=verification
P01: BLOCKED
kernel_code_authorized: false
business_feature_code_authorized: false
```

Do not infer implementation permission from the word `FROZEN`. Architecture is frozen; implementation is still locked.

## Mandatory read order

Before material work, read:

1. `AGENTS.md`
2. `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
3. `docs/governance/FOUNDATION_FREEZE.json`
4. `docs/governance/P01_ENTRY_GATE.md`
5. Product Constitution + system/module architecture
6. glossary, naming, domain ownership and dependency matrix
7. identifier/money/time/locale/error standards
8. API and Event standards
9. Security Standard, Data Classification, Security Control Matrix and Threat Model
10. Testing, CI, Release and Quality Gate standards
11. repository/local-development/toolchain/configuration/developer-command standards
12. SLO, Incident and Reliability standards
13. roadmap `MASTER_PLAN.md`, `STATUS.md`, `STATE.json`
14. AI Execution Policy, Change Control and Definition of Done
15. active CI exception if present
16. relevant accepted ADRs, especially ADR-0010.

If canonical documents conflict, stop implementation and resolve through change control.

## P00.10 freeze rule

P00.01–P00.09 are accepted as **Omnexa Foundation Architecture v1**. Material reinterpretation requires a superseding accepted ADR and reconciliation.

P00.10 remains active in exit verification. The only authorized work is:

- P00 exit verification;
- branch-protection/CI entry-gate remediation;
- narrow governance reconciliation;
- explicit correction of a discovered frozen-baseline contradiction.

Do not start P01 kernel code while the P01 entry gate is blocked.

## P01 entry blockers

### Issue #3 — protected integration path

Before any executable P01 merge, `main` protection must be applied and verified: PR-based integration, blocked force-push/deletion, controlled bypass, conversation resolution and required checks once an executable CI lane exists.

### Issue #14 — executable verification lane

Before any executable P01 merge, an approved CI/self-hosted/provider lane must be able to run repository-owned verification semantics in a clean reproducible environment. GitHub Actions itself is not architecturally mandatory; a compliant equivalent is acceptable.

### Issue #4 — external distribution only

Licensing/IP/trademark resolution is a hard gate before public/external distribution or self-hosted customer delivery. It does not block private internal P01 engineering after Issue #3/#14 are cleared.

## Temporary hosted-CI exception

ADR-0006 applies only while P00 remains open. Hosted Actions is `BLOCKED`/`NOT RUN`, never `PASS`. The exception cannot authorize executable P01 work and expires at P00 exit.

## Architecture invariants

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

## Frozen primitive/API/event/security rules

- UUIDv7 canonical IDs; PostgreSQL native `uuid`;
- exact-decimal money + explicit currency; no binary floating point;
- UTC/`timestamptz` instants + IANA civil-time semantics;
- BCP 47 locale and first-class RTL;
- stable safe structured errors;
- stable HTTP routes `/api/v{major}/{domain}/{resources}` with OpenAPI 3.2.0;
- cursor pagination, explicit idempotency/concurrency and explicit lifecycle actions;
- events are producer-owned versioned facts with CloudEvents-compatible envelope, at-least-once assumptions, idempotent consumers, outbox/inbox, bounded retry/DLQ and governed replay;
- data classes: `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`;
- authn never substitutes for authz;
- authz = RBAC + relationships + contextual policy + bounded capabilities;
- tenant isolation spans all persisted/derived/AI stores;
- secrets/private keys/auth equivalents are `RESTRICTED`;
- privileged support/export/purge/replay/financial/high-impact AI actions require explicit policy/audit;
- prompt/tool injection cannot create authority.

## Quality and release rules

Gate classes: `G0` Governance, `G1` Static, `G2` Unit/Component, `G3` Contract/Integration, `G4` Data/Migration, `G5` Security/Tenancy, `G6` Lifecycle/Resilience, `G7` Build/Package, `G8` Supply Chain/Release.

Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`. Never treat blocked/unrun/N/A as PASS.

Affected risk requires positive and negative evidence: cross-tenant denial, authz deny, duplicate/idempotency/replay, fresh+upgrade migrations, optional-module lifecycle/degradation and security boundaries. Flaky tests are defects. Releases prefer immutable build-once/promote artifacts with source SHA and gate evidence.

## Repository/local-development rules

Canonical roots: `apps/`, `kernel/`, `modules/`, `platform/`, `shared/`, `infrastructure/`, `scripts/`, `docs/`, `generated/`.

- folder != microservice;
- module private code/schema/migrations stay with owner;
- generated output is derivative, not source of truth;
- default local infra = containerized PostgreSQL + Redis-compatible + NATS/JetStream + S3-compatible storage;
- Kubernetes is not required for ordinary local development;
- local and CI use the same semantic verification rules;
- toolchains/dependencies are repository-pinned;
- secrets are separate from committed config;
- production sensitive data is prohibited by default locally;
- use synthetic deterministic multi-tenant fixtures;
- Linux is canonical backend environment; Windows backend development prefers WSL2; native Windows is a separate certification target where required;
- supported workflows must not depend on hidden manual SQL/file/UI steps.

## Threat/reliability rules

Every future material trust boundary/provider/privileged capability requires a threat-model delta.

Operational criticality: `TIER_0`, `TIER_1`, `TIER_2`, `TIER_3` with initial mature-production availability objectives 99.99%, 99.95%, 99.9%, 99.5%.

Recovery targets: A <=5m RPO/<=30m RTO; B <=15m/<=2h; C <=24h/<=8h; D rebuild-based. These are targets until recovery rehearsal proves them.

Zero-tolerance conditions include cross-tenant disclosure, unauthorized privileged mutation, duplicate protected financial side effects, material financial/ledger integrity violation and lost acknowledged durable work. Error budgets never excuse these conditions.

Incident model is `SEV0`–`SEV3`. Fail closed: never bypass authorization because a dependency is down; never treat unverified payment as success; never drop protected audit; never silently lose durable work.

## Technology baseline

Go backend/core; TypeScript+React web/admin/builder/SDK; Rust justified edge/native/security; Python justified AI/data; PostgreSQL OLTP; Redis-compatible cache; S3-compatible storage; NATS/JetStream-class messaging; OpenTelemetry observability.

## Required work protocol

For every material change:

1. verify the active phase/package and implementation locks;
2. inspect `FOUNDATION_FREEZE.json`, `P01_ENTRY_GATE.md` and `STATE.json`;
3. identify canonical terminology and authoritative owner;
4. apply frozen primitive/API/event/security/quality/development/operations rules;
5. map risk to G0-G8 and operational criticality/recovery class where relevant;
6. preserve module/repository boundaries;
7. implement only authorized scope;
8. add positive + negative evidence and threat-model delta where needed;
9. execute canonical verification when executable work is authorized;
10. record evidence accurately and reconcile STATUS/STATE;
11. document architecture change through ADR/change control before implementation.

## Forbidden behavior

Do not start P01 while blocked; silently add domains; duplicate ownership; invent conflicting contracts/security/quality/toolchain/SLO semantics; cross-write/import module-private internals; bypass tenancy/authz/audit/classification; grant AI private write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; call blocked CI PASS; claim untested RPO/RTO as achieved; spend availability error budget on security/integrity violations; or mix unrelated project code.

## Exact next transition

Only after Issue #3 and #14 have verified evidence may a narrow governance PR:

- mark P00.10 and P00 done;
- expire ADR-0006;
- activate P01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- record entry evidence;
- define the first P01 kernel work package.

Do not combine this transition with unrelated kernel feature implementation.

## Scope drift

Useful work outside the active gate is recorded for later or proposed through change control. The objective is a platform that can grow without architectural decay.