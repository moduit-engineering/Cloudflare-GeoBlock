// Package cfgeoblock block or allow traffic from Cloudflare Geolocation.
package cfgeoblock

import (
	"context"
	"net/http"
)

const (
	forwardedFor   = "X-Forwarded-For"
	cfConnectingIP = "Cf-Connecting-Ip"
	ipCountry      = "Cf-Ipcountry"
)

// Config for interact with traefik config.
type Config struct {
	WhitelistCountry []string `json:"whitelistCountry" toml:"whitelistCountry" yaml:"whitelistCountry"`
	Disabled         bool     `json:"disabled,omitempty" toml:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// CreateConfig create config data for the plugin.
func CreateConfig() *Config {
	return &Config{}
}

// CloudflareRules config struct for the plugin.
type CloudflareRules struct {
	next             http.Handler
	WhitelistCountry []string
	Disabled         bool
}

// New constructor for this plugin.
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

	var geoLocation, realIP string
	geoLocation = req.Header.Get(ipCountry)

	if geoLocation == "" {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	if len(a.WhitelistCountry) > 0 && !contains(a.WhitelistCountry, geoLocation) {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	realIP = req.Header.Get(cfConnectingIP)
	req.Header.Set(forwardedFor, realIP)

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
