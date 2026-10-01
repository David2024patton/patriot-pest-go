package marketing

// main_test.go — seeds a realistic content catalog (pests + posts) once per
// test run so SEO tests exercise the real handlers and templates. The data
// package caches the catalog process-wide after the first Load, so seeding
// here in TestMain is deterministic; no existing test depends on an empty
// catalog.

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	_ "modernc.org/sqlite"
)

var seedPests = [][2]string{
	{"ants", "Ants"}, {"spiders", "Spiders"}, {"rodents", "Rodents"},
	{"bed-bugs", "Bed Bugs"}, {"termites", "Termites"}, {"mosquitoes", "Mosquitoes"},
	{"wasps", "Wasps & Hornets"}, {"cockroaches", "Cockroaches"},
	{"fleas-ticks", "Fleas & Ticks"}, {"silverfish", "Silverfish"},
	{"pantry-pests", "Pantry Pests"}, {"bees", "Bees"},
	{"squirrels", "Squirrels"}, {"raccoons", "Raccoons"}, {"bats", "Bats"},
	{"scorpions", "Scorpions"}, {"crickets", "Crickets"}, {"pack-rats", "Pack Rats"},
	{"box-elder-bugs", "Boxelder Bugs"}, {"stink-bugs", "Stink Bugs"},
	{"fruit-flies", "Fruit Flies"}, {"gophers", "Gophers"}, {"moles", "Moles"},
	{"hornets", "Hornets"}, {"yellow-jackets", "Yellow Jackets"},
}

var seedPosts = [][3]string{
	{"why-ants-invade-in-spring", "Why Ants Invade in Spring", "spring"},
	{"spiders-fall-guide", "Fall Spider Guide", "fall"},
	{"rodent-proof-your-home", "Rodent-Proof Your Home", "winter"},
}

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ppc-seo-catalog")
	if err != nil {
		fmt.Fprintln(os.Stderr, "seed temp dir:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, "catalog.db")
	if err := seedCatalog(dbPath); err != nil {
		fmt.Fprintln(os.Stderr, "seed catalog:", err)
		os.Exit(1)
	}
	if _, err := data.Load(dbPath); err != nil {
		fmt.Fprintln(os.Stderr, "data.Load:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func seedCatalog(dbPath string) error {
	db, err := sql.Open("sqlite", "file:"+dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE pest_photos (
		id INTEGER PRIMARY KEY, slug TEXT, name TEXT, scientific_name TEXT,
		filename TEXT, description TEXT, category TEXT,
		threat_level INTEGER, sort_order INTEGER)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE posts (
		slug TEXT PRIMARY KEY, title TEXT, excerpt TEXT, body_html TEXT,
		author TEXT, published_at TEXT, season TEXT, pest_category TEXT,
		pest_photo_id INTEGER)`); err != nil {
		return err
	}
	for i, p := range seedPests {
		if _, err := db.Exec(
			`INSERT INTO pest_photos (slug, name, filename, description, category, threat_level, sort_order)
			 VALUES (?, ?, 'test.jpg', 'Seed description for tests.', 'insect', 50, ?)`,
			p[0], p[1], i); err != nil {
			return err
		}
	}
	for i, p := range seedPosts {
		if _, err := db.Exec(
			`INSERT INTO posts (slug, title, excerpt, body_html, author, published_at, season, pest_category)
			 VALUES (?, ?, 'Seed excerpt.', '<p>Seed body.</p>', 'Test', ?, ?, 'insect')`,
			p[0], p[1], fmt.Sprintf("2026-0%d-01", i+1), p[2]); err != nil {
			return err
		}
	}
	return nil
}
