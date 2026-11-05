package view

import "strings"

// Pop UI host: pop.patriotpest.pro serves the same site, same data, same
// handlers, wearing the Pop UI build (tile rails, chunky cards, popup-first)
// in Patriot's own olive and orange. The main hosts never see this layout.

// IsPopHost reports whether the request arrived on the Pop UI host.
// Exact match, case-insensitive, port stripped (mirrors the report host).
func IsPopHost(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	if i := strings.LastIndex(h, ":"); i > -1 {
		h = h[:i]
	}
	return h == "pop.patriotpest.pro"
}

// layoutPopHTML is the Pop UI shell: pill header on the olive ground,
// chunky footer card, white phone bottom bar. Page bodies render inside
// <main> exactly as in the main layout; pop.css styles both the shell and
// the shared body classes.
const layoutPopHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  {{template "analytics" .}}
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover, interactive-widget=resizes-content">
  <title>{{.Title}}</title>
  {{if .Description}}<meta name="description" content="{{.Description}}">{{end}}
  <meta name="robots" content="{{.Robots}}">
  <link rel="canonical" href="{{.Canonical}}">
  <meta name="theme-color" content="#262f14">
  <meta property="og:type" content="website">
  <meta property="og:site_name" content="Patriot Pest Control">
  <meta property="og:title" content="{{.Title}}">
  <meta property="og:description" content="{{.Description}}">
  <meta property="og:url" content="{{.Canonical}}">
  <meta property="og:image" content="{{.OGImage}}">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Black+Ops+One&family=Barlow:wght@400;500;600;700;800;900&display=swap" rel="stylesheet">
  <link rel="stylesheet" href="/assets/pop.css?v=20261105-pop21">
  <link rel="stylesheet" href="/assets/beacon.css">
  <link rel="icon" href="/assets/img/pests/ants.jpg" type="image/jpeg">
</head>
<body class="pop{{if .BodyClass}} {{.BodyClass}}{{end}}">

<header class="pop-head">
  <div class="pop-head-in">
    <a class="pop-brand" href="/"><span class="pop-star">★</span> PATRIOT PEST <span class="pop-brand-accent">CONTROL</span></a>
  <nav class="pop-nav" aria-label="Main navigation">
    <a class="pop-pill" href="/#plans" data-scroll="plans"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.6 2.9 21.1 11.4a2 2 0 0 1 0 2.8l-6.9 6.9a2 2 0 0 1-2.8 0L2.9 12.6A2 2 0 0 1 2.3 11.2V4a1.7 1.7 0 0 1 1.7-1.7h7.2a2 2 0 0 1 1.4.6z"/><circle cx="7.5" cy="7.5" r="1.4"/></svg>Prices</a>
    <a class="pop-pill" href="/#services" data-scroll="services"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="13" r="5"/><path d="M12 8V4M8.5 9.5 5 6M15.5 9.5 19 6M7 13H3M21 13h-4M9 17l-3 4M15 17l3 4"/></svg>Services</a>
    <a class="pop-pill" href="/#about" data-scroll="about"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3l7 3v5c0 4.6-3 8.1-7 10-4-1.9-7-5.4-7-10V6z"/><path d="m9.3 11.8 2 2 3.6-4"/></svg>About</a>
    <a class="pop-pill" href="/#areas" data-scroll="areas"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 21s-7-5.6-7-11a7 7 0 0 1 14 0c0 5.4-7 11-7 11z"/><circle cx="12" cy="10" r="2.6"/></svg>Areas</a>
    <a class="pop-pill" href="/#blog" data-scroll="blog"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="15" height="16" rx="2"/><path d="M18 9h1a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2"/><path d="M7 9h7M7 13h7M7 17h5"/></svg>Blog</a>
    <a class="pop-pill" href="/#contact" data-scroll="contact"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="m22 7-10 6L2 7"/></svg>Contact</a>
  </nav>
  </div>
</header>

<main>

{{.PageBody}}
</main>

<!-- ============ POPUP LAYER 1: QUICK VIEW ============ -->
<div class="pop-scrim" id="pop-qv" hidden>
  <div class="pop-popup" role="dialog" aria-modal="true" aria-label="Quick view">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-qv-body"></div>
    <div class="pop-popup-row">
      <a class="pop-btn primary" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
      <button type="button" class="pop-btn dark pop-qv-more">More Info</button>
    </div>
  </div>
</div>

<!-- ============ POPUP LAYER 2: MORE INFO ============ -->
<div class="pop-scrim" id="pop-more" hidden>
  <div class="pop-popup" role="dialog" aria-modal="true" aria-label="More info">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-more-body"></div>
    <div class="pop-popup-row">
      <a class="pop-btn primary" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
    </div>
  </div>
</div>

