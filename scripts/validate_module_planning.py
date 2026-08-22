#!/usr/bin/env python3
"""Validate the mandatory future module/submodule planning catalog.

This validator verifies planning completeness only. It never authorizes a future
phase; executable authority remains docs/roadmap/STATE.json.
"""

from __future__ import annotations

import json
import pathlib
import re

ROOT = pathlib.Path(__file__).resolve().parents[1]
CATALOG = ROOT / "docs/roadmap/modules/SUBMODULE_CATALOG.json"
INDEX = ROOT / "docs/roadmap/modules/README.md"
DOSSIER_STANDARD = ROOT / "docs/roadmap/modules/DOSSIER_STANDARD.md"
EXECUTION_PROFILES = ROOT / "docs/roadmap/modules/EXECUTION_PROFILES.md"
EXECUTION_PROFILE_MAP = ROOT / "docs/roadmap/modules/EXECUTION_PROFILE_MAP.json"
BLUEPRINT = ROOT / "docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md"
AI_PROTOCOL = ROOT / "docs/governance/AI_MODULE_EXECUTION_PROTOCOL.md"
AI_POLICY = ROOT / "docs/governance/AI_EXECUTION_POLICY.md"
AGENTS = ROOT / "AGENTS.md"
MODULE_STANDARD = ROOT / "docs/architecture/MODULE_STANDARD.md"

for path in (
    CATALOG,
    INDEX,
    DOSSIER_STANDARD,
    EXECUTION_PROFILES,
    EXECUTION_PROFILE_MAP,
    BLUEPRINT,
    AI_PROTOCOL,
    AI_POLICY,
    AGENTS,
    MODULE_STANDARD,
):
    if not path.is_file():
        raise SystemExit(f"ERROR: missing module planning artifact: {path.relative_to(ROOT)}")

catalog = json.loads(CATALOG.read_text(encoding="utf-8"))
if catalog.get("schema_version") != 1:
    raise SystemExit("ERROR: unsupported SUBMODULE_CATALOG schema_version")
if catalog.get("status") != "mandatory_planning_baseline":
    raise SystemExit("ERROR: submodule catalog status must remain mandatory_planning_baseline")
if catalog.get("implementation_authority") != "docs/roadmap/STATE.json":
    raise SystemExit("ERROR: future planning must not replace STATE.json implementation authority")
if catalog.get("execution_blueprint") != "docs/roadmap/MODULE_SUBMODULE_EXECUTION_BLUEPRINT.md":
    raise SystemExit("ERROR: submodule catalog execution blueprint mismatch")

phases = catalog.get("phases")
if not isinstance(phases, list):
    raise SystemExit("ERROR: submodule catalog phases must be a list")

expected = [f"P{i:02d}" for i in range(2, 28)]
actual = [phase.get("id") for phase in phases]
if actual != expected:
    raise SystemExit(f"ERROR: module planning phases must be exactly {expected}; got {actual}")

phase_ids: set[str] = set()
submodule_ids: set[str] = set()
submodule_pattern = re.compile(r"^P(?:0[2-9]|1[0-9]|2[0-7])\.[A-Z][A-Z0-9]*$")

for phase in phases:
    phase_id = phase.get("id")
    if phase_id in phase_ids:
        raise SystemExit(f"ERROR: duplicate phase in submodule catalog: {phase_id}")
    phase_ids.add(phase_id)

    name = phase.get("name")
    dossier_rel = phase.get("dossier")
    submodules = phase.get("submodules")
    if not isinstance(name, str) or not name.strip():
        raise SystemExit(f"ERROR: {phase_id} missing phase name")
    if not isinstance(dossier_rel, str) or not dossier_rel:
        raise SystemExit(f"ERROR: {phase_id} missing dossier")
    dossier = ROOT / dossier_rel
    if not dossier.is_file():
        raise SystemExit(f"ERROR: {phase_id} dossier missing: {dossier_rel}")
    dossier_text = dossier.read_text(encoding="utf-8")

    if not isinstance(submodules, list) or not submodules:
        raise SystemExit(f"ERROR: {phase_id} must pre-plan at least one submodule")

    for item in submodules:
        if not isinstance(item, list) or len(item) != 2:
            raise SystemExit(f"ERROR: {phase_id} submodule entry must be [id, name]")
        sub_id, sub_name = item
        if not isinstance(sub_id, str) or not submodule_pattern.match(sub_id):
            raise SystemExit(f"ERROR: invalid submodule identifier: {sub_id!r}")
        if not sub_id.startswith(phase_id + "."):
            raise SystemExit(f"ERROR: submodule {sub_id} does not belong to {phase_id}")
        if sub_id in submodule_ids:
            raise SystemExit(f"ERROR: duplicate submodule identifier: {sub_id}")
        submodule_ids.add(sub_id)
        if not isinstance(sub_name, str) or not sub_name.strip():
            raise SystemExit(f"ERROR: {sub_id} missing name")
        if sub_id not in dossier_text:
            raise SystemExit(f"ERROR: {sub_id} is cataloged but not documented in {dossier_rel}")

