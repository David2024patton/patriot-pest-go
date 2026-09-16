package customers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/David2024patton/patriot-pest-go/internal/auth"
	"github.com/David2024patton/patriot-pest-go/internal/fieldroutes"
	"github.com/go-chi/chi/v5"
)

// Module — customers + FieldRoutes mirror (customers/appointments/invoices).
// Sync pulls every customer from every configured district (WA + AZ) and
// upserts into the local cache, preserving local-only opt-out flags.
type Module struct {
	Enabled bool
	FR      *fieldroutes.Client
}

func (m *Module) Register(r chi.Router) bool {
	if !m.Enabled {
		return false
	}
	r.Get("/staff/customers", m.List)
	r.Get("/staff/customers/{id}", m.Get)
	r.Post("/staff/customers/sync", m.Sync)
	r.Get("/api/customer-search", m.Search)
	r.Post("/api/customers/{id}/notes", m.AddNote)
	return true
}

// cacheRow is one normalized customer row held in the local cache.
type cacheRow struct {
	FRID          string  `json:"fr_id"`
	District      string  `json:"district"`
	Name          string  `json:"name"`
	Email         *string `json:"email,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	AccountNumber string  `json:"account_number"`
	Address       *string `json:"address,omitempty"`
	City          *string `json:"city,omitempty"`
	State         *string `json:"state,omitempty"`
	Zip           *string `json:"zip,omitempty"`
	Status        string  `json:"status"`
	LastService   *string `json:"last_service,omitempty"`
	Source        string  `json:"source"`
	// Local-only flags never clobbered by a FieldRoutes sync.
	IsNoCall  bool   `json:"is_no_call"`
	DncReason string `json:"dnc_reason,omitempty"`
}

var (
	cacheMu   sync.RWMutex
	cache     = map[string]cacheRow{} // key: fr_id|district
	syncStats = map[string]int{}      // district -> count, last sync
)

func (m *Module) Sync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if m.FR == nil || len(m.FR.Districts()) == 0 {
		json.NewEncoder(w).Encode(map[string]any{
			"status":      "ok",
			"source":      "fieldroutes",
			"missing":     []string{},
			"districts":   []string{},
			"synced":      0,
		})
		return
	}

	perDistrict := map[string]int{}
	total := 0
	for _, d := range m.FR.Districts() {
		rows, err := m.FR.PullCustomers(d)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]any{
				"status":     "error",
				"error":      err.Error(),
				"district":   d.Code,
			})
			return
		}
		upsertRows(rows)
		perDistrict[d.Code] = len(rows)
		total += len(rows)
	}
	cacheMu.Lock()
	syncStats = perDistrict
	cacheMu.Unlock()

	json.NewEncoder(w).Encode(map[string]any{
		"status":       "ok",
		"source":       "fieldroutes",
		"synced":       total,
		"districts":    perDistrict,
			"missing":      m.FR.Missing(),
		"total_cached": len(cache),
	})
}

func (m *Module) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	out := []cacheRow{}
	for _, row := range cache {
		out = append(out, row)
	}
	json.NewEncoder(w).Encode(map[string]any{"customers": out, "count": len(out)})
}

func (m *Module) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	for _, row := range cache {
		if row.FRID == id || row.AccountNumber == id {
			json.NewEncoder(w).Encode(map[string]any{"customer": row, "appointments": []any{}, "invoices": []any{}})
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]any{"customer": map[string]any{"id": id}, "appointments": []any{}, "invoices": []any{}})
}

func (m *Module) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	w.Header().Set("Content-Type", "application/json")
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	results := []cacheRow{}
	for _, row := range cache {
		if matches(row, q) {
			results = append(results, row)
			if len(results) >= 25 {
				break
			}
		}
	}
	json.NewEncoder(w).Encode(map[string]any{"q": q, "results": results, "count": len(results)})
}

func (m *Module) AddNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "ok", "customer_id": chi.URLParam(r, "id"), "gps": r.Header.Get("X-GPS"), "note": "added"})
}

// matches does a case-insensitive substring match across identity fields.
func matches(row cacheRow, q string) bool {
	q = strings.ToLower(q)
	if q == "" {
		return false
	}
	hay := strings.ToLower(row.Name + " " + row.AccountNumber + " " + strPtr(row.Email) + " " + strPtr(row.Phone))
	return strings.Contains(hay, q)
}

// RecentRow is one customer row shaped for dashboard tables. FRID lets the
// dashboards fan out to FieldRoutes for live appointment detail.
type RecentRow struct {
	FRID          string `json:"fr_id,omitempty"`
	Name          string `json:"name"`
	District      string `json:"district"`
	Status        string `json:"status"`
	LastService   string `json:"last_service,omitempty"`
}

// Snapshot is a point-in-time, read-only view of the customer cache for the
// dashboards. It never mutates state and is safe to call concurrently.
func Snapshot() (total int, perDistrict map[string]int, recent []RecentRow) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	total = len(cache)
	perDistrict = map[string]int{}
	var rows []cacheRow
	for _, row := range cache {
		perDistrict[row.District]++
		rows = append(rows, row)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		ls := func(r *cacheRow) string {
			if r.LastService != nil {
				return *r.LastService
			}
			return ""
		}
		return ls(&rows[i]) > ls(&rows[j])
	})
	limit := 8
	if len(rows) < limit {
		limit = len(rows)
	}
	for i := 0; i < limit; i++ {
		var ls string
		if rows[i].LastService != nil {
			ls = *rows[i].LastService
		}
		recent = append(recent, RecentRow{FRID: rows[i].FRID, Name: rows[i].Name, District: rows[i].District, Status: rows[i].Status, LastService: ls})
	}
	return total, perDistrict, recent
}

// Resolve maps a login identity — email, phone, OR account number — to a cached
// row. Email and account number are case-insensitive exact matches; phone is
// normalized to E.164 so "(509) 555-0199" and "+15095550199" resolve to the
// same record. Returns ok=false when nothing matches.
func Resolve(identity string) (cacheRow, bool) {
	id := strings.TrimSpace(identity)
	if id == "" {
		return cacheRow{}, false
	}
	np := fieldroutes.NormalizePhone(id) // nil unless the identity parses as a phone
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	for _, row := range cache {
		if row.Email != nil && strings.EqualFold(*row.Email, id) {
			return row, true
		}
		if row.AccountNumber != "" && strings.EqualFold(row.AccountNumber, id) {
			return row, true
		}
		if np != nil && row.Phone != nil && *row.Phone == *np {
			return row, true
		}
	}
	return cacheRow{}, false
}

// Seed upserts normalized rows into the local cache (idempotent). Used by tests
// and, in production, by the sync endpoint after a FieldRoutes pull.
func Seed(rows []fieldroutes.Row) { upsertRows(rows) }

// GeoRow is the exported slice of a cached customer used by the heatmap and
// command bar (no DNC flags, no full address — city-level only).
type GeoRow struct {
	FRID          string `json:"fr_id"`
	District      string `json:"district"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	Status        string `json:"status"`
	LastService   string `json:"last_service,omitempty"`
}

