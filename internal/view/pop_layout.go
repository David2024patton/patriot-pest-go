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
  <meta name="viewport" content="width=device-width, initial-scale=1.0, interactive-widget=resizes-content">
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
  <link rel="stylesheet" href="/assets/pop.css?v=20261104-pop5">
  <link rel="stylesheet" href="/assets/beacon.css">
  <link rel="icon" href="/assets/img/pests/ants.jpg" type="image/jpeg">
</head>
<body class="pop{{if .BodyClass}} {{.BodyClass}}{{end}}">

<header class="pop-head">
  <div class="pop-head-row">
    <a class="pop-brand" href="/"><span class="pop-star">★</span> PATRIOT PEST <span class="pop-brand-accent">CONTROL</span></a>
    <a class="pop-pill pop-call" href="{{.PhoneHref}}"><span class="call-long">Call or Text </span>{{.PhoneDisplay}}</a>
  </div>
  <nav class="pop-nav" aria-label="Main navigation">
    <a class="pop-pill{{if or (eq .Page "home-pop") (eq .Page "home")}} active{{end}}" href="/">Home</a>
    <a class="pop-pill{{if eq .Page "services"}} active{{end}}" href="/services">Services</a>
    <a class="pop-pill{{if eq .Page "prices"}} active{{end}}" href="/prices">Prices</a>
    <a class="pop-pill{{if or (eq .Page "areas") (eq .Page "area-detail")}} active{{end}}" href="/service-areas">Areas</a>
    <a class="pop-pill{{if or (eq .Page "blog-index") (eq .Page "blog-post")}} active{{end}}" href="/blogs">Blog</a>
    <a class="pop-pill{{if eq .Page "faqs"}} active{{end}}" href="/faqs">FAQs</a>
    <a class="pop-pill{{if eq .Page "contact"}} active{{end}}" href="/contact">Contact</a>
    <a class="pop-pill" href="/about">About</a>
    <a class="pop-pill" href="/signup">Sign Up</a>
  </nav>
</header>

<main>
{{if and .Crumb (gt (len .Crumb) 1)}}
  <div class="wrap" style="padding-top:1.1rem">
    <nav class="crumb" aria-label="Breadcrumb">
      {{range $i, $c := .Crumb}}{{if $i}}<span class="sep">/</span> {{end}}{{if lt $i (sub (len $.Crumb) 1)}}<a href="{{index $c 1}}">{{index $c 0}}</a>{{else}}<span>{{index $c 0}}</span>{{end}}{{end}}
    </nav>
  </div>
{{end}}
{{.PageBody}}
</main>

<footer class="pop-foot">
  <div class="wrap">
    <div class="pop-foot-grid">
      <div>
        <div class="pop-foot-brand"><span class="pop-star">★</span> PATRIOT PEST CONTROL</div>
        <p>Veteran-owned pest control for homes &amp; businesses across Washington, Idaho &amp; Oregon. Founded by U.S. Military Veteran Skyler Rose. Eco-friendly, family &amp; pet safe, 100% satisfaction guaranteed.</p>
      </div>
      <div>
        <h3>Quick Links</h3>
        <a href="/about">About Us</a><a href="/services">Services</a>
        <a href="/prices">Pricing</a><a href="/blogs">Blog</a>
        <a href="/faqs">FAQs</a><a href="/referral">Referral Program</a>
        <a href="/help">Help Center</a><a href="/links">All Links</a>
      </div>
      <div>
        <h3>Top Services</h3>
        <a href="/pest/ants">Ant Control</a><a href="/pest/termites">Termite Control</a>
        <a href="/pest/bed-bugs">Bed Bug Treatment</a><a href="/pest/rodents">Rodent Control</a>
        <a href="/pest/mosquitoes">Mosquito Control</a><a href="/pest/wasps">Wasp Removal</a>
      </div>
      <div>
        <h3>Contact</h3>
        <a href="{{.PhoneHref}}">{{.PhoneDisplay}} - {{.PhoneLabel}}</a>
        <a href="mailto:info@patriotpest.pro">info@patriotpest.pro</a>
        <a href="/contact">Deer Park, WA 99006, United States</a>
        <a href="/socials">Social Media</a>
      </div>
    </div>
    <div class="pop-foot-bottom">
      <span>© {{.Year}} PATRIOT PEST CONTROL CO · ALL RIGHTS RESERVED</span>
      <span><a href="/privacy-policy">PRIVACY</a> · <a href="/terms-of-use">TERMS</a> · <a href="/sitemap">SITEMAP</a></span>
      <span>VETERAN-OWNED AMERICAN COMPANY</span>
    </div>
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

<script src="/assets/js/pop.js?v=20261104-pop1" defer></script>
{{if .Track}}<script src="/assets/js/tracker.js" defer></script>{{end}}
</body>
</html>
`
