package data

import (
	"path/filepath"
	"testing"
)

func TestStoreAndListSMSConsent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "consent.db")
	rec := SMSConsentRecord{
		Name:           "Test Person",
		Phone:          "+15095550100",
		Source:         "/free-quote",
		Wording:        "Yes, Patriot Pest Control may text me.",
		WordingVersion: "2026-11 v1",
		IP:             "203.0.113.9",
	}
	if err := StoreSMSConsent(dbPath, rec); err != nil {
		t.Fatalf("store: %v", err)
	}
	got, err := ListSMSConsents(dbPath, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 record, got %d", len(got))
	}
	g := got[0]
	if g.Phone != rec.Phone || g.Source != rec.Source || g.Wording != rec.Wording ||
		g.WordingVersion != rec.WordingVersion || g.Name != rec.Name || g.IP != rec.IP {
		t.Fatalf("round trip mismatch: %+v", g)
	}
	if g.ConsentedAt == 0 {
		t.Fatal("consented_at was not stamped")
	}
}
