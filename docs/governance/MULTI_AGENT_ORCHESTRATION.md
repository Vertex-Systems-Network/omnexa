# Omnexa Multi-Agent Development Orchestration

Status: **Mandatory development-governance contract**  
Scope: repository development agents/contributors; this is **not** Omnexa P20 product-agent runtime  
Strategic alignment: `XQ-100 — AI-Native Engineering Governance & Quality OS`

## 1. Purpose

Omnexa may use multiple human or AI development agents concurrently to increase delivery speed without allowing shared-state races, hidden scope expansion, cross-module writes, migration collisions or stale-head merges.

Parallelism is a delivery technique. It does not create roadmap authority. `docs/roadmap/STATE.json` remains the only machine-readable execution cursor.

## 2. Current concurrency envelope

Until automated lease enforcement and conflict-aware merge orchestration are proven in canonical CI:

- recommended active team: **4-6 agents**;
- maximum concurrent **write agents: 3**;
- additional agents should be read-only planners, reviewers, security analyzers, test designers or evidence coordinators;
- all concurrent writers must remain inside the currently authorized work package;
- foundation phases do not become phase-parallel merely because multiple agents exist.

For the current P04.04 pattern, a safe split is:

| Role | Typical authority |
|---|---|
| Orchestrator / planner | read-only; decomposes work and dependencies |
| Runtime/source writer | bounded owning kernel source paths |
| Persistence/migration writer | exact reserved owner/path/version/data budget only |
| Test writer | non-overlapping focused test paths |
| Security/architecture reviewer | read-only |
| CI/evidence/documentation coordinator | verifier/evidence/docs only |

The limit may be increased only after measured evidence shows path leases, scope validation, stale-head handling and merge ordering are reliable.

## 3. Three-level execution model

```text
Phase
  -> Work Package
       -> Parallel Task
            -> Agent / branch / path lease / evidence
```

The roadmap cursor remains at phase/work-package level. Parallel tasks are subordinate execution units and may never mark a package complete independently.

## 4. Required task identity

Every material concurrent task must record:

- `task_id`;
- `agent_id` or role identity;
- exact `base_sha`;
- phase/work-package;
- owning module/domain/kernel capability;
- task purpose;
- dependencies;
- read paths;
- write paths;
- forbidden paths;
- shared/exclusive paths;
- migration reservation if applicable;
- contract/event reservation if applicable;
- risk tier;
- required tests/gates;
- expected handoff/merge order.

The planning schema is `docs/ai/AGENT_TASK_SCHEMA.json`.

## 5. Path and scope lease rules

A lease is a bounded permission to mutate a declared path or logical resource for one task. It is not broader product authority.

### 5.1 Exclusive module/task write lease

Use for module-private or task-private source/test paths. Two active writers must not hold overlapping exclusive leases.

### 5.2 Shared-read lease

Multiple agents may read the same architecture, contract, source or evidence.

### 5.3 Exclusive shared-surface lease

The following are shared/high-conflict surfaces and require serialized ownership when modified:

- `AGENTS.md`;
- `docs/roadmap/STATE.json`;
- package-sequence/state-transition files;
- `.github/workflows/**`;
- `go.mod`, `go.sum`, `go.work`, `go.work.sum`;
- global registries;
- kernel public contracts;
- global event/capability/permission schemas;
- migration version namespaces;
- cross-module shared test fixtures;
- release/promotion control files.

### 5.4 Forbidden cross-owner write

A module/task agent may not directly write:

- another module's private source;
- another module's private tables or migrations;
- undocumented private contracts;
- shared governance/runtime paths outside the assigned budget.

When such a change is required, create a dependency/change request and serialize it through the owning task.

The lease planning schema is `docs/ai/AGENT_LEASE_SCHEMA.json`.

## 6. Dependency DAG

Parallel tasks must form an acyclic dependency graph.

Example:

```text
kernel contract
   -> module adapter
      -> integration test
```

Independent tasks may run together. A dependent task may analyze ahead but must not implement against an unaccepted guessed dependency contract.

## 7. SHA freshness and stale-task handling

Every task starts from an exact base SHA.

Before material mutation, PR review and merge:

1. re-read current protected main;
2. compare current assumptions with the recorded base;
3. identify changed owned/shared paths and contracts;
4. if the assumptions are stale, stop the write path;
5. rebase/re-plan/review as required;
6. obtain fresh exact-head CI evidence.

A previously green run is not merge authority after a material head change.

## 8. Migration reservation

Before any schema mutation, the task must reserve and document:

- authoritative schema owner;
- exact migration path;
- exact version/sequence;
- affected tables/indexes/constraints;
- data/backfill budget;
- fresh-install path;
- supported upgrade path;
- rollback or forward-recovery strategy;
- tenant/owner-isolation proof.

