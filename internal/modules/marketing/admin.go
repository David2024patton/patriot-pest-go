// Admin analytics console: password-protected dashboard for the first-party
// website analytics (visits, retention, top pages, click tracking, traffic
// sources). Access is gated on two env vars:
//
//	ADMIN_EMAILS         comma-separated allowlist (e.g. owner emails)
//	ADMIN_PASSWORD_HASH  bcrypt hash of the shared admin password
//
// When either is empty the /admin routes are not registered at all, so the
// console 404s fail-closed instead of rendering a login form nobody can use.
// Configure the vars in Dokploy; this repo never carries the values.
package marketing

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	custommw "github.com/David2024patton/patriot-pest-go/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

// Admin session + tracker-exclusion cookies.
const (
	adminSessionCookie = "ppc_admin_sess" // HttpOnly, Path=/admin
	adminFlagCookie    = "ppc_admin"      // readable by tracker.js, Path=/
	adminSessionTTL    = 24 * time.Hour
)

// adminLoginLimiter: 8 login attempts per 15 minutes per IP, then 429.
var adminLoginLimiter = custommw.NewRateLimiter(8, 15*time.Minute)

type adminCtxKey string

const adminEmailCtxKey adminCtxKey = "admin_email"

// initAdmin reads the admin allowlist + password hash once. Called from
// Register; tests may set the fields directly instead.
func (m *Module) initAdmin() {
	m.adminEmails = nil
	for _, e := range strings.Split(os.Getenv("ADMIN_EMAILS"), ",") {
		if e = strings.ToLower(strings.TrimSpace(e)); e != "" {
			m.adminEmails = append(m.adminEmails, e)
		}
	}
	m.adminHash = []byte(os.Getenv("ADMIN_PASSWORD_HASH"))
}

// adminConfigured reports whether the console can serve anyone.
func (m *Module) adminConfigured() bool {
	return len(m.adminEmails) > 0 && len(m.adminHash) > 0
}

// randomHex returns n random bytes as hex.
func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// sha256Hex hashes a session token for storage (the raw token only ever
// lives in the HttpOnly cookie, never in the DB).
func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// adminDB opens the analytics DB (same file as the catalog) or fails the
// request with 500 + a structured log.
func (m *Module) adminDB(w http.ResponseWriter) *sql.DB {
	db, err := data.AnalyticsDB(m.DBPath)
	if err != nil {
		slog.Error("admin: analytics db unavailable", "err", err.Error())
		http.Error(w, "Analytics storage unavailable.", http.StatusInternalServerError)
		return nil
	}
	return db
}

// adminEmailFromSession validates the session cookie and returns the
// logged-in email, or false when the session is missing/expired.
func (m *Module) adminEmailFromSession(r *http.Request) (string, bool) {
	c, err := r.Cookie(adminSessionCookie)
	if err != nil || c.Value == "" {
		return "", false
	}
	db, err := data.AnalyticsDB(m.DBPath)
	if err != nil {
		return "", false
	}
	var email string
	err = db.QueryRow(`SELECT email FROM admin_sessions WHERE token_hash=? AND expires_at>?`,
		sha256Hex(c.Value), time.Now().Unix()).Scan(&email)
	if err != nil {
		return "", false
	}
	return email, true
}

// requireAdmin redirects unauthenticated visitors to the login page and
// injects the admin email into the request context.
func (m *Module) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := m.adminEmailFromSession(r)
		if !ok {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), adminEmailCtxKey, email)))
	})
}

// purgeExpiredSessions deletes dead admin sessions (opportunistic, on login).
func purgeExpiredSessions(db *sql.DB) {
	if _, err := db.Exec(`DELETE FROM admin_sessions WHERE expires_at<=?`, time.Now().Unix()); err != nil {
		slog.Warn("admin: session purge failed", "err", err.Error())
	}
}

// adminLoginGet renders the login form. Already authed -> straight to /admin.
func (m *Module) adminLoginGet(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.adminEmailFromSession(r); ok {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}
	m.renderLogin(w, r, "")
}

