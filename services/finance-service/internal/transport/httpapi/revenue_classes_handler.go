package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// ListSocialClasses handles GET /v1/revenue/classes. {"items":[...]}, never null.
func (h *Handler) ListSocialClasses(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListSocialClasses(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.SocialClass{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// socialClassDetailResponse bundles a class with the income sources and tendencies that reference it.
type socialClassDetailResponse struct {
	Class         revenue.SocialClass         `json:"class"`
	IncomeSources []revenue.ClassIncomeSource `json:"income_sources"`
	Tendencies    []revenue.ClassTendency     `json:"tendencies"`
}

// GetSocialClass handles GET /v1/revenue/classes/{id}. 404 when no class has the given
// ID, 200 + {class, income_sources, tendencies} otherwise.
func (h *Handler) GetSocialClass(w http.ResponseWriter, r *http.Request) {
	id := revenue.SocialClassID(r.PathValue("id"))
	class, err := h.Store.GetSocialClass(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	allSources, err := h.Store.ListClassIncomeSources(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	allTendencies, err := h.Store.ListClassTendencies(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	sources := []revenue.ClassIncomeSource{}
	for _, s := range allSources {
		if s.ClassID == string(id) {
			sources = append(sources, s)
		}
	}
	tendencies := []revenue.ClassTendency{}
	for _, t := range allTendencies {
		if t.ClassID == string(id) {
			tendencies = append(tendencies, t)
		}
	}
	writeJSON(w, http.StatusOK, socialClassDetailResponse{Class: class, IncomeSources: sources, Tendencies: tendencies})
}

// ListClassAmbiguities handles GET /v1/revenue/class-ambiguities. {"items":[...]}, never null.
func (h *Handler) ListClassAmbiguities(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListClassAmbiguities(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.ClassAmbiguity{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
