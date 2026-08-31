# Omnexa Program Status

Last reconciled: **2026-08-31 post P04.02 closure / P04.03 activation read-back**

## Current authoritative position

- Program: **Kernel Program**
- Phase: **P04 — Data, Jobs & Event Fabric**
- Phase state: **ACTIVE**
- P04 progress: **2 / 10 done**
- P04.01: **DONE — retained accepted evidence**
- P04.02: **DONE — accepted implementation/completion evidence retained**
- Current work package: **P04.03 — Durable Stream/Consumer Baseline & Checkpoint Model**
- P04.03: **SOLE ACTIVE PACKAGE**
- P04.04-P04.10: **PLANNED / LOCKED**
- P04 entry gate: **SATISFIED**
- P03: **DONE — 11 / 11**; exit **SATISFIED / HISTORICAL**
- P02: **DONE — 10 / 10**; exit **SATISFIED**
- P01: **DONE — 12 / 12**; exit **SATISFIED**
- Foundation Architecture v1: **FROZEN**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Canonical CI: **GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation authority: **P04.03 ONLY**
- Business-feature implementation: **NOT AUTHORIZED**
- Strategic X-program runtime: **NOT AUTHORIZED**
- AI/model/agent runtime: **NOT AUTHORIZED**

Accepted P04.03 activation read-back on protected `main` is `d74375cd0a2952ab8622117089e5eb43043e6e78`. Canonical `docs/roadmap/STATE.json` on that SHA confirms `current_phase=P04`, `current_work_package=P04.03`, `done_work_packages=2`, `kernel_code_authorized=true` for the bounded P04.03 contract, and `business_feature_code_authorized=false`.

This status file is a subordinate continuity snapshot. It never overrides `AGENTS.md`, canonical `STATE.json`, accepted ADRs, governance policy or live GitHub evidence.

## Accepted P04.02 evidence chain

P04.02 implementation/completion remains immutable predecessor evidence:

- P04.01 closure / P04.02 activation merge: `226f5236d80c806a7b07d65e7870981de320c05c`;
- source implementation PR: `#155`;
- promotion implementation PR: `#156`;
- final exact implementation head: `5ec1de746eebc8734f86ec3aa2f311daae0dc18a`;
- source Governance: `#563` / run `33350757489`;
- promotion Governance: `#564` / run `33351220746` / job `99364850230`;
- accepted implementation merge: `1b378a2f44c6e3cba87b936e39e05f1a18da94cc`;
- completion evidence: `docs/roadmap/evidence/P04.02_COMPLETION_2026-08-31.md`;
- completion-evidence carrier PR: `#157`;
- completion-evidence merge/read-back: `8334bed24e793e200540bf953f7249f11242badf`.

Historical P01.01-P01.12, P02.01-P02.10, P03.01-P03.11 and P04.01 evidence remains immutable and all retained regression verifiers remain mandatory.

## Accepted P04.03 preparation and activation chain

P04.03 contract/handoff preparation was accepted before activation:

- preparation source PR: `#158`;
- preparation exact head: `a20b37cbe0401626cd29e51a08efc399ee2d72e2`;
- source Governance: `#568` / run `33364934121`;
- unchanged preparation promotion PR: `#159`;
- preparation promotion Governance: `#569` / run `33368684906` / job `99414689966`;
- preparation merge/read-back: `2d454a87e03f404f081b6a87f216d0cfa8c7608d`;
- accepted contract: `docs/roadmap/work-packages/P04.03.md`;
- handoff: `docs/ai/handoffs/P04.03.md`.

The separate P04.02 closure / P04.03 activation transaction is now accepted:

- activation source PR: `#160`;
- exact source/promotion head: `452858a3ab9bfa827697105bd5168cf660bd62ba`;
- source Governance: `#571` / run `33370681216` / job `99420856278` — **PASS**;
- unchanged promotion PR: `#161`;
- promotion Governance: `#572` / run `33371203708` / job `99422521576` — **PASS**;
- zero unresolved review threads before merge;
- explicit source and promotion `SELF REVIEW` recorded without fabricating independent review;
- expected-head guarded promotion merge / protected-main activation read-back: `d74375cd0a2952ab8622117089e5eb43043e6e78`.

