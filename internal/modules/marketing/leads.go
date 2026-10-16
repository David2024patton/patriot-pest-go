package marketing

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

// leads.go — website lead capture into the AlphaFlux CRM (Option C).
//
// Every contact and signup form submission posts to the console's
// existing capture endpoint (POST /api/marketing/capture, channel
// website_form) with a Patriot-tenant API key, so the lead becomes a
// crm_contact with M1 source attribution (first touch never
// overwritten) and the console fires its instant SMS alert to the
// tenant's configured alert number. The key lives only in the app env
// (ALPHAFLUX_API_KEY); an empty key disables capture and the form
// handlers fall back to the local store so a lead is never discarded.

// leadForm is one form submission on its way to the CRM.
type leadForm struct {
	Name   string
	Phone  string
	Email  string
	Detail string // which form, pest, ZIP: "Contact form | Pest: ants | ZIP: 99006"
	Notes  string // the form's message; stored on the new CRM contact
}

// errLeadCaptureOff is returned when no AlphaFlux key is configured.
var errLeadCaptureOff = errors.New("lead capture not configured (ALPHAFLUX_API_KEY empty)")

// leadHTTP is the capture transport: short timeout, a form post must
// never hang on the CRM being slow.
var leadHTTP = &http.Client{Timeout: 5 * time.Second}

// captureLead posts the lead to the AlphaFlux capture endpoint. A nil
// error means the CRM accepted it; anything else is the caller's cue
// to fall back to the local store.
func (m *Module) captureLead(l leadForm) error {
	if m.AlphaFluxKey == "" {
		return errLeadCaptureOff
	}
	base := strings.TrimSuffix(m.AlphaFluxURL, "/")
	if base == "" {
		base = "https://alphaflux.net"
	}
	payload, err := json.Marshal(map[string]string{
		"name":    l.Name,
		"phone":   l.Phone,
		"email":   l.Email,
		"channel": "website_form",
		"detail":  l.Detail,
		"notes":   l.Notes,
	})
	if err != nil {
		return fmt.Errorf("lead capture encode: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, base+"/api/marketing/capture", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("lead capture request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.AlphaFluxKey)
	resp, err := leadHTTP.Do(req)
	if err != nil {
		return fmt.Errorf("lead capture post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lead capture status %d", resp.StatusCode)
	}
	return nil
}

// leadDetail builds the capture detail string: the form name plus the
// parts the visitor actually supplied.
func leadDetail(form, pest, zip string) string {
	parts := []string{form}
	if pest != "" {
		parts = append(parts, "Pest: "+pest)
	}
	if zip != "" {
		parts = append(parts, "ZIP: "+zip)
	}
	return strings.Join(parts, " | ")
}

// pestOptions is the pest picker list. The form templates render the
// same list through pestOptionsHTML; the handlers accept only these
// values (or empty), so a tampered post cannot invent a category.
var pestOptions = []string{
	"Ants", "Spiders", "Wasps & hornets", "Rodents", "Cockroaches",
	"Bed bugs", "Fleas & ticks", "Mosquitoes", "Box elder bugs",
	"Stink bugs", "Pantry pests", "Termites", "Bees", "Other",
}

// validPest reports whether p may be stored: empty (picker untouched)
// or one of pestOptions.
func validPest(p string) bool {
	if p == "" {
		return true
	}
	for _, o := range pestOptions {
		if p == o {
			return true
		}
	}
	return false
}

// pestOptionsHTML renders the picker's <option> elements with the
// submitted value pre-selected on a validation re-render.
func pestOptionsHTML(selected string) template.HTML {
	var b strings.Builder
	b.WriteString(`<option value="">Choose one (optional)</option>`)
	for _, o := range pestOptions {
		b.WriteString(`<option value="`)
		b.WriteString(template.HTMLEscapeString(o))
		b.WriteString(`"`)
		if o == selected {
			b.WriteString(` selected`)
		}
		b.WriteString(`>`)
		b.WriteString(template.HTMLEscapeString(o))
		b.WriteString(`</option>`)
	}
	return template.HTML(b.String())
}
