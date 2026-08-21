# Omnexa Platform Security Standard

Status: **Canonical v1**  
Work package: **P00.06**

This standard defines the security invariants that all Omnexa kernel capabilities, modules, APIs, events, workflows, integrations, extensions, AI tools, admin surfaces, developer tooling and deployment models must respect.

Security is a platform property. A module cannot opt out because it is “internal,” because a UI hides an action, or because another module already authenticated the caller.

## 1. Security objectives

Omnexa security architecture protects:

- tenant isolation;
- identity and session integrity;
- authorization correctness;
- confidentiality, integrity and availability of business data;
- financial/business mutation integrity;
- auditability and non-repudiation appropriate to the action;
- secrets and cryptographic material;
- integration boundaries;
- extension/module supply chain;
- operational/recovery systems;
- AI capability execution.

## 2. Trust model

Assume every trust boundary can fail independently.

Canonical trust boundaries include:

```text
Internet / untrusted client
    -> edge / gateway
    -> API/service capability
    -> owning domain
    -> database/cache/object store

Tenant A <-> Tenant B
Organization A <-> Organization B
User/service principal <-> authorization policy
Module <-> module public contract
Omnexa <-> external provider
Cloud control plane <-> regional data plane
Cloud <-> edge/POS/local agent
Human/AI agent <-> governed capability
CI/build system <-> artifact/signing/release system
Support/admin operator <-> customer data
```

Crossing a boundary requires explicit authentication/identity where applicable, authorization, validation, data classification handling and observability.

## 3. Zero implicit trust

Do not trust solely because traffic originates from:

- the same process;
- the same VPC/network;
- another Omnexa module;
- an internal admin UI;
- a background worker;
- a webhook source IP;
- an AI agent;
- an edge/POS device;
- a previously authenticated browser.

Trust is derived from authenticated identity, validated execution context, scoped capabilities and policy.

## 4. Identity types

Omnexa distinguishes at least:

- human `User`;
- `Service Account`;
- workload/service identity;
- device/edge identity;
- integration/OAuth client;
- support/operator identity;
- AI agent execution identity.

Non-human identities must never be represented as fake human users merely to reuse role logic.

Every principal has:

- stable identity;
- lifecycle state;
- authentication mechanism appropriate to its type;
- bounded scopes/capabilities;
- tenant/organization relationships where applicable;
- auditable credential lifecycle.

## 5. Authentication

Authentication proves identity; it does not grant business authority by itself.

Requirements:

- server-side verification at protected boundaries;
- strong password hashing through an approved adaptive password-hashing algorithm when passwords are supported;
- password plaintext/reversible password storage prohibited;
- MFA/passkeys supported by architecture for privileged and policy-required use cases;
- external identity/SSO claims validated for issuer, audience, signature, time and configured trust relationship;
- service credentials are scoped and rotatable;
- authentication failures do not disclose unnecessary account existence/details;
- recovery/reset flows are treated as authentication flows, not ordinary profile updates.

## 6. Sessions and tokens

Session/token architecture must provide:

- explicit expiry;
- rotation/revocation strategy;
- secure token transport/storage appropriate to client type;
- replay/theft mitigation appropriate to risk;
- device/session inventory for interactive users where supported;
- session invalidation after material account/security changes according to policy;
- tenant/organization context re-authorization rather than blindly trusting stale client context.

Rules:

- access tokens are short-lived relative to refresh/session credentials;
- refresh/session secrets are `RESTRICTED`;
- tokens are never logged;
- bearer possession is not enough to bypass current authorization/policy state;
- token contents must not be used as an excuse to skip authoritative permission/relationship checks when those may have changed.

## 7. Authorization model

Omnexa authorization combines:

```text
RBAC
+ relationship-based authorization
+ contextual/policy checks
+ capability boundaries
```

Rules:

