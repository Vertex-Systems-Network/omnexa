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
- **P04.04** is the sole active package, but runtime implementation remains at 0%.
- **P04.05-P04.10** remain planned/locked.
- P04.04 runtime work must begin on fresh isolated task branches from the latest protected-main SHA after re-reading canonical state and accepted continuity surfaces.

## Accepted P04 evidence chain

- **P04.01 — Event Envelope & Identity Contract:** accepted completion evidence is retained in `docs/roadmap/evidence/P04.01_COMPLETION_2026-08-31.md`.
- **P04.02 — Publish/Subscribe Abstraction & Ownership Boundaries:** accepted completion evidence is retained in `docs/roadmap/evidence/P04.02_COMPLETION_2026-08-31.md`.
- **P04.03 — Durable Stream/Consumer Baseline & Checkpoint Model:** DONE from source PR #165 / promotion PR #166, exact implementation head `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`, promotion Governance run/job `33405463251 / 99531835998`, implementation merge `b94189873bef11f4870935205398f1ef44f160bf`, and completion evidence `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`.
- **P04.04 — Transactional Outbox Reliability Primitive:** contract/handoff preparation accepted through source PR #168 and promotion PR #169; preparation promotion Governance run/job `33410873382 / 99549818529`; preparation merge/read-back `962a62c7c111079ca6f2047fa748deea97c84534`. P04.03 closure / P04.04 activation then passed source PR #171 / Governance `33546204586 / 99984199589`, unchanged promotion PR #172 / Governance `33547080086 / 99987132610`, and merged/read back as `50edeae03ad52e435d142b2f22c803f08a5c7f1a`.
- **Post-activation continuity:** source PR #173 final head `3a9ccc6ce165022279c52bfd4f7bedf2d739950d` passed Governance `33549921833 / 99996561115`; unchanged promotion PR #174 passed Governance `33550803632 / 99999473766` and merged as protected-main checkpoint `6c7937b3d94177604b03c4872a48504ac60034ce`. Earlier source run `33549594705 / 99995490009` remains diagnostic FAIL evidence and is not acceptance authority.
- **Multi-agent development foundation:** source PR #175 final head `c74e9d116880a9fde69d9220c1907a67e2cf86eb` passed Governance `33553700739 / 100009408168`; unchanged promotion PR #176 passed promotion Governance run `33554787584` and merged/read back as protected-main checkpoint `8fd737b318947c9d0f3cc0e5c5a0931636c2c40c`.

## P04.04 bounded implementation scope

P04.04 may implement only the accepted provider-neutral transactional outbox boundary: local PostgreSQL atomic mutation + canonical event persistence, recoverable pending state, relay through P04.02, duplicate-safe crash/restart behavior, concurrency safety, and tenant/owner isolation.

The activation/continuity carriers add **no migration** and grant no blanket schema authority. If durable outbox persistence requires schema mutation, the implementation plan must record the exact `kernel.events` migration path/version/data budget and prove fresh-install, upgrade, rollback/forward-recovery, and tenant/owner isolation under the retained P01 migration rules.

Still unauthorized: concrete broker selection, end-to-end exactly-once claims, P04.05 inbox/deduplication, P04.06 retry/DLQ, P04.07 schema registry, P04.08 background-job changes, business features, strategic X-program product runtime, and AI/model/agent product runtime.

## Multi-agent development operating model

Omnexa development may use parallel AI/human agents to accelerate delivery, but parallelism stays subordinate to the single canonical phase/work-package cursor.

**Current foundation target:** **4-6 active agents with no more than 3 concurrent write agents** until machine-enforced leases/scope validation are proven. Other agents should be read-only planners, reviewers, security analyzers or CI/evidence coordinators.

Safe current pattern:

```text
Supervisor / Integration     review + governed merge + own isolated verifier task
Runtime Source Agent          bounded source lease
Persistence/Migration Agent   exact reserved migration/data lease
Test Agent                    non-overlapping test lease
Security/Architecture Agent   read-only
CI/Evidence Agent             verifier/docs/evidence only
```

Business-domain agents such as CRM, Finance, Inventory, HR and Commerce may become broadly parallel only after canonical roadmap gates open those independent streams.

Mandatory concurrency rules are in `docs/governance/MULTI_AGENT_ORCHESTRATION.md`. Supervisor interruption/merge/sync behavior is in `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`. The reusable acceleration plan is `docs/roadmap/XQ_100_MULTI_AGENT_DEVELOPMENT_PLAN.md`. Machine-readable task/lease/signal shapes are under `docs/ai/`.

### Active P04.04 multi-agent wave

Wave: `P04.04-WAVE-20260902-01`  
Coordination channel: GitHub issue `#177`  
Branch bootstrap base: `8fd737b318947c9d0f3cc0e5c5a0931636c2c40c`

