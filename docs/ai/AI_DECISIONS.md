# Omnexa AI Continuity Decision Index

Status: **reference index / not an ADR registry**

This file preserves concise rationale that a future AI must not forget. It references authoritative decisions instead of replacing them. Any material change to frozen architecture must follow `docs/governance/CHANGE_CONTROL.md` and, when required, a new/superseding ADR.

## D-001 — Strict modular monolith first

**Authoritative source:** `docs/adr/ADR-0001-platform-architecture-baseline.md`, `docs/architecture/SYSTEM_ARCHITECTURE.md`.

Omnexa uses a strict modular monolith/service-ready baseline. Premature microservice extraction, shared-table ERP coupling and unrelated standalone products are rejected unless later evidence plus an ADR changes the architecture.

## D-002 — One authoritative domain owner; no private cross-domain writes

**Authoritative source:** ADR-0001, `docs/architecture/DOMAIN_OWNERSHIP.md`, `docs/architecture/MODULE_STANDARD.md`.

One owner controls each write model/capability. Other domains use governed capabilities, events, workflows or approved projections. Caches, projections and private imports are never alternate write authority.

## D-003 — AI has no alternate authority path

**Authoritative source:** ADR-0001, `docs/architecture/SYSTEM_ARCHITECTURE.md`, `docs/security/SECURITY_STANDARD.md`, ADR-0005.

Future AI acts only through authenticated/authorized, policy-controlled and auditable capabilities. Model output, prompts and retrieval never grant authority. Unrestricted DB/object-store/business-state writes remain prohibited.

## D-004 — Security and classification are platform invariants

**Authoritative source:** ADR-0005, `docs/security/SECURITY_STANDARD.md`, `docs/security/DATA_CLASSIFICATION.md`, `docs/security/THREAT_MODEL.md`.

Tenant isolation, deny-by-default authorization, secrets handling, audit authority and the four classification levels remain cross-cutting invariants and require a superseding ADR if materially changed.

## D-005 — Foundation Architecture v1 is frozen

**Authoritative source:** `docs/adr/ADR-0010-foundation-architecture-freeze.md`.

Foundation Architecture v1 is frozen. Contradictions are reopened through change control rather than normalized through local implementation convenience or CI bypasses.

## D-006 — P01 execution was strictly sequential and is complete

**Authoritative source:** `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`, `docs/governance/P01_EXIT_GATE.md`.

Each P01 kernel primitive was verified before the next package activated. While P01 was executing, exactly one package was active. Final P01.12 implementation PR #65 passed canonical run/job `32629072886 / 97168916985` and merged as `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

**Historical implication:** P01 completion created a terminal checkpoint with no active package and did not implicitly activate P02. P02 was later activated only through its separate governed readiness/activation transition.

## D-007 — Repository continuity is subordinate, not authority

**Authoritative source:** canonical governance plus this continuity convention.

`docs/ai/*` records verified snapshots, handoff evidence and rationale. Every session must re-verify `STATE.json`, branch/head, PR and CI. Conflicts make continuity stale; continuity never overrides canonical state.

## D-008 — Future AI architecture is direction only until its roadmap phases

**Authoritative source:** `docs/roadmap/MASTER_PLAN.md` P19/P20/P27 and system architecture.

Model gateway, retrieval, planner, policy/approval, capability broker, verification/replanning, governed agents and autonomous orchestration remain future direction only. Compatibility may be preserved; implementation is forbidden until canonical roadmap state authorizes it.

## D-009 — Observability is diagnostic infrastructure, not authority

**Authoritative source:** P01.07 spec/evidence and `AGENTS.md`.

Logs/traces/metrics explain runtime behavior but are not business-state, authorization or audit authority. P01.11 protected audit remains separate. P01.12 diagnostics preserve safe projections and must not expose protected payloads.

## D-010 — Health/readiness diagnostics are bounded operational signals

**Authoritative source:** P01.08 spec/evidence and `AGENTS.md`.

Liveness/readiness are distinct; dependency checks are classified and bounded; required/security-critical failures fail closed; optional failures may degrade; raw provider errors/secrets are not exposed. Diagnostics never create tenancy, authorization or business authority.

## D-011 — Completed phases do not auto-activate the next phase

**Authoritative source:** `docs/governance/DEFINITION_OF_DONE.md`, `docs/roadmap/MASTER_PLAN.md`, phase exit gates and `docs/roadmap/STATE.json`.

A completed phase can be truthfully represented at a terminal checkpoint with no active package/phase implementation scope while the next phase remains planned. This prevents evidence-based phase completion from silently granting next-phase authority.

Validators preserve strict one-active-package semantics during active execution and allow zero active packages only for a fully completed phase with its exit evidence reconciled and implementation locks false. This is a terminal-state representation, not a phase-order change or permission bypass.

## D-012 — P02 is complete; P03 requires separate readiness and activation

**Authoritative source:** `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/governance/P02_EXIT_GATE.md`, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md`.

P02 executed strictly sequentially through P02.10. Terminal implementation PR #88 final exact head `975e4925060a035780ca13b68c5437634ed0f4ea` passed canonical GitHub-hosted run/job `32904678957 / 97986011269` and merged as `88799aa41da8ce8c22540146d157d488565e2ce9`.

**Current implication after terminal closure:** P02 is `done` at 10 / 10, P02 exit is `SATISFIED`, there is no active work package, P03 remains `planned`, and both `kernel_code_authorized` and `business_feature_code_authorized` are false. P03 specification/readiness preparation and an explicit activation transition are separate governed work; P03 implementation cannot begin before that activation is accepted.
