# Omnexa Naming Standard

Status: **Canonical v1**  
Work package: **P00.02**

This standard prevents different teams and AI systems from inventing incompatible names for the same concept.

## 1. General rules

- Use English for code, contracts, schema, events, permissions and technical documentation identifiers.
- Product/UI translations are localization concerns and do not rename canonical technical identifiers.
- Prefer explicit domain nouns over generic words such as `data`, `item`, `manager`, `service`, `record` or `transaction`.
- Names must expose ownership and scope when ambiguity is possible.
- Do not encode implementation technology into domain names.
- Do not create abbreviations unless they are globally established or explicitly registered in this document.

## 2. Repository/domain IDs

Canonical module/domain IDs use lowercase dot-separated names:

```text
kernel.identity
kernel.organization
kernel.authorization
kernel.audit
platform.workflow
platform.integration
business.party
crm.sales
finance.ledger
commerce.catalog
commerce.orders
inventory.stock
payments.core
experience.cms
```

Rules:

- first segment identifies family/ownership class;
- IDs are stable public identifiers and are versioned only when semantics break;
- display names may change without changing IDs;
- never use tenant/customer names in module IDs;
- third-party extensions use an approved vendor namespace.

## 3. Go naming

- Packages: short lowercase names with one clear responsibility.
- Exported types/functions: PascalCase.
- Internal identifiers: idiomatic Go camelCase.
- Interfaces describe capability, not implementation (`Authorizer`, `EventPublisher`).
- Avoid `Manager`, `Helper`, `Utils`, `Common`, `Misc` as architectural package names.

## 4. TypeScript naming

- Components/classes/types/interfaces: PascalCase.
- Functions/variables: camelCase.
- Constants: use descriptive names; SCREAMING_SNAKE_CASE only for true process/environment constants.
- React components match file/component names.
- Shared UI packages must not import domain internals.

## 5. Database naming

- PostgreSQL schemas belong to owning modules/domains.
- Tables/columns: lowercase `snake_case`.
- Primary keys: `id` and follow `IDENTIFIER_STANDARD.md` unless an accepted ADR defines a stronger domain-specific composite identity.
- Foreign references within an owning schema: `<entity>_id`.
- Cross-domain database foreign keys are forbidden by default; use owned references/contracts according to architecture rules.
- Tenant-owned rows use canonical `tenant_id`; organization/legal-entity/business-unit/branch/team/user scope fields follow `IDENTIFIER_STANDARD.md`.
- Identifier, money, time and locale representation follow the P00.03 foundation standards; modules must not invent alternative primitives.

## 6. Capability naming

Capabilities use stable action-oriented identifiers:

```text
crm.customer.read
crm.lead.create
commerce.order.create
commerce.order.cancel
inventory.stock.reserve
payments.payment.authorize
payments.payment.refund
finance.invoice.issue
```

Rules:

- `<domain>.<resource>.<action>`;
- action is explicit and bounded;
- avoid broad powers such as `manage_all`;
- protected mutations require authorization and audit semantics;
- AI tools expose governed capabilities rather than raw CRUD/database access.

## 7. Permission naming

Permissions normally mirror capability scope:

```text
crm.lead.read
crm.lead.create
crm.lead.update
crm.lead.delete
finance.invoice.issue
payments.refund.execute
```

Do not use role names as permissions (`admin`, `manager`). Roles compose permissions.

## 8. Event naming

Canonical event names use past-tense facts plus major version:

```text
commerce.order.created.v1
commerce.order.cancelled.v1
inventory.stock.reserved.v1
payments.payment.authorized.v1
finance.invoice.issued.v1
```

Rules:

- event is a fact that already happened;
- producer owns the event contract;
- version is explicit;
- commands must not masquerade as events (`create_order` is not an event);
- consumers must not depend on undocumented producer internals.

## 9. API path naming

Until P00.04 freezes the full API standard:

- use plural resource nouns where resource APIs are appropriate;
- do not expose database table names merely because they exist;
- domain ownership must be visible in routing/versioning strategy;
- actions that are not CRUD must be explicit business operations;
- no unversioned public contract may be treated as stable.

## 10. Configuration naming

Configuration keys use stable dot-separated paths:

```text
platform.locale.default
security.session.max_age
commerce.checkout.guest_enabled
inventory.reservation.timeout
```

Secrets must never be stored in ordinary configuration values or committed files.

## 11. File and document naming

- Canonical architecture/governance documents: uppercase snake-style Markdown names where already established.
- ADRs: `ADR-NNNN-short-kebab-description.md`.
- Work packages: `Pxx.yy` identifiers are immutable once published.
- Migration naming convention will be fixed with implementation/repository standards before migrations exist.

## 12. UI naming

UI labels may be human-friendly but must map to one canonical concept. A UI label must not silently redefine ownership. Example: a screen may show `Company`, while the underlying canonical concept is `Legal Entity`; that mapping must be explicit in the owning module/product copy.

## 13. Forbidden naming patterns

Do not introduce:

- `GlobalAdmin` as an authorization bypass;
- `SuperUser` with implicit unrestricted access;
- shared `CommonModel`/`BaseBusinessEntity` that accumulates unrelated domain behavior;
- `MasterData` as a dumping-ground domain;
- `CoreCRM`, `NewCRM`, `CRM2` parallel implementations;
- unqualified `Account`, `Client`, `Member`, `Transaction`, `Status` as cross-platform canonical entities;
- technology names in domain contracts (`RedisOrder`, `PostgresCustomer`).

## 14. Naming conflict protocol

If a proposed name conflicts with this standard:

1. reuse the canonical term when semantics match;
2. qualify the term when semantics differ;
3. update glossary/ownership documentation if a genuinely new concept is required;
4. use an ADR when the change alters an established architectural concept.

Implementation must not proceed by silently introducing a synonym.
