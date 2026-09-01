#!/usr/bin/env python3
"""Shared dependency-free helpers for Omnexa multi-agent CI enforcement."""

from __future__ import annotations

import fnmatch
import json
import os
import pathlib
import re
import subprocess
import sys
from typing import Any

ROOT = pathlib.Path(__file__).resolve().parents[1]
PLAN_PATH = ROOT / "docs/ai/ACTIVE_MULTI_AGENT_PLAN.json"
STATE_PATH = ROOT / "docs/roadmap/STATE.json"
SHA_RE = re.compile(r"^[0-9a-f]{40}$")


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


def load_json(path: pathlib.Path, label: str) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        fail(f"{label} is unreadable/invalid JSON: {exc}")
    if not isinstance(value, dict):
        fail(f"{label} root must be an object")
    return value


def load_plan() -> dict[str, Any]:
    return load_json(PLAN_PATH, "ACTIVE_MULTI_AGENT_PLAN.json")


def load_state() -> dict[str, Any]:
    return load_json(STATE_PATH, "STATE.json")


def require_sha(value: str, label: str) -> str:
    if not SHA_RE.fullmatch(value or ""):
        fail(f"{label} must be a 40-character lowercase Git SHA")
    return value


def pr_context() -> tuple[str, str, str] | None:
    if os.environ.get("GITHUB_EVENT_NAME") != "pull_request":
        return None
    branch = os.environ.get("OMNEXA_PR_HEAD_REF") or os.environ.get("GITHUB_HEAD_REF") or ""
    base_sha = os.environ.get("OMNEXA_PR_BASE_SHA") or ""
    head_sha = os.environ.get("OMNEXA_PR_HEAD_SHA") or ""
    event_path = os.environ.get("GITHUB_EVENT_PATH")
    if event_path and (not branch or not base_sha or not head_sha):
        event = load_json(pathlib.Path(event_path), "GitHub event payload")
        pr = event.get("pull_request") or {}
        branch = branch or ((pr.get("head") or {}).get("ref") or "")
        base_sha = base_sha or ((pr.get("base") or {}).get("sha") or "")
        head_sha = head_sha or ((pr.get("head") or {}).get("sha") or "")
    if not branch:
        fail("pull_request event is missing head branch identity")
    return branch, require_sha(base_sha, "PR base SHA"), require_sha(head_sha, "PR head SHA")


def normalize_path(value: str) -> str:
    return value.strip().replace("\\", "/").lstrip("./")


def _prefix_pattern(pattern: str) -> str | None:
    pattern = normalize_path(pattern)
    if pattern.endswith("/**"):
        return pattern[:-3].rstrip("/")
    if pattern.endswith("/*"):
        return pattern[:-2].rstrip("/")
    return None


def path_matches(pattern: str, path: str) -> bool:
    pattern = normalize_path(pattern)
    path = normalize_path(path)
    prefix = _prefix_pattern(pattern)
    if prefix is not None:
        return path == prefix or path.startswith(prefix + "/")
    if any(ch in pattern for ch in "*?["):
        return fnmatch.fnmatchcase(path, pattern)
    return path == pattern


def patterns_overlap(left: str, right: str) -> bool:
    left = normalize_path(left)
    right = normalize_path(right)
    if left == right:
        return True
    lp = _prefix_pattern(left)
    rp = _prefix_pattern(right)
    if lp is not None and rp is not None:
        return lp == rp or lp.startswith(rp + "/") or rp.startswith(lp + "/")
    if lp is not None:
        return right == lp or right.startswith(lp + "/") or fnmatch.fnmatchcase(lp, right)
    if rp is not None:
        return left == rp or left.startswith(rp + "/") or fnmatch.fnmatchcase(rp, left)
    if any(ch in left for ch in "*?[") or any(ch in right for ch in "*?["):
        return fnmatch.fnmatchcase(left, right) or fnmatch.fnmatchcase(right, left)
    return False


def task_records(plan: dict[str, Any]) -> list[dict[str, Any]]:
    tasks = plan.get("tasks") or []
    if not isinstance(tasks, list):
        fail("active plan tasks must be an array")
    for task in tasks:
        if not isinstance(task, dict):
            fail("every active plan task must be an object")
    return tasks


def slot_records(plan: dict[str, Any]) -> list[dict[str, Any]]:
    slots = plan.get("worker_slots") or []
    if not isinstance(slots, list):
        fail("active plan worker_slots must be an array")
    for slot in slots:
        if not isinstance(slot, dict):
            fail("every worker slot must be an object")
    return slots


def all_write_records(plan: dict[str, Any]) -> list[dict[str, Any]]:
    records = list(task_records(plan))
    supervisor = plan.get("supervisor")
    if isinstance(supervisor, dict):
        records.append(supervisor)
    return records


def find_branch_record(plan: dict[str, Any], branch: str) -> dict[str, Any] | None:
    for record in all_write_records(plan):
        if record.get("branch") == branch:
            return record
    return None


def git(*args: str, check: bool = True) -> str:
    result = subprocess.run(
        ["git", *args],
        cwd=ROOT,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    if check and result.returncode != 0:
        fail(f"git {' '.join(args)} failed: {result.stderr.strip()}")
    return result.stdout.strip()


def ensure_commit(sha: str) -> None:
    require_sha(sha, "commit SHA")
    result = subprocess.run(
        ["git", "cat-file", "-e", f"{sha}^{{commit}}"],
        cwd=ROOT,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    if result.returncode == 0:
        return
    git("fetch", "--no-tags", "--depth=1", "origin", sha)


def remote_main_sha() -> str:
    output = git("ls-remote", "origin", "refs/heads/main")
    if not output:
        fail("could not resolve protected main from origin")
    sha = output.split()[0]
    return require_sha(sha, "origin/main SHA")


def changed_paths(base_sha: str, head_sha: str) -> list[str]:
    ensure_commit(base_sha)
    ensure_commit(head_sha)
    output = git("diff", "--name-only", f"{base_sha}...{head_sha}")
    return [normalize_path(line) for line in output.splitlines() if line.strip()]


def allowed_path(path: str, write_paths: list[str], shared_paths: list[str] | None = None) -> bool:
    patterns = list(write_paths)
    if shared_paths:
        patterns.extend(shared_paths)
    return any(path_matches(pattern, path) for pattern in patterns)
