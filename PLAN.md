# PLAN: Patriot Pest Control

The whole plan for the Patriot Pest Control public website, step by step. Each step names the TODO.md phase it maps to. When the plan changes, this file changes first, then TODO.md follows. Facts only: what is decided is stated, what is open is named as an open question.

## What it is

This repository is the public front door for Patriot Pest Control at patriotpest.pro. It serves marketing pages, the pest library, service areas, blog and FAQ content, search, signup and contact forms, the PWA shell, sitemap and RSS, health endpoints, and first-party analytics intake.

It is not the business platform. Customer accounts, staff tools, admin dashboards, the CRM, FieldRoutes work, Twilio work, and private APIs belong to AlphaFlux. This site sends sign-in traffic to AlphaFlux and sends website leads into the AlphaFlux CRM.

Patriot Pest Control serves Eastern Washington and North Idaho. Arizona is out after the operation was sold. The only public website number is (509) 471-5767.

## Architecture

| Module | What it owns |
| :--- | :--- |
| `cmd/server` | Process start, route composition, legacy redirects, `/api/*` retirement, graceful shutdown, and the container health command |
| `internal/config` | Environment and feature configuration with defaults |
| `internal/modules/marketing` | Public pages, search, sitemap, RSS, contact and signup handling, lead capture, and the first-party beacon sink |
| `internal/modules/health` | `/health`, `/ready`, and `/metrics` |
| `internal/view` | Page modules, templates, embedded assets, PWA files, SEO helpers, and pest and area copy |
| `internal/data` | SQLite content catalog, analytics storage and reporting inputs, drafts, and local signup fallback |
| `internal/middleware` | Shared request behavior |
| `pkg/crypto` | Small cryptographic helpers used by the service |
| AlphaFlux boundary | Accounts, dashboards, CRM records, source attribution, alert delivery, FieldRoutes, and Twilio; reached by redirect or the lead capture API, never duplicated here |

## Build order

1. Hold the foundation and repo standard: one Go service, a Docker image, a clear AlphaFlux boundary, complete configuration notes, and a rollback path. (TODO Phase 0)
2. Keep the public site fast, live, and easy to deploy from main, with smoke checks after changes. (TODO Phase 1)
3. Grow the content grid one page module at a time: pests, services, areas, blog, FAQs, prices, and referral content, all checked for full company name, market, phone, and veteran-owned positioning. (TODO Phase 2)
4. Make every form a reliable lead path: CRM first, local fallback if capture fails, then delivery-proven alerts. (TODO Phase 3)
5. Measure what the public site does with first-party analytics, search reporting, and conversion checks that feed AlphaFlux without moving private tools into this repo. (TODO Phase 4)
6. Raise the release gate until format, tests, lint, vulnerability scanning, secret scanning, accessibility, links, and restore drills all pass before a risky change ships. (TODO Phase 5)

## Operating rules for this site

- Push to main auto-deploys through Dokploy, so docs and code changes on main are production changes.
- A lead is not done when the form accepts it. It is done when the CRM record exists and the alert is proven delivered.
- Twilio acceptance is not delivery. SMS proof requires a delivered status or recipient confirmation.
- No secret belongs in this repo. The AlphaFlux website lead key lives only in the deployment environment.

## Open questions

- Which single pricing lineup replaces the differing website and phone menu prices? (owner: Skyler and David)
- What email address should receive lead alerts? (owner: David or Skyler)
- How should call tracking be added without replacing the only public number? (owner: David and Skyler)
- Should the “Book online” wording point to a booking flow, and which one? (owner: Skyler and David)
