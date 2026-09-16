package customers

import (
	"testing"

	"github.com/David2024patton/patriot-pest-go/internal/fieldroutes"
)

func ptr(s string) *string { return &s }

func TestResolveByEmailPhoneAccount(t *testing.T) {
	upsertRows([]fieldroutes.Row{
		{FRID: "fr-1", District: "wa", Name: "Jane", AccountNumber: "ACCT1", Email: ptr("Jane@Example.com"), Phone: ptr("+15095550199")},
		{FRID: "fr-2", District: "az", Name: "NoMail", AccountNumber: "ACCT2"}, // no email
	})

	// email — case-insensitive
	row, ok := Resolve("jane@example.com")
	if !ok || row.FRID != "fr-1" {
		t.Fatalf("email resolve: ok=%v row=%+v", ok, row)
	}
	// account number
	row, ok = Resolve("acct1")
	if !ok || row.FRID != "fr-1" {
		t.Fatalf("account resolve: ok=%v", ok)
	}
	// phone with different formatting resolves to the same record
	row, ok = Resolve("(509) 555-0199")
	if !ok || row.FRID != "fr-1" {
		t.Fatalf("phone resolve: ok=%v row=%+v", ok, row)
	}
	// unknown
	if _, ok := Resolve("nobody"); ok {
		t.Fatalf("expected unknown identity")
	}
}

func TestLookupReturnsProfile(t *testing.T) {
	upsertRows([]fieldroutes.Row{{FRID: "fr-9", District: "wa", Name: "Sam", AccountNumber: "ACCT9", Email: ptr("sam@x.com")}})
	p, ok := Lookup("sam@x.com")
	if !ok || p == nil {
		t.Fatalf("lookup failed: ok=%v", ok)
	}
	if p.FRID != "fr-9" || p.Email != "sam@x.com" || p.AccountNumber != "ACCT9" {
		t.Fatalf("bad profile: %+v", p)
	}
	if p.District != "wa" {
		t.Fatalf("expected district wa, got %q", p.District)
	}
}
