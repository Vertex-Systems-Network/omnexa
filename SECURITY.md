# Omnexa Security Policy

Omnexa treats security, tenant isolation, authorization and auditability as platform requirements.

## Reporting security issues

Do not publish exploitable security details in public issues or pull requests. Use the repository/organization's private security-reporting channel when configured. Until a dedicated security contact/channel is formally published, repository maintainers must handle reports privately.

Never include production credentials, secrets, customer data, access tokens, private keys, payment data or other sensitive information in issues, logs, screenshots or commits.

## Security principles

- least privilege by default;
- deny by default for protected operations;
- explicit tenant/organization scope;
- server-side authorization enforcement;
- attributable audit trails for security-sensitive and business-significant mutations;
- secrets outside source control;
- encryption in transit and at rest where applicable;
- versioned external/public contracts;
- input validation at trust boundaries;
- idempotency/replay protection where required;
- secure-by-default module lifecycle;
- AI/automation uses governed capabilities and cannot bypass policy;
- dependencies and build provenance are treated as supply-chain security concerns.

## Prohibited practices

The following are security defects, not shortcuts:

- bypassing authorization for admin/super-user code paths;
- trusting tenant IDs supplied by a client without authorization context;
- direct cross-tenant queries without enforced scope;
- logging credentials/tokens/secrets;
- committing `.env` or equivalent secret files;
- disabling TLS/security checks for production behavior;
- raw AI access to transactional database write credentials;
- hidden default passwords;
- unaudited privilege escalation;
- using unverified external webhook payloads for protected actions;
- silently weakening a security control to make a test pass.

## Security changes

Changes to authentication, tenancy, authorization, cryptography, secret management, payment security boundaries, audit semantics, data classification or trust boundaries require architecture/security review and may require an ADR under `docs/governance/CHANGE_CONTROL.md`.

## Vulnerability handling lifecycle

1. Receive privately.
2. Triage severity and affected versions/components.
3. Reproduce without exposing sensitive data.
4. Define containment/fix/compatibility impact.
5. Add regression tests where possible.
6. Review tenant and authorization impact.
7. Release through the governed release process.
8. Publish disclosure/advisory only when safe and appropriate.

## Supported versions

No production release exists yet. Supported-version policy will be established before the first supported release under the release/security governance work packages.
