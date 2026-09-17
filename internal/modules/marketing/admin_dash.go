// Admin analytics mission control: sidebar dashboard with per-section
// stats (overview, pages, clicks, traffic, devices, geography, visitors,
// activity). Server-rendered HTML, inline SVG charts, inline vanilla JS.
// Zero external JS/CSS/fonts. All admin routes stay behind requireAdmin.
package marketing

import (
	"database/sql"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

var adminSections = []struct{ id, label, icon string }{
	{"overview", "Overview", "◎"},
	{"pages", "Pages", "▤"},
	{"clicks", "Clicks", "➤"},
	{"traffic", "Traffic", "⇄"},
	{"devices", "Devices", "▣"},
	{"geo", "Geography", "◉"},
	{"visitors", "Visitors", "◍"},
	{"activity", "Activity", "≋"},
}

// Extra CSS for the mission-control layout (appended after adminCSS).
const adminDashCSS = `
.layout{display:flex;min-height:100vh}
aside.side{width:230px;flex:none;background:var(--olive-900);border-right:2px solid var(--olive-700);padding:1rem .8rem;position:sticky;top:0;height:100vh;overflow-y:auto}
.side .brand{font-weight:800;letter-spacing:.14em;font-size:.95rem;margin-bottom:.2rem}
.side .brand .star{color:var(--orange)}
.side .sub{color:var(--olive-300);font-size:.72rem;letter-spacing:.08em;margin-bottom:1rem}
.side nav{display:flex;flex-direction:column;gap:.25rem}
.side nav a{color:var(--khaki);text-decoration:none;border:1px solid transparent;border-radius:6px;padding:.5rem .7rem;font-size:.85rem;display:flex;gap:.6rem;align-items:center}
.side nav a .ic{width:1.2rem;text-align:center;color:var(--olive-300)}
.side nav a.on{background:var(--olive-800);border-color:var(--orange);color:var(--cream)}
.side nav a:hover{border-color:var(--olive-500);color:var(--cream)}
.side .days{display:flex;gap:.3rem;margin:1rem 0;flex-wrap:wrap}
.side .days a{color:var(--olive-300);text-decoration:none;border:1px solid var(--olive-700);border-radius:6px;padding:.3rem .55rem;font-size:.75rem}
.side .days a.on{border-color:var(--orange);color:var(--cream)}
.side .me{margin-top:1rem;padding-top:1rem;border-top:1px solid var(--olive-700);font-size:.75rem;color:var(--olive-300)}
.side .me a{color:var(--khaki)}
main.main{flex:1;min-width:0;padding:1.2rem;max-width:1200px}
.topbar{display:none}
.grid2{display:grid;grid-template-columns:1fr 1fr;gap:.8rem}
.grid3{display:grid;grid-template-columns:repeat(3,1fr);gap:.8rem}
@media(max-width:900px){.grid2,.grid3{grid-template-columns:1fr}}
.burger{display:none;background:var(--olive-900);border:1px solid var(--olive-700);color:var(--cream);border-radius:6px;padding:.5rem .8rem;font-size:1rem;cursor:pointer;margin-bottom:1rem}
@media(max-width:760px){
  .layout{flex-direction:column}
  aside.side{position:fixed;z-index:50;left:0;top:0;transform:translateX(-105%);transition:transform .2s ease;width:250px}
  body.nav-open aside.side{transform:none}
  .burger{display:inline-block}
  main.main{padding:.9rem}
}
th.sortable{cursor:pointer;user-select:none;white-space:nowrap}
th.sortable:hover{color:var(--cream)}
th.sortable .arr{color:var(--orange);font-size:.7rem}
tr.drill{display:none}
tr.drill.open{display:table-row}
tr.drill td{background:var(--olive-950)}
.clickable{cursor:pointer}
.badge{display:inline-block;font-size:.68rem;letter-spacing:.1em;border:1px solid var(--olive-700);border-radius:4px;padding:.15rem .45rem;color:var(--khaki)}
.badge.phone{border-color:var(--orange);color:var(--orange-hot)}
.badge.desktop{border-color:var(--olive-300);color:var(--olive-300)}
form.filters{display:flex;gap:.6rem;flex-wrap:wrap;align-items:flex-end;margin-bottom:1rem}
form.filters label{font-size:.68rem;letter-spacing:.14em;color:var(--olive-300);display:block;margin-bottom:.25rem}
form.filters select,form.filters input{background:var(--olive-950);border:1px solid var(--olive-700);border-radius:6px;color:var(--cream);padding:.45rem .6rem;font-size:.85rem}
form.filters button{background:var(--orange);border:0;border-radius:6px;color:#141a10;font-weight:800;padding:.5rem 1rem;cursor:pointer}
#tip{position:fixed;z-index:99;pointer-events:none;background:#0a0d07;border:1px solid var(--orange);color:var(--cream);font-size:.75rem;border-radius:6px;padding:.35rem .6rem;display:none;max-width:260px;white-space:pre-line}
.spark{display:block}
.legend{display:flex;gap:1rem;flex-wrap:wrap;font-size:.75rem;color:var(--khaki);margin-top:.5rem}
.legend i{display:inline-block;width:.8rem;height:.8rem;border-radius:2px;margin-right:.3rem;vertical-align:-1px}
.toggle{display:flex;align-items:center;gap:.5rem;font-size:.8rem;color:var(--khaki)}
.toggle input{accent-color:var(--orange)}
.empty{color:var(--olive-300);font-size:.85rem}`

// ---- small format helpers ----

func esc(s string) string { return html.EscapeString(s) }

// fnum formats with thousands separators.
func fnum(n int64) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}

