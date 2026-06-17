package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GitHubUser is the subset of the GitHub /user response we need.
type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

// OAuthClient talks to GitHub. Endpoints are fields so tests can point them at
// an httptest server.
type OAuthClient struct {
	HTTP        *http.Client
	AuthorizeEP string
	TokenEP     string
	UserEP      string
}

// NewOAuthClient returns a client wired to GitHub's real endpoints.
func NewOAuthClient() *OAuthClient {
	return &OAuthClient{
		HTTP:        http.DefaultClient,
		AuthorizeEP: "https://github.com/login/oauth/authorize",
		TokenEP:     "https://github.com/login/oauth/access_token",
		UserEP:      "https://api.github.com/user",
	}
}

// AuthorizeURL builds the GitHub authorize redirect URL.
func (c *OAuthClient) AuthorizeURL(clientID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", "read:user")
	q.Set("state", state)
	return c.AuthorizeEP + "?" + q.Encode()
}

// ExchangeCode swaps an authorization code for an access token.
func (c *OAuthClient) ExchangeCode(ctx context.Context, clientID, clientSecret, code, redirectURI string) (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.TokenEP, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github token exchange: status %d", resp.StatusCode)
	}
	var out struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if out.AccessToken == "" {
		return "", fmt.Errorf("github token exchange: no token (%s)", out.Error)
	}
	return out.AccessToken, nil
}

// FetchUser retrieves the authenticated user's numeric id and login.
func (c *OAuthClient) FetchUser(ctx context.Context, token string) (GitHubUser, error) {
	var u GitHubUser
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.UserEP, nil)
	if err != nil {
		return u, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return u, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return u, fmt.Errorf("github user: status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return u, err
	}
	return u, nil
}
