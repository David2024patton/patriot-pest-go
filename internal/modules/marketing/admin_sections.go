// Dashboard section renderers, part 1: overview, pages, clicks, traffic.
package marketing

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// ---- OVERVIEW ----

func (m *Module) secOverview(db *sql.DB, days int, now int64) (string, error) {
	st, err := data.QueryAnalytics(db, days, now)
	if err != nil {
		return "", err
	}
	ss, err := data.QuerySessionStats(db, days, now)
	if err != nil {
		return "", err
	}
	dev, err := data.QueryDevices(db, days, now)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(sectionHead("Overview", "The whole funnel at a glance. Last "+strconv.Itoa(days)+" days."))
	sb.WriteString(`<div class="cards">`)
	sb.WriteString(statCard("PAGEVIEWS", fnum(st.Pageviews), "pages loaded"))
	sb.WriteString(statCard("UNIQUE VISITORS", fnum(st.Visitors), "distinct browsers"))
	sb.WriteString(statCard("SESSIONS", fnum(ss.Sessions), "visits, 30 min timeout"))
	sb.WriteString(statCard("PAGES / SESSION", fmt.Sprintf("%.2f", ss.PagesPerSess), "depth per visit"))
	sb.WriteString(statCard("AVG SESSION", fdur(ss.AvgDurationSec), "mean visit length"))
	sb.WriteString(statCard("BOUNCE RATE", fpct(ss.BounceRate), "single-page sessions"))
	sb.WriteString(statCard("NEW VISITORS", fnum(st.NewVisitors), "first seen in range"))
	sb.WriteString(statCard("RETURNING", fnum(st.ReturningVisitors), "seen before range"))
	sb.WriteString(`</div>`)

	// Daily pageviews.
	var items []barItem
	for _, d := range st.Daily {
		items = append(items, barItem{Label: d.Date, Value: float64(d.Count)})
	}
	sb.WriteString(`<div class="panel"><h2>DAILY PAGEVIEWS</h2>` + svgBars(items, 900, 160, "#f4772e"))
	if len(st.Daily) > 0 {
		sb.WriteString(`<div class="chart-x"><span>` + esc(st.Daily[0].Date) + `</span><span>` +
			esc(st.Daily[len(st.Daily)-1].Date) + `</span></div>`)
	}
	sb.WriteString(`</div>`)

	sb.WriteString(`<div class="grid2">`)
	// Top pages mini list.
	sb.WriteString(`<div class="panel"><h2>TOP PAGES</h2>`)
	if len(st.TopPages) == 0 {
		sb.WriteString(`<p class="empty">Nothing yet.</p>`)
	} else {
		sb.WriteString(`<table>`)
		n := len(st.TopPages)
		if n > 8 {
			n = 8
		}
		for _, p := range st.TopPages[:n] {
			sb.WriteString(`<tr><td class="path">` + esc(p.Path) + `</td><td class="num">` + fnum(p.Count) + `</td></tr>`)
		}
		sb.WriteString(`</table><p class="muted"><a href="/admin?s=pages&days=` + strconv.Itoa(days) +
			`" style="color:var(--khaki)">Full page report &rarr;</a></p>`)
	}
	sb.WriteString(`</div>`)
	// Top sources mini list.
	sb.WriteString(`<div class="panel"><h2>TOP SOURCES</h2>`)
	if len(st.TopSources) == 0 {
		sb.WriteString(`<p class="empty">Nothing yet.</p>`)
	} else {
		sb.WriteString(`<table>`)
		n := len(st.TopSources)
		if n > 8 {
			n = 8
		}
		for _, s := range st.TopSources[:n] {
			sb.WriteString(`<tr><td>` + esc(s.Source) + `</td><td class="num">` + fnum(s.Count) + `</td></tr>`)
		}
		sb.WriteString(`</table><p class="muted"><a href="/admin?s=traffic&days=` + strconv.Itoa(days) +
			`" style="color:var(--khaki)">Full traffic report &rarr;</a></p>`)
	}
	sb.WriteString(`</div></div>`)

	// Device split donut.
	sb.WriteString(`<div class="panel"><h2>DEVICE SPLIT</h2><div style="display:flex;gap:1.5rem;align-items:center;flex-wrap:wrap">`)
	sb.WriteString(svgDonut(dev.Types, 170))
	sb.WriteString(`<div style="flex:1;min-width:220px">` + donutLegend(dev.Types) +
		`<p class="muted"><a href="/admin?s=devices&days=` + strconv.Itoa(days) +
		`" style="color:var(--khaki)">Full device report &rarr;</a></p></div>`)
	sb.WriteString(`</div></div>`)
	return sb.String(), nil
}

// ---- PAGES ----

