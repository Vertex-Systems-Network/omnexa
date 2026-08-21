# Omnexa Security Control Matrix

Status: **Canonical v1 control map**  
Work package: **P00.06**

This matrix assigns platform security responsibilities so modules and AI systems do not duplicate, omit or bypass controls. `Kernel` means a shared platform capability that future implementation must provide; `Domain` means the owning business/platform module remains responsible for its business-specific decision.

| Control area | Primary owner | Domain responsibility | Required implementation evidence later |
|---|---|---|---|
| Authentication | `kernel.identity` | declare allowed principal/login flows | auth flow + negative tests |
| Session/token lifecycle | `kernel.identity` | no custom token stores | expiry/rotation/revocation tests |
| Authorization engine | `kernel.authorization` | declare capabilities/relationships/policies | allow + deny + cross-scope tests |
| Tenant context | `kernel.organization` + authorization | persist/use `tenant_id` correctly | cross-tenant negative tests |
| Organization scope | `kernel.organization` + domain | declare record scope | cross-org negative tests |
| Audit platform | `kernel.audit` | emit business/security audit facts | immutable/privileged mutation tests |
| Secrets | platform security/runtime | reference secrets only | secret scan + rotation tests |
| Key management | platform security/runtime | no custom crypto/key stores | key access/rotation tests |
| TLS/transport security | platform runtime/edge | no TLS bypass | deployment/integration tests |
| Data classification | owning domain + security governance | classify fields/resources | registry/schema validation |
| Data masking | API/UI + owning domain | declare sensitive fields | field disclosure tests |
| Bulk export | owning domain + platform export capability | explicit export permission | scope/audit/field tests |
| Retention/deletion | owning domain + platform policy | map records to retention policy | retention/purge/legal-hold tests |
| File security | platform files/media | declare classification/publication | upload/download authorization tests |
| API security | gateway + owning domain | object/action authorization | contract/security tests |
| Event security | event fabric + producer/consumer | tenant context/idempotent handler | schema/tenant/replay tests |
| Webhooks | integration platform + connector | provider authenticity and mapping | signature/replay tests |
| Outbound egress/SSRF | integration/runtime platform | declare external destinations | SSRF/allowlist tests |
| Search isolation | platform search + domain projection | authorize indexed/result fields | cross-tenant/result leakage tests |
| Analytics isolation | data platform + source domain | minimize/classify data | tenant/export/delete propagation tests |
| AI retrieval | intelligence platform + source domain | expose only authorized sources | tenant/object retrieval tests |
| AI tool execution | intelligence platform + capability owner | capability/policy enforcement | tool allow/deny/approval tests |
| Workflow execution | workflow platform + capability owner | identity/approval/idempotency rules | delayed-step/current-policy tests |
| Module permissions | module runtime | declare permissions/dependencies | install/upgrade permission diff tests |
| Package integrity | developer/marketplace platform | no bypass | signature/provenance tests |
| CI credentials | CI/release platform | least-privilege workflow permissions | permission/secret exposure checks |
| Database privilege | platform runtime + schema owner | no cross-domain private writes | role/query boundary tests |
| Cache/broker isolation | platform runtime | safe key/subject/data use | auth/tenant/retention tests |
| POS/edge identity | edge platform | store/device capability policy | enrollment/revocation/offline tests |
| Support impersonation | governance/identity | no hidden bypass | reason/duration/audit/restriction tests |
| Security exceptions | governance | request named exception only | expiry/owner/compensating-control check |

## Required rule

A domain may implement stricter controls than this matrix, but it may not replace a shared kernel control with a private weaker alternative. If a control has no clear owner, implementation stops until ownership is resolved.

## Deny-by-default areas

The following are deny/prohibit by default unless an explicit capability/policy says otherwise:

- cross-tenant access;
- production data use in lower environments;
- `RESTRICTED` data in AI/model input;
- `RESTRICTED` data in logs/search/analytics;
- generic export of restricted fields;
- module private cross-domain writes;
- external network destinations;
- support impersonation;
- destructive purge;
- event replay with protected side effects;
- privileged AI actions without required approval;
- unsigned/untrusted extension execution once marketplace signing exists.

## Evidence principle

A security requirement is not considered implemented because configuration or code appears plausible. Later implementation phases must produce executable positive and negative evidence appropriate to the control, and P00.07 defines the CI/release gate taxonomy.
