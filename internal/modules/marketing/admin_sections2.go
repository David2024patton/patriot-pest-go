// Dashboard section renderers, part 2: devices, geography, visitors,
// activity, plus the live activity feed endpoint.
package marketing

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// ---- DEVICES ----

func (m *Module) secDevices(db *sql.DB, days int, now int64) (string, error) {
	dev, err := data.QueryDevices(db, days, now)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(sectionHead("Devices", "Phone vs computer vs tablet, OS, browser, and screen sizes. Parsed from User-Agent at ingest."))
	var total int64
	for _, t := range dev.Types {
		total += t.Count
	}
	sb.WriteString(`<div class="grid2">`)
	sb.WriteString(`<div class="panel"><h2>DEVICE TYPE</h2><div style="display:flex;gap:1.5rem;align-items:center;flex-wrap:wrap">` +
		svgDonut(dev.Types, 170) + `<div style="flex:1;min-width:200px">` + donutLegend(dev.Types) + `</div></div></div>`)
	sb.WriteString(`<div class="panel"><h2>SCREEN WIDTH BUCKETS</h2>`)
	if len(dev.Buckets) == 0 {
		sb.WriteString(`<p class="empty">No screen data yet. New pageviews carry it.</p>`)
	} else {
		var bt int64
		for _, b := range dev.Buckets {
			bt += b.Count
		}
		for i, b := range dev.Buckets {
			sb.WriteString(hbarRow(b.Name, b.Count, bt, chartPalette[i%len(chartPalette)]))
		}
	}
	sb.WriteString(`<p class="muted">Buckets by CSS pixel width. Mobile &lt;768, tablet 768-1279, desktop 1280+.</p></div>`)
	sb.WriteString(`</div><div class="grid2">`)
	sb.WriteString(`<div class="panel"><h2>OPERATING SYSTEM</h2>`)
	if len(dev.OS) == 0 {
		sb.WriteString(`<p class="empty">Nothing yet.</p>`)
	} else {
		for i, o := range dev.OS {
			sb.WriteString(hbarRow(o.Name, o.Count, total, chartPalette[i%len(chartPalette)]))
		}
	}
	sb.WriteString(`</div><div class="panel"><h2>BROWSER</h2>`)
	if len(dev.Browsers) == 0 {
		sb.WriteString(`<p class="empty">Nothing yet.</p>`)
	} else {
		for i, b := range dev.Browsers {
			sb.WriteString(hbarRow(b.Name, b.Count, total, chartPalette[i%len(chartPalette)]))
		}
	}
	sb.WriteString(`</div></div>`)
	sb.WriteString(`<div class="panel"><h2>TOP SCREEN RESOLUTIONS</h2>`)
	if len(dev.Resolutions) == 0 {
		sb.WriteString(`<p class="empty">No screen data yet. New pageviews carry it.</p>`)
	} else {
		sb.WriteString(`<table class="sortable"><thead><tr><th class="s">RESOLUTION</th><th class="s num" data-num>VISITS</th><th class="s num" data-num>SHARE</th><th style="width:40%">BAR</th></tr></thead><tbody>`)
		var rt int64
		for _, x := range dev.Resolutions {
			rt += x.Count
		}
		for _, x := range dev.Resolutions {
			wPct := pctOf(x.Count, rt)
			sb.WriteString(`<tr><td data-v="` + esc(x.Name) + `">` + esc(x.Name) +
				`</td><td class="num" data-v="` + strconv.FormatInt(x.Count, 10) + `">` + fnum(x.Count) +
				`</td><td class="num">` + wPct + `</td><td><div class="bar"><i style="width:` + wPct + `"></i></div></td></tr>`)
		}
		sb.WriteString(`</tbody></table>`)
	}
	sb.WriteString(`</div>`)
	return sb.String(), nil
}

// ---- GEOGRAPHY ----