<!-- ============ POPUP: FREE QUOTE (the real /contact form) ============ -->
<div class="pop-scrim" id="pop-quote" {{if .Errors}}data-autopen{{end}} hidden>
  <div class="pop-popup" role="dialog" aria-modal="true" aria-label="Free quote">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-lab">Free Quote</div>
    <h3>Get your free quote</h3>
    <p class="pop-form-note">We respond within one business day, usually much faster. Same-day service available across all three states.</p>
    <div class="pop-promo-line" id="pop-promo-line" hidden></div>
    <form method="post" action="/contact" class="form-stack">
      {{.Csrf}}
      <div class="field"><label for="q-name">Full Name *</label><input type="text" id="q-name" name="name" required autocomplete="name"></div>
      <div class="form-row">
        <div class="field"><label for="q-email">Email Address *</label><input type="email" id="q-email" name="email" required autocomplete="email"></div>
        <div class="field"><label for="q-phone">Phone Number *</label><input type="tel" id="q-phone" name="phone" required autocomplete="tel"></div>
      </div>
      <div class="form-row">
        <div class="field"><label for="q-pest">Pest Problem</label><select id="q-pest" name="pest">{{raw .PestOptions}}</select></div>
        <div class="field"><label for="q-zip">ZIP Code</label><input type="text" id="q-zip" name="zip" maxlength="10" autocomplete="postal-code"></div>
      </div>
      <div class="field"><label for="q-message">Message *</label><textarea id="q-message" name="message" required placeholder="Tell us about the pest, your location, and preferred timing..."></textarea></div>
      ` + smsConsent + `
      <div class="form-actions"><button type="submit" class="pop-btn primary">Send Message</button></div>
    </form>
  </div>
</div>


<footer class="pop-foot">
  <div class="wrap">
    <p class="pop-copy">&copy; {{.Year}} Patriot Pest Control Co. All rights reserved. &middot; <a href="/privacy-policy">Privacy</a> &middot; <a href="/terms-of-use">Terms</a> &middot; <a href="/sitemap">Sitemap</a><br><span>Veteran-Owned American Company</span></p>
  </div>
</footer>
<div class="pop-bottom-spacer"></div>

<nav class="pop-bottombar" aria-label="Mobile quick navigation">
  <a href="/"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 10.5 12 3l9 7.5"/><path d="M5 9.5V21h14V9.5"/></svg>Home</a>
  <a href="/prices"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v20"/><path d="M17 6.5c-.8-1.2-2.6-2-5-2-3 0-5 1.3-5 3.2 0 4.6 10 2.4 10 6.6 0 1.9-2 3.2-5 3.2-2.4 0-4.2-.8-5-2"/></svg>Prices</a>
  <a class="pop-bottom-call" href="{{.PhoneHref}}"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.79 19.79 0 0 1 2.12 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z"/></svg>Call</a>
  <a href="/services"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="13" r="5"/><path d="M12 8V4M8.5 9.5 5 6M15.5 9.5 19 6M7 13H3M21 13h-4M9 17l-3 4M15 17l3 4"/></svg>Pests</a>
  <a href="/contact"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="m22 7-10 6L2 7"/></svg>Quote</a>
</nav>


<button type="button" class="pop-fab gear" id="pop-gear" aria-label="Display settings"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg></button>

<button type="button" class="pop-fab ask" id="pop-ask" aria-label="Ask Emily"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg><span>Ask Emily</span></button>

<div class="pop-scrim" id="pop-settings" hidden>
  <div class="pop-popup settings-card" role="dialog" aria-modal="true" aria-label="Display settings">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-lab">Settings</div>
    <h3>Make it yours</h3>
    <div class="set-row"><span class="set-lab">Text size</span><div class="set-seg" id="set-text"><button type="button" data-fs="">Normal</button><button type="button" data-fs="lg">Large</button><button type="button" data-fs="xl">Extra large</button></div></div>
    <div class="set-row"><span class="set-lab">Tile size</span><div class="set-seg" id="set-tiles"><button type="button" data-ts="">Normal</button><button type="button" data-ts="lg">Large</button></div></div>
    <div class="set-row"><span class="set-lab">Jump</span><button type="button" class="pop-btn" id="set-top">Back to top</button></div>
    <p class="pop-form-note">Settings stay on this device.</p>
  </div>
</div>

<div class="em-backdrop" id="em-backdrop" hidden></div>
<div class="em-panel" id="em-panel" hidden>
  <div class="em-head"><div class="em-id"><span class="em-ava">E</span><div><b>Emily</b><small>Patriot Pest Control assistant</small></div></div><button type="button" class="em-x" id="em-close" aria-label="Close">&#10005;</button></div>
  <div class="em-body" id="em-body"></div>
  <div class="em-chips" id="em-chips"></div>
  <form class="em-input" id="em-form"><input type="text" id="em-text" placeholder="Ask about pests, prices, booking..." autocomplete="off"><button type="submit" class="pop-btn primary">Send</button></form>
</div>

<div class="pop-scrim" id="pop-promo" hidden>
  <div class="pop-popup promo-card" role="dialog" aria-modal="true" aria-label="Special offer">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-scroll">
    <img class="promo-art" src="/assets/img/illus/hero.jpg" alt="Illustrated neighborhood of homes protected by Patriot Pest Control">
    <div class="pop-lab" id="promo-lab">Special Offer</div>
    <h3 id="promo-title"></h3>
    <p class="promo-text" id="promo-text"></p>
    <div class="promo-flex">
      <div class="promo-codebox"><small>Your code</small><b id="promo-code"></b></div>
      <div class="promo-qr"><img id="promo-qr" src="" alt="QR code: scan to book with this code applied"><small>SCAN ME to book with the code applied</small></div>
    </div>
    </div>
    <div class="pop-popup-row"><button type="button" class="pop-btn primary" id="promo-book">Book Now</button><button type="button" class="pop-btn dark" data-pop-close>Maybe Later</button></div>
  </div>
