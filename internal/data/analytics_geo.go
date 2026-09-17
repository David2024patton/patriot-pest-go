// Coarse IP geolocation for first-party analytics. The MaxMind
// GeoLite2-City database is downloaded at Docker build time into
// /app/geo/GeoLite2-City.mmdb; when it is missing (build offline, local
// dev) lookups return an empty GeoInfo and everything else keeps working.
//
// Privacy: the raw IP is looked up and then discarded. Only the coarse
// result (country code/name, region, city) is stored. Raw IPs are never
// persisted anywhere in the analytics pipeline.
package data

import (
	"net"
	"strings"
	"sync"

	"github.com/oschwald/maxminddb-golang"
)

// GeoInfo is the coarse location attached to an analytics event.
type GeoInfo struct {
	CountryCode string
	CountryName string
	Region      string // subdivision (state/province), English name
	City        string
}

// Empty reports whether no location was resolved.
func (g GeoInfo) Empty() bool { return g.CountryCode == "" && g.City == "" }

var (
	geoMu sync.RWMutex
	geoDB *maxminddb.Reader
)

// InitGeoDB opens the GeoLite2-City database at path. A missing or broken
// file is not fatal: lookups simply return empty GeoInfo.
func InitGeoDB(path string) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return
	}
	geoMu.Lock()
	if geoDB != nil {
		geoDB.Close()
	}
	geoDB = db
	geoMu.Unlock()
}

// GeoConfigured reports whether a geo database is loaded.
func GeoConfigured() bool {
	geoMu.RLock()
	defer geoMu.RUnlock()
	return geoDB != nil
}

type mmdbCity struct {
	Country struct {
		ISOCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Subdivisions []struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
}

func enName(m map[string]string) string {
	if m == nil {
		return ""
	}
	if n := m["en"]; n != "" {
		return n
	}
	for _, n := range m {
		return n
	}
	return ""
}

// LookupGeo resolves an IP string to a coarse GeoInfo. Private, loopback,
// unparsable, or missing-DB inputs yield an empty GeoInfo. The caller must
// not retain the input IP beyond this call.
func LookupGeo(ipStr string) GeoInfo {
	geoMu.RLock()
	db := geoDB
	geoMu.RUnlock()
	if db == nil {
		return GeoInfo{}
	}
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil || ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() {
		return GeoInfo{}
	}
	var rec mmdbCity
	if err := db.Lookup(ip, &rec); err != nil {
		return GeoInfo{}
	}
	g := GeoInfo{
		CountryCode: strings.ToUpper(rec.Country.ISOCode),
		CountryName: enName(rec.Country.Names),
		City:        enName(rec.City.Names),
	}
	if len(rec.Subdivisions) > 0 {
		g.Region = enName(rec.Subdivisions[len(rec.Subdivisions)-1].Names)
	}
	if g.CountryCode == "" {
		return GeoInfo{}
	}
	return g
}
