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
bash scripts/verify_go_quality.sh
bash scripts/verify_p01_01.sh
bash scripts/verify_p01_02.sh
bash scripts/verify_p01_03.sh
P01_04_TEST_DATABASE_URL=<synthetic PostgreSQL test DSN> bash scripts/verify_p01_04.sh
P01_05_TEST_CACHE_ADDRESS=127.0.0.1:6379 P01_05_TEST_VALKEY_IMAGE=valkey/valkey:9.1.1 bash scripts/verify_p01_05.sh
P01_06_TEST_S3_ENDPOINT=http://127.0.0.1:9090 P01_06_TEST_S3_BUCKET=omnexa-p01-06 P01_06_TEST_S3_IMAGE=adobe/s3mock:5.1.0 bash scripts/verify_p01_06.sh
bash scripts/verify_p01_07.sh
bash scripts/verify_p01_08.sh
```

`verify_go_quality.sh` is a permanent repository-wide fail-closed Go quality gate. It verifies `gofmt`, pinned `golangci-lint v2.12.2` and pinned `govulncheck v1.7.0` against `./kernel/...`. It does not use `@latest`, silently modify source or convert required findings into warnings.

P01.01 covers the pinned Go workspace/build skeleton. P01.02 covers typed configuration, precedence, race tests, secret-safe configuration failures and startup behavior. P01.03 covers the transport-neutral structured failure contract, private-cause/public-redaction behavior and error/result conventions. P01.04 covers the PostgreSQL connection/migration substrate against PostgreSQL 18.6. P01.05 covers the Redis-compatible cache foundation against Valkey 9.1.1, including deterministic keys, bounded TTL/value semantics, serialization, miss/error distinction, provider outage/cancellation, flush non-authority and provider restart/reconnect behavior. P01.06 covers the S3-compatible object/file storage foundation, including deterministic namespaced/versioned keys, bounded streaming, untrusted metadata validation, integrity checks, missing-object semantics, provider timeout/cancellation/unavailability, concurrent integration and provider restart behavior. P01.07 covers the structured `log/slog` and OpenTelemetry baseline, including stable fields, bounded correlation/W3C trace propagation, classification/sensitive-key redaction, isolated trace/metric providers, vendor-neutral exporter injection, deterministic capture and bounded flush/shutdown behavior. P01.08 covers portable health/readiness/diagnostic primitives, including distinct liveness/readiness, criticality-aware dependency results, startup lifecycle state, deterministic ordering, safe diagnostic projection, timeout/cancellation/panic resilience and P01.07 observability integration.

Canonical GitHub-hosted governance executes repository Go quality and P01.01 through P01.08 in sequence on `ubuntu-24.04`. P01.04 provisions PostgreSQL 18.6, P01.05 provisions Valkey 9.1.1 and P01.06 provisions `adobe/s3mock:5.1.0` as governed synthetic services. P01.07 and P01.08 use deterministic in-process test doubles/primitives. The final `omnexa verify ...` CLI remains owned by P01.12.

## Completed P01.05

P01.05 is complete with canonical evidence in `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Cache/provider mappings retain lower-level causes privately while public failures remain provider/credential safe.

## Completed P01.06

P01.06 is complete with canonical evidence in `docs/roadmap/evidence/P01.06_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Storage/provider mappings retain lower-level causes privately while public failures remain provider/credential safe; object identity does not imply tenancy or authorization.

## Completed P01.07

P01.07 is complete with canonical evidence in `docs/roadmap/evidence/P01.07_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Observability remains diagnostic infrastructure, does not become a correctness dependency or business source of truth, preserves vendor-neutral exporter boundaries and redacts prohibited secret/classified content.

## Completed P01.08

