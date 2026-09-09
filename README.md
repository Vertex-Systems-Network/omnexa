# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI are governed domain families on one platform foundation.

> **Architecture state:** Foundation Architecture v1 is **FROZEN**; P00, P01, P02 and P03 are complete.

> **Current canonical state:** **P04 — Data, Jobs & Event Fabric is ACTIVE at 5 / 10 completed packages. P04.01-P04.05 are DONE and P04.06 Retry/Backoff, Terminal Failure & Dead-Letter/Quarantine Policy is the sole ACTIVE package.** Wave 2A, Wave 2B and Wave 2C T01-T04 are accepted through protected `main@afc1bb6d47cc988e2731bd16da3467d72d55af77`. P04.06 is **not complete**. Fresh-main audit Issue #248 identifies two explicit remaining acceptance-evidence gaps: P04.05 `already_applied` retry-resolution composition and real-PostgreSQL restart preservation of scheduled/quarantined/resolved retry state. `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` is the canonical machine-readable execution cursor. This README is a human-readable status and working-instruction mirror only; it never grants implementation authority by itself.

## Project progress dashboard

Phase-count view: **4 of 28 phases completed**, with **P04 active**. Current P04 package completion is **5 / 10 (50%)**.

| Phase | Status | Current package state |
|---|---|---|
| P00 | DONE | Product Constitution & Architecture Freeze |
| P01 | DONE | Omnexa Kernel |
| P02 | DONE | Identity, Tenancy & Organization |
| P03 | DONE | Module Runtime |
| P04 | ACTIVE | 5/10 packages done; P04.06 sole ACTIVE package |
| P05-P27 | PLANNED / LOCKED | Not started unless a later governed transition activates them |

## Current P04 package status

| Package | Status | Current truth |
|---|---|---|
| P04.01 | DONE | Event envelope & identity contract accepted |
| P04.02 | DONE | Publish/subscribe abstraction & ownership accepted |
| P04.03 | DONE | Durable consumer/checkpoint baseline accepted |
| P04.04 | DONE | Transactional outbox accepted |
| P04.05 | DONE | Consumer inbox/idempotency accepted |
| P04.06 | ACTIVE | Wave 2A + Wave 2B + Wave 2C T01-T04 accepted; Issue #248 closure gaps remain |
| P04.07-P04.10 | PLANNED / LOCKED | Not started |

### Accepted P04.06 evidence chain

- **Wave 2A — durable retry/quarantine state:** source #227 / promotion #228; protected merge/read-back `c0311bd79b4d60e9a71fcf0934d9a58f561e1c88`; immutable `kernel.events` migration 3 accepted.
- **Wave 2B — PostgreSQL retry-state CAS/lease persistence:** source #230 / promotion #231; protected merge/read-back `3e6bccf8164a105f92c0b576a99f20a498ec9025`.
- **Wave 2C T01 — retry execution composition:** source #238 / promotion #239; protected merge/read-back `920048584ba039d7c29bd25398c11a2138b2c0b3`.
- **Wave 2C T02 — bounded PostgreSQL due discovery:** source #240 / promotion #241; protected merge/read-back `75cc0f0f590a049639ef0ea5c19a56b22ffab7d1`.
- **Wave 2C T03 — quarantine/checkpoint crash-gap recovery:** source #242 / promotion #243; protected merge/read-back `66e1a330c42a37d252e1f32221020d1460aec48f`.
- **Wave 2C T04 — canonical P04.06 verifier:** source #244 / promotion #245; protected merge/read-back `a049164ca3264eff4b71a97dca8e1dbba0c276ae`.
- **Wave 2C lease/status closure:** source #246 / promotion #247; protected merge/read-back `afc1bb6d47cc988e2731bd16da3467d72d55af77`; all Wave 2C worker/Supervisor leases released.
- **Wave 2D closure-gap control candidate:** Issue #248. The candidate plan reserves T05/T06 and isolated Supervisor T07 only. Until this control source, unchanged promotion and protected-main read-back are accepted, those reservations grant **no runtime write authority**.

P04.07-P04.10 remain locked. No broker/provider-specific DLQ, migration 4, business feature, exactly-once/global-ordering claim or AI/model/agent product runtime is authorized by P04.06 closure work.

## P04.06 Wave 2D candidate leases

Control branch base: protected `main@afc1bb6d47cc988e2731bd16da3467d72d55af77`  
Coordination: GitHub Issue #248  
Current control state: **candidate awaiting Governance — no runtime mutation yet**

| Merge order | Task | Branch | Exact write lease | Dependency |
|---:|---|---|---|---|
| 1 | P04.06-T05 / Agent-01 — inbox/retry composition | `agent/20260909-p04-06-inbox-retry-closure` | `kernel/internal/events/retry_inbox.go`, `kernel/internal/events/retry_inbox_test.go` | none |
| 2 | P04.06-T06 / Agent-02 — PostgreSQL restart evidence | `agent/20260909-p04-06-restart-persistence` | `kernel/internal/events/retry_state_postgres_restart_integration_test.go` | none |
| 3 | P04.06-T07 / Supervisor — closure verifier/evidence | `supervisor/20260909-p04-06-closure-verifier` | `scripts/verify_p04_06.sh`, P04.06 handoff/completion evidence, README mirror | accepted T05 + T06 |

The logical worker-slot ledger contains **2 reserved/occupied candidate slots and 0 open slots**. These are not executable runtime authority until the Wave 2D control plan is accepted through source Governance, unchanged promotion Governance, guarded protected merge and protected-main read-back. No additional worker slot may be invented.

