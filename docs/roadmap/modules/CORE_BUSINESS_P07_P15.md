# Core Business Module Dossiers — P07 to P15

Status: **Mandatory future planning baseline**

These dossiers define intended submodule boundaries, flows and option surfaces before implementation. They do not authorize business-feature code while `business_feature_code_authorized=false`.

## P07 — CRM, Sales & Customer 360

Architecture: CRM owns sales/customer relationship state; finance owns accounting, commerce owns orders, communications/connectors own provider transport. Customer 360 is a read projection, not a second write owner.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P07.A | Leads | capture/import -> qualify -> convert/disqualify | lead sources, statuses, assignment rules |
| P07.B | Contacts & Organizations | create/update party-linked CRM profile -> merge/archive | duplicate policy, required fields, ownership model |
| P07.C | Opportunities & Pipelines | create -> stage progression -> won/lost | pipeline/stage definitions, probability, required gates |
| P07.D | Activities | schedule/log call/email/meeting/task -> complete/cancel | activity types, reminders, visibility |
| P07.E | Territory & Ownership | rule/manual assignment -> owner/team/territory | assignment priority, round-robin hooks |
| P07.F | Quotes/Proposals Boundary | opportunity -> request commercial quote/proposal capability -> track reference | template/provider selection, approval thresholds |
| P07.G | Customer 360 Projection | consume authorized domain facts -> project timeline/summary | source visibility, freshness, retention |
| P07.H | Scoring | inputs -> score model -> bounded recommendation | model/version, thresholds, explainability fields |
| P07.I | Sequences | enroll -> scheduled governed actions -> pause/complete | cadence, stop rules, channel availability |

Permissions: lead/contact/opportunity read/create/update/delete/assign; pipeline/admin; scoring/manage; sequence/manage. Events: lead.created/converted, contact.changed, opportunity.stage.changed/won/lost, activity.completed, owner.changed.

Delivery: B foundation -> A/C/D -> E -> F/G -> H/I. Tests include ownership, merge/duplicate, pipeline invariants, permission/tenant isolation, projection no-write rule, sequence idempotency.

## P08 — Finance & ERP Core

Architecture: finance is authoritative for accounting facts. Posted ledger history is immutable except governed reversal/correction entries.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P08.A | Chart of Accounts | define account -> validate hierarchy/type -> activate/archive | account classes, code rules, control accounts |
| P08.B | General Ledger | approved source -> balanced posting -> immutable ledger | fiscal periods, posting controls, dimensions |
| P08.C | Journal Engine | draft journal -> validate balance -> approve/post/reverse | approval thresholds, numbering, reversal rules |
| P08.D | Accounts Receivable | invoice receivable -> allocation -> aging/settlement | payment terms, aging buckets, credit controls |
| P08.E | Accounts Payable | supplier obligation -> approval -> payment allocation | approval routing, due-date policy |
| P08.F | Invoices & Credit Notes | draft -> calculate -> approve/post -> issue -> credit/reverse | numbering, tax/display policy, document templates |
| P08.G | Expenses | submit -> policy validation -> approval -> posting/reimbursement | categories, limits, receipt requirements |
| P08.H | Cash/Bank & Reconciliation | import/record transaction -> match -> reconcile | matching tolerances, statement formats |
| P08.I | Tax Integration | taxable transaction context -> tax engine/provider -> tax lines | jurisdiction/provider, rounding/reporting policy |
| P08.J | Budgeting | budget version -> allocate -> approve -> compare actuals | periods, dimensions, revision policy |
| P08.K | Fixed Assets | acquire -> capitalize -> depreciate -> dispose | depreciation methods, useful-life policies |
| P08.L | Procurement Accounting Boundary | procurement fact -> payable/commitment/posting request | posting maps, accrual policy |

Invariants: debits=credits, closed period controls, currency/rounding exactness, posted-history immutability, idempotent source posting. Events: journal.posted/reversed, invoice.posted/paid/credited, expense.approved, reconciliation.completed.

