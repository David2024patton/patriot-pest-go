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
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/David2024patton/patriot-pest-go/internal/config"
	"github.com/David2024patton/patriot-pest-go/internal/data"
	custommw "github.com/David2024patton/patriot-pest-go/internal/middleware"
	"github.com/David2024patton/patriot-pest-go/internal/modules/health"
	"github.com/David2024patton/patriot-pest-go/internal/modules/marketing"
	"github.com/David2024patton/patriot-pest-go/internal/view"
)

// movedPaths are the URLs that used to serve a dashboard, login flow or admin
// console from this app. They are kept as redirects so bookmarks, ads and
// inbound links never 404.
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
	// The landing pages render from the SQLite catalog (pest library, posts, areas).
	data.Load(cfg.DBPath)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// One place owns the "sign in" destination: the AlphaFlux platform.
	view.SetLoginURL(cfg.LoginURL)

	r := chi.NewRouter()
	r.Use(custommw.RequestID)
	r.Use(custommw.SlogLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(custommw.Timeout(15 * time.Second))
	r.Use(custommw.SecurityHeaders)
	r.Use(custommw.CORS)

	if (&health.Module{}).Register(r) {
		logger.Info("module enabled", "module", "health")
	}

	// Marketing — the public site. Flag-gated like every module.
	if mkt := (&marketing.Module{Enabled: cfg.MarketingEnabled}); mkt.Register(r) {
		logger.Info("module enabled", "module", "marketing")
	}

	loginRedirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, cfg.LoginURL, http.StatusFound)
	})
	for _, p := range movedPaths {
		r.Handle(p, loginRedirect)
	}
	// Admin console, staff tools and the removed API surface.
	r.Handle("/admin", loginRedirect)
	r.Handle("/admin/*", loginRedirect)
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
		logger.Info("listening", "addr", cfg.Addr, "env", cfg.Env, "login_url", cfg.LoginURL)
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
	_ = srv.Shutdown(ctx)
	logger.Info("shutdown complete")
}
