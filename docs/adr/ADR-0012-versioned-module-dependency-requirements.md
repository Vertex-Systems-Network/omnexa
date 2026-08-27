# ADR-0012 — Versioned Module Dependency Requirements

Status: **proposed**  
Date: 2026-08-27  
Owners: Architecture / `kernel.modules`

## Context

P03.01 established manifest schema v1 and is complete. Schema v1 represents required and optional module dependencies as module-ID strings:

- `dependencies: []string`
- `optional_dependencies: []string`

It separately represents the module's own SemVer version and a constrained `required_platform_version`.

P03.02 then established deterministic discovery/registry semantics. The current registry accepts exactly one discovered record per module ID; duplicate identities or conflicting versions fail discovery before dependency resolution.

P03.03 is now the sole active work package and requires version-aware required and optional module dependency resolution, required dependency presence/version enforcement, and deterministic incompatible-version / incompatible-constraint rejection.

There is currently no canonical field in schema v1 that can express the version constraint a module requires from another module. Inferring compatibility (for example, same-major, any installed version, or an encoded `module@constraint` string) would create unapproved module-lifecycle semantics and conflict with mandatory change control.

GitHub issue #96 records this contract gap.

## Problem

P03.03 cannot satisfy its acceptance criteria deterministically from the accepted P03.01 schema without changing a module contract or introducing a second compatibility authority.

The prerequisite manifest-contract change is therefore **Class C** under `docs/governance/CHANGE_CONTROL.md`: it changes module dependency/lifecycle semantics and requires an accepted ADR plus atomic reconciliation before implementation. P03.03 itself remains the active Class B implementation package; this ADR governs the Class C prerequisite needed before its version-aware resolver portion can proceed.

## Decision drivers

- fail closed for missing or incompatible required dependencies;
- deterministic results for identical manifests/registry state;
- explicit version constraints rather than inferred compatibility;
- preserve exact P03.01 schema-v1 parsing/validation evidence;
- preserve P03.02 one-record-per-module-ID registry semantics unless separately changed by a future ADR;
- avoid ambiguous string overloading or parser fallthrough between schema versions;
- keep optional dependency failure selective rather than global;
- minimize dependency-constraint grammar complexity in the first resolver contract;
- bound all untrusted constraint input and parser work;
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

The same top-level field names remain canonical:

```json
{
  "schema_version": 2,
  "dependencies": [
    {"id": "omnexa.example", "constraint": ">=1.2.0 <2.0.0"}
  ],
  "optional_dependencies": []
}
```

Each dependency record contains exactly `id` and `constraint`; unknown fields are rejected by the schema-v2 decoder.

The initial constraint grammar is deliberately bounded and deterministic:

- one to sixteen comparators;
- comparators are separated by exactly one ASCII space;
- no leading or trailing whitespace;
- comparator operators: `=`, `>`, `>=`, `<`, `<=`;
- comparator operands are strict SemVer 2.0.0 versions;
- maximum encoded constraint length: 256 bytes;
- every comparator in the sequence must match (logical AND);
- no OR, wildcard, caret, tilde, implicit range, comma-separated, locale-dependent or package-manager-specific syntax in the initial contract;
- SemVer precedence determines ordering and comparator equality; build metadata does not affect precedence;
- prerelease versions are not silently excluded by npm/Cargo-style hidden rules: they participate only according to explicit comparator expressions and SemVer precedence;
- malformed or over-bound constraints are validation failures, never resolver warnings.

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
2. schema dispatch must first read a bounded top-level `schema_version` discriminator, then route to a version-specific strict decoder; schema-v2 structured dependency records must never be decoded through or silently coerced into the schema-v1 `Manifest` shape;
3. unsupported schema versions fail with a stable safe schema error before lifecycle/resolution work;
4. schema v2 represents required and optional module dependencies as exact `{id, constraint}` records and rejects unknown record fields;
5. schema v2 rejects empty/invalid/over-bound constraints, duplicate dependency IDs, self-dependencies, and cross-class dependency conflicts deterministically;
6. required dependency resolution checks both presence and the declared constraint and fails closed on absence/incompatibility;
7. optional dependency absence or incompatibility yields explicit selective-degradation metadata for that optional relation and does not globally invalidate unrelated modules;
8. optional dependency edges do not create install/enable ordering authority and do not participate in the required-dependency topological cycle gate; required edges alone determine the release-blocking dependency order/cycle result;
9. schema-v1 manifests with no module dependencies remain resolver-eligible subject to all other gates;
10. schema-v1 manifests with required module dependencies remain parseable but are not install/enable eligible under P03.03 because the required version contract is unavailable; the resolver returns a stable fail-closed compatibility error requiring migration to schema v2;
11. schema-v1 optional module dependencies remain parseable; under P03.03 they produce explicit unresolved/degraded optional-dependency metadata until migrated to schema v2 rather than an inferred compatibility result;
12. P03.02's one-record-per-module-ID registry invariant remains authoritative: P03.03 resolves one discovered version for an ID against one declared constraint; selecting among multiple installed/discovered versions is out of scope and requires separate architecture authority;
13. no dependency declaration grants authorization, capability access, private package access or database access;
14. deterministic ordering, diagnostics and degradation metadata remain stable independent of discovery enumeration order.

