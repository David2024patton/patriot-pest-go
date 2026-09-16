package auth

// WebFlow — the browser-facing passwordless login surface (port of app/Controllers/AuthController.php):
//   GET  /login          identifier form (email / phone / account number)
//   POST /login          rate-limit, resolve identity, issue OTP by email
//   GET  /login/verify   code-entry form
//   POST /login/verify   verify code, start staff/customer session
//   GET  /logout         destroy session + pending state
// Sessions are the same "session" cookie the API handlers mint (HttpOnly, SameSite=Lax).

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/David2024patton/patriot-pest-go/internal/view"
	"github.com/go-chi/chi/v5"
)

// csrfCookieName is shared with the marketing module's contact form so one
// token covers both surfaces (PHP used a single session CSRF store).
const csrfCookieName = "_csrf"

// WebConfig carries the OTP/session knobs from config so handlers stay pure.
type WebConfig struct {
	Dev             bool // APP_ENV=local: drop Secure on cookies for http://localhost
	OTPTTL          int  // seconds (OTP_TTL, default 600)
	MaxAttempts     int  // OTP_MAX_ATTEMPTS (5)
	StaffSessionTTL int  // SESSION_LIFETIME_STAFF (7200)
	CustomerTTL     int  // SESSION_LIFETIME_CUSTOMER (900)
}

// WebFlow bundles config + mailer behind the login routes.
type WebFlow struct {
	cfg    WebConfig
	mailer *Mailer
}

// NewWebFlow builds the web flow. Mailer dev mode is decided by NewMailer.
func NewWebFlow(cfg WebConfig, mailer *Mailer) *WebFlow {
	return &WebFlow{cfg: cfg, mailer: mailer}
}

// RegisterWebRoutes mounts the login surface on the main router.
func RegisterWebRoutes(r chi.Router, f *WebFlow) {
	r.Get("/login", f.HandleLoginForm)
	r.Post("/login", f.HandleLoginRequest)
	r.Get("/login/verify", f.HandleVerifyForm)
	r.Post("/login/verify", f.HandleVerify)
	r.Get("/logout", f.HandleLogout)
}

// ---- CSRF (cookie-bound, SameSite=Lax) ----

func newCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "dev-csrf-token"
	}
	return hex.EncodeToString(b)
}

// csrfField renders the hidden input and sets the cookie when absent.
func (f *WebFlow) csrfField(w http.ResponseWriter, r *http.Request) template.HTML {
	ck, err := r.Cookie(csrfCookieName)
	tok := ""
	if err == nil {
		tok = ck.Value
	}
	if len(tok) < 32 {
		tok = newCSRFToken()
		http.SetCookie(w, &http.Cookie{Name: csrfCookieName, Value: tok, Path: "/", SameSite: http.SameSiteLaxMode})
	}
	return template.HTML(`<input type="hidden" name="_csrf" value="` + tok + `">`)
}

// verifyCSRF checks the posted field against the cookie (constant time).
func verifyCSRF(r *http.Request) bool {
	expected := ""
	if ck, err := r.Cookie(csrfCookieName); err == nil {
		expected = ck.Value
	}
	got := r.FormValue("_csrf")
	return got != "" && subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1
}

// ---- cookie state: pending login + flash (PRG pattern) ----

type pendingState struct {
	Email string `json:"email"` // canonical email the code was issued to
	Type  string `json:"type"`  // 'staff' | 'customer'
}

type flashMsg struct {
	Kind string `json:"k"` // 'error' | 'sent'
	Msg  string `json:"m,omitempty"`
	To   string `json:"to,omitempty"`
}

const (
	pendingCookieName = "ppc_pending"
	flashCookieName   = "ppc_flash"
)

// Cookie values carry base64(JSON): Go 1.26's net/http strips raw '"' from
// cookie values, which would corrupt JSON stored directly.
func setPending(w http.ResponseWriter, email, typ string) {
	b, _ := json.Marshal(pendingState{Email: email, Type: typ})
	http.SetCookie(w, &http.Cookie{Name: pendingCookieName, Value: base64.RawStdEncoding.EncodeToString(b), Path: "/", SameSite: http.SameSiteLaxMode})
}

