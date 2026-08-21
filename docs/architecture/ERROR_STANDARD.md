# Omnexa Error Contract Standard

Status: **Canonical v1**  
Work package: **P00.03**

This standard defines stable, diagnosable and safe error semantics across APIs, services, jobs, workflows and user interfaces.

## 1. Principles

Errors are contracts, not arbitrary strings.

Every material failure has a **stable machine error code**. Every material error must separate:

- stable machine meaning;
- safe human-facing explanation;
- transport status where applicable;
- correlation/diagnostic identity;
- internal developer/operational detail;
- retryability where relevant.

A user-visible message is never the primary machine identifier.

## 2. Canonical machine error code

Stable error codes use lowercase dot-separated identifiers.

Examples:

```text
validation.failed
authentication.required
authorization.denied
resource.not_found
conflict.version_mismatch
commerce.order.not_found
commerce.order.invalid_state
payments.payment.declined
inventory.stock.insufficient
integration.provider.unavailable
```

Rules:

- generic platform errors may use a platform-level namespace;
- domain errors start with the owning domain when the meaning is domain-specific;
- codes describe stable semantics, not implementation exceptions;
- do not embed localized text, HTTP status, class names, database names or transient infrastructure details in codes;
- changing the meaning of a published code is a breaking contract change.

## 3. Problem representation

HTTP-facing error responses should be compatible with the Problem Details model and use this canonical Omnexa shape unless P00.04 refines transport details:

```json
{
  "type": "urn:omnexa:error:commerce.order.invalid_state:v1",
  "title": "Order state does not allow this operation",
  "status": 409,
  "code": "commerce.order.invalid_state",
  "detail": "The order cannot be cancelled in its current state.",
  "instance": "/orders/018f47a6-7b7e-7c14-a847-0af56b2a44fe",
  "request_id": "018f47b1-f3a8-79da-8d9d-4a64d65c2fd5",
  "trace_id": "<trace identifier when available>",
  "retryable": false
}
```

Rules:

- `type` is a stable URI/URN identifying the error type/version;
- `code` is the stable Omnexa machine code;
- `title` is short and safe;
- `detail` is optional, safe and may be localized at the presentation boundary;
- `status` reflects HTTP semantics only at HTTP boundaries;
- `request_id` follows the Identifier Standard;
- `trace_id` may be present when tracing is available;
- internal stack traces/debug fields are forbidden in public responses.

## 4. Validation errors

Field/object validation failures use a stable top-level code plus structured violations.

Example:

```json
{
  "type": "urn:omnexa:error:validation.failed:v1",
  "title": "Validation failed",
  "status": 422,
  "code": "validation.failed",
  "violations": [
    {
      "path": "customer.email",
      "code": "validation.email.invalid",
      "message": "Enter a valid email address."
    }
  ]
}
```

Rules:

- field paths use contract field names, not database column names;
- each violation has its own stable code;
- messages are presentation text and may be localized;
- sensitive rejected values must not be echoed unless explicitly safe.

## 5. Canonical categories

At minimum, Omnexa distinguishes:

| Category | Meaning |
|---|---|
| validation | Input/contract/business validation failed before protected mutation. |
| authentication | No valid principal/session/credential. |
| authorization | Principal exists but lacks permission/policy approval. |
| not_found | Requested resource/capability object is not visible/found. |
| conflict | Current state/version/idempotency condition conflicts with request. |
| rate_limit | Caller exceeded an enforced limit. |
| dependency | Required downstream/internal dependency failed. |
| timeout | Operation exceeded an explicit deadline. |
| unavailable | Capability temporarily cannot serve the request. |
| invariant | Internal/domain invariant was violated; normally not caller-correctable. |
| internal | Unexpected failure without a safer specific public classification. |

Domain-specific codes refine these categories.

## 6. Security and information disclosure

Public errors must not disclose:

- stack traces;
- SQL queries/schema internals;
- filesystem paths;
- secrets/tokens/credentials;
- provider secret payloads;
- internal IPs/hosts unless explicitly public;
- existence of protected resources when authorization policy requires non-disclosure;
- raw PII beyond what the caller is already authorized to view.

Detailed diagnostics belong in protected structured logs/traces linked through request/correlation identifiers.

## 7. Retry semantics

`retryable` describes whether retry may be safe/useful; it does not authorize blind retry.

Rules:

- retry policy also depends on idempotency and operation type;
- rate-limit/unavailable responses may expose an explicit retry-after hint;
- validation/authentication/authorization errors are normally non-retryable without changing caller state/input;
- payment/financial mutations must never be automatically repeated unless idempotency guarantees protect the operation.

## 8. Internal errors and wrapping

Internal services may wrap lower-level failures to add context, but they must preserve the original cause for protected diagnostics.

Do not expose raw implementation exception names as cross-domain contracts.

At a domain boundary:

1. map implementation failure to a stable owned error code;
2. retain safe causal context in logs/traces;
3. preserve correlation identity;
4. expose only contract-approved fields.

## 9. Localization

Machine error codes and `type` identifiers are never localized.

Human-facing `title`, `detail` and violation `message` may be localized according to the Locale Standard.

Do not branch business logic on translated messages.

## 10. Observability

Material failures must be diagnosable with structured telemetry.

Logs should carry, where applicable:

- error code;
- request/correlation ID;
- trace/span context;
- tenant/organization context subject to data-classification rules;
- owning domain/capability;
- retry attempt/job/workflow context;
- protected internal cause.

Do not double-log the same failure at every layer without adding ownership/context; this creates noise and can leak data.

## 11. HTTP status guidance

P00.04 owns the definitive API mapping. Until then, common mappings are:

- 400 malformed request/contract syntax;
- 401 authentication required/invalid;
- 403 authorization denied where disclosure is safe;
- 404 not found/non-disclosed resource;
- 409 state/version/idempotency conflict;
- 422 semantically valid request with validation failures;
- 429 rate limited;
- 500 unexpected internal failure;
- 502 dependency/gateway failure where appropriate;
- 503 temporarily unavailable;
- 504 upstream/dependency timeout.

Do not choose a status code merely to match an existing framework exception.

## 12. Versioning and compatibility

Published error semantics are part of the contract.

Breaking examples include:

- reusing a code for a different meaning;
- removing fields consumers rely on;
- changing retryability semantics in a way that can cause duplicate mutations;
- changing authorization errors to reveal protected resource existence.

Such changes require contract versioning/change control under P00.04 and the ADR process where architectural.

## 13. Prohibited patterns

Do not:

- return `{"error":"Something went wrong"}` as the only machine contract;
- use HTTP 200 for failed protected operations;
- expose exception/class/database names as public error codes;
- localize machine codes;
- return stack traces in production responses;
- retry financial mutations without idempotency protection;
- use one global `500` response for all known business errors;
- silently swallow errors that change business state or auditability.
