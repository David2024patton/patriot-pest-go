package view

import (
	"net/http"
	"os"
	"strings"
)

// Region lines, mirroring Core\Geo. Patriot has ONE phone line,
// (509) 471-5767, shown to every visitor. The Arizona operation was sold
// in 2026-10, so there is no Arizona line or region anymore.
type regionLine struct {
	Display string
	Tel     string
	Label   string
}

var regions = map[string]regionLine{
	"wa": {"(509) 471-5767", "+15094715767", "WA, ID, OR"},
}

const geoCookie = "ppc_geo"

// resolveRegion picks the visitor's region: session cookie, then GEO_FORCE_REGION,
// else the default main line. Fail-open — never blocks a page.
func resolveRegion(r *http.Request) string {
	if c, err := r.Cookie(geoCookie); err == nil {
		if strings.ToLower(c.Value) == "wa" {
			return "wa"
		}
	}
	if f := os.Getenv("GEO_FORCE_REGION"); strings.ToLower(f) == "wa" {
		return "wa"
	}
	return "wa"
}

// Phone is the localized line for this visitor (display, tel, label, isAZ).
type Phone struct {
	Display  string
	Tel      string // "tel:+1..."
	Label    string
	IsAZ     bool
	Other    regionLine // the other line, for footers showing both
}

// LineFor returns the line for a state code. There is only one line now,
// so every state gets it. Used by /areas/{slug} pages.
func LineFor(stateCode string) regionLine {
	return regions["wa"]
}

// PhoneFor resolves the visitor's phone context.
func PhoneFor(r *http.Request) Phone {
	reg := resolveRegion(r)
	p := regions[reg]
	other := regions["wa"]
	return Phone{Display: p.Display, Tel: "tel:" + p.Tel, Label: p.Label, IsAZ: false, Other: other}
}