// fdur formats seconds as a human duration.
func fdur(sec float64) string {
	if sec < 1 {
		return "<1s"
	}
	s := int64(sec)
	if s < 60 {
		return strconv.FormatInt(s, 10) + "s"
	}
	if s < 3600 {
		return strconv.FormatInt(s/60, 10) + "m " + strconv.FormatInt(s%60, 10) + "s"
	}
	return strconv.FormatInt(s/3600, 10) + "h " + strconv.FormatInt((s%3600)/60, 10) + "m"
}

// fpct formats a 0..1 ratio as a percent string.
func fpct(f float64) string {
	return strconv.FormatInt(int64(f*100+0.5), 10) + "%"
}

func pctOf(part, total int64) string {
	if total <= 0 {
		return "0%"
	}
	return strconv.FormatInt(part*100/total, 10) + "%"
}

// ago renders a unix timestamp as "5m ago".
func ago(ts, now int64) string {
	d := now - ts
	if d < 0 {
		d = 0
	}
	switch {
	case d < 60:
		return strconv.FormatInt(d, 10) + "s ago"
	case d < 3600:
		return strconv.FormatInt(d/60, 10) + "m ago"
	case d < 86400:
		return strconv.FormatInt(d/3600, 10) + "h ago"
	default:
		return time.Unix(ts, 0).Format("Jan 02 15:04")
	}
}

var chartPalette = []string{"#f4772e", "#c8b98c", "#8fa05e", "#5c6f3a", "#d4a24e", "#7ba05b", "#a05b5b", "#5b7ba0"}

// ---- inline SVG charts (with data-tip for the JS tooltip) ----

type barItem struct {
	Label string
	Value float64
	Tip   string
}

// svgBars renders a vertical bar chart.
func svgBars(items []barItem, w, h int, color string) string {
	var sb strings.Builder
	if len(items) == 0 {
		return `<p class="empty">No data in this range.</p>`
	}
	maxV := 1.0
	for _, it := range items {
		if it.Value > maxV {
			maxV = it.Value
		}
	}
	n := float64(len(items))
	fmt.Fprintf(&sb, `<svg class="chart" viewBox="0 0 %d %d" role="img">`, w, h)
	for i, it := range items {
		bw := float64(w) / n
		bh := it.Value / maxV * float64(h-20)
		if bh < 1 && it.Value > 0 {
			bh = 1
		}
		x := float64(i)*bw + 1
		y := float64(h) - bh
		ww := bw - 2
		if ww < 1 {
			ww = 1
		}
		tip := it.Tip
		if tip == "" {
			tip = it.Label + ": " + fnum(int64(it.Value+0.5))
		}
		fmt.Fprintf(&sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" data-tip="%s"><title>%s</title></rect>`,
			x, y, ww, bh, color, esc(tip), esc(tip))
	}
	sb.WriteString(`</svg>`)
	return sb.String()
}

