# Omnexa Main Branch Protection Administration Runbook

Status: **P00.10 exit-verification tooling — currently plan-blocked**  
Tracker: **Issue #3**

This runbook is the deterministic administration path for applying and verifying the hosted GitHub `main` protection required by `REPOSITORY_HARDENING.md` and `P01_ENTRY_GATE.md`.

It does not authorize P01 by itself. P01 remains locked until the hosted setting is applied and verified, or EG-02 is deliberately superseded through accepted governance change control.

## 1. Current hosted-plan status

On 2026-08-22 the owner executed the merged apply script from an authenticated admin workstation against the private organization repository. GitHub returned:

```text
gh: Upgrade to GitHub Pro or make this repository public to enable this feature. (HTTP 403)
```

Verified facts:

- repository: `Vertex-Systems-Network/omnexa`;
- owner type: organization;
- visibility: `private`;
- authenticated account has repository `admin` permission;
- script syntax/policy validation passes on `LOCAL-WIN-4`;
- `main` remains unprotected.

This is a product-plan entitlement failure. Do not keep retrying the same API call while these conditions are unchanged.

GitHub's current feature model provides protected branches/rulesets for public repositories on GitHub Free; private repositories require an eligible paid plan. For this organization-owned private repository, resolve the plan entitlement before rerunning this runbook unless EG-02 is superseded through an accepted governance ADR.

Changing repository visibility to public merely to bypass the plan limitation is not an automatic remediation. It requires explicit owner/legal/security review and must account for Issue #4 licensing/IP/trademark concerns.

## 2. Security model

Use one of these owner-controlled authentication paths only after the plan/visibility prerequisite is satisfied:

1. an interactive `gh auth login` session for an account with repository Administration permission; or
2. a short-lived fine-grained PAT exposed only to the current process as `OMNEXA_GITHUB_ADMIN_TOKEN`.

For a fine-grained PAT:

- select only `Vertex-Systems-Network/omnexa`;
- grant repository **Administration: read and write**;
- keep the lifetime short;
- do not commit, echo, paste into issue comments or store it in repository files;
- remove/rotate it immediately after the hosted setting and evidence are complete.

The ordinary Actions `GITHUB_TOKEN` is not an acceptable substitute for repository Administration authority.

## 3. Prerequisites

Before running the apply script again, all of these must be true:

- the repository remains intentionally private or has an explicitly approved visibility change;
- the GitHub plan supports branch protection/rulesets for the repository's visibility;
- the authenticated account has repository Administration permission;
- the canonical `governance` check is stable and green.

On the administration workstation (Windows/PowerShell):

```powershell
gh --version
gh auth status -h github.com
```

If using a process-only PAT instead of an interactive GitHub CLI session:

```powershell
$env:OMNEXA_GITHUB_ADMIN_TOKEN = '<short-lived-token>'
```

Never print the variable value.

## 4. Apply the current Omnexa policy

From the repository root:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\apply_main_protection.ps1
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
powershell -ExecutionPolicy Bypass -File .\scripts\apply_main_protection.ps1 -RequiredApprovals 1 -RequireCodeOwnerReview
```

If GitHub again returns HTTP 403 with an upgrade/public-repository message, stop. Record the evidence in Issue #3 and resolve plan/visibility/governance policy; do not troubleshoot the script as if it were an authentication failure.

## 5. Verify live GitHub configuration

Run:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\verify_main_protection.ps1
```

Expected final line:

```text
Omnexa main branch protection verification: PASS
```

For the future independent-reviewer policy:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\verify_main_protection.ps1 -MinimumApprovals 1 -RequireCodeOwnerReview
```

The verifier fails closed if any required protection property is absent.

## 6. Required Issue #3 evidence

Configuration verification is necessary but not sufficient. Before Issue #3 can close under the current EG-02 policy, record evidence for all of the following:

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

## 7. P00 -> P01 transition after verification

Once Issue #3 evidence is complete, or EG-02 has been deliberately replaced by an accepted superseding governance ADR, a separate narrow governance transition must:

- mark P00.10 `done`;
- mark P00 `done`;
- activate P01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- retire ADR-0006 from active use while retaining it as historical evidence;
- record the final branch-protection/compensating-control and LOCAL-WIN-4 CI evidence;
- define/activate the first P01 kernel work package.

Do not combine that state transition with unrelated kernel feature code.

## 8. Cleanup

After administration:

```powershell
Remove-Item Env:OMNEXA_GITHUB_ADMIN_TOKEN -ErrorAction SilentlyContinue
Remove-Item Env:GH_TOKEN -ErrorAction SilentlyContinue
```

If a short-lived PAT was used, revoke/rotate it after verification evidence is recorded.
