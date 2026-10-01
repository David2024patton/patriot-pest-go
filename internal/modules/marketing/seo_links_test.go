package marketing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	"github.com/David2024patton/patriot-pest-go/internal/view"
)

var (
	hrefRe    = regexp.MustCompile(`href="([^"]+)"`)
	assetRe   = regexp.MustCompile(`content="(/assets/[^"]+)"`)
	locRe     = regexp.MustCompile(`<loc>([^<]+)</loc>`)
	lastmodRe = regexp.MustCompile(`<lastmod>(\d{4}-\d{2}-\d{2})</lastmod>`)
	titleRe   = regexp.MustCompile(`<title>([^<]*)</title>`)
	ogImageRe = regexp.MustCompile(`<meta property="og:image" content="([^"]+)"`)
)

// getOK GETs path (same Host discipline as seo_test.go) and fails on non-2xx/3xx.
func getOK(t *testing.T, r interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Host = "www.patriotpest.pro"
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

// TestSitemapXMLNoPriorityHasLastmod — the SEO pass removed cargo-cult
// <priority> and added honest <lastmod> per URL.
func TestSitemapXMLNoPriorityHasLastmod(t *testing.T) {
	r := seoRouter(t)
	rec := getOK(t, r, "/sitemap.xml")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /sitemap.xml = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<priority>") {
		t.Error("sitemap.xml still contains <priority> tags")
	}
	urlBlocks := strings.Count(body, "<url>")
	lastmods := lastmodRe.FindAllStringSubmatch(body, -1)
	if urlBlocks == 0 {
		t.Fatal("sitemap.xml has no <url> entries")
	}
	if len(lastmods) != urlBlocks {
		t.Errorf("sitemap.xml: %d <url> entries but %d <lastmod> dates", urlBlocks, len(lastmods))
	}
	for _, u := range []string{"/privacy-policy", "/terms-of-use"} {
		if !strings.Contains(body, "<loc>https://www.patriotpest.pro"+u+"</loc>") {
			t.Errorf("sitemap.xml missing %s", u)
		}
	}
}

// TestBreadcrumbJSONLD — pest and area pages emit a valid BreadcrumbList.
func TestBreadcrumbJSONLD(t *testing.T) {
	r := seoRouter(t)
	for path, wantNames := range map[string][]string{
		"/pest/ants":     {"Home", "Pest Library", "Ants"},
		"/areas/spokane": {"Home", "Service Areas", "Spokane"},
	} {
		rec := getOK(t, r, path)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s = %d, want 200", path, rec.Code)
		}
		var found bool
		for _, b := range ldRe.FindAllStringSubmatch(rec.Body.String(), -1) {
			var v map[string]any
			if err := json.Unmarshal([]byte(b[1]), &v); err != nil {
				t.Fatalf("GET %s: JSON-LD does not parse: %v", path, err)
			}
			if v["@type"] != "BreadcrumbList" {
				continue
			}
			found = true
			items, _ := v["itemListElement"].([]any)
			if len(items) != len(wantNames) {
				t.Fatalf("GET %s: BreadcrumbList has %d items, want %d", path, len(items), len(wantNames))
			}
			for i, it := range items {
				m, _ := it.(map[string]any)
				if m["name"] != wantNames[i] {
					t.Errorf("GET %s: item %d name = %v, want %q", path, i, m["name"], wantNames[i])
				}
				if int(m["position"].(float64)) != i+1 {
					t.Errorf("GET %s: item %d position = %v, want %d", path, i, m["position"], i+1)
				}
				if s, _ := m["item"].(string); !strings.HasPrefix(s, "https://www.patriotpest.pro/") {
					t.Errorf("GET %s: item %d has non-absolute URL %q", path, i, s)
				}
			}
		}
		if !found {
			t.Errorf("GET %s: no BreadcrumbList JSON-LD block found", path)
		}
	}
}

