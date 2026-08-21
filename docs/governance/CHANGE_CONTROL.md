# Omnexa Change Control

Status: **Mandatory**

## 1. Purpose

Omnexa is intended to be developed over many phases by multiple humans and AI systems. This process prevents local convenience from silently changing the platform architecture or roadmap.

## 2. Change classes

### Class A — Local implementation
Does not alter public contracts, data ownership, architecture, roadmap or external behavior.

Examples:
- internal refactor;
- test improvement;
- compatible bug fix.

Requires normal PR review and quality gates.

### Class B — Contract-compatible product change
Adds behavior inside an approved active work package without breaking existing architecture.

Examples:
- additive API field;
- new permission under an approved module;
- new workflow action with stable ownership.

Requires contract/test/document updates and PR evidence.

### Class C — Architectural change
Changes a baseline decision or cross-cutting rule.

Examples:
- technology/runtime baseline;
- tenancy hierarchy;
- authorization model;
- data ownership rule;
- module lifecycle semantics;
- service extraction strategy;
- event or API compatibility model;
- phase ordering/gates.

Requires ADR + architecture/roadmap reconciliation before implementation.

### Class D — Destructive or compatibility-breaking change
Removes or breaks data/contracts/capabilities or performs irreversible migration.

Requires Class C controls plus explicit migration, compatibility, backup and rollback/forward-fix strategy.

## 3. ADR process

For Class C/D changes, create `docs/adr/ADR-NNNN-short-title.md` with:

- Status: proposed / accepted / superseded / rejected
- Context
- Problem
- Decision
- Alternatives considered
- Consequences
- Compatibility impact
- Migration impact
- Security/tenancy impact
- Operational impact
- Rollback or forward-fix strategy
- Documents/work packages affected

An ADR is not accepted merely because the file exists. The PR containing the decision must be reviewed/merged according to repository governance.

## 4. Atomic reconciliation

When an architectural decision changes, the same PR should reconcile all affected sources of truth, including as applicable:

- `PRODUCT_CONSTITUTION.md`
- `SYSTEM_ARCHITECTURE.md`
- `MODULE_STANDARD.md`
- `MASTER_PLAN.md`
- `STATUS.md`
- `STATE.json`
- API/event/module specifications
- tests/validators

Do not leave contradictory governance documents after merge.

## 5. Roadmap amendment

To change phase scope/order:

1. state why the current plan cannot safely satisfy the new requirement;
2. identify affected dependencies;
3. identify whether already-completed work becomes invalid;
4. update phase acceptance gates;
5. update machine-readable state;
6. record the decision in an ADR when architecture is affected.

New ideas should normally enter a future phase/backlog, not the currently active package.

## 6. Emergency fixes

A production/security blocker may require work outside the active phase.

Emergency work must be:

- narrowly scoped;
- labeled as an exception;
- linked to the impacted invariant/risk;
- fully tested to the extent possible;
- followed by status/plan reconciliation;
- prevented from becoming a hidden feature expansion.

Emergency status does not permit bypassing tenant or authorization controls.

## 7. Breaking contract policy

Before breaking a public capability/event/API:

- establish a new version;
- define coexistence or migration window;
- identify consumers;
- add compatibility tests;
- publish deprecation behavior;
- define removal conditions.

A private internal refactor should not force public consumers to change.

## 8. Destructive data policy

Destructive changes require:

- explicit data owner;
- impact analysis;
- backup/restore path;
- tenant scope validation;
- dry-run or rehearsal where practical;
- audit record for runtime destructive actions;
- rollback or approved forward-fix plan.

Normal module disable is never a destructive-data operation.

## 9. Dependency introduction

A new external dependency must justify:

- capability gained;
- license compatibility;
- security/supply-chain implications;
- maintenance health;
- runtime/deployment burden;
- exit/replacement strategy for critical dependencies.

Do not add overlapping libraries/services when an approved platform capability already exists.

## 10. Decision precedence

If documents disagree, implementation must not guess.

Resolution order:

1. accepted, non-superseded ADR relevant to the question;
2. Product Constitution;
3. System Architecture / Module Standard;
4. Master Plan;
5. work-package documentation;
6. implementation details.

Any discovered contradiction should be fixed as governance debt before using it as precedent.
