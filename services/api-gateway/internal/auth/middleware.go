package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// IdentityFromRequest extracts and verifies the session identity, returning a
// zero (guest) Identity when no valid session cookie is present.
func (cfg Config) IdentityFromRequest(r *http.Request) Identity {
	c, err := r.Cookie(SessionCookie)
	if err != nil {
		return Identity{}
	}
	id, err := VerifyIdentity(c.Value, cfg.SigningKey, time.Now())
	if err != nil {
		return Identity{}
	}
	return id
}

// RequireWrite wraps next, allowing reads and compute-allowlisted writes for
// everyone, but restricting all other writes to the owner (else 403).
func (cfg Config) RequireWrite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.Disabled {
			next.ServeHTTP(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		if IsComputePath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if cfg.IdentityFromRequest(r).IsOwner {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "sign in as the owner to make changes",
		})
	})
}
