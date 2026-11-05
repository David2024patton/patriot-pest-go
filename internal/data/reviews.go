package data

// Google reviews snapshot — the pop homepage's Field Reports rail renders
// from these tables so it always shows the newest real 5-star Google
// reviews instead of hardcoded testimonial copy.
//
// The site itself holds no Google credentials and never calls Google. The
// snapshot is replaced wholesale by the marketing module's token-gated
// sync endpoint (POST /api/reviews/sync), which the daily Google Business
// Profile review check feeds with the newest reviews it reads in the
// profile UI. Reviews are stored newest-first: position 0 is the newest
// review captured, and LatestGoogleReviews preserves that order.

import (
	"database/sql"
	"fmt"
	"time"
)

// GoogleReview is one captured Google review, verbatim as displayed.
type GoogleReview struct {
	Reviewer  string
	Rating    int // 1..5
	Text      string
	DateLabel string // as displayed by Google, e.g. "2 weeks ago"
}

// GoogleReviewsProfile is the profile-level aggregate shown beside the
// rail ("4.8 stars · 111 Google reviews").
type GoogleReviewsProfile struct {
	Total   int
	Average float64
}

// EnsureReviewsTables creates the snapshot schema when missing. Safe to
// run on every read and write; production databases pick the tables up
// with no migration step.
func EnsureReviewsTables(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS google_reviews (
			position INTEGER PRIMARY KEY,
			reviewer TEXT NOT NULL,
			rating INTEGER NOT NULL,
			body TEXT NOT NULL DEFAULT '',
			date_label TEXT NOT NULL DEFAULT '',
			synced_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS google_reviews_meta (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			total INTEGER NOT NULL DEFAULT 0,
			average REAL NOT NULL DEFAULT 0,
			synced_at INTEGER NOT NULL
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("reviews schema: %w", err)
		}
	}
	return nil
}

// ReplaceGoogleReviews swaps the whole snapshot in one transaction: the
// rail can never show a half-updated mix of two captures. reviews must
// already be ordered newest-first; index in the slice becomes position.
func ReplaceGoogleReviews(dbPath string, profile GoogleReviewsProfile, reviews []GoogleReview) error {
	db, err := AnalyticsDB(dbPath)
	if err != nil {
		return err
	}
	if err := EnsureReviewsTables(db); err != nil {
		return err
	}
	now := time.Now().Unix()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("reviews replace: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM google_reviews`); err != nil {
		return fmt.Errorf("reviews replace clear: %w", err)
	}
	stmt, err := tx.Prepare(`INSERT INTO google_reviews (position, reviewer, rating, body, date_label, synced_at) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("reviews replace prepare: %w", err)
	}
	defer stmt.Close()
	for i, rv := range reviews {
		if _, err := stmt.Exec(i, rv.Reviewer, rv.Rating, rv.Text, rv.DateLabel, now); err != nil {
			return fmt.Errorf("reviews replace insert %d: %w", i, err)
		}
	}
	if _, err := tx.Exec(`INSERT INTO google_reviews_meta (id, total, average, synced_at) VALUES (1, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET total = excluded.total, average = excluded.average, synced_at = excluded.synced_at`,
		profile.Total, profile.Average, now); err != nil {
		return fmt.Errorf("reviews replace meta: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("reviews replace commit: %w", err)
	}
	return nil
}

// LatestGoogleReviews returns up to limit stored reviews at or above
// minRating, newest first (position order). A database with no snapshot
// yet returns empty results and no error, so the rail simply hides until
// the first sync lands.
func LatestGoogleReviews(dbPath string, minRating, limit int) ([]GoogleReview, GoogleReviewsProfile, error) {
	db, err := AnalyticsDB(dbPath)
	if err != nil {
		return nil, GoogleReviewsProfile{}, err
	}
	if err := EnsureReviewsTables(db); err != nil {
		return nil, GoogleReviewsProfile{}, err
	}
	var profile GoogleReviewsProfile
	row := db.QueryRow(`SELECT total, average FROM google_reviews_meta WHERE id = 1`)
	if err := row.Scan(&profile.Total, &profile.Average); err != nil && err != sql.ErrNoRows {
		return nil, GoogleReviewsProfile{}, fmt.Errorf("reviews meta read: %w", err)
	}
	rows, err := db.Query(`SELECT reviewer, rating, body, date_label FROM google_reviews
		WHERE rating >= ? ORDER BY position ASC LIMIT ?`, minRating, limit)
	if err != nil {
		return nil, GoogleReviewsProfile{}, fmt.Errorf("reviews read: %w", err)
	}
	defer rows.Close()
	var out []GoogleReview
	for rows.Next() {
		var rv GoogleReview
		if err := rows.Scan(&rv.Reviewer, &rv.Rating, &rv.Text, &rv.DateLabel); err != nil {
			return nil, GoogleReviewsProfile{}, fmt.Errorf("reviews scan: %w", err)
		}
		out = append(out, rv)
	}
	if err := rows.Err(); err != nil {
		return nil, GoogleReviewsProfile{}, fmt.Errorf("reviews rows: %w", err)
	}
	return out, profile, nil
}
