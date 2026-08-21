# Omnexa Identifier Standard

Status: **Canonical v1**  
Work package: **P00.03**

This standard defines stable identity and reference rules for entities, requests, events and distributed operations across Omnexa.

## 1. Canonical entity identifier

The default Omnexa entity identifier is **UUID version 7 (UUIDv7)**.

Rules:

- UUIDv7 is the default for new platform and business entities unless an accepted ADR proves a different identifier is required.
- Identifiers are immutable, opaque and never recycled.
- A database sequence, row position, timestamp or customer-visible code must never be treated as an entity's canonical identity.
- IDs must not encode tenant, organization, region, shard, role, security classification or business meaning.
- External-system identifiers are mappings/aliases owned by an integration boundary; they do not replace the Omnexa canonical ID.
- Random/ordered identifier generation must use a standards-compliant implementation and cryptographically suitable randomness for the random component.

## 2. Representation by boundary

### PostgreSQL

Store UUID identifiers using PostgreSQL's native `uuid` type.

Do not store canonical UUIDs as `varchar`, integer hashes or binary formats invented by individual modules.

### JSON/API/contracts

Represent UUIDs as lowercase canonical hyphenated strings. Contract consumers must treat the value as an opaque string. Lexicographic ordering must never be used as a business ordering guarantee.

### Logs/traces

Identifiers may be emitted as canonical strings when the data classification permits it. Sensitive identifiers must still respect logging/redaction policy.

## 3. Scope fields

Canonical scope reference names are:

```text
tenant_id
organization_id
legal_entity_id
business_unit_id
branch_id
team_id
user_id
```

Rules:

- tenant-owned persisted records use `tenant_id` when tenant ownership applies;
- do not invent `workspace_id`, `company_id`, `account_id` or another alias as a hidden substitute for tenant scope;
- organization-level ownership uses `organization_id` in addition to `tenant_id` when both scopes are material;
- more specific scope fields are additive and must not remove the primary tenant isolation boundary;
- globally shared platform/reference records must be explicitly classified as platform-scoped rather than silently using a nullable `tenant_id` convention.

## 4. Cross-domain references

A domain may retain the stable identifier of an entity owned by another domain, but it does not gain ownership of that entity.

Rules:

- cross-domain references are logical contract references, not permission to write another domain's schema;
- cross-domain database foreign keys are forbidden by default under the module architecture standard;
- cached labels/names are snapshots or projections and must be named/treated as such;
- disabling or removing an optional source module must not invalidate unrelated historical records that legitimately retain its identifier/snapshot.

## 5. Human-facing business numbers

Human-facing identifiers are separate from canonical IDs. Examples include invoice numbers, order numbers and employee numbers.

Requirements:

- uniqueness scope must be explicit;
- renumbering rules must be explicit;
- gaps must be allowed unless a legal/country rule states otherwise;
- legal numbering sequences belong to the owning domain/country pack;
- never expose a database sequence merely because it already exists.

## 6. Request, trace and correlation identity

Use UUIDv7 for Omnexa-generated request IDs, correlation IDs, event IDs, workflow execution IDs, job execution IDs and audit record IDs.

A trace provider may also carry its own trace/span identifiers. Omnexa correlation identity and OpenTelemetry trace identity may coexist; neither should be silently substituted for the other.

## 7. Idempotency keys

Idempotency keys are not entity IDs.

Rules:

- accept/generate opaque keys at an explicit operation scope;
- never derive them from mutable request payload order;
- store the owning tenant/principal, operation/capability, key and protected result/fingerprint according to the later API/event standards;
- expiry/retention must be explicit;
- reusing the same key with materially different input must fail safely rather than execute a second mutation.

The full HTTP/API semantics are defined in P00.04.

## 8. Temporary/local identifiers

Client/offline systems such as POS may create canonical UUIDv7 IDs locally when the local runtime can guarantee compliant generation. Synchronization must not replace those IDs after upload.

UI-only temporary keys must be clearly scoped and never persisted as canonical identities.

## 9. Prohibited patterns

Do not use as canonical entity identity:

- auto-increment integers exposed outside an owning private implementation;
- email address, phone number, username or domain name;
- mutable SKU/order/invoice code;
- database natural key selected only for convenience;
- tenant-prefixed composite string IDs;
- timestamps alone;
- hashes of PII;
- random strings with undocumented entropy/format.

## 10. Compatibility rule

Once a canonical ID is externally observable through a stable contract, changing its type/semantics is a breaking architecture/contract change and requires an ADR plus versioned migration/compatibility strategy.
