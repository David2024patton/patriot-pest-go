// Package report serves the weekly status dashboard on the
// report.patriotpest.pro host, gated by an email allowlist login.
//
// Only addresses at the patriotpest.pro domain may sign in (plus David's
// itak.live admin address so he cannot lock himself out). Typing an allowed
// email signs straight in with no code step, by owner decision. Sessions
// mirror the /admin hardening: random token, hash-only storage, 24h expiry,
// Secure HttpOnly SameSite cookies, CSRF on POSTs, per-IP throttling. Raw IPs
// are never persisted; the analytics beacon is not rendered on this host.
package report

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	"github.com/David2024patton/patriot-pest-go/internal/view"
)

// reportHost is the only host that serves the dashboard. Matching is
// case-insensitive and ignores any port suffix.
const reportHost = "report.patriotpest.pro"

const (
	sessionCookieName = "ppc_report_sess"
	sessionTTL        = 24 * time.Hour
)

// infraPaths pass through to the normal router on every host, including the
// report host, so health checks and metrics never get gated.
var infraPaths = map[string]bool{
	"/health":  true,
	"/ready":   true,
	"/metrics": true,
}

// Module serves the report dashboard. DBPath points at the analytics SQLite
// file, which also carries the report_sessions table.
type Module struct {
	DBPath string

	loginLimiter *keyLimiter // 10 login attempts per IP per hour
}

// Register adds the host-gating middleware. On any other host the middleware
// is a no-op and the existing site is completely untouched.
func (m *Module) Register(r chi.Router) {
	m.ensureTables()
	if m.loginLimiter == nil {
		m.loginLimiter = newKeyLimiter(10, time.Hour)
	}
	r.Use(m.Middleware)
}

// Middleware intercepts report.patriotpest.pro and serves the report router.
// Everything else falls through to the normal site.
func (m *Module) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isReportHost(r.Host) {
			next.ServeHTTP(w, r)
			return
		}
		if infraPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}
		m.route(w, r)
	})
}

// isReportHost matches report.patriotpest.pro case-insensitively, ignoring
// any port suffix. Subdomains of it and lookalikes do not match.
func isReportHost(host string) bool {
	h := strings.ToLower(host)
	if i := strings.LastIndex(h, ":"); i >= 0 {
		// Strip a port suffix, but not an IPv6 literal (bracketed).
		if !strings.HasPrefix(h, "[") {
			h = h[:i]
		}
	}
	return h == reportHost
}

// route dispatches inside the report host. Every response carries
// X-Robots-Tag: noindex, including the login pages.
func (m *Module) route(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Robots-Tag", "noindex")
	switch {
	case r.URL.Path == "/login" && r.Method == http.MethodGet:
		m.loginGet(w, r)
	case r.URL.Path == "/login" && r.Method == http.MethodPost:
		m.loginPost(w, r)
	case r.URL.Path == "/logout":
		m.logout(w, r)
	case r.URL.Path == "/login" || r.URL.Path == "/verify":
		http.Redirect(w, r, "/login", http.StatusFound)
	default:
		m.dashboard(w, r)
	}
}

// ---- email allowlist ----

// emailAllowed reports whether the address may request a sign-in code: the
// domain must be exactly patriotpest.pro, or the full address must be
// david@itak.live (David's admin exception).
func emailAllowed(raw string) bool {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" || strings.ContainsAny(email, " \t\r\n") {
		return false
	}
	if email == "david@itak.live" {
		return true
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return parts[1] == "patriotpest.pro"
}

// normalizeEmail lowercases and trims; used for limiter keys and storage.
func normalizeEmail(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

// ---- client IP (never persisted) ----

// clientIP prefers X-Forwarded-For (the app sits behind Traefik in Dokploy)
// and falls back to the direct remote address, port stripped. Used only for
// in-memory rate limiting; raw IPs are never written to storage.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.Index(xff, ","); i >= 0 {
			xff = xff[:i]
		}
		if ip := strings.TrimSpace(xff); ip != "" {
			return ip
		}
	}
	if h, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return h
	}
	return r.RemoteAddr
}

// ---- keyed rate limiter (in-memory only) ----

type keyLimiter struct {
	mu     sync.Mutex
	count  map[string]int
	reset  map[string]time.Time
	limit  int
	window time.Duration
}

func newKeyLimiter(limit int, window time.Duration) *keyLimiter {
	return &keyLimiter{count: map[string]int{}, reset: map[string]time.Time{}, limit: limit, window: window}
}

// allow records one hit for key and reports whether it is within the limit.
func (l *keyLimiter) allow(key string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if now.After(l.reset[key]) {
		l.count[key] = 0
		l.reset[key] = now.Add(l.window)
	}
	l.count[key]++
	n := l.count[key]
	for k, exp := range l.reset {
		if now.After(exp) {
			delete(l.reset, k)
			delete(l.count, k)
		}
	}
	return n <= l.limit
}

// ---- storage ----

func (m *Module) ensureTables() {
	db, err := data.AnalyticsDB(m.DBPath)
	if err != nil {
		slog.Error("report: analytics db unavailable", "err", err.Error())
		return
	}
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS report_sessions (
			token_hash TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			expires_at INTEGER NOT NULL
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			slog.Error("report: schema ensure failed", "err", err.Error())
		}
	}
}

func (m *Module) db() (*sql.DB, error) {
	return data.AnalyticsDB(m.DBPath)
}

// sha256Hex hashes tokens for storage. Raw values only ever live
// in the cookie, never in the database.
func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// ---- sessions (mirrors /admin hardening) ----

