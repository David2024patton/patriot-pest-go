package data

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestAttributeSource(t *testing.T) {
	cases := []struct {
		name       string
		referrer   string
		utmSource  string
		utmMedium  string
		wantSource string
		wantMedium string
	}{
		{"facebook referrer", "https://www.facebook.com/somepage", "", "", "Facebook", ""},
		{"instagram referrer", "https://instagram.com/patriot_pest/", "", "", "Facebook", ""},
		{"google referrer", "https://www.google.com/search?q=pest+control", "", "", "Google", ""},
		{"google co uk", "https://www.google.co.uk/search?q=x", "", "", "Google", ""},
		{"bing", "https://www.bing.com/search?q=x", "", "", "Bing", ""},
		{"duckduckgo", "https://duckduckgo.com/?q=x", "", "", "DuckDuckGo", ""},
		{"youtube", "https://www.youtube.com/watch?v=abc", "", "", "YouTube", ""},
		{"tiktok", "https://www.tiktok.com/@user/video/1", "", "", "TikTok", ""},
		{"x", "https://x.com/user/status/1", "", "", "X", ""},
		{"twitter", "https://twitter.com/user/status/1", "", "", "X", ""},
		{"nextdoor", "https://nextdoor.com/news_feed/", "", "", "Nextdoor", ""},
		{"yelp", "https://www.yelp.com/biz/abc", "", "", "Yelp", ""},
		{"utm overrides referrer", "https://www.google.com/", "facebook", "cpc", "facebook", "cpc"},
		{"utm only", "", "newsletter", "email", "newsletter", "email"},
		{"empty is direct", "", "", "", "Direct", ""},
		{"unknown domain passthrough", "https://blog.somecontractor.net/post/1", "", "", "somecontractor.net", ""},
		{"subdomain passthrough", "https://news.localpaper.example.com/a", "", "", "example.com", ""},
		{"medium passthrough", "", "google", "organic", "google", "organic"},
		{"utm trimmed", "  ", "  Facebook  ", "  ", "Facebook", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotSource, gotMedium := AttributeSource(c.referrer, c.utmSource, c.utmMedium)
			if gotSource != c.wantSource || gotMedium != c.wantMedium {
				t.Errorf("AttributeSource(%q,%q,%q) = (%q,%q), want (%q,%q)",
					c.referrer, c.utmSource, c.utmMedium, gotSource, gotMedium, c.wantSource, c.wantMedium)
			}
		})
	}
}

func TestIsBot(t *testing.T) {
	bots := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"DuckDuckBot/1.0; (+http://duckduckgo.com/duckduckbot.html)",
		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
		"GPTBot/1.0; +https://openai.com/gptbot",
		"ChatGPT-User/1.0; +https://openai.com/bot",
		"ClaudeBot/1.0; +https://claude.com/bot",
		"anthropic-ai/1.0; +http://www.anthropic.com/",
		"PerplexityBot/1.0; +https://perplexity.ai/perplexitybot",
		"Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
		"Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)",
		"MJ12bot/v1.4.8 (http://mj12bot.com/)",
	}
	for _, ua := range bots {
		if !IsBot(ua) {
			t.Errorf("IsBot(%q) = false, want true", ua)
		}
	}
	humans := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
		"",
	}
	for _, ua := range humans {
		if IsBot(ua) {
			t.Errorf("IsBot(%q) = true, want false", ua)
		}
	}
}

func openAnalyticsTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// :memory: databases are per-connection; pin the pool to one so schema
	// and rows stay on the same connection.
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatalf("EnsureAnalyticsTables: %v", err)
	}
	return db
}