func readPending(r *http.Request) *pendingState {
	ck, err := r.Cookie(pendingCookieName)
	if err != nil || ck.Value == "" {
		return nil
	}
	b, _ := base64.RawStdEncoding.DecodeString(ck.Value)
	var p pendingState
	if err := json.Unmarshal(b, &p); err != nil {
		return nil
	}
	return &p
}

func setFlash(w http.ResponseWriter, fl flashMsg) {
	b, _ := json.Marshal(fl)
	http.SetCookie(w, &http.Cookie{Name: flashCookieName, Value: base64.RawStdEncoding.EncodeToString(b), Path: "/", SameSite: http.SameSiteLaxMode})
}

func readFlash(r *http.Request) flashMsg {
	var fl flashMsg
	if ck, err := r.Cookie(flashCookieName); err == nil && ck.Value != "" {
		b, _ := base64.RawStdEncoding.DecodeString(ck.Value)
		_ = json.Unmarshal(b, &fl)
	}
	return fl
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Path: "/", MaxAge: -1})
}

func redirectTo(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
}

// ---- helpers ----

// minutes renders seconds as a whole-minute count for flash copy (min 1).
func minutes(sec int) int {
	m := (sec + 59) / 60
	if m < 1 {
		return 1
	}
	return m
}

// maskEmail mirrors AuthController::maskEmail: first char + bullets + domain.
func maskEmail(s string) string {
	i := strings.Index(s, "@")
	if i < 0 {
		return s
	}
	local, dom := s[:i], s[i:]
	bullets := len(local) - 1
	if bullets < 3 {
		bullets = 3
	}
	return local[:1] + strings.Repeat("•", bullets) + dom
}

// dashboardFor mirrors AuthController::dashboardFor.
func dashboardFor(role, roleLabel string) string {
	switch role {
	case "customer":
		return "/customer-dashboard"
	case "SuperAdmin":
		return "/admin"
	default:
		if roleLabel == "admin" || roleLabel == "super-user" {
			return "/admin"
		}
		return "/staff-dashboard"
	}
}

func sessionFromCookie(r *http.Request) *Session {
	ck, err := r.Cookie("session")
	if err != nil || ck.Value == "" {
		return nil
	}
	s, ok := GetSession(ck.Value)
	if !ok {
		return nil
	}
	return s
}

func (f *WebFlow) setSessionCookie(w http.ResponseWriter, id string, ttlSec int) {
	http.SetCookie(w, &http.Cookie{
		Name: "session", Value: id, Path: "/", HttpOnly: true,
		Secure: !f.cfg.Dev, SameSite: http.SameSiteLaxMode,
		Expires: time.Now().Add(time.Duration(ttlSec) * time.Second),
	})
}

// ---- handlers ----

// HandleLoginForm — GET /login. Authenticated visitors go straight to their dashboard.
func (f *WebFlow) HandleLoginForm(w http.ResponseWriter, r *http.Request) {
	if sess := sessionFromCookie(r); sess != nil {
		redirectTo(w, r, dashboardFor(sess.Role, sess.RoleLabel))
		return
	}
	fl := readFlash(r)
	clearCookie(w, flashCookieName)
	view.Page(w, r, "auth-login", "Sign In | Patriot Pest Control",
		"One secure sign-in for customers and staff. No password — we email you a code and send you to the right dashboard.", "",
		map[string]any{
			"AppUI":      true, // login surfaces load admin.css (parity with PHP $__appUi)
			"Csrf":       f.csrfField(w, r),
			"FlashError": flashErr(fl),
		})
}

// HandleVerifyForm — GET /login/verify. Renders the code-entry step.
func (f *WebFlow) HandleVerifyForm(w http.ResponseWriter, r *http.Request) {
	if sess := sessionFromCookie(r); sess != nil {
		redirectTo(w, r, dashboardFor(sess.Role, sess.RoleLabel))
		return
	}
	fl := readFlash(r)
	clearCookie(w, flashCookieName)
	var sent string
	if fl.Kind == "sent" && fl.To != "" {
		sent = fl.To
	} else if p := readPending(r); p != nil && p.Email != "" {
		sent = maskEmail(p.Email)
	}
	view.Page(w, r, "auth-verify", "Enter Your Code | Patriot Pest Control",
		"Type the 6-digit code we emailed you. No password to remember.", "",
		map[string]any{
			"AppUI":      true,
			"Csrf":       f.csrfField(w, r),
			"FlashError": flashErr(fl),
			"SentTo":     sent,
		})
}

