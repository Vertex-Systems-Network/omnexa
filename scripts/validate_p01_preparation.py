#!/usr/bin/env python3
"""Validate activated P01 handoff, entry controls, and bounded current-package authorization."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
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
if state.get("current_phase") != "P01":
    raise SystemExit("ERROR: current phase must be P01")

current = state.get("current_work_package")
expected_ids = [f"P01.{i:02d}" for i in range(1, 13)]
if current not in expected_ids:
    raise SystemExit(f"ERROR: invalid active P01 package: {current}")
current_index = expected_ids.index(current)
active_spec = ROOT / f"docs/roadmap/work-packages/{current}.md"
if not active_spec.is_file():
    raise SystemExit(f"ERROR: missing active P01 specification: {active_spec.relative_to(ROOT)}")

lock = state.get("implementation_lock") or {}
if lock.get("kernel_code_authorized") is not True:
    raise SystemExit("ERROR: kernel code must be authorized for active P01")
if lock.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: business feature code must remain locked")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "verified_protected" or branch.get("live_protected") is not True:
    raise SystemExit("ERROR: Issue #3 protection evidence must remain verified")
if branch.get("issue_state") != "closed" or branch.get("required_check") != "governance":
    raise SystemExit("ERROR: Issue #3 must remain closed and governance must remain the required check")

entry = tracking.get("p01_entry_gate") or {}
if entry.get("state") != "SATISFIED" or entry.get("blockers") != []:
    raise SystemExit("ERROR: P01 entry gate must remain SATISFIED with no blockers")

ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_github_hosted" or ci.get("routing_mode") != "github_hosted_only":
    raise SystemExit("ERROR: governance CI must remain GitHub-hosted only")
if ci.get("runner_label") != "ubuntu-24.04" or ci.get("self_hosted_allowed") is not False:
    raise SystemExit("ERROR: canonical runner policy mismatch")

phases = {item.get("id"): item for item in state.get("phases") or []}
if (phases.get("P00") or {}).get("state") != "done" or (phases.get("P01") or {}).get("state") != "active":
    raise SystemExit("ERROR: P00 must be done and P01 active")
if (phases.get("P01") or {}).get("active_work_package") != current:
    raise SystemExit("ERROR: phases[] P01 active_work_package mismatch")

prep = state.get("p01_preparation") or {}
for key, expected in {
    "state": "activated",
    "phase_state": "active",
    "next_work_package": current,
    "work_package_state": "active",
    "work_package_spec": f"docs/roadmap/work-packages/{current}.md",
    "package_sequence": "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json",
}.items():
    if prep.get(key) != expected:
        raise SystemExit(f"ERROR: p01_preparation.{key} must be {expected}")
if prep.get("prepared_spec_count") != 12 or prep.get("mandatory_spec_count") != 12:
    raise SystemExit("ERROR: P01 must retain 12 / 12 prepared specifications")
if prep.get("blocking_gate") is not None:
    raise SystemExit("ERROR: active P01 must have no blocking entry gate")

package = active_spec.read_text(encoding="utf-8")
for marker in [
    current,
    "State: `active`",
    "kernel_code_authorized=true",
    "business_feature_code_authorized=false",
    "Acceptance criteria",
    "Completion evidence",
    "GitHub-hosted",
]:
    if marker.lower() not in package.lower():
        raise SystemExit(f"ERROR: active {current} spec missing marker: {marker}")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
base_workflow_markers = [
    "name: governance",
    "runs-on: ubuntu-24.04",
    "RUNNER_ENVIRONMENT",
    "github-hosted",
    "python scripts/validate_governance.py",
    "python scripts/validate_freeze_review.py",
    "python scripts/validate_p01_preparation.py",
    "python scripts/validate_p01_package_specs.py",
    "bash scripts/verify_go_quality.sh",
]
for marker in base_workflow_markers:
    if marker not in workflow:
        raise SystemExit(f"ERROR: governance workflow missing marker: {marker}")
if "self-hosted" in workflow or "LOCAL-WIN-" in workflow:
    raise SystemExit("ERROR: local/self-hosted governance runners are prohibited")

# Every completed sequential package must retain immutable completion evidence and
# its regression verifier in both the repository and canonical workflow.
for index in range(current_index):
    package_number = index + 1
    evidence = ROOT / f"docs/roadmap/evidence/P01.{package_number:02d}_COMPLETION_2026-08-22.md"
    verifier = ROOT / f"scripts/verify_p01_{package_number:02d}.sh"
    if not evidence.is_file():
        raise SystemExit(f"ERROR: completed P01.{package_number:02d} missing canonical completion evidence")
    if not verifier.is_file():
        raise SystemExit(f"ERROR: completed P01.{package_number:02d} missing regression verifier")
    workflow_marker = f"bash scripts/verify_p01_{package_number:02d}.sh"
    if workflow_marker not in workflow:
        raise SystemExit(f"ERROR: completed P01.{package_number:02d} verifier missing from governance workflow")

# Active executable packages with implemented verification must fail closed in
# canonical CI. This condition is intentionally bounded to the current package
# and drops when governance advances to the next package.
if current in {"P01.03", "P01.04", "P01.05"}:
    package_number = current.split(".")[1]
    active_verifier = ROOT / f"scripts/verify_p01_{package_number}.sh"
    if not active_verifier.is_file():
        raise SystemExit(f"ERROR: active {current} verifier is missing")
    if f"bash scripts/verify_p01_{package_number}.sh" not in workflow:
        raise SystemExit(f"ERROR: active {current} verifier is not enforced by governance workflow")

print("Omnexa P01 activation/readiness validation: PASS")
print("P00: DONE")
print(f"Completed P01 packages: {current_index} / 12")
print(f"Active package: {current}")
print("Governance runner: GITHUB-HOSTED ONLY / ubuntu-24.04")
print(f"Kernel code: AUTHORIZED FOR {current}")
print("Business feature code: LOCKED")
