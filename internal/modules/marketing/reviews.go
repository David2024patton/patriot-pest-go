package marketing

// Google reviews sync — feeds the pop homepage's Field Reports rail.
//
// The site holds no Google credentials. The daily Google Business Profile
// review check (agent-side, owner-account browser session) reads the
// newest reviews in the profile UI and POSTs them here; this endpoint
// replaces the stored snapshot (internal/data/reviews.go) so the rail
// always shows the newest real 5-star reviews, verbatim.
//
// Fail closed by construction: the route is registered only when
// REVIEWS_SYNC_TOKEN is set (see Register), and every call must present
// that token in the X-Reviews-Token header (constant-time compare).

import (
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// reviewsSyncMaxReviews caps one capture. The profile UI read takes the
// newest dozen; a rail never needs more than the newest few 5-star ones.
const reviewsSyncMaxReviews = 20

type reviewsSyncPayload struct {
	Profile struct {
		Total   int     `json:"total"`
		Average float64 `json:"average"`
	} `json:"profile"`
	Reviews []struct {
		Reviewer  string `json:"reviewer"`
		Rating    int    `json:"rating"`
		Text      string `json:"text"`
		DateLabel string `json:"date"`
	} `json:"reviews"`
}

// reviewsSync replaces the stored Google reviews snapshot. Body is the
// newest reviews first, exactly as read from the profile.
func (m *Module) reviewsSync(w http.ResponseWriter, r *http.Request) {
	if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Reviews-Token")), []byte(m.ReviewsSyncToken)) != 1 {
		slog.Warn("reviews sync: bad token", "ip", clientIP(r))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var p reviewsSyncPayload
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10))
	if err := dec.Decode(&p); err != nil {
		slog.Warn("reviews sync: bad body", "err", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if len(p.Reviews) > reviewsSyncMaxReviews {
		http.Error(w, "too many reviews", http.StatusBadRequest)
		return
	}
	reviews := make([]data.GoogleReview, 0, len(p.Reviews))
	for _, rv := range p.Reviews {
		name := strings.TrimSpace(rv.Reviewer)
		if name == "" || len(name) > 100 || rv.Rating < 1 || rv.Rating > 5 {
			http.Error(w, "invalid review entry", http.StatusBadRequest)
			return
		}
		text := strings.TrimSpace(rv.Text)
		if len(text) > 1200 {
			text = text[:1200]
		}
		date := strings.TrimSpace(rv.DateLabel)
		if len(date) > 40 {
			date = date[:40]
		}
		reviews = append(reviews, data.GoogleReview{Reviewer: name, Rating: rv.Rating, Text: text, DateLabel: date})
	}
	profile := data.GoogleReviewsProfile{Total: p.Profile.Total, Average: p.Profile.Average}
	if err := data.ReplaceGoogleReviews(m.DBPath, profile, reviews); err != nil {
		slog.Error("reviews sync: store failed", "err", err.Error())
		http.Error(w, "store failed", http.StatusInternalServerError)
		return
	}
	slog.Info("reviews sync: snapshot replaced", "reviews", len(reviews), "total", profile.Total, "average", profile.Average)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "stored": len(reviews)})
}

// reviewCard is the Field Reports rail view model: one real Google
// review, stars and avatar initials pre-rendered so the template
// stays dumb.
type reviewCard struct {
	Reviewer  string
	Text      string
	DateLabel string
	Stars     string
	Initials  string
}

// initialsOf renders the avatar initials for a reviewer name: first
// letters of the first and last name parts, uppercased.
func initialsOf(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return ""
	}
	first := func(s string) string {
		r := []rune(s)
		if len(r) == 0 {
			return ""
		}
		return strings.ToUpper(string(r[0]))
	}
	if len(parts) == 1 {
		return first(parts[0])
	}
	return first(parts[0]) + first(parts[len(parts)-1])
}

// reviewCards reads the newest stored 5-star reviews for the pop home
// rail. A read failure degrades to an empty rail (the section hides)
// rather than breaking the homepage; the failure is logged.
func (m *Module) reviewCards(limit int) ([]reviewCard, data.GoogleReviewsProfile) {
	reviews, profile, err := data.LatestGoogleReviews(m.DBPath, 5, limit)
	if err != nil {
		slog.Error("reviews rail: read failed", "err", err.Error())
		return nil, data.GoogleReviewsProfile{}
	}
	cards := make([]reviewCard, 0, len(reviews))
	for _, rv := range reviews {
		cards = append(cards, reviewCard{
			Reviewer:  rv.Reviewer,
			Text:      rv.Text,
			DateLabel: rv.DateLabel,
			Stars:     strings.Repeat("★", rv.Rating),
			Initials:  initialsOf(rv.Reviewer),
		})
	}
	return cards, profile
}
