# Omnexa Definition of Done

Status: **Mandatory baseline**

A work package is not complete because code exists, a screen renders, or one happy-path test passes. Completion is evidence-based.

## 1. Universal completion gates

Every applicable change must satisfy:

### 1.1 Scope
- work maps to an approved active work package;
- in-scope/out-of-scope is explicit;
- no unrelated project files or features are mixed in;
- dependencies are declared.

### 1.2 Architecture
- ownership boundary is correct;
- no forbidden cross-module database write/import is introduced;
- kernel capabilities are reused instead of duplicated;
- public contracts/events are versioned where applicable;
- architecture changes have an accepted ADR.

### 1.3 Correctness
- build succeeds;
- required static/type checks succeed;
- unit tests succeed;
- integration/contract tests succeed where applicable;
- error paths are covered for material behavior.

### 1.4 Data and migrations
- fresh install/migration path succeeds from zero where schema is affected;
- upgrade migration succeeds from supported baseline;
- migrations are deterministic and tenant-safe;
- destructive changes have explicit approval and strategy;
- seed/reference data remains reproducible where used.

### 1.5 Tenancy and authorization
- tenant-owned operations are scoped explicitly;
- cross-tenant negative tests exist where relevant;
- authorization is enforced server-side;
- privileged operations are audited.

### 1.6 Async/event behavior
Where applicable:
- retries are bounded;
- duplicate delivery is safe;
- idempotency is tested;
- event schema/version validation passes;
- poison-message/dead-letter behavior is known;
- correlation/tenant context propagates.

### 1.7 Module lifecycle
Where applicable:
- install succeeds;
- enable succeeds;
- disable is non-destructive;
- re-enable succeeds;
- optional dependencies degrade safely;
- upgrade path succeeds;
- purge behavior is explicit and protected.

### 1.8 Security
- secrets are not committed;
- sensitive data classification is respected;
- input/output validation exists at trust boundaries;
- external webhook/signature verification is implemented where relevant;
- dependency/security scans pass when configured;
- new privileged capability has explicit permission scope.

### 1.9 Observability
Material operations expose enough telemetry to diagnose failure:
- structured logs;
- correlation/trace ID;
- relevant metrics/traces;
- health status for long-lived service/module dependencies.

### 1.10 Documentation and state
- code behavior and docs agree;
- PR records acceptance evidence;
- `STATUS.md` is updated if progress changed;
- `STATE.json` is updated only after evidence supports the transition;
- ADRs and affected architecture documents are reconciled.

## 2. Evidence vocabulary

Use only these evidence states:

- **PASS** — observed successful execution with command/run/reference recorded;
- **FAIL** — observed failure with actionable reference;
- **BLOCKED** — cannot execute because a named dependency/environment is unavailable;
- **NOT RUN** — not executed;
- **N/A** — gate demonstrably does not apply.

Do not use vague states such as "should pass", "looks good" or "probably fine" as completion evidence.

## 3. Required PR evidence block

Each implementation PR should contain a compact evidence section similar to:

```text
Build: PASS — <command or CI job>
Static checks: PASS — <command/job>
Unit: PASS — <command/job>
Integration: PASS/N/A — <command/job>
Migration fresh: PASS/N/A — <command/job>
Migration upgrade: PASS/N/A — <command/job>
Tenant isolation: PASS/N/A — <test/job>
Authorization: PASS/N/A — <test/job>
Module lifecycle: PASS/N/A — <test/job>
Contract/event validation: PASS/N/A — <test/job>
Security scan: PASS/N/A — <job>
```

## 4. Phase completion

A phase may be marked `done` only when:

1. every mandatory work package is done;
2. phase exit-gate tests/evidence pass;
3. unresolved blockers are zero, or explicitly moved by approved plan change;
4. documentation and machine-readable state reconcile;
5. the next phase prerequisites are demonstrably satisfied.

## 5. Reopening work

If later evidence proves a supposedly completed package violates an invariant or required gate, reopen it. Accuracy of project state is more important than preserving a progress percentage.

## 6. Zero-error interpretation

"Zero errors" means zero known failing required gates in the verified target environment. It does not mean the absence of all possible defects. Claims must always be tied to the exact tests/build/migration/runtime evidence executed.

## 7. Strategic AI-native applicability gate

When an accepted future strategic program from `docs/roadmap/STRATEGIC_PROGRAMS.json` applies, the work package must explicitly decide the following rather than silently omitting them:

- Research/problem/VOC/CTQ and intended business outcome;
- architecture/ADR and domain write owner;
- DataFlow, classification, lineage, retention/delete/export and derived-store behavior;
- Business Graph contribution or `N/A`;
- Process Graph contribution or `N/A`;
- System Graph/Flow contribution and evidence provenance or `N/A`;
- security/threat model and FMEA for high/critical material failure paths;
- tenancy/authorization/SoD/approval implications;
- AI/model/agent/tool context and evaluation implications;
- performance/capacity/cost budget;
- reliability/SLO/recovery/reconciliation behavior;
- configuration/environment promotion impact;
- migration/onboarding/rollback impact;
- tests/evals/target evidence;
- post-release outcome/control evidence.

`N/A` requires a rationale. This section does not activate proposed X-program implementation ahead of canonical phase/work-package authority.

## 8. AI-authored work and evidence independence

For material AI-assisted implementation:

- AI-authored PASS prose is not machine evidence;
- the reviewed/tested source SHA must be identifiable;
- deleted/skipped/weakened tests or reduced security/performance thresholds are explicit review signals;
- a code-authoring AI must not silently weaken governance/workflow/security controls to make its own change pass;
- high/critical work requires a distinct review context plus objective evidence appropriate to the risk;
- stale review is invalid when the reviewed head changes materially;
- modelled/AI-inferred graph/process/business outcomes are not observed facts;
- target/provider/runtime claims require the target/provider/runtime evidence specified by the owning gate.

## 9. Graph and derived-intelligence completion rule

When Business Graph, Process Graph, System Graph, analytics, search, simulation or AI projections are involved:

- authoritative domain ownership remains unchanged;
- source/provenance is recorded;
- tenant/org/field authorization is preserved;
- rebuild/consistency behavior for derived stores is defined;
- evidence class remains explicit (`declared`, `static`, `observed`, `tested`, `production-observed`, `modelled`, `ai-inferred` as applicable);
- sensitive topology/data is redacted and access-controlled;
- a generated diagram/dashboard does not substitute for execution evidence.
