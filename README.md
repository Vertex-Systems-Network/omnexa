# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Foundation Architecture v1 is **FROZEN**; P00, P01, P02 and P03 are complete.

> **Current canonical state:** **P04 — Data, Jobs & Event Fabric is ACTIVE at 5 / 10 completed packages. P04.01-P04.05 are DONE and P04.06 Retry/Backoff, Terminal Failure & Dead-Letter/Quarantine Policy is the sole ACTIVE package.** Accepted activation is protected `main@4c9f60843f2612bc4c9a10b4efca7b6a20826be3`. P04.06 runtime work remains blocked until the separate post-activation continuity carrier is governed, promoted unchanged, merged and read back. `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` is the canonical machine-readable execution cursor. This README is a human-readable status mirror only and never grants implementation authority by itself.

## Project progress dashboard

Phase-count view: **4 of 28 phases completed**, with **P04 active**. Current P04 package completion is **5 / 10 (50%)**; P04.06 is active under the required post-activation continuity gate.

| Phase | Program / Module Family | Status | Work-package progress | Progress | Start date | End date |
|---|---|---:|---:|---|---|---|
| P00 | Product Constitution & Architecture Freeze | ✅ DONE | 10/10 (100%) | `██████████` | 2026-08-21 | 2026-08-22 |
| P01 | Omnexa Kernel | ✅ DONE | 12/12 (100%) | `██████████` | 2026-08-22 | 2026-08-23 |
| P02 | Identity, Tenancy & Organization | ✅ DONE | 10/10 (100%) | `██████████` | 2026-08-23 | 2026-08-26 |
| P03 | Module Runtime | ✅ DONE | 11/11 (100%) | `██████████` | 2026-08-26 | 2026-08-29 |
| P04 | Data, Jobs & Event Fabric | 🟡 ACTIVE | 5/10 done (50%) | `█████░░░░░` | 2026-08-30 | TBD (active) |
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
| P04.04 | Transactional outbox reliability primitive | ✅ DONE | 100% accepted implementation | `██████████` | 2026-09-01 | 2026-09-04 |
| P04.05 | Consumer inbox/deduplication & idempotency primitive | ✅ DONE | 100% accepted implementation | `██████████` | 2026-09-04 | 2026-09-05 |
| P04.06 | Retry/backoff, terminal failure & dead-letter/quarantine policy | 🟡 ACTIVE / CONTINUITY GATED | 0% runtime implementation | `░░░░░░░░░░` | 2026-09-06 | TBD |
| P04.07 | Event schema registry, compatibility & validation | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.08 | Background-job ownership, tenant context & correlation | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.09 | Reliability observability, diagnostics & operator recovery | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |
| P04.10 | Replay/duplicate/failure/restart/poison-event acceptance | 🔒 PLANNED / LOCKED | Not started | `░░░░░░░░░░` | TBD | TBD |

### P04 status interpretation

- **P04.01-P04.05** have accepted completion evidence.
- **P04.06** is the sole active package on protected main.
- **P04.07-P04.10** remain planned/locked.
- No P04.06 runtime branch/task/lease/migration reservation exists in the post-activation continuity carrier.
- Runtime implementation may begin only from a fresh branch after the continuity source + unchanged promotion are accepted/read back.

## Accepted P04 / development-governance evidence chain

