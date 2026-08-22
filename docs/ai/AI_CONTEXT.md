# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This is the first file to read inside `docs/ai/`. It exists so a new AI session can recover project context without previous chat history. It is not a replacement for `AGENTS.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards, or the active work-package specification.

## Project identity

**Omnexa** is a Composable Enterprise Business Operating System being built architecture-first as a strict modular monolith with service-ready boundaries. The current development model is sequential, gated and evidence-based: foundation/kernel capabilities are implemented and verified before business modules are authorized.

## Core architectural principles

- one authoritative owner per write model and public capability;
- domain modules consume published capabilities, events, workflows or approved projections rather than private tables/packages;
- tenant isolation, authorization, auditability, classification and versioned contracts are platform invariants;
- infrastructure remains below capability/domain boundaries;
- cache, projections, analytics and AI are never alternate sources of business truth;
- AI may plan or propose only through governed identity, policy/approval and versioned capabilities when the roadmap authorizes those systems;
- AI must never receive unrestricted database, object-store, payment-provider or business-state authority;
- architecture changes require change control and, where applicable, a superseding ADR.

The intended future direction is:

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

This is architectural direction only. It does **not** authorize the AI platform, agents, planner, memory/RAG or autonomous execution today.

## Current snapshot

Snapshot verified: `2026-08-22T20:27:00+05:00`.

- canonical branch: `main`
- snapshot base main commit: `753fc54dbb35071807bcc3fd51743f289c400b98`
- phase: `P01 — Omnexa Kernel`
- canonical work-package state: `P01.06 — Object & File Storage Abstraction` is `active`
- completed P01 packages: `5 / 12`
- active implementation PR: `#50`
- active work branch: `feat/p01-06-storage-clean`
- active work head: `7028145965aa38194bdb0524ec511b98f206461c`
- latest observed CI: run `32581123619`, job `97050592422`, **FAIL**
- failure point: `Verify P01.01 Go workspace and build skeleton`
- failure reason: the committed `go.sum` checksum for `gopkg.in/yaml.v3@v3.0.1` does not match the canonical downloaded checksum; repository Go quality passed before this failure, while P01.02-P01.06 were skipped.

This section is a snapshot and can become stale. Always re-read `STATE.json`, the active work-package spec, the current branch/PR and latest CI before acting.

## Current restrictions

Until canonical state changes through governed evidence, do not:

- mark P01.06 done;
- activate or implement P01.07;
- start P02+;
- implement business modules/features;
- implement model gateway, context engine, planner, memory/RAG, semantic business state, task graphs, simulation, risk engine, approval engine, capability broker, governed agents or autonomous business OS;
- infer tenant/authorization/business ownership from infrastructure object keys;
- weaken governance/quality gates to obtain green CI;
- update roadmap state merely because implementation code exists.

## Exact next authorized action

On PR `#50`, correct the `gopkg.in/yaml.v3@v3.0.1` checksum recorded in `go.sum` to the canonical checksum exposed by run `32581123619`, then rerun the full GitHub-hosted `governance` workflow. Do not start P01.07.

## Authority and references

For live execution state, `docs/roadmap/STATE.json` is authoritative. For architectural conflicts, follow `docs/governance/CHANGE_CONTROL.md` and accepted non-superseded ADRs. `docs/ai/*` is always subordinate.

Mandatory references:

- `AGENTS.md`
- `docs/roadmap/STATE.json`
- `docs/roadmap/MASTER_PLAN.md`
- `docs/roadmap/work-packages/P01.06.md`
- `docs/architecture/SYSTEM_ARCHITECTURE.md`
- `docs/architecture/MODULE_STANDARD.md`
- `docs/architecture/DOMAIN_OWNERSHIP.md`
- `docs/security/SECURITY_STANDARD.md`
- `docs/security/DATA_CLASSIFICATION.md`
- `docs/security/THREAT_MODEL.md`
- `docs/governance/CHANGE_CONTROL.md`
- `docs/adr/ADR-0001-platform-architecture-baseline.md`
- `docs/adr/ADR-0005-security-data-classification-baseline.md`
- `docs/adr/ADR-0010-foundation-architecture-freeze.md`
- `docs/ai/AI_STATE.yaml`
- `docs/ai/AI_EXECUTION_PROTOCOL.md`
- `docs/ai/AI_DECISIONS.md`
- `docs/ai/handoffs/P01.06.md`
