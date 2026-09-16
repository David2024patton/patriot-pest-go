package fieldroutes

import "strings"

// NormalizePhone normalizes a phone to E.164 for US/Canada (+1XXXXXXXXXX).
// Mirrors the PHP normalizePhone: returns nil when unparseable.
func NormalizePhone(raw string) *string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	digits := digitsOnly(raw)
	n := len(digits)

	if strings.HasPrefix(raw, "+") {
		if n >= 11 && n <= 15 {
			v := "+" + digits
			return &v
		}
		return nil
	}
	if n == 11 && strings.HasPrefix(digits, "1") {
		v := "+" + digits
		return &v
	}
	if n == 10 {
		v := "+1" + digits
		return &v
	}
	if n >= 7 && n <= 15 {
		v := "+" + digits
		return &v
	}
	return nil
}

// digitsOnly strips non-digits.
func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// trim wraps strings.TrimSpace for a concise call-site.
func trim(s string) string { return strings.TrimSpace(s) }
