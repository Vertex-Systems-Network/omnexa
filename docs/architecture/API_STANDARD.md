# Omnexa HTTP API Contract Standard

Status: **Canonical v1**  
Work package: **P00.04**

This standard defines the HTTP/JSON contract rules for public, partner and cross-boundary Omnexa APIs. It is intentionally independent of internal implementation language/framework.

Normative external baselines:

- HTTP semantics: RFC 9110
- Problem Details for HTTP APIs: RFC 9457
- API description: OpenAPI Specification 3.2.0

P00.03 foundation primitives remain authoritative for identifiers, money, time, locale and errors.

## 1. API architecture principles

1. APIs expose owned capabilities/resources, never database tables or private module internals.
2. Every stable external contract is versioned.
3. Domain ownership is visible in the route/operation identity.
4. HTTP method/status/cache/concurrency semantics must follow HTTP rather than framework convenience.
5. Mutations are tenant-aware, authorization-aware, auditable and idempotency-aware.
6. Additive evolution is preferred; breaking evolution requires a new major contract version.
7. OpenAPI is the machine-readable source of truth for HTTP shape; prose documents architecture/policy.
8. Generated clients/servers/tests may derive from the contract, but generated output never becomes the canonical contract.
9. API contracts must not require consumers to know whether Omnexa is a modular monolith or extracted service.
10. AI tools consume the same governed capabilities/contracts as other authorized clients; there is no hidden AI bypass API.

## 2. Contract classes

Omnexa distinguishes:

- **Public API** — documented contract for external customers/developers.
- **Partner API** — external contract scoped to approved partners/connectors.
- **Internal cross-domain API** — versioned contract consumed across ownership boundaries inside Omnexa.
- **Private implementation endpoint** — implementation detail within one owning boundary; never consumed as a stable cross-domain contract.

Changing an endpoint from private to cross-domain/public makes it a governed stable contract.

## 3. Base path and major version

Canonical stable HTTP route form:

```text
/api/v{major}/{domain}/{resources}
```

Examples:

```text
/api/v1/crm/leads
/api/v1/commerce/orders
/api/v1/finance/invoices
/api/v1/inventory/stock-items
```

Rules:

- public/cross-boundary contract major version is in the path;
- `v1`, `v2` represent breaking contract generations, not deployment versions;
- minor/patch-compatible evolution occurs inside the same major path;
- do not expose service names, database schema names, framework controllers or deployment topology in URLs;
- domain segment uses canonical ownership/naming terminology;
- resource path segments are lowercase plural kebab-case nouns (`price-lists`, `service-accounts`).

## 4. Operation identity

Every OpenAPI operation has a globally unique stable `operationId`.

Canonical pattern:

```text
<domain>_<resource>_<action>
```

Examples:

```text
crm_lead_list
crm_lead_get
crm_lead_create
commerce_order_cancel
payments_payment_refund
```

Operation IDs are SDK/code-generation identities and must not be casually renamed.

## 5. JSON naming and representation

JSON property names use lowercase `snake_case`.

Examples:

```json
{
  "id": "018f47a6-7b7e-7c14-a847-0af56b2a44fe",
  "tenant_id": "018f47a6-7b7e-7c14-a847-0af56b2a44fe",
  "created_at": "2026-08-21T16:30:45.123Z"
}
```

All P00.03 rules apply:

- UUIDv7 identifiers are strings;
- exact decimals are strings;
- money always carries currency;
- instants are RFC 3339 timestamps with explicit offset/canonical `Z` output;
- business dates use `YYYY-MM-DD`;
- locale tags use BCP 47;
- country/currency/timezone remain separate concepts.

## 6. Success response envelope

Except where HTTP semantics deliberately require an empty body, stable Omnexa JSON success responses use an envelope.

Single resource/result:

```json
{
  "data": {
    "id": "018f47a6-7b7e-7c14-a847-0af56b2a44fe"
  },
  "meta": {}
}
```

Collection:

```json
{
  "data": [],
  "page": {
    "next_cursor": null,
    "has_more": false
  },
  "meta": {}
}
```

Rules:

- `data` contains the primary result;
- `meta` is optional contract metadata, not a dumping ground for business fields;
- list pagination metadata uses `page`;
- clients must ignore unknown additive envelope metadata fields unless the contract explicitly marks them closed;
- `204 No Content` has no response body.

## 7. Resource methods

Canonical mapping:

