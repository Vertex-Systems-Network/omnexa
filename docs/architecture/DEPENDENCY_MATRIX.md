# Omnexa Architecture Dependency Matrix

Status: **Canonical dependency baseline v1**  
Work package: **P00.02**

This document defines which dependency directions are allowed. It complements `DOMAIN_OWNERSHIP.md` and `MODULE_STANDARD.md`.

## 1. Dependency symbols

| Symbol | Meaning |
|---|---|
| `K` | May depend on a stable kernel capability. |
| `C` | May consume a versioned public capability/API. |
| `E` | May consume a versioned event. |
| `W` | May participate through a governed workflow. |
| `P` | May consume an approved read projection. |
| `X` | Direct dependency is forbidden. |
| `-` | Same owner/domain; internal dependency rules apply. |

A row consuming a column never grants direct access to the column owner's database tables or private implementation.

## 2. High-level matrix

| Consumer → / Owner ↓ | Kernel | Party | CRM | Finance | Commerce | Inventory | Payments | Workflow | Experience | AI |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Kernel | - | X | X | X | X | X | X | X | X | X |
| Party | K | - | X | X | X | X | X | W | X | X |
| CRM | K | C/E | - | C/E/P | C/E/P | E/P | E/P | W | C | C |
| Finance | K | C/E | X | - | E/C | E/C | C/E | W | X | C |
| Commerce | K | C/E | C/E | C/E | - | C/E | C/E | W | C | C |
| Inventory | K | X | X | E/C | C/E | - | X | W | X | C |
| Payments | K | C | X | C/E | C/E | X | - | W | C | C |
| Workflow | K | C | C | C | C | C | C | - | C | C |
| Experience | K | C/P | C/P | C/P | C/P | C/P | C/P | W | - | C |
| AI | K | C/P | C/P | C/P | C/P | C/P | C/P | W | C/P | - |

The matrix shows permitted integration mechanisms, not mandatory dependencies. The smallest required dependency must be chosen.

## 3. Non-negotiable forbidden paths

The following are always forbidden unless an approved ADR changes the architecture:

- kernel importing business-domain packages;
- a module writing another module's database tables;
- a consumer calling another module's private/internal handlers/classes;
- UI code bypassing capability authorization by writing directly to storage;
- AI agents writing transactional tables directly;
- a workflow engine taking ownership of domain business state;
- analytics/read models becoming authoritative write sources;
- connectors encoding business rules that belong to a domain owner;
- optional module absence causing unrelated module boot failure;
- circular required module dependencies.

## 4. Preferred integration order

Choose the least coupled mechanism that preserves correctness:

1. **Kernel capability** for truly shared platform primitives.
2. **Domain capability/API** when the caller requires synchronous authoritative behavior/result.
3. **Event** when reacting to a fact asynchronously.
4. **Workflow** when coordinating multi-step, long-running or compensating business processes.
5. **Projection** for cross-domain read/query experiences where eventual consistency is acceptable.

Do not use an event merely to perform a synchronous command, and do not use a synchronous API for every side-effect when an event is the correct semantic boundary.

## 5. Representative flows

### Commerce order → inventory/payment/finance

```text
Commerce creates Order
  -> capability: Inventory.reserve
  -> capability: Payments.authorize (when required synchronously)
  -> event: commerce.order.confirmed.v1
       -> Finance reacts / workflow coordinates invoice policy
       -> CRM projection/timeline may react
```

Commerce does not update inventory quantity, finance ledger or payment tables directly.

### Customer creation → CRM

```text
business.party owns customer relationship
  -> event: party.customer.created.v1
       -> CRM creates/updates its owned engagement projection
```

CRM may request party-owned mutations only through party capabilities.

### AI action

```text
AI Agent
  -> policy/authorization check
  -> declared Omnexa capability
  -> owning domain validates business rules
  -> mutation + audit/event
```

No raw SQL or unrestricted internal service access is permitted.

## 6. Dependency declaration

Every module manifest must eventually declare:

- required kernel capabilities;
- required module dependencies;
- optional module dependencies;
- capabilities consumed;
- capabilities provided;
- events published;
- events consumed;
- workflow actions/triggers provided;
- UI extension points used/provided.

P03 will make this machine-enforceable. Until then, this document is the canonical architecture rule.

## 7. Conflict rule

If an implementation appears to require an `X` path, do not create the dependency. Either:

1. expose a proper owner capability/event/projection;
2. move misplaced responsibility to the correct owner; or
3. raise an ADR when the ownership model itself is genuinely wrong.
