#!/usr/bin/env python3
"""Validate active multi-agent dependency graph and deterministic merge order."""

from __future__ import annotations

from agent_orchestration_common import fail, load_plan, task_records


def main() -> None:
    plan = load_plan()
    tasks = task_records(plan)
    supervisor = plan.get("supervisor") or {}
    records = list(tasks)
    if isinstance(supervisor, dict) and supervisor.get("task_id"):
        records.append(supervisor)

    by_id = {}
    for record in records:
        task_id = record.get("task_id")
        if not isinstance(task_id, str) or not task_id:
            fail("every dependency record must have task_id")
        if task_id in by_id:
            fail(f"duplicate task_id in dependency graph: {task_id}")
        by_id[task_id] = record

    for task_id, record in by_id.items():
        order = record.get("merge_order")
        if not isinstance(order, int) or order < 1:
            fail(f"{task_id} must have positive integer merge_order")
        for dependency in record.get("depends_on") or []:
            if dependency not in by_id:
                fail(f"{task_id} depends on unknown task {dependency}")
            dependency_order = by_id[dependency].get("merge_order")
            if not isinstance(dependency_order, int) or dependency_order >= order:
                fail(f"{task_id} dependency {dependency} must have an earlier merge_order")

    visiting: set[str] = set()
    visited: set[str] = set()

    def visit(task_id: str) -> None:
        if task_id in visited:
            return
        if task_id in visiting:
            fail(f"dependency cycle detected at {task_id}")
        visiting.add(task_id)
        for dependency in by_id[task_id].get("depends_on") or []:
            visit(dependency)
        visiting.remove(task_id)
        visited.add(task_id)

    for task_id in by_id:
        visit(task_id)

    configured = (plan.get("merge_strategy") or {}).get("deterministic_order") or []
    actual = [task_id for task_id, _ in sorted(by_id.items(), key=lambda item: item[1].get("merge_order"))]
    if configured != actual:
        fail(f"deterministic_order does not match task merge_order: expected {actual}, got {configured}")

    print(f"PASS: dependency DAG valid for {len(by_id)} tasks with deterministic merge order")


if __name__ == "__main__":
    main()
