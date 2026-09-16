// Search — site search across the pest library, blog guides, and service
// areas (Go port of PHP's SearchController). GET /search?q=...; results render
// through the shared layout like every other marketing page.
package marketing

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	"github.com/David2024patton/patriot-pest-go/internal/view"
)

// searchD — shared meta description for the search page.
const searchD = "Search Patriot Pest Control - blog guides, pest library, and service areas."

// searchPest / searchCity are display rows for the search template.
type searchPest struct {
	Slug       string
	Name       string
	Filename   string
	Desc       string // truncated for card display (80 chars + ellipsis)
}

type searchCity struct {
	City string
	Code string
}

// search — GET /search?q=... Empty query renders the intro copy; a query
// searches pest names/descriptions, blog titles/excerpts/bodies, and city
// names (case-insensitive substring), mirroring the live PHP behavior.
func (m *Module) search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	lq := strings.ToLower(q)

	var pests []searchPest
	var posts []data.Post
	var cities []searchCity
	if q != "" {
		for _, p := range data.AllPests() {
			if containsFold(p.Name, lq) || containsFold(p.Description, lq) {
				pests = append(pests, searchPest{Slug: p.Slug, Name: p.Name, Filename: p.Filename, Desc: truncateDesc(p.Description)})
			}
		}
		for _, p := range data.AllPosts() {
			if containsFold(p.Title, lq) || containsFold(p.Excerpt, lq) || containsFold(p.BodyHTML, lq) {
				posts = append(posts, p)
			}
		}
		for _, st := range data.States() {
			for _, c := range st.Cities {
				if containsFold(c, lq) {
					cities = append(cities, searchCity{City: c, Code: st.Code})
				}
			}
		}
	}

	title := "Patriot Pest Control"
	if q != "" {
		title = fmt.Sprintf("Search: %s | Patriot Pest Control", q)
	}
	view.Page(w, r, "search", title, searchD, metaKeywords, m.base(map[string]any{
		"Q":           q,
		"HasQuery":    q != "",
		"HasResults":  len(pests)+len(posts)+len(cities) > 0,
		"Pests":       pests,
		"Posts":       posts,
		"Cities":      cities,
	}))
}

// containsFold — case-insensitive substring match (lowercased needle).
func containsFold(haystack, needle string) bool { return strings.Contains(strings.ToLower(haystack), needle) }

// truncateDesc — card display truncation matching the live site: first 80
// characters plus an ellipsis when the description is longer.
func truncateDesc(s string) string {
	r := []rune(s)
	if len(r) > 80 {
		return string(r[:80]) + "…"
	}
	return s
}
