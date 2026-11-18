package data

// SMS consent records: the structured audit trail for text-message
// opt-ins given on the public forms (/free-quote, /sms-opt-in).
//
// A consent written only as prose inside a lead's Notes field cannot
// answer a carrier or TCPA audit: who consented, when, from which
// form, and to exactly what wording. Every checked consent box on a
// public form therefore also lands here as its own row, with the
// exact wording the visitor saw and the wording version in force at
// the time. The lead note stays for staff readability; this table is
// the record of proof.

import (
	"database/sql"
	"fmt"
	"time"
)

// SMSConsentRecord is one stored opt-in.
type SMSConsentRecord struct {
	ConsentedAt    int64 // Unix seconds, UTC
	Name           string
	Phone          string
	Source         string // the form path, e.g. "/free-quote"
	Wording        string // the exact checkbox label the visitor agreed to
	WordingVersion string // version tag for the wording, e.g. "2026-11 v1"
	IP             string // submitting IP, as seen by the server
}

// EnsureConsentTables creates the consent schema when missing. Safe to
// run on every write; production databases pick the table up with no
// migration step.
func EnsureConsentTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS sms_consent_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		consented_at INTEGER NOT NULL,
		name TEXT NOT NULL DEFAULT '',
		phone TEXT NOT NULL,
		source TEXT NOT NULL DEFAULT '',
		wording TEXT NOT NULL,
		wording_version TEXT NOT NULL DEFAULT '',
		ip TEXT NOT NULL DEFAULT ''
	)`)
	if err != nil {
		return fmt.Errorf("consent schema: %w", err)
	}
	return nil
}

// StoreSMSConsent appends one consent record. ConsentedAt is stamped
// here when the caller leaves it zero, so the stored time is the
// server's receive time, never a client-supplied value.
func StoreSMSConsent(dbPath string, rec SMSConsentRecord) error {
	db, err := AnalyticsDB(dbPath)
	if err != nil {
		return err
	}
	if err := EnsureConsentTables(db); err != nil {
		return err
	}
	if rec.ConsentedAt == 0 {
		rec.ConsentedAt = time.Now().Unix()
	}
	_, err = db.Exec(`INSERT INTO sms_consent_records
		(consented_at, name, phone, source, wording, wording_version, ip)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		rec.ConsentedAt, rec.Name, rec.Phone, rec.Source,
		rec.Wording, rec.WordingVersion, rec.IP)
	if err != nil {
		return fmt.Errorf("consent insert: %w", err)
	}
	return nil
}

// ListSMSConsents returns stored consent records, newest first. Used
// by audits and tests; the public site never renders it.
func ListSMSConsents(dbPath string, limit int) ([]SMSConsentRecord, error) {
	db, err := AnalyticsDB(dbPath)
	if err != nil {
		return nil, err
	}
	if err := EnsureConsentTables(db); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	rows, err := db.Query(`SELECT consented_at, name, phone, source, wording, wording_version, ip
		FROM sms_consent_records ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("consent list: %w", err)
	}
	defer rows.Close()
	var out []SMSConsentRecord
	for rows.Next() {
		var rec SMSConsentRecord
		if err := rows.Scan(&rec.ConsentedAt, &rec.Name, &rec.Phone, &rec.Source,
			&rec.Wording, &rec.WordingVersion, &rec.IP); err != nil {
			return nil, fmt.Errorf("consent scan: %w", err)
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}
