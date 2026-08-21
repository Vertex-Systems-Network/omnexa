#!/usr/bin/env python3
"""Validate that P01 is implementation-ready without prematurely authorizing kernel code."""

from __future__ import annotations

import json
import pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED_FILES = [
    "docs/roadmap/work-packages/P01.01.md",
    "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json",
    "docs/governance/P00_P01_TRANSITION_CHECKLIST.md",
    "docs/governance/LICENSING_DECISION.md",
    "docs/governance/LICENSING_DECISION_BRIEF.md",
    "docs/governance/P01_ENTRY_GATE.md",
    "docs/roadmap/STATE.json",
    ".github/workflows/governance.yml",
    ".github/workflows/main-protection-admin.yml",
]

for relative in REQUIRED_FILES:
    if not (ROOT / relative).is_file():
        raise SystemExit(f"ERROR: missing P01 preparation artifact: {relative}")

state = json.loads((ROOT / "docs/roadmap/STATE.json").read_text(encoding="utf-8"))

if state.get("current_phase") != "P00" or state.get("current_work_package") != "P00.10":
    raise SystemExit("ERROR: preparation validator expects canonical P00.10 exit-verification state")

lock = state.get("implementation_lock") or {}
if lock.get("kernel_code_authorized") is not False:
    raise SystemExit("ERROR: kernel code must remain locked during P01 preparation")
if lock.get("business_feature_code_authorized") is not False:
    raise SystemExit("ERROR: business feature code must remain locked during P01 preparation")

tracking = state.get("governance_tracking") or {}
branch = tracking.get("main_branch_protection") or {}
if branch.get("state") != "actionable_unprotected":
    raise SystemExit("ERROR: Issue #3 must remain actionable_unprotected until live protection is verified")
if branch.get("issue") != 3:
    raise SystemExit("ERROR: branch protection must remain tracked by Issue #3")
if branch.get("repository_visibility") != "public":
    raise SystemExit("ERROR: current repository visibility must be public")
if branch.get("live_protected") is not False:
    raise SystemExit("ERROR: P01 cannot remain blocked while state claims protection already verified")

ci = tracking.get("github_actions_ci") or {}
if ci.get("state") != "operational_github_hosted":
    raise SystemExit("ERROR: governance CI must be operational_github_hosted")
if ci.get("routing_mode") != "github_hosted_only":
    raise SystemExit("ERROR: governance CI routing must be github_hosted_only")
if ci.get("runner_label") != "ubuntu-24.04":
    raise SystemExit("ERROR: canonical runner label must be ubuntu-24.04")
if ci.get("self_hosted_allowed") is not False:
    raise SystemExit("ERROR: local/self-hosted governance runners must remain prohibited")
if ci.get("workflow_run") != 32537207455 or ci.get("job") != 96940269306:
    raise SystemExit("ERROR: hosted migration proof run/job mismatch")
if ci.get("evidence_environment") != "github-hosted":
    raise SystemExit("ERROR: hosted CI evidence must prove github-hosted environment")
if ci.get("final_check") != "governance":
    raise SystemExit("ERROR: canonical required check must remain governance")

phases = {item.get("id"): item for item in state.get("phases") or []}
p01 = phases.get("P01") or {}
if p01.get("state") != "blocked":
    raise SystemExit("ERROR: P01 must remain blocked during preparation")
if "issue:#3" not in (p01.get("blocked_by") or []):
    raise SystemExit("ERROR: P01 must remain blocked by Issue #3")

prep = state.get("p01_preparation") or {}
if prep.get("state") != "prepared_blocked":
    raise SystemExit("ERROR: p01_preparation.state must be prepared_blocked")
if prep.get("next_work_package") != "P01.01":
    raise SystemExit("ERROR: next prepared package must be P01.01")
if prep.get("work_package_spec") != "docs/roadmap/work-packages/P01.01.md":
    raise SystemExit("ERROR: P01.01 work-package spec path mismatch")
if prep.get("package_sequence") != "docs/roadmap/work-packages/P01_PACKAGE_SEQUENCE.json":
    raise SystemExit("ERROR: P01 package sequence path mismatch")
if prep.get("prepared_spec_count") != 12 or prep.get("mandatory_spec_count") != 12:
    raise SystemExit("ERROR: P01 preparation must report 12 / 12 specifications")
