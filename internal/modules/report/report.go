// Package report serves the weekly status dashboard on the
// report.patriotpest.pro host, gated by a passwordless email OTP login.
//
// Only addresses at the patriotpest.pro domain may sign in (plus David's
// itak.live admin address so he cannot lock himself out). Sessions mirror the
// /admin hardening: random token, hash-only storage, 24h expiry, Secure
// HttpOnly SameSite cookies, CSRF on POSTs, per-IP throttling. Raw IPs are
// never persisted; the analytics beacon is not rendered on this host.
package report

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	"errors"
	"html"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
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
	otpTTL            = 10 * time.Minute
	otpMaxAttempts    = 5
)

// infraPaths pass through to the normal router on every host, including the
// report host, so health checks and metrics never get gated.
var infraPaths = map[string]bool{
	"/health":  true,
	"/ready":   true,
	"/metrics": true,
}

var errSMTPNotConfigured = errors.New("report: REPORT_SMTP_HOST is not set")

// Module serves the report dashboard. DBPath points at the analytics SQLite
// file, which also carries the report_otps and report_sessions tables.
type Module struct {
	DBPath string

	otpReqLimiter *keyLimiter // 5 code requests per email per hour
	verifyLimiter *keyLimiter // 10 verify attempts per IP per hour
}

// Register adds the host-gating middleware. On any other host the middleware
// is a no-op and the existing site is completely untouched.
func (m *Module) Register(r chi.Router) {
	m.ensureTables()
	if m.otpReqLimiter == nil {
		m.otpReqLimiter = newKeyLimiter(5, time.Hour)
	}
	if m.verifyLimiter == nil {
		m.verifyLimiter = newKeyLimiter(10, time.Hour)
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
	case r.URL.Path == "/verify" && r.Method == http.MethodPost:
		m.verifyPost(w, r)
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
		`CREATE TABLE IF NOT EXISTS report_otps (
			code_hash TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			expires_at INTEGER NOT NULL,
			attempts INTEGER NOT NULL DEFAULT 0
		)`,
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

// sha256Hex hashes tokens and codes for storage. Raw values only ever live
// in the cookie or the email, never in the database.
func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// ---- OTP codes ----

// genCode returns a cryptographically random 6-digit code.
func genCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(n.Int64()+100000, 10), nil
}

// storeOTP replaces any pending code for the email with a fresh one.
func (m *Module) storeOTP(db *sql.DB, email, code string) error {
	now := time.Now()
	if _, err := db.Exec(`DELETE FROM report_otps WHERE email=?`, email); err != nil {
		return err
	}
	_, err := db.Exec(`INSERT INTO report_otps (code_hash, email, created_at, expires_at, attempts)
		VALUES (?, ?, ?, ?, 0)`, sha256Hex(code), email, now.Unix(), now.Add(otpTTL).Unix())
	return err
}

type otpRow struct {
	expiresAt int64
	attempts  int
}

// checkOTP validates a code for an email. At most one code is pending per
// email, and failures increment that row's attempt counter no matter which
// wrong code was tried, so the "5 attempts per code" cap cannot be dodged by
// guessing different values. A used, expired, or over-attempt code never
// validates twice.
func (m *Module) checkOTP(db *sql.DB, email, code string) bool {
	var row otpRow
	var hash string
	err := db.QueryRow(`SELECT code_hash, expires_at, attempts FROM report_otps WHERE email=?`,
		email).Scan(&hash, &row.expiresAt, &row.attempts)
	if err != nil {
		return false
	}
	fail := func() bool {
		row.attempts++
		if row.attempts >= otpMaxAttempts || time.Now().Unix() > row.expiresAt {
			_, _ = db.Exec(`DELETE FROM report_otps WHERE email=?`, email)
		} else {
			_, _ = db.Exec(`UPDATE report_otps SET attempts=? WHERE email=?`, row.attempts, email)
		}
		return false
	}
	if time.Now().Unix() > row.expiresAt {
		return fail()
	}
	if row.attempts >= otpMaxAttempts {
		return fail()
	}
	if subtle.ConstantTimeCompare([]byte(hash), []byte(sha256Hex(code))) != 1 {
		return fail()
	}
	_, _ = db.Exec(`DELETE FROM report_otps WHERE email=?`, email)
	return true
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
	if _, err := db.Exec(`DELETE FROM report_otps WHERE expires_at<=?`, time.Now().Unix()); err != nil {
		slog.Warn("report: otp purge failed", "err", err.Error())
	}
}

// ---- SMTP ----

// smtpConfigured reports whether the REPORT_SMTP_* env vars are present.
func smtpConfigured() bool {
	return strings.TrimSpace(os.Getenv("REPORT_SMTP_HOST")) != ""
}

// sendOTPEmail delivers the sign-in code. It returns errSMTPNotConfigured
// when the REPORT_SMTP_* env vars are unset.
func sendOTPEmail(to, code string) error {
	host := strings.TrimSpace(os.Getenv("REPORT_SMTP_HOST"))
	if host == "" {
		return errSMTPNotConfigured
	}
	port := strings.TrimSpace(os.Getenv("REPORT_SMTP_PORT"))
	if port == "" {
		port = "587"
	}
	user := os.Getenv("REPORT_SMTP_USER")
	pass := os.Getenv("REPORT_SMTP_PASS")
	from := strings.TrimSpace(os.Getenv("REPORT_SMTP_FROM"))
	if from == "" {
		from = user
	}
	body := "Your sign-in code is: " + code + "\r\n\r\n" +
		"It expires in 10 minutes. If you did not request this, you can ignore this email.\r\n"
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Your Patriot report sign-in code\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" + body)
	addr := net.JoinHostPort(host, port)
	if port == "465" {
		return sendImplicitTLS(addr, host, user, pass, from, to, msg)
	}
	return sendSTARTTLS(addr, host, user, pass, from, to, msg)
}

const smtpTimeout = 10 * time.Second

func smtpClient(addr, host string, implicitTLS bool) (*smtp.Client, net.Conn, error) {
	dialer := &net.Dialer{Timeout: smtpTimeout}
	var conn net.Conn
	var err error
	if implicitTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: host})
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return nil, nil, err
	}
	_ = conn.SetDeadline(time.Now().Add(smtpTimeout))
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	return c, conn, nil
}

func smtpDeliver(c *smtp.Client, conn net.Conn, host, user, pass, from, to string, msg []byte, startTLS bool) error {
	defer conn.Close()
	defer c.Close()
	if startTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
				return err
			}
		}
	}
	if user != "" {
		if ok, _ := c.Extension("AUTH"); ok {
			if err := c.Auth(smtp.PlainAuth("", user, pass, host)); err != nil {
				return err
			}
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func sendSTARTTLS(addr, host, user, pass, from, to string, msg []byte) error {
	c, conn, err := smtpClient(addr, host, false)
	if err != nil {
		return err
	}
	return smtpDeliver(c, conn, host, user, pass, from, to, msg, true)
}

func sendImplicitTLS(addr, host, user, pass, from, to string, msg []byte) error {
	c, conn, err := smtpClient(addr, host, true)
	if err != nil {
		return err
	}
	return smtpDeliver(c, conn, host, user, pass, from, to, msg, false)
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
		`Enter your Patriot Pest email and we will send you a sign-in code.`,
		"", emailFormHTML("", w, r))
}

// loginPost validates the email, rate-limits code requests, sends the OTP,
// and renders the code entry form.
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
	email := normalizeEmail(r.FormValue("email"))
	if !m.otpReqLimiter.allow("otp:"+email) {
		w.WriteHeader(http.StatusTooManyRequests)
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "", "Too many attempts. Wait a bit and try again.", emailFormHTML("", w, r))
		return
	}
	if !emailAllowed(email) {
		slog.Warn("report: unauthorized login attempt", "email", email)
		w.WriteHeader(http.StatusForbidden)
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "", "This email is not authorized for these reports.", emailFormHTML("", w, r))
		return
	}
	if !smtpConfigured() {
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "",
			"Email sending is not configured yet. Contact the administrator to finish setup.", emailFormHTML(email, w, r))
		return
	}
	db, err := m.db()
	if err != nil {
		slog.Error("report: db unavailable", "err", err.Error())
		http.Error(w, "Could not start sign-in. Try again later.", http.StatusInternalServerError)
		return
	}
	purgeExpiredReportSessions(db)
	code, err := genCode()
	if err != nil {
		slog.Error("report: code generation failed", "err", err.Error())
		http.Error(w, "Could not start sign-in. Try again later.", http.StatusInternalServerError)
		return
	}
	if err := m.storeOTP(db, email, code); err != nil {
		slog.Error("report: otp store failed", "err", err.Error())
		http.Error(w, "Could not start sign-in. Try again later.", http.StatusInternalServerError)
		return
	}
	if err := sendOTPEmail(email, code); err != nil {
		slog.Error("report: otp email failed", "err", err.Error())
		_, _ = db.Exec(`DELETE FROM report_otps WHERE code_hash=?`, sha256Hex(code))
		writeAuthPage(w, "Sign in | Patriot Status Reports", "PATRIOT STATUS",
			"Weekly reports for the Patriot Pest Control team.", "", "Could not send the code. Try again later.", emailFormHTML(email, w, r))
		return
	}
	slog.Info("report: otp sent", "email", email)
	writeAuthPage(w, "Check your email | Patriot Status Reports", "CHECK YOUR EMAIL",
		"We sent a 6-digit code to "+html.EscapeString(email)+". It expires in 10 minutes.",
		"", "", codeFormHTML(email, w, r))
}

