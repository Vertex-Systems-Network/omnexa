# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Foundation Architecture v1 is **FROZEN**; P00, P01, P02 and P03 are complete.

> **Current governed checkpoint:** **P03 — Module Runtime: DONE 11 / 11; P03 exit SATISFIED; current work package NONE; P04 remains PLANNED / NOT ACTIVATED.** `kernel_code_authorized=false`; `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` is the canonical machine-readable execution cursor. Completed P03 does not grant P04 implementation authority.

## Project progress

```text
P00  Product Constitution & Architecture Freeze  [██████████] 10/10  DONE
P01  Omnexa Kernel                               [██████████] 12/12  DONE
      └─ Exit: SATISFIED
P02  Identity, Tenancy & Organization            [██████████] 10/10  DONE
      └─ Exit: SATISFIED
P03  Module Runtime                              [██████████] 11/11  DONE
      ├─ P03.01 — Module Manifest Schema: DONE
      ├─ P03.02 — Registry & Deterministic Discovery: DONE
      ├─ P03.03 — Dependency Graph Resolver: DONE
      ├─ P03.04 — Module Lifecycle State Machine: DONE
      ├─ P03.05 — Module Settings & Feature Flags: DONE
      ├─ P03.06 — Capability Registry: DONE
      ├─ P03.07 — Permission Registration: DONE
      ├─ P03.08 — UI Contribution Registry Contract: DONE
      ├─ P03.09 — Migration Ownership Registry: DONE
      ├─ P03.10 — Module Health Reporting: DONE
      └─ P03.11 — Package Trust Hooks & P03 Exit Proof: DONE
      └─ Exit: SATISFIED
P04+ Future phases                               [░░░░░░░░░░]        PLANNED / LOCKED
```

## P03.11 terminal evidence

P03.11 is accepted through implementation issue #132, draft PR #133, promotion PR #134, exact head `a083a8a86ec3a51309fa479ee49c79e1b6ec9f10`, draft Governance `33258092323 / 99115191521`, promotion Governance `33258456851 / 99116152701`, and implementation merge `b3b9b61f963df6a05ea45cbd3c562e12974d92d0`.

The package retains already-validated publisher/provenance/SBOM/data/security declarations in immutable registry-bound snapshots and exposes deterministic typed/versioned `metadata_only` profiles. Metadata presence never means trusted/certified, secret locators/values are not surfaced, and package code is not executed for metadata discovery.

The aggregate P03 exit proof covers required dependency enforcement, optional dependency degradation, safe disable/re-enable, supported upgrade/migration path, forbidden cross-module dependency detection, health/state accuracy and no unrelated module corruption.

## Mandatory contributor / AI start here

Read `AGENTS.md` first, then `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json`, `docs/roadmap/evidence/P03.11_COMPLETION_2026-08-29.md`, `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml` and `docs/ai/handoffs/P03.11.md`.

P04 must not be implemented until a later separate readiness/preparation and explicit activation transition passes canonical governance and merges to protected `main`.

## Core laws

- Kernel before business modules.
- One authoritative owner per write model/capability.
- Cross-module direct DB writes/private implementation imports are forbidden.
- Cross-domain communication uses governed APIs/capabilities/events/workflows/read projections.
- Tenant scope, authorization, audit, observability and contract versioning are mandatory.
- Optional modules fail/degrade independently.
- AI acts only through governed authorized capabilities; no unrestricted raw DB/object-store/business-state authority.
- Strict modular monolith first; service extraction requires evidence and ADR.
- Architecture/roadmap changes require change control and reconciliation.

## Protected GitHub integration and executable CI

Issue #3 is closed and `main` remains protected with PR-only integration, strict required `governance`, failed-check merge rejection, conversation resolution and up-to-date enforcement. Canonical CI uses GitHub-hosted `ubuntu-24.04` only. Local/self-hosted governance runners are prohibited.

## External distribution gate

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains the external distribution/public-launch licensing, IP and trademark decision gate. It grants no implementation authority.