// HandleLoginRequest — POST /login. Rate-limits per IP, resolves the identifier
// (staff first, then customer by email/phone/account number), issues the OTP.
// Unknown identifiers still get the same "sent" response (enumeration defense).
func (f *WebFlow) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	if !verifyCSRF(r) {
		setFlash(w, flashMsg{Kind: "error", Msg: "Security check failed. Please try again."})
		redirectTo(w, r, "/login")
		return
	}
	identifier := strings.TrimSpace(r.FormValue("identifier"))
	if identifier == "" {
		setFlash(w, flashMsg{Kind: "error", Msg: "Please enter your email, phone number, or account number."})
		redirectTo(w, r, "/login")
		return
	}

	ip := clientIp(r)
	loginKey := "login:" + ip
	if rl.tooMany(loginKey, 5, 60) {
		wait := rl.retryAfter(loginKey, 5, 60)
		setFlash(w, flashMsg{Kind: "error", Msg: fmt.Sprintf("Too many login attempts. Please wait %d minute(s) before trying again.", minutes(wait))})
		redirectTo(w, r, "/login")
		return
	}
	rl.hit(loginKey)

	if s := FindStaffForLogin(identifier); s != nil {
		rl.clear(loginKey)
		setPending(w, s.Email, "staff")
		if f.issueAndEmail(w, r, s.Email) {
			setFlash(w, flashMsg{Kind: "sent", To: maskEmail(s.Email)})
			redirectTo(w, r, "/login/verify")
		}
		return
	}
	if p := resolve(identifier); p != nil && p.Email != "" {
		rl.clear(loginKey)
		setPending(w, p.Email, "customer")
		if f.issueAndEmail(w, r, p.Email) {
			setFlash(w, flashMsg{Kind: "sent", To: maskEmail(p.Email)})
			redirectTo(w, r, "/login/verify")
		}
		return
	}
	// Enumeration defense: look identical whether or not we know you; no code issued.
	setPending(w, identifier, "customer")
	setFlash(w, flashMsg{Kind: "sent", To: maskEmail(identifier)})
	redirectTo(w, r, "/login/verify")
}

// HandleVerify — POST /login/verify. Checks the code and starts a session.
func (f *WebFlow) HandleVerify(w http.ResponseWriter, r *http.Request) {
	if !verifyCSRF(r) {
		setFlash(w, flashMsg{Kind: "error", Msg: "Security check failed. Please try again."})
		redirectTo(w, r, "/login/verify")
		return
	}
	p := readPending(r)
	code := strings.TrimSpace(r.FormValue("code"))
	if p == nil || p.Email == "" || p.Type == "" || code == "" {
		clearCookie(w, pendingCookieName)
		setFlash(w, flashMsg{Kind: "error", Msg: "Please start over and request a new code."})
		redirectTo(w, r, "/login")
		return
	}
	key := "otp:login:" + p.Email
	if rl.tooMany(key, f.cfg.MaxAttempts, f.cfg.OTPTTL) {
		setFlash(w, flashMsg{Kind: "error", Msg: fmt.Sprintf("Too many attempts. Try again in about %d minute(s).", minutes(rl.retryAfter(key, f.cfg.MaxAttempts, f.cfg.OTPTTL)))})
		redirectTo(w, r, "/login/verify")
		return
	}
	ok, err := Verify(r.Context(), p.Email, "login", code, f.cfg.MaxAttempts)
	if !ok {
		rl.hit(key)
		var msg string
		if err != nil && strings.Contains(err.Error(), "too many") {
			msg = fmt.Sprintf("Too many attempts. Try again in about %d minute(s).", minutes(rl.retryAfter(key, f.cfg.MaxAttempts, f.cfg.OTPTTL)))
		} else {
			msg = "That code is incorrect or has expired. Please try again."
		}
		setFlash(w, flashMsg{Kind: "error", Msg: msg})
		redirectTo(w, r, "/login/verify")
		return
	}
	rl.clear(key)
	var dest string
	if p.Type == "staff" {
		dest = f.startStaff(w, r, p.Email)
	} else {
		dest = f.startCustomer(w, r, p.Email)
	}
	clearCookie(w, pendingCookieName)
	redirectTo(w, r, dest)
}

