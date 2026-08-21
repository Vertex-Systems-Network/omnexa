# ADR-0009 — Threat Model, SLO and Reliability Baseline

Status: **accepted**

## Context

Omnexa is a multi-tenant business operating system expected to carry identity, financial, operational, commerce, integration, edge and AI workloads. Before kernel implementation begins, the platform needs explicit threat and operational reliability semantics so implementation teams do not invent incompatible security/recovery expectations later.

## Problem

Without a shared baseline:

- tenant/security threats may be considered too late;
- availability may be optimized at the expense of integrity;
- RPO/RTO claims may be made without restore evidence;
- incident severity may vary by team;
- error budgets may hide zero-tolerance security/financial failures;
- operational readiness may depend on undocumented local knowledge.

## Decision

Adopt:

1. `docs/security/THREAT_MODEL.md` as the minimum cross-platform threat model;
2. `docs/operations/SLO_STANDARD.md` for criticality tiers, initial SLO targets, RPO/RTO classes and error-budget policy;
3. `docs/operations/INCIDENT_STANDARD.md` for SEV0-SEV3 incident semantics and response lifecycle;
4. `docs/operations/RELIABILITY_STANDARD.md` for observability, retries/backpressure, degradation, capacity, dependency and recovery readiness;
5. `docs/contracts/operations/operational-targets.schema.json` for future machine-readable capability declarations.

## Key invariants

- tenant escape, unauthorized privilege, duplicate protected financial side effects and lost acknowledged durable business work are zero-tolerance conditions;
- security/integrity is not traded for availability;
- Tier 0/Tier 1/Tier 2/Tier 3 classify operational criticality;
- Recovery classes A-D define initial RPO/RTO objectives;
- RPO/RTO are targets until restore/recovery rehearsal proves them;
- error-budget consumption may restrict risky releases but cannot excuse security/integrity breaches;
- every material production capability eventually has an owner, SLI/SLO, recovery treatment, observability, runbook and threat-model delta;
- provider failure remains visible in user-facing reliability semantics;
- optional degradation may preserve core functionality only when correctness and security remain intact.

## Alternatives considered

### Defer threat modeling until modules exist
Rejected. Platform-level trust boundaries and invariants must shape kernel design.

### Define only uptime targets
Rejected. Availability without integrity, recoverability and tenant isolation is insufficient.

### Treat all services as the same criticality
Rejected. It either over-engineers optional workloads or under-protects identity/financial/control paths.

### Make SLOs provider-specific
Rejected. Operational semantics are repository/product contracts; telemetry providers may change.

## Consequences

Positive:

- P01+ teams inherit explicit security/reliability objectives;
- operational claims require measurable evidence;
- incident language is consistent;
- service extraction/scaling decisions can reference real SLO pressure later.

Costs:

- every production capability must maintain ownership and operational metadata;
- stricter targets require telemetry, recovery rehearsal and on-call maturity;
- zero-tolerance integrity events may intentionally stop availability-preserving shortcuts.

## Compatibility impact

No runtime compatibility impact. This is a pre-implementation architecture baseline.

## Migration impact

None during P00. Future services/modules must adopt the baseline as they are implemented.

## Security/tenancy impact

Strengthens mandatory threat modeling, tenant isolation, privileged-action, data recovery and incident controls.

## Operational impact

Establishes target availability classes, RPO/RTO classes, error-budget rules, incident severity and operational-readiness requirements.

## Rollback / forward-fix

Material changes require a superseding ADR and reconciliation of security/operations/quality/roadmap documents. Existing production claims must not be weakened silently.