Source and promotion evidence are retained separately. The activation merge does not rewrite P04.02 implementation evidence and does not authorize P04.04+.

## P04.03 authoritative implementation boundary

Owner: `kernel.events`.

A **new separate implementation branch from exact current protected main** may implement only the accepted `docs/roadmap/work-packages/P04.03.md` contract after re-verifying the live cursor and this post-activation continuity state.

Authorized P04.03 scope is limited to:

- stable durable consumer identity compatible with P04.02 owner/consumer identity;
- explicit stream/shard/partition-equivalent checkpoint scope;
- monotonic checkpoint advancement within one scope;
- deterministic stale/regressive/conflicting checkpoint rejection;
- restart/resume from the last accepted checkpoint;
- no advancement past failed/cancelled unacknowledged work;
- fail-closed owner/tenant/scope collision handling;
- at-least-once-compatible duplicate delivery around crash/restart;
- preservation of P04.01 envelope metadata and P04.02 ownership laws;
- focused positive/adversarial tests and dedicated package verification/evidence.

Explicitly unauthorized under P04.03:

- Kafka, NATS, RabbitMQ, Redis Streams or another concrete broker/provider selection;
- end-to-end exactly-once delivery or exactly-once business-mutation claims;
- automatic migration authority merely because P04.03 is active;
- transactional outbox (`P04.04`);
- inbox/deduplication persistence (`P04.05`);
- retry/backoff/DLQ/quarantine runtime (`P04.06`);
- schema-registry runtime (`P04.07`);
- background-job execution changes (`P04.08`);
- business-domain handlers or cross-module private writes;
- P04.04+ implementation;
- business-feature code;
- strategic X-program runtime;
- AI/model/agent runtime.

If a later P04.03 implementation genuinely requires persistence, the exact persistence/migration path must be justified inside the accepted P04.03 contract budget with rollback, tenant-isolation and compatibility evidence; activation itself grants no blanket migration authority.

## Sequential activation invariant

`scripts/validate_p04_activation.py` now generically requires every completed P04 predecessor to retain its accepted specification, completion evidence, matching `governance_tracking.p04_NN_completion` PASS record, exact implementation head/merge/run/job and GitHub-hosted `ubuntu-24.04` evidence. P04.02 is therefore a fail-closed predecessor for active P04.03, and future transitions must preserve the same generic rule rather than add package-specific exceptions.

## State / continuity rule

Canonical and subordinate continuity must agree on the effective cursor:

```text
P04 ACTIVE — 2 / 10 done
P04.01 DONE
P04.02 DONE
P04.03 SOLE ACTIVE PACKAGE
P04.04-P04.10 PLANNED / LOCKED
kernel_code_authorized=true — P04.03 only
business_feature_code_authorized=false
```

Any future contradictory cursor/authority state is a stop-the-line reconciliation defect. Historical transaction documents may describe pre-activation conditions, but their final acceptance status must not be mistaken for current implementation authority.

## Protected integration / CI

`main` remains governed through PR-only integration. Canonical acceptance evidence is GitHub-hosted `ubuntu-24.04` only. Source Governance #571 and promotion Governance #572 are the accepted P04.02→P04.03 transition evidence; future P04.03 implementation must obtain its own exact-head Governance and later separate completion/next-package transition evidence.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions and grants no kernel/business implementation authority.

## Exact next work

1. Before any P04.03 runtime mutation, verify protected `main`, canonical `STATE.json`, this status, AI continuity and the active handoff all agree that P04.03 is sole active and P04.04+ remain locked.
2. If this post-activation continuity correction is not yet on protected `main`, finish its exact-head Governance, review, governed promotion/merge and protected-main read-back first.
3. Create a **new isolated P04.03 implementation branch** from that exact verified protected-main SHA; never reuse an activation, promotion or continuity branch for runtime code.
4. Implement only the accepted P04.03 contract with focused crash/restart, monotonicity, duplicate, cancellation, malformed-state and tenant/scope isolation evidence.
5. Do not auto-advance to P04.04 after implementation; completion evidence and any P04.04 preparation/activation require later separate governed transactions.
