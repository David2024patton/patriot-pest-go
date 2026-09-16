package view

// pageSearch — Go port of templates/pages/search.php. Site search form with
// server-rendered results across the pest library, blog guides, and service
// areas. Sections render only when they have matches; an empty query shows
// the intro copy instead.
const pageSearch = `<section class="block">
  <div class="wrap" style="max-width:860px">
    <div class="eyebrow">SEARCH // SITE</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Search</h1>
    <form action="/search" method="get" style="display:flex;gap:.6rem;margin-bottom:2rem">
      <input type="text" name="q" value="{{.Q}}" placeholder="Pests, guides, cities…" aria-label="Search" style="flex:1;min-height:48px;background:var(--olive-800);border:1px solid var(--olive-700);color:var(--cream);padding:0 1rem;font-family:var(--mono);font-size:.9rem">
      <button type="submit" style="min-height:48px;background:var(--orange);color:var(--ink);border:0;padding:0 1.4rem;font-family:var(--display);text-transform:uppercase;cursor:pointer">Search</button>
    </form>

{{if .HasQuery}}
  {{if .HasResults}}
    {{if .Pests}}
      <div class="eyebrow">PEST LIBRARY</div>
      <div class="grid g3" style="margin:1rem 0 2rem">
        {{range .Pests}}<a class="card" href="/pest/{{.Slug}}" style="text-decoration:none;color:inherit;display:flex;gap:.9rem;align-items:center">
          <span class="pphoto" style="width:70px;height:52px"><img src="{{pestimg .Filename}}" alt="{{.Name}}" loading="lazy"></span>
          <div>
            <h3 style="font-family:var(--display);color:var(--cream);font-size:.95rem">{{.Name}}</h3>
            <p style="color:var(--khaki);font-size:.8rem;line-height:1.4">{{.Desc}}</p>
          </div>
        </a>{{end}}
      </div>
    {{end}}
    {{if .Posts}}
      <div class="eyebrow">BLOG GUIDES</div>
      <div class="grid g3" style="margin:1rem 0 2rem">
        {{range .Posts}}<a class="card" href="/blogs/{{.Slug}}" style="text-decoration:none;color:inherit">
          <h3 style="font-family:var(--display);color:var(--cream);font-size:1rem;line-height:1.3">{{.Title}}</h3>
          <p style="color:var(--khaki);font-size:.85rem;line-height:1.55;margin-top:.4rem">{{.Excerpt}}</p>
        </a>{{end}}
      </div>
    {{end}}
    {{if .Cities}}
      <div class="eyebrow">SERVICE AREAS</div>
      <div style="display:flex;flex-wrap:wrap;gap:.6rem;margin:1rem 0">
        {{range .Cities}}<a href="/areas/{{cityslug .City}}" style="background:var(--olive-800);border:1px solid var(--olive-700);color:var(--cream);padding:.5rem .9rem;text-decoration:none;font-family:var(--mono);font-size:.8rem">{{.City}} · {{.Code}}</a>{{end}}
      </div>
    {{end}}
  {{else}}
    <p class="lead" style="color:var(--khaki)">No results for "{{.Q}}". Try a different term or <a href="/contact" style="color:var(--orange)">ask us directly</a>.</p>
  {{end}}
{{else}}
  <p class="lead" style="color:var(--khaki)">Type a pest, a topic, or a city to search the blog, pest library, and service areas.</p>
{{end}}
  </div>
</section>`
