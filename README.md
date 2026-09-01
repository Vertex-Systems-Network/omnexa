# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Foundation Architecture v1 is **FROZEN**; P00, P01, P02 and P03 are complete.

> **Current governed checkpoint:** **P04 — Data, Jobs & Event Fabric is ACTIVE at 3 / 10 completed packages; P04.04 Transactional Outbox Reliability Primitive is the sole ACTIVE work package.** `kernel_code_authorized=true` for P04.04 only; `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` is the canonical machine-readable execution cursor. This README is a human-readable status mirror only and never grants implementation authority by itself.

## Project progress dashboard

Phase-count view: **4 of 28 phases completed**, with **P04 active**. Canonical P04 package completion is **3 / 10 (30%)**; P04.04 is the sole active package and its runtime implementation is not yet complete.

| Phase | Program / Module Family | Status | Work-package progress | Progress | Start date | End date |
|---|---|---:|---:|---|---|---|
| P00 | Product Constitution & Architecture Freeze | ✅ DONE | 10/10 (100%) | `██████████` | 2026-08-21 | 2026-08-22 |
| P01 | Omnexa Kernel | ✅ DONE | 12/12 (100%) | `██████████` | 2026-08-22 | 2026-08-23 |
| P02 | Identity, Tenancy & Organization | ✅ DONE | 10/10 (100%) | `██████████` | 2026-08-23 | 2026-08-26 |
| P03 | Module Runtime | ✅ DONE | 11/11 (100%) | `██████████` | 2026-08-26 | 2026-08-29 |
| P04 | Data, Jobs & Event Fabric | 🟡 ACTIVE | 3/10 canonical done (30%) | `███░░░░░░░` | 2026-08-30 | TBD (active) |
| P05 | Omnexa Flow / Workflow OS | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P06 | Universal Business Foundation | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P07 | CRM, Sales & Customer 360 | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P08 | Finance & ERP Core | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P09 | Commerce OS | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P10 | Payment Fabric | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P11 | POS & Edge | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P12 | Experience Builder & CMS | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P13 | Portal Platform | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P14 | HR, Projects & Service Operations | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P15 | Supply Chain, Warehouse & Manufacturing | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P16 | Omnexa Connect / Integration Fabric | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P17 | Low-code App Builder | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P18 | Data, Reporting & BI | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P19 | Omnexa Intelligence Platform | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P20 | Governed AI Agents | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P21 | Developer Platform | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P22 | Omnexa Exchange / Marketplace | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P23 | Globalization & Country Packs | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P24 | Enterprise Governance, Security & Compliance | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P25 | Scale Fabric | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P26 | Industry Packs | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P27 | Autonomous Business OS | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |

**Date semantics:** completed/current dates are governance-history dates. Future dates remain `TBD` until governed planning or activation establishes them; the roadmap does not invent delivery deadlines.

## Current P04 package / module status

| Package | Module / Contract | Status | Progress | Progress bar | Start date | End date |
|---|---|---:|---:|---|---|---|
| P04.01 | Event envelope & identity contract | ✅ DONE | 100% | `██████████` | 2026-08-30 | 2026-08-31 |
| P04.02 | Publish/subscribe abstraction & ownership boundaries | ✅ DONE | 100% | `██████████` | 2026-08-31 | 2026-08-31 |
| P04.03 | Durable stream/consumer baseline & checkpoint model | ✅ DONE | 100% | `██████████` | 2026-08-31 | 2026-08-31 |
| P04.04 | Transactional outbox reliability primitive | 🟡 ACTIVE | 0% runtime implementation | `░░░░░░░░░░` | 2026-09-01 | TBD |
| P04.05 | Consumer inbox/deduplication & idempotency primitive | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.06 | Retry/backoff, terminal failure & dead-letter/quarantine policy | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.07 | Event schema registry, compatibility & validation | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.08 | Background-job ownership, tenant context & correlation | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.09 | Reliability observability, diagnostics & operator recovery | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.10 | Replay/duplicate/failure/restart/poison-event acceptance | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |

### P04 status interpretation

- **P04.01-P04.03** are canonically complete.
- **P04.04** is the sole active package, but runtime implementation remains at 0% on this transition carrier.
- **P04.05-P04.10** remain planned/locked.
- P04.04 runtime work must begin on a fresh separate branch created only after this closure/activation state is accepted and read back from protected `main`.

## Accepted P04 evidence chain
- **P04.01 — Event Envelope & Identity Contract:** accepted completion evidence is retained in `docs/roadmap/evidence/P04.01_COMPLETION_2026-08-31.md`.
- **P04.02 — Publish/Subscribe Abstraction & Ownership Boundaries:** accepted completion evidence is retained in `docs/roadmap/evidence/P04.02_COMPLETION_2026-08-31.md`.
- **P04.03 — Durable Stream/Consumer Baseline & Checkpoint Model:** DONE from source PR #165 / promotion PR #166, exact implementation head `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`, promotion Governance run/job `33405463251 / 99531835998`, implementation merge `b94189873bef11f4870935205398f1ef44f160bf`, and completion evidence `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`.
- **P04.04 — Transactional Outbox Reliability Primitive:** contract/handoff preparation accepted through source PR #168 and promotion PR #169; promotion Governance run/job `33410873382 / 99549818529`; preparation merge/read-back `962a62c7c111079ca6f2047fa748deea97c84534`. The closure/activation transaction makes P04.04 active; runtime remains unimplemented on this carrier.

## P04.04 bounded implementation scope

A fresh P04.04 implementation branch created from the accepted post-transition protected-main SHA may implement only the accepted provider-neutral transactional outbox boundary: local PostgreSQL atomic mutation + canonical event persistence, recoverable pending state, relay through P04.02, duplicate-safe crash/restart behavior, concurrency safety, and tenant/owner isolation.

This closure/activation carrier adds **no migration** and grants no blanket schema authority. If durable outbox persistence requires schema mutation, the implementation branch must first record the exact `kernel.events` migration path/version/data budget and prove fresh-install, upgrade, rollback/forward-recovery, and tenant/owner isolation under the retained P01 migration rules.

Still unauthorized: concrete broker selection, end-to-end exactly-once claims, P04.05 inbox/deduplication, P04.06 retry/DLQ, P04.07 schema registry, P04.08 background-job changes, business features, strategic X-program runtime, and AI/model/agent runtime.

## Mandatory contributor / AI start here

Read `AGENTS.md` first, then `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P04_ENTRY_GATE.md`, `docs/governance/P04_03_P04_04_TRANSITION_TRANSACTION.md`, `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`, retained P04.01-P04.03 completion evidence, `docs/roadmap/work-packages/P04.04.md`, `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml` and `docs/ai/handoffs/P04.04.md`.

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
