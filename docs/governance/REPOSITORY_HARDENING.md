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
- require the Omnexa governance CI check once the check has run successfully and is stable;
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
- `scripts/validate_governance.py`;
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

## 5. Required checks

The initial governance check validates repository/state consistency only. P00.07 will define the complete build/test/security/release check set. Do not falsely treat this initial workflow as the final CI architecture.

## 6. Current configuration evidence

At the time this baseline was written, GitHub reported `main` as unprotected with required status checks disabled. The repository must retain a tracked issue until the hosted branch/ruleset settings match this specification.

## 7. Verification procedure

After hosted settings are configured:

1. verify direct non-exempt push to `main` is rejected;
2. verify force push is rejected;
3. verify protected branch deletion is rejected;
4. open a test PR and verify governance CI is required;
5. modify a CODEOWNERS path and verify required owner review behavior when enabled;
6. record evidence in the tracking issue and execution ledger if this changes a formal gate.

## 8. Deterministic administration tooling

P00.10 provides repository-managed administration helpers so the hosted setting can be applied and verified without reconstructing the policy manually:

- `docs/governance/BRANCH_PROTECTION_ADMIN_RUNBOOK.md`;
- `scripts/apply_main_protection.ps1`;
- `scripts/verify_main_protection.ps1`.

The apply script uses GitHub's branch-protection Administration API through the GitHub CLI and requires either an owner-controlled authenticated `gh` session or a short-lived process-only `OMNEXA_GITHUB_ADMIN_TOKEN` with repository Administration read/write permission.

The current single-maintainer policy uses required approval count `0` while still requiring PR-based integration, because GitHub does not allow an author to approve their own PR. Once an independent reviewer exists, the baseline must be tightened to one or more approvals and Code Owner review according to the contributor model.

The verifier is fail-closed and checks the hosted configuration, but it does not replace the controlled negative tests listed above. Issue #3 remains open until live GitHub evidence and the required rejection tests are complete.
