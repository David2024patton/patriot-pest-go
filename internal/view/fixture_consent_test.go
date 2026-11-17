package view

import (
	"bytes"
	"html/template"
	"os"
	"testing"
)

/*
	TestWriteConsentFixture renders the gated analytics define with

plainly fake IDs (UA-TEST, 0000000000) into a minimal page for the
browser proof in proof_consent.py. Runs only when CONSENT_FIXTURE_OUT
names the path.
*/
func TestWriteConsentFixture(t *testing.T) {
	out := os.Getenv("CONSENT_FIXTURE_OUT")
	if out == "" {
		t.Skip("CONSENT_FIXTURE_OUT not set")
	}
	tmpl, err := template.New("fx").Parse(namedTemplates)
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]any{"AnalyticsOn": true, "GTag": "UA-TEST", "LoaderID": "UA-TEST", "FBPixel": "0000000000"}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "analytics", data); err != nil {
		t.Fatal(err)
	}
	page := "<!doctype html><html><head><title>Fixture</title>" + buf.String() +
		"</head><body><h1>Patriot Pest Control fixture</h1></body></html>"
	if err := os.WriteFile(out, []byte(page), 0o644); err != nil {
		t.Fatal(err)
	}
}
