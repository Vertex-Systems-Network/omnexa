#!/usr/bin/env python3
"""Validate prepared, active, terminal, or historical P02 strict sequential package state."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))
MANIFEST = json.loads((ROOT / "docs/roadmap/work-packages/P02_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))

EXPECTED = [f"P02.{i:02d}" for i in range(1, 11)]
OWNERS = {
    "P02.01": "kernel.identity",
    "P02.02": "kernel.tenancy",
    "P02.03": "kernel.organization",
    "P02.04": "kernel.identity",
    "P02.05": "kernel.authorization",
    "P02.06": "kernel.authorization",
    "P02.07": "kernel.identity",
    "P02.08": "kernel.identity",
    "P02.09": "kernel.configuration",
    "P02.10": "kernel.audit",
}

if MANIFEST.get("phase") != "P02" or MANIFEST.get("name") != "Identity, Tenancy & Organization":
    raise SystemExit("ERROR: P02 package manifest identity mismatch")
if MANIFEST.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P02 activation policy must be strict_sequential_one_active_package")
if MANIFEST.get("entry_gate") != "docs/governance/P02_ENTRY_GATE.md":
    raise SystemExit("ERROR: P02 entry gate path mismatch")
if MANIFEST.get("exit_gate") != "docs/governance/P02_EXIT_GATE.md":
    raise SystemExit("ERROR: P02 exit gate path mismatch")

packages = MANIFEST.get("packages") or []
if [item.get("id") for item in packages] != EXPECTED:
    raise SystemExit("ERROR: P02 package IDs/order must be exactly P02.01-P02.10")

current_phase = STATE.get("current_phase")
current = STATE.get("current_work_package")
phase = STATE.get("phase") or {}
phases = {item.get("id"): item for item in STATE.get("phases") or []}
planning = current_phase == "P01" and current is None and phase.get("state") == "done"
active = current_phase == "P02" and phase.get("state") == "active" and current in EXPECTED
terminal = current_phase == "P02" and phase.get("state") == "done" and current is None
historical = current_phase == "P03" and (phases.get("P02") or {}).get("state") == "done"
completed = terminal or historical
if not (planning or active or completed):
    raise SystemExit("ERROR: P02 specs may be validated only at completed-P01 planning, active P02, completed P02, or P03 historical checkpoint")

if planning:
    if MANIFEST.get("state") != "planned" or MANIFEST.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: P02 readiness manifest must remain planned with implementation_authorized=false")
    expected_states = ["planned"] * len(EXPECTED)
elif completed:
    if MANIFEST.get("state") != "done" or MANIFEST.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: completed P02 manifest must remain done with implementation_authorized=false")
    expected_states = ["done"] * len(EXPECTED)
else:
    if MANIFEST.get("state") != "active" or MANIFEST.get("implementation_authorized") is not True:
        raise SystemExit("ERROR: active P02 manifest must be active with implementation_authorized=true")
    current_index = EXPECTED.index(current)
    expected_states = ["done" if i < current_index else "active" if i == current_index else "planned" for i in range(len(EXPECTED))]

for index, item in enumerate(packages):
    pid = EXPECTED[index]
    expected_dep = [] if index == 0 else [EXPECTED[index - 1]]
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
if (planning or completed) and active_packages:
    raise SystemExit(f"ERROR: non-executing P02 record must have no active package, got {active_packages}")
if active and active_packages != [current]:
    raise SystemExit(f"ERROR: active P02 must have exactly current package {current}, got {active_packages}")

p02_row = phases.get("P02") or {}
p03_row = phases.get("P03") or {}
if planning and p02_row.get("state") != "planned":
    raise SystemExit("ERROR: phases[] P02 must remain planned during readiness preparation")
if active and (p02_row.get("state") != "active" or p02_row.get("active_work_package") != current):
    raise SystemExit("ERROR: phases[] P02 must be active and identify current package")
if completed:
    if p02_row.get("state") != "done" or p02_row.get("active_work_package") is not None:
        raise SystemExit("ERROR: completed phases[] P02 must remain done with no active package")
    if terminal:
        if p03_row.get("state") != "planned":
            raise SystemExit("ERROR: P03 must remain planned until a separate governed activation")
        lock = STATE.get("implementation_lock") or {}
        if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
            raise SystemExit("ERROR: terminal P02 checkpoint must lock kernel and business implementation")
    elif (STATE.get("implementation_lock") or {}).get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: business-feature implementation must remain locked during P03")

if active or terminal:
    phase_packages = phase.get("work_packages") or []
    if [item.get("id") for item in phase_packages] != EXPECTED:
        raise SystemExit("ERROR: STATE.phase package order must match P02 manifest")
    for manifest_item, state_item in zip(packages, phase_packages):
        if manifest_item.get("state") != state_item.get("state"):
            raise SystemExit(f"ERROR: STATE/manifest state mismatch for {manifest_item.get('id')}")
    done_count = sum(item.get("state") == "done" for item in packages)
    if phase.get("done_work_packages") != done_count:
        raise SystemExit("ERROR: P02 done_work_packages must match manifest completion")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
for index, item in enumerate(packages, start=1):
    if item.get("state") != "done":
        continue
    verifier = ROOT / f"scripts/verify_p02_{index:02d}.sh"
    marker = f"bash scripts/verify_p02_{index:02d}.sh"
    if not verifier.is_file():
        raise SystemExit(f"ERROR: completed {item.get('id')} missing regression verifier")
    if marker not in workflow:
        raise SystemExit(f"ERROR: completed {item.get('id')} verifier missing from governance workflow")

if active:
    active_index = EXPECTED.index(current) + 1
    active_verifier = ROOT / f"scripts/verify_p02_{active_index:02d}.sh"
    active_marker = f"bash scripts/verify_p02_{active_index:02d}.sh"
    if active_verifier.is_file() and active_marker not in workflow:
        raise SystemExit(f"ERROR: implemented active {current} verifier missing from governance workflow")

p02_10 = (ROOT / "docs/roadmap/work-packages/P02.10.md").read_text(encoding="utf-8").lower()
for marker in ["cross-tenant", "object/scope", "role", "service-account", "session invalidation", "p02 exit"]:
    if marker not in p02_10:
        raise SystemExit(f"ERROR: P02.10 must retain phase-exit marker: {marker}")

if completed:
    exit_gate = (ROOT / "docs/governance/P02_EXIT_GATE.md").read_text(encoding="utf-8")
    if "Status: **SATISFIED**" not in exit_gate:
        raise SystemExit("ERROR: completed P02 requires SATISFIED P02 exit gate")
    evidence = ROOT / "docs/roadmap/evidence/P02.10_COMPLETION_2026-08-26.md"
    if not evidence.is_file():
        raise SystemExit("ERROR: completed P02 requires P02.10 completion evidence")

print("Omnexa P02 package specification validation: PASS")
print("Prepared specs: 10 / 10")
print("Activation policy: STRICT SEQUENTIAL / ONE ACTIVE PACKAGE WHILE EXECUTING")
print(f"Mode: {'PLANNING' if planning else 'ACTIVE' if active else 'COMPLETED / HISTORICAL'}")
print(f"Active P02 package: {current if active else 'NONE'}")
print("Business feature code: LOCKED")
