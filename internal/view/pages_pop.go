package view

// pageHomePop is the Pop UI homepage, served only on pop.patriotpest.pro.
// Same real data as the main homepage (DB pests, latest posts, states) plus
// the prices-page plan facts. Everything tappable opens a centered quick
// view popup first; More info goes to the real page. The quote form is the
// real /contact form (same fields, same CSRF), living in a popup per Pop UI.

func init() {
	RegisterPageTemplate("home-pop", pageHomePop)
}

const pageHomePop = `<!-- ============ POP HERO ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-scene">
      <img class="pop-scene-bg" src="{{pestimg "rodents.jpg"}}" alt="">
      <div class="pop-scene-shade"></div>
      <div class="pop-scene-in">
        <div class="pop-tagpill">A PROUD VETERAN-OWNED AMERICAN COMPANY</div>
        <h1 class="pop-scene-h">Hostile pests?<br><span>Neutralized.</span></h1>
        <p class="pop-scene-sub">Military-precision pest control for homes and businesses across <b>Washington, Idaho &amp; Oregon</b>. Eco-friendly treatments safe for your family and pets, with <b>same-day service</b> when it can't wait.</p>
        <div class="pop-cta-row">
          <button type="button" class="pop-btn primary" data-pop-open="quote">Get a Free Quote</button>
          <a class="pop-btn outline" href="/prices">See the Plans</a>
          <a class="pop-btn light" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
        </div>
      </div>
    </div>
    <div class="pop-stats">
      <div class="pop-stat"><div class="pop-stat-n">4.8</div><div class="pop-stat-l">Google Rating</div></div>
      <div class="pop-stat"><div class="pop-stat-n">110</div><div class="pop-stat-l">Google Reviews</div></div>
      <div class="pop-stat"><div class="pop-stat-n">90-Day</div><div class="pop-stat-l">Warranty</div></div>
      <div class="pop-stat"><div class="pop-stat-n">{{.PestCount}}</div><div class="pop-stat-l">Pest Files</div></div>
    </div>
  </div>
</section>

<!-- ============ PLANS RAIL (pickers stay on the page) ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Our Plans</h2><p class="pop-sub">Transparent pricing. No hidden fees. Tap a plan for the quick view.</p></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-plans" aria-label="Scroll plans back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-plans" aria-label="Scroll plans forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-plans">
      <button type="button" class="pop-tile planp tone-bronze" data-qv="qv-plan-bronze">
        <span class="pop-plan-name">Bronze</span>
        <span class="pop-plan-price">$49<small> /MO</small></span>
        <span class="pop-plan-line">Exterior-Only Protection</span>
      </button>
      <button type="button" class="pop-tile planp tone-silver" data-qv="qv-plan-silver">
        <span class="pop-plan-name">Silver</span>
        <span class="pop-plan-price">$59<small> /MO</small></span>
        <span class="pop-plan-line">Interior + Exterior</span>
      </button>
      <button type="button" class="pop-tile planp tone-gold" data-qv="qv-plan-gold">
        <span class="pop-plan-tag">Most Popular</span>
        <span class="pop-plan-name">Gold</span>
        <span class="pop-plan-price">$79<small> /MO</small></span>
        <span class="pop-plan-line">Priority Interior &amp; Exterior</span>
      </button>
      <a class="pop-tile planp tone-more" href="/prices">
        <span class="pop-plan-name">All Plans</span>
        <span class="pop-plan-price">Compare</span>
        <span class="pop-plan-line">Full pricing, side by side</span>
      </a>
    </div>
  </div>
</section>

<!-- ============ THREAT FILES RAIL (DB-driven, every pest) ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Threat Files</h2><p class="pop-sub">{{.PestCount}} hostile categories operate in your region right now. Tap one for the quick view.</p></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-pests" aria-label="Scroll pests back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-pests" aria-label="Scroll pests forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-pests">
      {{range $p := .Pests}}
      <button type="button" class="pop-tile poster" data-qv="qv-pest-{{$p.Slug}}">
        <img class="ph" src="{{pestimg $p.Filename}}" alt="{{$p.Name}}" loading="lazy" decoding="async">
        <span class="pop-cap"><b>{{$p.Name}}</b><small>{{if $p.ScientificName}}{{$p.ScientificName}} · {{end}}Threat {{$p.ThreatLevel}}%</small></span>
      </button>
      {{end}}
    </div>
  </div>
</section>

<!-- ============ OPERATOR ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-split">
      <div class="pop-card">
        <div class="pop-lab">Personnel File 001</div>
        <dl class="pop-dl">
          <div><dt>Name</dt><dd>SKYLER ROSE</dd></div>
          <div><dt>Role</dt><dd>FOUNDER &amp; OPERATOR</dd></div>
          <div><dt>Service</dt><dd>U.S. MILITARY VETERAN</dd></div>
          <div><dt>Theater</dt><dd>WA · ID · OR</dd></div>
          <div><dt>Clearance</dt><dd>LICENSED · BONDED · INSURED</dd></div>
          <div><dt>Status</dt><dd>ACTIVE - SAME-DAY RESPONSE</dd></div>
        </dl>
      </div>
      <div class="pop-card">
        <h2 class="pop-h2 dark">The operator behind the operation.</h2>
        <p class="pop-lead dark">Patriot Pest Control was founded by <b>U.S. Military Veteran Skyler Rose</b>, bringing military discipline, integrity, and uncompromising excellence to pest control across three states. Over a decade of field experience. Thousands of homes and businesses protected.</p>
        <blockquote class="pop-quote">
          <p>After serving our country, I founded Patriot Pest Control to continue serving American families and businesses with the same dedication, precision, and integrity I learned in the military. We're not just eliminating pests. We're protecting what matters most.</p>
          <div class="pop-quote-foot"><b>Skyler Rose</b><small>FOUNDER &amp; VETERAN · PATRIOT PEST CONTROL</small></div>
        </blockquote>
      </div>
    </div>
  </div>
</section>

<!-- ============ AREAS ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Area of Operations</h2><p class="pop-sub">Three states. One call. Tap a town for its page.</p></div></div>
    <div class="pop-area-grid">
      {{range $st := .States}}
      <div class="pop-card area">
        <h3>{{$st.Code}} · {{$st.Name}}</h3>
        <div class="pop-cities">{{range $c := $st.Cities}}<a class="pop-city" href="/areas/{{cityslug $c}}">{{$c}}</a>{{end}}</div>
        <span class="pop-area-note">{{len $st.Cities}} zones · Same-day available</span>
      </div>
      {{end}}
    </div>
  </div>
</section>

<!-- ============ GUARANTEE ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">The Guarantee</h2><p class="pop-sub">We stand behind every mission.</p></div></div>
    <div class="pop-stats five">
      <div class="pop-stat"><div class="pop-stat-n">100%</div><div class="pop-stat-l">Satisfaction Guaranteed</div></div>
      <div class="pop-stat"><div class="pop-stat-n">90-Day</div><div class="pop-stat-l">Warranty on Treatments</div></div>
      <div class="pop-stat"><div class="pop-stat-n">24/7</div><div class="pop-stat-l">Appointment Scheduling</div></div>
      <div class="pop-stat"><div class="pop-stat-n">3</div><div class="pop-stat-l">States: WA · ID · OR</div></div>
      <div class="pop-stat"><div class="pop-stat-n">10+</div><div class="pop-stat-l">Years of Experience</div></div>
    </div>
    <div class="pop-card promise">
      <p><b>Our promise:</b> if pests return between scheduled visits, we re-treat at no additional cost. No hassles, no excuses. Licensed, bonded, and insured, with eco-friendly, low-toxicity products safe for kids, pets, and the environment.</p>
    </div>
  </div>
</section>

<!-- ============ FIELD REPORTS ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Field Reports</h2><p class="pop-sub">Debriefs from the front line.</p></div></div>
    <div class="pop-grid3">
      <div class="pop-card review"><div class="pop-stars">★★★★★</div><p>"Patriot Pest Control saved our home from a serious ant infestation. Professional, thorough, and the results were immediate."</p><div class="pop-who"><b>Sarah M.</b><small>WASHINGTON</small></div></div>
      <div class="pop-card review"><div class="pop-stars">★★★★★</div><p>"We've used Patriot for our restaurant for 2 years. Reliable, discreet, always on time. Our health inspections have never been better."</p><div class="pop-who"><b>James T.</b><small>LIBERTY LAKE, WA</small></div></div>
      <div class="pop-card review"><div class="pop-stars">★★★★★</div><p>"Fast response and eco-friendly products safe for my kids and pets. The technician was knowledgeable and explained everything."</p><div class="pop-who"><b>Lisa R.</b><small>COEUR D'ALENE, ID</small></div></div>
    </div>
  </div>
</section>

<!-- ============ FROM THE FIELD (blog rail) ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">From the Field</h2><p class="pop-sub">The latest Patriot Pest Control field guides.</p></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-blog" aria-label="Scroll articles back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-blog" aria-label="Scroll articles forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-blog">
      {{range $post := .Posts}}
      <a class="pop-tile poster" href="/blogs/{{$post.Slug}}">
        {{if $post.Photo}}<img class="ph" src="{{pestimg $post.Photo}}" alt="{{$post.Title}}" loading="lazy" decoding="async">{{end}}
        <span class="pop-cap"><b class="clamp">{{$post.Title}}</b><small>{{upper $post.PestCategory}}</small></span>
      </a>
      {{end}}
      <a class="pop-tile planp tone-more" href="/blogs">
        <span class="pop-plan-name">The Blog</span>
        <span class="pop-plan-price">Every field guide</span>
        <span class="pop-plan-line">All topics, all pests</span>
      </a>
    </div>
  </div>
</section>

<!-- ============ BUG BLASTER ============ -->
<section class="pop-sec" id="bugblaster-sec">
  <div class="wrap">
    <div class="pop-card game">
      <div class="pop-lab">Field Exercise // Bug Blaster</div>
      <h2 class="pop-h2 dark">Take a break. Blow up some bugs.</h2>
      <p class="pop-lead dark">Hostiles crawl in from the treeline, you click (or tap) to drop a pest bomb, and everything in the blast radius goes down. This is the fun version. The real one involves our technicians.</p>
      <div class="bb-arena">
        <canvas id="bugblaster" aria-label="Bug Blaster minigame: click to drop bombs on incoming bugs"></canvas>
        <div class="bb-top"><span id="bugblaster-hud">BUGS ELIMINATED: 0</span></div>
        <div class="bb-overlay" id="bugblaster-start">
          <div class="bb-card">
            <div class="bb-title">BUG BLASTER</div>
            <p>Click / tap anywhere to drop a bomb.<br>Blast radius does the rest.</p>
            <button type="button" class="btn btn-primary">START MISSION</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>
<script src="/assets/js/bugblaster.js" defer></script>

<!-- ============ FINAL CTA ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-hero final center">
      <div>
        <h2 class="pop-h1">Ready to go <span>pest-free?</span></h2>
        <p class="pop-lead">Book online in minutes or call the line. Same-day service available across all three states. Free quotes, transparent pricing, zero hidden fees.</p>
        <div class="pop-cta-row center">
          <a class="pop-btn primary" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
          <button type="button" class="pop-btn" data-pop-open="quote">Free Quote</button>
          <a class="pop-btn" href="/prices">View Plans</a>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- ============ QUICK VIEW SOURCES (hidden; pop.js clones these) ============ -->
<div class="pop-qv-sources" hidden>
  <div id="qv-plan-bronze" data-more="/prices">
    <div class="pop-lab">Plan · Bronze</div>
    <h3>Exterior-Only Protection</h3>
    <div class="pop-price">$49<small> /MO</small></div>
    <ul><li>Exterior treatment (interior on request)</li><li>Ants, spiders, wasps, beetles, earwigs, silverfish &amp; more</li><li>Free re-treatments and a 90-day warranty</li><li>Same/next-day starts</li></ul>
  </div>
  <div id="qv-plan-silver" data-more="/prices">
    <div class="pop-lab">Plan · Silver</div>
    <h3>Interior + Exterior</h3>
    <div class="pop-price">$59<small> /MO</small></div>
    <ul><li>Interior and exterior treatment</li><li>Ants, spiders, wasps, earwigs, silverfish &amp; more</li><li>Free re-treatments and a 90-day warranty</li><li>Priority scheduling and fast response</li></ul>
  </div>
  <div id="qv-plan-gold" data-more="/prices">
    <div class="pop-lab">Plan · Gold · Most Popular</div>
    <h3>Priority Interior &amp; Exterior</h3>
    <div class="pop-price">$79<small> /MO</small></div>
    <ul><li>Interior and exterior treatment, every visit</li><li>Everything in Silver, plus cockroaches and mice</li><li>Seasonal deep-dives and 24-hour callbacks</li><li>Priority scheduling and fast response</li></ul>
  </div>
  {{range $p := .Pests}}
  <div id="qv-pest-{{$p.Slug}}" data-more="/pest/{{$p.Slug}}">
    <div class="pop-ph wide"><img src="{{pestimg $p.Filename}}" alt="{{$p.Name}}"></div>
    <div class="pop-lab">{{upper $p.Category}} · Threat {{$p.ThreatLevel}}%</div>
    <h3>{{$p.Name}}</h3>
    {{if $p.ScientificName}}<div class="pop-sci">{{$p.ScientificName}}</div>{{end}}
    <p>{{$p.Description}}</p>
    {{$c := pestcopy2 $p.Slug $p.Category}}
    {{if $c.Signs}}<div class="pop-lab">Signs of activity</div><ul>{{range $c.Signs}}<li>{{.}}</li>{{end}}</ul>{{end}}
  </div>
  {{end}}
</div>

<!-- ============ POPUP LAYER 1: QUICK VIEW ============ -->
<div class="pop-scrim" id="pop-qv" hidden>
  <div class="pop-popup" role="dialog" aria-modal="true" aria-label="Quick view">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">X</button>
    <div class="pop-qv-body"></div>
    <div class="pop-popup-row">
      <a class="pop-btn primary" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
      <a class="pop-btn pop-qv-more" href="#">More Info</a>
    </div>
  </div>
</div>

<!-- ============ POPUP: FREE QUOTE (the real /contact form) ============ -->
<div class="pop-scrim" id="pop-quote" hidden>
  <div class="pop-popup" role="dialog" aria-modal="true" aria-label="Free quote">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">X</button>
    <div class="pop-lab">Free Quote</div>
    <h3>Get your free quote</h3>
    <p class="pop-form-note">We respond within one business day, usually much faster. Same-day service available across all three states.</p>
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
`
