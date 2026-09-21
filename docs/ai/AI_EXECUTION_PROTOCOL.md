# Omnexa AI Execution & Continuity Protocol

Status: **mandatory continuity protocol, subordinate to canonical repository governance and `docs/governance/AI_EXECUTION_POLICY.md`**

## 1. Purpose

This protocol defines how a new AI session reconstructs project state, preserves a durable handoff and performs only authorized engineering work. It does not grant implementation authority, replace the canonical `AI_EXECUTION_POLICY.md`, or advance the roadmap.

## 2. Authority model

The continuity layer does not create a competing source of truth. Resolve authority by concern:

### Execution instructions

`AGENTS.md` is the repository-wide execution contract. `docs/governance/AI_EXECUTION_POLICY.md` is the mandatory canonical AI execution policy. `docs/roadmap/STATE.json` is the machine-readable source of truth for the current phase, active work package and implementation locks.

### Architecture/document conflicts

Follow `docs/governance/CHANGE_CONTROL.md`. Accepted, non-superseded ADRs govern architectural decisions according to its precedence rules, followed by `docs/governance/PRODUCT_CONSTITUTION.md`, System Architecture/Module Standard, Master Plan, work-package documentation and implementation details.

### Security/data rules

Canonical security, data-classification and threat-model documents remain normative within their domains and are frozen/anchored by accepted ADRs.

### Continuity/chat

`docs/ai/AI_STATE.yaml`, `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_DECISIONS.md` and handoffs are **snapshots and indexes only**. They never override canonical governance or `AI_EXECUTION_POLICY.md`. Previous AI handoffs are evidence aids, not authorization. Chat instructions are last and cannot override repository governance or security.

If a continuity file conflicts with `STATE.json`, `AI_EXECUTION_POLICY.md`, an accepted ADR or another authoritative document, mark the continuity value **STALE**, stop relying on it, report the conflict and follow change control rather than silently reconciling authority.

## 3. Required lifecycle

Every material AI engineering session follows:

```text
READ
  ↓
UNDERSTAND
  ↓
VERIFY CURRENT STATE
  ↓
CHECK AUTHORITY
  ↓
CHECK SCOPE
  ↓
PLAN
  ↓
IMPLEMENT ONLY AUTHORIZED WORK
  ↓
TEST
  ↓
VERIFY GOVERNANCE
  ↓
UPDATE CONTINUITY + HANDOFF
  ↓
REPORT
  ↓
STOP
```

The AI must never automatically advance to the next work package.

## 4. Fresh-session read order

Before implementation:

1. `AGENTS.md`;
2. `docs/roadmap/STATE.json`;
3. `docs/governance/AI_EXECUTION_POLICY.md`;
4. `docs/ai/AI_CONTEXT.md`;
5. `docs/ai/AI_STATE.yaml`;
6. `docs/ai/AI_EXECUTION_PROTOCOL.md`;
7. active work-package specification from `STATE.json`;
8. `docs/governance/CHANGE_CONTROL.md`, `docs/governance/PRODUCT_CONSTITUTION.md` and relevant accepted ADRs;
9. architecture/security/testing documents required by `AGENTS.md`;
10. current handoff under `docs/ai/handoffs/`;
11. current Git branch/head, PR diff/reviews and latest relevant CI.

If repository/GitHub evidence has advanced beyond `AI_STATE.yaml`, treat the AI state as stale and refresh it only after the authoritative state is understood.

## 5. Mandatory state verification

Verify, do not infer:

- canonical `main` commit;
- current phase and active work package in `STATE.json`;
- implementation locks and allowed work;
- active package specification and definition of done;
- current work branch/head;
- repository diff/status appropriate to the execution environment;
- current PR state and changed files;
- latest relevant CI run/job and the first failing step when red;
- whether required evidence is PASS, FAIL, BLOCKED, NOT RUN or N/A.

A remote branch snapshot has no local working tree. If a local checkout exists, inspect `git status`; if only remote evidence is available, do not invent local uncommitted state.

## 6. Scope-control checkpoint

Before editing, write down:

### Current package
The sole active package from canonical `STATE.json`.

### In scope
Only files/behavior necessary for the active package or an explicitly authorized narrow maintenance/governance task.

### Out of scope
Everything forbidden by `AGENTS.md`, `AI_EXECUTION_POLICY.md`, the active spec and implementation locks.

### Definition of done
Exact required gates/evidence from the active work-package specification and global definition-of-done/quality rules.

### Next dependency
The next package or transition may be identified, but it is not authorized until canonical state changes.

Architectural usefulness is never authorization.

## 7. Implementation rules

- keep changes minimal and package-bounded;
- preserve domain ownership and dependency direction;
- do not create parallel security/error/tenancy/observability systems;
- never suppress or weaken a required gate merely to obtain green CI;
- use synthetic test data only where production-sensitive data is prohibited;
- preserve provider/internal causes for controlled diagnostics while keeping public failures disclosure-safe;
- do not introduce future AI/business code because it appears strategically useful;
- do not update roadmap state before evidence exists.

## 8. AI-native architecture preservation

