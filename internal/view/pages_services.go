package view

// pageServices — templates/pages/services.php. DB-driven off the pest library.
const pageServices = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">CAPABILITIES // SERVICES</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Every pest. <em>One team.</em></h1>
    <p class="lead">From ants to wildlife, we identify, treat, and prevent {{len .Pests}}+ pest categories across Washington, Idaho, Oregon &amp; Arizona, using eco-friendly, family- and pet-safe methods.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="grid g3">
      {{range $p := .Pests}}
      <a class="card" href="/pest/{{$p.Slug}}" style="text-decoration:none;color:inherit">
        <span class="pphoto"><img src="{{pestimg $p.Filename}}" alt="{{$p.Name}}" loading="lazy"><i class="ret" aria-hidden="true"></i></span>
        <h3 style="font-family:var(--display);color:var(--cream);margin-top:.7rem">{{$p.Name}} ▸</h3>
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
      <a class="btn btn-primary" href="tel:+15098180993">☎ (509) 818-0993</a>
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
    <p class="lead">Patriot Pest Control was founded by U.S. Military Veteran Skyler Rose to bring military discipline, integrity, and uncompromising excellence to pest control across four states.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="split" style="align-items:start">
      <div>
        <h2 style="font-family:var(--display);color:var(--cream)">Our <em>story.</em></h2>
        <p style="color:var(--khaki);line-height:1.8;margin-top:.8rem">After serving our country, our founder brought the same dedication and precision learned in the military to protecting American homes and businesses. What started as a one-operator mission has grown into a trusted, licensed, bonded, and insured team serving Washington, Idaho, Oregon, and Arizona.</p>
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
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">🎖️ Integrity</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Honest assessments, transparent pricing, and no upselling. We treat your home like our own.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">🎯 Precision</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Targeted treatments that eliminate pests at the source, not just the symptoms you see.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">🌿 Safety</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Eco-friendly, low-toxicity products that are tough on pests and safe for kids and pets.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">⚡ Responsiveness</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Same-day service when it can't wait, and a 24/7 line that's always open.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">🛡️ Accountability</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">If pests return between visits, we re-treat free. No hassles, no excuses.</p></div>
      <div class="card"><h3 style="font-family:var(--display);color:var(--cream)">🇺🇸 Service</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">We continue serving American families with the same dedication we showed in uniform.</p></div>
    </div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Join the families who <em>trust Patriot.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="tel:+15098180993">☎ (509) 818-0993</a>
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>
`
