package marketing

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	"github.com/David2024patton/patriot-pest-go/internal/view"
)

const (
	signupT = "Sign Up - Online Account | Patriot Pest Control"
	signupD = "Create a free Patriot Pest Control account online. Request service, track visits, and get member pricing. Veteran-owned, WA/ID/OR/AZ."
)

// signupGet — the public account signup page (no auth required).
func (m *Module) signupGet(w http.ResponseWriter, r *http.Request) {
	view.Page(w, r, "signup", signupT, signupD, metaKeywords, m.base(map[string]any{
		"Csrf":        m.csrfField(r, w),
		"PestOptions": pestOptionsHTML(""),
	}))
}

// signupPost validates the form, stores the signup as a website-sourced
// customer record (data.CreateSignup), captures it as an AlphaFlux CRM
// lead (leads.go), and fires generate_lead on success.
func (m *Module) signupPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Warn("signup form parse failed", "err", err.Error())
	}
	token := r.FormValue("_csrf")
	var expected string
	if c, err := r.Cookie(csrfCookieName); err == nil && c != nil {
		expected = c.Value
	}
	if !csrfOK(token, expected) {
		slog.Warn("CSRF check failed", "path", r.URL.Path, "ip", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(419)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or missing CSRF token. Please refresh and retry."})
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))
	city := strings.TrimSpace(r.FormValue("city"))
	state := strings.ToUpper(strings.TrimSpace(r.FormValue("state")))
	zip := strings.TrimSpace(r.FormValue("zip"))
	pest := strings.TrimSpace(r.FormValue("pest"))

	var errs []string
	if name == "" {
		errs = append(errs, "Name is required.")
	} else if runeLen(name) > 120 {
		errs = append(errs, "Name must be at most 120 characters.")
	}
	if email == "" {
		errs = append(errs, "Email is required.")
	} else if !emailRe.MatchString(email) {
		errs = append(errs, "Email must be a valid email address.")
	} else if runeLen(email) > 254 {
		errs = append(errs, "Email must be at most 254 characters.")
	}
	if phone == "" {
		errs = append(errs, "Phone is required.")
	} else if phoneDigits(phone) < 7 || phoneDigits(phone) > 15 {
		errs = append(errs, "Phone must be a valid phone number.")
	}
	if state != "" && (len(state) != 2 || !isUpperLetters(state)) {
		errs = append(errs, "State must be a 2-letter code.")
	}
	if zip != "" && (runeLen(zip) > 10 || !isZipSafe(zip)) {
		errs = append(errs, "ZIP must be a valid ZIP code.")
	}
	if !validPest(pest) {
		errs = append(errs, "Please choose a pest from the list.")
	}

	if len(errs) > 0 {
		view.Page(w, r, "signup", signupT, signupD, metaKeywords, m.base(map[string]any{
			"Errors":  errs,
			"OldName": name, "OldEmail": email, "OldPhone": phone,
			"OldCity": city, "OldState": state, "OldZip": zip, "OldPest": pest,
			"Csrf": m.csrfField(r, w), "PestOptions": pestOptionsHTML(pest),
		}))
		return
	}

	res := data.CreateSignup(m.DBPath, name, email, phone, city, state, zip)
	if !res.OK {
		// The user gets the generic message; the log gets the real error.
		slog.Error("signup failed", "email", email, "err", res.Err)
		view.PageStatus(w, r, 422, "signup", signupT, signupD, metaKeywords, m.base(map[string]any{
			"Errors":  []string{res.Message},
			"OldName": name, "OldEmail": email, "OldPhone": phone,
			"OldCity": city, "OldState": state, "OldZip": zip, "OldPest": pest,
			"Csrf": m.csrfField(r, w), "PestOptions": pestOptionsHTML(pest),
		}))
		return
	}

	// Option C: the signup also becomes a CRM lead in AlphaFlux with
	// website_form attribution (and the console's instant alert). The
	// local record above stays as the site's own copy; a capture
	// failure is logged, never shown to the visitor.
	notes := ""
	if city != "" || state != "" {
		notes = "From " + strings.TrimSpace(city+" "+state)
	}
	if err := m.captureLead(leadForm{
		Name: name, Phone: phone, Email: email,
		Detail: leadDetail("Signup form", pest, zip),
		Notes:  notes,
	}); err != nil {
		slog.Error("signup lead capture failed (local record kept)", "err", err, "email", email)
	}

	slog.Info("Website signup", "email", email, "name", name, "city", city, "state", state)
	view.Page(w, r, "signup", signupT, signupD, metaKeywords, m.base(map[string]any{
		"Success":        res.Message,
		"AnalyticsEvent": "generate_lead",
		"PestOptions":    pestOptionsHTML(""),
	}))
}

func isUpperLetters(s string) bool {
	for _, c := range s {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

func isZipSafe(s string) bool {
	digits := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			digits++
		} else if c != '-' {
			return false
		}
	}
	return digits > 0
}