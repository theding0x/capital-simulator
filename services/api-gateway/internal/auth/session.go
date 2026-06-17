// Package auth provides GitHub OAuth login and owner-only write enforcement
// for the api-gateway. Sessions are stateless HMAC-signed cookies.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Cookie names set by the gateway.
const (
	SessionCookie = "cs_session"
	StateCookie   = "cs_oauth_state"
)

// Identity is the authenticated principal carried in the session cookie.
type Identity struct {
	UserID  int64  `json:"user_id"`
	Login   string `json:"login"`
	IsOwner bool   `json:"is_owner"`
	Exp     int64  `json:"exp"` // unix seconds
}

var (
	ErrBadSignature = errors.New("auth: bad signature")
	ErrMalformed    = errors.New("auth: malformed token")
	ErrExpired      = errors.New("auth: token expired")
)

// SignIdentity encodes id as "base64url(json).base64url(hmac-sha256)".
func SignIdentity(id Identity, key []byte) (string, error) {
	payload, err := json.Marshal(id)
	if err != nil {
		return "", err
	}
	body := base64.RawURLEncoding.EncodeToString(payload)
	return body + "." + sign(body, key), nil
}

func sign(body string, key []byte) string {
	m := hmac.New(sha256.New, key)
	m.Write([]byte(body))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// VerifyIdentity verifies the signature and expiry, returning the Identity.
func VerifyIdentity(token string, key []byte, now time.Time) (Identity, error) {
	var id Identity
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return id, ErrMalformed
	}
	if !hmac.Equal([]byte(sign(parts[0], key)), []byte(parts[1])) {
		return id, ErrBadSignature
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return id, ErrMalformed
	}
	if err := json.Unmarshal(payload, &id); err != nil {
		return id, ErrMalformed
	}
	if now.Unix() > id.Exp {
		return id, ErrExpired
	}
	return id, nil
}
