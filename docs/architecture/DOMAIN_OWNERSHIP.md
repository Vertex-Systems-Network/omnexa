# Omnexa Domain Ownership Registry

Status: **Canonical ownership baseline v1**  
Work package: **P00.02**

Every authoritative write model, business rule and public capability has exactly one owning domain. Other domains consume contracts, events, workflows or approved projections.

## 1. Ownership rules

1. One authoritative write owner per concept.
2. Consumers never write another domain's tables directly.
3. Shared vocabulary does not imply shared write ownership.
4. A projection/cache is not authoritative source data.
5. Ownership changes require an ADR, migration/compatibility plan and registry update.
6. Kernel domains may own platform primitives, but must not absorb business-domain behavior.
7. AI and automation act through owner-published capabilities.

## 2. Platform/kernel ownership

| Capability / concept | Owner | Notes |
|---|---|---|
| identity.user | `kernel.identity` | Authentication identity, not business Person. |
| identity.service_account | `kernel.identity` | Non-human principal. |
| tenant.tenant | `kernel.tenancy` | Primary isolation boundary. |
| organization.organization | `kernel.organization` | Organization hierarchy. |
| organization.legal_entity | `kernel.organization` | Legal/accounting organizational identity. |
| organization.branch | `kernel.organization` | Organizational operating branch. |
| authorization.policy | `kernel.authorization` | Central enforcement semantics. |
| authorization.role | `kernel.authorization` | Role composition. |
| audit.record | `kernel.audit` | Immutable attributable audit record. |
| file.object | `kernel.files` | Generic managed binary/file primitive. |
| notification.delivery | `kernel.notifications` | Cross-platform delivery primitive. |
| module.manifest | `kernel.modules` | Module identity/lifecycle metadata. |
| module.lifecycle | `kernel.modules` | Install/enable/disable/etc. |
| configuration.setting | `kernel.configuration` | Scoped platform/module settings. |
| feature_flag | `kernel.configuration` | Governed feature state. |
| job.execution | `kernel.jobs` | Technical background execution. |
| workflow.definition | `platform.workflow` | Durable workflow model. |
| workflow.execution | `platform.workflow` | Workflow runtime state. |
| integration.connector | `platform.integration` | External integration metadata/runtime contract. |
| developer.application | `platform.developer` | OAuth/API application identity. |

## 3. Universal business foundation ownership

| Capability / concept | Owner | Notes |
|---|---|---|
| party.party | `business.party` | Business-reference abstraction. |
| party.person | `business.party` | Natural-person business record. |
| party.organization | `business.party` | Party representation; distinct from tenant organization hierarchy. |
| party.address | `business.party` | Reusable business address/contact primitive. |
| party.contact_point | `business.party` | Email/phone/etc. business contact data. |
| customer.relationship | `business.party` | Canonical customer role/relationship; CRM enriches engagement. |
| supplier.relationship | `business.party` | Canonical supplier role/relationship; procurement enriches operations. |

## 4. Domain ownership baseline

| Capability / concept | Owner |
|---|---|
| crm.lead | `crm.sales` |
| crm.opportunity | `crm.sales` |
| crm.pipeline | `crm.sales` |
| crm.activity | `crm.sales` |
| crm.customer_360_projection | `crm.sales` |
| finance.account | `finance.ledger` |
| finance.journal | `finance.ledger` |
| finance.ledger_entry | `finance.ledger` |
| finance.invoice | `finance.billing` |
| finance.credit_note | `finance.billing` |
| finance.receivable | `finance.billing` |
| finance.payable | `finance.payables` |
| commerce.product | `commerce.catalog` |
| commerce.variant | `commerce.catalog` |
| commerce.collection | `commerce.catalog` |
| commerce.price_list | `commerce.pricing` |
| commerce.cart | `commerce.checkout` |
| commerce.checkout | `commerce.checkout` |
| commerce.order | `commerce.orders` |
| commerce.return | `commerce.orders` |
| inventory.stock | `inventory.stock` |
| inventory.reservation | `inventory.stock` |
| inventory.stock_movement | `inventory.stock` |
| procurement.purchase_order | `procurement.core` |
| warehouse.operation | `warehouse.core` |
| logistics.shipment | `logistics.core` |
| payments.payment | `payments.core` |
| payments.authorization | `payments.core` |
| payments.refund | `payments.core` |
| payments.settlement | `payments.core` |
| pos.sale_session | `pos.core` |
| pos.shift | `pos.core` |
| hr.employment | `hr.core` |
| hr.leave | `hr.core` |
| payroll.compensation | `payroll.core` |
| project.project | `projects.core` |
| project.task | `projects.core` |
| service.case | `service.core` |
| cms.content_type | `experience.cms` |
| cms.content_entry | `experience.cms` |
| website.site | `experience.builder` |
| website.page | `experience.builder` |
| website.component_definition | `experience.builder` |
| portal.configuration | `experience.portal` |
| analytics.metric_definition | `data.analytics` |
| analytics.semantic_model | `data.analytics` |
| ai.agent_definition | `intelligence.agents` |
| ai.knowledge_base | `intelligence.knowledge` |
| ai.model_route | `intelligence.gateway` |

## 5. Important ownership distinctions

### User vs Person

`kernel.identity` owns authenticated users. `business.party` owns business persons. Linking is explicit; neither domain writes the other's internal state.

### Tenant Organization vs Party Organization

`kernel.organization` describes the customer's internal access/operating hierarchy. `business.party` represents organizations participating in business relationships. A legal entity may have an explicit reference to a party record, but the concepts are not merged into one universal table.

### Customer vs CRM

Customer relationship authority belongs to `business.party`; CRM owns lead/opportunity/engagement behavior and customer-360 projections. CRM is not allowed to become the sole identity/customer database for the platform.

### Order vs Invoice vs Payment

`commerce.orders` owns orders, `finance.billing` owns invoices/receivables, and `payments.core` owns payment execution/state. Each references the others through public identifiers/contracts/events; no domain writes another's state directly.

### Product vs Stock

`commerce.catalog` owns sellable catalog/product definitions. `inventory.stock` owns stock quantities/reservations/movements. Inventory references product/catalog identifiers without taking catalog ownership.

## 6. Registry extension rule

Before a new module/entity is implemented, its authoritative owner must be added here or to a later machine-readable ownership registry. Duplicate ownership is a blocking architecture defect.
