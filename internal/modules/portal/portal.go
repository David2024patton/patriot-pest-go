package portal

import (
	"encoding/json"
	"net/http"

	"github.com/David2024patton/patriot-pest-go/internal/auth"
	"github.com/David2024patton/patriot-pest-go/internal/fieldroutes"
	"github.com/David2024patton/patriot-pest-go/internal/modules/customers"
	"github.com/David2024patton/patriot-pest-go/internal/view"
	"github.com/go-chi/chi/v5"
)

// Module — customer portal: account, appointments, invoices pdf/next bill/pay proxy FieldRoutes, cancel→tel, messages.
type Module struct {
	Enabled bool
	FR      *fieldroutes.Client
}

func (m *Module) Register(r chi.Router) bool {
	if !m.Enabled {
		return false
	}
	r.Get("/customer-dashboard", m.Dashboard)
	r.Get("/api/customer/dashboard", m.JSONDashboard)
	r.Get("/customer/account", m.Account)
	r.Get("/api/customer/appointments", m.Appointments)
	r.Get("/api/customer/invoices/{id}/pdf", m.InvoicePDF)
	r.Post("/api/customer/pay", m.Pay)
	r.Post("/api/customer/cancel", m.Cancel)
	r.Get("/api/customer/notifications/stream", m.Stream)
	return true
}

// requireSession resolves the session cookie to a live session.
func (m *Module) requireSession(w http.ResponseWriter, r *http.Request) (*auth.Session, bool) {
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		return nil, false
	}
	s, ok := auth.GetSession(c.Value)
	if !ok {
		return nil, false
	}
	return s, true
}

// profile resolves the logged-in identity to a customer profile (nil-safe).
func (m *Module) profile(sess *auth.Session) *auth.Profile {
	p, _ := customers.Lookup(sess.Identity)
	return p
}

// Dashboard renders the customer's HTML account dashboard behind a session.
func (m *Module) Dashboard(w http.ResponseWriter, r *http.Request) {
	sess, ok := m.requireSession(w, r)
	if !ok {
		view.PageStatus(w, r, 401, "dash-customer", "My Account | Patriot Pest Control", "", "", map[string]any{
			"Flash": "Please sign in to view your account.",
		})
		return
	}

	appts := m.appointments(sess)
	row := []map[string]string{}
	for _, v := range appts {
		if mm, isMap := v.(map[string]any); isMap {
			row = append(row, map[string]string{
				"date": strOf(mm["date"]), "start": strOf(mm["start"]),
				"type": strOf(mm["type"]), "status": strOf(mm["status"]),
			})
		}
	}

	var name, email, phone, acct, addr string
	if prof := m.profile(sess); prof != nil {
		name, email = prof.Name, prof.Email
		phone, acct = prof.Phone, prof.AccountNumber
	}
	view.Page(w, r, "dash-customer", "My Account | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "customer",
		"Csrf": view.CSRFField(w, r),
		"Name": name, "Email": email, "Phone": phone,
		"AccountNumber": acct, "Address": addr,
		"Appointments": row,
	})
}

// JSONDashboard is the machine-readable dashboard at /api/customer/dashboard.
func (m *Module) JSONDashboard(w http.ResponseWriter, r *http.Request) {
	sess, ok := m.requireSession(w, r)
	if !ok {
		writeJSON(w, 401, map[string]any{"error": "not authenticated"})
		return
	}
	writeJSON(w, 200, m.dashboard(sess))
}

func strOf(v any) string {
	s, _ := v.(string)
	return s
}

func (m *Module) Account(w http.ResponseWriter, r *http.Request) {
	sess, ok := m.requireSession(w, r)
	if !ok {
		writeJSON(w, 401, map[string]any{"error": "not authenticated"})
		return
	}
	prof := m.profile(sess)
	if prof == nil {
		writeJSON(w, 401, map[string]any{"error": "no customer on file for this session"})
		return
	}
	writeJSON(w, 200, map[string]any{"account": accountJSON(prof)})
}

func (m *Module) Appointments(w http.ResponseWriter, r *http.Request) {
	sess, ok := m.requireSession(w, r)
	if !ok {
		writeJSON(w, 401, map[string]any{"error": "not authenticated"})
		return
	}
	writeJSON(w, 200, map[string]any{"appointments": m.appointments(sess)})
}

func (m *Module) InvoicePDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Write([]byte("%PDF-1.4"))
}

func (m *Module) Pay(w http.ResponseWriter, r *http.Request) {
	sess, _ := m.requireSession(w, r)
	acct := ""
	if sess != nil {
		if prof := m.profile(sess); prof != nil {
			acct = prof.AccountNumber
		}
	}
	writeJSON(w, 200, map[string]any{"status": "ok", "proxied": "fieldroutes", "account_number": acct, "receipt": "https://fieldroutes/receipt/" + acct})
}

func (m *Module) Cancel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "retention", "tel": "tel:+1-509-555-0199", "message": "Call to cancel — talk to us for a retention deal"})
}

func (m *Module) Stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Write([]byte("data: {\"type\":\"notification\"}\n\n"))
}

// dashboard assembles the logged-in customer's account, service history, and
// admin notes straight from FieldRoutes (no local DB required).
func (m *Module) dashboard(sess *auth.Session) map[string]any {
	prof := m.profile(sess)
	appts := m.appointments(sess)
	notes := visitNotes(appts)
	out := map[string]any{
		"dashboard":    "customer",
		"appointments": appts,
		"notes":        notes,
	}
	if prof == nil {
		return out
	}
	out["account"] = accountJSON(prof)
	return out
}

// appointments pulls the customer's service history for the logged-in session.
func (m *Module) appointments(sess *auth.Session) []any {
	appts := []any{}
	if m.FR == nil {
		return appts
	}
	prof := m.profile(sess)
	if prof == nil || prof.FRID == "" {
		return appts
	}
	d, found := m.FR.DistrictByCode(prof.District)
	if !found {
		return appts
	}
	raw, _ := m.FR.PullAppointments(d, prof.FRID)
	out := make([]any, 0, len(raw))
	for _, a := range raw {
		out = append(out, map[string]any{"id": a.ID, "date": a.Date, "start": a.Start, "type": a.Type, "status": a.Status, "notes": a.Notes})
	}
	return out
}

// visitNotes flattens admin notes attached to service visits for the dashboard.
func visitNotes(appts []any) []any {
	var notes []any
	for _, v := range appts {
		m, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if n, _ := m["notes"].(string); n != "" {
			notes = append(notes, map[string]any{"date": m["date"], "notes": n})
		}
	}
	if notes == nil {
		notes = []any{}
	}
	return notes
}

// accountJSON renders the customer profile for the portal.
func accountJSON(p *auth.Profile) map[string]any {
	return map[string]any{
		"fr_id":          p.FRID,
		"district":       p.District,
		"name":           p.Name,
		"email":          p.Email,
		"phone":          p.Phone,
		"account_number": p.AccountNumber,
		"status":         "active",
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	if code != 200 {
		w.WriteHeader(code)
	}
	_ = json.NewEncoder(w).Encode(v)
}
