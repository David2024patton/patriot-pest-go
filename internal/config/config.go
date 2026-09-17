// Package config implements env-based configuration for the public landing site.
//
// Integration settings (FieldRoutes, Twilio, mail, staff auth) deliberately do
// not live here any more: those surfaces moved to the AlphaFlux platform.
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
