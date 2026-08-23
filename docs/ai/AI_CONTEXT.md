# Omnexa AI Project Context

Status: **continuity snapshot / subordinate to canonical governance**

This is the first file to read inside `docs/ai/`. It never overrides `AGENTS.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/roadmap/STATE.json`, accepted ADRs, architecture/security standards or live GitHub evidence.

## Project identity

**Omnexa** is a Composable Enterprise Business Operating System being built architecture-first as a strict modular monolith with service-ready boundaries. Development is sequential, gated and evidence-based: platform foundations are implemented and verified before later phases or business modules are authorized.

Core architecture retains one authoritative owner per capability/write model, governed cross-domain boundaries, mandatory isolation/authorization/audit/classification/versioning, non-authoritative infrastructure caches/telemetry, and future AI execution only through governed capabilities rather than direct database/object-store/payment/business-state authority.

## Verified P01.12 implementation result

Live GitHub evidence retained for the final P01 package:

- implementation PR: `#65` — merged;
- final exact source head: `2ee9a619f3bf828a4c38f8f3af7277fe8c7634f9`;
- canonical implementation run/job: `32629072886 / 97168916985` — **PASS**;
- runner: `GitHub Actions 1000013010`, GitHub-hosted Ubuntu 24.04.4 LTS / X64;
- runner image: `ubuntu-24.04 / 20260816.277.1`;
- Go: `1.26.7`;
- repository Go quality: PASS (`gofmt` 66 files, `golangci-lint v2.12.2` 0 issues, `govulncheck v1.7.0` no vulnerabilities);
- P01.01-P01.11 regressions: PASS;
- real developer `db migrate`: PASS;
- real `verify all`: PASS;
- P01.12 G0-G8: PASS;
- implementation squash merge: `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

Initial run `32628865023 / 97168401358` remains retained as FAIL history for one gofmt alignment, three gosec G204 findings and one staticcheck ST1005 finding. They were fixed directly without linter exclusions, governance weakening or scope expansion.

## P01.12 delivered boundary

P01.12 established only `kernel.developer`: deterministic help/version, safe structured health diagnostics, guarded non-production migration, fail-closed governed verification orchestration, exact executable-plus-argument allowlisting, verification environment isolation, focused positive/negative/race tests and the P01 fresh-install exit proof.

It did **not** build P02 identity/tenancy/role administration, P03 module runtime administration, P04+ domain/event/workflow commands, deployment/Kubernetes authority, business modules or AI/model/agent runtime.

## P01 completion / exit snapshot

Closure branch: `chore/p01-12-close-p01-exit`  
Base: P01.12 implementation merge `eeebaf5ae3817588b014ddf4c9911bca52c97ed7`.

The closure reconciles exactly:

- P01.01-P01.12 `done`;
- P01 progress `12 / 12 done`;
- P01 exit gate `SATISFIED`;
- active P01 work package: none;
- P02 remains `planned / not active`;
- `kernel_code_authorized=false`;
- `business_feature_code_authorized=false`;
- no P02 runtime implementation in the closure transition.

Canonical final package evidence is `docs/roadmap/evidence/P01.12_COMPLETION_2026-08-23.md`. Phase-exit evidence is `docs/governance/P01_EXIT_GATE.md`.

The terminal state is intentionally between phases: P01 is complete, but P02 has not been activated implicitly. A future session must re-read protected `main`, `STATE.json` and latest CI before doing anything material.

## Exact next authorized action after closure

After the P01 completion closure merges, **STOP that execution session**.

In a new governed execution session, only P02 specification/readiness preparation and an explicit P02 activation transition may proceed under the non-implementation allowance in `STATE.json`. P02 implementation must not begin until that later transition merges and canonical state explicitly authorizes it.

Do not automatically start P02/P03/business features or AI/model/agent implementation.

## Authority and references

Mandatory sources remain `AGENTS.md`, `docs/roadmap/STATE.json`, `docs/governance/P01_EXIT_GATE.md`, `docs/governance/AI_EXECUTION_POLICY.md`, `docs/governance/CHANGE_CONTROL.md`, `docs/governance/DEFINITION_OF_DONE.md`, accepted ADRs, architecture/security/quality standards, `docs/ai/AI_STATE.yaml`, `docs/ai/AI_EXECUTION_PROTOCOL.md` and the relevant handoff.
