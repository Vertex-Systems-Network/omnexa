# Omnexa Repository Hardening Baseline

Status: **Applied and verified for current single-maintainer model**

This document defines repository settings that protect architecture/governance from accidental or unauthorized drift. File controls are versioned in-repo; hosted rules are enforced by GitHub rulesets/branch protection.

## 1. Required `main` policy

Current required baseline:

- pull request required before merge;
- required `governance` status check;
- conversation resolution required;
- direct pushes blocked;
- force pushes blocked;
- branch deletion blocked;
- ordinary bypass disabled; break-glass must be explicit and documented;
- single-maintainer mode: required approvals `0`, required Code Owner review disabled to avoid self-review deadlock;
- once an independent reviewer exists: require at least one independent approval and Code Owner review, with stale-approval policy appropriate to contributor model.

## 2. Verified 2026-08-22 evidence — Issue #3

Issue #3 is **closed/completed**.

Evidence:

- live GitHub API reports `main protected:true`;
- failed-check PR #34 / run `32540836431` was rejected because required `governance` was failing;
- direct fast-forward update probe commit `44ca19e80c5fccccebfd8d4f96dde6dc5af14bc2` was rejected because changes must be made through a PR and `governance` was expected;
- force-update probe was rejected with `Cannot force-push to this branch`;
- CODEOWNERS-path PR #37 / hosted run `32541439589` was blocked while an inline conversation remained unresolved; after resolution the same green PR merged as `866646f5a2db444fc668dd62b8d1ff824b6359bc`;
- valid green PR #35 merged as `843c615170058ab900ba69516dbed80a47f26973`;
- deletion of `main` remains blocked by configured ruleset. A destructive default-branch deletion test was intentionally not performed because it would create unnecessary recovery risk.

## 3. Break-glass rule

Emergency bypass must be exceptional and explicitly authorized. Any use must record reason/incident, exact commits/files, security/architecture impact, skipped gates, follow-up verification and ledger reconciliation. Never use bypass for convenience or speed.

## 4. Repository-managed controls

- `.github/CODEOWNERS`;
- `.github/pull_request_template.md`;
- `.github/workflows/governance.yml`;
- `.github/workflows/main-protection-admin.yml`;
- `.github/dependabot.yml`;
- `scripts/validate_governance.py`;
- `scripts/validate_freeze_review.py`;
- `scripts/validate_p01_preparation.py`;
- `scripts/validate_p01_package_specs.py`;
- `AGENTS.md` and canonical governance/architecture/state documents.

Critical execution-control paths include `AGENTS.md`, `docs/governance/**`, `docs/architecture/**`, `docs/roadmap/**`, `docs/adr/**` and `.github/**`. `STATE.json` changes are execution-control changes, not ordinary prose edits.

## 5. CI and runner policy

The canonical required check is `governance` and runs only on GitHub-hosted `ubuntu-24.04`. The workflow must fail unless `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64.

Local/self-hosted Actions runners are prohibited for canonical governance CI and repository administration automation. Future jobs may expand quality coverage but may not silently weaken P00.07 semantics or hosted-only execution.

## 6. Deterministic administration tooling

Repository-managed helpers remain for audit/reapplication:

- `docs/governance/BRANCH_PROTECTION_ADMIN_RUNBOOK.md`;
- `.github/workflows/main-protection-admin.yml`;
- `scripts/apply_main_protection.ps1`;
- `scripts/verify_main_protection.ps1`.

The admin workflow is manual (`workflow_dispatch`) and GitHub-hosted only. It requires an owner-controlled short-lived `OMNEXA_GITHUB_ADMIN_TOKEN` with repository Administration read/write. Ordinary `GITHUB_TOKEN` is not Administration authority.

## 7. Verification rule

Protection must be demonstrated through live metadata plus controlled negative/positive behavior; documentation alone is not proof. If protection is changed later, rerun equivalent direct-push, required-check, conversation and valid-merge probes without using destructive tests unnecessarily.
