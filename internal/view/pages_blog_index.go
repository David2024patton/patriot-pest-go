package view

// pageBlogIndex — templates/pages/blog-index.php. Blog landing page.
// Data: .Posts ([]data.Post, newest first).
const pageBlogIndex = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">FIELD INTEL // THE BLOG</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Pest Control <em>Intel &amp; Guides.</em></h1>
    <p class="lead">Expert identification, seasonal guides, and prevention tips from our licensed technicians, written for homeowners across Washington, Idaho, Oregon &amp; Arizona.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    {{if .Posts}}
    <div class="blog-grid">
      {{range $p := .Posts}}
      <a class="card blog-card" href="/blogs/{{$p.Slug}}" style="text-decoration:none;color:inherit;display:flex;flex-direction:column">
        {{if $p.Photo}}<span class="pphoto"><img src="{{pestimg $p.Photo}}" alt="{{if $p.PestName}}{{$p.PestName}}{{else}}{{$p.Title}}{{end}}" loading="lazy"><i class="ret" aria-hidden="true"></i></span>{{end}}
        <div class="post-meta" style="margin-top:.8rem">
          {{if $p.Season}}<span>{{ucfirst $p.Season}}</span>{{end}}
          {{if $p.PublishedAt}}<span>{{dateMDY $p.PublishedAt}}</span>{{end}}
        </div>
        <h3 style="font-family:var(--display);color:var(--cream);font-size:1.1rem;line-height:1.25;margin:.2rem 0 .5rem">{{$p.Title}}</h3>
        <p style="color:var(--khaki);font-size:.9rem;line-height:1.6;flex:1">{{$p.Excerpt}}</p>
        <span class="more" style="color:var(--orange);font-family:var(--mono);font-size:.8rem;margin-top:.8rem">READ REPORT ▸</span>
      </a>
      {{end}}
    </div>
    {{else}}
    <p class="empty">New field reports are being prepared. Check back soon.</p>
    {{end}}
  </div>
</section>`