- roles are permission compositions, not hard-coded bypasses;
- `admin`, `owner` or `superuser` names do not imply unrestricted platform bypass;
- all protected mutations enforce server-side authorization;
- object/tenant/organization relationships are checked where relevant;
- field-level and export-level authorization may be stricter than resource read permission;
- background jobs/events/workflows carry trusted execution context or resolve an approved service identity;
- cross-domain calls invoke public capabilities that enforce owning-domain policy;
- authorization failures are auditable when material and disclosure-safe.

## 8. Capability security

Every protected capability defines:

- owning domain;
- allowed principal types;
- required permission/policy relationship;
- tenant/organization scope;
- input validation;
- data classification exposure;
- audit requirement;
- idempotency/concurrency semantics where relevant;
- approval requirement for high-risk actions;
- rate/abuse controls where exposed externally.

AI, workflows and integrations invoke the same governed capabilities; they do not receive alternate “internal” write paths.

## 9. Tenant isolation

Tenant isolation is a mandatory invariant.

Requirements:

- tenant ownership is explicit in persisted tenant-owned state through canonical `tenant_id` where applicable;
- tenant context is resolved from trusted identity/policy/execution context;
- client-provided tenant IDs are never authorization authority;
- repository/data access abstractions must make tenant scoping hard to omit accidentally;
- cross-tenant operations require explicit platform-level capability and audit;
- caches, search, analytics, files, queues, events and vector indexes preserve tenant boundaries;
- tenant export/backup/restore operations cannot accidentally include another tenant;
- negative cross-tenant tests are mandatory for affected implementation paths.

“No row returned” and “not authorized” may deliberately be indistinguishable where resource-existence disclosure is sensitive.

## 10. Organization and sub-scope isolation

Inside a tenant, Organization, Legal Entity, Business Unit, Branch, Team and other scopes remain distinct.

Tenant membership is not sufficient authority for every child scope.

Cross-organization operations require explicit policy/relationships, and domain records must not silently fall back to “any organization within this tenant.”

## 11. Input validation

Every trust boundary validates input independently.

Validation covers:

- type/shape/schema;
- identifiers and ownership references;
- size/count/range limits;
- enum/state transitions;
- filenames/content types where relevant;
- URLs/redirect targets;
- external resource identifiers;
- locale/time/money primitives;
- business invariants.

Do not depend solely on frontend validation.

Validation does not replace output encoding, authorization or parameterized data access.

## 12. Injection and output safety

Implementation must use safe APIs and context-appropriate encoding to prevent injection classes, including:

- SQL/query injection;
- shell/command injection;
- template injection;
- HTML/script injection;
- header injection;
- path traversal;
- unsafe deserialization;
- expression/workflow injection;
- prompt/tool injection leading to unauthorized capability execution.

Raw user-controlled values must not be interpolated into privileged interpreter/query contexts.

## 13. Secrets management

Secrets include credentials, API keys, private keys, signing keys, encryption keys, refresh tokens and other authentication-equivalent material.

Rules:

- secrets are never committed to source control;
- production secrets are stored through an approved secret-management/KMS mechanism;
- configuration references secrets; it does not embed them;
- secrets are scoped to the smallest practical environment/service/tenant purpose;
- rotation and revocation are supported;
- secret values are redacted from logs, traces, errors, support tools and CI output;
- secrets are not passed to AI models;
- developer `.env` files or local secret stores are excluded from version control;
- example configuration uses non-secret placeholders only.

## 14. Encryption in transit

Protected network communication uses authenticated encryption in transit.

- Internet-facing production endpoints use TLS.
- Service-to-service transport uses protected channels appropriate to deployment trust boundaries.
- database/cache/broker/storage connections use encryption where crossing untrusted/shared boundaries or where deployment policy requires it.
- certificate/peer verification must not be disabled in production to “fix connectivity.”
- webhook signatures/authentication are verified independently of TLS.

## 15. Encryption at rest and key management

Encryption at rest protects storage media but does not replace authorization.

Requirements:

