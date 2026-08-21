# ADR-0003 — HTTP API Contract Baseline

Status: **Accepted**  
Date: **2026-08-21**  
Work package: **P00.04**

## Context

Omnexa needs stable public, partner and cross-domain APIs that survive module extraction, SDK generation, multi-tenant authorization and long-term platform evolution. Without one contract policy, modules would invent incompatible routes, pagination, idempotency, error and versioning behavior.

P00.03 already froze identifiers, money, time, locale and error primitives. P00.04 must define their HTTP/JSON representation and compatibility rules before any runtime API is implemented.

## Decision

### Protocol and description

- Stable HTTP APIs follow RFC 9110 semantics.
- HTTP error payloads use RFC 9457-compatible Problem Details plus Omnexa extensions.
- Canonical machine-readable API descriptions use OpenAPI 3.2.0.
- The canonical description is contract source-of-truth; generated SDK/server/documentation output is derivative.

### Versioning and routes

Stable routes use:

```text
/api/v{major}/{domain}/{resources}
```

Major version changes only for breaking contract generations. Domain ownership remains visible in routing while deployment topology remains hidden.

### JSON shape

- JSON properties use lowercase `snake_case`.
- Success bodies use `data` plus optional `meta`; collection responses add `page`.
- P00.03 decimal/time/identifier/locale/error representations are mandatory.

### Business operations

- HTTP methods retain their standard semantics.
- Domain lifecycle operations such as cancel/issue/refund use explicit action endpoints rather than arbitrary `status` field mutation.
- Asynchronous acceptance uses `202` only with a defined status/outcome reference.

### Idempotency and concurrency

- Protected retriable mutations use `Idempotency-Key` when required by contract.
- Duplicate keys are scoped to tenant + principal/client + operation and protected by request fingerprints/recorded outcome.
- Lost-update-sensitive resources use `ETag`/`If-Match` optimistic concurrency where applicable.

### Collections

- Cursor pagination is the default scalable collection pattern using `page_size` and `page_cursor`.
- Simple filters use allowlisted `filter[field]` semantics.
- Sorting uses allowlisted comma-separated `sort`, with leading `-` for descending.
- Arbitrary database/query exposure is forbidden.

### Tenancy/security boundary

- tenant/organization selectors are never authorization authority;
- resolved server-side principal/scope controls access;
- credentials do not appear in query strings;
- P00.06 will define authentication/security schemes without changing resource semantics.

### Compatibility

Additive compatible changes are preferred. Field removal/rename/type changes, new required inputs, primitive semantic changes, changed idempotency/concurrency guarantees or changed authorization visibility require breaking-change handling/versioning.

## Consequences

### Positive

- one predictable API style across modules;
- SDK/codegen and automated contract testing can be standardized;
- module extraction does not change consumer-facing routes/contracts;
- multi-tenant scope remains authorization-governed;
- financial mutations have explicit retry/concurrency semantics;
- cursor pagination supports large datasets without exposing persistence implementation.

### Costs / constraints

- OpenAPI 3.2.0 tooling must be evaluated/pinned during P00.07/P01; compatibility generation may be needed for lagging tools;
- action endpoints require domain modeling rather than generic CRUD shortcuts;
- idempotency outcome storage and optimistic concurrency require real platform support later;
- strict contract compatibility requires diff/lint gates before releases.

## Rejected alternatives

### Unversioned APIs

Rejected because long-lived external/partner clients need explicit breaking-change boundaries.

### Version only in headers/content negotiation

Rejected as the sole platform strategy because path major versions are more visible, routable, cache/tool friendly and easier for broad ecosystem support. Media/header negotiation may still be used for representation concerns.

### Database-table-shaped APIs

Rejected because they leak implementation and create cross-domain coupling.

### Offset pagination as universal default

Rejected because large/changing datasets become inefficient and unstable. Offset pagination remains an explicitly justified option for small/reference/admin collections.

### Arbitrary GraphQL/query layer over every module

Rejected as the baseline because unrestricted relationship traversal would bypass domain ownership/authorization/complexity controls. GraphQL or specialized query surfaces may be introduced later as governed capabilities with an ADR where justified.

### One generic batch endpoint

Rejected because heterogeneous arbitrary commands weaken authorization, auditability, idempotency and failure semantics. Domains may expose bounded bulk operations.

## Compliance

Detailed rules are defined by:

- `docs/architecture/API_STANDARD.md`
- `docs/contracts/http/openapi-template.yaml`
- P00.03 foundation standards

P00.05/P00.06/P00.07 may add event/security/CI enforcement details but may not contradict this baseline without formal change control.