if prep.get("transition_checklist") != "docs/governance/P00_P01_TRANSITION_CHECKLIST.md":
    raise SystemExit("ERROR: transition checklist path mismatch")

package = (ROOT / "docs/roadmap/work-packages/P01.01.md").read_text(encoding="utf-8")
for marker in [
    "# P01.01 — Go Workspace / Build Skeleton",
    "State: `planned`",
    "kernel_code_authorized=true",
    "P01.02",
    "P02",
    "P03",
    "G0",
    "G1",
    "G2",
    "G7",
    "GitHub-hosted",
    "no business",
]:
    if marker.lower() not in package.lower():
        raise SystemExit(f"ERROR: P01.01 spec missing marker: {marker}")

transition = (ROOT / "docs/governance/P00_P01_TRANSITION_CHECKLIST.md").read_text(encoding="utf-8")
for marker in [
    "EG-02",
    "P00.10.state",
    "kernel_code_authorized",
    "business_feature_code_authorized",
    "P01.01",
    "governance",
    "superseded",
    "GitHub-hosted",
    "ubuntu-24.04",
]:
    if marker.lower() not in transition.lower():
        raise SystemExit(f"ERROR: transition checklist missing marker: {marker}")

workflow = (ROOT / ".github/workflows/governance.yml").read_text(encoding="utf-8")
for marker in [
    "name: governance",
    "runs-on: ubuntu-24.04",
    "RUNNER_ENVIRONMENT",
    "github-hosted",
    "python scripts/validate_governance.py",
    "python scripts/validate_development_spec.py",
    "python scripts/validate_operations_spec.py",
    "python scripts/validate_freeze_review.py",
    "python scripts/validate_p01_preparation.py",
    "python scripts/validate_p01_package_specs.py",
]:
    if marker not in workflow:
        raise SystemExit(f"ERROR: governance workflow missing GitHub-hosted marker: {marker}")
for forbidden in ["self-hosted", "LOCAL-WIN-", "runs-on: [self-hosted"]:
    if forbidden in workflow:
        raise SystemExit(f"ERROR: governance workflow must not use local/self-hosted runners: {forbidden}")

admin_workflow = (ROOT / ".github/workflows/main-protection-admin.yml").read_text(encoding="utf-8")
for marker in [
    "workflow_dispatch",
    "runs-on: ubuntu-24.04",
    "RUNNER_ENVIRONMENT",
    "github-hosted",
    "OMNEXA_GITHUB_ADMIN_TOKEN",
    "apply_main_protection.ps1",
    "verify_main_protection.ps1",
]:
    if marker not in admin_workflow:
        raise SystemExit(f"ERROR: main protection admin workflow missing hosted-admin marker: {marker}")
if "self-hosted" in admin_workflow or "LOCAL-WIN-" in admin_workflow:
    raise SystemExit("ERROR: branch-protection administration workflow must not use local/self-hosted runners")

licensing = (ROOT / "docs/governance/LICENSING_DECISION_BRIEF.md").read_text(encoding="utf-8")
for marker in [
    "Issue #4",
    "GPLv3",
    "distribution model",
    "core licensing strategy",
    "contributor IP",
    "marketplace",
    "trademark",
    "LICENSE replacement: NOT AUTHORIZED",
]:
    if marker.lower() not in licensing.lower():
        raise SystemExit(f"ERROR: licensing decision brief missing marker: {marker}")

PREMATURE_EXECUTABLE_PATHS = [
    "go.work",
    "go.mod",
    "kernel/cmd/omnexa/main.go",
]
for relative in PREMATURE_EXECUTABLE_PATHS:
    if (ROOT / relative).exists():
        raise SystemExit(
            f"ERROR: premature kernel implementation path exists while kernel_code_authorized=false: {relative}"
        )

print("Omnexa P01 preparation validation: PASS")
print("Governance runner policy: GITHUB-HOSTED ONLY / ubuntu-24.04")
print("Branch-protection admin runner policy: GITHUB-HOSTED ONLY / workflow_dispatch")
print("P01 specifications: 12 / 12 PREPARED / PLANNED")
print("P01 implementation: BLOCKED BY ISSUE #3")
print("Kernel code: LOCKED")
print("Business feature code: LOCKED")
