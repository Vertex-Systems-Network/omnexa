# Omnexa Foundation Architecture Freeze Review

Status: **Architecture baseline accepted; P00 exit in verification**  
Work package: **P00.10**  
Review date: **2026-08-22**

## 1. Purpose

This is the formal P00 architecture-freeze review. It verifies that Omnexa has enough governed foundation architecture to begin kernel implementation **once P01 implementation-entry blockers are cleared**.

This review does not claim executable product readiness. P00 intentionally contains architecture, governance and specifications only.

## 2. Review result

### Architecture decision

**ACCEPTED FOR FREEZE.**

The P00.01–P00.09 foundation is coherent enough to freeze as Omnexa Foundation Architecture v1. Future material changes require change control and an accepted superseding ADR.

### Program transition decision

**P00 EXIT NOT YET AUTHORIZED.**

P00.10 remains in `verification` because the repository is not yet safe to admit executable P01 kernel merges:

- Issue #3: `main` branch/ruleset protection is not applied/verified.
- Issue #14: an executable CI lane is unavailable because GitHub Actions allowance/runner execution is blocked.

These are not architecture-design defects; they are implementation-entry controls. They must be satisfied before P00 can exit and before `kernel_code_authorized` becomes true.

Issue #4 (licensing/IP/trademark) is classified separately as a hard gate before external distribution/public launch, not as a blocker for private internal P01 engineering.

## 3. Frozen foundation inventory

### P00.01 — Governance

Frozen:
- Product Constitution;
- AI execution policy;
- change control and Definition of Done;
- roadmap/state/ledger discipline;
- module and system architecture baseline.

### P00.02 — Terminology, ownership and dependencies

Frozen:
- canonical glossary/naming;
- domain ownership registry;
- dependency matrix;
- CODEOWNERS/contributor/security governance;
- repository hardening target.

### P00.03 — Foundation primitives

Frozen:
- UUIDv7 canonical identifiers;
- exact decimal money and explicit currency;
- UTC/`timestamptz`, business-date and IANA civil-time semantics;
- BCP 47 locale and RTL;
- stable structured error contracts.

### P00.04 — Stable HTTP API

Frozen:
- `/api/v{major}/{domain}/{resources}`;
- OpenAPI 3.2.0;
- lowercase `snake_case` JSON;
- Problem Details error transport;
- cursor pagination;
- idempotency and optimistic concurrency rules;
- client tenant IDs never become authorization authority.

### P00.05 — Events

Frozen:
- producer-owned versioned past-tense event facts;
- CloudEvents-compatible envelope;
- at-least-once delivery assumptions;
- idempotent consumers;
- outbox/inbox durability;
- bounded retry/DLQ;
- governed replay semantics.

### P00.06 — Security and data classification

Frozen:
- `PUBLIC`, `INTERNAL`, `CONFIDENTIAL`, `RESTRICTED`;
- authentication separate from authorization;
- RBAC + relationships + contextual policy + capability authorization;
- cross-store tenant isolation;
- secrets/KMS/crypto rules;
- audit/privileged/support/AI controls;
- integration/module/edge trust boundaries.

### P00.07 — Testing, CI and release

Frozen:
- G0–G8 quality-gate taxonomy;
- evidence states `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`;
- provider-independent local/CI semantics;
- negative tenant/security/idempotency/migration/lifecycle evidence;
- semantic versioning and build-once/promote release model;
- SBOM/provenance/signing expectations where applicable.

### P00.08 — Repository/local development

Frozen:
- governed monorepo ownership roots;
- strict modular monolith as initial deployment architecture;
- module-private code/schema/migration boundaries;
- repository-pinned toolchains;
- containerized local PostgreSQL/Redis-compatible/NATS/S3-compatible dependencies;
- config/secrets separation;
- Linux canonical backend environment and WSL2-preferred Windows backend workflow;
- Kubernetes not required for default local development.

### P00.09 — Threat/reliability

Frozen:
- platform threat model T01–T24;
- operational criticality TIER_0–TIER_3;
- recovery classes A–D;
- initial SLO/RPO/RTO targets;
- zero-tolerance security/financial/tenant integrity SLIs;
- error-budget semantics;
- SEV0–SEV3 incident model;
- observability/retry/backpressure/degradation/capacity/recovery-readiness rules.

## 4. Cross-document consistency review

Reviewed source-of-truth relationships:

