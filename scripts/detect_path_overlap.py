#!/usr/bin/env python3
"""Fail when active multi-agent write budgets overlap."""

from __future__ import annotations

from agent_orchestration_common import all_write_records, fail, load_plan, patterns_overlap


def main() -> None:
    records = all_write_records(load_plan())
    overlaps: list[str] = []
    for index, left in enumerate(records):
        left_label = left.get("task_id") or left.get("agent_id") or left.get("branch")
        for right in records[index + 1 :]:
            right_label = right.get("task_id") or right.get("agent_id") or right.get("branch")
            for left_path in left.get("write_paths") or []:
                for right_path in right.get("write_paths") or []:
                    if patterns_overlap(left_path, right_path):
                        overlaps.append(f"{left_label}:{left_path} <-> {right_label}:{right_path}")
    if overlaps:
        fail("overlapping active write paths: " + "; ".join(overlaps))
    print("PASS: no active multi-agent write-path overlaps detected")


if __name__ == "__main__":
    main()
