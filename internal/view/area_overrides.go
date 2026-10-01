package view

// AreaCopy carries genuinely local content for one city page. The shared
// area-detail template renders these fields when present for the city slug;
// cities without an entry keep the generic template copy.
type AreaCopy struct {
	Intro     string // 2-4 sentences, genuinely local (researched, never invented)
	LocalNote string // neighborhoods, geography, landmarks (researched)
	PestNote  string // local pest pressure notes (researched)
}

var areaCopies = map[string]AreaCopy{}

// RegisterAreaCopy records local copy for one city slug. Called from init()
// in area_content.go. Later registrations win on collision.
func RegisterAreaCopy(slug string, c AreaCopy) {
	areaCopies[slug] = c
}

// AreaCopyFor returns the local copy for a city slug, if one is registered.
func AreaCopyFor(slug string) (AreaCopy, bool) {
	c, ok := areaCopies[slug]
	return c, ok
}
