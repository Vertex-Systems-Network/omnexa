#!/usr/bin/env python3
"""Validate active multi-agent task/slot identity and canonical authority alignment."""

from __future__ import annotations

from agent_orchestration_common import fail, load_plan, load_state, slot_records, task_records


def main() -> None:
    plan = load_plan()
    state = load_state()
    authority = plan.get("authority") or {}
    if authority.get("phase") != state.get("current_phase"):
        fail("active plan phase does not match canonical STATE.json")
    if authority.get("work_package") != state.get("current_work_package"):
        fail("active plan work package does not match canonical STATE.json")
    lock = state.get("implementation_lock") or {}
    if authority.get("business_feature_code_authorized") != lock.get("business_feature_code_authorized"):
        fail("active plan business-feature authority differs from canonical STATE.json")
    if authority.get("kernel_code_authorized") != lock.get("kernel_code_authorized"):
        fail("active plan kernel authority differs from canonical STATE.json")

    slots = slot_records(plan)
    tasks = task_records(plan)
    onboarding = plan.get("onboarding") or {}
    if onboarding.get("current_worker_slot_capacity") != len(slots):
        fail("current_worker_slot_capacity does not match worker slot count")
    open_count = sum(slot.get("status") == "open" for slot in slots)
    if onboarding.get("current_open_worker_slots") != open_count:
        fail("current_open_worker_slots does not match actual open slot count")
    if len(slots) > 3:
        fail("current M2 rollout permits at most 3 worker write slots until real-worker evidence proves scaling")

    slot_ids: set[str] = set()
    slot_branches: set[str] = set()
    slots_by_id = {}
    for slot in slots:
        slot_id = slot.get("slot_id")
        branch = slot.get("branch")
        if not isinstance(slot_id, str) or not slot_id:
            fail("worker slot missing slot_id")
        if not isinstance(branch, str) or not branch:
            fail(f"worker slot {slot_id} missing branch")
        if slot_id in slot_ids:
            fail(f"duplicate worker slot id: {slot_id}")
        if branch in slot_branches:
            fail(f"duplicate worker slot branch: {branch}")
        slot_ids.add(slot_id)
        slot_branches.add(branch)
        slots_by_id[slot_id] = slot
        if slot.get("status") == "occupied" and not slot.get("agent_id"):
            fail(f"occupied slot {slot_id} has no agent_id")
        if slot.get("status") == "open" and slot.get("agent_id") is not None:
            fail(f"open slot {slot_id} must not have an agent_id")

    task_ids: set[str] = set()
    task_branches: set[str] = set()
    for task in tasks:
        task_id = task.get("task_id")
        branch = task.get("branch")
        slot_id = task.get("slot_id")
        if not isinstance(task_id, str) or not task_id:
            fail("active task missing task_id")
        if task_id in task_ids:
            fail(f"duplicate active task id: {task_id}")
        if not isinstance(branch, str) or not branch:
            fail(f"active task {task_id} missing branch")
        if branch in task_branches:
            fail(f"duplicate active task branch: {branch}")
        task_ids.add(task_id)
        task_branches.add(branch)
        if slot_id not in slots_by_id:
            fail(f"task {task_id} references unknown worker slot {slot_id}")
        slot = slots_by_id[slot_id]
        for field in ("branch", "agent_id", "module"):
            if task.get(field) != slot.get(field):
                fail(f"task {task_id} {field} does not match slot {slot_id}")
        if list(task.get("write_paths") or []) != list(slot.get("write_paths") or []):
            fail(f"task {task_id} write_paths do not match slot {slot_id}")

    supervisor = plan.get("supervisor") or {}
    if supervisor.get("branch") in task_branches:
        fail("Supervisor branch must be isolated from worker branches")
    if supervisor.get("task_id") in task_ids:
        fail("Supervisor task_id must be unique")

    print(f"PASS: multi-agent task/slot authority valid ({len(tasks)} workers, {open_count} open slots)")


if __name__ == "__main__":
    main()