// TestWaspsTitle — the wasp-specific page must not carry the "Wasps & Hornets"
// catalog name in its <title> (hornets have their own page).
func TestWaspsTitle(t *testing.T) {
	r := seoRouter(t)
	rec := getOK(t, r, "/pest/wasps")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /pest/wasps = %d, want 200", rec.Code)
	}
	m := titleRe.FindStringSubmatch(rec.Body.String())
	if m == nil {
		t.Fatal("no <title> on /pest/wasps")
	}
	if m[1] != "Wasp Control | Patriot Pest Control" {
		t.Errorf("wasps title = %q, want %q", m[1], "Wasp Control | Patriot Pest Control")
	}
	// The sibling page keeps its own title.
	rec2 := getOK(t, r, "/pest/hornets")
	m2 := titleRe.FindStringSubmatch(rec2.Body.String())
	if m2 == nil || m2[1] != "Hornets Control | Patriot Pest Control" {
		t.Errorf("hornets title = %v, want it unchanged", m2)
	}
}

// TestOGImagePerSection — pest, area, and blog pages each get their branded
// OG image; everything else keeps the default.
func TestOGImagePerSection(t *testing.T) {
	r := seoRouter(t)
	for path, want := range map[string]string{
		"/pest/ants":                       "og-pest.png",
		"/areas/spokane":                   "og-area.png",
		"/blogs/why-ants-invade-in-spring": "og-blog.png",
		"/":                                "og.png",
	} {
		rec := getOK(t, r, path)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s = %d, want 200", path, rec.Code)
		}
		m := ogImageRe.FindStringSubmatch(rec.Body.String())
		if m == nil {
			t.Fatalf("GET %s: no og:image meta", path)
		}
		if !strings.HasSuffix(m[1], "/assets/img/"+want) {
			t.Errorf("GET %s: og:image = %q, want suffix /assets/img/%s", path, m[1], want)
		}
	}
}

// TestServiceSchemaKeepsProviderID — the pre-existing Service schema on pest
// pages (with provider @id cross-link) must survive the SEO pass untouched.
func TestServiceSchemaKeepsProviderID(t *testing.T) {
	r := seoRouter(t)
	rec := getOK(t, r, "/pest/mosquitoes")
	var found bool
	for _, b := range ldRe.FindAllStringSubmatch(rec.Body.String(), -1) {
		var v map[string]any
		if err := json.Unmarshal([]byte(b[1]), &v); err != nil {
			t.Fatalf("JSON-LD does not parse: %v", err)
		}
		if v["@type"] != "Service" {
			continue
		}
		found = true
		prov, _ := v["provider"].(map[string]any)
		if prov["@id"] != "https://www.patriotpest.pro/#business" {
			t.Errorf("Service provider @id = %v, want the business cross-link", prov["@id"])
		}
		if v["serviceType"] != "Mosquitoes Control" {
			t.Errorf("Service serviceType = %v", v["serviceType"])
		}
	}
	if !found {
		t.Error("no Service JSON-LD on /pest/mosquitoes")
	}
}

// TestInternalLinkTargetsResolve — every entry in the code-level link tables
// must point at a route that exists: no invented slugs.
func TestInternalLinkTargetsResolve(t *testing.T) {
	staticOK := map[string]bool{"/contact": true, "/services": true}
	for postSlug, links := range view.BlogRelatedLinks {
		if _, ok := data.PublishedPostBySlug(postSlug); !ok {
			t.Errorf("BlogRelatedLinks: unknown post slug %q", postSlug)
		}
		for _, l := range links {
			href := l[1]
			switch {
			case strings.HasPrefix(href, "/pest/"):
				if _, ok := data.PestBySlug(strings.TrimPrefix(href, "/pest/")); !ok {
					t.Errorf("BlogRelatedLinks[%q]: unknown pest target %q", postSlug, href)
				}
			case strings.HasPrefix(href, "/areas/"):
				if _, _, _, ok := data.FindCity(strings.TrimPrefix(href, "/areas/")); !ok {
					t.Errorf("BlogRelatedLinks[%q]: unknown city target %q", postSlug, href)
				}
			default:
				if !staticOK[href] {
					t.Errorf("BlogRelatedLinks[%q]: unknown static target %q", postSlug, href)
				}
			}
		}
	}
	for code, slugs := range view.AreaPestLinks {
		if len(slugs) == 0 {
			t.Errorf("AreaPestLinks[%q]: empty", code)
		}
		for _, slug := range slugs {
			if _, ok := data.PestBySlug(slug); !ok {
				t.Errorf("AreaPestLinks[%q]: unknown pest slug %q", code, slug)
			}
		}
		if got := view.AreaPestsFor(code); len(got) != len(slugs) {
			t.Errorf("AreaPestsFor(%q) resolved %d of %d slugs", code, len(got), len(slugs))
		}
	}
}

