# Omnexa Supervisor Multi-Agent Repository Workflow

Status: **Mandatory development-governance extension**  
Parent contract: `docs/governance/MULTI_AGENT_ORCHESTRATION.md`  
Scope: repository development agents only; this is not P20 product-agent runtime

## 1. Purpose

This workflow defines how a Supervisor coordinates multiple development agents across isolated Git branches while preserving Omnexa's canonical roadmap cursor, module ownership, path leases, exact-head CI evidence, protected-main integration and stale-work safety.

The Supervisor coordinates and reviews. The Supervisor does **not** gain authority to bypass `AGENTS.md`, `docs/roadmap/STATE.json`, branch protection, required Governance, promotion rules, module/data ownership or accepted ADRs.

## 2. Roles

### Supervisor

The agent operating the main repository/integration workflow acts as **Supervisor**.

The Supervisor is responsible for:

- establishing the parallel-work wave from a verified protected-main SHA;
- creating the worker branches before task work starts;
- maintaining the AI-Native active-wave plan;
- assigning logical agent/task identities and leases;
- monitoring the coordination channel for completion signals;
- pausing its own task when a submission requires review;
- reviewing submitted diffs, tests, CI, leases, dependency freshness and instructions;
- sequencing accepted merges according to dependencies;
- using the repository's governed source/promotion path before protected-main merge;
- reading back protected main after merge;
- broadcasting the required merge alert;
- requiring every active worker to synchronize the new protected-main SHA before resuming;
- resuming its own paused task only after the review/merge action reaches a stable state.

### Worker agent

A worker agent owns exactly one bounded task/branch at a time unless the active plan explicitly records otherwise. A worker may read broadly but writes only inside its declared budget and active leases.

### Reviewer / security / architecture agent

A reviewer may inspect multiple branches but is read-only unless a separately declared task/lease grants a narrow write budget.

## 3. Supervisor bootstrap gate — first repository mutation

For every new parallel development wave:

1. mandatory repository authority readback happens first: protected `main`, `STATE.json`, `AGENTS.md`, applicable work-package contract and current orchestration policy;
2. determine which tasks/modules are **currently authorized** to run in parallel;
3. the **first repository mutation** after that mandatory readback is to create a separate branch for every worker task/module and a separate branch for the Supervisor's own task;
4. every branch must be created from the same recorded protected-main `base_sha` unless the plan explicitly records a dependency requiring a later base;
5. verify every created branch ref;
6. only after branch bootstrap is complete may the Supervisor write the active AI-Native plan, assign tasks, or start its own implementation work.

The phrase “first action” never authorizes skipping mandatory canonical-state readback. It means **branch creation is the first mutation**.

The Supervisor must not create implementation branches for future locked phases/modules merely to appear parallel. Only currently authorized work is eligible.

## 4. AI-Native active-wave plan

Every live wave must have a machine-readable plan. Current location:

`docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`

At minimum it records:

- wave ID;
- phase/work package;
- protected-main bootstrap SHA;
- current required main SHA;
- coordination channel;
- Supervisor identity and branch;
- every worker agent/task/module branch;
- exact read/write/forbidden paths;
- active/exclusive leases;
- migration or contract reservations;
- dependencies;
- merge order;
- task state;
- last synchronized main SHA;
- completion signal requirements;
- merge-alert requirements.

The active plan is coordination metadata. It cannot grant work outside canonical `STATE.json` authority.

## 5. Required task completion signal

When any worker or the Supervisor finishes the task it intends to submit, it must announce exactly:

`Work Done and Submitted`

The signal must also carry:

- `task_id`;
- agent/role identity;
- branch;
- exact submitted head SHA;
- PR number or submission reference;
- test/CI state using honest PASS/FAIL/BLOCKED/NOT RUN semantics;
- latest synchronized protected-main SHA;
- instruction-check result;
- lease/dependency notes;
- whether the task is safe for Supervisor review now.

A phrase without the required metadata is not a valid submission signal.

The machine-readable signal shape is `docs/ai/AGENT_SIGNAL_SCHEMA.json`.

## 6. Supervisor interrupt handling

When the Supervisor receives a valid `Work Done and Submitted` signal:

1. checkpoint the Supervisor's own branch/head and mark its current task `paused-for-review` in coordination state;
2. stop making new material changes on the Supervisor task until the submitted task has been triaged;
3. read current protected main and the submitted exact head;
4. verify the worker branch matches its task/lease/path budget;
5. review code, tests, security/tenancy implications, migration/contract reservations and instruction synchronization;
6. verify dependency and main-SHA freshness;
7. if changes are required, publish `Review Changes Required` with concrete findings and resume the Supervisor's own task after the review state is recorded;
8. if approved, enter the governed merge pipeline;
9. after merge/readback or a recorded non-merge outcome, resume the Supervisor's own paused task from its checkpoint, after synchronizing it if main moved.

