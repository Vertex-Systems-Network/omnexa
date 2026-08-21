# Omnexa CI Standard

Status: **Canonical v1**  
Work package: **P00.07**

CI is an execution environment for canonical quality commands, not the owner of quality semantics. Omnexa quality gates must be runnable locally and by any approved CI provider so hosted-runner quotas or provider changes do not redefine correctness.

## 1. Provider-independent rule

The authoritative gates are repository commands/contracts. GitHub Actions, GitLab CI, Buildkite, Jenkins or another runner may orchestrate them, but may not invent weaker semantics.

Future repository tooling must expose stable commands for at least:

```text
verify:governance
verify:format
verify:lint
verify:static
verify:unit
verify:contracts
verify:integration
verify:migrations
verify:security
verify:module-lifecycle
verify:build
verify:release
verify:all
```

Exact implementation may use Make, Task, scripts or a cross-platform CLI, but CI and local developers call the same underlying checks.

## 2. Required CI lanes

### Fast PR lane

Runs on every material pull request and should provide rapid feedback:

- governance/state checks;
- formatting/lint/static/type checks;
- unit tests;
- contract/schema validation;
- generated-artifact drift checks;
- secret scanning;
- dependency/license policy checks where configured.

### Integration lane

Runs when affected paths require real infrastructure:

- PostgreSQL integration;
- Redis-compatible integration;
- event/broker integration;
- object storage integration;
- migration fresh/upgrade paths;
- module lifecycle tests;
- cross-tenant/authorization negative tests.

### Release-certification lane

Runs on release candidate/tag/promotion and may include:

- full test matrix;
- production-like build;
- SBOM/provenance generation;
- artifact signing/verification;
- vulnerability scans;
- packaging/install/upgrade certification;
- supported OS/browser/device/provider certification where applicable;
- backup/restore or resilience rehearsal for releases that affect those paths.

### Scheduled assurance lane

For expensive or ecosystem-sensitive checks:

- dependency vulnerability refresh;
- longer fuzz/property tests;
- performance regressions;
- provider sandbox/live certification;
- restore drills;
- compatibility matrix refresh.

## 3. Change-based selection

Path filtering may skip irrelevant expensive jobs, but may not skip invariant checks just because a file path looks unrelated.

A change-impact map eventually links paths/domains to required gates. When impact cannot be proven safely, run the broader gate.

## 4. Fail-closed semantics

Required gates have only these outcomes:

- `PASS` — observed successful execution;
- `FAIL` — executed and failed;
- `BLOCKED` — infrastructure/dependency unavailable;
- `NOT RUN` — not executed;
- `N/A` — demonstrably irrelevant.

`BLOCKED` and `NOT RUN` are never green equivalents. Release/merge exceptions require explicit governance and expiry.

## 5. Reproducibility

- Pin toolchain/runtime major/minor/patch as appropriate to release policy.
- Lock dependency graphs.
- Use reproducible container/service versions for CI dependencies.
- Avoid mutable `latest` images/tags for required gates.
- Cache only performance artifacts; correctness must not depend on cache presence.
- CI should be able to start from a clean checkout with no hidden machine state.

## 6. Least-privilege CI

The canonical CI credential and permission rule is `least-privilege`: grant only the minimum permissions needed for the specific job and trust boundary.

- Default workflow/token permissions are read-only/minimal.
- Secrets are scoped by environment/job and unavailable to untrusted fork code.
- Pull requests do not receive production deploy/signing credentials.
- Release credentials live behind protected environments/approval boundaries.
- CI logs/artifacts obey data classification rules.
- Third-party CI/actions/plugins are dependencies subject to review/pinning.

## 7. Dependency and supply-chain gates

As tooling becomes available, required lanes must address:

- dependency lockfile integrity;
- vulnerability scanning;
- secret scanning;
- license-policy checks;
- SBOM generation for release artifacts;
- provenance/build identity;
- signature/checksum verification;
- malicious/compromised dependency response;
- review of executable install/build scripts.

Scanner output is triaged by severity/exploitability/context; silently suppressing findings is prohibited.

## 8. Generated artifacts

Generated OpenAPI clients, schemas, codegen, migrations metadata or documentation must be reproducible.

CI either:

- generates and compares for zero diff; or
- generates release artifacts from canonical inputs.

Hand-edited generated output without updating its source is invalid.

## 9. Matrix strategy

Matrices must reflect supported environments, not every imaginable environment.

Potential dimensions include:

- Go/runtime versions;
- Node/browser versions;
- PostgreSQL supported majors;
- OS/architecture;
- browser engines;
- POS/device platform;
- deployment topology;
- connector/provider sandbox.

The support policy controls required combinations and may evolve through release policy.

## 10. Concurrency and cancellation

Superseded PR runs may be cancelled to conserve capacity, but release certification must not be silently cancelled by a newer unrelated run.

CI concurrency controls must not permit two release jobs to mutate the same release/promotion environment unsafely.

## 11. Artifact retention

Retain enough evidence for diagnosis and release traceability without retaining sensitive data unnecessarily.

Release evidence should include:

- commit SHA;
- source tree identity;
- toolchain versions;
- gate results;
- artifact digests;
- SBOM/provenance references where applicable;
- release approver/promotion reference.

## 12. Branch/PR enforcement

Target policy for `main`:

- PR required;
- required quality checks from canonical gates;
- required conversations resolved;
- force push/delete blocked;
- bypass restricted to explicit break-glass actors;
- CODEOWNERS/approval policy as contributor model permits.

Hosted configuration remains tracked separately because repository files cannot alone enforce GitHub settings.

## 13. CI outage / quota behavior

CI capacity exhaustion is an infrastructure state, not a product pass.

- record `BLOCKED`;
- use approved temporary exceptions only if scope and risk permit;
- preserve the commands/configuration that would have run;
- rerun once capacity returns;
- never allow a temporary documentation exception to become precedent for executable release bypass.

ADR-0006 is one such temporary P00-only exception and expires before P01 implementation.

## 14. No CI-only fixes

If code passes CI only because of CI-specific environment hacks that local/release environments do not share, the pipeline is masking a defect. Canonical commands must remain portable and environment assumptions explicit.
