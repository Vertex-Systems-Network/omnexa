#!/usr/bin/env python3
"""Validate active write leases, migration reservations and forbidden-path safety."""

from __future__ import annotations

from agent_orchestration_common import all_write_records, fail, load_plan, patterns_overlap


def main() -> None:
    plan = load_plan()
    records = all_write_records(plan)

    for record in records:
        label = record.get("task_id") or record.get("agent_id") or record.get("branch") or "unknown-record"
        writes = list(record.get("write_paths") or [])
        forbidden = list(record.get("forbidden_paths") or [])
        if not writes:
            fail(f"{label} has no write_paths")
        for write in writes:
            for blocked in forbidden:
                if patterns_overlap(write, blocked):
                    fail(f"{label} write path {write} overlaps its forbidden path {blocked}")

    for index, left in enumerate(records):
        left_label = left.get("task_id") or left.get("agent_id") or left.get("branch")
        for right in records[index + 1 :]:
            right_label = right.get("task_id") or right.get("agent_id") or right.get("branch")
            for left_path in left.get("write_paths") or []:
                for right_path in right.get("write_paths") or []:
                    if patterns_overlap(left_path, right_path):
                        fail(
                            f"active write lease overlap: {left_label}:{left_path} <-> {right_label}:{right_path}"
                        )

    seen_versions: set[tuple[str, str]] = set()
    seen_paths: set[str] = set()
    for record in records:
        reservation = record.get("migration_reservation")
        if not isinstance(reservation, dict):
            continue
        owner = reservation.get("owner")
        version = str(reservation.get("version") or "")
        path = reservation.get("path")
        key = (str(owner or ""), version)
        if key in seen_versions:
            fail(f"duplicate migration owner/version reservation: {key[0]} / {key[1]}")
        if path in seen_paths:
            fail(f"duplicate migration path reservation: {path}")
        if path not in (record.get("write_paths") or []):
            fail(f"migration reservation path is outside task write_paths: {path}")
        seen_versions.add(key)
        seen_paths.add(path)

    print(f"PASS: {len(records)} active write records have non-overlapping leases/reservations")


if __name__ == "__main__":
    main()
