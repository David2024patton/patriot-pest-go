package report

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestIsReportHost(t *testing.T) {
	ok := []string{
		"report.patriotpest.pro",
		"REPORT.PATRIOTPEST.PRO",
		"Report.PatriotPest.Pro",
		"report.patriotpest.pro:443",
		"report.patriotpest.pro:80",
	}
	for _, h := range ok {
		if !isReportHost(h) {
			t.Errorf("isReportHost(%q) = false, want true", h)
		}
	}
	bad := []string{
		"patriotpest.pro",
		"www.patriotpest.pro",
		"www.report.patriotpest.pro",
		"evilreport.patriotpest.pro",
		"report.patriotpest.pro.evil.com",
		"",
	}
	for _, h := range bad {
		if isReportHost(h) {
			t.Errorf("isReportHost(%q) = true, want false", h)
		}
	}
}

func TestEmailAllowed(t *testing.T) {
	ok := []string{
		"skyler@patriotpest.pro",
		"Skyler@PatriotPest.Pro ",
		" AVA@PATRIOTPEST.PRO",
		"david@itak.live",
		"David@ITAK.Live",
	}
	for _, e := range ok {
		if !emailAllowed(e) {
			t.Errorf("emailAllowed(%q) = false, want true", e)
		}
	}
	bad := []string{
		"",
		"notanemail",
		"x@gmail.com",
		"x@patriotpest.pro.evil.com",
		"x@evipatriotpest.pro",
		"@patriotpest.pro",
		"skyler@patriotpest.pro\nBcc:evil@x.com",
		"david@itak.live.evil.com",
	}
	for _, e := range bad {
		if emailAllowed(e) {
			t.Errorf("emailAllowed(%q) = true, want false", e)
		}
	}
}

func TestGenCode(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		c, err := genCode()
		if err != nil {
			t.Fatal(err)
		}
		if len(c) != 6 {
			t.Fatalf("code %q is not 6 chars", c)
		}
		for _, ch := range c {
			if ch < '0' || ch > '9' {
				t.Fatalf("code %q is not all digits", c)
			}
		}
		seen[c] = true
	}
	if len(seen) < 40 {
		t.Errorf("codes look non-random: %d unique of 50", len(seen))
	}
}

func testModule(t *testing.T) *Module {
	t.Helper()
	m := &Module{DBPath: t.TempDir() + "/report_test.db"}
	m.ensureTables()
	m.otpReqLimiter = newKeyLimiter(1000, time.Hour)
	m.verifyLimiter = newKeyLimiter(1000, time.Hour)
	return m
}

