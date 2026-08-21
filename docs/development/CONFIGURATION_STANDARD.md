# Omnexa Configuration Standard

Status: **Canonical v1**  
Work package: **P00.08**

Configuration is explicit, typed/validated where practical, environment-scoped and separate from secrets. Hidden machine state is not a supported configuration source.

## Configuration precedence

Future runtime configuration follows an explicit precedence model such as:

```text
compiled/default safe values
 -> versioned non-secret config
 -> environment/deployment config
 -> secret references/resolution
 -> narrowly scoped runtime overrides/feature flags
```

The exact loader may evolve, but precedence must be deterministic and documented.

## Rules

- Configuration keys use stable lowercase `snake_case` or environment-variable equivalents according to the owning runtime convention.
- Every required key has an owner, description, type, default/required semantics and sensitivity classification.
- Unknown configuration should fail or warn according to strictness policy; silently accepting misspelled security-critical keys is prohibited.
- Secrets are references resolved from approved local/production secret stores, not committed values.
- Tenant business settings are not global process environment variables; they belong in governed tenant configuration capabilities/data.
- Feature flags are governed configuration, not hidden code forks.

## Environment classes

Canonical environment intents:

```text
local
ci
preview/test
staging
production
```

Names may map to deployment-specific identifiers, but production credentials/data never become defaults for lower environments.

## Local configuration

Repository may provide examples such as `.env.example` or structured example config containing only safe placeholders.

Local configuration files containing secrets are ignored by version control. Bootstrap validates required fields and reports missing/invalid keys.

## Secrets

Secrets include passwords, API keys, private/signing/encryption keys, tokens, webhook credentials and authentication-equivalent values.

- never commit them;
- never expose them in generated configuration docs/output;
- never log them;
- production secrets use approved secret/KMS systems;
- local development uses separate non-production values;
- rotation/revocation behavior is part of the owning capability.

## URLs and external providers

External endpoints/providers are explicit configuration with environment-specific allowlists/policy where security requires. A developer-provided URL does not bypass SSRF/egress controls.

## Database/cache/broker/storage configuration

Connection configuration separates:

- endpoint;
- logical database/bucket/stream namespace;
- authentication reference;
- TLS/security options;
- tenant-independent runtime identity.

No domain module invents a second arbitrary connection mechanism when a platform-owned infrastructure client exists.

## Configuration validation

Application/service startup validates configuration before accepting traffic or work. Security-sensitive invalid configuration fails closed.

Future generated config reference documentation should derive from canonical schemas/registries where practical.

## No config-as-code ownership bypass

A config file cannot redefine domain ownership, create undocumented permissions or silently activate a future phase/module. Module enablement/entitlements/lifecycle remain governed platform capabilities.
