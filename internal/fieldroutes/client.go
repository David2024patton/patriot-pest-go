// Package fieldroutes is a read/write client for the FieldRoutes CRM API.
//
// FieldRoutes is the SOURCE OF TRUTH for customers. We pull customers from
// EVERY configured district (WA covers WA-ID-OR; AZ is separate) into a
// normalized row. The API contract is: GET {base}/api/{entity}/{action} with
// authenticationKey + authenticationToken as query params. The envelope is
// {success, ...IDs, customers|appointments|subscriptions, count, errorMessage}.
// customer/search returns IDs; customer/get returns records (batched at 200).
//
// Every network call degrades to an error rather than throwing so a dead
// endpoint never surfaces as a panic in a handler.
package fieldroutes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	batchSize = 200
	timeout   = 30 * time.Second
)

// District is one FieldRoutes account (one base URL, its own key + token).
type District struct {
	Code  string // "wa" | "az"
	Base  string // e.g. https://patriotpestc.fieldroutes.com
	Key   string
	Token string
}

// Row is the normalized customer row shape shared with the local cache.
type Row struct {
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
}

// Client holds one district per credential set. The district list is mutable
// at runtime (the admin keys page adds new districts live), so all access goes
// through the mutex.
type Client struct {
	mu        sync.RWMutex
	districts []District
	http      *http.Client
}

// New builds a client from the configured districts (skips incomplete ones).
func New(districts []District) *Client {
	var good []District
	for _, d := range districts {
		if d.Base != "" && d.Key != "" && d.Token != "" {
			good = append(good, d)
		}
	}
	return &Client{
		districts: good,
		http:      &http.Client{Timeout: timeout},
	}
}

// SetDistricts live-replaces the configured districts (used when a super-admin
// adds a new district key from /admin/keys). Incomplete entries are dropped.
func (c *Client) SetDistricts(ds []District) {
	var good []District
	for _, d := range ds {
		if d.Base != "" && d.Key != "" && d.Token != "" {
			good = append(good, d)
		}
	}
	c.mu.Lock()
	c.districts = good
	c.mu.Unlock()
}

// Districts returns a copy of the configured districts.
func (c *Client) Districts() []District {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]District, len(c.districts))
	copy(out, c.districts)
	return out
}

// DistrictByCode returns the configured district for a code, or ok=false.
func (c *Client) DistrictByCode(code string) (District, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, d := range c.districts {
		if strings.EqualFold(d.Code, code) {
			return d, true
		}
	}
	return District{}, false
}

// Missing returns the codes of districts NOT fully configured (for a friendly
// "what's needed" message). Codes are lowercased.
func (c *Client) Missing() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	have := map[string]bool{}
	for _, d := range c.districts {
		have[strings.ToLower(d.Code)] = true
	}
	var out []string
	for _, code := range []string{"wa", "az"} {
		if !have[code] {
			out = append(out, code)
		}
	}
	return out
}

// PullCustomers pulls every customer from one district as normalized rows.
func (c *Client) PullCustomers(d District) ([]Row, error) {
	ids, err := c.SearchCustomerIDs(d)
	if err != nil {
		return nil, err
	}
	var rows []Row
	for start := 0; start < len(ids); start += batchSize {
		end := start + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		recs, err := c.GetCustomers(d, ids[start:end])
		if err != nil {
			return nil, err
		}
		for _, r := range recs {
			rows = append(rows, Normalize(r, d.Code))
		}
	}
	return rows, nil
}

// SearchCustomerIDs returns the list of customer IDs for a district.
func (c *Client) SearchCustomerIDs(d District) ([]string, error) {
	data, err := c.request(d, "customer/search", nil)
	if err != nil {
		return nil, err
	}
	raw := data["customerIDs"]
	if raw == nil {
		raw = data["result"]
	}
	ids, ok := toStrings(raw)
	if !ok {
		return nil, nil
	}
	return ids, nil
}

// GetCustomers fetches customer records for a batch of IDs (one call, <=200).
func (c *Client) GetCustomers(d District, ids []string) ([]map[string]any, error) {
	data, err := c.request(d, "customer/get", map[string]string{"customerIDs": strings.Join(ids, ",")})
	if err != nil {
		return nil, err
	}
	raw, _ := data["customers"].([]any)
	var out []map[string]any
	for _, v := range raw {
		if m, ok := v.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// Normalize maps one FieldRoutes customer record onto a normalized cache row.
func Normalize(c map[string]any, district string) Row {
	name := trim(str(c["fname"]) + " " + str(c["lname"]))
	if name == "" {
		name = trim(str(c["companyName"]))
	}
	id := str(c["customerID"])

	status := "active"
	lower := strings.ToLower(str(c["statusText"]))
	if strings.Contains(lower, "cancel") {
		status = "cancelled"
	} else if strings.Contains(lower, "inactive") {
		status = "inactive"
	}
	if raw := str(c["dateCancelled"]); raw != "" && !strings.HasPrefix(raw, "0000-00-00") {
		status = "cancelled"
	}

	var lastService *string
	if raw := c["lastCompleted"]; raw != nil && strings.TrimSpace(str(raw)) != "" {
		if ts := parseTime(str(raw)); ts != "" {
			lastService = &ts
		}
	} else if raw := c["lastService"]; raw != nil && strings.TrimSpace(str(raw)) != "" {
		if ts := parseTime(str(raw)); ts != "" {
			lastService = &ts
		}
	}

	row := Row{
		FRID:          id,
		District:      district,
		Name:          nameIfEmpty(name, "Customer "+id),
		AccountNumber: id,
		Status:        status,
		LastService:   lastService,
	}
	row.Email = strPtr(c["email"])
	row.Phone = NormalizePhone(str(c["phone1"]))
	row.Address = strPtr(c["address"])
	row.City = strPtr(c["city"])
	row.State = strPtr(c["state"])
	row.Zip = strPtr(c["zip"])
	return row
}

// request issues one GET against the FieldRoutes API and returns the decoded
// envelope, or an error on transport/HTTP/JSON failure.
func (c *Client) request(d District, endpoint string, params map[string]string) (map[string]any, error) {
	q := url.Values{}
	q.Set("authenticationKey", d.Key)
	q.Set("authenticationToken", d.Token)
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u := strings.TrimRight(d.Base, "/") + "/api/" + strings.TrimLeft(endpoint, "/") + "?" + q.Encode()

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "PatriotPest/1.0 (+fieldroutes-sync)")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("FR HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("FR bad JSON: %w", err)
	}
	return data, nil
}

// toStrings converts a JSON value to []string, returning ok=false on mismatch.
func toStrings(v any) ([]string, bool) {
	arr, ok := v.([]any)
	if !ok {
		return nil, false
	}
	var out []string
	for _, e := range arr {
		s := fmt.Sprintf("%v", e)
		if strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out, true
}

// str returns the string form of a JSON value (nil/absent -> "").
func str(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	if f, ok := v.(float64); ok {
		return fmt.Sprintf("%v", f)
	}
	return fmt.Sprintf("%v", v)
}

// strPtr returns a *string pointer or nil for empty values.
func strPtr(v any) *string {
	s := str(v)
	if s == "" {
		return nil
	}
	return &s
}

// parseTime converts a FR date to "Y-m-d" for the last_service signal.
func parseTime(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if t, err := time.Parse("2006-01-02", raw[:10]); err == nil {
		return t.Format("2006-01-02")
	}
	if t, err := time.Parse("2006-01-02 15:04:05", raw); err == nil {
		return t.Format("2006-01-02")
	}
	return raw
}

// nameIfEmpty falls back to a synthetic name.
func nameIfEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
