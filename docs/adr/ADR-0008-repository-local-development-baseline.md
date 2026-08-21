# ADR-0008 — Repository & Local Development Baseline

Status: **Accepted**  
Date: **2026-08-22**  
Work package: **P00.08**

## Context

Omnexa needs a repository/developer model that supports a large modular platform without forcing premature microservices, hidden local dependencies or CI-only quality workflows. Multiple runtimes and future deployment targets increase the risk of inconsistent developer environments and cross-module architectural leakage.

## Decision

1. Omnexa remains one governed monorepo while the architecture begins as a strict modular monolith.
2. Canonical root ownership categories are `apps`, `kernel`, `modules`, `platform`, `shared`, `infrastructure`, `scripts`, `docs` and `generated`.
3. Directory boundaries express ownership, not automatic deployment boundaries.
4. Modules keep private implementation/schema/migrations with their owner and expose governed contracts only.
5. Default local infrastructure is containerized PostgreSQL, Redis-compatible cache, NATS/JetStream and S3-compatible object storage, adding other services only when active phases require them.
6. Kubernetes is not required for the default local developer loop.
7. Local verification is first-class and uses the same semantic `verify:*` gates as CI.
8. Toolchain and dependency versions are repository-owned/pinned; unversioned global tooling is not a supported prerequisite.
9. JavaScript/TypeScript uses one monorepo package manager and one lockfile strategy once implementation begins.
10. Local configuration is explicit; secrets are separate, ignored from source control and replaced by safe examples/references.
11. Development/test data is synthetic/deterministic by default; production sensitive data is prohibited from ordinary local environments.
12. Linux is canonical for backend execution; macOS is supported where upstream tooling permits; Windows backend development prefers WSL2, while native Windows is a separate certification target where POS/edge requirements need it.
13. Future bootstrap/dev/db/verify/module command semantics are governed by `DEVELOPER_COMMANDS.md` and must not rely on hidden manual steps.
14. Generated artifacts are reproducible derived output and never become an alternative source of truth.

Normative details:

- `docs/development/REPOSITORY_STRUCTURE.md`
- `docs/development/LOCAL_DEVELOPMENT.md`
- `docs/development/TOOLCHAIN_STANDARD.md`
- `docs/development/CONFIGURATION_STANDARD.md`
- `docs/development/DEVELOPER_COMMANDS.md`
- `docs/contracts/development/workspace.schema.json`

## Consequences

### Positive

- ownership remains visible as the repository grows;
- local development is independent of GitHub Actions availability;
- onboarding can converge on one reproducible bootstrap path;
- premature Kubernetes/microservice complexity is avoided;
- local/CI quality semantics remain aligned.

### Costs

- multi-runtime toolchain/bootstrap maintenance is centralized;
- module boundaries need enforcement tooling later;
- WSL2 is a recommended Windows backend prerequisite rather than pretending all shell/filesystem behavior is identical natively;
- local container dependencies consume developer resources.

## Rejected alternatives

### Polyrepo from day one

Rejected because it increases contract/version/coordination overhead before deployment/team boundaries justify it.

### Kubernetes as mandatory local runtime

Rejected because it adds cost/complexity without being needed to prove modular architecture.

### CI-only verification

Rejected because CI provider quotas/outages must not define developer correctness.

### Unpinned global tools

Rejected because they create invisible environment drift.

### Production database snapshots for normal development

Rejected under the P00.06 classification/security baseline.

## Compatibility

This ADR implements the development environment required to execute P00.07 quality semantics and preserves ADR-0001 modular-monolith-first architecture plus P00.06 security/data rules. P00.09 may add reliability/threat-driven local simulation requirements without weakening these rules.

## Supersession

Material changes to monorepo ownership model, local execution baseline, toolchain ownership/pinning, canonical developer command semantics or default platform support require a superseding ADR.