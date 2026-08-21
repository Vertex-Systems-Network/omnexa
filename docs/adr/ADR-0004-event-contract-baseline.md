# ADR-0004 — Event Contract Baseline

Status: **Accepted**  
Date: **2026-08-21**  
Work package: **P00.05**

## Context

Omnexa modules must communicate asynchronously without becoming coupled through database tables, implementation classes or broker-specific behavior. The platform also requires tenant isolation, replay, traceability, safe retries and independently evolvable modules.

Without a frozen event contract baseline, different teams or AI systems could invent incompatible envelopes, event IDs, naming/versioning rules, delivery assumptions and retry semantics. That would make modularity unreliable and create duplicate side effects in finance, payments, inventory and workflow domains.

## Decision

Omnexa adopts the following platform-wide baseline for published domain/platform events:

1. Events are immutable facts and use past-tense, versioned names: `<domain>.<subject>.<fact>.v<major>`.
2. The authoritative producing domain owns the event type, meaning, schema and publication conditions.
3. The canonical event envelope is CloudEvents-compatible structured JSON using `specversion: 1.0` and standard CloudEvents attribute names.
4. Event IDs use Omnexa UUIDv7 identifier semantics.
5. Tenant-owned events carry producer-derived `tenantid`; organization context is explicit where relevant.
6. Correlation, causation and distributed trace context are separate concepts and are propagated explicitly.
7. Cross-domain event delivery is treated as at-least-once by default; consumers must be idempotent.
8. Business-significant publication uses a transactional outbox or an equivalent mechanism with demonstrably equivalent consistency guarantees.
9. Business-significant consumption uses an inbox/deduplication mechanism or equivalent durable idempotency control.
10. Omnexa provides no global event ordering guarantee. Subject-scoped ordering is explicit and may use `subjectsequence` where required.
11. Retries are bounded; poison/permanent failures move to governed dead-letter/quarantine handling.
12. Replay preserves the immutable event identity and business payload; consumer side effects remain protected by idempotency/approval policy.
13. Events are transport-independent. Broker subjects/streams do not become canonical business identifiers.
14. Publishing integration events does not imply that a domain or the entire platform uses event sourcing.
15. Event payloads minimize sensitive data and may never contain secrets/credentials.

The detailed normative rules live in `docs/architecture/EVENT_STANDARD.md`; the common machine-readable envelope lives in `docs/contracts/events/event-envelope.schema.json`.

## Consequences

### Positive

- modules can evolve independently behind stable event contracts;
- duplicate delivery/replay can be handled safely;
- tenant and observability context remain explicit;
- producer ownership remains clear;
- broker technology can evolve without changing business semantics;
- operational failure handling becomes diagnosable and auditable;
- future workflow/agent systems can consume governed facts instead of private database state.

### Costs

- outbox/inbox infrastructure and contract validation add implementation work;
- consumers cannot assume global order or exactly-once delivery;
- version migration may require dual publishing for bounded periods;
- event schema governance becomes a mandatory part of CI/release work.

## Rejected alternatives

### Direct cross-module database reads/writes

Rejected because they violate data ownership and failure isolation.

### Unversioned event names

Rejected because consumers cannot safely evolve independently.

### Exactly-once as a platform assumption

Rejected because end-to-end exactly-once side effects are generally not guaranteed merely by a broker. Omnexa instead requires durable idempotency.

### Global event ordering

Rejected because it imposes unnecessary scalability and availability constraints. Ordering is declared only where business semantics require it.

### Broker-specific event identity

Rejected because NATS/Kafka/other transport details must remain replaceable infrastructure rather than business contracts.

### Platform-wide event sourcing by default

Rejected because integration events and event-sourced aggregates solve different problems. Event sourcing requires explicit domain-specific justification and ADR approval.

## Compatibility

This ADR depends on:

- ADR-0001 platform architecture baseline;
- ADR-0002 foundation identifiers/time/error primitives;
- ADR-0003 HTTP API contract baseline where HTTP-to-event boundaries interact.

Future P00.06 security/data-classification rules may strengthen event handling but must not silently weaken producer ownership, tenant isolation, idempotency or replay guarantees.

## Supersession

Changes to the canonical envelope, naming/versioning model, delivery guarantee baseline, outbox/inbox consistency model, ordering semantics or replay semantics require a superseding ADR and roadmap/state reconciliation before implementation.
