package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testHandlers(cfg Config) *Handlers {
	h := NewHandlers(cfg, NewOAuthClient())
	h.now = func() time.Time { return time.Unix(1_000_000, 0) }
	return h
}

func TestHandleMeGuest(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
	testHandlers(Config{SigningKey: []byte("k")}).handleMe(rec, req)

	var body map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["authenticated"] != false || body["is_owner"] != false {
		t.Fatalf("guest /me = %v", body)
	}
}

func TestHandleMeOwner(t *testing.T) {
	t.Parallel()
	cfg := Config{SigningKey: []byte("k"), OwnerUserID: 522224}
	tok, _ := SignIdentity(Identity{UserID: 522224, Login: "theding0x", IsOwner: true, Exp: time.Now().Add(time.Hour).Unix()}, cfg.SigningKey)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookie, Value: tok})
	testHandlers(cfg).handleMe(rec, req)

	var body map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["is_owner"] != true || body["login"] != "theding0x" {
		t.Fatalf("owner /me = %v", body)
	}
}

func TestHandleLoginRedirectsWhenConfigured(t *testing.T) {
	t.Parallel()
	cfg := Config{ClientID: "cid", ClientSecret: "s", SigningKey: []byte("k"), OAuthConfigured: true, RedirectBaseURL: "https://app.daskap.io/api"}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/github/login", nil)
	testHandlers(cfg).handleLogin(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("login status = %d, want 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc == "" {
		t.Fatal("login: no Location header")
	}
}

func TestHandleLoginUnavailableWhenUnconfigured(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/github/login", nil)
	testHandlers(Config{}).handleLogin(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("login status = %d, want 503", rec.Code)
	}
}

func TestHandleLogoutClearsCookie(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
	testHandlers(Config{}).handleLogout(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("logout status = %d, want 204", rec.Code)
	}
	if sc := rec.Result().Cookies(); len(sc) == 0 || sc[0].MaxAge >= 0 {
		t.Fatalf("logout did not expire cookie: %v", sc)
	}
}
