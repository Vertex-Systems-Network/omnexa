#!/usr/bin/env python3
"""Validate P04 readiness and the bounded P04.01 activation boundary."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE_PATH = ROOT / "docs/roadmap/STATE.json"
SEQUENCE_PATH = ROOT / "docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json"
ENTRY_PATH = ROOT / "docs/governance/P04_ENTRY_GATE.md"
READINESS_PATH = ROOT / "docs/governance/P03_P04_TRANSITION_READINESS.md"
P04_01_PATH = ROOT / "docs/roadmap/work-packages/P04.01.md"

REQUIRED_FILES = [
    "docs/governance/P03_EXIT_GATE.md",
    "docs/governance/P03_P04_TRANSITION_READINESS.md",
    "docs/governance/P04_ACTIVATION_GAP_LEDGER.md",
    "docs/governance/P04_ACTIVATION_TRANSACTION_TEMPLATE.md",
    "docs/governance/P04_ENTRY_GATE.md",
    "docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json",
    "docs/roadmap/work-packages/P04.01.md",
]

for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing P04 activation artifact: {relative}")

state = json.loads(STATE_PATH.read_text(encoding="utf-8"))
sequence = json.loads(SEQUENCE_PATH.read_text(encoding="utf-8"))
entry = ENTRY_PATH.read_text(encoding="utf-8")
readiness = READINESS_PATH.read_text(encoding="utf-8")
p04_01 = P04_01_PATH.read_text(encoding="utf-8")
phase_rows = {row.get("id"): row for row in state.get("phases") or []}
lock = state.get("implementation_lock") or {}

if (phase_rows.get("P03") or {}).get("state") != "done":
    raise SystemExit("ERROR: P04 requires completed P03")

p03_exit = (ROOT / "docs/governance/P03_EXIT_GATE.md").read_text(encoding="utf-8")
if "Status: **SATISFIED**" not in p03_exit:
    raise SystemExit("ERROR: P04 requires SATISFIED P03 exit gate")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "verified_protected" or branch.get("live_protected") is not True:
    raise SystemExit("ERROR: P04 requires verified protected main")
if branch.get("required_check") != "governance":
    raise SystemExit("ERROR: P04 requires governance as the protected required check")
ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_github_hosted" or ci.get("routing_mode") != "github_hosted_only":
    raise SystemExit("ERROR: P04 requires canonical GitHub-hosted CI")
if ci.get("runner_label") != "ubuntu-24.04" or ci.get("self_hosted_allowed") is not False:
    raise SystemExit("ERROR: P04 canonical runner policy mismatch")

if sequence.get("schema_version") != 1 or sequence.get("phase") != "P04":
    raise SystemExit("ERROR: invalid P04 package sequence identity")
if sequence.get("activation_policy") != "strict_sequential_one_active_package":
    raise SystemExit("ERROR: P04 must retain strict sequential single-package activation")

packages = sequence.get("packages") or []
expected_ids = [f"P04.{i:02d}" for i in range(1, 11)]
if [pkg.get("id") for pkg in packages] != expected_ids:
    raise SystemExit("ERROR: P04 package sequence must be exactly P04.01-P04.10 in order")

by_id = {pkg["id"]: pkg for pkg in packages}
expected_dependencies = {
    "P04.01": [],
    "P04.02": ["P04.01"],
    "P04.03": ["P04.01", "P04.02"],
    "P04.04": ["P04.01", "P04.03"],
    "P04.05": ["P04.01", "P04.04"],
    "P04.06": ["P04.03", "P04.05"],
    "P04.07": ["P04.01", "P04.06"],
    "P04.08": ["P04.01", "P04.07"],
    "P04.09": ["P04.03", "P04.08"],
    "P04.10": ["P04.01", "P04.09"],
}
for package_id, dependencies in expected_dependencies.items():
    if by_id[package_id].get("depends_on") != dependencies:
        raise SystemExit(f"ERROR: {package_id} dependency contract drift")

if by_id["P04.01"].get("spec") != "docs/roadmap/work-packages/P04.01.md":
    raise SystemExit("ERROR: P04.01 must have the accepted bounded specification")
for package_id in expected_ids[1:]:
    if by_id[package_id].get("spec") is not None:
        raise SystemExit(f"ERROR: {package_id} must not gain an implementation spec before predecessor acceptance")

for marker in [
    "duplicate publish/delivery can occur",
    "no global ordering guarantee exists",
    "tenant context is explicit",
    "must not select a broker",
    "database migration: none",
    "change business-feature authorization",
    "`P04.02` must remain locked",
]:
    if marker.lower() not in p04_01.lower():
        raise SystemExit(f"ERROR: P04.01 specification missing marker: {marker}")

for forbidden in ["kafka dependency", "rabbitmq dependency", "nats dependency", "redis streams dependency"]:
    if forbidden in p04_01.lower():
        raise SystemExit(f"ERROR: P04.01 prematurely selects transport technology: {forbidden}")

for marker in [
    "P03 exit = SATISFIED",
    "duplicate publish",
    "consumer crash before commit",
    "cross-tenant event attempt",
    "separate P04 activation transaction",
]:
    if marker.lower() not in readiness.lower():
        raise SystemExit(f"ERROR: P04 readiness missing marker: {marker}")

current_phase = state.get("current_phase")
current_package = state.get("current_work_package")
phase = state.get("phase") or {}
planning = (
    current_phase == "P03"
    and current_package is None
    and phase.get("id") == "P03"
    and phase.get("state") == "done"
)
active = (
    current_phase == "P04"
    and current_package == "P04.01"
    and phase.get("id") == "P04"
    and phase.get("state") == "active"
)

if not (planning or active):
    raise SystemExit("ERROR: invalid P04 readiness/activation checkpoint")

if planning:
    if (phase_rows.get("P04") or {}).get("state") != "planned":
        raise SystemExit("ERROR: P04 must remain planned before activation")
    if lock.get("kernel_code_authorized") is not False or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: readiness mode must keep implementation locked")
    if sequence.get("state") != "ready" or sequence.get("implementation_authorized") is not False:
        raise SystemExit("ERROR: readiness sequence must be ready but unauthorized")
    if by_id["P04.01"].get("state") != "ready":
        raise SystemExit("ERROR: P04.01 must be ready before activation")
    if any(by_id[p].get("state") != "planned" for p in expected_ids[1:]):
        raise SystemExit("ERROR: P04.02-P04.10 must remain planned before activation")
    if "NOT YET SATISFIED" not in entry:
        raise SystemExit("ERROR: readiness mode P04 entry gate must remain NOT YET SATISFIED")
else:
    if (phase_rows.get("P04") or {}).get("state") != "active":
        raise SystemExit("ERROR: active P04 requires phases[] P04 active")
    if (phase_rows.get("P04") or {}).get("active_work_package") != "P04.01":
        raise SystemExit("ERROR: phases[] P04 must identify P04.01 as the sole active package")
    if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: P04.01 activation must authorize bounded kernel code only")
    if sequence.get("state") != "active" or sequence.get("implementation_authorized") is not True:
        raise SystemExit("ERROR: active P04 sequence must authorize P04.01 implementation")
    if by_id["P04.01"].get("state") != "active":
        raise SystemExit("ERROR: P04.01 must be the sole active package")
    if any(by_id[p].get("state") != "planned" for p in expected_ids[1:]):
        raise SystemExit("ERROR: P04.02-P04.10 must remain planned while P04.01 is active")
    if "Status: **SATISFIED**" not in entry:
        raise SystemExit("ERROR: active P04 requires SATISFIED entry gate")

print("Omnexa P04 activation validation: PASS")
print(f"Mode: {'PLANNING / READY' if planning else 'ACTIVE / P04.01'}")
print(f"P04 sequence: {len(packages)} packages")
print(f"Kernel code: {'LOCKED' if planning else 'AUTHORIZED FOR P04.01 ONLY'}")
print("Business feature code: LOCKED")
