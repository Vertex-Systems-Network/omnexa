# ADR-0012 — Versioned Module Dependency Requirements

Status: **proposed**  
Date: 2026-08-27  
Owners: Architecture / `kernel.modules`

## Context

P03.01 established manifest schema v1 and is complete. Schema v1 represents required and optional module dependencies as module-ID strings:

- `dependencies: []string`
- `optional_dependencies: []string`

It separately represents the module's own SemVer version and a constrained `required_platform_version`.

P03.03 is now the sole active work package and requires version-aware required and optional module dependency resolution, required dependency presence/version enforcement, and deterministic incompatible-version / incompatible-constraint rejection.

There is currently no canonical field in schema v1 that can express the version constraint a module requires from another module. Inferring compatibility (for example, same-major, any installed version, or an encoded `module@constraint` string) would create unapproved module-lifecycle semantics and conflict with mandatory change control.

GitHub issue #96 records this contract gap.

## Problem

P03.03 cannot satisfy its acceptance criteria deterministically from the accepted P03.01 schema without changing a module contract or introducing a second compatibility authority.

The change is therefore Class C under `docs/governance/CHANGE_CONTROL.md`: it changes module dependency/lifecycle semantics and requires an accepted ADR plus atomic reconciliation before implementation.

## Decision drivers

- fail closed for missing or incompatible required dependencies;
- deterministic results for identical manifests/registry state;
- explicit version constraints rather than inferred compatibility;
- preserve exact P03.01 schema-v1 parsing/validation evidence;
- avoid ambiguous string overloading;
- keep optional dependency failure selective rather than global;
- minimize dependency-constraint grammar complexity in the first resolver contract;
- maintain a single canonical source for module dependency requirements.

## Options considered

### Option A — Add structured dependency requirements in manifest schema v2

Introduce a new schema version whose required and optional module dependency declarations use structured records:

```json
{
  "id": "omnexa.example",
  "constraint": ">=1.2.0 <2.0.0"
}
```

The initial constraint grammar is deliberately bounded:

- one or more whitespace-separated comparators;
- comparator operators: `=`, `>`, `>=`, `<`, `<=`;
- comparator operands are valid SemVer versions;
- every comparator in the sequence must match (logical AND);
- no OR, wildcard, caret, tilde or implicit-range syntax in the initial contract;
- SemVer precedence determines ordering; build metadata does not affect precedence;
- prerelease versions participate according to normal SemVer precedence and only satisfy the explicit comparator sequence they actually match.

`platform_dependencies` remains a distinct platform-capability identifier class. `required_platform_version` remains the platform-version compatibility declaration and is not repurposed for module dependencies.

Benefits:
- version requirements live beside the dependency they constrain;
- deterministic validation and resolver behavior;
- no secondary compatibility source;
- future schema evolution remains explicit.

Costs/risks:
- requires schema-v2 parser/validator work and fixtures;
- P03.01 contract documentation/evidence must be reconciled without rewriting historical evidence;
- consumers must migrate dependency-bearing manifests to schema v2 before they can be resolver-eligible.

### Option B — Keep schema v1 and create an external compatibility matrix

Store module-version compatibility outside manifests and have P03.03 consult that authority.

Benefits:
- leaves the schema-v1 wire shape untouched.

Costs/risks:
- introduces a second source of truth;
- creates identity/version drift risk between a manifest and the matrix;
- complicates package portability and deterministic offline validation;
- requires separate ownership, versioning and reconciliation rules.

### Option C — Encode version constraints into existing strings

Examples: `omnexa.example@>=1.2.0` or another delimiter convention.

Benefits:
- superficially avoids a new JSON field shape.

Costs/risks:
- breaks the existing module-ID validator and schema-v1 meaning;
- overloads identity and policy into one string;
- makes compatibility and migration implicit;
- contradicts P03.01's stable machine-validatable ID contract.

### Option D — Infer compatibility when no constraint is declared

Examples include accepting any installed version or assuming same-major compatibility.

Benefits:
- no schema change.

Costs/risks:
- cannot prove P03.03 incompatible-version enforcement;
- silently invents lifecycle semantics;
- fails the repository's conflict and fail-closed rules.

## Proposed decision

Adopt **Option A**: introduce explicit structured module dependency requirements in **manifest schema v2** and keep module dependency compatibility in the manifest as the single canonical declaration.

This ADR is **proposed**, not accepted. No schema-v2 or resolver implementation is authorized by this file alone. P03.03 implementation remains stopped at the contract boundary until the ADR is reviewed/accepted and the affected canonical documents are reconciled in the governed change.

If accepted, the normative contract is:

