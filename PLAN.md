# PLAN: Patriot Pest Control

The whole plan for the Patriot Pest Control public website, step by step. Sections are numbered S1, S2, ... and each build step names the TODO.md phase it maps to; every TODO item points back with (PLAN S<n>). When the plan changes, this file changes first, then TODO.md follows. Facts only: what is decided is stated, what is open is named as an open question.

## S1: What it is

This repository is the public front door for Patriot Pest Control at patriotpest.pro. It serves marketing pages, the pest library, service areas, blog and FAQ content, search, signup and contact forms, the PWA shell, sitemap and RSS, health endpoints, and first-party analytics intake.

It is not the business platform. Customer accounts, staff tools, admin dashboards, the CRM, FieldRoutes work, Twilio work, and private APIs belong to AlphaFlux. This site sends sign-in traffic to AlphaFlux and sends website leads into the AlphaFlux CRM.

Patriot Pest Control serves Eastern Washington and North Idaho. Arizona is out after the operation was sold. The only public website number is (509) 471-5767.

## S2: Architecture

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

## S3: Hold the foundation (TODO Phase 0)

One Go service, a Docker image, a clear AlphaFlux boundary, complete configuration notes, and a rollback path. The repo standard files (TODO, IDEA, PLAN, ENDGOAL, docs/RULES.md) stay current so the project can be picked up cold.

## S4: Keep the public site fast, live, and easy to deploy (TODO Phase 1)

Deploy from main through Dokploy, with smoke checks after changes and verified rendering at desktop, tablet, and phone sizes. Uptime and error alerting runs on cloudprober (from the take-shelf, google/cloudprober, Apache-2.0), run as-is, so a dead page alerts the operations channel instead of waiting for a person to notice.

## S5: Grow the content grid one page module at a time (TODO Phase 2)

Pests, services, areas, blog, FAQs, prices, and referral content, all checked for the full company name Patriot Pest Control, the Eastern Washington and North Idaho market, the single public phone number, and veteran-owned positioning. Service locations get Plus Codes from open-location-code (google/open-location-code, Apache-2.0) where a street address describes a service location badly, so a code pastes straight into Google Maps, tickets, and dashboards.

## S6: Make every form a reliable lead path (TODO Phase 3)

CRM first, local fallback if capture fails, then delivery-proven alerts. On the AlphaFlux side of this path, langextract (google/langextract, Apache-2.0) turns lead emails, form submissions, and texts into clean CRM fields with every field traceable to its source text; it can run against David's local inference via Ollama. The alert leg is only done when delivery is proven, never on Twilio acceptance alone.

## S7: Measure what the public site does (TODO Phase 4)

First-party analytics, search reporting, and conversion checks feed AlphaFlux without moving private tools into this repo. Static assets are compressed at build time with zopfli (google/zopfli, Apache-2.0) for smallest-possible gzip. Marketing spend is reviewed with meridian (google/meridian, Apache-2.0) for marketing mix modeling on the GPU rig, and doubleml-pipeline (google/doubleml-pipeline, Apache-2.0) for causal estimates of which ad spend actually produces leads, instead of platform-reported attribution.

## S8: Raise the release gate (TODO Phase 5)

Format, tests, lint, vulnerability scanning, secret scanning, accessibility, links, and restore drills all pass before a risky change ships. One gate command runs the whole set.

## S9: Operating rules for this site

- Push to main auto-deploys through Dokploy, so docs and code changes on main are production changes.
- A lead is not done when the form accepts it. It is done when the CRM record exists and the alert is proven delivered.
- Twilio acceptance is not delivery. SMS proof requires a delivered status or recipient confirmation.
- No secret belongs in this repo. The AlphaFlux website lead key lives only in the deployment environment.

## S10: Open questions

- Which single pricing lineup replaces the differing website and phone menu prices? (owner: Skyler and David)
- What email address should receive lead alerts? (owner: David or Skyler)
- How should call tracking be added without replacing the only public number? (owner: David and Skyler)
- Should the "Book online" wording point to a booking flow, and which one? (owner: Skyler and David)
