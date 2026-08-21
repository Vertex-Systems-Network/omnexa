# Omnexa Time & Calendar Standard

Status: **Canonical v1**  
Work package: **P00.03**

This standard separates absolute instants, local business dates/times, durations and calendar periods so timezone/DST behavior remains deterministic.

## 1. Absolute instants

An instant is a unique point on the global timeline.

Rules:

- persist absolute instants in UTC;
- PostgreSQL storage uses `timestamptz` for instants;
- API/JSON representation uses RFC 3339 timestamps with an explicit offset, with canonical Omnexa output normalized to `Z` where the contract represents an instant;
- never persist a server-local timezone as implicit meaning;
- system/audit timestamps such as `created_at`/`updated_at` are instants when those fields exist.

Example:

```text
2026-08-21T16:30:45.123Z
```

## 2. Timezone identity

Use IANA timezone identifiers for business/user timezone configuration, for example:

```text
Asia/Karachi
Asia/Dubai
Europe/Istanbul
Europe/London
```

Rules:

- raw numeric offsets such as `+05:00` are not durable timezone identities because DST/rules can change;
- timezone IDs are configuration/reference values, not locale values;
- tenant, organization, legal entity, location and user may each have timezone context where explicitly supported;
- timezone fallback must be deterministic and documented by the owning capability.

## 3. Business date

A business date is a calendar date with no time-of-day or UTC conversion semantics.

Use PostgreSQL `date` and ISO `YYYY-MM-DD` representation.

Examples include invoice date, accounting period date, birthday and contractual effective date where only the civil date matters.

Do not convert a business date to midnight UTC and treat it as equivalent.

## 4. Local date-time

A local date-time without timezone may be used only when the business meaning explicitly depends on a civil clock value whose timezone is carried separately.

Examples:

- store opens at 09:00 in `Asia/Dubai`;
- a recurring payroll cutoff occurs at local midnight;
- an appointment requested for a location's civil time before an absolute instant is resolved.

Persist/contract the timezone identity with the local value whenever interpretation depends on it.

## 5. Scheduling and DST

Recurring schedules must define civil-time semantics rather than repeatedly adding fixed elapsed hours.

Rules:

- recurrence anchored to local time must retain its IANA timezone;
- DST gaps/overlaps require a deterministic policy defined by the scheduling/workflow owner;
- converting recurring `09:00 Europe/London` into a fixed UTC offset is forbidden;
- schedule evaluation must use a maintained timezone database provided by the runtime/environment.

P00.05/P00.08 will refine event/runtime implementation requirements.

## 6. Durations vs calendar periods

A duration is fixed elapsed time. A calendar period is calendar-relative.

Examples:

```text
Duration: 90 seconds
Duration: 24 hours
Calendar period: 1 month
Calendar period: 1 business day
```

Do not model one month as 30 days or one year as 365 days unless the domain explicitly defines that equivalence.

Cross-language contracts must distinguish fixed durations from calendar periods. P00.04 will define the API schema.

## 7. Precision

Default persisted/application timestamp precision must be sufficient for deterministic ordering and audit without inventing nanosecond accuracy the infrastructure cannot preserve.

Baseline:

- contracts may carry fractional seconds;
- PostgreSQL/runtime precision must be documented and round/truncate consistently;
- business correctness must not depend on two unrelated events never sharing the same timestamp;
- use identifiers/sequences/event metadata for tie-breaking when required.

## 8. Clock source

Runtime code must use an injectable/testable clock abstraction for business behavior that depends on "now".

Do not scatter direct wall-clock calls through domain logic when deterministic tests, time travel, retries or replay matter.

The system clock is not a source of business identity.

## 9. Expiry and TTL

Expiry values must distinguish:

- an absolute expiry instant; or
- a duration used to calculate an expiry from a named start instant.

Never persist ambiguous values such as `expires_in = 30` without a unit/contract.

## 10. Historical/timezone-rule changes

Historical instants remain absolute. Local display is derived using the relevant timezone rules.

When a legal/business record requires preservation of the originally presented local time or offset, store that snapshot in addition to the canonical instant rather than rewriting the instant.

## 11. Prohibited patterns

Do not:

- store UTC instants in timezone-less timestamp columns;
- use server timezone as business timezone;
- infer timezone from language or currency;
- represent a business date as midnight UTC;
- represent months/years as fixed seconds;
- rely on timestamps as unique IDs;
- silently discard timezone information from user input;
- assume DST transitions are always one hour.
