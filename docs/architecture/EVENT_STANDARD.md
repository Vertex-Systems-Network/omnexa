# Omnexa Event Contract Standard

Status: **Canonical v1**  
Work package: **P00.05**

This standard defines how domain and platform events are named, published, versioned, transported, consumed, retried and replayed across Omnexa. It is a contract standard, not an authorization shortcut and not a license for modules to expose internal persistence details.

## 1. Core principles

1. An event is an immutable fact that already happened.
2. The producing domain owns the event contract.
3. Events cross module boundaries only through published, versioned contracts.
4. Delivery is assumed to be **at least once** unless a transport explicitly proves a stronger guarantee; consumers therefore must be idempotent.
5. There is no platform-wide global ordering guarantee.
6. Tenant/security context is explicit and producer-derived.
7. Event payloads contain the minimum data consumers need, not a dump of the producer database row.
8. Durable publication and consumption use outbox/inbox patterns where business consistency requires them.
9. Replay is a first-class operational capability and must not create duplicate business side effects.
10. Events are not commands, RPC calls, audit logs or an implicit shared database.

## 2. Event naming

Canonical event types use lower-case dot-separated past-tense facts with an explicit major version:

```text
commerce.order.created.v1
commerce.order.cancelled.v1
inventory.stock.reserved.v1
payments.payment.authorized.v1
finance.invoice.issued.v1
identity.user.invited.v1
```

Rules:

- pattern: `<domain>.<subject>.<past_tense_fact>.v<major>`;
- the name describes a fact, not an instruction;
- producer ownership is visible from the domain prefix;
- public event type strings are stable identifiers;
- semantic meaning of an existing event type must never be silently changed;
- commands such as `create_order`, `reserve_stock`, or `send_email` are not events.

## 3. Canonical event envelope

Omnexa event messages use a **CloudEvents-compatible structured JSON envelope**. Standard CloudEvents attribute names remain unchanged; Omnexa-specific context is carried through named extension attributes.

Canonical example:

```json
{
  "specversion": "1.0",
  "id": "018f47d1-35ca-7d66-9a86-9348c574f162",
  "source": "urn:omnexa:module:commerce.orders",
  "type": "commerce.order.created.v1",
  "subject": "order/018f47cf-8f27-7e43-8967-92f7ce14d436",
  "time": "2026-08-21T17:40:12.438Z",
  "datacontenttype": "application/json",
  "dataschema": "urn:omnexa:event-schema:commerce.order.created:v1",
  "tenantid": "018f47a4-2b15-77d7-b6a5-f4ed1f5c7862",
  "organizationid": "018f47a6-d8ca-7499-a512-a79bf3e37a31",
  "correlationid": "018f47d0-01f8-7c75-8b32-c062119ee86f",
  "causationid": "018f47d0-a642-7bef-92d9-19db524bbd32",
  "traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
  "data": {
    "order_id": "018f47cf-8f27-7e43-8967-92f7ce14d436",
    "customer_id": "018f47b0-3845-7848-86dd-beb4c34f20e2",
    "currency": "AED",
    "grand_total": "1299.95"
  }
}
```

The canonical JSON Schema for the common envelope is `docs/contracts/events/event-envelope.schema.json`.

## 4. Required envelope attributes

All published Omnexa domain events require:

| Attribute | Rule |
|---|---|
| `specversion` | `1.0` for the current CloudEvents-compatible baseline. |
| `id` | Globally unique event ID using the Identifier Standard; UUIDv7 by default. |
| `source` | Stable URI identifying the producing module/domain boundary. |
| `type` | Canonical event name including explicit major version. |
| `time` | UTC event occurrence instant using RFC 3339-compatible representation. |
| `datacontenttype` | Normally `application/json` for canonical JSON payloads. |
| `dataschema` | Stable URI/URN identifying the event payload contract/version. |
| `data` | Event-specific payload conforming to the owning schema. |