func (m *Module) secPages(db *sql.DB, days int, now int64) (string, error) {
	pages, err := data.QueryPageStats(db, days, now)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(sectionHead("Pages", "What people read. Entrances = first page of a visit, exits = last. Time on page measured between pageviews in a session."))
	sb.WriteString(`<div class="panel"><h2>ALL PAGES &middot; CLICK A HEADER TO SORT</h2>`)
	if len(pages) == 0 {
		sb.WriteString(`<p class="empty">No pageviews in this range.</p></div>`)
		return sb.String(), nil
	}
	sb.WriteString(`<div style="overflow-x:auto"><table class="sortable"><thead><tr>` +
		`<th class="s">PAGE</th><th class="s num" data-num>VIEWS</th><th class="s num" data-num>UNIQUES</th>` +
		`<th class="s num" data-num>ENTRANCES</th><th class="s num" data-num>EXITS</th>` +
		`<th class="s num" data-num>AVG TIME</th><th>TREND</th></tr></thead><tbody>`)
	for _, p := range pages {
		t := "-"
		tv := -1.0
		if p.AvgTimeSec > 0 {
			t = fdur(p.AvgTimeSec)
			tv = p.AvgTimeSec
		}
		tip := p.Path + "\nViews: " + fnum(p.Views) + "\nUniques: " + fnum(p.Uniques) +
			"\nEntrances: " + fnum(p.Entrances) + "\nExits: " + fnum(p.Exits)
		sb.WriteString(`<tr data-tip="` + esc(tip) + `">`)
		sb.WriteString(`<td class="path" data-v="` + esc(p.Path) + `">` + esc(p.Path) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(p.Views, 10) + `">` + fnum(p.Views) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(p.Uniques, 10) + `">` + fnum(p.Uniques) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(p.Entrances, 10) + `">` + fnum(p.Entrances) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(p.Exits, 10) + `">` + fnum(p.Exits) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatFloat(tv, 'f', 1, 64) + `">` + t + `</td>`)
		sb.WriteString(`<td>` + svgSpark(p.SparkDaily, 90, 24, "#f4772e") + `</td>`)
		sb.WriteString(`</tr>`)
	}
	sb.WriteString(`</tbody></table></div></div>`)
	return sb.String(), nil
}

// ---- CLICKS ----

func (m *Module) secClicks(r *http.Request, db *sql.DB, days int, now int64) (string, error) {
	pageFilter := strings.TrimSpace(r.URL.Query().Get("page"))
	clicks, err := data.QueryClicks(db, days, now, pageFilter)
	if err != nil {
		return "", err
	}
	// Page options for the filter dropdown.
	pages, err := data.QueryPageStats(db, days, now)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(sectionHead("Clicks", "What visitors actually tapped. Filter by page to see exactly what got clicked where."))
	sb.WriteString(`<form class="filters" method="get" action="/admin">` +
		`<input type="hidden" name="s" value="clicks"><input type="hidden" name="days" value="` + strconv.Itoa(days) + `">` +
		`<div><label>PAGE</label><select name="page"><option value="">All pages</option>`)
	for _, p := range pages {
		sel := ""
		if p.Path == pageFilter {
			sel = " selected"
		}
		sb.WriteString(`<option value="` + esc(p.Path) + `"` + sel + `>` + esc(p.Path) + ` (` + fnum(p.Views) + ` views)</option>`)
	}
	sb.WriteString(`</select></div><div><button type="submit">FILTER</button></div></form>`)
	sb.WriteString(`<div class="panel"><h2>CLICKED ELEMENTS` + (map[bool]string{true: " &middot; " + esc(pageFilter), false: ""}[pageFilter != ""]) + `</h2>`)
	if len(clicks) == 0 {
		sb.WriteString(`<p class="empty">No clicks tracked yet.</p></div>`)
		return sb.String(), nil
	}
	sb.WriteString(`<div style="overflow-x:auto"><table class="sortable"><thead><tr>` +
		`<th class="s">CLICKED</th><th class="s">LINK TARGET</th><th class="s">PAGE</th>` +
		`<th class="s num" data-num>CLICKS</th><th class="s num" data-num>CLICKERS</th><th class="s">ELEMENT</th>` +
		`</tr></thead><tbody>`)
	for _, c := range clicks {
		label := c.Label
		if label == "" {
			label = "(no label)"
		}
		href := c.Href
		if href == "" {
			href = "-"
		}
		sb.WriteString(`<tr>`)
		sb.WriteString(`<td data-v="` + esc(label) + `">` + esc(label) + `</td>`)
		sb.WriteString(`<td class="path" data-v="` + esc(href) + `">` + esc(href) + `</td>`)
		sb.WriteString(`<td class="path" data-v="` + esc(c.Path) + `">` + esc(c.Path) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(c.Clicks, 10) + `">` + fnum(c.Clicks) + `</td>`)
		sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(c.Clickers, 10) + `">` + fnum(c.Clickers) + `</td>`)
		sb.WriteString(`<td data-v="` + esc(c.ElemHint) + `"><span class="badge">` + esc(c.ElemHint) + `</span></td>`)
		sb.WriteString(`</tr>`)
	}
	sb.WriteString(`</tbody></table></div></div>`)
	return sb.String(), nil
}

