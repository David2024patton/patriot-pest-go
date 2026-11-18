<!--
About this file (a description of what it holds and the conventions it follows).

This file is the build rules for this repo in one page, at docs/RULES.md:
how every module in this project gets built, gated, and proven, recorded
from David Patton's standing build orders. Project-specific numbers and
additions live in a "Project specifics" section at the end; the sections
above it come from the github-repo-standard skill template.

Conventions:
- Where a project practice conflicts with this page, the page records
  the current rule until David changes it.
- The file holds no secrets, passwords, or tokens.
-->
# RULES: how this project gets built

The build law for this repo, in one page. These rules come from David
Patton and apply to every module, every language, every time.

## Modularity: the grid and the tiles
- The main program or page is a GRID. Every section is a MODULE plugged
  into a grid slot: header, footer, each content block, each tile.
  Modular to the full max, never partial.
- Every feature is its own isolated module with a defined interface.
  One feature, one module, one commit, one push.
- Modules are independently replaceable and updatable on the fly. A
  broken module gets fixed or swapped alone; the system never rebuilds
  around it.

## Build 20 years ahead
- Design for the bottleneck that is COMING, not just today's.
- Interface seams at every change-prone boundary: model providers,
  executors, verifiers, storage, transports.
- Operational limits (timeouts, caps, retry counts, quotas) are config
  with sane defaults, never constants compiled into logic.
- Data formats carry version fields so future versions read today's data.

## Debugging and linting, in every module, from birth
- They are not the same thing. Linting reads the code for mistakes
  without running it. Debugging is finding out why the program
  misbehaved when it ran, backed by diagnostics built into the module.
- Every module ships with BOTH: structured diagnostics, a way to pull
  the bug information out, and lint gates that run clean. A module that
  cannot say what it did and why it refused is not done.

## Check the shelf before building
- Before building anything, check the take-shelf
  (~/workspace/google-org-harvest/USE-IT.md) and the wider harvest: if
  Google or another trusted source already built it under a usable
  license, take it or retrofit it. Rebuilding what exists is waste.
- For any new market or feature area, research the top 10 competitors
  (the customer's list and the platform's list) before designing.

## Security
- No secrets in the repo, ever: keys, tokens, and passwords live in
  secure storage only. Secret scanning (gitleaks or betterleaks) runs
  in the gates.
- Every refusal, overage, and verdict emits through the project's
  diagnostic channel. Nothing fails silently.
- Untrusted input gets a bounded parser and a fuzz target, not trust.

## Verification
- Done means proven: a test with asserted values, a log line, a bench
  note, or a live check. "It compiles" is not proof.
- Never blindly trust delegated work; re-run the gates before merging.
- Emulation-proven is not bench-proven; label which one a claim is.

## Writing
- Plain, direct, technical. No em dashes in prose, docs, or code.
- Brand casing is exact: AlphaVision, AlphaFlux, iTaK, Mimir,
  Patriot Pest Control (full name wherever there are words).
- Project docs: TODO.md tracks the work, PLAN.md holds the plan,
  IDEA.md holds David's ideas in his words.

## The relational diagram (David's rule, 2026-11-18)
- This project always has a relational diagram: docs/MAP.md, nodes and
  wires for every module, page, and external connection, drawn from
  the real code in this repo, in readable markdown.
- After every feature build, and after any change to any part of the
  project, the diagram is updated in the same commit and the same
  push/deploy cycle.
- A feature is not fully landed until the diagram matches it. Where
  the code and the diagram disagree, the code wins and the diagram is
  fixed in that same pass.
