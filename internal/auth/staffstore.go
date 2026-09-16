package auth

// StaffStore — the staff table behind passwordless login. Mirrors the PHP
// `staff` table (database/schema.sql): email is the unique login key, role
// gates the dashboard, active gates access. No passwords by design.
//
// Backed by an in-memory map loaded from the SQLite file at boot so the web
// flow has zero DB dependencies of its own (the Go app ships the same SQLite
// content catalog as the PHP app). A known seed set is embedded as a
// fail-open fallback when the DB file is absent, so login still works on a
// bare checkout.

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"

	_ "modernc.org/sqlite" // pure-Go driver (same as internal/data)
)

// Staff is one row of the staff table. Title is a display-only role title
// (e.g. "Senior Technician"); it is not part of the SQLite schema yet and is
// managed at runtime through People management (/admin/people).
type Staff struct {
	ID     int
	Email  string
	Name   string
	Role   string // 'staff' | 'admin' | 'super-user'
	Title  string
	Active bool
}

var (
	staffMu     sync.RWMutex
	staffByEmail = map[string]*Staff{} // key: lowercased email
)

// seedStaff is the embedded fallback, seeded from database/patriot.db
// (2026-08 snapshot). Only used when the DB file has no staff table.
var seedStaff = []*Staff{
	{ID: 1, Email: "ppc_info@patriotpest.pro", Name: "Patriot Pest Control", Role: "admin", Active: true},
	{ID: 2, Email: "mrose@patriotpest.pro", Name: "M. Rose", Role: "admin", Active: true},
	{ID: 3, Email: "jordan@patriotpest.pro", Name: "Jordan", Role: "staff", Active: true},
	{ID: 4, Email: "david.richard.patton@gmail.com", Name: "David Patton", Role: "admin", Active: true},
	{ID: 5, Email: "david@itak.net", Name: "David Patton", Role: "admin", Active: true},
}

// LoadStaff replaces the in-memory staff map from the SQLite file. Fail-open:
// a missing file or empty table falls back to seedStaff so login never blocks.
// Returns the number of staff loaded.
func LoadStaff(dbPath string) int {
	rows := make(map[string]*Staff)
	loaded := false
	if dbPath != "" {
		db, err := sql.Open("sqlite", "file:"+dbPath+"?_mode=ro")
		if err == nil {
			defer db.Close()
			var rs []row
			rs, err = queryStaff(db)
			if err == nil && len(rs) > 0 {
				for _, r := range rs {
					rows[strings.ToLower(r.Email)] = &Staff{ID: r.id, Email: r.Email, Name: r.Name, Role: r.role, Active: r.active}
				}
				loaded = true
			}
		}
	}
	if !loaded {
		for _, s := range seedStaff {
			rows[strings.ToLower(s.Email)] = s
		}
	}
	staffMu.Lock()
	staffByEmail = rows
	staffMu.Unlock()
	slog.Info("staff loaded", "count", len(rows), "source", dbPath)
	return len(rows)
}

type row struct {
	id     int
	Email  string
	Name   string
	role   string
	active bool
}

// queryStaff reads the staff table. The driver is modernc.org/sqlite (pure Go).
func queryStaff(db *sql.DB) ([]row, error) {
	rows, err := db.Query("SELECT id, email, name, role, active FROM staff")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []row
	for rows.Next() {
		var r row
		var active int
		if err := rows.Scan(&r.id, &r.Email, &r.Name, &r.role, &active); err != nil {
			return nil, err
		}
		r.active = active == 1
		out = append(out, r)
	}
	return out, rows.Err()
}

// StaffByEmail returns the staff row for an email (case-insensitive), or nil.
func StaffByEmail(email string) *Staff {
	staffMu.RLock()
	defer staffMu.RUnlock()
	return staffByEmail[strings.ToLower(strings.TrimSpace(email))]
}

// ListStaff returns a name-sorted snapshot of every staff row.
func ListStaff() []*Staff {
	staffMu.RLock()
	defer staffMu.RUnlock()
	out := make([]*Staff, 0, len(staffByEmail))
	for _, s := range staffByEmail {
		out = append(out, s)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

// AddStaff upserts a person into the staff store (People management). Email is
// the unique key; role must be one of 'staff' | 'admin' | 'super-user'. A new
// row gets the next ID and starts active. Returns the stored copy.
func AddStaff(email, name, role, title string) (*Staff, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)
	title = strings.TrimSpace(title)
	if email == "" || !strings.Contains(email, "@") {
		return nil, fmt.Errorf("email required")
	}
	if name == "" {
		return nil, fmt.Errorf("name required")
	}
	switch role {
	case "staff", "admin", "super-user":
	default:
		return nil, fmt.Errorf("role must be staff, admin or super-user")
	}

	staffMu.Lock()
	defer staffMu.Unlock()
	if s := staffByEmail[email]; s != nil {
		s.Name = name
		s.Role = role
		s.Title = title
		s.Active = true
		c := *s
		return &c, nil
	}
	maxID := 0
	for _, s := range staffByEmail {
		if s.ID > maxID {
			maxID = s.ID
		}
	}
	s := &Staff{ID: maxID + 1, Email: email, Name: name, Role: role, Title: title, Active: true}
	staffByEmail[email] = s
	return s, nil
}

// FindStaffForLogin resolves a login identifier to an active, non-super-user
// staff member. Mirrors AuthController::findStaffForLogin: super-user accounts
// are excluded so the elevated /su surface stays isolated from /login.
func FindStaffForLogin(identifier string) *Staff {
	s := StaffByEmail(identifier)
	if s == nil || !s.Active || s.Role == "super-user" {
		return nil
	}
	return s
}
