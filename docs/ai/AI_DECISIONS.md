# Omnexa AI Continuity Decision Index

Status: **reference index / not an ADR registry**

This file preserves concise rationale that a future AI must not forget. It references authoritative decisions instead of replacing them. Any material change to frozen architecture must follow `docs/governance/CHANGE_CONTROL.md` and, when required, a new/superseding ADR.

## D-001 — Strict modular monolith first

**Authoritative source:** `docs/adr/ADR-0001-platform-architecture-baseline.md`, `docs/architecture/SYSTEM_ARCHITECTURE.md`.

**Why:** Omnexa needs coherent shared platform invariants and domain boundaries without paying premature microservice complexity.

**Rejected:** traditional shared-table ERP monolith; microservices from day one; unrelated standalone products as the primary architecture.

**Must not change without ADR:** strict modular monolith/service-ready baseline and evidence-based service extraction criteria.

## D-002 — One authoritative domain owner; no private cross-domain writes

**Authoritative source:** ADR-0001, `docs/architecture/DOMAIN_OWNERSHIP.md`, `docs/architecture/MODULE_STANDARD.md`.

**Why:** duplicate write ownership destroys lifecycle independence, auditability and deterministic business truth.

**Rejected:** modules writing another module's private tables/packages or treating projections/caches as alternate authority.

**Future dependency:** events, workflows, integrations, analytics and AI all depend on stable owner-published capabilities.

## D-003 — AI has no alternate authority path

**Authoritative source:** ADR-0001, `docs/architecture/SYSTEM_ARCHITECTURE.md`, `docs/security/SECURITY_STANDARD.md`, `docs/security/DATA_CLASSIFICATION.md`, ADR-0005.

**Why:** model output and retrieved content are untrusted; prompt/tool injection must not create privilege.

**Required direction:** AI acts through authenticated/authorized, policy-controlled, auditable capabilities owned by domains.

**Rejected:** unrestricted DB writes, direct object-store/business-state authority, relaxed AI permissions, prompts as authorization.

**Future dependency:** P19 Intelligence Platform, P20 Governed AI Agents and P27 Autonomous Business OS.

## D-004 — Security and classification are platform invariants

**Authoritative source:** `docs/adr/ADR-0005-security-data-classification-baseline.md`, `docs/security/SECURITY_STANDARD.md`, `docs/security/DATA_CLASSIFICATION.md`, `docs/security/THREAT_MODEL.md`.

**Why:** tenant boundaries, secrets, business records, external integrations and future AI need consistent handling across every surface.

**Must not change without superseding ADR:** tenant isolation model, authorization baseline, classification levels, secrets handling, audit authority and AI execution authority.

## D-005 — Foundation Architecture v1 is frozen

**Authoritative source:** `docs/adr/ADR-0010-foundation-architecture-freeze.md`.

**Why:** architecture acceptance must be durable and separate from implementation convenience. Contradictions are reopened through change control, not bypassed locally.

**Rejected:** treating local-only work or unavailable gates as permission to skip canonical governance.

## D-006 — P01 execution is strictly sequential

**Authoritative source:** `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json`.

**Why:** each kernel primitive becomes a verified dependency for the next package; multiple active foundation packages make evidence and authority ambiguous.

**Current implication:** P01.10 implementation is merged with exact-head canonical PASS evidence. The separate closure branch `chore/p01-10-close-p01-11-activate` reconciles P01.01-P01.10 `done`, P01.11 as the sole active package and P01.12 planned. P01.11 runtime must not start until that closure merges; after the closure merge the execution session must stop, and P01.11 implementation starts only in a new governed session.

## D-007 — Repository continuity is a subordinate snapshot, not a new source of truth

**Authoritative source:** this continuity convention plus existing governance authority; architectural conflict rules remain in `docs/governance/CHANGE_CONTROL.md`.

**Why:** chat/session memory is temporary, but duplicating canonical state as equal authority would create stale-state conflicts.

**Decision:** `docs/ai/*` records verified snapshots, handoff evidence and rationale. Every new session must re-verify `STATE.json`, branch/head, PR and CI. Conflicts make the continuity value stale; they never authorize overwriting canonical state.

**When an ADR is required:** only if this operational convention is later changed in a way that alters frozen architecture, roadmap gates, security authority or public contracts.

## D-008 — Future AI architecture is direction only until its roadmap phases

**Authoritative source:** `docs/roadmap/MASTER_PLAN.md` P19/P20/P27 and `docs/architecture/SYSTEM_ARCHITECTURE.md` AI architecture.

**Why:** documenting the target now prevents foundations from creating incompatible shortcuts, while scope gates prevent premature implementation.

**Direction:** model gateway, context/retrieval, planner, risk/policy, approval, capability broker, verification/replanning, governed agents and autonomous business orchestration belong to future authorized phases.

**Current rule:** compatibility may be preserved; implementation is forbidden until canonical roadmap state authorizes it.

## D-009 — Observability is diagnostic infrastructure, not authority

**Authoritative source:** `docs/roadmap/work-packages/P01.07.md`, `docs/roadmap/evidence/P01.07_COMPLETION_2026-08-22.md`, `AGENTS.md`.

**Why:** logs/traces/metrics must help explain runtime behavior without becoming a second correctness path, business-state owner, authorization signal or audit substitute.

**Implemented P01.07 direction:** standard structured logging plus isolated OpenTelemetry trace/metric SDK providers, vendor-neutral exporter injection, bounded lifecycle and fail-closed redaction. Correlation/trace identifiers remain diagnostic identifiers only.

**Rejected:** proprietary telemetry coupling in the kernel contract, global-provider mutation as hidden process state, raw secret/classified/error leakage, audit semantics in ordinary logs, or treating telemetry availability as application correctness.

**Future dependency:** P01.11 audit transport must remain a separate protected contract from ordinary observability.

## D-010 — Health/readiness diagnostics are bounded operational signals, not authority

**Authoritative source:** `docs/roadmap/work-packages/P01.08.md`, `docs/roadmap/evidence/P01.08_COMPLETION_2026-08-22.md`, `AGENTS.md`.

**Why:** a process can be alive while temporarily unsuitable for work; optional dependency degradation must not be conflated with process death, and operational diagnostics must never become an authorization or business-state channel.

**Implemented P01.08 direction:** liveness and readiness are distinct; dependency checks have explicit required/optional/security-critical classification; required/security-critical failures fail readiness closed; optional failures can degrade; checks are timeout/cancellation bounded and panic-safe; reports expose stable machine states/reasons and build identity without retaining raw probe/provider errors.

**Rejected:** one generic health boolean, unbounded probes, raw secret/provider error exposure, treating object keys/connection details as diagnostics, module/tenant aggregation before P03, public business status semantics, or Kubernetes-specific architecture as the kernel contract.

**Future dependency:** P01.11 transport health may expose safe operational state but must not log protected audit payloads or create authorization/business authority.