// adminLoginPost checks CSRF, the email allowlist and the bcrypt password.
func (m *Module) adminLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		m.renderLogin(w, r, "Could not read the form. Try again.")
		return
	}
	var expected string
	if c, err := r.Cookie(csrfCookieName); err == nil && c != nil {
		expected = c.Value
	}
	if !csrfOK(r.FormValue("_csrf"), expected) {
		w.WriteHeader(http.StatusForbidden)
		m.renderLogin(w, r, "Security token expired. Reload the page and try again.")
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")

	db := m.adminDB(w)
	if db == nil {
		return
	}
	purgeExpiredSessions(db)

	allowed := false
	for _, a := range m.adminEmails {
		if subtle.ConstantTimeCompare([]byte(email), []byte(a)) == 1 {
			allowed = true
			break
		}
	}
	okPass := bcrypt.CompareHashAndPassword(m.adminHash, []byte(password)) == nil
	if !allowed || !okPass {
		slog.Warn("admin: failed login", "email", email)
		w.WriteHeader(http.StatusUnauthorized)
		m.renderLogin(w, r, "Invalid email or password.")
		return
	}

	token, err := randomHex(32)
	if err != nil {
		slog.Error("admin: session token generation failed", "err", err.Error())
		http.Error(w, "Could not start a session.", http.StatusInternalServerError)
		return
	}
	now := time.Now()
	if _, err := db.Exec(`INSERT INTO admin_sessions (token_hash, email, created_at, expires_at)
		VALUES (?, ?, ?, ?)`, sha256Hex(token), email, now.Unix(), now.Add(adminSessionTTL).Unix()); err != nil {
		slog.Error("admin: session insert failed", "err", err.Error())
		http.Error(w, "Could not start a session.", http.StatusInternalServerError)
		return
	}
	secure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	http.SetCookie(w, &http.Cookie{
		Name: adminSessionCookie, Value: token, Path: "/admin",
		MaxAge: int(adminSessionTTL.Seconds()), HttpOnly: true,
		Secure: secure, SameSite: http.SameSiteLaxMode,
	})
	// Tracker exclusion flag: tracker.js skips beaconing when this is set,
	// so owner visits never pollute the dashboard.
	http.SetCookie(w, &http.Cookie{
		Name: adminFlagCookie, Value: "1", Path: "/",
		MaxAge: int(adminSessionTTL.Seconds()), HttpOnly: false,
		Secure: secure, SameSite: http.SameSiteLaxMode,
	})
	slog.Info("admin: login", "email", email)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// adminLogout destroys the session and clears both cookies.
func (m *Module) adminLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(adminSessionCookie); err == nil && c.Value != "" {
		if db, err := data.AnalyticsDB(m.DBPath); err == nil {
			_, _ = db.Exec(`DELETE FROM admin_sessions WHERE token_hash=?`, sha256Hex(c.Value))
		}
	}
	clearCookie := func(name, path string, httpOnly bool) {
		http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: path, MaxAge: -1, HttpOnly: httpOnly})
	}
	clearCookie(adminSessionCookie, "/admin", true)
	clearCookie(adminFlagCookie, "/", false)
	http.Redirect(w, r, "/", http.StatusFound)
}

// ---- login page ----

