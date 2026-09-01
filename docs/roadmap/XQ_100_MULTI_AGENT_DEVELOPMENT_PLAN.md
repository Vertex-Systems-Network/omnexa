# XQ-100 Multi-Agent Development Acceleration Plan

Status: **PLANNED / governance semantics partially active**  
Parent strategic program: `XQ-100 — AI-Native Engineering Governance & Quality OS`  
Execution authority: subordinate to `AGENTS.md`, `docs/roadmap/STATE.json` and mandatory governance  
Current roadmap cursor: unchanged; P04.04 remains the sole active work package

## Objective

Increase Omnexa development throughput by decomposing an authorized work package into non-overlapping parallel tasks while preserving module ownership, contract integrity, migration safety, exact-head evidence and protected-main integration.

This plan does not make foundation phases parallel. It makes safe task-level parallelism possible inside the currently authorized package and becomes reusable for independent business modules later.

## Design principles

1. **One roadmap cursor, many subordinate tasks.**
2. **One authoritative writer per leased path/resource.**
3. **Read widely, write narrowly.**
4. **No guessed dependency contracts.**
5. **No cross-module private writes.**
6. **Every task starts from an exact SHA.**
7. **Shared contracts/migrations serialize.**
8. **Review and CI attach to an exact final head.**
9. **Instruction changes are synchronized to README.**
10. **More agents are added only when measured throughput improves.**
11. **A new parallel wave creates all worker/Supervisor branches before task work begins.**
12. **Supervisor approval never bypasses protected-main Governance/promotion rules.**
13. **Every accepted main merge forces active branches to synchronize before resuming.**

## Immediate operating target

For P04/P05 foundation work:

- 4-6 active agents;
- no more than 3 concurrent writers;
- one Supervisor/integration coordinator with its own isolated task branch;
- one or two bounded source writers where paths do not overlap;
- one focused persistence/migration writer where authorized;
- one test writer;
- one read-only security/architecture reviewer;
- one CI/evidence/documentation role as needed.

A typical P04.04 decomposition:

```text
P04.04
├── T01 outbox domain/transaction contract implementation
├── T02 PostgreSQL persistence + reserved migration
├── T03 relay/reliability/concurrency tests after required contracts
└── T04 Supervisor integration/verifier/Governance hook
```

The exact split must be derived from the accepted P04.04 contract and actual repository paths. This example grants no path authority by itself.

## Supervisor-led repository workflow

The mandatory detailed protocol is `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`.

For each new wave:

```text
mandatory canonical readback
  -> create every worker branch + Supervisor task branch
  -> verify common base SHA
  -> write/refresh AI-Native active-wave plan
  -> assign tasks/leases/dependencies
  -> parallel work
  -> worker: Work Done and Submitted
  -> Supervisor pauses own task
  -> Supervisor reviews exact submitted head
  -> governed source/promotion/main merge
  -> protected-main readback
  -> merge alert to all workers
  -> every active branch syncs new main
  -> agents resume
```

The Supervisor may execute the final merge operation but cannot waive exact-head CI, promotion, review-thread or protected-main freshness requirements.

## Active P04.04 wave — P04.04-WAVE-20260902-01

Machine-readable authority mirror: `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`  
Persistent coordination channel: GitHub issue `#177`  
Bootstrap/current required main for this initial plan: `8fd737b318947c9d0f3cc0e5c5a0931636c2c40c`

All branches were created before task work and synchronized to that SHA.

| Merge order | Task | Logical agent | Branch | Bounded responsibility |
|---:|---|---|---|---|
| 1 | P04.04-T01 | Agent-01 | `agent/20260902-p04-04-outbox-core` | additive outbox core contract/runtime |
| 2 | P04.04-T02 | Agent-02 | `agent/20260902-p04-04-persistence-migration` | PostgreSQL adapter + owner-local migration |
| 3 | P04.04-T03 | Agent-03 | `agent/20260902-p04-04-reliability-tests` | crash/restart/concurrency/tenant isolation tests |
| 4 | P04.04-T04 | Supervisor | `supervisor/20260902-p04-04-integration` | verifier, Governance hook, integration sequencing |

### Current write-path budgets

**T01 — Outbox Core**

- `kernel/internal/events/outbox.go`
- `kernel/internal/events/outbox_core_test.go`

**T02 — Persistence / Migration**

- `kernel/internal/events/outbox_postgres.go`
- `kernel/internal/events/outbox_postgres_integration_test.go`
- provisional migration reservation `kernel/migrations/kernel.events/1_create_transactional_outbox.sql`

Migration reservation is owner `kernel.events`, version `1`, new outbox persistence only. It forbids business-table mutation, business-data backfill, destructive data operations and P04.05/P04.06 semantics.

**T03 — Reliability Tests**

- `kernel/internal/events/outbox_reliability_test.go`
- `kernel/internal/events/outbox_concurrency_test.go`

**T04 — Supervisor Integration / Verifier**

- `scripts/verify_p04_04.sh`
- `.github/workflows/governance.yml` additive P04.04 verifier integration only

Accepted P04.01-P04.03 implementation files such as `envelope.go`, `bus.go` and `durable.go` are treated as shared/forbidden for these initial worker write budgets unless a later governed re-plan explicitly changes ownership.

