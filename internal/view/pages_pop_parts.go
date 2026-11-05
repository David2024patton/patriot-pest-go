package view

// Shared Pop UI building blocks (pop.patriotpest.pro). David's streaming
// rule: the surface carries illustrated tiles only; all information lives
// behind a tile tap, in the quick view and its More Info layer. These parts
// are registered as named templates so every pop page composes the same
// rails instead of laying content out on the page. A part that renders
// tiles also renders their hidden quick view sources right after them.

func init() {
	// Plans rail: the three plan posters, with the plan quick view
	// sources attached. Dot is the page data.
	RegisterPageTemplate("part-plans-rail", `
<section class="pop-sec" id="plans">
  <div class="wrap">
    <div class="pop-rail-head"><div><h2 class="pop-h2">Our Plans</h2></div>
    <div class="pop-arrows"><button type="button" class="pop-arrow" data-rail-prev="rail-plans" aria-label="Scroll plans back">&#8249;</button><button type="button" class="pop-arrow" data-rail-next="rail-plans" aria-label="Scroll plans forward">&#8250;</button></div></div>
    <div class="pop-rail" id="rail-plans">
      <button type="button" class="pop-tile poster" data-qv="qv-plan-bronze">
        <img class="ph" src="/assets/img/illus/plan-bronze.jpg" alt="Bronze plan" width="2048" height="1152" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Bronze &middot; $49/mo</b><small>Exterior-Only Protection</small></span>
      </button>
      <button type="button" class="pop-tile poster" data-qv="qv-plan-silver">
        <img class="ph" src="/assets/img/illus/plan-silver.jpg" alt="Silver plan" width="2048" height="1152" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Silver &middot; $59/mo</b><small>Interior + Exterior</small></span>
      </button>
      <button type="button" class="pop-tile poster" data-qv="qv-plan-gold">
        <span class="pop-plan-tag">Most Popular</span>
        <img class="ph" src="/assets/img/illus/plan-gold.jpg" alt="Gold plan" width="2048" height="1152" loading="lazy" decoding="async">
        <span class="pop-cap"><b>Gold &middot; $79/mo</b><small>Priority Interior &amp; Exterior</small></span>
      </button>
    </div>
  </div>
  <div class="pop-qv-sources" hidden>
    <div id="qv-plan-bronze">
      <div class="pop-qv-art"><img src="/assets/img/illus/plan-bronze.jpg" alt=""></div>
      <div class="pop-lab">Plan · Bronze</div>
      <h3>Exterior-Only Protection</h3>
      <div class="pop-price">$49<small> /MO</small></div>
      <p>A protective barrier around the outside of your home, refreshed every visit. Interior service available on request.</p>
      <div class="qv-more" hidden>
        <div class="pop-lab">What is included</div>
        <ul><li>Exterior treatment (interior on request)</li><li>Ants, spiders, wasps, beetles, earwigs, silverfish &amp; more</li><li>Free re-treatments and a 90-day warranty</li><li>Same/next-day starts</li></ul>
      </div>
    </div>
    <div id="qv-plan-silver">
      <div class="pop-qv-art"><img src="/assets/img/illus/plan-silver.jpg" alt=""></div>
      <div class="pop-lab">Plan · Silver</div>
      <h3>Interior + Exterior</h3>
      <div class="pop-price">$59<small> /MO</small></div>
      <p>Full interior and exterior treatment on every visit, with priority scheduling when something shows up between visits.</p>
      <div class="qv-more" hidden>
        <div class="pop-lab">What is included</div>
        <ul><li>Interior and exterior treatment</li><li>Ants, spiders, wasps, earwigs, silverfish &amp; more</li><li>Free re-treatments and a 90-day warranty</li><li>Priority scheduling and fast response</li></ul>
      </div>
    </div>
    <div id="qv-plan-gold">
      <div class="pop-qv-art"><img src="/assets/img/illus/plan-gold.jpg" alt=""></div>
      <div class="pop-lab">Plan · Gold · Most Popular</div>
      <h3>Priority Interior &amp; Exterior</h3>
      <div class="pop-price">$79<small> /MO</small></div>
      <p>Our highest level of cover: interior and exterior every visit, cockroach and mice coverage, and 24-hour callbacks.</p>
      <div class="qv-more" hidden>
        <div class="pop-lab">What is included</div>
        <ul><li>Interior and exterior treatment, every visit</li><li>Everything in Silver, plus cockroaches and mice</li><li>Seasonal deep-dives and 24-hour callbacks</li><li>Priority scheduling and fast response</li></ul>
      </div>
    </div>
  </div>
</section>`)

	// Pest poster tiles for a rail or grid. Dot is a []data.Pest slice.
	RegisterPageTemplate("part-pest-tiles", `
{{range $p := .}}
<button type="button" class="pop-tile poster" data-qv="qv-pest-{{$p.Slug}}">
  <img class="ph" src="{{pestillus $p.Filename}}" alt="{{$p.Name}}" width="2048" height="1152" loading="lazy" decoding="async">
  <span class="pop-cap"><b>{{$p.Name}}</b><small>{{if $p.ScientificName}}{{$p.ScientificName}} · {{end}}Threat {{$p.ThreatLevel}}%</small></span>
</button>
{{end}}`)

	// Quick view sources matching part-pest-tiles. Dot is a []data.Pest.
	RegisterPageTemplate("part-pest-sources", `
<div class="pop-qv-sources" hidden>
{{range $p := .}}
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
</div>`)

	// State poster tiles with town chips inside the quick view.
	// Dot is a []data.State slice.
	RegisterPageTemplate("part-state-tiles", `
{{range $st := .}}
<button type="button" class="pop-tile poster" data-qv="qv-state-{{$st.Code}}">
  <img class="ph" src="{{if eq $st.Code "WA"}}/assets/img/illus/state-wa.jpg{{else if eq $st.Code "ID"}}/assets/img/illus/state-id.jpg{{else}}/assets/img/illus/state-or.jpg{{end}}" alt="{{$st.Name}}" width="2048" height="1152" loading="lazy" decoding="async">
  <span class="pop-cap"><b>{{$st.Name}}</b><small>{{len $st.Cities}} zones &middot; Same-day available</small></span>
</button>
{{end}}
<div class="pop-qv-sources" hidden>
{{range $st := .}}
<div id="qv-state-{{$st.Code}}">
  <div class="pop-qv-art"><img src="{{if eq $st.Code "WA"}}/assets/img/illus/state-wa.jpg{{else if eq $st.Code "ID"}}/assets/img/illus/state-id.jpg{{else}}/assets/img/illus/state-or.jpg{{end}}" alt=""></div>
  <div class="pop-lab">Service Area &middot; {{$st.Code}}</div>
  <h3>{{$st.Name}}</h3>
  <p>Patriot covers {{len $st.Cities}} zones across {{$st.Name}}, with same-day service when it can't wait. Tap your town for its page.</p>
  <div class="pop-cities">{{range $c := $st.Cities}}<a class="pop-city" href="/areas/{{cityslug $c}}">{{$c}}</a>{{end}}</div>
</div>
{{end}}
</div>`)

	// Blog poster tiles. Dot is a []data.Post slice.
	RegisterPageTemplate("part-blog-tiles", `
{{range $post := .}}
<button type="button" class="pop-tile poster" data-qv="qv-post-{{$post.Slug}}">
  {{if $post.Photo}}<img class="ph" src="{{pestillus $post.Photo}}" alt="{{$post.Title}}" loading="lazy" decoding="async">{{end}}
  <span class="pop-cap"><b class="clamp">{{$post.Title}}</b><small>{{upper $post.PestCategory}}</small></span>
</button>
{{end}}`)

	// Quick view sources matching part-blog-tiles. Dot is a []data.Post.
	RegisterPageTemplate("part-blog-sources", `
<div class="pop-qv-sources" hidden>
{{range $post := .}}
<div id="qv-post-{{$post.Slug}}">
  {{if $post.Photo}}<div class="pop-qv-art"><img src="{{pestillus $post.Photo}}" alt=""></div>{{end}}
  <div class="pop-lab">{{upper $post.PestCategory}} &middot; Field Guide</div>
  <h3>{{$post.Title}}</h3>
  <p>{{$post.Excerpt}}</p>
  <div class="qv-more" hidden>
    <p>{{$post.Excerpt}}</p>
    <p><a class="pop-btn dark" href="/blogs/{{$post.Slug}}">Read the Full Article</a></p>
  </div>
</div>
{{end}}
</div>`)

	// Compact closing scene for inner pages: actions only, no laid-out info.
	RegisterPageTemplate("part-cta", `
<section class="pop-sec">
  <div class="wrap">
    <div class="pop-hero final center">
      <div>
        <h2 class="pop-h1">Ready to go <span>pest-free?</span></h2>
        <div class="pop-cta-row center">
          <a class="pop-btn primary" href="{{.PhoneHref}}">Call or Text {{.PhoneDisplay}}</a>
          <button type="button" class="pop-btn" data-pop-open="quote">Free Quote</button>
        </div>
      </div>
    </div>
  </div>
</section>`)
}
