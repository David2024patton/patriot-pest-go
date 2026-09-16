// Package staffdash — the staff-facing operations dashboard at /staff-dashboard.
// It is cache-driven for speed (customer counts + recent customers come straight
// from the customers module's in-memory mirror) and fans out to FieldRoutes for a
// live "upcoming appointments" strip across the most recent customers. Every FR
// call degrades gracefully so an offline district never blanks the page.
package staffdash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/David2024patton/patriot-pest-go/internal/auth"
	"github.com/David2024patton/patriot-pest-go/internal/fieldroutes"
	"github.com/David2024patton/patriot-pest-go/internal/modules/customers"
	"github.com/David2024patton/patriot-pest-go/internal/view"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	Enabled bool
	FR      *fieldroutes.Client
}

func (m *Module) Register(r chi.Router) bool {
	if !m.Enabled {
		return false
	}
	r.Get("/staff-dashboard", m.Page)
	r.Get("/api/staff/overview", m.JSONOverview)
	return true
}

// requireStaff resolves the session cookie to a live, non-customer session.
func (m *Module) requireStaff(w http.ResponseWriter, r *http.Request) (*auth.Session, bool) {
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		return nil, false
	}
	s, ok := auth.GetSession(c.Value)
	if !ok || s.Role == "customer" {
		return nil, false
	}
	return s, true
}

// Page renders the staff operations dashboard behind a staff session.
func (m *Module) Page(w http.ResponseWriter, r *http.Request) {
	_, ok := m.requireStaff(w, r)
	if !ok {
		view.PageStatus(w, r, 403, "dash-staff", "Staff Dashboard | Patriot Pest Control", "", "", map[string]any{
			"AppUI": true, "UserType": "staff",
			"Flash": "Please sign in with a staff account.",
		})
		return
	}
	view.Page(w, r, "dash-staff", "Staff Dashboard | Patriot Pest Control", "", "", m.overview())
}

// JSONOverview is the machine-readable overview at /api/staff/overview.
func (m *Module) JSONOverview(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := m.overview()
	_ = json.NewEncoder(w).Encode(map[string]any{
		"stats":        d["Stats"],
		"appointments": d["Appointments"],
	})
}

// overview assembles the stat cards + live appointment strip for the page.
func (m *Module) overview() map[string]any {
	total, perDistrict, recent := customers.Snapshot()

	stats := []map[string]string{{"v": fmtInt(total), "k": "Customers"}}
	for _, code := range m.districtCodes() {
		stats = append(stats, map[string]string{"v": fmtInt(perDistrict[code]), "k": code + " District"})
	}

	recentRows := make([]map[string]string, 0, len(recent))
	for _, row := range recent {
		recentRows = append(recentRows, map[string]string{
			"name": row.Name, "district": row.District, "status": row.Status, "date": row.LastService,
		})
	}

	return map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
		"Stats": stats,
		"Appointments": m.liveAppointments(recent),
		"Recent": recentRows,
	}
}

// liveAppointments fans out to FieldRoutes for up to 3 recent customers and
// merges their service history into a single most-recent-first list (capped).
func (m *Module) liveAppointments(recent []customers.RecentRow) []map[string]string {
	if m.FR == nil || len(m.FR.Districts()) == 0 {
		return []map[string]string{}
	}
	out := []map[string]string{}
	seen := map[string]bool{}
	var mu sync.Mutex

	limit := 3
	if len(recent) < limit {
		limit = len(recent)
	}
	var wg sync.WaitGroup
	for i := 0; i < limit; i++ {
		row := recent[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			d, found := m.FR.DistrictByCode(row.District)
			if !found {
				return
			}
			appts, err := m.FR.PullAppointments(d, row.FRID)
			if err != nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			for _, a := range appts {
				key := a.ID + "|" + a.Date
				if seen[key] {
					continue
				}
				seen[key] = true
				out = append(out, map[string]string{
					"date": a.Date, "start": a.Start, "type": a.Type,
					"district": row.District, "status": a.Status,
				})
			}
		}()
	}
	wg.Wait()

	sortRows(out)
	if len(out) > 25 {
		out = out[:25]
	}
	return out
}

func (m *Module) districtCodes() []string {
	codes := []string{}
	if m.FR == nil {
		return codes
	}
	for _, d := range m.FR.Districts() {
		codes = append(codes, d.Code)
	}
	return codes
}

func fmtInt(n int) string { return fmt.Sprintf("%d", n) }

// sortRows orders appointment rows by date descending (best-effort).
func sortRows(rows []map[string]string) {
	for i := 1; i < len(rows); i++ {
		for j := i; j > 0 && rows[j]["date"] > rows[j-1]["date"]; j-- {
			rows[j], rows[j-1] = rows[j-1], rows[j]
		}
	}
}
