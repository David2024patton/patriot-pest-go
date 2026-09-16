package fieldroutes

import (
	"fmt"
	"sort"
	"strings"
)

// Appointment is one normalized service visit for the customer portal.
type Appointment struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	Start  string `json:"start,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
	Notes  string `json:"notes,omitempty"`
}

// PullAppointments returns a customer's service history from one district,
// most-recent first, capped at 25. Degrades to an error on transport failure
// so a dead endpoint never surfaces as a panic in a handler.
func (c *Client) PullAppointments(d District, frID string) ([]Appointment, error) {
	frID = strings.TrimSpace(frID)
	if frID == "" {
		return nil, nil
	}

	raw, err := c.request(d, "appointment/search", map[string]string{"customerID": frID, "includeData": "1"})
	if err != nil {
		return nil, err
	}

	// includeData=1 returns records inline; otherwise fall back to /appointment/get.
	recs, ok := raw["appointments"].([]any)
	if !ok {
		var ids []string
		for _, e := range toArr(raw["appointmentIDs"]) {
			if s := fmt.Sprintf("%v", e); strings.TrimSpace(s) != "" {
				ids = append(ids, s)
			}
		}
		if len(ids) == 0 {
			return nil, nil
		}
		got, err := c.request(d, "appointment/get", map[string]string{"appointmentIDs": strings.Join(ids, ",")})
		if err != nil {
			return nil, err
		}
		recs, ok = got["appointments"].([]any)
	}
	if !ok {
		return nil, nil
	}

	out := make([]Appointment, 0, len(recs))
	for _, v := range recs {
		if m, ok := v.(map[string]any); ok {
			out = append(out, normalizeAppointment(d.Code, m))
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Date > out[j].Date })
	if len(out) > 25 {
		out = out[:25]
	}
	return out, nil
}

// normalizeAppointment maps a raw FR appointment record to a display row.
func normalizeAppointment(district string, r map[string]any) Appointment {
	a := Appointment{
		Date:   parseTime(str(r["date"])),
		Start:  str(r["start"]),
		Type:   str(r["type"]),
		Status: str(r["status"]),
		Notes:  str(r["notes"]),
	}
	if id := str(r["appointmentID"]); id != "" {
		a.ID = id
	}
	if a.ID == "" {
		a.ID = "district:" + district + ":" + fmt.Sprintf("%v", r["customerID"])
	}
	return a
}

// toArr coerces an envelope value to []any (nil-safe).
func toArr(v any) []any {
	if a, ok := v.([]any); ok {
		return a
	}
	return nil
}
