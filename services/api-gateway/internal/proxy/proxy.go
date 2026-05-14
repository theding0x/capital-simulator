// Package proxy builds reverse-proxy handlers for fanning gateway requests
// out to downstream services.
package proxy

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// New returns a reverse-proxy handler whose target is the given URL string
// (e.g. http://commodity-service:8081). Failed dials produce a 502.
func New(target string, logger *slog.Logger) (http.Handler, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("proxy: parse target %q: %w", target, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("proxy: target %q must include scheme and host", target)
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Warn("upstream error",
			slog.String("target", target),
			slog.String("path", r.URL.Path),
			slog.String("err", err.Error()),
		)
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
	}
	originalDirector := rp.Director
	rp.Director = func(r *http.Request) {
		originalDirector(r)
		// Hide upstream Host so the downstream service sees the gateway.
		r.Host = u.Host
	}
	return rp, nil
}
