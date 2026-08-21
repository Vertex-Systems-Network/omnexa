#!/usr/bin/env python3
"""Fail-fast consistency checks for Omnexa architecture/governance state.

Dependency-free by design so it can run before the application toolchain exists.
The canonical quality semantics are defined by P00.07 and are CI-provider independent.
"""

from __future__ import annotations

import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]
STATE_PATH = ROOT / "docs/roadmap/STATE.json"
STATUS_PATH = ROOT / "docs/roadmap/STATUS.md"
CLASSIFICATION_SCHEMA_PATH = ROOT / "docs/contracts/security/data-classification.schema.json"
QUALITY_SCHEMA_PATH = ROOT / "docs/contracts/quality/quality-gates.schema.json"

REQUIRED_FILES = [
    "AGENTS.md",
    "CONTRIBUTING.md",
    "SECURITY.md",
    ".github/CODEOWNERS",
    ".github/pull_request_template.md",
    ".github/workflows/governance.yml",
    "docs/governance/PRODUCT_CONSTITUTION.md",
    "docs/governance/AI_EXECUTION_POLICY.md",
    "docs/governance/CHANGE_CONTROL.md",
    "docs/governance/DEFINITION_OF_DONE.md",
    "docs/governance/REPOSITORY_HARDENING.md",
    "docs/governance/LICENSING_DECISION.md",
    "docs/governance/CI_EVIDENCE_EXCEPTION_2026-08-22.md",
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
    "docs/architecture/API_STANDARD.md",
    "docs/architecture/EVENT_STANDARD.md",
    "docs/security/SECURITY_STANDARD.md",
    "docs/security/DATA_CLASSIFICATION.md",
    "docs/security/SECURITY_CONTROL_MATRIX.md",
    "docs/quality/TESTING_STANDARD.md",
    "docs/quality/CI_STANDARD.md",
    "docs/quality/RELEASE_STANDARD.md",
    "docs/quality/QUALITY_GATE_MATRIX.md",
    "docs/contracts/http/openapi-template.yaml",
    "docs/contracts/events/event-envelope.schema.json",
    "docs/contracts/security/data-classification.schema.json",
    "docs/contracts/quality/quality-gates.schema.json",
    "docs/roadmap/MASTER_PLAN.md",
    "docs/roadmap/STATUS.md",
    "docs/roadmap/STATE.json",
    "docs/roadmap/EXECUTION_LEDGER.md",
    "docs/adr/ADR-0001-platform-architecture-baseline.md",
    "docs/adr/ADR-0002-foundation-data-conventions.md",
    "docs/adr/ADR-0003-http-api-contract-baseline.md",
    "docs/adr/ADR-0004-event-contract-baseline.md",
    "docs/adr/ADR-0005-security-data-classification-baseline.md",
    "docs/adr/ADR-0006-temporary-p00-ci-evidence-exception.md",
    "docs/adr/ADR-0007-testing-ci-release-baseline.md",
    "docs/adr/TEMPLATE.md",
]

EVIDENCE_REQUIREMENTS = {
    "P00.03": {
        "docs/architecture/IDENTIFIER_STANDARD.md",
        "docs/architecture/MONEY_STANDARD.md",
        "docs/architecture/TIME_STANDARD.md",
        "docs/architecture/LOCALE_STANDARD.md",
        "docs/architecture/ERROR_STANDARD.md",
        "docs/adr/ADR-0002-foundation-data-conventions.md",
    },
    "P00.04": {
        "docs/architecture/API_STANDARD.md",
        "docs/contracts/http/openapi-template.yaml",
        "docs/adr/ADR-0003-http-api-contract-baseline.md",
    },
    "P00.05": {
        "docs/architecture/EVENT_STANDARD.md",
        "docs/contracts/events/event-envelope.schema.json",
        "docs/adr/ADR-0004-event-contract-baseline.md",
    },
    "P00.06": {
        "docs/security/SECURITY_STANDARD.md",
        "docs/security/DATA_CLASSIFICATION.md",
        "docs/security/SECURITY_CONTROL_MATRIX.md",
        "docs/contracts/security/data-classification.schema.json",
        "docs/adr/ADR-0005-security-data-classification-baseline.md",
    },
    "P00.07": {
        "docs/quality/TESTING_STANDARD.md",
        "docs/quality/CI_STANDARD.md",
        "docs/quality/RELEASE_STANDARD.md",
        "docs/quality/QUALITY_GATE_MATRIX.md",
        "docs/contracts/quality/quality-gates.schema.json",
        "docs/adr/ADR-0007-testing-ci-release-baseline.md",
    },
}

