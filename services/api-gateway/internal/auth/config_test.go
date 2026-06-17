package auth

import "testing"

func TestConfigFromEnv(t *testing.T) {
	// not parallel: mutates process env
	t.Setenv("GITHUB_OAUTH_CLIENT_ID", "cid")
	t.Setenv("GITHUB_OAUTH_CLIENT_SECRET", "secret")
	t.Setenv("SESSION_SIGNING_KEY", "key")
	t.Setenv("OWNER_GITHUB_USER_ID", "522224")
	t.Setenv("OAUTH_REDIRECT_BASE_URL", "https://app.daskap.io/api")
	t.Setenv("AUTH_DISABLED", "")

	cfg := ConfigFromEnv()
	if cfg.OwnerUserID != 522224 {
		t.Errorf("OwnerUserID = %d, want 522224", cfg.OwnerUserID)
	}
	if !cfg.OAuthConfigured {
		t.Error("OAuthConfigured = false, want true")
	}
	if cfg.Disabled {
		t.Error("Disabled = true, want false")
	}
}

func TestConfigDisabled(t *testing.T) {
	t.Setenv("AUTH_DISABLED", "true")
	if !ConfigFromEnv().Disabled {
		t.Error("Disabled = false, want true")
	}
}
