# Omnexa Licensing, IP & Trademark Owner Decision Brief

Status: **OWNER / LEGAL DECISION REQUIRED BEFORE EXTERNAL DISTRIBUTION**  
Tracker: **Issue #4**  
Effect on internal P01 engineering: **non-blocking**  
Effect on public/external distribution: **blocking**

This brief converts the broad licensing gate into explicit owner decisions. It does **not** change `LICENSE`, grant rights, clear a trademark or constitute legal advice.

## 1. Current repository fact

The repository currently contains GPLv3. Under Omnexa governance, that file is historical/current repository state only and must not be treated as proof that GPLv3 is the intended long-term commercial model.

Until an owner-authorized decision is accepted:

- do not replace `LICENSE`;
- do not advertise Omnexa as proprietary/open-source/open-core/dual-licensed;
- do not accept external contributions under an unstated contributor-IP model;
- do not externally distribute downloadable/self-hosted builds based on an assumed license strategy;
- do not represent `Omnexa` as trademark-cleared.

## 2. Owner decision A — distribution model

Choose the intended commercial delivery modes.

```text
[ ] SaaS / hosted service only
[ ] SaaS + customer self-hosted/on-prem
[ ] Self-hosted/on-prem primary
[ ] Community downloadable edition + commercial services
[ ] Other: ______________________________
```

Decision notes:

```text
____________________________________________________________
____________________________________________________________
```

Why this matters: distribution mode materially changes which source/license obligations, customer contract terms, update mechanisms and packaging boundaries need legal review.

## 3. Owner decision B — core licensing strategy

Select **one intended direction for legal review**, not an automatic repository change.

### Option B1 — proprietary/commercial core

Potential fit when control of redistribution, source access and commercial licensing is a primary requirement.

Questions to answer:

- Will customers receive binaries only, source escrow, or negotiated source access?
- Will on-prem customers receive perpetual or subscription rights?
- What modification/redistribution rights, if any, are allowed?

### Option B2 — open-source core

Potential fit when broad source availability/community redistribution is a deliberate product strategy.

Questions to answer:

- Strong copyleft, weak copyleft or permissive?
- How will commercial differentiation work?
- How will trademark and official-build control be maintained?

### Option B3 — open-core

Potential fit when a community core is public but enterprise modules/capabilities remain commercial.

Questions to answer:

- Which capabilities are permanently community vs commercial?
- Where is the technical/legal boundary?
- Can community and enterprise modules be distributed independently?

### Option B4 — dual-license

Potential fit when the same owner-controlled codebase is offered under an open/source-available license plus a separate commercial license.

Questions to answer:

- Does Omnexa control sufficient copyright to offer both licenses?
- What contributor agreement is required to preserve relicensing ability?
- Which customers/extensions require the commercial path?

### Option B5 — source-available commercial

Potential fit when source inspection/modification is desirable but unrestricted redistribution/competitive hosted use is not.

Questions to answer:

- What source access rights are granted?
- Is production use restricted by organization/revenue/hosted-service terms?
- Are the restrictions compatible with intended ecosystem positioning?

Owner-selected direction for counsel review:

```text
[ ] B1 Proprietary/commercial
[ ] B2 Open source
[ ] B3 Open-core
[ ] B4 Dual-license
[ ] B5 Source-available commercial
[ ] Other: ______________________________
```

## 4. Owner decision C — contributor IP model

Before accepting outside code contributions, choose a legally reviewed contributor model.

```text
[ ] No external code contributions initially
[ ] DCO-style contribution certification
[ ] Contributor License Agreement (CLA)
[ ] Copyright assignment
[ ] Other: ______________________________
```

Required policy questions:

- Who owns copyright in employee/contractor contributions?
- Does the model preserve future commercial/dual licensing if required?
- Are generated/code-assistance contributions covered by provenance policy?
- What representations are required from third-party contributors?

## 5. Owner decision D — marketplace / extension boundary

Omnexa core licensing must not silently determine every extension's license.

Choose initial policy direction:

```text
[ ] Official extensions must use the core license
[ ] Extensions may choose approved licenses independently
[ ] Marketplace supports proprietary + approved OSS extensions
[ ] Marketplace licensing deferred until P22, with no external marketplace before then
```

Regardless of choice, future module packages should declare:

- publisher identity;
- package license identifier/text reference;
- dependency licenses;
- redistribution rights;
- commercial entitlement requirements;
- trademark/branding restrictions where applicable.

## 6. Owner decision E — third-party dependency policy

Before P01 dependency growth, adopt a conservative interim rule:

- prefer standard library and dependencies with clear, reviewed license metadata;
- every foundational dependency must have provenance + license recorded;
- unknown/no-license dependencies are prohibited;
- copying source snippets from unrelated projects is prohibited without provenance/compatibility review;
- dependency license compatibility is checked before external distribution certification.

Final automated license-policy allow/deny rules belong to the P00.07/P01+ supply-chain tooling path and must follow the approved core strategy.

## 7. Owner decision F — Omnexa name/trademark

`Omnexa` is a product name in repository architecture, not evidence of trademark clearance.

Before public commercial launch:

```text
[ ] Perform professional trademark clearance in target jurisdictions/classes
[ ] Confirm company/product naming conflicts
[ ] Confirm primary domains/social handles as commercially appropriate
[ ] Define trademark usage policy for community/partners/extensions
[ ] Define official vs unofficial build/branding rules
```

A name change, if legally/business required, must be handled as governed product/architecture reconciliation rather than ad-hoc search/replace.

## 8. Suggested decision sequence

1. Owner selects distribution model.
2. Owner selects intended licensing strategy for counsel review.
3. Counsel reviews current GPLv3 repository history/ownership and intended transition path.
4. Contributor IP policy is chosen before external contributions.
5. Core/extension/marketplace license boundaries are recorded.
6. Dependency policy is configured against the approved model.
7. Trademark/name clearance is completed for target launch markets.
8. Owner authorizes exact repository license/trademark changes.
9. Accepted ADR/legal-decision record is merged.
10. Issue #4 closes only after resulting repository/public-launch actions are complete.

## 9. Formal approval record template

Do not fill this section by AI inference.

```text
Owner decision date: __________________
Approved distribution model: ______________________________
Approved core licensing direction: ________________________
Approved contributor IP model: ____________________________
Approved extension/marketplace policy: ____________________
Trademark/name decision status: ___________________________
Legal reviewer / reference: _______________________________
Authorized repository LICENSE change: YES / NO / DEFERRED
Authorized external distribution: YES / NO / CONDITIONS
Owner approver: ___________________________________________
```

## 10. Current gate result

```text
Internal architecture/specification: ALLOWED
Internal P01 engineering after P01 entry gate: ALLOWED
External downloadable/self-hosted distribution: BLOCKED
Public launch licensing claims: BLOCKED
External contribution program: BLOCKED until contributor-IP model approved
Trademark clearance: NOT ESTABLISHED
LICENSE replacement: NOT AUTHORIZED
```

Issue #4 remains open until an explicit owner/legal decision supplies closing evidence.
