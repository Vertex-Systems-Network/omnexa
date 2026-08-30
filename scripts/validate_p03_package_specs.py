#!/usr/bin/env python3
"""Validate prepared, active, terminal, or P04-era historical P03 package state."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
MANIFEST = json.loads((ROOT / "docs/roadmap/work-packages/P03_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))

EXPECTED = [f"P03.{i:02d}" for i in range(1, 12)]
NAMES = {
    "P03.01": "Module manifest schema",
    "P03.02": "Registry & deterministic discovery",
    "P03.03": "Dependency graph resolver",
    "P03.04": "Module lifecycle state machine",
    "P03.05": "Module settings & feature flags",
    "P03.06": "Capability registry",
    "P03.07": "Permission registration",
    "P03.08": "UI contribution registry contract",
    "P03.09": "Migration ownership registry",
    "P03.10": "Module health reporting",
    "P03.11": "Package trust hooks & P03 exit proof",
}
OWNERS = {pid: "kernel.modules" for pid in EXPECTED}

if MANIFEST.get("phase") != "P03" or MANIFEST.get("name") != "Module Runtime":
    raise SystemExit("ERROR: P03 package manifest identity mismatch")
if MANIFEST.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P03 activation policy must be strict_sequential_one_active_package")
if MANIFEST.get("entry_gate") != "docs/governance/P03_ENTRY_GATE.md":
    raise SystemExit("ERROR: P03 entry gate path mismatch")
if MANIFEST.get("exit_gate") != "docs/governance/P03_EXIT_GATE.md":
    raise SystemExit("ERROR: P03 exit gate path mismatch")
if MANIFEST.get("transition_checklist") != "docs/governance/P02_P03_TRANSITION_CHECKLIST.md":
    raise SystemExit("ERROR: P03 transition checklist path mismatch")
if MANIFEST.get("ai_native_alignment") != "docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md":
    raise SystemExit("ERROR: P03 AI-native alignment path mismatch")

packages = MANIFEST.get("packages") or []
if [item.get("id") for item in packages] != EXPECTED:
    raise SystemExit("ERROR: P03 package IDs/order must be exactly P03.01-P03.11")

current_phase = STATE.get("current_phase")
current = STATE.get("current_work_package")
phase = STATE.get("phase") or {}
planning = current_phase == "P02" and current is None and phase.get("id") == "P02" and phase.get("state") == "done"
active = current_phase == "P03" and phase.get("id") == "P03" and phase.get("state") == "active" and current in EXPECTED
terminal = current_phase == "P03" and phase.get("id") == "P03" and phase.get("state") == "done" and current is None
historical = current_phase == "P04" and phase.get("id") == "P04" and phase.get("state") == "active" and isinstance(current, str) and current.startswith("P04.")
if not (planning or active or terminal or historical):
    raise SystemExit("ERROR: P03 specs may be validated only at completed-P02 planning, active P03, completed-P03 terminal, or active-P04 historical checkpoint")

if planning:
    if MANIFEST.get("state") != "planned" or MANIFEST.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: P03 readiness manifest must remain planned with implementation_authorized=false")
    expected_states = ["planned"] * len(EXPECTED)
elif terminal or historical:
    if MANIFEST.get("state") != "done" or MANIFEST.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: completed P03 manifest must remain done with implementation_authorized=false")
    expected_states = ["done"] * len(EXPECTED)
else:
    if MANIFEST.get("state") != "active" or MANIFEST.get("implementation_authorized") is not True:
        raise SystemExit("ERROR: active P03 manifest must be active with implementation_authorized=true")
    current_index = EXPECTED.index(current)
    expected_states = ["done" if i < current_index else "active" if i == current_index else "planned" for i in range(len(EXPECTED))]

for index, item in enumerate(packages):
    pid = EXPECTED[index]
    expected_dep = [] if index == 0 else [EXPECTED[index - 1]]
    if item.get("name") != NAMES[pid]:
        raise SystemExit(f"ERROR: {pid} name must be {NAMES[pid]}")
    if item.get("depends_on") != expected_dep:
        raise SystemExit(f"ERROR: {pid} dependency must be {expected_dep}")
    if item.get("owner") != OWNERS[pid]:
        raise SystemExit(f"ERROR: {pid} owner must be {OWNERS[pid]}")
    if item.get("state") != expected_states[index]:
        raise SystemExit(f"ERROR: {pid} must be {expected_states[index]}, got {item.get('state')}")
    spec = item.get("spec")
    if spec != f"docs/roadmap/work-packages/{pid}.md" or not (ROOT / spec).is_file():
        raise SystemExit(f"ERROR: missing/canonical spec path mismatch for {pid}")
    text = (ROOT / spec).read_text(encoding="utf-8")
    required = [
        pid,
        f"State: `{expected_states[index]}`",
        "Owner/domain:",
        "Objective",
        "In scope",
        "Out of scope",
        "Security and architecture invariants",
        "Acceptance criteria",
        "Completion evidence",
        "State transition",
        "GitHub-hosted",
        "business_feature_code_authorized=false",
    ]
    for marker in required:
        if marker.lower() not in text.lower():
            raise SystemExit(f"ERROR: {pid} spec missing marker: {marker}")

active_packages = [item.get("id") for item in packages if item.get("state") == "active"]
if (planning or terminal or historical) and active_packages:
    raise SystemExit(f"ERROR: non-executing/historical P03 checkpoint must have no active P03 package, got {active_packages}")
if active and active_packages != [current]:
    raise SystemExit(f"ERROR: active P03 must have exactly current package {current}, got {active_packages}")

p03_row = next((item for item in STATE.get("phases") or [] if item.get("id") == "P03"), {})
p04_row = next((item for item in STATE.get("phases") or [] if item.get("id") == "P04"), {})
if planning and p03_row.get("state") != "planned":
    raise SystemExit("ERROR: phases[] P03 must remain planned during readiness preparation")
if active and (p03_row.get("state") != "active" or p03_row.get("active_work_package") != current):
    raise SystemExit("ERROR: phases[] P03 must be active and identify current package")
if terminal or historical:
    if p03_row.get("state") != "done" or p03_row.get("active_work_package") is not None:
        raise SystemExit("ERROR: completed phases[] P03 must remain done with no active P03 package")

if terminal:
    if p04_row.get("state") != "planned":
        raise SystemExit("ERROR: P04 must remain planned until a separate governed activation")
    lock = STATE.get("implementation_lock") or {}
    if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: completed P03 terminal checkpoint must lock kernel and business implementation")

if historical:
    if p04_row.get("state") != "active" or p04_row.get("active_work_package") != current:
        raise SystemExit("ERROR: active-P04 historical checkpoint must identify the current P04 package")
    lock = STATE.get("implementation_lock") or {}
    if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: active P04 must authorize bounded kernel code and keep business code locked")

if active or terminal:
    phase_packages = phase.get("work_packages") or []
    if [item.get("id") for item in phase_packages] != EXPECTED:
        raise SystemExit("ERROR: STATE.phase package order must match P03 manifest")
    for manifest_item, state_item in zip(packages, phase_packages):
        if manifest_item.get("state") != state_item.get("state"):
            raise SystemExit(f"ERROR: STATE/manifest state mismatch for {manifest_item.get('id')}")
    done_count = sum(item.get("state") == "done" for item in packages)
    if phase.get("done_work_packages") != done_count:
        raise SystemExit("ERROR: P03 done_work_packages must match manifest completion")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
for index, item in enumerate(packages, start=1):
    if item.get("state") != "done":
        continue
    verifier = ROOT / f"scripts/verify_p03_{index:02d}.sh"
    marker = f"bash scripts/verify_p03_{index:02d}.sh"
    if not verifier.is_file():
        raise SystemExit(f"ERROR: completed {item.get('id')} missing regression verifier")
    if marker not in workflow:
        raise SystemExit(f"ERROR: completed {item.get('id')} verifier missing from governance workflow")

if active:
    active_index = EXPECTED.index(current) + 1
    active_verifier = ROOT / f"scripts/verify_p03_{active_index:02d}.sh"
    active_marker = f"bash scripts/verify_p03_{active_index:02d}.sh"
    if active_verifier.is_file() and active_marker not in workflow:
        raise SystemExit(f"ERROR: implemented active {current} verifier missing from governance workflow")

p03_11 = (ROOT / "docs/roadmap/work-packages/P03.11.md").read_text(encoding="utf-8").lower()
for marker in [
    "required dependency enforcement",
    "optional dependency degradation",
    "safe disable/re-enable",
    "upgrade/migration path",
    "forbidden cross-module dependency detection",
    "health/state accuracy",
    "no unrelated module corruption",
    "p03 exit",
]:
    if marker not in p03_11:
        raise SystemExit(f"ERROR: P03.11 must retain phase-exit marker: {marker}")

alignment = (ROOT / "docs/roadmap/P03_AI_NATIVE_ALIGNMENT.md").read_text(encoding="utf-8").lower()
for marker in ["xq-100", "xsg-100", "xtrust-100", "xpf-200", "xperf-100", "planning-only"]:
    if marker not in alignment:
        raise SystemExit(f"ERROR: P03 AI-native alignment missing marker: {marker}")

if terminal or historical:
    exit_gate = (ROOT / "docs/governance/P03_EXIT_GATE.md").read_text(encoding="utf-8")
    if "Status: **SATISFIED**" not in exit_gate:
        raise SystemExit("ERROR: completed P03 requires SATISFIED P03 exit gate")
    evidence = list((ROOT / "docs/roadmap/evidence").glob("P03.11_COMPLETION_*.md"))
    if not evidence:
        raise SystemExit("ERROR: completed P03 requires P03.11 completion evidence")

mode = "PLANNING" if planning else "ACTIVE" if active else "COMPLETED / HISTORICAL" if historical else "COMPLETED"
print("Omnexa P03 package specification validation: PASS")
print("Prepared specs: 11 / 11")
print("Activation policy: STRICT SEQUENTIAL / ONE ACTIVE PACKAGE WHILE EXECUTING")
print(f"Mode: {mode}")
print(f"Active P03 package: {current if active else 'NONE'}")
print("Business feature code: LOCKED")
