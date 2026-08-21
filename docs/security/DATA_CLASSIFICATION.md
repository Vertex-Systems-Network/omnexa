# Omnexa Data Classification Standard

Status: **Canonical v1**  
Work package: **P00.06**

This standard defines how Omnexa classifies, handles, exposes, logs, stores, transmits, exports, retains and deletes data. Classification follows the data itself across modules, APIs, events, files, analytics, support tooling and AI systems.

## 1. Principles

1. Data is classified by the highest-impact characteristic present, not by the screen or table that contains it.
2. Classification is inherited by derived copies unless an approved irreversible transformation demonstrably lowers sensitivity.
3. Tenant isolation is independent of sensitivity classification: even `INTERNAL` tenant data must not cross tenants without explicit authorized purpose.
4. Data minimization is preferred over broad access plus masking after the fact.
5. Production data must not be copied into lower environments by default.
6. Logs, traces, queues, dead-letter stores, analytics and AI prompts are data stores and must honor classification rules.
7. A module cannot lower the classification of data owned by another domain.
8. Retention and deletion are governed by data purpose, legal/business requirements, tenant policy and system integrity—not by arbitrary table cleanup.

## 2. Canonical classification levels

Omnexa uses four confidentiality classes:

| Class | Meaning | Typical examples |
|---|---|---|
| `PUBLIC` | Intended for unrestricted external disclosure. | Published website content, public product description, public help documentation. |
| `INTERNAL` | Non-public operational data with limited impact if disclosed. | Internal feature configuration, non-sensitive workflow metadata, internal documentation. |
| `CONFIDENTIAL` | Business, personal or tenant data whose unauthorized disclosure could cause meaningful harm. | Customer/contact data, invoices, contracts, non-public prices, employee records, tenant business data. |
| `RESTRICTED` | Highest-sensitivity data requiring exceptional protection or narrow handling. | Credentials/secrets, authentication factors, sensitive payment/security material, private keys, recovery codes, high-risk regulated data when supported. |

The default for tenant-owned business records is **at least `CONFIDENTIAL`** unless the owning domain explicitly publishes a `PUBLIC` projection.

## 3. Separate handling tags

Confidentiality class alone is not sufficient. Data may also carry one or more handling tags:

```text
PII
AUTH_SECRET
CRYPTO_KEY
PAYMENT_SENSITIVE
FINANCIAL_RECORD
EMPLOYEE_DATA
LEGAL_RECORD
HEALTH_SENSITIVE
CHILD_DATA
BIOMETRIC
LOCATION_PRECISE
SECURITY_TELEMETRY
CUSTOMER_CONTENT
MODEL_INPUT
MODEL_OUTPUT
```

Tags describe handling obligations and may trigger stricter controls than the base confidentiality class.

A domain must not introduce a new high-risk handling tag silently. New tags require security review and registry update.

## 4. Classification registry

Each authoritative domain must maintain a machine-readable field/data-category declaration before implementation reaches production maturity. At minimum it identifies:

- domain/module owner;
- resource/field or data category;
- confidentiality class;
- handling tags;
- tenant/organization ownership;
- permitted storage classes;
- logging policy;
- export policy;
- retention policy identifier;
- masking/redaction requirements;
- AI eligibility;
- indexing/search eligibility.

P00 defines the platform model; domain-specific registries are created with their modules.

## 5. Classification inheritance

Derived data follows these rules:

- exact copy: inherits source classification and tags;
- aggregation: remains at least as sensitive as required to prevent re-identification;
- pseudonymization: does not automatically become `PUBLIC`;
- encryption: protects data but does not lower its classification;
- irreversible anonymization: may lower classification only after an approved re-identification risk assessment;
- generated AI output: inherits relevant source sensitivity if it contains or reveals source information.

## 6. PUBLIC

`PUBLIC` data:

- may be exposed through explicitly public capabilities;
- must still preserve integrity and provenance where business-significant;
- must not be confused with “no authorization required” for write operations;
- must not contain hidden confidential metadata in files, HTML, APIs or generated artifacts.

