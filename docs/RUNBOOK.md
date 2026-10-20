<!--
INSTRUCTIONS FOR THIS FILE (user and AI: read before editing. Keep this header.)

WHAT THIS FILE IS: the operations book for this project, at docs/RUNBOOK.md.
How to build it, run it, deploy it, verify it, and where everything lives,
so anyone can operate the project without tribal knowledge.

EXACT STRUCTURE:
- "## Build": the real commands, copied from the working build.
- "## Run locally": how to start it and what you should see.
- "## Deploy": where it is hosted, how a deploy happens, how to roll back.
- "## Verify": the checks that prove it works after a change or deploy.
- "## Where things live": hosts, services, dashboards, and WHERE
  credentials are stored (names of vault entries or connectors only).
- "## Gotchas": the traps already learned the hard way, each with its fix.

RULES:
- NEVER a credential value in this file: no passwords, keys, or tokens.
  Pointers to secure storage only.
- Commands are copied from what actually works; if a command changes,
  this file changes in the same commit.
- No secrets, passwords, or tokens in this file, ever.
-->
# RUNBOOK: Patriot Pest Control

## Build
```bash
go build -o /tmp/patriot-landing ./cmd/server
go test ./...
```

## Run locally
```bash
go run ./cmd/server
```
The server serves the public site; check `/health` and `/ready` answer and the home page renders. Sign-in paths redirect to AlphaFlux by design; this repo does not authenticate anyone.

## Deploy
- Hosted on Dokploy (project PPC, app "Patriot") on David's home server.
- Push to `main` auto-deploys from the root `Dockerfile`. Domains: patriotpest.pro and www.patriotpest.pro, Let's Encrypt through Dokploy Traefik.
- Rollback: redeploy the prior commit from Dokploy, or revert the bad commit on main and let the auto-deploy rebuild. The recorded rollback-path write-up is TODO T-0.6.

## Verify
- Open https://patriotpest.pro and confirm the home page, `/contact`, `/signup`, `/search`, `/sitemap.xml`, `/robots.txt`, `/health`, and `/ready` respond.
- Confirm (509) 471-5767 is the only phone number shown and the company name reads "Patriot Pest Control" in full.
- Submit a labeled TEST through the contact form and confirm the AlphaFlux CRM record exists (delivery proven, not just accepted).

## Where things live
- Site: this repo, deployed by Dokploy (app "Patriot").
- Business platform (CRM, dashboards, FieldRoutes, Twilio): AlphaFlux, a separate project and repo.
- Data: SQLite catalog bind-mounted on the host at /home/server/go-patriot/data/patriot.db.
- Credentials: the AlphaFlux lead-capture key and any API credentials live only in the Dokploy deployment environment (names in the README environment table). Values are never in this repo.

## Gotchas
- Push to main is a production deploy; docs-only pushes still rebuild the live site.
- Twilio acceptance is not delivery: an SMS alert is only proven by a delivered status or recipient confirmation. The SMS alert leg is blocked on A2P 10DLC approval (TODO T-3.7).
- Do not put phone numbers, links, or email in Google Business Profile post text or artwork; GBP posting is paused under appeal in any case (see DECISIONS.md).
- Leave Oregon pages and the Arizona Google profile alone until David says otherwise.
