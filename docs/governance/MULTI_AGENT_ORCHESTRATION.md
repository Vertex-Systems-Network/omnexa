# Omnexa Multi-Agent Development Orchestration

Status: **Mandatory development-governance contract**  
Scope: repository development agents/contributors; this is **not** Omnexa P20 product-agent runtime  
Strategic alignment: `XQ-100 — AI-Native Engineering Governance & Quality OS`

## 1. Purpose

Omnexa may use multiple human or AI development agents concurrently to increase delivery speed without allowing shared-state races, hidden scope expansion, cross-module writes, migration collisions or stale-head merges.

Parallelism is a delivery technique. It does not create roadmap authority. `docs/roadmap/STATE.json` remains the only machine-readable execution cursor.

## 2. Current concurrency envelope

Current safe operating envelope:

- recommended active team: **4-6 agents**;
- maximum concurrent **write agents: 3** while the first real worker PR proves the new M2 worker-specific scope/freshness gates;
- additional agents should be read-only planners, reviewers, security analyzers, test designers or evidence coordinators;
- all concurrent writers must remain inside the currently authorized work package;
- foundation phases do not become phase-parallel merely because multiple agents exist.

For the current P04.04 pattern:

| Role | Typical authority |
|---|---|
| Orchestrator / planner | read-only; decomposes work and dependencies |
| Runtime/source writer | bounded owning kernel source paths |
| Persistence/migration writer | exact reserved owner/path/version/data budget only |
| Test writer | non-overlapping focused test paths |
| Security/architecture reviewer | read-only |
| CI/evidence/documentation coordinator | verifier/evidence/docs only |

The writer limit may increase only after measured evidence shows the M2 task/lease/scope/freshness enforcement works on real registered worker PRs without increasing stale work, merge conflicts or escaped defects.

## 3. Three-level execution model

```text
Phase
  -> Work Package
       -> Parallel Task
            -> Agent / branch / path lease / evidence
```

The roadmap cursor remains at phase/work-package level. Parallel tasks are subordinate execution units and may never mark a package complete independently.

## 4. Required task identity

Every material concurrent task records at least:

- `task_id`;
- `agent_id` or role identity;
- exact base/synchronization SHA evidence;
- branch/workspace identity;
- phase/work-package;
- owning module/domain/kernel capability;
- task purpose;
- dependencies;
- read/write/forbidden/shared paths;
- migration reservation if applicable;
- contract/event reservation if applicable;
- risk/gate expectations;
- expected handoff/merge order.

`docs/ai/AGENT_TASK_SCHEMA.json` is the standalone task-record schema. `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json` may use a compact active-wave assignment representation, but the M2 validators still enforce the effective identity, slot, authority, path, dependency and reservation invariants required for the live wave.

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

A module/task agent may not directly write another module's private source/tables/migrations, undocumented private contracts, or shared governance/runtime paths outside its assigned budget.

When such a change is required, create a dependency/change request and serialize it through the owning task.

The lease planning schema is `docs/ai/AGENT_LEASE_SCHEMA.json`.

## 6. Dependency DAG

Parallel tasks must form an acyclic dependency graph. Independent tasks may run together. A dependent task may analyze ahead but must not implement against an unaccepted guessed dependency contract.

M2 validates that dependencies exist, are acyclic, and have an earlier deterministic merge order than their dependents.

## 7. Protected-main freshness and stale-task handling

Every task starts/synchronizes from a known protected-main state, but **live protected `main` is the freshness authority**. SHA fields stored in repository plans are audit snapshots, not permanent locks.

Before material mutation, onboarding assignment, PR review and merge:

1. resolve current protected `main`;
2. compare current assumptions with the branch/task state;
3. identify changed owned/shared paths and contracts;
4. if assumptions are stale, stop the write path;
5. synchronize/re-plan/review as required;
6. obtain fresh exact-head CI evidence.

For registered worker PRs, M2 fails required `governance` when the PR base is not the current protected-main SHA or the head does not descend from that base.

A previously green run is not merge authority after a material head change.

## 8. Migration reservation

Before schema mutation, the task reserves and documents authoritative schema owner, exact migration path/version, data budget, fresh-install/upgrade path, recovery strategy and tenant/owner-isolation expectations.

M2 fails duplicate active owner/version reservations, duplicate migration paths, or a reservation whose path is outside the owning task's write budget.

## 9. Contract/event reservation

A task changing a public API, capability or event declares exact contract identity/version. Two agents may not independently redefine the same public contract. The owning contract task lands first; dependents refresh their base and consume the accepted version.

Current M2 covers the active-wave path/migration/dependency invariants. Broader machine enforcement for future public-contract registries may be extended as those registries become active; that is not a blocker for the current P04.04 wave because no competing public-contract reservation exists in the active plan.

