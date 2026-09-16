package legacy

// Phase 1 handlers — real implementations replacing the JSON stubs, per
// PLANS/MASTERPLAN_GO_REWRITE.md and master_plan.md "STILL TO DO Phase 1".
// Conventions mirror the admin module: session gate -> view.Page with
// AppUI/UserType/Csrf/Flash; POSTs verified with view.VerifyCSRF.

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/David2024patton/patriot-pest-go/internal/auth"
	"github.com/David2024patton/patriot-pest-go/internal/data"
	"github.com/David2024patton/patriot-pest-go/internal/modules/api"
	"github.com/David2024patton/patriot-pest-go/internal/modules/customers"
	"github.com/David2024patton/patriot-pest-go/internal/modules/inbox"
	"github.com/David2024patton/patriot-pest-go/internal/rbac"
	"github.com/David2024patton/patriot-pest-go/internal/view"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

// ---- cost (/cost) — standalone valuation receipt ----

func (m *Module) Cost(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(view.CostPage))
}

func (m *Module) CostData(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(view.CostPricingJSON))
}

// ---- legacy auth aliases -> real surfaces ----

func redirect(to string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusFound)
	}
}

// DashboardAlias sends a session holder to the right dashboard, else /login.
func (m *Module) DashboardAlias(w http.ResponseWriter, r *http.Request) {
	if s := m.sess(r); s != nil {
		if s.Role == "staff" {
			http.Redirect(w, r, "/staff-dashboard", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/customer-dashboard", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// ---- session helpers (same cookie contract as portal/admin) ----

func (m *Module) sess(r *http.Request) *auth.Session {
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		return nil
	}
	s, ok := auth.GetSession(c.Value)
	if !ok {
		return nil
	}
	return s
}

func (m *Module) isSuper(s *auth.Session) bool {
	return rbac.IsSuperUser(s.Identity) || s.RoleLabel == "super-user"
}

// gateAdmin renders a 403 app shell unless the session is admin-or-above.
func (m *Module) gateAdmin(w http.ResponseWriter, r *http.Request, page string) (*auth.Session, bool) {
	s := m.sess(r)
	if s == nil || s.Role != "staff" || !(m.isSuper(s) || s.RoleLabel == "admin") {
		view.PageStatus(w, r, 403, page, "Admin | Patriot Pest Control", "", "", map[string]any{
			"AppUI": true, "UserType": "staff", "Flash": "Needs an admin or super-user session.",
		})
		return nil, false
	}
	return s, true
}

// ---- /account — self-service profile for the signed-in user ----

func (m *Module) Account(w http.ResponseWriter, r *http.Request) {
	s := m.sess(r)
	if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	d := map[string]any{"AppUI": true, "Flash": ""}
	if s.Role == "staff" {
		d["UserType"] = "staff"
		d["Type"] = "staff"
		d["IsAdmin"] = m.isSuper(s) || s.RoleLabel == "admin"
		rec := map[string]string{}
		if st := auth.StaffByEmail(s.Identity); st != nil {
			rec = map[string]string{"name": st.Name, "email": st.Email, "role": st.Role}
			d["Name"] = st.Name
			d["RoleLabel"] = strings.Title(strings.ReplaceAll(st.Role, "-", " "))
		} else {
			d["Name"] = s.Identity
			d["RoleLabel"] = "Staff"
		}
		d["Record"] = rec
	} else {
		d["UserType"] = "customer"
		d["Type"] = "customer"
		d["Name"] = s.Identity
		d["RoleLabel"] = "Customer"
		if p, ok := customers.Lookup(s.Identity); ok {
			d["Name"] = p.Name
			d["Record"] = map[string]string{
				"account_number": p.AccountNumber, "name": p.Name, "email": p.Email,
				"phone": p.Phone, "district": strings.ToUpper(p.District), "status": "Active",
			}
		} else {
			d["Record"] = nil
		}
	}
	view.Page(w, r, "account", "My Account | Patriot Pest Control", "", "", d)
}

// ---- /customer-portal + /customer/messages ----

func (m *Module) CustomerPortal(w http.ResponseWriter, r *http.Request) {
	view.Page(w, r, "customer-portal", "Customer Portal | Patriot Pest Control", "", "", map[string]any{"AppUI": true})
}

type custMsg struct {
	From string `json:"from"`
	Text string `json:"text"`
	At   string `json:"at"`
	Mine bool   `json:"mine"`
}

var (
	msgMu sync.Mutex
	msgs  = map[string][]custMsg{}
)

func (m *Module) CustomerMessages(w http.ResponseWriter, r *http.Request) {
	s := m.sess(r)
	if s == nil || s.Role != "customer" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	msgMu.Lock()
	list := append([]custMsg{}, msgs[s.Identity]...)
	msgMu.Unlock()
	rows := make([]map[string]any, 0, len(list))
	for _, mm := range list {
		rows = append(rows, map[string]any{"from": mm.From, "text": mm.Text, "at": mm.At, "mine": mm.Mine})
	}
	view.Page(w, r, "customer-messages", "Messages | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "customer",
		"Csrf": view.CSRFField(w, r), "Messages": rows,
	})
}

func (m *Module) CustomerMessagePost(w http.ResponseWriter, r *http.Request) {
	s := m.sess(r)
	if s == nil || s.Role != "customer" {
		writeJSONStatus(w, 401, map[string]any{"error": "not authenticated"})
		return
	}
	if !view.VerifyCSRF(r) {
		writeJSONStatus(w, 403, map[string]any{"error": "bad csrf token"})
		return
	}
	text := strings.TrimSpace(r.FormValue("message"))
	if text == "" {
		writeJSONStatus(w, 422, map[string]any{"error": "message required"})
		return
	}
	if len(text) > 2000 {
		text = text[:2000]
	}
	name := s.Identity
	if p, ok := customers.Lookup(s.Identity); ok && p.Name != "" {
		name = p.Name
	}
	msgMu.Lock()
	msgs[s.Identity] = append(msgs[s.Identity], custMsg{From: name, Text: text, At: time.Now().Format("Jan 2, 15:04"), Mine: true})
	msgMu.Unlock()
	inbox.PublishCustomerMessage(name, text) // feeds Unified Inbox + SSE HUD
	// Accept both form posts (redirect) and XHR posts (JSON).
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		writeJSONStatus(w, 200, map[string]any{"status": "ok"})
		return
	}
	http.Redirect(w, r, "/customer/messages", http.StatusSeeOther)
}

// ---- /customer/invoices/{id}/download — FieldRoutes proxy ----

// InvoiceDownload streams an invoice PDF via the portal's FieldRoutes proxy.
// Zero card data ever touches this server (PCI rule 4).
func (m *Module) InvoiceDownload(w http.ResponseWriter, r *http.Request) {
	s := m.sess(r)
	if s == nil || s.Role != "customer" {
		writeJSONStatus(w, 401, map[string]any{"error": "not authenticated"})
		return
	}
	id := chi.URLParam(r, "id")
	if m.FR == nil || len(m.FR.Districts()) == 0 {
		writeJSONStatus(w, 502, map[string]any{"error": "FieldRoutes not configured — invoice download unavailable"})
		return
	}
	// FieldRoutes exposes billing documents per customer; until the billing
	// endpoint contract is wired we surface a clear, actionable response.
	writeJSONStatus(w, 501, map[string]any{
		"error":   "invoice streaming pending FieldRoutes billing endpoint wiring",
		"invoice": id,
		"hint":    "view invoices in the customer dashboard or call the office",
	})
}

// ---- /tech/ask — field tech copilot over the local knowledge base ----

func (m *Module) TechAsk(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	hits := []map[string]string{}
	if q != "" {
		ql := strings.ToLower(q)
		for _, p := range data.AllPests() {
			if strings.Contains(strings.ToLower(p.Name), ql) ||
				strings.Contains(strings.ToLower(p.Description), ql) ||
				strings.Contains(strings.ToLower(p.Category), ql) {
				hits = append(hits, map[string]string{
					"title":   p.Name + " — " + p.ScientificName,
					"url":     "/pest/" + p.Slug,
					"snippet": clip(p.Description, 160),
				})
			}
		}
		for _, p := range data.AllPosts() {
			if strings.Contains(strings.ToLower(p.Title), ql) ||
				strings.Contains(strings.ToLower(p.Excerpt), ql) ||
				strings.Contains(strings.ToLower(p.BodyHTML), ql) {
				hits = append(hits, map[string]string{
					"title":   p.Title,
					"url":     "/blogs/" + p.Slug,
					"snippet": clip(p.Excerpt, 160),
				})
			}
		}
		if len(hits) > 12 {
			hits = hits[:12]
		}
	}
	view.Page(w, r, "tech-ask", "Field Tech Copilot | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "Q": q, "Hits": hits,
	})
}

