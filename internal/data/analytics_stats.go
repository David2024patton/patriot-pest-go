// Extended analytics aggregates for the mission-control dashboard:
// sessions, pages, clicks, devices, geo, visitor retention. All functions
// take a day range and a unix "now" so tests can pin time.
package data

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

// normDays pins the day range to the supported set.
func normDays(days int) int {
	switch days {
	case 1, 7, 30, 90:
		return days
	}
	return 30
}

func rangeStart(days int, now int64) int64 {
	return now - int64(normDays(days))*86400
}

// ---- sessions ----

// SessionStats holds engagement numbers for the range.
type SessionStats struct {
	Sessions       int64
	BounceRate     float64 // single-pageview sessions / sessions
	AvgDurationSec float64 // mean session length over pageviews
	PagesPerSess   float64
}

// QuerySessionStats computes bounce rate, average session duration and
// pages per session from pageview events in the range.
func QuerySessionStats(db *sql.DB, days int, now int64) (*SessionStats, error) {
	start := rangeStart(days, now)
	rows, err := db.Query(`SELECT session_id, COUNT(*), MIN(ts), MAX(ts)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND session_id<>''
		GROUP BY session_id`, start, now)
	if err != nil {
		return nil, fmt.Errorf("session stats: %w", err)
	}
	defer rows.Close()
	var n, bounces, totalViews int64
	var totalDur int64
	for rows.Next() {
		var sid string
		var views, mn, mx int64
		if err := rows.Scan(&sid, &views, &mn, &mx); err != nil {
			return nil, fmt.Errorf("session stats scan: %w", err)
		}
		n++
		totalViews += views
		if views == 1 {
			bounces++
		}
		if mx > mn {
			totalDur += mx - mn
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("session stats: %w", err)
	}
	st := &SessionStats{Sessions: n}
	if n > 0 {
		st.BounceRate = float64(bounces) / float64(n)
		st.AvgDurationSec = float64(totalDur) / float64(n)
		st.PagesPerSess = float64(totalViews) / float64(n)
	}
	return st, nil
}

// ---- pages ----

// PageStat is one page's funnel numbers for the range.
type PageStat struct {
	Path        string
	Views       int64
	Uniques     int64
	Entrances   int64 // sessions whose first pageview was here
	Exits       int64 // sessions whose last pageview was here
	AvgTimeSec  float64
	SparkDaily  []int64 // daily view counts for the sparkline
	SparkLabels []string
}

// QueryPageStats returns per-page stats for the range, most viewed first.
func QueryPageStats(db *sql.DB, days int, now int64) ([]PageStat, error) {
	start := rangeStart(days, now)
	rows, err := db.Query(`SELECT path, COUNT(*), COUNT(DISTINCT visitor_id)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=?
		GROUP BY path ORDER BY COUNT(*) DESC LIMIT 50`, start, now)
	if err != nil {
		return nil, fmt.Errorf("page stats: %w", err)
	}
	var pages []PageStat
	for rows.Next() {
		var p PageStat
		if err := rows.Scan(&p.Path, &p.Views, &p.Uniques); err != nil {
			rows.Close()
			return nil, fmt.Errorf("page stats scan: %w", err)
		}
		if p.Path == "" {
			p.Path = "/"
		}
		pages = append(pages, p)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("page stats: %w", err)
	}
	rows.Close()

	// Entrances and exits per session: one pass each.
	enterRows, err := db.Query(`SELECT path, COUNT(*) FROM (
			SELECT session_id, path FROM analytics_events a
			WHERE kind='pageview' AND ts>=? AND ts<=? AND session_id<>''
			AND ts = (SELECT MIN(ts) FROM analytics_events
				WHERE kind='pageview' AND session_id=a.session_id AND ts>=? AND ts<=?)
		) GROUP BY path`, start, now, start, now)
	if err != nil {
		return nil, fmt.Errorf("page entrances: %w", err)
	}
	entrances := scanCountMap(enterRows)
	exitRows, err := db.Query(`SELECT path, COUNT(*) FROM (
			SELECT session_id, path FROM analytics_events a
			WHERE kind='pageview' AND ts>=? AND ts<=? AND session_id<>''
			AND ts = (SELECT MAX(ts) FROM analytics_events
				WHERE kind='pageview' AND session_id=a.session_id AND ts>=? AND ts<=?)
		) GROUP BY path`, start, now, start, now)
	if err != nil {
		return nil, fmt.Errorf("page exits: %w", err)
	}
	exits := scanCountMap(exitRows)

	// Average time on page: diff between a pageview and the next pageview
	// in the same session (window function), 1s..30min only, averaged.
	timeRows, err := db.Query(`SELECT path, AVG(d) FROM (
			SELECT path, LEAD(ts) OVER (PARTITION BY session_id ORDER BY ts) - ts AS d
			FROM analytics_events
			WHERE kind='pageview' AND ts>=? AND ts<=? AND session_id<>''
		) WHERE d >= 1 AND d <= 1800 GROUP BY path`, start, now)
	if err != nil {
		return nil, fmt.Errorf("page time: %w", err)
	}
	times := map[string]float64{}
	for timeRows.Next() {
		var path string
		var avg float64
		if err := timeRows.Scan(&path, &avg); err != nil {
			timeRows.Close()
			return nil, fmt.Errorf("page time scan: %w", err)
		}
		times[path] = avg
	}
	if err := timeRows.Err(); err != nil {
		timeRows.Close()
		return nil, fmt.Errorf("page time: %w", err)
	}
	timeRows.Close()

	// Daily sparklines for the top pages only (keeps the query cheap).
	dayStart := start - (start % 86400)
	nDays := int((now-dayStart)/86400) + 1
	for i := range pages {
		if i >= 12 {
			break
		}
		sr, err := db.Query(`SELECT CAST(ts/86400 AS INTEGER)*86400, COUNT(*)
			FROM analytics_events
			WHERE kind='pageview' AND path=? AND ts>=? AND ts<=?
			GROUP BY 1 ORDER BY 1`, pages[i].Path, start, now)
		if err != nil {
			return nil, fmt.Errorf("page spark: %w", err)
		}
		m := scanCountMapInt(sr)
		pages[i].SparkDaily = make([]int64, nDays)
		pages[i].SparkLabels = make([]string, nDays)
		for d := 0; d < nDays; d++ {
			ds := dayStart + int64(d)*86400
			pages[i].SparkDaily[d] = m[ds]
			pages[i].SparkLabels[d] = dayLabel(ds)
		}
		pages[i].Entrances = entrances[pages[i].Path]
		pages[i].Exits = exits[pages[i].Path]
		pages[i].AvgTimeSec = times[pages[i].Path]
	}
	for i := range pages {
		if pages[i].SparkDaily == nil {
			pages[i].Entrances = entrances[pages[i].Path]
			pages[i].Exits = exits[pages[i].Path]
			pages[i].AvgTimeSec = times[pages[i].Path]
		}
	}
	return pages, nil
}

func scanCountMap(rows *sql.Rows) map[string]int64 {
	defer rows.Close()
	m := map[string]int64{}
	for rows.Next() {
		var k string
		var v int64
		if err := rows.Scan(&k, &v); err != nil {
			break
		}
		if k == "" {
			k = "/"
		}
		m[k] = v
	}
	return m
}

func scanCountMapInt(rows *sql.Rows) map[int64]int64 {
	defer rows.Close()
	m := map[int64]int64{}
	for rows.Next() {
		var k, v int64
		if err := rows.Scan(&k, &v); err != nil {
			break
		}
		m[k] = v
	}
	return m
}

// ---- clicks ----

// ClickDetail is one clicked element aggregated for the range.
type ClickDetail struct {
	Label    string
	Href     string
	Path     string
	ElemHint string
	Clicks   int64
	Clickers int64 // distinct visitors
}

// QueryClicks aggregates clicks, optionally filtered to one page. Empty
// pathFilter means all pages.
func QueryClicks(db *sql.DB, days int, now int64, pathFilter string) ([]ClickDetail, error) {
	start := rangeStart(days, now)
	q := `SELECT label, href, path, elem_hint, COUNT(*), COUNT(DISTINCT visitor_id)
		FROM analytics_events
		WHERE kind='click' AND ts>=? AND ts<=?`
	args := []any{start, now}
	if pathFilter != "" {
		q += ` AND path=?`
		args = append(args, pathFilter)
	}
	q += ` GROUP BY label, href, path, elem_hint ORDER BY COUNT(*) DESC LIMIT 100`
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("clicks: %w", err)
	}
	defer rows.Close()
	var out []ClickDetail
	for rows.Next() {
		var c ClickDetail
		if err := rows.Scan(&c.Label, &c.Href, &c.Path, &c.ElemHint, &c.Clicks, &c.Clickers); err != nil {
			return nil, fmt.Errorf("clicks scan: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ---- traffic ----

// SourceStat is one traffic source for the range.
type SourceStat struct {
	Source  string
	Views   int64
	Uniques int64
}

// QuerySources aggregates pageviews per attributed source.
func QuerySources(db *sql.DB, days int, now int64) ([]SourceStat, error) {
	start := rangeStart(days, now)
	rows, err := db.Query(`SELECT COALESCE(NULLIF(source,''),'Direct'), COUNT(*), COUNT(DISTINCT visitor_id)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=?
		GROUP BY 1 ORDER BY COUNT(*) DESC LIMIT 30`, start, now)
	if err != nil {
		return nil, fmt.Errorf("sources: %w", err)
	}
	defer rows.Close()
	var out []SourceStat
	for rows.Next() {
		var s SourceStat
		if err := rows.Scan(&s.Source, &s.Views, &s.Uniques); err != nil {
			return nil, fmt.Errorf("sources scan: %w", err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// CampaignStat is one UTM campaign under a source.
type CampaignStat struct {
	Campaign string
	Views    int64
	Uniques  int64
}

// QueryCampaigns breaks one source down by utm_campaign.
func QueryCampaigns(db *sql.DB, days int, now int64, source string) ([]CampaignStat, error) {
	start := rangeStart(days, now)
	rows, err := db.Query(`SELECT COALESCE(NULLIF(campaign,''),'(none)'), COUNT(*), COUNT(DISTINCT visitor_id)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND COALESCE(NULLIF(source,''),'Direct')=?
		GROUP BY 1 ORDER BY COUNT(*) DESC LIMIT 20`, start, now, source)
	if err != nil {
		return nil, fmt.Errorf("campaigns: %w", err)
	}
	defer rows.Close()
	var out []CampaignStat
	for rows.Next() {
		var c CampaignStat
		if err := rows.Scan(&c.Campaign, &c.Views, &c.Uniques); err != nil {
			return nil, fmt.Errorf("campaigns scan: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ReferrerStat is one external referrer domain.
type ReferrerStat struct {
	Domain string
	Views  int64
}

// QueryReferrers aggregates distinct referrer hosts (excludes our own
// domain and empty referrers).
func QueryReferrers(db *sql.DB, days int, now int64, ownHost string) ([]ReferrerStat, error) {
	start := rangeStart(days, now)
	rows, err := db.Query(`SELECT referrer, COUNT(*)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND referrer<>''
		GROUP BY referrer ORDER BY COUNT(*) DESC LIMIT 200`, start, now)
	if err != nil {
		return nil, fmt.Errorf("referrers: %w", err)
	}
	defer rows.Close()
	agg := map[string]int64{}
	for rows.Next() {
		var ref string
		var n int64
		if err := rows.Scan(&ref, &n); err != nil {
			return nil, fmt.Errorf("referrers scan: %w", err)
		}
		host := ""
		if u, err := url.Parse(ref); err == nil {
			host = strings.ToLower(u.Hostname())
		}
		if host == "" || host == strings.ToLower(ownHost) {
			continue
		}
		agg[host] += n
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("referrers: %w", err)
	}
	var out []ReferrerStat
	for d, n := range agg {
		out = append(out, ReferrerStat{Domain: d, Views: n})
	}
	// Small result set: sort in Go.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Views > out[i].Views {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if len(out) > 30 {
		out = out[:30]
	}
	return out, nil
}

// ---- devices ----

// DeviceStats holds the device/OS/browser breakdown for the range.
type DeviceStats struct {
	Types       []NameCount // phone | tablet | desktop
	OS          []NameCount
	Browsers    []NameCount
	Resolutions []NameCount // top raw WxH strings
	Buckets     []NameCount // Mobile <768 | Tablet 768-1279 | Desktop 1280+
}

// NameCount is a labeled counter for bar charts.
type NameCount struct {
	Name  string
	Count int64
}

func groupCount(db *sql.DB, col string, start, now int64, limit int) ([]NameCount, error) {
	rows, err := db.Query(`SELECT COALESCE(NULLIF(`+col+`,''),'Other'), COUNT(*)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=?
		GROUP BY 1 ORDER BY COUNT(*) DESC LIMIT ?`, start, now, limit)
	if err != nil {
		return nil, fmt.Errorf("device group %s: %w", col, err)
	}
	defer rows.Close()
	var out []NameCount
	for rows.Next() {
		var n NameCount
		if err := rows.Scan(&n.Name, &n.Count); err != nil {
			return nil, fmt.Errorf("device group %s scan: %w", col, err)
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// QueryDevices returns the device breakdown for the range.
func QueryDevices(db *sql.DB, days int, now int64) (*DeviceStats, error) {
	start := rangeStart(days, now)
	st := &DeviceStats{}
	var err error
	if st.Types, err = groupCount(db, "device_type", start, now, 5); err != nil {
		return nil, err
	}
	if st.OS, err = groupCount(db, "os_name", start, now, 10); err != nil {
		return nil, err
	}
	if st.Browsers, err = groupCount(db, "browser", start, now, 10); err != nil {
		return nil, err
	}
	// Top raw resolutions.
	rows, err := db.Query(`SELECT screen_w || 'x' || screen_h, COUNT(*)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND screen_w>0 AND screen_h>0
		GROUP BY 1 ORDER BY COUNT(*) DESC LIMIT 12`, start, now)
	if err != nil {
		return nil, fmt.Errorf("device resolutions: %w", err)
	}
	for rows.Next() {
		var n NameCount
		if err := rows.Scan(&n.Name, &n.Count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("device resolutions scan: %w", err)
		}
		st.Resolutions = append(st.Resolutions, n)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("device resolutions: %w", err)
	}
	rows.Close()
	// Width buckets.
	brows, err := db.Query(`SELECT
		CASE WHEN screen_w < 768 THEN 'Mobile <768'
		     WHEN screen_w < 1280 THEN 'Tablet 768-1279'
		     ELSE 'Desktop 1280+' END, COUNT(*)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND screen_w>0
		GROUP BY 1`, start, now)
	if err != nil {
		return nil, fmt.Errorf("device buckets: %w", err)
	}
	for brows.Next() {
		var n NameCount
		if err := brows.Scan(&n.Name, &n.Count); err != nil {
			brows.Close()
			return nil, fmt.Errorf("device buckets scan: %w", err)
		}
		st.Buckets = append(st.Buckets, n)
	}
	if err := brows.Err(); err != nil {
		brows.Close()
		return nil, fmt.Errorf("device buckets: %w", err)
	}
	brows.Close()
	return st, nil
}

// ---- geography ----

// CountryStat is one country for the range.
type CountryStat struct {
	Code    string
	Name    string
	Views   int64
	Uniques int64
}

// RegionStat is one region/city inside a country drill-down.
type RegionStat struct {
	Name    string
	Views   int64
	Uniques int64
}

// QueryCountries aggregates pageviews by country.
func QueryCountries(db *sql.DB, days int, now int64) ([]CountryStat, error) {
	start := rangeStart(days, now)
	rows, err := db.Query(`SELECT
			COALESCE(NULLIF(country_code,''),'(unknown)'),
			COALESCE(NULLIF(country_name,''),'(unknown)'),
			COUNT(*), COUNT(DISTINCT visitor_id)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=?
		GROUP BY 1, 2 ORDER BY COUNT(*) DESC LIMIT 30`, start, now)
	if err != nil {
		return nil, fmt.Errorf("countries: %w", err)
	}
	defer rows.Close()
	var out []CountryStat
	for rows.Next() {
		var c CountryStat
		if err := rows.Scan(&c.Code, &c.Name, &c.Views, &c.Uniques); err != nil {
			return nil, fmt.Errorf("countries scan: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// QueryCountryDrill breaks one country into regions (and top cities).
func QueryCountryDrill(db *sql.DB, days int, now int64, code string) (regions []RegionStat, cities []RegionStat, err error) {
	start := rangeStart(days, now)
	rrows, err := db.Query(`SELECT COALESCE(NULLIF(region,''),'(unknown)'), COUNT(*), COUNT(DISTINCT visitor_id)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND COALESCE(NULLIF(country_code,''),'(unknown)')=?
		GROUP BY 1 ORDER BY COUNT(*) DESC LIMIT 20`, start, now, code)
	if err != nil {
		return nil, nil, fmt.Errorf("regions: %w", err)
	}
	for rrows.Next() {
		var r RegionStat
		if err := rrows.Scan(&r.Name, &r.Views, &r.Uniques); err != nil {
			rrows.Close()
			return nil, nil, fmt.Errorf("regions scan: %w", err)
		}
		regions = append(regions, r)
	}
	if err := rrows.Err(); err != nil {
		rrows.Close()
		return nil, nil, fmt.Errorf("regions: %w", err)
	}
	rrows.Close()
	crows, err := db.Query(`SELECT COALESCE(NULLIF(city,''),'(unknown)'), COUNT(*), COUNT(DISTINCT visitor_id)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND COALESCE(NULLIF(country_code,''),'(unknown)')=? AND city<>''
		GROUP BY 1 ORDER BY COUNT(*) DESC LIMIT 20`, start, now, code)
	if err != nil {
		return nil, nil, fmt.Errorf("cities: %w", err)
	}
	for crows.Next() {
		var c RegionStat
		if err := crows.Scan(&c.Name, &c.Views, &c.Uniques); err != nil {
			crows.Close()
			return nil, nil, fmt.Errorf("cities scan: %w", err)
		}
		cities = append(cities, c)
	}
	if err := crows.Err(); err != nil {
		crows.Close()
		return nil, nil, fmt.Errorf("cities: %w", err)
	}
	crows.Close()
	return regions, cities, nil
}

// ---- visitors / retention ----

// VisitorStats holds retention numbers for the range.
type VisitorStats struct {
	NewTrend       []DailyCount // new visitors per day
	ReturningTrend []DailyCount // returning visitors per day
	Return7d       int64        // visitors back within 7 days of first visit
	Return7dPct    float64
	Cohort1d       float64 // % of visitors first seen 1+ days ago who returned within 1 day
	Cohort7d       float64
	Cohort30d      float64
	TopVisitors    []ActiveVisitor
}

// ActiveVisitor is an anonymized heavy user (truncated id, no PII).
type ActiveVisitor struct {
	IDPrefix string
	Views    int64
	Sessions int64
	LastSeen int64
}

// QueryVisitors computes retention stats for the range.
func QueryVisitors(db *sql.DB, days int, now int64) (*VisitorStats, error) {
	start := rangeStart(days, now)
	st := &VisitorStats{}

	// First-seen map for visitors active in the range.
	frows, err := db.Query(`SELECT visitor_id, MIN(ts) FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND visitor_id<>''
		GROUP BY visitor_id`, start, now)
	if err != nil {
		return nil, fmt.Errorf("visitors first-seen: %w", err)
	}
	firstSeen := map[string]int64{}
	for frows.Next() {
		var vid string
		var mn int64
		if err := frows.Scan(&vid, &mn); err != nil {
			frows.Close()
			return nil, fmt.Errorf("visitors first-seen scan: %w", err)
		}
		firstSeen[vid] = mn
	}
	if err := frows.Err(); err != nil {
		frows.Close()
		return nil, fmt.Errorf("visitors first-seen: %w", err)
	}
	frows.Close()

	// Anyone seen before the range at all counts as previously known.
	knownRows, err := db.Query(`SELECT DISTINCT visitor_id FROM analytics_events
		WHERE kind='pageview' AND ts<? AND visitor_id<>''`, start)
	if err != nil {
		return nil, fmt.Errorf("visitors known: %w", err)
	}
	knownBefore := map[string]bool{}
	for knownRows.Next() {
		var vid string
		if err := knownRows.Scan(&vid); err != nil {
			knownRows.Close()
			return nil, fmt.Errorf("visitors known scan: %w", err)
		}
		knownBefore[vid] = true
	}
	if err := knownRows.Err(); err != nil {
		knownRows.Close()
		return nil, fmt.Errorf("visitors known: %w", err)
	}
	knownRows.Close()

	// Daily new vs returning trend.
	dayStart := start - (start % 86400)
	nDays := int((now-dayStart)/86400) + 1
	newByDay := make([]int64, nDays)
	retByDay := make([]int64, nDays)
	seenDay := map[string]map[int]bool{}
	drows, err := db.Query(`SELECT visitor_id, CAST(ts/86400 AS INTEGER)*86400
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND visitor_id<>''
		GROUP BY visitor_id, 2`, start, now)
	if err != nil {
		return nil, fmt.Errorf("visitors daily: %w", err)
	}
	for drows.Next() {
		var vid string
		var ds int64
		if err := drows.Scan(&vid, &ds); err != nil {
			drows.Close()
			return nil, fmt.Errorf("visitors daily scan: %w", err)
		}
		idx := int((ds - dayStart) / 86400)
		if idx < 0 || idx >= nDays {
			continue
		}
		if seenDay[vid] == nil {
			seenDay[vid] = map[int]bool{}
		}
		if seenDay[vid][idx] {
			continue
		}
		seenDay[vid][idx] = true
		if knownBefore[vid] || firstSeen[vid] < ds {
			retByDay[idx]++
		} else {
			newByDay[idx]++
		}
	}
	if err := drows.Err(); err != nil {
		drows.Close()
		return nil, fmt.Errorf("visitors daily: %w", err)
	}
	drows.Close()
	for d := 0; d < nDays; d++ {
		lbl := dayLabel(dayStart + int64(d)*86400)
		st.NewTrend = append(st.NewTrend, DailyCount{Date: lbl, Count: newByDay[d]})
		st.ReturningTrend = append(st.ReturningTrend, DailyCount{Date: lbl, Count: retByDay[d]})
	}

	// Cohort return rates: visitors first seen long enough ago that the
	// window has fully elapsed, who came back within the window.
	cohort := func(windowDays int64) (float64, error) {
		var eligible, returned int64
		for vid, first := range firstSeen {
			if knownBefore[vid] {
				continue // not a first visit, skip
			}
			if now-first < windowDays*86400 {
				continue // window not elapsed yet
			}
			eligible++
			var n int64
			if err := db.QueryRow(`SELECT COUNT(*) FROM analytics_events
				WHERE kind='pageview' AND visitor_id=? AND ts>? AND ts<=?`,
				vid, first, first+windowDays*86400).Scan(&n); err != nil {
				return 0, fmt.Errorf("cohort: %w", err)
			}
			if n > 0 {
				returned++
			}
		}
		if eligible == 0 {
			return 0, nil
		}
		return float64(returned) / float64(eligible), nil
	}
	if st.Cohort1d, err = cohort(1); err != nil {
		return nil, err
	}
	if st.Cohort7d, err = cohort(7); err != nil {
		return nil, err
	}
	if st.Cohort30d, err = cohort(30); err != nil {
		return nil, err
	}

	// Returned within 7 days (any visitor first seen in range, window may
	// still be open; labeled as such on the dashboard).
	var r7, r7e int64
	for vid, first := range firstSeen {
		if knownBefore[vid] || now-first < 7*86400 {
			continue
		}
		r7e++
		var n int64
		if err := db.QueryRow(`SELECT COUNT(*) FROM analytics_events
			WHERE kind='pageview' AND visitor_id=? AND ts>? AND ts<=?`,
			vid, first, first+7*86400).Scan(&n); err != nil {
			return nil, fmt.Errorf("return7d: %w", err)
		}
		if n > 0 {
			r7++
		}
	}
	st.Return7d = r7
	if r7e > 0 {
		st.Return7dPct = float64(r7) / float64(r7e)
	}

	// Most active anonymous visitors.
	arows, err := db.Query(`SELECT visitor_id, COUNT(*), COUNT(DISTINCT session_id), MAX(ts)
		FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND visitor_id<>''
		GROUP BY visitor_id ORDER BY COUNT(*) DESC LIMIT 10`, start, now)
	if err != nil {
		return nil, fmt.Errorf("top visitors: %w", err)
	}
	defer arows.Close()
	for arows.Next() {
		var vid string
		var v ActiveVisitor
		if err := arows.Scan(&vid, &v.Views, &v.Sessions, &v.LastSeen); err != nil {
			return nil, fmt.Errorf("top visitors scan: %w", err)
		}
		if len(vid) > 8 {
			vid = vid[:8]
		}
		v.IDPrefix = vid
		st.TopVisitors = append(st.TopVisitors, v)
	}
	return st, arows.Err()
}

// ---- activity feed ----

// FeedHit is one event for the live activity feed.
type FeedHit struct {
	TS         int64
	Kind       string
	Path       string
	Source     string
	Label      string
	Href       string
	DeviceType string
	Browser    string
	Country    string
}

// QueryFeed returns the most recent events for the activity feed,
// optionally filtered by kind and page substring.
func QueryFeed(db *sql.DB, limit int, kindFilter, pageFilter string) ([]FeedHit, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	q := `SELECT ts, kind, path, COALESCE(NULLIF(source,''),'Direct'), label, href,
			COALESCE(NULLIF(device_type,''),'desktop'),
			COALESCE(NULLIF(browser,''),'Other'),
			COALESCE(NULLIF(country_code,''),'??')
		FROM analytics_events WHERE 1=1`
	var args []any
	if kindFilter == "pageview" || kindFilter == "click" {
		q += ` AND kind=?`
		args = append(args, kindFilter)
	}
	if pageFilter != "" {
		q += ` AND path LIKE ?`
		args = append(args, "%"+pageFilter+"%")
	}
	q += ` ORDER BY ts DESC, id DESC LIMIT ?`
	args = append(args, limit)
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("feed: %w", err)
	}
	defer rows.Close()
	var out []FeedHit
	for rows.Next() {
		var h FeedHit
		if err := rows.Scan(&h.TS, &h.Kind, &h.Path, &h.Source, &h.Label, &h.Href,
			&h.DeviceType, &h.Browser, &h.Country); err != nil {
			return nil, fmt.Errorf("feed scan: %w", err)
		}
		out = append(out, h)
	}
	return out, rows.Err()
}
