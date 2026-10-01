package view

// PestOverride carries code-level corrections for one pest catalog row.
// The catalog (name, scientific name, description, photo) lives in the
// production database; these overrides fix region accuracy and duplicate
// copy without touching the DB. ThreatLevel is never overridden here.
type PestOverride struct {
	ScientificName string   // region-accurate scientific name (researched, never invented)
	Description    string   // unique, region-accurate description (researched, never invented)
	Signs          []string // pest-specific signs of activity (replaces category generic)
	Treat          string   // pest-specific treatment notes
	Prev           []string // pest-specific prevention tips
}

var pestOverrides = map[string]PestOverride{}

// RegisterPestOverride records the override for one pest slug. Called from
// init() in pest_copy_*.go files. Later registrations win on collision.
func RegisterPestOverride(slug string, o PestOverride) {
	pestOverrides[slug] = o
}

// PestOverrideForSlug returns the override for a slug, if one is registered.
func PestOverrideForSlug(slug string) (PestOverride, bool) {
	o, ok := pestOverrides[slug]
	return o, ok
}

// pestCopyForSlug is the template-facing lookup used by the threat-file
// page: per-slug override first, category generic as fallback.
func pestCopyForSlug(slug, cat string) PestCopy {
	if o, ok := pestOverrides[slug]; ok && len(o.Signs) > 0 {
		return PestCopy{Signs: o.Signs, Treat: o.Treat, Prev: o.Prev}
	}
	return PestCopyFor(cat)
}