- **P04.01 — Event Envelope & Identity Contract:** `docs/roadmap/evidence/P04.01_COMPLETION_2026-08-31.md`.
- **P04.02 — Publish/Subscribe Abstraction & Ownership Boundaries:** `docs/roadmap/evidence/P04.02_COMPLETION_2026-08-31.md`.
- **P04.03 — Durable Stream/Consumer Baseline & Checkpoint Model:** source PR #165 / promotion PR #166, exact implementation head `ea13d171290fc580cfa8b8ff59cd3ea0f8e26cfe`, promotion Governance `33405463251 / 99531835998`, implementation merge `b94189873bef11f4870935205398f1ef44f160bf`, evidence `docs/roadmap/evidence/P04.03_COMPLETION_2026-08-31.md`.
- **P04.04 preparation/activation:** preparation source PR #168 / promotion #169, promotion Governance `33410873382 / 99549818529`, preparation merge `962a62c7c111079ca6f2047fa748deea97c84534`; activation source #171 Governance `33546204586 / 99984199589`, promotion #172 Governance `33547080086 / 99549818529`, activation merge/readback `50edeae03ad52e435d142b2f22c803f08a5c7f1a`.
- **P04.04 implementation/completion:** final Supervisor head `ef09b878577d25a4a1186cb8fe84205b08a24851`, promotion #193 Governance `33810095507 / 100829646792`, implementation merge `66c072b5caf42ceecb88d30cd1a1ee4e910322e6`, evidence PR #194 merge `4445c21f1e6b03e84859d31ce7b32169b9c4cccc`.
- **P04.05 preparation/activation:** preparation source #195 / unchanged promotion #196 at exact head `211fea2077d7a1bf94be48f32f047b27273a4515`, preparation merge/read-back `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`; activation source #197 / promotion #198 at exact head `6907253d375125a7ff096fb434c3433dbc17b331`, activation merge/read-back `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a`.
- **P04.05 implementation/completion:** final Supervisor head `dd713fe3217a0d092ab3ff31115ac031ae8c0303`, source #210 Governance `33927932705`, promotion #211 Governance `33985111334 / 101357077040`, implementation merge/read-back `0c66a3371dbf2fa942a95b7d0475b06235392474`, completion evidence PR #213 merge/read-back `e44ece77ddf7b821c03997266ca0c68c07162910`.
- **P04.06 preparation:** source #215 / unchanged promotion #216 at exact head `7babb9c39185636b3af5184d5a7bd31cedbc37a0`, source Governance `33987003924`, promotion Governance `33987472967`, preparation merge/read-back `3f547180eb5e839439834eb2ce7977324803df18`.
- **P04.05 closure / P04.06 activation:** source #218 / unchanged promotion #219 at exact head `b718ad7316dba6fca0cafccb514df6da653abe13`; source Governance #677 / `33990309561`; promotion Governance #678; protected merge/read-back `4c9f60843f2612bc4c9a10b4efca7b6a20826be3`.
- **P04.05 post-activation continuity:** source #173 final head `3a9ccc6ce165022279c52bfd4f7bedf2d739950d` Governance `33549921833 / 99996561115`; promotion #174 Governance `33550803632 / 99999473766`; merge `6c7937b3d94177604b03c4872a48504ac60034ce`. Earlier source run `33549594705 / 99995490009` remains diagnostic FAIL evidence.
- **Multi-agent foundation:** source #175 head `c74e9d116880a9fde69d9220c1907a67e2cf86eb` Governance `33553700739 / 100009408168`; promotion #176 Governance run `33554787584`; merge/readback `8fd737b318947c9d0f3cc0e5c5a0931636c2c40c`.
- **Supervisor-led multi-agent workflow:** source #178 head `5bb2e2d83bb5b6d64117a8784a9b38ef326a6a84` Governance `33560130753 / 100030373290`; promotion #179 Governance `33560753173 / 100032373974`; merge/readback `6556023c3f07fffa6a81776dd25188a685d2033c`.
- **New-agent slot onboarding:** source #180 head `16cf3c0325ce1fdd43cb2d9afcf5124806110eb7` Governance `33562800663 / 100039010267`; promotion #181 Governance `33563512576 / 100041307298`; merge/readback `34b9989825e85573aa6d38e782f841132e25f041`.

## P04.06 bounded implementation scope after continuity acceptance

P04.06 may later implement only the accepted provider-neutral retry/backoff/failure-disposition/quarantine boundary: stable structured failure classification, finite deterministic attempts/backoff, authoritative UTC eligibility, one-at-a-time claim/CAS/lease semantics, stale-retry suppression after P04.05 completion, terminal quarantine-before-checkpoint progress, crash-gap recovery and bounded classification-safe evidence.

This continuity carrier adds no migration or runtime branch and grants no schema authority. A later fresh implementation wave must record exact `kernel.events` migration/path/data budget before any schema mutation.

Still unauthorized: concrete broker selection, provider-native DLQ, external-side-effect/end-to-end exactly-once claims, P04.07 schema registry, P04.08 background-job changes, business features, strategic X-program product runtime, and AI/model/agent product runtime.

## Multi-agent development operating model

Omnexa development may use parallel AI/human agents to accelerate delivery, but parallelism stays subordinate to the single canonical phase/work-package cursor.

**Current safe envelope:** **4-6 active agents with no more than 3 concurrent write agents.** The XQ-100 M2 carrier wires machine enforcement into the already-required `governance` job, but the writer cap remains 3 until at least one real registered worker PR proves worker-specific scope and live-main freshness enforcement in canonical CI.

