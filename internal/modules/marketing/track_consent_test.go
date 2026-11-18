package marketing

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// The track endpoints store nothing unless the payload carries the
// consent marker the first-party trackers add after a banner accept.
func TestTrackEndpointRequiresConsentMarker(t *testing.T) {
	m := &Module{DBPath: filepath.Join(t.TempDir(), "track.db")}
	h := m.trackEndpoint("view")
	post := func(body string) {
		req := httptest.NewRequest(http.MethodPost, "/api/track/view", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", browserUA)
		rec := httptest.NewRecorder()
		h(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("status = %d, want 202 (the endpoint always answers 202)", rec.Code)
		}
	}
	count := func() int {
		db, err := data.AnalyticsDB(m.DBPath)
		if err != nil {
			t.Fatal(err)
		}
		if err := data.EnsureAnalyticsTables(db); err != nil {
			t.Fatal(err)
		}
		var n int
		if err := db.QueryRow(`SELECT COUNT(*) FROM analytics_events`).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}

	post(`{"vid":"v1","sid":"s1","path":"/no-consent"}`)
	if n := count(); n != 0 {
		t.Fatalf("event without consent marker was stored (%d rows)", n)
	}
	post(`{"consent":false,"vid":"v1","sid":"s1","path":"/declined"}`)
	if n := count(); n != 0 {
		t.Fatalf("event with consent:false was stored (%d rows)", n)
	}
	post(`{"consent":true,"vid":"v1","sid":"s1","path":"/accepted"}`)
	if n := count(); n != 1 {
		t.Fatalf("event with consent:true not stored (%d rows)", n)
	}
}
