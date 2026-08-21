# ADR-0010 — Foundation Architecture Freeze v1

Status: **accepted**

## Context

P00.01 through P00.09 established Omnexa's governance, terminology/ownership, primitive semantics, API/events, security, quality/release, repository/local-development and threat/reliability foundations.

The program now needs an explicit decision separating **architecture freeze acceptance** from **permission to begin executable kernel implementation**.

## Problem

Treating architecture completion and implementation readiness as the same state would allow P01 code to begin while repository protection and executable verification are unavailable. Conversely, keeping the architecture perpetually unfrozen because of hosted operational blockers would make architecture state inaccurate.

## Decision

1. Accept and freeze P00.01–P00.09 as **Omnexa Foundation Architecture v1**.
2. Keep P00.10 / P00 exit in verification until implementation-entry controls are satisfied.
3. Classify Issue #3 (`main` branch/ruleset protection) as a **P01 implementation-entry blocker**.
4. Classify Issue #14 (executable CI lane unavailable) as a **P01 implementation-entry blocker**.
5. Classify Issue #4 (licensing/IP/trademark) as an **external distribution/public-launch blocker**, not an internal private P01 engineering blocker.
6. Keep `kernel_code_authorized = false` and `business_feature_code_authorized = false` while P00 exit is in verification.
7. Permit P01 activation only through a narrow entry-transition change that records verified protection + executable CI evidence, expires ADR-0006, closes P00.10/P00 and activates P01.
8. P00 temporary CI exception cannot be used as precedent for P01 executable work.

## Frozen architecture scope

Foundation v1 includes accepted ADR-0001 through ADR-0009 and the normative documents referenced by `docs/governance/FOUNDATION_FREEZE_REVIEW.md`.

## Alternatives considered

### Mark P00 done immediately and mark P01 blocked
Rejected for the current machine-state model because it would create a phase transition without implementation-entry controls and could invite AI agents to infer kernel authorization from P00 completion.

### Keep architecture unfrozen until hosted blockers are solved
Rejected. Issue #3/#14 do not invalidate the architecture decisions themselves.

### Allow local-only P01 coding without merge protection/CI
Rejected as canonical repository policy. Experimental work outside the governed repository is not P01 progress; governed executable P01 changes require the entry gate.

## Consequences

Positive:
- architecture state is truthful: frozen but implementation remains locked;
- AI systems have a machine-readable P01 entry gate;
- hosted blockers cannot silently become exceptions for executable work;
- licensing/trademark work does not unnecessarily stop private engineering.

Costs:
- P00 remains programmatically active in exit verification until issue #3/#14 are resolved;
- P01 cannot begin in the canonical repository until a viable executable CI lane and branch protection exist.

## Compatibility impact

No runtime impact; no runtime exists in P00.

## Migration impact

None.

## Security/tenancy impact

Strengthens repository-level protection before executable security/tenant code enters `main`.

## Operational impact

Requires executable CI capacity and protected integration policy before P01.

## Rollback / forward-fix

If a frozen P00 baseline proves contradictory or unsafe, reopen the affected package through a superseding ADR rather than bypassing the freeze. If entry blockers are satisfied, use the documented P01 transition procedure; do not supersede this ADR merely to avoid a gate.