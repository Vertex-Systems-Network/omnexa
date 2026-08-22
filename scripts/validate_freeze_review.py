#!/usr/bin/env python3
"""Validate frozen P00 evidence and the completed P00 -> P01 handoff."""

from __future__ import annotations

import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/governance/FOUNDATION_FREEZE_REVIEW.md",
    "docs/governance/P01_ENTRY_GATE.md",
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
    raise SystemExit("ERROR: P00 exit must be DONE after transition")

expected_frozen = {f"P00.{i:02d}" for i in range(1, 10)}
if set(manifest.get("frozen_packages") or []) != expected_frozen:
    raise SystemExit("ERROR: frozen architecture set must remain exactly P00.01-P00.09")

entry = manifest.get("p01_entry_gate") or {}
if entry.get("state") != "SATISFIED":
    raise SystemExit("ERROR: P01 entry gate must be SATISFIED")
if entry.get("kernel_code_authorized") is not True:
    raise SystemExit("ERROR: kernel code must be authorized after P01 activation")
if entry.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: business feature code must remain unauthorized")

controls = {item.get("tracker"): item.get("state") for item in entry.get("blockers") or []}
if controls.get("issue:#3") != "SATISFIED" or controls.get("issue:#14") != "SATISFIED":
    raise SystemExit("ERROR: EG-02/Issue #3 and EG-03/Issue #14 must be SATISFIED")

if state.get("current_phase") != "P01" or state.get("current_work_package") != "P01.01":
    raise SystemExit("ERROR: canonical transition state must be P01 / P01.01")
lock = state.get("implementation_lock") or {}
if lock.get("kernel_code_authorized") is not True or lock.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: transition implementation locks are inconsistent")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "verified_protected" or branch.get("live_protected") is not True:
    raise SystemExit("ERROR: main protection must be verified_protected / live_protected=true")
if branch.get("issue") != 3 or branch.get("issue_state") != "closed":
    raise SystemExit("ERROR: main protection must preserve closed Issue #3 tracking")
if branch.get("required_check") != "governance":
    raise SystemExit("ERROR: required status check must be governance")
if branch.get("repository_visibility") != "public":
    raise SystemExit("ERROR: repository visibility must be public")
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
for marker in ["status: **satisfied", "issue #3", "protected: true", "32540836431", "44ca19e80c5fccccebfd8d4f96dde6dc5af14bc2", "32541439589", "github-hosted", "p01.01"]:
    if marker.lower() not in p01_gate:
        raise SystemExit(f"ERROR: P01 entry gate missing verified marker: {marker}")

hardening = (ROOT / "docs/governance/REPOSITORY_HARDENING.md").read_text(encoding="utf-8")
for marker in ["Issue #3 is **closed/completed**", "Required `governance`", "Cannot force-push", "conversation", "BRANCH_PROTECTION_ADMIN_RUNBOOK.md"]:
    if marker.lower() not in hardening.lower():
        raise SystemExit(f"ERROR: hardening record missing marker: {marker}")

adr6 = (ROOT / "docs/adr/ADR-0006-temporary-p00-ci-evidence-exception.md").read_text(encoding="utf-8")
if "Expired — historical evidence only" not in adr6 or "cannot authorize a present or future CI bypass" not in adr6:
    raise SystemExit("ERROR: ADR-0006 must be expired/historical-only")

external = manifest.get("external_distribution_gate") or {}
if external.get("tracker") != "issue:#4" or external.get("state") != "BLOCKED":
    raise SystemExit("ERROR: Issue #4 must remain the external distribution blocker")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
if "runs-on: ubuntu-24.04" not in workflow or "RUNNER_ENVIRONMENT" not in workflow or "github-hosted" not in workflow:
    raise SystemExit("ERROR: canonical governance workflow must prove hosted ubuntu execution")
if "self-hosted" in workflow or "LOCAL-WIN-" in workflow:
    raise SystemExit("ERROR: canonical governance workflow must not use local/self-hosted runners")

print("Omnexa foundation freeze / P00 exit validation: PASS")
print("Architecture: FROZEN")
print("P00 exit: DONE")
print("P01 entry: SATISFIED")
print("P01.01: ACTIVE")
print("Kernel code: AUTHORIZED FOR ACTIVE P01 PACKAGE")
print("Business feature code: LOCKED")