## 7. INTERNAL

`INTERNAL` data:

- requires authenticated or service-authorized access unless a specific public projection exists;
- may be logged when useful and safe;
- must not be made public merely because it has no obvious PII;
- remains tenant-scoped where tenant-owned.

## 8. CONFIDENTIAL

`CONFIDENTIAL` data:

- requires explicit authorized business purpose;
- uses encryption in transit and approved encryption at rest;
- is excluded from public logs and traces by default;
- is masked/redacted in support/admin interfaces unless full value is required and authorized;
- is subject to tenant-aware export and deletion controls;
- must not be placed in unrestricted analytics, test fixtures or prompt logs;
- should be minimized in events/webhooks and external integrations.

## 9. RESTRICTED

`RESTRICTED` data:

- follows least-privilege access with narrow capabilities;
- must not be written to ordinary logs, analytics events, crash payloads, traces, dead-letter bodies or support screenshots;
- requires approved secret/key/specialized storage where applicable;
- is encrypted in transit and at rest with key-management controls appropriate to the data type;
- requires explicit audit of privileged access where technically feasible;
- must not be used as AI model input unless a separately approved policy explicitly allows the data class/use case;
- must not be exported through generic bulk export endpoints;
- may require shorter retention or cryptographic destruction according to policy.

Passwords are never stored reversibly. Authentication secrets are handled according to the Security Standard, not as ordinary encrypted profile fields.

## 10. Storage handling matrix

| Surface | PUBLIC | INTERNAL | CONFIDENTIAL | RESTRICTED |
|---|---|---|---|---|
| Primary OLTP | allowed | allowed | allowed with controls | only when approved for data type |
| Cache | allowed | allowed | minimized/TTL | exceptional; avoid by default |
| Object storage | allowed | controlled | private/encrypted | dedicated approved handling |
| Search index | allowed | controlled | explicit field allowlist | prohibited by default |
| Analytics warehouse | allowed | controlled | minimized/pseudonymized where possible | prohibited by default |
| Logs/traces | allowed | controlled | metadata/redacted only by default | prohibited |
| Events/webhooks | allowed | controlled | minimum necessary | prohibited by default |
| Lower environments | allowed | synthetic preferred | synthetic/anonymized only by default | prohibited |
| AI prompts/context | allowed | policy-controlled | explicit purpose/policy | prohibited by default |

“Allowed” never overrides tenant scope, authorization or purpose limitation.

## 11. Logging and observability

Structured telemetry must classify fields before they become common logging helpers.

Do not log by default:

- passwords or password equivalents;
- access/refresh tokens;
- API secrets;
- session secrets;
- private keys;
- recovery codes;
- full payment credentials;
- raw authorization headers;
- secret webhook signatures;
- unnecessary customer content;
- sensitive document bodies.

For `CONFIDENTIAL` values, log stable IDs, classification-safe metadata or redacted representations when diagnostics require context.

## 12. API and UI exposure

- API schemas must expose only fields required by the capability.
- Read permission for a resource does not automatically authorize every sensitive field.
- Sensitive fields may require field-level capability/policy checks.
- UI masking is defense-in-depth and never substitutes for server-side authorization.
- list endpoints should avoid returning high-sensitivity fields merely because detail endpoints can.
- exports are distinct privileged capabilities and must not inherit broad UI-read assumptions.

## 13. Event and webhook handling

Events/webhooks follow `EVENT_STANDARD.md` plus classification rules:

- publish the minimum immutable business context;
- prefer stable references when payload duplication creates unnecessary sensitivity;
- external webhooks are explicit disclosure boundaries;
- payload schemas must eventually annotate classified fields;
- dead-letter/quarantine stores must apply the same or stronger protection as the payload source;
- replay tooling must preserve access controls and auditability.

## 14. Files and media

Private files require:

- non-public object storage by default;
- authorization at access/download time;
- short-lived signed access where appropriate;
- malware/content validation workflow where untrusted uploads are supported;
- metadata scrubbing where published assets could leak private metadata;
- classification inherited from the containing business record unless explicitly different.

