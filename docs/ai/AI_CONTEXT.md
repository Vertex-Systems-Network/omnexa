# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This file never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Current governed checkpoint

- Foundation Architecture v1: FROZEN.
- P00: DONE.
- P01: DONE — 12 / 12; exit gate SATISFIED.
- P02: ACTIVE — 0 / 10 done.
- sole active package: `P02.01 — Principal & user identity foundation`.
- owner: `kernel.identity`.
- `kernel_code_authorized=true` only for P02.01.
- `business_feature_code_authorized=false`.
- P02.02-P02.10, P03+, business modules and AI/model/agent runtime remain unauthorized.

## P02 readiness evidence

Preparation PR #67 created the 10-package P02 sequence, entry/exit gates, transition checklist, P02 readiness/package validators and permanent completed-P01 prerequisite guards.

- readiness final exact head: `ec936b12ecc23f56c7dbd56d7bdb440d0c5a13b9`;
- canonical readiness run/job: `32632920772 / 97178312240` — PASS;
- preparation merge: `c6301ca4a5eec5dd62bcb75481d900e40ad968bd`;
- canonical runner policy: GitHub-hosted `ubuntu-24.04`, Linux/X64 only;
- repository Go quality and P01.01-P01.12 regressions remained green.

Initial readiness run `32632854943 / 97178139290` remains retained as diagnostic FAIL for a missing literal `session invalidation` marker in the newly introduced exit-gate document. The gate text was aligned to the canonical Master Plan vocabulary; no validator or quality gate was weakened.

## Active P02.01 boundary

P02.01 may establish only the canonical human principal/User identity foundation:

- stable UUIDv7 principal/User identifiers;
- deterministic User lifecycle states/transitions;
- classification-safe identity attributes;
- `kernel.identity` owned repository/persistence boundary where required;
- safe structured errors;
- focused positive/negative tests;
- applicable fresh/upgrade migration evidence.

P02.01 must not pull forward tenant lifecycle/membership authority, organization hierarchy, authentication/session behavior, RBAC or relationship policy, MFA/passkeys, service-account/API credential lifecycle, tenant settings, P03 module runtime, business domains, deployment authority or AI/model/agent runtime.

User remains authentication identity and is not the business `Person` model. Non-human identities are not fake users. Identity attributes do not create authorization.

## Retained P01 evidence

Final P01.12 implementation evidence remains PR #65, exact head `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`, canonical run/job `32629072886 / 97168916985`, and merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`. P01 regressions remain enforced by governance during P02.

## Exact next authorized action

This snapshot accompanies the P02.01 activation transition. After that activation PR merges, verify protected `main` and `STATE.json`, identify P02.01 as the sole authorized implementation scope, then **STOP the execution session**.

In the next governed session, create a fresh implementation branch from the exact protected-main SHA and implement only P02.01. Do not auto-advance to P02.02.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/roadmap/STATUS.md`, `docs/governance/P02_ENTRY_GATE.md`, `docs/governance/P02_EXIT_GATE.md`, `docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json`, `docs/roadmap/work-packages/P02.01.md`, `docs/governance/AI_EXECUTION_POLICY.md`, Change Control, Definition of Done, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and `docs/ai/handoffs/P02.01.md`.
