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

Snapshot verified during P01.07 closure on `2026-08-23`.

Authoritative live facts verified before this snapshot:

- protected `main` after P01.07 implementation merge: `245fd67d60c02a5ed546f0cd4dc934345b5b4d42`;
- P01.07 implementation PR `#54`: merged;
- P01.07 final strict-up-to-date implementation head: `f716bd8ce8b57394ce52462e0d3ec15ecaf93bad`;
- P01.07 final implementation run/job: `32595413156` / `97085413083` — PASS;
- P01.07 implementation merge: `245fd67d60c02a5ed546f0cd4dc934345b5b4d42`;
- protected main still records P01.07 `active` until the separate closure PR merges;
- closure/state-transition PR: `#55`;
- closure branch: `chore/p01-07-close-p01-08-activate`;
- closure branch proposes P01.07 `done`, P01 progress `7 / 12 done`, and P01.08 `active` only;
- no P01.08 runtime implementation is part of closure PR #55.

Because this is a closure-branch snapshot, a future session must re-read live protected `main`, `docs/roadmap/STATE.json`, PR #55 and its latest CI. If PR #55 has merged, the canonical active package becomes P01.08. If it has not merged, protected main remains P01.07 active despite the proposed branch state.

## Current restrictions

Until closure PR #55 actually merges:

- do not claim canonical P01.07 `done` solely from the closure branch snapshot;
- do not begin P01.08 runtime implementation;
- do not start P01.09+ or P02+;
- do not implement business modules/features;
- do not implement model gateway, context engine, planner, memory/RAG, semantic business state, task graphs, simulation, risk/approval engine, capability broker, governed agents or autonomous business OS;
- do not weaken strict `governance` or main protection to obtain a merge.

After PR #55 merges, only P01.08 becomes authorized under `docs/roadmap/work-packages/P01.08.md`; all later restrictions remain.

## Exact next authorized action

Verify the complete GitHub-hosted `governance` lane for closure PR #55 on its exact final head. If all required validators/quality/P01.01-P01.07 regressions pass, the diff remains governance/documentation-only, reviews/conversations are clean and the branch is current with protected `main`, squash-merge PR #55. Do not start P01.08 runtime implementation before that merge.

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
- `docs/roadmap/work-packages/P01.07.md`
- `docs/roadmap/work-packages/P01.08.md`
- `docs/roadmap/evidence/P01.07_COMPLETION_2026-08-22.md`
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
- `docs/ai/handoffs/P01.07.md`
