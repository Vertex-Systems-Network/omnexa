#!/usr/bin/env python3
"""Focused dependency-free tests for multi-agent path enforcement helpers."""

from __future__ import annotations

import unittest

from agent_orchestration_common import allowed_path, normalize_path, path_matches, patterns_overlap


class AgentOrchestrationPathTests(unittest.TestCase):
    def test_exact_path_match(self) -> None:
        self.assertTrue(path_matches("kernel/internal/events/outbox.go", "kernel/internal/events/outbox.go"))
        self.assertFalse(path_matches("kernel/internal/events/outbox.go", "kernel/internal/events/bus.go"))

    def test_recursive_prefix_match(self) -> None:
        self.assertTrue(path_matches("modules/crm/**", "modules/crm/internal/service.go"))
        self.assertFalse(path_matches("modules/crm/**", "modules/finance/ledger.go"))

    def test_hidden_github_path_is_preserved(self) -> None:
        self.assertEqual(normalize_path("./.github/workflows/governance.yml"), ".github/workflows/governance.yml")
        self.assertTrue(path_matches(".github/workflows/**", ".github/workflows/governance.yml"))

    def test_exact_overlap(self) -> None:
        self.assertTrue(patterns_overlap("README.md", "README.md"))

    def test_prefix_overlap(self) -> None:
        self.assertTrue(patterns_overlap("kernel/internal/events/**", "kernel/internal/events/outbox.go"))
        self.assertFalse(patterns_overlap("modules/crm/**", "modules/finance/**"))

    def test_allowed_path_uses_shared_budget(self) -> None:
        self.assertTrue(
            allowed_path(
                "README.md",
                ["modules/crm/**"],
                ["README.md"],
            )
        )
        self.assertFalse(
            allowed_path(
                "modules/finance/ledger.go",
                ["modules/crm/**"],
                ["README.md"],
            )
        )


if __name__ == "__main__":
    unittest.main()
