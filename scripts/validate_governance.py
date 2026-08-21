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
    "docs/architecture/IDENTIFIER_STANDARD.md",
    "docs/architecture/MONEY_STANDARD.md",
    "docs/architecture/TIME_STANDARD.md",
    "docs/architecture/LOCALE_STANDARD.md",
    "docs/architecture/ERROR_STANDARD.md",
    "docs/roadmap/MASTER_PLAN.md",
    "docs/roadmap/STATUS.md",
    "docs/roadmap/STATE.json",
    "docs/roadmap/EXECUTION_LEDGER.md",
    "docs/adr/ADR-0001-platform-architecture-baseline.md",
    "docs/adr/TEMPLATE.md",
]

P00_03_EVIDENCE = {
    "docs/architecture/IDENTIFIER_STANDARD.md",
    "docs/architecture/MONEY_STANDARD.md",
    "docs/architecture/TIME_STANDARD.md",
    "docs/architecture/LOCALE_STANDARD.md",
    "docs/architecture/ERROR_STANDARD.md",
}

P00_03_MARKERS = {
    "docs/architecture/IDENTIFIER_STANDARD.md": ["UUIDv7", "tenant_id", "PostgreSQL's native `uuid` type"],
    "docs/architecture/MONEY_STANDARD.md": ["NUMERIC(38,18)", "ISO 4217", "round half to even"],
    "docs/architecture/TIME_STANDARD.md": ["timestamptz", "IANA timezone", "business date"],
    "docs/architecture/LOCALE_STANDARD.md": ["BCP 47", "ISO 3166-1 alpha-2", "RTL"],
    "docs/architecture/ERROR_STANDARD.md": ["stable machine error code", "request_id", "retryable"],
}

ALLOWED_STATES = {
    "planned",
    "ready",
    "active",
    "blocked",
    "verification",
    "done",
    "superseded",
}

DEPENDENCY_SATISFIED_STATES = {"ready", "active", "verification", "done"}


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


def validate_p00_03_standards() -> None:
    for path, markers in P00_03_MARKERS.items():
        text = (ROOT / path).read_text(encoding="utf-8")
        for marker in markers:
            if marker not in text:
                fail(f"{path} is missing canonical marker: {marker}")


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

    packages_by_id = {pkg.get("id"): pkg for pkg in work_packages}

    for pkg in work_packages:
        pkg_id = pkg.get("id")
        pkg_state = pkg.get("state")
        if pkg_state not in ALLOWED_STATES:
            fail(f"invalid state for work package {pkg_id}")

        for dependency in pkg.get("depends_on", []):
            if dependency not in ids:
                fail(f"unknown dependency {dependency} in {pkg_id}")
            if pkg_state in DEPENDENCY_SATISFIED_STATES and packages_by_id[dependency].get("state") != "done":
                fail(f"{pkg_id} is {pkg_state} before dependency {dependency} is done")

        if pkg_state == "done":
            evidence = pkg.get("evidence") or []
            if not evidence:
                fail(f"done work package {pkg_id} has no evidence")
            missing_evidence = [path for path in evidence if not (ROOT / path).is_file()]
            if missing_evidence:
                fail(f"done work package {pkg_id} references missing evidence: {', '.join(missing_evidence)}")

    p00_03 = packages_by_id.get("P00.03")
    if p00_03 and p00_03.get("state") == "done":
        if not P00_03_EVIDENCE.issubset(set(p00_03.get("evidence") or [])):
            fail("P00.03 done state is missing one or more mandatory foundation-convention evidence files")

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
    validate_p00_03_standards()
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