1. schema v1 remains parseable and validated exactly according to its accepted P03.01 contract;
2. schema v2 represents required and optional module dependencies as `{id, constraint}` records;
3. schema v2 rejects empty/invalid constraints, duplicate dependency IDs, and cross-class dependency conflicts deterministically;
4. required dependency resolution checks both presence and the declared constraint and fails closed on absence/incompatibility;
5. optional dependency absence or incompatibility yields explicit selective-degradation metadata for that optional relation and does not globally invalidate unrelated modules;
6. schema-v1 manifests with no module dependencies remain resolver-eligible subject to all other gates;
7. schema-v1 manifests with required module dependencies remain parseable but are not install/enable eligible under P03.03 because the required version contract is unavailable; the resolver returns a stable fail-closed compatibility error requiring migration to schema v2;
8. schema-v1 optional module dependencies remain parseable; under P03.03 they produce explicit unresolved/degraded optional-dependency metadata until migrated to schema v2 rather than an inferred compatibility result;
9. no dependency declaration grants authorization, capability access, private package access or database access;
10. deterministic ordering/error reporting remains required independent of discovery enumeration order.

## Consequences

### Positive

- P03.03 can enforce version compatibility from canonical machine-readable input.
- Version compatibility is explicit and testable.
- Existing schema-v1 evidence remains meaningful instead of being silently reinterpreted.
- Optional dependency behavior stays selectively degradable.

### Negative / trade-offs

- A schema-v2 compatibility path must be added before the resolver can complete.
- Dependency-bearing schema-v1 manifests require migration for full resolver eligibility.
- P03.03 cannot proceed directly to implementation until governance reconciliation is accepted.

### Risks

- An underspecified constraint parser could create divergent compatibility results.
- Treating v1 dependency strings as unconstrained would accidentally weaken fail-closed behavior.
- Changing existing P03.01 completion evidence instead of adding forward evidence could corrupt historical provenance.

## Architecture impact

- affected domains/modules: `kernel.modules`, module manifest contract, dependency resolver;
- tenancy/security impact: none to tenancy; dependency declarations remain metadata and never authorize access;
- data ownership impact: none;
- API/event impact: no HTTP/event contract change; module manifest contract gains schema v2;
- deployment/operations impact: dependency-bearing manifests must migrate to v2 before install/enable eligibility under the resolver;
- compatibility/migration impact: dual v1/v2 parsing during transition; no silent reinterpretation of v1 dependency strings.

## Implementation constraints

If this ADR is accepted:

- do not overload schema-v1 dependency strings with constraint syntax;
- do not infer same-major, latest, or any-version compatibility;
- do not add an external compatibility matrix unless this ADR is superseded;
- keep comparison deterministic and bounded;
- keep stable safe validation/resolution error codes;
- retain all P03.01 regression fixtures and add explicit v1/v2 compatibility fixtures;
- do not weaken forbidden/private dependency checks;
- do not let resolver results grant runtime authority;
- justify any new external SemVer/constraint library separately under dependency-introduction policy; a library choice is not made by this ADR.

## Migration / transition plan

1. Review/accept this ADR or replace it with another explicit architectural decision.
2. In the accepting/reconciliation PR, update all affected canonical sources of truth, including as applicable `MODULE_STANDARD.md`, P03.01/P03.03 contract documentation, manifest specifications/fixtures, and roadmap/state references required by change control.
3. Preserve the historical P03.01 completion evidence as evidence of schema-v1 completion; add forward compatibility evidence rather than rewriting the historical run.
4. Implement dual schema-version parsing/validation with explicit v1 and v2 behavior.
5. Add migration fixtures/examples for dependency-bearing v1 manifests.
6. Only after the contract is accepted and reconciled, implement P03.03 resolver semantics against the accepted model.
7. Roll back by rejecting/superseding this proposal before implementation; after implementation, use a superseding ADR and forward-fix migration rather than silently changing constraint meaning.

## Verification

Acceptance/implementation evidence must include:

- governance validation for the ADR and reconciled canonical docs;
- retained P03.01 schema-v1 regression suite passing unchanged where semantics are unchanged;
- schema-v2 positive/negative parser and validation fixtures;
- deterministic constraint tests covering exact, lower/upper bound, conflict, prerelease and build-metadata precedence cases;
- required dependency absent/incompatible fail-closed fixtures;
- optional dependency absent/incompatible selective-degradation fixtures;
- duplicate and cross-class dependency rejection;
- deterministic results independent of registry/discovery ordering;
- all retained P01/P02/P03.01/P03.02 regressions plus P03.03 dedicated verification on the exact implementation head.

## Documents/work packages affected

- `docs/architecture/MODULE_STANDARD.md`
- `docs/architecture/DEPENDENCY_MATRIX.md` if it currently states compatibility semantics that need reconciliation
- `docs/roadmap/work-packages/P03.01.md` (forward compatibility note only; historical completion evidence remains immutable)
- `docs/roadmap/work-packages/P03.03.md`
- `docs/ai/handoffs/P03.03.md`
- manifest schema/specification, parser/validator and fixtures
- GitHub issue #96

## Supersedes / superseded by

- Supersedes: none
- Superseded by: none

## References

- GitHub issue #96 — P03.03 dependency-version contract gap
- `docs/governance/CHANGE_CONTROL.md`
- `docs/governance/AI_EXECUTION_POLICY.md`
- `docs/roadmap/work-packages/P03.01.md`
- `docs/roadmap/work-packages/P03.03.md`
- `docs/architecture/MODULE_STANDARD.md`
- `docs/architecture/DEPENDENCY_MATRIX.md`