- `GET /resources` — list/search within documented filters.
- `POST /resources` — create where resource creation is meaningful.
- `GET /resources/{id}` — retrieve.
- `PATCH /resources/{id}` — partial update where field mutation is valid.
- `PUT /resources/{id}` — full replacement only when the resource semantics genuinely support replacement; not a default update method.
- `DELETE /resources/{id}` — only when deletion is a valid owned business operation; retention/archive/purge semantics must not be hidden behind generic DELETE.

Business lifecycle commands are explicit actions, not fake field updates.

Example:

```text
POST /api/v1/commerce/orders/{id}/actions/cancel
POST /api/v1/finance/invoices/{id}/actions/issue
POST /api/v1/payments/payments/{id}/actions/refund
```

Do not model `cancel`, `approve`, `issue`, `refund`, `close`, `post` or similar protected transitions as arbitrary status PATCHes when domain rules/actions exist.

## 8. HTTP method semantics

- `GET`/`HEAD` are safe and must not cause business mutations.
- `PUT`/`DELETE` use idempotent HTTP semantics when those methods are exposed.
- `POST` may be non-idempotent at HTTP-method level, so protected/retryable mutations use the Omnexa idempotency contract.
- `PATCH` behavior and accepted media type must be explicit per operation.
- Framework limitations do not justify violating method semantics.

## 9. Creation semantics

Successful synchronous creation normally returns:

- `201 Created`;
- the created representation or documented minimal result;
- `Location` pointing to the canonical resource URL when meaningful.

If creation/operation is accepted for asynchronous execution, return `202 Accepted` with a stable operation/status reference or `Location` defined by the owning contract. `202` must not imply eventual success.

## 10. Error responses

HTTP API errors use RFC 9457-compatible `application/problem+json` and the canonical `ERROR_STANDARD.md` extensions.

Required Omnexa machine fields for material API errors:

```text
type
status
code
request_id
```

Optional/conditional fields include:

```text
title
detail
instance
trace_id
retryable
retry_after_ms
violations
```

Rules:

- machine `code`/`type` are stable and not localized;
- `title`/`detail` may be localized;
- status code and problem body must not contradict each other;
- public errors never expose stack traces, SQL, secrets or protected implementation internals;
- authorization policy may intentionally map protected resource existence to 404 rather than revealing 403.

## 11. Canonical HTTP status guidance

Common baseline:

- `200` successful read/update/action with body;
- `201` synchronous resource creation;
- `202` accepted asynchronous processing;
- `204` successful operation with intentionally no body;
- `304` conditional cache result where applicable;
- `400` malformed request/unsupported syntax;
- `401` authentication required/invalid;
- `403` authorization denied where disclosure is permitted;
- `404` resource/capability not found or deliberately non-disclosed;
- `409` domain state, duplicate or idempotency conflict;
- `412` failed optimistic concurrency/precondition;
- `415` unsupported media type;
- `422` semantically valid request with validation errors;
- `429` rate limited;
- `500` unexpected internal failure;
- `502` upstream/dependency gateway failure where semantically appropriate;
- `503` temporarily unavailable;
- `504` upstream timeout.

Do not return `200` for a failed protected operation.

## 12. Validation

Contract syntax validation and domain/business validation are distinct.

Validation responses:

- use top-level code `validation.failed`;
- use structured `violations[]`;
- use contract field paths, not DB/internal field names;
- never echo secrets or unsafe submitted values;
- reject unknown fields when the specific schema is intentionally closed; otherwise compatibility policy controls additive fields.

## 13. Tenant and organization context

Every tenant-owned operation resolves one authorized tenant context before business logic executes.

Rules:

- a credential/session scoped to one tenant resolves that tenant without trusting arbitrary client data;
- a principal authorized for multiple tenants may select a tenant using a canonical request context defined by the authentication/gateway layer;
- any client-supplied tenant/organization selector is untrusted input and must be authorized server-side;
- resource body `tenant_id` must never be accepted as authority to cross tenant boundaries;
- organization/legal-entity/branch scope is separately validated when relevant;
- query filters cannot be used to escape resolved tenant scope.

P00.06 will freeze authentication/authorization transport/security details.

## 14. Idempotency

Protected retriable mutations use the `Idempotency-Key` request header where the operation contract requires it.

Baseline behavior:

- key is opaque and client-generated or SDK-generated;
- canonical scope includes resolved tenant, authenticated client/principal, operation identity and key;
- server stores a request fingerprint and protected result/outcome for the retention window;
- same key + equivalent protected request returns the recorded/equivalent result without repeating the mutation;
- same key + materially different protected request fails with `409` and a stable idempotency-conflict code;
- default retention target is at least 24 hours unless a domain contract defines a longer period;
- payment/financial mutations may require stronger/longer retention;
- idempotency storage/claim must be atomic enough to protect concurrent duplicates.

