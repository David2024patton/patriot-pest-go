package view

// CSRF helpers for the marketing forms (contact, signup). Fail closed: if
// randomness is unavailable the request 500s instead of issuing a weak token.

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
// Fail-closed: if the RNG fails the request gets a 500, never a guessable token.
func CSRFField(w http.ResponseWriter, r *http.Request) template.HTML {
	ck, err := r.Cookie(CSRFCookieName)
	tok := ""
	if err == nil {
		tok = ck.Value
	}
	if len(tok) < 32 {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return template.HTML("")
		}
		tok = hex.EncodeToString(b)
		http.SetCookie(w, csrfCookie(CSRFCookieName, tok, r))
	}
	return template.HTML(`<input type="hidden" name="_csrf" value="` + tok + `">`)
}

// csrfCookie builds the double-submit cookie: HttpOnly so JS cannot read it,
// Secure on TLS so it never travels in the clear, SameSite=Lax against CSRF.
func csrfCookie(name, value string, r *http.Request) *http.Cookie {
	secure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   secure,
	}
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
