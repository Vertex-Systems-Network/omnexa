#!/usr/bin/env python3
"""Validate P01.01-P01.12 sequence with exactly one active package."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
MANIFEST = json.loads((ROOT / "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))

if STATE.get("current_phase") != "P01" or STATE.get("current_work_package") != "P01.01":
    raise SystemExit("ERROR: current P01 package must be P01.01")
lock = STATE.get("implementation_lock") or {}
if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: P01 locks must authorize kernel and prohibit business features")

if MANIFEST.get("phase") != "P01" or MANIFEST.get("state") != "active":
    raise SystemExit("ERROR: P01 package manifest must be active")
if MANIFEST.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P01 activation policy mismatch")
if MANIFEST.get("implementation_authorized") is not True:
    raise SystemExit("ERROR: P01 manifest must authorize bounded implementation")

packages = MANIFEST.get("packages") or []
expected_ids = [f"P01.{i:02d}" for i in range(1, 13)]
if [item.get("id") for item in packages] != expected_ids:
    raise SystemExit("ERROR: P01 package order mismatch")

active = [item.get("id") for item in packages if item.get("state") == "active"]
if active != ["P01.01"]:
    raise SystemExit(f"ERROR: exactly P01.01 must be active, got {active}")

for index, item in enumerate(packages):
    pid = expected_ids[index]
    expected_state = "active" if index == 0 else "planned"
    if item.get("state") != expected_state:
        raise SystemExit(f"ERROR: {pid} must be {expected_state}")
    expected_dep = [] if index == 0 else [expected_ids[index - 1]]
    if item.get("depends_on") != expected_dep:
        raise SystemExit(f"ERROR: {pid} dependency must be {expected_dep}")
    spec = item.get("spec")
    if not isinstance(spec, str) or not (ROOT / spec).is_file():
        raise SystemExit(f"ERROR: missing spec for {pid}: {spec}")
    text = (ROOT / spec).read_text(encoding="utf-8")
    state_marker = "State: `active`" if index == 0 else "State: `planned`"
    for marker in [pid, state_marker, "Acceptance criteria", "Completion evidence", "State transition"]:
        if marker.lower() not in text.lower():
            raise SystemExit(f"ERROR: {pid} spec missing marker: {marker}")

prep = STATE.get("p01_preparation") or {}
if prep.get("next_work_package") != "P01.01" or prep.get("work_package_state") != "active":
    raise SystemExit("ERROR: P01.01 must be the active controlling package")

phase = STATE.get("phase") or {}
phase_packages = phase.get("work_packages") or []
if [item.get("id") for item in phase_packages] != expected_ids:
    raise SystemExit("ERROR: STATE phase package order must match P01 manifest")
if [item.get("id") for item in phase_packages if item.get("state") == "active"] != ["P01.01"]:
    raise SystemExit("ERROR: STATE must contain exactly one active P01 package")

print("Omnexa P01 package sequence validation: PASS")
print("Prepared specs: 12 / 12")
print("Activation policy: STRICT SEQUENTIAL / ONE ACTIVE PACKAGE")
print("Active package: P01.01")
print("P01.02-P01.12: PLANNED")
print("Business feature code: LOCKED")
