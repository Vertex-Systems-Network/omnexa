# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This is the first file to read inside `docs/ai/`. It never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or the active work-package specification.

## Project identity

**Omnexa** is a Composable Enterprise Business Operating System being built architecture-first as a strict modular monolith with service-ready boundaries. Development is sequential, gated and evidence-based: kernel capabilities are implemented and verified before business modules are authorized.

Core architecture retains one authoritative owner per capability/write model, governed cross-domain boundaries, mandatory isolation/authorization/audit/classification/versioning, non-authoritative infrastructure caches/telemetry, and future AI execution only through governed capabilities rather than direct database/object-store/payment/business-state authority.

## Verified P01.09 implementation result

Live GitHub evidence verified during this session:

- implementation PR: `#59` — merged;
- final strict-up-to-date source head: `61e9c1115d05300ac9aedf5a555138c6a5a5be1e`;
- canonical implementation run/job: `32605309150` / `97109396616` — PASS;
- runner: GitHub-hosted `ubuntu-24.04`, `GitHub Actions 1000011549`, Ubuntu 24.04.4 LTS / X64;
- Go: `1.26.7`;
- repository Go quality: PASS (`gofmt`, `golangci-lint v2.12.2` 0 issues, `govulncheck v1.7.0` no vulnerabilities);
- completed P01.01-P01.08 regressions: PASS;
- P01.09 G1/G2/G3/G5/G6/G7: PASS;
- implementation squash merge: `0bcafbfc52324acba1df9d8eff84a264dda0f233`.

Failed diagnostic runs are retained and are not completion evidence: `32604728124 / 97108038830` exposed formatting/shadowing findings; `32604931754 / 97108518734` exposed deterministic unknown-job/lifecycle admission ordering. Both were fixed without weakening required gates before the final full PASS.

## P01.09 delivered boundary

P01.09 established process-local `kernel.jobs` primitives only: deterministic type registry, UUIDv7 execution IDs, bounded synchronous/queued execution, bounded retry/backoff with explicit idempotency for retries, in-memory duplicate/conflict protection, repeatable completion handles, cancellation/deadline propagation, graceful queued/synchronous drain and deadline cancellation, one-shot/fixed-interval UTC schedules, safe P01.07 observability propagation, panic containment and race-oriented evidence.

It did **not** build durable messaging, outbox/inbox, distributed workflow timers, tenant runtime, business jobs, later P01 capabilities or AI behavior.

## Closure transition snapshot

Closure branch: `chore/p01-09-close-p01-10-activate`  
Base: P01.09 implementation merge `0bcafbfc52324acba1df9d8eff84a264dda0f233`.

The closure proposes exactly:

- P01.01-P01.09 `done`;
- P01 progress `9 / 12 done`;
- P01.10 — Feature flag & configuration registry `active` as the sole executable kernel package;
- P01.11-P01.12 `planned`;
- `kernel_code_authorized=true` only for P01.10;
- `business_feature_code_authorized=false`;
- no runtime/dependency implementation in the closure transition.

A future session must re-read protected `main`, `STATE.json`, this closure PR and latest CI. If the closure PR has merged, P01.10 is canonical active scope. If it has not merged, protected `main` still records P01.09 active even though implementation is merged and completion-eligible.

## P01.10 bounded direction after closure

P01.10 owner/domain is `kernel.configuration`. Authorized scope after the closure is limited to typed flag/config definitions, stable identifiers/descriptions/owner metadata, deterministic defaults/fallbacks, runtime evaluation distinct from static P01.02 configuration, future-scope-aware evaluation inputs without P02 identity, version/change metadata hooks, bounded refresh/invalidation, explicitly declared operational kill switches and a deterministic test provider.

Flags cannot grant authority or bypass authorization/data isolation. P01.10 must not implement experimentation/analytics, tenant admin UI, pricing/entitlement/licensing, business-module flags before owners exist, P01.11/P01.12 behavior, P02+ identity/business behavior or AI runtime.

## Exact next authorized action

If the P01.09 closure PR is still open, verify its exact-head full GitHub-hosted `governance` lane, branch freshness and review/conversation state, then merge only if all required evidence is green. If the closure has already merged, re-verify `STATE.json` and begin **P01.10 only** in a new executable PR under `docs/roadmap/work-packages/P01.10.md`.

Do not automatically start P01.11+, P02+ or business/AI implementation.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/CHANGE_CONTROL.md`, `docs/governance/DEFINITION_OF_DONE.md`, active work-package spec, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the current handoff.