// svgStacked renders two stacked series (a on bottom, b on top).
func svgStacked(labels []string, a, b []float64, aLabel, bLabel string, w, h int) string {
	var sb strings.Builder
	if len(labels) == 0 {
		return `<p class="empty">No data in this range.</p>`
	}
	maxV := 1.0
	for i := range labels {
		if t := a[i] + b[i]; t > maxV {
			maxV = t
		}
	}
	n := float64(len(labels))
	fmt.Fprintf(&sb, `<svg class="chart" viewBox="0 0 %d %d" role="img">`, w, h)
	for i := range labels {
		bw := float64(w) / n
		ha := a[i] / maxV * float64(h-20)
		hb := b[i] / maxV * float64(h-20)
		if ha < 1 && a[i] > 0 {
			ha = 1
		}
		if hb < 1 && b[i] > 0 {
			hb = 1
		}
		x := float64(i)*bw + 1
		ww := bw - 2
		if ww < 1 {
			ww = 1
		}
		ya := float64(h) - ha
		yb := ya - hb
		tip := labels[i] + "\n" + aLabel + ": " + fnum(int64(a[i]+0.5)) + "\n" + bLabel + ": " + fnum(int64(b[i]+0.5))
		fmt.Fprintf(&sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#5c6f3a" data-tip="%s"><title>%s</title></rect>`,
			x, ya, ww, ha, esc(tip), esc(tip))
		fmt.Fprintf(&sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#f4772e" data-tip="%s"><title>%s</title></rect>`,
			x, yb, ww, hb, esc(tip), esc(tip))
	}
	sb.WriteString(`</svg>`)
	return sb.String()
}

// svgDonut renders a donut chart from name/count pairs.
func svgDonut(items []data.NameCount, size int) string {
	var sb strings.Builder
	var total int64
	for _, it := range items {
		total += it.Count
	}
	if total == 0 {
		return `<p class="empty">No data in this range.</p>`
	}
	r := 54.0
	c := 2 * 3.14159265 * r
	fmt.Fprintf(&sb, `<svg viewBox="0 0 140 140" width="%d" height="%d" role="img" style="max-width:100%%;height:auto">`, size, size)
	off := 0.0
	for i, it := range items {
		frac := float64(it.Count) / float64(total)
		length := frac * c
		color := chartPalette[i%len(chartPalette)]
		tip := it.Name + ": " + fnum(it.Count) + " (" + fpct(frac) + ")"
		fmt.Fprintf(&sb, `<circle cx="70" cy="70" r="%.1f" fill="none" stroke="%s" stroke-width="22" stroke-dasharray="%.2f %.2f" stroke-dashoffset="%.2f" transform="rotate(-90 70 70)" data-tip="%s"><title>%s</title></circle>`,
			r, color, length-1.5, c-length+1.5, -off, esc(tip), esc(tip))
		off += length
	}
	fmt.Fprintf(&sb, `<text x="70" y="66" text-anchor="middle" fill="#f5f1e4" font-size="20" font-weight="800">%s</text>`, fnum(total))
	fmt.Fprintf(&sb, `<text x="70" y="84" text-anchor="middle" fill="#8fa05e" font-size="10">TOTAL</text>`)
	sb.WriteString(`</svg>`)
	return sb.String()
}

// svgSpark renders a tiny area sparkline.
func svgSpark(vals []int64, w, h int, color string) string {
	var sb strings.Builder
	if len(vals) == 0 {
		return ""
	}
	maxV := int64(1)
	for _, v := range vals {
		if v > maxV {
			maxV = v
		}
	}
	n := len(vals)
	var pts []string
	for i, v := range vals {
		var x float64
		if n > 1 {
			x = float64(i) / float64(n-1) * float64(w)
		} else {
			x = float64(w) / 2
		}
		y := float64(h) - float64(v)/float64(maxV)*float64(h-2) - 1
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", x, y))
	}
	line := strings.Join(pts, " ")
	fmt.Fprintf(&sb, `<svg class="spark" width="%d" height="%d" viewBox="0 0 %d %d"><polygon points="0,%d %s %d,%d" fill="%s" opacity="0.25"/><polyline points="%s" fill="none" stroke="%s" stroke-width="1.5"/></svg>`,
		w, h, w, h, h, line, w, h, color, line, color)
	return sb.String()
}

// hbarRow renders one horizontal bar row for bar-list charts.
func hbarRow(label string, count, total int64, color string) string {
	wPct := pctOf(count, total)
	return `<div style="display:flex;align-items:center;gap:.6rem;margin:.35rem 0">` +
		`<div style="width:11rem;flex:none;font-size:.8rem;color:var(--khaki);overflow:hidden;text-overflow:ellipsis;white-space:nowrap" title="` + esc(label) + `">` + esc(label) + `</div>` +
		`<div class="bar" style="flex:1"><i style="width:` + wPct + `;background:` + color + `"></i></div>` +
		`<div class="num" style="width:4.5rem;text-align:right;font-size:.8rem;font-variant-numeric:tabular-nums">` + fnum(count) + `</div>` +
		`<div style="width:3rem;text-align:right;font-size:.75rem;color:var(--olive-300)">` + wPct + `</div></div>`
}