## Multi-agent development operating model

Omnexa development may use parallel AI/human agents to accelerate delivery, but parallelism stays subordinate to the single canonical phase/work-package cursor.

**Framework safe envelope:** up to **4-6 active agents with no more than 3 concurrent write agents** when the governed active plan actually opens those slots. The numerical framework cap is not permission to invent workers, branches, packages or authority.

Mandatory concurrency rules: `docs/governance/MULTI_AGENT_ORCHESTRATION.md`. Supervisor interruption/merge/sync/onboarding: `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`. Active machine plan: `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`.

### M2 required CI enforcement

The M2 gate executes inside the required `governance` job through `scripts/verify_go_quality.sh` before expensive Go tooling. It validates:

- active-plan authority vs canonical `STATE.json`;
- worker-slot/task identity and capacity;
- active write paths vs forbidden paths;
- pairwise write-path overlap;
- migration owner/version/path reservation collisions;
- dependency DAG and deterministic merge order;
- registered worker PR changed paths vs declared budget;
- registered worker PR base/ancestry vs live protected `main`;
- helper behavior through focused dependency-free unit tests.

Unknown `agent/*` branches, scope violations, stale worker PR bases, authority mismatches, path overlaps, migration collisions and dependency-order/cycle violations fail closed. A governance/control branch is not falsely treated as a worker task; it still must pass the normal protected Governance/promotion flow.

### Supervisor submission / interrupt protocol

A completed task announces exactly:

`Work Done and Submitted`

with task, branch, exact head SHA, PR, CI/evidence and instruction-check metadata.

When a valid submission arrives, the Supervisor checkpoints/pauses its own task, reviews the submitted exact head, and either requests changes or routes the approved head through required source/Governance/unchanged-promotion/protected-main integration. Supervisor approval never bypasses protected integration.

After protected-main merge/readback the Supervisor announces exactly:

`New changes have been merged — please merge these changes into your branch first, then resume your own work.`

Every active worker resolves and synchronizes the new protected `main`, re-checks instructions/leases/dependencies, invalidates stale CI when its head changes, and only then resumes. Recommended acknowledgement:

`Sync Complete — Resuming Work`

The active plan's coordination issue is the authoritative signal channel. For Wave 2D it is Issue #248.

### New Agent Onboarding

A newly arriving development agent always starts from protected `main`. Arrival alone grants no branch, module, task or lease.

The Supervisor immediately checks the AI-Native worker-slot ledger:

- if an authorized slot is `open`, resolve live main, synchronize the slot branch, assign the agent, mark `occupied`, record identity/SHA/start status, then allow work;
- if multiple slots are open, dependency/merge order decides priority;
- never invent a new module, exceed the concurrency cap or activate a locked phase because an agent arrived;
- if no slot is open, stop onboarding and say exactly: `Go Home Come Back Next Time`.

Wave 2D has no open slot. Its two candidate slots are already logically reserved for T05/T06 and remain non-executable until control acceptance. An unrelated arriving worker therefore receives `Go Home Come Back Next Time`.

### Live protected-main freshness rule

`required_main_ref = main` is live freshness authority. Stored SHA fields are audit snapshots only. This developer orchestration is not P20 product-agent runtime and does not activate a future business phase.

## Agent Working Instructions

Every material agent task checks these instructions **at task start and again before PR submission**:

1. inspect live open Issues and all open PRs/MRs before new implementation; verify exact heads, draft/mergeability, required reviews and CI, merge every safe approved green ready PR/MR first, and resolve or explicitly document blocked/draft/red/stale items rather than bypassing them;
2. re-read protected `main` and canonical `docs/roadmap/STATE.json` after any accepted merge;
3. confirm active phase/work package and owning module/domain/kernel capability;
4. read `AGENTS.md`, applicable work-package spec/handoff, AI execution policy, multi-agent orchestration contract, Supervisor workflow and active plan;
5. newly arriving agents start from protected `main` and require an authorized `open` slot before switching to a task branch;
6. if no slot is open, Supervisor responds exactly `Go Home Come Back Next Time` and grants no work authority;
7. record task/agent identity, slot, branch and current live-main/base/sync evidence;
8. declare read/write/forbidden/shared paths and check active overlap;
9. resolve dependencies; do not code against guessed future contracts;
10. reserve exact owner/path/version/data budget before schema mutation;
11. resolve current protected `main` and the active coordination issue before each material mutation;
12. after another accepted merge, sync new main before resuming and rerun stale tests/CI;
13. registered worker branches must remain registered in the active plan; renaming/creating an `agent/*` branch does not bypass M2 — unknown `agent/*` PRs fail required Governance;
14. do not write outside the declared task budget; registered worker PR scope is machine-enforced by M2;
15. registered worker PRs must be based on current protected `main`; stale worker PRs fail M2 and must synchronize/resubmit;
16. when ready for review, announce `Work Done and Submitted` with exact-head metadata;
17. keep cross-module private writes/imports forbidden;
18. require exact-final-head tests/CI/review and re-check protected-main freshness before merge;
19. compare effective working instructions with this README;
20. **if instructions changed materially, update this section in the same PR**;
21. if unchanged, record: `Agent instructions checked — README instruction delta: none`.

Material instruction changes include phase/package, role, owner, slot availability/assignment, branch/base/sync strategy, path leases, migration budget, dependency/contract assumptions, required Issues/PR intake, M2/tests/gates, Supervisor/merge/onboarding process, coordination channel, tool/network/secret restrictions or stop conditions.

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