## Parser/version-dispatch contract

The existing P03.01 parser is intentionally strict: it rejects unknown fields, requires all schema-v1 fields and validates `SchemaVersion == 1`. Schema v2 must not weaken that evidence by changing schema-v1 decoding semantics in place.

If this ADR is accepted, parsing must use an explicit version-dispatch boundary:

1. enforce the existing manifest byte-size bound before decoding;
2. decode only enough bounded top-level JSON metadata to obtain `schema_version`;
3. dispatch `1` to the preserved schema-v1 decoder/validator and `2` to a separate schema-v2 decoder/validator;
4. reject missing, malformed or unsupported versions with stable safe errors;
5. keep unknown-field rejection version-specific;
6. never attempt parser fallback (for example, v2 -> v1) after a version-specific validation failure;
7. never execute module code, filesystem scanning or network access during parsing/validation.

This avoids parser ambiguity, prevents a structured v2 manifest from surfacing as an accidental generic JSON/type error in the v1 path, and preserves historical v1 verification as a real compatibility contract.

## Registry and graph semantics

P03.02 discovery currently fails the whole registry when the same module ID appears more than once, whether as a duplicate or version conflict. This ADR does not introduce package-manager-style multi-version solving.

For P03.03 under this decision:

- one module ID resolves to at most one `RegistryRecord`;
- required dependency compatibility compares that record's strict SemVer version to the declaring module's constraint;
- a module may not depend on itself in either required or optional dependency classes;
- only required dependency edges participate in deterministic topological ordering;
- a cycle of required edges is release-blocking invalid state;
- optional relations are evaluated independently for availability/degradation and cannot turn an otherwise valid required graph into a global cycle failure;
- forbidden/private dependency checks remain orthogonal and may still invalidate a module regardless of version compatibility.

Multi-version selection, SAT solving, automatic upgrades/downgrades and remote package acquisition are explicitly out of scope.

## Consequences

### Positive

- P03.03 can enforce version compatibility from canonical machine-readable input.
- Version compatibility is explicit, bounded and testable.
- Existing schema-v1 evidence remains meaningful instead of being silently reinterpreted.
- P03.02 registry cardinality remains stable.
- Optional dependency behavior stays selectively degradable without corrupting required topological ordering.

### Negative / trade-offs

- A schema-v2 compatibility path and explicit parser dispatch must be added before the resolver can complete.
- Dependency-bearing schema-v1 manifests require migration for full resolver eligibility.
- The initial comparator grammar is intentionally less expressive than common package-manager grammars.
- P03.03 cannot proceed directly to implementation until governance reconciliation is accepted.

### Risks

- An underspecified constraint parser could create divergent compatibility results.
- Reusing the existing permissive SemVer regex as if it were a full SemVer 2.0.0 parser could accept invalid prerelease identifiers or compare versions incorrectly; schema-v2 version/constraint comparison therefore requires one deterministic strict SemVer implementation path.
- Treating v1 dependency strings as unconstrained would accidentally weaken fail-closed behavior.
- Allowing parser fallback across schema versions could make malformed manifests ambiguous or bypass version-specific validation.
- Changing existing P03.01 completion evidence instead of adding forward evidence could corrupt historical provenance.

## Architecture impact

