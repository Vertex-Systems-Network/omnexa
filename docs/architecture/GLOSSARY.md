# Omnexa Product & Domain Glossary

Status: **Canonical terminology v1**  
Work package: **P00.02**

This glossary defines the language that all humans, AI systems, code, schemas, APIs, events, documentation and UI concepts must use unless a later ADR explicitly supersedes a term.

## 1. Platform terms

| Term | Canonical meaning |
|---|---|
| Omnexa | The complete Composable Enterprise Business Operating System. |
| Platform | The full Omnexa runtime, kernel, domain modules, developer surfaces and operational infrastructure. |
| Kernel | Shared platform capabilities required by multiple modules; not a business domain. |
| Domain | A bounded business or platform responsibility with explicit ownership. |
| Module | An independently governed package that contributes capabilities, data, events, permissions and UI surfaces. |
| Capability | A stable, authorized operation exposed by an owning domain/module. |
| Contract | A versioned interface consumed outside the owning implementation boundary. |
| Extension | Additive functionality that consumes public contracts without owning the base domain. |
| Connector | An integration adapter between Omnexa and an external system. |
| Country Pack | Localized tax, document, regulatory, payment or business rules for a country/region. |
| Industry Pack | A governed composition of existing modules/configuration for a vertical. |
| Tenant | The primary isolation boundary for customer-controlled data and configuration. |
| Organization | A business structure that exists inside a tenant. |
| Legal Entity | An organization unit that has legal/accounting identity. |
| Business Unit | Operational grouping under an organization/legal entity. |
| Branch | A managed operating branch/location under an organization. |
| Team | A people/access grouping; not a legal entity. |
| Workspace | UI/product context only when a feature explicitly needs a collaborative workspace; it must never silently replace Tenant or Organization semantics. |

## 2. Identity and party terms

| Term | Canonical meaning |
|---|---|
| Identity | Authentication/security principal recognized by Omnexa. |
| User | A human identity with interactive access. |
| Service Account | A non-human principal used by services/integrations. |
| Party | A business-reference abstraction for a person or organization participating in business relationships. |
| Person | A natural person record in the business model; not automatically an authenticated User. |
| Organization Party | An external/internal organization represented in the business model. |
| Customer | A role/relationship of a Party with a selling organization; not a universal duplicate person table. |
| Supplier | A role/relationship of a Party providing goods/services. |
| Employee | A workforce relationship between a Person and an Organization. |
| Contact | A context-specific person/contact relationship; it must not become a second canonical identity model. |

## 3. Commerce and finance terms

| Term | Canonical meaning |
|---|---|
| Product | A catalog item offered or managed by a business. |
| Variant | A sellable/configurable variation of a Product. |
| Service | A non-stock or service-type commercial offering. |
| Catalog | A governed collection of products/services and commercial metadata. |
| Price List | Contextual commercial pricing rules/amounts. |
| Cart | Pre-order commerce state. |
| Order | A committed commercial request/transaction document owned by Commerce. |
| Invoice | A financial receivable/payable document owned by Finance/Billing, not a synonym for Order. |
| Payment | Movement/authorization of funds through the Payments domain. |
| Transaction | Generic term only when domain-qualified; avoid using it as an ambiguous business entity name. |
| Ledger Entry | Atomic accounting posting owned by Finance. |
| Inventory | Quantity/state of stockable items. |
| Stock Movement | Inventory-owned quantity transition. |
| Shipment | Logistics fulfillment movement; not an Order state substitute. |
| Subscription | Recurring commercial agreement/state with explicit billing lifecycle. |

## 4. Workflow and integration terms

| Term | Canonical meaning |
|---|---|
| Event | Immutable statement that something meaningful already happened. |
| Command | Request for an owning capability to attempt an action. |
| Query | Read request that does not mutate business state. |
| Workflow | Durable, versioned orchestration of steps, waits, approvals and compensations. |
| Automation | User/configuration-driven trigger/action behavior; may be backed by workflows. |
| Job | Background unit of technical work; not a business event. |
| Webhook | External HTTP delivery mechanism for a versioned event/notification contract. |
| Projection | Derived/read-optimized view built from owned source data/contracts/events. |
| Outbox | Reliable publication pattern coupling domain commit with later message delivery. |
| Inbox | Reliable/idempotent consumption record for external/domain messages. |

## 5. AI terms

| Term | Canonical meaning |
|---|---|
| AI Gateway | Governed model/provider routing and policy layer. |
| Agent | AI-driven actor that plans/executes only through authorized Omnexa capabilities. |
| Tool | Capability exposed to an AI/automation runtime. |
| Knowledge Base | Governed retrievable content source for AI/search. |
| RAG | Retrieval-augmented generation over approved data sources. |
| Human Approval | Explicit human gate before a protected or policy-defined action executes. |

## 6. Lifecycle terms

| Term | Meaning |
|---|---|
| Install | Register schema/contracts/configuration required for a module. |
| Enable | Make an installed module available for permitted use. |
| Disable | Stop normal module participation while preserving its owned data/state. |
| Suspend | Temporarily block operational behavior without declaring archival/removal. |
| Archive | Preserve historical state while ending active use. |
| Detach | Remove active platform linkage after safe export/archive requirements are met. |
| Purge | Explicit destructive deletion according to retention/legal rules. |
| Upgrade | Move a module/schema/contracts forward through a supported version transition. |

## 7. Terms that must not be used ambiguously

The following words require qualification or are forbidden as primary model names because they cause architecture drift:

- `account` — qualify as user account, financial account, customer account, merchant account, etc.;
- `company` — use Organization or Legal Entity unless the specific UI/business term is intentionally Company;
- `transaction` — qualify by domain;
- `record`, `item`, `data` — never use as canonical business entity names;
- `workspace` — never use as a hidden synonym for tenant;
- `client` — use Customer for commercial role, API Client for software, or explicitly defined domain meaning;
- `member` — qualify as team member, organization member, subscription member, etc.;
- `status` — domain states must be named and versioned, not shared as one global enum;
- `admin` — a role/permission composition, never a bypass of authorization;
- `global` — must state whether it means platform-wide, tenant-wide, organization-wide or locale-independent.

## 8. Naming rule for new concepts

Before creating a new domain noun, an implementer must answer:

1. Does the concept already exist in this glossary or another owning domain?
2. Is it an entity, role, capability, event, projection, setting or UI-only concept?
3. Which domain owns the authoritative write model?
4. What scope applies: platform, tenant, organization, legal entity, branch, team, user or external party?
5. Is the proposed name unambiguous outside its current screen/module?

If ownership or meaning remains ambiguous, implementation stops and the glossary/ADR is resolved first.
