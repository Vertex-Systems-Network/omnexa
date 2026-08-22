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
- cache, storage, projections, analytics and AI are never alternate sources of business truth;
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

Snapshot verified during P01.06 closure on `2026-08-22`.

Authoritative live facts verified before this snapshot:

- canonical `main` before the closure transition: `f7867d9e1c570e3abbed90740970acf7b5a30bd7`;
- P01.06 implementation PR `#50`: merged;
- P01.06 final strict-up-to-date implementation run/job: `32588244996` / `97067784835` — PASS;
- P01.06 implementation merge: `f7867d9e1c570e3abbed90740970acf7b5a30bd7`;
- canonical main still records P01.06 active until the closure PR merges;
- closure/state-transition PR: `#53`;
- closure branch: `chore/p01-06-close-p01-07-activate`;
- closure branch state proposes P01.06 `done`, P01 progress `6 / 12 done`, and P01.07 `active` only;
- no P01.07 runtime implementation is part of the closure PR.

Because this is a branch snapshot, a future session must re-read live `main`, PR #53 and its latest CI. If PR #53 has merged, the canonical active package becomes P01.07. If it has not merged, canonical main remains P01.06 active despite the proposed branch state.

## Current restrictions

Until the closure transition is actually merged:

- do not claim canonical P01.06 `done` solely from the branch snapshot;
- do not begin P01.07 runtime implementation;
- do not start P01.08+ or P02+;
- do not implement business modules/features;
- do not implement model gateway, context engine, planner, memory/RAG, semantic business state, task graphs, simulation, risk/approval engine, capability broker, governed agents or autonomous business OS;
- do not weaken strict `governance` or main protection to obtain a merge.

After PR #53 merges, only P01.07 becomes authorized under its active work-package specification; all later restrictions remain.

## Exact next authorized action

Verify the complete GitHub-hosted `governance` lane for closure PR `#53`. If the exact current closure head is green, scope/reviews/mergeability are clean and the branch is current with protected `main`, squash-merge PR #53. Do not start P01.07 implementation before that merge.

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
- `docs/roadmap/work-packages/P01.06.md`
- `docs/roadmap/work-packages/P01.07.md`
- `docs/roadmap/evidence/P01.06_COMPLETION_2026-08-22.md`
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
- `docs/ai/handoffs/P01.06.md`
