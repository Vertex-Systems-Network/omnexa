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
]

for path in REQUIRED_FILES:
    if not (ROOT / path).is_file():
        print(f"ERROR: missing freeze artifact: {path}", file=sys.stderr)
        raise SystemExit(1)

manifest = json.loads((ROOT / "docs/governance/FOUNDATION_FREEZE.json").read_text(encoding="utf-8"))

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

p01_gate = (ROOT / "docs/governance/P01_ENTRY_GATE.md").read_text(encoding="utf-8")
for marker in [
    "EG-02",
    "Issue #3",
    "EG-03",
    "SATISFIED",
    "LOCAL-WIN-4",
    "32528329184",
    "1a14362e2ed52a20d66cec6f28b93a2ee457f9a9",
]:
    if marker not in p01_gate:
        raise SystemExit(f"ERROR: P01 entry gate missing reconciliation marker: {marker}")

external = manifest.get("external_distribution_gate") or {}
if external.get("tracker") != "issue:#4" or external.get("state") != "BLOCKED":
    raise SystemExit("ERROR: issue #4 must remain external-distribution blocker")

review = (ROOT / "docs/governance/FOUNDATION_FREEZE_REVIEW.md").read_text(encoding="utf-8")
for marker in ["ACCEPTED FOR FREEZE", "Issue #3", "Issue #14", "Issue #4", "P01 implementation-entry blockers"]:
    if marker not in review:
        raise SystemExit(f"ERROR: freeze review missing historical marker: {marker}")

print("Omnexa P00.10 foundation freeze review validation: PASS")
print("Architecture: FROZEN")
print("Executable CI gate: SATISFIED ON LOCAL-WIN-4")
print("P00 exit: VERIFICATION")
print("P01 entry: BLOCKED BY ISSUE #3")
