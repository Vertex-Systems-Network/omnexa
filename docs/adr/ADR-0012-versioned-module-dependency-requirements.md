# ADR-0012 — Versioned Module Dependency Requirements

- Status: **Accepted**
- Date: 2026-08-27
- Decision class: Class C — module manifest / dependency lifecycle contract
- Owners: Architecture / `kernel.modules`

> This accepted status becomes authoritative only when the accepting/reconciliation PR merges to protected `main`. The ADR does not by itself mark P03.03 complete and does not authorize P03.04+.

## Context

P03.01 completed manifest schema v1. Schema v1 represents required and optional module dependencies as module-ID strings:

- `dependencies: []string`
- `optional_dependencies: []string`

It separately represents the module's own SemVer version and `required_platform_version`.

P03.02 completed deterministic discovery/registry semantics. Discovery accepts at most one record per module ID; duplicate identities or conflicting versions fail closed. The public `RegistryRecord` retains module identity/version/owner/source metadata, but not dependency declarations.

P03.03 is the sole active package and requires deterministic version-aware required and optional dependency resolution, incompatible-version detection, stable required-graph ordering, cycle detection and selective optional degradation.

Schema v1 cannot express the version constraint that one module requires from another. P03.03 also must not reparse or separately pair arbitrary raw manifests after registry creation because doing so could evaluate dependency policy against declarations different from the manifest snapshot that discovery validated.

GitHub issue #96 records this contract gap.

## Problem

P03.03 cannot satisfy its accepted version-aware resolver criteria without either:

1. explicitly evolving the manifest dependency contract; or
2. introducing a second compatibility authority or inferred compatibility rule.

Inferred compatibility, string overloading and independently re-paired manifest inputs are rejected because they are ambiguous and can fail open.

## Decision drivers

- fail closed for missing or incompatible required dependencies;
- deterministic output for identical validated discovery state;
- explicit constraints rather than inferred compatibility;
- preserve historical P03.01 schema-v1 parsing evidence;
- preserve P03.02 one-record-per-module-ID public registry semantics;
- bind resolver declarations to the exact validated discovery snapshot;
- avoid parser fallthrough and ambiguous string overloading;
- keep optional dependency failure selective rather than global;
- keep the first constraint grammar intentionally bounded;
- bound all untrusted parser work;
- keep the manifest as the single canonical dependency-requirement declaration.

## Decision

Adopt **Option A: manifest schema v2 with structured module dependency requirements**.

Schema v2 uses the same top-level dependency field names, but required and optional module dependencies become exact structured records:

```json
{
  "schema_version": 2,
  "dependencies": [
    {
      "id": "omnexa.example",
      "constraint": ">=1.2.0 <2.0.0"
    }
  ],
  "optional_dependencies": []
}
```

Each module dependency record contains exactly:

- `id`
- `constraint`

Unknown record fields are rejected.

`platform_dependencies` remains a distinct platform-capability identifier class. `required_platform_version` remains the platform-version declaration and is not repurposed for module dependencies.

### Constraint grammar

The initial module dependency constraint grammar is deliberately small and deterministic:

- one to sixteen comparators;
- maximum encoded constraint length: 256 bytes;
- comparators separated by exactly one ASCII space;
- no leading or trailing whitespace;
- operators: `=`, `>`, `>=`, `<`, `<=`;
- operands: strict SemVer 2.0.0 versions;
- all comparators in a sequence must match (logical AND);
- no OR, wildcard, caret, tilde, comma-separated, implicit range, locale-specific or package-manager-specific shorthand;
- SemVer precedence determines ordering and comparator equality;
- build metadata does not affect precedence;
- prerelease versions participate only according to the explicit comparator sequence and SemVer precedence;
- malformed, empty or over-bound constraints are validation failures, never resolver warnings.

## Normative contract