profile_map = json.loads(EXECUTION_PROFILE_MAP.read_text(encoding="utf-8"))
if profile_map.get("schema_version") != 1:
    raise SystemExit("ERROR: unsupported EXECUTION_PROFILE_MAP schema_version")
if profile_map.get("profiles_document") != "docs/roadmap/modules/EXECUTION_PROFILES.md":
    raise SystemExit("ERROR: execution profile document mismatch")
if profile_map.get("implementation_authority") != "docs/roadmap/STATE.json":
    raise SystemExit("ERROR: execution profile map must not replace STATE.json authorization")

valid_profiles = {f"EP{i:02d}" for i in range(1, 13)}
profile_text = EXECUTION_PROFILES.read_text(encoding="utf-8")
for profile in sorted(valid_profiles):
    if profile not in profile_text:
        raise SystemExit(f"ERROR: execution profile is mapped but undocumented: {profile}")

def validate_profile_list(owner: str, values: object) -> None:
    if not isinstance(values, list) or not values:
        raise SystemExit(f"ERROR: {owner} requires at least one execution profile")
    if len(values) > 3:
        raise SystemExit(f"ERROR: {owner} has too many execution profiles; refine the planning boundary")
    seen: set[str] = set()
    for profile in values:
        if profile not in valid_profiles:
            raise SystemExit(f"ERROR: {owner} references unknown execution profile: {profile}")
        if profile in seen:
            raise SystemExit(f"ERROR: {owner} repeats execution profile: {profile}")
        seen.add(profile)

phase_defaults = profile_map.get("phase_defaults")
if not isinstance(phase_defaults, dict) or list(phase_defaults.keys()) != expected:
    raise SystemExit("ERROR: execution profile phase defaults must cover P02-P27 in exact order")
for phase_id, profiles in phase_defaults.items():
    validate_profile_list(phase_id, profiles)

overrides = profile_map.get("submodule_overrides")
if not isinstance(overrides, dict):
    raise SystemExit("ERROR: execution profile submodule_overrides must be an object")
for sub_id, profiles in overrides.items():
    if sub_id not in submodule_ids:
        raise SystemExit(f"ERROR: execution profile override references unknown submodule: {sub_id}")
    validate_profile_list(sub_id, profiles)

index_text = INDEX.read_text(encoding="utf-8")
for required in [
    "DOSSIER_STANDARD.md",
    "EXECUTION_PROFILES.md",
    "EXECUTION_PROFILE_MAP.json",
    "FOUNDATION_P02_P06.md",
    "CORE_BUSINESS_P07_P15.md",
    "PLATFORM_P16_P18.md",
    "INTELLIGENCE_P19_P20.md",
    "ECOSYSTEM_P21_P22.md",
    "GLOBAL_ENTERPRISE_P23_P25.md",
    "INDUSTRY_AUTONOMY_P26_P27.md",
    "SUBMODULE_CATALOG.json",
]:
    if required not in index_text:
        raise SystemExit(f"ERROR: module planning index missing reference: {required}")

required_cross_document_markers = {
    AI_POLICY: [
        "SUBMODULE_CATALOG.json",
        "AI_MODULE_EXECUTION_PROTOCOL.md",
        "next incomplete authorized task",
    ],
    AGENTS: [
        "SUBMODULE_CATALOG.json",
        "AI_MODULE_EXECUTION_PROTOCOL.md",
        "validate_module_planning.py",
    ],
    MODULE_STANDARD: [
        "SUBMODULE_CATALOG.json",
        "Option, setting and policy contract",
        "Flow contract",
    ],
    AI_PROTOCOL: [
        "SUBMODULE_CATALOG.json",
        "No-replanning rule",
        "Required submodule task card",
    ],
    DOSSIER_STANDARD: [
        "Primary flows",
        "Options/settings/policies",
        "S01",
        "S10",
    ],
    EXECUTION_PROFILES: [
        "EP01",
        "EP06",
        "EP11",
        "Profile composition rule",
    ],
}
for path, markers in required_cross_document_markers.items():
    text = path.read_text(encoding="utf-8")
    for marker in markers:
        if marker not in text:
            raise SystemExit(f"ERROR: {path.relative_to(ROOT)} missing planning marker: {marker}")

print("Omnexa module/submodule planning validation: PASS")
print(f"Planned future phases: {len(phases)}")
print(f"Preplanned submodules/families: {len(submodule_ids)}")
print(f"Execution profiles: {len(valid_profiles)}")
print(f"Submodule profile overrides: {len(overrides)}")
print("Architecture/flow/options/task-profile system: ENFORCED")
print("Implementation authority: docs/roadmap/STATE.json")
