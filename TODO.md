<!--
INSTRUCTIONS FOR THIS FILE (user and AI: read before editing. Keep this header.)

WHAT THIS FILE IS: the super-detailed task list for this project. It is the
only place work items live. PLAN.md says why and in what order; ENDGOAL.md
says what finished looks like; this file says exactly what remains to do.

EXACT STRUCTURE:
- Phases in order: "## Phase <n>: <name> (PLAN S<m>)".
- Parked work in its own phase at the end: "## Parked".
- Every item is one checkbox line: "- [ ] T-<phase>.<n> <task> (proof: <named
  proof>) (PLAN S<m>)". Parked items use T-P.<n> and name what blocks them.

RULES:
- A box is checked ONLY when its named proof exists (a test, a log line, a
  bench note, a live check). Never check a box on intention or partial work.
- Emulation-proven and bench-proven are different claims; write which one.
- Item IDs and (PLAN S<n>) tags are stable once published; never renumber
  casually. A deliberate renumber updates PLAN.md in the same pass.
- New work enters as new items at the end of their phase. Completed items
  stay, checked, as the record.
- No secrets, passwords, or tokens in this file, ever.
-->
# TODO: Patriot Pest Control

The detailed task list for the Patriot Pest Control public website. Rules: a box is checked only when its proof exists and is named (a test, a log line, a live check, a rendered crawl, or a delivery record). Plan steps live in PLAN.md; each phase here names its PLAN section.

## Phase 0: Foundations and repo standard (PLAN S3)
- [x] T-0.1 Keep the public site as a separate Go service from the AlphaFlux platform (proof: `cmd/server/main.go` composes the health and marketing modules; README records the platform boundary). (PLAN S3)
- [x] T-0.2 Keep production packaging in the repo (proof: `Dockerfile`, `go.mod`, and `go.sum` are present on main). (PLAN S3)
- [x] T-0.3 Add the repo standard files (proof: `TODO.md`, `IDEA.md`, `PLAN.md`, and `docs/RULES.md` are present after this rollout). (PLAN S3)
- [ ] T-0.4 Complete the environment contract in one operations document, covering every variable read by the code and marking which values are secret (proof: a code search matches every documented variable and no secret value appears). (PLAN S3)
- [ ] T-0.5 Add a repeatable local startup check that uses a disposable SQLite catalog and proves the server reaches ready state (proof: script output plus `/ready` response). (PLAN S3)
- [ ] T-0.6 Record the Dokploy rollback path for a bad main push, including the prior image or commit to restore (proof: rollback note reviewed against the live app settings). (PLAN S3)

## Phase 1: Public site and deployment (PLAN S4)
- [x] T-1.1 Serve the public marketing site from this repo (proof: `internal/modules/marketing`, `internal/view`, and the production section in README). (PLAN S4)
- [x] T-1.2 Keep sign-in, account, staff, and dashboard paths pointed at AlphaFlux instead of serving private tools here (proof: route redirects in `cmd/server/main.go` and README). (PLAN S4)
- [x] T-1.3 Return `410 Gone` for legacy `/api/*` paths owned by AlphaFlux (proof: route handler in `cmd/server/main.go`). (PLAN S4)
- [x] T-1.4 Deploy by Dokploy from main (proof: production record in README; live site verified on 2026-10-14). (PLAN S4)
- [ ] T-1.5 Add a post-deploy smoke checklist for `/`, `/contact`, `/signup`, `/search`, `/sitemap.xml`, `/robots.txt`, `/health`, and `/ready` (proof: completed checklist attached to the next production change). (PLAN S4)
- [ ] T-1.6 Verify desktop, tablet, and phone rendering after the next template change (proof: screenshots at all three sizes with no clipped calls to action). (PLAN S4)
- [ ] T-1.7 Add uptime and error-rate alerting for the public site that does not depend on a person noticing a failed page (proof: test alert delivered to the approved operations channel). (PLAN S4)

## Phase 2: Content, local SEO, and brand rules (PLAN S5)
- [x] T-2.1 Maintain the pest library, service pages, service areas, blog, FAQs, prices, referral, search, sitemap, and RSS surfaces in the repo (proof: page and content files under `internal/view` and `internal/data`). (PLAN S5)
- [x] T-2.2 Remove Arizona from the public website after the Arizona operation was sold (proof: live-site verification recorded on 2026-10-14). (PLAN S5)
- [x] T-2.3 Make (509) 471-5767 the only public website number (proof: live-site verification recorded on 2026-10-14). (PLAN S5)
- [ ] T-2.4 Crawl every rendered page and confirm company references use the full name Patriot Pest Control (proof: crawl report with zero unresolved company-name violations). (PLAN S5)
- [ ] T-2.5 Crawl every rendered page, tel link, schema block, and social image and confirm (509) 471-5767 is the only public number (proof: crawl report). (PLAN S5)
- [ ] T-2.6 Keep new service and location content limited to Eastern Washington and North Idaho unless David changes the market decision (proof: content review checklist on each new page). (PLAN S5)
- [ ] T-2.7 Leave Oregon pages unchanged until David gives a new instruction (proof: no Oregon content diff in the market cleanup work). (PLAN S5)
- [ ] T-2.8 Resolve the single pricing lineup with Skyler and David; the website and phone menu must not publish different plan prices (proof: approved lineup recorded and both surfaces updated). (PLAN S5)
- [ ] T-2.9 Decide how call tracking will work without replacing the only public number (proof: decision recorded before any tracking code is added). (PLAN S5)
- [ ] T-2.10 Verify the homepage satisfaction statistic in raw HTML and rendered output so crawlers and visitors see the same value (proof: before and after HTML capture). (PLAN S5)
- [ ] T-2.11 Use one sign-in label in the header and account links (proof: rendered header crawl). (PLAN S5)
- [ ] T-2.12 Make every “Book online” statement link to a real booking path or revise the wording with David and Skyler approval (proof: link crawl and approved copy). (PLAN S5)
- [ ] T-2.13 Add a content review gate for the veteran-owned line, full company name, market, phone number, and claims before publication (proof: checklist in the publishing workflow). (PLAN S5)

