# Omnexa Repository Hardening Baseline

Status: **Required governance configuration**

This document defines repository settings that protect the architecture/governance process from accidental or unauthorized drift. File-level controls are versioned in the repository; GitHub-hosted settings must be configured in repository rulesets/branch protection.

## 1. Required `main` protection

Target configuration for `main`:

- require pull request before merge;
- require at least one approving review for protected architecture/governance paths once the contributor model expands;
- require review from Code Owners where GitHub plan/settings support it;
- dismiss stale approvals when protected files change;
- require conversation resolution before merge;
- prohibit force pushes;
- prohibit branch deletion;
- block direct pushes except explicitly authorized emergency/break-glass actors;
- require the Omnexa `governance` CI check;
- require branch to be up to date before merge when CI/build topology makes that necessary;
- do not allow ordinary administrators/maintainers to bypass the ruleset unless a documented break-glass process explicitly grants it.

## 2. Break-glass rule

Emergency bypass must be exceptional. Any bypass must produce:

1. reason and incident reference;
2. exact files/commits changed;
3. security/architecture impact statement;
4. follow-up PR/verification if normal gates were skipped;
5. ledger entry when execution state or architecture evidence changes.

Emergency access must never be used for convenience or speed.

## 3. File-level controls

Repository-managed controls include:

- `.github/CODEOWNERS`;
- `.github/pull_request_template.md`;
- `.github/workflows/governance.yml`;
- `.github/workflows/main-protection-admin.yml`;
- `scripts/validate_governance.py`;
- `scripts/validate_p01_preparation.py`;
- `AGENTS.md`;
- canonical architecture/governance/state documents.

## 4. Critical paths

The following paths require heightened review because they can redefine what future AI systems are allowed to do:

```text
AGENTS.md
docs/governance/**
docs/architecture/**
docs/roadmap/**
docs/adr/**
.github/**
```

Changes to `STATE.json` must be treated as execution-control changes, not ordinary documentation edits.

## 5. Required checks and runner policy

The canonical required check is named `governance` and runs on GitHub-hosted `ubuntu-24.04`. The workflow must fail unless `RUNNER_ENVIRONMENT=github-hosted`, Linux and X64. Local/self-hosted Actions runners are prohibited for canonical governance CI and repository administration automation.

Future P01+ quality gates may add jobs/checks, but must preserve the P00.07 semantics and hosted-only execution policy unless a later explicit governance change replaces it.

## 6. Current configuration evidence

The repository is now public. The former private-plan entitlement blocker is historical. Live GitHub metadata still reports `main` as unprotected with required status checks disabled, so the repository must retain Issue #3 until the hosted branch/ruleset settings match this specification.

Hosted CI proof: run `32537207455`, job `96940269306`, `RUNNER_ENVIRONMENT=github-hosted`, Ubuntu 24.04.4 LTS / X64, with all canonical validators PASS.

## 7. Verification procedure

After hosted settings are configured:

1. verify GitHub reports `protected:true`;
2. verify strict required check `governance` is configured;
3. verify direct non-exempt push to `main` is rejected;
4. verify force push is rejected;
5. verify protected branch deletion is rejected;
6. open a test PR and verify governance CI is required;
7. modify a CODEOWNERS path and verify required owner review behavior when enabled;
8. record evidence in Issue #3 and the execution ledger if this changes a formal gate.

## 8. Deterministic administration tooling

P00.10 provides repository-managed administration helpers:

- `docs/governance/BRANCH_PROTECTION_ADMIN_RUNBOOK.md`;
- `.github/workflows/main-protection-admin.yml`;
- `scripts/apply_main_protection.ps1`;
- `scripts/verify_main_protection.ps1`.

The administration workflow is `workflow_dispatch` only, runs on GitHub-hosted `ubuntu-24.04`, requires `RUNNER_ENVIRONMENT=github-hosted`, and consumes a short-lived owner-controlled `OMNEXA_GITHUB_ADMIN_TOKEN` GitHub Actions secret with repository Administration read/write permission.

The ordinary Actions `GITHUB_TOKEN` is not sufficient repository Administration authority and must not be used as a substitute.

The current single-maintainer policy uses required approval count `0` while still requiring PR-based integration, because GitHub does not allow an author to approve their own PR. Once an independent reviewer exists, the baseline must be tightened to one or more approvals and Code Owner review according to the contributor model.

The verifier is fail-closed and checks the hosted configuration, but it does not replace the controlled negative tests listed above. Issue #3 remains open until live GitHub evidence and the required rejection tests are complete.