Delivery: A/B/C -> D/E/F -> G/H/I -> J/K/L. Evidence includes representative multi-currency, period close, reversal, duplicate-post prevention and reconciliation tests.

## P09 — Commerce OS

Architecture: commerce owns catalog-selling/order lifecycle. Inventory owns stock movement, payments own payment state, finance owns ledger/tax postings.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P09.A | Catalog & Products | create product -> publish/channel eligibility -> archive | product types, status, channels |
| P09.B | Variants & Options | option definitions -> variant matrix/SKU -> availability | option types, variant limits, SKU generation policy |
| P09.C | Collections/Categories | rule/manual membership -> navigation/filter exposure | hierarchy, dynamic rules, sort policy |
| P09.D | Price Lists | product/variant + context -> resolve price | currency, customer/channel/quantity scope |
| P09.E | Promotions & Discounts | cart context -> eligible rule -> calculated adjustment | coupon/automatic, stacking, limits |
| P09.F | Cart | add/update/remove -> price/availability snapshot -> ready checkout | expiry, guest/customer policy, quantity limits |
| P09.G | Checkout | cart -> identity/address/shipping/tax/payment orchestration -> order | required steps, guest checkout, terms/consent |
| P09.H | Orders | place -> confirm -> fulfillment/cancel/complete | status transitions, edit/cancel windows |
| P09.I | Returns/Refund Orchestration | request -> approve/receive -> inventory/payment/finance commands | return windows, reasons, disposition |
| P09.J | Subscriptions | plan -> subscribe -> renew/pause/cancel | billing cadence, proration hooks, retry policy |
| P09.K | Bundles/Kits | define bundle -> price/availability resolution -> order expansion | fixed/dynamic contents, pricing policy |
| P09.L | B2B/Wholesale | account terms -> catalog/price/quantity/order rules | MOQ, credit/payment terms, approval |
| P09.M | Marketplace/Multi-vendor Boundary | seller listing/order allocation -> settlement references | seller approval, commission policy hooks |

Events: product.published, cart.updated, checkout.completed/failed, order.placed/cancelled/fulfilled, return.approved, subscription.renewed/cancelled. Tests prove no direct inventory/payment/finance writes.

Delivery: A/B/C/D -> E/F -> G/H -> I -> J/K/L -> M boundary.

## P10 — Payment Fabric

Architecture: provider-neutral state machine around external payment providers. Raw card secrets remain outside Omnexa where possible; token/provider references are not reusable secrets unless explicitly classified.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P10.A | Provider Registry | configure provider adapter -> validate capabilities -> activate | provider priority, environment, supported methods |
| P10.B | Payment Intent/Authorization | create intent -> provider authorize -> normalized state | capture mode, expiry, idempotency |
| P10.C | Capture/Void | authorized payment -> capture/void -> normalized result | partial capture, capture deadline |
| P10.D | Refunds | refund request -> policy -> provider -> normalized result | partial/multiple refunds, approval threshold |
| P10.E | Token/Reference Boundary | provider token/reference -> store restricted reference -> revoke | retention, provider scope |
| P10.F | Webhooks | receive -> authenticate -> dedupe -> normalize -> process | signature schemes, replay window |
| P10.G | Recurring/Mandates | mandate setup -> activate -> charge/use -> revoke | consent, mandate types, expiry |
| P10.H | Settlement | provider batch -> ingest -> normalize -> settlement record | statement cadence, currency handling |
| P10.I | Reconciliation | internal payments vs settlement -> match/exceptions | tolerances, auto-match rules |
| P10.J | Disputes/Chargebacks | provider dispute -> case lifecycle -> evidence/status | deadlines, evidence workflow |
| P10.K | Payouts | approved payable balance -> payout -> provider status | schedule, minimums, destination validation |

Required flows: auth/capture, auth/void, partial refund, duplicate webhook, provider timeout then webhook reconciliation, settlement mismatch. At least two adapters must prove provider neutrality.