- Product Constitution -> System Architecture -> Module Standard;
- Glossary/Naming -> Domain Ownership -> Dependency Matrix;
- primitive semantics -> API/Event contracts;
- Security/Data Classification -> Threat Model;
- Quality Gates -> Local Developer Commands/CI semantics;
- Release Standard -> supply-chain/provenance expectations;
- SLO/Reliability -> Quality/Incident/Threat requirements;
- STATE/STATUS/README/AGENTS execution lock.

No intentional contradictory architecture baseline is accepted by this review. If a later contradiction is discovered, P00.10 or the affected package must be reopened rather than guessed around.

## 5. Explicit non-decisions

P00 does **not** freeze prematurely:

- exact cloud vendor;
- Kubernetes as mandatory runtime;
- day-one microservices;
- exact analytics warehouse/search/vector vendor;
- exact identity provider implementation;
- exact payment gateways;
- exact AI model/provider;
- exact country tax/legal packs;
- final commercial licensing/trademark strategy.

These remain later governed decisions within frozen architectural constraints.

## 6. P01 implementation-entry blockers

### BLOCKER-01 — Protected integration path
Tracked by Issue #3.

Before any P01 executable merge to `main`:

- `main` requires PR-based integration;
- force push and deletion are blocked;
- conversation-resolution policy is enforced;
- direct/bypass access is restricted to explicit break-glass actors;
- required quality check policy is configured once a viable executable CI lane exists;
- configured protections are verified rather than assumed.

### BLOCKER-02 — Executable verification lane
Tracked by Issue #14.

Before any P01 executable merge:

- at least one approved CI execution lane can actually run repository verification commands; GitHub Actions is not mandatory if another approved provider/self-hosted lane satisfies P00.07;
- `omnexa verify`/equivalent canonical commands can run in a clean reproducible environment;
- G0/G1/G2/G7 minimum kernel bootstrap gates can execute;
- later affected gates G3–G8 are added as implementation capability appears;
- CI credentials/secrets remain least-privileged;
- a blocked runner is never represented as a green check.

## 7. External distribution blocker

### BLOCKER-EXT-01 — License/IP/trademark
Tracked by Issue #4.

Before public launch, external source/binary distribution, self-hosted customer delivery or marketplace commercialization:

- final license model is owner/legal approved;
- third-party license policy is defined;
- contributor/IP policy is defined;
- `Omnexa` name/trademark/domain clearance is completed for intended markets;
- resulting LICENSE/trademark decisions are committed through governed change control.

This does not block private internal architecture or kernel development.

## 8. P01 entry contract

When BLOCKER-01 and BLOCKER-02 are satisfied, a narrow governance change may:

1. mark P00.10 `done`;
2. mark P00 phase `done`;
3. expire ADR-0006 temporary P00 CI exception;
4. activate P01;
5. set `kernel_code_authorized = true`;
6. keep `business_feature_code_authorized = false` until its later phase gate;
7. create the first P01 work package from the canonical kernel plan;
8. record verification evidence for branch protection and executable CI.

Do not combine this entry transition with unrelated kernel feature implementation.

## 9. P01 implementation starting boundary

P01 is kernel/platform skeleton only. It does not authorize CRM, ERP, commerce, POS, payments, website builder, HR, manufacturing or AI business features.

Initial P01 capabilities should establish boot/configuration, database abstraction/migrations, event/job foundations, observability, health, module bootloader skeleton, feature flags, secrets interface, caching/files interfaces, scheduler/jobs, testing/CLI scaffolding and repository-local canonical verification commands according to the frozen contracts.

## 10. Freeze rule

From this review onward, P00.01–P00.09 are **frozen foundation v1**. A contributor or AI may not silently reinterpret them. Material changes require:

- explicit problem statement;
- affected baseline identification;
- accepted ADR;
- compatibility/security/operational impact;
- roadmap/state reconciliation;
- relevant tests/validators once executable infrastructure exists.

## 11. Final review state

```text
Architecture baseline: ACCEPTED / FROZEN
P00.10 package: VERIFICATION
P00 phase: ACTIVE pending entry controls
P01: NOT AUTHORIZED
Kernel code: LOCKED
Business feature code: LOCKED
Issue #3: P01 ENTRY BLOCKER
Issue #14: P01 ENTRY BLOCKER
Issue #4: EXTERNAL DISTRIBUTION BLOCKER
```
