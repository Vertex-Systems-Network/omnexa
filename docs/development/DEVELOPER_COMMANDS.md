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

Completed P01 packages retain their verification wrappers as mandatory regressions:

```text
bash scripts/verify_p01_01.sh
bash scripts/verify_p01_02.sh
bash scripts/verify_p01_03.sh
P01_04_TEST_DATABASE_URL=<synthetic PostgreSQL test DSN> bash scripts/verify_p01_04.sh
```

P01.01 covers the pinned Go workspace/build skeleton. P01.02 covers typed configuration, precedence, race tests, secret-safe configuration failures and startup behavior. P01.03 covers the transport-neutral structured failure contract, private-cause/public-redaction behavior and error/result conventions. P01.04 covers the PostgreSQL connection/migration substrate, bounded provider behavior, transaction semantics, fresh/upgrade/failure migrations, immutable ledger/drift checks and advisory-lock coordination against PostgreSQL 18.6.

Canonical GitHub-hosted governance currently executes P01.01 through P01.04 in sequence on `ubuntu-24.04`. P01.04 additionally provisions the governed PostgreSQL 18.6 synthetic test service. The final `omnexa verify ...` CLI remains owned by P01.12.

### Active P01.05

P01.05 — Cache abstraction is the sole active executable kernel package. Its implementation PR must add a fail-closed `scripts/verify_p01_05.sh` and a governed Redis-compatible synthetic test service before P01.05 can complete.

The active P01.05 verifier must map the package requirements to:

- `verify format` / `verify static` -> pinned toolchain, formatting/lint/static checks and cache-package dependency/scope boundaries;
- `verify unit` -> deterministic key namespace/version, TTL, serialization and miss/error behavior;
- `verify integration` -> Redis-compatible get/set/delete and any justified atomic primitive against a real synthetic provider;
- `verify security` -> provider-secret redaction and no business/session/tenant/later-capability pull-forward;
- lifecycle/resilience evidence -> provider flush/restart/eviction/outage proving cache is non-authoritative;
- `verify build` -> complete kernel package build with pinned cache dependency metadata.

Until that active verifier/service exists and passes on the canonical GitHub-hosted lane, P01.05 cannot be marked done.

## Go result convention — completed P01.03

P01.03 did **not** introduce a custom generic `Result[T]` container. The canonical Go call convention remains:

```go
value, err := operation()
```

The structured `kernel/internal/failure` primitive governs the `error` side when a stable Omnexa failure contract is required. This preserves normal Go composition, `errors.Is`/`errors.As` behavior and avoids forcing a second result abstraction throughout the kernel.

P01.04 database/provider mappings retain lower-level causes privately while public failure projections remain provider/credential safe. P01.05 cache/provider mappings must follow the same safe public/private cause boundary. Transport adapters and telemetry emission remain owned by their later packages.

## Configuration startup contract — completed P01.02 / P01.04 database extension

The kernel process accepts only explicit governed configuration sources. Application-level controls remain:

```text
OMNEXA_CONFIG_FILE=<explicit JSON file path>   # optional; no implicit file discovery
OMNEXA_ENVIRONMENT=local|ci|preview|test|staging|production
```

P01.04 composes these database settings only at the database boundary; it does not make the existing kernel startup path database-dependent:

```text
OMNEXA_DATABASE_URL=<RESTRICTED PostgreSQL DSN>
OMNEXA_DATABASE_CONNECT_TIMEOUT=5s
OMNEXA_DATABASE_MAX_CONNECTIONS=10
OMNEXA_DATABASE_MIN_CONNECTIONS=0
OMNEXA_DATABASE_MAX_CONNECTION_LIFETIME=30m
OMNEXA_DATABASE_MAX_CONNECTION_IDLE_TIME=5m
```

`OMNEXA_DATABASE_URL` is RESTRICTED and must never be printed in logs, failures or artifacts. The remaining values are bounded by the P01.04 adapter before a pool is constructed.

P01.05 may add only cache-specific configuration keys justified by its active specification. Cache credentials/endpoints are sensitive and must follow the same no-leak rule.

Configuration precedence remains:

```text
default -> explicit JSON config file -> environment variable -> explicit in-process test override
```

Unknown `OMNEXA_*` configuration variables fail in strict application mode. Configuration errors identify the key/problem without printing raw sensitive values. Test overrides are instance-local and are not a global runtime feature-flag/config registry.

## Future web UI quality commands

When a future authorized browser-UI package implements `docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md`, developer/CI tooling should expose reproducible operations equivalent to:

```text
omnexa verify web-standards
omnexa verify accessibility
```

Their semantics must cover rendered-output W3C validation, WAVE evaluation where the owning package requires it, and manual-evidence hooks for keyboard/focus/screen-reader/zoom-reflow checks. A required WAVE API/license dependency that is unavailable is `BLOCKED`, not PASS. These commands are planning semantics only during P01 and do not authorize UI implementation.

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
