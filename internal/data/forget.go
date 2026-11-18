package data

// Visitor data erasure: the deletion path behind the privacy page's
// "delete my tracking data" control (critic cycle 1, finding 4).
//
// The first-party analytics store rows under a persistent visitor id.
// A visitor who asks for erasure gets every analytics row carrying
// their id deleted. The id itself is never logged by the caller; the
// returned count is the receipt.

import (
	"fmt"
)

// DeleteVisitorAnalytics removes every analytics event stored under
// visitorID and returns how many rows were deleted.
func DeleteVisitorAnalytics(dbPath, visitorID string) (int64, error) {
	db, err := AnalyticsDB(dbPath)
	if err != nil {
		return 0, err
	}
	if err := EnsureAnalyticsTables(db); err != nil {
		return 0, err
	}
	res, err := db.Exec(`DELETE FROM analytics_events WHERE visitor_id = ?`, visitorID)
	if err != nil {
		return 0, fmt.Errorf("visitor delete: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("visitor delete count: %w", err)
	}
	return n, nil
}
