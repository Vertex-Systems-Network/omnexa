# Omnexa Developer Command Contract

Status: **Canonical v1 semantic contract**  
Work package: **P00.08**

This document freezes developer command semantics. P01.12 is authorized to implement the first repository-owned CLI surface, but local development and CI must continue to call the same governed underlying operations.

## Development lifecycle

```text
omnexa dev bootstrap
omnexa dev up
omnexa dev down
omnexa dev status
omnexa dev reset
```

`dev bootstrap` verifies pinned toolchains, prepares non-secret local config, provisions disposable local dependencies, waits for health, installs dependencies and applies migrations. `dev reset` must fail closed if target identity is ambiguous.

## Database lifecycle

```text
omnexa db status
omnexa db migrate
omnexa db fresh
omnexa db seed
```

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
```

`verify_go_quality.sh` is a permanent repository-wide fail-closed Go quality gate. It verifies `gofmt`, pinned `golangci-lint v2.12.2` and pinned `govulncheck v1.7.0` against `./kernel/...`. It does not use `@latest`, silently modify source or convert required findings into warnings.

Canonical GitHub-hosted governance executes repository Go quality and completed P01 regressions in sequence on `ubuntu-24.04`. P01.04 provisions PostgreSQL 18.6, P01.05 provisions Valkey 9.1.1 and P01.06 provisions `adobe/s3mock:5.1.0` as governed synthetic services. Later completed packages use deterministic in-process primitives/test doubles where appropriate.

## Completed P01.11

P01.11 is complete with canonical evidence in `docs/roadmap/evidence/P01.11_COMPLETION_2026-08-23.md`. Its verifier remains a required regression gate. Protected audit remains separate from ordinary logs; required-audit failures do not silently succeed; generic audit rejects secret/auth/key/payment-sensitive metadata; audit write does not imply read/export authority; and descriptive actor/scope metadata does not grant identity, tenancy or authorization authority.

## Active P01.12

P01.12 — Developer CLI Baseline is the sole active executable kernel package and owns the P01 exit proof.

The active P01.12 implementation must map package requirements to the existing semantic command contract rather than create a competing verification path:

- `help` / `version` -> deterministic repository/build identity with no machine-secret leakage;
- `verify governance` / `verify all` -> compose the canonical governance/readiness/package and repository Go-quality gates with fail-closed non-zero exit behavior;
- `verify format`, `verify lint`, `verify static`, `verify unit`, `verify contracts`, `verify integration`, `verify migrations`, `verify security`, `verify build`, `verify release` -> invoke governed underlying operations without silently weakening or reimplementing their semantics;
- configuration-facing commands -> compose P01.02 startup configuration rules without turning the CLI into a secrets store;
- migration-facing commands -> compose P01.04 boundaries with explicit environment/resource identity and no hidden production reset semantics;
- diagnostic/health commands -> compose P01.08 safe projections without leaking provider errors, credentials or RESTRICTED payloads;
- P01 exit proof -> from a clean environment prove configuration, build/start, fresh migration, cache/storage provider contracts, safe logs/telemetry, readiness/diagnostics, job/configuration/audit primitive tests, canonical developer verification and required security/supply-chain/build gates reproducibly.

P01.12 must add focused positive/negative tests and a fail-closed P01.12 verifier during its implementation PR. This closure transition does not create that verifier or any P01.12 runtime code.

P01.12 may not implement production super-admin CLI, P02 tenant/user/role administration, P03 module install/runtime commands, P04+ event/workflow/domain commands, deployment/Kubernetes orchestration, hidden SQL/file mutation, business modules or AI runtime.

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

A required WAVE API/license dependency that is unavailable is `BLOCKED`, not PASS. These commands are planning semantics only during P01.

## Exit-code contract

- `0`: operation succeeded;
- non-zero: operation failed or a required prerequisite was unavailable.

A blocked external dependency must not return success merely to keep pipelines green. Rich machine output may distinguish `FAIL`, `BLOCKED`, `NOT RUN` and `N/A`, while process exit status remains fail-closed for required operations.

## Safety guardrails

Destructive commands require explicit environment/resource identity. `local` is not inferred from hostname alone. Commands must never silently target production because credentials happen to be present. CLI convenience does not create authority; privileged future operations must use governed authentication, authorization and audit paths.

## No hidden manual steps

If a documented supported workflow repeatedly requires an undocumented command, file edit, SQL statement or UI action, that is a tooling/documentation defect and must be incorporated into the governed developer contract rather than becoming hidden manual steps or tribal knowledge.
