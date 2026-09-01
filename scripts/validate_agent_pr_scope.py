#!/usr/bin/env python3
"""Enforce declared write/shared path budgets for active worker PRs."""

from __future__ import annotations

import os

from agent_orchestration_common import allowed_path, changed_paths, fail, find_branch_record, load_plan, require_sha


def main() -> None:
    if os.environ.get("GITHUB_EVENT_NAME") != "pull_request":
        print("PASS: agent PR scope check not applicable to non-PR event")
        return

    branch = os.environ.get("OMNEXA_PR_HEAD_REF") or os.environ.get("GITHUB_HEAD_REF") or ""
    plan = load_plan()
    record = find_branch_record(plan, branch)
    if record is None:
        if branch.startswith("agent/"):
            fail(f"agent branch is not registered in ACTIVE_MULTI_AGENT_PLAN.json: {branch}")
        print(f"PASS: control/non-agent branch is outside active worker scope enforcement: {branch}")
        return

    base_sha = require_sha(os.environ.get("OMNEXA_PR_BASE_SHA", ""), "PR base SHA")
    head_sha = require_sha(os.environ.get("OMNEXA_PR_HEAD_SHA", ""), "PR head SHA")
    paths = changed_paths(base_sha, head_sha)
    write_paths = list(record.get("write_paths") or [])
    shared_paths = list(record.get("shared_paths") or [])
    violations = [path for path in paths if not allowed_path(path, write_paths, shared_paths)]
    if violations:
        label = record.get("task_id") or record.get("agent_id") or branch
        fail(f"{label} PR writes outside declared budget: {', '.join(violations)}")
    print(f"PASS: {branch} changed {len(paths)} paths inside its declared write/shared budget")


if __name__ == "__main__":
    main()
