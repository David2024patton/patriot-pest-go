package view

// Pop UI variants of the inner pages (pop.patriotpest.pro). view.render
// swaps each page to its "-pop" variant on the pop host. Same rule as the
// pop home: illustrated tiles only on the surface, every fact behind a tap.
// The main site keeps the original templates untouched.

func init() {
	// ---------- Services: mini scene + the full pest rail + plans ----------
	RegisterPageTemplate("services-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="{{pestillus "ants.jpg"}}" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Capabilities // Services</span>
    <h1 class="pop-scene-h">Every pest. <span>One team.</span></h1>
    <p class="pop-scene-sub">The full catalog, grouped like a movie library. Tap any service for the briefing.</p>
  </div>
</section>
{{range $cat := .Categories}}
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">{{$cat.Name}}</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-cat-{{$cat.Name}}" aria-label="Scroll {{$cat.Name}} back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-cat-{{$cat.Name}}" aria-label="Scroll {{$cat.Name}} forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-cat-{{$cat.Name}}">{{template "part-pest-tiles" $cat.Pests}}</div>
  </div>
</section>
{{end}}
{{template "part-pest-sources" .Pests}}
{{template "part-plans-rail" .}}
{{template "part-cta" .}}`)

	// ---------- Prices: plans rail only (compare removed 2026-11-05) ----------
	RegisterPageTemplate("prices-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="{{pestillus "wasps.jpg"}}" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Supply Manifest // Pricing</span>
    <h1 class="pop-scene-h">Protect your <span>home.</span></h1>
    <p class="pop-scene-sub">Tap a plan for what it covers.</p>
  </div>
</section>
{{template "part-plans-rail" .}}
{{template "part-cta" .}}`)

	// ---------- Service areas index: state tiles in a grid ----------
	RegisterPageTemplate("areas-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="/assets/img/illus/state-wa.jpg" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Theater of Operations</span>
    <h1 class="pop-scene-h">Where we <span>operate.</span></h1>
    <p class="pop-scene-sub">Tap a state for its towns.</p>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Area of Operations</h2></div></div>
    <div class="pop-grid">{{template "part-state-tiles" .States}}</div>
  </div>
</section>
{{template "part-cta" .}}`)

	// ---------- FAQs: question tiles, answers in the popup ----------
	RegisterPageTemplate("faqs-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="{{pestillus "bees.jpg"}}" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Intel // FAQ</span>
    <h1 class="pop-scene-h">Questions, <span>answered.</span></h1>
    <p class="pop-scene-sub">Tap a question for the answer.</p>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Safety &amp; Methods</h2></div></div>
    <div class="pop-grid">
      <button type="button" class="pop-tile poster" data-qv="qv-faq-kids"><img class="ph" src="/assets/img/illus/faq-safety.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">Are your treatments safe for kids and pets?</b><small>Safety &amp; Methods</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-faq-eco"><img class="ph" src="/assets/img/illus/faq-safety.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">Are your products eco-friendly?</b><small>Safety &amp; Methods</small></span></button>
    </div>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Pricing &amp; Guarantees</h2></div></div>
    <div class="pop-grid">
      <button type="button" class="pop-tile poster" data-qv="qv-faq-cost"><img class="ph" src="/assets/img/illus/faq-price.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">How much does pest control cost?</b><small>Pricing &amp; Guarantees</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-faq-guarantee"><img class="ph" src="/assets/img/illus/faq-price.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">Do you offer a guarantee?</b><small>Pricing &amp; Guarantees</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-faq-fees"><img class="ph" src="/assets/img/illus/faq-price.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">Are there any hidden fees?</b><small>Pricing &amp; Guarantees</small></span></button>
    </div>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Service &amp; Scheduling</h2></div></div>
    <div class="pop-grid">
      <button type="button" class="pop-tile poster" data-qv="qv-faq-sameday"><img class="ph" src="/assets/img/illus/faq-sched.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">Do you offer same-day service?</b><small>Service &amp; Scheduling</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-faq-areas"><img class="ph" src="/assets/img/illus/faq-sched.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">What areas do you serve?</b><small>Service &amp; Scheduling</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-faq-prep"><img class="ph" src="/assets/img/illus/faq-sched.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">Do I need to prepare my home before treatment?</b><small>Service &amp; Scheduling</small></span></button>
    </div>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Getting Started</h2></div></div>
    <div class="pop-grid">
      <button type="button" class="pop-tile poster" data-qv="qv-faq-start"><img class="ph" src="/assets/img/illus/faq-start.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">How do I get started?</b><small>Getting Started</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-faq-commercial"><img class="ph" src="/assets/img/illus/faq-start.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b class="clamp">Do you handle commercial properties?</b><small>Getting Started</small></span></button>
    </div>
  </div>
  <div class="pop-qv-sources" hidden>
    <div id="qv-faq-kids"><div class="pop-lab">Safety &amp; Methods</div><h3>Are your treatments safe for kids and pets?</h3><p>Yes. We use low-toxicity products chosen to be tough on pests and easy on households, and we apply them with precision in cracks, crevices, and other spots kids and pets cannot reach, rather than broadcasting them across open surfaces. Baits go into tamper-resistant stations placed out of reach. Your technician will tell you about any brief re-entry wait for the specific treatment used and what to do if your pet shows interest in a treated area. If anyone in the home has allergies, asthma, or chemical sensitivities, mention it when you book so we can plan the treatment around it.</p></div>
    <div id="qv-faq-eco"><div class="pop-lab">Safety &amp; Methods</div><h3>Are your products eco-friendly?</h3><p>Absolutely. We work from an integrated approach: inspection, exclusion, and sanitation come first, and products are used only where they are needed, in the smallest effective amount. Applications are targeted at pest harborages and entry points instead of blanket spraying the whole property, which keeps pollinators, soil, and waterways out of the picture. It still gets the job done, because the product goes where the pests actually live. Ask your technician what is being applied and why, and you will get a straight answer.</p></div>
    <div id="qv-faq-cost"><div class="pop-lab">Pricing &amp; Guarantees</div><h3>How much does pest control cost?</h3><p>It depends on the pest, how far the infestation has spread, the size of your home, and whether you need a one-time knockdown or ongoing protection. A single wasp nest and a whole-house rodent exclusion are very different jobs, so we do not do one-size-fits-all pricing. Every quote is free, itemized, and given before any work starts, so you see exactly what is included. Call or request a quote online and you will get a firm number, not a range that moves later.</p></div>
    <div id="qv-faq-guarantee"><div class="pop-lab">Pricing &amp; Guarantees</div><h3>Do you offer a guarantee?</h3><p>Yes. Every treatment is backed by our 90-day warranty and 100% satisfaction guarantee. If pests return between scheduled visits, we re-treat at no additional cost.</p></div>
    <div id="qv-faq-fees"><div class="pop-lab">Pricing &amp; Guarantees</div><h3>Are there any hidden fees?</h3><p>Never. The price we quote is the price you pay, and we confirm it before work starts so there is nothing to argue about later. Quotes are free, and re-treatments between scheduled visits are covered, not billed as extras. No fuel surcharges, no surprise line items on the invoice. If a job ever needs something outside the original quote, we tell you first and get your OK before we do it.</p></div>
    <div id="qv-faq-sameday"><div class="pop-lab">Service &amp; Scheduling</div><h3>Do you offer same-day service?</h3><p>Yes. When pests cannot wait, like a wasp nest by the front door or mice in the pantry, we offer same-day visits across Washington, Idaho, and Oregon. Scheduling depends on technician availability and routing, so call as early in the day as you can to get a slot. If the day is fully booked, we will set the next available visit and tell you what to do in the meantime.</p></div>
    <div id="qv-faq-areas"><div class="pop-lab">Service &amp; Scheduling</div><h3>What areas do you serve?</h3><p>We cover three states. In Washington: Spokane, Spokane Valley, Cheney, Liberty Lake, Airway Heights, Medical Lake, Deer Park, and Mead. In Idaho: Coeur d'Alene, Post Falls, Hayden, and Rathdrum. In Oregon: Hermiston and Milton-Freewater. See our service areas page for the full list, and if you are just outside these areas, call and ask.</p></div>
    <div id="qv-faq-prep"><div class="pop-lab">Service &amp; Scheduling</div><h3>Do I need to prepare my home before treatment?</h3><p>Usually very little. When you schedule, we give you prep instructions matched to your pest and treatment, for example clearing under-sink cabinets before an ant treatment or trimming grass back from the foundation before an exterior visit. Most interior work just needs clear access along baseboards and to the problem areas. Your technician walks you through anything left on arrival, so you are never guessing.</p></div>
    <div id="qv-faq-start"><div class="pop-lab">Getting Started</div><h3>How do I get started?</h3><p>Call us at {{.PhoneDisplay}} or request a free quote through the contact page. We will ask what you are seeing, where, and how long it has been going on, then recommend a plan and a price before any work happens. Many jobs can be scheduled the same day you call. There is no obligation and no pressure; if all you need is advice, we will say so.</p></div>
    <div id="qv-faq-commercial"><div class="pop-lab">Getting Started</div><h3>Do you handle commercial properties?</h3><p>Yes. We treat restaurants, offices, retail spaces, warehouses, and multi-unit housing as well as homes. Visits are scheduled around your operating hours so customers and staff are not disrupted. We also provide the service records that health inspections and audits ask for, so compliance stays simple.</p></div>
  </div>
</section>
{{template "part-cta" .}}`)
}