// TestBrokenLinkCrawl — crawl every sitemap URL plus every internal href
// found while rendering, and require each to resolve (200 or redirect).
// This is the broken-link pass as a regression test.
func TestBrokenLinkCrawl(t *testing.T) {
	r := seoRouter(t)

	// Seed the crawl with every sitemap URL.
	rec := getOK(t, r, "/sitemap.xml")
	var queue []string
	for _, m := range locRe.FindAllStringSubmatch(rec.Body.String(), -1) {
		p := strings.TrimPrefix(m[1], "https://www.patriotpest.pro")
		queue = append(queue, p)
	}
	// Plus the HTML sitemap page, which links everything.
	queue = append(queue, "/sitemap")

	seen := map[string]bool{}
	statusOf := map[string]int{}
	referrer := map[string]string{}
	okStatus := func(code int) bool {
		return code == http.StatusOK || code == http.StatusMovedPermanently || code == http.StatusFound
	}
	const maxVisits = 600
	for len(queue) > 0 && len(seen) < maxVisits {
		path := queue[0]
		queue = queue[1:]
		if seen[path] {
			continue
		}
		seen[path] = true
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Host = "www.patriotpest.pro"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		statusOf[path] = w.Code
		if w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound {
			if loc := w.Header().Get("Location"); strings.HasPrefix(loc, "/") && !seen[loc] {
				referrer[loc] = path + " (redirect)"
				queue = append(queue, loc)
			}
			continue
		}
		if w.Code != http.StatusOK || !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
			continue
		}
		for _, m := range hrefRe.FindAllStringSubmatch(w.Body.String(), -1) {
			addLink(t, path, m[1], seen, &queue, referrer)
		}
		for _, m := range assetRe.FindAllStringSubmatch(w.Body.String(), -1) {
			addLink(t, path, m[1], seen, &queue, referrer)
		}
	}
	var bad []string
	for path, code := range statusOf {
		if !okStatus(code) {
			bad = append(bad, fmt.Sprintf("%s -> %d (linked from %s)", path, code, referrer[path]))
		}
	}
	if len(bad) > 0 {
		t.Errorf("broken internal links:\n%s", strings.Join(bad, "\n"))
	}
	t.Logf("crawled %d internal URLs, all resolved", len(seen))
}

// addLink queues an internal href for the crawl; external/tel/anchor links
// are out of scope.
func addLink(t *testing.T, from, href string, seen map[string]bool, queue *[]string, referrer map[string]string) {
	t.Helper()
	if href == "" || strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "javascript:") {
		return
	}
	var path string
	switch {
	case strings.HasPrefix(href, "/") && !strings.HasPrefix(href, "//"):
		path = href
	case strings.HasPrefix(href, "https://www.patriotpest.pro/"):
		path = strings.TrimPrefix(href, "https://www.patriotpest.pro")
	default:
		return // external or relative — out of scope
	}
	if i := strings.Index(path, "#"); i >= 0 {
		path = path[:i]
	}
	if i := strings.Index(path, "?"); i >= 0 {
		path = path[:i]
	}
	if path == "" {
		path = "/"
	}
	if !seen[path] {
		if _, ok := referrer[path]; !ok {
			referrer[path] = from
		}
		*queue = append(*queue, path)
	}
}
