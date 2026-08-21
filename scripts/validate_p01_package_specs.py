#!/usr/bin/env python3
"""Validate P01.01-P01.12 specification completeness/order while implementation remains locked."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
MANIFEST_PATH = ROOT / "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json"
MANIFEST = json.loads(MANIFEST_PATH.read_text(encoding="utf-8"))

if STATE.get("current_phase") != "P00" or STATE.get("current_work_package") != "P00.10":
    raise SystemExit("ERROR: P01 specification preparation expects P00.10 exit-verification state")

lock = STATE.get("implementation_lock") or {}
if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: P01 specification preparation requires both implementation locks to remain false")

if MANIFEST.get("phase") != "P01" or MANIFEST.get("state") != "blocked":
    raise SystemExit("ERROR: P01 manifest must remain blocked during preparation")
if MANIFEST.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P01 activation policy mismatch")
if MANIFEST.get("implementation_authorized") is not False:
    raise SystemExit("ERROR: P01 manifest must not authorize implementation")

packages = MANIFEST.get("packages") or []
expected_ids = [f"P01.{i:02d}" for i in range(1, 13)]
ids = [item.get("id") for item in packages]
if ids != expected_ids:
    raise SystemExit(f"ERROR: P01 package order mismatch: {ids}")

for index, item in enumerate(packages):
    pid = expected_ids[index]
    if item.get("state") != "planned":
        raise SystemExit(f"ERROR: {pid} must remain planned")
    expected_dep = [] if index == 0 else [expected_ids[index - 1]]
    if item.get("depends_on") != expected_dep:
        raise SystemExit(f"ERROR: {pid} dependency must be {expected_dep}")
    spec = item.get("spec")
    if not isinstance(spec, str) or not (ROOT / spec).is_file():
        raise SystemExit(f"ERROR: missing spec for {pid}: {spec}")
    text = (ROOT / spec).read_text(encoding="utf-8")
    for marker in [pid, "State: `planned`", "Specification-only", "Acceptance criteria", "Completion evidence", "State transition"]:
        if marker.lower() not in text.lower():
            raise SystemExit(f"ERROR: {pid} spec missing marker: {marker}")

# Confirm the controlling next package remains P01.01 and future specs do not change state.
prep = STATE.get("p01_preparation") or {}
if prep.get("next_work_package") != "P01.01" or prep.get("work_package_state") != "planned":
    raise SystemExit("ERROR: P01.01 must remain the prepared next package")

for forbidden in ["go.work", "go.mod", "kernel/cmd/omnexa/main.go"]:
    if (ROOT / forbidden).exists():
        raise SystemExit(f"ERROR: premature executable P01 path exists while kernel code is locked: {forbidden}")

print("Omnexa P01 package specification validation: PASS")
print("Prepared specs: 12 / 12")
print("Activation policy: STRICT SEQUENTIAL / ONE ACTIVE PACKAGE")
print("P01 implementation: BLOCKED")
print("Kernel code: LOCKED")