func clip(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ---- /admin/fieldroutes — district sync health ----

func (m *Module) AdminFieldRoutes(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "dash-fieldroutes"); !ok {
		return
	}
	rows := []map[string]any{}
	if m.FR != nil {
		ds := m.FR.Districts()
		type res struct {
			idx int
			row map[string]any
		}
		out := make(chan res, len(ds))
		var wg sync.WaitGroup
		for i, d := range ds {
			wg.Add(1)
			go func(i int, code, base, key, token string) {
				defer wg.Done()
				row := map[string]any{
					"code": code, "base": base,
					"configured": key != "" && token != "",
					"healthy":    false, "latency": "—", "customers": 0,
				}
				if key != "" && token != "" && m.FR != nil {
					if dd, ok := m.FR.DistrictByCode(code); ok {
						t0 := time.Now()
						rowsPull, err := m.FR.PullCustomers(dd)
						if err == nil {
							row["healthy"] = true
							row["latency"] = time.Since(t0).Round(time.Millisecond).String()
							row["customers"] = len(rowsPull)
						}
					}
				}
				out <- res{i, row}
			}(i, d.Code, d.Base, d.Key, d.Token)
		}
		wg.Wait()
		close(out)
		ordered := make([]map[string]any, len(ds))
		for rr := range out {
			ordered[rr.idx] = rr.row
		}
		rows = ordered
	}
	view.Page(w, r, "dash-fieldroutes", "FieldRoutes Health | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true, "Districts": rows,
	})
}

