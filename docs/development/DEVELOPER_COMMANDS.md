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

Completed P01 packages retain their verification wrappers as mandatory regressions, while the sole active package adds its verifier before final hosted verification:

```text
bash scripts/verify_go_quality.sh
bash scripts/verify_p01_01.sh
bash scripts/verify_p01_02.sh
bash scripts/verify_p01_03.sh
P01_04_TEST_DATABASE_URL=<synthetic PostgreSQL test DSN> bash scripts/verify_p01_04.sh
P01_05_TEST_CACHE_ADDRESS=127.0.0.1:6379 P01_05_TEST_VALKEY_IMAGE=valkey/valkey:9.1.1 bash scripts/verify_p01_05.sh
P01_06_TEST_S3_ENDPOINT=http://127.0.0.1:9090 P01_06_TEST_S3_BUCKET=omnexa-p01-06 P01_06_TEST_S3_IMAGE=adobe/s3mock:5.1.0 bash scripts/verify_p01_06.sh
```

`verify_go_quality.sh` is a permanent repository-wide fail-closed Go quality gate. It verifies `gofmt`, pinned `golangci-lint v2.12.2` and pinned `govulncheck v1.7.0` against `./kernel/...`. It does not use `@latest`, silently modify source or convert required findings into warnings.

P01.01 covers the pinned Go workspace/build skeleton. P01.02 covers typed configuration, precedence, race tests, secret-safe configuration failures and startup behavior. P01.03 covers the transport-neutral structured failure contract, private-cause/public-redaction behavior and error/result conventions. P01.04 covers the PostgreSQL connection/migration substrate against PostgreSQL 18.6. P01.05 covers the Redis-compatible cache foundation against Valkey 9.1.1, including deterministic keys, bounded TTL/value semantics, serialization, miss/error distinction, provider outage/cancellation, flush non-authority and provider restart/reconnect behavior. P01.06 covers the S3-compatible object/file storage foundation, including deterministic keys, bounded metadata/size rules, streamed put/open, head/delete/missing semantics, integrity verification, provider outage/cancellation and restart/reconnect behavior.

Canonical GitHub-hosted governance is configured to execute repository Go quality and P01.01 through P01.06 in sequence on `ubuntu-24.04`. P01.04 provisions PostgreSQL 18.6, P01.05 provisions Valkey 9.1.1 and P01.06 provisions Adobe S3Mock 5.1.0 as governed synthetic services. The final `omnexa verify ...` CLI remains owned by P01.12.

## Runner-deferred implementation workflow

For an already-authorized package, implementation work may be completed before consuming the hosted runner. The normal sequence is:

1. source/tests/docs/verifier implementation;
2. deterministic static/unit/self-review preparation;
3. final executable PR;
4. GitHub-hosted canonical governance verification;
5. fix any discovered defects without weakening checks;
6. merge only on green evidence;
7. immutable completion/state/ledger reconciliation.

This workflow reduces repeated runner usage; it does not permit unverified `PASS`, `done` or protected merge claims.

## Completed P01.05

