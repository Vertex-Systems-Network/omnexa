# Omnexa Module Architecture Dossiers

Status: **Mandatory pre-development planning baseline**

This directory is the canonical architecture-planning layer for roadmap modules that have not yet entered executable implementation. It supplements `docs/roadmap/MASTER_PLAN.md`, `docs/architecture/MODULE_STANDARD.md`, `docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md`, and the machine state in `docs/roadmap/STATE.json`.

These documents are **planning only**. They do not activate future phases or authorize code outside the sole active package in `STATE.json`.

## Required read order for future module work

Before implementing any P02-P27 module or submodule, an AI/human contributor must read:

1. `AGENTS.md` and `docs/roadmap/STATE.json`;
2. the active work-package specification;
3. `docs/governance/AI_EXECUTION_POLICY.md`;
4. `docs/architecture/MODULE_STANDARD.md` and ownership/dependency standards;
5. `docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md`;
6. `docs/roadmap/modules/SUBMODULE_CATALOG.json`;
7. the owning program dossier in this directory;
8. applicable API/event/security/data/accessibility standards.

If an approved submodule decomposition already exists here, implementation agents **must execute it rather than restart generic architecture planning**. A change to ownership, phase placement, public contracts, persistence ownership, or dependency class requires change control/ADR as applicable.

## Dossier set

- `FOUNDATION_P02_P06.md` — identity/tenancy, module runtime, data/event fabric, workflow OS, universal business foundation.
- `CORE_BUSINESS_P07_P15.md` — CRM, finance/ERP, commerce, payments, POS, experience/CMS, portals, HR/projects/service, supply chain/manufacturing.
- `PLATFORM_P16_P18.md` — integration fabric, low-code app builder, reporting/BI.
- `INTELLIGENCE_P19_P20.md` — AI platform and governed agents.
- `ECOSYSTEM_P21_P22.md` — developer platform and marketplace.
- `GLOBAL_ENTERPRISE_P23_P25.md` — globalization, enterprise governance/security/compliance, scale fabric.
- `INDUSTRY_AUTONOMY_P26_P27.md` — industry packs and autonomous business OS.

## Mandatory dossier structure

Every phase/submodule plan must define:

- stable submodule ID and owner;
- architecture boundary and authoritative write model;
- required/optional/platform/forbidden dependencies;
- primary user/system flows;
- configurable options/settings/feature flags and who may change them;
- permissions and tenant/org scope;
- APIs/capabilities/events/workflow/UI contribution points;
- persistence/migrations, storage/cache/search/reporting implications;
- security/data classification/secrets;
- failure, retry, idempotency and degradation behavior;
- lifecycle, import/export and compatibility behavior;
- UI accessibility/localization/RTL requirements when applicable;
- ordered implementation tasks and evidence gates.

## Planning state versus implementation state

A submodule may be `planned` in `SUBMODULE_CATALOG.json` while its phase is inactive. That means its boundary is preplanned, **not executable**. The only implementation authority remains `docs/roadmap/STATE.json`.

## No hidden feature creation

If implementation discovers a missing capability, classify it as one of:

- an already-owned subtask in the active submodule;
- a prerequisite defect in an earlier platform contract;
- a new future submodule/change-control item.

Do not silently create a new independent write model, auth stack, event transport, builder runtime, or cross-module schema while coding.