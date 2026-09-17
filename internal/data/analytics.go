// First-party website analytics storage. Events land in the same SQLite
// catalog file the marketing module already uses (m.DBPath), in dedicated
// tables created on demand. Raw IPs are never stored — only a SHA-256 hash
// (one-way, for rough dedup/counting; not a secret, not reversible).
package data

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// Event is one first-party analytics hit: a pageview, a click, or a
// generic named event. Timestamps are unix seconds (UTC).
type Event struct {
	TS        int64
	VisitorID string
	SessionID string
	Kind      string // pageview | click | event | session_end
	Path      string
	Referrer  string
	Source    string // attributed traffic source (Facebook, Google, Direct, ...)
	Medium    string // utm_medium when present
	Campaign  string // utm_campaign when present
	Element   string // tag name for clicks
	Label     string // click label / event detail
	UA        string
	IPHash    string
}

var (
	analyticsMu  sync.Mutex
	analyticsDBs = map[string]*sql.DB{}
)

// AnalyticsDB opens (once per path) a handle on dbPath for the analytics
// tables and ensures the schema exists. A separate handle from the catalog
// DB is fine: SQLite WAL mode supports concurrent readers and one writer.
func AnalyticsDB(dbPath string) (*sql.DB, error) {
	if strings.ContainsAny(dbPath, "?#") {
		return nil, fmt.Errorf("refusing db path with url metacharacters: %q", dbPath)
	}
	analyticsMu.Lock()
	defer analyticsMu.Unlock()
	if db, ok := analyticsDBs[dbPath]; ok {
		return db, nil
	}
	db, err := sql.Open("sqlite", "file:"+dbPath+"?_pragma=journal_mode(WAL)&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open analytics db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping analytics db: %w", err)
	}
	if err := EnsureAnalyticsTables(db); err != nil {
		db.Close()
		return nil, err
	}
	analyticsDBs[dbPath] = db
	return db, nil
}

// EnsureAnalyticsTables creates the analytics schema when missing.
func EnsureAnalyticsTables(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS analytics_events (
			id INTEGER PRIMARY KEY,
			ts INTEGER NOT NULL,
			visitor_id TEXT,
			session_id TEXT,
			kind TEXT NOT NULL,
			path TEXT,
			referrer TEXT,
			source TEXT,
			medium TEXT,
			campaign TEXT,
			element TEXT,
			label TEXT,
			ua TEXT,
			ip_hash TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_ae_ts ON analytics_events (ts)`,
		`CREATE INDEX IF NOT EXISTS idx_ae_path ON analytics_events (path)`,
		`CREATE INDEX IF NOT EXISTS idx_ae_source ON analytics_events (source)`,
		`CREATE INDEX IF NOT EXISTS idx_ae_session ON analytics_events (session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ae_visitor ON analytics_events (visitor_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ae_kind_ts ON analytics_events (kind, ts)`,
		`CREATE TABLE IF NOT EXISTS admin_sessions (
			token_hash TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			expires_at INTEGER NOT NULL
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("analytics schema: %w", err)
		}
	}
	return nil
}

// capStr truncates s to n runes so overlong client input cannot bloat rows.
func capStr(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n])
}

