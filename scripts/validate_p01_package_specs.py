#!/usr/bin/env python3
"""Validate the completed P01 strict-sequential package record across later phases."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
MANIFEST = json.loads((ROOT / "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))
EXPECTED = [f"P01.{i:02d}" for i in range(1, 13)]

if STATE.get("current_phase") not in {"P01", "P02", "P03", "P04"}:
    raise SystemExit("ERROR: P01 historical validator must be reviewed before advancing beyond P04")
if MANIFEST.get("phase") != "P01" or MANIFEST.get("name") != "Omnexa Kernel":
    raise SystemExit("ERROR: P01 package manifest identity mismatch")
if MANIFEST.get("state") != "done" or MANIFEST.get("implementation_authorized") is not False:
    raise SystemExit("ERROR: completed P01 manifest must remain done and implementation_authorized=false")
if MANIFEST.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P01 activation policy mismatch")

packages = MANIFEST.get("packages") or []
if [item.get("id") for item in packages] != EXPECTED:
    raise SystemExit("ERROR: P01 package order mismatch")
if any(item.get("state") != "done" for item in packages):
    raise SystemExit("ERROR: every P01 package must remain done")
if any(item.get("state") == "active" for item in packages):
    raise SystemExit("ERROR: completed P01 must have zero active packages")

for index, item in enumerate(packages):
    pid = EXPECTED[index]
    expected_dep = [] if index == 0 else [EXPECTED[index - 1]]
    if item.get("depends_on") != expected_dep:
        raise SystemExit(f"ERROR: {pid} dependency must remain {expected_dep}")
    spec = item.get("spec")
    if spec != f"docs/roadmap/work-packages/{pid}.md" or not (ROOT / spec).is_file():
        raise SystemExit(f"ERROR: missing/canonical spec path mismatch for {pid}")
    text = (ROOT / spec).read_text(encoding="utf-8")
    for marker in [pid, "State: `done`", "Acceptance criteria", "Completion evidence", "State transition"]:
        if marker.lower() not in text.lower():
            raise SystemExit(f"ERROR: {pid} completed spec missing marker: {marker}")

phases = {item.get("id"): item for item in STATE.get("phases") or []}
p01_row = phases.get("P01") or {}
if p01_row.get("state") != "done" or p01_row.get("active_work_package") is not None:
    raise SystemExit("ERROR: phases[] P01 must remain done with no active package")

prep = STATE.get("p01_preparation") or {}
if prep.get("state") != "completed" or prep.get("phase_state") != "done":
    raise SystemExit("ERROR: p01_preparation must remain completed/done")
if prep.get("next_work_package") is not None or prep.get("work_package_state") != "done":
    raise SystemExit("ERROR: completed P01 preparation must have no next P01 package")
if prep.get("work_package_spec") != "docs/roadmap/work-packages/P01.12.md":
    raise SystemExit("ERROR: completed P01 preparation must retain P01.12 as final spec")
if prep.get("prepared_spec_count") != 12 or prep.get("mandatory_spec_count") != 12:
    raise SystemExit("ERROR: P01 must retain 12 / 12 package specifications")

exit_gate = ROOT / "docs/governance/P01_EXIT_GATE.md"
if not exit_gate.is_file() or "Status: **SATISFIED**" not in exit_gate.read_text(encoding="utf-8"):
    raise SystemExit("ERROR: completed P01 requires SATISFIED P01 exit gate")

print("Omnexa completed P01 package sequence validation: PASS")
print("Prepared specs: 12 / 12")
print("Completed packages: 12 / 12")
print("Active P01 package: NONE")
print("P01 implementation: LOCKED / HISTORICAL")