// ---- superuser /su flow (SUPERUSER_ENABLED-gated) ----

const suPendingCookie = "su_pending"

func (m *Module) suEnabled(w http.ResponseWriter, r *http.Request) bool {
	if m.SuperuserEnabled {
		return true
	}
	http.NotFound(w, r)
	return false
}

func (m *Module) SuLogin(w http.ResponseWriter, r *http.Request) {
	if !m.suEnabled(w, r) {
		return
	}
	view.Page(w, r, "su-login", "Superuser Sign In | Patriot Pest Control", "", "", map[string]any{
		"Csrf": view.CSRFField(w, r),
	})
}

func (m *Module) SuRequest(w http.ResponseWriter, r *http.Request) {
	if !m.suEnabled(w, r) {
		return
	}
	if !view.VerifyCSRF(r) {
		view.PageStatus(w, r, 403, "su-login", "Superuser Sign In | Patriot Pest Control", "", "", map[string]any{
			"Csrf": view.CSRFField(w, r), "FlashError": "Session expired — try again.",
		})
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	// Enumeration defense: only the seeded super-user identity may receive a
	// code; everyone else sees the same "code sent" screen (no code issued).
	authorized := email != "" && (rbac.IsSuperUser(email) || isSuperStaff(email))
	sentTo := ""
	if authorized {
		code, err := auth.Issue(r.Context(), email, "super_login", 8, 300)
		if err == nil && m.Mailer != nil {
			m.Mailer.Send(email, "Patriot Pest — superuser sign-in code",
				auth.MailTemplate("Superuser sign-in code",
					"<p style='font-size:28px;font-weight:800;letter-spacing:.2em'>"+code+"</p><p>Expires in 5 minutes. Works once.</p>"))
		}
		setPending(r, w, email)
		sentTo = email
	} else if email != "" {
		setPending(r, w, email)
		sentTo = email
	}
	view.Page(w, r, "su-verify", "Verify Code | Patriot Pest Control", "", "", map[string]any{
		"Csrf": view.CSRFField(w, r), "SentTo": sentTo,
	})
}

func (m *Module) SuVerifyForm(w http.ResponseWriter, r *http.Request) {
	if !m.suEnabled(w, r) {
		return
	}
	view.Page(w, r, "su-verify", "Verify Code | Patriot Pest Control", "", "", map[string]any{
		"Csrf": view.CSRFField(w, r), "SentTo": getPending(r),
	})
}

func (m *Module) SuVerify(w http.ResponseWriter, r *http.Request) {
	if !m.suEnabled(w, r) {
		return
	}
	if !view.VerifyCSRF(r) {
		view.PageStatus(w, r, 403, "su-verify", "Verify Code | Patriot Pest Control", "", "", map[string]any{
			"Csrf": view.CSRFField(w, r), "FlashError": "Session expired — try again.", "SentTo": getPending(r),
		})
		return
	}
	email := getPending(r)
	code := strings.TrimSpace(r.FormValue("code"))
	if email == "" || code == "" {
		view.PageStatus(w, r, 400, "su-verify", "Verify Code | Patriot Pest Control", "", "", map[string]any{
			"Csrf": view.CSRFField(w, r), "FlashError": "Request a new code first.",
		})
		return
	}
	ok, _ := auth.Verify(r.Context(), email, "super_login", code, 3)
	if !ok || !isSuperStaff(email) && !rbac.IsSuperUser(email) {
		view.PageStatus(w, r, 401, "su-verify", "Verify Code | Patriot Pest Control", "", "", map[string]any{
			"Csrf": view.CSRFField(w, r), "FlashError": "That code didn't match. Check the latest email.", "SentTo": email,
		})
		return
	}
	s := auth.CreateSession(email, "staff", 7200)
	s.SetRoleLabel("super-user")
	m.setSessionCookie(w, s.ID)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func isSuperStaff(email string) bool {
	st := auth.StaffByEmail(email)
	return st != nil && st.Active && st.Role == "super-user"
}

func setPending(r *http.Request, w http.ResponseWriter, email string) {
	v := base64.RawURLEncoding.EncodeToString([]byte(email))
	http.SetCookie(w, &http.Cookie{
		Name: suPendingCookie, Value: v, Path: "/", HttpOnly: true,
		SameSite: http.SameSiteLaxMode, Secure: !m_Dev, MaxAge: 600,
	})
}

func getPending(r *http.Request) string {
	c, err := r.Cookie(suPendingCookie)
	if err != nil {
		return ""
	}
	b, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return ""
	}
	return string(b)
}

// m_Dev mirrors APP_ENV=local (set from main via SetDev) so su cookies drop
// Secure for http://localhost testing.
var m_Dev = false

func SetDev(dev bool) { m_Dev = dev }

func (m *Module) setSessionCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name: "session", Value: id, Path: "/", HttpOnly: true,
		SameSite: http.SameSiteLaxMode, Secure: !m_Dev, MaxAge: 7200,
	})
}

