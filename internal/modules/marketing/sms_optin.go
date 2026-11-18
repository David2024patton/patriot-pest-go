package marketing

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/David2024patton/patriot-pest-go/internal/view"
)

const (
	smsOptInT = "Text Message Opt-In | Patriot Pest Control"
	smsOptInD = "How Patriot Pest Control customers consent to text messages: by phone, when starting a service plan, or with the opt-in form on this page. Message frequency varies. Reply STOP to end, HELP for help."
)

// smsOptInGet — the public SMS opt-in proof page (/sms-opt-in). It exists
// so A2P 10DLC reviewers can verify exactly how consent is given: the
// program, the opt-in paths, the verbal script, and this working form.
func (m *Module) smsOptInGet(w http.ResponseWriter, r *http.Request) {
	view.Page(w, r, "sms-opt-in", smsOptInT, smsOptInD, metaKeywords, m.base(map[string]any{
		"Csrf": m.csrfField(r, w),
	}))
}

// smsOptInPost — CSRF-protected web opt-in. The consent checkbox is the
// opt-in instrument, so it is enforced server-side: an unchecked box is
// an error, never a silent pass. A successful opt-in is captured as a
// website lead through the same CRM path as the contact form, with the
// consent wording recorded in the lead notes. If the CRM capture fails
// the visitor is told plainly and given the phone number; the sign-up is
// never shown as saved when it was not.
func (m *Module) smsOptInPost(w http.ResponseWriter, r *http.Request) {
	render := func(extra map[string]any) {
		view.Page(w, r, "sms-opt-in", smsOptInT, smsOptInD, metaKeywords, m.base(extra))
	}
	if err := r.ParseForm(); err != nil {
		slog.Warn("sms opt-in form parse failed", "err", err.Error())
		render(map[string]any{
			"Errors": []string{"Could not read the form. Please try again."},
			"Csrf":   m.csrfField(r, w),
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
	phone := strings.TrimSpace(r.FormValue("phone"))
	consent := r.FormValue("sms_consent")

	var errs []string
	if name == "" {
		errs = append(errs, "Name is required.")
	} else if runeLen(name) > 120 {
		errs = append(errs, "Name must be at most 120 characters.")
	}
	if phone == "" {
		errs = append(errs, "Mobile number is required.")
	} else if phoneDigits(phone) < 7 || phoneDigits(phone) > 15 {
		errs = append(errs, "Mobile number must be a valid phone number.")
	}
	if consent != "yes" {
		errs = append(errs, "Please check the box to say yes to text messages.")
	}

	if len(errs) > 0 {
		render(map[string]any{
			"Errors":   errs,
			"OldName":  name,
			"OldPhone": phone,
			"Csrf":     m.csrfField(r, w),
		})
		return
	}

	// The audit-grade consent record (who, when, which form, exact
	// wording) is stored on its own; the lead note below stays for
	// staff readability.
	m.recordSMSConsent(r, name, phone, "/sms-opt-in", smsOptInConsentWording)

	lead := leadForm{
		Name:   name,
		Phone:  phone,
		Detail: leadDetail("SMS opt-in form (/sms-opt-in)", "", ""),
		Notes:  "Web opt-in at /sms-opt-in: the customer checked the consent box agreeing to text messages from Patriot Pest Control about appointments and service. Message frequency varies. Message and data rates may apply. Reply STOP to end, HELP for help. Consent is not a condition of purchase.",
	}
	if err := m.captureLead(lead); err != nil {
		slog.Error("sms opt-in lead capture failed", "err", err, "phone", phone)
		render(map[string]any{
			"Errors":   []string{"We could not save your sign-up just now. Please call (509) 471-5767 and we will add you to the texting list."},
			"OldName":  name,
			"OldPhone": phone,
			"Csrf":     m.csrfField(r, w),
		})
		return
	}
	slog.Info("SMS opt-in submitted", "phone", phone)
	render(map[string]any{
		"Success":        "Thank you. You are signed up for text messages from Patriot Pest Control. Reply STOP at any time to end them.",
		"AnalyticsEvent": "generate_lead",
		"Csrf":           m.csrfField(r, w),
	})
}
