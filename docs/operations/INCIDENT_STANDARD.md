# Omnexa Incident Severity and Response Standard

Status: **Canonical foundation v1**  
Work package: **P00.09**

## 1. Purpose

Define a common incident language for security, availability, integrity, privacy and operational failures.

## 2. Severity classes

### SEV0 — Crisis / existential integrity event
Examples:
- confirmed broad cross-tenant compromise;
- compromise of signing/root key material;
- uncontrolled unauthorized financial movement at platform scale;
- destructive corruption with no known safe recovery path;
- active compromise of the release/update trust chain.

Response expectations:
- immediate incident command;
- stop unsafe releases/automation;
- contain before restore;
- executive/security/legal/compliance escalation as applicable;
- preserve forensic evidence;
- customer/regulatory communication decisions owned by authorized humans.

### SEV1 — Critical
Examples:
- material tenant data exposure;
- privilege escalation with protected access;
- payment/ledger integrity breach;
- Tier 0/Tier 1 widespread outage;
- verified backup/restore failure affecting critical recovery posture.

Target acknowledgement: **<= 5 minutes** once production on-call exists.

### SEV2 — Major
Examples:
- partial Tier 1 outage;
- significant degraded mode;
- repeated data-processing failure without confirmed disclosure/corruption;
- large customer-impacting integration outage.

Target acknowledgement: **<= 15 minutes**.

### SEV3 — Moderate
Examples:
- localized non-critical degradation;
- Tier 2/Tier 3 incident with workaround;
- limited delayed processing;
- security weakness without evidence of exploitation but requiring prompt remediation.

Target acknowledgement: **<= 4 business hours**, with 24x7 handling when the risk requires it.

Lower-severity defects remain tracked through normal engineering workflow.

## 3. Severity rules

- severity is based on actual/potential impact, not team ownership;
- security/privacy/integrity may outrank availability impact;
- uncertainty should bias toward temporary higher severity until bounded;
- multiple-tenant impact raises severity;
- `RESTRICTED` data exposure is at least SEV1 unless evidence proves materially lower impact;
- zero-tolerance SLI breach is not downgraded because availability remains healthy.

## 4. Incident command roles

For SEV0-SEV2, establish as applicable:

- Incident Commander — coordinates decisions and timeline;
- Operations Lead — containment/recovery execution;
- Security Lead — threat/forensics/security decisions;
- Communications Lead — status/customer/internal communications;
- Scribe — timestamped decision/action log;
- Domain Owner — authoritative business/data semantics.

One person may hold multiple roles early in company maturity, but responsibilities remain explicit.

## 5. Response lifecycle

```text
DETECT -> TRIAGE -> DECLARE -> CONTAIN -> MITIGATE -> RECOVER -> VERIFY -> CLOSE -> REVIEW
```

Containment may intentionally reduce availability to protect tenant, security or financial integrity.

## 6. Required incident evidence

Capture:

- incident ID and severity;
- detection source/time;
- affected tenants/capabilities/data classes;
- actor/principal when known;
- correlation/trace/audit references;
- containment decisions;
- customer/provider dependencies;
- recovery verification;
- unresolved risk;
- timeline and decision owners.

Secrets and unnecessary sensitive data must not be copied into incident chat/tickets.

## 7. Communications

- internal status must distinguish confirmed facts from hypotheses;
- external communication must not invent root cause or recovery time;
- regulatory/customer notification is governed by applicable law/contract and authorized business/legal owners;
- status updates should give impact, mitigation state and next meaningful checkpoint without exposing exploitable details.

## 8. Post-incident review

SEV0/SEV1 always require a written review. SEV2 requires review when recurring, high-impact or architecture-relevant.

Review includes:

- impact;
- timeline;
- root/contributing causes;
- controls that worked/failed;
- detection gaps;
- recovery performance against RPO/RTO;
- action items with owner and due condition;
- whether threat model/SLO/runbook/test/architecture must change.

Postmortems are blameless about individuals but accountable about systems, decisions and follow-up.

## 9. Security incident special handling

Potential compromise requires evidence preservation, credential/key revocation/rotation strategy, scope validation and clean recovery. Do not destroy evidence merely to restore service faster.

## 10. Incident closure

An incident is closed only when immediate unsafe conditions are resolved, recovery is verified and remaining actions are explicitly tracked. Closure does not mean all long-term corrective work is finished.