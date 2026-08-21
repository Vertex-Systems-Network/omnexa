# Omnexa Quality Gate Matrix

Status: **Canonical v1 control map**  
Work package: **P00.07**

This matrix maps change risk to mandatory evidence. A gate may be `N/A` only when the PR explains why the affected invariant cannot be exercised by that gate.

| Change/risk | Required minimum evidence |
|---|---|
| Governance/state/docs contract | governance/state/schema checks + source-of-truth reconciliation |
| Go/backend logic | format/lint/static + unit + affected integration/contract tests |
| TypeScript/web logic | format/lint/type + unit/component + affected browser/E2E tests |
| Rust/edge/native | format/lint/static + unit + target-platform integration/certification |
| Database schema | fresh migration + supported upgrade + deterministic seed/reference checks + rollback/forward-fix analysis |
| Tenant-owned data access | same-tenant success + cross-tenant denial + cross-org/object checks where applicable |
| Authorization/role/policy | allow + deny + privilege-escalation negative tests + audit where material |
| Public HTTP contract | OpenAPI validation + implementation contract tests + compatibility/deprecation analysis |
| Published event | event schema validation + producer tests + duplicate-delivery/replay/compatibility tests |
| Money/financial logic | exact precision + rounding + boundary + invariants + reversal/compensation tests where relevant |
| Time/scheduling | timezone/DST/date-boundary + recurrence semantics |
| Module/package | manifest/dependency/permission validation + install/enable/disable/re-enable/upgrade lifecycle tests |
| Optional module dependency | degradation/isolation tests proving unrelated domains remain healthy |
| External webhook | signature/authenticity + replay + malformed payload + tenant/provider binding tests |
| Configurable outbound URL | SSRF/allowlist/redirect/private-network negative tests |
| Files/media | upload validation + authorization + type/path/content handling + malware workflow where relevant |
| Secret/auth material | secret scan + redaction + rotation/revocation semantics where implemented |
| Bulk export | authorization + tenant/field/classification + audit + expiry/access tests |
| Search/analytics/vector projection | tenant/result authorization + classification + deletion/retention propagation tests |
| AI retrieval/tool | authorized retrieval + cross-tenant/object denial + prompt-injection/tool-denial + approval boundary tests |
| Payment/financial side effect | idempotency + duplicate retry + provider error + reconciliation + audit + reversal/compensation tests |
| Async job/workflow | retry + idempotency + timeout + compensation + current-authorization recheck where required |
| Performance/SLO-sensitive path | benchmark/load/resilience evidence against defined budget |
| Release packaging | build/package + digest + install/upgrade + SBOM/provenance/signature gates when available |

## Gate classes

### G0 — Governance

- canonical state/status consistency;
- ADR/change-control requirements;
- ownership/dependency validation;
- required evidence files.

### G1 — Static

- formatting;
- lint;
- type/static analysis;
- generated artifact drift;
- forbidden dependency checks.

### G2 — Unit/component

- fast deterministic logic/component tests.

### G3 — Contract/integration

- API/event/schema/module contracts;
- database/cache/broker/storage/provider boundaries.

### G4 — Data/migration

- fresh install;
- supported upgrade;
- deterministic seed/reference data;
- destructive-change controls.

### G5 — Security/tenancy

- authentication/authorization;
- cross-tenant/org negative tests;
- classification/secrets/webhook/SSRF/privileged behavior;
- configured security scans.

### G6 — Lifecycle/resilience

- module lifecycle;
- retry/idempotency/replay;
- failure isolation;
- recovery/backup where applicable.

### G7 — Build/package

- production build;
- target platform packaging;
- install/upgrade certification.

### G8 — Supply chain/release

- dependency/license policy;
- secret/vulnerability scans;
- SBOM/provenance/signature;
- release notes/compatibility/promotion record.

## PR rule

Every material PR maps affected risk to gate classes and records one of:

```text
PASS
FAIL
BLOCKED
NOT RUN
N/A
```

`N/A` requires rationale. `BLOCKED` requires a named blocker. Neither is silently equivalent to PASS.

## P00 exception

ADR-0006 may temporarily permit P00 documentation/specification PRs to merge while hosted CI execution is `BLOCKED`, provided manual evidence satisfies the documented exception. G2-G8 executable product gates remain N/A during P00 because implementation is locked; the exception cannot carry into P01 runtime merges.