P01.05 is complete with canonical evidence in `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Cache/provider mappings retain lower-level causes privately while public failures remain provider/credential safe.

## Active P01.06

P01.06 — Object & file storage abstraction is the sole active executable kernel package. The implementation branch now provides fail-closed `scripts/verify_p01_06.sh`, the pinned AWS SDK/S3 client baseline and a governed Adobe S3Mock 5.1.0 synthetic provider configuration. **Canonical hosted P01.06 evidence remains pending until the final implementation PR run; P01.06 is not done before that run passes.**

The active P01.06 verifier maps package requirements to:

- `verify format` / `verify static` -> pinned toolchain/client versions, repository Go quality, formatting/static checks and storage-package dependency/scope boundaries;
- `verify unit` -> deterministic object-key namespace/version behavior, metadata/key validation, streaming length/integrity helpers and secret-safe failure behavior;
- `verify integration` -> S3-compatible put/open/head/delete, missing-object behavior and provider contract tests against Adobe S3Mock 5.1.0;
- streaming evidence -> 8 MiB deterministic synthetic object with bounded read buffers rather than whole-object buffering;
- `verify security` -> provider-secret redaction, path-traversal containment, untrusted metadata validation and no module/cache/database/telemetry coupling;
- lifecycle/resilience evidence -> unavailable provider, timeout/cancellation, deliberate content-integrity mismatch and provider restart/reconnect;
- `verify build` -> complete kernel package build with pinned storage dependency metadata;
- completed P01.01-P01.05 regression preservation through the same governance job.

## Go result convention — completed P01.03

P01.03 did **not** introduce a custom generic `Result[T]` container. The canonical Go call convention remains:

```go
value, err := operation()
```

The structured `kernel/internal/failure` primitive governs the `error` side when a stable Omnexa failure contract is required. This preserves normal Go composition, `errors.Is`/`errors.As` behavior and avoids forcing a second result abstraction throughout the kernel.

P01.04 database/provider mappings and P01.05 cache/provider mappings retain lower-level causes privately while public failure projections remain provider/credential safe. P01.06 storage/provider mappings follow the same boundary. Transport adapters and telemetry emission remain owned by their later packages.

## Configuration startup contract

The kernel process accepts only explicit governed configuration sources. Application-level controls remain:

```text
OMNEXA_CONFIG_FILE=<explicit JSON file path>   # optional; no implicit file discovery
OMNEXA_ENVIRONMENT=local|ci|preview|test|staging|production
```

P01.04 database settings remain a completed boundary extension:

```text
OMNEXA_DATABASE_URL=<RESTRICTED PostgreSQL DSN>
OMNEXA_DATABASE_CONNECT_TIMEOUT=5s
OMNEXA_DATABASE_MAX_CONNECTIONS=10
OMNEXA_DATABASE_MIN_CONNECTIONS=0
OMNEXA_DATABASE_MAX_CONNECTION_LIFETIME=30m
OMNEXA_DATABASE_MAX_CONNECTION_IDLE_TIME=5m
```

P01.05 cache settings remain a completed boundary extension:

```text
OMNEXA_CACHE_ADDRESS=<sensitive Redis-compatible endpoint>
OMNEXA_CACHE_USERNAME=<sensitive provider username, optional>
OMNEXA_CACHE_PASSWORD=<RESTRICTED provider password, optional>
OMNEXA_CACHE_CONNECT_TIMEOUT=3s
OMNEXA_CACHE_OPERATION_TIMEOUT=2s
OMNEXA_CACHE_KEY_PREFIX=omnexa
OMNEXA_CACHE_MAX_VALUE_BYTES=1048576
OMNEXA_CACHE_MAX_TTL=24h
```

P01.06 composes only these storage-specific settings at the storage boundary:

```text
OMNEXA_STORAGE_ENDPOINT=<sensitive S3-compatible endpoint>
OMNEXA_STORAGE_ACCESS_KEY=<RESTRICTED provider access key>
OMNEXA_STORAGE_SECRET_KEY=<RESTRICTED provider secret key>
OMNEXA_STORAGE_REGION=us-east-1
OMNEXA_STORAGE_BUCKET=<provider bucket/container name>
OMNEXA_STORAGE_USE_PATH_STYLE=true
OMNEXA_STORAGE_CONNECT_TIMEOUT=3s
OMNEXA_STORAGE_OPERATION_TIMEOUT=5s
OMNEXA_STORAGE_KEY_PREFIX=omnexa
OMNEXA_STORAGE_MAX_OBJECT_BYTES=1073741824
```

P01.06 validates the endpoint scheme/shape, bucket/region, timeout bounds, namespace prefix and object-size cap before provider construction. Endpoint/access/secret values are marked sensitive and must never be printed in public failures or artifacts. The default object-size bound is 1 GiB and the simple single-object adapter rejects configuration above the 5 GiB S3 PutObject boundary; later multipart behavior requires its own governed scope.

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

`module create` uses the governed module manifest/ownership/naming standard and `docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md`; it may not generate direct cross-module dependencies or restart an already-defined module plan from scratch.

## Contract commands

Future tooling should support explicit validation/generation for OpenAPI contracts, event schemas, data-classification declarations, module manifests, quality evidence records and SDK/generated artifacts. Generation and validation must be reproducible from pinned tool versions.

## Exit-code contract

- `0`: operation succeeded;
- non-zero: operation failed or required prerequisite unavailable.

A blocked external dependency must not return success merely to keep pipelines green. Rich machine output may distinguish `FAIL`, `BLOCKED`, `NOT RUN` and `N/A`, while process exit status remains fail-closed for required operations.

## Safety guardrails

Destructive commands require explicit environment/resource identity. `local` is not inferred from hostname alone. Commands must never silently target production because credentials happen to be present.

## No hidden manual steps

If a documented supported workflow repeatedly requires an undocumented command, file edit, SQL statement or UI action, that is a tooling/documentation defect and must be incorporated into the governed developer contract rather than becoming tribal knowledge.
