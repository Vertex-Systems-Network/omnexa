# ADR-0014 — Continuous AI-Native Development Invocation

Status: **Accepted**

> This ADR becomes authoritative only after the change-control carrier that contains it passes exact-head Governance, any required unchanged promotion, guarded protected merge and protected-main readback.

## Context

Omnexa already authorizes AI systems to perform bounded development work under canonical roadmap, security, exact-head CI, review, path-lease and protected-main rules.

Several interaction-layer instructions nevertheless converted ordinary engineering checkpoints into mandatory user-turn boundaries. Examples included one-turn/one-milestone wording, `REPORT -> STOP`, mandatory next-action choices after every response, pending-CI stop behavior and worker-slot wording that could be interpreted as stopping the whole Supervisor invocation.

Issue #325 records the owner requirement that AI-Native repository development continue autonomously through safe authorized work without repeatedly asking for routine engineering confirmation.

## Problem

Interaction-driven stopping causes avoidable development stalls even when the next action is already authorized and technically routine. It also conflates four different concepts:

1. stopping one blocked implementation path;
2. denying one worker-slot admission;
3. checkpointing a completed milestone;
4. ending the entire AI development invocation.

Those are not equivalent.

## Decision

Omnexa development uses a **continuous authorized invocation** model.

The execution loop is:

```text
READ
  -> RECONCILE
  -> SELECT HIGHEST-PRIORITY SAFE AUTHORIZED TASK
  -> IMPLEMENT
  -> TEST
  -> REVIEW / REPAIR / RE-PLAN
  -> INTEGRATE WHEN REQUIRED GATES PASS
  -> READ BACK
  -> CHECKPOINT / REPORT
  -> SELECT NEXT SAFE AUTHORIZED TASK
  -> repeat
```

A logical milestone is a durable checkpoint, not a mandatory end-of-turn boundary.

Within standing repository authority, the AI autonomously performs routine engineering actions such as:

- diagnose/fix/retry failed tests or CI;
- resolve merge conflicts and stale-base assumptions;
- apply review findings;
- re-plan the same-scope implementation strategy;
- maintain task Issues/PRs required by the governed workflow;
- merge through the protected path when exact-head required gates pass;
- checkpoint a blocked task and continue another non-conflicting authorized task.

The AI must not ask broad `continue?`, `retry?`, `fix?` or `merge?` questions for those already-authorized actions.

A blocked path does not end the whole invocation when another safe authorized task exists.

No open worker slot denies assignment to that arriving worker only. It does not terminate Supervisor review, integration, planning, preflight or other authorized work.

Pending CI is recorded truthfully and is not tight-polled. If another non-conflicting authorized task exists, the Supervisor continues that work.

The invocation exits only when:

1. requested/canonical authorized scope is complete;
2. the current tool/runtime/token budget is exhausted;
3. a genuine external/owner decision is the sole remaining path and no other safe authorized work exists;
4. safety/governance prohibits all remaining work.

## Owner / external decision boundary

Owner input is required only for a decision existing repository authority cannot supply, such as:

- unavailable credentials, permissions or production access;
- legally meaningful licensing/IP/trademark decisions;
- irreversible destructive approval;
- provider/business decisions not already fixed by canonical governance;
- mandatory external/human approval actually enforced by policy/platform.

The request must identify the exact missing input rather than asking for generic permission to continue.

## Authority preservation

Continuous invocation does **not**:

- bypass branch protection, exact-head Governance, security or review requirements;
- infer future work-package or phase authority;
- infer migration, provider, production, destructive or business-feature authority;
- permit a worker to write without a valid task/path lease;
- weaken tests, security controls or evidence requirements;
- change product runtime, roadmap package state or domain ownership.

## Interaction handoff

`.ai/NEXT-ACTION-OPTIONS.md` is not a mandatory gate after every development checkpoint.

Options remain appropriate for:

- URL-only/read-only repository entry;
- an explicit user request for choices;
- a genuine owner/external decision boundary when no autonomous safe work remains.

During an active authorized run, the next safe task is selected automatically.

## Compatibility impact

No product API, runtime contract, database schema, roadmap package state or business behavior changes.

This changes only repository-development orchestration and interaction semantics.

## Migration impact

None.

## Security / tenancy impact

Security and tenancy controls are unchanged. The policy remains fail-closed at actual authority boundaries.

## Operational impact

AI development sessions can progress through multiple dependency-related checkpoints without routine user confirmation. External waits and blocked surfaces remain durably recorded.

## Rollback / forward-fix

Before protected merge, close/revert the candidate. After protected merge, a later governed change may supersede this ADR. Do not restore interaction-driven stopping by weakening repository safety controls.

## Documents affected

- `AGENTS.md`
- `docs/ai/AI_EXECUTION_PROTOCOL.md`
- `docs/governance/AI_EXECUTION_POLICY.md`
- `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`
- `.ai/NEXT-ACTION-OPTIONS.md`
- `README.md`

Change-control Issue: #325
