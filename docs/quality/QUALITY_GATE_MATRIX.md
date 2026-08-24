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
| Business Graph / master-data projection | source owner + provenance + tenant/field authorization + temporal/identity semantics + rebuild/consistency + no cross-domain write authority |
| Data lineage / derived store | source/transform/owner + classification + retention/delete/export propagation + rebuild/reconciliation + cross-tenant denial |
| System Graph / Flow evidence | node/edge schema + provider/source identity + evidence-class separation + expected-vs-observed drift + sensitive-topology authorization/redaction |
| Process intelligence/mining | event/activity provenance + tenant/business-scope authorization + conformance/variant correctness + modelled-vs-observed separation |
| AI model/provider/prompt/tool change | exact version/lineage + offline evals + risk/safety/grounding/task-success + latency/cost + fallback/rollback + affected online monitoring |
| AI agent / agent-team change | identity/scope/tool permissions + prompt/tool-injection tests + approval boundaries + multi-agent delegation limits + eval/regression + audit/replay/recovery |
| AI development governance/tooling | exact base/source identity + scope/allowed-path checks + test-oracle integrity + governance self-modification guard + independent review for high/critical changes |
| Financial SoD / continuous control | conflict/allow/deny matrix + maker-checker + privilege-escalation + affected transaction/config monitoring + immutable evidence |
| EPM/planning/simulation | model/version/input provenance + financial arithmetic + scenario isolation + actual-vs-plan separation + modelled label + approval before consequential execution |
| Environment/config promotion | semantic diff + dependency validation + secret-reference check + DEV/TEST/UAT/PROD promotion + rollback/drift evidence |
| Migration/onboarding | source profiling + mapping/transform validation + dry run + idempotent resume + duplicate handling + reconciliation + exception/audit evidence |
| Cost/FinOps attribution | source/provider pricing identity + allocation rules + tenant isolation + reconciliation/tolerance + no billing-authority inference |
| Business outcome/AI value claim | predefined metric/CTQ + observation window/source + baseline/comparison + uncertainty/limitations + no fabricated post-release result |

## Gate classes

### G0 — Governance

- canonical state/status consistency;
- ADR/change-control requirements;
- ownership/dependency validation;
- required evidence files;
- active strategic-program/work-package authority when applicable;
- AI-development scope/evidence authority checks when implemented.

### G1 — Static

- formatting;
- lint;
- type/static analysis;
- generated artifact drift;
- forbidden dependency checks;
- declared System Graph/package capability drift where implemented.

### G2 — Unit/component

- fast deterministic logic/component tests.

### G3 — Contract/integration

- API/event/schema/module contracts;
- database/cache/broker/storage/provider boundaries;
- graph/data/process/AI tool contracts where applicable.

### G4 — Data/migration

- fresh install;
- supported upgrade;
- deterministic seed/reference data;
- destructive-change controls;
- lineage/master-data/rebuild/reconciliation checks where applicable.

### G5 — Security/tenancy

- authentication/authorization;
- cross-tenant/org negative tests;
- classification/secrets/webhook/SSRF/privileged behavior;
- AI prompt/tool/delegation security where applicable;
- graph/analytics/vector/retrieval authorization;
- configured security scans.

### G6 — Lifecycle/resilience

- module lifecycle;
- retry/idempotency/replay;
- failure isolation;
- recovery/backup where applicable;
- agent/model/provider/graph/rebuild/kill-switch lifecycle where applicable.

### G7 — Build/package

- production build;
- target platform packaging;
- install/upgrade certification;
- exact reviewed/tested/built artifact identity where available.

### G8 — Supply chain/release

- dependency/license policy;
- secret/vulnerability scans;
- SBOM/provenance/signature;
- release notes/compatibility/promotion record;
- third-party model/agent/MCP/package provenance and revocation policy where applicable.

## Strategic AI-native evidence rules

Future strategic programs introduced through ADR-0011 use these additional rules:

1. Business Graph, Process Graph and System Graph remain separate evidence planes with explicit provenance.
2. `modelled` and `ai-inferred` evidence never satisfy an `observed`, `tested` or `production-observed` requirement.
3. A graph/dashboard/AI explanation does not become business write authority.
4. High/critical AI-assisted implementation cannot rely only on self-authored tests/prose; independent review and objective evidence are required.
5. Deleted/skipped/weakened tests, lowered thresholds or changed security policies are review-significant even when final CI is green.
6. Provider/payment/model/agent/supply-chain certification is profile-specific; a generic package PASS does not imply specialized security/compliance PASS.
7. Configuration, workflow, low-code and agent assets promoted between environments require version/diff/dependency/secret-reference/rollback evidence rather than hidden manual recreation.

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