`subject` is required when the event concerns an addressable aggregate/entity and must be stable enough to support ordering/diagnostics.

## 5. Omnexa extension attributes

The following extension attributes are canonical where applicable:

| Attribute | Meaning |
|---|---|
| `tenantid` | Tenant isolation context derived by the producer. Required for tenant-owned events. |
| `organizationid` | Organization scope when the fact is organization-bound. |
| `correlationid` | Groups messages/actions belonging to one logical business interaction. |
| `causationid` | ID of the event/command/request that directly caused this event. |
| `traceparent` | Distributed-tracing propagation compatible with W3C Trace Context when available. |
| `actorid` | Stable actor/principal ID when attribution is material and disclosure is allowed. |
| `actortype` | Qualified actor type such as `user`, `service_account`, `system`, `agent`. |
| `subjectsequence` | Optional monotonic sequence for one subject/aggregate when consumers require producer-defined ordering. |

Standard CloudEvents attributes deliberately keep their standard spellings. Event payload object properties inside `data` follow Omnexa canonical `snake_case` naming unless a registered external contract requires otherwise.

## 6. Event identity

- `id` is the immutable identity of one event publication fact.
- Re-delivery/replay of the same logical event preserves the same `id`.
- Re-emitting a genuinely new business fact uses a new event ID.
- Consumers deduplicate on event ID within the owning consumer/inbox scope.
- Do not generate a new ID merely because a broker redelivers a message.

## 7. Producer ownership

The producer that owns the authoritative write model owns:

- event type/name;
- payload schema;
- semantic meaning;
- major version evolution;
- publication conditions;
- producer-side tenant/security context;
- source identifier.

Consumers must not redefine or fork an upstream event contract under the same type name.

If a consumer needs a different durable fact, it creates its own derived event under its own domain only after performing an owned state transition; it must not simply rename another producer's message.

## 8. Payload design

Events should contain enough immutable business context to let intended consumers react without querying producer internals synchronously for every field.

Prefer:

- stable IDs;
- small immutable snapshots required by consumers;
- explicit money/currency and time primitives from P00.03;
- normalized, contract-owned states;
- references to large/sensitive resources where direct payload inclusion is unsafe.

Do not publish:

- ORM/database row dumps;
- internal table/column implementation detail solely because it exists;
- secrets, password hashes, tokens or credentials;
- unbounded blobs/media;
- unnecessary PII;
- mutable derived fields presented as immutable truth without clear semantics.

## 9. Tenant and organization context

For tenant-owned events:

- `tenantid` is mandatory;
- it is taken from trusted producer execution context, not accepted blindly from an untrusted incoming payload;
- consumers propagate the trusted tenant context into owned processing/audit records;
- a missing/malformed tenant context on a tenant-required event is a contract/security failure, not an invitation to infer a tenant globally;
- `organizationid` is included when organization scope affects business meaning.

Event transport context is not permission proof. A consumer performing a new protected action still applies the owning capability/policy rules appropriate to that action.

## 10. Correlation and causation

`correlationid` and `causationid` are distinct:

- correlation links a wider business interaction/process;
- causation points to the immediate triggering message/request/event.

Rules:

- preserve an existing correlation ID across a logical chain;
- create a new correlation ID at a new root interaction;
- set causation ID to the immediate parent message/request identity where one exists;
- do not overwrite correlation identity merely because a new service/module processes the event;
- propagate tracing independently through `traceparent` where available.

## 11. Schema location and registry

Canonical event payload schemas live under:

```text
docs/contracts/events/<domain>/<event>/v<major>.schema.json
```

Example:

```text
docs/contracts/events/commerce/order-created/v1.schema.json
```

The envelope schema is shared; each event payload schema owns only the `data` contract for its event unless a generated combined schema is intentionally produced.

Future runtime registries may publish these schemas, but generated/runtime copies must be traceable back to the canonical contract source.

## 12. Versioning and compatibility

Major version is embedded in the event type and schema identifier.

