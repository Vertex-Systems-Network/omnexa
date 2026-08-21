# Contributing to Omnexa

Omnexa is architecture-first and phase-gated. Contribution speed is never allowed to override domain ownership, tenancy, authorization, auditability, contract compatibility or execution-state rules.

## Mandatory read order

Before changing anything, read:

1. `AGENTS.md`
2. `docs/governance/PRODUCT_CONSTITUTION.md`
3. `docs/architecture/SYSTEM_ARCHITECTURE.md`
4. `docs/architecture/MODULE_STANDARD.md`
5. `docs/architecture/GLOSSARY.md`
6. `docs/architecture/NAMING_STANDARD.md`
7. `docs/architecture/DOMAIN_OWNERSHIP.md`
8. `docs/architecture/DEPENDENCY_MATRIX.md`
9. `docs/roadmap/MASTER_PLAN.md`
10. `docs/roadmap/STATUS.md`
11. `docs/roadmap/STATE.json`
12. `docs/governance/CHANGE_CONTROL.md`
13. `docs/governance/DEFINITION_OF_DONE.md`
14. relevant ADRs.

## Allowed work

The machine-readable `docs/roadmap/STATE.json` controls which work package is active. Do not start future phases early. Valuable out-of-scope work must be recorded/proposed rather than silently implemented.

## Branches and pull requests

- Do not intentionally develop on `main`.
- Use one focused branch per work package/change.
- Keep unrelated changes out of the PR.
- Use the repository PR template completely.
- Architecture changes require an ADR before implementation.
- State/progress may move forward only when evidence exists.

Suggested branch forms:

```text
feat/pNN-NN-short-description
fix/pNN-NN-short-description
docs/pNN-NN-short-description
chore/pNN-NN-short-description
```

## Architecture rules

- One authoritative owner per write model.
- No direct cross-module database writes.
- No circular required module dependencies.
- Cross-domain calls use public capabilities, events, workflows or approved projections.
- Tenant and authorization context are mandatory where relevant.
- Protected mutations must be auditable.
- Public contracts are versioned.
- AI acts through governed capabilities, never raw database authority.

## Change workflow

1. Identify phase/work package.
2. Inspect current repository state.
3. Identify owner domain and affected contracts.
4. Confirm scope/non-scope.
5. Add/update ADR if architecture changes.
6. Implement the smallest complete change.
7. Add tests/validators required by the phase.
8. Run quality gates.
9. Reconcile `STATUS.md` and `STATE.json` only when evidence supports transition.
10. Record execution evidence in the ledger when a package/phase transition occurs.

## Commit guidance

Prefer conventional intent prefixes:

```text
feat:
fix:
docs:
chore:
test:
ci:
refactor:
```

Commit messages must describe the actual change, not claim broader completion than evidence supports.

## Definition of done

A change is not done because files exist. Completion means the applicable gates in `docs/governance/DEFINITION_OF_DONE.md` are satisfied. `NOT RUN`, `BLOCKED`, `FAIL` and `N/A` must be reported truthfully.

## Generated/AI contributions

AI-generated changes are held to the same standards as human changes. AI systems must not invent architecture, dependencies, phase transitions or hidden compatibility assumptions. Repository governance is the source of truth.
