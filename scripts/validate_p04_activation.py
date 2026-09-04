#!/usr/bin/env python3
"""Validate P04 readiness and strict sequential package activation boundaries."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE_PATH = ROOT / "docs/roadmap/STATE.json"
SEQUENCE_PATH = ROOT / "docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json"
ENTRY_PATH = ROOT / "docs/governance/P04_ENTRY_GATE.md"
READINESS_PATH = ROOT / "docs/governance/P03_P04_TRANSITION_READINESS.md"
P04_01_PATH = ROOT / "docs/roadmap/work-packages/P04.01.md"
P04_02_PATH = ROOT / "docs/roadmap/work-packages/P04.02.md"
P04_03_PATH = ROOT / "docs/roadmap/work-packages/P04.03.md"
P04_04_PATH = ROOT / "docs/roadmap/work-packages/P04.04.md"
P04_05_PATH = ROOT / "docs/roadmap/work-packages/P04.05.md"

REQUIRED_FILES = [
    "docs/governance/P03_EXIT_GATE.md",
    "docs/governance/P03_P04_TRANSITION_READINESS.md",
    "docs/governance/P04_ACTIVATION_GAP_LEDGER.md",
    "docs/governance/P04_ACTIVATION_TRANSACTION_TEMPLATE.md",
    "docs/governance/P04_ENTRY_GATE.md",
    "docs/roadmap/work-packages/P04_PACKAGE_SEQUENCE.json",
    "docs/roadmap/work-packages/P04.01.md",
    "docs/roadmap/work-packages/P04.02.md",
    "docs/roadmap/work-packages/P04.03.md",
    "docs/roadmap/work-packages/P04.04.md",
    "docs/roadmap/work-packages/P04.05.md",
]

for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing P04 activation artifact: {relative}")

state = json.loads(STATE_PATH.read_text(encoding="utf-8"))
sequence = json.loads(SEQUENCE_PATH.read_text(encoding="utf-8"))
entry = ENTRY_PATH.read_text(encoding="utf-8")
readiness = READINESS_PATH.read_text(encoding="utf-8")
p04_01 = P04_01_PATH.read_text(encoding="utf-8")
p04_02 = P04_02_PATH.read_text(encoding="utf-8")
p04_03 = P04_03_PATH.read_text(encoding="utf-8")
p04_04 = P04_04_PATH.read_text(encoding="utf-8")
p04_05 = P04_05_PATH.read_text(encoding="utf-8")
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
    "provider-neutral publish/subscribe boundary",
    "duplicate publish/delivery can occur",
    "no global ordering guarantee exists",
    "receipt of an event is not an authorization credential",
    "P04.02 owns no concrete broker/provider",
    "does not implement a worker loop",
    "database migration: `none`",
    "`P04.03`",
]:
    if marker.lower() not in p04_02.lower():
        raise SystemExit(f"ERROR: P04.02 specification missing marker: {marker}")

for marker in [
    "checkpoint state represents **consumption progress only**",
    "at-least-once-compatible semantics",
    "no global ordering guarantee exists",
    "checkpoint or transport offset is never an authorization credential",
    "no provider or storage feature may be described as proving end-to-end exactly-once business mutation",
    "P04.04",
]:
    if marker.lower() not in p04_03.lower():
        raise SystemExit(f"ERROR: P04.03 specification missing marker: {marker}")

for marker in [
    "same local PostgreSQL transaction",
    "duplicate publication remains explicitly possible",
    "no global event ordering guarantee",
    "does not mean a downstream protected mutation happened exactly once",
    "preparation carrier adds no migration",
    "P04.05",
]:
    if marker.lower() not in p04_04.lower():
        raise SystemExit(f"ERROR: P04.04 specification missing marker: {marker}")

for marker in [
    "same local PostgreSQL transaction",
    "EventID alone must not become a global cross-consumer lock",
    "P04.03 checkpoints and P04.05 inbox completion represent different facts",
    "external/non-transactional side effects are not made exactly-once",
    "P04.06",
]:
    if marker.lower() not in p04_05.lower():
        raise SystemExit(f"ERROR: P04.05 specification missing marker: {marker}")

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
    and current_package in expected_ids
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
    mode = "PLANNING / READY"
    authority = "LOCKED"
else:
    active_index = expected_ids.index(current_package)
    done_ids = expected_ids[:active_index]
    future_ids = expected_ids[active_index + 1 :]

    if (phase_rows.get("P04") or {}).get("state") != "active":
        raise SystemExit("ERROR: active P04 requires phases[] P04 active")
    if (phase_rows.get("P04") or {}).get("active_work_package") != current_package:
        raise SystemExit("ERROR: phases[] P04 active_work_package must match the canonical current package")
    if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
        raise SystemExit("ERROR: active P04 must authorize bounded kernel code only")
    if sequence.get("state") != "active" or sequence.get("implementation_authorized") is not True:
        raise SystemExit("ERROR: active P04 sequence must remain implementation-authorized")

    for package_id in done_ids:
        package = by_id[package_id]
        if package.get("state") != "done":
            raise SystemExit(f"ERROR: completed predecessor {package_id} must be done")
        expected_spec = f"docs/roadmap/work-packages/{package_id}.md"
        if package.get("spec") != expected_spec or not (ROOT / expected_spec).is_file():
            raise SystemExit(f"ERROR: completed predecessor {package_id} must retain its accepted specification")
        evidence = package.get("evidence") or []
        if not evidence or any(not (ROOT / item).is_file() for item in evidence):
            raise SystemExit(f"ERROR: completed predecessor {package_id} must retain completion evidence")

        tracking_key = f"p04_{package_id.split('.')[1]}_completion"
        completion = tracking.get(tracking_key) or {}
        if completion.get("state") != "PASS":
            raise SystemExit(f"ERROR: completed predecessor {package_id} requires retained PASS tracking evidence")
        if completion.get("completion_evidence") not in evidence:
            raise SystemExit(f"ERROR: completed predecessor {package_id} tracking/evidence reference drift")
        for field in ["final_exact_head", "implementation_merge", "workflow_run", "job"]:
            if not completion.get(field):
                raise SystemExit(f"ERROR: completed predecessor {package_id} missing canonical {field} tracking")
        if completion.get("evidence_environment") != "github-hosted":
            raise SystemExit(f"ERROR: completed predecessor {package_id} evidence must remain GitHub-hosted")
        if completion.get("runner_image") != "ubuntu-24.04":
            raise SystemExit(f"ERROR: completed predecessor {package_id} runner-image evidence drift")

    active_package = by_id[current_package]
    if active_package.get("state") != "active":
        raise SystemExit(f"ERROR: {current_package} must be the sole active package")
    expected_active_spec = f"docs/roadmap/work-packages/{current_package}.md"
    if active_package.get("spec") != expected_active_spec or not (ROOT / expected_active_spec).is_file():
        raise SystemExit(f"ERROR: active package {current_package} must have its accepted specification")

    if any(by_id[p].get("state") != "planned" for p in future_ids):
        raise SystemExit("ERROR: future P04 packages must remain planned")
    if any(by_id[p].get("spec") is not None for p in future_ids):
        raise SystemExit("ERROR: future P04 package sequence entries must remain without accepted implementation specs")

    phase_packages = {pkg.get("id"): pkg for pkg in phase.get("work_packages") or []}
    if phase.get("done_work_packages") != len(done_ids):
        raise SystemExit("ERROR: P04 done_work_packages must equal the strict completed predecessor count")
    for package_id in done_ids:
        if (phase_packages.get(package_id) or {}).get("state") != "done":
            raise SystemExit(f"ERROR: STATE phase mirror must mark {package_id} done")
    if (phase_packages.get(current_package) or {}).get("state") != "active":
        raise SystemExit("ERROR: STATE phase mirror must mark the canonical current package active")
    if any((phase_packages.get(p) or {}).get("state") != "planned" for p in future_ids):
        raise SystemExit("ERROR: STATE phase mirror must keep future P04 packages planned")

    preparation = state.get("p04_preparation") or {}
    if preparation.get("next_work_package") != current_package:
        raise SystemExit("ERROR: P04 preparation cursor must match the canonical current package")
    if preparation.get("work_package_state") != "active":
        raise SystemExit("ERROR: P04 preparation must identify the current package as active")
    if preparation.get("work_package_spec") != expected_active_spec:
        raise SystemExit("ERROR: P04 preparation spec must match the canonical current package")
    if preparation.get("prepared_spec_count") != active_index + 1:
        raise SystemExit("ERROR: P04 prepared spec count must include all completed predecessors plus the active package")

    if "Status: **SATISFIED**" not in entry:
        raise SystemExit("ERROR: active P04 requires SATISFIED entry gate")

    mode = f"ACTIVE / {current_package}"
    authority = f"AUTHORIZED FOR {current_package} ONLY"

print("Omnexa P04 sequential activation validation: PASS")
print(f"Mode: {mode}")
print(f"P04 sequence: {len(packages)} packages")
print(f"Kernel code: {authority}")
print("Business feature code: LOCKED")
