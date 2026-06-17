package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthorizeURL(t *testing.T) {
	t.Parallel()
	c := NewOAuthClient()
	got := c.AuthorizeURL("cid", "https://app.daskap.io/api/v1/auth/github/callback", "xyz")
	for _, want := range []string{"client_id=cid", "state=xyz", "scope=read%3Auser"} {
		if !strings.Contains(got, want) {
			t.Errorf("AuthorizeURL missing %q in %q", want, got)
		}
	}
}

func TestExchangeCodeAndFetchUser(t *testing.T) {
	t.Parallel()
	token := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"gho_abc"}`))
	}))
	defer token.Close()
	user := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer gho_abc" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":522224,"login":"theding0x"}`))
	}))
	defer user.Close()

	c := &OAuthClient{HTTP: http.DefaultClient, AuthorizeEP: "x", TokenEP: token.URL, UserEP: user.URL}
	tok, err := c.ExchangeCode(context.Background(), "cid", "sec", "code123", "redir")
	if err != nil || tok != "gho_abc" {
		t.Fatalf("ExchangeCode = %q, %v", tok, err)
	}
	u, err := c.FetchUser(context.Background(), tok)
	if err != nil {
		t.Fatalf("FetchUser: %v", err)
	}
	if u.ID != 522224 || u.Login != "theding0x" {
		t.Fatalf("FetchUser = %+v", u)
	}
}
