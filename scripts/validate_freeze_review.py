#!/usr/bin/env python3
"""Validate frozen P00 evidence and completed P01 prerequisites across P02/P03/P04."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/governance/FOUNDATION_FREEZE_REVIEW.md",
    "docs/governance/P01_ENTRY_GATE.md",
    "docs/governance/P01_EXIT_GATE.md",
    "docs/governance/P00_P01_TRANSITION_CHECKLIST.md",
    "docs/governance/FOUNDATION_FREEZE.json",
    "docs/governance/REPOSITORY_HARDENING.md",
    "docs/governance/BRANCH_PROTECTION_ADMIN_RUNBOOK.md",
    "docs/contracts/governance/foundation-freeze.schema.json",
    "docs/adr/ADR-0010-foundation-architecture-freeze.md",
    "docs/adr/ADR-0006-temporary-p00-ci-evidence-exception.md",
    "scripts/apply_main_protection.ps1",
    "scripts/verify_main_protection.ps1",
]
for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing freeze/transition artifact: {relative}")

manifest = json.loads((ROOT / "docs/governance/FOUNDATION_FREEZE.json").read_text(encoding="utf-8"))
state = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))

if manifest.get("version") != "foundation-v1" or manifest.get("architecture_status") != "FROZEN":
    raise SystemExit("ERROR: Foundation v1 must remain FROZEN")
if manifest.get("p00_exit_status") != "DONE":
    raise SystemExit("ERROR: P00 exit must remain DONE")
expected_frozen = {f"P00.{i:02d}" for i in range(1, 10)}
if set(manifest.get("frozen_packages") or []) != expected_frozen:
    raise SystemExit("ERROR: frozen architecture set must remain exactly P00.01-P00.09")

entry = manifest.get("p01_entry_gate") or {}
if entry.get("state") != "SATISFIED" or entry.get("kernel_code_authorized") is not True or entry.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: historical P01 entry gate evidence is inconsistent")
controls = {item.get("tracker"): item.get("state") for item in entry.get("blockers") or []}
if controls.get("issue:#3") != "SATISFIED" or controls.get("issue:#14") != "SATISFIED":
    raise SystemExit("ERROR: EG-02/Issue #3 and EG-03/Issue #14 must remain SATISFIED")

if state.get("current_phase") not in {"P01", "P02", "P03", "P04"}:
    raise SystemExit("ERROR: freeze validator must be reviewed before advancing beyond P04")
phases = {item.get("id"): item for item in state.get("phases") or []}
if (phases.get("P00") or {}).get("state") != "done":
    raise SystemExit("ERROR: phases[] P00 must remain done")
p01 = phases.get("P01") or {}
if p01.get("state") != "done" or p01.get("active_work_package") is not None:
    raise SystemExit("ERROR: phases[] P01 must remain done with no active work package")

p01_manifest = json.loads((ROOT / "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json").read_text(encoding="utf-8"))
if p01_manifest.get("state") != "done" or p01_manifest.get("implementation_authorized") is not False:
    raise SystemExit("ERROR: completed P01 package manifest must remain locked")
if len(p01_manifest.get("packages") or []) != 12 or any(item.get("state") != "done" for item in p01_manifest.get("packages") or []):
    raise SystemExit("ERROR: P01 must remain complete at 12 / 12")

p01_exit = (ROOT / "docs/governance/P01_EXIT_GATE.md").read_text(encoding="utf-8")
if "Status: **SATISFIED**" not in p01_exit:
    raise SystemExit("ERROR: P01 exit must remain SATISFIED")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "verified_protected" or branch.get("live_protected") is not True:
    raise SystemExit("ERROR: main protection must remain verified_protected / live_protected=true")
if branch.get("issue") != 3 or branch.get("issue_state") != "closed" or branch.get("required_check") != "governance":
    raise SystemExit("ERROR: protected integration evidence is inconsistent")
if branch.get("repository_visibility") != "public":
    raise SystemExit("ERROR: repository visibility must remain public")
if "HTTP 403" not in str(branch.get("historical_private_attempt_result") or ""):
    raise SystemExit("ERROR: historical private-repository 403 evidence must be retained")

ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_github_hosted" or ci.get("routing_mode") != "github_hosted_only":
    raise SystemExit("ERROR: canonical CI must remain operational_github_hosted / github_hosted_only")
if ci.get("runner_label") != "ubuntu-24.04" or ci.get("self_hosted_allowed") is not False:
    raise SystemExit("ERROR: canonical runner must be ubuntu-24.04 and self-hosted must remain prohibited")
if ci.get("evidence_environment") != "github-hosted" or ci.get("final_check") != "governance":
    raise SystemExit("ERROR: CI evidence must prove github-hosted governance")

p01_gate = (ROOT / "docs/governance/P01_ENTRY_GATE.md").read_text(encoding="utf-8").lower()
for marker in ["status: **satisfied", "issue #3", "protected: true", "32540836431", "44ca19e80c5fccccebfd8d4f96dde6dc5af14bc2", "32541439589", "github-hosted"]:
    if marker.lower() not in p01_gate:
        raise SystemExit(f"ERROR: P01 entry gate missing verified historical marker: {marker}")

hardening = (ROOT / "docs/governance/REPOSITORY_HARDENING.md").read_text(encoding="utf-8")
for marker in ["Issue #3 is **closed/completed**", "Required `governance`", "Cannot force-push", "conversation", "BRANCH_PROTECTION_ADMIN_RUNBOOK.md"]:
    if marker.lower() not in hardening.lower():
        raise SystemExit(f"ERROR: hardening record missing marker: {marker}")

adr6 = (ROOT / "docs/adr/ADR-0006-temporary-p00-ci-evidence-exception.md").read_text(encoding="utf-8")
if "Expired — historical evidence only" not in adr6 or "cannot authorize a present or future CI bypass" not in adr6:
    raise SystemExit("ERROR: ADR-0006 must remain expired/historical-only")

external = manifest.get("external_distribution_gate") or {}
if external.get("tracker") != "issue:#4" or external.get("state") != "BLOCKED":
    raise SystemExit("ERROR: Issue #4 must remain the external distribution blocker")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
if "runs-on: ubuntu-24.04" not in workflow or "RUNNER_ENVIRONMENT" not in workflow or "github-hosted" not in workflow:
    raise SystemExit("ERROR: canonical governance workflow must prove hosted ubuntu execution")
if "self-hosted" in workflow or "LOCAL-WIN-" in workflow:
    raise SystemExit("ERROR: canonical governance workflow must not use local/self-hosted runners")

lock = state.get("implementation_lock") or {}
if lock.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: business-feature implementation must remain locked")
current_phase = state.get("current_phase")
phase = state.get("phase") or {}
if current_phase == "P01" and lock.get("kernel_code_authorized") is not False:
    raise SystemExit("ERROR: completed-P01 readiness checkpoint must keep kernel implementation locked")
if current_phase == "P02":
    if phase.get("state") == "active" and lock.get("kernel_code_authorized") is not True:
        raise SystemExit("ERROR: active P02 must explicitly authorize bounded kernel implementation")
    if phase.get("state") == "done":
        if lock.get("kernel_code_authorized") is not False:
            raise SystemExit("ERROR: completed P02 checkpoint must lock kernel implementation")
        p02 = phases.get("P02") or {}
        p03 = phases.get("P03") or {}
        if p02.get("state") != "done" or p02.get("active_work_package") is not None:
            raise SystemExit("ERROR: completed phases[] P02 must be done with no active package")
        if p03.get("state") != "planned":
            raise SystemExit("ERROR: P03 must remain planned until separate activation")
if current_phase == "P03":
    p02 = phases.get("P02") or {}
    p03 = phases.get("P03") or {}
    if p02.get("state") != "done" or p02.get("active_work_package") is not None:
        raise SystemExit("ERROR: P02 must remain completed before/during P03")
    if phase.get("id") != "P03" or p03.get("state") != phase.get("state"):
        raise SystemExit("ERROR: P03 current phase/state mismatch")
    if phase.get("state") == "active":
        if lock.get("kernel_code_authorized") is not True:
            raise SystemExit("ERROR: active P03 must explicitly authorize bounded kernel implementation")
        if not isinstance(state.get("current_work_package"), str) or not state.get("current_work_package").startswith("P03."):
            raise SystemExit("ERROR: active P03 must identify one current P03 package")
        if p03.get("active_work_package") != state.get("current_work_package"):
            raise SystemExit("ERROR: phases[] P03 active package mismatch")
    elif phase.get("state") == "done":
        if lock.get("kernel_code_authorized") is not False:
            raise SystemExit("ERROR: completed P03 checkpoint must lock kernel implementation")
    else:
        raise SystemExit("ERROR: P03 current phase must be active or done")
if current_phase == "P04":
    p02 = phases.get("P02") or {}
    p03 = phases.get("P03") or {}
    p04 = phases.get("P04") or {}
    if p02.get("state") != "done" or p02.get("active_work_package") is not None:
        raise SystemExit("ERROR: P02 must remain completed during P04")
    if p03.get("state") != "done" or p03.get("active_work_package") is not None:
        raise SystemExit("ERROR: P03 must remain completed before/during P04")
    if phase.get("id") != "P04" or phase.get("state") != "active" or p04.get("state") != "active":
        raise SystemExit("ERROR: P04 current phase/state mismatch")
    current = state.get("current_work_package")
    if not isinstance(current, str) or not current.startswith("P04."):
        raise SystemExit("ERROR: active P04 must identify one current P04 package")
    if p04.get("active_work_package") != current:
        raise SystemExit("ERROR: phases[] P04 active package mismatch")
    if lock.get("kernel_code_authorized") is not True:
        raise SystemExit("ERROR: active P04 must explicitly authorize bounded kernel implementation")

print("Omnexa foundation freeze / completed P01 prerequisite validation: PASS")
print("Architecture: FROZEN")
print("P00 exit: DONE")
print("P01: DONE / 12 OF 12")
print("P01 exit: SATISFIED")
print(f"Current phase: {state.get('current_phase')}")
print("Business feature code: LOCKED")