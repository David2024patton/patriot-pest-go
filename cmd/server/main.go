// Command server is the public landing site for patriotpest.pro.
//
// Scope: marketing pages, content catalog (pests / service areas / blog),
// lead capture, and the PWA shell. Customer, staff and admin dashboards — and
// all FieldRoutes / Twilio integration — live in the AlphaFlux platform, so the
// site links out to it and every URL that used to serve a dashboard now
// redirects there.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/David2024patton/patriot-pest-go/internal/config"
	"github.com/David2024patton/patriot-pest-go/internal/data"
	custommw "github.com/David2024patton/patriot-pest-go/internal/middleware"
	"github.com/David2024patton/patriot-pest-go/internal/modules/health"
	"github.com/David2024patton/patriot-pest-go/internal/modules/marketing"
	"github.com/David2024patton/patriot-pest-go/internal/modules/report"
	"github.com/David2024patton/patriot-pest-go/internal/view"
)

// version is the build stamp, injected with:
//
//	go build -ldflags "-X main.version=$(git rev-parse --short HEAD)"
var version = "dev"

// movedPaths are the URLs that used to serve a dashboard, login flow or admin
// console from this app. They are kept as redirects so bookmarks, ads and
// inbound links never 404. They point at the marketing home page until the
// AlphaFlux platform is ready to host accounts.
var movedPaths = []string{
	"/login", "/login/verify", "/logout",
	"/customer-auth", "/customer-verify", "/account",
	"/customer-dashboard", "/customer-portal",
	"/staff", "/staff-verify", "/staff-logout", "/staff-dashboard", "/dashboard",
	"/su", "/su/verify",
}

func main() {
	// Health-only mode: the container HEALTHCHECK runs "/app/patriot-server -health".
	// Check the catalog and exit so the probe never binds :3000.
	if len(os.Args) > 1 && os.Args[1] == "-health" {
		cfg := config.Load()
		if _, err := os.Stat(cfg.DBPath); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel(cfg)}))
	// Coarse IP geolocation for analytics. The MaxMind database is
	// downloaded at Docker build time; when it is missing (local dev,
	// offline build) InitGeoDB is a no-op and geo fields stay empty.
	geoPath := os.Getenv("GEOIP_DB_PATH")
	if geoPath == "" {
		geoPath = "geo/GeoLite2-City.mmdb"
	}
	data.InitGeoDB(geoPath)
	if !data.GeoConfigured() {
		logger.Info("geo: no GeoIP database, country/city will report unknown", "path", geoPath)
	}
	slog.SetDefault(logger)

	// The landing pages render from the SQLite catalog (pest library, posts,
	// areas). Fail-open on a broken catalog — empty pages, not a dead site —
	// but a broken catalog must be loud in the logs, never silent.
	if _, err := data.Load(cfg.DBPath); err != nil {
		logger.Error("catalog load failed, running with empty catalog", "err", err, "db", cfg.DBPath)
	}

	r := chi.NewRouter()
	r.Use(custommw.RequestID)
	r.Use(custommw.SlogLogger(logger))
	r.Use(recoverer(logger))
	r.Use(custommw.Timeout(15 * time.Second))
	r.Use(custommw.SecurityHeaders)
	r.Use(custommw.CORS)

	// Status dashboard: served only on the report.patriotpest.pro host,
	// behind the report module's email-OTP login. Every other host falls
	// through to the normal site untouched.
	(&report.Module{DBPath: cfg.DBPath}).Register(r)

	if (&health.Module{}).Register(r) {
		logger.Info("module enabled", "module", "health")
	}

	// Marketing — the public site. Flag-gated like every module.
	if mkt := (&marketing.Module{Enabled: cfg.MarketingEnabled, DBPath: cfg.DBPath}); mkt.Register(r) {
		logger.Info("module enabled", "module", "marketing")
	}

	loginRedirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/", http.StatusFound)
	})
	for _, p := range movedPaths {
		r.Handle(p, loginRedirect)
	}
	// The removed API surface. NOTE: /admin is the live analytics console
	// (registered by the marketing module when ADMIN_EMAILS +
	// ADMIN_PASSWORD_HASH are configured); it is intentionally not redirected.
	r.Handle("/api/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte(`{"error":"moved","message":"This API has moved to the AlphaFlux platform."}`))
	}))

	// Project valuation tool — public page, kept with the site.
	r.Get("/cost", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(view.CostPage))
	})
	costData := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(view.CostPricingJSON))
	}
	r.Get("/cost/data/pricing.json", costData)
	r.Get("/data/pricing.json", costData)

	srv := &http.Server{Addr: cfg.Addr, Handler: r}
	go func() {
		logger.Info("listening", "addr", cfg.Addr, "env", cfg.Env, "version", version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown failed", "err", err)
	}
	logger.Info("shutdown complete")
}

// logLevel honors APP_LOG_LEVEL when set, else APP_DEBUG, else info.
func logLevel(cfg config.Config) slog.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_LOG_LEVEL"))) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	if cfg.Debug {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

// recoverer converts panics into 500s and logs them as structured JSON with
// a stack trace and request ID — chi's Recoverer only prints to stderr.
func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"err", fmt.Sprint(rec),
						"stack", string(debug.Stack()),
						"path", r.URL.Path,
						"request_id", custommw.RequestIDFromContext(r.Context()),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