func TestAnalyticsRoundtrip(t *testing.T) {
	db := openAnalyticsTestDB(t)
	now := int64(1_750_000_000)
	// Visitor A: two pageviews on two days, one click. Visitor B: one pageview.
	if err := InsertEvent(db, Event{TS: now - 10*86400, VisitorID: "a", SessionID: "s1", Kind: "pageview", Path: "/", Source: "Facebook", UA: "x", IPHash: "h"}); err != nil {
		t.Fatal(err)
	}
	if err := InsertEvent(db, Event{TS: now, VisitorID: "a", SessionID: "s2", Kind: "pageview", Path: "/prices", Source: "Facebook", UA: "x", IPHash: "h"}); err != nil {
		t.Fatal(err)
	}
	if err := InsertEvent(db, Event{TS: now, VisitorID: "a", SessionID: "s2", Kind: "click", Path: "/prices", Source: "Facebook", Element: "A", Label: "Call now", UA: "x", IPHash: "h"}); err != nil {
		t.Fatal(err)
	}
	if err := InsertEvent(db, Event{TS: now, VisitorID: "b", SessionID: "s3", Kind: "pageview", Path: "/", Source: "Direct", UA: "x", IPHash: "h2"}); err != nil {
		t.Fatal(err)
	}

	st, err := QueryAnalytics(db, 30, now)
	if err != nil {
		t.Fatalf("QueryAnalytics: %v", err)
	}
	if st.Pageviews != 3 {
		t.Errorf("Pageviews = %d, want 3", st.Pageviews)
	}
	if st.Visitors != 2 {
		t.Errorf("Visitors = %d, want 2", st.Visitors)
	}
	if st.Sessions != 3 {
		t.Errorf("Sessions = %d, want 3", st.Sessions)
	}
	if st.NewVisitors != 2 || st.ReturningVisitors != 0 {
		t.Errorf("New/Returning = %d/%d, want 2/0", st.NewVisitors, st.ReturningVisitors)
	}
	if len(st.TopPages) != 2 || st.TopPages[0].Path != "/" || st.TopPages[0].Count != 2 {
		t.Errorf("TopPages = %+v, want / first with 2", st.TopPages)
	}
	if len(st.TopSources) != 2 {
		t.Errorf("TopSources = %+v, want 2 sources", st.TopSources)
	}
	if len(st.TopClicks) != 1 || st.TopClicks[0].Label != "Call now" {
		t.Errorf("TopClicks = %+v, want the Call now click", st.TopClicks)
	}
	if len(st.Recent) != 4 {
		t.Errorf("Recent = %d rows, want 4", len(st.Recent))
	}
	// Returning-visitor math: shrink the window so visitor A was seen before it.
	st2, err := QueryAnalytics(db, 7, now)
	if err != nil {
		t.Fatalf("QueryAnalytics(7): %v", err)
	}
	if st2.Pageviews != 2 || st2.Visitors != 2 {
		t.Errorf("7d: pageviews=%d visitors=%d, want 2/2", st2.Pageviews, st2.Visitors)
	}
	if st2.ReturningVisitors != 1 || st2.NewVisitors != 1 {
		t.Errorf("7d: new=%d returning=%d, want 1/1", st2.NewVisitors, st2.ReturningVisitors)
	}
}

func TestInsertEventCaps(t *testing.T) {
	db := openAnalyticsTestDB(t)
	long := ""
	for i := 0; i < 3000; i++ {
		long += "x"
	}
	if err := InsertEvent(db, Event{TS: 1, Kind: "pageview", Path: long, Referrer: long, Label: long, UA: long}); err != nil {
		t.Fatal(err)
	}
	var path, ref, label, ua string
	if err := db.QueryRow(`SELECT path, referrer, label, ua FROM analytics_events`).Scan(&path, &ref, &label, &ua); err != nil {
		t.Fatal(err)
	}
	if len([]rune(path)) != 512 || len([]rune(label)) != 512 || len([]rune(ua)) != 512 {
		t.Errorf("text caps not applied: path=%d label=%d ua=%d", len([]rune(path)), len([]rune(label)), len([]rune(ua)))
	}
	if len([]rune(ref)) != 2048 {
		t.Errorf("referrer cap not applied: %d", len([]rune(ref)))
	}
}