const adminCSS = `:root{--ink:#141a10;--olive-950:#0e130a;--olive-900:#1c2415;--olive-800:#26301c;--olive-700:#334024;--olive-500:#5c6f3a;--olive-300:#8fa05e;--khaki:#c8b98c;--paper:#ece4cd;--orange:#f4772e;--orange-hot:#ff8c3b;--red:#c8402a;--cream:#f5f1e4}
*{box-sizing:border-box}
body{margin:0;background:var(--olive-950);color:var(--cream);font-family:system-ui,-apple-system,"Segoe UI",Roboto,sans-serif}
.wrap{max-width:1080px;margin:0 auto;padding:1.2rem}
header.top{display:flex;align-items:center;justify-content:space-between;gap:1rem;flex-wrap:wrap;border-bottom:2px solid var(--olive-700);padding-bottom:1rem;margin-bottom:1.4rem}
.brand{font-weight:800;letter-spacing:.14em;font-size:1.05rem}
.brand .star{color:var(--orange)}
.sub{color:var(--olive-300);font-size:.8rem;letter-spacing:.08em}
nav.links{display:flex;gap:.5rem;align-items:center;flex-wrap:wrap}
nav.links a{color:var(--khaki);text-decoration:none;border:1px solid var(--olive-700);border-radius:6px;padding:.35rem .7rem;font-size:.8rem}
nav.links a.on,nav.links a:hover{border-color:var(--orange);color:var(--cream)}
nav.links .who{color:var(--olive-300);font-size:.78rem;margin-right:.4rem}
.cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:.8rem;margin-bottom:1.4rem}
.card{background:var(--olive-900);border:1px solid var(--olive-700);border-left:4px solid var(--orange);border-radius:8px;padding:.9rem 1rem}
.card .k{font-size:.68rem;letter-spacing:.18em;color:var(--olive-300)}
.card .v{font-size:1.7rem;font-weight:800;margin-top:.2rem}
.card .s{font-size:.75rem;color:var(--khaki);margin-top:.15rem}
.panel{background:var(--olive-900);border:1px solid var(--olive-700);border-radius:8px;padding:1rem;margin-bottom:1.4rem}
.panel h2{margin:0 0 .8rem;font-size:.8rem;letter-spacing:.2em;color:var(--khaki)}
table{width:100%;border-collapse:collapse;font-size:.85rem}
th{text-align:left;color:var(--olive-300);font-weight:600;font-size:.68rem;letter-spacing:.14em;padding:.4rem .5rem;border-bottom:1px solid var(--olive-700)}
td{padding:.45rem .5rem;border-bottom:1px solid var(--olive-800);vertical-align:top}
tr:last-child td{border-bottom:0}
td.num,th.num{text-align:right;font-variant-numeric:tabular-nums;white-space:nowrap}
td.path{font-family:ui-monospace,Menlo,monospace;font-size:.78rem;word-break:break-all}
.bar{height:8px;background:var(--olive-800);border-radius:4px;overflow:hidden;min-width:60px}
.bar i{display:block;height:100%;background:var(--orange)}
.muted{color:var(--olive-300);font-size:.8rem}
.foot{margin:2rem 0 3rem;color:var(--olive-300);font-size:.75rem;text-align:center}
form.login{max-width:380px;margin:12vh auto;background:var(--olive-900);border:1px solid var(--olive-700);border-radius:10px;padding:1.6rem}
form.login h1{margin:0 0 .3rem;font-size:1rem;letter-spacing:.2em}
form.login p{color:var(--olive-300);font-size:.8rem;margin:0 0 1.2rem}
form.login label{display:block;font-size:.7rem;letter-spacing:.14em;color:var(--khaki);margin:.9rem 0 .3rem}
form.login input{width:100%;padding:.6rem .7rem;background:var(--olive-950);border:1px solid var(--olive-700);border-radius:6px;color:var(--cream);font-size:1rem}
form.login input:focus{outline:2px solid var(--orange);border-color:var(--orange)}
form.login button{width:100%;margin-top:1.2rem;padding:.7rem;background:var(--orange);border:0;border-radius:6px;color:#141a10;font-weight:800;letter-spacing:.1em;cursor:pointer;font-size:.95rem}
form.login button:hover{background:var(--orange-hot)}
.err{background:#3a1512;border:1px solid var(--red);color:#f2b8a8;border-radius:6px;padding:.6rem .8rem;font-size:.82rem;margin-bottom:1rem}
svg.chart{width:100%;height:auto;display:block}
svg.chart rect{fill:var(--orange)}
svg.chart rect:hover{fill:var(--orange-hot)}
.chart-x{display:flex;justify-content:space-between;color:var(--olive-300);font-size:.68rem;margin-top:.3rem}`

// renderLogin writes the standalone login page (no site chrome: this URL
// should not look like public marketing).
func (m *Module) renderLogin(w http.ResponseWriter, r *http.Request, errMsg string) {
	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8">` +
		`<meta name="viewport" content="width=device-width,initial-scale=1">` +
		`<meta name="robots" content="noindex,nofollow">` +
		`<title>Admin Login | Patriot Pest Control</title><style>` + adminCSS + `</style></head><body>`)
	sb.WriteString(`<form class="login" method="post" action="/admin/login">`)
	sb.WriteString(`<h1><span style="color:var(--orange)">&#9733;</span> PATRIOT ANALYTICS</h1>`)
	sb.WriteString(`<p>Restricted. Authorized owners only.</p>`)
	if errMsg != "" {
		sb.WriteString(`<div class="err">` + html.EscapeString(errMsg) + `</div>`)
	}
	sb.WriteString(string(m.csrfField(r, w)))
	sb.WriteString(`<label for="email">EMAIL</label><input id="email" name="email" type="email" autocomplete="username" required autofocus>`)
	sb.WriteString(`<label for="pw">PASSWORD</label><input id="pw" name="password" type="password" autocomplete="current-password" required>`)
	sb.WriteString(`<button type="submit">SIGN IN</button></form></body></html>`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	_, _ = w.Write([]byte(sb.String()))
}