1. Schema v1 remains parseable and validated according to its accepted P03.01 contract.
2. Manifest parsing first reads a bounded top-level `schema_version` discriminator and then routes to a version-specific strict decoder.
3. Schema-v2 structured dependency records are never decoded through or silently coerced into the schema-v1 manifest shape.
4. Missing, malformed or unsupported schema versions fail with stable safe schema errors before lifecycle or resolver work.
5. Schema v2 represents required and optional module dependencies as exact `{id, constraint}` records.
6. Schema v2 rejects empty/invalid/over-bound constraints, duplicate dependency IDs, self-dependencies and cross-class dependency conflicts deterministically.
7. Required dependency resolution checks both presence and declared constraint and fails closed on absence or incompatibility.
8. Optional dependency absence or incompatibility produces explicit selective-degradation metadata for that relation and does not globally invalidate unrelated modules.
9. Only required dependency edges create install/enable ordering authority and participate in the release-blocking topological cycle gate.
10. Optional edges do not participate in the required dependency topological order or global required-cycle result.
11. Schema-v1 manifests with no module dependencies remain resolver-eligible subject to all other gates.
12. Schema-v1 manifests with required module dependencies remain parseable but are not install/enable eligible under P03.03 because no required version contract exists; the resolver returns a stable fail-closed compatibility error requiring migration to schema v2.
13. Schema-v1 optional module dependencies remain parseable but produce explicit unresolved/degraded optional metadata under P03.03 until migrated to schema v2; compatibility is not inferred.
14. P03.02's one-record-per-module-ID invariant remains authoritative. P03.03 resolves one discovered module version against one declared constraint; multi-version selection is out of scope.
15. Resolver dependency declarations come only from the exact normalized validated manifest snapshot atomically associated with the registry record during discovery.
16. Callers cannot pair a registry created from one manifest set with dependency declarations reparsed from another set.
17. No dependency declaration or resolver result grants authorization, capability access, private-package access, tenant authority or database authority.
18. Ordering, diagnostics and degradation metadata remain deterministic independent of discovery enumeration order.

## Parser and version-dispatch contract

The accepted P03.01 schema-v1 parser remains a compatibility contract.

Implementation must:

1. enforce the existing manifest byte-size bound before decoding;
2. decode only enough bounded top-level JSON metadata to obtain `schema_version`;
3. dispatch version `1` to the preserved schema-v1 decoder/validator;
4. dispatch version `2` to a separate strict schema-v2 decoder/validator;
5. reject unsupported versions with stable safe errors;
6. keep unknown-field rejection version-specific;
7. never fall back from one schema decoder to another after validation failure;
8. never execute module code, filesystem scanning or network access during parsing/validation.

One strict SemVer 2.0.0 parser/comparator path must be used for schema-v2 module versions and dependency constraint operands. The existing permissive regex must not be treated as a complete SemVer precedence engine.

## Resolver input and discovery binding

P03.02 parses manifests during `Discover` and currently exposes public registry records without dependency declarations. P03.03 must preserve a single provenance path instead of reparsing raw manifests later.

The accepted design is:

1. version-specific parsing/validation produces a normalized, immutable-by-convention validated manifest snapshot containing schema version, identity/version/owner and normalized dependency declarations;
2. discovery atomically derives the existing public `RegistryRecord` and stores the corresponding validated snapshot in the same registry construction result, keyed by the same unique module ID;
3. public P03.02 `List`/`Lookup` record semantics remain unchanged unless a later governed decision says otherwise;
4. the normalized dependency snapshot may remain package-private to `kernel.modules`;
5. the resolver reads dependency declarations only from this registry-bound validated snapshot;
6. missing snapshots, identity/version/owner mismatches, duplicate associations or impossible schema/dependency state fail closed with stable diagnostics;
7. snapshots exposed across internal paths are copied/read-only so callers cannot mutate registry evidence after discovery;
8. resolver execution performs no filesystem/network rediscovery and executes no package code.

The manifest remains the source of truth. The registry-bound snapshot is only the validated normalized representation of the exact manifest used by discovery.

## Registry and graph semantics

P03.02's existing public registry behavior remains authoritative:

- one module ID resolves to at most one public `RegistryRecord` and one bound validated manifest snapshot;
- duplicate identities fail discovery;
- conflicting versions fail discovery;
- no last-write-wins behavior;
- no multi-version registry selection.

P03.03 graph semantics are:

- self-dependency is invalid in both required and optional classes;
- required compatibility compares the discovered module's strict SemVer version to the declaring module's constraint;
- required edges alone determine deterministic topological order;
- a cycle of required edges is release-blocking invalid state;
- optional edges are evaluated independently for availability/degradation and cannot make an otherwise valid required graph globally cyclic;
- forbidden/private dependency checks remain orthogonal and may invalidate a module regardless of version compatibility.

Multi-version selection, SAT solving, automatic upgrades/downgrades, remote package acquisition and a second raw-manifest resolver input are explicitly out of scope.

## Alternatives considered

### Option B — external compatibility matrix

Rejected because it creates a second source of truth, identity/version drift risk, portability problems and separate reconciliation/versioning obligations.

### Option C — encode constraints into schema-v1 strings

Rejected because it overloads identity and policy, breaks the existing module-ID meaning/validator and silently changes schema-v1 semantics.

### Option D — infer compatibility

Rejected. Accept-any-version, same-major or similar implicit behavior cannot prove incompatible-version enforcement and violates fail-closed change control.

## Compatibility and migration

- Historical P03.01 schema-v1 completion evidence remains immutable.
- Historical P03.02 registry/discovery completion evidence remains immutable.
- Forward P03.03 implementation adds schema-v2 and discovery-binding evidence without relabeling historical runs.
- Dependency-bearing schema-v1 manifests must migrate to schema v2 for full required-dependency resolver eligibility.
- Schema-v1 manifests without module dependencies remain eligible subject to other gates.
- Schema-v1 optional dependency declarations degrade explicitly until migrated; they are not interpreted as unconstrained compatible relations.
- Public P03.02 registry cardinality, lookup/list ordering and duplicate/version-conflict behavior remain unchanged.

## Security and authority impact

- Tenancy model: unchanged.
- Authorization model: unchanged.
- Data ownership: unchanged.
- HTTP/event contracts: unchanged.
- Dependency declarations remain metadata only.
- Resolver results are eligibility/degradation metadata only.
- Kernel-to-business and private cross-module dependency prohibitions remain unchanged.
- No manifest or resolver path grants permissions, capabilities or storage authority.

## Implementation constraints

Under this accepted ADR, P03.03 implementation must not:

- overload schema-v1 dependency strings with constraint syntax;
- infer same-major, latest or any-version compatibility;
- introduce an external compatibility matrix unless this ADR is superseded;
- add multi-version solving or package selection;
- accept independently reparsed raw manifests as resolver authority when a registry snapshot exists;
- let optional edges affect required graph topological order;
- mutate public P03.02 registry behavior merely to simplify resolver implementation;
- weaken forbidden/private dependency checks;
- let resolver results grant runtime authority.

Implementation must:

- keep comparison deterministic, strict and bounded;
- retain existing manifest byte/list bounds;
- add explicit dependency-record, constraint-byte and comparator-count bounds;
- use one strict SemVer 2.0.0 implementation path;
- preserve stable safe validation/resolution error codes and deterministic diagnostic ordering;
- retain P03.01/P03.02 regressions;
- add explicit v1/v2 compatibility and registry-binding fixtures;
- justify any new external SemVer/constraint dependency separately under repository dependency-introduction policy.

## Verification requirements

The accepting/reconciliation PR must pass canonical governance with existing P01/P02/P03.01/P03.02 regressions.

The later P03.03 implementation must additionally prove:

- schema-version dispatch with strict v1/v2 decoders and no fallthrough;
- schema-v2 positive/negative parser fixtures;
- unknown dependency-record field rejection;
- self-dependency rejection;
- dependency-record/list/constraint resource bounds;
- strict SemVer conformance, including invalid numeric prerelease identifiers, prerelease ordering and build-metadata precedence equality;
- exact/lower/upper/conflicting comparator behavior and malformed whitespace rejection;
- registry-to-normalized-manifest binding and fail-closed mismatch/missing/mutable association cases;
- required dependency absent/incompatible fail-closed behavior;
- optional dependency absent/incompatible selective-degradation behavior;
- required-cycle rejection;
- optional-cycle fixtures proving optional edges do not globally invalidate the required graph;
- duplicate and cross-class dependency rejection;
- retained P03.02 duplicate/version-conflict and public lookup/list behavior;
- deterministic results independent of discovery ordering;
- dedicated `scripts/verify_p03_03.sh` plus repository Go quality and all retained P01/P02/P03.01/P03.02 regressions on the exact implementation head.

## Transition and rollback

This ADR is accepted as the governing architecture decision for the P03.03 dependency-version prerequisite. It becomes authoritative when the accepting/reconciliation PR merges to protected `main`.

After that merge:

1. rehydrate protected `main` and canonical P03.03 state;
2. create a separate implementation branch from the exact new main SHA;
3. implement schema-version dispatch, schema v2, registry-bound normalized manifests and resolver behavior only inside P03.03 scope;
4. verify on exact implementation head;
5. use a later separate closure/state-transition PR to mark P03.03 done and activate P03.04 only if all gates pass.

Before merge, this decision can be abandoned by closing/revising the accepting PR. After merge, changing the meaning requires a superseding ADR and forward compatibility/migration plan; do not silently reinterpret constraints.

## Documents/work packages reconciled by the accepting PR

- `docs/architecture/MODULE_STANDARD.md`
- `docs/architecture/DEPENDENCY_MATRIX.md`
- `docs/roadmap/work-packages/P03.01.md` — forward compatibility note only
- `docs/roadmap/work-packages/P03.02.md` — forward compatibility note only
- `docs/roadmap/work-packages/P03.03.md`
- `docs/ai/handoffs/P03.03.md`
- GitHub issue #96 / Linear mirror

Implementation will later touch the manifest parser/validator, discovery/registry internals, fixtures and dedicated P03.03 verification under this accepted contract.

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
