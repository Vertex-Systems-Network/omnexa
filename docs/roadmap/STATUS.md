# Omnexa Program Status

Last reconciled: **2026-08-28**

## Current position

- Program: **Kernel Program**
- Phase: **P03 — Module Runtime**
- Phase state: **ACTIVE**
- Current work package after this closure merges: **P03.08 — UI Contribution Registry Contract**
- P03 progress after this closure merges: **7 / 11 done**
- P03.01: **DONE**
- P03.02: **DONE**
- P03.03: **DONE**
- P03.04: **DONE**
- P03.05: **DONE**
- P03.06: **DONE**
- P03.07: **DONE — canonical implementation accepted**
- P03.08: **ACTIVE only after this closure transition merges**
- P03.09-P03.11: **PLANNED / LOCKED**
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
- Kernel implementation authority after this closure merges: **P03.08 ONLY**
- Business-feature implementation: **NOT AUTHORIZED**

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This status file mirrors the P03.07 completion / P03.08 activation closure candidate and grants no P03.08 implementation authority before the closure carrier itself passes exact-final-head canonical governance, merges to protected `main`, and protected-main state is re-read.

## P03.07 canonical completion boundary

P03.07 — Permission Registration completed through implementation PR #111 with exact final evidence:

- implementation issue: #110 — completed
- final exact implementation head: `28e36b3ac3183f28ec500f1e70b1fefe02c0c325`
- canonical Governance run/job: `33195104185 / 98930123416` — **PASS**
- implementation merge / protected-main base for this closure: `66f8b4cc630f6cd865e440a62478df365e042a31`
- canonical environment: GitHub-hosted `ubuntu-24.04`, Linux/X64
- repository-pinned Go: `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.07_COMPLETION_2026-08-28.md`
- retained verifier: `scripts/verify_p03_07.sh`

The accepted exact-head run passed repository Go quality, all retained P01/P02/P03.01-P03.06 regressions and the dedicated P03.07 verifier.

Diagnostic-only history remains unchanged: Governance #474 / `33192567020 / 98921494281` failed repository Go quality; Governance #476 / `33194438411 / 98927852853` failed only a brittle dedicated-verifier source marker after all prior regressions passed. Both were corrected without semantic or gate weakening. Only Governance #477 / `33195104185 / 98930123416` is completion authority.

## P03.07 delivered boundary retained

P03.07 added stable module-declared permission registration while preserving `kernel.authorization` as the sole deny-by-default enforcement/policy authority. Invalid/reserved/colliding ownership fails closed; lifecycle availability is explicit; unknown/unavailable permissions deny; registration cannot create or widen grants; role names create no bypass; trusted tenant/organization scope remains governed by P02; optional capability association remains descriptive and non-invoking; role/history/audit references are preserved across lifecycle changes.

The additive migration `kernel/migrations/kernel.authorization/4_add_module_permission_registration.sql` does not authorize destructive cleanup. Historical identity and references remain protected.

## P03.08 activation boundary

After this separate closure/state-transition carrier passes exact-final-head canonical governance and merges to protected `main`, P03.08 becomes the sole active implementation package.

Authorized P03.08 scope is limited to `P03.08 — UI Contribution Registry Contract`:

- stable contribution ID/module owner/slot metadata;
- bounded declarative navigation/page/widget/settings/builder-slot contribution metadata;
- independent permission requirement metadata without grants;
- lifecycle/module availability conditions;
- feature-flag conditions without authorization authority;
- optional-dependency conditions and selective fallback/degradation metadata;
- deterministic duplicate/conflict/unknown-slot validation;
- versioned metadata suitable for future Experience/Portal/Product Federation consumers without implementing those runtimes;
- focused positive/adversarial tests and a dedicated P03.08 verifier.

UI visibility never substitutes for backend authorization. Contribution metadata cannot grant permission, widen tenant/org scope, embed secrets or unrestricted executable code, expose private handlers as execution shortcuts, or create cross-module database-write authority.

P03.09-P03.11, P04+, business features, rendering/CMS/Experience Builder runtime, portal runtime, product federation runtime, generic RPC/service mesh, workflow orchestration, package acquisition/marketplace runtime, strategic X-program runtime and AI/model/agent runtime remain unauthorized.

## Retained prerequisite chain

P03.01-P03.07 remain DONE with canonical completion evidence and retained verifiers. P01.01-P01.12 and P02.01-P02.10 remain complete historical prerequisites. All retained P01/P02/P03.01-P03.07 regression verifiers remain mandatory during later P03.08 implementation.

## Protected integration / CI

`main` remains the protected PR-only integration authority. Canonical governance evidence is GitHub-hosted `ubuntu-24.04` only; local/self-hosted governance evidence is prohibited.

This closure carrier must remain current with protected `main` and its exact final head must pass canonical governance, including repository Go quality, all retained P01/P02/P03.01-P03.07 regression verification, P03 preparation/package validators and conversation-resolution requirements before merge.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. It does not broaden P03 implementation authority.

## Exact next work

1. Complete the atomic P03.07 closure / P03.08 activation reconciliation under GitHub issue #112 without P03.08 runtime code.
2. Keep the exact carrier within the approved governance/evidence/continuity boundary.
3. Require a fresh exact-final-head canonical Governance PASS; stale or diagnostic runs are not merge authority.
4. Merge the closure carrier only if it remains current with protected `main`, repository merge gates permit it and no unresolved review/conversation blocker exists.
5. Re-read protected `main`, `STATE.json`, this status, the package sequence and P03.08 handoff after closure merge and record the exact new main SHA.
6. Identify a new separate P03.08 implementation branch from that exact post-closure SHA as the next authorized action and stop without implementing P03.08 in this closure transaction.
7. Keep P03.09+ locked and use a later separate governed closure for P03.08 completion/P03.09 activation.
