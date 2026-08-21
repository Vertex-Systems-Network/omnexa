# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active**
- Current work package: **P00.10 — Foundation architecture freeze review**
- Business-feature implementation: **NOT AUTHORIZED YET**
- Kernel implementation: **NOT AUTHORIZED YET**
- P00 progress: **9 / 10 done**

## P00 work packages

| ID | Work package | State |
|---|---|---|
| P00.01 | Repository governance baseline | done |
| P00.02 | Product/domain glossary and naming | done |
| P00.03 | ID/money/time/locale/error conventions | done |
| P00.04 | API contract standard | done |
| P00.05 | Event contract standard | done |
| P00.06 | Security/data classification | done |
| P00.07 | Testing/CI/release standard | done |
| P00.08 | Local developer/repository structure | done |
| P00.09 | Threat model and operational SLO targets | done |
| P00.10 | Foundation architecture freeze review | **active** |

## P00.09 frozen threat/reliability baseline

The foundation threat model covers cross-tenant escape, authorization/authentication abuse, privilege escalation, injection, SSRF, webhook spoofing/replay, event/job replay, financial duplication/integrity loss, module/supply-chain compromise, CI/release credential theft, POS/edge compromise, backup/export/search/vector leakage, AI prompt/tool abuse, insider/support misuse, noisy-neighbor/resource exhaustion, DDoS/provider outage, migration corruption, region failure, audit tampering, secrets exposure and misconfiguration.

Operational criticality tiers:

```text
TIER_0  integrity-critical control paths
TIER_1  core transactions
TIER_2  interactive supporting capabilities
TIER_3  optional/background capabilities
```

Initial mature-production availability objectives are 99.99%, 99.95%, 99.9% and 99.5% respectively, subject to capability-specific refinement. Security and integrity are not traded for uptime.

Recovery classes:

```text
A  critical identity/security/financial state: target RPO <= 5m, RTO <= 30m
B  core transactional state:              target RPO <= 15m, RTO <= 2h
C  supporting/derived state:              target RPO <= 24h, RTO <= 8h
D  reproducible caches/indexes:            rebuild-based objective
```

These are architecture targets until restore/recovery rehearsal proves them.

Zero-tolerance SLIs include cross-tenant disclosure, unauthorized privileged mutations, duplicate protected financial side effects caused by replay/retry, lost acknowledged durable business work and material financial/ledger integrity violations.

Incident model is **SEV0–SEV3**, with security/privacy/integrity allowed to outrank pure availability impact. Reliability rules cover observability, golden signals plus domain-integrity signals, bounded timeouts/retries, circuit breakers, backpressure, graceful degradation, capacity/noisy-neighbor protection, dependency resilience, backup/restore rehearsal and production-readiness ownership.

Normative P00.09 evidence:

- `docs/security/THREAT_MODEL.md`
- `docs/operations/SLO_STANDARD.md`
- `docs/operations/INCIDENT_STANDARD.md`
- `docs/operations/RELIABILITY_STANDARD.md`
- `docs/contracts/operations/operational-targets.schema.json`
- `docs/adr/ADR-0009-threat-model-slo-reliability-baseline.md`
- `scripts/validate_operations_spec.py`

## Existing frozen baselines

P00.03–P00.08 remain binding: foundation primitives, HTTP/event contracts, security/data classification, provider-independent quality/release gates, and governed monorepo/local-development structure.

## Temporary GitHub Actions exception

GitHub Actions allowance is exhausted/disabled. ADR-0006 permits only P00 documentation/specification manual evidence while hosted execution is unavailable. Hosted CI is **BLOCKED / NOT RUN**, never PASS.

**The exception expires at P00 exit and cannot authorize P01 implementation.**

## Outstanding governance/business blockers

1. Issue #3 — `main` branch/ruleset protection still requires hosted admin configuration.
2. Issue #14 — hosted Actions quota/runner remains unavailable.
3. Issue #4 — licensing/IP/trademark strategy remains unresolved before external distribution/public launch.

P00.10 must explicitly classify these as P01 entry blockers, external-launch blockers, or acceptable post-P00 debt; it may not silently ignore them.

## Execution lock

P00.10 is the only authorized work package. Do not begin kernel, database model, CRM, ERP, commerce, POS, website builder, payments or AI product implementation until the final freeze review is complete **and P01 entry prerequisites are satisfied**.