// HandleLogout — GET /logout. Destroys session + login state, then home.
func (f *WebFlow) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if sess := sessionFromCookie(r); sess != nil {
		slog.Info("web logout", "identity", sess.Identity)
	}
	clearCookie(w, "session")
	clearCookie(w, pendingCookieName)
	clearCookie(w, flashCookieName)
	redirectTo(w, r, "/")
}

// ---- issue + session internals ----

// issueAndEmail issues the login OTP under the canonical email and sends it.
// Returns false (having redirected to /login) when something failed.
func (f *WebFlow) issueAndEmail(w http.ResponseWriter, r *http.Request, to string) bool {
	key := "otp_issue:" + to
	if rl.tooMany(key, 3, 300) {
		wait := rl.retryAfter(key, 3, 300)
		setFlash(w, flashMsg{Kind: "error", Msg: fmt.Sprintf("Too many code requests. Please wait %d minute(s) before requesting another.", minutes(wait))})
		redirectTo(w, r, "/login")
		return false
	}
	rl.hit(key)
	code, err := Issue(r.Context(), to, "login", 6, f.cfg.OTPTTL)
	if err != nil {
		setFlash(w, flashMsg{Kind: "error", Msg: "Could not issue your code. Please try again."})
		redirectTo(w, r, "/login")
		return false
	}
	expiresIn := f.cfg.OTPTTL / 60
	body := MailTemplate("Your Patriot Pest Control sign-in code",
		`<p style="font-size:32px;letter-spacing:8px;font-weight:bold;color:#c8a24a">`+code+`</p>`+
			fmt.Sprintf(`<p>This code expires in %d minutes and works once. If you didn't request it, you can safely ignore this email.</p>`, expiresIn))
	if !f.mailer.Send(to, "Your Patriot Pest Control sign-in code", body) {
		setFlash(w, flashMsg{Kind: "error", Msg: "Could not send your code. Please try again in a moment."})
		redirectTo(w, r, "/login")
		return false
	}
	return true
}

// startStaff starts a staff session and returns the dashboard destination
// ("" when it redirected back to /login).
func (f *WebFlow) startStaff(w http.ResponseWriter, r *http.Request, email string) string {
	s := StaffByEmail(email)
	if s == nil {
		setFlash(w, flashMsg{Kind: "error", Msg: "This account is no longer active."})
		redirectTo(w, r, "/login")
		return ""
	}
	if s.Role == "super-user" {
		slog.Info("staff login via /login requires the elevated /su surface", "email", email)
		setFlash(w, flashMsg{Kind: "error", Msg: "Super-user accounts must use the dedicated /su login surface."})
		redirectTo(w, r, "/login")
		return ""
	}
	sess := CreateSession(email, "staff", f.cfg.StaffSessionTTL)
	sess.SetRoleLabel(s.Role)
	f.setSessionCookie(w, sess.ID, f.cfg.StaffSessionTTL)
	return dashboardFor("staff", s.Role)
}

// startCustomer starts a customer session and returns the destination.
func (f *WebFlow) startCustomer(w http.ResponseWriter, r *http.Request, email string) string {
	sess := CreateSession(email, "customer", f.cfg.CustomerTTL)
	f.setSessionCookie(w, sess.ID, f.cfg.CustomerTTL)
	return "/customer-dashboard"
}

// flashErr extracts the error message from a flash cookie (sent flashes show nothing).
func flashErr(fl flashMsg) string {
	if fl.Kind == "error" {
		return fl.Msg
	}
	return ""
}