- approved storage encryption for production data stores/object stores/backups;
- key material separated from encrypted data where architecture permits;
- key access follows least privilege;
- rotation/versioning supported without destructive re-encryption assumptions;
- sensitive field-level/application encryption used only when threat/data requirements justify it;
- password hashing is not encryption;
- cryptographic algorithms/key sizes are centralized approved policy, not invented per module.

## 16. Cryptographic ownership

Modules do not implement custom cryptography.

The platform provides governed cryptographic capabilities for:

- secure random generation;
- hashing/HMAC;
- encryption/decryption where approved;
- signing/verification;
- key versioning/rotation;
- token generation.

Protocol/algorithm changes are security architecture changes and require review.

## 17. Audit logging

Audit is separate from application debug logging.

Business/security-significant operations record, where applicable:

- audit event ID;
- UTC timestamp;
- tenant/organization scope;
- actor/principal type and ID;
- effective capability/action;
- target resource/type/ID;
- outcome;
- request/correlation/trace context;
- origin/channel/device/integration context when useful;
- approval/impersonation/delegation context;
- classification-safe change summary.

Audit records must not store secrets and must be protected from ordinary user mutation.

## 18. Privileged and dangerous operations

Examples include:

- role/permission changes;
- tenant ownership transfer;
- MFA/authentication policy change;
- secret/key operations;
- bulk export;
- destructive purge;
- payment refund/payout/settlement actions;
- financial close/reversal;
- impersonation/support access;
- event replay affecting side effects;
- marketplace/module trust changes;
- AI high-impact action approval.

Such operations may require stronger authentication, explicit approval, reason capture, dual control or additional audit based on risk policy.

## 19. Impersonation and support access

Support/admin access to a tenant is never an invisible superuser shortcut.

If impersonation/support access is supported it requires:

- dedicated capability;
- explicit operator identity;
- target tenant/user context;
- reason/ticket reference where policy requires;
- bounded duration;
- clear UI indication;
- audit trail;
- restrictions on especially sensitive actions;
- no exposure of user passwords/secrets.

## 20. API security

In addition to `API_STANDARD.md`:

- authentication/authorization occurs before protected data exposure;
- object-level authorization is explicit;
- mass assignment is prevented through contract allowlists;
- pagination/query complexity is bounded;
- rate/abuse controls are capability-aware;
- idempotency keys cannot be reused across different principals/tenant scopes/payloads;
- CORS is explicit, not wildcard-by-default for credentialed APIs;
- redirect/callback URLs use registered allowlists where applicable;
- public errors follow disclosure-safe error standard.

## 21. Event security

In addition to `EVENT_STANDARD.md`:

- publishers are authenticated/authorized to publish on producer-owned subjects/streams;
- consumers subscribe only to required event classes;
- tenant context is producer-derived and validated;
- events are not authorization proof for new privileged actions;
- dead-letter/replay tooling is privileged and audited;
- sensitive payload handling follows Data Classification Standard;
- external event ingestion is treated as untrusted until verified/validated.

## 22. Webhooks

Inbound webhooks require provider-appropriate authenticity verification, such as signature validation, mTLS or another approved mechanism.

Requirements:

- verify signature/authentication before business processing;
- prevent replay where provider semantics support timestamps/nonces/event IDs;
- enforce payload size/schema limits;
- deduplicate delivery;
- retain raw payload only when justified/classified;
- do not trust provider-supplied tenant/account mapping without configured integration ownership.

Outbound webhooks:

- sign/authenticate deliveries where supported;
- minimize payloads;
- prevent SSRF through destination registration/validation;
- use bounded retries;
- expose delivery status without leaking secrets.

## 23. SSRF and outbound network controls

User/admin-configurable URLs are not automatically safe destinations.

Capabilities that fetch/call URLs must defend against:

- loopback/private/link-local/metadata endpoints where not explicitly allowed;
- DNS rebinding;
- unsafe redirects;
- unsupported schemes;
- credential leakage through URLs/headers;
- unrestricted port/protocol access.