// ---- TRAFFIC ----

func (m *Module) secTraffic(r *http.Request, db *sql.DB, days int, now int64) (string, error) {
	sources, err := data.QuerySources(db, days, now)
	if err != nil {
		return "", err
	}
	var total int64
	for _, s := range sources {
		total += s.Views
	}
	var sb strings.Builder
	sb.WriteString(sectionHead("Traffic", "Where visitors come from. Click a source row to expand its UTM campaigns."))
	sb.WriteString(`<div class="panel"><h2>SOURCES &middot; CLICK A ROW FOR CAMPAIGNS</h2>`)
	if len(sources) == 0 {
		sb.WriteString(`<p class="empty">Nothing yet.</p></div>`)
	} else {
		sb.WriteString(`<div style="overflow-x:auto"><table class="sortable"><thead><tr>` +
			`<th class="s">SOURCE</th><th class="s num" data-num>VIEWS</th><th class="s num" data-num>UNIQUES</th>` +
			`<th class="s num" data-num>SHARE</th><th style="width:30%">BAR</th></tr></thead><tbody>`)
		for i, s := range sources {
			if i >= 12 {
				break
			}
			camps, err := data.QueryCampaigns(db, days, now, s.Source)
			if err != nil {
				return "", err
			}
			drillID := "camp" + strconv.Itoa(i)
			wPct := pctOf(s.Views, total)
			sb.WriteString(`<tr class="clickable" data-drill="` + drillID + `" data-tip="Click to expand UTM campaigns">`)
			sb.WriteString(`<td data-v="` + esc(s.Source) + `"><strong>` + esc(s.Source) + `</strong> <span style="color:var(--orange)">+</span></td>`)
			sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(s.Views, 10) + `">` + fnum(s.Views) + `</td>`)
			sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(s.Uniques, 10) + `">` + fnum(s.Uniques) + `</td>`)
			sb.WriteString(`<td class="num" data-v="` + strconv.FormatInt(s.Views, 10) + `">` + wPct + `</td>`)
			sb.WriteString(`<td><div class="bar"><i style="width:` + wPct + `"></i></div></td></tr>`)
			sb.WriteString(`<tr class="drill" id="` + drillID + `"><td colspan="5">`)
			if len(camps) == 0 {
				sb.WriteString(`<p class="empty">No campaign data.</p>`)
			} else {
				sb.WriteString(`<table><tr><th>CAMPAIGN</th><th class="num">VIEWS</th><th class="num">UNIQUES</th></tr>`)
				for _, c := range camps {
					sb.WriteString(`<tr><td>` + esc(c.Campaign) + `</td><td class="num">` + fnum(c.Views) +
						`</td><td class="num">` + fnum(c.Uniques) + `</td></tr>`)
				}
				sb.WriteString(`</table>`)
			}
			sb.WriteString(`</td></tr>`)
		}
		sb.WriteString(`</tbody></table></div>`)
		sb.WriteString(`<p class="muted">UTM parameters on ad links override referrer detection. Untagged links fall back to the referrer domain.</p></div>`)
	}

	refs, err := data.QueryReferrers(db, days, now, r.Host)
	if err != nil {
		return "", err
	}
	sb.WriteString(`<div class="panel"><h2>REFERRER DOMAINS</h2>`)
	if len(refs) == 0 {
		sb.WriteString(`<p class="empty">No external referrers in this range.</p>`)
	} else {
		sb.WriteString(`<table class="sortable"><thead><tr><th class="s">DOMAIN</th><th class="s num" data-num>VIEWS</th><th style="width:40%">BAR</th></tr></thead><tbody>`)
		var rTotal int64
		for _, x := range refs {
			rTotal += x.Views
		}
		for _, x := range refs {
			wPct := pctOf(x.Views, rTotal)
			sb.WriteString(`<tr><td class="path" data-v="` + esc(x.Domain) + `">` + esc(x.Domain) +
				`</td><td class="num" data-v="` + strconv.FormatInt(x.Views, 10) + `">` + fnum(x.Views) +
				`</td><td><div class="bar"><i style="width:` + wPct + `"></i></div></td></tr>`)
		}
		sb.WriteString(`</tbody></table>`)
	}
	sb.WriteString(`</div>`)
	return sb.String(), nil
}
