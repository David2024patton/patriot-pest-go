package legacy

import (
	"net/http"

	"github.com/David2024patton/patriot-pest-go/internal/auth"
	"github.com/David2024patton/patriot-pest-go/internal/fieldroutes"
	"github.com/go-chi/chi/v5"
)

// Module — final-phase routes from the 103-route spec. Phase 1 of the master
// plan replaced every JSON stub with a real handler (pages_phase1.go +
// phase1.go). Routes owned by other modules stay absent here: chi's
// InsertRoute lets the LAST registration win, and legacy registers last.
type Module struct {
	Enabled          bool
	SuperuserEnabled bool // SUPERUSER_ENABLED gates the /su surface
	FR               *fieldroutes.Client
	Mailer           *auth.Mailer
	DBPath           string // SQLite catalog path for CMS writes
}

func (m *Module) Register(r chi.Router) bool {
	if !m.Enabled {
		return false
	}
	// FR-002 cost — standalone valuation receipt (real page).
	r.Get("/cost", m.Cost)
	r.Get("/cost/data/pricing.json", m.CostData)
	r.Get("/data/pricing.json", m.CostData) // cost.js resolves ./data/ from /cost

	// FR-004 auth legacy aliases -> real surfaces.
	r.Get("/customer-auth", redirect("/login"))
	r.Get("/customer-verify", redirect("/login/verify"))
	r.Get("/staff", redirect("/login"))
	r.Get("/staff-verify", redirect("/login/verify"))
	r.Get("/staff-logout", redirect("/logout"))
	r.Get("/dashboard", m.DashboardAlias)

	// Super-user elevated surface (SUPERUSER_ENABLED-gated; 404 otherwise).
	r.Get("/su", m.SuLogin)
	r.Post("/su", m.SuRequest)
	r.Get("/su/verify", m.SuVerifyForm)
	r.Post("/su/verify", m.SuVerify)

	// FR-006 account (real profile page for staff or customer sessions).
	r.Get("/account", m.Account)

	// FR-007 admin CMS (SQLite-backed posts + settings + staff CRUD).
	r.Get("/admin/posts/new", m.AdminPostNew)
	r.Get("/admin/posts/{id}", m.AdminPostEdit)
	r.Post("/admin/posts/{id}", m.AdminPostSave)
	r.Post("/admin/posts", m.AdminPostSave)
	r.Get("/admin/settings", m.AdminSettingsPage)
	r.Post("/admin/settings", m.AdminSettingsSave)
	r.Get("/admin/staff/new", m.AdminStaffNew)
	r.Post("/admin/staff", m.AdminStaffCreate)
	r.Get("/admin/staff/{id}", m.AdminStaffEdit)

	// FR-008 api-keys lifecycle.
	r.Get("/admin/api-keys/audit", m.ApiKeyAudit)
	r.Post("/admin/api-keys/{id}/revoke", m.ApiKeyRevoke)
	r.Post("/admin/api-keys/{id}/rotate", m.ApiKeyRotate)
	r.Post("/admin/api-keys/{id}/scopes", m.ApiKeyScopes)

	// FR-010 retention settings save.
	r.Post("/admin/retention/settings", m.RetentionSettingsSave)

	// FR-026-028 portal extras.
	r.Get("/customer-portal", m.CustomerPortal)
	r.Get("/customer/invoices/{id}/download", m.InvoiceDownload)
	r.Get("/customer/messages", m.CustomerMessages)
	r.Post("/api/customer/messages", m.CustomerMessagePost)
	r.Get("/admin/fieldroutes", m.AdminFieldRoutes)

	// FR-029 AI copilot search.
	r.Get("/tech/ask", m.TechAsk)

	// Phase 2 — Salesforce-killer dashboards.
	r.Get("/api/command", m.CommandSearch)
	r.Get("/api/staff/events", m.StaffEvents)
	r.Get("/staff/customers/{id}/timeline", m.CustomerTimeline)
	r.Get("/admin/heatmap", m.Heatmap)
	r.Get("/api/admin/heatmap.json", m.HeatmapJSON)

	return true
}

var _ = http.StatusOK
