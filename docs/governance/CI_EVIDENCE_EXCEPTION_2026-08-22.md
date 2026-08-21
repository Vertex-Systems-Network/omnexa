# Temporary P00 CI Evidence Exception — 2026-08-22

Status: **Approved, temporary, narrowly scoped**  
Owner authorization: **explicitly approved in the project conversation on 2026-08-22**  
Applies to: **P00 architecture/governance/specification-only work**  
Does not apply to: **P01+ runtime/kernel/business implementation, migrations, production releases, security-sensitive executable changes**

## Reason

The repository's GitHub Actions allowance is exhausted/disabled for the current billing period. PR #13 demonstrated that GitHub Actions jobs fail before runner steps execute. Issue #14 records the infrastructure evidence.

The project owner explicitly instructed the project to skip GitHub Actions while this quota condition exists and continue the remaining P00 specification work.

This exception does **not** reinterpret a failed executable test as a pass. Hosted CI is recorded as `BLOCKED`/`NOT RUN` because the runner cannot execute.

## Temporary evidence model

For P00 documentation/specification changes only, a package may progress while hosted Actions are unavailable when all of the following are true:

1. The change is strictly architecture, governance, contract-schema or documentation work permitted by the P00 implementation lock.
2. Repository diff is inspected and contains no kernel/runtime/business implementation.
3. Mandatory package evidence files exist on the feature branch.
4. `STATE.json` dependency ordering, done-count, exactly-one-active-package rule and implementation lock are manually verified.
5. `STATUS.md` agrees with `STATE.json`.
6. Machine-readable JSON/JSON Schema/YAML artifacts are structurally reviewed; any executable validation unavailable due the quota is recorded as `BLOCKED`/`NOT RUN`, never `PASS`.
7. Relevant ADR and canonical entrypoints are reconciled.
8. PR remains the integration vehicle; no direct feature implementation is pushed to `main`.
9. The merge record explicitly states that GitHub Actions evidence was unavailable because of quota and that the owner authorized this temporary exception.
10. The exception is re-evaluated before P01 and must not be used to waive executable runtime/build/test/migration/security gates.

## P00.06 manual acceptance record

PR: **#13 — P00.06 Security and Data Classification**

Observed evidence:

- Scope/diff inspection: **PASS** — documentation, governance, schema and validator/workflow files only; no runtime/kernel/business implementation.
- Security standard evidence files: **PASS by repository inspection** — `SECURITY_STANDARD.md`, `DATA_CLASSIFICATION.md`, `SECURITY_CONTROL_MATRIX.md`, classification JSON Schema and ADR-0005 exist.
- Canonical state dependency check: **PASS by repository inspection** — P00.01 through P00.06 are ordered; P00.04/P00.05/P00.06 prerequisites satisfy P00.07 activation.
- State/status reconciliation: **PASS by repository inspection** — branch declares `6 / 10 done` and `P00.07 active` consistently.
- Runtime build/unit/integration/migration gates: **N/A** — P00 implementation lock prohibits runtime/kernel/business code.
- GitHub Actions hosted validation: **BLOCKED** — quota exhausted/disabled; issue #14 contains runs showing jobs fail before runner steps start.
- Full governance validator execution: **NOT RUN** — hosted runner unavailable; this must be rerun when Actions returns.

## Expiry / revocation

This exception expires at the earliest of:

- GitHub Actions quota/runner execution becoming available again;
- the P00 architecture freeze exit review;
- before any P01 implementation merge;
- explicit owner revocation.

When Actions return, the governance validator must be run again against the then-current canonical state. Any failure reopens the affected package according to the Definition of Done.

## Non-precedent rule

This exception is operational debt, not a permanent weakening of Omnexa's quality model. Future AI systems must not cite it as permission to skip tests, builds, migrations, security scans or release evidence for executable software.