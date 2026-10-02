package report

// dashboardHTML is the single-page status dashboard served on
// report.patriotpest.pro to signed-in users. All styling and behavior are
// inline so the report host needs no asset pipeline. The <!--USER_EMAIL-->
// placeholder is replaced server-side with the signed-in address.
const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="robots" content="noindex,nofollow">
<title>Weekly Status Report | Patriot Pest Control</title>
<style>
:root{
  --bg:#0b0f1a; --bg2:#10162a; --panel:#141b33; --panel2:#182040;
  --line:#26305c; --ink:#f2f5ff; --muted:#9aa6cf; --faint:#6b7699;
  --red:#e63946; --red-dk:#b02331; --blue:#3b82f6; --blue-dk:#1d4ed8;
  --green:#2ea043; --green-bg:rgba(46,160,67,.14);
  --amber:#d29922; --amber-bg:rgba(210,153,34,.14);
  --gray:#8b949e; --gray-bg:rgba(139,148,158,.14);
  --radius:12px;
}
*{box-sizing:border-box}
html,body{margin:0;padding:0}
body{background:var(--bg);color:var(--ink);font-family:system-ui,-apple-system,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;line-height:1.55;-webkit-font-smoothing:antialiased}
.app{display:flex;min-height:100vh}

/* ---------- sidebar ---------- */
.sidebar{width:280px;flex:0 0 280px;background:linear-gradient(180deg,#0d1330 0%,#0b0f1a 100%);border-right:1px solid var(--line);position:sticky;top:0;height:100vh;display:flex;flex-direction:column;padding:1.4rem 1.1rem;z-index:40}
.brand{display:flex;align-items:center;gap:.7rem;margin-bottom:2rem}
.brand .mark{width:42px;height:42px;border-radius:10px;background:linear-gradient(135deg,var(--red) 0%,var(--red-dk) 55%,var(--blue-dk) 130%);display:flex;align-items:center;justify-content:center;font-weight:900;font-size:1.3rem;color:#fff;box-shadow:0 4px 14px rgba(230,57,70,.35)}
.brand .t1{font-weight:800;letter-spacing:.12em;font-size:.95rem}
.brand .t2{color:var(--muted);font-size:.72rem;letter-spacing:.22em}
.nav-label{font-size:.68rem;letter-spacing:.24em;color:var(--faint);margin:0 0 .6rem .2rem}
.dates{display:flex;flex-direction:column;gap:.5rem}
.date{display:flex;align-items:center;justify-content:space-between;gap:.5rem;width:100%;text-align:left;background:transparent;border:1px solid transparent;border-radius:9px;color:var(--ink);padding:.7rem .8rem;font-size:.92rem;cursor:pointer;transition:background .15s,border-color .15s}
.date:hover{background:rgba(59,130,246,.08)}
.date.active{background:rgba(59,130,246,.14);border-color:var(--blue)}
.date .dot{width:8px;height:8px;border-radius:50%;background:var(--blue);flex:0 0 8px}
.date.upcoming{color:var(--faint);cursor:not-allowed}
.date.upcoming:hover{background:transparent}
.date .soon{font-size:.62rem;letter-spacing:.16em;border:1px solid var(--line);border-radius:20px;padding:.15rem .5rem;color:var(--faint)}
.side-foot{margin-top:auto;padding-top:1rem;border-top:1px solid var(--line);font-size:.8rem;color:var(--muted)}
.side-foot .who{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;margin-bottom:.4rem}
.side-foot a{color:var(--muted);text-decoration:none}
.side-foot a:hover{color:var(--ink)}

/* ---------- main ---------- */
.main{flex:1;min-width:0;display:flex;flex-direction:column}
.topbar{display:none}
.tabs{position:sticky;top:0;z-index:30;background:rgba(11,15,26,.92);backdrop-filter:blur(8px);border-bottom:1px solid var(--line);padding:.9rem 2rem;display:flex;gap:.6rem}
.tab{border:1px solid var(--line);background:transparent;color:var(--muted);border-radius:999px;padding:.55rem 1.3rem;font-size:.92rem;font-weight:700;letter-spacing:.02em;cursor:pointer;transition:all .15s}
.tab:hover{color:var(--ink);border-color:var(--blue)}
.tab.active{background:linear-gradient(135deg,var(--red) 0%,var(--red-dk) 100%);border-color:transparent;color:#fff;box-shadow:0 4px 14px rgba(230,57,70,.35)}
.tab.active.blue{background:linear-gradient(135deg,var(--blue) 0%,var(--blue-dk) 100%);box-shadow:0 4px 14px rgba(59,130,246,.35)}
.content{padding:2rem;max-width:1080px;width:100%;margin:0 auto}
.report-head{margin:0 0 1.6rem}
.report-head h1{margin:0;font-size:1.7rem;letter-spacing:.01em}
.report-head .date-line{color:var(--muted);font-size:.95rem;margin-top:.25rem}
.report-head .date-line b{color:var(--ink)}

/* stat strip */
.stats{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:.8rem;margin-bottom:1.8rem}
.stat{background:linear-gradient(180deg,var(--panel) 0%,var(--bg2) 100%);border:1px solid var(--line);border-radius:var(--radius);padding:1rem 1.1rem;position:relative;overflow:hidden}
.stat::before{content:"";position:absolute;left:0;top:0;bottom:0;width:4px;background:linear-gradient(180deg,var(--red),var(--blue))}
.stat .v{font-size:1.65rem;font-weight:800;letter-spacing:-.01em}
.stat .k{font-size:.78rem;color:var(--muted);margin-top:.2rem;line-height:1.4}

/* cards */
.card{background:linear-gradient(180deg,var(--panel) 0%,var(--bg2) 100%);border:1px solid var(--line);border-radius:var(--radius);padding:1.4rem 1.5rem;margin-bottom:1.2rem;box-shadow:0 6px 24px rgba(0,0,0,.25)}
.card-head{display:flex;align-items:center;gap:.8rem;flex-wrap:wrap;margin-bottom:.9rem}
.card-head h2{margin:0;font-size:1.15rem}
.card-head .note{width:100%;font-size:.82rem;color:var(--muted);margin-top:-.3rem}
.card ul{margin:0;padding:0;list-style:none}
.card li{position:relative;padding:.45rem 0 .45rem 1.4rem;font-size:.93rem;color:#dbe2ff}
.card li::before{content:"";position:absolute;left:.2rem;top:.95em;width:7px;height:7px;border-radius:50%;background:var(--blue)}
.card.money li::before{background:var(--green)}
.card.money{border-left:4px solid var(--green)}
.pill{display:inline-flex;align-items:center;gap:.35rem;font-size:.68rem;font-weight:800;letter-spacing:.14em;border-radius:999px;padding:.28rem .75rem;text-transform:uppercase}
.pill::before{content:"";width:7px;height:7px;border-radius:50%;background:currentColor}
.pill.done{color:#3fb950;background:var(--green-bg);border:1px solid rgba(46,160,67,.4)}
.pill.progress{color:#e3b341;background:var(--amber-bg);border:1px solid rgba(210,153,34,.4)}
.pill.pending{color:#a7b0be;background:var(--gray-bg);border:1px solid rgba(139,148,158,.35)}
.foot{margin:2.5rem 0 1rem;color:var(--faint);font-size:.78rem;text-align:center}

/* drawer scrim */
.scrim{display:none}

/* ---------- mobile ---------- */
@media (max-width:860px){
  .sidebar{position:fixed;left:0;top:0;transform:translateX(-105%);transition:transform .22s ease;box-shadow:8px 0 30px rgba(0,0,0,.5)}
  body.nav-open .sidebar{transform:translateX(0)}
  body.nav-open .scrim{display:block;position:fixed;inset:0;background:rgba(0,0,0,.55);z-index:35}
  .topbar{display:flex;align-items:center;gap:.8rem;position:sticky;top:0;z-index:30;background:rgba(11,15,26,.94);backdrop-filter:blur(8px);border-bottom:1px solid var(--line);padding:.7rem 1rem}
  .menu{background:transparent;border:1px solid var(--line);border-radius:8px;color:var(--ink);font-size:1.1rem;padding:.35rem .6rem;cursor:pointer}
  .topbar .tt{font-weight:800;letter-spacing:.1em;font-size:.85rem}
  .tabs{padding:.7rem 1rem;overflow-x:auto}
  .tab{white-space:nowrap}
  .content{padding:1.2rem 1rem}
  .report-head h1{font-size:1.35rem}
  .stats{grid-template-columns:repeat(2,1fr)}
  .card{padding:1.1rem}
}
@media (prefers-reduced-motion:reduce){*{transition:none!important}}
</style>
</head>
<body>
<div class="app">
  <aside class="sidebar" id="sidebar" aria-label="Report navigation">
    <div class="brand">
      <div class="mark">P</div>
      <div><div class="t1">PATRIOT STATUS</div><div class="t2">WEEKLY REPORTS</div></div>
    </div>
    <p class="nav-label">REPORT DATES</p>
    <nav class="dates" id="dateNav">
      <button class="date active" data-date="2026-10-01"><span>October 1, 2026</span><span class="dot"></span></button>
      <button class="date upcoming" disabled><span>Week of October 8</span><span class="soon">SOON</span></button>
    </nav>
    <div class="side-foot">
      <span class="who">Signed in as <!--USER_EMAIL--></span>
      <a href="/logout">Sign out</a>
    </div>
  </aside>
  <div class="scrim" id="scrim"></div>

  <div class="main">
    <header class="topbar">
      <button class="menu" id="menuBtn" aria-label="Open report dates">&#9776;</button>
      <div class="tt">PATRIOT STATUS</div>
    </header>
    <div class="tabs" role="tablist">
      <button class="tab active" data-tab="patriot" role="tab">Patriot Pest Control</button>
      <button class="tab blue" data-tab="alphaflux" role="tab">AlphaFlux</button>
    </div>
    <div class="content" id="content"></div>
    <p class="foot">Weekly status report for Skyler Rose. New entries are added after each weekly call.</p>
  </div>
</div>

<script>
const REPORTS = {
"2026-10-01": {
  label: "October 1, 2026",
  patriot: {
    stats: [
      {v:"6", k:"website pushes shipped and verified live"},
      {v:"90", k:"Google reviews answered (86 positive + 4 negative)"},
      {v:"2", k:"Google Business Profiles cleaned up"},
      {v:"25", k:"pest pages rewritten"},
      {v:"15", k:"city pages expanded"},
      {v:"$4,000", k:"Keystone proposal avoided by doing the work in-house"}
    ],
    cards: [
      {title:"Website rebuild", pill:"done", cls:"", note:"",
       points:[
        "Six updates shipped to patriotpest.pro, each verified live.",
        "Fixed the behind-the-scenes business data Google reads: valid structured data, real privacy policy and terms pages, removed a test page that was showing in search, bigger tap buttons on phones, cleaner page addresses.",
        "Set the official phone numbers everywhere: (509) 818-0993 for WA/ID/OR and (602) 755-8414 for AZ. The old number is gone from the site.",
        "Rewrote all 25 pest pages with unique copy and correct local species. Added local-area write-ups to all 15 city pages. Expanded thin blog posts. Fixed small text and low contrast.",
        "SEO extras: real dates on the sitemap, custom preview images for pest, city, and blog pages, breadcrumb navigation data, better internal links, plus an automated test that crawls 63 pages so nothing breaks again.",
        "Hours standardized everywhere: Mon to Sat 7am to 7pm, Sunday closed, online scheduling 24/7.",
        "Official business name set to \u201CPatriot Pest Control CO\u201D in the footer, business data, and legal pages."
      ]},
      {title:"Google Business Profiles", pill:"done", cls:"", note:"A few items are still in Google review.",
       points:[
        "Two profiles: Washington/Idaho (98 reviews, 4.7 stars) and Arizona (Mesa/Tempe, 0 reviews).",
        "Audited both profiles end to end against the live Google data.",
        "Killed a bad pending edit that was trying to switch the WA phone back to the old number.",
        "Fixed the website link to the secure https address.",
        "Added the veteran-owned attribute to the WA profile.",
        "Rewrote the Arizona description. It was talking about Spokane and Northwest pests. Now it covers Mesa, Tempe, and desert pests.",
        "Fixed Arizona Saturday hours. Was 9 to 5, now 7am to 7pm like the rest of the week.",
        "Renamed the Arizona profile to \u201CPatriot Pest Control CO\u201D to match everywhere else.",
        "Corrected the opening date to 2016.",
        "Turned on chat messaging on both profiles.",
        "Replied to 86 positive reviews in a warm, plain-spoken style, each one personalized.",
        "Replied to 4 negative reviews with calm, factual owner responses (billing and service complaints).",
        "Reported one fake 1-star review from 2019 to Google for removal. The reviewer was never a customer. Decision expected within about 3 business days.",
        "Left 2 old negative reviews alone on purpose: one from 2022 that replying would only wake up, and the reported fake one."
      ]},
      {title:"Sameday AI phone agent", pill:"progress", cls:"", note:"Review request text campaign launched and active.",
       points:[
        "Decision: keeping Sameday for the full year already paid for. No cancellation. We keep improving it: prompts, call routing, transfers, and matching callers to FieldRoutes.",
        "This is Skyler's company, so transfers were never routed to David. The AI handles calls itself and books the job; when a human is needed, it goes to Skyler's number.",
        "Launched Oct 2 at 9:36 AM EDT: the \u201CReview Request - Post Service\u201D campaign is active and sent review request texts to 34 customers from September's completed appointments. Texts land around 9:36 AM Phoenix time asking happy customers for a Google review.",
        "Campaign goal: collect Google reviews from recent service customers.",
        "Note: Sameday cannot send customer emails, so review emails will need to come from another system later."
      ]},
      {title:"The $4,000 proposal", pill:"done", cls:"money", note:"Reviewed and already handled in-house.",
       points:[
        "Skyler shared Keystone's $4,000 \u201CFoundation Package\u201D proposal documents.",
        "Within an hour of getting the documents, their entire technical scope had been reviewed and confirmed: nearly everything they proposed was already built, fixed, and live in-house.",
        "Website rebuild, Google profiles, phones, hours, review replies, and sitemap: done before their plan would have even started.",
        "Verdict: do not pay it. The work is done."
      ]},
      {title:"Up next", pill:"pending", cls:"", note:"",
       points:[
        "Get the Arizona profile its first reviews.",
        "Fresh photos for both profiles (some old graphics still show the old phone number).",
        "Clean up the Washington services list (Google added junk entries like Lice).",
        "Publish fall posts (yellow jacket and rodent season).",
        "Confirm the pending Google edits all clear review.",
        "External listings cleanup: Yelp, BBB, Facebook, Apple, Bing, Nextdoor.",
        "Website: answer-style content blocks for AI search, full rewrites of the Deer Park and Spokane pages."
      ]}
    ]
  },
  alphaflux: {
    stats: [
      {v:"Sep 28", k:"AlphaFlux LLC formed and active in Kentucky"},
      {v:"Sep 30", k:"EIN issued by the IRS"},
      {v:"14", k:"commits: inventory module shipped live"},
      {v:"16", k:"social media connectors merged"},
      {v:"27", k:"Twilio phone API endpoints wired"}
    ],
    cards: [
      {title:"Everything done", pill:"done", cls:"", note:"",
       points:[
        "Business: AlphaFlux LLC filed and approved in Kentucky (Sep 28, veteran fee waiver, $0). EIN issued by the IRS (Sep 30). Company logo created.",
        "Console: split into an admin shell and a customer shell with audit logging and a help mode. Mobile navigation restyled. Dashboard restyled. Call flow builder redesigned with flow, kanban, and code views.",
        "Twilio coverage: 27 phone API endpoints wired into the Phone tab, including conferences, queues, SIP, Verify, Lookup, Studio flows, TrustHub, and usage alerts.",
        "Inventory module: full build with districts, warehouses, trucks, barcode scanning with offline retry, chemical sign in/out ledger, QR pages, low stock alerts, and purchase orders. Live in production.",
        "Social: 16-platform connector framework merged (OAuth, token refresh, scheduled publishing).",
        "Research: 8 competitor reports covering every console tab, each with top 10 competitors, pricing, and real customer complaints.",
        "Bank: Relay business bank application started (Mercury had declined with no reason).",
        "Stripe: account created as AlphaFlux LLC. Activation in progress.",
        "Twilio: parent business profile filed and in manual review (Oct 1).",
        "Photos: Immich self-hosted photo platform live and updated.",
        "Product: MowStrip weed barrier, 20+ product renders in the approved style."
      ]},
      {title:"What is left", pill:"pending", cls:"", note:"",
       points:[
        "Relay bank account: needs date of birth, a password set on the secure card, and the final review click.",
        "Stripe activation: needs date of birth plus last 4 of SSN, then bank payout details and possible ID check.",
        "Twilio: identity verification email (check david@itak.live), then the Patriot brand and campaign filing (needs Skyler\u2019s EIN), plus cleanup of 11 junk draft brands.",
        "Mail server: the real certificate fix is still pending (needs a DNS record or container logs).",
        "Logins: OAuth app registrations stalled waiting on current passwords.",
        "Console: database scaling and backup, remaining tab build-out, Field Services photo app, Sales/Leads dashboard section.",
        "Mimir coding harness: two patches built and tested, waiting on the go-ahead. The full harness spec is written and awaiting approval.",
        "The big goal: first dollar of revenue needs live Stripe charges plus one paying tenant."
      ]}
    ]
  }
}};

let curDate = "2026-10-01";
let curTab = "patriot";
const PILL = {done:"Done", progress:"In Progress", pending:"Pending"};

function render(){
  const rep = REPORTS[curDate];
  const tab = rep[curTab];
  const el = document.getElementById("content");
  let h = '<div class="report-head"><h1>' + (curTab === "patriot" ? "Patriot Pest Control" : "AlphaFlux") + '</h1>' +
          '<div class="date-line">Week of <b>' + rep.label + '</b></div></div>';
  h += '<div class="stats">';
  tab.stats.forEach(function(s){
    h += '<div class="stat"><div class="v">' + s.v + '</div><div class="k">' + s.k + '</div></div>';
  });
  h += '</div>';
  tab.cards.forEach(function(c){
    h += '<section class="card ' + c.cls + '"><div class="card-head"><h2>' + c.title + '</h2>' +
         '<span class="pill ' + c.pill + '">' + PILL[c.pill] + '</span>';
    if (c.note) h += '<span class="note">' + c.note + '</span>';
    h += '</div><ul>';
    c.points.forEach(function(p){ h += '<li>' + p + '</li>'; });
    h += '</ul></section>';
  });
  el.innerHTML = h;
  document.querySelectorAll(".tab").forEach(function(t){
    t.classList.toggle("active", t.dataset.tab === curTab);
  });
  document.querySelectorAll(".date").forEach(function(d){
    d.classList.toggle("active", d.dataset.date === curDate);
  });
}

document.querySelectorAll(".tab").forEach(function(t){
  t.addEventListener("click", function(){ curTab = t.dataset.tab; render(); });
});
document.querySelectorAll(".date[data-date]").forEach(function(d){
  d.addEventListener("click", function(){
    curDate = d.dataset.date; render();
    document.body.classList.remove("nav-open");
  });
});
document.getElementById("menuBtn").addEventListener("click", function(){
  document.body.classList.toggle("nav-open");
});
document.getElementById("scrim").addEventListener("click", function(){
  document.body.classList.remove("nav-open");
});
render();
</script>
</body>
</html>`

// authCSS styles the standalone sign-in pages.
const authCSS = `:root{--bg:#0b0f1a;--panel:#141b33;--line:#26305c;--ink:#f2f5ff;--muted:#9aa6cf;--red:#e63946;--red-dk:#b02331;--blue:#3b82f6}
*{box-sizing:border-box}
body{margin:0;background:radial-gradient(1200px 600px at 50% -10%,#16204a 0%,var(--bg) 60%);color:var(--ink);font-family:system-ui,-apple-system,"Segoe UI",Roboto,sans-serif;min-height:100vh}
.auth{min-height:100vh;display:flex;align-items:center;justify-content:center;padding:1.2rem}
.card{width:100%;max-width:400px;background:var(--panel);border:1px solid var(--line);border-radius:14px;padding:2rem;box-shadow:0 12px 40px rgba(0,0,0,.45)}
.card h1{margin:0;font-size:1.05rem;letter-spacing:.22em}
.card h1::before{content:"";display:inline-block;width:10px;height:10px;border-radius:3px;background:linear-gradient(135deg,var(--red),var(--blue));margin-right:.6rem}
.sub{color:var(--muted);font-size:.85rem;margin:.4rem 0 0}
.intro{font-size:.92rem;margin:1.2rem 0 0}
label{display:block;font-size:.68rem;letter-spacing:.18em;color:var(--muted);margin:1.1rem 0 .35rem}
input{width:100%;padding:.7rem .8rem;background:var(--bg);border:1px solid var(--line);border-radius:8px;color:var(--ink);font-size:1rem}
input:focus{outline:2px solid var(--blue);border-color:var(--blue)}
button{width:100%;margin-top:1.3rem;padding:.75rem;background:linear-gradient(135deg,var(--red),var(--red-dk));border:0;border-radius:8px;color:#fff;font-weight:800;letter-spacing:.1em;cursor:pointer;font-size:.95rem}
button:hover{filter:brightness(1.1)}
.err{background:rgba(230,57,70,.12);border:1px solid var(--red);color:#f2b8a8;border-radius:8px;padding:.65rem .85rem;font-size:.85rem;margin-top:1.1rem}
.alt{margin-top:1rem;font-size:.82rem;text-align:center}
.alt a{color:var(--muted)}
`