// verifyPost checks the code and, on success, creates the session.
func (m *Module) verifyPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if !view.VerifyCSRF(r) {
		w.WriteHeader(http.StatusForbidden)
		writeAuthPage(w, "Check your email | Patriot Status Reports", "CHECK YOUR EMAIL", "",
			"", "Security token expired. Reload the page and try again.", emailFormHTML("", w, r))
		return
	}
	if !m.verifyLimiter.allow("verify:"+clientIP(r)) {
		w.WriteHeader(http.StatusTooManyRequests)
		writeAuthPage(w, "Check your email | Patriot Status Reports", "CHECK YOUR EMAIL", "",
			"", "Too many attempts. Wait a bit and try again.", emailFormHTML("", w, r))
		return
	}
	email := normalizeEmail(r.FormValue("email"))
	code := strings.TrimSpace(r.FormValue("code"))
	if email == "" || code == "" || !emailAllowed(email) {
		w.WriteHeader(http.StatusForbidden)
		writeAuthPage(w, "Check your email | Patriot Status Reports", "CHECK YOUR EMAIL", "",
			"", "This email is not authorized for these reports.", emailFormHTML("", w, r))
		return
	}
	db, err := m.db()
	if err != nil {
		slog.Error("report: db unavailable", "err", err.Error())
		http.Error(w, "Could not verify the code. Try again later.", http.StatusInternalServerError)
		return
	}
	if !m.checkOTP(db, email, code) {
		slog.Warn("report: bad otp attempt", "email", email)
		w.WriteHeader(http.StatusUnauthorized)
		writeAuthPage(w, "Check your email | Patriot Status Reports", "CHECK YOUR EMAIL",
			"We sent a 6-digit code to "+html.EscapeString(email)+". It expires in 10 minutes.",
			"", "That code did not work. Check the code and try again.", codeFormHTML(email, w, r))
		return
	}
	purgeExpiredReportSessions(db)
	if err := m.createSession(db, w, r, email); err != nil {
		slog.Error("report: session create failed", "err", err.Error())
		http.Error(w, "Could not start a session. Try again later.", http.StatusInternalServerError)
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
		`<button type="submit">SEND CODE</button></form>`
}

func codeFormHTML(email string, w http.ResponseWriter, r *http.Request) string {
	return `<form method="post" action="/verify" novalidate>` +
		string(view.CSRFField(w, r)) +
		`<input type="hidden" name="email" value="` + html.EscapeString(email) + `">` +
		`<label for="code">CODE</label>` +
		`<input id="code" name="code" inputmode="numeric" autocomplete="one-time-code" maxlength="6" required autofocus>` +
		`<button type="submit">VERIFY</button></form>` +
		`<p class="alt"><a href="/login">Wrong email? Start over.</a></p>`
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
