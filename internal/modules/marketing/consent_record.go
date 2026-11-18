package marketing

// Structured SMS consent recording (critic cycle 1, finding 2).
//
// The consent checkbox labels below are quoted verbatim from the
// public forms: freeQuoteConsentWording is the label beside the box
// on /free-quote (internal/view/pages_freequote.go) and
// smsOptInConsentWording is the label on /sms-opt-in
// (internal/view/pages_sms_optin.go). If a label ever changes, change
// its constant in the same commit and bump smsConsentWordingVersion:
// the stored record must always carry the exact words the visitor
// agreed to, or the record proves nothing.

import (
	"log/slog"
	"net/http"

	"github.com/David2024patton/patriot-pest-go/internal/data"
)

// smsConsentWordingVersion tags the wording set currently live on
// the public forms.
const smsConsentWordingVersion = "2026-11 v1"

const freeQuoteConsentWording = "Yes, Patriot Pest Control may text me about my quote, my appointments, and my service. Message frequency varies. Message and data rates may apply. Reply STOP to end, HELP for help. Consent is not a condition of purchase."

const smsOptInConsentWording = "Yes, Patriot Pest Control may text me about my appointments and service. Message frequency varies. Message and data rates may apply. Reply STOP to end, HELP for help. Consent is not a condition of purchase."

// recordSMSConsent stores one structured consent record for a checked
// box. A store failure is logged loudly with the phone number so staff
// can reconcile by hand; it never blocks the visitor's own request,
// because the lead itself still carries the consent note.
func (m *Module) recordSMSConsent(r *http.Request, name, phone, source, wording string) {
	rec := data.SMSConsentRecord{
		Name:           name,
		Phone:          phone,
		Source:         source,
		Wording:        wording,
		WordingVersion: smsConsentWordingVersion,
		IP:             clientIP(r),
	}
	if err := data.StoreSMSConsent(m.DBPath, rec); err != nil {
		slog.Error("sms consent record store failed", "err", err, "phone", phone, "source", source)
		return
	}
	slog.Info("sms consent recorded", "source", source, "phone", phone, "wording_version", smsConsentWordingVersion)
}
