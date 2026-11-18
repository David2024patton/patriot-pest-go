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
	freeQuoteT = "Free Quote | Patriot Pest Control"
	freeQuoteD = "Get a free, no-obligation quote from Patriot Pest Control for your home in Washington, Idaho or Oregon. Tell us the pest and your ZIP and we call you back with a straight price."
)

// freeQuoteGet — the dedicated free quote page (/free-quote). Every
// "Get a Free Quote" button on the site points here since 2026-11-12;
// before that they landed on the contact form and there was no quote
// page at all (Skyler asked, David ordered it built).
func (m *Module) freeQuoteGet(w http.ResponseWriter, r *http.Request) {
	view.Page(w, r, "free-quote", freeQuoteT, freeQuoteD, metaKeywords, m.base(map[string]any{
		"Csrf":        m.csrfField(r, w),
		"PestOptions": pestOptionsHTML(""),
	}))
}

// freeQuotePost — CSRF-protected quote request. The lead goes through
// the same capture as the contact form (CRM contact with website_form
// attribution, and the console fires its instant SMS alert to staff),
// with the contact form's local-store fallback so a request is never
// discarded; if both stores fail the visitor is told plainly to call.
// The texting checkbox is optional on purpose: consent is never a
// condition of getting a quote. It is recorded in the lead notes with
// the exact wording when checked, and the notes say call-only when it
// is not, so staff never text a quote lead who did not opt in.
func (m *Module) freeQuotePost(w http.ResponseWriter, r *http.Request) {
	render := func(extra map[string]any) {
		view.Page(w, r, "free-quote", freeQuoteT, freeQuoteD, metaKeywords, m.base(extra))
	}
	if err := r.ParseForm(); err != nil {
		slog.Warn("free quote form parse failed", "err", err.Error())
		render(map[string]any{
			"Errors":      []string{"Could not read the form. Please try again."},
			"Csrf":        m.csrfField(r, w),
			"PestOptions": pestOptionsHTML(""),
		})
		return
	}
	token := r.FormValue("_csrf")
	var expected string
	if c, err := r.Cookie(csrfCookieName); err == nil && c != nil {
		expected = c.Value
	}
	if !csrfOK(token, expected) {
		slog.Warn("CSRF check failed", "path", r.URL.Path, "ip", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(419) // PHP Csrf::verifyOrDie uses 419 for bad tokens
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or missing CSRF token. Please refresh and retry."})
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))
	address := strings.TrimSpace(r.FormValue("address"))
	city := strings.TrimSpace(r.FormValue("city"))
	zip := strings.TrimSpace(r.FormValue("zip"))
	pest := strings.TrimSpace(r.FormValue("pest"))
	message := strings.TrimSpace(r.FormValue("message"))
	consent := r.FormValue("sms_consent")

	var errs []string
	if name == "" {
		errs = append(errs, "Name is required.")
	} else if runeLen(name) > 120 {
		errs = append(errs, "Name must be at most 120 characters.")
	}
	if email != "" && !emailRe.MatchString(email) {
		errs = append(errs, "Email must be a valid email address.")
	} else if runeLen(email) > 254 {
		errs = append(errs, "Email must be at most 254 characters.")
	}
	if phone == "" {
		errs = append(errs, "Phone is required.")
	} else if phoneDigits(phone) < 7 || phoneDigits(phone) > 15 {
		errs = append(errs, "Phone must be a valid phone number.")
	}
	if runeLen(address) > 200 {
		errs = append(errs, "Address must be at most 200 characters.")
	}
	if runeLen(city) > 80 {
		errs = append(errs, "City must be at most 80 characters.")
	}
	if zip != "" && (runeLen(zip) > 10 || !isZipSafe(zip)) {
		errs = append(errs, "ZIP must be a valid ZIP code.")
	}
	if !validPest(pest) {
		errs = append(errs, "Please choose a pest from the list.")
	}
	if runeLen(message) > 2000 {
		errs = append(errs, "Notes must be at most 2000 characters.")
	}

	if len(errs) > 0 {
		render(map[string]any{
			"Errors":  errs,
			"OldName": name, "OldEmail": email, "OldPhone": phone,
			"OldAddress": address, "OldCity": city, "OldZip": zip,
			"OldMessage":  message,
			"Csrf":        m.csrfField(r, w),
			"PestOptions": pestOptionsHTML(pest),
		})
		return
	}

	var notes strings.Builder
	notes.WriteString("Free quote request from /free-quote.")
	if address != "" || city != "" {
		notes.WriteString("\nService address: ")
		notes.WriteString(strings.Trim(strings.Join([]string{address, city}, ", "), ", "))
		if zip != "" {
			notes.WriteString(" " + zip)
		}
	}
	if message != "" {
		notes.WriteString("\nCustomer notes: " + message)
	}
	if consent == "yes" {
		notes.WriteString("\nTexting consent: the customer checked the box on /free-quote agreeing to text messages from Patriot Pest Control about their quote, appointments, and service. Message frequency varies. Message and data rates may apply. Reply STOP to end, HELP for help. Consent is not a condition of purchase.")
		// The note above is for staff; the audit-grade record (who,
		// when, which form, exact wording) is stored separately.
		m.recordSMSConsent(r, name, phone, "/free-quote", freeQuoteConsentWording)
	} else {
		notes.WriteString("\nTexting consent: NOT given on the form. Call this lead; do not text until they opt in.")
	}

	lead := leadForm{
		Name:   name,
		Phone:  phone,
		Email:  email,
		Detail: leadDetail("Free quote form (/free-quote)", pest, zip),
		Notes:  notes.String(),
	}
	if err := m.captureLead(lead); err != nil {
		slog.Error("free quote lead capture failed, saving locally", "err", err, "phone", phone)
		res := data.CreateSignup(m.DBPath, name, email, phone, city, "", zip)
		if !res.OK {
			slog.Error("free quote lead local save failed", "err", res.Err, "phone", phone)
			render(map[string]any{
				"Errors":  []string{"We could not save your request just now. Please call (509) 471-5767 and we will price it on the spot."},
				"OldName": name, "OldEmail": email, "OldPhone": phone,
				"OldAddress": address, "OldCity": city, "OldZip": zip,
				"OldMessage":  message,
				"Csrf":        m.csrfField(r, w),
				"PestOptions": pestOptionsHTML(pest),
			})
			return
		}
	}
	slog.Info("Free quote form submitted", "phone", phone, "pest", pest, "zip", zip)
	render(map[string]any{
		"Success":        "Thank you. We got your quote request and will call you within one business day with your price, usually much faster.",
		"AnalyticsEvent": "generate_lead",
		"Csrf":           m.csrfField(r, w),
		"PestOptions":    pestOptionsHTML(""),
	})
}
