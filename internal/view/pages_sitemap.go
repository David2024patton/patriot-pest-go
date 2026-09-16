package view

// pageSitemap — templates/pages/sitemap.php. HTML sitemap.
// Data: .States ([]data.State).
const pageSitemap = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">NAVIGATION // SITEMAP</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Sitemap</h1>
    <p class="lead">Every page on the Patriot Pest Control website.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="grid g3">
      <div class="card">
        <h3 style="font-family:var(--display);color:var(--orange);margin-bottom:.8rem">Pages</h3>
        <div style="display:flex;flex-direction:column;gap:.5rem">
          <a href="/" style="color:var(--cream);text-decoration:none">Home ▸</a>
          <a href="/about" style="color:var(--cream);text-decoration:none">About Us ▸</a>
          <a href="/services" style="color:var(--cream);text-decoration:none">Services ▸</a>
          <a href="/prices" style="color:var(--cream);text-decoration:none">Pricing ▸</a>
          <a href="/service-areas" style="color:var(--cream);text-decoration:none">Service Areas ▸</a>
          <a href="/blogs" style="color:var(--cream);text-decoration:none">Blog ▸</a>
          <a href="/faqs" style="color:var(--cream);text-decoration:none">FAQs ▸</a>
          <a href="/contact" style="color:var(--cream);text-decoration:none">Contact ▸</a>
          <a href="/referral" style="color:var(--cream);text-decoration:none">Referral Program ▸</a>
          <a href="/socials" style="color:var(--cream);text-decoration:none">Social Media ▸</a>
          <a href="/help" style="color:var(--cream);text-decoration:none">Help Center ▸</a>
          <a href="/links" style="color:var(--cream);text-decoration:none">All Links ▸</a>
          <a href="/privacy-policy" style="color:var(--cream);text-decoration:none">Privacy Policy ▸</a>
          <a href="/terms-of-use" style="color:var(--cream);text-decoration:none">Terms of Use ▸</a>
        </div>
      </div>
      {{range $st := .States}}
      <div class="card">
        <h3 style="font-family:var(--display);color:var(--orange);margin-bottom:.8rem">{{.Name}} ({{.Code}})</h3>
        <div style="display:flex;flex-direction:column;gap:.5rem">
          {{range $c := $st.Cities}}<a href="/areas/{{cityslug $c}}" style="color:var(--cream);text-decoration:none">{{.}} ▸</a>{{end}}
        </div>
      </div>
      {{end}}
    </div>
  </div>
</section>`