// ---- admin CMS: posts new/edit (SQLite-backed) ----

func (m *Module) openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", m.DBPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (m *Module) AdminPostNew(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-post-edit"); !ok {
		return
	}
	view.Page(w, r, "admin-post-edit", "New Post | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
		"Csrf": view.CSRFField(w, r), "IsNew": true,
		"Post": map[string]string{"published_at": time.Now().Format("2006-01-02")},
	})
}

func (m *Module) AdminPostEdit(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-post-edit"); !ok {
		return
	}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	post := m.postByID(id)
	if post == nil {
		http.NotFound(w, r)
		return
	}
	view.Page(w, r, "admin-post-edit", "Edit Post | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
		"Csrf": view.CSRFField(w, r), "IsNew": false, "Post": post,
	})
}

func (m *Module) AdminPostSave(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-post-edit"); !ok {
		return
	}
	if !view.VerifyCSRF(r) {
		http.Error(w, "bad csrf token", http.StatusForbidden)
		return
	}
	_ = r.ParseForm()
	post := map[string]string{
		"title":         strings.TrimSpace(r.FormValue("title")),
		"slug":          slugify(r.FormValue("slug")),
		"excerpt":       strings.TrimSpace(r.FormValue("excerpt")),
		"body_html":     r.FormValue("body_html"),
		"author":        strings.TrimSpace(r.FormValue("author")),
		"season":        strings.TrimSpace(r.FormValue("season")),
		"pest_category": strings.TrimSpace(r.FormValue("pest_category")),
		"published_at":  strings.TrimSpace(r.FormValue("published_at")),
	}
	if post["title"] == "" || post["slug"] == "" {
		view.PageStatus(w, r, 422, "admin-post-edit", "Edit Post | Patriot Pest Control", "", "", map[string]any{
			"AppUI": true, "UserType": "staff", "IsAdmin": true,
			"Csrf": view.CSRFField(w, r), "IsNew": false, "Post": post,
			"Flash": "Title and slug are required.",
		})
		return
	}
	if err := m.upsertPost(post); err != nil {
		view.PageStatus(w, r, 500, "admin-post-edit", "Edit Post | Patriot Pest Control", "", "", map[string]any{
			"AppUI": true, "UserType": "staff", "IsAdmin": true,
			"Csrf": view.CSRFField(w, r), "IsNew": false, "Post": post,
			"Flash": "Save failed: " + err.Error(),
		})
		return
	}
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (m *Module) postByID(id int) map[string]string {
	db, err := m.openDB()
	if err != nil {
		return nil
	}
	defer db.Close()
	row := db.QueryRow(`SELECT slug, title, COALESCE(excerpt,''), COALESCE(body_html,''),
		COALESCE(author,''), COALESCE(season,''), COALESCE(pest_category,''),
		COALESCE(substr(published_at,1,10),'') FROM posts WHERE id = ?`, id)
	var p map[string]string
	var slug, title, excerpt, body, author, season, cat, pub string
	if err := row.Scan(&slug, &title, &excerpt, &body, &author, &season, &cat, &pub); err != nil {
		return nil
	}
	p = map[string]string{
		"id": strconv.Itoa(id), "slug": slug, "title": title, "excerpt": excerpt,
		"body_html": body, "author": author, "season": season,
		"pest_category": cat, "published_at": pub,
	}
	return p
}

