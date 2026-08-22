# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.10 — Feature flag & configuration registry**
- P01 progress: **9 / 12 done**
- P01.01–P01.09: **DONE**
- P01.10: **ACTIVE**
- P01.11–P01.12: **PLANNED / NOT ACTIVE**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**

## Completed P01 packages

Canonical completion evidence is retained for P01.01 through P01.09 under `docs/roadmap/evidence/`.

Latest completed package: **P01.09 — Job & scheduler primitives**. Implementation PR #59 passed final strict-up-to-date GitHub-hosted run `32605309150` / job `97109396616` and merged as `0bcafbfc52324acba1df9d8eff84a264dda0f233`. Canonical evidence: `docs/roadmap/evidence/P01.09_COMPLETION_2026-08-22.md`.

P01.09 established deterministic process-local job registration/execution, UUIDv7 execution IDs, bounded worker/queue concurrency, bounded retry/backoff with explicit idempotency protection, duplicate-safe in-memory execution hooks, repeatable completion handles, graceful queued/synchronous drain/cancel semantics, UTC-normalized one-shot/fixed-interval schedules, and P01.07 observability propagation. Scheduler/job identity remains operational identity rather than authorization, tenancy or business authority.

## Active P01 package — P01.10

P01.10 is the sole active package and owns only the governed feature flag/runtime configuration registry in `docs/roadmap/work-packages/P01.10.md`:

- typed flag/config definitions;
- stable identifiers, descriptions and owner metadata;
- deterministic defaults/fallback behavior;
- runtime evaluation boundary;
- future-scope-aware evaluation inputs without implementing P02 identity;
- version/change metadata hooks;
- bounded refresh/invalidation semantics;
- explicitly declared operational kill switches;
- deterministic test provider;
- P01.01-P01.09 regression preservation and repository Go quality.

P01.10 must not implement product experimentation/analytics, tenant admin UI, pricing/entitlement/licensing, authorization based solely on flags, business-module flags before their owners exist, P01.11 audit transport, P01.12 developer CLI, P02 identity/tenancy, later business behavior, or AI/model/agent/planner functionality.

Flags never grant authority or bypass authorization/data isolation. Security controls fail closed and cannot be disabled by undeclared generic flags. Sensitive configuration remains governed by classification/secrets policy rather than the runtime registry.

P01.11 may activate only after P01.10 reaches `done` with required canonical evidence.

## Repository Go quality

`docs/quality/GO_CODE_QUALITY.md` defines the permanent Go quality gate. Canonical `governance` executes `bash scripts/verify_go_quality.sh` before package regressions using pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

## Protected integration / CI

`main` remains protected with PR-only integration, strict required `governance`, blocked direct/force updates and required conversation resolution. Canonical governance CI remains GitHub-hosted only on `ubuntu-24.04`; local/self-hosted runners are prohibited. Completed P01.01-P01.09 verifiers remain regression gates in the same required job.

## Future web UI quality/accessibility plan

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` remains mandatory planning input whenever a future package actually authorizes browser UI. It does not authorize business/UI implementation during P01.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block authorized P01 kernel engineering.

## Exact next work

After this governed P01.09 closure/state-transition PR merges, implement **P01.10 only** in a separate executable PR. Add its fail-closed verifier, preserve repository Go quality and P01.01-P01.09 regressions, obtain GitHub-hosted G0/G1/G2/G3/G5/G6/G7 evidence, then use a separate closure transition before activating P01.11. Keep `business_feature_code_authorized=false` and P02+ implementation locked.
