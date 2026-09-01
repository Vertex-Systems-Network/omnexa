#!/usr/bin/env python3
"""Fail active worker PRs that are based on stale protected main."""

from __future__ import annotations

import os

from agent_orchestration_common import ensure_commit, fail, find_branch_record, git, load_plan, remote_main_sha, require_sha


def main() -> None:
    if os.environ.get("GITHUB_EVENT_NAME") != "pull_request":
        print("PASS: agent base-SHA check not applicable to non-PR event")
        return

    branch = os.environ.get("OMNEXA_PR_HEAD_REF") or os.environ.get("GITHUB_HEAD_REF") or ""
    record = find_branch_record(load_plan(), branch)
    if record is None:
        if branch.startswith("agent/"):
            fail(f"agent branch is not registered in ACTIVE_MULTI_AGENT_PLAN.json: {branch}")
        print(f"PASS: control/non-agent branch is outside active worker freshness enforcement: {branch}")
        return

    base_sha = require_sha(os.environ.get("OMNEXA_PR_BASE_SHA", ""), "PR base SHA")
    head_sha = require_sha(os.environ.get("OMNEXA_PR_HEAD_SHA", ""), "PR head SHA")
    live_main = remote_main_sha()
    if base_sha != live_main:
        fail(f"stale worker PR base: PR base {base_sha} != live protected main {live_main}; sync main and resubmit")

    ensure_commit(base_sha)
    ensure_commit(head_sha)
    result = git("merge-base", "--is-ancestor", base_sha, head_sha, check=False)
    # git merge-base --is-ancestor prints nothing; inspect return code separately for correctness.
    import subprocess
    from agent_orchestration_common import ROOT

    proc = subprocess.run(
        ["git", "merge-base", "--is-ancestor", base_sha, head_sha],
        cwd=ROOT,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    if proc.returncode != 0:
        fail("worker PR head does not descend from the current protected-main base")

    print(f"PASS: {branch} is based on current protected main {live_main}")


if __name__ == "__main__":
    main()
