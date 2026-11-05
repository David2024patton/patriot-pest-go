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
      <img class="pop-scene-bg" src="/assets/img/illus/hero.jpg" alt="">
      <div class="pop-scene-shade"></div>
      <div class="pop-scene-in">
        <div class="pop-tagpill">A PROUD VETERAN-OWNED AMERICAN COMPANY</div>
        <h1 class="pop-scene-h">Hostile pests.<br><span>Neutralized.</span></h1>
        <p class="pop-scene-sub">Military-precision pest control across <b>Washington, Idaho &amp; Oregon</b>, with <b>same-day service</b> when it can't wait.</p>
        <div class="pop-cta-row">
          <button type="button" class="pop-btn primary" data-pop-open="quote">Book Now</button>
          <a class="pop-btn outline" href="/#plans" data-scroll="plans">See the Plans</a>
          <a class="pop-btn light" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- ============ PLANS RAIL (shared part) ============ -->
{{template "part-plans-rail" .}}

<!-- ============ THREAT FILES RAIL (DB-driven, every pest) ============ -->
<section class="pop-sec" id="services">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Services</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-pests" aria-label="Scroll pests back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-pests" aria-label="Scroll pests forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-pests">
      {{range $p := .TopPests}}
      <button type="button" class="pop-tile poster" data-qv="qv-pest-{{$p.Slug}}">
        <img class="ph" src="{{pestillus $p.Filename}}" alt="{{$p.Name}}" width="2048" height="1152" loading="lazy" decoding="async">
        <span class="pop-cap"><b>{{$p.Name}}</b><small>{{if $p.ScientificName}}{{$p.ScientificName}} · {{end}}Threat {{$p.ThreatLevel}}%</small></span>
      </button>
      {{end}}
      <a class="pop-tile poster" href="/services">
        <img class="ph" src="{{pestillus "crickets.jpg"}}" alt="View all pests" width="2096" height="1184" loading="lazy" decoding="async">
        <span class="pop-cap"><b>View All Pests</b><small>{{.PestCount}} services &middot; Full catalog</small></span>
      </a>
    </div>
  </div>
</section>

<!-- ============ OPERATOR (tile; the file lives in the popup) ============ -->
<section class="pop-sec" id="about">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">The Operator</h2></div></div>
    <div class="pop-feature">
      <img class="pop-feature-art" src="/assets/img/illus/operator.jpg?v=20261105-op2" alt="Skyler Rose, founder and operator">
      <div class="pop-feature-in">
        <h3 class="pop-feature-h">Skyler Rose</h3>
        <p class="pop-feature-p">Founder &amp; Operator &middot; U.S. Military Veteran. The man who answers for every job Patriot runs.</p>
        <div class="pop-scene-cta"><button type="button" class="pop-btn" data-qv="qv-operator">Meet Skyler</button><a class="pop-btn primary" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a></div>
      </div>
    </div>
  </div>
</section>

<!-- ============ AREAS (state tiles; towns live in the popup) ============ -->
<section class="pop-sec" id="areas">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Area of Operations</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-areas" aria-label="Scroll areas back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-areas" aria-label="Scroll areas forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-areas">{{template "part-state-tiles" .States}}</div>
  </div>
</section>

<!-- ============ GUARANTEE (tile; the numbers live in the popup) ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">The Guarantee</h2></div></div>
    <button type="button" class="pop-tile poster wide" data-qv="qv-guarantee">
      <img class="ph" src="/assets/img/illus/guarantee.jpg" alt="The Patriot guarantee" width="2048" height="1152" loading="lazy" decoding="async">
      <span class="pop-cap"><b>90-Day Mission Guarantee</b><small>We stand behind every mission</small></span>
    </button>
  </div>
</section>

<!-- ============ FIELD REPORTS ============ -->
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Field Reports</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-reviews" aria-label="Scroll reviews back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-reviews" aria-label="Scroll reviews forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-reviews">
      <div class="pop-card review"><div class="pop-stars">★★★★★</div><p>"Patriot Pest Control saved our home from a serious ant infestation. Professional, thorough, and the results were immediate."</p><div class="pop-who"><b>Sarah M.</b><small>WASHINGTON</small></div></div>
      <div class="pop-card review"><div class="pop-stars">★★★★★</div><p>"We've used Patriot for our restaurant for 2 years. Reliable, discreet, always on time. Our health inspections have never been better."</p><div class="pop-who"><b>James T.</b><small>LIBERTY LAKE, WA</small></div></div>
      <div class="pop-card review"><div class="pop-stars">★★★★★</div><p>"Fast response and eco-friendly products safe for my kids and pets. The technician was knowledgeable and explained everything."</p><div class="pop-who"><b>Lisa R.</b><small>COEUR D'ALENE, ID</small></div></div>
    </div>
  </div>
</section>

<!-- ============ FROM THE FIELD (blog rail) ============ -->
<section class="pop-sec" id="blog">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">From the Field</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-blog" aria-label="Scroll articles back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-blog" aria-label="Scroll articles forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-blog">
      {{range $post := .Posts}}
      <button type="button" class="pop-tile poster" data-qv="qv-post-{{$post.Slug}}">
        {{if $post.Photo}}<img class="ph" src="{{pestillus $post.Photo}}" alt="{{$post.Title}}" loading="lazy" decoding="async">{{end}}
        <span class="pop-cap"><b class="clamp">{{$post.Title}}</b><small>{{upper $post.PestCategory}}</small></span>
      </button>
      {{end}}
      <a class="pop-tile planp tone-more" href="/blogs">
        <span class="pop-plan-name">The Blog</span>
        <span class="pop-plan-price">Every field guide</span>
        <span class="pop-plan-line">All topics, all pests</span>
      </a>
    </div>
  </div>
