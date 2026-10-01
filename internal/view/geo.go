package view

import (
	"net/http"
	"os"
	"strings"
)

// Region lines, mirroring Core\Geo. Patriot runs two phone lines: the main
// WA/ID/OR line and a dedicated Arizona line.
type regionLine struct {
	Display string
	Tel     string
	Label   string
}

var regions = map[string]regionLine{
	"wa": {"(509) 818-0993", "+15098180993", "WA, ID, OR"},
	"az": {"(602) 755-8414", "+16027558414", "ARIZONA"},
}

const geoCookie = "ppc_geo"

// resolveRegion picks the visitor's region: session cookie, then GEO_FORCE_REGION,
// else the default main line. Fail-open — never blocks a page.
func resolveRegion(r *http.Request) string {
	if c, err := r.Cookie(geoCookie); err == nil {
		if strings.ToLower(c.Value) == "az" || strings.ToLower(c.Value) == "wa" {
			return strings.ToLower(c.Value)
		}
	}
	if f := os.Getenv("GEO_FORCE_REGION"); strings.ToLower(f) == "az" || strings.ToLower(f) == "wa" {
		return strings.ToLower(f)
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

// LineFor returns the regional line for a state code (AZ gets the Arizona
// line; everything else the WA/ID/OR main line). Used by /areas/{slug} pages.
func LineFor(stateCode string) regionLine {
	if strings.EqualFold(stateCode, "AZ") {
		return regions["az"]
	}
	return regions["wa"]
}

// PhoneFor resolves the visitor's phone context.
func PhoneFor(r *http.Request) Phone {
	reg := resolveRegion(r)
	p := regions[reg]
	other := regions["wa"]
	if reg == "wa" {
		other = regions["az"]
	}
	return Phone{Display: p.Display, Tel: "tel:" + p.Tel, Label: p.Label, IsAZ: reg == "az", Other: other}
}
