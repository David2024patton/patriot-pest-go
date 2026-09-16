package auth

import (
	"context"
	"strings"
	"testing"
)

// TestOTPRoundTrip proves a brand-new identity can issue + verify an OTP and
// obtain a session, with single-use semantics — i.e. "create a new account for
// testing" works without any pre-seeded record.
func TestOTPRoundTrip(t *testing.T) {
	id := "test-fr-portal@new.example.com"
	ctx := context.Background()

	code, err := Issue(ctx, id, "login", 6, 600)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("expected 6-digit code, got %q", code)
	}

	ok, err := Verify(ctx, id, "login", code, 5)
	if !ok || err != nil {
		t.Fatalf("verify: ok=%v err=%v", ok, err)
	}
	sess := CreateSession(id, "Customer", 900)
	got, found := GetSession(sess.ID)
	if !found || got.Identity != id || got.Role != "Customer" {
		t.Fatalf("session not retrievable: found=%v got=%+v", found, got)
	}

	// Single-use: a second verify of the same code must fail.
	if ok, _ := Verify(ctx, id, "login", code, 5); ok {
		t.Fatalf("code was single-use but verified twice")
	}
}

// TestOTPBadCode proves a wrong code is rejected (constant-time compare path).
func TestOTPBadCode(t *testing.T) {
	id := "bad@new.example.com"
	ctx := context.Background()
	code, _ := Issue(ctx, id, "login", 6, 600)
	bad := strings.Repeat("0", 6)
	if code != bad {
		if ok, _ := Verify(ctx, id, "login", bad, 5); ok {
			t.Fatalf("wrong code should not verify")
		}
	}
}

// TestSuperLoginDigits proves the super-login path issues 8 digits.
func TestSuperLoginDigits(t *testing.T) {
	id := "su@new.example.com"
	ctx := context.Background()
	code, _ := Issue(ctx, id, "super_login", 8, 300)
	if len(code) != 8 {
		t.Fatalf("expected 8-digit super code, got %q", code)
	}
}
