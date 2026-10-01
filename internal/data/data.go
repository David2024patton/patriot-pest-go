// Package data loads the content catalog (pest library + blog posts + site
// settings) from the SQLite database that ships with the template. It mirrors
// the PHP app's pest_photos / posts tables so the Go port renders identical
// content. Loading is fail-open: a missing DB yields an empty catalog rather
// than blocking a page.
package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	_ "modernc.org/sqlite" // pure-Go SQLite driver
)

// Pest — one row of pest_photos.
type Pest struct {
	ID             int
	Slug           string
	Name           string
	ScientificName string
	Filename       string
	Description    string
	Category       string
	ThreatLevel    int
	SortOrder      int
}

// Post — one row of posts joined with its pest photo.
type Post struct {
	Slug         string
	Title        string
	Excerpt      string
	BodyHTML     string
	Photo        string // pest_photos.filename
	PestName     string
	PestSlug     string
	Season       string
	PestCategory string
	Author       string
	PublishedAt  string
}

// Store — the cached catalog, loaded once at boot.
type Store struct {
	loaded     bool
	Pests      []Pest
	Posts      []Post
	pestBySlug map[string]Pest
	postBySlug map[string]Post
}

var (
	mu   sync.RWMutex
	load = &Store{}

	// catalogDB stays open for the life of the process so /ready can verify
	// the catalog is reachable. Guarded by mu.
	catalogDB *sql.DB
)

// htmlSanitizer strips scripts, event handlers and javascript: URLs from blog
// body HTML at load time. Content is sanitized on save, but the save path
// lives outside this repo — read-time sanitizing is the backstop.
var htmlSanitizer = bluemonday.UGCPolicy()