// ---- dashboard ----

// adminDashboard renders the traffic dashboard for the selected day range.
func (m *Module) adminDashboard(w http.ResponseWriter, r *http.Request) {
	days := 30
	if d, err := strconv.Atoi(r.URL.Query().Get("days")); err == nil {
		days = d
	}
	db := m.adminDB(w)
	if db == nil {
		return
	}
	st, err := data.QueryAnalytics(db, days, time.Now().Unix())
	if err != nil {
		slog.Error("admin: dashboard query failed", "err", err.Error())
		http.Error(w, "Could not load analytics.", http.StatusInternalServerError)
		return
	}
	email, _ := r.Context().Value(adminEmailCtxKey).(string)
	m.renderDashboard(w, st, email)
}

func pct(part, total int64) string {
	if total <= 0 {
		return "0%"
	}
	return strconv.FormatInt(part*100/total, 10) + "%"
}

func (m *Module) renderDashboard(w http.ResponseWriter, st *data.AnalyticsStats, email string) {
	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8">` +
		`<meta name="viewport" content="width=device-width,initial-scale=1">` +
		`<meta name="robots" content="noindex,nofollow">` +
		`<title>Analytics | Patriot Pest Control</title><style>` + adminCSS + `</style></head><body><div class="wrap">`)

	// Header.
	sb.WriteString(`<header class="top"><div><div class="brand"><span class="star">&#9733;</span> PATRIOT ANALYTICS</div>` +
		`<div class="sub">FIRST-PARTY TRAFFIC &middot; LAST ` + strconv.Itoa(st.Days) + ` DAYS</div></div><nav class="links">`)
	for _, d := range []int{7, 30, 90} {
		cls := ""
		if d == st.Days {
			cls = ` class="on"`
		}
		sb.WriteString(`<a` + cls + ` href="/admin?days=` + strconv.Itoa(d) + `">` + strconv.Itoa(d) + `d</a>`)
	}
	sb.WriteString(`<span class="who">` + html.EscapeString(email) + `</span><a href="/admin/logout">LOGOUT</a></nav></header>`)

	// Stat cards.
	card := func(k, v, s string) {
		sb.WriteString(`<div class="card"><div class="k">` + k + `</div><div class="v">` + v + `</div><div class="s">` + s + `</div></div>`)
	}
	sb.WriteString(`<div class="cards">`)
	card("PAGEVIEWS", strconv.FormatInt(st.Pageviews, 10), "pages loaded")
	card("UNIQUE VISITORS", strconv.FormatInt(st.Visitors, 10), "distinct browsers")
	card("SESSIONS", strconv.FormatInt(st.Sessions, 10), "visits (30 min timeout)")
	card("NEW VISITORS", strconv.FormatInt(st.NewVisitors, 10), "first seen in range")
	card("RETURNING", strconv.FormatInt(st.ReturningVisitors, 10), "seen before range")
	sb.WriteString(`</div>`)

	// Daily bar chart (inline SVG, no JS).
	sb.WriteString(`<div class="panel"><h2>DAILY PAGEVIEWS</h2>`)
	if len(st.Daily) == 0 || st.Pageviews == 0 {
		sb.WriteString(`<p class="muted">No traffic recorded in this range yet.</p>`)
	} else {
		var maxC int64 = 1
		for _, d := range st.Daily {
			if d.Count > maxC {
				maxC = d.Count
			}
		}
		const W, H = 900, 150
		n := float64(len(st.Daily))
		sb.WriteString(`<svg class="chart" viewBox="0 0 ` + strconv.Itoa(W) + ` ` + strconv.Itoa(H) + `" role="img" aria-label="Daily pageviews bar chart">`)
		for i, d := range st.Daily {
			bw := float64(W) / n
			h := float64(d.Count) / float64(maxC) * float64(H-24)
			if h < 1 && d.Count > 0 {
				h = 1
			}
			x := float64(i)*bw + 1
			y := float64(H) - h
			ww := bw - 2
			if ww < 1 {
				ww = 1
			}
			sb.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f"><title>%s: %d</title></rect>`,
				x, y, ww, h, html.EscapeString(d.Date), d.Count))
		}
		sb.WriteString(`</svg><div class="chart-x"><span>` + html.EscapeString(st.Daily[0].Date) +
			`</span><span>` + html.EscapeString(st.Daily[len(st.Daily)-1].Date) + `</span></div>`)
	}
	sb.WriteString(`</div>`)

	// Top pages.
	sb.WriteString(`<div class="panel"><h2>TOP PAGES</h2>`)
	if len(st.TopPages) == 0 {
		sb.WriteString(`<p class="muted">Nothing yet.</p>`)
	} else {
		sb.WriteString(`<table><tr><th>PAGE</th><th class="num">VIEWS</th><th class="num">SHARE</th><th style="width:30%"></th></tr>`)
		for _, p := range st.TopPages {
			wPct := pct(p.Count, st.Pageviews)
			sb.WriteString(`<tr><td class="path">` + html.EscapeString(p.Path) + `</td><td class="num">` +
				strconv.FormatInt(p.Count, 10) + `</td><td class="num">` + wPct +
				`</td><td><div class="bar"><i style="width:` + wPct + `"></i></div></td></tr>`)
		}
		sb.WriteString(`</table>`)
	}
	sb.WriteString(`</div>`)

	// Traffic sources.
	sb.WriteString(`<div class="panel"><h2>TRAFFIC SOURCES</h2>`)
	if len(st.TopSources) == 0 {
		sb.WriteString(`<p class="muted">Nothing yet.</p>`)
	} else {
		sb.WriteString(`<table><tr><th>SOURCE</th><th class="num">VIEWS</th><th class="num">SHARE</th><th style="width:30%"></th></tr>`)
		for _, s := range st.TopSources {
			wPct := pct(s.Count, st.Pageviews)
			sb.WriteString(`<tr><td>` + html.EscapeString(s.Source) + `</td><td class="num">` +
				strconv.FormatInt(s.Count, 10) + `</td><td class="num">` + wPct +
				`</td><td><div class="bar"><i style="width:` + wPct + `"></i></div></td></tr>`)
		}
		sb.WriteString(`</table><p class="muted">Facebook / Google / Direct and friends. UTM parameters on ad links override referrer detection.</p>`)
	}
	sb.WriteString(`</div>`)

	// Top clicks.
	sb.WriteString(`<div class="panel"><h2>TOP CLICKS</h2>`)
	if len(st.TopClicks) == 0 {
		sb.WriteString(`<p class="muted">No clicks tracked yet.</p>`)
	} else {
		sb.WriteString(`<table><tr><th>CLICKED</th><th>PAGE</th><th class="num">CLICKS</th></tr>`)
		for _, c := range st.TopClicks {
			label := c.Label
			if label == "" {
				label = "(no label)"
			}
			sb.WriteString(`<tr><td>` + html.EscapeString(label) + `</td><td class="path">` +
				html.EscapeString(c.Path) + `</td><td class="num">` + strconv.FormatInt(c.Count, 10) + `</td></tr>`)
		}
		sb.WriteString(`</table>`)
	}
	sb.WriteString(`</div>`)

	// Recent activity.
	sb.WriteString(`<div class="panel"><h2>RECENT ACTIVITY</h2>`)
	if len(st.Recent) == 0 {
		sb.WriteString(`<p class="muted">Nothing yet.</p>`)
	} else {
		sb.WriteString(`<table><tr><th>TIME</th><th>WHAT</th><th>PAGE</th><th>SOURCE</th><th>DETAIL</th></tr>`)
		for _, h := range st.Recent {
			detail := h.Label
			if h.Kind == "pageview" {
				detail = ""
			}
			sb.WriteString(`<tr><td class="num">` + time.Unix(h.TS, 0).Format("Jan 02 15:04") + `</td><td>` +
				html.EscapeString(h.Kind) + `</td><td class="path">` + html.EscapeString(h.Path) + `</td><td>` +
				html.EscapeString(h.Source) + `</td><td>` + html.EscapeString(detail) + `</td></tr>`)
		}
		sb.WriteString(`</table>`)
	}
	sb.WriteString(`</div>`)

	sb.WriteString(`<div class="foot">TRACKING ACTIVE &middot; YOUR OWN VISITS ARE EXCLUDED (ADMIN COOKIE) &middot; BOTS FILTERED &middot; DNT RESPECTED</div>`)
	sb.WriteString(`</div></body></html>`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	_, _ = w.Write([]byte(sb.String()))
}
