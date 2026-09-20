package view

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCSRFFieldRejectsPlantedCookie covers the regression where the hidden input
// was rendered from any cookie at least 32 bytes long. The old guard was
// `len(tok) < 32`, which is a length check rather than validation, so a value
// planted by a sibling subdomain or an XSS foothold was written into the
// attribute unescaped.
func TestCSRFFieldRejectsPlantedCookie(t *testing.T) {
	valid := strings.Repeat("ab", 32) // 64 characters, matches what we issue
	hostile := `"><script>alert(1)</script>` + strings.Repeat("a", 40)

	cases := []struct {
		name      string
		cookie    string
		wantFresh bool
	}{
		{"a token we issued is reused", valid, false},
		{"a planting attempt is replaced", hostile, true},
		{"a short value is replaced", "abc", true},
		{"garbage of valid length is replaced", strings.Repeat("<", 64), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/contact", nil)
			r.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: tc.cookie})
			w := httptest.NewRecorder()

			out := string(CSRFField(w, r))

			if strings.Contains(out, "<script>") {
				t.Fatalf("planted value reached the page: %s", out)
			}

			set := w.Result().Cookies()
			if tc.wantFresh {
				if len(set) != 1 {
					t.Fatalf("expected exactly one Set-Cookie, got %d", len(set))
				}
				if !csrfTokenRe.MatchString(set[0].Value) {
					t.Fatalf("fresh token is not a shape we issue: %q", set[0].Value)
				}
				if !strings.Contains(out, `value="`+set[0].Value+`"`) {
					t.Fatalf("rendered field does not carry the fresh token: %s", out)
				}
				return
			}

			if len(set) != 0 {
				t.Fatalf("a valid token should be reused, got %d Set-Cookie", len(set))
			}
			if !strings.Contains(out, `value="`+valid+`"`) {
				t.Fatalf("valid token missing from %s", out)
			}
		})
	}
}

// TestCSRFRoundTrip proves the tightened guard did not break the forms: the
// token rendered into the page must still verify against the cookie that
// carried it, and a mismatched token must not.
func TestCSRFRoundTrip(t *testing.T) {
	w := httptest.NewRecorder()
	out := string(CSRFField(w, httptest.NewRequest(http.MethodGet, "/contact", nil)))

	set := w.Result().Cookies()
	if len(set) != 1 {
		t.Fatalf("expected one Set-Cookie, got %d", len(set))
	}
	tok := set[0].Value

	post := func(value string) *http.Request {
		r := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader("_csrf="+value))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: tok})
		return r
	}

	if !VerifyCSRF(post(tok)) {
		t.Fatal("a freshly issued token failed verification")
	}
	if VerifyCSRF(post("deadbeef")) {
		t.Fatal("a mismatched token verified")
	}
	if VerifyCSRF(post("")) {
		t.Fatal("an empty token verified")
	}
	if !strings.Contains(out, tok) {
		t.Fatalf("rendered field and issued cookie disagree: %s", out)
	}
}
