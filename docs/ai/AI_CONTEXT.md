# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This is the first file to read inside `docs/ai/`. It exists so a new AI session can recover project context without previous chat history. It never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or the active work-package specification.

## Project identity

**Omnexa** is a Composable Enterprise Business Operating System being built architecture-first as a strict modular monolith with service-ready boundaries. Development is sequential, gated and evidence-based: foundation/kernel capabilities are implemented and verified before business modules are authorized.

## Core architectural principles

- one authoritative owner per write model/public capability;
- domain modules consume published capabilities, events, workflows or approved projections rather than private tables/packages;
- tenant isolation, authorization, auditability, classification and versioned contracts are platform invariants;
- infrastructure remains below capability/domain boundaries;
- cache, storage, telemetry, projections, analytics and AI are never alternate sources of business truth;
- AI may act only through governed identity/context, policy/approval and versioned capabilities when future roadmap phases authorize those systems;
- AI must never receive unrestricted database, object-store, payment-provider or business-state authority;
- architecture changes require change control and, when applicable, a superseding ADR.

Future AI direction remains:

```text
AI / Model / Agent
        ↓
Context / Planner
        ↓
Policy / Risk / Approval
        ↓
Governed Versioned Capability
        ↓
Owning Domain
        ↓
Infrastructure
```

This is architecture direction only. It does **not** authorize the AI platform, agents, planner, memory/RAG or autonomous execution during P01.

## Current snapshot

Snapshot reconciled during P01.08 closure on `2026-08-23`.

Authoritative live facts verified before this snapshot:

- P01.08 implementation PR `#56`: merged;
- P01.08 final strict-up-to-date implementation head: `988ef3673d49f54bbd105d3e0067ba134c66b236`;
- P01.08 final implementation run/job: `32601741049` / `97100949202` — PASS;
- P01.08 implementation merge: `a2a93454bb283464bfa144bb5a38539041e40069`;
- protected `main` subsequently advanced only through maintenance PR `#57` (`.github/dependabot.yml`) to `7e92b8acc5f5465774a2f5ef9bdc13c7dda2d2a8`;
- closure branch absorbed that main-only maintenance delta through explicit two-parent sync commit `0fce489f98a65639afa227fab9702e2423ad11b3`;
- after synchronization the closure branch compared `ahead` with `behind_by=0` against protected main and the Dependabot file was no longer part of PR #58's diff;
- protected main still records P01.08 `active` until the separate closure PR merges;
- closure/state-transition PR: `#58`;
- closure branch: `chore/p01-08-close-p01-09-activate`;
- closure branch proposes P01.08 `done`, P01 progress `8 / 12 done`, and P01.09 `active` only;
- no P01.09 runtime implementation is part of closure PR #58.

Because this is a closure-branch snapshot, a future session must re-read live protected `main`, `docs/roadmap/STATE.json`, PR #58 and its latest CI. If PR #58 has merged, the canonical active package becomes P01.09. If it has not merged, protected main remains P01.08 active despite the proposed branch state.

## P01.08 completed implementation result

P01.08 established the portable `kernel.operations` health/readiness/diagnostic boundary:

- process liveness is distinct from dependency readiness;
- startup/ready/stopping/failed lifecycle state is explicit;
- dependency checks are deterministic and classify required, optional and security-critical dependencies;
- required/security-critical failure produces `unready`; optional failure can produce `degraded`;
- each check is timeout/cancellation bounded and panic-safe;
- diagnostic reports expose stable safe states/reasons only, never raw probe/provider errors;
- build identity comes from the completed P01.01 boundary;
- health evaluation integrates with the completed P01.07 observability boundary;
- no public business status page, tenant/module aggregation, scheduler or later capability was introduced.

Final implementation verification also retained repository Go quality and P01.01-P01.07 regressions. The initial diagnostic run `32601550204 / 97100473606` remains FAIL and is not relabeled PASS.

## Current restrictions

Until closure PR #58 actually merges:

- do not claim canonical P01.08 `done` solely from the closure branch snapshot;
- do not begin P01.09 runtime implementation;
- do not start P01.10+ or P02+;
- do not implement business modules/features;
- do not implement NATS/JetStream durable event/job fabric, transactional outbox/inbox or distributed workflow timers under P01.09;
- do not implement model gateway, context engine, planner, memory/RAG, semantic business state, task graphs, simulation, risk/approval engine, capability broker, governed agents or autonomous business OS;
- do not weaken strict `governance` or main protection to obtain a merge.

After PR #58 merges, only P01.09 becomes authorized under `docs/roadmap/work-packages/P01.09.md`; all later restrictions remain.

## P01.09 bounded direction after closure

The next package is **P01.09 — Job & scheduler primitives** (`kernel.jobs`). Its allowed scope is minimal kernel-local background work only:

- deterministic job identity/type registry and unknown-job safe failure;
- enqueue/execute result model;
- bounded worker concurrency and bounded graceful shutdown/drain/cancel;
- cancellation/deadline propagation;
- explicit bounded retry/backoff policy;
- idempotency-key hook and duplicate-safe handler contract;
- simple recurring/one-shot schedules for kernel maintenance work;
- correlation/observability propagation through the completed P01.07 boundary;
- deterministic in-memory/test harness.

It must not implement P04 durable messaging/event streams, transactional outbox/inbox, P05 distributed workflow timers, business jobs, tenant-context runtime, P01.10+ capabilities or AI behavior.

## Exact next authorized action

Verify the complete GitHub-hosted `governance` lane for closure PR #58 on its exact final head. If all required validators/quality/P01.01-P01.08 regressions pass, the diff remains governance/documentation-only, reviews/conversations are clean and the branch is current with protected `main`, squash-merge PR #58. Do not start P01.09 runtime implementation before that merge.

## Authority and references

For live execution state, `docs/roadmap/STATE.json` on protected `main` is authoritative. `docs/governance/AI_EXECUTION_POLICY.md` is the mandatory AI execution policy. Architecture conflicts follow `docs/governance/CHANGE_CONTROL.md` and accepted non-superseded ADRs. `docs/ai/*` is always subordinate and may become stale.

Mandatory references:

- `AGENTS.md`
- `docs/governance/AI_EXECUTION_POLICY.md`
- `docs/governance/PRODUCT_CONSTITUTION.md`
- `docs/governance/CHANGE_CONTROL.md`
- `docs/governance/DEFINITION_OF_DONE.md`
- `docs/roadmap/STATE.json`
- `docs/roadmap/MASTER_PLAN.md`
- `docs/roadmap/work-packages/P01.08.md`
- `docs/roadmap/work-packages/P01.09.md`
- `docs/roadmap/evidence/P01.08_COMPLETION_2026-08-22.md`
- `docs/architecture/SYSTEM_ARCHITECTURE.md`
- `docs/architecture/MODULE_STANDARD.md`
- `docs/architecture/DOMAIN_OWNERSHIP.md`
- `docs/security/SECURITY_STANDARD.md`
- `docs/security/DATA_CLASSIFICATION.md`
- `docs/security/THREAT_MODEL.md`
- `docs/adr/ADR-0001-platform-architecture-baseline.md`
- `docs/adr/ADR-0005-security-data-classification-baseline.md`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`
- `docs/ai/AI_STATE.yaml`
- `docs/ai/AI_EXECUTION_PROTOCOL.md`
- `docs/ai/AI_DECISIONS.md`
- `docs/ai/handoffs/P01.08.md`
