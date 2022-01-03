// Package cloudflare_geoblock block or allow traffic from Cloudflare Geolocation.
package cloudflare_geoblock

import (
	"context"
	"net/http"
)

const (
	forwardedFor   = "X-Forwarded-For"
	cfConnectingIP = "Cf-Connecting-Ip"
	ipCountry      = "Cf-Ipcountry"
)

// Config the plugin configuration.
type Config struct {
	WhitelistCountry []string `json:"whitelistCountry" toml:"whitelistCountry" yaml:"whitelistCountry"`
	Disabled         bool     `json:"disabled" toml:"disabled" yaml:"disabled"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// CloudflareRules a Demo plugin.
type CloudflareRules struct {
	next             http.Handler
	WhitelistCountry []string
	Disabled         bool
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &CloudflareRules{
		next:             next,
		WhitelistCountry: config.WhitelistCountry,
		Disabled:         config.Disabled,
	}, nil
}

func (a *CloudflareRules) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if a.Disabled {
		a.next.ServeHTTP(rw, req)
		return
	}

	var geoLocation, realIp string
	geoLocation = req.Header.Get(ipCountry)

	if geoLocation == "" {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	if len(a.WhitelistCountry) > 0 && !contains(a.WhitelistCountry, geoLocation) {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	realIp = req.Header.Get(cfConnectingIP)
	req.Header.Set(forwardedFor, realIp)

	a.next.ServeHTTP(rw, req)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
