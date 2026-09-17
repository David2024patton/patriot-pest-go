package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	custommw "github.com/David2024patton/patriot-pest-go/internal/middleware"
)

// Module is always on — /health no auth + /ready DB ping + Prometheus /metrics.
type Module struct{}

func (m *Module) Register(r chi.Router) bool {
	r.Get("/health", m.health)
	r.Get("/ready", m.readyH)
	r.Get("/metrics", m.metrics)
	return true
}

func (m *Module) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":     "ok",
		"service":    "patriot-pest-go",
		"time":       time.Now().UTC().Format(time.RFC3339),
		"request_id": custommw.RequestIDFromContext(r.Context()),
	})
}

// readyH pings the catalog database: 503 when the catalog is unreachable so
// the orchestrator stops sending traffic instead of serving empty pages.
func (m *Module) readyH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := data.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "not ready", "db": "unreachable"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (m *Module) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = fmt.Fprintf(w, "# HELP patriot_requests_total\n# TYPE patriot_requests_total counter\npatriot_requests_total %d\n", custommw.RequestsTotal())
}