</section>

<!-- ============ BUG BLASTER (tile; the game lives in its popup) ============ -->
<section class="pop-sec" id="bugblaster-sec">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Play Our Game</h2></div></div>
    <button type="button" class="pop-tile poster wide" data-pop-open="game">
      <img class="ph" src="/assets/img/illus/game-art.jpg" alt="Bug Blaster game" width="2048" height="1152" loading="lazy" decoding="async">
      <span class="pop-cap"><b>Bug Blaster</b><small>Tap to play Bug Blaster</small></span>
    </button>
  </div>
</section>
<div class="pop-scrim" id="pop-game" hidden>
  <div class="pop-sheet"></div>
  <div class="pop-popup wide-pop" role="dialog" aria-modal="true" aria-label="Bug Blaster game">
    <button type="button" class="pop-x" data-pop-close aria-label="Close">&#10005;</button>
<div class="pop-card game">
      <div class="pop-lab">Play Our Game // Bug Blaster</div>
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
</div>
<script src="/assets/js/bugblaster.js?v=20261105-bb2" defer></script>

<!-- ============ QUICK VIEW SOURCES (hidden; pop.js clones these) ============ -->
<div class="pop-qv-sources" hidden>
  <div id="qv-operator">
    <div class="pop-qv-art"><img src="/assets/img/illus/operator.jpg?v=20261105-op2" alt=""></div>
    <div class="pop-lab">Personnel File 001</div>
    <h3>Skyler Rose</h3>
    <p>Patriot Pest Control was founded by U.S. Military Veteran Skyler Rose, bringing military discipline, integrity, and uncompromising excellence to pest control across three states. Over a decade of field experience. Thousands of homes and businesses protected.</p>
    <div class="qv-more" hidden>
      <div class="pop-lab">Personnel File</div>
      <dl class="pop-dl">
        <div><dt>Name</dt><dd>SKYLER ROSE</dd></div>
        <div><dt>Role</dt><dd>FOUNDER &amp; OPERATOR</dd></div>
        <div><dt>Service</dt><dd>U.S. MILITARY VETERAN</dd></div>
        <div><dt>Theater</dt><dd>WA &middot; ID &middot; OR</dd></div>
        <div><dt>Clearance</dt><dd>LICENSED &middot; BONDED &middot; INSURED</dd></div>
        <div><dt>Status</dt><dd>ACTIVE - SAME-DAY RESPONSE</dd></div>
      </dl>
      <blockquote class="pop-quote">
        <p>After serving our country, I founded Patriot Pest Control to continue serving American families and businesses with the same dedication, precision, and integrity I learned in the military. We're not just eliminating pests. We're protecting what matters most.</p>
        <div class="pop-quote-foot"><b>Skyler Rose</b><small>FOUNDER &amp; VETERAN</small></div>
      </blockquote>
    </div>
  </div>
  <div id="qv-guarantee">
    <div class="pop-qv-art"><img src="/assets/img/illus/guarantee.jpg" alt=""></div>
    <div class="pop-lab">The Guarantee</div>
    <h3>90-Day Mission Guarantee</h3>
    <dl class="pop-dl">
      <div><dt>Satisfaction</dt><dd>100% GUARANTEED</dd></div>
      <div><dt>Warranty</dt><dd>90 DAYS ON TREATMENTS</dd></div>
      <div><dt>Scheduling</dt><dd>24/7 APPOINTMENTS</dd></div>
      <div><dt>Coverage</dt><dd>3 STATES: WA &middot; ID &middot; OR</dd></div>
      <div><dt>Experience</dt><dd>10+ YEARS IN THE FIELD</dd></div>
    </dl>
    <div class="qv-more" hidden>
      <div class="pop-lab">Our Promise</div>
      <p>If pests return between scheduled visits, we re-treat at no additional cost. No hassles, no excuses. Licensed, bonded, and insured, with eco-friendly, low-toxicity products safe for kids, pets, and the environment.</p>
    </div>
  </div>
  {{range $p := .TopPests}}
  <div id="qv-pest-{{$p.Slug}}" data-more="/pest/{{$p.Slug}}">
    <div class="pop-qv-art"><img src="{{pestillus $p.Filename}}" alt=""></div>
    <div class="pop-lab">{{upper $p.Category}} · Threat {{$p.ThreatLevel}}%</div>
    <h3>{{$p.Name}}</h3>
    {{if $p.ScientificName}}<div class="pop-sci">{{$p.ScientificName}}</div>{{end}}
    <p>{{$p.Description}}</p>
    {{$c := pestcopy2 $p.Slug $p.Category}}
    {{if $c.Signs}}<div class="qv-more" hidden><div class="pop-lab">Signs of activity</div><ul>{{range $c.Signs}}<li>{{.}}</li>{{end}}</ul><p>The full service briefing, treatment plan, and pricing live on this pest's page.</p></div>{{end}}
  </div>
  {{end}}
  {{range $post := .Posts}}
  <div id="qv-post-{{$post.Slug}}">
    {{if $post.Photo}}<div class="pop-qv-art"><img src="{{pestillus $post.Photo}}" alt=""></div>{{end}}
    <div class="pop-lab">{{upper $post.PestCategory}} · Field Guide</div>
    <h3>{{$post.Title}}</h3>
    <p>{{$post.Excerpt}}</p>
    <div class="qv-more" hidden>
      <p>{{$post.Excerpt}}</p>
      <p><a class="pop-btn dark" href="/blogs/{{$post.Slug}}">Read the Full Article</a></p>
    </div>
  </div>
  {{end}}
</div>
`