High-risk outbound connectivity may require centralized egress policy/proxy.

## 24. Files and uploads

Untrusted uploads require:

- authenticated/authorized upload capability;
- size/type limits;
- generated storage keys rather than trusting filenames;
- path traversal prevention;
- private storage by default for non-public content;
- malware/content scanning where relevant;
- safe image/document processing isolation where parsers are complex;
- download-time authorization;
- content-disposition/content-type hardening;
- prevention of executable active content in contexts not designed for it.

## 25. Browser and frontend security

Frontend security is defense-in-depth; backend policy is authoritative.

The web platform should support:

- secure cookie/session attributes where cookies are used;
- CSRF protection for cookie-authenticated state changes;
- content security policy strategy;
- clickjacking protection;
- safe external-link/redirect handling;
- sensitive-data masking;
- avoiding secrets in local storage/source maps/client logs;
- dependency integrity and controlled third-party scripts.

## 26. POS / edge / device security

Edge/POS/local agents require a distinct device identity and lifecycle.

Architecture must support:

- enrollment/provisioning;
- device credential rotation/revocation;
- encrypted local sensitive state;
- secure update/signature verification;
- least-privilege device capabilities;
- offline authorization policy with explicit risk limits;
- tamper/replay-aware sync;
- tenant/store/device binding;
- remote disable/wipe of Omnexa-managed secrets where feasible.

A device being physically inside a store/network is not sufficient trust.

## 27. Integrations and OAuth apps

Third-party connectors are separate principals with explicit scopes.

Requirements:

- least-privilege scopes/capabilities;
- tenant-specific authorization/consent;
- secure token storage/refresh/revocation;
- redirect URI allowlists;
- provider account binding validated against the tenant configuration;
- webhook authenticity validation;
- classification-aware data disclosure;
- removal/revocation behavior on uninstall.

## 28. Module and marketplace security

A module/extension declares permissions and external access requirements before installation.

Future marketplace/runtime controls must support:

- signed/verifiable packages;
- immutable package/version identity;
- SBOM/dependency metadata;
- declared permissions/capabilities;
- declared network/external services;
- install/upgrade review gates;
- isolation/sandboxing where extension type requires it;
- revocation/quarantine of compromised versions;
- prohibition on private cross-module imports/tables.

Installing a module does not grant it unrestricted tenant data access.

## 29. Dependency and software supply chain

Build/release process must eventually enforce:

- lockfiles/reproducible dependency intent;
- vulnerability scanning;
- secret scanning;
- provenance/build identity;
- artifact integrity/signing where release model requires;
- dependency update policy;
- review of high-risk native/build scripts;
- minimal CI token permissions;
- protected release credentials;
- no secrets available to untrusted fork code paths.

P00.07 freezes executable CI/release gates.

## 30. Database security

- applications use least-privilege database roles appropriate to runtime/migration/ops duties;
- runtime credentials do not receive arbitrary administrative rights;
- parameterized access is mandatory;
- direct human production access is exceptional and audited;
- backups are encrypted/access-controlled;
- tenant-safe query design is mandatory;
- cross-domain direct writes remain forbidden even when database permission technically permits them.

## 31. Cache, broker and search security

Redis-compatible, NATS/JetStream-class and search systems are protected infrastructure, not trusted public endpoints.

- authenticate where supported/required;
- network exposure restricted;
- tenant context preserved in data/routing design;
- no secrets in cache/message keys;
- TTL/retention aligned to classification/purpose;
- administrative consoles not exposed publicly by default;
- payload/log handling follows classification rules.

## 32. Backup, restore and disaster recovery security

Backups inherit source data classification.

Requirements:

- encryption and access control;
- inventory/retention policy;
- restore authorization;
- tenant isolation validation for tenant-scoped restore/export mechanisms;
- restoration into lower environments prohibited by default for sensitive production data;
- restore procedures preserve secrets/key/version dependencies;
- backup deletion/expiry supports retention obligations.