MARKER_REQUIREMENTS = {
    "P00.03": {
        "docs/architecture/IDENTIFIER_STANDARD.md": ["UUIDv7", "tenant_id", "PostgreSQL's native `uuid` type"],
        "docs/architecture/MONEY_STANDARD.md": ["NUMERIC(38,18)", "ISO 4217", "round half to even"],
        "docs/architecture/TIME_STANDARD.md": ["timestamptz", "IANA timezone", "business date"],
        "docs/architecture/LOCALE_STANDARD.md": ["BCP 47", "ISO 3166-1 alpha-2", "RTL"],
        "docs/architecture/ERROR_STANDARD.md": ["stable machine error code", "request_id", "retryable"],
    },
    "P00.04": {
        "docs/architecture/API_STANDARD.md": [
            "OpenAPI Specification 3.2.0",
            "/api/v{major}/{domain}/{resources}",
            "Idempotency-Key",
            "ETag",
            "page_cursor",
            "application/problem+json",
        ],
        "docs/contracts/http/openapi-template.yaml": [
            "openapi: 3.2.0",
            "Money:",
            "Problem:",
            "IdempotencyKey:",
            "PageCursor:",
        ],
    },
    "P00.05": {
        "docs/architecture/EVENT_STANDARD.md": [
            "CloudEvents-compatible",
            "at least once",
            "transactional outbox",
            "inbox/deduplication",
            "subjectsequence",
            "dead-letter",
            "Replay",
        ],
        "docs/contracts/events/event-envelope.schema.json": [
            '"specversion"',
            '"tenantid"',
            '"correlationid"',
            '"causationid"',
            '"subjectsequence"',
        ],
    },
    "P00.06": {
        "docs/security/DATA_CLASSIFICATION.md": [
            "PUBLIC",
            "INTERNAL",
            "CONFIDENTIAL",
            "RESTRICTED",
            "AI and model handling",
            "Lower environments",
            "Retention",
            "Deletion and erasure",
        ],
        "docs/security/SECURITY_STANDARD.md": [
            "Zero implicit trust",
            "Authorization model",
            "Tenant isolation",
            "Secrets management",
            "Audit logging",
            "Impersonation and support access",
            "Module and marketplace security",
            "AI",
        ],
        "docs/security/SECURITY_CONTROL_MATRIX.md": [
            "Deny-by-default areas",
            "cross-tenant access",
            "AI tool execution",
            "Support impersonation",
            "Security exceptions",
        ],
    },
    "P00.07": {
        "docs/quality/TESTING_STANDARD.md": [
            "Security/negative tests",
            "Migration tests",
            "Module lifecycle tests",
            "Flaky test policy",
            "BLOCKED",
        ],
        "docs/quality/CI_STANDARD.md": [
            "Provider-independent rule",
            "verify:all",
            "Fail-closed semantics",
            "least-privilege",
            "CI outage / quota behavior",
        ],
        "docs/quality/RELEASE_STANDARD.md": [
            "MAJOR.MINOR.PATCH",
            "Build once, promote",
            "SBOM",
            "Rollback and forward fix",
            "No production release during P00",
        ],
        "docs/quality/QUALITY_GATE_MATRIX.md": [
            "G0 — Governance",
            "G5 — Security/tenancy",
            "G8 — Supply chain/release",
            "PASS",
            "BLOCKED",
        ],
    },
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
EXPECTED_CLASSIFICATIONS = {"PUBLIC", "INTERNAL", "CONFIDENTIAL", "RESTRICTED"}
EXPECTED_GATE_CLASSES = {f"G{i}" for i in range(9)}
EXPECTED_EVIDENCE_STATES = {"PASS", "FAIL", "BLOCKED", "NOT RUN", "N/A"}


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


def require_files() -> None:
    missing = [path for path in REQUIRED_FILES if not (ROOT / path).is_file()]
    if missing:
        fail("missing required governance files: " + ", ".join(missing))


def load_json(path: pathlib.Path, label: str) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        fail(f"{label} is unreadable/invalid JSON: {exc}")
    if not isinstance(value, dict):
        fail(f"{label} root must be an object")
    return value


def validate_markers() -> None:
    for package_id, files in MARKER_REQUIREMENTS.items():
        for path, markers in files.items():
            text = (ROOT / path).read_text(encoding="utf-8")
            for marker in markers:
                if marker not in text:
                    fail(f"{package_id}: {path} is missing canonical concept marker: {marker}")


def validate_classification_schema() -> None:
    schema = load_json(CLASSIFICATION_SCHEMA_PATH, "data-classification schema")
    properties = schema.get("properties") or {}
    classification = properties.get("classification") or {}
    if set(classification.get("enum") or []) != EXPECTED_CLASSIFICATIONS:
        fail("data-classification schema must define exactly PUBLIC/INTERNAL/CONFIDENTIAL/RESTRICTED")

    required = set(schema.get("required") or [])
    for field in {"owner", "resource", "classification", "tenant_scope", "logging", "export", "retention_policy", "ai_eligibility"}:
        if field not in required:
            fail(f"data-classification schema is missing required field: {field}")

    handling_tags = ((properties.get("handling_tags") or {}).get("items") or {}).get("enum") or []
    for tag in {"PII", "AUTH_SECRET", "CRYPTO_KEY", "PAYMENT_SENSITIVE", "MODEL_INPUT", "MODEL_OUTPUT"}:
        if tag not in handling_tags:
            fail(f"data-classification schema is missing canonical handling tag: {tag}")


def validate_quality_schema() -> None:
    schema = load_json(QUALITY_SCHEMA_PATH, "quality-gates schema")
    gates = ((schema.get("properties") or {}).get("gates") or {})
    item = gates.get("items") or {}
    props = item.get("properties") or {}
    if set((props.get("class") or {}).get("enum") or []) != EXPECTED_GATE_CLASSES:
        fail("quality-gates schema must define exactly G0-G8")
    if set((props.get("state") or {}).get("enum") or []) != EXPECTED_EVIDENCE_STATES:
        fail("quality-gates schema must define exact evidence vocabulary")
    required = set(item.get("required") or [])
    if not {"id", "class", "state"}.issubset(required):
        fail("quality-gates schema gate item must require id/class/state")


def validate_state(state: dict) -> None:
    if state.get("schema_version") != 1:
        fail("unsupported or missing STATE.json schema_version")

    if set(state.get("state_vocabulary", [])) != ALLOWED_STATES:
        fail("state_vocabulary does not exactly match the canonical state set")

    current_phase = state.get("current_phase")
    phase = state.get("phase") or {}
    if phase.get("id") != current_phase or phase.get("state") != "active":
        fail("current phase identity/state is inconsistent")

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
            if dependency not in packages_by_id:
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

    for package_id, required in EVIDENCE_REQUIREMENTS.items():
        pkg = packages_by_id.get(package_id)
        if pkg and pkg.get("state") == "done" and not required.issubset(set(pkg.get("evidence") or [])):
            fail(f"{package_id} done state is missing mandatory evidence")

    done_count = sum(pkg.get("state") == "done" for pkg in work_packages)
    if phase.get("done_work_packages") != done_count:
        fail("done_work_packages does not match actual done count")

    active_packages = [pkg for pkg in work_packages if pkg.get("state") == "active"]
    if len(active_packages) != 1:
        fail("foundation execution requires exactly one active work package")
    if active_packages[0].get("id") != state.get("current_work_package"):
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
    validate_markers()
    validate_classification_schema()
    validate_quality_schema()
    state = load_json(STATE_PATH, "STATE.json")
    validate_state(state)
    validate_status(state)
    print("Omnexa governance validation: PASS")
    print(f"Current phase: {state['current_phase']}")
    print(f"Current work package: {state['current_work_package']}")
    print(f"Progress: {state['phase']['done_work_packages']} / {state['phase']['mandatory_work_packages']} done")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
