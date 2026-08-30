#!/usr/bin/env python3
"""Validate P03 readiness, bounded activation, terminal completion, and P04-era historical retention."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/governance/P02_EXIT_GATE.md",
    "docs/governance/P03_ENTRY_GATE.md",
    "docs/governance/P03_EXIT_GATE.md",
    "docs/governance/P02_P03_TRANSITION_CHECKLIST.md",
    "docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md",
    "docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json",
    *[f"docs/roadmap/work-packages/P03.{i:02d}.md" for i in range(1, 12)],
    "docs/ai/handoffs/P03.01.md",
    "scripts/validate_p03_package_specs.py",
    ".github/workflows/governance.yml",
]
for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing P03 readiness artifact: {relative}")

state = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
manifest = json.loads((ROOT / "docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))
phase_rows = {item.get("id"): item for item in state.get("phases") or []}
lock = state.get("implementation_lock") or {}

for completed in ["P00", "P01", "P02"]:
    if (phase_rows.get(completed) or {}).get("state") != "done":
        raise SystemExit(f"ERROR: {completed} must remain done before/after P03")

p02_exit = (ROOT / "docs/governance/P02_EXIT_GATE.md").read_text(encoding="utf-8")
if "Status: **SATISFIED**" not in p02_exit:
    raise SystemExit("ERROR: P03 readiness requires SATISFIED P02 exit gate")

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
    raise SystemExit("ERROR: P03 canonical runner policy mismatch")

current_phase = state.get("current_phase")
current = state.get("current_work_package")
phase = state.get("phase") or {}
planning = current_phase == "P02" and current is None and phase.get("id") == "P02" and phase.get("state") == "done"
active = current_phase == "P03" and phase.get("id") == "P03" and phase.get("state") == "active" and isinstance(current, str) and current.startswith("P03.")
terminal = current_phase == "P03" and current is None and phase.get("id") == "P03" and phase.get("state") == "done"
historical = current_phase == "P04" and phase.get("id") == "P04" and phase.get("state") == "active" and isinstance(current, str) and current.startswith("P04.")
if not (planning or active or terminal or historical):
    raise SystemExit("ERROR: invalid P03 readiness/activation/completion/historical checkpoint")

entry = (ROOT / "docs/governance/P03_ENTRY_GATE.md").read_text(encoding="utf-8")
for marker in [
    "P02 exit satisfied",
    "GitHub-hosted",
    "strict sequential",
    "kernel.modules",
    "required, optional, platform and forbidden",
    "business_feature_code_authorized=false",
]:
    if marker.lower() not in entry.lower():
        raise SystemExit(f"ERROR: P03 entry gate missing marker: {marker}")

exit_gate = (ROOT / "docs/governance/P03_EXIT_GATE.md").read_text(encoding="utf-8")
for marker in [
    "required dependency enforcement",
    "optional dependency degradation",
    "disable and re-enable",
    "upgrade and migration path",
    "forbidden dependency detection",
    "health and state accuracy",
    "unrelated module isolation",
]:
    if marker.lower() not in exit_gate.lower():
        raise SystemExit(f"ERROR: P03 exit gate missing marker: {marker}")

alignment = (ROOT / "docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md").read_text(encoding="utf-8")
for marker in ["XQ-100", "XSG-100", "XTRUST-100", "XPF-200", "XPERF-100", "implementation_authorized=false"]:
    if marker.lower() not in alignment.lower():
        raise SystemExit(f"ERROR: P03 AI-native alignment missing marker: {marker}")

if planning:
    if (phase_rows.get("P03") or {}).get("state") != "planned":
        raise SystemExit("ERROR: P03 must remain planned during readiness preparation")
    if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: P03 readiness must not authorize implementation")
    if manifest.get("state") != "planned" or manifest.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: P03 readiness manifest must remain planned/unauthorized")
    if "Status: **READY — NOT ACTIVATED**" not in entry:
        raise SystemExit("ERROR: planning-mode P03 entry gate must be READY — NOT ACTIVATED")
elif active:
    expected_ids = [f"P03.{i:02d}" for i in range(1, 12)]
    if current not in expected_ids:
        raise SystemExit(f"ERROR: invalid active P03 package: {current}")
    if (phase_rows.get("P03") or {}).get("state") != "active" or (phase_rows.get("P03") or {}).get("active_work_package") != current:
        raise SystemExit("ERROR: phases[] P03 activation mismatch")
    if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: active P03 must authorize bounded kernel code and keep business code locked")
    if manifest.get("state") != "active" or manifest.get("implementation_authorized") is not True:
        raise SystemExit("ERROR: active P03 manifest must authorize implementation")
    if "Status: **SATISFIED**" not in entry:
        raise SystemExit("ERROR: active P03 requires SATISFIED entry gate")
    prep = state.get("p03_preparation") or {}
    expected_prep = {
        "state": "activated",
        "phase_state": "active",
        "next_work_package": current,
        "work_package_state": "active",
        "work_package_spec": f"docs/roadmap/work-packages/{current}.md",
        "package_sequence": "docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json",
        "prepared_spec_count": 11,
        "mandatory_spec_count": 11,
        "entry_gate": "docs/governance/P03_ENTRY_GATE.md",
        "exit_gate": "docs/governance/P03_EXIT_GATE.md",
        "transition_checklist": "docs/governance/P02_P03_TRANSITION_CHECKLIST.md",
        "ai_native_alignment": "docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md",
    }
    for key, expected in expected_prep.items():
        if prep.get(key) != expected:
            raise SystemExit(f"ERROR: p03_preparation.{key} must be {expected}")
    if prep.get("blocking_gate") is not None:
        raise SystemExit("ERROR: active P03 must have no unresolved entry blocker")
else:
    if (phase_rows.get("P03") or {}).get("state") != "done" or (phase_rows.get("P03") or {}).get("active_work_package") is not None:
        raise SystemExit("ERROR: completed P03 must remain done with no active P03 package")
    if manifest.get("state") != "done" or manifest.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: completed P03 manifest must remain done/unauthorized")
    if "Status: **SATISFIED**" not in exit_gate:
        raise SystemExit("ERROR: completed P03 requires SATISFIED exit gate")

    if terminal:
        if (phase_rows.get("P04") or {}).get("state") != "planned":
            raise SystemExit("ERROR: P04 must remain planned until a separate governed activation")
        if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
            raise SystemExit("ERROR: completed P03 terminal checkpoint must lock kernel and business implementation")
    else:
        p04_row = phase_rows.get("P04") or {}
        if p04_row.get("state") != "active" or p04_row.get("active_work_package") != current:
            raise SystemExit("ERROR: P04 historical checkpoint must identify the current P04 package")
        if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
            raise SystemExit("ERROR: active P04 must keep business code locked while P03 remains historical")

    prep = state.get("p03_preparation") or {}
    expected_prep = {
        "state": "completed",
        "phase_state": "done",
        "next_work_package": None,
        "work_package_state": "done",
        "work_package_spec": "docs/roadmap/work-packages/P03.11.md",
        "package_sequence": "docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json",
        "prepared_spec_count": 11,
        "mandatory_spec_count": 11,
        "entry_gate": "docs/governance/P03_ENTRY_GATE.md",
        "exit_gate": "docs/governance/P03_EXIT_GATE.md",
        "transition_checklist": "docs/governance/P02_P03_TRANSITION_CHECKLIST.md",
        "ai_native_alignment": "docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md",
    }
    for key, expected in expected_prep.items():
        if prep.get(key) != expected:
            label = "historical" if historical else "terminal"
            raise SystemExit(f"ERROR: {label} p03_preparation.{key} must be {expected}")
    if prep.get("blocking_gate") is not None:
        raise SystemExit("ERROR: completed P03 must have no unresolved P03 blocker")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
for marker in [
    "runs-on: ubuntu-24.04",
    "RUNNER_ENVIRONMENT",
    "github-hosted",
    "python scripts/validate_p03_preparation.py",
    "python scripts/validate_p03_package_specs.py",
    "bash scripts/verify_go_quality.sh",
    "bash scripts/verify_p01_12.sh",
    "bash scripts/verify_p02_10.sh",
]:
    if marker not in workflow:
        raise SystemExit(f"ERROR: governance workflow missing P03 readiness marker: {marker}")
if "self-hosted" in workflow or "LOCAL-WIN-" in workflow:
    raise SystemExit("ERROR: local/self-hosted governance runners are prohibited")

mode = (
    "PLANNING / NOT ACTIVATED"
    if planning
    else "ACTIVE"
    if active
    else "COMPLETED / HISTORICAL — P04 ACTIVE"
    if historical
    else "COMPLETED / NOT ADVANCED"
)
kernel_mode = (
    f"AUTHORIZED FOR {current}"
    if active
    else "P03 LOCKED / HISTORICAL"
    if historical
    else "LOCKED"
)
print("Omnexa P03 preparation/readiness validation: PASS")
print("P02 exit: SATISFIED")
print("P03 specifications: 11 / 11")
print(f"Mode: {mode}")
print(f"Active P03 package: {current if active else 'NONE'}")
print("Governance runner: GITHUB-HOSTED ONLY / ubuntu-24.04")
print(f"P03 kernel authority: {kernel_mode}")
print("Business feature code: LOCKED")