func (m *Module) createSession(db *sql.DB, w http.ResponseWriter, r *http.Request, email string) error {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	token := hex.EncodeToString(b)
	now := time.Now()
	if _, err := db.Exec(`INSERT INTO report_sessions (token_hash, email, created_at, expires_at)
		VALUES (?, ?, ?, ?)`, sha256Hex(token), email, now.Unix(), now.Add(sessionTTL).Unix()); err != nil {
		return err
	}
	secure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookieName, Value: token, Path: "/",
		MaxAge: int(sessionTTL.Seconds()), HttpOnly: true,
		Secure: secure, SameSite: http.SameSiteLaxMode,
	})
	return nil
}

// emailFromSession validates the session cookie and returns the signed-in
// email, or false when the session is missing or expired.
func (m *Module) emailFromSession(r *http.Request) (string, bool) {
	c, err := r.Cookie(sessionCookieName)
	if err != nil || c.Value == "" {
		return "", false
	}
	db, err := m.db()
	if err != nil {
		return "", false
	}
	var email string
	err = db.QueryRow(`SELECT email FROM report_sessions WHERE token_hash=? AND expires_at>?`,
		sha256Hex(c.Value), time.Now().Unix()).Scan(&email)
	if err != nil {
		return "", false
	}
	return email, true
}

func (m *Module) destroySession(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
		if db, err := m.db(); err == nil {
			_, _ = db.Exec(`DELETE FROM report_sessions WHERE token_hash=?`, sha256Hex(c.Value))
		}
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
}

func purgeExpiredReportSessions(db *sql.DB) {
	if _, err := db.Exec(`DELETE FROM report_sessions WHERE expires_at<=?`, time.Now().Unix()); err != nil {
		slog.Warn("report: session purge failed", "err", err.Error())
	}
}

// ---- handlers ----

// loginGet renders the email form. Already signed in -> straight to the dashboard.
func (m *Module) loginGet(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.emailFromSession(r); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
		"Weekly reports for the Patriot Pest Control team.",
		`Enter your Patriot Pest email to sign in.`,
		"", emailFormHTML("", w, r))
}

// loginPost validates the email against the allowlist and, when allowed,
// signs straight in with a session. No code step, by owner decision.
func (m *Module) loginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "", "Could not read the form. Try again.", emailFormHTML("", w, r))
		return
	}
	if !view.VerifyCSRF(r) {
		w.WriteHeader(http.StatusForbidden)
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "", "Security token expired. Reload the page and try again.", emailFormHTML("", w, r))
		return
	}
	if !m.loginLimiter.allow("login:" + clientIP(r)) {
		w.WriteHeader(http.StatusTooManyRequests)
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "", "Too many attempts. Wait a bit and try again.", emailFormHTML("", w, r))
		return
	}
	email := normalizeEmail(r.FormValue("email"))
	if !emailAllowed(email) {
		slog.Warn("report: unauthorized login attempt", "email", email)
		w.WriteHeader(http.StatusForbidden)
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "", "This email is not authorized for these reports.", emailFormHTML("", w, r))
		return
	}
	db, err := m.db()
	if err != nil {
		slog.Error("report: db unavailable", "err", err.Error())
		http.Error(w, "Could not sign in. Try again later.", http.StatusInternalServerError)
		return
	}
	purgeExpiredReportSessions(db)
	if err := m.createSession(db, w, r, email); err != nil {
		slog.Error("report: session create failed", "err", err.Error())
		http.Error(w, "Could not sign in. Try again later.", http.StatusInternalServerError)
		return
	}
	slog.Info("report: login", "email", email)
	http.Redirect(w, r, "/", http.StatusFound)
}

// logout destroys the session and returns to the login page.
func (m *Module) logout(w http.ResponseWriter, r *http.Request) {
	m.destroySession(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// dashboard serves the status dashboard to signed-in users only.
func (m *Module) dashboard(w http.ResponseWriter, r *http.Request) {
	email, ok := m.emailFromSession(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	page := strings.Replace(dashboardHTML, "<!--USER_EMAIL-->", html.EscapeString(email), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}

// ---- auth page rendering ----

func emailFormHTML(email string, w http.ResponseWriter, r *http.Request) string {
	return `<form method="post" action="/login" novalidate>` +
		string(view.CSRFField(w, r)) +
		`<label for="email">EMAIL</label>` +
		`<input id="email" name="email" type="email" autocomplete="email" required autofocus value="` + html.EscapeString(email) + `">` +
		`<button type="submit">SIGN IN</button></form>`
}

// writeAuthPage renders a standalone auth page (no dashboard chrome).
func writeAuthPage(w http.ResponseWriter, title, heading, sub, intro, errMsg, formHTML string) {
	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8">` +
		`<meta name="viewport" content="width=device-width,initial-scale=1">` +
		`<meta name="robots" content="noindex,nofollow">` +
		`<title>` + html.EscapeString(title) + `</title><style>` + authCSS + `</style></head><body>`)
	sb.WriteString(`<div class="auth"><div class="card">`)
	sb.WriteString(`<h1>` + html.EscapeString(heading) + `</h1>`)
	if sub != "" {
		sb.WriteString(`<p class="sub">` + sub + `</p>`)
	}
	if intro != "" {
		sb.WriteString(`<p class="intro">` + html.EscapeString(intro) + `</p>`)
	}
	if errMsg != "" {
		sb.WriteString(`<div class="err">` + html.EscapeString(errMsg) + `</div>`)
	}
	sb.WriteString(formHTML)
	sb.WriteString(`</div></div></body></html>`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(sb.String()))
}
