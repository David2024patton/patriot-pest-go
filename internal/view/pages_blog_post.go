package view

// pageBlogPost — templates/pages/blog-post.php. The unified single-post template.
// Data: .Post (data.Post), .Related ([]data.Post). BodyHTML is sanitized at
// catalog load (bluemonday UGC policy) as well as on save, so raw() is safe.
const pageBlogPost = `<article class="block">
  <div class="wrap" style="max-width:820px">
    <div class="post-meta">
      {{if .Post.Season}}<span>{{ucfirst .Post.Season}}</span>{{end}}
      {{if .Post.PublishedAt}}<span>{{dateFM .Post.PublishedAt}}</span>{{end}}
      {{if .Post.Author}}<span>By {{.Post.Author}}</span>{{end}}
    </div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(1.8rem,5vw,2.8rem);line-height:1.1;margin:.4rem 0 1rem">{{.Post.Title}}</h1>

    {{if .Post.Photo}}<span class="pphoto s-lg" style="margin-bottom:1.6rem"><img src="{{pestimg .Post.Photo}}" alt="{{if .Post.PestName}}{{.Post.PestName}}{{else}}{{.Post.Title}}{{end}}" loading="eager"><i class="ret" aria-hidden="true"></i></span>{{end}}

    {{if .Post.Excerpt}}<p class="lead" style="margin-bottom:1.6rem">{{.Post.Excerpt}}</p>{{end}}

    <div class="prose">
      {{raw .Post.BodyHTML}}
    </div>

    {{if .Post.PestSlug}}
    <div class="panel" style="margin-top:2rem">
      <strong style="color:var(--cream)">Dealing with {{lower .Post.PestName}}?</strong>
      <span class="muted"> Get the full threat file and a free quote.</span>
      <a href="/pest/{{.Post.PestSlug}}" style="color:var(--orange)">View {{.Post.PestName}} Control ▸</a>
    </div>
    {{end}}
  </div>
</article>

{{if .Related}}
<section class="block alt">
  <div class="wrap">
    <div class="eyebrow">RELATED INTEL</div>
    <h2 style="font-family:var(--display);color:var(--cream);margin:.4rem 0 1.4rem">Keep <em>reading.</em></h2>
    <div class="grid g3">
      {{range $r := .Related}}
      <a class="card" href="/blogs/{{$r.Slug}}" style="text-decoration:none;color:inherit">
        <h3 style="font-family:var(--display);color:var(--cream);font-size:1.05rem;line-height:1.3">{{$r.Title}}</h3>
        <p style="color:var(--khaki);font-size:.88rem;line-height:1.6;margin-top:.5rem">{{$r.Excerpt}}</p>
        <span class="more" style="color:var(--orange);font-family:var(--mono);font-size:.78rem">READ ▸</span>
      </a>
      {{end}}
    </div>
  </div>
</section>
{{end}}`
