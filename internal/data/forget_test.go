package data

import (
	"path/filepath"
	"testing"
)

func TestDeleteVisitorAnalytics(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "forget.db")
	db, err := AnalyticsDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := EnsureAnalyticsTables(db); err != nil {
		t.Fatal(err)
	}
	for _, e := range []Event{
		{TS: 100, VisitorID: "visitor-a", SessionID: "s1", Kind: "pageview", Path: "/"},
		{TS: 101, VisitorID: "visitor-a", SessionID: "s1", Kind: "click", Path: "/prices"},
		{TS: 102, VisitorID: "visitor-b", SessionID: "s2", Kind: "pageview", Path: "/"},
	} {
		if err := InsertEvent(db, e); err != nil {
			t.Fatal(err)
		}
	}
	n, err := DeleteVisitorAnalytics(dbPath, "visitor-a")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if n != 2 {
		t.Fatalf("deleted %d rows, want 2", n)
	}
	var left int
	if err := db.QueryRow(`SELECT COUNT(*) FROM analytics_events`).Scan(&left); err != nil {
		t.Fatal(err)
	}
	if left != 1 {
		t.Fatalf("%d rows left, want 1 (the other visitor)", left)
	}
}
