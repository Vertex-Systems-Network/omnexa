# Omnexa Local Development Standard

Status: **Canonical v1**  
Work package: **P00.08**

A contributor must be able to create a clean, reproducible local Omnexa development environment without depending on GitHub Actions or hidden machine state.

## Principles

1. Local verification is first-class; CI executes the same canonical commands.
2. A clean checkout plus documented prerequisites/bootstrap is sufficient.
3. Local infrastructure is disposable and reproducible.
4. Production secrets/data are not required for ordinary development.
5. Kubernetes is not required for the default developer loop.
6. Linux is the canonical execution environment; macOS and Windows/WSL2 are supported developer hosts where tooling permits.
7. Domain/module ownership remains identical locally and in production.

## Default local dependency stack

The initial developer profile uses containerized dependencies with pinned versions:

```text
PostgreSQL
Redis-compatible cache
NATS + JetStream
S3-compatible object storage emulator/server
mail sink / test SMTP when needed
```

Additional services are introduced only by the phase/module that requires them. Search/analytics/vector systems are not mandatory day-one local dependencies unless their active phase makes them necessary.

## Bootstrap contract

Future executable repository tooling must expose one documented bootstrap entrypoint equivalent to:

```text
omnexa dev bootstrap
```

or a repository-native wrapper with the same semantics. Bootstrap must:

- verify toolchain prerequisites;
- create/copy non-secret local configuration from examples;
- start required local infrastructure;
- wait for health readiness;
- create/reset development databases only with explicit local scope;
- run migrations;
- install deterministic reference/development data when requested;
- report local service URLs/health;
- not require manual database edits.

Bootstrap is idempotent where practical and fails with actionable diagnostics.

## Canonical developer commands

P00.07 defines semantic command families. P00.08 reserves these developer operations for future CLI/task wrappers:

```text
omnexa dev bootstrap
omnexa dev up
omnexa dev down
omnexa dev status
omnexa dev reset
omnexa db migrate
omnexa db fresh
omnexa db seed
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

Exact binary/wrapper syntax may evolve before implementation, but semantics must remain available locally and in CI.

## Database workflow

Supported local paths:

```text
fresh database -> migrations -> reference/dev seed -> boot
existing supported dev schema -> migrations -> boot
```

`dev reset` is destructive only to explicitly identified local-development resources and must never target production-like environments by inference.

## Developer seed model

Synthetic deterministic fixtures should provide:

- at least two tenants;
- multiple organizations/sub-scopes;
- principals with different roles/relationships;
- enough reference data to exercise cross-tenant negative tests later.

Seeds are not production demo data and never contain real customer secrets/data.

## Network and ports

Local service ports are documented centrally and configurable through non-secret local config. Services must bind to loopback by default when public exposure is unnecessary.

Avoid undocumented host-file/domain requirements. Friendly local domains may be optional convenience layers, not correctness dependencies.

## Cross-platform expectations

### Linux

Canonical native development environment.

### macOS

Supported through the same container/toolchain contracts where upstream runtimes support it.

### Windows

Preferred baseline is WSL2 for backend/platform development so shell/filesystem/container behavior remains close to Linux. Native Windows may be required later for Windows-specific POS/edge certification, but it is a distinct target lane rather than the default backend environment.

Line endings, path separators and executable-bit assumptions must not make the repository unusable across supported hosts.

## No hidden global tools

If a generator/linter/formatter/schema tool is required, it must be version-pinned through repository/toolchain configuration or a documented bootstrap-managed mechanism. "Install latest globally" is not a reproducible prerequisite.

## Offline/limited connectivity

The first bootstrap may require fetching dependencies/images. Once dependencies are cached, ordinary local verification should avoid uncontrolled internet calls. Tests use sandboxes/fakes/containers instead of live third-party production APIs unless explicitly running a certification lane.

## Developer security

- `.env`/local secret files are ignored;
- example config contains placeholders only;
- no production credentials in developer setup;
- local object/database data is treated as non-production synthetic data;
- local services use non-production credentials and are not exposed broadly by default;
- downloaded third-party fixtures must be reviewed/classified before repository use.

## Troubleshooting contract

Bootstrap/verify commands should surface:

- failed prerequisite/version;
- failed service and health endpoint;
- migration failure;
- port conflict;
- config key missing/invalid;
- exact failing quality gate.

A developer must not need undocumented tribal knowledge to recover from a clean-start failure.