func (m *Module) upsertPost(p map[string]string) error {
	db, err := m.openDB()
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT, slug TEXT NOT NULL UNIQUE, title TEXT NOT NULL,
		excerpt TEXT, body_html TEXT, pest_photo_id INTEGER, season TEXT, pest_category TEXT,
		status TEXT NOT NULL DEFAULT 'published', author TEXT, views INTEGER NOT NULL DEFAULT 0,
		published_at TEXT, date_modified TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now')), updated_at TEXT NOT NULL DEFAULT (datetime('now')))`); err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (slug, title, excerpt, body_html, author, season, pest_category, status, published_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'published', ?, datetime('now'))
		ON CONFLICT(slug) DO UPDATE SET title=excluded.title, excerpt=excluded.excerpt,
		body_html=excluded.body_html, author=excluded.author, season=excluded.season,
		pest_category=excluded.pest_category, published_at=excluded.published_at, updated_at=datetime('now')`,
		p["slug"], p["title"], p["excerpt"], p["body_html"], p["author"], p["season"], p["pest_category"], p["published_at"])
	return err
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, rn := range s {
		switch {
		case rn >= 'a' && rn <= 'z' || rn >= '0' && rn <= '9':
			b.WriteRune(rn)
		case rn == '-' || rn == '_' || rn == ' ':
			if b.Len() > 0 && !strings.HasSuffix(b.String(), "-") {
				b.WriteByte('-')
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// ---- admin staff new/edit/create ----

func (m *Module) AdminStaffNew(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-staff-edit"); !ok {
		return
	}
	view.Page(w, r, "admin-staff-edit", "Add Staff | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
		"Csrf": view.CSRFField(w, r), "IsNew": true,
		"StaffRow": map[string]string{"role": "staff"},
	})
}

func (m *Module) AdminStaffEdit(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-staff-edit"); !ok {
		return
	}
	email := chi.URLParam(r, "id")
	st := auth.StaffByEmail(email)
	row := map[string]string{"email": email}
	if st != nil {
		row = map[string]string{"email": st.Email, "name": st.Name, "role": st.Role, "title": st.Title}
	}
	view.Page(w, r, "admin-staff-edit", "Edit Staff | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
		"Csrf": view.CSRFField(w, r), "IsNew": st == nil, "StaffRow": row,
	})
}

func (m *Module) AdminStaffCreate(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-staff-edit"); !ok {
		return
	}
	if !view.VerifyCSRF(r) {
		http.Error(w, "bad csrf token", http.StatusForbidden)
		return
	}
	_ = r.ParseForm()
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	name := strings.TrimSpace(r.FormValue("name"))
	role := strings.TrimSpace(r.FormValue("role"))
	title := strings.TrimSpace(r.FormValue("title"))
	if _, err := auth.AddStaff(email, name, role, title); err != nil {
		view.PageStatus(w, r, 422, "admin-staff-edit", "Edit Staff | Patriot Pest Control", "", "", map[string]any{
			"AppUI": true, "UserType": "staff", "IsAdmin": true,
			"Csrf": view.CSRFField(w, r), "IsNew": true,
			"StaffRow": map[string]string{"email": email, "name": name, "role": role, "title": title},
			"Flash":    "Could not save: " + err.Error(),
		})
		return
	}
	http.Redirect(w, r, "/admin/people", http.StatusSeeOther)
}

// ---- admin settings (persisted, hot-reloadable) ----

const settingsPath = "storage/settings.json"

func loadSettings() map[string]string {
	out := map[string]string{
		"site_name": "Patriot Pest Control", "fb_pixel": "", "gtag_id": "", "gads_id": "", "clarity_id": "",
	}
	raw, err := os.ReadFile(settingsPath)
	if err != nil {
		return out
	}
	_ = json.Unmarshal(raw, &out)
	return out
}

func saveSettings(s map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return err
	}
	raw, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(settingsPath, raw, 0o644)
}

func (m *Module) AdminSettingsPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-settings"); !ok {
		return
	}
	view.Page(w, r, "admin-settings", "Site Settings | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
		"Csrf": view.CSRFField(w, r), "Settings": loadSettings(),
	})
}

func (m *Module) AdminSettingsSave(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-settings"); !ok {
		return
	}
	if !view.VerifyCSRF(r) {
		http.Error(w, "bad csrf token", http.StatusForbidden)
		return
	}
	_ = r.ParseForm()
	s := loadSettings()
	for _, k := range []string{"site_name", "fb_pixel", "gtag_id", "gads_id", "clarity_id"} {
		if v, ok2 := r.PostForm[k]; ok2 && len(v) > 0 {
			s[k] = strings.TrimSpace(v[0])
		}
	}
	if err := saveSettings(s); err != nil {
		http.Error(w, "save failed: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
}

// ---- api-keys lifecycle: audit / revoke / rotate / scopes ----

const auditPath = "storage/apikeys-audit.jsonl"

func auditLog(action, label, detail string) {
	_ = os.MkdirAll(filepath.Dir(auditPath), 0o755)
	f, err := os.OpenFile(auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	row, _ := json.Marshal(map[string]string{
		"at": time.Now().Format("2006-01-02 15:04:05"), "action": action, "label": label, "detail": detail,
	})
	f.Write(append(row, '\n'))
}

func auditRead() []map[string]string {
	raw, err := os.ReadFile(auditPath)
	if err != nil {
		return nil
	}
	var out []map[string]string
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var row map[string]string
		if json.Unmarshal([]byte(line), &row) == nil {
			out = append(out, row)
		}
	}
	// newest first, cap 100
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	if len(out) > 100 {
		out = out[:100]
	}
	return out
}

func (m *Module) ApiKeyAudit(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "apikey-audit"); !ok {
		return
	}
	keys := []map[string]string{}
	for i, k := range api.ListKeys() {
		prefix := k.Token
		if len(prefix) > 13 {
			prefix = prefix[:13]
		}
		keys = append(keys, map[string]string{
			"id": strconv.Itoa(i), "label": k.Label, "prefix": prefix,
			"scopes": strings.Join(k.Scopes, ", "), "created": k.Created.Format("2006-01-02 15:04"),
		})
	}
	view.Page(w, r, "apikey-audit", "API Key Audit | Patriot Pest Control", "", "", map[string]any{
		"AppUI": true, "UserType": "staff", "IsAdmin": true,
		"Csrf": view.CSRFField(w, r), "Keys": keys, "Audit": auditRead(),
	})
}

func (m *Module) keyByID(idStr string) *api.Key {
	i, err := strconv.Atoi(idStr)
	if err != nil {
		return nil
	}
	keys := api.ListKeys()
	if i < 0 || i >= len(keys) {
		return nil
	}
	return keys[i]
}

func (m *Module) ApiKeyRevoke(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "apikey-audit"); !ok {
		return
	}
	if !view.VerifyCSRF(r) {
		http.Error(w, "bad csrf token", http.StatusForbidden)
		return
	}
	k := m.keyByID(chi.URLParam(r, "id"))
	if k == nil || !api.RevokeKey(k.Token) {
		http.Error(w, "unknown key", http.StatusNotFound)
		return
	}
	auditLog("revoke", k.Label, "key revoked via /admin/api-keys")
	http.Redirect(w, r, "/admin/api-keys/audit", http.StatusSeeOther)
}

func (m *Module) ApiKeyRotate(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "apikey-audit"); !ok {
		return
	}
	if !view.VerifyCSRF(r) {
		http.Error(w, "bad csrf token", http.StatusForbidden)
		return
	}
	k := m.keyByID(chi.URLParam(r, "id"))
	if k == nil {
		http.Error(w, "unknown key", http.StatusNotFound)
		return
	}
	nk, err := api.IssueKey(k.Label, strings.Join(k.Scopes, ","))
	if err != nil {
		http.Error(w, "rotate failed: "+err.Error(), 500)
		return
	}
	api.RevokeKey(k.Token)
	auditLog("rotate", k.Label, "token rotated; new token issued")
	_ = nk
	http.Redirect(w, r, "/admin/api-keys/audit", http.StatusSeeOther)
}

func (m *Module) ApiKeyScopes(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "apikey-audit"); !ok {
		return
	}
	if !view.VerifyCSRF(r) {
		http.Error(w, "bad csrf token", http.StatusForbidden)
		return
	}
	k := m.keyByID(chi.URLParam(r, "id"))
	if k == nil {
		http.Error(w, "unknown key", http.StatusNotFound)
		return
	}
	scopes := strings.TrimSpace(r.FormValue("scopes"))
	nk, err := api.IssueKey(k.Label, scopes)
	if err != nil {
		http.Error(w, "scope update failed: "+err.Error(), 422)
		return
	}
	api.RevokeKey(k.Token)
	auditLog("scopes", k.Label, "scopes set to "+strings.Join(nk.Scopes, ","))
	http.Redirect(w, r, "/admin/api-keys/audit", http.StatusSeeOther)
}

// ---- retention settings save ----

func (m *Module) RetentionSettingsSave(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.gateAdmin(w, r, "admin-settings"); !ok {
		return
	}
	if !view.VerifyCSRF(r) {
		http.Error(w, "bad csrf token", http.StatusForbidden)
		return
	}
	_ = r.ParseForm()
	s := loadSettings()
	for _, k := range []string{"retention_enabled", "retention_intervals", "retention_credit"} {
		if v, ok2 := r.PostForm[k]; ok2 && len(v) > 0 {
			s[k] = strings.TrimSpace(v[0])
		}
	}
	if err := saveSettings(s); err != nil {
		http.Error(w, "save failed: "+err.Error(), 500)
		return
	}
	auditLog("retention.settings", "system", "retention settings updated")
	http.Redirect(w, r, "/admin/retention", http.StatusSeeOther)
}

// ---- shared JSON helper ----

func writeJSONStatus(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// silence unused import when context is only used via auth calls
var _ = context.Background
var _ = fmt.Sprintf
