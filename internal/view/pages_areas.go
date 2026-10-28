package view

// pageAreas — templates/pages/areas.php. Service-area overview.
// Data: .States ([]data.State{Code, Name, Cities}).
const pageAreas = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">AREA OF OPERATIONS</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Three states. <em>One call.</em></h1>
    <p class="lead">Same-day pest control across Washington, Idaho &amp; Oregon. Find your community below. If we're not listed, call us; we likely still cover you.</p>
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
      <a class="btn btn-primary" href="tel:+15094715767"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.79 19.79 0 0 1 2.12 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z"/></svg> (509) 471-5767</a>
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
      <a class="btn btn-primary" href="{{.AreaPhoneHref}}"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.79 19.79 0 0 1 2.12 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z"/></svg> Call {{.AreaPhoneDisplay}}</a>
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>

{{if .HasAreaCopy}}<section class="block alt">
  <div class="wrap">
    <div class="eyebrow">{{upper .CityStateCode}} LOCAL BRIEFING</div>
    <h2 style="font-family:var(--display);color:var(--cream);font-size:clamp(1.4rem,4vw,2rem);margin:.4rem 0 .8rem">About pest control in {{.CityName}}</h2>
    <p class="lead">{{.AreaCopy.Intro}}</p>
    <p style="color:var(--khaki);line-height:1.7">{{.AreaCopy.LocalNote}}</p>
    <p style="color:var(--khaki);line-height:1.7">{{.AreaCopy.PestNote}}</p>
  </div>
</section>{{end}}

<section class="block alt">
  <div class="wrap">
    <div class="grid g3">
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic fill" viewBox="0 0 24 24" aria-hidden="true"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10"/></svg> Same-Day Service</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Fast response across {{.CityName}} and surrounding communities when pests can't wait.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M11 20A7 7 0 0 1 9.8 6.1C15.5 5 17 4.48 19 2c1 2 2 4.18 2 8 0 5.5-4.78 10-10 10Z"/><path d="M2 21c0-3 1.85-5.36 5.08-6C9.5 14.52 12 13 13 12"/></svg> Family &amp; Pet Safe</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Low-toxicity, eco-friendly treatments that protect your household, not just eliminate pests.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1 1 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"/><path d="m9 12 2 2 4-4"/></svg> 90-Day Warranty</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">If pests come back between visits, we re-treat free. No hassles, no excuses.</p></div>
    </div>
    <div class="promise" style="margin-top:1.6rem">Serving {{.CityName}}, {{.CityStateName}} with ants, spiders, rodents, wasps, bed bugs, termites &amp; more. <a href="/services" style="color:var(--orange)">See every pest we treat ▸</a></div>
  </div>
</section>

{{if .AreaPests}}
<section class="block">
  <div class="wrap">
    <div class="eyebrow">POPULAR SERVICES // {{upper .CityStateCode}}</div>
    <h2 style="font-family:var(--display);color:var(--cream);font-size:clamp(1.4rem,4vw,2rem);margin:.4rem 0 1.2rem">Popular services <em>in {{.CityName}}.</em></h2>
    <div class="grid g3">
      {{range $p := .AreaPests}}
      <a class="card" href="/pest/{{$p.Slug}}" style="text-decoration:none;color:inherit">
        <h3 style="font-family:var(--display);color:var(--cream)">{{$p.Name}} Control ▸</h3>
        <p style="color:var(--khaki);font-size:.88rem;line-height:1.6;margin-top:.5rem">Same-day {{lower $p.Name}} treatment in {{$.CityName}} with a 90-day warranty.</p>
      </a>
      {{end}}
    </div>
    <div class="hero-ctas" style="margin-top:1.4rem">
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>
{{end}}

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Ready in {{.CityName}}? <em>Call now.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="{{.AreaPhoneHref}}"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.79 19.79 0 0 1 2.12 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z"/></svg> {{.AreaPhoneDisplay}}</a>
    </div>
  </div>
</section>`