</div>

<div class="pop-scrim" id="pop-page" hidden>
  <div class="pop-popup page-card" role="dialog" aria-modal="true" aria-label="Page viewer">
    <button type="button" class="pop-x big" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-page-body"></div>
  </div>
</div>

<div class="pop-scrim" id="pop-booking" hidden>
  <div class="pop-popup booking-card" role="dialog" aria-modal="true" aria-label="Book your plan">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-scroll">
    <div class="bk-main">
    <div class="pop-lab" id="bk-lab">Book Your Plan</div>
    <h3>Pick your plan. Lock in the year.</h3>
    <div class="pop-promo-line" id="bk-offer-line" hidden></div>
    <form id="bk-form" method="post" action="/signup" novalidate>
      {{.Csrf}}
      <input type="hidden" name="plan" id="bk-plan-input" value="Gold">
      <input type="hidden" name="billing" id="bk-billing-input" value="yearly">
      <input type="hidden" name="promo" id="bk-promo-input" value="">
      <div class="bk-plans" id="bk-plans">
        <label class="bk-plan"><input type="radio" name="bk-plan-radio" value="Bronze"><span><b>Bronze</b><small>Exterior-Only Protection</small></span><span class="bk-price">$49/mo</span></label>
        <label class="bk-plan"><input type="radio" name="bk-plan-radio" value="Silver"><span><b>Silver</b><small>Interior + Exterior</small></span><span class="bk-price">$59/mo</span></label>
        <label class="bk-plan on"><input type="radio" name="bk-plan-radio" value="Gold" checked><span><b>Gold</b><small>Priority Interior &amp; Exterior</small></span><span class="bk-price">$79/mo</span></label>
      </div>
      <div class="bk-bills">
        <button type="button" class="bk-bill on" data-billing="yearly"><b>Pay for the year</b><span class="bk-big" id="bk-yearly-big"></span><small class="bk-note" id="bk-yearly-note"></small></button>
        <button type="button" class="bk-bill" data-billing="monthly"><b>Pay monthly</b><span class="bk-big" id="bk-monthly-big"></span><small class="bk-note" id="bk-monthly-note"></small></button>
      </div>
      <p class="bk-note">Every plan is a 12-month plan with the 90-day guarantee. No hidden fees, ever.</p>
      <div class="bk-errors" id="bk-errors" hidden></div>
      <div class="field"><label for="bk-name">Full Name *</label><input type="text" id="bk-name" name="name" required autocomplete="name"></div>
      <div class="form-row">
        <div class="field"><label for="bk-email">Email Address *</label><input type="email" id="bk-email" name="email" required autocomplete="email"></div>
        <div class="field"><label for="bk-phone">Phone Number *</label><input type="tel" id="bk-phone" name="phone" required autocomplete="tel"></div>
      </div>
      <div class="field"><label for="bk-address">Service Address</label><input type="text" id="bk-address" name="address" autocomplete="street-address"></div>
      <div class="form-row">
        <div class="field"><label for="bk-city">City</label><input type="text" id="bk-city" name="city" autocomplete="address-level2"></div>
        <div class="field"><label for="bk-state">State</label><input type="text" id="bk-state" name="state" maxlength="2" placeholder="WA" autocomplete="address-level1"></div>
        <div class="field"><label for="bk-zip">ZIP Code</label><input type="text" id="bk-zip" name="zip" maxlength="10" autocomplete="postal-code"></div>
      </div>
      <div class="form-row">
        <div class="field"><label for="bk-pest">Main Pest Concern</label><select id="bk-pest" name="pest">{{raw .PestOptions}}</select></div>
        <div class="field"><label for="bk-start">Preferred Start Date</label><input type="date" id="bk-start" name="start"></div>
      </div>
      ` + smsConsent + `
      <div class="form-actions"><button type="submit" class="pop-btn primary" id="bk-submit">Book My Plan</button></div>
      <p class="bk-note">Booking sends your plan straight to the Patriot team. They confirm your first visit and set up payment with you. Nothing is charged on this page.</p>
    </form>
    </div>
    <div class="bk-done" id="bk-done" hidden>
      <div class="bk-check">&#10003;</div>
      <h3>You are booked in.</h3>
      <p id="bk-done-text">Your plan is with the Patriot team.</p>
      <p class="bk-note">Questions right now? Call or text (509) 471-5767.</p>
      <div class="pop-cta-row" style="justify-content:center"><button type="button" class="pop-btn primary" data-pop-close>Done</button></div>
    </div>
    </div>
  </div>
</div>

<script src="/assets/js/pop.js?v=20261105-pop21" defer></script>
{{if .Track}}<script src="/assets/js/tracker.js" defer></script>{{end}}
</body>
</html>
`
