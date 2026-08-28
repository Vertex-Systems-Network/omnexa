# Omnexa Program Status

Last reconciled: **2026-08-28**

## Current position

- Program: **Kernel Program**
- Phase: **P03 — Module Runtime**
- Phase state: **ACTIVE**
- Current work package after this closure merges: **P03.06 — Capability Registry**
- P03 progress after this closure merges: **5 / 11 done**
- P03.01: **DONE**
- P03.02: **DONE**
- P03.03: **DONE**
- P03.04: **DONE**
- P03.05: **DONE — canonical implementation accepted**
- P03.06: **ACTIVE only after this closure transition merges**
- P03.07-P03.11: **PLANNED / LOCKED**
- P03 entry gate: **SATISFIED**
- P03 exit gate: **NOT SATISFIED**
- P02: **DONE — 10 / 10**
- P02 exit gate: **SATISFIED**
- P01: **DONE — 12 / 12**
- P01 exit gate: **SATISFIED**
- Foundation Architecture v1: **FROZEN**
- P00: **DONE — 10 / 10**
- Repository visibility: **PUBLIC**
- Main integration protection / Issue #3: **SATISFIED / CLOSED**
- Executable CI / Issue #14: **SATISFIED — GITHUB-HOSTED ONLY / ubuntu-24.04**
- Local/self-hosted governance runners: **PROHIBITED**
- Kernel implementation authority after this closure merges: **P03.06 ONLY**
- Business-feature implementation: **NOT AUTHORIZED**

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This status file mirrors the P03.05 completion / P03.06 activation closure candidate and grants no P03.06 implementation authority before the closure carrier itself passes exact-final-head canonical governance, merges to protected `main`, and protected-main state is re-read.

## P03.05 canonical completion boundary

P03.05 — Module Settings & Feature Flags completed through implementation PR #103 with exact final evidence:

- activation/base main: `f38daddf2a4b71290cadef02cbad8c1afaaca15f`
- final exact implementation head: `c52b48be1a82eb27670f03bdd4e1be4df6eb9f54`
- canonical Governance run/job: `33132237120 / 98724184966` — **PASS**
- implementation merge / protected-main base for this closure: `0c6b075c272aeac5a6e5f9d4210b1c5a30a040ce`
- canonical environment: GitHub-hosted `ubuntu-24.04`
- repository-pinned Go: `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.05_COMPLETION_2026-08-28.md`
- retained verifier: `scripts/verify_p03_05.sh`

The accepted exact-head run passed repository Go quality, all retained P01/P02/P03.01/P03.02/P03.03/P03.04 regressions and the dedicated P03.05 verifier.

Earlier failed/cancelled/stale candidates remain diagnostic history only. Governance #457 failed a test-call typecheck shape and was corrected without weakening production semantics or governance; only the exact final head/run above is completion authority.

## P03.05 delivered boundary retained

P03.05 preserved manifest schema v1/v2 and bound exact manifest-declared setting/feature-flag keys to existing typed `kernel.configuration.Definition` registrations. It introduced explicit global/scoped registration scope, reused the existing P02.09 scoped policy validator, retained trusted tenant/organization scope construction, preserved non-destructive configuration history across lifecycle states, rejected ownership/class/scope/policy collisions, and kept `kernel.configuration.NewRegistry` as type/default/duplicate authority.

Settings and feature flags create no permission, capability, tenant or database authority. No duplicate configuration subsystem was introduced.

## P03.06 activation boundary

After this separate closure/state-transition carrier passes exact-final-head canonical governance and merges to protected `main`, P03.06 becomes the sole active implementation package.

Authorized P03.06 scope is limited to `P03.06 — Capability Registry`:

- stable capability ID and major-version metadata;
- provider module/owner and consumer declarations;
- lifecycle-derived availability state;
- authorization / tenant-org scope requirement metadata and contract references;
- collision/version compatibility validation;
- lifecycle-aware registration/withdrawal availability;
- deterministic lookup suitable for later attribution/graph consumers;
- focused positive/adversarial tests and a dedicated P03.06 verifier.

Capability registration is metadata/availability only. It never grants invocation permission. The owning capability boundary remains responsible for validation, authorization, trusted tenant/org enforcement and required audit. Registry metadata must not expose private handlers, implementation objects, raw database tables, credentials, secrets or restricted values. Disabled/unhealthy/unavailable providers cannot be represented as active merely because a declaration exists.

P03.07-P03.11, P04+, business features, generic remote RPC/service mesh, workflow orchestration, product federation runtime, package acquisition/marketplace runtime, strategic X-program runtime and AI/model/agent runtime remain unauthorized.

## Retained prerequisite chain

P03.01-P03.05 remain DONE with canonical completion evidence and retained verifiers. P01.01-P01.12 and P02.01-P02.10 remain complete historical prerequisites. All retained P01/P02/P03.01-P03.05 regression verifiers remain mandatory during later P03.06 implementation.

## Protected integration / CI

`main` remains the protected PR-only integration authority. Canonical governance evidence is GitHub-hosted `ubuntu-24.04` only; local/self-hosted governance evidence is prohibited.

This closure carrier must remain current with protected `main` and its exact final head must pass canonical governance, including repository Go quality, all retained P01/P02/P03.01-P03.05 regression verification, P03 preparation/package validators and conversation-resolution requirements before merge.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. It does not broaden P03 implementation authority.

## Exact next work

1. Complete the atomic P03.05 closure / P03.06 activation reconciliation under GitHub issue #104 without P03.06 runtime code.
2. Keep the exact carrier within the approved governance/evidence/continuity boundary.
3. Require a fresh exact-final-head canonical Governance PASS; stale or diagnostic runs are not merge authority.
4. Merge the closure carrier only if it remains current with protected `main`, repository merge gates permit it and no unresolved review/conversation blocker exists.
5. Re-read protected `main`, `STATE.json`, this status, the package sequence and P03.06 handoff after closure merge and record the exact new main SHA.
6. Create a new separate P03.06 implementation branch from that exact post-closure SHA.
7. Implement only P03.06; keep P03.07+ locked and use a later separate governed closure for P03.06 completion/P03.07 activation.
