package marketing

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// The structured consent record is the audit artifact for a checked
// box: it must land with the form's exact wording and version.
func TestRecordSMSConsentStoresStructuredRecord(t *testing.T) {
	m := &Module{DBPath: filepath.Join(t.TempDir(), "consent.db")}
	r := httptest.NewRequest("POST", "/free-quote", nil)
	r.RemoteAddr = "203.0.113.9:5555"
	m.recordSMSConsent(r, "Test Person", "+15095550100", "/free-quote", freeQuoteConsentWording)

	got, err := data.ListSMSConsents(m.DBPath, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 consent record, got %d", len(got))
	}
	rec := got[0]
	if rec.Wording != freeQuoteConsentWording {
		t.Fatalf("stored wording drifted from the form label: %q", rec.Wording)
	}
	if rec.WordingVersion != smsConsentWordingVersion || rec.Source != "/free-quote" ||
		rec.Phone != "+15095550100" || rec.Name != "Test Person" {
		t.Fatalf("record fields wrong: %+v", rec)
	}
	if rec.ConsentedAt == 0 {
		t.Fatal("consent record has no timestamp")
	}
}
