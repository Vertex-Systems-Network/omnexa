# AI Execution Policy

Status: **Mandatory**  
Audience: AI coding agents, autonomous agents, human reviewers, maintainers

## 1. Purpose

This policy ensures that any AI system working on Omnexa follows the same architecture, roadmap, evidence and change-control rules. AI may accelerate implementation; it may not redefine the product implicitly.

## 2. Source of truth

The following artifacts are authoritative:

- `AGENTS.md` — repository-wide execution contract
- `docs/governance/PRODUCT_CONSTITUTION.md` — product and architecture laws
- `docs/architecture/SYSTEM_ARCHITECTURE.md` — technical boundaries
- `docs/architecture/MODULE_STANDARD.md` — module contract
- `docs/roadmap/MASTER_PLAN.md` — phase sequence and gates
- `docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md` — mandatory preplanned module/submodule decomposition and ordered execution template
- `docs/roadmap/modules/SUBMODULE_CATALOG.json` — machine-readable P02-P27 submodule/family identity
- `docs/roadmap/modules/` — owning architecture, flow and option dossiers
- `docs/governance/AI_MODULE_EXECUTION_PROTOCOL.md` — mandatory continuation/no-replanning/handoff protocol
- `docs/roadmap/STATUS.md` — human-readable progress
- `docs/roadmap/STATE.json` — machine-readable canonical execution state and sole implementation authorization
- `docs/governance/CHANGE_CONTROL.md` — architecture-change process
- `docs/governance/DEFINITION_OF_DONE.md` — completion evidence
- `docs/adr/*` — approved architectural decisions

AI systems must read the applicable artifacts before implementation. A preplanned future submodule is not executable unless `STATE.json` authorizes its owning phase/work package.

## 3. Required behavior before coding

For each request, an AI system must determine:

1. which roadmap phase/work package owns the request;
2. whether that work package is currently permitted by `STATE.json`;
3. which canonical module/submodule ID owns the behavior when P02-P27 planning applies;
4. what existing module owns the relevant data/capability;
5. which preplanned S01-S10 work-package/task owns the requested behavior;
6. what contracts/events/workflows/UI slots may be affected;
7. what tenant, authorization, audit and migration implications exist;
8. what settings/options/feature flags are affected, including scope and change permission;
9. what acceptance tests/evidence are required;
10. what the next incomplete authorized task is.

If a requested feature is outside the active scope, the AI should record/refine it in the owning future module/submodule plan rather than silently implementing it.

When a module/submodule decomposition already exists, the AI must continue from the next incomplete planned task. It must not restart a generic planning cycle merely because a new implementation session or agent begins. Replanning is appropriate only for a genuine architecture conflict, missing owner/dependency, changed requirement, security/regulatory constraint or approved change-control/ADR event.

## 4. Prohibited autonomous decisions

Without an approved ADR and plan reconciliation, AI must not:

- change the primary backend language;
- replace PostgreSQL as the primary transactional data model;
- replace the tenancy hierarchy;
- create independent authentication stacks inside modules;
- create direct cross-module database writes;
- change module lifecycle semantics;
- convert the project to microservices broadly;
- add a new business domain to the active phase;
- bypass policy checks for convenience;
- introduce a breaking event/API version silently;
- remove auditability, retry safety or rollback behavior;
- change roadmap phase order;
- mark a phase done because code compiles only;
- invent a second ownership model because an existing submodule plan is inconvenient;
- treat a preplanned future module/submodule as implementation authorization.

## 5. Allowed autonomy

Within an active work package, AI may autonomously choose implementation details when they do not violate approved contracts or architecture. Examples include:

- local function decomposition;
- internal naming consistent with repository conventions;
- test organization;
- non-breaking refactoring;
- implementation of already-approved interfaces;
- performance improvements that preserve semantics;
- documentation corrections that do not alter architecture.

## 6. Runner-deferred implementation workflow

Canonical hosted verification may be deferred until the implementation branch is substantially complete. For an already-authorized work package, AI should normally:

