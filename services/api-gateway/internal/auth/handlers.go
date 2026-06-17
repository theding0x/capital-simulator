package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

const sessionTTL = 7 * 24 * time.Hour

// Handlers serves the OAuth + session endpoints.
type Handlers struct {
	cfg    Config
	client *OAuthClient
	now    func() time.Time
}

// NewHandlers builds the auth handler set.
func NewHandlers(cfg Config, client *OAuthClient) *Handlers {
	return &Handlers{cfg: cfg, client: client, now: time.Now}
}

// Register attaches the auth routes via the given registrar (srv.HandleFunc).
func (h *Handlers) Register(reg func(pattern string, handler http.HandlerFunc)) {
	reg("GET /v1/auth/github/login", h.handleLogin)
	reg("GET /v1/auth/github/callback", h.handleCallback)
	reg("POST /v1/auth/logout", h.handleLogout)
	reg("GET /v1/auth/me", h.handleMe)
}

func (h *Handlers) redirectURI() string {
	return h.cfg.RedirectBaseURL + "/v1/auth/github/callback"
}

func randomState() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (h *Handlers) setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.OAuthConfigured {
		http.Error(w, "github login not configured", http.StatusServiceUnavailable)
		return
	}
	state := randomState()
	h.setCookie(w, StateCookie, state, 600)
	http.Redirect(w, r, h.client.AuthorizeURL(h.cfg.ClientID, h.redirectURI(), state), http.StatusFound)
}

func (h *Handlers) handleCallback(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.OAuthConfigured {
		http.Error(w, "github login not configured", http.StatusServiceUnavailable)
		return
	}
	stateCookie, err := r.Cookie(StateCookie)
	if err != nil || stateCookie.Value == "" || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}
	token, err := h.client.ExchangeCode(r.Context(), h.cfg.ClientID, h.cfg.ClientSecret, code, h.redirectURI())
	if err != nil {
		http.Error(w, "oauth exchange failed", http.StatusBadGateway)
		return
	}
	user, err := h.client.FetchUser(r.Context(), token)
	if err != nil {
		http.Error(w, "oauth user fetch failed", http.StatusBadGateway)
		return
	}
	id := Identity{
		UserID:  user.ID,
		Login:   user.Login,
		IsOwner: user.ID == h.cfg.OwnerUserID,
		Exp:     h.now().Add(sessionTTL).Unix(),
	}
	signed, err := SignIdentity(id, h.cfg.SigningKey)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}
	h.setCookie(w, SessionCookie, signed, int(sessionTTL.Seconds()))
	h.setCookie(w, StateCookie, "", -1)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handlers) handleLogout(w http.ResponseWriter, _ *http.Request) {
	h.setCookie(w, SessionCookie, "", -1)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) handleMe(w http.ResponseWriter, r *http.Request) {
	id := h.cfg.IdentityFromRequest(r)
	resp := map[string]any{
		"authenticated": id.UserID != 0,
		"login":         id.Login,
		"is_owner":      id.IsOwner,
		"oauth_enabled": h.cfg.OAuthConfigured && !h.cfg.Disabled,
	}
	if h.cfg.Disabled {
		// Local dev: gateway treats everyone as the owner.
		resp["authenticated"] = true
		resp["is_owner"] = true
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
