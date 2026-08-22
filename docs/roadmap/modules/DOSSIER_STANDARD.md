# Module/Submodule Dossier Standard

Status: **Mandatory planning format**

Use this format whenever an activated phase refines a catalog submodule into executable work packages. The program dossiers already provide the baseline; activation may add detail but must not silently change ownership or architecture.

## 1. Header

Every executable submodule plan records:

```text
Phase: Pxx
Module: <domain/module>
Submodule: Pxx.X — <name>
State: planned|ready|active|blocked|verification|done
Owner/write boundary: <authoritative owner>
Depends on: <platform/required/optional dependencies>
Dossier baseline: <path>
Implementation authority: <STATE.json work package>
```

## 2. Purpose and outcomes

Document:

- user/business/platform problem solved;
- measurable outcome;
- what remains explicitly outside this submodule;
- why this capability belongs to this owner rather than another module.

## 3. Architecture

Required architecture sections:

### 3.1 Ownership

- authoritative write model;
- read projections/references;
- prohibited foreign writes/imports;
- immutable snapshots when cross-domain history requires them.

### 3.2 Components/layers

Use a text diagram such as:

```text
transport/UI
 -> application/use case
 -> domain invariants
 -> repository/provider boundary
 -> owned persistence or governed external capability
 -> events/audit/projections
```

Builder modules use:

```text
definition/schema
 -> validator/versioning
 -> runtime renderer/executor
 -> registry
 -> authoring API
 -> visual editor
 -> preview/simulation
 -> publish/activate lifecycle
```

### 3.3 Dependencies

Classify each as:

- platform;
- required module;
- optional module;
- external provider;
- forbidden dependency.

For optional dependencies, document exact degraded behavior.

### 3.4 Data/persistence

- entities/value objects;
- owned tables/schema;
- indexes/uniqueness;
- migration sequence;
- fresh-install/upgrade/backfill;
- retention/archive/purge;
- data classification.

## 4. Primary flows

Every plan documents flow IDs using this minimum set:

- `F01 Happy path`;
- `F02 Validation rejection`;
- `F03 Authorization/tenant denial` when applicable;
- `F04 Missing/disabled dependency`;
- `F05 Provider/service unavailable`;
- `F06 Timeout/cancellation`;
- `F07 Duplicate/retry/idempotency` when applicable;
- `F08 Upgrade/restart/recovery`;
- `F09 Rollback/compensation` when partial state can occur.

Represent important flows explicitly:

```text
actor/request
 -> context + authorization
 -> validation
 -> owning capability
 -> domain mutation/read
 -> persistence/provider
 -> event/audit
 -> response/projection
```

A flow must identify where failure becomes a structured public error and which steps are retryable.

## 5. Options/settings/policies

Every option row uses:

| Key | Type | Default | Scope | Change permission | Sensitive | Audit | Effective | Validation/rollback |
|---|---|---|---|---|---|---|---|---|

Rules:

- secret != feature flag;
- permission/invariant != ordinary setting;
- accounting/tenancy/security ownership constraints cannot be disabled by a convenience toggle;
- changing a public behavior option requires compatibility analysis;
- tenant options cannot bypass platform hard limits.

## 6. Contracts

List as applicable:

- capabilities provided/consumed;
- HTTP/OpenAPI endpoints;
- events produced/consumed;
- workflow triggers/actions;
- UI slots/navigation/pages/widgets/blocks;
- import/export formats;
- provider adapter contracts;
- search/reporting projections.

Every public contract records version, permission, scope, validation/failure semantics and compatibility policy.

## 7. Permissions and security

Document:

- stable permission names;
- server-side enforcement boundary;
- tenant/org/resource scope;
- privileged actions/approvals;
- data classes;
- secrets;
- upload/download/network risks;
- abuse/rate limits;
- audit events.

## 8. UI/UX plan

Only when UI is authorized:

- route/navigation/slot;
- information hierarchy;
- responsive behavior;
- loading/empty/error/success/disabled states;
- keyboard/focus/drag-drop alternatives;
- localization/RTL;
- reduced motion;
- W3C/WAVE/WCAG 2.2 AA/manual accessibility evidence;
- performance budgets.

## 9. Ordered work packages

Map implementation to the shared sequence:

- `S01` ownership/contracts;
- `S02` data/migrations;
- `S03` domain/application behavior;
- `S04` provider/integration adapters;
- `S05` permissions/tenancy/security;
- `S06` UI/UX contributions;
- `S07` lifecycle/resilience;
- `S08` tests/quality;
- `S09` docs/AI handoff;
- `S10` final hosted verification/closure.

Each S-step may be split into `Sxx.T01`, `Sxx.T02`, etc. Every task has owner, dependency, acceptance and evidence.

## 10. Task card template

```markdown
### S03.T02 — <task title>

State: planned|active|blocked|verification|done

In scope:
- ...

Out of scope:
- ...

Depends on:
- ...

Changes:
- contracts/data/options/UI affected

Acceptance:
- ...

Evidence:
- unit/integration/security/accessibility/build/CI identifiers

Handoff:
- next task / blocker / exact pending verification
```

## 11. Definition of done

A submodule is not done until:

- all mandatory task cards are done;
- ownership/contracts/options/flows match the canonical dossier;
- required migrations/lifecycle paths pass;
- permission/tenant/security negatives pass;
- UI accessibility gates pass when applicable;
- repository code quality and package gates pass;
- canonical hosted evidence exists;
- immutable completion evidence and state/ledger reconciliation are complete.

## 12. Change rule

If implementation proves the dossier wrong, do not hide the divergence in code. Record the conflict, use change control/ADR when required, update catalog/dossier/contracts, then continue from the reconciled baseline.