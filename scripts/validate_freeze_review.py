#!/usr/bin/env python3
"""Dependency-free structural validation for the P00.10 foundation freeze review."""

from __future__ import annotations

import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/governance/FOUNDATION_FREEZE_REVIEW.md",
    "docs/governance/P01_ENTRY_GATE.md",
    "docs/governance/FOUNDATION_FREEZE.json",
    "docs/contracts/governance/foundation-freeze.schema.json",
    "docs/adr/ADR-0010-foundation-architecture-freeze.md",
    "docs/governance/BRANCH_PROTECTION_ADMIN_RUNBOOK.md",
    "scripts/apply_main_protection.ps1",
    "scripts/verify_main_protection.ps1",
]

for path in REQUIRED_FILES:
    if not (ROOT / path).is_file():
        print(f"ERROR: missing freeze artifact: {path}", file=sys.stderr)
        raise SystemExit(1)

manifest = json.loads((ROOT / "docs/governance/FOUNDATION_FREEZE.json").read_text(encoding="utf-8"))
state = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))

if manifest.get("version") != "foundation-v1":
    raise SystemExit("ERROR: freeze version must be foundation-v1")
if manifest.get("architecture_status") != "FROZEN":
    raise SystemExit("ERROR: architecture_status must be FROZEN")
if manifest.get("p00_exit_status") not in {"VERIFICATION", "DONE"}:
    raise SystemExit("ERROR: invalid P00 exit status")

expected_packages = {f"P00.{i:02d}" for i in range(1, 10)}
if set(manifest.get("frozen_packages") or []) != expected_packages:
    raise SystemExit("ERROR: freeze manifest must contain exactly P00.01-P00.09")

entry = manifest.get("p01_entry_gate") or {}
if entry.get("state") != "BLOCKED":
    raise SystemExit("ERROR: P01 entry must remain BLOCKED while issue #3 is unresolved")
if entry.get("kernel_code_authorized") is not False:
    raise SystemExit("ERROR: kernel code must remain unauthorized")
if entry.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: business feature code must remain unauthorized")

entry_states = {item.get("tracker"): item.get("state") for item in entry.get("blockers") or []}
if entry_states.get("issue:#3") != "BLOCKED":
    raise SystemExit("ERROR: issue #3 must remain the P01 entry blocker")
if entry_states.get("issue:#14") != "SATISFIED":
    raise SystemExit("ERROR: issue #14 executable CI gate must be SATISFIED")

tracking = state.get("governance_tracking") or {}
main_protection = tracking.get("main_branch_protection") or {}
if main_protection.get("state") != "blocked_by_plan":
    raise SystemExit("ERROR: main branch protection must be classified blocked_by_plan until hosted entitlement changes")
if main_protection.get("issue") != 3:
    raise SystemExit("ERROR: main branch protection must remain tracked by issue #3")
if main_protection.get("repository_visibility") != "private":
    raise SystemExit("ERROR: plan-blocked evidence expects the repository to remain private")
if "HTTP 403" not in str(main_protection.get("attempt_result") or ""):
    raise SystemExit("ERROR: plan-blocked branch protection must retain HTTP 403 evidence")

ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_self_hosted":
    raise SystemExit("ERROR: executable CI gate must remain operational_self_hosted")
if ci.get("routing_mode") != "any_available_windows_x64_self_hosted":
    raise SystemExit("ERROR: governance CI must route to any available Windows/X64 self-hosted runner")
if ci.get("final_check") != "governance":
    raise SystemExit("ERROR: final required CI check must be governance")

p01_gate = (ROOT / "docs/governance/P01_ENTRY_GATE.md").read_text(encoding="utf-8")
for marker in [
    "EG-02",
    "Issue #3",
    "BLOCKED_BY_PLAN",
    "HTTP 403",
    "Issue #14",
    "SATISFIED",
    "any available Windows/X64 self-hosted runner",
    "32535324900",
    "LOCAL-WIN-02",
]:
    if marker.lower() not in p01_gate.lower():
        raise SystemExit(f"ERROR: P01 entry gate missing reconciliation marker: {marker}")

external = manifest.get("external_distribution_gate") or {}
if external.get("tracker") != "issue:#4" or external.get("state") != "BLOCKED":
    raise SystemExit("ERROR: issue #4 must remain external-distribution blocker")

review = (ROOT / "docs/governance/FOUNDATION_FREEZE_REVIEW.md").read_text(encoding="utf-8")
for marker in ["ACCEPTED FOR FREEZE", "Issue #3", "Issue #14", "Issue #4", "P01 implementation-entry blockers"]:
    if marker not in review:
        raise SystemExit(f"ERROR: freeze review missing historical marker: {marker}")

hardening = (ROOT / "docs/governance/REPOSITORY_HARDENING.md").read_text(encoding="utf-8")
for marker in [
    "BRANCH_PROTECTION_ADMIN_RUNBOOK.md",
    "apply_main_protection.ps1",
    "verify_main_protection.ps1",
    "Issue #3 remains open",
]:
    if marker not in hardening:
        raise SystemExit(f"ERROR: repository hardening missing admin-tooling marker: {marker}")

runbook = (ROOT / "docs/governance/BRANCH_PROTECTION_ADMIN_RUNBOOK.md").read_text(encoding="utf-8")
for marker in [
    "HTTP 403",
    "product-plan entitlement failure",
    "Issue #4",
    "Do not keep retrying",
]:
    if marker not in runbook:
        raise SystemExit(f"ERROR: branch protection runbook missing plan-limitation marker: {marker}")

apply_script = (ROOT / "scripts/apply_main_protection.ps1").read_text(encoding="utf-8")
for marker in [
    "required_status_checks",
    "governance",
    "enforce_admins",
    "required_pull_request_reviews",
    "required_conversation_resolution",
    "allow_force_pushes",
    "allow_deletions",
    "OMNEXA_GITHUB_ADMIN_TOKEN",
]:
    if marker not in apply_script:
        raise SystemExit(f"ERROR: apply_main_protection.ps1 missing policy marker: {marker}")

verify_script = (ROOT / "scripts/verify_main_protection.ps1").read_text(encoding="utf-8")
for marker in [
    "protected=true",
    "required status check",
    "conversation resolution",
    "force pushes are blocked",
    "branch deletion is blocked",
    "administrators",
]:
    if marker not in verify_script:
        raise SystemExit(f"ERROR: verify_main_protection.ps1 missing verification marker: {marker}")

print("Omnexa P00.10 foundation freeze review validation: PASS")
print("Architecture: FROZEN")
print("Executable CI gate: SATISFIED ON ANY AVAILABLE WINDOWS/X64 SELF-HOSTED RUNNER")
print("Branch-protection admin tooling: PRESENT")
print("Branch-protection hosted entitlement: BLOCKED_BY_PLAN (HTTP 403)")
print("P00 exit: VERIFICATION")
print("P01 entry: BLOCKED BY ISSUE #3")
