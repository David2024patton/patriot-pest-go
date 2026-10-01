package view

import "github.com/David2024patton/patriot-pest-go/internal/data"

// seo_links.go — code-level internal-link tables for the SEO pass.
//
// Same pattern as the pest_copy overrides: hand-curated {label, path} pairs
// checked into the repo, no DB writes. Every target is verified by
// TestInternalLinkTargetsResolve in the marketing package, so a bad slug can
// never ship a dead link.

// BlogRelatedLinks maps a blog post slug to hand-picked internal links shown
// under the article: the pest threat file it covers, a relevant service-area
// page, and the conversion page.
var BlogRelatedLinks = map[string][][2]string{
	"why-ants-invade-in-spring": {
		{"Ant Control", "/pest/ants"},
		{"Pest Control in Spokane, WA", "/areas/spokane"},
		{"Get a Free Quote", "/contact"},
	},
	"spiders-fall-guide": {
		{"Spider Control", "/pest/spiders"},
		{"Pest Control in Coeur d'Alene, ID", "/areas/coeur-d-alene"},
		{"Get a Free Quote", "/contact"},
	},
	"rodent-proof-your-home": {
		{"Rodent Control", "/pest/rodents"},
		{"Pest Control in Phoenix, AZ", "/areas/phoenix"},
		{"Get a Free Quote", "/contact"},
	},
}

// AreaPestLinks maps a state code to the pest slugs featured on that state's
// city pages ("Popular services in {city}"). Slugs are resolved against the
// live catalog by AreaPestsFor, so a catalog change can never render a dead
// link — unknown slugs are skipped.
var AreaPestLinks = map[string][]string{
	"WA": {"ants", "spiders", "rodents", "wasps"},
	"ID": {"ants", "spiders", "rodents", "wasps"},
	"OR": {"ants", "spiders", "rodents", "wasps"},
	"AZ": {"scorpions", "ants", "spiders", "rodents"},
}

// AreaPestsFor resolves the featured pest slugs for a state code to catalog
// rows, skipping any slug the catalog does not know.
func AreaPestsFor(code string) []data.Pest {
	var out []data.Pest
	for _, slug := range AreaPestLinks[code] {
		if p, ok := data.PestBySlug(slug); ok {
			out = append(out, p)
		}
	}
	return out
}
