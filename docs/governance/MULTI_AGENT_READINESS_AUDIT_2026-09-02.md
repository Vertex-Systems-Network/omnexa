# Multi-Agent Repository Readiness Audit — 2026-09-02

Status: **M2 enforcement candidate; canonical acceptance requires this carrier to pass source + unchanged-promotion Governance and merge to protected main.**

Scope: development-agent orchestration only. This audit does not grant P04.05+, business-feature, product-agent or future-module authority.

## Executive conclusion

The current P04.04 wave can safely operate with **3 concurrent write agents plus Supervisor/review support** after this M2 carrier is accepted. No known repository-coordination blocker remains for that current envelope.

The remaining M3 developer CLI and M4 automated merge-queue work are scale/convenience improvements, not prerequisites for the current three-writer wave. Increasing the writer cap above 3 remains intentionally blocked until a real registered worker PR proves M2 worker-specific scope and live-main freshness enforcement in canonical CI.

## Audit matrix

| Area | Result | Evidence / control |
|---|---|---|
| Canonical execution authority | PASS | `STATE.json` remains sole execution cursor; P04.04 only |
| Module/task isolation | PASS | active task branches + bounded write/forbidden paths |
| Branch-first wave bootstrap | PASS | worker + Supervisor branches created before task work |
| Protected-main integration | PASS | active repository ruleset requires PR, strict `governance`, conversation resolution, no bypass |
| New-agent admission | PASS | protected-main-first onboarding + machine-readable slot ledger |
| No-capacity behavior | PASS | exact `Go Home Come Back Next Time`; no task/branch/lease granted |
| Current capacity | PASS | 3/3 worker slots occupied; 0 open |
| Main freshness | PASS | live `main` ref is authority; stored SHAs are audit snapshots |
| Completion signal | PASS | exact `Work Done and Submitted` contract |
| Supervisor interrupt/review | PASS | checkpoint/pause/review/merge/resume protocol |
| Merge alert + resync | PASS | exact all-agent alert + `Sync Complete — Resuming Work` |
| Deterministic dependency order | M2 ENFORCED | dependency DAG + merge-order validator |
| Write lease overlap | M2 ENFORCED | pairwise path overlap fails required Governance |
| Forbidden self-overlap | M2 ENFORCED | task write vs forbidden path collision fails |
| Migration reservation collision | M2 ENFORCED | duplicate owner/version/path fails |
| Worker PR path budget | M2 ENFORCED | registered worker changed paths must fit declared write/shared budget |
| Unknown agent branch | M2 ENFORCED | unregistered `agent/*` PR fails |
| Worker stale-main PR | M2 ENFORCED | PR base must equal live protected main and head must descend from it |
| Helper correctness | M2 ENFORCED | dependency-free focused unittest runs in required Governance |
| Canonical CI integration | M2 ENFORCED | validators run inside existing `governance` job before expensive Go tooling |
| Runtime/product authority | PASS / unchanged | no P04.05+, business, provider or P20 authority added |

## Protected-main control audit

Repository ruleset `21174858` (`main`) was read during this audit and is active. It requires pull-request integration, strict required status `governance`, review-thread resolution, rejects non-fast-forward/deletion and has no bypass actor. Therefore M2 validators wired into the existing required Governance path become merge-blocking controls once this carrier is accepted.

## Current active wave

`P04.04-WAVE-20260902-01` remains bounded to `kernel.events` transactional-outbox work:

1. `P04.04-T01` / Agent-01 / Outbox Core;
2. `P04.04-T02` / Agent-02 / PostgreSQL Persistence + reserved migration;
3. `P04.04-T03` / Agent-03 / Reliability tests;
4. `P04.04-T04` / Supervisor / integration verifier.

The three worker write budgets are non-overlapping. T02 depends on T01 for accepted adapter assumptions; T03 depends on T01/T02; Supervisor integration follows all three. The current migration reservation is only `kernel.events / 1 / kernel/migrations/kernel.events/1_create_transactional_outbox.sql`.

## M2 implementation boundary

M2 uses only Python standard library + Git and is run from `scripts/verify_go_quality.sh`, which is already part of the required `governance` job.

Implemented validators:

- `scripts/validate_agent_task.py`
- `scripts/validate_agent_leases.py`
- `scripts/detect_path_overlap.py`
- `scripts/validate_task_dependencies.py`
- `scripts/validate_agent_pr_scope.py`
- `scripts/validate_agent_base_sha.py`

Shared/test support:

- `scripts/agent_orchestration_common.py`
- `scripts/test_agent_orchestration_common.py`

The plan-level checks run on all Governance executions. Worker-specific PR path and stale-main checks apply to registered active worker/Supervisor task branches; unknown `agent/*` branches fail. Governance/control carriers are not misclassified as worker tasks, but they still cannot merge without the repository's normal protected Governance pipeline.

## Deliberately retained limits

### Writer cap remains 3

M2 being present is not enough evidence to immediately raise concurrency. At least one actual registered worker PR must demonstrate:

- exact declared-path acceptance;
- live-main base acceptance;
- canonical Governance success;
- no coordination repair caused by M2 itself.

Only after that evidence should the plan consider 4-5 concurrent writers.

### M3 is not a blocker

A developer CLI/task registry would improve ergonomics, but the active JSON plan, Supervisor workflow, issue #177 and required CI already provide the current wave's required coordination state.

### M4 is not a blocker

An automatic merge queue/webhook orchestrator would reduce manual Supervisor work at larger scale. Current main ruleset + deterministic Supervisor ordering + exact-head source/promotion Governance safely serialize the current three-writer wave.

## Final readiness decision

After this M2 carrier passes both Governance runs and merges to protected main:

**READY — current P04.04 multi-agent development may proceed with up to 3 concurrent write agents under the recorded slots/branches.**

No additional repository-governance preparation is required before starting those three tasks. Any future writer-cap increase, new module stream, or future phase still requires evidence/canonical authorization; it must not be inferred from this audit.
