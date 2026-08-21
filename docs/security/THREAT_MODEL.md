# Omnexa Initial Threat Model

Status: **Canonical foundation v1**  
Work package: **P00.09**

This threat model defines the minimum threats every future Omnexa implementation must consider. It is architecture-level and must be refined per module/deployment/industry pack before production exposure.

## 1. Protected assets

High-value assets include:

- tenant and organization boundaries;
- identities, sessions, authorization relationships and policy;
- customer/employee/vendor/business data;
- financial ledgers, invoices, orders, payments and settlements;
- secrets, keys, signing material and integration tokens;
- audit/security records;
- module/package trust and provenance;
- event/workflow state and replay authority;
- backups, exports, search/analytics/vector projections;
- POS/edge device identity and offline state;
- AI retrieval context and tool/capability authority;
- release pipeline credentials and artifacts.

## 2. Trust boundaries

Primary boundaries:

```text
Internet/client -> edge/gateway -> capability/domain
Tenant A <-> Tenant B
Organization/scope A <-> B
Human/service/device/AI principal -> authorization engine
Module -> module public contract
Runtime -> database/cache/broker/object store/search
Omnexa -> external provider/webhook destination
Cloud -> POS/edge/local agent
Developer/CI -> build/signing/release
Support operator -> tenant/customer data
Control plane -> regional data plane
Backup/recovery system -> production data
```

Every boundary requires explicit identity/trust, validation, authorization, data-classification handling and observability appropriate to the risk.

## 3. Threat taxonomy

Omnexa evaluates at least spoofing, tampering, repudiation/audit loss, information disclosure, denial of service and elevation of privilege, plus multi-tenant, supply-chain and AI-specific abuse.

## 4. Priority threat scenarios

### T01 — Cross-tenant data escape

Examples: missing tenant filter, cache key collision, search result leakage, analytics projection error, signed URL reuse, backup/export mixing.

Required controls/evidence:

- trusted tenant context;
- tenant-scoped persistence/cache/search/file/event design;
- same-tenant positive + cross-tenant negative tests;
- explicit privileged cross-tenant capability only;
- audit of approved cross-tenant operations.

Severity baseline: **critical**.

### T02 — Object/relationship authorization bypass

Examples: IDOR/BOLA, tenant member reading another organization's records, guessed resource IDs, field-level disclosure, bulk export using ordinary read permission.

Controls: owning-domain authorization, relationship/context policy, field/export-specific capabilities, negative tests.

Severity baseline: **critical/high** depending on asset.

### T03 — Authentication/session takeover

Examples: credential stuffing, stolen refresh/session token, MFA bypass, reset-flow abuse, stale authorization encoded in long-lived tokens.

Controls: adaptive password hashing where used, MFA/passkeys architecture, short-lived access credentials, rotation/revocation, current policy checks, anomaly telemetry.

### T04 — Privilege escalation / hidden administrator bypass

Examples: role assignment without authority, mass assignment of privileges, support impersonation without audit, hard-coded superuser shortcuts.

Controls: explicit privileged capabilities, stronger authentication/approval where required, immutable audit, no name-based bypass.

### T05 — Injection and interpreter abuse

Examples: SQL/query injection, command injection, template/script injection, path traversal, unsafe deserialization, expression/workflow injection.

Controls: typed/parameterized APIs, context-safe encoding, allowlisted operations, input/schema validation, sandboxing where interpreter capability exists.

### T06 — SSRF / unsafe outbound connectivity

Examples: integrations fetching cloud metadata/loopback/private networks, redirect/DNS-rebinding abuse, credential leakage in URLs.

Controls: centralized egress policy, destination validation, scheme/port rules, redirect/DNS/private-address controls, credential isolation.

### T07 — Webhook spoofing/replay

Examples: unsigned callback, stolen webhook secret, replayed payment/shipping event, wrong provider-account-to-tenant binding.

Controls: signature/mTLS/provider auth, timestamp/nonce/event deduplication, payload validation, configured tenant/provider binding.

### T08 — Event/job/workflow replay causing duplicate side effects

Examples: duplicate charge/refund, duplicate inventory reservation, repeated notification or destructive workflow action.

Controls: at-least-once assumption, idempotency keys/inbox, state-machine guards, correlation/causation, privileged audited replay tooling.

### T09 — Financial/payment duplication or integrity loss

Examples: double capture/refund/payout, rounding drift, reconciliation mismatch, tampered ledger transition.

Controls: exact-money semantics, idempotency, immutable/auditable financial transitions, provider reconciliation, approval/separation-of-duty where risk requires.

Severity baseline: **critical** for unauthorized money movement/integrity loss.

### T10 — Module/marketplace compromise

Examples: malicious package, dependency confusion, overbroad declared permissions, update takeover, direct private-table access.

Controls: signed/versioned package identity, provenance/SBOM, permission manifest, review/scanning, runtime isolation where appropriate, revocation/quarantine, dependency boundaries.

### T11 — Software supply-chain / CI credential theft

