package data

import (
	"database/sql"
	"os"
	"testing"
)

// ---- device parsing ----

func TestParseDevice(t *testing.T) {
	cases := []struct {
		name, ua, dtype, os, browser string
	}{
		{"iPhone Safari",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
			"phone", "iOS", "Safari"},
		{"Android Chrome phone",
			"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36",
			"phone", "Android", "Chrome"},
		{"Android tablet",
			"Mozilla/5.0 (Linux; Android 13; SM-X810) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
			"tablet", "Android", "Chrome"},
		{"Windows Chrome",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
			"desktop", "Windows", "Chrome"},
		{"Mac Safari",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15",
			"desktop", "macOS", "Safari"},
		{"iPad",
			"Mozilla/5.0 (iPad; CPU OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
			"tablet", "iOS", "Safari"},
		{"Linux Firefox",
			"Mozilla/5.0 (X11; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0",
			"desktop", "Linux", "Firefox"},
		{"Edge",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
			"desktop", "Windows", "Edge"},
		{"Samsung Internet",
			"Mozilla/5.0 (Linux; Android 14; SAMSUNG SM-S921B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/26.0 Chrome/126.0.0.0 Mobile Safari/537.36",
			"phone", "Android", "Samsung Internet"},
		{"Chromebook",
			"Mozilla/5.0 (X11; CrOS x86_64 15699.58.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
			"desktop", "ChromeOS", "Chrome"},
		{"iOS Chrome",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/126.0.6367.66 Mobile/15E148 Safari/604.1",
			"phone", "iOS", "Chrome"},
		{"empty",
			"", "desktop", "Other", "Other"},
		{"bot UA still parses without crashing",
			"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			"desktop", "Other", "Other"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := ParseDevice(c.ua)
			if d.DeviceType != c.dtype || d.OS != c.os || d.Browser != c.browser {
				t.Errorf("ParseDevice = %+v, want %s/%s/%s", d, c.dtype, c.os, c.browser)
			}
		})
	}
}

// ---- geo ----

func TestGeoFallbackWhenDBMissing(t *testing.T) {
	if GeoConfigured() {
		t.Skip("geo db already configured in this process")
	}
	g := LookupGeo("8.8.8.8")
	if !g.Empty() {
		t.Errorf("LookupGeo without DB = %+v, want empty", g)
	}
	for _, ip := range []string{"", "not-an-ip", "127.0.0.1", "192.168.1.1", "10.0.0.5"} {
		if g := LookupGeo(ip); !g.Empty() {
			t.Errorf("LookupGeo(%q) without DB = %+v, want empty", ip, g)
		}
	}
}

func TestGeoLookupRealDB(t *testing.T) {
	dbPath := os.Getenv("GEO_TEST_DB")
	if dbPath == "" {
		dbPath = "/tmp/geo/GeoLite2-City.mmdb"
	}
	if _, err := os.Stat(dbPath); err != nil {
		t.Skipf("no test geo db at %s", dbPath)
	}
	InitGeoDB(dbPath)
	if !GeoConfigured() {
		t.Fatal("GeoConfigured false after loading a real DB")
	}
	g := LookupGeo("8.8.8.8")
	if g.CountryCode != "US" {
		t.Errorf("LookupGeo(8.8.8.8) = %+v, want country US", g)
	}
	// Private/loopback/garbage never resolve.
	for _, ip := range []string{"127.0.0.1", "10.1.2.3", "junk"} {
		if g := LookupGeo(ip); !g.Empty() {
			t.Errorf("LookupGeo(%q) = %+v, want empty", ip, g)
		}
	}
}

