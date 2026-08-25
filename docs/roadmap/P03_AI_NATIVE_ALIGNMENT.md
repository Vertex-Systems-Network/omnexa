# P03 AI-Native Contract Alignment

Status: **PLANNING-ONLY / NO STRATEGIC PROGRAM ACTIVATION**  
Owner phase: **P03 — Module Runtime**

This document maps P03 contracts to proposed strategic X-program requirements so P03 does not create avoidable future incompatibility. `STRATEGIC_PROGRAMS.json` and `STRATEGIC_ACCEPTANCE_GATES.md` remain proposed under ADR-0011 and do not authorize implementation.

## Guardrails

- `STATE.json` remains execution authority; P03 is currently planned and inactive.
- `XQ-100`, `XSG-100`, `XTRUST-100`, `XPF-200` and `XPERF-100` remain `implementation_authorized=false`.
- P03 may define stable typed identities/metadata/hooks required by its own module-runtime scope.
- P03 must not implement System Graph storage/runtime, package trust enforcement/sandbox brokers, product federation runtime, performance-intelligence runtime, AI/model/agent runtime or generic future-framework machinery.
- Compatibility hooks must be minimal, typed, versioned, testable and bounded to canonical owners.

## Compatibility matrix

| P03 contract | Future strategic requirement | Design implication required now | Deferred implementation | Evidence / risk |
|---|---|---|---|---|
| Engineering/readiness workflow | `XQ-100` run/scope identity, exact-head evidence, test-oracle integrity, bounded retry, independent review | Keep P03 packages strict-sequential, exact-state bound, machine-validated and exact-head CI evidenced; do not weaken tests/governance to obtain green | XQ orchestration runtime and automated scope leases | Risk: AI work drifting beyond active package or self-certifying fake evidence |
| Manifest identity/version/owner | `XSG-100` typed module nodes; `XPF-200` product/module version identity; `XPERF-100` module attribution | Stable module ID, version, contract version and owner fields must be machine-readable and usable as correlation keys | System Graph, federation registry, performance analytics | Risk: future telemetry/graph/federation cannot join runtime identity reliably |
| Declared dependencies | `XSG-100` dependency edges; `XTRUST-100` dependency/advisory provenance | Preserve required/optional/platform/forbidden classes and deterministic version constraints; keep declaration separate from observed evidence | Graph drift analysis, advisory/license/trust policy | Risk: undeclared/private coupling becomes invisible technical debt |
| Lifecycle state/version | `XSG-100` lifecycle history; `XTRUST-100` revoke/kill/recovery; `XPF-200` disconnect/deprecate | Lifecycle transitions need stable state/reason/version/source identity and idempotent semantics | Trust kill switch, federated disconnect/revocation | Risk: unsafe disable/purge or untraceable lifecycle changes |
| Settings & feature flags | `XPF-200` attachment configuration; future policy/entitlement separation | Reuse `kernel.configuration`; settings/flags never grant authorization and lifecycle availability remains explicit | Federation entitlements/config promotion | Risk: feature flags become hidden privilege path |
| Capability registry | `XSG-100` capability nodes/edges; `XPF-200` versioned product capabilities; `XPERF-100` capability attribution | Stable capability ID/version/owner/consumer metadata with authorization requirements and availability state | Federation capability routing, graph/runtime attribution | Risk: anonymous capability calls cannot be governed or measured |
| Permission registration | `XSG-100` permission/security paths; `XTRUST-100` declared scopes; `XPF-200` entitlement-vs-authorization separation | Preserve stable permission IDs and `kernel.authorization` enforcement authority; registration is declaration, not grant | Trust scope consent, federated entitlement mapping | Risk: module install silently widens authority |
| UI contribution registry | `XPF-200` unified navigation/work surfaces | Contributions must declare slot, permission, availability and fallback; UI never authorizes backend action | Product federation/unified work runtime | Risk: hidden UI coupling or UI-only authorization assumptions |
| Migration ownership | `XSG-100` data topology/history; `XTRUST-100` package provenance | Bind migration to module/version/owner and reject cross-owner schema mutation unless explicitly approved | Graph lineage, trust package migration policy | Risk: lifecycle operation corrupts unrelated domain data |
| Health reporting | `XSG-100` observed/tested evidence; `XPF-200` health/SLO/version identity; `XPERF-100` attribution | Health must identify module/version/state/dependency status accurately and safely, with stable attribution identifiers | System Graph observations, federation SLO surfaces, performance budgets | Risk: operational state cannot be diagnosed or compared across versions |
| Package trust hooks | `XTRUST-100` publisher/signature/provenance/SBOM/scope profile | Reserve typed optional metadata/verification hook interfaces without declaring trust or executing untrusted package code | Full trust runtime, sandbox/network/secret/file brokers, quotas, revocation | Risk: premature pseudo-security creates false confidence |
| Reference modules / exit proof | All five programs need representative evidence and negative fixtures | P03 fixtures must prove dependency, lifecycle, migration, forbidden-coupling, health and isolation behaviors using stable source/version identity | Strategic-program-specific acceptance suites | Risk: architecture remains descriptive rather than executable |

## P03-specific CTQs

The following critical-to-quality properties are frozen for P03 planning:

1. **Isolation:** optional or failing modules cannot corrupt unrelated modules.
2. **Determinism:** discovery, dependency resolution and lifecycle transition validation are deterministic.
3. **Ownership:** no private cross-module imports/writes and no cross-owner migration authority.
4. **Security:** registration cannot grant authority; destructive lifecycle operations are explicit, authorized and auditable.
5. **Recoverability:** disable/re-enable and failed lifecycle operations preserve a recoverable state.
6. **Diagnosability:** lifecycle/dependency/health state is accurate and classification-safe.
7. **Compatibility:** stable IDs/versions/metadata can later feed graph, trust, federation and performance systems without those systems becoming P03 dependencies.

## Deferred scope by program

- `XQ-100`: automated AI run manifests, leases, evidence attestation and cost/circuit-breaker orchestration.
- `XSG-100`: graph storage/query, collectors, evidence-level inference, drift and blast-radius runtime.
- `XTRUST-100`: publisher verification, signature enforcement, SBOM/advisory/license policy, sandbox/brokers, quotas and kill/revoke runtime.
- `XPF-200`: native/embedded/federated/edge product attachment, SSO federation, entitlement and disconnect runtime.
- `XPERF-100`: performance budgets, load/capacity baselines, noisy-neighbor analysis and cross-layer timing aggregation.

## Acceptance for this planning artifact

This alignment is acceptable when:

- all 11 P03 packages reference only compatibility implications necessary to P03's own contract;
- no strategic program is marked active/done or implementation-authorized;
- `kernel.modules`, `kernel.configuration` and `kernel.authorization` ownership boundaries remain explicit;
- future graph/trust/federation/performance consumers can use stable module/capability/version/state identifiers without requiring P03 to implement those future runtimes;
- canonical readiness validators enforce the planning/activation boundary on GitHub-hosted `ubuntu-24.04`.
