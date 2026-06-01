package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// ListProductionRelationBases handles GET /v1/revenue/production-bases. {"items":[...]}, never null.
func (h *Handler) ListProductionRelationBases(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListProductionRelationBases(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.ProductionRelationBasis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// ListDistributionRelations handles GET /v1/revenue/distribution-relations. {"items":[...]}, never null.
func (h *Handler) ListDistributionRelations(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListDistributionRelations(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.DistributionRelation{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// ListHistoricalTransiences handles GET /v1/revenue/historical-transiences. {"items":[...]}, never null.
func (h *Handler) ListHistoricalTransiences(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListHistoricalTransiences(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.HistoricalTransience{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createProductiveForceContradictionRequest is the body for POST /v1/revenue/productive-force-contradictions.
type createProductiveForceContradictionRequest struct {
	ProductiveForce     string `json:"productive_force"`
	SocialForm          string `json:"social_form"`
	ContradictionNature string `json:"contradiction_nature"`
	CrisisSignal        string `json:"crisis_signal"`
}

// CreateProductiveForceContradiction handles POST /v1/revenue/productive-force-contradictions.
// 400 bad JSON, 422 empty productive_force/social_form, 201 + Location.
func (h *Handler) CreateProductiveForceContradiction(w http.ResponseWriter, r *http.Request) {
	var req createProductiveForceContradictionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ProductiveForce == "" || req.SocialForm == "" {
		writeError(w, http.StatusUnprocessableEntity, "productive_force and social_form must not be empty")
		return
	}
	saved, err := h.Store.CreateProductiveForceContradiction(r.Context(), revenue.ProductiveForceContradiction{
		ProductiveForce:     req.ProductiveForce,
		SocialForm:          req.SocialForm,
		ContradictionNature: req.ContradictionNature,
		CrisisSignal:        req.CrisisSignal,
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/revenue/productive-force-contradictions")
	writeJSON(w, http.StatusCreated, saved)
}

// ListProductiveForceContradictions handles GET /v1/revenue/productive-force-contradictions. {"items":[...]}, never null.
func (h *Handler) ListProductiveForceContradictions(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListProductiveForceContradictions(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.ProductiveForceContradiction{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
