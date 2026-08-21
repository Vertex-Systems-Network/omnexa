# ADR-0007 — Testing, CI & Release Baseline

Status: **Accepted**  
Date: **2026-08-22**  
Work package: **P00.07**

## Context

Omnexa will span backend services, web/admin, modules, migrations, events, workflows, POS/edge, integrations, marketplace packages and AI capabilities. A generic "tests must pass" rule is insufficient because different risks require different evidence, and CI-provider outages must not redefine correctness.

## Decision

Omnexa adopts the following baseline:

1. Quality semantics are repository-owned and CI-provider independent.
2. Local and CI environments execute the same canonical quality commands/gates.
3. Testing is risk-based and includes negative evidence for affected tenancy, authorization, security, replay/idempotency and lifecycle invariants.
4. Canonical test layers include static, unit, component, contract, integration, migration, security/negative, module lifecycle, E2E, performance/resilience, compatibility and recovery/rehearsal testing.
5. Required tests are deterministic; flaky tests are defects, not accepted background noise.
6. Database changes require fresh-install and supported-upgrade migration evidence.
7. Public HTTP/event/module contracts require machine validation and compatibility evidence.
8. Every material change maps affected risk to gate classes G0-G8.
9. Evidence states are exactly `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`; blocked/unrun gates are never green equivalents.
10. CI uses least privilege, reproducible toolchains/dependencies, controlled secrets and supply-chain checks.
11. Release flow prefers build-once/promote-immutable-artifacts, semantic versioning, immutable tags/digests and explicit compatibility/migration notes.
12. Production-capable releases require applicable testing, migration, security, packaging and supply-chain gates.
13. SBOM/provenance/signing become required as the package/deployment model reaches implementation maturity.
14. Rollback/forward-fix semantics must account for irreversible database and external business side effects.
15. Hosted CI outages/quotas are infrastructure states and require explicit exceptions; ADR-0006 remains P00-only and cannot authorize executable release bypass.

Normative details live in:

- `docs/quality/TESTING_STANDARD.md`
- `docs/quality/CI_STANDARD.md`
- `docs/quality/RELEASE_STANDARD.md`
- `docs/quality/QUALITY_GATE_MATRIX.md`
- `docs/contracts/quality/quality-gates.schema.json`

## Consequences

### Positive

- quality requirements remain stable across CI providers;
- modules cannot hide behind aggregate coverage or happy-path tests;
- tenant/security/contract failures receive first-class negative evidence;
- release artifacts become traceable to source and verification;
- Actions quotas/outages do not become architecture decisions.

### Costs

- implementation phases must create local canonical verification commands;
- test matrices and release certification add engineering/compute cost;
- flaky tests must be fixed/quarantined with governance instead of ignored;
- migration and module lifecycle testing requires real infrastructure fixtures.

## Rejected alternatives

### CI-provider-specific quality semantics

Rejected because changing provider/quota must not change what `correct` means.

### Coverage-percentage-only quality gate

Rejected because high coverage can coexist with missing tenant/security/business invariants.

### E2E-heavy strategy

Rejected because large E2E suites are slow/brittle and cannot replace deterministic unit/contract/integration evidence.

### Rebuild for each environment

Rejected where immutable artifact promotion is possible because different bytes weaken release traceability.

### Treat blocked CI as green

Rejected. `BLOCKED` is a distinct state and requires explicit governance if progression is allowed.

## Compatibility

This ADR consumes the P00.03 primitive, P00.04 API, P00.05 event and P00.06 security contracts. P00.08 will define repository/local-development structure capable of implementing the canonical commands. P00.09 will add concrete threat/SLO-driven test budgets. P01+ executable work must follow this baseline.

## Supersession

Material changes to test-layer taxonomy, evidence states, gate classes, CI-provider independence, release artifact identity/promotion or production release gate semantics require a superseding ADR.