// InsertEvent stores one analytics hit. Text fields are capped server-side
// (512 runes, 2048 for referrer) regardless of what the client sent.
func InsertEvent(db *sql.DB, e Event) error {
	if e.Path == "" {
		e.Path = "/"
	}
	_, err := db.Exec(
		`INSERT INTO analytics_events
		 (ts, visitor_id, session_id, kind, path, referrer, source, medium, campaign, element, label, ua, ip_hash)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.TS,
		capStr(e.VisitorID, 512), capStr(e.SessionID, 512), capStr(e.Kind, 64),
		capStr(e.Path, 512), capStr(e.Referrer, 2048),
		capStr(e.Source, 64), capStr(e.Medium, 64), capStr(e.Campaign, 128),
		capStr(e.Element, 16), capStr(e.Label, 512),
		capStr(e.UA, 512), capStr(e.IPHash, 128),
	)
	if err != nil {
		return fmt.Errorf("insert analytics event: %w", err)
	}
	return nil
}

// AttributeSource turns a referrer URL plus UTM params into a human traffic
// source. utm_source always wins (it is the advertiser's own labeling);
// otherwise the referrer host is mapped to the known platform, falling back
// to the bare registrable domain. Empty referrer means Direct.
func AttributeSource(referrer, utmSource, utmMedium string) (source, medium string) {
	utmSource = strings.TrimSpace(utmSource)
	medium = capStr(strings.TrimSpace(utmMedium), 64)
	if utmSource != "" {
		return capStr(utmSource, 64), medium
	}
	referrer = strings.TrimSpace(referrer)
	if referrer == "" {
		return "Direct", medium
	}
	host := ""
	if u, err := url.Parse(referrer); err == nil {
		host = strings.ToLower(u.Hostname())
	}
	if host == "" {
		return "Direct", medium
	}
	switch {
	case strings.Contains(host, "facebook.com") || strings.Contains(host, "instagram.com") || host == "fb.com":
		return "Facebook", medium
	case strings.Contains(host, "google."):
		return "Google", medium
	case strings.Contains(host, "bing.com"):
		return "Bing", medium
	case strings.Contains(host, "duckduckgo.com"):
		return "DuckDuckGo", medium
	case strings.Contains(host, "yahoo.com"):
		return "Yahoo", medium
	case strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be"):
		return "YouTube", medium
	case strings.Contains(host, "tiktok.com"):
		return "TikTok", medium
	case host == "x.com" || strings.Contains(host, "twitter.com"):
		return "X", medium
	case strings.Contains(host, "nextdoor.com"):
		return "Nextdoor", medium
	case strings.Contains(host, "yelp.com"):
		return "Yelp", medium
	}
	// Registrable domain fallback: last two labels.
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "." + parts[len(parts)-1], medium
	}
	return host, medium
}

// botSubstrings — any UA containing one of these is a bot/crawler and is
// excluded from the dashboard. AI crawlers are in here too: they read the
// site for search answers, they are not human visitors.
var botSubstrings = []string{
	"bot", "crawler", "spider", "slurp", "mediapartners",
	"baidu", "yandex", "sogou", "exabot", "facebot", "facebookexternalhit", "externalhit", "ia_archiver",
	"ahrefs", "semrush", "mj12bot", "dotbot", "petalbot", "bytespider",
	"gptbot", "chatgpt-user", "claudebot", "anthropic-ai", "perplexitybot",
	"applebot", "mediacodec", "curl", "wget", "python-requests", "go-http-client",
}

// IsBot reports whether a User-Agent looks automated.
func IsBot(ua string) bool {
	l := strings.ToLower(ua)
	for _, b := range botSubstrings {
		if strings.Contains(l, b) {
			return true
		}
	}
	return false
}

// HashIP returns the hex SHA-256 of an IP address. One-way: the dashboard
// never needs the raw IP, only a stable per-visitor fingerprint input.
func HashIP(ip string) string {
	sum := sha256.Sum256([]byte("patriot-analytics|" + strings.TrimSpace(ip)))
	return hex.EncodeToString(sum[:])
}

// ---- Dashboard aggregates ----

// DailyCount is one day bucket of pageviews.
type DailyCount struct {
	Date  string // YYYY-MM-DD
	Count int64
}

// PathCount is a page path with its view count.
type PathCount struct {
	Path  string
	Count int64
}

// SourceCount is a traffic source with its view count.
type SourceCount struct {
	Source string
	Count  int64
}

// ClickCount is a clicked element label with its click count and page.
type ClickCount struct {
	Label string
	Path  string
	Count int64
}

// RecentHit is one recent event for the activity feed.
type RecentHit struct {
	TS     int64
	Kind   string
	Path   string
	Source string
	Label  string
}

// AnalyticsStats is everything the admin dashboard renders for a range.
type AnalyticsStats struct {
	Days              int
	Pageviews         int64
	Visitors          int64 // distinct visitor_id on pageviews
	Sessions          int64 // distinct session_id on pageviews
	NewVisitors       int64
	ReturningVisitors int64
	Daily             []DailyCount
	TopPages          []PathCount
	TopSources        []SourceCount
	TopClicks         []ClickCount
	Recent            []RecentHit
}

// QueryAnalytics aggregates the dashboard stats for the last `days` days.
// days must be one of 7, 30, 90; anything else becomes 30.
func QueryAnalytics(db *sql.DB, days int, now int64) (*AnalyticsStats, error) {
	switch days {
	case 7, 30, 90:
	default:
		days = 30
	}
	start := now - int64(days)*86400
	st := &AnalyticsStats{Days: days}

	q1 := func(query, argDesc string, dest *int64, args ...any) error {
		if err := db.QueryRow(query, args...).Scan(dest); err != nil {
			return fmt.Errorf("analytics %s: %w", argDesc, err)
		}
		return nil
	}
	if err := q1(`SELECT COUNT(*) FROM analytics_events WHERE kind='pageview' AND ts>=? AND ts<=?`,
		"pageviews", &st.Pageviews, start, now); err != nil {
		return nil, err
	}
	if err := q1(`SELECT COUNT(DISTINCT visitor_id) FROM analytics_events WHERE kind='pageview' AND ts>=? AND ts<=? AND visitor_id<>''`,
		"visitors", &st.Visitors, start, now); err != nil {
		return nil, err
	}
	if err := q1(`SELECT COUNT(DISTINCT session_id) FROM analytics_events WHERE kind='pageview' AND ts>=? AND ts<=? AND session_id<>''`,
		"sessions", &st.Sessions, start, now); err != nil {
		return nil, err
	}
	// Returning = visitors in range who were first seen before the range.
	if err := q1(`SELECT COUNT(DISTINCT visitor_id) FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? AND visitor_id<>''
		AND visitor_id IN (SELECT DISTINCT visitor_id FROM analytics_events WHERE kind='pageview' AND ts<=?)`,
		"returning", &st.ReturningVisitors, start, now, start); err != nil {
		return nil, err
	}
	st.NewVisitors = st.Visitors - st.ReturningVisitors
	if st.NewVisitors < 0 {
		st.NewVisitors = 0
	}

	// Daily buckets (fill gaps with zeros so the chart is continuous).
	dayMap := map[string]int64{}
	rows, err := db.Query(`SELECT CAST(ts/86400 AS INTEGER)*86400, COUNT(*)
		FROM analytics_events WHERE kind='pageview' AND ts>=? AND ts<=? GROUP BY 1 ORDER BY 1`, start, now)
	if err != nil {
		return nil, fmt.Errorf("analytics daily: %w", err)
	}
	for rows.Next() {
		var dayStart, c int64
		if err := rows.Scan(&dayStart, &c); err != nil {
			rows.Close()
			return nil, fmt.Errorf("analytics daily scan: %w", err)
		}
		dayMap[dayLabel(dayStart)] = c
	}
	rows.Close()
	for d := start - (start % 86400); d <= now; d += 86400 {
		lbl := dayLabel(d)
		st.Daily = append(st.Daily, DailyCount{Date: lbl, Count: dayMap[lbl]})
	}
	// Trim leading all-zero buckets outside the exact window.
	for len(st.Daily) > days+1 {
		st.Daily = st.Daily[1:]
	}

	top := func(query string, limit int) ([]PathCount, error) {
		r, err := db.Query(query, start, now, limit)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		var out []PathCount
		for r.Next() {
			var pc PathCount
			if err := r.Scan(&pc.Path, &pc.Count); err != nil {
				return nil, err
			}
			if pc.Path == "" {
				pc.Path = "/"
			}
			out = append(out, pc)
		}
		return out, r.Err()
	}
	if st.TopPages, err = top(`SELECT path, COUNT(*) FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? GROUP BY path ORDER BY COUNT(*) DESC LIMIT ?`, 15); err != nil {
		return nil, fmt.Errorf("analytics top pages: %w", err)
	}
	sr, err := db.Query(`SELECT source, COUNT(*) FROM analytics_events
		WHERE kind='pageview' AND ts>=? AND ts<=? GROUP BY source ORDER BY COUNT(*) DESC LIMIT 15`, start, now)
	if err != nil {
		return nil, fmt.Errorf("analytics top sources: %w", err)
	}
	for sr.Next() {
		var sc SourceCount
		if err := sr.Scan(&sc.Source, &sc.Count); err != nil {
			sr.Close()
			return nil, fmt.Errorf("analytics top sources scan: %w", err)
		}
		if sc.Source == "" {
			sc.Source = "Direct"
		}
		st.TopSources = append(st.TopSources, sc)
	}
	if err := sr.Err(); err != nil {
		sr.Close()
		return nil, fmt.Errorf("analytics top sources: %w", err)
	}
	sr.Close()

	cr, err := db.Query(`SELECT label, path, COUNT(*) FROM analytics_events
		WHERE kind='click' AND ts>=? AND ts<=? GROUP BY label, path ORDER BY COUNT(*) DESC LIMIT 15`, start, now)
	if err != nil {
		return nil, fmt.Errorf("analytics top clicks: %w", err)
	}
	for cr.Next() {
		var cc ClickCount
		if err := cr.Scan(&cc.Label, &cc.Path, &cc.Count); err != nil {
			cr.Close()
			return nil, fmt.Errorf("analytics top clicks scan: %w", err)
		}
		st.TopClicks = append(st.TopClicks, cc)
	}
	if err := cr.Err(); err != nil {
		cr.Close()
		return nil, fmt.Errorf("analytics top clicks: %w", err)
	}
	cr.Close()

	rr, err := db.Query(`SELECT ts, kind, path, source, label FROM analytics_events
		WHERE ts>=? AND ts<=? AND kind IN ('pageview','click') ORDER BY ts DESC, id DESC LIMIT 50`, start, now)
	if err != nil {
		return nil, fmt.Errorf("analytics recent: %w", err)
	}
	for rr.Next() {
		var h RecentHit
		if err := rr.Scan(&h.TS, &h.Kind, &h.Path, &h.Source, &h.Label); err != nil {
			rr.Close()
			return nil, fmt.Errorf("analytics recent scan: %w", err)
		}
		st.Recent = append(st.Recent, h)
	}
	if err := rr.Err(); err != nil {
		rr.Close()
		return nil, fmt.Errorf("analytics recent: %w", err)
	}
	rr.Close()
	return st, nil
}

// dayLabel formats a unix day-start as YYYY-MM-DD (UTC).
func dayLabel(dayStart int64) string {
	return time.Unix(dayStart, 0).UTC().Format("2006-01-02")
}
