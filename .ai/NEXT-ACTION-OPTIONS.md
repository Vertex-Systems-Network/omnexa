# VSN Organization Next-Action Options Contract

This repository adopts the Vertex Systems Network interactive AI-development handoff standard.

## User-facing handoff

During an already-authorized AI-Native development invocation, do **not** require a user choice after each milestone. The Supervisor selects the next highest-priority safe authorized task automatically and continues while authority, safety and the current tool/runtime budget remain available.

Expose 1 to 3 next-action options only when one of these applies:

- the user entered through the URL-only/read-only flow below;
- the user explicitly asks for choices/options;
- a genuine owner/external decision is required and no other safe authorized work remains.

When options are appropriate:

- include the canonical/recommended next action and mark it **Recommended**;
- when two or more valid options exist, numbering may be reshuffled;
- a reply containing only an option number selects that option, subject to fresh repository revalidation;
- interactive buttons may be used when the host supports them; otherwise numbered one-line options are acceptable.

Do not use options as a routine `continue?`, `retry?`, `fix?` or `merge?` confirmation gate for actions already inside standing development authority.

## URL-only repository entry

When the user's message contains only this repository's canonical GitHub URL (optionally with surrounding whitespace), treat it as a read-only development entry request.

1. Resolve the repository and default/protected branch.
2. Read this repository's durable/current state and governing instructions.
3. Reconcile open Issues first, then open PRs, then any repository-specific coordination/runner state required by local rules.
4. Do **not** create a branch, commit, PR, merge, deployment, provider call, destructive action, or other mutation from the URL alone.
5. Respond with 1 to 3 shuffled valid next-action options and mark the canonical one **Recommended**.
6. The user's subsequent number selection initiates the normal fully revalidated development turn.

## Safety and local authority

Repository-specific governance, security, exact-head CI, approval, migration, production/provider and release rules remain authoritative and may be stricter than this interaction contract. This file never grants execution authority and never permits bypassing an accepted actionable Issue/PR or deferred work boundary. Continuous invocation changes interaction cadence only; it does not expand scope or authority.
