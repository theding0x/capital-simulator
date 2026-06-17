package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func ownerCookie(t *testing.T, cfg Config) *http.Cookie {
	t.Helper()
	tok, err := SignIdentity(Identity{UserID: 522224, Login: "theding0x", IsOwner: true, Exp: time.Now().Add(time.Hour).Unix()}, cfg.SigningKey)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Cookie{Name: SessionCookie, Value: tok}
}

func baseCfg() Config {
	return Config{SigningKey: []byte("k"), OwnerUserID: 522224}
}

func TestMiddlewareAllowsGuestReads(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/commodities", nil)
	baseCfg().RequireWrite(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("guest GET: got %d, want 200", rec.Code)
	}
}

func TestMiddlewareBlocksGuestWrites(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/commodities", nil)
	baseCfg().RequireWrite(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("guest POST: got %d, want 403", rec.Code)
	}
}

func TestMiddlewareAllowsGuestComputeWrites(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", nil)
	baseCfg().RequireWrite(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("guest compute POST: got %d, want 200", rec.Code)
	}
}

func TestMiddlewareAllowsOwnerWrites(t *testing.T) {
	t.Parallel()
	cfg := baseCfg()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/commodities", nil)
	req.AddCookie(ownerCookie(t, cfg))
	cfg.RequireWrite(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("owner POST: got %d, want 200", rec.Code)
	}
}

func TestMiddlewareDisabledBypassesAll(t *testing.T) {
	t.Parallel()
	cfg := baseCfg()
	cfg.Disabled = true
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/commodities", nil)
	cfg.RequireWrite(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("disabled POST: got %d, want 200", rec.Code)
	}
}