Within one major version, producers may normally make only backward-compatible additive changes.

Safe examples, subject to validation:

- add an optional field;
- add documentation/description;
- loosen a minimum/maximum only when consumers remain safe;
- add metadata that consumers are required to ignore when unknown.

Potentially breaking changes require a new major version, including:

- remove/rename a field;
- change field type/precision/meaning;
- make an optional field required;
- repurpose an enum/state value;
- change tenant/ownership meaning;
- change event occurrence semantics;
- change retry/replay meaning such that side effects could duplicate.

Enum additions are treated as potentially breaking unless the contract explicitly defines the enum as open and consumers are required to tolerate unknown future values.

## 13. Major-version migration

When publishing `...v2`:

- do not silently stop `v1` while supported consumers still depend on it;
- support a documented migration/deprecation window;
- dual-publish only when the producer can guarantee semantic consistency and operational cost is justified;
- consumers migrate deliberately and prove compatibility;
- final retirement is a controlled contract change.

## 14. Publication consistency and outbox

When a domain state mutation and event publication represent one business transaction, the producer must not rely on an unsafe sequence such as:

```text
commit database
then attempt broker publish
```

without a recovery mechanism.

The baseline pattern is a **transactional outbox** or an equivalent mechanism that provides the same failure guarantees:

1. commit owned domain state and outbox intent atomically in the same authoritative transaction where feasible;
2. publish asynchronously;
3. mark publication state idempotently;
4. retry publication safely;
5. retain enough diagnostics to detect stuck publication.

An equivalent technology may replace an outbox only through evidence/ADR if it preserves consistency guarantees.

## 15. Consumer idempotency and inbox

Consumers assume duplicate delivery is possible.

For business-significant consumers:

- use an inbox/deduplication record or equivalent durable mechanism keyed by event ID + consumer identity;
- do not mark an event processed before owned side effects are safely committed;
- duplicate delivery after success must be a no-op or return the same owned result;
- idempotency retention must cover the broker/replay horizon required by the contract;
- financial/payment/inventory side effects need explicit duplicate tests before release.

## 16. Ordering

Omnexa does **not** guarantee global event order.

A producer may guarantee order only within a documented partition/subject boundary.

Rules:

- consumers must not depend on order across unrelated subjects;
- when subject order matters, partition/key by stable subject identity where the transport supports it;
- `subjectsequence` may be emitted when a monotonic producer-owned sequence is required;
- consumers detecting a sequence gap must follow a documented recovery/reconciliation path rather than guessing missing state;
- timestamps do not replace a sequence guarantee.

## 17. Retries

Transient processing failures use bounded retry policy.

Baseline behavior:

- exponential/backoff strategy with jitter where appropriate;
- bounded attempts/time horizon;
- idempotent handler semantics;
- retry metadata/attempt count observable;
- permanent contract/validation/security failures must not spin indefinitely;
- provider/financial side effects must honor their own idempotency contract before retry.

Exact runtime timings belong to P04/P05 implementation configuration and P00.09 operational SLO work; business semantics remain governed here.

## 18. Dead-letter handling

A message that cannot be safely processed after the defined retry policy moves to a dead-letter/quarantine mechanism with:

- original immutable event;
- consumer identity;
- failure category/code;
- attempt count;
- first/last failure timestamps;
- correlation/trace identifiers where available;
- protected diagnostics sufficient for remediation.

DLQ messages are not automatically discarded or blindly replayed. Operators/workflows must classify the cause, correct it where necessary, then replay through governed tooling.

## 19. Replay

Replay is a supported operational action, not an exceptional hack.

Replay rules:

- preserve original event `id`, `type`, `time`, producer and business payload;
- transport/runtime may attach separate delivery-attempt metadata outside immutable event semantics;
- consumers must remain idempotent;
- replay tools require scoped authorization and audit logging;
- replay ranges/filters must be explicit and bounded;
- side-effecting consumers may require dry-run, approval or suppression policies depending on risk;
- replay must never modify the historical event to make old data look current.