A guessable object URL is not authorization.

## 15. Search and indexing

Search is a disclosure surface.

- indexing requires field allowlists;
- result visibility is re-authorized at query/result boundary;
- tenant scope is mandatory;
- restricted secrets are not indexed;
- cached/search projections must be invalidated or filtered after permission/scope changes;
- autocomplete/suggestions must not leak hidden records.

## 16. Analytics and BI

Analytics systems must not become a second ungoverned source of truth.

- ingest minimum fields needed for the metric/use case;
- keep tenant boundaries explicit;
- pseudonymize where detailed identity is unnecessary;
- restrict raw row-level exports;
- propagate deletion/retention requirements where technically and legally applicable;
- document exceptions for immutable financial/legal records.

## 17. AI and model handling

AI use requires both authorization and classification eligibility.

- an agent sees only data available to its executing principal/service policy;
- retrieval applies tenant and object authorization before context assembly;
- prompts do not receive secrets by default;
- `RESTRICTED` data is prohibited as model input by default;
- model/provider routing must respect data residency/retention/provider policy once configured;
- prompt/output logging obeys source classification;
- embeddings/vector representations inherit source sensitivity;
- model output is treated as untrusted until validated for protected actions;
- AI-generated exports cannot bypass ordinary export permissions.

## 18. Data export

Exports are high-impact operations.

Bulk export capabilities require:

- explicit authorization;
- tenant/organization scope;
- classification-aware field selection;
- audit record;
- bounded generation/access lifetime;
- secure delivery;
- no inclusion of secrets/restricted fields unless a specialized approved export exists.

## 19. Retention

Retention policies are identified by policy IDs rather than scattered hard-coded durations.

A retention policy defines:

- data category/classification;
- business/legal basis;
- active retention period;
- archive period if any;
- deletion/anonymization method;
- legal-hold behavior;
- backup expiration implications;
- downstream copy propagation.

Domains reference approved policies; they do not silently invent retention periods.

## 20. Deletion and erasure

Deletion semantics distinguish:

```text
soft_delete
archive
anonymize
purge
crypto_shred
legal_hold
```

Rules:

- user-facing deletion does not automatically mean physical purge;
- purge requires authorization, audit and referential/business-integrity checks;
- immutable accounting/audit/legal records may require retention even when related profile data is anonymized;
- deletion requests must propagate to governed replicas/projections/search/AI indexes according to policy;
- backups age out under backup retention rather than being rewritten unsafely unless legal policy demands another mechanism.

## 21. Lower environments

Production `CONFIDENTIAL`/`RESTRICTED` data must not be copied to development/test by default.

Preferred order:

1. synthetic fixtures;
2. generated representative datasets;
3. approved irreversibly anonymized samples;
4. exceptional controlled production subsets only through explicit security approval.

Developer convenience is not sufficient justification.

## 22. Third-party disclosure

Sending data to a payment gateway, email provider, AI provider, analytics provider, storage provider or other integration is an external disclosure boundary.

The connector/integration must declare:

- classes/tags transmitted;
- purpose;
- destination/provider;
- authentication method;
- encryption expectations;
- provider retention/data-use expectations;
- webhook/callback verification;
- failure/retry storage;
- tenant configuration/consent requirements where applicable.

## 23. Classification changes

Raising classification may be done defensively when new risk is discovered.

Lowering classification requires evidence that the data semantics/handling justify it. Encryption, access restriction or moving data to another table does not by itself lower classification.

Material platform-wide classification changes require security review and change control.

## 24. Prohibited patterns

Do not:

- label all data `INTERNAL` to avoid classification work;
- treat encrypted data as non-sensitive;
- send production customer data to test/dev by default;
- log tokens, passwords, private keys or full payment secrets;
- put restricted data in generic search/vector indexes;
- make bulk export equivalent to ordinary list permission;
- allow a module to downgrade another domain's data classification;
- infer tenant authorization from a tenant ID inside a payload;
- assume AI provider input/output is exempt from normal classification rules.
