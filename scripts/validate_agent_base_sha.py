#!/usr/bin/env python3
"""Fail active worker PRs that are based on stale protected main."""

from __future__ import annotations

import subprocess

from agent_orchestration_common import ROOT, ensure_commit, fail, find_branch_record, load_plan, pr_context, remote_main_sha


def main() -> None:
    context = pr_context()
    if context is None:
        print("PASS: agent base-SHA check not applicable to non-PR event")
        return

    branch, base_sha, head_sha = context
    record = find_branch_record(load_plan(), branch)
    if record is None:
        if branch.startswith("agent/"):
            fail(f"agent branch is not registered in ACTIVE_MULTI_AGENT_PLAN.json: {branch}")
        print(f"PASS: control/non-agent branch is outside active worker freshness enforcement: {branch}")
        return

    live_main = remote_main_sha()
    if base_sha != live_main:
        fail(f"stale worker PR base: PR base {base_sha} != live protected main {live_main}; sync main and resubmit")

    ensure_commit(base_sha)
    ensure_commit(head_sha)
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
