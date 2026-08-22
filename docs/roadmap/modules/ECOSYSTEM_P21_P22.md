# Ecosystem Module Dossiers — P21 to P22

Status: **Mandatory future planning baseline**

## P21 — Developer Platform

Architecture: external developers receive versioned, least-privilege contracts around the same runtime/module/capability/event rules used internally. Generated SDKs/docs are derivative artifacts; canonical schemas/contracts remain source of truth.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P21.A | Public API Surface | approved capability -> API contract -> auth/policy -> execution | version policy, pagination, rate limits |
| P21.B | API Credentials/Apps | developer app -> scopes -> credential/OAuth registration -> rotate/revoke | grant types, scopes, redirect policy |
| P21.C | SDK Generation | canonical OpenAPI/event schemas -> generator -> language SDK -> package test | supported languages, generator versions |
| P21.D | Developer CLI | auth/context -> generate/validate/test/package/deploy-safe commands | profiles, output modes, no hidden prod defaults |
| P21.E | Module Generator | requested module ID -> governed skeleton/manifest/tests -> validate | templates, runtime, optional UI scaffolding |
| P21.F | Sandbox Tenant | provision isolated sandbox -> seed fixtures -> expire/reset | lifetime, quotas, allowed providers |
| P21.G | Contract/Event Explorer | registry -> searchable docs/examples -> version diff | visibility, deprecation indicators |
| P21.H | Local Harness | pinned dependencies -> disposable platform -> module/connector tests | service versions, reset policy |
| P21.I | Developer Docs | canonical contracts -> examples/tutorials/reference -> versioned publication | doc version, language, examples |
| P21.J | Compatibility Tooling | package/contract -> target platform versions -> compatibility report | supported baselines, strict/warn classes |
| P21.K | Webhook/Test Utilities | synthetic event/request -> signed delivery/simulation -> trace | event types, destination allowlist |

Security: developer credentials are restricted, sandbox isolation is mandatory, generated examples contain synthetic data only, API docs never expose private/internal endpoints by accident.

## P22 — Omnexa Exchange / Marketplace

Architecture: marketplace packages are declared, validated and signed artifacts installed through P03/P16/P12/P17/P20 extension contracts. Marketplace never grants packages undeclared platform authority.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P22.A | Publisher Identity | register publisher -> verify -> status/risk lifecycle | verification levels, organization linkage |
| P22.B | Package Manifest | upload package -> parse scopes/dependencies/version -> validate | package types, size, supported platform versions |
| P22.C | Signing & Integrity | package digest -> publisher/platform signature -> verification chain | algorithms, key rotation, trust roots |
| P22.D | Automated Validation | static/contracts/security/license/test sandbox -> report | required gates by package type |
| P22.E | Review & Approval | validation result -> human/policy review -> approve/reject | risk classes, reviewer requirements |
| P22.F | Listing & Discovery | approved release -> metadata/categories/search -> listing | categories, regions, visibility |
| P22.G | Install/Consent | tenant selects package -> show declared scopes -> approve -> P03 install | scope consent, admin role, dependencies |
| P22.H | Upgrade/Rollback | compatible update -> preview changes/scopes -> migrate -> activate/rollback | auto-update policy, maintenance window |
| P22.I | Revoke/Suspend | security/legal/platform signal -> suspend listing/install/runtime policy | severity, tenant notification |
| P22.J | Ratings/Release Metadata | verified install/use context -> rating/review/release notes | moderation, eligibility |
| P22.K | Commercial Entitlement Boundary | purchase/license fact -> entitlement -> install/use check | pricing/license models only when authorized |
| P22.L | Package-type Profiles | module/connector/theme/workflow/AI tool/country pack -> specialized validation | allowed manifests/scopes per type |

Package flow:

```text
publisher verification
 -> package + manifest
 -> digest/signature
 -> automated gates
 -> policy/human review
 -> listing
 -> tenant scope/dependency consent
 -> governed install
 -> health/usage
 -> update/revoke lifecycle
```

Forbidden: arbitrary installer scripts with unrestricted host access, undeclared network destinations/secrets, hidden permissions, direct marketplace DB writes into domain schemas, unsigned production releases once signing is required.

## Common ecosystem options

Every developer/marketplace option declares platform-vs-tenant scope, security implications, compatibility range, auditability and revocation behavior. Convenience settings cannot weaken package validation, tenant consent or signature checks.