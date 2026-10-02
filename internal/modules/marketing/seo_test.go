package marketing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// seoRouter builds a marketing router without admin configured (public pages
// only), for SEO regression tests.
func seoRouter(t *testing.T) chi.Router {
	t.Helper()
	m := &Module{Enabled: true, DBPath: t.TempDir() + "/seo.db"}
	r := chi.NewRouter()
	if !m.Register(r) {
		t.Fatal("Register returned false")
	}
	return r
}

var ldRe = regexp.MustCompile(`<script type="application/ld\+json">([\s\S]*?)</script>`)

// TestJSONLDParsesAsObject guards the layout's JSON-LD against the
// html/template double-encoding bug: every ld+json block must parse as a JSON
// object, never as a quoted string.
func TestJSONLDParsesAsObject(t *testing.T) {
	r := seoRouter(t)
	for _, path := range []string{"/", "/services", "/contact", "/service-areas", "/blogs"} {
		req := httptest.NewRequest("GET", path, nil)
		req.Host = "www.patriotpest.pro"
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: status %d", path, rec.Code)
		}
		blocks := ldRe.FindAllStringSubmatch(rec.Body.String(), -1)
		if len(blocks) == 0 {
			t.Fatalf("GET %s: no JSON-LD blocks found", path)
		}
		for _, b := range blocks {
			var v any
			if err := json.Unmarshal([]byte(b[1]), &v); err != nil {
				t.Fatalf("GET %s: JSON-LD block does not parse: %v", path, err)
			}
			obj, ok := v.(map[string]any)
			if !ok {
				t.Fatalf("GET %s: JSON-LD block is not an object (double-encoded?): %T", path, v)
			}
			if obj["@context"] != "https://schema.org" {
				t.Errorf("GET %s: JSON-LD missing schema.org @context", path)
			}
		}
	}
}

// TestLegalShortURLsRedirect verifies the /privacy and /terms aliases 301 to
// the canonical legal pages.
func TestLegalShortURLsRedirect(t *testing.T) {
	r := seoRouter(t)
	for from, to := range map[string]string{
		"/privacy": "/privacy-policy",
		"/terms":   "/terms-of-use",
	} {
		req := httptest.NewRequest("GET", from, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusMovedPermanently {
			t.Errorf("GET %s: status %d, want 301", from, rec.Code)
		}
		if loc := rec.Header().Get("Location"); loc != to {
			t.Errorf("GET %s: Location %q, want %q", from, loc, to)
		}
	}
}

// TestLegalPagesServe verifies the canonical legal pages render with real
// content (not thin stubs).
func TestLegalPagesServe(t *testing.T) {
	r := seoRouter(t)
	for path, needle := range map[string]string{
		"/privacy-policy": "Patriot Pest Control CO",
		"/terms-of-use":   "Patriot Pest Control CO",
	} {
		req := httptest.NewRequest("GET", path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s: status %d, want 200", path, rec.Code)
		}
		if body := rec.Body.String(); len(body) < 1000 || !strings.Contains(body, needle) {
			t.Errorf("GET %s: legal page looks thin or missing content", path)
		}
	}
}