Two active tasks may not independently claim the same migration version or owner namespace.

## 9. Contract/event reservation

A task changing a public API, capability or event must declare the exact contract identity/version.

Two agents may not independently redefine the same public contract. The owning contract task lands first; dependents refresh their base and consume the accepted version.

## 10. Conflict detection

Before a task is assigned and before PR submission, compare:

- planned write paths vs active write leases;
- migration reservations;
- public contract/event identities;
- module/data ownership;
- dependency graph;
- current protected-main SHA.

Any material overlap requires one of:

- task decomposition;
- explicit dependency ordering;
- exclusive lease handoff;
- task merge into one writer;
- change-control escalation.

"Agents will probably edit different lines" is not an acceptable conflict strategy for shared authoritative files.

## 11. Merge model

Preferred order:

```text
task branches
  -> task CI/review
  -> dependency-aware integration order
  -> exact-head source acceptance
  -> unchanged promotion when required
  -> protected main
```

An integration/orchestrator role sequences accepted work but does not receive unrestricted permission to modify every module.

## 12. Agent Working Instructions synchronization

Every material task must check the effective instructions at task start and again before PR submission.

Required comparison:

- current phase/work-package;
- owner/module;
- role/task objective;
- base SHA;
- branch/workspace strategy;
- read/write/forbidden/shared paths;
- active leases and dependencies;
- migration/data budget;
- public contract/event assumptions;
- tests and evidence gates;
- CI/review/promotion process;
- allowed tools/network/secrets;
- stop conditions.

The human-readable mirror is the `README.md` section **Agent Working Instructions**.

If any effective working instruction changes materially, update that README section in the same governed PR.

If no instruction changed, the PR must record:

`Agent instructions checked — README instruction delta: none`

README is a mirror, not authority. `AGENTS.md`, `STATE.json`, mandatory governance policy and accepted ADRs win on conflict.

## 13. Current-vs-future scaling

### Current foundation stages

Use 4-6 agents with no more than 3 concurrent writers. Parallelize tasks inside the sole active work package.

### After modular business-domain parallelism is canonically opened

Independent domain teams may run concurrently where ownership and contracts are stable, for example CRM, Finance, Inventory, HR and Commerce.

Recommended scale progression:

| Maturity | Total agents | Concurrent writers |
|---|---:|---:|
| Manual lease discipline | 4-6 | <=3 |
| CI-enforced task/lease validation | 6-8 | 3-5 |
| Stable independent business modules | 8-12 | 5-8 |
| Proven automated orchestration/merge queue | 12-20+ | evidence-based |

Agent count is never a target by itself. Increase only while throughput rises without increasing stale work, conflicts, review debt or escaped defects.

## 14. Planned enforcement roadmap

### M0 — Manual semantics — effective immediately

- isolated branch/workspace;
- task identity;
- declared path budget;
- manual overlap check;
- exact base SHA;
- dependency ordering;
- README instruction synchronization;
- PR declaration and exact-head CI.

### M1 — Machine-readable planning contracts — defined now

- `docs/ai/AGENT_TASK_SCHEMA.json`;
- `docs/ai/AGENT_LEASE_SCHEMA.json`.

These schemas define the target record shape. They do not themselves create a distributed lock service.

### M2 — CI enforcement

Planned repository tooling:

- `scripts/validate_agent_task.py`;
- `scripts/validate_agent_leases.py`;
- `scripts/detect_path_overlap.py`;
- `scripts/validate_agent_pr_scope.py`;
- `scripts/validate_agent_base_sha.py`;
- `scripts/validate_task_dependencies.py`;
- PR/governance checks that fail undeclared path writes or conflicting active leases.

M2 implementation requires a separately authorized governance/tooling change; this document does not silently authorize it inside unrelated work.

### M3 — Developer CLI / task registry

Target commands:

```text
omnexa dev task create
omnexa dev task claim
omnexa dev task status
omnexa dev task release
omnexa dev task graph
```

### M4 — Conflict-aware merge orchestration

Target behavior:

- dependency-aware merge queue;
- automatic stale-task invalidation;
- lease expiry/release;
- contract/migration reservation checks;
- measured concurrency tuning.

## 15. Completion metrics

Measure whether multi-agent development is actually faster:

- lead time per work package;
- cycle time per task;
- CI queue time;
- stale-head invalidations;
- merge conflicts;
- overlapping-path violations;
- rework after integration;
- escaped regression count;
- reviewer wait time;
- percentage of tasks merged without coordination repair.

Increase concurrency only when these metrics remain healthy.

## 16. Non-authority statement

This contract changes development coordination only.

It does **not**:

- activate P04.05+;
- activate P05+;
- authorize business-feature code;
- authorize product AI/agent runtime;
- relax module/data ownership;
- relax PR/CI/review requirements;
- grant a task rights outside `STATE.json`.
