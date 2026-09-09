# Omnexa

**Composable Enterprise Business Operating System**

Omnexa is a governed modular platform above the scope of a conventional ERP. ERP, CRM, finance, commerce, POS, payments, website/CMS, portals, workflow, integrations, analytics, low-code and AI remain governed domain families on one platform foundation.

> **Architecture state:** Foundation Architecture v1 is **FROZEN**; P00, P01, P02 and P03 are complete.

> **Current canonical state:** **P04 — Data, Jobs & Event Fabric is ACTIVE at 6 / 10 completed packages. P04.01-P04.06 are DONE and P04.07 is the sole ACTIVE package.** The P04.06→P04.07 transition was accepted through source #265 / unchanged promotion #266 at exact head `02b58b4245da5eef7a3ab1090698cd10a4832d90`, source Governance `34411353055`, promotion Governance `34412007672`, and protected-main merge/read-back `5af9383c3e973c4055eea66d48462b7c9a2a5858`. P04.07 runtime has **not** started; worker slots/tasks/runtime branches and migration reservations remain zero. P04.08-P04.10 remain locked. `business_feature_code_authorized=false`.

`docs/roadmap/STATE.json` is the canonical machine-readable execution cursor. This README is a human-readable status/instruction mirror only and never grants package authority by itself.

## Project progress dashboard

Phase-count view: **4 of 28 phases completed**, with **P04 active**. Current P04 package completion is **6 / 10 (60%)**.

| Phase | Status | Current package state |
|---|---|---|
| P00 | DONE | Product Constitution & Architecture Freeze |
| P01 | DONE | Omnexa Kernel |
| P02 | DONE | Identity, Tenancy & Organization |
| P03 | DONE | Module Runtime |
| P04 | ACTIVE | 6/10 packages done; P04.07 sole ACTIVE package; runtime lease not yet opened |
| P05-P27 | PLANNED / LOCKED | Not started unless a later governed transition activates them |

## Current P04 package status

| Package | Status | Current truth |
|---|---|---|
| P04.01 | DONE | Event envelope & identity contract accepted |
| P04.02 | DONE | Publish/subscribe abstraction & ownership accepted |
| P04.03 | DONE | Durable consumer/checkpoint baseline accepted |
| P04.04 | DONE | Transactional outbox accepted |
| P04.05 | DONE | Consumer inbox/idempotency accepted |
| P04.06 | DONE | Retry/backoff/quarantine implementation and completion evidence accepted; migration 3 historical/immutable |
| P04.07 | ACTIVE | Schema registry/compatibility/validation contract active; runtime still gated by post-activation continuity and a later implementation-plan/worker-plan |
| P04.08-P04.10 | PLANNED / LOCKED | Not started |

## Accepted P04.06 implementation chain

- **Wave 2A — durable retry/quarantine state:** source #227 / promotion #228 -> `c0311bd79b4d60e9a71fcf0934d9a58f561e1c88`; immutable `kernel.events` migration 3 accepted.
- **Wave 2B — PostgreSQL CAS/lease persistence:** source #230 / promotion #231 -> `3e6bccf8164a105f92c0b576a99f20a498ec9025`.
- **Wave 2C T01 — retry execution composition:** source #238 / promotion #239 -> `920048584ba039d7c29bd25398c11a2138b2c0b3`.
- **Wave 2C T02 — bounded PostgreSQL due discovery:** source #240 / promotion #241 -> `75cc0f0f590a049639ef0ea5c19a56b22ffab7d1`.
- **Wave 2C T03 — quarantine/checkpoint crash-gap recovery:** source #242 / promotion #243 -> `66e1a330c42a37d252e1f32221020d1460aec48f`.
- **Wave 2C T04 — canonical verifier:** source #244 / promotion #245 -> `a049164ca3264eff4b71a97dca8e1dbba0c276ae`.
- **Wave 2C status closure:** source #246 / promotion #247 -> `afc1bb6d47cc988e2731bd16da3467d72d55af77`.
- **Wave 2D control:** source #249 / promotion #250 -> `43943c91c7491248d3ba1573b72d800da5825f96`.
- **Wave 2D T05 — P04.05 applied/already-applied → claimed retry resolution:** source #251 / promotion #253; source Governance #730, promotion Governance #732; protected merge/read-back `180d741ef81a1c5e7d33c6bafbc9805d1c801757`.
- **Wave 2D T06 — real PostgreSQL restart persistence:** final source #254 / promotion #255; source Governance #735 / `34396694586`, promotion Governance #736 / `34397580962`; protected merge/read-back `df60c3f406dfbc03e4156c60b8103a23ef257e32`.
- **Wave 2D T07 — closure verifier / implementation evidence:** source #257 / promotion #258, exact head `b0d047c93e8c53e9da62a96e36c4a1a8a7ce634f`; source Governance #738 / `34400644888`, promotion Governance #739 / `34401644140`; protected merge/read-back `bdb96bb7f0dabf5b103acd78335599cc59e92b69`.
- **Wave 2D ledger closure:** source #259 / promotion #260 -> protected merge/read-back `6050bc2d970b5cf83118115557e952e3e91e0f96`; all worker/Supervisor leases released.