Do not implement idempotency by simply checking whether a similarly shaped business row exists.

## 15. Optimistic concurrency

Mutable resources that can suffer lost updates should expose HTTP entity versioning with `ETag`.

Baseline:

- retrieval may return `ETag`;
- protected update/delete/action contracts that require current-state matching accept `If-Match`;
- mismatch returns `412 Precondition Failed` with a stable conflict/precondition code;
- clients must not use timestamps as concurrency tokens unless the owning contract proves the semantics;
- domain actions still re-check business invariants inside the transaction.

## 16. Pagination

Cursor pagination is the default for scalable collections.

Canonical query parameters:

```text
page_size
page_cursor
```

Canonical response:

```json
{
  "data": [],
  "page": {
    "next_cursor": "opaque-token-or-null",
    "has_more": true
  }
}
```

Rules:

- cursors are opaque to clients;
- cursor contents/encoding are not a public contract;
- default/max `page_size` is operation-specific and documented;
- ordering used by the cursor must be deterministic;
- page tokens must be integrity-protected when tampering could alter query/scope semantics;
- cursor must not allow a caller to escape tenant/filter authorization scope;
- offset pagination may be explicitly used for small/admin/reference collections where stable scale requirements justify it, but it is not the platform default.

## 17. Filtering, search and sorting

Filters are allowlisted per operation.

Canonical filtering syntax for simple filters:

```text
filter[status]=active
filter[customer_id]=018f...
```

Canonical sort parameter:

```text
sort=-created_at,name
```

Rules:

- `-field` means descending;
- undocumented DB columns are never automatically filterable/sortable;
- free-form SQL/operator injection is forbidden;
- filter/sort authorization applies to sensitive fields;
- full-text/search endpoints may define an explicit `q` parameter with owned semantics;
- query complexity/limits must be bounded.

## 18. Sparse fields and related resources

Sparse fieldsets/includes are optional capabilities, not automatic behavior.

When exposed, use explicit allowlisted semantics documented by the operation. A request must not use includes/fields to access data the caller cannot otherwise read.

Avoid building a universal arbitrary relationship-expansion language that bypasses module ownership.

## 19. Request/response size and bulk operations

- request/response limits are explicit platform/gateway policy;
- domains define bounded bulk operations when business semantics require them;
- no generic endpoint may execute arbitrary heterogeneous commands in one unauditable batch;
- bulk mutation results identify per-item success/failure where partial behavior is allowed;
- all-or-nothing vs partial semantics must be explicit.

## 20. Content types

Baseline JSON content:

```text
application/json
application/problem+json
```

PATCH media type is declared per operation. File/media upload/download contracts may use multipart/binary/presigned mechanisms explicitly documented by the owning module.

Do not base64-embed large files in ordinary JSON by default.

## 21. Request/correlation tracing

- server generates/propagates canonical `request_id` per Identifier Standard;
- W3C trace context/OpenTelemetry trace propagation may coexist with request IDs;
- externally supplied correlation/trace headers are untrusted and normalized/validated by the gateway;
- response/support diagnostics expose only safe identifiers.

## 22. Caching and conditional requests

Caching is explicit.

- sensitive/tenant data must not become publicly cacheable through defaults;
- use standard HTTP cache directives/validators where valid;
- `ETag` may serve caching and optimistic concurrency semantics when documented;
- mutation endpoints must not be cached as successful reads;
- gateway/CDN caching policy must respect authorization/tenant variance.

## 23. Authentication and credentials boundary

P00.06 owns the security baseline. Until then:

- credentials/tokens must never appear in query strings;
- authentication is enforced before protected capability execution;
- APIs never trust role names or client-supplied authorization claims without server verification;
- service accounts/API clients are principals and are scoped/authorized like other identities;
- API design must permit future OAuth/OIDC/service-account strategies without changing resource semantics.

## 24. Rate limiting and quotas

When limits apply:

- return `429 Too Many Requests`;
- include standard `Retry-After` when meaningful;
- stable error code identifies limit class;
- quota/limit identity may be tenant, principal, client, capability or a combination;
- do not expose sensitive infrastructure capacity details;
- rate limiting must not replace authorization/business quotas.

## 25. Asynchronous operations

Long-running work must not hold HTTP requests open indefinitely.

