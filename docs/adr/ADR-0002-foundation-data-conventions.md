# ADR-0002 — Foundation Data & Contract Conventions

Status: **Accepted**  
Date: **2026-08-21**  
Work package: **P00.03**

## Context

Omnexa will span finance, commerce, POS/offline sync, global tenants, integrations, workflows and AI-driven operations across multiple languages/runtimes. Without one platform convention for identity, money, time, locale and errors, modules would create incompatible primitives and hidden conversion bugs.

P00.02 established canonical terminology and ownership. P00.03 must now freeze cross-platform primitive semantics before APIs, events, schemas or application code exist.

## Decision

### Identifiers

- UUIDv7 is the default canonical identifier for new Omnexa entities and Omnexa-generated request/correlation/event/workflow/job/audit identities.
- PostgreSQL uses native `uuid`.
- Canonical tenant scope is `tenant_id`.
- Human-facing document/business numbers remain separate attributes.
- External IDs remain integration mappings rather than replacing Omnexa identity.

### Money and numeric precision

- Money is exact decimal amount + explicit currency.
- Binary floating point is forbidden for monetary values.
- JSON contracts serialize exact decimal monetary/rate values as strings.
- PostgreSQL general high-precision money baseline is `NUMERIC(38,18)`.
- Fiat codes use ISO 4217 where applicable.
- Rounding and currency conversion are explicit, auditable operations; generic fallback rounding is half-even when no governing business/legal policy overrides it.

### Time

- Absolute instants are persisted in UTC using PostgreSQL `timestamptz`.
- Business dates remain date-only values.
- IANA timezone identifiers carry civil-time meaning.
- Recurring local schedules retain timezone/DST semantics rather than fixed offsets.
- Durations and calendar periods are distinct concepts.

### Locale/regionalization

- UI/content locales use BCP 47 language tags.
- Country identity, locale, timezone and currency are independent configuration/data concepts.
- RTL is a first-class layout/accessibility requirement.
- Translatable content must scale through locale-keyed content/translation models rather than one-column-per-language schemas.

### Errors

- Errors expose stable machine codes plus safe structured problem details.
- Public errors carry request/correlation diagnostics but never stack traces, SQL, secrets or protected internals.
- Machine codes/type identifiers are never localized.
- Retry semantics are explicit and must respect idempotency, especially for financial mutations.

## Consequences

### Positive

- consistent cross-language contracts between Go, TypeScript, Rust and Python;
- safer distributed/offline identity generation;
- deterministic finance/accounting behavior;
- global timezone/locale support without data corruption;
- stable errors usable by UI, SDKs, workflows and AI tools;
- API/event standards can now build on frozen primitives.

### Costs / constraints

- UUIDv7-capable libraries/runtimes must be selected/validated during implementation;
- exact decimal libraries/types are mandatory for financial work;
- providers using integer minor units need explicit adapter conversion;
- time/calendar logic must carry timezone context rather than relying on server defaults;
- API clients must treat decimal strings and stable error codes correctly.

## Rejected alternatives

### Auto-increment integer IDs as platform identity

Rejected because they expose persistence coupling, collide with offline/distributed creation patterns and encourage cross-domain assumptions.

### UUIDv4 everywhere

Rejected as default because UUIDv7 provides a standardized time-ordered identifier suitable for database locality while retaining opaque distributed identity. UUIDv4 may still exist where an explicit contract/compatibility reason requires it.

### Money as floating point

Rejected due to non-decimal representation and reconciliation/accounting risk.

### Money universally as integer minor units

Rejected as the platform canonical representation because calculation precision, exchange rates, allocations, non-standard exponents and asset/payment adapters require richer exact-decimal semantics. Providers may still use minor units at their boundary.

### Store all date/time values as UTC timestamps

Rejected because civil dates, recurring local schedules and business calendars are not equivalent to absolute instants.

### Locale determines country/currency/timezone

Rejected because those relationships are not universally true and would corrupt persisted business meaning.

### Arbitrary error strings

Rejected because they are unstable, non-localizable as machine contracts, hard to automate and unsafe for diagnostics.

## Compliance

Canonical detailed rules live in:

- `docs/architecture/IDENTIFIER_STANDARD.md`
- `docs/architecture/MONEY_STANDARD.md`
- `docs/architecture/TIME_STANDARD.md`
- `docs/architecture/LOCALE_STANDARD.md`
- `docs/architecture/ERROR_STANDARD.md`

P00.04/P00.05 may refine transport/event representations but may not contradict these primitive semantics without an ADR/change-control update.