Examples: malicious dependency/build script/action, poisoned artifact, stolen signing/deploy token, untrusted PR accessing secrets.

Controls: pinned dependencies/actions, least-privilege CI, protected release credentials, SBOM/provenance/signing, immutable artifacts, no secrets in untrusted PR paths.

### T12 — POS/edge device loss or compromise

Examples: stolen terminal, cloned credentials, modified local queue, replayed offline transactions, malicious update.

Controls: device identity, encrypted local sensitive state, credential rotation/revocation, signed updates, store/tenant/device binding, replay-aware sync, remote disable/wipe of managed secrets where feasible.

### T13 — Backup/export leakage

Examples: broadly accessible backups, wrong-tenant export, old backup copied to dev, unbounded signed export URL.

Controls: source classification inheritance, encryption/access control, tenant-safe restore/export, short-lived delivery, audit, lower-environment prohibition.

### T14 — Search/analytics/vector leakage

Examples: hidden documents in autocomplete, stale authorization projection, cross-tenant embeddings, deleted data retained in vector index.

Controls: tenant-aware indexing, field allowlists, result re-authorization, classification propagation, deletion/retention propagation.

### T15 — AI prompt injection / tool abuse

Examples: retrieved content tells an agent to bypass policy, model invents a privileged action, tool parameters target another tenant, model output triggers unsafe automation.

Controls: model is not authority; authorize retrieval before context assembly; tools use governed capabilities; validate parameters; approval for high-impact actions; restricted data excluded by default; audit model/tool decisions.

### T16 — Insider/support misuse

Examples: employee browses customer data, hidden impersonation, bulk export without purpose, key access abuse.

Controls: least privilege, dedicated support capabilities, reason/duration, tenant-visible cues where appropriate, audit/detection, separation of duty for high-risk operations.

### T17 — Resource exhaustion / noisy neighbor

Examples: abusive API tenant, huge reports/search queries, retry storms, unbounded workflows/events, storage explosion.

Controls: quotas, bounded pagination/query complexity, rate limits, backpressure, circuit breakers, tenant/resource budgets, retry limits, workload isolation.

### T18 — DDoS / availability attack

Controls: edge protection/rate limiting, load shedding, caching where safe, capacity/runbooks, graceful degradation of optional capabilities.

### T19 — Dependency/provider outage

Examples: payment/email/AI/shipping provider unavailable, NATS/cache/object store partial outage.

Controls: timeouts, bounded retry/backoff, circuit breakers, durable queues where valid, degraded-mode semantics, provider fallback only when contract/security supports it.

### T20 — Data corruption / migration failure

Examples: bad schema migration, partial backfill, wrong tenant migration, incompatible rollback.

Controls: fresh+upgrade migration tests, backups, expand/contract where needed, migration observability, explicit forward-fix/rollback limits.

### T21 — Region/control-plane failure

Controls: deployment topology appropriate to stage, regional isolation, documented failover/recovery, data-residency-safe restoration, control-plane degradation design.

### T22 — Audit tampering / repudiation

Examples: privileged actor deletes audit record, request correlation lost, support action unattributable.

Controls: protected append-oriented audit semantics, independent authorization, immutable/retained storage strategy as maturity requires, correlation/actor context.

### T23 — Secrets exposure

Examples: git commit, CI log, error trace, client bundle, AI prompt, support screenshot.

Controls: secret stores, scanning, redaction, client/server separation, rotation/revocation incident process.

### T24 — Misconfiguration

Examples: wildcard CORS, disabled TLS verification, public database/bucket, debug mode/stack traces, overly broad OAuth scope.

Controls: secure defaults, typed config validation, fail-closed security config, deployment policy/scanning, config change audit.

## 5. Abuse cases that must remain impossible

```text
A tenant admin reads another tenant by changing tenant_id.
A module writes another module's private table.
An AI agent bypasses permission because its prompt says so.
A replay tool repeats a payment side effect without idempotency/approval.
A support operator impersonates invisibly.
A generic export returns RESTRICTED secrets.
A local/dev bootstrap connects to production by inference.
An untrusted PR receives signing/deployment credentials.
```

## 6. Threat ownership

- shared platform threats -> owning kernel/platform capability;
- business-specific threats -> owning domain with platform controls;
- integration/provider threats -> integration platform + connector owner;
- module/package threats -> module runtime/developer platform/marketplace;
- operational threats -> platform operations/SRE/security;
- AI threats -> intelligence platform + source/tool capability owner.

No threat may be marked "handled elsewhere" without a named owner/control.

## 7. Risk treatment

Allowed dispositions:

```text
mitigate
avoid
transfer
accept_with_explicit_owner
```

Risk acceptance requires documented owner, rationale, compensating controls where applicable and review/expiry for material risks.

## 8. Refinement rule

Every implementation phase/module produces a delta threat model for new trust boundaries, data classes, external providers and privileged capabilities. This foundation model is the minimum, not a substitute for feature-specific threat modeling.