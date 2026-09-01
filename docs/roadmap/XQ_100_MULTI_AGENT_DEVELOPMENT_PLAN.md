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

## Immediate operating target

For P04/P05 foundation work:

- 4-6 active agents;
- no more than 3 concurrent writers;
- one orchestrator/planner;
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
├── T03 relay integration against accepted P04.02 abstraction
├── T04 crash/restart/concurrency tests
├── T05 tenant/owner isolation and security review
└── T06 verifier/evidence/documentation
```

The exact split must be derived from the accepted P04.04 contract and actual repository paths. This example grants no path authority by itself.

## Future business-module model

After canonical roadmap gates open domain parallelism:

```text
CRM agent/team ----------\
Finance agent/team -------\
Inventory agent/team ------> governed contracts/events -> integration queue -> main
HR agent/team ------------/
Commerce agent/team -----/
```

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

### XQ-MA.01 — Manual task/lease discipline

Deliverables:

- mandatory orchestration contract;
- README working-instruction mirror;
- PR checklist;
- task/lease machine-readable schema definitions.

Success:

- every concurrent PR declares task/base/path/dependency information;
- no undeclared overlapping writer paths;
- README instruction delta explicitly checked.

### XQ-MA.02 — CI scope and overlap enforcement

Planned:

- task manifest validator;
- path-overlap detector;
- PR-diff vs declared-write-path validator;
- stale-base detector;
- lease/dependency validator.

Success:

- CI blocks out-of-budget writes and incompatible concurrent leases.

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
- dependency visualization.

### XQ-MA.05 — Conflict-aware merge queue

Planned:

- merge sequencing from dependency DAG;
- automatic stale-head invalidation;
- affected-task notification/revalidation;
- concurrency metrics.

## Agent instruction lifecycle

At both task start and PR submission:

```text
read authority
 -> resolve active scope
 -> resolve owner
 -> resolve role/task
 -> resolve base SHA
 -> resolve path budget/leases
 -> resolve dependencies
 -> resolve tests/gates
 -> compare working instructions
 -> update README if changed
 -> execute/review
```

If instructions changed, README is updated in the same PR. If not, PR records `Agent instructions checked — README instruction delta: none`.

## Metrics and scaling gate

Do not increase concurrent writers based on perceived speed alone.

Scale only if:

- median package cycle time improves;
- stale-work rate does not materially rise;
- merge conflicts remain low;
- CI queue time remains acceptable;
- review quality does not fall;
- rework/regression rate remains stable or improves.

## Roadmap compatibility

This plan operationalizes the existing XQ-100 concepts of run identity, scope leases, concurrent-writer conflict protection and multi-agent coordination.

It does not alter P00-P27 numbering, does not activate a future phase, and does not authorize product AI agents or business-feature implementation.
