<!--
INSTRUCTIONS FOR THIS FILE (user and AI: read before editing. Keep this header.)

WHAT THIS FILE IS: the living right-now snapshot of this project. Anyone
opening the repo cold reads this first and is current in two minutes.
Updated at the end of every work session on the project; a stale STATUS
is worse than none, so when you finish work here, update it in the same
pass.

EXACT STRUCTURE:
- "## As of <YYYY-MM-DD>": the date this snapshot was written.
- "## Proven state": what is done, with its proof named.
- "## In flight": work currently underway and where it stands.
- "## Blocked": each blocker and exactly who or what it waits on.
- "## Next moves": the next three moves, in order, referencing TODO ids.

RULES:
- Claims here carry the same proof standard as TODO checkboxes.
- Keep it short: a snapshot, not a history (history is CHANGELOG.md).
- No secrets, passwords, or tokens in this file, ever.
-->
# STATUS: Patriot Pest Control

## As of 2026-10-20

## Proven state
- Public site live at patriotpest.pro from this repo (proof: live site verified 2026-10-14; Dokploy deploys from main).
- Arizona removed from the website and (509) 471-5767 is the only public number (proof: live-site verification 2026-10-14, TODO T-2.2 and T-2.3).
- Website forms create AlphaFlux CRM leads with source attribution, phone required and server-validated, local fallback if capture fails (proof: live TEST submission and CRM record verification 2026-10-16, TODO T-3.1 to T-3.5).
- Repo standard set complete: TODO, IDEA, PLAN, ENDGOAL, docs/RULES.md, plus this jump-back-in set (proof: files on main).

## In flight
- Repo standard jump-back-in backfill (this file set) landing 2026-10-20.
- Lead alert completion: SMS leg waits on A2P approval, email leg waits on the alert address (TODO Phase 3).

## Blocked
- SMS lead alerts (waiting on: approved A2P 10DLC brand and campaign, TODO T-3.7).
- Lead alert email address (waiting on: David or Skyler confirming the address; never guess it, TODO T-3.9).
- Single pricing lineup (waiting on: Skyler and David, TODO T-P.1).

## Next moves
1. Complete the environment contract in one operations document (TODO T-0.4).
2. Add the repeatable local startup check proving ready state (TODO T-0.5).
3. Record the Dokploy rollback path for a bad main push (TODO T-0.6).
