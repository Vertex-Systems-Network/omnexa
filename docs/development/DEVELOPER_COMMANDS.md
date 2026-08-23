# Omnexa Developer Command Contract

Status: **Canonical v1 semantic contract**  
Work package: **P00.08**

This document freezes developer command semantics. P01.12 implemented the first repository-owned CLI surface by composing existing governed operations. P01 is now complete; the CLI remains a kernel capability but creates no future-phase authority.

## Development lifecycle

```text
omnexa dev bootstrap
omnexa dev up
omnexa dev down
omnexa dev status
omnexa dev reset
```

`dev bootstrap` verifies pinned toolchains, prepares non-secret local config, provisions disposable local dependencies, waits for health, installs dependencies and applies migrations. `dev reset` must fail closed if target identity is ambiguous.

These lifecycle names remain the canonical future command contract. P01.12 did **not** pull deployment orchestration or destructive reset/bootstrap behavior forward merely to fill these names; they remain unimplemented until a later governed scope authorizes them.

## Database lifecycle

```text
omnexa db status
omnexa db migrate
omnexa db fresh
omnexa db seed
```

Production migration execution is a separate release/operations concern and must not inherit unsafe local-reset semantics.

P01.12 implemented only the safe `db migrate` boundary. `db status`, `db fresh` and `db seed` remain contract names for later governed implementation and are not silently approximated.

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

Each command maps to the P00.07 quality gate classes and returns a non-zero exit status on failure. Machine-readable output may distinguish exact evidence states. `verify all` is a deterministic aggregate, not a secret CI-only path.

## P01 executable mappings