// ---- migration ----

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+t.TempDir()+"/t.db?_pragma=journal_mode(WAL)&_busy_timeout=5000")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// legacySchemaDB builds a database with the pre-expansion schema.
func legacySchemaDB(t *testing.T) *sql.DB {
	t.Helper()
	db := testDB(t)
	_, err := db.Exec(`CREATE TABLE analytics_events (
		id INTEGER PRIMARY KEY, ts INTEGER NOT NULL, visitor_id TEXT,
		session_id TEXT, kind TEXT NOT NULL, path TEXT, referrer TEXT,
		source TEXT, medium TEXT, campaign TEXT, element TEXT, label TEXT,
		ua TEXT, ip_hash TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO analytics_events
		(ts, visitor_id, session_id, kind, path, source, ua, ip_hash)
		VALUES (1000, 'oldv', 'olds', 'pageview', '/old', 'Direct', 'ua', 'h')`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestMigrationIdempotent(t *testing.T) {
	db := legacySchemaDB(t)
	for i := 0; i < 2; i++ {
		if err := EnsureAnalyticsTables(db); err != nil {
			t.Fatalf("ensure pass %d: %v", i, err)
		}
	}
	// Old row survived.
	var n int64
	if err := db.QueryRow(`SELECT COUNT(*) FROM analytics_events WHERE path='/old'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("legacy row lost: n=%d err=%v", n, err)
	}
	// New columns exist and accept writes.
	e := Event{TS: 2000, VisitorID: "v2", SessionID: "s2", Kind: "click",
		Path: "/p", Label: "Buy now", Href: "/prices", ElemHint: "#buy",
		DeviceType: "phone", OS: "iOS", Browser: "Safari",
		ScreenW: 390, ScreenH: 844,
		CountryCode: "US", CountryName: "United States", Region: "Kentucky", City: "Louisville"}
	if err := InsertEvent(db, e); err != nil {
		t.Fatalf("insert with new fields: %v", err)
	}
	var got ClickDetail
	err := db.QueryRow(`SELECT label, href, elem_hint FROM analytics_events WHERE visitor_id='v2'`).
		Scan(&got.Label, &got.Href, &got.ElemHint)
	if err != nil {
		t.Fatal(err)
	}
	if got.Label != "Buy now" || got.Href != "/prices" || got.ElemHint != "#buy" {
		t.Errorf("roundtrip = %+v", got)
	}
	var dt, osn, br, cc, city string
	var sw, sh int
	err = db.QueryRow(`SELECT device_type, os_name, browser, screen_w, screen_h, country_code, city
		FROM analytics_events WHERE visitor_id='v2'`).Scan(&dt, &osn, &br, &sw, &sh, &cc, &city)
	if err != nil {
		t.Fatal(err)
	}
	if dt != "phone" || osn != "iOS" || br != "Safari" || sw != 390 || sh != 844 || cc != "US" || city != "Louisville" {
		t.Errorf("device/geo roundtrip wrong: %q %q %q %d %d %q %q", dt, osn, br, sw, sh, cc, city)
	}
	// Screen clamping: garbage numbers stay in range.
	if err := InsertEvent(db, Event{TS: 2001, Kind: "pageview", ScreenW: -5, ScreenH: 999999}); err != nil {
		t.Fatal(err)
	}
	var cw, ch int
	if err := db.QueryRow(`SELECT screen_w, screen_h FROM analytics_events WHERE ts=2001`).Scan(&cw, &ch); err != nil {
		t.Fatal(err)
	}
	if cw != 0 || ch != 10000 {
		t.Errorf("clamp = %d/%d, want 0/10000", cw, ch)
	}
}

// ---- new queries ----

func seedEvents(t *testing.T, db *sql.DB, now int64) {
	t.Helper()
	evs := []Event{
		// Visitor A: two-page session (not a bounce), iPhone Safari, US.
		{TS: now - 3600, VisitorID: "visA", SessionID: "sessA", Kind: "pageview", Path: "/", Source: "Google", DeviceType: "phone", OS: "iOS", Browser: "Safari", CountryCode: "US", CountryName: "United States"},
		{TS: now - 3500, VisitorID: "visA", SessionID: "sessA", Kind: "pageview", Path: "/services", Source: "Google", DeviceType: "phone", OS: "iOS", Browser: "Safari", CountryCode: "US", CountryName: "United States"},
		{TS: now - 3490, VisitorID: "visA", SessionID: "sessA", Kind: "click", Path: "/services", Label: "Get a quote", Href: "/contact", ElemHint: "#quote", DeviceType: "phone", OS: "iOS", Browser: "Safari", CountryCode: "US"},
		// Visitor B: single-page session (bounce), Windows Chrome, Direct.
		{TS: now - 7200, VisitorID: "visB", SessionID: "sessB", Kind: "pageview", Path: "/", Source: "Direct", DeviceType: "desktop", OS: "Windows", Browser: "Chrome", CountryCode: "US", CountryName: "United States", ScreenW: 1920, ScreenH: 1080},
		// Visitor C: first visit 10 days ago, back the next day (returning).
		{TS: now - 10*86400, VisitorID: "visC", SessionID: "sessC1", Kind: "pageview", Path: "/", Source: "Facebook", DeviceType: "desktop", OS: "macOS", Browser: "Safari", CountryCode: "CA", CountryName: "Canada"},
		{TS: now - 9*86400, VisitorID: "visC", SessionID: "sessC2", Kind: "pageview", Path: "/prices", Source: "Facebook", DeviceType: "desktop", OS: "macOS", Browser: "Safari", CountryCode: "CA", CountryName: "Canada"},
	}
	for _, e := range evs {
		if err := InsertEvent(db, e); err != nil {
			t.Fatal(err)
		}
	}
}

func TestQuerySessionStats(t *testing.T) {
	db := testDB(t)
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	now := int64(1758000000)
	seedEvents(t, db, now)
	st, err := QuerySessionStats(db, 30, now)
	if err != nil {
		t.Fatal(err)
	}
	if st.Sessions != 4 { // sessA, sessB, sessC1, sessC2
		t.Errorf("sessions = %d, want 4", st.Sessions)
	}
	if st.BounceRate != 0.75 { // 3 of 4 single-page
		t.Errorf("bounce = %v, want 0.75", st.BounceRate)
	}
	if st.PagesPerSess != 1.25 { // 5 pageviews / 4 sessions
		t.Errorf("pages/session = %v, want 1.25", st.PagesPerSess)
	}
	if st.AvgDurationSec != 25 { // sessA lasted 100s, others 0: 100/4
		t.Errorf("avg duration = %v, want 25", st.AvgDurationSec)
	}
}

func TestQueryPageStats(t *testing.T) {
	db := testDB(t)
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	now := int64(1758000000)
	seedEvents(t, db, now)
	pages, err := QueryPageStats(db, 30, now)
	if err != nil {
		t.Fatal(err)
	}
	byPath := map[string]PageStat{}
	for _, p := range pages {
		byPath[p.Path] = p
	}
	home := byPath["/"]
	if home.Views != 3 { // visA, visB, visC(old)
		t.Errorf("/ views = %d, want 3", home.Views)
	}
	if home.Entrances != 3 {
		t.Errorf("/ entrances = %d, want 3", home.Entrances)
	}
	svcs := byPath["/services"]
	if svcs.Exits != 1 { // sessA ended on /services
		t.Errorf("/services exits = %d, want 1", svcs.Exits)
	}
	// Time on "/" is the gap to the next pageview in the same session.
	if byPath["/"].AvgTimeSec != 100 {
		t.Errorf("/ avg time = %v, want 100", byPath["/"].AvgTimeSec)
	}
}

func TestQueryClicks(t *testing.T) {
	db := testDB(t)
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	now := int64(1758000000)
	seedEvents(t, db, now)
	clicks, err := QueryClicks(db, 30, now, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(clicks) != 1 {
		t.Fatalf("clicks = %d, want 1", len(clicks))
	}
	c := clicks[0]
	if c.Label != "Get a quote" || c.Href != "/contact" || c.Path != "/services" || c.ElemHint != "#quote" || c.Clicks != 1 || c.Clickers != 1 {
		t.Errorf("click = %+v", c)
	}
	filtered, err := QueryClicks(db, 30, now, "/nope")
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 0 {
		t.Errorf("filtered clicks = %d, want 0", len(filtered))
	}
}

func TestQueryDevicesAndGeo(t *testing.T) {
	db := testDB(t)
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	now := int64(1758000000)
	seedEvents(t, db, now)
	dev, err := QueryDevices(db, 30, now)
	if err != nil {
		t.Fatal(err)
	}
	typeCount := map[string]int64{}
	for _, x := range dev.Types {
		typeCount[x.Name] = x.Count
	}
	if typeCount["phone"] != 2 || typeCount["desktop"] != 3 {
		t.Errorf("device types = %v", typeCount)
	}
	countries, err := QueryCountries(db, 30, now)
	if err != nil {
		t.Fatal(err)
	}
	cc := map[string]int64{}
	for _, c := range countries {
		cc[c.Code] = c.Views
	}
	if cc["US"] != 3 || cc["CA"] != 2 {
		t.Errorf("countries = %v", cc)
	}
	regions, cities, err := QueryCountryDrill(db, 30, now, "US")
	if err != nil {
		t.Fatal(err)
	}
	if len(regions) == 0 {
		t.Error("US drill returned no regions")
	}
	_ = cities
}

func TestQueryVisitorsAndFeed(t *testing.T) {
	db := testDB(t)
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	now := int64(1758000000)
	seedEvents(t, db, now)
	v, err := QueryVisitors(db, 30, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(v.TopVisitors) == 0 {
		t.Error("no top visitors")
	}
	if v.TopVisitors[0].Views < 2 {
		t.Errorf("top visitor views = %d, want >= 2", v.TopVisitors[0].Views)
	}
	if len(v.TopVisitors[0].IDPrefix) > 8 {
		t.Errorf("visitor id not truncated: %q", v.TopVisitors[0].IDPrefix)
	}
	// Cohort: visC first seen 10 days ago, back the next day. Nobody is
	// 30 days old yet, so the day-30 cohort is empty by design.
	if v.Cohort1d != 1 || v.Cohort7d != 1 || v.Cohort30d != 0 {
		t.Errorf("cohorts = %v/%v/%v, want 1/1/0", v.Cohort1d, v.Cohort7d, v.Cohort30d)
	}
	feed, err := QueryFeed(db, 10, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed) == 0 {
		t.Fatal("empty feed")
	}
	if feed[0].DeviceType == "" || feed[0].Browser == "" {
		t.Errorf("feed missing device info: %+v", feed[0])
	}
	clicks, err := QueryFeed(db, 10, "click", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range clicks {
		if h.Kind != "click" {
			t.Errorf("kind filter leaked %q", h.Kind)
		}
	}
}

func TestReferrerDomains(t *testing.T) {
	db := testDB(t)
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	now := int64(1758000000)
	for _, ref := range []string{"https://www.facebook.com/ads", "https://m.facebook.com/x", "https://example.com/p"} {
		if err := InsertEvent(db, Event{TS: now - 100, Kind: "pageview", Path: "/", Referrer: ref, Source: "Facebook"}); err != nil {
			t.Fatal(err)
		}
	}
	refs, err := QueryReferrers(db, 30, now, "www.patriotpest.pro")
	if err != nil {
		t.Fatal(err)
	}
	byDom := map[string]int64{}
	for _, r := range refs {
		byDom[r.Domain] = r.Views
	}
	if byDom["www.facebook.com"] != 1 || byDom["m.facebook.com"] != 1 || byDom["example.com"] != 1 {
		t.Errorf("referrers = %v", byDom)
	}
}
