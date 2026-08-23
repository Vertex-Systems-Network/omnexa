# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This is the first file to read inside `docs/ai/`. It never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or the active work-package specification.

## Project identity

**Omnexa** is a Composable Enterprise Business Operating System being built architecture-first as a strict modular monolith with service-ready boundaries. Development is sequential, gated and evidence-based: kernel capabilities are implemented and verified before business modules are authorized.

Core architecture retains one authoritative owner per capability/write model, governed cross-domain boundaries, mandatory isolation/authorization/audit/classification/versioning, non-authoritative infrastructure caches/telemetry, and future AI execution only through governed capabilities rather than direct database/object-store/payment/business-state authority.

## Verified P01.10 implementation result

Live GitHub evidence retained for P01.10:

- implementation PR: `#61` — merged;
- final exact source head: `4c9914e4641d0d6e94a895d0fcd16c3a6bf4d962`;
- canonical implementation run/job: `32609018028` / `97118796940` — PASS;
- runner: GitHub-hosted `ubuntu-24.04`, `GitHub Actions 1000012097`, Ubuntu 24.04.4 LTS / X64;
- runner image: `ubuntu-24.04 / 20260816.277.1`;
- Go: `1.26.7`;
- repository Go quality: PASS (`gofmt`, `golangci-lint v2.12.2` 0 issues, `govulncheck v1.7.0` no vulnerabilities);
- completed P01.01-P01.09 regressions: PASS;
- P01.10 G1/G2/G3/G5/G6/G7: PASS;
- implementation squash merge: `9d11b9250eb74ca2ade531ee58e8f905468cf103`.

Initial run `32608872763 / 97118409671` remains retained as FAIL history for one `nilnil` lint finding in cache-miss invalidation. It was fixed through explicit `(Change, bool, error)` semantics without weakening any gate.

## P01.10 delivered boundary

P01.10 established only `kernel.configuration`: typed runtime definitions, stable owner/version metadata, deterministic default/provider/fallback evaluation, future tenant/org/user UUIDv7 references as opaque metadata, bounded non-authoritative cache/refresh/invalidation, value-free change metadata, disable-only fail-closed kill switches, deterministic in-memory provider, provider panic/timeout containment and caller cancellation/deadline propagation.

It did **not** build experimentation/analytics, entitlement/licensing, tenant admin UI, business-module flags, P01.11 audit transport, P01.12 CLI, P02 identity/tenancy, durable messaging/workflow runtime, business behavior or AI runtime.

## P01.10 closure transition snapshot

Closure branch: `chore/p01-10-close-p01-11-activate`  
Base: P01.10 implementation merge `9d11b9250eb74ca2ade531ee58e8f905468cf103`.

The closure reconciles exactly:

- P01.01-P01.10 `done`;
- P01 progress `10 / 12 done`;
- P01.11 — Audit Transport Foundation `active` as the sole executable kernel package;
- P01.12 `planned`;
- `kernel_code_authorized=true` only for P01.11;
- `business_feature_code_authorized=false`;
- no P01.11 runtime implementation in the closure transition.

Canonical P01.10 evidence is `docs/roadmap/evidence/P01.10_COMPLETION_2026-08-23.md`.

A future session must re-read protected `main`, `STATE.json`, the closure PR and latest CI. If the closure has merged, P01.11 is canonical active scope. If it has not merged, protected `main` still records P01.10 active even though implementation is merged and completion-eligible.

## P01.11 bounded direction after closure

P01.11 owner/domain is `kernel.audit`. Authorized scope after closure is limited to a stable classification-aware audit envelope; actor/action/target/scope/outcome/correlation/reason/approval metadata without P02 identity; append-oriented sink interface; explicit required-audit failure semantics; classification/redaction enforcement; immutable UUIDv7/timestamp conventions; impersonation/privileged metadata; deterministic local/test sink; and bounded protected-payload-safe transport health observability.

Audit write capability does not imply read/export authority. Actor/scope metadata does not grant authority or create identity/tenancy. Required-audit failure cannot silently claim success. P01.11 must not implement P02 identity/role catalogs, business audit events, compliance/reporting UI, legal hold/retention systems, P01.12 CLI, durable messaging/outbox/inbox, later business behavior or AI/model/agent runtime.

## Exact next authorized action

If the P01.10 closure PR is still open, verify its exact-head full GitHub-hosted `governance` lane, branch freshness and review/conversation state, then merge only if all required evidence is green. After that merge, **STOP**. P01.11 implementation belongs to a new governed execution session.

Do not automatically start P01.12, P02+ or business/AI implementation.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/CHANGE_CONTROL.md`, `docs/governance/DEFINITION_OF_DONE.md`, active work-package spec, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current handoff.