Completed P01 package verifiers remain mandatory regressions:

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
bash scripts/verify_p01_10.sh
bash scripts/verify_p01_11.sh
P01_12_TEST_DATABASE_URL=<synthetic PostgreSQL test DSN> bash scripts/verify_p01_12.sh
```

`verify_go_quality.sh` is a permanent repository-wide fail-closed Go quality gate. It verifies `gofmt`, pinned `golangci-lint v2.12.2` and pinned `govulncheck v1.7.0` against `./kernel/...`. It does not use `@latest`, silently modify source or convert required findings into warnings.

Canonical GitHub-hosted governance executes repository Go quality and P01 regressions in sequence on `ubuntu-24.04`. P01.04 uses PostgreSQL 18.6, P01.05 uses Valkey 9.1.1 and P01.06 uses `adobe/s3mock:5.1.0` as governed synthetic services. Later completed packages use deterministic in-process primitives/test doubles where appropriate.

## Completed P01.12 developer CLI

P01.12 — Developer CLI Baseline is complete with canonical evidence in `docs/roadmap/evidence/P01.12_COMPLETION_2026-08-23.md` and the satisfied phase exit gate in `docs/governance/P01_EXIT_GATE.md`.

The repository-owned executable surface is intentionally bounded to:

```text
omnexa                         # minimal kernel startup/smoke path
omnexa help
omnexa version
omnexa health
omnexa db migrate
omnexa verify <target>
```

Implementation ownership is `kernel.developer`. The CLI composes completed kernel capabilities and approved repository tooling; it is not a business/domain owner or privileged administration path.

### Help and version

`help` is static deterministic command metadata. `version` renders only the bounded P01.01 build identity (`version` and `commit`) and does not expose timestamps, usernames, working directories or machine metadata.

### Health

`health` composes the P01.02/P01.04 configuration boundary with the P01.08 safe diagnostic projection. Output is JSON containing bounded build/lifecycle/liveness/readiness/dependency state only. Database URLs, provider errors, credentials and RESTRICTED values are not rendered.

### Database migration guard

`db migrate`:

- requires an **explicit** configured environment source rather than silently accepting the default `local` identity;
- requires the governed P01.04 database configuration/resource identity;
- is allowed only for `local`, `ci`, `preview` and `test` developer environments;
- fails closed for `staging` and `production` because production/staging release migration authority is outside this developer CLI baseline;
- invokes the existing P01.04 connection/migration foundation and does not add hidden SQL discovery, reset, seed or destructive mutation semantics;
- emits only the safe environment label on success and never the database URL.

### Verification target mappings

P01.12 does not reimplement quality logic. It invokes fixed allowlisted executable/argument tuples (`python`, `bash`, `go`) without shell-string expansion and maps targets to existing governed operations:

- `verify governance` -> canonical governance/development/operations/freeze/P01 readiness/package validators;
- `verify format`, `verify lint`, `verify static` -> `scripts/verify_go_quality.sh`;
- `verify unit` -> repository kernel Go tests;
- `verify contracts` -> completed contract-oriented P01 verifier subset;
- `verify integration` -> PostgreSQL/cache/storage integration verifier subset;
- `verify migrations` -> P01.04 verifier;
- `verify security` -> repository Go quality plus completed security-relevant P01 verifier subset;
- `verify build` -> kernel build;
- `verify release` -> repository Go quality + module checksum verification + kernel build;
- `verify module-lifecycle` -> explicit `N/A` because P03 module runtime is not active; it is never relabeled PASS;
- `verify all` -> canonical governance validators + repository Go quality + **P01.01-P01.11** regression verifiers + `go mod verify` + kernel build.

`verify all` intentionally does not recursively invoke `scripts/verify_p01_12.sh`. Canonical governance runs real `omnexa verify all` and the focused P01.12 verifier as separate sequential gates, preventing recursive verification while proving local/CI semantic parity.

### P01.12 focused verifier and exit path

`scripts/verify_p01_12.sh` is fail-closed and covers applicable G0-G8 P01.12 risk. It verifies governance/readiness, formatting/static analysis, unit/race behavior, deterministic help/version/health, guarded PostgreSQL migration, secret-safe negatives, command allowlisting, build/package and module checksum integrity.

Canonical GitHub-hosted governance invokes:

```text
omnexa db migrate
omnexa verify all
bash scripts/verify_p01_12.sh
```

with synthetic PostgreSQL/Valkey/S3-compatible providers provisioned by the governed `ubuntu-24.04` job. Final implementation run/job `32629072886 / 97168916985` passed this lane, establishing the P01 exit proof for configuration resolution, kernel start/build, migration, cache/storage contracts, safe observability, readiness/diagnostics, jobs/configuration/audit primitives, developer verification and required quality/security/supply-chain/build checks without hidden manual steps.

The CLI does not authorize production super-admin behavior, P02 tenant/user/role administration, P03 module install/runtime commands, P04+ event/workflow/domain commands, deployment/Kubernetes orchestration, hidden SQL/file mutation, business modules or AI runtime.

## Go result convention

Canonical Go call convention remains:

```go
value, err := operation()
```

The structured `kernel/internal/failure` primitive governs the `error` side when a stable Omnexa failure contract is required. Provider/sink mappings retain lower-level causes privately while public failure projections remain safe.

## Configuration startup contract

Application-level static configuration remains:

```text
OMNEXA_CONFIG_FILE=<explicit JSON file path>
OMNEXA_ENVIRONMENT=local|ci|preview|test|staging|production
```

Static precedence remains:

```text
default -> explicit JSON config file -> environment variable -> explicit in-process test override
```

Unknown `OMNEXA_*` configuration variables fail in strict application mode. Configuration errors identify the key/problem without printing raw sensitive values. P01.10 runtime configuration is a distinct boundary and must not be collapsed back into static startup configuration.

## Future web UI quality commands

When a future authorized browser-UI package implements `docs/quality/WEB_UI_ACCESSIBILITY_PLAN.md`, tooling should expose reproducible operations equivalent to:

```text
omnexa verify web-standards
omnexa verify accessibility
```

A required WAVE API/license dependency that is unavailable is `BLOCKED`, not PASS. These commands are planning semantics only until an authorized browser-UI package exists.

## Exit-code contract

- `0`: operation succeeded;
- non-zero: operation failed or a required prerequisite was unavailable.

A blocked external dependency must not return success merely to keep pipelines green. Rich machine output may distinguish `FAIL`, `BLOCKED`, `NOT RUN` and `N/A`, while process exit status remains fail-closed for required operations.

## Safety guardrails

Destructive commands require explicit environment/resource identity. `local` is not inferred from hostname alone. Commands must never silently target production because credentials happen to be present. CLI convenience does not create authority; privileged future operations must use governed authentication, authorization and audit paths.

At the P01-complete checkpoint both `kernel_code_authorized` and `business_feature_code_authorized` are false. No future command surface may be implemented until a governed phase activation authorizes it.

## No hidden manual steps

If a documented supported workflow repeatedly requires an undocumented command, file edit, SQL statement or UI action, that is a tooling/documentation defect and must be incorporated into the governed developer contract rather than becoming hidden manual steps or tribal knowledge.
