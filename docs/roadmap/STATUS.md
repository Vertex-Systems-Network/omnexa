# Omnexa Roadmap Status

Last reconciled: 2026-09-05 — **P04.05 ACTIVE after accepted P04.04 closure / P04.05 activation merge/read-back**

## Authoritative checkpoint

Current protected main at continuity start: `35b438a77838518400758e8b877170918cc3278f`.

Accepted P04.05 activation read-back: `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a`.

The later `35b438a77838518400758e8b877170918cc3278f` checkpoint contains only the separately governed checkout-v7 maintenance merge #199 after activation; it does not change the P04.05 runtime/state boundary.

Canonical state:

- Foundation Architecture v1 is FROZEN.
- P00 is DONE — 10 / 10.
- P01 is DONE — 12 / 12; exit SATISFIED.
- P02 is DONE — 10 / 10; exit SATISFIED.
- P03 is DONE — 11 / 11; exit SATISFIED.
- P04 is ACTIVE — **4 / 10 done**.
- P04.01-P04.04 are DONE with retained accepted evidence.
- P04.05 is the **sole ACTIVE package**.
- P04.06-P04.10 remain planned/locked.
- `kernel_code_authorized=true` only for P04.05.
- `business_feature_code_authorized=false`.
- strategic X-program and AI/model/agent runtime remain unauthorized.
- canonical CI is GitHub-hosted `ubuntu-24.04` only.

`docs/roadmap/STATE.json` remains the canonical machine-readable cursor. This continuity file is subordinate and grants no authority beyond that state.

## Accepted P04.04 implementation/completion

P04.04 is accepted through its bounded multi-agent chain:

- T01 core: source #184 / promotion #186;
- T02 PostgreSQL persistence/migration: source #185 / promotion #187;
- T03 reliability/concurrency: source #188 / promotion #189;
- T02 regression-fixture repair: source #191 / promotion #192;
- T04 Supervisor verifier: source #190 / promotion #193;
- final exact implementation head `ef09b878577d25a4a1186cb8fe84205b08a24851`;
- promotion Governance `33810095507 / 100829646792` — PASS;
- implementation merge/read-back `66c072b5caf42ceecb88d30cd1a1ee4e910322e6`;
- completion evidence `docs/roadmap/evidence/P04.04_COMPLETION_2026-09-04.md`;
- evidence carrier #194 merge/read-back `4445c21f1e6b03e84859d31ce7b32169b9c4cccc`.

The accepted implementation preserves provider-neutral, at-least-once-compatible producer semantics. It does not select a broker, grant business authority, or claim end-to-end exactly once.

## Accepted P04.05 preparation

- preparation source #195;
- unchanged preparation promotion #196;
- exact preparation head `211fea2077d7a1bf94be48f32f047b27273a4515`;
- merge/read-back `fa53b01cd92c8e0dd59026abff06f5f95f642d2d`;
- contract `docs/roadmap/work-packages/P04.05.md`;
- handoff `docs/ai/handoffs/P04.05.md`.

## Accepted P04.04 closure / P04.05 activation

- source PR #197;
- exact source/promotion head `6907253d375125a7ff096fb434c3433dbc17b331`;
- source Governance `33819597433 / 100859169990` — PASS;
- source review: SELF REVIEW only; independent approval not claimed;
- source unresolved review threads: 0;
- unchanged promotion PR #198;
- promotion Governance `33820312475 / 100861375770` — PASS;
- promotion review: SELF REVIEW only; independent approval not claimed;
- promotion unresolved review threads: 0;
- expected-head guarded merge/read-back `3402cf7a8b2b1370aca99543d47a33dee3dc0c5a`.

The accepted activation transaction changed governance/state/continuity only. It introduced no P04.05 runtime source, migration, retry/DLQ, schema-registry, background-job, provider, business or AI/model/agent runtime.

## P04.05 runtime boundary

Owner: `kernel.events`.

After this post-activation continuity correction is itself governed, promoted unchanged, merged and read back, a fresh implementation wave may implement only:

- stable processing identity bound to canonical EventID plus consumer/owner/tenant/route scope;
- no global cross-consumer EventID lock;
- inbox completion and the protected local PostgreSQL mutation in the same local transaction;
- no completion before the protected mutation commits;
- deterministic already-applied outcome for committed redelivery;
- restart/checkpoint-gap duplicate safety;
- concurrency, conflicting identity/content reuse and tenant/owner/consumer rebinding that fail closed;
- checkpoint progress and inbox completion as separate facts;
- payload-safe failures and focused verification.

Still unauthorized:

- P04.06 retry/backoff/terminal failure/DLQ/quarantine;
- P04.07 schema registry;
- P04.08 background-job runtime;
- broad P04.09 operator recovery;
- concrete broker/provider selection;
- external-side-effect or end-to-end exactly-once claims;
- business features, strategic X runtime or AI/model/agent runtime.

## Persistence / migration preflight

Production P04.05 semantics require durable local PostgreSQL inbox persistence, but activation added no migration and grants no blanket schema authority. Before the first schema mutation on the later fresh implementation branch, record the exact `kernel.events` owner/path/version/table/index/constraint/data budget and prove retained P01 fresh-install, upgrade, ledger identity, forward-recovery and owner/tenant-isolation rules. Do not pull P04.06+ schema semantics forward.

## Exact next work

1. Govern this **post-activation continuity-only** source carrier without changing canonical `STATE.json`, package sequence, validators or runtime source.
2. Inspect exact diff/review state, promote the unchanged source head and require fresh promotion Governance.
3. Merge with expected-head protection while current with protected main and re-read all continuity surfaces.
4. Confirm P04 remains 4 / 10, P04.05 sole ACTIVE, P04.06+ locked and business/AI runtime unauthorized.
5. Only then create a fresh P04.05 implementation wave from exact protected main and record the runtime/migration budget before source mutation.
6. Do not auto-advance to P04.06.
