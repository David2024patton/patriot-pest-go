package legacy

// Phase 2 handlers — the Salesforce-killer dashboard layer:
//   - Global Command Bar (Ctrl+K) backed by /api/command
//   - Customer 360 interactive timeline (/staff/customers/{id}/timeline)
//   - USA customer density heatmap (/admin/heatmap + /api/admin/heatmap.json)
//   - SSE live feed HUD (/api/staff/events) with a process-wide event hub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/David2024patton/patriot-pest-go/internal/data"
	"github.com/David2024patton/patriot-pest-go/internal/events"
	"github.com/David2024patton/patriot-pest-go/internal/modules/customers"
	"github.com/David2024patton/patriot-pest-go/internal/view"
	"github.com/go-chi/chi/v5"
)

// ---- SSE event hub (live feed HUD) ----

// PublishEvent fans an event out to every connected SSE client (delegates to
// the shared events hub so other modules can emit too). Safe from anywhere.
func PublishEvent(typ, text string) { events.Publish(typ, text) }

// StaffEvents streams live HUD events (calls, messages, callbacks) over SSE.
func (m *Module) StaffEvents(w http.ResponseWriter, r *http.Request) {
	fl, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", 500)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	ch := events.Subscribe()
	defer events.Unsubscribe(ch)
	fmt.Fprint(w, "retry: 3000\n\n")
	fl.Flush()
	ctx := r.Context()
	keep := time.NewTicker(25 * time.Second)
	defer keep.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-keep.C:
			fmt.Fprint(w, ": keepalive\n\n")
			fl.Flush()
		case ev := <-ch:
			b, _ := json.Marshal(ev)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, b)
			fl.Flush()
		}
	}
}

// ---- Global Command Bar (Ctrl+K) search API ----

// CommandSearch merges customers, pests, posts and nav pages into one
// grouped result set for the spotlight modal.
func (m *Module) CommandSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	out := map[string]any{"q": q, "customers": []any{}, "pests": []any{}, "posts": []any{}, "pages": []any{}}
	if q == "" {
		writeJSONStatus(w, 200, out)
		return
	}
	ql := strings.ToLower(q)

	custs := []map[string]string{}
	for _, c := range customers.SearchRows(q, 8) {
		custs = append(custs, map[string]string{
			"title": c.Name, "sub": c.AccountNumber + " · " + strings.ToUpper(c.District),
			"url": "/staff/customers/" + c.FRID + "/timeline",
		})
	}
	out["customers"] = custs

	pests := []map[string]string{}
	for _, p := range data.AllPests() {
		if strings.Contains(strings.ToLower(p.Name), ql) || strings.Contains(strings.ToLower(p.Category), ql) {
			pests = append(pests, map[string]string{"title": p.Name, "sub": p.ScientificName, "url": "/pest/" + p.Slug})
			if len(pests) >= 5 {
				break
			}
		}
	}
	out["pests"] = pests

	posts := []map[string]string{}
	for _, p := range data.AllPosts() {
		if strings.Contains(strings.ToLower(p.Title), ql) {
			posts = append(posts, map[string]string{"title": p.Title, "sub": "Blog", "url": "/blogs/" + p.Slug})
			if len(posts) >= 5 {
				break
			}
		}
	}
	out["posts"] = posts

	pages := []map[string]string{}
	for _, pg := range []map[string]string{
		{"title": "Home", "url": "/"}, {"title": "Services", "url": "/services"},
		{"title": "Prices", "url": "/prices"}, {"title": "Service Areas", "url": "/service-areas"},
		{"title": "Blog", "url": "/blogs"}, {"title": "FAQs", "url": "/faqs"},
		{"title": "Contact", "url": "/contact"}, {"title": "Admin Console", "url": "/admin"},
		{"title": "Staff Dashboard", "url": "/staff-dashboard"}, {"title": "FieldRoutes Health", "url": "/admin/fieldroutes"},
		{"title": "Customer Heatmap", "url": "/admin/heatmap"}, {"title": "Kanban Board", "url": "/admin/board"},
	} {
		if strings.Contains(strings.ToLower(pg["title"]), ql) {
			pages = append(pages, pg)
		}
	}
	out["pages"] = pages
	writeJSONStatus(w, 200, out)
}

// ---- Customer 360 timeline ----

