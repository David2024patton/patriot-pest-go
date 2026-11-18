package marketing

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// The forget endpoint erases exactly the rows stored under the
// presented visitor id and rejects ids that match neither tracker
// format.
func TestTrackForgetDeletesVisitorRows(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "forget.db")
	m := &Module{DBPath: dbPath}
	db, err := data.AnalyticsDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := data.EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	vid := "123e4567-e89b-12d3-a456-426614174000"
	for _, e := range []data.Event{
		{TS: 100, VisitorID: vid, SessionID: "s1", Kind: "pageview", Path: "/"},
		{TS: 101, VisitorID: "other-visitor", SessionID: "s2", Kind: "pageview", Path: "/"},
	} {
		if err := data.InsertEvent(db, e); err != nil {
			t.Fatal(err)
		}
	}

	post := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/track/forget", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		m.trackForget(rec, req)
		return rec
	}

	if rec := post(`{"visitor_id":"not-an-id"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad id: status = %d, want 400", rec.Code)
	}
	rec := post(`{"visitor_id":"` + vid + `"}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"deleted":1`) {
		t.Fatalf("forget: status = %d body = %s, want 200 with deleted 1", rec.Code, rec.Body.String())
	}
	var left int
	if err := db.QueryRow(`SELECT COUNT(*) FROM analytics_events`).Scan(&left); err != nil {
		t.Fatal(err)
	}
	if left != 1 {
		t.Fatalf("%d rows left, want 1 (the other visitor)", left)
	}
}