1. complete the planned source/tests/docs/verifier work first;
2. perform all available deterministic static/unit/self-review preparation;
3. open the executable PR only when it is ready for end-stage integration verification;
4. use the required GitHub-hosted governance lane as the authoritative final integration gate;
5. fix failures rather than weakening checks;
6. merge only after required checks are green;
7. reconcile completion evidence/state only after merge evidence exists.

Deferring the runner does **not** permit an unverified `PASS`, `done` or protected merge claim.

## 7. Evidence requirement

AI-generated work is not complete until evidence exists for all applicable gates:

- formatting/lint/type/static analysis;
- unit tests;
- integration/contract tests;
- fresh database/migration validation;
- tenant isolation tests;
- authorization tests;
- module lifecycle tests;
- build/package checks;
- security scans where configured;
- W3C/WAVE/manual accessibility evidence when authorized browser UI is affected;
- CI result.

The evidence must be summarized in the PR and reflected in status files only after it exists.

## 8. No invented success

An AI system must distinguish:

- `not run`
- `run and failed`
- `run and passed`
- `blocked by unavailable environment`

It must never report `passed`, `green`, `done` or equivalent without observable evidence.

## 9. Repository hygiene

AI must not copy files from unrelated repositories or projects unless the task explicitly requires a reviewed import. Every created file must have an Omnexa purpose.

Temporary experiments, generated artifacts, local secrets and debug output must not enter version control.

## 10. Scope control

Each work item must have explicit:

- **In scope**
- **Out of scope**
- **Dependencies**
- **Acceptance criteria**

Every future module/submodule task must also resolve:

- canonical phase/module/submodule ID;
- authoritative owner/write boundary;
- owning dossier architecture and primary flow;
- options/settings/policies affected;
- permissions/tenant/org scope;
- contract/event/workflow/UI contribution impact;
- required S01-S10 task/evidence position.

If implementation discovers a necessary dependency, classify it as:

- same-work-package prerequisite;
- blocker requiring upstream fix;
- future enhancement.

Do not expand scope invisibly.

## 11. Architecture conflict protocol

If current code conflicts with governance documents:

1. do not normalize the conflict silently;
2. document the conflict;
3. determine whether code or governance is intended to change;
4. if architecture changes, create an ADR;
5. update affected architecture/roadmap/catalog/dossier/state documents atomically with the approved change;
6. only then implement the new baseline.

## 12. Data ownership rule

An AI system must ask "who owns this write?" before adding persistence.

Examples:

- CRM owns opportunity writes.
- Finance owns ledger writes.
- Inventory owns stock movement writes.
- Payments owns payment state transitions.
- Commerce may request or react to those capabilities but may not mutate their private tables.

Cross-domain reads should use capability APIs, events, projections or explicitly approved read models.

## 13. AI feature implementation rule

AI features inside Omnexa itself must use governed capability interfaces.

Required path:

```text
AI request
  -> authenticated actor/service identity
  -> tenant/context resolution
  -> policy evaluation
  -> approved capability/tool
  -> domain validation
  -> state change
  -> audit/event
```

Direct unrestricted database mutation by an AI agent is prohibited.

## 14. Module/submodule continuation rule

For P02-P27 follow `docs/governance/AI_MODULE_EXECUTION_PROTOCOL.md`:

```text
request
 -> phase
 -> module/submodule ID from SUBMODULE_CATALOG
 -> STATE.json authorization
 -> owning dossier architecture/flows/options
 -> S01-S10 task sequence
 -> next incomplete authorized task
 -> implementation/evidence
```

The end-of-session branch/document state must make the active submodule/task, completed work, next task, blockers and evidence status recoverable without relying on chat memory.

## 15. Status update rule

`STATUS.md` is explanatory. `STATE.json` is canonical machine state.

Any transition to `verification` or `done` must include evidence references. If the two files disagree, the safest lower-completion state wins until reconciled.

## 16. Stop conditions

An AI system must stop implementation and surface the issue when:

- required governance documents are missing or contradictory;
- a requested change requires an unapproved architectural decision;
- a destructive migration has no approved migration/rollback strategy;
- tenant isolation cannot be demonstrated;
- a required dependency would violate module boundaries;
- the submodule owner/write boundary cannot be resolved;
- completion cannot be verified.

Stopping means refusing to invent progress; it does not mean abandoning analysis or a safe planning fix.
