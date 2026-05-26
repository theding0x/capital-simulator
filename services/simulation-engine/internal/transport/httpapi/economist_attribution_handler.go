package httpapi

// Vol. II Ch. 10 — Theories of Fixed and Circulating Capital:
// The Physiocrats and Adam Smith.

import "net/http"

// ListEconomistAttributions handles GET /v1/circulation/economist-attributions.
// Optional query param: theorist (e.g. ?theorist=Quesnay or ?theorist=Smith).
// Returns an empty array (never null) when no records match.
func (h *Handler) ListEconomistAttributions(w http.ResponseWriter, r *http.Request) {
	theorist := r.URL.Query().Get("theorist")
	attrs, err := h.EconomistAttributions.ListEconomistAttributions(r.Context(), theorist)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, attrs)
}
