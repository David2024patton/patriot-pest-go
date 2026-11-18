package data

// CSRF token registry: the server-side half of the double-submit
// cookie pattern (critic cycle 1, finding 7).
//
// A cookie-only double-submit check proves nothing when an attacker
// can plant the cookie (a sibling subdomain, a network foothold):
// they set the cookie and submit the matching field, and the pair
// agrees because both halves are theirs. So every token the server
// issues is also registered here, by hash, and verification requires
// the token to be one the server actually issued within the last 24
// hours. A planted cookie value was never registered, so it fails.
//
// The registry is process memory on purpose: tokens are short-lived
// form credentials, and a restart simply makes open forms reload.

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// csrfTokenMaxAge bounds how long an issued token stays verifiable.
const csrfTokenMaxAge = 24 * time.Hour

var csrfRegistry = struct {
	sync.Mutex
	byHash    map[string]int64 // sha256 hex of token -> issued-at Unix seconds
	lastPrune int64
}{byHash: map[string]int64{}}

func csrfHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// RegisterCSRFToken records a freshly issued token. Called by every
// code path that mints a CSRF token (the marketing module's form
// field and the view package's field share this registry, because
// they share the cookie).
func RegisterCSRFToken(token string) {
	if token == "" {
		return
	}
	now := time.Now().Unix()
	csrfRegistry.Lock()
	defer csrfRegistry.Unlock()
	csrfRegistry.byHash[csrfHash(token)] = now
	if now-csrfRegistry.lastPrune > 3600 {
		csrfRegistry.lastPrune = now
		cutoff := now - int64(csrfTokenMaxAge.Seconds())
		for h, at := range csrfRegistry.byHash {
			if at < cutoff {
				delete(csrfRegistry.byHash, h)
			}
		}
	}
}

// CSRFTokenRegistered reports whether the server issued this token
// within the acceptance window.
func CSRFTokenRegistered(token string) bool {
	if token == "" {
		return false
	}
	csrfRegistry.Lock()
	defer csrfRegistry.Unlock()
	at, ok := csrfRegistry.byHash[csrfHash(token)]
	if !ok {
		return false
	}
	return time.Now().Unix()-at <= int64(csrfTokenMaxAge.Seconds())
}