// donutLegend renders a color legend for a donut.
func donutLegend(items []data.NameCount) string {
	var sb strings.Builder
	var total int64
	for _, it := range items {
		total += it.Count
	}
	sb.WriteString(`<div class="legend">`)
	for i, it := range items {
		color := chartPalette[i%len(chartPalette)]
		fmt.Fprintf(&sb, `<span><i style="background:%s"></i>%s %s (%s)</span>`,
			color, esc(it.Name), fnum(it.Count), pctOf(it.Count, total))
	}
	sb.WriteString(`</div>`)
	return sb.String()
}

// ---- page shell ----

// adminDashboard routes /admin?s=<section>&days=<n>.
func (m *Module) adminDashboard(w http.ResponseWriter, r *http.Request) {
	days := 30
	if d, err := strconv.Atoi(r.URL.Query().Get("days")); err == nil {
		days = d
	}
	switch days {
	case 1, 7, 30, 90:
	default:
		days = 30
	}
	section := r.URL.Query().Get("s")
	okSec := false
	for _, sec := range adminSections {
		if sec.id == section {
			okSec = true
			break
		}
	}
	if !okSec {
		section = "overview"
	}
	db := m.adminDB(w)
	if db == nil {
		return
	}
	email, _ := r.Context().Value(adminEmailCtxKey).(string)
	body, err := m.renderSection(r, db, section, days)
	if err != nil {
		slog.Error("admin: dashboard section failed", "section", section, "err", err.Error())
		http.Error(w, "Could not load analytics.", http.StatusInternalServerError)
		return
	}
	m.renderShell(w, r, section, days, email, body)
}

// renderSection dispatches to the per-section renderer.
func (m *Module) renderSection(r *http.Request, db *sql.DB, section string, days int) (string, error) {
	now := time.Now().Unix()
	switch section {
	case "pages":
		return m.secPages(db, days, now)
	case "clicks":
		return m.secClicks(r, db, days, now)
	case "traffic":
		return m.secTraffic(r, db, days, now)
	case "devices":
		return m.secDevices(db, days, now)
	case "geo":
		return m.secGeo(r, db, days, now)
	case "visitors":
		return m.secVisitors(db, days, now)
	case "activity":
		return m.secActivity(r, now)
	default:
		return m.secOverview(db, days, now)
	}
}

// renderShell writes the full page: sidebar + section body + inline JS.
func (m *Module) renderShell(w http.ResponseWriter, r *http.Request, section string, days int, email, body string) {
	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8">` +
		`<meta name="viewport" content="width=device-width,initial-scale=1">` +
		`<meta name="robots" content="noindex,nofollow">` +
		`<title>Analytics Mission Control | Patriot Pest Control</title><style>` +
		adminCSS + adminDashCSS + `</style></head><body>`)
	sb.WriteString(`<div class="layout"><aside class="side">`)
	sb.WriteString(`<div class="brand"><span class="star">&#9733;</span> PATRIOT</div><div class="sub">ANALYTICS MISSION CONTROL</div>`)
	sb.WriteString(`<nav>`)
	for _, sec := range adminSections {
		cls := ""
		if sec.id == section {
			cls = ` class="on"`
		}
		sb.WriteString(`<a` + cls + ` href="/admin?s=` + sec.id + `&days=` + strconv.Itoa(days) + `"><span class="ic">` + sec.icon + `</span>` + sec.label + `</a>`)
	}
	sb.WriteString(`</nav><div class="days">`)
	for _, d := range []int{1, 7, 30, 90} {
		cls := ""
		if d == days {
			cls = ` class="on"`
		}
		lbl := strconv.Itoa(d) + "d"
		if d == 1 {
			lbl = "24h"
		}
		sb.WriteString(`<a` + cls + ` href="/admin?s=` + section + `&days=` + strconv.Itoa(d) + `">` + lbl + `</a>`)
	}
	sb.WriteString(`</div><div class="me">` + esc(email) + `<br><a href="/admin/logout">LOGOUT</a></div>`)
	sb.WriteString(`</aside><main class="main">`)
	sb.WriteString(`<button class="burger" id="burger" aria-label="Menu">&#9776; MENU</button>`)
	sb.WriteString(body)
	sb.WriteString(`<div class="foot">FIRST-PARTY TRACKING &middot; YOUR VISITS EXCLUDED (ADMIN COOKIE) &middot; BOTS FILTERED &middot; DNT RESPECTED &middot; NO RAW IPS STORED</div>`)
	sb.WriteString(`</main></div><div id="tip"></div>`)
	sb.WriteString(`<script>` + adminDashJS + `</script></body></html>`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	_, _ = w.Write([]byte(sb.String()))
}

