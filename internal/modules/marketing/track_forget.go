package marketing

// POST /api/track/forget: the visitor-facing erasure path for the
// first-party analytics (critic cycle 1, finding 4). The privacy
// page's tracking-data control sends the visitor ids this browser
// holds; the server deletes every analytics row stored under them.
// Knowing the id is the credential: ids are random 122-bit or
// 256-bit values that are never published, so a request can only
// erase the rows of a browser that actually holds the id.

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// The two id shapes the first-party trackers issue: the beacon's
// UUID and the shell tracker's 64-character hex string.
var (
	visitorUUIDRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	visitorHexRe  = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

func validVisitorID(id string) bool {
	return visitorUUIDRe.MatchString(id) || visitorHexRe.MatchString(id)
}

func (m *Module) trackForget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(io.LimitReader(r.Body, 16*1024))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok":false,"error":"could not read request"}`))
		return
	}
	var p map[string]any
	if err := json.Unmarshal(body, &p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok":false,"error":"invalid request"}`))
		return
	}
	id, _ := p["visitor_id"].(string)
	if id == "" {
		id, _ = p["vid"].(string)
	}
	if !validVisitorID(id) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok":false,"error":"unknown visitor id"}`))
		return
	}
	n, err := data.DeleteVisitorAnalytics(m.DBPath, id)
	if err != nil {
		slog.Error("track forget: delete failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ok":false,"error":"delete failed"}`))
		return
	}
	// The id is a credential-grade random value; only the count is
	// logged, never the id itself.
	slog.Info("track forget: visitor rows deleted", "deleted", n)
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "deleted": n})
}
