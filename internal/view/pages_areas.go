package view

// pageAreas — templates/pages/areas.php. Service-area overview.
// Data: .States ([]data.State{Code, Name, Cities}).
const pageAreas = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">AREA OF OPERATIONS</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Four states. <em>One call.</em></h1>
    <p class="lead">Same-day pest control across Washington, Idaho, Oregon &amp; Arizona. Find your community below. If we're not listed, call us; we likely still cover you.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="area-grid" style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:1.2rem">
      {{range $st := .States}}
      <div class="area-col card">
        <h3 style="font-family:var(--display);color:var(--orange);font-size:1.6rem">{{upper $st.Code}}</h3>
        <div style="font-family:var(--mono);font-size:.72rem;letter-spacing:.1em;color:var(--khaki);text-transform:uppercase;margin-bottom:.8rem">{{.Name}}</div>
        <div class="area-cities" style="display:flex;flex-direction:column;gap:.4rem">
          {{range $c := $st.Cities}}<a href="/areas/{{cityslug $c}}" style="color:var(--cream);text-decoration:none">{{.}} ▸</a>
          {{end}}
        </div>
        <span class="area-note" style="display:block;margin-top:.9rem;font-family:var(--mono);font-size:.7rem;color:var(--olive-300)">● {{len $st.Cities}} ZONE{{if ne (len $st.Cities) 1}}S{{end}} · SAME-DAY AVAILABLE</span>
      </div>
      {{end}}
    </div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">In the region? <em>Let's talk.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="tel:+15098180993">☎ (509) 818-0993</a>
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>`

// pageAreaDetail — templates/pages/area-detail.php. One city landing page.
// Data: .CityName, .CityStateCode, .CityStateName, plus the localized line
// precomputed by the handler (.AreaPhoneDisplay / .AreaPhoneHref) so the
// template stays branch-free.
const pageAreaDetail = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">AREA OF OPERATIONS // {{upper .CityStateCode}}</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(1.8rem,5vw,2.8rem);margin:.4rem 0 .8rem">Pest Control in {{.CityName}}, {{.CityStateName}}</h1>
    <p class="lead">Same-day, eco-friendly pest control for {{.CityName}} homes and businesses. Veteran-owned, licensed and insured, backed by a 90-day warranty.</p>
    <div class="hero-ctas" style="margin-top:1.4rem">
      <a class="btn btn-primary" href="{{.AreaPhoneHref}}">☎ Call {{.AreaPhoneDisplay}}</a>
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="grid g3">
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">⚡ Same-Day Service</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Fast response across {{.CityName}} and surrounding communities when pests can't wait.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">🌿 Family &amp; Pet Safe</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Low-toxicity, eco-friendly treatments that protect your household, not just eliminate pests.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">🛡️ 90-Day Warranty</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">If pests come back between visits, we re-treat free. No hassles, no excuses.</p></div>
    </div>
    <div class="promise" style="margin-top:1.6rem">Serving {{.CityName}}, {{.CityStateName}} with ants, spiders, rodents, wasps, bed bugs, termites &amp; more. <a href="/services" style="color:var(--orange)">See every pest we treat ▸</a></div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Ready in {{.CityName}}? <em>Call now.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="{{.AreaPhoneHref}}">☎ {{.AreaPhoneDisplay}}</a>
    </div>
  </div>
</section>`