func (m *Module) secGeo(r *http.Request, db *sql.DB, days int, now int64) (string, error) {
	country := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("country")))
	var sb strings.Builder
	if country != "" {
		regions, cities, err := data.QueryCountryDrill(db, days, now, country)
		if err != nil {
			return "", err
		}
		sb.WriteString(sectionHead("Geography: "+country, "Region and city drill-down. Coarse IP geo only, raw IPs are never stored."))
		sb.WriteString(`<p><a href="/admin?s=geo&days=` + strconv.Itoa(days) + `" style="color:var(--khaki)">&larr; All countries</a></p>`)
		sb.WriteString(`<div class="grid2"><div class="panel"><h2>REGIONS</h2>`)
		if len(regions) == 0 {
			sb.WriteString(`<p class="empty">Nothing yet.</p>`)
		} else {
			sb.WriteString(`<table class="sortable"><thead><tr><th class="s">REGION</th><th class="s num" data-num>VIEWS</th><th class="s num" data-num>UNIQUES</th></tr></thead><tbody>`)
			for _, x := range regions {
				sb.WriteString(`<tr><td data-v="` + esc(x.Name) + `">` + esc(x.Name) +
					`</td><td class="num" data-v="` + strconv.FormatInt(x.Views, 10) + `">` + fnum(x.Views) +
					`</td><td class="num" data-v="` + strconv.FormatInt(x.Uniques, 10) + `">` + fnum(x.Uniques) + `</td></tr>`)
			}
			sb.WriteString(`</tbody></table>`)
		}
		sb.WriteString(`</div><div class="panel"><h2>CITIES</h2>`)
		if len(cities) == 0 {
			sb.WriteString(`<p class="empty">Nothing yet.</p>`)
		} else {
			sb.WriteString(`<table class="sortable"><thead><tr><th class="s">CITY</th><th class="s num" data-num>VIEWS</th><th class="s num" data-num>UNIQUES</th></tr></thead><tbody>`)
			for _, x := range cities {
				sb.WriteString(`<tr><td data-v="` + esc(x.Name) + `">` + esc(x.Name) +
					`</td><td class="num" data-v="` + strconv.FormatInt(x.Views, 10) + `">` + fnum(x.Views) +
					`</td><td class="num" data-v="` + strconv.FormatInt(x.Uniques, 10) + `">` + fnum(x.Uniques) + `</td></tr>`)
			}
			sb.WriteString(`</tbody></table>`)
		}
		sb.WriteString(`</div></div>`)
		return sb.String(), nil
	}
	countries, err := data.QueryCountries(db, days, now)
	if err != nil {
		return "", err
	}
	sb.WriteString(sectionHead("Geography", "Where in the world. Coarse IP geolocation (country/region/city). Raw IPs are never stored. Click a country to drill down."))
	sb.WriteString(`<div class="panel"><h2>COUNTRIES</h2>`)
	if len(countries) == 0 {
		sb.WriteString(`<p class="empty">Nothing yet.</p></div>`)
		return sb.String(), nil
	}
	var total int64
	for _, c := range countries {
		total += c.Views
	}
	sb.WriteString(`<div style="overflow-x:auto"><table class="sortable"><thead><tr>` +
		`<th class="s">COUNTRY</th><th class="s num" data-num>VIEWS</th><th class="s num" data-num>UNIQUES</th>` +
		`<th class="s num" data-num>SHARE</th><th style="width:30%">BAR</th></tr></thead><tbody>`)
	for _, c := range countries {
		name := c.Name
		if name == "" || name == "(unknown)" {
			name = "Unknown"
		} else {
			name = name + " (" + c.Code + ")"
		}
		wPct := pctOf(c.Views, total)
		link := "/admin?s=geo&days=" + strconv.Itoa(days) + "&country=" + esc(c.Code)
		sb.WriteString(`<tr>`)
		sb.WriteString(`<td data-v="` + esc(name) + `"><a href="` + link + `" style="color:var(--khaki)">` + esc(name) + ` &rarr;</a></td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(c.Views, 10) + `">` + fnum(c.Views) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(c.Uniques, 10) + `">` + fnum(c.Uniques) + `</td>`)
		sb.WriteString(`<td class="num">` + wPct + `</td>`)
		sb.WriteString(`<td><div class="bar"><i style="width:` + wPct + `"></i></div></td></tr>`)
	}
	sb.WriteString(`</tbody></table></div></div>`)
	return sb.String(), nil
}

// ---- VISITORS ----

func (m *Module) secVisitors(db *sql.DB, days int, now int64) (string, error) {
	v, err := data.QueryVisitors(db, days, now)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(sectionHead("Visitors", "Who comes back. Retention cohorts and the most active anonymous browsers (truncated IDs, no PII)."))
	sb.WriteString(`<div class="cards">`)
	sb.WriteString(statCard("BACK WITHIN 7 DAYS", fnum(v.Return7d)+" ("+fpct(v.Return7dPct)+")", "of visitors old enough to measure"))
	sb.WriteString(statCard("DAY-1 RETURN", fpct(v.Cohort1d), "came back the next day"))
	sb.WriteString(statCard("DAY-7 RETURN", fpct(v.Cohort7d), "came back within a week"))
	sb.WriteString(statCard("DAY-30 RETURN", fpct(v.Cohort30d), "came back within a month"))
	sb.WriteString(`</div>`)
	// Stacked new vs returning chart.
	var labels []string
	var nv, rv []float64
	for i := range v.NewTrend {
		labels = append(labels, v.NewTrend[i].Date)
		nv = append(nv, float64(v.NewTrend[i].Count))
		rv = append(rv, float64(v.ReturningTrend[i].Count))
	}
	sb.WriteString(`<div class="panel"><h2>NEW VS RETURNING PER DAY</h2>` +
		svgStacked(labels, nv, rv, "New", "Returning", 900, 170) +
		`<div class="legend"><span><i style="background:#5c6f3a"></i>New</span><span><i style="background:#f4772e"></i>Returning</span></div>`)
	if len(labels) > 0 {
		sb.WriteString(`<div class="chart-x"><span>` + esc(labels[0]) + `</span><span>` + esc(labels[len(labels)-1]) + `</span></div>`)
	}
	sb.WriteString(`</div>`)
	sb.WriteString(`<div class="panel"><h2>MOST ACTIVE VISITORS</h2>`)
	if len(v.TopVisitors) == 0 {
		sb.WriteString(`<p class="empty">Nothing yet.</p>`)
	} else {
		sb.WriteString(`<table class="sortable"><thead><tr><th class="s">VISITOR</th><th class="s num" data-num>VIEWS</th><th class="s num" data-num>SESSIONS</th><th class="s" data-num>LAST SEEN</th></tr></thead><tbody>`)
		for _, t := range v.TopVisitors {
			sb.WriteString(`<tr><td class="path" data-v="` + esc(t.IDPrefix) + `">` + esc(t.IDPrefix) + `&hellip;</td>` +
				`<td class="num" data-v="` + strconv.FormatInt(t.Views, 10) + `">` + fnum(t.Views) + `</td>` +
				`<td class="num" data-v="` + strconv.FormatInt(t.Sessions, 10) + `">` + fnum(t.Sessions) + `</td>` +
				`<td data-v="` + strconv.FormatInt(t.LastSeen, 10) + `">` + esc(ago(t.LastSeen, now)) + `</td></tr>`)
		}
		sb.WriteString(`</tbody></table><p class="muted">Anonymous IDs truncated to 8 chars. No personal data.</p>`)
	}
	sb.WriteString(`</div>`)
	return sb.String(), nil
}

// ---- ACTIVITY ----

func (m *Module) secActivity(r *http.Request, now int64) (string, error) {
	db, err := data.AnalyticsDB(m.DBPath)
	if err != nil {
		return "", err
	}
	kind := r.URL.Query().Get("kind")
	if kind != "pageview" && kind != "click" {
		kind = ""
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	feed, err := data.QueryFeed(db, 40, kind, q)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(sectionHead("Activity", "Live event stream. Toggle auto-refresh to watch visits land in real time."))
	sb.WriteString(`<form class="filters" id="feedform">` +
		`<div><label>EVENT TYPE</label><select id="feedkind" name="kind">` +
		`<option value="">All</option>` +
		`<option value="pageview"` + map[bool]string{true: " selected", false: ""}[kind == "pageview"] + `>Pageviews</option>` +
		`<option value="click"` + map[bool]string{true: " selected", false: ""}[kind == "click"] + `>Clicks</option>` +
		`</select></div>` +
		`<div><label>PAGE CONTAINS</label><input id="feedq" name="q" value="` + esc(q) + `" placeholder="/services"></div>` +
		`<div><button type="submit">APPLY</button></div>` +
		`<div class="toggle"><input type="checkbox" id="feedauto"><label for="feedauto" style="margin:0">AUTO-REFRESH 5s</label></div>` +
		`</form>`)
	sb.WriteString(`<div class="panel"><h2>LIVE FEED</h2><div style="overflow-x:auto"><table><thead><tr>` +
		`<th>WHEN</th><th>WHAT</th><th>PAGE</th><th>DETAIL</th><th>DEVICE</th><th>SOURCE</th><th>GEO</th>` +
		`</tr></thead><tbody id="feedbody">`)
	sb.WriteString(feedRows(feed, now))
	sb.WriteString(`</tbody></table></div></div>`)
	return sb.String(), nil
}

// feedRows renders activity rows (shared by the section and the AJAX feed).
func feedRows(feed []data.FeedHit, now int64) string {
	var sb strings.Builder
	for _, h := range feed {
		detail := h.Label
		if h.Kind == "click" && h.Href != "" {
			if detail != "" {
				detail += " -> " + h.Href
			} else {
				detail = h.Href
			}
		}
		devCls := ""
		devLabel := strings.ToUpper(h.DeviceType)
		switch h.DeviceType {
		case "phone":
			devCls = "phone"
		case "desktop":
			devCls = "desktop"
		}
		sb.WriteString(`<tr>`)
		sb.WriteString(`<td class="num" style="white-space:nowrap">` + esc(ago(h.TS, now)) + `</td>`)
		sb.WriteString(`<td><span class="badge">` + esc(h.Kind) + `</span></td>`)
		sb.WriteString(`<td class="path">` + esc(h.Path) + `</td>`)
		sb.WriteString(`<td>` + esc(detail) + `</td>`)
		sb.WriteString(`<td><span class="badge ` + devCls + `">` + esc(devLabel) + `</span> <span class="muted">` + esc(h.Browser) + `</span></td>`)
		sb.WriteString(`<td>` + esc(h.Source) + `</td>`)
		sb.WriteString(`<td>` + esc(h.Country) + `</td>`)
		sb.WriteString(`</tr>`)
	}
	if len(feed) == 0 {
		sb.WriteString(`<tr><td colspan="7" class="empty">No events match.</td></tr>`)
	}
	return sb.String()
}

// adminFeed serves the activity table rows for auto-refresh (HTML fragment).
func (m *Module) adminFeed(w http.ResponseWriter, r *http.Request) {
	kind := r.URL.Query().Get("kind")
	if kind != "pageview" && kind != "click" {
		kind = ""
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	db := m.adminDB(w)
	if db == nil {
		return
	}
	feed, err := data.QueryFeed(db, 40, kind, q)
	if err != nil {
		slog.Error("admin: feed query failed", "err", err.Error())
		http.Error(w, "feed unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	_, _ = w.Write([]byte(feedRows(feed, time.Now().Unix())))
}
