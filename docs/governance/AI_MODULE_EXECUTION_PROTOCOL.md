# AI Module & Submodule Execution Protocol

Status: **Mandatory**
Audience: AI coding agents, human implementers and reviewers

## Purpose

This protocol prevents repeated redesign of Omnexa modules during implementation. The roadmap module/submodule architecture is preplanned before activation; an implementation agent should resolve the active authorized submodule and execute its next incomplete task instead of generating a new generic plan every session.

## Canonical planning sources

For P02-P27 work, use these sources in order:

1. `docs/roadmap/STATE.json` — sole implementation authorization;
2. active work-package specification — executable scope and acceptance gates;
3. `docs/roadmap/modules/SUBMODULE_CATALOG.json` — machine-readable phase/submodule identity;
4. owning dossier under `docs/roadmap/modules/` — architecture, flows, options and boundaries;
5. `docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md` — standard ordered S01-S10 task decomposition;
6. `docs/architecture/MODULE_STANDARD.md` and ownership/dependency/API/event/security standards;
7. applicable UI/accessibility, developer-command and quality documents.

A planning source never overrides `STATE.json`. A future submodule can be fully documented and still remain non-executable.

## Required AI resolution algorithm

For each requested feature/change:

```text
request
 -> identify owning phase
 -> identify owning module/submodule ID from SUBMODULE_CATALOG
 -> verify STATE.json authorization
 -> load owning dossier
 -> locate the relevant flow/options/contracts
 -> map work to S01-S10 task sequence
 -> select next incomplete authorized task
 -> implement + test + document
 -> update explicit task/evidence state
 -> defer hosted runner until implementation-complete when permitted
 -> final hosted verification -> exact-head merge -> closure
```

If the request belongs to an inactive phase, record/refine planning only. Do not create executable business/module code.

## No-replanning rule

Do not restart architecture planning merely because:

- a new chat/session starts;
- a different AI agent takes over;
- the implementation branch contains many files;
- the agent would personally choose another framework or folder split;
- a submodule such as Page Builder, Template Builder, Form Builder, Dashboard Builder, Checkout, General Ledger or Agent Runtime is large.

Replanning is allowed only when there is a concrete reason:

- current requirement materially conflicts with the canonical dossier;
- no owner exists for a required write model/capability;
- a required dependency is missing or forbidden;
- a public contract/phase boundary must change;
- new regulatory/security evidence invalidates the existing design;
- approved change control/ADR explicitly changes the baseline.

When replanning is justified, update the canonical dossier/catalog atomically before implementation; do not keep a private alternate plan only in chat.

## Required submodule task card

Before executable coding, resolve a task card containing:

- phase/module/submodule ID;
- task ID/order under S01-S10;
- owner/write boundary;
- dependencies and forbidden dependencies;
- in-scope/out-of-scope;
- affected capabilities/APIs/events/workflows/UI slots;
- persistence/migration implications;
- permissions/tenant/org scope;
- settings/options changed;
- security/data-classification implications;
- tests/evidence required;
- completion state: `planned|active|blocked|verification|done`.

The card can live in the activated work-package plan/checklist; it need not create one file per tiny task.

## Options/configuration rule

Every configurable option added by a submodule must specify:

- stable key/name;
- type and allowed values/bounds;
- default;
- scope (`platform`, `tenant`, `organization`, `module`, `user`, or explicit domain scope);
- who may read/change it;
- sensitivity/classification;
- whether changing it is audited;
- effective timing (immediate, next request, restart, migration, scheduled effective date);
- compatibility/rollback behavior;
- whether it is a setting, policy or feature flag.

Security, accounting, tenancy and ownership invariants must not be disguised as ordinary user-editable settings.

## Flow rule

Every non-trivial submodule must document at least:

- happy path;
- validation failure;
- authorization/tenant denial when applicable;
- dependency/provider unavailable path;
- cancellation/timeout;
- duplicate/retry/idempotency behavior where applicable;
- disable/upgrade/restart behavior;
- rollback/compensation where state can partially commit.

Builder submodules also document author -> validate -> preview -> version -> publish/activate -> restore/degrade flows.

## UI rule

When an activated task includes browser UI, the task inherits `docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md`: semantic HTML, WCAG 2.2 AA target, keyboard/focus behavior, W3C validation, WAVE evaluation and required manual evidence. Automated accessibility output alone cannot be called full compliance.

## Handoff rule

At the end of each implementation session, leave the branch/doc state so another agent can determine without chat history:

- active phase/submodule/task;
- completed task IDs;
- next incomplete task;
- known blockers;
- exact tests already run and their state;
- whether hosted verification is still pending;
- files/contracts changed.

Do not leave critical planning only in conversational memory.

## Final verification rule

Hosted runner consumption may be deferred until source/tests/docs/verifier work is implementation-complete. That does not change evidence semantics: no `PASS`, `done`, protected merge or phase transition without the required canonical hosted evidence.