## Phase 3: Lead capture and alert delivery (PLAN S6)
- [x] T-3.1 Stop discarding contact form leads; valid submissions now create an AlphaFlux CRM record with website form attribution (proof: live TEST submission and CRM record verification on 2026-10-16). (PLAN S6)
- [x] T-3.2 Preserve first-touch attribution while appending later touches (proof: signup follow-up appended a second touch without replacing the first on 2026-10-16). (PLAN S6)
- [x] T-3.3 Require phone on the contact and signup forms and validate it on the server (proof: no-phone live POST rejected on 2026-10-16). (PLAN S6)
- [x] T-3.4 Add ZIP to the contact form and the approved pest picker to both forms (proof: live form test on 2026-10-16). (PLAN S6)
- [x] T-3.5 Fall back to the local store if AlphaFlux capture fails, so a CRM outage does not discard a lead (proof: fallback path in `internal/modules/marketing/leads.go`). (PLAN S6)
- [ ] T-3.6 Add an operator alert and reconciliation job for leads that land in the local fallback store (proof: a forced capture failure produces an alert and the fallback lead is reconciled to CRM). (PLAN S6)
- [ ] T-3.7 Complete Patriot A2P 10DLC registration for the lead-alert sending path (blocked on: Patriot account registration work; proof required: Twilio shows an approved brand and campaign). (PLAN S6)
- [ ] T-3.8 Re-test SMS lead alerts after A2P approval and prove delivery, not only Twilio acceptance (proof: Twilio delivery status is delivered and Skyler confirms receipt of a labeled TEST). (PLAN S6)
- [ ] T-3.9 Supply and configure Skyler's lead alert email in AlphaFlux (blocked on: David or Skyler confirming the address; never guess it). (PLAN S6)
- [ ] T-3.10 Build the email alert leg in AlphaFlux after the address is configured, with audit rows and tests (proof: labeled TEST email received and audit record found). (PLAN S6)
- [ ] T-3.11 Run a scheduled synthetic lead test with cleanup so form, CRM, attribution, and alert health are proven without waiting for a real customer (proof: dated test report and deleted test contact). (PLAN S6)

## Phase 4: Analytics, search, and conversion (PLAN S7)
- [x] T-4.1 Keep first-party analytics and the beacon sink in the site codebase (proof: analytics files under `internal/data` and `internal/view`, plus the marketing module). (PLAN S7)
- [ ] T-4.2 Publish an event dictionary for page views, calls, form starts, form errors, form submissions, search, and outbound platform redirects (proof: dictionary matches emitted event names in code). (PLAN S7)
- [ ] T-4.3 Add a search report that lists zero-result searches and the page or pest content created in response (proof: first monthly report). (PLAN S7)
- [ ] T-4.4 Measure image weight and Core Web Vitals on the home, pest, area, and blog templates (proof: dated report at phone and desktop sizes). (PLAN S7)
- [ ] T-4.5 Add conversion checks for call buttons and form calls to action on phone layouts (proof: screenshot walkthrough and analytics test events). (PLAN S7)
- [ ] T-4.6 Keep customer and staff reporting in AlphaFlux; this repo only emits the public-site events AlphaFlux needs (proof: integration contract reviewed with the console). (PLAN S7)

## Phase 5: Quality, security, and release gates (PLAN S8)
- [x] T-5.1 Keep Go tests beside the data, marketing, report, view, CSRF, SEO, and dashboard-redirect code (proof: `_test.go` files in the repo tree). (PLAN S8)
- [ ] T-5.2 Make one gate command run format, vet, tests, lint, vulnerability scan, and secret scan (proof: clean gate log on main). (PLAN S8)
- [ ] T-5.3 Add rendered-page accessibility checks for labels, focus order, contrast, keyboard use, and reduced motion (proof: audit report with no critical findings). (PLAN S8)
- [ ] T-5.4 Add a link and sitemap consistency check to the gate (proof: checker fails on a deliberately broken test link and passes on main). (PLAN S8)
- [ ] T-5.5 Review security headers, CSRF coverage, form sanitization, rate limits, and log redaction as one release checklist (proof: signed checklist with code references). (PLAN S8)
- [ ] T-5.6 Add tests for AlphaFlux timeout, non-200 response, malformed form input, and fallback-store behavior (proof: named tests in the marketing package). (PLAN S8)
- [ ] T-5.7 Add a backup and restore drill for the bind-mounted SQLite catalog and settings file (proof: restored copy serves the expected page count in a local run). (PLAN S8)

## Parked (waiting on a decision, hardware, or hands)
- T-P.1 Pricing lineup and booking wording (blocked on: Skyler and David). (PLAN S5)
- T-P.2 Skyler lead alert email address (blocked on: David or Skyler; do not guess). (PLAN S6)
- T-P.3 SMS alert delivery (blocked on: approved A2P brand and campaign). (PLAN S6)
- T-P.4 Any customer, staff, or admin dashboard work (blocked on: belongs in AlphaFlux, not this repo). (PLAN S2)
