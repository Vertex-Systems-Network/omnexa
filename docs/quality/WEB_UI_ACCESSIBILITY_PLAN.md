# Web UI Standards & Accessibility Execution Plan

Status: **Canonical planning requirement — implementation activates only with an authorized web-UI work package**

This plan tells human contributors and AI coding systems how browser-delivered Omnexa interfaces must be designed, implemented and verified once UI implementation is authorized. It does **not** authorize P12/P13/P17 or any other future UI/business package during P01.

## 1. Conformance target

All production browser UI must target:

- **WCAG 2.2 Level AA** as the default accessibility conformance target;
- standards-based semantic HTML rendered to the browser;
- valid CSS where applicable;
- native HTML semantics before ARIA;
- WAI-ARIA patterns only when a native semantic element cannot express the required interaction;
- keyboard, focus, screen-reader, zoom/reflow and contrast behavior appropriate to the affected interface.

The owning package may adopt a stricter target, but may not silently weaken WCAG 2.2 AA.

## 2. W3C validation implementation

### HTML

Rendered critical pages and representative component states must be checked with the W3C HTML/Nu validation tooling or an equivalent W3C-conformant validation path.

Rules:

- validate the **rendered DOM/output**, not raw React/TSX source as if it were HTML;
- no unresolved HTML conformance errors on required release surfaces;
- warnings must be reviewed and either fixed or documented with a standards-based rationale;
- duplicate IDs, invalid nesting, invalid attributes, broken landmark/heading structure and malformed form semantics are defects.

### CSS

CSS produced/owned by Omnexa should be validated with the W3C CSS Validation Service where practical.

Rules:

- standards errors owned by Omnexa are fixed rather than suppressed;
- framework/vendor-extension warnings may be documented when they are intentional and do not break interoperability/accessibility;
- generated third-party CSS is not rewritten blindly; ownership and upstream remediation are recorded.

## 3. WAVE implementation

WAVE is a required evaluation input for critical web surfaces, but it is **not** treated as a certification badge and does not replace human accessibility testing.

### Developer workflow

For local, authenticated, private or highly dynamic pages, contributors should use the WAVE browser extension against the rendered interface.

### CI / repeatable audit workflow

When the owning UI package reaches implementation, the repository must integrate one approved WAVE automation mode:

1. WAVE Subscription API for publicly reachable preview/test URLs; or
2. WAVE Stand-alone API / Testing Engine for private, authenticated or CI-hosted rendered pages.

Implementation rules:

- WAVE API/license credentials are secrets and must never be committed, logged or copied into artifacts;
- a required WAVE lane without an available approved key/license is `BLOCKED`, not `PASS`;
- CI must record the tested route/state, viewport, tool/engine version where available and summarized findings;
- raw reports containing sensitive rendered content must not be published as public artifacts;
- rate limits, credits/licensing and external-service availability must be handled explicitly rather than bypassing the gate.

### WAVE result policy

For critical release surfaces:

- unresolved WAVE **Errors** block completion;
- unresolved WAVE **Contrast Errors** block completion;
- WAVE **Alerts** require human review and recorded disposition; an alert is not automatically a defect and may not be silently ignored;
- ARIA findings require semantic review, not mechanical attribute addition;
- a clean automated WAVE scan alone is insufficient to claim WCAG conformance.

## 4. Mandatory manual accessibility checks

Automated tools cannot prove full accessibility. Applicable UI changes must include human verification of affected flows.

Minimum checks for interactive browser UI:

- complete keyboard operation without a mouse;
- visible, logical focus order;
- no keyboard trap;
- focus is not obscured by sticky headers, dialogs or overlays;
- dialogs, menus, tabs, comboboxes and other composite widgets follow appropriate interaction semantics;
- accessible names for controls, icons and links;
- labels, instructions and programmatic error association for forms;
- status/error/success messages are exposed appropriately to assistive technology;
- heading hierarchy and landmarks support navigation;
- skip/navigation mechanisms where repeated content warrants them;
- color is not the sole carrier of meaning;
- WCAG AA text/non-text contrast requirements are met;
- content remains usable under browser zoom/text scaling and narrow reflow conditions;
- pointer target size and dragging interactions have accessible alternatives where WCAG 2.2 requires them;
- at least one supported screen-reader smoke test for critical navigation/forms/dialog flows before a UI package can be considered release-ready.

## 5. Responsive, locale and RTL rule

Accessibility testing must cover the layouts actually shipped:

- desktop and mobile breakpoints;
- long translated strings;
- RTL presentation when the owning feature supports RTL/localization;
- validation and error states;
- loading, empty, disabled and destructive-confirmation states;
- high-content-density enterprise screens, not only marketing/landing pages.

Responsive visual correctness does not excuse semantic or keyboard regressions.

## 6. AI execution rules for UI work

When an AI coding system implements or reviews Omnexa UI, it must:

1. prefer native semantic HTML over `div`/ARIA reconstruction;
2. preserve keyboard and focus behavior when changing layout or animation;
3. never add ARIA merely to silence an automated checker;
4. never claim “WAVE passed”, “WCAG compliant” or “W3C compliant” solely because an automated scan has no errors;
5. fix source defects rather than disabling validator rules to obtain green CI;
6. preserve accessible names and labels when replacing icons/components;
7. treat contrast, focus visibility, hit-target size and zoom/reflow as implementation requirements, not final polish;
8. include validation evidence in the owning PR for affected critical routes/states;
9. record any tool limitation or blocked external WAVE dependency honestly;
10. keep accessibility requirements intact when generating templates, pages, portal surfaces or low-code output.

## 7. Quality-gate mapping

This plan uses the existing G0-G8 vocabulary; it does not invent a new gate class.

Typical mapping when web UI becomes active:

- **G0 Governance** — owning package explicitly names WCAG/W3C/WAVE evidence applicability;
- **G1 Static** — lint/type checks plus standards validation where deterministic/static output allows it;
- **G2 Unit/Component** — component semantics, keyboard/focus and state behavior;
- **G3 Contract/Integration** — rendered route/component integration, assistive-technology semantics and automated accessibility scans;
- **G5 Security/Tenancy** — only when accessibility changes affect authorization/tenant-safe UI behavior; accessibility itself is not reclassified as a security gate;
- **G7 Build/Package** — production build preserves required validation and accessibility behavior.

The owning package specification determines which classes are required.

## 8. Evidence format

A UI completion record should include, as applicable:

- source SHA and PR;
- critical routes/components tested;
- viewport(s) and locale/RTL state;
- W3C HTML validation result;
- W3C CSS validation result where applicable;
- WAVE tool mode (extension/API/stand-alone engine) and summarized Error/Contrast Error/Alert disposition;
- keyboard/focus result;
- screen-reader smoke-test result;
- zoom/reflow result;
- unresolved exceptions with issue, owner and expiry/review condition;
- CI run/job identifiers.

Required failures may not be relabeled PASS.

## 9. Roadmap applicability

This plan becomes especially relevant to:

- P12 — Experience Builder & CMS;
- P13 — Portal Platform;
- P17 — Low-code App Builder;
- P21 — Developer Platform examples/documentation UI;
- any earlier/later administrative, authentication, workflow, reporting, marketplace or enterprise web surface once its package authorizes browser UI.

Generated pages/components are held to the same standard as hand-authored first-party UI.

## 10. Current authorization boundary

At the time this planning requirement is introduced, P01 is the active phase and business-feature/UI implementation is locked. Therefore:

- this file is a future execution requirement only;
- no P12/P13/P17 frontend is authorized by this planning change;
- no WAVE secret/license is required during P01 kernel packages unless a future authorized package explicitly owns browser UI;
- the active P01 package remains the sole executable kernel scope.

## References

- W3C Web Content Accessibility Guidelines (WCAG) 2.2: `https://www.w3.org/TR/WCAG22/`
- W3C CSS Validation Service: `https://jigsaw.w3.org/css-validator/`
- WAVE evaluation tools: `https://wave.webaim.org/`
- WAVE API: `https://wave.webaim.org/api/`
- WAVE browser extensions: `https://wave.webaim.org/extension/`

These external references inform implementation; the repository's governed package/state rules remain the execution authority.