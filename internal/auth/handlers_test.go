package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func jsonReader(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func setTestLookup(t *testing.T, l Lookup) {
	t.Helper()
	SetLookup(l)
	t.Cleanup(func() { SetLookup(nil) })
}

func postHandler(h http.Handler, path string, body any) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", path, jsonReader(body))
	h.ServeHTTP(w, r)
	return w
}

// TestIssueRequiresEmailOnFile drives the OTP issue handler through chi-less
// dispatch so the email-required gate is exercised.
func TestIssueRequiresEmailOnFile(t *testing.T) {
	setTestLookup(t, func(identity string) (*Profile, bool) {
		if identity == "acct-nomail" {
			return &Profile{FRID: "fr-1", Email: ""}, true
		}
		return nil, false
	})
	w := postHandler(http.HandlerFunc(IssueHandler), "/api/auth/otp/issue", map[string]string{"identity": "acct-nomail", "purpose": "login"})
	var resp map[string]any
	_ = json.NewDecoder(w.Result().Body).Decode(&resp)
	if resp["status"] != "email_required" {
		t.Fatalf("expected email_required, got %v", resp["status"])
	}
	if msg, _ := resp["message"].(string); msg == "" {
		t.Fatalf("expected a message")
	}
}

func TestLoginRoundTripByPhone(t *testing.T) {
	setTestLookup(t, func(identity string) (*Profile, bool) {
		return &Profile{FRID: "fr-9", District: "wa", Email: "sam@x.com", AccountNumber: "ACCT9"}, true
	})
	os.Setenv("OTP_DEBUG", "1")
	defer os.Unsetenv("OTP_DEBUG")

	iss := postHandler(http.HandlerFunc(IssueHandler), "/api/auth/otp/issue", map[string]string{"identity": "+15095550199", "purpose": "login"})
	var issResp map[string]any
	_ = json.NewDecoder(iss.Result().Body).Decode(&issResp)
	if issResp["status"] != "sent" {
		t.Fatalf("issue: %v", issResp)
	}
	code, _ := issResp["debug_code"].(string)

	ver := postHandler(http.HandlerFunc(VerifyHandler), "/api/auth/otp/verify", map[string]string{"identity": "+15095550199", "purpose": "login", "code": code})
	var verResp map[string]any
	_ = json.NewDecoder(ver.Result().Body).Decode(&verResp)
	if verResp["status"] != "ok" {
		t.Fatalf("verify: %v", verResp)
	}
	if verResp["role"] != "Customer" {
		t.Fatalf("role: %v", verResp["role"])
	}
	if c := ver.Result().Cookies(); len(c) == 0 || c[0].Name != "session" {
		t.Fatalf("expected session cookie")
	}
}
