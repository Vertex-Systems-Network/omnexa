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

Each command maps to P00.07 quality gate classes and returns a non-zero exit status on failure. Machine-readable output should eventually be available for CI/release evidence. `verify all` is a deterministic aggregate, not a secret CI-only path.

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
bash scripts/verify_p01_09.sh
```

`verify_go_quality.sh` is a permanent repository-wide fail-closed Go quality gate. It verifies `gofmt`, pinned `golangci-lint v2.12.2` and pinned `govulncheck v1.7.0` against `./kernel/...`. It does not use `@latest`, silently modify source or convert required findings into warnings.

P01.01 covers the pinned Go workspace/build skeleton. P01.02 covers typed static configuration, precedence, race tests, secret-safe failures and startup behavior. P01.03 covers the transport-neutral structured failure contract. P01.04 covers PostgreSQL 18.6. P01.05 covers the Redis-compatible cache foundation against Valkey 9.1.1. P01.06 covers the S3-compatible object/file storage foundation. P01.07 covers structured `log/slog` and OpenTelemetry. P01.08 covers portable health/readiness/diagnostic primitives. P01.09 covers process-local job/scheduler primitives: deterministic registration, UUIDv7 execution IDs, bounded synchronous/queued execution, bounded retry/idempotency, repeatable completion handles, cancellation/deadlines, graceful drain/cancel, one-shot/fixed-interval schedules and safe observability propagation.

Canonical GitHub-hosted governance executes repository Go quality and P01.01 through P01.09 in sequence on `ubuntu-24.04`. P01.04 provisions PostgreSQL 18.6, P01.05 provisions Valkey 9.1.1 and P01.06 provisions `adobe/s3mock:5.1.0` as governed synthetic services. P01.07-P01.09 use deterministic in-process primitives/test doubles where appropriate. The final `omnexa verify ...` CLI remains owned by P01.12.

## Completed P01.05

P01.05 is complete with canonical evidence in `docs/roadmap/evidence/P01.05_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Cache/provider mappings retain lower-level causes privately while public failures remain provider/credential safe.

## Completed P01.06

P01.06 is complete with canonical evidence in `docs/roadmap/evidence/P01.06_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Storage/provider mappings retain lower-level causes privately while public failures remain provider/credential safe; object identity does not imply tenancy or authorization.

## Completed P01.07

P01.07 is complete with canonical evidence in `docs/roadmap/evidence/P01.07_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Observability remains diagnostic infrastructure, does not become a correctness dependency or business source of truth, preserves vendor-neutral exporter boundaries and redacts prohibited secret/classified content.

## Completed P01.08

P01.08 is complete with canonical evidence in `docs/roadmap/evidence/P01.08_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Liveness/readiness remain distinct, required/security-critical readiness fails closed, optional dependency failure may degrade, dependency probes remain bounded/panic-safe and diagnostics remain safe non-authoritative operational evidence.

## Completed P01.09

P01.09 is complete with canonical evidence in `docs/roadmap/evidence/P01.09_COMPLETION_2026-08-22.md`. Its verifier remains a required regression gate. Job/scheduler identity remains non-authoritative; retries remain bounded and idempotency-protected where required; execution/queue capacity is bounded; shutdown retains accepted-work drain/cancel semantics; schedules remain process-local maintenance primitives rather than durable workflow timers.

## Active P01.10

P01.10 — Feature flag & configuration registry is the sole active executable kernel package after the P01.09 closure transition.

The active P01.10 implementation must map package requirements to:

- `verify format` / `verify static` -> pinned toolchain, repository Go quality, `kernel.configuration` dependency/scope boundaries and no P02/business/P01.11+ pull-forward;
- `verify unit` -> typed definition validation, duplicate stable identifiers, deterministic defaults/fallbacks, evaluation behavior and deterministic test provider;
- `verify contracts` / `verify integration` -> provider-failure fallback, version/change metadata, bounded refresh/invalidation and explicit evaluation context behavior;
- `verify security` -> flags cannot grant authority or bypass authorization/data isolation, security controls fail closed, sensitive values remain governed by classification/secrets policy and future scoped inputs do not create P02 identity;
- lifecycle/resilience evidence -> provider outage/refresh/invalidation behavior remains bounded and deterministic;
- `verify build` -> complete kernel package build with canonical dependency metadata;
- completed P01.01-P01.09 regression preservation.

P01.10 may not implement product experimentation/analytics, tenant admin UI, pricing/entitlement/licensing, authorization based solely on flags, business-module flags before their owners exist, P01.11 audit transport, P01.12 developer CLI, P02+ behavior or AI runtime functionality. P01.11 becomes eligible only after P01.10 completion evidence and a separate governed transition.

## Go result convention — completed P01.03

P01.03 did **not** introduce a custom generic `Result[T]` container. The canonical Go call convention remains:

```go
value, err := operation()
```

The structured `kernel/internal/failure` primitive governs the `error` side when a stable Omnexa failure contract is required. This preserves normal Go composition and `errors.Is`/`errors.As` behavior.

Provider mappings retain lower-level causes privately while public failure projections remain provider/credential safe. P01.07 observability, P01.08 diagnostics, P01.09 jobs and P01.10 runtime configuration must reuse the existing safe failure/observability boundaries rather than inventing a second error authority.

## Configuration startup contract

The kernel process accepts only explicit governed static configuration sources. Application-level controls remain:

```text
OMNEXA_CONFIG_FILE=<explicit JSON file path>   # optional; no implicit file discovery
OMNEXA_ENVIRONMENT=local|ci|preview|test|staging|production
```

Completed P01.04 database, P01.05 cache, P01.06 storage and P01.07 observability configuration surfaces remain governed by their implemented typed configuration boundaries. P01.08 and P01.09 introduced no canonical environment-variable surface.

P01.10 is specifically a **runtime flag/configuration registry distinct from P01.02 static environment configuration**. No new P01.10 environment variable, provider credential, flag identifier or runtime configuration key is canonical until implemented and verified in the repository. Sensitive values remain subject to data classification/secrets rules and must not be turned into generic flags.

Configuration precedence for the static P01.02 boundary remains:

```text
default -> explicit JSON config file -> environment variable -> explicit in-process test override
```

Unknown `OMNEXA_*` configuration variables fail in strict application mode. Configuration errors identify the key/problem without printing raw sensitive values. Test overrides are instance-local and are not a substitute for the P01.10 governed runtime registry.

## Future web UI quality commands

When a future authorized browser-UI package implements `docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md`, developer/CI tooling should expose reproducible operations equivalent to:

```text
omnexa verify web-standards
omnexa verify accessibility
```

A required WAVE API/license dependency that is unavailable is `BLOCKED`, not PASS. These commands are planning semantics only during P01 and do not authorize UI implementation.

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