## 20. Commands vs events

Commands request an action; events report completed facts.

Examples:

```text
Command: commerce.order.cancel
Event:   commerce.order.cancelled.v1
```

Rules:

- a consumer cannot treat an event as permission to bypass the owning action policy;
- commands may fail/reject; facts have already occurred;
- do not name request queues/topics with past-tense facts unless they actually carry immutable events;
- workflow steps that require a response should use a governed command/capability contract rather than disguising RPC as an event.

## 21. Event bus subject/topic mapping

Logical event type is transport-independent.

NATS/JetStream-class mapping may use a subject convention such as:

```text
omnexa.events.<domain>.<subject>.<fact>.v<major>
```

but transport routing names are not the canonical semantic identity; the envelope `type` is.

Do not bake broker cluster/stream names into business code contracts.

## 22. Retention and event sourcing

Domain events are not automatically an eternal source-of-truth event store.

- retention is defined by operational, replay, legal and data-classification requirements;
- an event stream used for integration may be compacted/expired according to policy;
- audit history and accounting records follow their own retention requirements;
- adopting event sourcing for a domain requires an explicit ADR and domain-specific correctness model.

Do not infer “events exist” => “the entire platform is event sourced.”

## 23. Sensitive data and classification

Until P00.06 freezes the complete data-classification model:

- secrets/credentials/authentication material are forbidden in event payloads;
- minimize PII and sensitive commercial data;
- prefer stable references over duplication when consumers can safely resolve through governed capabilities;
- event schemas must declare or document sensitive fields once the classification standard exists;
- dead-letter and observability systems must not become uncontrolled copies of sensitive payloads.

P00.06 may add stricter requirements but must not weaken these controls without formal change control.

## 24. Observability

Publication and consumption telemetry should expose, where applicable:

- event type and version;
- event ID;
- source/producer;
- consumer identity;
- tenant/organization context subject to classification policy;
- correlation/causation IDs;
- trace context;
- attempt count;
- processing latency;
- publish/ack/failure/DLQ/replay outcomes.

Logs/traces must not dump full sensitive payloads by default.

## 25. Failure isolation

One consumer failure must not block unrelated consumers indefinitely.

- consumer state/offsets are independently managed;
- poison messages are quarantined according to policy;
- optional modules may be disabled without preventing producer domains from committing their own valid state;
- producer availability must not require synchronous acknowledgement from every downstream consumer.

## 26. Contract testing

Future implementation must include:

- schema validation of produced events;
- compatibility checks for schema evolution;
- duplicate-delivery tests;
- outbox publication recovery tests where applicable;
- inbox/idempotency tests;
- retry/DLQ tests;
- replay tests for side-effecting consumers;
- tenant-context negative tests;
- ordering/sequence tests where an ordering guarantee is declared.

P00.07 will define the complete executable CI/release gates.

## 27. Prohibited patterns

Do not:

- publish unversioned cross-domain events;
- publish commands under event names;
- use events as an authorization bypass;
- make consumers read producer private tables because the event “does not have enough data”;
- assume exactly-once delivery without proof;
- assume global ordering;
- generate a new event ID on redelivery/replay of the same event;
- retry indefinitely;
- silently discard poison messages;
- include secrets/tokens/password material;
- publish ORM row dumps as contracts;
- mutate historical events during replay;
- use timestamps as a substitute for a defined sequence guarantee;
- couple business semantics to NATS/Kafka-specific routing names.

## 28. Change control

Changes to any of the following are architectural and require ADR/change-control reconciliation:

- event naming/versioning semantics;
- canonical envelope shape;
- producer ownership rules;
- tenant-context semantics;
- outbox/inbox consistency model;
- delivery/idempotency assumptions;
- ordering guarantee model;
- retry/DLQ/replay semantics;
- event-sourcing adoption.

Implementation must not diverge first and document later.
