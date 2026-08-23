# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This is the first file to read inside `docs/ai/`. It never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or the active work-package specification.

## Project identity

**Omnexa** is a Composable Enterprise Business Operating System being built architecture-first as a strict modular monolith with service-ready boundaries. Development is sequential, gated and evidence-based: kernel capabilities are implemented and verified before later phases or business modules are authorized.

Core architecture retains one authoritative owner per capability/write model, governed cross-domain boundaries, mandatory isolation/authorization/audit/classification/versioning, non-authoritative infrastructure caches/telemetry, and future AI execution only through governed capabilities rather than direct database/object-store/payment/business-state authority.

## Verified P01.11 implementation result

Live GitHub evidence retained for P01.11:

- implementation PR: `#63` — merged;
- final exact source head: `1c1ab1f8d5120fb6b1e5908fdb93cffef9275940`;
- canonical implementation run/job: `32610902537` / `97123708250` — PASS;
- runner: GitHub-hosted `ubuntu-24.04`, `GitHub Actions 1000012523`, Ubuntu 24.04.4 LTS / X64;
- runner image: `ubuntu-24.04 / 20260816.277.1`;
- Go: `1.26.7`;
- repository Go quality: PASS (`gofmt`, `golangci-lint v2.12.2` 0 issues, `govulncheck v1.7.0` no vulnerabilities);
- completed P01.01-P01.10 regressions: PASS;
- P01.11 G1/G2/G3/G5/G6/G7: PASS;
- implementation squash merge: `10c94a638b89d47da05f5481fb2db298a2da6942`.

Initial run `32610720614 / 97123236672` remains retained as FAIL history for canonical gofmt alignment and govet test-variable shadow findings. Those findings were fixed directly in source/tests without adding linter exclusions or weakening governance.

## P01.11 delivered boundary

P01.11 established only `kernel.audit`: immutable classification-aware records; UUIDv7 IDs and UTC timestamps; tamper-evident integrity; descriptive actor/action/target/scope/outcome/correlation/reason/approval and privileged/impersonation metadata without P02 authority; append-only sink capability; explicit required-audit failure and best-effort degradation semantics; deterministic bounded memory sink; secret/prohibited-field rejection; and protected-payload-safe transport-health observability.

It did **not** build P02 identity/tenancy/role catalogs, business audit-event catalogs, audit read/export/reporting, legal hold/retention, database audit persistence, durable messaging/outbox/inbox, P01.12 CLI, business behavior or AI/model/agent runtime.

## P01.11 closure transition snapshot

Closure branch: `chore/p01-11-close-p01-12-activate`  
Base: P01.11 implementation merge `10c94a638b89d47da05f5481fb2db298a2da6942`.

The closure reconciles exactly:

- P01.01-P01.11 `done`;
- P01 progress `11 / 12 done`;
- P01.12 — Developer CLI Baseline `active` as the sole executable kernel package;
- P02+ remain `planned`;
- `kernel_code_authorized=true` only for P01.12;
- `business_feature_code_authorized=false`;
- no P01.12 runtime implementation in the closure transition.

Canonical P01.11 evidence is `docs/roadmap/evidence/P01.11_COMPLETION_2026-08-23.md`.

A future session must re-read protected `main`, `STATE.json`, the closure PR and latest CI. If the closure has merged, P01.12 is canonical active scope. If it has not merged, protected `main` still records P01.11 active even though implementation is merged and completion-eligible.

## P01.12 bounded direction after closure

P01.12 owner/domain is `kernel.developer`. Authorized scope after closure is limited to a stable repository-owned developer CLI/command surface for help/version/verify/build-test/approved diagnostics; deterministic fail-closed verification orchestration; explicit exit codes and structured-safe output; safe composition of P01 configuration/migration/diagnostic capabilities; no-secret/no-RESTRICTED output; clean-checkout/CI reproducibility; and the complete P01 fresh-install exit proof.

The P01 exit proof must accurately demonstrate configuration resolution, kernel build/start, fresh PostgreSQL migration, cache/storage provider contracts, safe logs/telemetry, readiness/diagnostics, jobs/configuration/audit primitive tests, canonical developer verification, and required security/supply-chain/build gates without hidden manual steps.

P01.12 must not implement production super-admin authority, P02 tenant/user/role administration, P03 module runtime administration, P04+ event/workflow/domain commands, deployment/Kubernetes orchestration, hidden SQL/file mutation, business modules or AI/model/agent runtime. CLI convenience does not create authority.

## Exact next authorized action

For this closure session, finish the P01.11 closure PR, obtain exact-head full GitHub-hosted `governance`, verify branch freshness and review/conversation state, and merge only if all required evidence is green. After the closure merge, **STOP**.

In the next governed execution session only, re-read protected `main` and `STATE.json`, then implement **P01.12 only** with focused positive/negative tests, a fail-closed P01.12 verifier, preserved P01.01-P01.11 regressions and the complete P01 exit proof.

Do not automatically start P02+, business features or AI/model/agent implementation.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/CHANGE_CONTROL.md`, `docs/governance/DEFINITION_OF_DONE.md`, active work-package spec, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current handoff.