- affected domains/modules: `kernel.modules`, module manifest contract, dependency resolver;
- tenancy/security impact: none to tenancy; dependency declarations remain metadata and never authorize access;
- data ownership impact: none;
- API/event impact: no HTTP/event contract change; module manifest contract gains schema v2;
- deployment/operations impact: dependency-bearing manifests must migrate to v2 before install/enable eligibility under the resolver;
- compatibility/migration impact: explicit version-dispatch with dual v1/v2 parsing during transition; no silent reinterpretation of v1 dependency strings;
- P03.02 impact: retain single-record-per-module-ID discovery and duplicate/version-conflict behavior; add only forward evidence showing v2 manifests participate without weakening those invariants.

## Implementation constraints

If this ADR is accepted:

- do not overload schema-v1 dependency strings with constraint syntax;
- do not infer same-major, latest, or any-version compatibility;
- do not add an external compatibility matrix unless this ADR is superseded;
- do not add multi-version solving or package selection to P03.03;
- do not let optional edges affect required graph topological order;
- keep comparison deterministic, strict and bounded;
- use one strict SemVer 2.0.0 parser/comparator path for module versions and constraint operands rather than duplicating subtly different version rules;
- retain the existing manifest byte/list bounds and add explicit dependency-record, constraint-byte and comparator-count bounds;
- keep stable safe validation/resolution error codes and deterministic diagnostic ordering;
- retain all P03.01/P03.02 regression fixtures and add explicit v1/v2 compatibility fixtures;
- do not weaken forbidden/private dependency checks;
- do not let resolver results grant runtime authority;
- justify any new external SemVer/constraint library separately under dependency-introduction policy; a library choice is not made by this ADR.

## Migration / transition plan

1. Review/accept this ADR or replace it with another explicit architectural decision.
2. In the accepting/reconciliation PR, update all affected canonical sources of truth, including as applicable `MODULE_STANDARD.md`, P03.01/P03.02/P03.03 contract documentation, manifest specifications/fixtures, and roadmap/state references required by change control.
3. Preserve the historical P03.01 completion evidence as evidence of schema-v1 completion; add forward compatibility evidence rather than rewriting the historical run.
4. Preserve P03.02 single-version registry semantics and add forward tests proving schema-v2 discovery feeds the same deterministic registry model.
5. Implement explicit schema-version dispatch plus separate strict v1/v2 parsing/validation with no fallback.
6. Add migration fixtures/examples for dependency-bearing v1 manifests.
7. Only after the contract is accepted and reconciled, implement P03.03 resolver semantics against the accepted model.
8. Roll back by rejecting/superseding this proposal before implementation; after implementation, use a superseding ADR and forward-fix migration rather than silently changing constraint meaning.

## Verification

Acceptance/implementation evidence must include:

- governance validation for the ADR and reconciled canonical docs;
- retained P03.01 schema-v1 regression suite passing unchanged where semantics are unchanged;
- parser-dispatch fixtures proving v1 and v2 use their own strict decoders and malformed/unsupported versions cannot fall through;
- schema-v2 positive/negative parser and validation fixtures, including unknown dependency-record fields, self-dependencies and resource-bound failures;
- strict SemVer conformance fixtures, including invalid numeric prerelease identifiers, prerelease ordering and build-metadata precedence equality;
- deterministic constraint tests covering exact, lower/upper bound, conflict, comparator-count/length bounds and malformed whitespace;
- required dependency absent/incompatible fail-closed fixtures;
- optional dependency absent/incompatible selective-degradation fixtures;
- required-cycle fixtures plus optional-cycle fixtures proving optional edges do not globally invalidate required graph order;
- duplicate and cross-class dependency rejection;
- P03.02 duplicate/version-conflict registry behavior retained unchanged;
- deterministic results independent of registry/discovery ordering;
- all retained P01/P02/P03.01/P03.02 regressions plus P03.03 dedicated verification on the exact implementation head.

## Documents/work packages affected

- `docs/architecture/MODULE_STANDARD.md`
- `docs/architecture/DEPENDENCY_MATRIX.md` if it currently states compatibility semantics that need reconciliation
- `docs/roadmap/work-packages/P03.01.md` (forward compatibility note only; historical completion evidence remains immutable)
- `docs/roadmap/work-packages/P03.02.md` (forward compatibility note only if parser/registry evidence needs extension; historical completion evidence remains immutable)
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
- `docs/roadmap/work-packages/P03.02.md`
- `docs/roadmap/work-packages/P03.03.md`
- `docs/architecture/MODULE_STANDARD.md`
- `docs/architecture/DEPENDENCY_MATRIX.md`
- `kernel/internal/modules/manifest.go`
- `kernel/internal/modules/validate.go`
- `kernel/internal/modules/registry.go`
