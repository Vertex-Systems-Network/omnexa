# Omnexa Program Status

Last reconciled: **2026-08-22**

## Current position

- Program: **Foundation Program**
- Phase: **P00 — Product Constitution & Architecture Freeze**
- Phase state: **active**
- Current work package: **P00.09 — Initial threat model and operational SLO targets**
- Business-feature implementation: **NOT AUTHORIZED YET**
- Kernel implementation: **NOT AUTHORIZED YET**
- P00 progress: **8 / 10 done**

## P00 work packages

| ID | Work package | State |
|---|---|---|
| P00.01 | Repository governance baseline | done |
| P00.02 | Product/domain glossary and naming | done |
| P00.03 | ID/money/time/locale/error conventions | done |
| P00.04 | API contract standard | done |
| P00.05 | Event contract standard | done |
| P00.06 | Security/data classification | done |
| P00.07 | Testing/CI/release standard | done |
| P00.08 | Local developer/repository structure | done |
| P00.09 | Threat model and operational SLO targets | **active** |
| P00.10 | Foundation architecture freeze review | planned |

## P00.08 frozen development baseline

Omnexa remains a **governed monorepo + strict modular monolith** initially. Canonical ownership categories are:

```text
apps/
kernel/
modules/
platform/
shared/
infrastructure/
scripts/
docs/
generated/
```

Directory boundaries express ownership, not automatic service/deployment boundaries. Modules keep private implementation/schema/migrations with their owner; other modules do not import private internals or write private tables.

Default local dependencies are containerized PostgreSQL, Redis-compatible cache, NATS/JetStream and S3-compatible object storage, plus a mail sink when needed. Kubernetes is **not** required for the default developer loop.

Local development and CI share the same canonical verification semantics. Future tooling reserves `omnexa dev`, `omnexa db`, `omnexa verify` and `omnexa module` command families. `omnexa verify all` must be runnable locally; CI is an orchestrator rather than the sole quality environment.

Toolchains/dependencies are repository-pinned. Unversioned global tools and multiple conflicting JS package-manager lockfiles are prohibited. Generated artifacts are reproducible derivatives, never alternate sources of truth.

Configuration is deterministic and separated from secrets. Production secrets/data are not required or allowed by default in ordinary local development. Synthetic deterministic fixtures should include multiple tenants/scopes for later negative-isolation testing.

Linux is the canonical backend execution environment; macOS is supported where tooling permits; Windows backend development prefers **WSL2**, while native Windows remains a distinct POS/edge certification target when needed.

Normative P00.08 evidence:

- `docs/development/REPOSITORY_STRUCTURE.md`
- `docs/development/LOCAL_DEVELOPMENT.md`
- `docs/development/TOOLCHAIN_STANDARD.md`
- `docs/development/CONFIGURATION_STANDARD.md`
- `docs/development/DEVELOPER_COMMANDS.md`
- `docs/contracts/development/workspace.schema.json`
- `docs/adr/ADR-0008-repository-local-development-baseline.md`
- `scripts/validate_development_spec.py`

## Existing frozen quality/security contracts

P00.03–P00.07 remain binding: UUIDv7/exact-money/time/locale/error semantics; governed HTTP and event contracts; tenant/security/data-classification rules; and provider-independent G0–G8 quality/release gates.

## Temporary GitHub Actions exception

GitHub Actions allowance is exhausted/disabled. ADR-0006 permits only P00 documentation/specification manual evidence while hosted execution is unavailable. Hosted CI is **BLOCKED / NOT RUN**, never PASS. The exception expires before any P01 implementation merge, at P00 exit, or sooner if Actions returns.

## Outstanding hosted/business governance

1. Issue #3 — `main` branch/ruleset protection still requires hosted admin configuration.
2. Issue #14 — hosted Actions quota/runner remains unavailable.
3. Issue #4 — licensing/IP/trademark strategy remains unresolved before external distribution/public launch.

## Execution lock

Until P00 exits, do not begin kernel, database model, CRM, ERP, commerce, POS, website builder, payments or AI product implementation.

P00.09 is now the only authorized work package.