## P11 — POS & Edge

Architecture: POS is an edge client/runtime backed by server contracts; offline state is a durable queue, not an alternate authoritative business database.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P11.A | POS Runtime | operator/session -> basket -> tender -> sale receipt | store/location profile, UI policy |
| P11.B | Device Registry/Adapters | discover/register device -> capability -> health/use | allowed device classes, driver/provider |
| P11.C | Local Transaction Queue | offline command -> durable local queue -> sync -> acknowledge | queue bounds, encryption, retention |
| P11.D | Offline Policy | connectivity state + risk rules -> allow/deny operation | offline amount/item limits, stale-data tolerance |
| P11.E | Sync & Reconciliation | reconnect -> ordered upload/download -> conflict handling | batching, conflict strategy |
| P11.F | Shift/Session | open shift -> cash/session operations -> close/reconcile | cash float, manager approval |
| P11.G | Receipt | sale result -> receipt model -> print/email/export | template, numbering, localization |
| P11.H | Barcode/Scanner | input -> resolve product/action -> validate | barcode types, fallback behavior |
| P11.I | Cash Drawer/Printer/Scale | normalized command -> adapter -> result | device routing, timeout |
| P11.J | Payment Terminal Integration | sale intent -> payment fabric/terminal adapter -> result | terminal selection, offline terminal policy |

Tests: offline sale/reconnect/replay, duplicate sync, conflict, power loss during queue write, shift discrepancy, device unavailable degradation.

## P12 — Experience Builder & CMS

Canonical detailed decomposition is in `docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md` P12.A-P12.L and remains authoritative. This dossier adds cross-submodule flow and option governance.

Architecture layers: content/page schema -> version validation -> runtime renderer -> block/component/template/theme registries -> authoring APIs -> visual builder -> preview -> publish lifecycle. Runtime behavior must never exist only in client-side editor code.

Primary authoring flow: site/channel -> choose/create page -> edit block tree -> validate data bindings/permissions/accessibility -> preview -> review/version -> publish/schedule -> cache/search/event hooks. Primary template flow: create template type -> define slots/data schema -> validate dependencies -> preview -> activate/assign -> version/restore. Primary CMS flow: define content type -> migrate schema/version -> edit localized content -> validate -> publish -> API/search projection.

Options must be schema-driven: design tokens, breakpoints, editor grid/guides, autosave interval, allowed block set, content-type fields, publishing approvals, preview policy, SEO defaults, locale/RTL, form limits, import/export collision policy. Options cannot bypass permissions or server validation.

UI gates: WCAG 2.2 AA target, semantic/native HTML, keyboard equivalent for drag/drop operations, focus management, reduced motion, W3C validation, WAVE evaluation and manual screen-reader/zoom-reflow evidence.

## P13 — Portal Platform

Architecture: one portal runtime with persona profiles; no duplicated auth stacks.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P13.A | Portal Runtime | authenticated scoped actor -> portal profile -> capabilities/navigation -> page | host/channel, session integration |
| P13.B | Portal Profile Builder | define persona/profile -> capability/menu/theme mapping -> activate | persona, branding, modules, locales |
| P13.C | Customer Portal | customer context -> documents/orders/requests/status capabilities | allowed capabilities, self-service actions |
| P13.D | Vendor Portal | vendor context -> procurement/docs/status capabilities | allowed organizations, workflows |
| P13.E | Employee Portal | employee context -> HR/project/service capabilities | role/organization scope |
| P13.F | Partner Portal | partner context -> shared capabilities | partner tier/entitlement |
| P13.G | Capability Navigation | available modules + permissions -> deterministic menu/routes | ordering, labels, feature flags |
| P13.H | Self-service Surfaces | capability query/command -> status/document/request projection | upload/download rules, action confirmation |

Tests: same runtime across personas, no capability leakage, module disable graceful navigation degradation, localization/accessibility.