The earlier T06 PR #252 / Governance #734 is diagnostic only: its stale PR-event base correctly triggered M2 scope rejection after T05 moved main. #252 was closed unmerged; no validator bypass or code weakening was used.

## P04.06 completed boundary retained

Accepted P04.06 implementation/verifier evidence proves:

- finite deterministic retry/backoff and safe terminal/interruption classification;
- durable scheduled/quarantined/resolved state with immutable migration 3;
- PostgreSQL due discovery, revision CAS, bounded lease claim and stale-transition protection;
- explicit quarantine-before-checkpoint ordering and restart-safe checkpoint crash-gap recovery;
- P04.05 `InboxApplied` / `InboxAlreadyApplied` resolves stale retry scheduling only through an authoritative retry claim and cannot re-run the protected mutation;
- scheduled/quarantined/resolved retry evidence survives a complete writer-pool close + fresh reader-pool/store boundary on real PostgreSQL;
- post-restart scheduled work is still non-executable until authoritative due discovery + `ClaimDue` succeeds;
- the canonical P04.06 verifier checks Wave 2D closure evidence while retaining P04.01-P04.05 regression, race and build boundaries.

Checkpoint, inbox completion and retry/quarantine are distinct facts. No retry record, quarantine row, checkpoint, inbox row, EventID, attempt number, claim or provider receipt grants authorization or proves exactly-once processing.

Completion evidence: `docs/roadmap/evidence/P04.06_COMPLETION_2026-09-09.md`.

## Accepted P04.07 preparation and activation

P04.07 preparation was accepted through source #262 / unchanged promotion #263 at exact head `94789264824a34e3f2608283cb6b1c2f158e0b81`, source Governance `34407394506`, promotion Governance `34408137084`, and protected-main read-back `23f3dca090338b5debbf1af02b6db49b2f03f30b`.

The separate P04.06 closure / P04.07 activation then passed:

- source #265;
- unchanged promotion #266;
- exact source/promotion head `02b58b4245da5eef7a3ab1090698cd10a4832d90`;
- source Governance #747 / `34411353055` — PASS;
- promotion Governance #748 / `34412007672` — PASS;
- honest SELF/Supervisor reviews with no fabricated independent approval;
- zero unresolved review threads;
- expected-head guarded merge/read-back `5af9383c3e973c4055eea66d48462b7c9a2a5858`.

Issue #264 is complete. Activation changed governance/state/continuity only and did not create P04.07 runtime, migration 4 or provider/vendor registry authority.

## Current P04.07 active boundary

Owner: `kernel.events`.

P04.07 is canonically ACTIVE. Its accepted provider-neutral contract is `docs/roadmap/work-packages/P04.07.md` and handoff is `docs/ai/handoffs/P04.07.md`.

Retained P04.07 laws include:

- payload schema version is distinct from P04.01 envelope version;
- owner + stable event type + explicit payload schema version forms the core registry identity;
- accepted historical schema versions are immutable and fingerprint-bound;
- canonicalization/fingerprinting must be deterministic and bounded;
- compatibility modes are explicit `backward`, `forward`, `full`, `exact`;
- the exact predecessor comparison set must be frozen before runtime acceptance;
- unknown/unregistered/unsupported schema fails before protected handler mutation;
- payload validation is bounded and local-only;
- remote HTTP(S), DNS, registry-to-registry and arbitrary filesystem schema/reference fetching is forbidden;
- dynamic code/eval/plugin/shell/schema callback execution is forbidden;
- schema validation never grants authorization;
- tenant-specific schema forks/shadows are forbidden in V1;
- checkpoint, inbox, retry/quarantine and schema-validation evidence remain separate facts;
- no provider/vendor schema registry is selected;
- accepted migration 3 remains immutable P04.06 history and migration 4 is not assumed or reserved.

## Issue #267 continuity state

Coordination: GitHub Issue #267.