P01.08 is complete with canonical evidence in `docs/roadmap/evidence/P01.08_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Liveness/readiness remain distinct, required/security-critical readiness fails closed, optional dependency failure may degrade, dependency probes remain bounded/panic-safe and diagnostics remain safe non-authoritative operational evidence.

## Active P01.09

P01.09 — Job & scheduler primitives is the sole active executable kernel package after the P01.08 closure transition.

The active P01.09 implementation must map package requirements to:

- `verify format` / `verify static` -> pinned toolchain, repository Go quality, jobs package dependency/scope boundaries and no durable-messaging/workflow pull-forward;
- `verify unit` -> deterministic job registration/execution, unknown-job safe failure, retry/backoff bounds, schedule validation and handler-result behavior;
- `verify contracts` / `verify integration` -> explicit idempotency/duplicate-safe handler contract, bounded worker execution, observability/correlation propagation and deterministic in-memory harness;
- `verify security` -> scheduler/job identity grants no authority, future tenant/actor scope is not invented, no business/later-phase behavior is introduced and diagnostics/errors remain safe;
- lifecycle/resilience evidence -> bounded concurrency, deadline/cancellation propagation, bounded graceful shutdown/drain/cancel and no infinite retries;
- `verify build` -> complete kernel package build with canonical dependency metadata;
- completed P01.01-P01.08 regression preservation.

P01.09 may not implement NATS/JetStream durable streams/event consumers, transactional outbox/inbox, distributed workflow timers, business jobs, tenant-context runtime before P02, P01.10 feature registry, P01.11 audit transport, P01.12 developer CLI, P02+ behavior or AI runtime functionality. P01.10 becomes eligible only after P01.09 completion evidence and a separate governed transition.

## Go result convention — completed P01.03

P01.03 did **not** introduce a custom generic `Result[T]` container. The canonical Go call convention remains:

```go
value, err := operation()
```

The structured `kernel/internal/failure` primitive governs the `error` side when a stable Omnexa failure contract is required. This preserves normal Go composition, `errors.Is`/`errors.As` behavior and avoids forcing a second result abstraction throughout the kernel.

P01.04 database/provider, P01.05 cache/provider and P01.06 storage/provider mappings retain lower-level causes privately while public failure projections remain provider/credential safe. P01.07 observability must not expose those private causes or sensitive configuration merely because telemetry is available. P01.08 health/diagnostic output likewise remains safe and does not expose provider internals or sensitive dependency data. P01.09 job/scheduler failures must reuse the existing safe failure/observability boundaries rather than emitting handler payloads or inventing a second error authority.

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

P01.06 storage settings remain the completed boundary extension implemented by `kernel/internal/storage/store.go`:

```text
OMNEXA_STORAGE_ENDPOINT=<sensitive S3-compatible endpoint>
OMNEXA_STORAGE_ACCESS_KEY=<sensitive provider access key>
OMNEXA_STORAGE_SECRET_KEY=<RESTRICTED provider secret>
OMNEXA_STORAGE_REGION=us-east-1
OMNEXA_STORAGE_BUCKET=<governed bucket>
OMNEXA_STORAGE_USE_PATH_STYLE=true
OMNEXA_STORAGE_CONNECT_TIMEOUT=3s
OMNEXA_STORAGE_OPERATION_TIMEOUT=5s
OMNEXA_STORAGE_KEY_PREFIX=omnexa
OMNEXA_STORAGE_MAX_OBJECT_BYTES=1073741824
```

P01.07 observability settings are the completed boundary extension implemented by `kernel/internal/observability/config.go`:

```text
OMNEXA_OBSERVABILITY_ENABLED=true
OMNEXA_OBSERVABILITY_SERVICE_NAME=omnexa-kernel
OMNEXA_OBSERVABILITY_LOG_LEVEL=auto
OMNEXA_OBSERVABILITY_EXPORT_INTERVAL=30s
OMNEXA_OBSERVABILITY_EXPORT_TIMEOUT=3s
OMNEXA_OBSERVABILITY_SHUTDOWN_TIMEOUT=5s
```

`OMNEXA_OBSERVABILITY_LOG_LEVEL=auto` resolves to debug in local/test environments and info otherwise. Export interval is bounded to 1s–10m; export and shutdown timeouts are each bounded to 10ms–30s. The completed P01.07 configuration intentionally contains no exporter credentials or backend-specific configuration.

P01.08 introduced no canonical environment-variable surface. P01.09 may add only job/scheduler configuration justified by `docs/roadmap/work-packages/P01.09.md` and implemented through the existing typed configuration boundary. No P01.09 configuration key or environment variable is canonical until it exists in repository implementation and verification.

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

Future tooling should support explicit validation/generation for OpenAPI contracts, event schemas, data-classification declarations, module manifests, quality evidence records and SDK/generated artifacts. Generation and validation must be reproducible from pinned tool versions.

## Exit-code contract

- `0`: operation succeeded;
- non-zero: operation failed or required prerequisite unavailable.

A blocked external dependency must not return success merely to keep pipelines green. Rich machine output may distinguish `FAIL`, `BLOCKED`, `NOT RUN` and `N/A`, while process exit status remains fail-closed for required operations.

## Safety guardrails

Destructive commands require explicit environment/resource identity. `local` is not inferred from hostname alone. Commands must never silently target production because credentials happen to be present.

## No hidden manual steps

If a documented supported workflow repeatedly requires an undocumented command, file edit, SQL statement or UI action, that is a tooling/documentation defect and must be incorporated into the governed developer contract rather than becoming tribal knowledge.