“Pause” is a workflow state, not a promise that GitHub can suspend an arbitrary external AI process. Agent runtimes must honor this protocol.

## 7. Governed merge pipeline

Supervisor approval never means a direct unverified push to `main`.

For an approved submission:

```text
worker exact head
  -> source PR review
  -> exact-head Governance
  -> unchanged promotion carrier when required
  -> promotion-specific Governance
  -> zero unresolved review threads
  -> protected-main freshness / expected-head check
  -> guarded merge
  -> protected-main readback
```

If current repository governance evolves, the active accepted rule wins and this document/README must be synchronized.

The Supervisor may perform the final merge operation but may not waive required checks.

## 8. Merge alert to all active agents

After a worker/supervisor submission is accepted into protected main and read back, the Supervisor must publish exactly:

`New changes have been merged — please merge these changes into your branch first, then resume your own work.`

The alert must include:

- new protected-main SHA;
- merged task/PR;
- affected dependencies/contracts/migrations if known;
- which active tasks are definitely stale;
- whether any lease changed/released.

The alert is published to the active coordination channel and recorded in the active plan's required-main SHA.

## 9. Worker response to merge alert

Every active worker receiving or discovering a newer required-main SHA must, **before any further material mutation**:

1. stop/hold its current write path;
2. fetch/read the new protected main;
3. merge/pull the accepted main changes into its own branch using a non-destructive repository-approved method;
4. resolve conflicts without discarding another task's accepted behavior;
5. re-read effective instructions, leases, dependencies, migrations/contracts and path budget;
6. update `last_synced_main_sha` to the new protected-main SHA;
7. invalidate stale tests/review/CI evidence and rerun what is required;
8. announce `Sync Complete — Resuming Work` with task/branch/synced SHA;
9. only then continue/resume its assigned task.

An agent that missed the alert is still protected by the mandatory SHA check: the active plan's `required_main_sha` is authoritative coordination state for the wave.

## 10. Supervisor's own module/task

The Supervisor must also use an isolated task branch and declared path lease for its own work. It does not work directly on protected main.

When another submission interrupts the Supervisor:

- Supervisor task state becomes `paused-for-review`;
- its head/checkpoint is recorded;
- after the incoming review/merge is complete, the Supervisor synchronizes new main if required;
- its previous CI/review is treated as stale when its head changes;
- it resumes only within its original scope unless re-planned.

The Supervisor must not use its review role to silently broaden its own write authority. Where independent review is configured/available, the Supervisor's own material/high-risk branch should receive independent review rather than being self-approved solely by the same authority path.

## 11. Multiple simultaneous submissions

If multiple valid completion signals arrive while the Supervisor is busy, place them in a review queue.

Priority order:

1. security/correctness blocker that affects other tasks;
2. dependency-unblocking submission;
3. migration/public-contract owner submission needed by dependents;
4. otherwise oldest valid submitted signal.

The Supervisor processes one shared-surface merge at a time. Independent review can occur concurrently, but protected-main merge ordering remains deterministic.

## 12. Stale submission protocol

A submitted branch is stale when protected main has moved in a way that can affect its assumptions, paths, dependencies, contracts or migration namespace.

Before merge, Supervisor must require the worker to synchronize/revalidate if needed. A previously green CI run on a materially stale head is not merge authority.

The worker then posts a new `Work Done and Submitted` signal for its new exact head.

## 13. Review rejection / conflict handling

If review fails:

- do not merge;
- publish `Review Changes Required`;
- identify exact findings and whether leases/dependencies change;
- worker fixes only within the accepted task budget or requests re-plan;
- Supervisor resumes its own task unless another valid submission is queued.

If branch sync creates conflicts on an exclusive/shared authoritative surface, the owning task/Supervisor resolves the merge order; agents must not independently choose conflicting resolutions.

## 14. Branch and lease lifecycle

After a task is merged and protected main is read back:

- release its exclusive leases;
- mark the task merged/done only when required evidence exists;
- notify dependents;
- keep or delete the branch according to repository retention policy;
- never reuse a completed branch as implicit authority for unrelated work.

## 15. Coordination channel

The current P04.04 wave uses GitHub issue `#177` as its persistent coordination channel.

Issue comments/signals are **coordination data**, not implementation authority. They cannot override canonical repository documents or CI evidence.

A future external orchestrator may replace polling with webhooks/messages, but it must preserve the exact task/SHA/lease/signal semantics defined here.

## 16. Current-wave scope constraint

At the current canonical checkpoint only P04.04 is active. Therefore this workflow may split P04.04 into bounded parallel tasks, but it may not activate P04.05+, P05+, CRM, Finance, Inventory, HR, Commerce or other future business modules.

Broad independent module branching becomes legal only when canonical roadmap gates explicitly open those streams.
