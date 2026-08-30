# P04.01 Closure / P04.02 Activation Transaction

Status: **CANDIDATE ON ISOLATED GOVERNANCE BRANCH — NOT AUTHORITY UNTIL MERGED + MAIN READ-BACK**

Exact transaction base: `16dfca22cd0430b8b2727e55628cc6fba381b5ef`

This carrier is governance/state/continuity only. It must not contain P04.02 runtime implementation.

## Accepted prerequisite

P04.01 implementation is already accepted on `main`:

- implementation exact head: `3c4ee6e79b76a6f042f3ee77b518943f3d7064e4`;
- promotion implementation PR: `#147`;
- implementation merge: `6b1a01c009e9a08d45613a00e9b63c0c272bf020`;
- promotion Governance: `#548` / run `33331934481`;
- completion evidence: `docs/roadmap/evidence/P04.01_COMPLETION_2026-08-31.md`;
- prepared P04.02 contract/handoff promotion: PR `#149`;
- preparation merge / this transaction base: `16dfca22cd0430b8b2727e55628cc6fba381b5ef`.

## Atomic candidate result after this carrier merges

The merged transaction must make all authoritative/continuity surfaces agree that:

- `current_phase=P04`;
- P04 progress is `1 / 10 done`;
- P04.01 is `done` with its accepted completion evidence;
- P04.02 is the **sole active** package;
- P04.03-P04.10 remain `planned / locked`;
- `kernel_code_authorized=true` only for P04.02 and only after this carrier passes exact-final-head Governance, merges, and protected `main` is re-read;
- `business_feature_code_authorized=false`;
- strategic X-program runtime remains unauthorized;
- AI/model/agent runtime remains unauthorized.

## P04.02 authority boundary

The later implementation branch may implement only the accepted provider-neutral P04.02 contract in `docs/roadmap/work-packages/P04.02.md`.

This transition does **not** authorize:

- Kafka, NATS, RabbitMQ, Redis Streams or another broker/provider selection;
- durable stream/consumer/checkpoint runtime;
- database migrations;
- outbox/inbox persistence;
- retry/backoff/DLQ/quarantine runtime;
- schema-registry runtime;
- background-job execution changes;
- business-domain event handlers;
- business-feature code;
- AI/model/agent runtime.

## Mandatory candidate files

Before this branch may open for governed review, the candidate must reconcile at least:

1. `docs/roadmap/STATE.json`;
2. `docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json`;
3. `AGENTS.md`;
4. `docs/roadmap/STATUS.md`;
5. `docs/ai/AI_STATE.yaml`;
6. `docs/ai/AI_CONTEXT.md` only if its current cursor/authority language is stale.

No authoritative file may claim P04.01 is still active while another claims P04.02 is active.

## Merge protocol

1. Re-read the exact branch diff and current protected `main` before review.
2. Require exact-final-head Omnexa Governance PASS for the complete candidate.
3. Require zero unresolved review threads and honest review provenance.
4. Promote unchanged state through a fresh promotion-specific carrier if the repository process requires it.
5. Merge with expected-head guard.
6. Re-read protected `main` and canonical state.
7. Only then create a **new** P04.02 runtime implementation branch from the exact post-transition `main` SHA.

A branch-only candidate, source CI PASS, this document, or the prepared P04.02 spec never grants runtime authority by itself.
