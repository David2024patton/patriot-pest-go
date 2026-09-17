package marketing

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	adminTestEmail    = "owner@example.com"
	adminTestPassword = "correct horse battery staple"
	browserUA         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0 Safari/537.36"
)

var csrfRe = regexp.MustCompile(`<input type="hidden" name="_csrf" value="([0-9a-f]{64})">`)

// adminRouter builds a marketing router with the admin console configured.
func adminRouter(t *testing.T) (*Module, chi.Router) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(adminTestPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADMIN_EMAILS", adminTestEmail)
	t.Setenv("ADMIN_PASSWORD_HASH", string(hash))
	m := &Module{Enabled: true, DBPath: t.TempDir() + "/analytics.db"}
	r := chi.NewRouter()
	if !m.Register(r) {
		t.Fatal("Register returned false")
	}
	return m, r
}

// loginForm GETs /admin/login and returns the CSRF cookie + token.
func loginForm(t *testing.T, r chi.Router) (*http.Cookie, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /admin/login = %d, want 200", rec.Code)
	}
	var csrfCookie *http.Cookie
	for _, c := range rec.Result().Cookies() {
		if c.Name == csrfCookieName {
			csrfCookie = c
		}
	}
	if csrfCookie == nil {
		t.Fatal("no CSRF cookie set")
	}
	mm := csrfRe.FindStringSubmatch(rec.Body.String())
	if mm == nil {
		t.Fatal("no CSRF token in login form")
	}
	return csrfCookie, mm[1]
}

// tryLogin POSTs the login form and returns the recorder.
func tryLogin(t *testing.T, r chi.Router, email, password string) *httptest.ResponseRecorder {
	t.Helper()
	csrfCookie, token := loginForm(t, r)
	form := url.Values{"email": {email}, "password": {password}, "_csrf": {token}}
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(csrfCookie)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestAdminDisabledWhenUnconfigured(t *testing.T) {
	t.Setenv("ADMIN_EMAILS", "")
	t.Setenv("ADMIN_PASSWORD_HASH", "")
	m := &Module{Enabled: true, DBPath: t.TempDir() + "/x.db"}
	r := chi.NewRouter()
	if !m.Register(r) {
		t.Fatal("Register returned false")
	}
	for _, path := range []string{"/admin", "/admin/login"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d, want 404 when admin unconfigured", path, rec.Code)
		}
	}
}

func TestAdminLoginRejectsUnknownEmail(t *testing.T) {
	_, r := adminRouter(t)
	rec := tryLogin(t, r, "stranger@example.com", adminTestPassword)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("unknown email: status = %d, want 401", rec.Code)
	}
}

func TestAdminLoginRejectsWrongPassword(t *testing.T) {
	_, r := adminRouter(t)
	rec := tryLogin(t, r, adminTestEmail, "wrong password")
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("wrong password: status = %d, want 401", rec.Code)
	}
}

func TestAdminLoginAcceptsAndServesDashboard(t *testing.T) {
	_, r := adminRouter(t)

	// Seed one pageview through the real beacon endpoint.
	beaconBody := `{"vid":"v1","sid":"s1","path":"/prices-test","ref":"https://www.google.com/","utm_source":"","utm_medium":""}`
	breq := httptest.NewRequest(http.MethodPost, "/api/track/view", strings.NewReader(beaconBody))
	breq.Header.Set("Content-Type", "application/json")
	breq.Header.Set("User-Agent", browserUA)
	brec := httptest.NewRecorder()
	r.ServeHTTP(brec, breq)
	if brec.Code != http.StatusAccepted {
		t.Fatalf("beacon: status = %d, want 202", brec.Code)
	}

	// Unauthenticated dashboard access redirects to login.
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/admin/login" {
		t.Fatalf("unauthed /admin = %d -> %q, want 302 to /admin/login", rec.Code, rec.Header().Get("Location"))
	}

	// Correct credentials.
	lrec := tryLogin(t, r, adminTestEmail, adminTestPassword)
	if lrec.Code != http.StatusFound || lrec.Header().Get("Location") != "/admin" {
		t.Fatalf("login = %d -> %q, want 302 to /admin", lrec.Code, lrec.Header().Get("Location"))
	}
	var sessCookie *http.Cookie
	sawFlag := false
	for _, c := range lrec.Result().Cookies() {
		switch c.Name {
		case adminSessionCookie:
			sessCookie = c
			if !c.HttpOnly {
				t.Error("session cookie is not HttpOnly")
			}
		case adminFlagCookie:
			sawFlag = true
			if c.HttpOnly {
				t.Error("tracker-exclusion cookie must be readable by tracker.js (not HttpOnly)")
			}
		}
	}
	if sessCookie == nil {
		t.Fatal("no session cookie set on login")
	}
	if !sawFlag {
		t.Error("tracker-exclusion cookie (ppc_admin) not set on login")
	}

	// Authenticated dashboard renders and shows the seeded pageview.
	dreq := httptest.NewRequest(http.MethodGet, "/admin?days=30", nil)
	dreq.AddCookie(sessCookie)
	drec := httptest.NewRecorder()
	r.ServeHTTP(drec, dreq)
	if drec.Code != http.StatusOK {
		t.Fatalf("authed /admin = %d, want 200", drec.Code)
	}
	body := drec.Body.String()
	for _, want := range []string{"PATRIOT ANALYTICS", "/prices-test", "Google", "TRAFFIC SOURCES", "TOP CLICKS"} {
		if !strings.Contains(body, want) {
			t.Errorf("dashboard missing %q", want)
		}
	}

	// Logout kills the session: dashboard redirects to login again.
	oreq := httptest.NewRequest(http.MethodGet, "/admin/logout", nil)
	oreq.AddCookie(sessCookie)
	orec := httptest.NewRecorder()
	r.ServeHTTP(orec, oreq)
	if oreq2 := httptest.NewRequest(http.MethodGet, "/admin", nil); true {
		oreq2.AddCookie(sessCookie)
		rec2 := httptest.NewRecorder()
		r.ServeHTTP(rec2, oreq2)
		if rec2.Code != http.StatusFound {
			t.Errorf("after logout /admin = %d, want 302", rec2.Code)
		}
	}
	if orec.Code != http.StatusFound {
		t.Errorf("logout = %d, want 302", orec.Code)
	}
}

func TestAdminDashboardRejectsBadRange(t *testing.T) {
	_, r := adminRouter(t)
	lrec := tryLogin(t, r, adminTestEmail, adminTestPassword)
	var sessCookie *http.Cookie
	for _, c := range lrec.Result().Cookies() {
		if c.Name == adminSessionCookie {
			sessCookie = c
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/admin?days=banana", nil)
	req.AddCookie(sessCookie)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/admin?days=banana = %d, want 200 (defaults to 30)", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "LAST 30 DAYS") {
		t.Error("bad days param did not default to 30")
	}
}