| Merge order | Agent / Task | Module / responsibility | Branch | Current dependency |
|---:|---|---|---|---|
| 1 | Agent-01 / P04.04-T01 | Outbox core | `agent/20260902-p04-04-outbox-core` | none |
| 2 | Agent-02 / P04.04-T02 | PostgreSQL persistence + migration | `agent/20260902-p04-04-persistence-migration` | T01 contract before final adapter submission |
| 3 | Agent-03 / P04.04-T03 | crash/restart/concurrency/tenant-isolation tests | `agent/20260902-p04-04-reliability-tests` | T01 + T02 accepted contracts |
| 4 | Supervisor / P04.04-T04 | integration verifier + Governance hook | `supervisor/20260902-p04-04-integration` | T01 + T02 + T03 |

All four branches were created before task work and synchronized to the accepted bootstrap SHA. The machine-readable assignments, write budgets, dependencies and merge order are in `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`.

Provisional migration reservation for T02, pending acceptance of the active-wave governance plan:

`kernel/migrations/kernel.events/1_create_transactional_outbox.sql`

Owner/version: `kernel.events / 1`. It may contain only new P04.04 owner-scoped outbox persistence; no business-table mutation/backfill/destructive data operation or P04.05/P04.06 semantics.

### Supervisor submission / interrupt protocol

A completed task announces exactly:

`Work Done and Submitted`

with task, branch, exact head SHA, PR, CI/evidence and instruction-check metadata.

When a valid submission arrives, the Supervisor checkpoints/pauses its own task, reviews the submitted exact head, and either requests changes or routes the approved head through the repository's required source/Governance/unchanged-promotion/protected-main pipeline. Supervisor approval never bypasses protected integration.

After protected-main merge/readback the Supervisor announces exactly:

`New changes have been merged — please merge these changes into your branch first, then resume your own work.`

Every active worker must synchronize that new protected-main SHA, re-check instructions/leases/dependencies, invalidate stale CI when its head changes, and only then resume. Recommended acknowledgement:

`Sync Complete — Resuming Work`

The current persistent signal channel is issue `#177`. Missing the message is not an excuse to continue stale work: every material mutation must compare the branch against the active plan's `required_main_sha`.

This developer orchestration is not P20 product-agent runtime and does not activate a future business phase.

## Agent Working Instructions

Every material agent task must check these instructions **at task start and again before PR submission**:

1. re-read protected `main` and canonical `docs/roadmap/STATE.json`;
2. confirm active phase/work package and owning module/domain/kernel capability;
3. read `AGENTS.md`, applicable work-package spec/handoff, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/MULTI_AGENT_ORCHESTRATION.md`, and when a Supervisor wave is active, `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md` plus `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`;
4. for a **new** parallel wave, after mandatory canonical readback the Supervisor's first repository mutation is creation of all worker branches and its own isolated task branch from the recorded base SHA; planning/coding starts only after branch verification;
5. record task/agent identity, branch and exact base/last-synced/required-main SHAs;
6. declare read paths, write paths, forbidden paths and shared/exclusive paths;
7. check concurrent tasks/leases for path, migration, registry and public-contract overlap;
8. resolve dependency ordering; do not code against guessed future contracts;
9. for schema changes, reserve exact owner/path/version/data budget before mutation;
10. check issue `#177` or the active wave's recorded coordination channel and compare against `required_main_sha` before each material mutation;
11. when another accepted task moves protected main, **sync the new main into the branch before resuming**, then re-run stale tests/CI as required;
12. when task work is ready for Supervisor review, announce `Work Done and Submitted` with the required exact-head metadata;
13. keep cross-module private writes/imports forbidden;
14. require exact-final-head tests/CI/review and re-check protected-main freshness before merge;
15. compare the effective working instructions with this README snapshot;
16. **if working instructions changed materially, update this section in the same PR**;
17. if they did not change, record in the PR: `Agent instructions checked — README instruction delta: none`.

Material instruction changes include changes to phase/package, role, owner, branch/base/sync strategy, allowed/forbidden/shared paths, migration budget, dependency/contract assumptions, required tests/gates, Supervisor/merge process, coordination channel, tool/network/secret restrictions or stop conditions.

This README is only the human-readable mirror. `AGENTS.md`, `STATE.json`, mandatory governance policy and accepted ADRs remain higher authority.

## Mandatory contributor / AI start here

Read `AGENTS.md` first, then `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P03_EXIT_GATE.md`, `docs/governance/P04_ENTRY_GATE.md`, `docs/governance/P04_03_P04_04_TRANSITION_TRANSACTION.md`, `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`, retained P04.01-P04.03 completion evidence, `docs/roadmap/work-packages/P04.04.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/MULTI_AGENT_ORCHESTRATION.md`, `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`, `docs/roadmap/XQ_100_MULTI_AGENT_DEVELOPMENT_PLAN.md`, `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`, `docs/ai/AGENT_SIGNAL_SCHEMA.json`, `docs/ai/AI_CONTEXT.md`, `docs/ai/AI_STATE.yaml` and `docs/ai/handoffs/P04.04.md`.

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