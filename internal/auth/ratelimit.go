package auth

// RateLimiter — sliding-window attempt limiting for the web login flow. Port
// of app/Core/RateLimiter.php, backed by an in-memory map instead of the
// login_attempts table (the Go app has no DB yet; parity of semantics is what
// matters: count failed attempts per key inside a window, block past threshold).

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateLimiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time // key -> timestamps of failed attempts
}

var rl = rateLimiter{hits: map[string][]time.Time{}}

// clientIp returns the best-effort real client IP. Behind the reverse proxy
// REMOTE_ADDR is the proxy itself, so trust the leftmost X-Forwarded-For entry.
func clientIp(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		first := strings.TrimSpace(parts[0])
		if first != "" {
			return first
		}
	}
	parts := strings.Split(r.RemoteAddr, ":")
	host := parts[0]
	if host == "" {
		return "0.0.0.0"
	}
	return host
}

// tooMany reports whether $key has >= maxAttempts failures inside the window.
func (l *rateLimiter) tooMany(key string, maxAttempts, windowSec int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.recent(key, windowSec)) >= maxAttempts
}

// retryAfter returns seconds until the key is allowed again (0 = not locked).
// Mirrors PHP: wait from the newest recorded failure plus the full window.
func (l *rateLimiter) retryAfter(key string, maxAttempts, windowSec int) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	recent := l.recent(key, windowSec)
	if len(recent) < maxAttempts {
		return 0
	}
	wait := int(recent[len(recent)-1].Add(time.Duration(windowSec) * time.Second).Sub(time.Now()).Seconds())
	if wait < 0 {
		return 0
	}
	return wait
}

// hit records a failed attempt for $key.
func (l *rateLimiter) hit(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hits[key] = append(l.hits[key], time.Now())
}

// clear removes recorded failures for $key (called after a successful auth).
func (l *rateLimiter) clear(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.hits, key)
}

// recent returns hits inside the window (caller holds lock).
func (l *rateLimiter) recent(key string, windowSec int) []time.Time {
	return l.recentLocked(key, int64(time.Now().Add(-time.Duration(windowSec)*time.Second).Unix()))
}

// recentLocked drops timestamps older than since (caller holds lock).
func (l *rateLimiter) recentLocked(key string, since int64) []time.Time {
	cut := time.Unix(since, 0)
	var out []time.Time
	for _, t := range l.hits[key] {
		if !t.Before(cut) {
			out = append(out, t)
		}
	}
	return out
}
