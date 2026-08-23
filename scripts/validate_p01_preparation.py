#!/usr/bin/env python3
"""Validate completed P01 evidence as a permanent prerequisite for later phases."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json",
    "docs/governance/P00_P01_TRANSITION_CHECKLIST.md",
    "docs/governance/P01_ENTRY_GATE.md",
    "docs/governance/P01_EXIT_GATE.md",
    "docs/roadmap/STATE.json",
    ".github/workflows/governance.yml",
]
for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing completed-P01 artifact: {relative}")

state = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
manifest = json.loads((ROOT / "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))
phases = {item.get("id"): item for item in state.get("phases") or []}

if state.get("current_phase") not in {"P01", "P02"}:
    raise SystemExit("ERROR: completed P01 validator must be reviewed before advancing beyond P02")
if (phases.get("P00") or {}).get("state") != "done":
    raise SystemExit("ERROR: P00 must remain done")
p01_row = phases.get("P01") or {}
if p01_row.get("state") != "done" or p01_row.get("active_work_package") is not None:
    raise SystemExit("ERROR: phases[] P01 must remain done with no active work package")

if manifest.get("phase") != "P01" or manifest.get("state") != "done":
    raise SystemExit("ERROR: P01 package manifest must remain done")
if manifest.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P01 activation policy changed")
if manifest.get("implementation_authorized") is not False:
    raise SystemExit("ERROR: completed P01 manifest must keep implementation_authorized=false")

expected_ids = [f"P01.{i:02d}" for i in range(1, 13)]
packages = manifest.get("packages") or []
if [item.get("id") for item in packages] != expected_ids:
    raise SystemExit("ERROR: P01 package order must remain P01.01-P01.12")
if any(item.get("state") != "done" for item in packages):
    raise SystemExit("ERROR: every P01 package must remain done")

prep = state.get("p01_preparation") or {}
expected_prep = {
    "state": "completed",
    "phase_state": "done",
    "next_work_package": None,
    "work_package_state": "done",
    "work_package_spec": "docs/roadmap/work-packages/P01.12.md",
    "package_sequence": "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json",
    "prepared_spec_count": 12,
    "mandatory_spec_count": 12,
}
for key, expected in expected_prep.items():
    if prep.get(key) != expected:
        raise SystemExit(f"ERROR: p01_preparation.{key} must remain {expected}")
if prep.get("blocking_gate") is not None:
    raise SystemExit("ERROR: completed P01 must retain no unresolved blocking gate")

exit_gate = (ROOT / "docs/governance/P01_EXIT_GATE.md").read_text(encoding="utf-8")
if "Status: **SATISFIED**" not in exit_gate:
    raise SystemExit("ERROR: P01 exit gate must remain SATISFIED")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "verified_protected" or branch.get("live_protected") is not True:
    raise SystemExit("ERROR: Issue #3 protection evidence must remain verified")
if branch.get("issue_state") != "closed" or branch.get("required_check") != "governance":
    raise SystemExit("ERROR: protected integration policy mismatch")
ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_github_hosted" or ci.get("routing_mode") != "github_hosted_only":
    raise SystemExit("ERROR: canonical governance CI must remain GitHub-hosted only")
if ci.get("runner_label") != "ubuntu-24.04" or ci.get("self_hosted_allowed") is not False:
    raise SystemExit("ERROR: canonical runner policy mismatch")

completion_evidence = {
    **{f"P01.{number:02d}": f"docs/roadmap/evidence/P01.{number:02d}_COMPLETION_2026-08-22.md" for number in range(1, 10)},
    "P01.10": "docs/roadmap/evidence/P01.10_COMPLETION_2026-08-23.md",
    "P01.11": "docs/roadmap/evidence/P01.11_COMPLETION_2026-08-23.md",
    "P01.12": "docs/roadmap/evidence/P01.12_COMPLETION_2026-08-23.md",
}
workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
for number, package_id in enumerate(expected_ids, start=1):
    evidence = ROOT / completion_evidence[package_id]
    verifier = ROOT / f"scripts/verify_p01_{number:02d}.sh"
    if not evidence.is_file():
        raise SystemExit(f"ERROR: completed {package_id} missing canonical completion evidence")
    if not verifier.is_file():
        raise SystemExit(f"ERROR: completed {package_id} missing regression verifier")
    if f"bash scripts/verify_p01_{number:02d}.sh" not in workflow:
        raise SystemExit(f"ERROR: completed {package_id} verifier missing from governance workflow")

for marker in [
    "name: governance",
    "runs-on: ubuntu-24.04",
    "RUNNER_ENVIRONMENT",
    "github-hosted",
    "python scripts/validate_governance.py",
    "python scripts/validate_freeze_review.py",
    "bash scripts/verify_go_quality.sh",
]:
    if marker not in workflow:
        raise SystemExit(f"ERROR: governance workflow missing marker: {marker}")
if "self-hosted" in workflow or "LOCAL-WIN-" in workflow:
    raise SystemExit("ERROR: local/self-hosted governance runners are prohibited")

if state.get("current_phase") == "P01":
    lock = state.get("implementation_lock") or {}
    if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: completed P01 planning checkpoint must keep implementation locked")
else:
    if (state.get("implementation_lock") or {}).get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: business-feature implementation must remain locked during P02")

print("Omnexa completed P01 prerequisite validation: PASS")
print("P00: DONE")
print("P01: DONE / 12 OF 12")
print("P01 exit: SATISFIED")
print("P01 regressions: ENFORCED IN GOVERNANCE")
print("Governance runner: GITHUB-HOSTED ONLY / ubuntu-24.04")
print("Business feature code: LOCKED")
