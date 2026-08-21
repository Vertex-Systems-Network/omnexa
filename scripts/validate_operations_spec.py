#!/usr/bin/env python3
"""Dependency-free structural validation for P00.09 operational specifications."""

from __future__ import annotations

import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]

REQUIRED = {
    "docs/security/THREAT_MODEL.md": ["T01", "T24", "Cross-tenant", "AI prompt injection", "Risk treatment"],
    "docs/operations/SLO_STANDARD.md": ["Tier 0", "99.99%", "RPO", "RTO", "Error budgets", "zero-tolerance"],
    "docs/operations/INCIDENT_STANDARD.md": ["SEV0", "SEV1", "SEV2", "SEV3", "DETECT -> TRIAGE -> DECLARE"],
    "docs/operations/RELIABILITY_STANDARD.md": ["Observability minimum", "Golden signals", "Backpressure", "Graceful degradation", "Operational readiness gate"],
    "docs/adr/ADR-0009-threat-model-slo-reliability-baseline.md": ["Status: **accepted**", "zero-tolerance", "Recovery classes"],
}

SCHEMA_PATH = ROOT / "docs/contracts/operations/operational-targets.schema.json"
EXPECTED_TIERS = {"TIER_0", "TIER_1", "TIER_2", "TIER_3"}
EXPECTED_RECOVERY = {"A", "B", "C", "D"}


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


for path, markers in REQUIRED.items():
    target = ROOT / path
    if not target.is_file():
        fail(f"missing required P00.09 artifact: {path}")
    text = target.read_text(encoding="utf-8")
    for marker in markers:
        if marker not in text:
            fail(f"{path} missing marker: {marker}")

try:
    schema = json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))
except Exception as exc:
    fail(f"operational targets schema invalid/unreadable: {exc}")

props = schema.get("properties") or {}
if set((props.get("criticality_tier") or {}).get("enum") or []) != EXPECTED_TIERS:
    fail("operational schema must define exactly TIER_0..TIER_3")
if set((props.get("recovery_class") or {}).get("enum") or []) != EXPECTED_RECOVERY:
    fail("operational schema must define exactly recovery classes A-D")
if (props.get("incident_severity_model") or {}).get("enum") != ["SEV0-SEV3"]:
    fail("operational schema must declare SEV0-SEV3 incident model")

required = set(schema.get("required") or [])
for field in {"capability", "owner", "criticality_tier", "recovery_class", "slos", "incident_severity_model"}:
    if field not in required:
        fail(f"operational schema missing required field: {field}")

print("Omnexa P00.09 operations specification validation: PASS")
