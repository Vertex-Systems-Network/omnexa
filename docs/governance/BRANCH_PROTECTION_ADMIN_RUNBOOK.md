# Omnexa Main Branch Protection Administration Runbook

Status: **P00.10 exit-verification tooling**  
Tracker: **Issue #3**

This runbook is the deterministic administration path for applying and verifying the hosted GitHub `main` protection required by `REPOSITORY_HARDENING.md` and `P01_ENTRY_GATE.md`.

It does not authorize P01 by itself. P01 remains locked until the hosted setting is applied, live GitHub evidence is captured, the controlled negative tests are completed, and canonical state is reconciled.

## 1. Security model

Use one of these owner-controlled authentication paths:

1. an interactive `gh auth login` session for an account with repository Administration permission; or
2. a short-lived fine-grained PAT exposed only to the current process as `OMNEXA_GITHUB_ADMIN_TOKEN`.

For a fine-grained PAT:

- select only `Vertex-Systems-Network/omnexa`;
- grant repository **Administration: read and write**;
- keep the lifetime short;
- do not commit, echo, paste into issue comments or store it in repository files;
- remove/rotate it immediately after the hosted setting and evidence are complete.

The ordinary Actions `GITHUB_TOKEN` is not an acceptable substitute for repository Administration authority.

## 2. Prerequisites

On the administration workstation (Windows/PowerShell):

```powershell
pwsh --version
gh --version
gh auth status -h github.com
```

If using a process-only PAT instead of an interactive GitHub CLI session:

```powershell
$env:OMNEXA_GITHUB_ADMIN_TOKEN = '<short-lived-token>'
```

Never print the variable value.

## 3. Apply the current Omnexa policy

From the repository root:

```powershell
pwsh -File .\scripts\apply_main_protection.ps1
```

Default P00.10 policy:

- branch: `main`;
- PR-based integration required;
- required status check: `governance`;
- required status check is strict/up-to-date;
- stale approvals dismissed;
- conversation resolution required;
- force pushes blocked;
- branch deletion blocked;
- protection enforced for administrators;
- required approval count: `0` while the repository has a single independent maintainer;
- Code Owner review is not required until an independent reviewer exists.

The approval count of `0` is intentional: it keeps PR-only integration enforced without creating an impossible self-approval deadlock for a single-maintainer private repository. It is not a permanent relaxation.

When a second independent reviewer is available, tighten the policy:

```powershell
pwsh -File .\scripts\apply_main_protection.ps1 -RequiredApprovals 1 -RequireCodeOwnerReview
```

## 4. Verify live GitHub configuration

Run:

```powershell
pwsh -File .\scripts\verify_main_protection.ps1
```

Expected final line:

```text
Omnexa main branch protection verification: PASS
```

For the future independent-reviewer policy:

```powershell
pwsh -File .\scripts\verify_main_protection.ps1 -MinimumApprovals 1 -RequireCodeOwnerReview
```

The verifier fails closed if any required protection property is absent.

## 5. Required Issue #3 evidence

Configuration verification is necessary but not sufficient. Before Issue #3 can close, record evidence for all of the following:

1. GitHub branch metadata reports `protected: true`.
2. Required check `governance` is present and strict/up-to-date.
3. Conversation resolution is required.
4. Force pushes are disabled.
5. Deletion of `main` is disabled.
6. Administrator enforcement is enabled.
7. A non-exempt direct push attempt is rejected.
8. A controlled force-push attempt is rejected.
9. A controlled branch-deletion attempt is rejected.
10. A test PR cannot merge while `governance` fails/is incomplete.
11. CODEOWNERS/review behavior matches the current contributor model.

Negative tests must use a controlled non-break-glass actor or a safe test procedure. Do not deliberately risk an irreversible `main` mutation merely to prove a rejection.

## 6. P00 -> P01 transition after verification

Once Issue #3 evidence is complete, a separate narrow governance transition must:

- mark P00.10 `done`;
- mark P00 `done`;
- activate P01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- retire ADR-0006 from active use while retaining it as historical evidence;
- record the final branch-protection and LOCAL-WIN-4 CI evidence;
- define/activate the first P01 kernel work package.

Do not combine that state transition with unrelated kernel feature code.

## 7. Cleanup

After administration:

```powershell
Remove-Item Env:OMNEXA_GITHUB_ADMIN_TOKEN -ErrorAction SilentlyContinue
Remove-Item Env:GH_TOKEN -ErrorAction SilentlyContinue
```

If a short-lived PAT was used, revoke/rotate it after verification evidence is recorded.