## 10. Conflict detection

Before assignment and before PR submission compare planned write paths, active write leases, migration reservations, public-contract identities, ownership, dependency graph and live protected main.

M2 now fails closed on active write-path overlap. "Agents will probably edit different lines" is not an acceptable conflict strategy.

## 11. Merge model

```text
task branch
  -> task CI / Supervisor review
  -> dependency-aware source acceptance
  -> unchanged promotion when required
  -> promotion-specific Governance
  -> protected-main freshness / expected-head merge
  -> protected-main readback
  -> merge alert / worker synchronization
```

The Supervisor sequences accepted work but does not receive unrestricted permission to modify every module.

## 12. Agent Working Instructions synchronization

Every material task checks effective instructions at task start and again before PR submission, including phase/package, owner, role, branch/main freshness, path budget, leases, dependencies, migration/contract assumptions, tests/gates, Supervisor process, coordination channel and stop conditions.

The human-readable mirror is the README section **Agent Working Instructions**. If effective instructions change materially, update that README section in the same governed PR. If unchanged, the PR records:

`Agent instructions checked — README instruction delta: none`

README is a mirror, not authority. `AGENTS.md`, `STATE.json`, mandatory governance policy and accepted ADRs win on conflict.

## 13. Current-vs-future scaling

| Maturity | Total agents | Concurrent writers |
|---|---:|---:|
| Current M2 initial proof | 4-6 | <=3 |
| M2 proven on real worker PRs | 6-8 | 3-5, evidence-gated |
| Stable independent business modules | 8-12 | 5-8 |
| Proven automated orchestration/merge queue | 12-20+ | evidence-based |

Agent count is never a target by itself. Increase only while throughput rises without increasing stale work, conflicts, review debt or escaped defects.

## 14. Enforcement maturity

### M0 — Manual semantics — implemented

Isolated branches/workspaces, task identity, path budgets, base/main checks, dependency order, README instruction synchronization, PR declaration and exact-head CI are established.

### M1 — Machine-readable planning — implemented

Machine-readable task, lease, signal and slot contracts plus `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json` are established.

### M2 — Required CI enforcement — implemented by XQ-100 M2

The existing required `governance` job executes the M2 gate through `scripts/verify_go_quality.sh` before expensive Go tooling.

M2 implementation:

- `scripts/agent_orchestration_common.py` — shared dependency-free Git/plan/path helpers;
- `scripts/test_agent_orchestration_common.py` — focused helper behavior tests;
- `scripts/validate_agent_task.py` — canonical authority, task/slot identity and capacity checks;
- `scripts/validate_agent_leases.py` — write/forbidden-path and migration reservation checks;
- `scripts/detect_path_overlap.py` — explicit active writer overlap detector;
- `scripts/validate_task_dependencies.py` — DAG and deterministic merge-order checks;
- `scripts/validate_agent_pr_scope.py` — registered worker PR changed-path budget enforcement;
- `scripts/validate_agent_base_sha.py` — registered worker PR live-main freshness/ancestry enforcement.

Fail-closed semantics:

- unknown `agent/*` PR branch -> fail;
- active plan authority mismatch -> fail;
- active write-path overlap -> fail;
- migration reservation collision -> fail;
- dependency cycle/order error -> fail;
- registered worker PR path outside declared budget -> fail;
- registered worker PR based on stale protected main -> fail.

Control/governance branches still execute plan/slot/lease/DAG validation, but worker-specific PR scope/base enforcement applies to the registered active worker/Supervisor task branches. This allows governed coordination/tooling evolution without pretending a control PR is a worker task.

The writer cap stays at 3 until at least one real registered worker PR proves the worker-specific path/freshness gates in canonical CI.

### M3 — Developer CLI / task registry — future optimization, not a current blocker

Target commands remain `omnexa dev task create|claim|status|release|graph`. The current wave can operate safely through the machine-readable plan + required CI without this convenience layer.

### M4 — Conflict-aware automated merge orchestration — future scale optimization, not a current blocker

A future orchestrator may add automatic dependency-aware queueing, lease expiry/release and stale-task notification. The current <=3 writer wave already has deterministic Supervisor merge order, strict main protection and M2 CI enforcement.

## 15. Completion metrics

Track lead/cycle time, CI queue time, stale invalidations, merge conflicts, overlap violations, rework, escaped regressions, reviewer wait time and percentage of tasks merged without coordination repair. Scale concurrency only when evidence remains healthy.

## 16. Non-authority statement

This contract changes development coordination only. It does **not** activate P04.05+, P05+, business-feature code, product AI/agent runtime, cross-owner writes, or any authority outside `STATE.json`.
