# Kernel migration ownership convention

P01.04 establishes the PostgreSQL migration substrate only. It does **not** add tenant, module-runtime, event, cache, storage, telemetry, or business schemas.

## Layout

Runtime SQL migrations added by an authorized future work package must be committed under the authoritative owner's directory:

```text
kernel/migrations/<owner>/<version>_<name>.sql
```

Rules:

- `<owner>` is the canonical lowercase dotted owner identity recorded in the migration ledger (for example `kernel.data` when that owner has an authorized schema migration);
- `<version>` is a positive, monotonically increasing integer for that owner;
- `<name>` is lowercase snake_case and immutable after merge;
- a migration file belongs to exactly one owner and may mutate only objects that owner is authorized to own;
- cross-owner writes, hidden/manual SQL, production data, and secrets are forbidden;
- applied migration SQL is immutable: changing/removing an applied version must be detected as ledger drift, not silently accepted;
- test-only synthetic migrations may be declared in Go tests and must not be treated as production schema.

## P01.04 foundation object

The P01.04 runner creates only `omnexa_kernel.schema_migrations`, which is migration-foundation metadata. No production owner migration files are shipped by P01.04 because later schemas are not yet authorized.
