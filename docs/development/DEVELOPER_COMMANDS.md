# Omnexa Developer Command Contract

Status: **Canonical v1 semantic contract**  
Work package: **P00.08**

This document freezes developer command **semantics**, not the final CLI implementation. P01+ tooling may expose these through an `omnexa` CLI, Task/Make wrapper or equivalent, but local development and CI must call the same underlying operations.

## Development lifecycle

```text
omnexa dev bootstrap
omnexa dev up
omnexa dev down
omnexa dev status
omnexa dev reset
```

### `dev bootstrap`

Verifies pinned toolchains, prepares non-secret local config, provisions/starts required disposable local dependencies, waits for health, installs dependencies, applies migrations and optionally deterministic development seed data.

### `dev up`

Starts required local dependencies/application processes without destructive reset.

### `dev down`

Stops local resources without deleting persistent development volumes unless explicitly requested.

### `dev status`

Reports toolchain versions, local service state/health, configuration profile, database migration state and useful service endpoints without leaking secrets.

### `dev reset`

Explicitly destroys/recreates only resources proven to belong to the local Omnexa development profile. It must fail closed if target identity is ambiguous.

## Database lifecycle

```text
omnexa db status
omnexa db migrate
omnexa db fresh
omnexa db seed
```

- `db status`: show connection/profile and migration state safely;
- `db migrate`: apply pending migrations;
- `db fresh`: create a clean local/test schema from zero; destructive and environment-guarded;
- `db seed`: load deterministic approved reference/development fixtures.

Production migration execution is a separate release/operations concern and must not inherit unsafe local-reset semantics.

## Quality commands

```text
omnexa verify governance
omnexa verify format
omnexa verify lint
omnexa verify static
omnexa verify unit
omnexa verify contracts
omnexa verify integration
omnexa verify migrations
omnexa verify security
omnexa verify module-lifecycle
omnexa verify build
omnexa verify release
omnexa verify all
```

Each command maps to P00.07 quality gate classes and returns a non-zero exit status on failure. Machine-readable output should eventually be available for CI/release evidence.

`verify all` is a deterministic aggregate, not a secret CI-only path.

## P01 executable mappings

P01.01 established the first repository-owned verification wrapper:

```text
bash scripts/verify_p01_01.sh
```

It remains a required regression gate for the completed Go workspace/build skeleton and covers pinned Go version, format/static checks, unit tests, dependency boundary, build/source metadata and process smoke behavior.

The active P01.02 configuration/environment package adds:

```text
bash scripts/verify_p01_02.sh
```

P01.02 maps the currently applicable quality semantics as follows:

- `verify format` -> `gofmt` cleanliness across the kernel Go source;
- `verify static` -> exact pinned Go version, canonical workspace/module metadata, `go vet`, and Omnexa dependency-boundary validation;
- `verify unit` -> kernel unit tests plus race-enabled configuration tests;
- `verify security` -> invalid/unknown configuration rejection, secret redaction, raw-value non-disclosure and isolated-loader tests;
- `verify build` -> trimmed kernel build, default/config-file/environment precedence startup smoke and fail-closed invalid startup.

The GitHub-hosted `governance` job invokes both the completed P01.01 regression wrapper and the active P01.02 wrapper. P01.02 does **not** implement the final `omnexa verify ...` CLI; that broader developer CLI belongs to P01.12. Database migration, event replay, module lifecycle and other checks without an implemented owning capability remain governed `N/A`, not fabricated PASS.

## Configuration startup contract during P01.02

The kernel process accepts only explicit governed configuration sources. Current application-level controls are:

```text
OMNEXA_CONFIG_FILE=<explicit JSON file path>   # optional; no implicit file discovery
OMNEXA_ENVIRONMENT=local|ci|preview|test|staging|production
```

The configuration file contains lowercase `snake_case` keys. Current application schema exposes only `environment`; the generic loader supports typed/required/sensitive definitions for later governed package-owned keys without pre-creating later infrastructure or business settings.

Precedence is deterministic:

```text
default -> explicit JSON config file -> environment variable -> explicit in-process test override
```

Unknown `OMNEXA_*` configuration variables fail in strict application mode. Configuration errors must identify the key/problem without printing raw sensitive values. Test overrides are instance-local and are not a global runtime feature-flag/config registry.

## Module developer commands

Future module SDK/CLI should support semantically stable operations such as:

```text
omnexa module create <id>
omnexa module validate <id>
omnexa module test <id>
omnexa module package <id>
```

`module create` uses the governed module manifest/ownership/naming standard; it may not generate direct cross-module dependencies.

## Contract commands

Future tooling should support explicit validation/generation for:

- OpenAPI contracts;
- event schemas;
- data-classification declarations;
- module manifests;
- quality evidence records;
- SDK/generated artifacts.

Generation and validation must be reproducible from pinned tool versions.

## Exit-code contract

- `0`: operation succeeded;
- non-zero: operation failed or required prerequisite unavailable.

A blocked external dependency must not return success merely to keep pipelines green. Rich machine output may distinguish `FAIL`, `BLOCKED`, `NOT RUN` and `N/A`, while process exit status remains fail-closed for required operations.

## Safety guardrails

Destructive commands require explicit environment/resource identity. `local` is not inferred from hostname alone. Commands must never silently target production because credentials happen to be present.

## No hidden manual steps

If a documented supported workflow repeatedly requires an undocumented command, file edit, SQL statement or UI action, that is a tooling/documentation defect and must be incorporated into the governed developer contract rather than becoming tribal knowledge.
