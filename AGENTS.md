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
Repository visibility: PUBLIC
Executable CI gate: SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04
Hosted proof: run 32537207455 / job 96940269306
Local/self-hosted governance runners: PROHIBITED
P01: BLOCKED BY EG-02 / ISSUE #3
EG-02: ACTIONABLE_UNPROTECTED — live main protected=false
P01.01-P01.12: PREPARED / PLANNED / NOT ACTIVE
kernel_code_authorized: false
business_feature_code_authorized: false
```

Do not infer implementation permission from `FROZEN`, public visibility, hosted CI availability or prepared P01 specifications. Architecture/readiness can be prepared while executable kernel implementation remains locked.

## Mandatory read order

Before material work, read:

1. `AGENTS.md`
2. `docs/governance/FOUNDATION_FREEZE_REVIEW.md`
3. `docs/governance/FOUNDATION_FREEZE.json`
4. `docs/governance/P01_ENTRY_GATE.md`
5. `docs/governance/P00_P01_TRANSITION_CHECKLIST.md`
6. Product Constitution + system/module architecture
7. glossary, naming, domain ownership and dependency matrix
8. identifier/money/time/locale/error standards
9. API and Event standards
10. Security Standard, Data Classification, Security Control Matrix and Threat Model
11. Testing, CI, Release and Quality Gate standards
12. repository/local-development/toolchain/configuration/developer-command standards
13. SLO, Incident and Reliability standards
14. roadmap `MASTER_PLAN.md`, `STATUS.md`, `STATE.json`
15. `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` and the active/prepared P01 package specification when working on the P01 handoff
16. `docs/governance/LICENSING_DECISION.md` and `LICENSING_DECISION_BRIEF.md` for distribution/IP work
17. AI Execution Policy, Change Control and Definition of Done
18. historical/active CI exception record if relevant
19. relevant accepted ADRs, especially ADR-0010.

If canonical documents conflict, stop implementation and resolve through change control.

## P00.10 freeze rule

P00.01–P00.09 are accepted as **Omnexa Foundation Architecture v1**. Material reinterpretation requires a superseding accepted ADR and reconciliation.

P00.10 remains active in exit verification. Authorized work is limited to:

- P00 exit verification;
- branch-protection entry-gate remediation;
- narrow governance/CI reconciliation;
- P01 implementation-readiness **specification only**, with executable kernel code prohibited;
- licensing/IP/trademark decision preparation without changing `LICENSE` or claiming clearance;
- explicit correction of a discovered frozen-baseline contradiction.

Do not create `go.mod`, `go.work`, `kernel/cmd/omnexa/main.go` or equivalent P01 executable implementation while `kernel_code_authorized=false`.

## P01 entry gates

### EG-02 / Issue #3 — protected integration path — ACTIONABLE_UNPROTECTED

Hosted `main` protection remains required unless deliberately superseded through an owner-approved governance ADR.

The repository is now public, so the earlier private-repository plan entitlement blocker is no longer active. However live GitHub branch metadata still reports `protected=false`. P01 remains blocked until the required policy is applied and verified.

Required live evidence includes PR-only integration, strict required `governance` check, conversation resolution, blocked force push/deletion, administrator enforcement and contributor-model-appropriate review policy.

Never treat public visibility or the disappearance of the old HTTP 403 as proof that EG-02 is satisfied.

### EG-03 / Issue #14 — executable verification lane — SATISFIED

Canonical governance CI is **GitHub-hosted only**:

```yaml
runs-on: ubuntu-24.04
```

The job must fail closed unless `RUNNER_ENVIRONMENT=github-hosted`, `RUNNER_OS=Linux` and `RUNNER_ARCH=X64`.

Current evidence: PR #31 migration run `32537207455`, job `96940269306`, runner `GitHub Actions 1000006777`, Ubuntu 24.04.4 LTS / X64, image `ubuntu-24.04`; governance, development, operations, freeze-review, P01-preparation and P01-package-spec validators all PASS.

The canonical workflow must not contain `self-hosted`, `LOCAL-WIN-*`, runner-name discovery, local evidence aggregation or any local-runner fallback.

### Hosted branch-protection administration

Repository administration automation for EG-02 must also use GitHub-hosted compute. `.github/workflows/main-protection-admin.yml` is manual (`workflow_dispatch`) only and must run on `ubuntu-24.04` with an owner-controlled short-lived `OMNEXA_GITHUB_ADMIN_TOKEN` secret.

The ordinary Actions `GITHUB_TOKEN` is not repository Administration authority and must not be used as a privilege substitute.

## P01 readiness rule

`docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json` defines strict sequential one-active-package execution for P01.01–P01.12. All twelve specifications are prepared, but all remain `planned` while EG-02 is blocked.

P01.01 remains the only next package: **Go workspace/build skeleton**.

While blocked:

- implementation is prohibited;
- `scripts/validate_p01_preparation.py` must pass;
- `scripts/validate_p01_package_specs.py` must pass;
- only specification/acceptance/test-boundary refinement is allowed;
- P01.02+ may not be made active;
- P02/P03/business behavior may not be pulled into P01.

When EG-02 is cleared or deliberately superseded, follow `P00_P01_TRANSITION_CHECKLIST.md` in a governance-only PR before any kernel code is added.

## Issue #4 — public visibility / external distribution

The repository is now public while the current `LICENSE` remains GPLv3 and trademark clearance remains unresolved. Issue #4 therefore requires explicit owner/legal reconciliation for product launch, commercial licensing, contribution policy and trademark claims.

`LICENSING_DECISION_BRIEF.md` is an owner/legal decision worksheet only. It does not replace legal review, change `LICENSE`, grant new redistribution rights or establish trademark clearance.

Issue #4 does not automatically authorize or prohibit P01 engineering once the technical P01 entry gate is satisfied.

## Historical CI exception

ADR-0006 records the temporary P00 manual-evidence exception used while hosted Actions capacity was unavailable. GitHub-hosted executable CI is now operational, so ADR-0006 is historical evidence only and cannot authorize any current bypass.

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
- local development and CI must use the same semantic verification rules, but canonical repository CI runs on GitHub-hosted infrastructure only;
- toolchains/dependencies are repository-pinned;
- secrets are separate from committed config;
- production sensitive data is prohibited by default locally;
- use synthetic deterministic multi-tenant fixtures;
- Linux is canonical backend environment;
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
2. inspect `FOUNDATION_FREEZE.json`, `P01_ENTRY_GATE.md`, `P00_P01_TRANSITION_CHECKLIST.md` and `STATE.json`;
3. identify canonical terminology and authoritative owner;
4. apply frozen primitive/API/event/security/quality/development/operations rules;
5. map risk to G0-G8 and operational criticality/recovery class where relevant;
6. preserve module/repository boundaries;
7. implement only authorized scope;
8. add positive + negative evidence and threat-model delta where needed;
9. execute canonical verification through the required GitHub-hosted `governance` job;
10. execute P01 preparation/package validators while P01 remains blocked;
11. record evidence accurately and reconcile STATUS/STATE;
12. document architecture change through ADR/change control before implementation.

## Forbidden behavior

Do not start P01 while EG-02 is blocked; use local/self-hosted runners for canonical governance CI; silently add domains; duplicate ownership; invent conflicting contracts/security/quality/toolchain/SLO semantics; cross-write/import module-private internals; bypass tenancy/authz/audit/classification; grant AI private write authority; commit secrets; use production sensitive data locally; create hidden super-admin bypasses; weaken gates to get green; call blocked CI PASS; claim untested RPO/RTO as achieved; spend availability error budget on security/integrity violations; change `LICENSE` by inference; claim trademark clearance without evidence; or mix unrelated project code.

## Exact next transition

Only after EG-02 is either verified satisfied or explicitly superseded by an owner-approved governance ADR may a narrow governance PR:

- mark P00.10 and P00 done;
- retire ADR-0006 from active use;
- activate P01 and P01.01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- record applicable integration-protection/compensating-control evidence plus GitHub-hosted CI evidence;
- preserve P01.02–P01.12 as planned.

Do not combine this transition with kernel implementation. The first kernel-code PR follows only after the transition is merged and verified on `main`.

## Scope drift

Useful work outside the active gate is recorded for later or proposed through change control. The objective is a platform that can grow without architectural decay.