func TestOTPFlow(t *testing.T) {
	m := testModule(t)
	db, err := m.db()
	if err != nil {
		t.Fatal(err)
	}
	email := "skyler@patriotpest.pro"
	code, err := genCode()
	if err != nil {
		t.Fatal(err)
	}
	if err := m.storeOTP(db, email, code); err != nil {
		t.Fatal(err)
	}
	if !m.checkOTP(db, email, code) {
		t.Fatal("valid code rejected")
	}
	if m.checkOTP(db, email, code) {
		t.Fatal("single-use code validated twice")
	}
	// Wrong code burns attempts; after the cap the real code is dead.
	code2, _ := genCode()
	if err := m.storeOTP(db, email, code2); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < otpMaxAttempts; i++ {
		if m.checkOTP(db, email, "000000") {
			t.Fatal("wrong code accepted")
		}
	}
	if m.checkOTP(db, email, code2) {
		t.Fatal("code validated after attempt cap")
	}
	// Expired code never validates.
	code3, _ := genCode()
	if err := m.storeOTP(db, email, code3); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE report_otps SET expires_at=? WHERE code_hash=?`,
		time.Now().Add(-time.Minute).Unix(), sha256Hex(code3)); err != nil {
		t.Fatal(err)
	}
	if m.checkOTP(db, email, code3) {
		t.Fatal("expired code validated")
	}
	// Code bound to another email does not validate.
	code4, _ := genCode()
	if err := m.storeOTP(db, email, code4); err != nil {
		t.Fatal(err)
	}
	if m.checkOTP(db, "other@patriotpest.pro", code4) {
		t.Fatal("code validated for wrong email")
	}
}

func TestSessionFlow(t *testing.T) {
	m := testModule(t)
	db, err := m.db()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/verify", nil)
	r.Header.Set("X-Forwarded-Proto", "https")
	if err := m.createSession(db, w, r, "skyler@patriotpest.pro"); err != nil {
		t.Fatal(err)
	}
	var sessCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == sessionCookieName {
			sessCookie = c
		}
	}
	if sessCookie == nil {
		t.Fatal("no session cookie set")
	}
	if !sessCookie.HttpOnly || !sessCookie.Secure || sessCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("weak cookie flags: %+v", sessCookie)
	}
	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	r2.AddCookie(sessCookie)
	email, ok := m.emailFromSession(r2)
	if !ok || email != "skyler@patriotpest.pro" {
		t.Fatalf("session lookup failed: %q %v", email, ok)
	}
	// Expired sessions do not validate.
	if _, err := db.Exec(`UPDATE report_sessions SET expires_at=?`,
		time.Now().Add(-time.Minute).Unix()); err != nil {
		t.Fatal(err)
	}
	if _, ok := m.emailFromSession(r2); ok {
		t.Fatal("expired session validated")
	}
}

func TestKeyLimiter(t *testing.T) {
	l := newKeyLimiter(3, time.Hour)
	for i := 0; i < 3; i++ {
		if !l.allow("k") {
			t.Fatalf("hit %d denied, want allow", i+1)
		}
	}
	if l.allow("k") {
		t.Fatal("4th hit allowed over limit 3")
	}
	if !l.allow("other") {
		t.Fatal("fresh key denied")
	}
	fast := newKeyLimiter(1, time.Millisecond)
	if !fast.allow("k") {
		t.Fatal("first hit denied")
	}
	time.Sleep(5 * time.Millisecond)
	if !fast.allow("k") {
		t.Fatal("hit after window reset denied")
	}
}

func TestNoEmDashes(t *testing.T) {
	for _, s := range []string{dashboardHTML, authCSS} {
		if strings.Contains(s, "—") || strings.Contains(s, "–") {
			t.Error("dashboard copy contains an em/en dash")
		}
	}
}

func TestMiddlewareHostGating(t *testing.T) {
	m := testModule(t)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	h := m.Middleware(inner)

	// Other hosts pass through untouched.
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Host = "www.patriotpest.pro"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusTeapot {
		t.Fatalf("main site host: got %d, want passthrough", w.Code)
	}

	// Report host without a session redirects to /login.
	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	r2.Host = "report.patriotpest.pro"
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	if w2.Code != http.StatusFound {
		t.Fatalf("report host unauthenticated: got %d, want 302 to /login", w2.Code)
	}
	if loc := w2.Header().Get("Location"); loc != "/login" {
		t.Fatalf("redirect location = %q, want /login", loc)
	}
	if tag := w2.Header().Get("X-Robots-Tag"); tag != "noindex" {
		t.Fatalf("X-Robots-Tag = %q, want noindex", tag)
	}

	// Report host login page carries noindex and renders the form.
	r3 := httptest.NewRequest(http.MethodGet, "/login", nil)
	r3.Host = "report.patriotpest.pro"
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	if w3.Code != http.StatusOK {
		t.Fatalf("GET /login: got %d", w3.Code)
	}
	if tag := w3.Header().Get("X-Robots-Tag"); tag != "noindex" {
		t.Fatalf("login X-Robots-Tag = %q, want noindex", tag)
	}
	if !strings.Contains(w3.Body.String(), "SEND CODE") {
		t.Fatal("login page missing email form")
	}

	// Infra paths pass through on the report host too.
	for _, p := range []string{"/health", "/ready", "/metrics"} {
		rp := httptest.NewRequest(http.MethodGet, p, nil)
		rp.Host = "report.patriotpest.pro"
		wp := httptest.NewRecorder()
		h.ServeHTTP(wp, rp)
		if wp.Code != http.StatusTeapot {
			t.Fatalf("%s on report host: got %d, want passthrough", p, wp.Code)
		}
	}
}

func TestDashboardRenders(t *testing.T) {
	m := testModule(t)
	db, err := m.db()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/verify", nil)
	if err := m.createSession(db, w, r, "skyler@patriotpest.pro"); err != nil {
		t.Fatal(err)
	}
	var sessCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == sessionCookieName {
			sessCookie = c
		}
	}
	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	r2.Host = "report.patriotpest.pro"
	r2.AddCookie(sessCookie)
	w2 := httptest.NewRecorder()
	m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})).ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("dashboard: got %d", w2.Code)
	}
	body := w2.Body.String()
	for _, want := range []string{
		"Patriot Pest Control", "AlphaFlux", "October 1, 2026",
		"Week of October 8", "skyler@patriotpest.pro",
		"Website rebuild", "Google Business Profiles", "Sameday AI phone agent",
		"The $4,000 proposal", "Up next", "Everything done", "What is left",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("dashboard missing %q", want)
		}
	}
	if ct := w2.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("content-type = %q", ct)
	}
}

func TestRegisterAddsMiddleware(t *testing.T) {
	m := &Module{DBPath: t.TempDir() + "/reg_test.db"}
	r := chi.NewRouter()
	m.Register(r)
	// chi skips the middleware stack entirely when the router has zero
	// routes; one dummy route makes the test exercise the real path.
	r.Get("/dummy", func(w http.ResponseWriter, r *http.Request) {})
	// Walk the middleware stack indirectly: a report-host request must not 404.
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Host = "report.patriotpest.pro"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("after Register, GET /login on report host = %d", w.Code)
	}
}
