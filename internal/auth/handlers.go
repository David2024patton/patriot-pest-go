package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Post("/api/auth/otp/issue", IssueHandler)
	r.Post("/api/auth/otp/verify", VerifyHandler)
	r.Post("/api/auth/magic-link", MagicLinkHandler)
	r.Get("/auth/verify", MagicVerifyHandler)
}

// jsonResp writes a JSON object as the handler response.
func jsonResp(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	if code != 200 {
		w.WriteHeader(code)
	}
	_ = json.NewEncoder(w).Encode(v)
}

func IssueHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Identity string `json:"identity"`
		Purpose  string `json:"purpose"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Identity == "" {
		jsonResp(w, 400, map[string]any{"error": "identity required"})
		return
	}
	if body.Purpose == "" {
		body.Purpose = "login"
	}

	digits := 6
	ttl := 600
	role := "Customer"
	var prof *Profile
	if body.Purpose == "super_login" {
		digits = 8
		ttl = 300
		role = "SuperAdmin"
	} else {
		prof = resolve(body.Identity)
	}

	// Customer login resolves the identity to a real record first. The canonical
	// OTP key is always the email on file, so verification works no matter which
	// identifier (account number / phone / email) the customer typed.
	var delivery, display string
	switch {
	case prof == nil:
		// Unknown identity or a non-login purpose: issue under the typed identity.
		delivery, display = body.Identity, body.Identity
	case prof.Email == "":
		jsonResp(w, 200, map[string]any{
			"status":   "email_required",
			"message":  "You need an email on file to receive your code.",
			"identity": body.Identity,
		})
		return
	default:
		delivery, display = prof.Email, prof.Email
	}

	code, _ := Issue(r.Context(), delivery, body.Purpose, digits, ttl)
	resp := map[string]any{"status": "sent", "identity": display, "role": role, "expires_in": ttl, "hint": code[:1] + "*****"}
	if prof != nil && prof.Email != "" {
		resp["email"] = prof.Email
	}
	// Debug affordance (OTP_DEBUG=1) returns the plaintext code so the round-trip
	// can be driven without SMTP. Default off — never set in production.
	if v := os.Getenv("OTP_DEBUG"); v == "1" || v == "true" {
		resp["debug_code"] = code
	}
	jsonResp(w, 200, resp)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Identity string `json:"identity"`
		Purpose  string `json:"purpose"`
		Code     string `json:"code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Purpose == "" {
		body.Purpose = "login"
	}

	// Canonicalize the identity the same way issue did, so verification looks up
	// the OTP under the email-on-file key. Non-login purposes use the typed value.
	identity := body.Identity
	var prof *Profile
	if body.Purpose != "super_login" {
		prof = resolve(body.Identity)
		if prof != nil && prof.Email != "" {
			identity = prof.Email
		}
	}
	ok, err := Verify(r.Context(), identity, body.Purpose, body.Code, 5)
	if !ok {
		jsonResp(w, 401, map[string]any{"error": err.Error()})
		return
	}
	ttl := 900
	role := "Customer"
	if body.Purpose == "super_login" {
		ttl = 7200
		role = "SuperAdmin"
	}
	// Session identity is the canonical email so the portal can resolve it back.
	sessIdentity := identity
	if prof != nil && prof.Email != "" {
		sessIdentity = prof.Email
	}
	sess := CreateSession(sessIdentity, role, ttl)
	http.SetCookie(w, &http.Cookie{Name: "session", Value: sess.ID, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode, Expires: time.Now().Add(time.Duration(ttl) * time.Second)})

	resp := map[string]any{"status": "ok", "session": sess.ID, "role": role}
	if prof != nil {
		resp["customer"] = profileJSON(prof)
	}
	jsonResp(w, 200, resp)
}

// profileJSON renders a customer profile for the login response.
func profileJSON(p *Profile) map[string]any {
	return map[string]any{
		"fr_id":          p.FRID,
		"district":       p.District,
		"name":           p.Name,
		"email":          p.Email,
		"phone":          p.Phone,
		"account_number": p.AccountNumber,
	}
}

func MagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	IssueHandler(w, r)
}
func MagicVerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	identity := r.URL.Query().Get("identity")
	if token == "" || identity == "" {
		http.Error(w, `{"error":"bad link"}`, 400)
		return
	}
	ok, _ := Verify(r.Context(), identity, "magic", token, 5)
	if !ok {
		http.Error(w, `{"error":"expired"}`, 401)
		return
	}
	sess := CreateSession(identity, "Customer", 900)
	http.SetCookie(w, &http.Cookie{Name: "session", Value: sess.ID, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})
	http.Redirect(w, r, "/customer-dashboard", 302)
}
