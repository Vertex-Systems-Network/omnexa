# Omnexa Pull Request

## Work package

- Phase / package: `Pxx.xx`
- State before PR:
- Intended state after verified merge:

## Scope

### In scope
- 

### Out of scope
- 

## AI / concurrent-work coordination

- Wave ID:
- Coordination channel / issue:
- Supervisor identity:
- Task / run ID:
- Agent / role:
- Branch:
- Exact branch-bootstrap base SHA:
- Last synchronized protected-main SHA:
- Current required protected-main SHA:
- Declared read paths:
- Declared write paths:
- Forbidden paths:
- Shared / exclusive paths or leases:
- Task dependencies / required merge order:
- Concurrent related tasks / agents:
- Migration reservation (owner/path/version/data budget), if any:
- Public contract/event reservation, if any:
- Conflict / overlap check result:
- Protected-main freshness checked:

### Supervisor-led submission protocol

- Completion signal posted: `Work Done and Submitted` — Yes / No / N/A
- Completion signal exact submitted head SHA:
- Submitted PR / reference:
- Submission CI state: PASS / FAIL / BLOCKED / NOT RUN / N/A
- Supervisor review state: pending / changes-required / approved-for-governed-merge / N/A
- If protected main moved after task start, branch synchronization completed: Yes / No / N/A
- Synchronized protected-main SHA after latest merge alert:
- Sync acknowledgement posted: `Sync Complete — Resuming Work` — Yes / No / N/A

- [ ] Effective agent working instructions were re-evaluated at task start and before PR submission.
- [ ] `README.md` **Agent Working Instructions** was updated in this PR if the effective instructions changed.
- [ ] If no instruction changed, the PR states: `Agent instructions checked — README instruction delta: none`.
- [ ] No undeclared overlapping writer path, migration namespace, public contract or shared authoritative surface is being modified.
- [ ] Child/sub-agent authority stayed within this task's declared scope.
- [ ] The branch is synchronized to the active wave's required protected-main SHA, or a documented reason proves the older base remains valid.
- [ ] A valid worker submission will be reviewed by the Supervisor before protected integration.
- [ ] Supervisor approval will use required exact-head Governance/promotion/protected-main rules; no direct unverified main merge is requested.

See `docs/governance/MULTI_AGENT_ORCHESTRATION.md` and `docs/governance/SUPERVISOR_MULTI_AGENT_WORKFLOW.md`.

## Architecture impact

- Owning module/domain:
- Kernel capability used/changed:
- Cross-module contracts affected:
- ADR required: Yes / No
- ADR reference (if yes):

## Data / migration impact

- Schema owner:
- New/changed migrations:
- Fresh-install impact:
- Upgrade impact:
- Destructive operation: Yes / No
- Rollback / forward-fix strategy:

## Security / tenancy

- Tenant-scoped data affected:
- Authorization scopes affected:
- Audit/security events affected:
- Sensitive data classification affected:

## Contracts and events

- APIs/capabilities changed:
- Events changed:
- Compatibility/versioning notes:

## Module lifecycle

- Install/upgrade impact:
- Disable/re-enable behavior:
- Optional dependency behavior:
- Purge behavior (if applicable):

## Acceptance evidence

Use only PASS / FAIL / BLOCKED / NOT RUN / N/A.

- Build:
- Static/type checks:
- Unit tests:
- Integration tests:
- Fresh migration/install:
- Upgrade migration:
- Tenant isolation:
- Authorization:
- Module lifecycle:
- Contract/event validation:
- Security scan:
- CI run / job references:

## Observability

- Logs/traces/metrics added or changed:
- Correlation/tenant context verified:

## Documentation / state reconciliation

- [ ] Relevant architecture/docs updated
- [ ] `STATUS.md` updated if progress changes
- [ ] `STATE.json` updated only if evidence supports transition
- [ ] README agent instructions re-checked and synchronized when materially changed
- [ ] Active multi-agent plan / required-main SHA updated when the wave changed materially
- [ ] ADR added/updated if architecture changed
- [ ] No unrelated files or scope added

## Risk / rollback

- Primary risk:
- Rollback or forward-fix approach:
