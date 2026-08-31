# Omnexa Program Status

Last reconciled: **2026-08-31 P04.02 closure / P04.03 activation candidate**

## Current authoritative position vs candidate

Protected `main@2d454a87e03f404f081b6a87f216d0cfa8c7608d` remains authoritative until this complete transaction passes exact-final-head Governance, is promoted unchanged, merges through the protected path, and protected main plus canonical state are re-read.

Authoritative protected-main state before this transaction:

- P04: **ACTIVE — 1 / 10 done**;
- P04.01: **DONE**;
- P04.02: **SOLE ACTIVE PACKAGE** with accepted implementation/completion evidence;
- P04.03: **PREPARED / PLANNED / LOCKED**;
- P04.04-P04.10: **PLANNED / LOCKED**;
- `kernel_code_authorized=true` for P04.02 only;
- `business_feature_code_authorized=false`.

Candidate atomic result on this isolated governance branch:

- Program: **Kernel Program**;
- Phase: **P04 — Data, Jobs & Event Fabric**;
- Phase state: **ACTIVE**;
- P04 progress: **2 / 10 done**;
- P04.01: **DONE — retained accepted evidence**;
- P04.02: **DONE — accepted implementation/completion evidence retained**;
- Current work package after governed merge/read-back: **P04.03 — Durable Stream/Consumer Baseline & Checkpoint Model**;
- P04.03: **SOLE ACTIVE PACKAGE after governed merge/read-back**;
- P04.04-P04.10: **PLANNED / LOCKED**;
- P04 entry gate: **SATISFIED**;
- P03: **DONE — 11 / 11**; exit **SATISFIED / HISTORICAL**;
- P02: **DONE — 10 / 10**; exit **SATISFIED**;
- P01: **DONE — 12 / 12**; exit **SATISFIED**;
- Foundation Architecture v1: **FROZEN**;
- Main integration protection / Issue #3: **SATISFIED / CLOSED**;
- Canonical CI: **GITHUB-HOSTED ONLY / ubuntu-24.04**;
- Local/self-hosted governance runners: **PROHIBITED**;
- Kernel implementation authority after governed merge/read-back: **P04.03 ONLY**;
- Business-feature implementation: **NOT AUTHORIZED**.

This carrier contains governance/state/continuity and validator hardening only. It contains no P04.03 runtime implementation.

## Accepted P04.02 evidence chain

P04.02 implementation/completion facts are already accepted on protected main and are not created by this transition:

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

## Accepted P04.03 preparation chain

P04.03 contract/handoff preparation is also already landed but remains non-authorizing before this transition:

- preparation source PR: `#158`;
- preparation exact head: `a20b37cbe0401626cd29e51a08efc399ee2d72e2`;
- source Governance: `#568` / run `33364934121`;
- unchanged preparation promotion PR: `#159`;
- promotion Governance: `#569` / run `33368684906` / job `99414689966`;
- preparation merge/read-back: `2d454a87e03f404f081b6a87f216d0cfa8c7608d`;
- contract: `docs/roadmap/work-packages/P04.03.md`;
- handoff: `docs/ai/handoffs/P04.03.md`.

Preparation alone does not authorize durable consumer/checkpoint runtime.

## P04.03 candidate implementation boundary

Only after this closure/activation transaction itself passes exact-final-head Governance, is promoted unchanged, merges, and protected main is re-read may a **new separate implementation branch** implement P04.03.

Authorized P04.03 scope is limited to the accepted provider-neutral durability/checkpoint contract:

- stable durable consumer identity compatible with P04.02 owner/consumer identity;
- explicit stream/partition-equivalent checkpoint scope;
- monotonic checkpoint advancement within one scope;
- deterministic stale/regressive/conflicting checkpoint rejection;
- restart/resume from the last accepted checkpoint;
- cancellation/failure behavior that cannot advance past unacknowledged work;
- fail-closed owner/tenant/scope collision handling;
- explicit at-least-once-compatible duplicate delivery around crash/restart;
- preservation of P04.01 envelope metadata and P04.02 ownership laws;
- focused positive/adversarial tests and package verification.

Explicitly unauthorized under P04.03:

- Kafka, NATS, RabbitMQ, Redis Streams or another broker/provider selection;
- exactly-once delivery/business-mutation claims;
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

## Sequential activation validator hardening

`scripts/validate_p04_activation.py` is hardened in this transaction so advancement no longer relies on a P04.01-specific exception. For every completed predecessor it now requires:

- `state=done` in the package sequence;
- retained accepted spec and completion-evidence file;
- matching canonical `governance_tracking.p04_NN_completion` entry with `state=PASS`;
- matching evidence reference;
- exact implementation head, implementation merge, workflow run and job tracking;
- GitHub-hosted evidence on `ubuntu-24.04`;
- state/sequence mirrors and prepared-spec count consistent with strict sequential activation.

This makes P04.02 evidence a fail-closed prerequisite for P04.03 and generalizes the invariant for later P04 packages.

## State/continuity rule

On this candidate branch `STATE.json`, `P04_PACKAGE_SEQUENCE.json`, `AGENTS.md`, this status file and `docs/ai/` must converge on the same proposed atomic state:

```text
P04 ACTIVE — 2 / 10 done after governed merge/read-back
P04.01 DONE
P04.02 DONE
P04.03 SOLE ACTIVE PACKAGE after governed merge/read-back
P04.04-P04.10 PLANNED / LOCKED
kernel_code_authorized=true — P04.03 only after governed merge/read-back
business_feature_code_authorized=false
```

Any contradictory cursor/authority state is a stop-the-line reconciliation defect.

## Protected integration / CI

This source carrier is not merge authority. It must obtain fresh exact-final-head Omnexa Governance, substantive exact-head review, zero unresolved review threads, then be promoted unchanged through a fresh promotion-specific carrier with its own Governance pass. Merge must use the expected-head guard while current with protected main, followed by protected-main/state read-back.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions and grants no kernel/business implementation authority.

## Exact next work

1. Finish this P04.02 closure / P04.03 activation governance transaction without runtime source.
2. Require exact-final-head Omnexa Governance and substantive review.
3. Promote the unchanged reviewed head through a fresh promotion carrier and require promotion-specific Governance.
4. Merge only while current with protected main and with zero unresolved conversations.
5. Re-read protected main, canonical state, package sequence and continuity; confirm P04.03 is sole active and P04.04+ remain locked.
6. Only then create a new P04.03 implementation branch from that exact protected-main SHA. Do not implement P04.03 on this activation carrier.
