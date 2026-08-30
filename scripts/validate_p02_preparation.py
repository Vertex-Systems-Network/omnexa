#!/usr/bin/env python3
"""Validate P02 readiness, bounded activation, terminal completion, and later-phase retention."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/governance/P01_EXIT_GATE.md",
    "docs/governance/P02_ENTRY_GATE.md",
    "docs/governance/P02_EXIT_GATE.md",
    "docs/governance/P01_P02_TRANSITION_CHECKLIST.md",
    "docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json",
    *[f"docs/roadmap/work-packages/P02.{i:02d}.md" for i in range(1, 11)],
    "scripts/validate_p02_package_specs.py",
    ".github/workflows/governance.yml",
]
for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing P02 readiness artifact: {relative}")

state = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
manifest = json.loads((ROOT / "docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))
phase_rows = {item.get("id"): item for item in state.get("phases") or []}
lock = state.get("implementation_lock") or {}

if (phase_rows.get("P00") or {}).get("state") != "done":
    raise SystemExit("ERROR: P00 must remain done before/after P02")
if (phase_rows.get("P01") or {}).get("state") != "done":
    raise SystemExit("ERROR: P01 must remain done before/after P02")

p01_exit = (ROOT / "docs/governance/P01_EXIT_GATE.md").read_text(encoding="utf-8")
if "Status: **SATISFIED**" not in p01_exit:
    raise SystemExit("ERROR: P02 readiness requires SATISFIED P01 exit gate")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "verified_protected" or branch.get("live_protected") is not True:
    raise SystemExit("ERROR: main protection must remain verified")
if branch.get("issue_state") != "closed" or branch.get("required_check") != "governance":
    raise SystemExit("ERROR: Issue #3/governance protection evidence mismatch")
ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_github_hosted" or ci.get("routing_mode") != "github_hosted_only":
    raise SystemExit("ERROR: canonical CI must remain GitHub-hosted only")
if ci.get("runner_label") != "ubuntu-24.04" or ci.get("self_hosted_allowed") is not False:
    raise SystemExit("ERROR: P02 canonical runner policy mismatch")

current_phase = state.get("current_phase")
current = state.get("current_work_package")
phase = state.get("phase") or {}
planning = current_phase == "P01" and current is None and phase.get("state") == "done"
active = current_phase == "P02" and phase.get("state") == "active" and isinstance(current, str) and current.startswith("P02.")
terminal = current_phase == "P02" and current is None and phase.get("state") == "done"
historical = current_phase in {"P03", "P04"} and (phase_rows.get("P02") or {}).get("state") == "done"
completed = terminal or historical
if not (planning or active or completed):
    raise SystemExit("ERROR: invalid P02 readiness/activation/completion/later-phase checkpoint")

entry = (ROOT / "docs/governance/P02_ENTRY_GATE.md").read_text(encoding="utf-8")
for marker in ["P01 exit satisfied", "GitHub-hosted", "strict sequential", "kernel.identity", "kernel.tenancy", "kernel.authorization", "business_feature_code_authorized=false"]:
    if marker.lower() not in entry.lower():
        raise SystemExit(f"ERROR: P02 entry gate missing marker: {marker}")

exit_gate = (ROOT / "docs/governance/P02_EXIT_GATE.md").read_text(encoding="utf-8")
for marker in ["cross-tenant", "object/scope", "role", "service account", "session invalidation"]:
    if marker.lower() not in exit_gate.lower():
        raise SystemExit(f"ERROR: P02 exit gate missing marker: {marker}")

if planning:
    if (phase_rows.get("P02") or {}).get("state") != "planned":
        raise SystemExit("ERROR: P02 must remain planned during readiness preparation")
    if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: P02 readiness must not authorize implementation")
    if manifest.get("state") != "planned" or manifest.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: P02 readiness manifest must remain planned/unauthorized")
    if "Status: **READY — NOT ACTIVATED**" not in entry:
        raise SystemExit("ERROR: planning-mode P02 entry gate must be READY — NOT ACTIVATED")
elif active:
    expected_ids = [f"P02.{i:02d}" for i in range(1, 11)]
    if current not in expected_ids:
        raise SystemExit(f"ERROR: invalid active P02 package: {current}")
    if (phase_rows.get("P02") or {}).get("state") != "active" or (phase_rows.get("P02") or {}).get("active_work_package") != current:
        raise SystemExit("ERROR: phases[] P02 activation mismatch")
    if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: active P02 must authorize bounded kernel code and keep business code locked")
    if manifest.get("state") != "active" or manifest.get("implementation_authorized") is not True:
        raise SystemExit("ERROR: active P02 manifest must authorize implementation")
    if "Status: **SATISFIED**" not in entry:
        raise SystemExit("ERROR: active P02 requires SATISFIED entry gate")
    prep = state.get("p02_preparation") or {}
    expected_prep = {
        "state": "activated",
        "phase_state": "active",
        "next_work_package": current,
        "work_package_state": "active",
        "work_package_spec": f"docs/roadmap/work-packages/{current}.md",
        "package_sequence": "docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json",
        "prepared_spec_count": 10,
        "mandatory_spec_count": 10,
        "entry_gate": "docs/governance/P02_ENTRY_GATE.md",
        "exit_gate": "docs/governance/P02_EXIT_GATE.md",
        "transition_checklist": "docs/governance/P01_P02_TRANSITION_CHECKLIST.md",
    }
    for key, expected in expected_prep.items():
        if prep.get(key) != expected:
            raise SystemExit(f"ERROR: p02_preparation.{key} must be {expected}")
    if prep.get("blocking_gate") is not None:
        raise SystemExit("ERROR: active P02 must have no unresolved entry blocker")
else:
    p02_row = phase_rows.get("P02") or {}
    if p02_row.get("state") != "done" or p02_row.get("active_work_package") is not None:
        raise SystemExit("ERROR: completed phases[] P02 must remain done with no active package")
    if manifest.get("state") != "done" or manifest.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: completed P02 manifest must remain done/unauthorized")
    if "Status: **SATISFIED**" not in exit_gate:
        raise SystemExit("ERROR: completed P02 requires SATISFIED exit gate")
    prep = state.get("p02_preparation") or {}
    expected_prep = {
        "state": "completed",
        "phase_state": "done",
        "next_work_package": None,
        "work_package_state": "done",
        "work_package_spec": "docs/roadmap/work-packages/P02.10.md",
        "package_sequence": "docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json",
        "prepared_spec_count": 10,
        "mandatory_spec_count": 10,
        "entry_gate": "docs/governance/P02_ENTRY_GATE.md",
        "exit_gate": "docs/governance/P02_EXIT_GATE.md",
        "transition_checklist": "docs/governance/P01_P02_TRANSITION_CHECKLIST.md",
    }
    for key, expected in expected_prep.items():
        if prep.get(key) != expected:
            raise SystemExit(f"ERROR: completed p02_preparation.{key} must remain {expected}")
    if prep.get("blocking_gate") is not None:
        raise SystemExit("ERROR: completed P02 must retain no unresolved P02 blocker")
    if terminal:
        if (phase_rows.get("P03") or {}).get("state") != "planned":
            raise SystemExit("ERROR: P03 must remain planned until a separate governed activation")
        if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
            raise SystemExit("ERROR: terminal P02 checkpoint must lock kernel and business implementation")
    elif lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: business-feature implementation must remain locked during P03/P04")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
for marker in [
    "runs-on: ubuntu-24.04",
    "RUNNER_ENVIRONMENT",
    "github-hosted",
    "python scripts/validate_p02_preparation.py",
    "python scripts/validate_p02_package_specs.py",
    "bash scripts/verify_go_quality.sh",
    "bash scripts/verify_p01_12.sh",
]:
    if marker not in workflow:
        raise SystemExit(f"ERROR: governance workflow missing P02 readiness marker: {marker}")
if "self-hosted" in workflow or "LOCAL-WIN-" in workflow:
    raise SystemExit("ERROR: local/self-hosted governance runners are prohibited")

mode = "PLANNING / NOT ACTIVATED" if planning else "ACTIVE" if active else "COMPLETED / HISTORICAL"
print("Omnexa P02 preparation/readiness validation: PASS")
print("P01 exit: SATISFIED")
print("P02 specifications: 10 / 10")
print(f"Mode: {mode}")
print(f"Active P02 package: {current if active else 'NONE'}")
print("Governance runner: GITHUB-HOSTED ONLY / ubuntu-24.04")
print(f"Kernel code: {'AUTHORIZED FOR ' + current if active else 'P02 LOCKED / HISTORICAL'}")
print("Business feature code: LOCKED")