# Omnexa Main Branch Protection Administration Runbook

Status: **P00.10 exit-verification tooling — actionable**  
Tracker: **Issue #3**

This runbook is the deterministic administration path for applying and verifying the hosted GitHub `main` protection required by `REPOSITORY_HARDENING.md` and `P01_ENTRY_GATE.md`.

It does not authorize P01 by itself. P01 remains locked until the hosted setting is applied and verified, or EG-02 is deliberately superseded through accepted governance change control.

## 1. Current repository state

The repository is now public. The earlier private-repository attempt returned:

```text
gh: Upgrade to GitHub Pro or make this repository public to enable this feature. (HTTP 403)
```

That HTTP 403 remains historical evidence only. The former private-repository plan blocker is no longer the active reason EG-02 is blocked.

Current verified facts:

- repository: `Vertex-Systems-Network/omnexa`;
- owner type: organization;
- repository is now public;
- authenticated account has repository `admin` permission;
- canonical governance CI is GitHub-hosted only on `ubuntu-24.04`;
- hosted CI proof: run `32537207455`, job `96940269306`, `RUNNER_ENVIRONMENT=github-hosted`;
- live `main` metadata still reports `protected: false`.

Therefore the remaining task is configuration + verification, not plan troubleshooting.

## 2. Runner policy

All repository administration automation for this gate must use GitHub-hosted compute. Local/self-hosted Actions runners are prohibited.

The canonical validation workflow uses:

```yaml
runs-on: ubuntu-24.04
```

If administration is executed through Actions, use the separate manual GitHub-hosted administration workflow and an owner-controlled short-lived repository Administration token. The ordinary Actions `GITHUB_TOKEN` does not provide repository Administration authority and must not be treated as a substitute.

## 3. Security model

Use a short-lived fine-grained PAT exposed only as `OMNEXA_GITHUB_ADMIN_TOKEN`.

For the token:

- select only `Vertex-Systems-Network/omnexa`;
- grant repository **Administration: read and write**;
- keep the lifetime short;
- store it only as a protected GitHub Actions secret for the administration run;
- do not commit, echo, paste into issue comments or repository files;
- revoke/rotate it immediately after the hosted setting and evidence are complete.

No local workstation credential path is part of the canonical procedure anymore.

## 4. Required protection policy

Apply the policy encoded by `scripts/apply_main_protection.ps1`:

- branch: `main`;
- PR-based integration required;
- required status check: `governance`;
- required status check strict/up-to-date;
- stale approvals dismissed;
- conversation resolution required;
- force pushes blocked;
- branch deletion blocked;
- protection enforced for administrators;
- required approval count: `0` while the repository has a single independent maintainer;
- Code Owner review not required until an independent reviewer exists.

The approval count of `0` is intentional: PR-only integration remains enforced without creating an impossible self-approval deadlock. When a second independent reviewer is available, increase the required approvals and enable Code Owner review.

## 5. GitHub-hosted administration workflow

The repository contains `.github/workflows/main-protection-admin.yml` as a `workflow_dispatch`-only administrative workflow. It must:

1. run on `ubuntu-24.04`;
2. verify `RUNNER_ENVIRONMENT=github-hosted`;
3. require `secrets.OMNEXA_GITHUB_ADMIN_TOKEN`;
4. execute `scripts/apply_main_protection.ps1` using `pwsh`;
5. execute `scripts/verify_main_protection.ps1` using `pwsh`;
6. fail closed on any missing protection property.

Do not add push/pull-request triggers to that privileged workflow.

## 6. Verify live GitHub configuration

The verifier must finish with:

```text
Omnexa main branch protection verification: PASS
```

Verification must prove at least:

1. GitHub branch metadata reports `protected: true`.
2. Required check `governance` is present and strict/up-to-date.
3. Conversation resolution is required.
4. Force pushes are disabled.
5. Deletion of `main` is disabled.
6. Administrator enforcement is enabled.
7. Current review policy matches the contributor model.

## 7. Required Issue #3 negative evidence

Before Issue #3 can close, record controlled evidence for:

1. a non-exempt direct push to `main` being rejected;
2. a force-push attempt being rejected;
3. a branch-deletion attempt being rejected;
4. a test PR being unable to merge while `governance` fails/is incomplete;
5. CODEOWNERS/review behavior matching the current contributor model.

Negative tests must use a controlled non-break-glass actor or a safe test procedure. Do not risk an irreversible `main` mutation merely to prove rejection.

## 8. P00 -> P01 transition after verification

Once Issue #3 evidence is complete, or EG-02 has been deliberately replaced by an accepted superseding governance ADR, a separate narrow governance transition must:

- mark P00.10 `done`;
- mark P00 `done`;
- activate P01 and P01.01;
- set `kernel_code_authorized = true`;
- keep `business_feature_code_authorized = false`;
- retire ADR-0006 from active use while retaining it as historical evidence;
- record final branch-protection and GitHub-hosted CI evidence;
- keep P01.02–P01.12 planned.

Do not combine that state transition with kernel implementation.

## 9. Cleanup

After the administration run, revoke/rotate the short-lived PAT and remove the `OMNEXA_GITHUB_ADMIN_TOKEN` repository secret if it is no longer needed.