### Dependency and parallel strategy

- T01 may start immediately after this wave plan is accepted.
- T02 may prepare migration/schema work under the accepted P04.04 record law, but must not guess an unaccepted T01 Go API; synchronize after T01 before final adapter submission when assumptions change.
- T03 may design tests immediately, but concrete API-bound test implementation waits for/synchronizes accepted T01/T02 contracts.
- T04 may prepare verifier/review structure, but final integration follows T01/T02/T03.

Thus all branches exist and agents can make useful progress without pretending dependency contracts are already accepted.

### Required signals

Worker/supervisor completion:

`Work Done and Submitted`

Supervisor merge broadcast:

`New changes have been merged — please merge these changes into your branch first, then resume your own work.`

Recommended worker sync acknowledgement:

`Sync Complete — Resuming Work`

Review rejection signal:

`Review Changes Required`

Machine-readable signal requirements are in `docs/ai/AGENT_SIGNAL_SCHEMA.json`.

## Future business-module model

After canonical roadmap gates open domain parallelism:

```text
CRM agent/team ----------\
Finance agent/team -------\
Inventory agent/team ------> Supervisor + governed contracts/events -> integration queue -> main
HR agent/team ------------/
Commerce agent/team -----/
```

For each future wave the Supervisor first creates the branch for every authorized parallel module plus its own integration branch, then records branch→module→agent→lease→merge order in the AI-Native active plan.

Each module remains responsible for its own:

- private source;
- domain invariants;
- schema and migrations;
- public capabilities;
- events;
- permissions;
- tests.

Shared kernel/public-contract work is a separate dependency owned by the appropriate platform task.

## Milestones

### XQ-MA.01 — Manual task/lease/Supervisor discipline

Deliverables:

- mandatory orchestration contract;
- mandatory Supervisor workflow extension;
- README working-instruction mirror;
- PR checklist;
- task/lease/signal machine-readable schema definitions;
- active-wave machine-readable plan;
- durable coordination channel.

Success:

- every concurrent PR declares task/base/path/dependency information;
- all wave branches are created before work starts;
- no undeclared overlapping writer paths;
- exact `Work Done and Submitted` submissions are reviewable by Supervisor;
- every protected-main merge triggers all-active-branch synchronization;
- README instruction delta explicitly checked.

### XQ-MA.02 — CI scope and overlap enforcement

Planned:

- task manifest validator;
- path-overlap detector;
- PR-diff vs declared-write-path validator;
- stale-base / required-main detector;
- lease/dependency validator;
- signal/active-plan consistency validator.

Success:

- CI blocks out-of-budget writes, incompatible concurrent leases and stale required-main state.

### XQ-MA.03 — Migration and contract reservation

Planned:

- machine-readable migration reservations;
- event/API/capability version reservation;
- owner namespace validation.

Success:

- concurrent tasks cannot independently claim the same migration or public-contract identity.

### XQ-MA.04 — Task registry and developer CLI

Planned:

- create/claim/status/release/graph commands;
- task lifecycle and lease expiry;
- dependency visualization;
- Supervisor queue/status commands.

### XQ-MA.05 — Conflict-aware merge queue and broadcast

Planned:

- merge sequencing from dependency DAG;
- automatic stale-head invalidation;
- affected-task notification/revalidation;
- webhook/broadcast integration for external agent runtimes;
- automatic branch-sync acknowledgement tracking;
- concurrency metrics.

## Agent instruction lifecycle

At both task start and PR submission:

```text
read authority
 -> resolve active scope
 -> resolve active wave / required main SHA
 -> resolve owner
 -> resolve role/task
 -> resolve branch + base/sync SHA
 -> resolve path budget/leases
 -> resolve dependencies
 -> resolve tests/gates
 -> check coordination channel / Supervisor state
 -> compare working instructions
 -> update README if changed
 -> execute/review
```

If instructions changed, README is updated in the same PR. If not, PR records `Agent instructions checked — README instruction delta: none`.

After every accepted protected-main merge, each active agent synchronizes that main SHA before resuming and re-establishes exact-head evidence on its changed branch.

## Deep-audit record

The 2026-09-02 workflow audit, findings and dispositions are retained in:

`docs/governance/MULTI_AGENT_WORKFLOW_AUDIT_2026-09-02.md`

Important resolved risks include direct-main bypass, missed alerts, simultaneous submissions, stale completion signals, Supervisor self-review conflicts, branch/lease release and the fact that repository coordination cannot physically pause arbitrary external AI processes without their cooperation.

## Metrics and scaling gate

Do not increase concurrent writers based on perceived speed alone.

Scale only if:

- median package cycle time improves;
- stale-work rate does not materially rise;
- merge conflicts remain low;
- CI queue time remains acceptable;
- review quality does not fall;
- rework/regression rate remains stable or improves;
- Supervisor review queue does not become the new bottleneck;
- sync-alert-to-resume latency remains acceptable.

## Roadmap compatibility

This plan operationalizes the existing XQ-100 concepts of run identity, scope leases, concurrent-writer conflict protection and multi-agent coordination.

It does not alter P00-P27 numbering, does not activate a future phase, and does not authorize product AI agents or business-feature implementation.
