package auth

import (
	"os"
	"strconv"
)

// Config holds the gateway's auth settings, read from the environment.
type Config struct {
	Disabled        bool   // AUTH_DISABLED=true → treat every request as owner (local dev)
	ClientID        string // GitHub OAuth app client id
	ClientSecret    string // GitHub OAuth app client secret
	SigningKey       []byte // HMAC key for session + state cookies
	OwnerUserID     int64  // numeric GitHub id of the owner
	RedirectBaseURL string // e.g. https://app.daskap.io/api
	OAuthConfigured bool   // true when client id/secret/key are all set
}

// IsOwnerID reports whether the given GitHub numeric id is the configured
// owner. Returns false when no owner id is configured (fail-closed).
func (cfg Config) IsOwnerID(id int64) bool {
	return cfg.OwnerUserID != 0 && id == cfg.OwnerUserID
}

// ConfigFromEnv builds Config from process environment variables.
func ConfigFromEnv() Config {
	cfg := Config{
		Disabled:        os.Getenv("AUTH_DISABLED") == "true",
		ClientID:        os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
		ClientSecret:    os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
		SigningKey:       []byte(os.Getenv("SESSION_SIGNING_KEY")),
		RedirectBaseURL: os.Getenv("OAUTH_REDIRECT_BASE_URL"),
	}
	if v := os.Getenv("OWNER_GITHUB_USER_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			cfg.OwnerUserID = id
		}
	}
	cfg.OAuthConfigured = cfg.ClientID != "" && cfg.ClientSecret != "" && len(cfg.SigningKey) > 0
	return cfg
}