## P14 — HR, Projects & Service Operations

Architecture: HR, project operations and service operations remain distinct write models connected through stable employee/project/customer/reference capabilities and finance integration.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P14.A | Workforce Records | hire/create worker -> assignments -> status changes -> separation | worker types, required records |
| P14.B | Leave | request -> entitlement/policy -> approval -> balance/update | leave types, accrual, approval |
| P14.C | Attendance & Time | clock/event/manual entry -> validation -> timesheet/attendance record | schedules, rounding, grace rules |
| P14.D | Projects | create -> scope/team/budget refs -> active -> close/archive | project types, statuses, numbering |
| P14.E | Tasks | create/assign -> progress -> complete | workflow/status templates, priorities |
| P14.F | Timesheets | time entry -> validate -> approve -> costing/export | periods, approval, billable rules |
| P14.G | Resource Scheduling | demand -> availability/skills -> allocation -> conflict handling | capacity, working calendars |
| P14.H | Service Tickets/Cases | intake -> classify -> assign -> SLA work -> resolve/close | priorities, queues, SLA policy |
| P14.I | Service Scheduling | service job -> resource/time/location -> dispatch/complete | territory, travel buffers |
| P14.J | Payroll Integration Boundary | approved HR/time facts -> country/provider payroll adapter/export | pay periods, export mapping |

Delivery: A -> B/C -> D/E/F/G -> H/I -> J boundary. Tests emphasize policy dates/timezones, approval permissions, costing handoff and no finance direct writes.

## P15 — Supply Chain, Warehouse & Manufacturing

Architecture: inventory movement ledger is authoritative for quantity state; commerce/procurement/manufacturing request movements through capabilities. Finance consumes valuation/posting facts.

| ID | Submodule | Primary flow | Key options |
|---|---|---|---|
| P15.A | Inventory Ledger & Movements | authorized movement command -> validate -> append movement -> derive balance | costing/availability policy hooks, negative-stock policy |
| P15.B | Warehouses & Bins | define hierarchy -> capacity/status -> assign stock location | bin rules, putaway/pick zones |
| P15.C | Suppliers | supplier reference/profile -> sourcing status | qualification/status rules |
| P15.D | Purchasing | requisition -> approval -> PO -> receive/close | approval thresholds, terms |
| P15.E | Transfers | request -> reserve/ship -> receive -> reconcile | in-transit policy, partial receipt |
| P15.F | Fulfillment | order demand -> allocate/pick/pack/ship | allocation strategy, partial fulfillment |
| P15.G | Logistics/Shipping Boundary | shipment -> rate/label/provider -> tracking events | carrier/service selection, package rules |
| P15.H | BOM | versioned components/quantities -> approve/activate | effectivity, substitutes |
| P15.I | Routings | operation sequence/resources -> version/activate | work centers, setup/run standards |
| P15.J | Production Orders | demand -> plan/release -> consume/produce -> close | lot sizing, backflush policy |
| P15.K | Quality | inspection plan -> sample/result -> accept/hold/reject | sampling, disposition rules |
| P15.L | Maintenance | asset/equipment -> preventive/corrective work -> complete | schedules, meter/calendar triggers |
| P15.M | Lot/Serial Traceability | receive/produce -> assign lot/serial -> movements -> genealogy | tracking requirement by item class |

Tests: no quantity mutation outside movement ledger, duplicate receipt/movement prevention, transfer/production restart safety, lot genealogy, inventory-finance boundary.

## Common core-business configuration rule

Every option must declare owning module/submodule, scope, default, validation, permission, audit, effective date and migration/compatibility implications. Financial/security/authorization invariants cannot be downgraded to ordinary user-editable settings.

## Common UI/quality rule

All future browser surfaces use the web accessibility plan: WCAG 2.2 AA target, W3C/WAVE plus manual evidence. All business actions must expose loading/empty/error/success/disabled states and server-side authorization; client hiding is never authorization.