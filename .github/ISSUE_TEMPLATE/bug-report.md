---
name: Bug report
description: Report a reproducible defect without redefining architecture
title: "[BUG] "
labels: []
assignees: []
---

## Phase / work package

Identify the current or affected phase/work package if known.

## Summary

What is broken?

## Reproduction

1. ...
2. ...
3. ...

## Expected behavior

What should happen according to the current contract/specification?

## Actual behavior

What happened instead?

## Scope / environment

- deployment/environment:
- tenant/org context (no sensitive identifiers):
- module/version/commit:

## Evidence

Provide safe logs, test output or screenshots with secrets/customer data removed.

## Architecture/security check

- Does fixing this require a new dependency or architecture change? Yes / No / Unknown
- Tenant isolation affected? Yes / No / Unknown
- Authorization affected? Yes / No / Unknown
- Data integrity/migration affected? Yes / No / Unknown

If the fix requires changing an established architecture rule, use the Architecture Change template/ADR process instead of hiding the change inside a bug fix.