Safe pattern:

```text
Supervisor / Integration     review + governed merge + own isolated verifier task
Runtime Source Agent          bounded source lease
Persistence/Migration Agent   exact reserved migration/data lease
Test Agent                    non-overlapping test lease
Security/Architecture Agent   read-only
CI/Evidence Agent             verifier/docs/evidence only
```

Business-domain agents such as CRM, Finance, Inventory, HR and Commerce may become broadly parallel only after canonical roadmap gates open those independent streams.

Mandatory concurrency rules: `docs/governance/MULTI_AGENT_ORCHESTRATION.md`. Supervisor interruption/merge/sync/onboarding: `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`. Final readiness audit: `docs/governance/MULTI_AGENT_READINESS_AUDIT_2026-09-02.md`. Active machine plan: `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`.

### Completed P04.04 multi-agent wave

Historical wave: `P04.04-WAVE-20260902-01` — **COMPLETED; all leases released**  
Coordination channel: GitHub issue `#177`  
Historical branch-bootstrap base: `8fd737b318947c9d0f3cc0e5c5a0931636c2c40c`  
Last recorded protected-main sync receipt before this M2 carrier: `34b9989825e85573aa6d38e782f841132e25f041`  
**Live required main:** resolve protected `main` before every material mutation, new-wave creation, onboarding assignment, submission review and merge.

| Merge order | Agent / Task | Module / responsibility | Branch | Current dependency |
|---:|---|---|---|---|
| 1 | Agent-01 / P04.04-T01 | Outbox core | `agent/20260902-p04-04-outbox-core` | none |
| 2 | Agent-02 / P04.04-T02 | PostgreSQL persistence + migration | `agent/20260902-p04-04-persistence-migration` | T01 contract before final adapter submission |
| 3 | Agent-03 / P04.04-T03 | crash/restart/concurrency/tenant-isolation tests | `agent/20260902-p04-04-reliability-tests` | T01 + T02 accepted contracts |
| 4 | Supervisor / P04.04-T04 | integration verifier + Governance hook | `supervisor/20260902-p04-04-integration` | T01 + T02 + T03 |

The historical branches and their accepted merges are recorded in `docs/roadmap/evidence/P04.04_COMPLETION_2026-09-04.md`. They grant no current write authority.

Current worker-slot state: **0 active, 0 open**. A fresh P04.06 implementation wave may be created only after this post-activation continuity source + unchanged promotion are accepted and protected main is read back.

Historical accepted migrations include:

- `kernel/migrations/kernel.events/1_create_transactional_outbox.sql` — P04.04 outbox ownership;
- accepted `kernel.events` migration version 2 — P04.05 inbox ownership.

Both are immutable history, not live P04.06 reservations. No version 3 path is reserved by activation or continuity.

### M2 required CI enforcement

The M2 gate is executed inside the existing required `governance` job through `scripts/verify_go_quality.sh`, before expensive Go tooling.

It validates:

- active plan authority vs canonical `STATE.json`;
- worker-slot/task identity and current capacity;
- active write vs forbidden paths;
- pairwise active write-path overlap;
- migration owner/version/path reservation collisions;
- dependency DAG and deterministic merge order;
- registered worker PR changed paths vs declared budget;
- registered worker PR base/ancestry vs live protected `main`;
- helper behavior through focused dependency-free unit tests.

Fail-closed behavior includes unknown `agent/*` PR branches, scope violations, stale worker PR bases, active authority mismatches, path overlaps, migration collisions and dependency-order/cycle violations.

A governance/control branch is not falsely treated as a worker task; it still runs plan/slot/lease/DAG validation and must pass the normal protected Governance/promotion flow.

### Supervisor submission / interrupt protocol

A completed task announces exactly:

`Work Done and Submitted`

with task, branch, exact head SHA, PR, CI/evidence and instruction-check metadata.

When a valid submission arrives, the Supervisor checkpoints/pauses its own task, reviews the submitted exact head, and either requests changes or routes the approved head through required source/Governance/unchanged-promotion/protected-main integration. Supervisor approval never bypasses protected integration.

After protected-main merge/readback the Supervisor announces exactly:

`New changes have been merged — please merge these changes into your branch first, then resume your own work.`

Every active worker resolves and synchronizes the new protected `main`, re-checks instructions/leases/dependencies, invalidates stale CI when its head changes, and only then resumes. Recommended acknowledgement:

`Sync Complete — Resuming Work`

Issue #177 is the persistent multi-agent signal channel. Missing the message is not an excuse to continue stale work.

### New Agent Onboarding

A newly arriving development agent always starts from protected `main`. Arrival alone grants no branch, module, task or lease.

The Supervisor immediately checks the AI-Native worker-slot ledger:

- if an authorized slot is `open`, resolve live main, synchronize the slot branch, assign the agent, mark `occupied`, record identity/SHA/start status, then allow work;
- if multiple slots are open, dependency/merge order decides priority;
- never invent a new module, exceed the concurrency cap or activate a locked phase because an agent arrived;
- if no slot is open, stop onboarding and say exactly: `Go Home Come Back Next Time`.

No P04.06 runtime slot is currently authorized on the continuity carrier, so an arriving runtime worker receives `Go Home Come Back Next Time` until continuity acceptance/read-back and a later governed plan explicitly opens a valid slot.

### Live protected-main freshness rule

`required_main_ref = main` is live freshness authority. Stored SHA fields are audit snapshots only; this avoids recursive SHA-only coordination commits.

This developer orchestration is not P20 product-agent runtime and does not activate a future business phase.

## Agent Working Instructions

Every material agent task checks these instructions **at task start and again before PR submission**:

1. re-read protected `main` and canonical `docs/roadmap/STATE.json`;
2. confirm active phase/work package and owning module/domain/kernel capability;
3. read `AGENTS.md`, applicable work-package spec/handoff, AI execution policy, multi-agent orchestration contract, Supervisor workflow and active plan;
4. newly arriving agents start from protected `main` and require an authorized `open` slot before switching to a task branch;
5. if no slot is open, Supervisor responds exactly `Go Home Come Back Next Time` and grants no work authority;
6. record task/agent identity, slot, branch and current live-main/base/sync evidence;
7. declare read/write/forbidden/shared paths and check active overlap;
8. resolve dependencies; do not code against guessed future contracts;
9. reserve exact owner/path/version/data budget before schema mutation;
10. resolve current protected `main` and the active coordination issue before each material mutation;
11. after another accepted merge, sync new main before resuming and rerun stale tests/CI;
12. registered worker branches must remain registered in the active plan; renaming/creating an `agent/*` branch does not bypass M2 — unknown `agent/*` PRs fail required Governance;
13. do not write outside the declared task budget; registered worker PR scope is machine-enforced by M2;
14. registered worker PRs must be based on current protected `main`; stale worker PRs fail M2 and must synchronize/resubmit;
15. when ready for review, announce `Work Done and Submitted` with exact-head metadata;
16. keep cross-module private writes/imports forbidden;
17. require exact-final-head tests/CI/review and re-check protected-main freshness before merge;
18. compare effective working instructions with this README;
19. **if instructions changed materially, update this section in the same PR**;
20. if unchanged, record: `Agent instructions checked — README instruction delta: none`.

Material instruction changes include phase/package, role, owner, slot availability/assignment, branch/base/sync strategy, path leases, migration budget, dependency/contract assumptions, required M2/tests/gates, Supervisor/merge/onboarding process, coordination channel, tool/network/secret restrictions or stop conditions.

README is only the human-readable mirror. `AGENTS.md`, `STATE.json`, mandatory governance policy and accepted ADRs remain higher authority.

## Mandatory contributor / AI start here

Read `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, P04 entry/transition/package surfaces, `docs/roadmap/work-packages/P04.06.md`, `docs/governance/P04_05_P04_06_TRANSITION_TRANSACTION.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/MULTI_AGENT_ORCHESTRATION.md`, `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`, `docs/governance/MULTI_AGENT_READINESS_AUDIT_2026-09-02.md`, `docs/roadmap/XQ_100_MULTI_AGENT_DEVELOPMENT_PLAN.md`, `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`, task/lease/signal/slot schemas, AI context/state and `docs/ai/handoffs/P04.06.md` before material work.

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

Repository ruleset `21174858` for `main` is active. Protected integration requires PRs, strict required `governance`, conversation resolution, rejects non-fast-forward/deletion and has no bypass actor. Canonical CI uses GitHub-hosted `ubuntu-24.04` only; local/self-hosted governance runners are prohibited.

## External distribution gate

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains the external distribution/public-launch licensing, IP and trademark decision gate. It grants no implementation authority.