// ---- inline JS: sortable tables, tooltips, sidebar, live feed ----

const adminDashJS = `
(function(){
"use strict";
/* sortable tables: click a header to sort by that column */
document.querySelectorAll("table.sortable").forEach(function(tbl){
  var ths = tbl.querySelectorAll("thead th");
  ths.forEach(function(th, idx){
    if(!th.classList.contains("s")) return;
    th.classList.add("sortable");
    th.addEventListener("click", function(){
      var tbody = tbl.querySelector("tbody");
      var rows = Array.prototype.slice.call(tbody.rows);
      var num = th.hasAttribute("data-num");
      var asc = th.getAttribute("data-dir") === "asc";
      rows.sort(function(a, b){
        var av = a.cells[idx].getAttribute("data-v") || a.cells[idx].textContent;
        var bv = b.cells[idx].getAttribute("data-v") || b.cells[idx].textContent;
        var c;
        if(num){ c = parseFloat(av) - parseFloat(bv); }
        else { c = String(av).localeCompare(String(bv)); }
        return asc ? c : -c;
      });
      rows.forEach(function(r){ tbody.appendChild(r); });
      ths.forEach(function(h){ h.removeAttribute("data-dir"); var e=h.querySelector(".arr"); if(e) e.remove(); });
      th.setAttribute("data-dir", asc ? "desc" : "asc");
      var s = document.createElement("span"); s.className = "arr"; s.textContent = asc ? " \\u25BC" : " \\u25B2";
      th.appendChild(s);
    });
  });
});
/* tooltips for [data-tip] */
var tip = document.getElementById("tip");
document.addEventListener("mousemove", function(ev){
  var t = ev.target && ev.target.closest ? ev.target.closest("[data-tip]") : null;
  if(t){
    tip.textContent = t.getAttribute("data-tip");
    tip.style.display = "block";
    tip.style.left = (ev.clientX + 14) + "px";
    tip.style.top = (ev.clientY + 14) + "px";
  } else {
    tip.style.display = "none";
  }
});
/* drill-down rows */
document.querySelectorAll("[data-drill]").forEach(function(el){
  el.addEventListener("click", function(){
    var row = document.getElementById(el.getAttribute("data-drill"));
    if(row) row.classList.toggle("open");
  });
});
/* mobile sidebar */
var burger = document.getElementById("burger");
if(burger) burger.addEventListener("click", function(){ document.body.classList.toggle("nav-open"); });
/* live activity feed */
var feedBody = document.getElementById("feedbody");
var feedKind = document.getElementById("feedkind");
var feedQ = document.getElementById("feedq");
var feedAuto = document.getElementById("feedauto");
function feedURL(){
  var u = "/admin/feed?kind=" + encodeURIComponent(feedKind ? feedKind.value : "");
  if(feedQ && feedQ.value) u += "&q=" + encodeURIComponent(feedQ.value);
  return u;
}
function refreshFeed(){
  if(!feedBody) return;
  fetch(feedURL(), {credentials: "same-origin"}).then(function(r){ return r.text(); }).then(function(html){
    feedBody.innerHTML = html;
  }).catch(function(){});
}
var feedTimer = null;
if(feedAuto){
  feedAuto.addEventListener("change", function(){
    if(feedTimer){ clearInterval(feedTimer); feedTimer = null; }
    if(feedAuto.checked){ feedTimer = setInterval(refreshFeed, 5000); }
  });
}
var feedForm = document.getElementById("feedform");
if(feedForm) feedForm.addEventListener("submit", function(ev){ ev.preventDefault(); refreshFeed(); });
})();
`

// sectionHead renders the section title row.
func sectionHead(title, sub string) string {
	return `<div style="margin-bottom:1.2rem"><h1 style="margin:0;font-size:1.25rem;letter-spacing:.12em">` +
		esc(title) + `</h1><div class="sub">` + esc(sub) + `</div></div>`
}

// statCard renders one KPI card.
func statCard(k, v, s string) string {
	return `<div class="card"><div class="k">` + esc(k) + `</div><div class="v">` + esc(v) +
		`</div><div class="s">` + esc(s) + `</div></div>`
}
