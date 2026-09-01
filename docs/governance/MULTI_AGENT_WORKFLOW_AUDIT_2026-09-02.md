# Multi-Agent Repository Workflow Deep Audit — 2026-09-02

Status: **Governance audit / implementation input**  
Scope: XQ-100 development orchestration only  
Canonical execution cursor: unchanged; P04.04 remains the sole active work package

## Audit objective

Audit the existing Omnexa multi-agent development plan against a Supervisor-led repository workflow that requires branch-first parallel setup, explicit completion signals, interrupt-driven review, deterministic protected-main merging, all-agent synchronization after every accepted merge, and reusable behavior for future independent modules.

## Evidence reviewed

- `AGENTS.md`
- `docs/roadmap/STATE.json`
- `docs/governance/AI_EXECUTION_POLICY.md`
- `docs/governance/MULTI_AGENT_ORCHESTRATION.md`
- `docs/roadmap/XQ_100_MULTI_AGENT_DEVELOPMENT_PLAN.md`
- `docs/ai/AGENT_TASK_SCHEMA.json`
- `docs/ai/AGENT_LEASE_SCHEMA.json`
- `.github/pull_request_template.md`
- `README.md`
- active P04.04 contract and current `kernel/internal/events`/migration layout

## Executive result

The existing plan had a strong foundation: isolated branches, task identity, path leases, dependency DAGs, exact base SHA, migration/contract reservations, conflict-aware merge order and README instruction synchronization were already defined.

However, it did not yet define a complete **Supervisor runtime protocol** for a live repository wave. The most important missing controls were branch-first bootstrap, exact submission signals, Supervisor interruption/resume behavior, a persistent coordination channel, mandatory all-agent synchronization after each protected-main merge, and machine-readable active branch/agent/merge-order state.

These gaps are addressed by `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`, `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`, `docs/ai/AGENT_SIGNAL_SCHEMA.json`, README instructions and PR coordination fields.

## Findings and disposition

| ID | Severity | Finding | Risk | Disposition |
|---|---|---|---|---|
| MA-01 | HIGH | Supervisor role was implicit, not normative. | No single review/merge coordinator; inconsistent merge behavior. | Add mandatory Supervisor role and duties. |
| MA-02 | HIGH | No branch-first bootstrap gate. | Agents may begin from different bases or edit before isolation exists. | First repository mutation after canonical readback is creation of all wave branches. |
| MA-03 | HIGH | No exact task-completion signal. | Supervisor cannot deterministically distinguish ready work from partial progress. | Require exact `Work Done and Submitted` signal plus SHA/PR/CI metadata. |
| MA-04 | HIGH | No interrupt/resume protocol for Supervisor's own task. | Submitted work may sit waiting or Supervisor may merge without a stable checkpoint. | Supervisor pauses own task, reviews, merges/rejects, synchronizes and resumes. |
| MA-05 | HIGH | No mandatory post-merge synchronization across all active branches. | Workers continue on stale main and generate avoidable conflicts/rework. | Required global merge alert + branch sync before any further mutation. |
| MA-06 | HIGH | User-style direct merge wording could conflict with protected Omnexa promotion rules. | Governance bypass / stale or unverified merge. | “Supervisor merge” means execute accepted PR/Governance/promotion pipeline, never bypass it. |
| MA-07 | MEDIUM | No live AI-native branch/agent/merge-order document. | Assignments and dependencies exist only in prose/PR context. | Add machine-readable `ACTIVE_MULTI_AGENT_PLAN.json`. |
| MA-08 | MEDIUM | No durable coordination channel. | Completion/alert messages can be lost between sessions. | Current wave uses GitHub issue #177; future orchestrator may use webhooks. |
| MA-09 | MEDIUM | No missed-alert fallback. | Sleeping/offline agent may resume without seeing a broadcast. | Active plan carries `required_main_sha`; every material mutation rechecks it. |
| MA-10 | MEDIUM | No simultaneous-submission queue policy. | Supervisor starvation, arbitrary first-finished merge ordering. | Priority queue: blockers, dependency-unblockers, contract/migration owners, oldest. |
| MA-11 | MEDIUM | Supervisor's own branch/review conflict not explicit. | Same authority path can write and approve its own material work. | Supervisor uses isolated branch; independent review preferred/required where configured by risk. |
| MA-12 | MEDIUM | No stale-submission resignal rule. | Old `Work Done` signal may be reused after main/head changes. | Any material head sync requires fresh tests/review and a new completion signal. |
| MA-13 | LOW | No branch/lease release lifecycle. | Old branches/leases appear active and create false conflicts. | Release leases after protected-main readback; do not reuse completed branches for unrelated work. |
| MA-14 | INFO | Git alone cannot literally pause or message arbitrary external AI processes. | Protocol may be mistaken for process-control capability. | Define pause/alert as agent-observed workflow state; issue #177 + SHA gate is current durable mechanism. |

## Current P04.04 repository-path audit

Current event implementation is concentrated under `kernel/internal/events`, including accepted P04.01-P04.03 files such as `envelope.go`, `bus.go` and `durable.go`. P04.04 should prefer additive new files rather than concurrent edits to those accepted shared files unless re-planned.

Migration convention is owner-scoped:

`kernel/migrations/<owner>/<version>_<name>.sql`

There is no current `kernel.events` migration directory at the accepted base. Therefore the active wave may reserve, but not create before governed authorization, the first owner-local migration identity:

`kernel/migrations/kernel.events/1_create_transactional_outbox.sql`

Reservation constraints:

- owner: `kernel.events`;
- version: `1`;
- new P04.04 outbox persistence only;
- no business-table schema mutation;
- no backfill of business data;
- no destructive data operation;
- no P04.05/P04.06 semantics;
- retained P01 migration runner/ledger remains authoritative.

## Parallelism audit

The requested model is safe if “parallel” means branches/tasks are prepared concurrently but dependency gates are respected.

For the current wave:

- T01 Outbox Core can begin immediately within additive event files.
- T02 Persistence/Migration can prepare owner-local schema work, but its Go adapter must not guess T01's unaccepted API; it synchronizes after T01 if the interface changes/lands first.
- T03 Reliability Tests has its branch created immediately but implementation tests that bind concrete APIs wait for accepted T01/T02 contracts.
- Supervisor Integration/Verifier can prepare review/verifier structure, but final integration depends on accepted worker heads.

This produces useful parallel preparation without violating the existing “no guessed dependency contract” rule.

## Merge-order audit

Current deterministic merge order:

1. `P04.04-T01` — Outbox Core;
2. `P04.04-T02` — PostgreSQL persistence + owner migration;
3. `P04.04-T03` — crash/restart/concurrency/tenant-isolation reliability tests;
4. `P04.04-T04` — Supervisor integration/verifier/governance hook.

A dependency or security finding can change this order only by updating the active plan and README instructions when material.

## Signal audit

Required exact worker completion phrase:

`Work Done and Submitted`

Required exact Supervisor merge alert:

`New changes have been merged — please merge these changes into your branch first, then resume your own work.`

Recommended sync acknowledgement:

`Sync Complete — Resuming Work`

A signal is valid only when its machine metadata satisfies `docs/ai/AGENT_SIGNAL_SCHEMA.json`.

## Safety conclusion

The revised model is suitable for the current P04.04 wave and reusable later for CRM, Finance, Inventory, HR, Commerce and other independent modules **only after canonical roadmap gates open them**.

The plan accelerates development by allowing multiple isolated tasks and continuous review while preserving one protected-main truth. It does not make locked phases parallel, weaken CI, permit shared-path races, or grant the Supervisor unrestricted product-write authority.