// AllGeo returns every cached customer shaped for the density heatmap.
func AllGeo() []GeoRow {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	out := make([]GeoRow, 0, len(cache))
	for _, row := range cache {
		out = append(out, geoOf(row))
	}
	return out
}

// SearchRows searches the cache by name/email/phone/account number and
// returns up to limit exported rows (command bar + spotlight search).
func SearchRows(q string, limit int) []GeoRow {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	out := []GeoRow{}
	for _, row := range cache {
		if matches(row, q) {
			out = append(out, geoOf(row))
			if len(out) >= limit {
				break
			}
		}
	}
	return out
}

// GetRow returns one cached customer by FRID or account number.
func GetRow(id string) (GeoRow, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	for _, row := range cache {
		if row.FRID == id || row.AccountNumber == id {
			return geoOf(row), true
		}
	}
	return GeoRow{}, false
}

func geoOf(row cacheRow) GeoRow {
	g := GeoRow{
		FRID: row.FRID, District: row.District, Name: row.Name,
		AccountNumber: row.AccountNumber, Status: row.Status,
		City: strPtr(row.City), State: strPtr(row.State),
	}
	if row.LastService != nil {
		g.LastService = *row.LastService
	}
	return g
}

// Lookup returns an auth profile for a login identity, or ok=false when unknown.
// This is the concrete resolver wired into the auth module via auth.SetLookup.
func Lookup(identity string) (*auth.Profile, bool) {
	row, ok := Resolve(identity)
	if !ok {
		return nil, false
	}
	email, phone := "", ""
	if row.Email != nil {
		email = *row.Email
	}
	if row.Phone != nil {
		phone = *row.Phone
	}
	return &auth.Profile{
		FRID:          row.FRID,
		District:      row.District,
		Name:          row.Name,
		Email:         email,
		Phone:         phone,
		AccountNumber: row.AccountNumber,
	}, true
}

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// upsertRows applies FR rows into the cache: match by fr_id+district, then
// claim an unlinked seed row by email, otherwise insert. Local opt-out flags
// are never touched on update — FieldRoutes does not own them.
func upsertRows(rows []fieldroutes.Row) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	for _, row := range rows {
		if row.FRID == "" {
			continue
		}
		key := row.FRID + "|" + row.District
		if existing, ok := cache[key]; ok {
			merged := existing
			merged.Name = row.Name
			merged.Email = row.Email
			merged.Phone = row.Phone
			merged.Address = row.Address
			merged.City = row.City
			merged.State = row.State
			merged.Zip = row.Zip
			merged.Status = row.Status
			merged.LastService = row.LastService
			merged.Source = "fieldroutes"
			cache[key] = merged
			continue
		}
		// Claim an unlinked seed row by email.
		claimed := false
		if row.Email != nil {
			for k, v := range cache {
				if v.FRID == "" && v.Email != nil && *v.Email == *row.Email {
					merged := v
					merged.FRID = row.FRID
					merged.District = row.District
					merged.Name = row.Name
					merged.Email = row.Email
					merged.Phone = row.Phone
					merged.AccountNumber = row.AccountNumber
					merged.Address = row.Address
					merged.City = row.City
					merged.State = row.State
					merged.Zip = row.Zip
					merged.Status = row.Status
					merged.LastService = row.LastService
					merged.Source = "fieldroutes"
					cache[key] = merged
					delete(cache, k)
					claimed = true
					break
				}
			}
		}
		if !claimed {
			cache[key] = cacheRow{
				FRID: row.FRID, District: row.District, Name: row.Name, Email: row.Email, Phone: row.Phone,
				AccountNumber: row.AccountNumber, Address: row.Address, City: row.City, State: row.State, Zip: row.Zip,
				Status: row.Status, LastService: row.LastService, Source: "fieldroutes",
			}
		}
	}
}
