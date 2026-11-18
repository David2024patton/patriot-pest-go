package view

// CSRF helpers for the marketing forms (contact, signup). Fail closed: if
// randomness is unavailable the request 500s instead of issuing a weak token.

import (
	"github.com/David2024patton/patriot-pest-go/internal/data"

	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"html/template"
	"net/http"
	"regexp"
)

// CSRFCookieName is the double-submit cookie for the marketing forms. It is
// also read by internal/modules/marketing, which issues the same 64-character
// hex token, so both writers and this reader agree on the shape.
const CSRFCookieName = "_csrf"

// csrfTokenRe is the shape this package and the marketing module issue: 32
// random bytes, hex encoded. A cookie that does not match is replaced rather
// than rendered, so a value planted by a sibling subdomain or an XSS foothold
// cannot reach the attribute below.
var csrfTokenRe = regexp.MustCompile(`^[0-9a-f]{64}$`)

// CSRFField renders the hidden input, issuing a token when the cookie is
// missing or is not one we could have issued. Fail-closed: if the RNG fails the
// request gets a 500, never a guessable token.
func CSRFField(w http.ResponseWriter, r *http.Request) template.HTML {
	tok := ""
	if ck, err := r.Cookie(CSRFCookieName); err == nil && csrfTokenRe.MatchString(ck.Value) && data.CSRFTokenRegistered(ck.Value) {
		tok = ck.Value
	} else {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return template.HTML("")
		}
		tok = hex.EncodeToString(b)
		data.RegisterCSRFToken(tok)
		http.SetCookie(w, csrfCookie(CSRFCookieName, tok, r))
	}
	// Escaped as well as validated: the token is hex, but the attribute is
	// written as template.HTML, which turns escaping off entirely.
	return template.HTML(`<input type="hidden" name="_csrf" value="` + template.HTMLEscapeString(tok) + `">`)
}

// csrfCookie builds the double-submit cookie: HttpOnly so JS cannot read it,
// always Secure so it never travels in the clear (the site is HTTPS-only
// behind Traefik; X-Forwarded-Proto is client-supplied and never decides
// cookie security), SameSite=Lax against CSRF.
func csrfCookie(name, value string, r *http.Request) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
}

// VerifyCSRF checks the posted _csrf field against the cookie (constant
// time) and against the server registry: the token must be one this
// server issued, so a planted cookie and field pair cannot pass.
func VerifyCSRF(r *http.Request) bool {
	expected := ""
	if ck, err := r.Cookie(CSRFCookieName); err == nil {
		expected = ck.Value
	}
	got := r.FormValue("_csrf")
	if got == "" || subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
		return false
	}
	return data.CSRFTokenRegistered(got)
}