// CustomerTimeline renders the unified 360 feed for one customer: profile,
// FieldRoutes appointments, service history and portal messages.
func (m *Module) CustomerTimeline(w http.ResponseWriter, r *http.Request) {
	s := m.sess(r)
	if s == nil || s.Role != "staff" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	id := chi.URLParam(r, "id")
	row, found := customers.GetRow(id)

	events := []map[string]any{}
	if found {
		if row.LastService != "" {
			events = append(events, map[string]any{
				"icon": "🛠", "kind": "Service", "at": row.LastService,
				"text": "Last completed service on record.",
			})
		}
		if m.FR != nil {
			if d, ok := m.FR.DistrictByCode(row.District); ok {
				if appts, err := m.FR.PullAppointments(d, row.FRID); err == nil {
					for _, a := range appts {
						events = append(events, map[string]any{
							"icon": "📅", "kind": "Appointment", "at": a.Date + " " + a.Start,
							"text": a.Type + " — " + a.Status,
						})
					}
				}
			}
		}
		msgMu.Lock()
		for identity, list := range msgs {
			if p, ok := customers.Lookup(identity); ok && p.FRID == row.FRID {
				for _, mm := range list {
					events = append(events, map[string]any{
						"icon": "💬", "kind": "Message", "at": mm.At, "text": mm.Text,
					})
				}
			}
		}
		msgMu.Unlock()
	}

	view.Page(w, r, "customer-360", "Customer 360 | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff",
		"Found": found, "Customer": row, "Events": events,
	})
}

// ---- USA customer density heatmap ----

// cityGeo is the static geocoder: known metro centroids + state fallbacks.
var cityGeo = map[string][2]float64{
	"spokane":      {47.6588, -117.4260},
	"spokane valley": {47.6732, -117.2344},
	"seattle":      {47.6062, -122.3321},
	"tacoma":       {47.2529, -122.4443},
	"vancouver":    {45.6387, -122.6615},
	"bellevue":     {47.6101, -122.2015},
	"everett":      {47.9790, -122.2021},
	"yakima":       {46.6021, -120.5059},
	"kennewick":    {46.2112, -119.1372},
	"boise":        {43.6150, -116.2023},
	"meridian":     {43.6121, -116.3915},
	"nampa":        {43.5407, -116.5635},
	"idaho falls":  {43.4917, -112.0339},
	"pocatello":    {42.8713, -112.4455},
	"portland":     {45.5152, -122.6784},
	"gresham":      {45.5001, -122.4302},
	"eugene":       {44.0521, -123.0868},
	"salem":        {44.9429, -123.0351},
	"phoenix":      {33.4484, -112.0740},
	"mesa":         {33.4152, -111.8315},
	"scottsdale":   {33.4942, -111.9261},
	"tempe":        {33.4255, -111.9400},
	"glendale":     {33.5387, -112.1860},
	"chandler":     {33.3062, -111.8413},
	"gilbert":      {33.3528, -111.7890},
	"peoria":       {33.5806, -112.2374},
	"surprise":     {33.6292, -112.3679},
	"avondale":     {33.4356, -112.3496},
	"tucson":       {32.2226, -110.9747},
	"flagstaff":    {35.1983, -111.6513},
	"prescott":     {34.5400, -112.4685},
	"yuma":         {32.6927, -114.6277},
}

var stateGeo = map[string][2]float64{
	"WA": {47.7511, -120.7401}, "ID": {44.0682, -114.7420},
	"OR": {43.8041, -120.5542}, "AZ": {34.0489, -111.0937},
}

// HeatmapJSON aggregates customers per city for the Leaflet layer.
func (m *Module) HeatmapJSON(w http.ResponseWriter, r *http.Request) {
	type spot struct {
		City  string    `json:"city"`
		State string    `json:"state"`
		Lat   float64   `json:"lat"`
		Lng   float64   `json:"lng"`
		Count int       `json:"count"`
		Names []string  `json:"names"`
	}
	agg := map[string]*spot{}
	for _, c := range customers.AllGeo() {
		city := strings.ToLower(strings.TrimSpace(c.City))
		st := strings.ToUpper(strings.TrimSpace(c.State))
		if st == "" {
			st = strings.ToUpper(c.District)
		}
		geo, ok := cityGeo[city]
		if !ok {
			geo, ok = stateGeo[st]
		}
		if !ok {
			continue
		}
		key := city + "|" + st
		sp, ok := agg[key]
		if !ok {
			label := strings.Title(city)
			if city == "" {
				label = st + " (statewide)"
			}
			sp = &spot{City: label, State: st, Lat: geo[0], Lng: geo[1]}
			agg[key] = sp
		}
		sp.Count++
		if len(sp.Names) < 8 {
			sp.Names = append(sp.Names, c.Name)
		}
	}
	out := make([]*spot, 0, len(agg))
	for _, sp := range agg {
		out = append(out, sp)
	}
	writeJSONStatus(w, 200, map[string]any{"spots": out, "total": len(customers.AllGeo())})
}

// Heatmap renders the Leaflet density map page.
func (m *Module) Heatmap(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "dash-heatmap"); !ok {
		return
	}
	view.Page(w, r, "dash-heatmap", "Customer Heatmap | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
	})
}
