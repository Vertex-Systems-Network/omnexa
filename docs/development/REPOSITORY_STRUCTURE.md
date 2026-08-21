# Omnexa Repository Structure Standard

Status: **Canonical v1**  
Work package: **P00.08**

Omnexa uses one governed monorepo as the source of truth while beginning as a strict modular monolith. Directory names express architectural ownership, not deployment count. A folder does not imply a microservice.

## Canonical top-level layout

```text
omnexa/
├── apps/
│   ├── admin-web/
│   ├── portal-web/
│   ├── pos/
│   └── mobile/
├── kernel/
│   ├── identity/
│   ├── tenancy/
│   ├── authz/
│   ├── modules/
│   ├── events/
│   ├── workflow/
│   ├── files/
│   └── observability/
├── modules/
│   ├── crm/
│   ├── sales/
│   ├── finance/
│   ├── inventory/
│   ├── commerce/
│   └── ...
├── platform/
│   ├── api/
│   ├── sdk/
│   ├── cli/
│   └── module-sdk/
├── shared/
│   └── contracts/
├── infrastructure/
│   ├── local/
│   ├── docker/
│   ├── kubernetes/
│   └── terraform/
├── scripts/
├── docs/
└── generated/
```

Folders are created only when implementation reaches the phase that owns them. P00 documents the target layout; it does not authorize empty scaffolding or product code.

## Ownership boundaries

- `kernel/` contains shared platform capabilities only.
- `modules/<module>/` contains one domain/module's private implementation and owned schema/migrations.
- `apps/` contains user-facing composition shells; apps do not become alternate business-authority layers.
- `platform/` contains public developer surfaces such as API gateway composition, SDKs, CLI and module SDK.
- `shared/contracts/` contains intentionally shared public contracts/value primitives, never a dumping ground for domain logic.
- `infrastructure/` contains deployment/runtime orchestration, separated from product domain code.
- `generated/` contains reproducible derived artifacts only and is never a hand-maintained source of truth.

## Module internal layout

A future module should converge on a predictable shape such as:

```text
modules/<module>/
├── module.yaml
├── internal/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   └── transport/
├── contracts/
│   ├── api/
│   ├── events/
│   └── schemas/
├── migrations/
├── tests/
└── README.md
```

Exact language-specific folders may vary through implementation ADRs, but ownership rules may not.

## Import/dependency rules

- Another module may not import `modules/<other>/internal/**`.
- Cross-module direct DB writes are forbidden; another module may not read/write another module's private tables/migrations.
- Cross-domain code uses governed public contracts/capabilities/events/workflows/read projections.
- `shared/` may contain universal primitives only when ownership is truly platform-wide.
- Apps may compose modules but may not bypass their public capabilities for protected business mutations.
- Generated code may depend on canonical contracts; canonical contracts never depend on generated output.

## Schema and migration placement

- Kernel schema/migrations live with the owning kernel capability.
- Domain schema/migrations live with the owning module.
- Migration ordering across owners is coordinated by the platform migration engine/manifest contract, not arbitrary filename coupling.
- No central `misc` migration directory may become an ownership bypass.

## Contract locations

Repository-wide canonical contract templates live under `docs/contracts/` during P00. Executable source contracts later live with their owner and may generate/validate documentation artifacts.

Public/shared contracts must expose ownership, version and compatibility semantics.

## Generated artifact policy

Derived outputs may include:

- generated SDKs;
- OpenAPI renderings;
- event/schema bundles;
- generated client types;
- documentation indexes;
- release metadata.

They must be reproducible from canonical inputs. CI/local verification must detect drift when generated artifacts are committed.

## No architecture-by-folder

Creating `services/`, `microservices/`, separate databases or Kubernetes deployments does not authorize service extraction. Deployment boundaries require evidence and an ADR under the modular-monolith-first law.

## Repository hygiene

Prohibited:

- unrelated project files;
- copied vendor/build output unless explicitly required;
- secrets/credentials;
- local databases/logs/temp files;
- duplicate generated artifacts without an owner;
- `common`, `utils`, `helpers` dumping grounds that hide domain ownership;
- abandoned migration/test fixtures with no supported path.

## Future extraction

If a domain is extracted into a service, the logical owner and public contracts remain stable. The repository may retain one monorepo even when deployment topology evolves; physical extraction does not redefine business ownership.
