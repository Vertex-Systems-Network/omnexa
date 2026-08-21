# Omnexa Toolchain Standard

Status: **Canonical v1**  
Work package: **P00.08**

Omnexa pins development/build toolchains so a clean checkout can reproduce the supported environment. Repository configuration, not a developer's globally installed latest version, defines the supported toolchain.

## Runtime families

The architectural language baseline remains:

- Go for kernel/backend/domain services;
- TypeScript + React for admin/web/builder/SDK surfaces;
- Rust for justified edge/native/security-sensitive components;
- Python only for AI/data workloads where ecosystem value warrants it.

P00.08 does not freeze a specific patch forever. It freezes **how versions are selected and upgraded**.

## Version policy

Each active runtime/tool must have one repository-owned version declaration or lock source.

Preferred mechanisms include:

- Go: `go.mod` / workspace metadata plus pinned CI/bootstrap tool versions;
- Node: repository package-manager declaration + lockfile + runtime version declaration;
- Rust: `rust-toolchain.toml` plus Cargo lock policy appropriate to application/library package type;
- Python: explicit interpreter range/version plus locked environment/dependency artifact for each actual Python workload;
- containers: immutable version tags, preferably digest-pinned for release-critical infrastructure.

Do not maintain contradictory version numbers across README, CI and local scripts.

## Package manager policy

The JavaScript/TypeScript workspace uses **one** package manager for the monorepo. The exact selection is made before the first JS workspace implementation and recorded once; mixing npm/yarn/pnpm lockfiles is forbidden.

Go/Rust/Python use their ecosystem-native dependency tooling under repository policy.

## Lockfiles

- Application/deployable dependency graphs are locked.
- Lockfiles are reviewed source artifacts.
- Required CI/local builds use frozen/locked dependency resolution modes.
- Regenerating a lockfile is an intentional dependency change, not incidental noise.
- Dependency updates include compatibility/security evidence appropriate to risk.

## Tool installation

Formatters, linters, schema tools, generators and security utilities must be:

1. part of the runtime toolchain;
2. installed through a repository-pinned tool manifest/bootstrap; or
3. executed through a pinned container.

Unversioned `curl | sh` or `install latest` steps are not canonical build prerequisites.

## Upgrade process

A material runtime/toolchain upgrade requires:

- compatibility review;
- lockfile/regeneration review;
- affected quality gates;
- local/bootstrap update;
- CI/release environment update;
- documentation update;
- ADR only when architecture/support policy changes materially.

Toolchain upgrades should be isolated from unrelated feature work where practical.

## Supported-version policy

Omnexa should normally track maintained stable versions rather than indefinite legacy runtimes. Exact supported ranges are declared when implementation begins and release support commitments exist.

A runtime may not be upgraded simply because a newer version exists; upgrade evidence must show compatibility and operational benefit/risk.

## Compiler/linter strictness

Compiler/type/lint settings are repository-owned. Teams/modules may add stricter local rules but may not disable shared correctness/security rules without change control.

Suppressions require narrow scope and explanation; blanket ignore files are prohibited unless explicitly justified.

## Generated code tools

Code generation versions are pinned. Generated output must be reproducible from the same canonical inputs and generator version.

## Platform architecture vs toolchain

Tooling may support a future service extraction, but it must not force microservice topology. Workspace/build tooling must support the modular-monolith baseline efficiently.

## Toolchain evidence

Future `omnexa dev status` or equivalent should report the resolved versions of supported toolchains and important local infrastructure, making environment drift visible before debugging begins.