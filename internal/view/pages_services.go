package view

// pageServices — templates/pages/services.php. DB-driven off the pest library.
const pageServices = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">CAPABILITIES // SERVICES</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Every pest. <em>One team.</em></h1>
    <p class="lead">From ants to wildlife, we identify, treat, and prevent {{len .Pests}}+ pest categories across Washington, Idaho &amp; Oregon, using eco-friendly, family- and pet-safe methods.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="grid g3">
      {{range $p := .Pests}}
      <a class="card" href="/pest/{{$p.Slug}}" style="text-decoration:none;color:inherit">
        <span class="pphoto"><img src="{{pestimg $p.Filename}}" alt="{{$p.Name}}" loading="lazy" decoding="async"><i class="ret" aria-hidden="true"></i></span>
        <h2 style="font-family:var(--display);color:var(--cream);margin-top:.7rem;font-size:1.17rem">{{$p.Name}} ▸</h2>
        <span class="badge" style="margin-top:.4rem;color:var(--khaki)">{{ucfirst $p.Category}}</span>
      </a>
      {{end}}
    </div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Don't see your pest? <em>Call us anyway.</em></h2>
    <p class="lead">If it invades, we handle it. Free quotes, same-day service available.</p>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="tel:+15094715767"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.79 19.79 0 0 1 2.12 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z"/></svg> (509) 471-5767</a>
      <a class="btn btn-ghost" href="/contact">Request Service ▸</a>
    </div>
  </div>
</section>
`

// pageAbout — templates/pages/about.php.
const pageAbout = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">PERSONNEL FILE // ABOUT US</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Veteran-owned. <em>Mission-driven.</em></h1>
    <p class="lead">Patriot Pest Control was founded by U.S. Military Veteran Skyler Rose to bring military discipline, integrity, and uncompromising excellence to pest control across three states.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="split" style="align-items:start">
      <div>
        <h2 style="font-family:var(--display);color:var(--cream)">Our <em>story.</em></h2>
        <p style="color:var(--khaki);line-height:1.8;margin-top:.8rem">After serving our country, our founder brought the same dedication and precision learned in the military to protecting American homes and businesses. What started as a one-operator mission has grown into a trusted, licensed, bonded, and insured team serving Washington, Idaho, and Oregon.</p>
        <p style="color:var(--khaki);line-height:1.8;margin-top:1rem">We're not just eliminating pests. We're protecting what matters most. Every treatment is backed by our 90-day warranty and our 100% satisfaction guarantee.</p>
      </div>
      <div class="dossier-file">
        <div class="form-id"><span>FORM PPC-14 · COMPANY</span><span>FILE 001</span></div>
        <dl>
          <div class="drow"><dt>Founded</dt><dd>BY VETERAN SKYLER ROSE</dd></div>
          <div class="drow"><dt>Theater</dt><dd>WA · ID · OR · AZ</dd></div>
          <div class="drow"><dt>Clearance</dt><dd>LICENSED · BONDED · INSURED</dd></div>
          <div class="drow"><dt>Methods</dt><dd>ECO-FRIENDLY · FAMILY &amp; PET SAFE</dd></div>
          <div class="drow"><dt>Guarantee</dt><dd>90-DAY WARRANTY · 100% SATISFACTION</dd></div>
          <div class="drow"><dt>Status</dt><dd>ACTIVE - SAME-DAY RESPONSE</dd></div>
        </dl>
      </div>
    </div>
  </div>
</section>

<section class="block">
  <div class="wrap">
    <div class="eyebrow">CORE VALUES</div>
    <h2 style="font-family:var(--display);color:var(--cream);margin:.4rem 0 1.4rem">What we <em>stand for.</em></h2>
    <div class="grid g3">
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="15" r="5"/><path d="M8.9 10.9 6.2 3h4L12 7.2 13.8 3h4l-2.7 7.9"/></svg> Integrity</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Honest assessments, transparent pricing, and no upselling. We treat your home like our own.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="6"/><circle cx="12" cy="12" r="2"/></svg> Precision</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Targeted treatments that eliminate pests at the source, not just the symptoms you see.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M11 20A7 7 0 0 1 9.8 6.1C15.5 5 17 4.48 19 2c1 2 2 4.18 2 8 0 5.5-4.78 10-10 10Z"/><path d="M2 21c0-3 1.85-5.36 5.08-6C9.5 14.52 12 13 13 12"/></svg> Safety</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Eco-friendly, low-toxicity products that are tough on pests and safe for kids and pets.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic fill" viewBox="0 0 24 24" aria-hidden="true"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10"/></svg> Responsiveness</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Same-day service when it can't wait, and online scheduling that's always open.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1 1 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"/><path d="m9 12 2 2 4-4"/></svg> Accountability</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">If pests return between visits, we re-treat free. No hassles, no excuses.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)"><svg width="1.2em" height="1.2em" viewBox="0 0 24 24" style="vertical-align:-.22em" aria-hidden="true"><rect x="2" y="4" width="20" height="16" rx="1.5" fill="#b31942"/><path d="M2 5.85h20M2 8.3h20M2 10.76h20M2 13.2h20M2 15.66h20M2 18.1h20" stroke="#fff" stroke-width="1.23"/><rect x="2" y="4" width="9" height="8.6" fill="#0a3161"/></svg> Service</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">We continue serving American families with the same dedication we showed in uniform.</p></div>
    </div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Join the families who <em>trust Patriot.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="tel:+15094715767"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.79 19.79 0 0 1 2.12 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z"/></svg> (509) 471-5767</a>
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>
`
