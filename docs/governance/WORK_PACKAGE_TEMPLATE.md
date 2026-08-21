# Omnexa Work Package Template

Use this structure before implementation begins for every material work package.

---

# Pxx.xx — Work Package Name

- Phase: `Pxx`
- State: `planned | ready | active | blocked | verification | done | superseded`
- Owner/domain:
- Depends on:
- Blocks:
- Change class: `A | B | C | D`

## 1. Objective

State the outcome in one or two precise paragraphs. Describe the capability being established, not implementation activity alone.

## 2. Why this package exists

Explain the architectural/business problem and why it belongs in this phase.

## 3. In scope

- explicit deliverable;
- explicit deliverable;
- explicit deliverable.

## 4. Out of scope

- future feature intentionally excluded;
- adjacent domain intentionally excluded;
- optimization intentionally deferred.

## 5. Architecture boundary

- Owning domain/module:
- Kernel capabilities consumed:
- Capabilities provided:
- Capabilities consumed:
- Events published:
- Events consumed:
- Workflow actions/triggers:
- Data/schema owner:

## 6. Tenancy and authorization

Document:

- tenant scope;
- organization/company scope where applicable;
- permissions/capabilities;
- relationship/context rules;
- service-account behavior;
- cross-tenant negative cases.

## 7. Data model and migrations

Document:

- entities/value objects;
- ownership;
- IDs/references;
- migration sequence;
- fresh-install path;
- upgrade/backfill path;
- destructive changes, if any;
- retention/audit requirements.

## 8. API / event / contract changes

For each public contract record:

- name/version;
- request/payload;
- response;
- errors;
- compatibility policy;
- idempotency semantics;
- deprecation impact.

## 9. Failure and recovery behavior

Define:

- validation failure;
- dependency unavailable;
- retry/timeout behavior;
- duplicate/replay behavior;
- partial failure;
- rollback/compensation;
- restart/recovery.

## 10. Observability

Define required:

- logs;
- metrics;
- traces;
- audit records;
- health signals;
- correlation identifiers.

## 11. Security considerations

Document:

- trust boundaries;
- sensitive data;
- secret use;
- privileged operations;
- webhook/external-call security;
- abuse/rate risks;
- applicable threat-model items.

## 12. Acceptance criteria

Use numbered, objectively testable criteria.

1. AC-01 — ...
2. AC-02 — ...
3. AC-03 — ...

No acceptance criterion should depend on subjective wording such as "works well" or "looks correct" without measurable evidence.

## 13. Required tests

- unit;
- integration;
- contract;
- tenant isolation;
- authorization;
- migration fresh/upgrade;
- lifecycle;
- retry/idempotency;
- end-to-end where justified.

## 14. Completion evidence

Record exact commands, CI run/job IDs or test references using PASS / FAIL / BLOCKED / NOT RUN / N/A.

## 15. Risks and mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| | | |

## 16. ADRs

List any required or related ADRs.

## 17. State transition

When all acceptance criteria and Definition of Done gates pass:

- update this package state;
- reconcile `STATUS.md`;
- reconcile `STATE.json`;
- append an execution-ledger entry;
- activate only the next permitted package.
