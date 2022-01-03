package cloudflare_geoblock_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	cloudflarerules "traefik-plugin"
)

const (
	forwardedFor   = "X-Forwarded-For"
	cfConnectingIP = "Cf-Connecting-Ip"
	ipCountry      = "Cf-Ipcountry"
)

func allowedCountries() []string {
	return []string{"ID", "SG"}
}

func TestAllowedCountries(t *testing.T) {
	cfg := setupConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, _ := setupHandler(t, ctx, next, cfg)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(ipCountry, "ID")
	req.Header.Set(cfConnectingIP, "http://localhost")
	handler.ServeHTTP(recorder, req)
	assertStatusEqual(t, recorder, req, 200)
	assertHeaderForwardedFor(t, req)
}

func TestNotAllowedCountries(t *testing.T) {
	cfg := setupConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, _ := setupHandler(t, ctx, next, cfg)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(ipCountry, "IE")
	req.Header.Set(cfConnectingIP, "http://localhost")
	handler.ServeHTTP(recorder, req)
	assertStatusEqual(t, recorder, req, 403)
}

func TestDisabledConfig(t *testing.T) {
	cfg := setupConfig()
	cfg.Disabled = true

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, _ := setupHandler(t, ctx, next, cfg)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(ipCountry, "IE")
	req.Header.Set(cfConnectingIP, "http://localhost")
	handler.ServeHTTP(recorder, req)
	assertStatusEqual(t, recorder, req, 200)
}

func TestEmptyGeolocation(t *testing.T) {
	cfg := setupConfig()
	cfg.WhitelistCountry = []string{}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, _ := setupHandler(t, ctx, next, cfg)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(cfConnectingIP, "http://localhost")
	handler.ServeHTTP(recorder, req)
	assertStatusEqual(t, recorder, req, 403)
}

func setupHandler(t *testing.T, ctx context.Context, next http.HandlerFunc, cfg *cloudflarerules.Config) (http.Handler, error) {
	handler, err := cloudflarerules.New(ctx, next, cfg, "cloudflare-rules")
	if err != nil {
		t.Fatal(err)
	}
	return handler, err
}

func setupConfig() *cloudflarerules.Config {
	cfg := cloudflarerules.CreateConfig()
	cfg.WhitelistCountry = allowedCountries()
	return cfg
}

func assertStatusEqual(t *testing.T, recorder *httptest.ResponseRecorder, req *http.Request, status int) {
	res := recorder.Result()
	if res.StatusCode != status {
		t.Errorf("invalid status: %s", req.Response.Status)
	}
}

func assertHeaderForwardedFor(t *testing.T, req *http.Request) {
	fwd := req.Header.Get(forwardedFor)
	if fwd != "http://localhost" {
		t.Errorf("invalid forwarded for: %s", fwd)
	}
}