// Load reads the catalog from dbPath (database/patriot.db). Safe to call once;
// subsequent calls return the cached store. An error is returned when the
// catalog cannot be read — the caller decides whether to fail or run degraded,
// but a broken catalog must never be silent.
func Load(dbPath string) (*Store, error) {
	mu.Lock()
	defer mu.Unlock()
	if load.loaded {
		return load, nil
	}
	load.loaded = true
	if strings.ContainsAny(dbPath, "?#") {
		return load, fmt.Errorf("refusing db path with url metacharacters: %q", dbPath)
	}
	db, err := sql.Open("sqlite", "file:"+dbPath+"?_pragma=journal_mode(WAL)&_busy_timeout=5000")
	if err != nil {
		return load, fmt.Errorf("open catalog db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return load, fmt.Errorf("ping catalog db: %w", err)
	}
	catalogDB = db
	load.Pests = queryPests(db)
	load.Posts = queryPosts(db)
	load.pestBySlug = mapPests(load.Pests)
	load.postBySlug = mapPosts(load.Posts)
	return load, nil
}

// Ping verifies the catalog database is reachable. Used by /ready.
func Ping() error {
	mu.RLock()
	db := catalogDB
	mu.RUnlock()
	if db == nil {
		return fmt.Errorf("catalog db not loaded")
	}
	return db.Ping()
}

func mapPests(ps []Pest) map[string]Pest {
	m := make(map[string]Pest, len(ps))
	for _, p := range ps {
		m[p.Slug] = p
	}
	return m
}

func mapPosts(ps []Post) map[string]Post {
	m := make(map[string]Post, len(ps))
	for _, p := range ps {
		m[p.Slug] = p
	}
	return m
}

func queryPests(db *sql.DB) []Pest {
	rows, err := db.Query("SELECT id, slug, name, scientific_name, filename, description, category, threat_level, sort_order FROM pest_photos ORDER BY sort_order ASC, id ASC")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Pest
	for rows.Next() {
		var p Pest
		var sci *string
		if err := rows.Scan(&p.ID, &p.Slug, &p.Name, &sci, &p.Filename, &p.Description, &p.Category, &p.ThreatLevel, &p.SortOrder); err != nil {
			slog.Error("catalog pest row scan failed", "err", err.Error())
			continue
		}
		if sci != nil {
			p.ScientificName = *sci
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		slog.Error("catalog pest query failed", "err", err.Error())
	}
	return out
}

func queryPosts(db *sql.DB) []Post {
	rows, err := db.Query(`
		SELECT p.slug, p.title, p.excerpt, p.body_html, p.author, p.published_at, p.season, p.pest_category,
		       COALESCE(ph.filename,''), COALESCE(ph.name,''), COALESCE(ph.slug,'')
		FROM posts p
		LEFT JOIN pest_photos ph ON ph.id = p.pest_photo_id
		ORDER BY p.published_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Post
	for rows.Next() {
		var p Post
		var author *string
		if err := rows.Scan(&p.Slug, &p.Title, &p.Excerpt, &p.BodyHTML, &author, &p.PublishedAt, &p.Season, &p.PestCategory, &p.Photo, &p.PestName, &p.PestSlug); err != nil {
			slog.Error("catalog post row scan failed", "err", err.Error())
			continue
		}
		if author != nil {
			p.Author = *author
		}
		p.BodyHTML = htmlSanitizer.Sanitize(p.BodyHTML)
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		slog.Error("catalog post query failed", "err", err.Error())
	}
	return out
}

// ---------- accessors ----------

// AllPests returns the catalog (empty if DB missing).
func AllPests() []Pest { mu.RLock(); defer mu.RUnlock(); return load.Pests }

// AllPosts returns published posts newest first.
func AllPosts() []Post { mu.RLock(); defer mu.RUnlock(); return load.Posts }

// IsDraftPost reports whether a post is a test/draft entry. Drafts use a
// "test-" slug prefix and are hidden from the blog index, sitemap, RSS, and
// direct URLs (they 404) so test content never leaks to search engines.
func IsDraftPost(p Post) bool { return strings.HasPrefix(p.Slug, "test-") }

// PublishedPosts returns posts excluding drafts, newest first.
func PublishedPosts() []Post {
	mu.RLock()
	defer mu.RUnlock()
	var out []Post
	for _, p := range load.Posts {
		if !IsDraftPost(p) {
			out = append(out, p)
		}
	}
	return out
}

// PublishedPostBySlug returns the post with a slug, or ok=false when the
// slug is unknown or belongs to a draft.
func PublishedPostBySlug(slug string) (Post, bool) {
	p, ok := PostBySlug(slug)
	if !ok || IsDraftPost(p) {
		return Post{}, false
	}
	return p, true
}

// PestBySlug returns the pest with a slug, or ok=false.
func PestBySlug(slug string) (Pest, bool) { mu.RLock(); defer mu.RUnlock(); p, ok := load.pestBySlug[slug]; return p, ok }

// PostBySlug returns the post with a slug, or ok=false.
func PostBySlug(slug string) (Post, bool) { mu.RLock(); defer mu.RUnlock(); p, ok := load.postBySlug[slug]; return p, ok }

// RelatedPests returns up to n other pests for a pest page.
func RelatedPests(slug string, n int) []Pest {
	mu.RLock()
	defer mu.RUnlock()
	var out []Pest
	for _, p := range load.Pests {
		if p.Slug != slug {
			out = append(out, p)
			if len(out) >= n {
				break
			}
		}
	}
	return out
}

// RelatedPosts returns up to n other posts for a blog post page.
func RelatedPosts(slug string, n int) []Post {
	mu.RLock()
	defer mu.RUnlock()
	var out []Post
	for _, p := range load.Posts {
		if p.Slug != slug {
			out = append(out, p)
			if len(out) >= n {
				break
			}
		}
	}
	return out
}

// State — a service area (state + its cities).
type State struct {
	Code     string // WA / ID / OR / AZ
	Name     string
	Cities   []string
}

// States — the four-state service footprint, mirroring PageController::states().
func States() []State {
	return []State{
		{"WA", "Washington", []string{"Spokane", "Spokane Valley", "Cheney", "Liberty Lake", "Airway Heights", "Medical Lake", "Deer Park", "Mead"}},
		{"ID", "Idaho", []string{"Coeur d'Alene", "Post Falls", "Hayden", "Rathdrum"}},
		{"OR", "Oregon", []string{"Hermiston", "Milton-Freewater"}},
		{"AZ", "Arizona", []string{"Phoenix"}},
	}
}

// CitySlug slugs a city for /areas/{slug} routing.
func CitySlug(city string) string {
	var b []byte
	for i := 0; i < len(city); i++ {
		c := city[i]
		switch {
		case c >= 'A' && c <= 'Z':
			b = append(b, c-'A'+('a'))
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9':
			b = append(b, c)
		default:
			b = append(b, '-')
		}
	}
	s := string(b)
	return trimLeadTail(s, "-")
}

func trimLeadTail(s string, cut string) string {
	lo := 0
	for lo < len(s) && s[lo] == cut[0] {
		lo++
	}
	hi := len(s)
	for hi > lo && s[hi-1] == cut[0] {
		hi--
	}
	return s[lo:hi]
}

// FindCity resolves a city slug to its (city, state, stateName) or ok=false.
func FindCity(slug string) (city, code, stateName string, ok bool) {
	for _, st := range States() {
		for _, c := range st.Cities {
			if CitySlug(c) == slug {
				return c, st.Code, st.Name, true
			}
		}
	}
	return "", "", "", false
}

// _ keeps time used for future TTL cache.
var _ = time.Second
