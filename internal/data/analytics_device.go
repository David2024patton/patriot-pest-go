// Device detection for first-party analytics. A small hand-rolled
// User-Agent parser (no new dependencies): splits each hit into
// device type (phone/tablet/desktop), OS, and browser. Bots are filtered
// before this ever runs, so the parser only needs to handle human UAs.
package data

import "strings"

// DeviceInfo is the parsed shape of a visitor's User-Agent.
type DeviceInfo struct {
	DeviceType string // phone | tablet | desktop
	OS         string // Windows | macOS | Linux | Android | iOS | ChromeOS | Other
	Browser    string // Chrome | Safari | Firefox | Edge | Samsung Internet | Other
}

// ParseDevice extracts device/OS/browser from a User-Agent string.
// Empty or unrecognized UAs land in desktop/Other/Other rather than
// dropping the hit.
func ParseDevice(ua string) DeviceInfo {
	d := DeviceInfo{DeviceType: "desktop", OS: "Other", Browser: "Other"}
	if ua == "" {
		return d
	}
	l := strings.ToLower(ua)

	// Device type + OS first, since OS hints live next to device tokens.
	switch {
	case strings.Contains(l, "iphone"):
		d.DeviceType, d.OS = "phone", "iOS"
	case strings.Contains(l, "ipad"):
		d.DeviceType, d.OS = "tablet", "iOS"
	case strings.Contains(l, "ipod"):
		d.DeviceType, d.OS = "phone", "iOS"
	case strings.Contains(l, "android"):
		d.OS = "Android"
		if strings.Contains(l, "mobile") {
			d.DeviceType = "phone"
		} else {
			d.DeviceType = "tablet"
		}
	case strings.Contains(l, "windows nt") || strings.Contains(l, "windows phone"):
		d.DeviceType, d.OS = "desktop", "Windows"
		if strings.Contains(l, "windows phone") {
			d.DeviceType = "phone"
		}
	case strings.Contains(l, "cros"):
		d.DeviceType, d.OS = "desktop", "ChromeOS"
	case strings.Contains(l, "mac os x") || strings.Contains(l, "macintosh"):
		d.DeviceType, d.OS = "desktop", "macOS"
	case strings.Contains(l, "linux"):
		d.DeviceType, d.OS = "desktop", "Linux"
	}

	// Browser. Order matters: Edge and Samsung Internet both embed
	// "Chrome", and iOS Chrome/Firefox run on WebKit but still identify
	// themselves (CriOS / FxiOS).
	switch {
	case strings.Contains(l, "edg/") || strings.Contains(l, "edga") || strings.Contains(l, "edgios") || strings.Contains(l, "edge/"):
		d.Browser = "Edge"
	case strings.Contains(l, "samsungbrowser"):
		d.Browser = "Samsung Internet"
	case strings.Contains(l, "crios"):
		d.Browser = "Chrome"
	case strings.Contains(l, "fxios"):
		d.Browser = "Firefox"
	case strings.Contains(l, "firefox") || strings.Contains(l, "fxios"):
		d.Browser = "Firefox"
	case strings.Contains(l, "chrome") || strings.Contains(l, "chromium"):
		d.Browser = "Chrome"
	case strings.Contains(l, "safari") && strings.Contains(l, "version/"):
		d.Browser = "Safari"
	case strings.Contains(l, "safari"):
		d.Browser = "Safari"
	}
	return d
}
