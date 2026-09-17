package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// requestsTotal counts every request through SlogLogger. Read by /metrics.
var requestsTotal atomic.Uint64

// RequestsTotal returns the lifetime request count.
func RequestsTotal() uint64 { return requestsTotal.Load() }

// requestIDKey is the context key carrying the request ID. Typed to avoid
// collisions with other context values. Middleware is the only writer.
type requestIDKeyT string

// RequestIDKey is the context key under which RequestID stores the request ID.
const RequestIDKey requestIDKeyT = "request_id"

// RequestIDFromContext returns the request ID stashed by the RequestID
// middleware, or "" when absent.
func RequestIDFromContext(ctx context.Context) string {
	if v := ctx.Value(RequestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// SlogLogger logs each request as JSON with X-Request-ID.
func SlogLogger(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &wrapWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(ww, r)
			requestsTotal.Add(1)
			l.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", RequestIDFromContext(r.Context()),
				"remote", r.RemoteAddr,
			)
		})
	}
}

type wrapWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrapWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Flush promotes http.Flusher through the logging wrapper — embedding the
// http.ResponseWriter interface alone does NOT promote Flusher, which would
// silently break SSE handlers downstream.
func (w *wrapWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// SecurityHeaders adds HSTS + secure cookie hints + CSP basics.
// The CSP enumerates the script/style/image hosts the site actually uses —
// no scheme wildcards, so a compromised CDN or ad host cannot inject JS.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com https://unpkg.com https://www.googletagmanager.com https://connect.facebook.net https://www.clarity.ms; "+
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
				"img-src 'self' data: https://*.tile.openstreetmap.org; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"connect-src 'self'; "+
				"frame-ancestors 'self';")
		next.ServeHTTP(w, r)
	})
}

// CORS allows same-origin + go/www/test.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "https://go.patriotpest.pro" || origin == "https://www.patriotpest.pro" || origin == "https://test.patriotpest.pro" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-CSRF-Token,X-Request-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequestID ensures X-Request-ID on every request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 8)
			if _, err := rand.Read(b); err != nil {
				slog.Error("request id generation failed", "err", err.Error())
				id = "req-fallback"
			} else {
				id = hex.EncodeToString(b)
			}
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RateLimiter is a simple in-memory fixed window per client IP.
type RateLimiter struct {
	mu     sync.Mutex
	count  map[string]int
	reset  map[string]time.Time
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{count: make(map[string]int), reset: make(map[string]time.Time), limit: limit, window: window}
}

// clientIP strips the port from RemoteAddr so one client is one bucket.
func clientIP(r *http.Request) string {
	if h, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return h
	}
	return r.RemoteAddr
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		now := time.Now()
		rl.mu.Lock()
		if now.After(rl.reset[ip]) {
			rl.count[ip] = 0
			rl.reset[ip] = now.Add(rl.window)
		}
		rl.count[ip]++
		n := rl.count[ip]
		// Evict stale buckets so the maps cannot grow without bound.
		for k, exp := range rl.reset {
			if now.After(exp) {
				delete(rl.reset, k)
				delete(rl.count, k)
			}
		}
		rl.mu.Unlock()
		if n > rl.limit {
			http.Error(w, `{"error":"rate limited"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Timeout wrapper.
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, `{"error":"timeout"}`)
	}
}

// TimeoutExcept applies http.TimeoutHandler except on the given path prefixes.
// Streaming endpoints (SSE) must be exempt: http.TimeoutHandler swaps the
// ResponseWriter for one that does not implement http.Flusher, which breaks
// event streaming with a 500.
func TimeoutExcept(d time.Duration, skipPrefixes ...string) func(http.Handler) http.Handler {
	timeout := func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, `{"error":"timeout"}`)
	}
	return func(next http.Handler) http.Handler {
		tw := timeout(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, p := range skipPrefixes {
				if strings.HasPrefix(r.URL.Path, p) {
					next.ServeHTTP(w, r)
					return
				}
			}
			tw.ServeHTTP(w, r)
		})
	}
}
