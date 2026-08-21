# Omnexa Release Standard

Status: **Canonical v1**  
Work package: **P00.07**

A release is the promotion of an immutable, verified source/artifact identity through controlled environments. Rebuilding different bytes for each environment is not a release process.

## 1. Versioning

Omnexa uses semantic versioning semantics for platform and independently versioned modules/packages where public compatibility matters:

```text
MAJOR.MINOR.PATCH
```

- MAJOR: intentional incompatible contract/support boundary change;
- MINOR: backward-compatible capabilities/features;
- PATCH: backward-compatible fixes/security/maintenance.

Pre-release identifiers may represent alpha/beta/rc channels. Version numbers do not replace API/event/module contract versions; those remain explicit governed contracts.

## 2. Immutable source identity

Every release candidate maps to exactly one immutable source commit/tag and one set of generated artifact digests.

Required metadata eventually includes:

- release version;
- commit SHA;
- tag identity;
- build/toolchain identity;
- artifact digests;
- SBOM/provenance references;
- schema/migration baseline;
- supported upgrade-from range;
- release notes/changelog;
- known risks/exceptions;
- approval/promotion evidence.

## 3. Build once, promote

Preferred release flow:

```text
source commit
 -> verified build
 -> immutable artifacts
 -> test/staging verification
 -> promotion
 -> production
```

Do not rebuild application artifacts independently for staging and production when the artifact type supports promotion. Environment-specific configuration/secrets remain external to the immutable artifact.

## 4. Release branches/tags

- `main` represents the current integration truth.
- Release tags are immutable.
- Maintenance/release branches may exist only when supporting parallel maintained versions justifies them.
- No long-lived environment branches.
- Hotfixes are merged back into canonical development history.

## 5. Release gates

A production-capable release is blocked unless all applicable gates pass:

- governance/state consistency;
- format/lint/static/type;
- unit/component/contract/integration tests;
- fresh/upgrade migration certification;
- tenant/authorization/security negative tests;
- module lifecycle certification;
- production build/package verification;
- dependency/secret/vulnerability/license policy checks;
- SBOM/provenance/signature checks when implemented;
- release notes and compatibility assessment;
- deployment/rollback or forward-fix plan;
- required operational/SLO/readiness evidence.

`BLOCKED` or `NOT RUN` required release gates prohibit production promotion unless a separately approved emergency production policy explicitly applies. ADR-0006 is **not** such a policy.

## 6. Database compatibility

Schema evolution must support the deployment model.

Prefer expand/contract where zero/low-downtime deployments require old/new application versions to coexist:

1. additive compatible schema;
2. deploy code capable of both states;
3. backfill/migrate data safely;
4. switch reads/writes;
5. remove old structure only after compatibility window.

A rollback plan must account for irreversible data changes. When rollback is unsafe, define an explicit forward-fix/recovery plan.

## 7. API/event compatibility

Release notes identify public contract changes.

- breaking HTTP contract -> new governed major contract version;
- breaking event contract -> new event major version/type;
- supported old versions continue for their declared window;
- deprecation/removal follows documented policy;
- internal refactors must not force external consumer changes unnecessarily.

## 8. Release channels

Potential channels:

```text
internal -> alpha -> beta -> rc -> stable -> maintenance
```

Not every module needs every channel. Channel policy defines audience, stability expectations, data compatibility and support commitments.

## 9. Artifact integrity and provenance

Release artifacts should become progressively supply-chain verifiable:

- content digest/checksum;
- SBOM;
- build provenance/attestation;
- signed artifacts/packages where required;
- signature verification at install/deploy boundaries;
- immutable package/version identity.

Marketplace, edge/POS installers and self-hosted distributions require especially strong package identity/integrity.

## 10. Changelog and release notes

Release notes are user/operator-facing impact, not a raw commit dump.

Include as relevant:

- added capabilities;
- behavior changes;
- fixes/security updates;
- migrations/upgrade notes;
- deprecations/removals;
- compatibility/support changes;
- known issues;
- rollback/operational cautions.

Security disclosures may be staged to avoid unsafe premature detail.

## 11. Deployment promotion

Promotion requires explicit target environment/region/tenant cohort and artifact digest.

Strategies may include:

- rolling;
- blue/green;
- canary;
- tenant cohort;
- regional wave;
- feature-flagged activation.

Deployment strategy must match data/schema and service compatibility constraints.

## 12. Rollback and forward fix

Each material release declares:

- rollback feasibility;
- application rollback steps;
- schema compatibility;
- cache/event implications;
- external side-effect implications;
- when forward-fix is safer than rollback.

A release with irreversible financial/external side effects cannot pretend rollback erases those effects.

## 13. Emergency release

Emergency security/availability fixes may use an expedited path only if a dedicated emergency policy defines:

- authorized initiators/approvers;
- minimum non-waivable tests;
- security/tenant invariants that remain mandatory;
- evidence and post-incident reconciliation;
- time-bounded exception cleanup.

Normal deadline pressure is not an emergency.

## 14. Release certification record

Every stable release eventually produces a machine-readable certification record that can answer:

```text
what source?
what artifacts?
what tests/gates?
what migrations?
what dependencies/SBOM?
who/what approved promotion?
where deployed?
what rollback/forward-fix path?
```

## 15. No production release during P00

P00 defines standards only. The temporary hosted-CI exception allows P00 specification progress; it cannot be used to certify or release executable Omnexa software.