func init() {
	// ---------- About: story, dossier, and values as tiles ----------
	RegisterPageTemplate("about-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="/assets/img/illus/state-wa.jpg" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Personnel File // About Us</span>
    <h1 class="pop-scene-h">Veteran-owned. <span>Mission-driven.</span></h1>
    <p class="pop-scene-sub">Tap a file for the details.</p>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">The Company</h2></div></div>
    <div class="pop-grid">
      <button type="button" class="pop-tile poster" data-qv="qv-about-story">
        <img class="ph" src="/assets/img/illus/story.jpg" alt="Our story" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Our Story</b><small>How Patriot started</small></span>
      </button>
      <button type="button" class="pop-tile poster" data-qv="qv-about-file">
        <img class="ph" src="/assets/img/illus/dossier.jpg" alt="Company file" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Company File</b><small>Form PPC-14 &middot; File 001</small></span>
      </button>
    </div>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Core Values</h2></div></div>
    <div class="pop-grid">
      <button type="button" class="pop-tile poster" data-qv="qv-val-integrity"><img class="ph" src="/assets/img/illus/val-medal.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Integrity</b><small>Core Values</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-val-precision"><img class="ph" src="/assets/img/illus/val-medal.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Precision</b><small>Core Values</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-val-safety"><img class="ph" src="/assets/img/illus/val-medal.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Safety</b><small>Core Values</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-val-response"><img class="ph" src="/assets/img/illus/val-medal.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Responsiveness</b><small>Core Values</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-val-account"><img class="ph" src="/assets/img/illus/val-medal.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Accountability</b><small>Core Values</small></span></button>
      <button type="button" class="pop-tile poster" data-qv="qv-val-service"><img class="ph" src="/assets/img/illus/val-medal.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Service</b><small>Core Values</small></span></button>
    </div>
  </div>
  <div class="pop-qv-sources" hidden>
    <div id="qv-about-story">
      <div class="pop-qv-art"><img src="/assets/img/illus/story.jpg" alt=""></div>
      <div class="pop-lab">About &middot; Our Story</div>
      <h3>Our story</h3>
      <p>After serving our country, our founder brought the same dedication and precision learned in the military to protecting American homes and businesses. What started as a one-operator mission has grown into a trusted, licensed, bonded, and insured team serving Washington, Idaho, and Oregon.</p>
      <div class="qv-more" hidden><p>We're not just eliminating pests. We're protecting what matters most. Every treatment is backed by our 90-day warranty and our 100% satisfaction guarantee.</p></div>
    </div>
    <div id="qv-about-file">
      <div class="pop-qv-art"><img src="/assets/img/illus/dossier.jpg" alt=""></div>
      <div class="pop-lab">Form PPC-14 &middot; Company</div>
      <h3>Company file</h3>
      <dl class="pop-dl">
        <div><dt>Founded</dt><dd>BY VETERAN SKYLER ROSE</dd></div>
        <div><dt>Theater</dt><dd>WA &middot; ID &middot; OR</dd></div>
        <div><dt>Clearance</dt><dd>LICENSED &middot; BONDED &middot; INSURED</dd></div>
        <div><dt>Methods</dt><dd>ECO-FRIENDLY &middot; FAMILY &amp; PET SAFE</dd></div>
        <div><dt>Guarantee</dt><dd>90-DAY WARRANTY &middot; 100% SATISFACTION</dd></div>
        <div><dt>Status</dt><dd>ACTIVE - SAME-DAY RESPONSE</dd></div>
      </dl>
    </div>
    <div id="qv-val-integrity"><div class="pop-lab">Core Values</div><h3>Integrity</h3><p>Honest assessments, transparent pricing, and no upselling. We treat your home like our own.</p></div>
    <div id="qv-val-precision"><div class="pop-lab">Core Values</div><h3>Precision</h3><p>Targeted treatments that eliminate pests at the source, not just the symptoms you see.</p></div>
    <div id="qv-val-safety"><div class="pop-lab">Core Values</div><h3>Safety</h3><p>Eco-friendly, low-toxicity products that are tough on pests and safe for kids and pets.</p></div>
    <div id="qv-val-response"><div class="pop-lab">Core Values</div><h3>Responsiveness</h3><p>Same-day service when it can't wait, and online scheduling that's always open.</p></div>
    <div id="qv-val-account"><div class="pop-lab">Core Values</div><h3>Accountability</h3><p>If pests return between visits, we re-treat free. No hassles, no excuses.</p></div>
    <div id="qv-val-service"><div class="pop-lab">Core Values</div><h3>Service</h3><p>We continue serving American families with the same dedication we showed in uniform.</p></div>
  </div>
</section>
{{template "part-cta" .}}`)

	// ---------- Contact: action tiles; the form is the quote popup ----------
	RegisterPageTemplate("contact-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="{{pestillus "mosquitoes.jpg"}}" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Comms // Contact</span>
    <h1 class="pop-scene-h">Reach the <span>team.</span></h1>
    <p class="pop-scene-sub">Pick a channel. We respond within one business day, usually much faster.</p>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    {{if .Success}}<div class="pop-notice"><b>Message sent.</b><br>{{.Success}}</div>{{end}}
    {{if .Errors}}<div class="pop-notice error"><b>Please fix the following:</b><ul>{{range .Errors}}<li>{{.}}</li>{{end}}</ul></div>{{end}}
    <div class="pop-rail-head"><div><h2 class="pop-h2">Channels</h2></div></div>
    <div class="pop-grid">
      <a class="pop-tile poster" href="{{.PhoneHref}}">
        <img class="ph" src="/assets/img/illus/ch-call.jpg" alt="" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Call or Text {{.PhoneDisplay}}</b><small>Fastest &middot; Tap to call</small></span>
      </a>
      <button type="button" class="pop-tile poster" data-pop-open="quote">
        <img class="ph" src="/assets/img/illus/ch-quote.jpg" alt="" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Get a Free Quote</b><small>Free &middot; No obligation</small></span>
      </button>
      <a class="pop-tile poster" href="/service-areas">
        <img class="ph" src="/assets/img/illus/ch-map.jpg" alt="" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Service Areas</b><small>Coverage</small></span>
      </a>
      <a class="pop-tile poster" href="/prices">
        <img class="ph" src="/assets/img/illus/plan-compare.jpg" alt="" loading="lazy" decoding="async">
        <span class="pop-cap"><b>See the Plans</b><small>Pricing</small></span>
      </a>
      <button type="button" class="pop-tile poster" data-qv="qv-office">
        <img class="ph" src="/assets/img/illus/dossier.jpg" alt="" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Office &amp; Email</b><small>Headquarters</small></span>
      </button>
      <a class="pop-tile poster" href="/signup">
        <img class="ph" src="/assets/img/illus/ch-account.jpg" alt="" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Create an Account</b><small>Accounts</small></span>
      </a>
    </div>
    <div class="pop-qv-sources" hidden>
      <div id="qv-office">
        <div class="pop-qv-art"><img src="/assets/img/illus/dossier.jpg" alt=""></div>
        <div class="pop-lab">Headquarters</div>
        <h3>Office &amp; Email</h3>
        <dl class="pop-dl">
          <div><dt>Phone</dt><dd>{{.PhoneDisplay}} - WA, ID, OR</dd></div>
          <div><dt>Email</dt><dd>INFO@PATRIOTPEST.PRO</dd></div>
          <div><dt>Office</dt><dd>DEER PARK, WA 99006</dd></div>
        </dl>
      </div>
    </div>
  </div>
</section>
{{template "part-plans-rail" .}}`)

	// ---------- Blog index: every post as a poster tile ----------
	RegisterPageTemplate("blog-index-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="{{pestillus "fruit-flies.jpg"}}" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">From the Field // Blog</span>
    <h1 class="pop-scene-h">Field <span>guides.</span></h1>
    <p class="pop-scene-sub">Tap a guide for the briefing.</p>
  </div>
</section>
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">All Field Guides</h2></div></div>
    <div class="pop-grid">{{template "part-blog-tiles" .Posts}}</div>
  </div>
</section>
{{template "part-blog-sources" .Posts}}
{{template "part-cta" .}}`)

	// ---------- Pest file: scene header, section tiles, related rail ----------
	RegisterPageTemplate("pest-pop", `
<section class="pop-scene"><img class="pop-scene-art" src="{{pestillus .Pest.Filename}}" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Service // {{upper .Pest.Category}}</span>
    <h1 class="pop-scene-h">{{.Pest.Name}} <span>Control</span></h1>
    {{if .Pest.ScientificName}}<p class="pop-scene-sub"><i>{{.Pest.ScientificName}}</i> &middot; Regional threat level {{.Pest.ThreatLevel}}%</p>{{end}}
    <div class="pop-cta-row">
      <a class="pop-btn primary" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
      <button type="button" class="pop-btn" data-pop-open="quote">Free Quote</button>
    </div>
  </div>
</section>
{{$c := pestcopy2 .Pest.Slug .Pest.Category}}
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">The Briefing</h2></div></div>
    <div class="pop-grid">
      <button type="button" class="pop-tile poster" data-qv="qv-pest-brief"><img class="ph" src="{{pestillus .Pest.Filename}}" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>What you're up against</b><small>Services</small></span></button>
      {{if $c.Signs}}<button type="button" class="pop-tile poster" data-qv="qv-pest-signs"><img class="ph" src="/assets/img/illus/pb-signs.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Signs of Activity</b><small>Field Signs</small></span></button>{{end}}
      {{if $c.Treat}}<button type="button" class="pop-tile poster" data-qv="qv-pest-treat"><img class="ph" src="/assets/img/illus/pb-treat.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>How we treat it</b><small>Our Treatment</small></span></button>{{end}}
      {{if $c.Prev}}<button type="button" class="pop-tile poster" data-qv="qv-pest-prev"><img class="ph" src="/assets/img/illus/pb-prev.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>Keep them out</b><small>Prevention</small></span></button>{{end}}
      <button type="button" class="pop-tile poster" data-qv="qv-pest-warranty"><img class="ph" src="/assets/img/illus/pb-warranty.jpg" alt="" loading="lazy" decoding="async"><span class="pop-cap"><b>90-Day Warranty</b><small>Coverage</small></span></button>
    </div>
  </div>
  <div class="pop-qv-sources" hidden>
    <div id="qv-pest-brief">
      <div class="pop-qv-art"><img src="{{pestillus .Pest.Filename}}" alt=""></div>
      <div class="pop-lab">{{upper .Pest.Category}} &middot; Threat {{.Pest.ThreatLevel}}%</div>
      <h3>{{.Pest.Name}}</h3>
      {{if .Pest.ScientificName}}<div class="pop-sci">{{.Pest.ScientificName}}</div>{{end}}
      <p>{{.Pest.Description}}</p>
    </div>
    {{if $c.Signs}}<div id="qv-pest-signs"><div class="pop-lab">Field Signs</div><h3>Signs of activity</h3><ul>{{range $c.Signs}}<li>{{.}}</li>{{end}}</ul></div>{{end}}
    {{if $c.Treat}}<div id="qv-pest-treat"><div class="pop-lab">Our Treatment</div><h3>How we treat it</h3><p>{{$c.Treat}}</p></div>{{end}}
    {{if $c.Prev}}<div id="qv-pest-prev"><div class="pop-lab">Prevention</div><h3>Keep them out</h3><ul>{{range $c.Prev}}<li>{{.}}</li>{{end}}</ul></div>{{end}}
    <div id="qv-pest-warranty"><div class="pop-lab">Coverage</div><h3>90-Day Warranty</h3><p>If {{lower .Pest.Name}} return between scheduled visits, we re-treat at no additional cost. Licensed, bonded, and insured across Washington, Idaho &amp; Oregon.</p></div>
  </div>
</section>
{{if .Related}}
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Also in Your Region</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-related" aria-label="Scroll related back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-related" aria-label="Scroll related forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-related">{{template "part-pest-tiles" .Related}}</div>
  </div>
</section>
{{template "part-pest-sources" .Related}}
{{end}}
{{template "part-plans-rail" .}}
{{template "part-cta" .}}`)

	// ---------- Town page: scene, local briefing tile, local pests ----------
	RegisterPageTemplate("area-detail-pop", `
<section class="pop-scene slim"><img class="pop-scene-art" src="/assets/img/illus/state-wa.jpg" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Service Area // {{.CityStateName}}</span>
    <h1 class="pop-scene-h">Pest control in <span>{{.CityName}}, {{.CityStateCode}}</span></h1>
    <p class="pop-scene-sub">Same-day service, 90-day warranty, veteran-owned.</p>
    <div class="pop-cta-row">
      <a class="pop-btn primary" href="{{.AreaPhoneHref}}">Call or Text {{.AreaPhoneDisplay}}</a>
      <button type="button" class="pop-btn" data-pop-open="quote">Free Quote</button>
    </div>
  </div>
</section>
{{if .HasAreaCopy}}
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Local Briefing</h2></div></div>
    <button type="button" class="pop-tile poster wide" data-qv="qv-town-brief">
      <img class="ph" src="{{if eq .CityStateCode "WA"}}/assets/img/illus/state-wa.jpg{{else if eq .CityStateCode "ID"}}/assets/img/illus/state-id.jpg{{else}}/assets/img/illus/state-or.jpg{{end}}" alt="{{.CityName}} briefing" loading="lazy" decoding="async">
      <span class="pop-cap"><b>{{.CityName}} Briefing</b><small>Local conditions &amp; pest pressure</small></span>
    </button>
  </div>
  <div class="pop-qv-sources" hidden>
    <div id="qv-town-brief">
      <div class="pop-qv-art"><img src="{{if eq .CityStateCode "WA"}}/assets/img/illus/state-wa.jpg{{else if eq .CityStateCode "ID"}}/assets/img/illus/state-id.jpg{{else}}/assets/img/illus/state-or.jpg{{end}}" alt=""></div>
      <div class="pop-lab">Local Briefing &middot; {{.CityStateCode}}</div>
      <h3>{{.CityName}}, {{.CityStateName}}</h3>
      <p>{{.AreaCopy.Intro}}</p>
      <div class="qv-more" hidden>
        {{if .AreaCopy.LocalNote}}<div class="pop-lab">On the ground</div><p>{{.AreaCopy.LocalNote}}</p>{{end}}
        {{if .AreaCopy.PestNote}}<div class="pop-lab">Pest pressure</div><p>{{.AreaCopy.PestNote}}</p>{{end}}
      </div>
    </div>
  </div>
</section>
{{end}}
{{if .AreaPests}}
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Top Threats in {{.CityName}}</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-areapests" aria-label="Scroll pests back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-areapests" aria-label="Scroll pests forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-areapests">{{template "part-pest-tiles" .AreaPests}}</div>
  </div>
</section>
{{template "part-pest-sources" .AreaPests}}
{{end}}
{{template "part-plans-rail" .}}
{{template "part-cta" .}}`)
}

func init() {
	// ---------- Sign up: the account form lives in its popup ----------
	signupPop := `
<section class="pop-scene slim"><img class="pop-scene-art" src="{{pestillus "silverfish.jpg"}}" alt="">
  <div class="pop-scene-in">
    <span class="pop-tagpill">Account // Sign Up</span>
    <h1 class="pop-scene-h">Create your <span>account.</span></h1>
    <p class="pop-scene-sub">Request service, track visits, and get member pricing. No card required.</p>
  </div>
</section>
{{if .Success}}
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-notice"><b>You're on the list.</b><br>{{.Success}}</div>
  </div>
</section>
{{if eq .AnalyticsEvent "generate_lead"}}<script>if (window.gtag) { gtag('event', 'generate_lead', { 'event_category': 'conversion', 'event_label': 'Account Signup' }); } if (window.fbq) { fbq('track', 'Lead', { 'content_name': 'Account Signup' }); }</script>{{end}}
{{else}}
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Get Started</h2></div></div>
    <button type="button" class="pop-tile poster wide" data-pop-open="signup">
      <img class="ph" src="/assets/img/illus/signup-art.jpg" alt="Create your account" loading="lazy" decoding="async">
      <span class="pop-cap"><b>Create Your Account</b><small>Tap to open the sign-up form</small></span>
    </button>
  </div>
</section>
<div class="pop-scrim" id="pop-signup" {{if .Errors}}data-autopen{{end}} hidden>
  <div class="pop-popup" role="dialog" aria-modal="true" aria-label="Sign up">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
    <div class="pop-lab">Sign Up</div>
    <h3>Create your account</h3>
    {{if .Errors}}<div class="pop-notice error"><b>Please fix the following:</b><ul>{{range .Errors}}<li>{{.}}</li>{{end}}</ul></div>{{end}}
    <form method="post" action="/signup" class="form-stack" novalidate>
      {{raw .Csrf}}
      <div class="field"><label for="s-name">Full name</label><input type="text" id="s-name" name="name" required autocomplete="name" value="{{.OldName}}"></div>
      <div class="field"><label for="s-email">Email</label><input type="email" id="s-email" name="email" required autocomplete="email" value="{{.OldEmail}}"></div>
      <div class="field"><label for="s-phone">Phone (optional)</label><input type="tel" id="s-phone" name="phone" autocomplete="tel" value="{{.OldPhone}}"></div>
      <div class="form-row">
        <div class="field"><label for="s-city">City</label><input type="text" id="s-city" name="city" autocomplete="address-level2" value="{{.OldCity}}"></div>
        <div class="field"><label for="s-state">State</label><input type="text" id="s-state" name="state" maxlength="2" placeholder="WA" autocomplete="address-level1" value="{{.OldState}}"></div>
      </div>
      <div class="field"><label for="s-zip">ZIP (optional)</label><input type="text" id="s-zip" name="zip" maxlength="10" autocomplete="postal-code" value="{{.OldZip}}"></div>
      ` + smsConsent + `
      <div class="form-actions"><button type="submit" class="pop-btn primary">Sign Up</button></div>
      <p class="pop-form-note">Online accounts are coming soon. For now, call or message us and a real human will take care of you.</p>
    </form>
  </div>
</div>
{{end}}
{{template "part-plans-rail" .}}
{{template "part-cta" .}}`
	RegisterPageTemplate("signup-pop", signupPop)
}