This post-activation continuity carrier may reconcile only:

1. `AGENTS.md`
2. `README.md`
3. `docs/roadmap/STATUS.md`
4. `docs/ai/AI_STATE.yaml`
5. `docs/ai/AI_CONTEXT.md`
6. `docs/ai/handoffs/P04.07.md`
7. `docs/governance/P04_06_P04_07_TRANSITION_TRANSACTION.md`

It must not modify canonical `STATE.json`, the P04 package sequence, activation validators, `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`, runtime Go source, migrations, workflows or modules.

The machine ledger remains **0 worker slots, 0 open slots, 0 runtime tasks, 0 P04.07 runtime branches and no active Supervisor write lease**. No migration reservation exists. A later worker requires a fresh separately governed implementation-plan/worker-plan from current protected main.

## Multi-agent development operating model

Parallel AI/human development remains subordinate to the single canonical phase/work-package cursor.

**Framework safe envelope:** up to **4-6 active agents with no more than 3 concurrent write agents** only when a governed active plan actually opens those slots. The framework cap is not permission to invent workers, branches, packages or authority.

Mandatory concurrency rules: `docs/governance/MULTI_AGENT_ORCHESTRATION.md`. Supervisor workflow: `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`. Active machine plan: `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json`.

### M2 required CI enforcement

The required `governance` job validates active-plan/canonical-state alignment, task/slot identity, path leases, overlaps, migration reservations, dependency order, registered worker PR scope and registered worker base/ancestry against live protected main. Unknown `agent/*` branches, stale worker PR bases and out-of-budget paths fail closed.

### Supervisor submission / interrupt protocol

A completed task announces exactly:

`Work Done and Submitted`

After protected-main merge/readback the Supervisor announces exactly:

`New changes have been merged — please merge these changes into your branch first, then resume your own work.`

Workers synchronize new protected main, invalidate stale CI and acknowledge:

`Sync Complete — Resuming Work`

The active plan's coordination issue is the authoritative signal channel when a governed worker plan is live.

### New Agent Onboarding

A newly arriving agent starts from protected main and receives no authority merely by arriving. Supervisor checks the machine worker-slot ledger first. If an authorized slot is open, assignment follows dependency/merge order. If no slot is open, Supervisor says exactly:

`Go Home Come Back Next Time`

The current ledger has **0 open runtime slots**. No new package, migration, module or write lease may be invented to accommodate an arriving worker.

### Live protected-main freshness rule

`required_main_ref = main` is live freshness authority. Stored SHA fields in plans are audit snapshots only. Developer orchestration is not P20 product-agent runtime.

## Agent Working Instructions

Every material agent task checks these instructions **at task start and again before PR submission**:

1. inspect live open Issues and all open PRs/MRs before new implementation; verify exact heads, draft/mergeability, required reviews, unresolved threads and CI, merge every safe approved green ready PR/MR first, and resolve or explicitly document blocked/draft/red/stale items rather than bypassing them;
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

Read `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P04_06_P04_07_TRANSITION_TRANSACTION.md`, `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`, retained P04.01-P04.06 completion evidence, `docs/roadmap/work-packages/P04.07.md`, `docs/ai/handoffs/P04.07.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/MULTI_AGENT_ORCHESTRATION.md`, `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md` and `docs/ai/ACTIVE_MULTI_AGENT_PLAN.json` before material work.

## Separate P04.07 implementation-plan gate

After Issue #267 continuity is accepted/read back, do not start runtime immediately. A new separately governed P04.07 implementation-plan/worker-plan must freeze:

- exact implementation paths and non-overlapping worker leases;
- supported bounded schema representation/dialect or subset;
- canonicalization/fingerprint algorithm and encoding;
- compatibility predecessor-set semantics;
- validator size/depth/time/memory/collection bounds and safe error contract;
- local-reference policy and graph limits if supported;
- registry persistence form;
- only if durable persistence is required, exact next `kernel.events` migration owner/version/path/data budget after fresh-main migration preflight.

Only after that plan passes its own source/promotion Governance and protected read-back may the exact bounded P04.07 runtime lease begin.

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

Protected integration requires PRs, strict required `governance`, conversation resolution, fresh-main checks and guarded merge. Canonical CI uses GitHub-hosted `ubuntu-24.04`; local/self-hosted governance evidence is not substituted.

## External distribution gate

The repository is public and the current `LICENSE` remains GPLv3. Issue #4 remains the external distribution/public-launch licensing, IP and trademark decision gate. It grants no P04.07 runtime authority.