When using `202 Accepted`, the contract defines how to inspect outcome/status. The async operation must carry tenant/principal/correlation/audit context and have explicit terminal states.

A generic platform operation resource may be introduced only by its owning later platform capability; domains must not invent incompatible async envelopes meanwhile.

## 26. Webhooks

Outbound webhooks are an HTTP delivery of versioned events/notifications, not a second unversioned event model.

P00.05 defines event envelope/schema semantics. P00.06 defines signature/secret security requirements. HTTP webhook contracts must additionally define:

- destination/registration ownership;
- content type;
- delivery ID;
- retry/backoff behavior;
- timeout;
- success status range;
- replay/duplicate expectations;
- signature/version headers when security standard is frozen.

Consumers must assume duplicate delivery unless a specific contract proves otherwise.

## 27. Compatibility rules

Normally backward-compatible within a major version:

- adding an optional response field;
- adding a new endpoint/operation;
- adding an optional request field with unchanged default semantics;
- adding a new error code for a newly distinguishable failure where existing generic handling remains valid.

Potentially/breaking and requires explicit compatibility review/versioning:

- removing/renaming fields/endpoints;
- changing field type/meaning/format;
- making optional input required;
- narrowing accepted values/ranges;
- changing identifier/money/time semantics;
- changing pagination token semantics in a way that invalidates active clients without transition;
- changing idempotency/concurrency guarantees;
- changing authorization visibility behavior;
- changing enum behavior in a way strict clients cannot tolerate;
- reusing an error code for a different meaning.

Consumers should tolerate unknown additive response fields. SDK enum generation must support forward-compatible unknown values where the domain is extensible.

## 28. Deprecation and sunset

Deprecation is deliberate and observable.

A deprecated operation/field must have:

- contract documentation marking it deprecated;
- replacement/migration guidance where applicable;
- announced support window appropriate to customer impact;
- telemetry where feasible to measure remaining use;
- a removal plan that requires a new major version when removal is breaking.

Do not remove a stable API because internal code stopped using it.

## 29. OpenAPI source of truth

Canonical HTTP API descriptions use **OpenAPI 3.2.0**.

Rules:

- OpenAPI documents live under a governed contracts tree (baseline template under `docs/contracts/http/` during P00);
- `operationId` is stable and unique;
- schemas reuse P00.03 primitives/components;
- examples must be valid against the intended schema;
- security requirements are explicitly declared once P00.06 freezes schemes;
- generated documentation/SDKs come from the canonical description;
- hand-written runtime routes must not silently diverge from the contract;
- P00.07 will define lint/compatibility/diff gates and toolchain pinning.

If a required production tool does not yet support OpenAPI 3.2.0, compatibility output may be generated from the canonical contract or an ADR may temporarily pin an older description profile. Tooling limitation must not silently redefine API semantics.

## 30. API review checklist

Before a new/changed stable operation is accepted, verify:

1. owning domain/capability is unambiguous;
2. route/version/operationId follow this standard;
3. tenant/org scope is explicit and server-authorized;
4. request/response uses P00.03 primitives;
5. errors use canonical problem semantics;
6. mutation idempotency requirements are explicit;
7. lost-update/concurrency behavior is explicit;
8. pagination/filter/sort behavior is bounded and authorized;
9. sensitive fields do not leak through includes/errors/logging;
10. backward compatibility is classified;
11. OpenAPI contract and tests/implementation are reconciled;
12. audit/observability requirements are stated for protected mutations.

## 31. Prohibited API patterns

Do not:

- expose database tables/controllers as API design;
- create unversioned stable public endpoints;
- trust `tenant_id` from body/query as authorization;
- put secrets/tokens in URLs;
- use GET for mutations;
- return arbitrary error strings instead of stable problem contracts;
- use binary floats for monetary JSON values;
- return locale-formatted numbers/dates as canonical machine data;
- use offset pagination by default for high-scale feeds;
- implement a universal arbitrary query language over all module data;
- expose cross-domain internal endpoints as shortcuts around capability/event contracts;
- create undocumented business actions through status-field PATCH hacks;
- make SDK behavior the source of truth instead of the contract.

## 32. References

- OpenAPI Specification 3.2.0: `https://spec.openapis.org/oas/v3.2.0.html`
- HTTP Semantics (RFC 9110): `https://www.rfc-editor.org/rfc/rfc9110.html`
- Problem Details for HTTP APIs (RFC 9457): `https://www.rfc-editor.org/rfc/rfc9457.html`