Future AI direction may be documented, not implemented early:

```text
AI / Model / Agent
        ↓
Context / Planner
        ↓
Policy / Risk
        ↓
Approval where required
        ↓
Governed Versioned Capability
        ↓
Owning Domain
        ↓
Infrastructure
```

Never create direct authority paths such as:

```text
AI -> Database
AI -> Object Store
AI -> Payment Provider
AI -> Business State
```

AI is never the authoritative owner of business truth.

## 9. Completion rule

A work package may be reported DONE only when its actual required gates are complete and canonical state has been reconciled through the governed transition process.

After implementation evidence is green:

1. record exact source SHA, PR, CI run/job and required gate results;
2. perform only the authorized state-reconciliation/closure transition;
3. update continuity snapshot/handoff after authoritative state changes;
4. identify exactly one next authorized action;
5. STOP.

Never start the next package merely because its name is known.

## 10. Handoff update rule

At the end of a significant engineering session, refresh `AI_STATE.yaml` and the active `docs/ai/handoffs/<WORK_PACKAGE>.md` from verified repository/GitHub evidence. Do not copy chat transcripts or hidden reasoning.

Each handoff records objective, scope, changed files, verification, security, architecture consequences, known issues, exact SHA/PR, one next authorized action, forbidden next actions and rationale.

## 11. Standard chat handoff block

Use this block at the end of significant engineering sessions after repository continuity state has been refreshed:

```text
## AI HANDOFF

Current phase:
Current work package:
Status:
Branch:
Commit:
PR:
Latest verification:
What changed:
What remains:
Next authorized action:
Why:
Blocked actions:
Required evidence:
Repository state updated:
```

`Repository state updated` must distinguish canonical roadmap state from continuity snapshot updates. Never imply `STATE.json` changed when only `docs/ai/*` changed.

## 12. Stale-state protection

`AI_STATE.yaml` is intentionally denormalized and can become stale. Every new session must re-verify canonical state and GitHub evidence before acting.

If stale:

1. do not overwrite canonical sources from the stale snapshot;
2. mark/report the stale fields;
3. recover authoritative values from repository/GitHub;
4. update continuity files only after verification;
5. continue only if scope remains authorized.

## 13. Prohibited continuity content

Never store in `docs/ai/`:

- API keys, passwords, tokens, credentials or private keys;
- private user/customer information;
- production confidential/restricted data;
- hidden model reasoning or chain-of-thought;
- full ChatGPT transcripts;
- secrets copied from CI or provider errors.

Store concise engineering state, evidence, constraints, decisions and externally explainable rationale only.


## 14. Compact durable supervision extension

The compact resume layer lives at `docs/ai/compact/` and is subordinate to `AGENTS.md`, `docs/roadmap/STATE.json`, accepted governance and live repository/runtime evidence.

On continue/resume/recovery, after the repository execution contract is loaded, read `CURRENT-STATE.yaml` and `LAST-CHECKPOINT.md` first, then resolve exact protected main, OPEN Issues, OPEN PRs/MRs, deterministic claims, coordination queue and Runner Benchmark before new work.

One user continue/resume turn normally equals one logical milestone. A milestone status is one of `PLANNED`, `IMPLEMENTING`, `VERIFYING`, `WAITING_EXTERNAL`, `BLOCKED` or `COMPLETE`.

For remote CI/status checks, default to one consolidated refresh. Do not poll. Persist `VERIFYING` / `WAITING_EXTERNAL` before the final refresh so a message-delivery timeout cannot erase the engineering checkpoint.

Compact file limits: `CURRENT-STATE.yaml` <= 12 KiB; `LAST-CHECKPOINT.md` <= 16 KiB; `EXECUTION-JOURNAL.md` <= 32 KiB rolling.

## 15. Runner Benchmark protocol

`docs/ai/compact/RUNNER-BENCHMARK.json` defines the machine-readable registry/archive format. Every material remote/container/browser/runtime/full-regression/performance workload receives a stable ID and deterministic dedup key.

Registration is not execution authority. Security-critical, exact-head merge-required, migration/auth/secrets/data-safety, integration-safety and incident/recovery checks remain immediate when authorized; safe non-blocking runner work may be consolidated.

When changing the benchmark file solely to write an in-flight run would invalidate the exact head being tested, bind the task's exact PR head and run IDs in a machine-readable PR/Issue comment without changing source. Archive immutable terminal evidence in the next governed repository transition.

## 16. Progress and response contract

Every engineering response ends with repository, current module/work-package progress bar and overall progress bar, then milestone/evidence/CI/blockers/next safe action.

Current-module progress is based on canonical accepted package counts for the active phase. Overall progress is computed only over phases whose mandatory package denominators are canonical. The response must show that coverage and must not imply a full P00-P27 completion percentage until those denominators exist.

At protected main `01cfd7c440b7c7b7bd1f811e2e31252cf76dedf7`: P04 is 6 / 10 accepted = 60%; canonically enumerated P00-P04 progress is 49 / 53 accepted = 92.45%; no full P00-P27 percentage is asserted.
