package view

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

// TestDashboardsCompile ensures every dashboard page body parses and executes
// against representative data so template typos fail at test time, not runtime.
func TestDashboardsCompile(t *testing.T) {
	tpl, err := compile()
	if err != nil {
		t.Fatalf("compile: %v", err.Error())
	}

	cases := []struct {
		page string
		data map[string]any
	}{
		{"dash-customer", map[string]any{
			"Name": "Test Customer", "Email": "t@p.p", "Phone": "(555) 100", "AccountNumber": "1001", "Address": "1 Main St",
			"Appointments": []map[string]string{{"date": "2026-08-30", "start": "09:00", "type": "Treatment", "status": "Scheduled"}},
			"Invoices":       []map[string]string{{"id": "INV-1", "date": "2026-07-01", "amount": "$180", "status": "Paid"}},
		}},
		{"dash-staff", map[string]any{
			"Stats":        []map[string]string{{"v": "545", "k": "Customers"}, {"v": "12", "k": "Appointments"}},
			"Appointments": []map[string]string{{"date": "2026-08-30", "start": "09:00", "type": "Treatment", "district": "wa", "status": "Scheduled"}},
		}},
		{"dash-admin", map[string]any{
			"Stats":  []map[string]string{{"v": "545", "k": "Customers"}, {"v": "170", "k": "AZ District"}},
			"Recent": []map[string]string{{"name": "Test Co", "district": "wa", "status": "Active", "date": "2026-08-01"}},
		}},
		{"dash-people", map[string]any{
			"Csrf":  "<input>",
			"Staff": []map[string]string{{"email": "a@p.p", "name": "A", "role": "admin", "title": "Ops"}},
		}},
		{"dash-keys", map[string]any{
			"Districts": []map[string]string{{"code": "wa", "base": "https://api.fieldroutes.com", "key": "sk-...", "color": "#2563eb"}},
			"Twilio":    map[string]string{"sid": "tw1", "token": "tkt", "phone": "+1555000"},
		}},
		{"dash-apikeys", map[string]any{
			"Csrf": "<input>",
			"Keys": []map[string]string{{"label": "Reporting", "token": "ppc_live_abc", "scopes": "customer:read", "created": "2026-01-01"}},
		}},
	}

	for _, c := range cases {
		t.Run(c.page, func(tt *testing.T) {
			d := withBase(httptest.NewRequest("GET", "/"+c.page, nil), c.page, "T", "d", "")
			for k, v := range c.data {
				d[k] = v
			}
			var buf bytes.Buffer
			if err := tpl.ExecuteTemplate(&buf, c.page, d); err != nil {
				tt.Errorf("%s: %v", c.page, err.Error())
			}
		})
	}
}
