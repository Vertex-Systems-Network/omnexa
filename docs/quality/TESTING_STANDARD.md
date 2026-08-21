# Omnexa Testing Standard

Status: **Canonical v1**  
Work package: **P00.07**

Testing proves contracts and invariants, not merely code coverage. Every later implementation package must map changed risk to the smallest sufficient set of deterministic tests and must include negative evidence for security, tenancy and lifecycle boundaries where relevant.

## 1. Test layers

Omnexa recognizes these canonical layers:

1. **Static/structural checks** — formatting, lint, type checks, generated-contract drift, forbidden dependency rules, schema validation and repository governance.
2. **Unit tests** — pure/domain behavior with no network or external infrastructure.
3. **Component tests** — one deployable/module boundary with real internal dependencies where practical and controlled external doubles.
4. **Contract tests** — HTTP/OpenAPI, event/schema, module manifest, capability and connector compatibility.
5. **Integration tests** — real database/cache/broker/object-store boundaries and provider adapters where required.
6. **Migration tests** — fresh install, upgrade, rollback/forward-fix assumptions and deterministic seeds/reference data.
7. **Security/negative tests** — authentication, authorization, tenant/org isolation, sensitive-field disclosure, webhook authenticity, SSRF, replay/idempotency and abuse boundaries.
8. **Module lifecycle tests** — install, enable, disable, re-enable, upgrade, optional dependency degradation, export/detach/purge semantics.
9. **End-to-end tests** — business-critical cross-capability flows through public interfaces; use sparingly, never as a substitute for lower layers.
10. **Performance/resilience tests** — latency/throughput/resource budgets, retries, backpressure, failure isolation and recovery for defined SLO-sensitive paths.
11. **Compatibility tests** — supported API/event/schema/module/client versions and deprecation windows.
12. **Disaster-recovery/rehearsal tests** — backup/restore, tenant-safe recovery and operational runbooks where the phase requires them.

## 2. Risk-based rule

A changed capability must test its highest-risk failure modes. Examples:

- tenant-owned data -> positive same-tenant + negative cross-tenant test;
- privileged mutation -> allow + deny + audit test;
- retriable mutation -> idempotency duplicate/retry test;
- event consumer -> duplicate delivery + replay safety test;
- schema change -> fresh + upgrade migration test;
- optional module -> disable/re-enable/dependency degradation test;
- money -> precision/rounding/boundary tests;
- civil-time logic -> DST/timezone/date-boundary tests;
- localization -> fallback/RTL/locale-independence tests;
- external webhook -> authenticity, replay and malformed payload tests;
- AI tool -> permission denial, prompt-injection resistance and approval boundary tests.

A happy-path-only change is incomplete when a governed negative invariant is affected.

## 3. Determinism

Required tests must be reproducible.

- Freeze/inject clocks where time matters.
- Seed randomness deterministically where randomness is not the feature under test.
- Avoid uncontrolled internet dependencies.
- External provider tests use recorded/sandboxed fixtures unless an explicit live-certification lane is required.
- Tests do not depend on execution order.
- Parallelizable tests must not share mutable tenant/database/file state unsafely.
- Flaky tests are failures; they are not accepted as normal noise.

## 4. Data and fixtures

- Synthetic fixtures are the default.
- Production CONFIDENTIAL/RESTRICTED data is prohibited in ordinary test environments.
- Every tenant-scoped fixture declares tenant ownership.
- Reference data is deterministic and versioned.
- Factories/builders prefer semantic defaults over random meaningless values.
- Test fixtures must not hide invalid assumptions by bypassing public constructors/validation unless the test explicitly targets corrupted state.

## 5. Database and migration tests

When persistence is affected, required evidence includes as applicable:

```text
fresh database -> all migrations -> seeds/reference data -> application boot
supported previous schema -> upgrade migrations -> application boot
```

Rules:

- migration execution is deterministic and repeatable;
- runtime and migration credentials are separable;
- destructive migrations require explicit change control;
- cross-tenant integrity is tested where schema/query behavior could leak scope;
- migration tests validate indexes/constraints critical to correctness;
- no manual database edit is part of a supported installation path.

## 6. Contract tests

### HTTP

- implementation conforms to governed OpenAPI;
- unknown/forbidden fields follow contract policy;
- status/error/problem shapes are validated;
- pagination/idempotency/concurrency semantics are exercised where applicable.

### Events

- producer payload conforms to versioned schema;
- event identity/type/source/tenant/correlation fields conform to `EVENT_STANDARD.md`;
- consumers tolerate duplicate delivery;
- compatibility tests protect supported old event versions.

### Modules

- manifest/dependencies/permissions/capabilities are schema-valid;
- private cross-module imports/tables are rejected by architecture checks where technically enforceable.

## 7. Security test baseline

Later executable phases must provide negative tests for affected risks, including:

- unauthenticated access;
- unauthorized principal;
- wrong tenant;
- wrong organization/object relationship;
- restricted-field disclosure;
- privilege escalation/mass assignment;
- token/session revocation where relevant;
- CSRF/CORS/browser boundary where relevant;
- webhook signature/replay;
- SSRF/egress restrictions;
- file authorization/path/content-type boundaries;
- support impersonation restrictions;
- AI/tool capability denial and approval enforcement.

Security scans supplement these tests; they do not replace them.

## 8. Module lifecycle certification

A module cannot be considered mature until applicable lifecycle tests prove:

```text
install -> enable -> use -> disable -> re-enable -> upgrade
```

Optional modules additionally prove unrelated domains keep working when the module is unavailable. Purge/destructive removal is a separate privileged operation and must not be conflated with disable/uninstall.

## 9. Test doubles

Use the lowest-fidelity double that preserves the contract being tested:

- pure fake for deterministic unit behavior;
- adapter stub for provider error mapping;
- protocol-compatible test container/service for integration;
- provider sandbox/live certification only when the external system itself must be proven.

Never mock away the exact boundary whose behavior the test is supposed to prove.

## 10. Coverage policy

Line/branch coverage is a diagnostic metric, not a completion metric.

- No universal percentage is sufficient evidence by itself.
- Critical authorization, money, billing, ledger, payment, migration, module lifecycle and security code should approach exhaustive branch/invariant coverage.
- Untested high-risk paths block release even if aggregate coverage is high.
- Coverage regressions on critical packages require explicit review.

## 11. Flaky test policy

A test that fails intermittently is quarantinable only with:

- issue/reference;
- named owner;
- expiry/review date;
- proof that quarantine does not remove a critical release gate;
- root-cause/fix plan.

Retries may diagnose transient infrastructure but must not convert nondeterministic product behavior into green evidence.

## 12. Test naming and evidence

Tests should describe invariant + expected outcome, not implementation detail.

Evidence records the exact command/job, target environment and result. Allowed evidence states remain `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, `N/A`.

## 13. Release-blocking principle

If a required test for an affected invariant fails, the change does not ship. If infrastructure prevents execution, status is `BLOCKED`, not PASS, and any exception must follow explicit change control.
