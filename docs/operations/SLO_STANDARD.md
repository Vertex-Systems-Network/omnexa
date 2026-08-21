# Omnexa SLO, RPO and RTO Standard

Status: **Canonical foundation v1**  
Work package: **P00.09**

This document defines initial operational objectives for Omnexa. These are product architecture targets, not claims about an implementation that does not yet exist.

## 1. Principles

1. SLOs are defined per user-visible or business-critical capability, not per server.
2. Security, tenant isolation and financial integrity are never traded for availability.
3. `BLOCKED`, `NOT RUN` and degraded dependencies do not become hidden success.
4. SLOs use measured service-level indicators (SLIs) and explicit windows.
5. Every production SLO has an owner, telemetry source, alert policy and error-budget policy.
6. RPO/RTO are recovery objectives, not guaranteed outcomes until tested.
7. Stage-specific targets may be stricter, but may not silently weaken canonical minimums.

## 2. Capability criticality classes

### Tier 0 — Integrity-critical control capabilities
Examples:
- identity/session validation;
- tenant/authorization enforcement;
- secrets/key access;
- audit/security event intake;
- payment authorization/capture/refund controls;
- financial ledger integrity paths;
- module trust/signing decisions.

These capabilities must fail closed when authorization/integrity cannot be proven.

### Tier 1 — Core transactional capabilities
Examples:
- order creation;
- inventory reservation;
- invoice creation;
- workflow execution;
- core customer/vendor operations;
- POS cloud synchronization.

### Tier 2 — Interactive supporting capabilities
Examples:
- dashboards;
- search;
- reporting reads;
- notifications;
- CMS/admin management.

### Tier 3 — Background/optional capabilities
Examples:
- batch analytics;
- recommendations;
- non-critical exports;
- optional AI enrichment;
- non-critical scheduled jobs.

## 3. Initial availability objectives

These are foundation targets for mature production service classes and must be refined per deployment tier.

| Tier | Monthly availability objective | Notes |
|---|---:|---|
| Tier 0 | 99.99% | integrity-first; fail-closed security behavior may reduce availability without being treated as a security failure |
| Tier 1 | 99.95% | core business transactions |
| Tier 2 | 99.9% | interactive supporting functions |
| Tier 3 | 99.5% | optional/background functions |

Planned maintenance is not automatically excluded. Exclusion rules must be explicit in the service's SLO definition.

## 4. Latency objectives

Latency objectives are capability-specific. Initial interactive baseline:

- Tier 0/Tier 1 synchronous APIs: target p95 <= 400 ms and p99 <= 1 s for ordinary in-region requests excluding external-provider time where separately measured;
- Tier 2 interactive reads: target p95 <= 800 ms and p99 <= 2 s;
- asynchronous acceptance endpoints should acknowledge durable acceptance quickly and expose job/workflow state rather than blocking on long work;
- provider-dependent operations must expose provider latency separately from Omnexa processing latency.

Large reports, imports, exports, media processing and AI generation must not be forced into synchronous API latency targets; they use job completion objectives.

## 5. Correctness and integrity SLIs

Availability alone is insufficient. Critical SLIs include:

- cross-tenant disclosure rate: target **0**;
- unauthorized privileged mutation rate: target **0**;
- duplicate protected financial side effects caused by replay/retry: target **0**;
- unreconciled ledger/payment integrity violations: target **0**;
- lost acknowledged durable business events/jobs: target **0**;
- successful audit attribution for privileged actions: target **100%**;
- restore verification success for protected backups according to recovery class: target **100% of scheduled rehearsals**.

A zero-tolerance security/integrity SLI breach is an incident even if monthly availability remains within budget.

## 6. Recovery classes

### Recovery Class A — Critical identity/security/financial state
- target RPO: **<= 5 minutes**;
- target RTO: **<= 30 minutes** for an approved recovery architecture once implemented;
- stronger objectives may be required for payment/financial deployments;
- restoration must preserve tenant isolation, auditability and integrity.

### Recovery Class B — Core transactional business state
- target RPO: **<= 15 minutes**;
- target RTO: **<= 2 hours**.

### Recovery Class C — Supporting/derived state
- target RPO: **<= 24 hours** or reproducible from authoritative sources;
- target RTO: **<= 8 hours**.

### Recovery Class D — Rebuildable cache/index/derived artifacts
- no backup requirement when deterministic rebuild from authoritative data is proven;
- target RTO defined by operational need and rebuild capacity.

RPO/RTO classification is declared by the authoritative owner. A derived store may not claim a weaker recovery class if it becomes the only copy of business-significant data.

## 7. Error budgets

For an SLO objective `S`, error budget over the measurement window is `1 - S`.

Rules:

- budget burn is measured for applicable availability/reliability SLIs;
- zero-tolerance integrity/security violations bypass normal error-budget tradeoffs and trigger incident handling;
- sustained burn above policy threshold pauses reliability-risk-increasing releases for the affected capability;
- teams may spend budget on controlled change, not on known unbounded defects;
- optional feature degradation may preserve a core capability's SLO only if user-visible semantics remain honest.

Initial burn policy:

- > 50% of monthly error budget consumed in <= 7 days -> reliability review and release-risk restriction;
- > 80% consumed -> freeze non-essential risky changes for the affected capability until mitigation;
- 100% consumed -> no discretionary risk-increasing release until owner approves recovery plan.

## 8. Measurement windows

Default availability window: rolling 28 days, with monthly reporting also retained.

Fast-burn alerts should evaluate shorter windows such as 1 hour and 6 hours to detect catastrophic burn before the rolling window is exhausted.

## 9. Dependency and degraded-mode accounting

External provider failures are measured separately but do not disappear from user-visible reliability reporting. Capability definitions must state:

- which dependency is required;
- whether durable queue/degraded mode exists;
- whether the user request is accepted, rejected or pending;
- when compensation/reconciliation is required.

A dependency outage may not justify unsafe bypass of authorization, validation, financial controls or tenant isolation.

## 10. Recovery evidence

Before a deployment can claim a recovery objective, evidence must include as applicable:

- backup creation and encryption validation;
- restore into an isolated environment;
- tenant-scope verification;
- integrity/reconciliation checks;
- measured recovery duration;
- dependency/bootstrap recovery;
- documented operator steps/runbook;
- findings and follow-up actions.

Paper-only RPO/RTO is a target, not a verified capability.

## 11. Evolution

P01+ capability owners must define concrete SLIs and SLOs using this standard. Industry packs or regulated deployments may require stricter objectives, data residency or retention/recovery controls.