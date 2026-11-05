package marketing

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func syncRequest(token, body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/reviews/sync", strings.NewReader(body))
	if token != "" {
		req.Header.Set("X-Reviews-Token", token)
	}
	return req
}

const syncBody = `{"profile":{"total":111,"average":4.8},"reviews":[
	{"reviewer":"Robert Davis","rating":5,"text":"Excellent service. Highly recommended.","date":"a week ago"},
	{"reviewer":"Quiet One","rating":5,"text":"","date":"2 weeks ago"},
	{"reviewer":"Grumpy","rating":1,"text":"Nope.","date":"3 weeks ago"}
]}`

func TestReviewsSyncRequiresToken(t *testing.T) {
	m := &Module{DBPath: filepath.Join(t.TempDir(), "r.db"), ReviewsSyncToken: "sekret"}
	rec := httptest.NewRecorder()
	m.reviewsSync(rec, syncRequest("", syncBody))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("no token: got %d, want 401", rec.Code)
	}
	rec = httptest.NewRecorder()
	m.reviewsSync(rec, syncRequest("wrong", syncBody))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong token: got %d, want 401", rec.Code)
	}
}

func TestReviewsSyncStoresAndRailFilters(t *testing.T) {
	m := &Module{DBPath: filepath.Join(t.TempDir(), "r.db"), ReviewsSyncToken: "sekret"}
	rec := httptest.NewRecorder()
	m.reviewsSync(rec, syncRequest("sekret", syncBody))
	if rec.Code != http.StatusOK {
		t.Fatalf("sync: got %d (%s), want 200", rec.Code, rec.Body.String())
	}
	cards, profile := m.reviewCards(6)
	if profile.Total != 111 || profile.Average != 4.8 {
		t.Errorf("profile = %+v, want 111 / 4.8", profile)
	}
	if len(cards) != 2 {
		t.Fatalf("rail cards = %d, want 2 (only the 5-star reviews)", len(cards))
	}
	if cards[0].Reviewer != "Robert Davis" || cards[0].Stars != "★★★★★" {
		t.Errorf("first card = %+v, want Robert Davis with 5 stars", cards[0])
	}
	if cards[1].Reviewer != "Quiet One" || cards[1].Text != "" {
		t.Errorf("second card = %+v, want the textless 5-star review kept", cards[1])
	}
}

func TestReviewsSyncRejectsBadEntries(t *testing.T) {
	m := &Module{DBPath: filepath.Join(t.TempDir(), "r.db"), ReviewsSyncToken: "sekret"}
	for name, body := range map[string]string{
		"not json":    `{`,
		"no reviewer": `{"reviews":[{"reviewer":" ","rating":5}]}`,
		"bad rating":  `{"reviews":[{"reviewer":"A","rating":9}]}`,
	} {
		rec := httptest.NewRecorder()
		m.reviewsSync(rec, syncRequest("sekret", body))
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: got %d, want 400", name, rec.Code)
		}
	}
}
