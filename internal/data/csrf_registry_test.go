package data

import "testing"

func TestCSRFRegistry(t *testing.T) {
	if CSRFTokenRegistered("never-issued") {
		t.Fatal("unregistered token verified")
	}
	RegisterCSRFToken("issued-token-abc")
	if !CSRFTokenRegistered("issued-token-abc") {
		t.Fatal("registered token did not verify")
	}
	if CSRFTokenRegistered("issued-token-abd") {
		t.Fatal("a different token verified against the registry")
	}
	RegisterCSRFToken("")
	if CSRFTokenRegistered("") {
		t.Fatal("empty token verified")
	}
}
