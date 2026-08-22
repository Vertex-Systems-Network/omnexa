#!/usr/bin/env python3
"""Validate the activated P01 handoff and bounded P01.01 authorization."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/roadmap/work-packages/P01.01.md",
    "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json",
    "docs/governance/P00_P01_TRANSITION_CHECKLIST.md",
    "docs/governance/P01_ENTRY_GATE.md",
    "docs/roadmap/STATE.json",
    ".github/workflows/governance.yml",
]
for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing P01 activation artifact: {relative}")

state = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
if state.get("current_phase") != "P01" or state.get("current_work_package") != "P01.01":
    raise SystemExit("ERROR: active handoff must be P01 / P01.01")

lock = state.get("implementation_lock") or {}
if lock.get("kernel_code_authorized") is not True:
    raise SystemExit("ERROR: kernel code must be authorized for active P01")
if lock.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: business feature code must remain locked")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "verified_protected" or branch.get("live_protected") is not True:
    raise SystemExit("ERROR: Issue #3 protection evidence must be verified before P01")
if branch.get("issue_state") != "closed" or branch.get("required_check") != "governance":
    raise SystemExit("ERROR: Issue #3 must be closed and governance must be the required check")

entry = tracking.get("p01_entry_gate") or {}
if entry.get("state") != "SATISFIED" or entry.get("blockers") != []:
    raise SystemExit("ERROR: P01 entry gate must be SATISFIED with no blockers")

ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_github_hosted" or ci.get("routing_mode") != "github_hosted_only":
    raise SystemExit("ERROR: governance CI must remain GitHub-hosted only")
if ci.get("runner_label") != "ubuntu-24.04" or ci.get("self_hosted_allowed") is not False:
    raise SystemExit("ERROR: canonical runner policy mismatch")

phases = {item.get("id"): item for item in state.get("phases") or []}
if (phases.get("P00") or {}).get("state") != "done" or (phases.get("P01") or {}).get("state") != "active":
    raise SystemExit("ERROR: P00 must be done and P01 active")

prep = state.get("p01_preparation") or {}
for key, expected in {
    "state": "activated",
    "phase_state": "active",
    "next_work_package": "P01.01",
    "work_package_state": "active",
    "work_package_spec": "docs/roadmap/work-packages/P01.01.md",
    "package_sequence": "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json",
}.items():
    if prep.get(key) != expected:
        raise SystemExit(f"ERROR: p01_preparation.{key} must be {expected}")
if prep.get("prepared_spec_count") != 12 or prep.get("mandatory_spec_count") != 12:
    raise SystemExit("ERROR: P01 must retain 12 / 12 prepared specifications")
if prep.get("blocking_gate") is not None:
    raise SystemExit("ERROR: P01 activation must have no blocking gate")

package = (ROOT / "docs/roadmap/work-packages/P01.01.md").read_text(encoding="utf-8")
for marker in ["# P01.01 — Go Workspace / Build Skeleton", "State: `active`", "kernel_code_authorized=true", "business_feature_code_authorized=false", "P01.02", "P02", "P03", "G0", "G1", "G2", "G7", "GitHub-hosted"]:
    if marker.lower() not in package.lower():
        raise SystemExit(f"ERROR: active P01.01 spec missing marker: {marker}")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
for marker in ["name: governance", "runs-on: ubuntu-24.04", "RUNNER_ENVIRONMENT", "github-hosted", "python scripts/validate_governance.py", "python scripts/validate_freeze_review.py", "python scripts/validate_p01_preparation.py", "python scripts/validate_p01_package_specs.py"]:
    if marker not in workflow:
        raise SystemExit(f"ERROR: governance workflow missing marker: {marker}")
if "self-hosted" in workflow or "LOCAL-WIN-" in workflow:
    raise SystemExit("ERROR: local/self-hosted governance runners are prohibited")

print("Omnexa P01 activation/readiness validation: PASS")
print("P00: DONE")
print("P01.01: ACTIVE")
print("Governance runner: GITHUB-HOSTED ONLY / ubuntu-24.04")
print("Kernel code: AUTHORIZED FOR P01.01")
print("Business feature code: LOCKED")
