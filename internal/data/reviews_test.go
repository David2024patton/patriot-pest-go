package data

import (
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGoogleReviewsRoundtrip(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "reviews.db")
	in := []GoogleReview{
		{Reviewer: "New Person", Rating: 5, Text: "Great work.", DateLabel: "a day ago"},
		{Reviewer: "Older Person", Rating: 5, Text: "Also great.", DateLabel: "2 weeks ago"},
		{Reviewer: "Unhappy Person", Rating: 1, Text: "Bad.", DateLabel: "3 weeks ago"},
	}
	if err := ReplaceGoogleReviews(dbPath, GoogleReviewsProfile{Total: 111, Average: 4.8}, in); err != nil {
		t.Fatalf("replace: %v", err)
	}
	got, profile, err := LatestGoogleReviews(dbPath, 5, 6)
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if profile.Total != 111 || profile.Average != 4.8 {
		t.Errorf("profile = %+v, want total 111 average 4.8", profile)
	}
	if len(got) != 2 {
		t.Fatalf("latest returned %d reviews, want 2 (5-star filter)", len(got))
	}
	if got[0].Reviewer != "New Person" || got[1].Reviewer != "Older Person" {
		t.Errorf("order = %q, %q; want newest first (New Person, Older Person)", got[0].Reviewer, got[1].Reviewer)
	}
	if got[0].Text != "Great work." || got[0].DateLabel != "a day ago" {
		t.Errorf("first review = %+v, want text and date label preserved", got[0])
	}
}

func TestGoogleReviewsReplaceSwapsWholeSnapshot(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "reviews.db")
	first := []GoogleReview{{Reviewer: "First", Rating: 5, Text: "one", DateLabel: "a day ago"}}
	second := []GoogleReview{{Reviewer: "Second", Rating: 5, Text: "two", DateLabel: "an hour ago"}}
	if err := ReplaceGoogleReviews(dbPath, GoogleReviewsProfile{Total: 1, Average: 5}, first); err != nil {
		t.Fatalf("replace 1: %v", err)
	}
	if err := ReplaceGoogleReviews(dbPath, GoogleReviewsProfile{Total: 2, Average: 5}, second); err != nil {
		t.Fatalf("replace 2: %v", err)
	}
	got, _, err := LatestGoogleReviews(dbPath, 5, 6)
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if len(got) != 1 || got[0].Reviewer != "Second" {
		t.Errorf("after second replace got %+v, want only Second", got)
	}
}

func TestGoogleReviewsEmptyDBIsEmptyNotError(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "reviews.db")
	got, profile, err := LatestGoogleReviews(dbPath, 5, 6)
	if err != nil {
		t.Fatalf("latest on fresh db: %v", err)
	}
	if len(got) != 0 || profile.Total != 0 {
		t.Errorf("fresh db got %+v / %+v, want empty", got, profile)
	}
}
