// Package config implements env-based configuration for the public landing site.
//
// Integration settings (FieldRoutes, Twilio, mail, staff auth) deliberately do
// not live here any more: those surfaces moved to the AlphaFlux platform. The
// one exception is the AlphaFlux lead capture pair below: the site posts
// website form leads INTO the platform, so it needs the platform's API base
// URL and a tenant-scoped API key (crm:write only).
package config

import (
	"os"
	"strings"
)

// Config holds all runtime configuration. Never commit .env; inject via env.
type Config struct {
	Env    string
	Debug  bool
	AppURL string
	Addr   string
	DBPath string // SQLite content catalog: pests, posts, service areas

	// MarketingEnabled gates the public pages module.
	MarketingEnabled bool

	// AlphaFluxAPIURL is the AlphaFlux console API base the site posts
	// form leads to (POST /api/marketing/capture).
	AlphaFluxAPIURL string
	// AlphaFluxAPIKey is the Patriot-tenant API key for lead capture
	// (scope crm:write). Empty disables capture; form handlers then fall
	// back to the local signup store so a lead is never discarded.
	AlphaFluxAPIKey string
}

// Load reads env vars with defaults.
func Load() Config {
	return Config{
		Env:              env("APP_ENV", "local"),
		Debug:            envBool("APP_DEBUG", false),
		AppURL:           env("APP_URL", "https://patriotpest.pro"),
		Addr:             env("ADDR", ":3000"),
		DBPath:           env("DB_PATH", "database/patriot.db"),
		MarketingEnabled: envBool("MARKETING_ENABLED", true),
		AlphaFluxAPIURL:  env("ALPHAFLUX_API_URL", "https://alphaflux.net"),
		AlphaFluxAPIKey:  env("ALPHAFLUX_API_KEY", ""),
	}
}

func (c Config) IsProduction() bool { return strings.EqualFold(c.Env, "production") }
func (c Config) IsLocal() bool      { return strings.EqualFold(c.Env, "local") }

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envBool(k string, def bool) bool {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "yes")
}
