# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Kernel Program**
- Phase: **P01 — Omnexa Kernel**
- Phase state: **active**
- Current work package: **P01.09 — Job & scheduler primitives**
- P01 progress: **8 / 12 done**
- P01.01–P01.08: **DONE**
- P01.09: **ACTIVE**
- P01.10–P01.12: **PLANNED / NOT ACTIVE**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation: **AUTHORIZED ONLY FOR ACTIVE P01 PACKAGE**
- Business-feature implementation: **NOT AUTHORIZED**

## Completed P01 packages

Canonical completion evidence is retained for P01.01 through P01.08 under `docs/roadmap/evidence/`.

Latest completed package: **P01.08 — Health, readiness & diagnostics**. Implementation PR #56 passed final strict-up-to-date GitHub-hosted run `32601741049` / job `97100949202` and merged as `a2a93454bb283464bfa144bb5a38539041e40069`. Canonical evidence: `docs/roadmap/evidence/P01.08_COMPLETION_2026-08-22.md`.

P01.08 established portable process liveness/readiness semantics, criticality-aware dependency checks, bounded timeout/cancellation behavior, safe diagnostic projection, startup/stopping lifecycle states, build identity and P01.07 observability integration. Health output remains operational evidence rather than authorization, tenancy or business-state authority.

## Active P01 package — P01.09

P01.09 is the sole active package and owns only the bounded kernel-local job/scheduler primitives in `docs/roadmap/work-packages/P01.09.md`:

- deterministic job identity/type registration;
- kernel-local enqueue/execute result model;
- bounded worker concurrency and graceful shutdown;
- cancellation/deadline propagation;
- bounded retry/backoff metadata;
- idempotency-key hook and duplicate-safe handler contract;
- simple recurring/one-shot kernel maintenance schedules;
- correlation/observability propagation;
- deterministic in-memory/test harness;
- P01.01-P01.08 regression preservation and repository Go quality.

P01.09 must not implement NATS/JetStream durable streams, transactional outbox/inbox, P05 workflow timers, business jobs, tenant-context runtime, feature registry, audit transport, developer CLI, later P01/P02+ behavior, or AI/model/agent/planner functionality.

Job/scheduler identity never implies authority. Future tenant/actor context must be explicit and revalidated. Retry behavior must be bounded, duplicate-safe where protected effects are possible, and shutdown must stop accepting work with bounded drain/cancel semantics.

P01.10 may activate only after P01.09 reaches `done` with required canonical evidence.

## Repository Go quality

`docs/quality/GO_CODE_QUALITY.md` defines the permanent Go quality gate. Canonical `governance` executes `bash scripts/verify_go_quality.sh` before package regressions using pinned `golangci-lint v2.12.2` and `govulncheck v1.7.0`.

## Protected integration / CI

`main` remains protected with PR-only integration, strict required `governance`, blocked direct/force updates and required conversation resolution. Canonical governance CI remains GitHub-hosted only on `ubuntu-24.04`; local/self-hosted runners are prohibited. Completed P01 package verifiers remain regression gates in the same required job.

## Future web UI quality/accessibility plan

`docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md` remains mandatory planning input whenever a future package actually authorizes browser UI. It does not authorize business/UI implementation during P01.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. This does not block authorized P01 kernel engineering.

## Exact next work

After this governed P01.08 closure/state-transition PR merges, implement **P01.09 only**. Add its fail-closed verifier, preserve repository Go quality and P01.01-P01.08 regressions, obtain GitHub-hosted G0/G1/G2/G3/G5/G6/G7 evidence, then use a separate closure transition before activating P01.10. Keep `business_feature_code_authorized=false` and P02+ implementation locked.
