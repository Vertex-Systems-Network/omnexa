# Omnexa Program Status

Last reconciled: **2026-08-28**

## Current position

- Program: **Kernel Program**
- Phase: **P03 — Module Runtime**
- Phase state: **ACTIVE**
- Current work package after this closure merges: **P03.07 — Permission Registration**
- P03 progress after this closure merges: **6 / 11 done**
- P03.01: **DONE**
- P03.02: **DONE**
- P03.03: **DONE**
- P03.04: **DONE**
- P03.05: **DONE**
- P03.06: **DONE — canonical implementation accepted**
- P03.07: **ACTIVE only after this closure transition merges**
- P03.08-P03.11: **PLANNED / LOCKED**
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
- Kernel implementation authority after this closure merges: **P03.07 ONLY**
- Business-feature implementation: **NOT AUTHORIZED**

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This status file mirrors the P03.06 completion / P03.07 activation closure candidate and grants no P03.07 implementation authority before the closure carrier itself passes exact-final-head canonical governance, merges to protected `main`, and protected-main state is re-read.

## P03.06 canonical completion boundary

P03.06 — Capability Registry completed through implementation PR #107 with exact final evidence:

- activation/base main: `c3dedb6e6f3e602bec0410a27616e4a7f5cea218`
- initial implementation head: `28a6343ad10924869d02d666f7b625121f7ecb7d`
- final exact implementation head: `c895f44a1383d1c1d9c5fd23c95d7864810353c3`
- canonical Governance run/job: `33181421854 / 98883286556` — **PASS**
- implementation merge / protected-main base for this closure: `13dbe8a393c20cabeb8aac60d073a6c66775efd3`
- canonical environment: GitHub-hosted `ubuntu-24.04`
- repository-pinned Go: `1.26.7`
- completion evidence: `docs/roadmap/evidence/P03.06_COMPLETION_2026-08-28.md`
- retained verifier: `scripts/verify_p03_06.sh`

The accepted exact-head run passed repository Go quality, all retained P01/P02/P03.01-P03.05 regressions and the dedicated P03.06 verifier.

Earlier failed/stale candidates remain diagnostic history only. Governance #467 / run `33180840326`, job `98881325283` failed repository Go quality on one variable-shadow finding and was corrected without behavior or gate weakening; only the exact final head/run above is completion authority.

## P03.06 delivered boundary retained

P03.06 preserved manifest schema v1/v2 and retained validated provided/consumed capability declarations in the package-private discovery snapshot. It introduced deterministic stable capability + major-version identity, provider module/owner and consumer metadata, lifecycle-derived availability, descriptive authorization/scope/contract references, fail-closed duplicate/owner/version validation and deterministic lookup without creating invocation authority.

Capability registration/lookup grants no permission, accepts no raw tenant/org scope, exposes no private handler/table/secret, creates no cross-module write authority and does not duplicate `kernel.authorization`. Only lifecycle-enabled providers are active while unavailable lifecycle states retain historical identity without reporting active.

## P03.07 activation boundary

After this separate closure/state-transition carrier passes exact-final-head canonical governance and merges to protected `main`, P03.07 becomes the sole active implementation package.

Authorized P03.07 scope is limited to `P03.07 — Permission Registration`:

- stable permission name/owner/module metadata;
- declaration collision/namespace validation;
- optional capability association as descriptive metadata;
- lifecycle-derived permission availability;
- preservation of policy/role references/history across non-destructive disable/re-enable;
- fail-closed unknown/unavailable permission behavior;
- auditability of material registration/lifecycle changes where required;
- focused positive/adversarial tests and a dedicated P03.07 verifier.

Permission registration is declaration metadata only; it never creates grants. `kernel.authorization` remains deny-by-default enforcement/policy authority. Role names cannot create bypass authority, tenant/org scope comes only from trusted P02 context, registration cannot mutate role grants implicitly or widen principal scope, and disabled/unavailable module permissions cannot authorize behavior.

P03.08-P03.11, P04+, business features, role-editor/admin UI, a duplicate authorization engine, super-admin bypass, entitlements/licensing runtime, generic RPC/service mesh, workflow orchestration, product federation runtime, package acquisition/marketplace runtime, strategic X-program runtime and AI/model/agent runtime remain unauthorized.

## Retained prerequisite chain

P03.01-P03.06 remain DONE with canonical completion evidence and retained verifiers. P01.01-P01.12 and P02.01-P02.10 remain complete historical prerequisites. All retained P01/P02/P03.01-P03.06 regression verifiers remain mandatory during later P03.07 implementation.

## Protected integration / CI

`main` remains the protected PR-only integration authority. Canonical governance evidence is GitHub-hosted `ubuntu-24.04` only; local/self-hosted governance evidence is prohibited.

This closure carrier must remain current with protected `main` and its exact final head must pass canonical governance, including repository Go quality, all retained P01/P02/P03.01-P03.06 regression verification, P03 preparation/package validators and conversation-resolution requirements before merge.

## Issue #4 — external distribution gate

Issue #4 remains open for licensing/IP/trademark and public-launch decisions. Repository visibility is public and current `LICENSE` remains GPLv3. It does not broaden P03 implementation authority.

## Exact next work

1. Complete the atomic P03.06 closure / P03.07 activation reconciliation under GitHub issue #108 without P03.07 runtime code.
2. Keep the exact carrier within the approved governance/evidence/continuity boundary.
3. Require a fresh exact-final-head canonical Governance PASS; stale or diagnostic runs are not merge authority.
4. Merge the closure carrier only if it remains current with protected `main`, repository merge gates permit it and no unresolved review/conversation blocker exists.
5. Re-read protected `main`, `STATE.json`, this status, the package sequence and P03.07 handoff after closure merge and record the exact new main SHA.
6. Create a new separate P03.07 implementation branch from that exact post-closure SHA.
7. Implement only P03.07; keep P03.08+ locked and use a later separate governed closure for P03.07 completion/P03.08 activation.
