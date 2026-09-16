package view

// CSRF helpers for authenticated app pages (dashboards + admin surfaces).
// The cookie name matches internal/auth's login flow ("_csrf") so one token
// covers every surface in a single browser session.

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"html/template"
	"net/http"
)

// CSRFCookieName is shared with internal/auth's web flow.
const CSRFCookieName = "_csrf"

// CSRFField renders the hidden input and sets the cookie when absent or short.
func CSRFField(w http.ResponseWriter, r *http.Request) template.HTML {
	ck, err := r.Cookie(CSRFCookieName)
	tok := ""
	if err == nil {
		tok = ck.Value
	}
	if len(tok) < 32 {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			tok = "dev-csrf-token"
		} else {
			tok = hex.EncodeToString(b)
		}
		http.SetCookie(w, &http.Cookie{Name: CSRFCookieName, Value: tok, Path: "/", SameSite: http.SameSiteLaxMode})
	}
	return template.HTML(`<input type="hidden" name="_csrf" value="` + tok + `">`)
}

// VerifyCSRF checks the posted _csrf field against the cookie (constant time).
func VerifyCSRF(r *http.Request) bool {
	expected := ""
	if ck, err := r.Cookie(CSRFCookieName); err == nil {
		expected = ck.Value
	}
	got := r.FormValue("_csrf")
	return got != "" && subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1
}
