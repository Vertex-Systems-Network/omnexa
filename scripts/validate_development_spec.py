#!/usr/bin/env python3
"""Validate the P00.08 repository/local-development specification."""

from __future__ import annotations

import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED = [
    "docs/development/REPOSITORY_STRUCTURE.md",
    "docs/development/LOCAL_DEVELOPMENT.md",
    "docs/development/TOOLCHAIN_STANDARD.md",
    "docs/development/CONFIGURATION_STANDARD.md",
    "docs/development/DEVELOPER_COMMANDS.md",
    "docs/contracts/development/workspace.schema.json",
    "docs/adr/ADR-0008-repository-local-development-baseline.md",
]

MARKERS = {
    "docs/development/REPOSITORY_STRUCTURE.md": ["apps/", "kernel/", "modules/", "generated/", "modular monolith", "direct DB writes"],
    "docs/development/LOCAL_DEVELOPMENT.md": ["PostgreSQL", "NATS + JetStream", "Kubernetes is not required", "WSL2", "synthetic"],
    "docs/development/TOOLCHAIN_STANDARD.md": ["repository-owned", "lockfile", "Unversioned", "Generated code"],
    "docs/development/CONFIGURATION_STANDARD.md": ["Secrets", "environment", "fail", "Feature flags"],
    "docs/development/DEVELOPER_COMMANDS.md": ["omnexa dev bootstrap", "omnexa verify all", "non-zero exit", "hidden manual steps"],
}

EXPECTED_ROOTS = {"apps", "kernel", "modules", "platform", "shared", "infrastructure", "scripts", "docs", "generated"}
EXPECTED_SERVICES = {"postgresql", "redis_compatible", "nats_jetstream", "s3_compatible", "mail_sink"}


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


def main() -> int:
    missing = [p for p in REQUIRED if not (ROOT / p).is_file()]
    if missing:
        fail("missing P00.08 evidence: " + ", ".join(missing))

    for path, markers in MARKERS.items():
        text = (ROOT / path).read_text(encoding="utf-8")
        for marker in markers:
            if marker.lower() not in text.lower():
                fail(f"{path} missing concept marker: {marker}")

    schema_path = ROOT / "docs/contracts/development/workspace.schema.json"
    try:
        schema = json.loads(schema_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        fail(f"workspace schema invalid: {exc}")

    props = schema.get("properties") or {}
    workspace = props.get("workspace") or {}
    root_enum = (((workspace.get("properties") or {}).get("root_categories") or {}).get("items") or {}).get("enum") or []
    if set(root_enum) != EXPECTED_ROOTS:
        fail("workspace schema root categories do not match canonical monorepo roots")

    local_services = ((props.get("local_services") or {}).get("items") or {}).get("enum") or []
    if set(local_services) != EXPECTED_SERVICES:
        fail("workspace schema local service vocabulary is inconsistent")

    if (props.get("kubernetes_required_for_default_local_loop") or {}).get("const") is not False:
        fail("default local loop must not require Kubernetes")
    if (props.get("production_data_allowed_in_default_local_environment") or {}).get("const") is not False:
        fail("production data must not be allowed by default locally")

    print("Omnexa P00.08 development specification validation: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
