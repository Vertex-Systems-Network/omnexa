#!/usr/bin/env python3
"""Fail-fast consistency checks for Omnexa architecture/governance state.

This is intentionally dependency-free so it can run before the application toolchain exists.
It is not the final P00.07 CI/test architecture.
"""

from __future__ import annotations

import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE_PATH = ROOT / "docs/roadmap/STATE.json"
STATUS_PATH = ROOT / "docs/roadmap/STATUS.md"

REQUIRED_FILES = [
    "AGENTS.md",
    "CONTRIBUTING.md",
    "SECURITY.md",
    ".github/CODEOWNERS",
    ".github/pull_request_template.md",
    "docs/governance/PRODUCT_CONSTITUTION.md",
    "docs/governance/AI_EXECUTION_POLICY.md",
    "docs/governance/CHANGE_CONTROL.md",
    "docs/governance/DEFINITION_OF_DONE.md",
    "docs/governance/REPOSITORY_HARDENING.md",
    "docs/governance/LICENSING_DECISION.md",
    "docs/architecture/SYSTEM_ARCHITECTURE.md",
    "docs/architecture/MODULE_STANDARD.md",
    "docs/architecture/GLOSSARY.md",
    "docs/architecture/NAMING_STANDARD.md",
    "docs/architecture/DOMAIN_OWNERSHIP.md",
    "docs/architecture/DEPENDENCY_MATRIX.md",
    "docs/roadmap/MASTER_PLAN.md",
    "docs/roadmap/STATUS.md",
    "docs/roadmap/STATE.json",
    "docs/roadmap/EXECUTION_LEDGER.md",
    "docs/adr/ADR-0001-platform-architecture-baseline.md",
    "docs/adr/TEMPLATE.md",
]

ALLOWED_STATES = {
    "planned",
    "ready",
    "active",
    "blocked",
    "verification",
    "done",
    "superseded",
}


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


def require_files() -> None:
    missing = [path for path in REQUIRED_FILES if not (ROOT / path).is_file()]
    if missing:
        fail("missing required governance files: " + ", ".join(missing))


def load_state() -> dict:
    try:
        state = json.loads(STATE_PATH.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        fail(f"STATE.json is unreadable/invalid: {exc}")
    if not isinstance(state, dict):
        fail("STATE.json root must be an object")
    return state


def validate_state(state: dict) -> None:
    if state.get("schema_version") != 1:
        fail("unsupported or missing STATE.json schema_version")

    vocabulary = set(state.get("state_vocabulary", []))
    if vocabulary != ALLOWED_STATES:
        fail("state_vocabulary does not exactly match the canonical state set")

    current_phase = state.get("current_phase")
    phase = state.get("phase") or {}
    if phase.get("id") != current_phase:
        fail("current_phase and phase.id disagree")
    if phase.get("state") != "active":
        fail("current phase must be active")

    work_packages = phase.get("work_packages") or []
    if phase.get("mandatory_work_packages") != len(work_packages):
        fail("mandatory_work_packages does not match work_packages length")

    ids = [pkg.get("id") for pkg in work_packages]
    if len(ids) != len(set(ids)):
        fail("duplicate work-package IDs detected")

    for pkg in work_packages:
        if pkg.get("state") not in ALLOWED_STATES:
            fail(f"invalid state for work package {pkg.get('id')}")
        for dependency in pkg.get("depends_on", []):
            if dependency not in ids:
                fail(f"unknown dependency {dependency} in {pkg.get('id')}")

    done_count = sum(pkg.get("state") == "done" for pkg in work_packages)
    if phase.get("done_work_packages") != done_count:
        fail("done_work_packages does not match actual done count")

    active_packages = [pkg for pkg in work_packages if pkg.get("state") == "active"]
    if len(active_packages) != 1:
        fail("foundation execution requires exactly one active work package")

    current_package = state.get("current_work_package")
    if active_packages[0].get("id") != current_package:
        fail("current_work_package is not the active work package")

    phase_rows = state.get("phases") or []
    phase_ids = [item.get("id") for item in phase_rows]
    if len(phase_ids) != len(set(phase_ids)):
        fail("duplicate phase IDs detected")

    active_phases = [item for item in phase_rows if item.get("state") == "active"]
    if len(active_phases) != 1 or active_phases[0].get("id") != current_phase:
        fail("phases[] must contain exactly one active phase matching current_phase")

    if current_phase == "P00":
        lock = state.get("implementation_lock") or {}
        if lock.get("business_feature_code_authorized") is not False:
            fail("business feature code must remain locked during P00")
        if lock.get("kernel_code_authorized") is not False:
            fail("kernel code must remain locked during P00")


def validate_status(state: dict) -> None:
    try:
        status = STATUS_PATH.read_text(encoding="utf-8")
    except OSError as exc:
        fail(f"STATUS.md unreadable: {exc}")

    current_package = state["current_work_package"]
    phase = state["phase"]
    expected_progress = f"{phase['done_work_packages']} / {phase['mandatory_work_packages']} done"

    if current_package not in status:
        fail(f"STATUS.md does not mention current work package {current_package}")
    if expected_progress not in status:
        fail(f"STATUS.md does not contain canonical progress '{expected_progress}'")


def main() -> int:
    require_files()
    state = load_state()
    validate_state(state)
    validate_status(state)
    print("Omnexa governance validation: PASS")
    print(f"Current phase: {state['current_phase']}")
    print(f"Current work package: {state['current_work_package']}")
    print(
        "Progress: "
        f"{state['phase']['done_work_packages']} / "
        f"{state['phase']['mandatory_work_packages']} done"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
