# Omnexa Go Code Quality Gate

Status: **active repository quality control**  
Applies to: executable Go code under `kernel/` and later governed Go packages

This control supplements package-specific tests. It does not replace unit, integration, migration, security, lifecycle or build evidence.

## Canonical command

```text
bash scripts/verify_go_quality.sh
```

The same command runs in the required GitHub-hosted `governance` job.

## Pinned tools

- Go toolchain: repository `.go-version` exact pin;
- `golangci-lint v2.12.2`;
- `govulncheck v1.7.0`.

Tool installation must use the exact versions above. `@latest` is prohibited in canonical evidence.

## golangci-lint policy

`.golangci.yml` uses configuration version 2 and an explicit allow-list rather than enabling every available linter. The gate currently includes:

- `govet`;
- `staticcheck`;
- `errcheck`;
- `ineffassign`;
- `unused`;
- `errorlint`;
- `gosec`;
- `bodyclose`;
- `noctx`;
- `nilerr`;
- `nilnil`;
- `copyloopvar`;
- `misspell`;
- `unconvert`.

The list may change only through normal review when a rule materially improves signal. A noisy rule must not be disabled merely to make a failing change green; the defect is fixed first unless an explicit, narrowly justified exclusion is the correct contract.

## Vulnerability analysis

`govulncheck` runs against `./kernel/...` and analyzes reachable known Go vulnerabilities using the Go vulnerability database. A reachable known vulnerability fails the required gate. If the vulnerability database cannot be reached and the check is required, evidence is not PASS.

## Formatting and tests

`gofmt` cleanliness remains fail-closed through the completed P01 regression verifiers. Package-specific `go vet`, unit, race, integration, migration, security, lifecycle and build checks remain required where applicable.

No universal line-coverage percentage is introduced by this control. Omnexa's testing standard remains risk/invariant based; a high percentage cannot compensate for an untested critical path.

## AI / contributor rules

Human and AI contributors must not:

- use lint auto-fix as a substitute for reviewing semantic changes;
- add broad `//nolint` directives merely to silence CI;
- weaken or skip a quality tool because a new change fails it;
- replace pinned versions with floating `latest` installs;
- relabel tool outage, vulnerability-database outage, skipped analysis or unsupported execution as PASS;
- treat a clean static-analysis run as proof that runtime behavior, security, tenancy or business invariants are correct.

Any necessary suppression must be narrow, specific to a named linter/check, justified next to the code or in governed configuration, and reviewable.

## Evidence

A canonical quality PASS records at minimum:

- source SHA;
- Go version;
- golangci-lint version;
- govulncheck version;
- configuration validation result;
- lint result;
- vulnerability-analysis result;
- GitHub-hosted workflow run/job.
