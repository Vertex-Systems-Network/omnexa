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

- Task / run ID:
- Agent / role:
- Exact base SHA:
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

- [ ] Effective agent working instructions were re-evaluated at task start and before PR submission.
- [ ] `README.md` **Agent Working Instructions** was updated in this PR if the effective instructions changed.
- [ ] If no instruction changed, the PR states: `Agent instructions checked — README instruction delta: none`.
- [ ] No undeclared overlapping writer path, migration namespace, public contract or shared authoritative surface is being modified.
- [ ] Child/sub-agent authority stayed within this task's declared scope.

See `docs/governance/MULTI_AGENT_ORCHESTRATION.md`.

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
- [ ] ADR added/updated if architecture changed
- [ ] No unrelated files or scope added

## Risk / rollback

- Primary risk:
- Rollback or forward-fix approach:
