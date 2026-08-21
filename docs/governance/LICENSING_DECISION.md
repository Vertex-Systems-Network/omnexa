# Omnexa Licensing & IP Decision Gate

Status: **Decision required before external distribution / public launch**

This document prevents the repository's existing license from becoming an accidental long-term product/business-model decision.

## 1. Current state

The repository currently contains the GNU General Public License version 3 (GPLv3). That is a strong copyleft license and may materially affect distribution obligations for covered derivative works.

No contributor or AI system may assume that GPLv3 is the intended final commercial licensing model merely because the file exists.

## 2. Required business decision

Before Omnexa is distributed externally, offered as downloadable/self-hosted software, opened to outside contributors, or publicly released, project ownership must explicitly choose and document one of the following strategies (or another legally reviewed strategy):

- GPL/open-source distribution;
- permissive open source;
- source-available commercial license;
- proprietary/commercial license;
- dual licensing;
- open-core/community + commercial editions.

This is a business/legal decision, not an AI architecture optimization.

## 3. Decision inputs

The decision must consider:

- SaaS-only vs downloadable/self-hosted distribution;
- enterprise/on-prem requirements;
- marketplace/plugin ecosystem goals;
- whether third parties may redistribute modified builds;
- contributor copyright ownership/assignment policy;
- compatibility with third-party dependencies;
- ability to sell commercial licenses/support;
- customer requirements around source availability;
- trademark/brand control;
- jurisdiction-specific legal advice.

## 4. Dependency license policy

Until P00/P01 governance establishes automated dependency controls:

- do not introduce dependencies with unclear or incompatible licenses;
- record license metadata for foundational dependencies;
- avoid copying code from external projects into Omnexa without explicit provenance/license compatibility;
- AI-generated code must not reproduce licensed source from unrelated repositories;
- marketplace/extensions must eventually declare their license independently from core platform licensing.

## 5. Name/trademark gate

`Omnexa` must not be treated as legally cleared for public commercial launch solely because the repository uses the name. Trademark/domain/company-name clearance is a separate pre-launch decision.

## 6. Governance rule

Changing the repository license requires explicit owner authorization and an ADR/legal decision record. An AI agent must never replace `LICENSE` autonomously.

## 7. Blocking scope

This decision does **not** block internal architecture/specification work. It becomes a hard gate before external distribution/public launch and should be resolved before licensing assumptions become embedded in contributor, marketplace or packaging architecture.

> This repository document is product governance, not legal advice. Formal commercial licensing should be reviewed by qualified counsel for the intended jurisdictions and distribution model.
