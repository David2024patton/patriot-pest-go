package data

import (
	"database/sql"
	"fmt"
	"strings"
)

// SignupResult reports the outcome of a website signup. Err carries the
// underlying failure for server logs only — Message is what the user sees.
type SignupResult struct {
	OK      bool
	Message string
	Err     error
}

// CreateSignup writes a website signup into the customers table with
// source='website' so marketing tracking can attribute it. Dedupes on email:
// a returning email updates the row instead of creating a duplicate.
// dbPath is the configured catalog path — never a second hardcoded location.
func CreateSignup(dbPath, name, email, phone, city, state, zip string) SignupResult {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return SignupResult{false, "Email is required.", nil}
	}

	db, err := sql.Open("sqlite", "file:"+dbPath+"?_pragma=journal_mode(WAL)&_busy_timeout=5000")
	if err != nil {
		return SignupResult{false, "Could not open the database.", fmt.Errorf("signup open db: %w", err)}
	}
	defer db.Close()

	var id int
	err = db.QueryRow(`SELECT id FROM customers WHERE email = ? COLLATE NOCASE LIMIT 1`, email).Scan(&id)
	if err == sql.ErrNoRows {
		_, err = db.Exec(`INSERT INTO customers (name, email, phone, city, state, zip, status, source, district)
			VALUES (?, ?, ?, ?, ?, ?, 'active', 'website', 'wa')`,
			strings.TrimSpace(name), email, strings.TrimSpace(phone),
			strings.TrimSpace(city), strings.TrimSpace(state), strings.TrimSpace(zip))
		if err != nil {
			return SignupResult{false, "Could not save your account.", fmt.Errorf("signup insert: %w", err)}
		}
		return SignupResult{true, "Account created. Check your email for what happens next.", nil}
	}
	if err != nil {
		return SignupResult{false, "Could not save your account.", fmt.Errorf("signup lookup: %w", err)}
	}
	// Existing email: refresh contact info rather than duplicate.
	_, err = db.Exec(`UPDATE customers SET name = COALESCE(NULLIF(?, ''), name), phone = COALESCE(NULLIF(?, ''), phone),
		city = COALESCE(NULLIF(?, ''), city), state = COALESCE(NULLIF(?, ''), state), zip = COALESCE(NULLIF(?, ''), zip),
		updated_at = datetime('now') WHERE id = ?`,
		strings.TrimSpace(name), strings.TrimSpace(phone), strings.TrimSpace(city), strings.TrimSpace(state), strings.TrimSpace(zip), id)
	if err != nil {
		return SignupResult{false, "Could not update your account.", fmt.Errorf("signup update: %w", err)}
	}
	return SignupResult{true, "Welcome back. Your account details were updated.", nil}
}
