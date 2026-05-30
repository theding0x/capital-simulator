package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// GetPartIISummary handles GET /v1/avgprofit/summary — a read-only snapshot
// closing Part II. It gathers the production spheres, prices of production, and
// wage-effect analyses recorded across Part II and computes the conservation
// law: total price of production equals total value, and for the
// average-composition capital surplus-value produced equals average profit.
func (h *Handler) GetPartIISummary(w http.ResponseWriter, r *http.Request) {
	spheres, err := h.Store.ListProductionSpheres(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	prices, err := h.Store.ListPricesOfProduction(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	wageAnalyses, err := h.Store.ListWageEffectAnalyses(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	summary := avgprofit.ComputePartIISummary(spheres, prices, wageAnalyses)
	writeJSON(w, http.StatusOK, summary)
}
