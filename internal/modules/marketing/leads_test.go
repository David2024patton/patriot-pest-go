package marketing

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCaptureLeadPostsToConsole(t *testing.T) {
	var gotAuth, gotPath string
	var gotBody map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"id":"crm_contact:x","created":true}`))
	}))
	defer srv.Close()

	m := &Module{AlphaFluxURL: srv.URL, AlphaFluxKey: "afx_test"}
	err := m.captureLead(leadForm{
		Name: "Sam Lead", Phone: "5095550199", Email: "sam@example.com",
		Detail: "Contact form | Pest: ants | ZIP: 99006", Notes: "Ants in the kitchen",
	})
	if err != nil {
		t.Fatalf("captureLead: %v", err)
	}
	if gotPath != "/api/marketing/capture" {
		t.Errorf("path = %q", gotPath)
	}
	if gotAuth != "Bearer afx_test" {
		t.Errorf("auth header = %q", gotAuth)
	}
	if gotBody["channel"] != "website_form" {
		t.Errorf("channel = %q, want website_form", gotBody["channel"])
	}
	if gotBody["name"] != "Sam Lead" || gotBody["phone"] != "5095550199" || gotBody["email"] != "sam@example.com" {
		t.Errorf("body identity fields wrong: %v", gotBody)
	}
	if !strings.Contains(gotBody["detail"], "Pest: ants") || gotBody["notes"] != "Ants in the kitchen" {
		t.Errorf("detail/notes wrong: %v", gotBody)
	}
}

func TestCaptureLeadOffWithoutKey(t *testing.T) {
	m := &Module{AlphaFluxURL: "https://alphaflux.net"}
	if err := m.captureLead(leadForm{Name: "Sam"}); err != errLeadCaptureOff {
		t.Errorf("err = %v, want errLeadCaptureOff", err)
	}
}

func TestCaptureLeadSurfacesHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()
	m := &Module{AlphaFluxURL: srv.URL, AlphaFluxKey: "afx_test"}
	if err := m.captureLead(leadForm{Name: "Sam"}); err == nil {
		t.Error("expected an error on a 403 from the console")
	}
}

func TestLeadDetailAndPests(t *testing.T) {
	if got := leadDetail("Contact form", "Ants", "99006"); got != "Contact form | Pest: Ants | ZIP: 99006" {
		t.Errorf("detail = %q", got)
	}
	if got := leadDetail("Signup form", "", ""); got != "Signup form" {
		t.Errorf("bare detail = %q", got)
	}
	if !validPest("") || !validPest("Ants") || !validPest("Other") {
		t.Error("picker values must validate")
	}
	if validPest("Unicorns") {
		t.Error("invented pest must not validate")
	}
	html := string(pestOptionsHTML("Ants"))
	for _, o := range pestOptions {
		// Values render HTML-escaped (& as &amp;); the browser decodes
		// them back to the exact pestOptions string on submit.
		if !strings.Contains(html, `value="`+template.HTMLEscapeString(o)+`"`) {
			t.Errorf("options HTML missing %q", o)
		}
	}
	if !strings.Contains(html, `value="Ants" selected`) {
		t.Error("submitted pest must render selected")
	}
}