## 33. Security telemetry and incident response

Security-significant events should produce structured telemetry suitable for detection and investigation, including:

- authentication anomalies;
- privileged role/policy changes;
- cross-tenant denial attempts;
- secret/key lifecycle events;
- bulk exports;
- support impersonation;
- module install/trust changes;
- high-risk AI actions;
- repeated webhook/signature failures;
- suspicious rate/abuse patterns.

Telemetry must not leak the secrets it is trying to protect.

Incident-response runbooks/operational readiness are expanded in later operational phases.

## 34. AI and agent security

AI is an untrusted planner/interpreter operating through governed capabilities.

Rules:

- model output never directly mutates the database;
- tool/capability execution uses an explicit execution principal and current authorization;
- prompts/retrieval are tenant/object scoped before data reaches the model;
- tool arguments are validated as untrusted input;
- prompt injection cannot grant new capabilities;
- retrieved customer content cannot redefine system/security policy;
- sensitive/high-impact actions may require human approval;
- model/provider cannot receive `RESTRICTED` data by default;
- prompt/output/tool audit follows classification rules;
- AI cannot disable logging, authorization, policy or approval controls;
- agents receive least-privilege toolsets per task, not the entire platform API.

## 35. Workflow security

Workflow definitions are executable policy-adjacent artifacts.

- editing/publishing workflows is a protected capability;
- each action executes under a defined service/user/delegated identity model;
- delayed steps re-evaluate authorization/policy where business semantics require current authority;
- secrets are referenced, never embedded in workflow JSON;
- external HTTP actions use integration/egress controls;
- retries/compensations remain idempotent;
- privileged workflows are versioned and audited.

## 36. Rate limiting and abuse prevention

Rate limits are layered by risk/context, potentially including:

- source/IP;
- principal/client;
- tenant;
- capability/resource;
- expensive query/export;
- authentication/recovery flow;
- webhook/integration provider.

Rate limiting does not replace authorization or business quota enforcement.

## 37. Environment separation

Production, staging and development are distinct trust environments.

- credentials/secrets are environment-specific;
- production secrets unavailable to lower environments by default;
- production data does not flow down by convenience;
- release promotion preserves artifact identity rather than rebuilding ad hoc where provenance matters;
- test/dev bypasses must not compile silently into production behavior.

## 38. Secure defaults

New capabilities default to:

- private/authenticated rather than public;
- deny rather than allow when authorization policy is undefined;
- tenant scoped rather than global;
- minimal fields rather than full object dump;
- no sensitive logging;
- no external network access unless declared;
- no destructive purge unless explicitly enabled/protected;
- no AI access unless classification/policy permits.

## 39. Security exceptions

A security exception requires:

- named owner;
- exact control being bypassed;
- business/technical justification;
- affected scope/data classification;
- compensating controls;
- expiry/review date;
- approval authority;
- tracked issue/risk.

“Temporary” without an expiry/review is not an accepted exception.

## 40. Prohibited patterns

Do not:

- create an authorization-bypass superuser path;
- trust client tenant IDs as authority;
- log/access-token/password/private-key values;
- store passwords reversibly;
- disable TLS/certificate verification in production;
- embed secrets in code/workflows/module manifests;
- allow modules direct private cross-domain data writes;
- make support impersonation invisible;
- give AI unrestricted database/network/tool access;
- treat events/webhooks as proof of permission;
- expose internal admin consoles publicly by default;
- copy production restricted/confidential data to dev/test for convenience;
- invent custom cryptographic algorithms in business modules;
- rely on frontend hiding/disabled buttons as authorization.

## 41. Change control

Changes to tenant isolation, authentication/session semantics, authorization model, audit integrity, secret/key handling, cryptographic baseline, privileged-operation controls, support impersonation, module trust model or AI execution boundary require formal security review and an ADR/change-control reconciliation before implementation.
