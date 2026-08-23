#!/usr/bin/env python3
"""Validate P01.01-P01.12 strict sequential execution and terminal completion."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
MANIFEST = json.loads((ROOT / "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))

if STATE.get("current_phase") != "P01":
    raise SystemExit("ERROR: current phase checkpoint must be P01")
current = STATE.get("current_work_package")
expected_ids = [f"P01.{i:02d}" for i in range(1, 13)]
phase = STATE.get("phase") or {}
terminal = phase.get("state") == "done"

if terminal:
    if current is not None:
        raise SystemExit("ERROR: completed P01 must have no current work package")
    current_index = 12
else:
    if phase.get("state") != "active" or current not in expected_ids:
        raise SystemExit(f"ERROR: invalid current P01 package: {current}")
    current_index = expected_ids.index(current)

lock = STATE.get("implementation_lock") or {}
if terminal:
    if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: completed P01 locks must prohibit kernel and business-feature implementation")
else:
    if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: active P01 locks must authorize kernel code and prohibit business-feature code")

if MANIFEST.get("phase") != "P01":
    raise SystemExit("ERROR: P01 package manifest phase mismatch")
if MANIFEST.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P01 activation policy mismatch")
if terminal:
    if MANIFEST.get("state") != "done" or MANIFEST.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: completed P01 manifest must be done and implementation_authorized=false")
else:
    if MANIFEST.get("state") != "active" or MANIFEST.get("implementation_authorized") is not True:
        raise SystemExit("ERROR: active P01 manifest must authorize bounded kernel implementation")

packages = MANIFEST.get("packages") or []
if [item.get("id") for item in packages] != expected_ids:
    raise SystemExit("ERROR: P01 package order mismatch")

active = [item.get("id") for item in packages if item.get("state") == "active"]
if terminal:
    if active:
        raise SystemExit(f"ERROR: completed P01 must have no active package, got {active}")
else:
    if active != [current]:
        raise SystemExit(f"ERROR: exactly current package {current} must be active, got {active}")

for index, item in enumerate(packages):
    pid = expected_ids[index]
    if terminal:
        expected_state = "done"
    else:
        expected_state = "done" if index < current_index else "active" if index == current_index else "planned"
    if item.get("state") != expected_state:
        raise SystemExit(f"ERROR: {pid} must be {expected_state}, got {item.get('state')}")
    expected_dep = [] if index == 0 else [expected_ids[index - 1]]
    if item.get("depends_on") != expected_dep:
        raise SystemExit(f"ERROR: {pid} dependency must be {expected_dep}")
    spec = item.get("spec")
    if not isinstance(spec, str) or not (ROOT / spec).is_file():
        raise SystemExit(f"ERROR: missing spec for {pid}: {spec}")
    text = (ROOT / spec).read_text(encoding="utf-8")
    state_marker = f"State: `{expected_state}`"
    for marker in [pid, state_marker, "Acceptance criteria", "Completion evidence", "State transition"]:
        if marker.lower() not in text.lower():
            raise SystemExit(f"ERROR: {pid} spec missing marker: {marker}")

prep = STATE.get("p01_preparation") or {}
if terminal:
    if prep.get("state") != "completed" or prep.get("phase_state") != "done":
        raise SystemExit("ERROR: terminal p01_preparation must be completed/done")
    if prep.get("next_work_package") is not None or prep.get("work_package_state") != "done":
        raise SystemExit("ERROR: terminal p01_preparation must have no next P01 package and P01.12 done")
    if prep.get("work_package_spec") != "docs/roadmap/work-packages/P01.12.md":
        raise SystemExit("ERROR: terminal p01_preparation spec must be P01.12")
else:
    if prep.get("next_work_package") != current or prep.get("work_package_state") != "active":
        raise SystemExit("ERROR: p01_preparation must identify the current active package")
    if prep.get("work_package_spec") != f"docs/roadmap/work-packages/{current}.md":
        raise SystemExit("ERROR: p01_preparation active spec path mismatch")

phase_packages = phase.get("work_packages") or []
if [item.get("id") for item in phase_packages] != expected_ids:
    raise SystemExit("ERROR: STATE phase package order must match P01 manifest")
phase_states = {item.get("id"): item.get("state") for item in phase_packages}
for item in packages:
    if phase_states.get(item.get("id")) != item.get("state"):
        raise SystemExit(f"ERROR: STATE/manifest state mismatch for {item.get('id')}")
if phase.get("done_work_packages") != current_index:
    raise SystemExit("ERROR: P01 done_work_packages must equal the completed sequential prefix length")

phase_rows = {item.get("id"): item for item in STATE.get("phases") or []}
p01_row = phase_rows.get("P01") or {}
p02_row = phase_rows.get("P02") or {}
if terminal:
    if p01_row.get("state") != "done" or p01_row.get("active_work_package") is not None:
        raise SystemExit("ERROR: phases[] P01 row must be done with no active package")
    if p02_row.get("state") != "planned":
        raise SystemExit("ERROR: P02 must remain planned until separate activation")
else:
    if p01_row.get("state") != "active" or p01_row.get("active_work_package") != current:
        raise SystemExit("ERROR: phases[] P01 row must identify the current active package")

print("Omnexa P01 package sequence validation: PASS")
print("Prepared specs: 12 / 12")
print("Activation policy: STRICT SEQUENTIAL / ONE ACTIVE PACKAGE WHILE EXECUTING")
print(f"Completed packages: {current_index} / 12")
print(f"Active package: {current or 'NONE'}")
print("Business feature code: LOCKED")
