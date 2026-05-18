package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 33 — The Modern Theory of Colonisation.

type colonialMarketRequest struct {
	Colony          string `json:"colony"`
	FreeLabourers   int64  `json:"free_labourers"`
	AnnualWagePence int64  `json:"annual_wage_pence"`
	LandAvailable   bool   `json:"land_available"`
}

type colonialMarketResponse struct {
	ID                       string `json:"id"`
	Colony                   string `json:"colony"`
	FreeLabourers            int64  `json:"free_labourers"`
	AnnualWagePence          int64  `json:"annual_wage_pence"`
	LandAvailable            bool   `json:"land_available"`
	WakefieldSchemeApplied   bool   `json:"wakefield_scheme_applied"`
	IndependenceYears        int64  `json:"independence_years"`
	SurplusLabourExtractable bool   `json:"surplus_labour_extractable"`
	CreatedAt                string `json:"created_at,omitempty"`
}

type sufficientPriceRequest struct {
	PricePerAcrePence int64 `json:"price_per_acre_pence"`
	DesiredAcres      int64 `json:"desired_acres"`
}

type regulateRequest struct {
	SufficientPrice sufficientPriceRequest `json:"sufficient_price"`
}

type independenceRequest struct {
	WorkerID        string                 `json:"worker_id"`
	SufficientPrice sufficientPriceRequest `json:"sufficient_price"`
	YearsLimit      int64                  `json:"years_limit,omitempty"`
}

type independenceResponse struct {
	WorkerID           string `json:"worker_id"`
	AnnualSavingsPence int64  `json:"annual_savings_pence"`
	YearsWorked        int64  `json:"years_worked"`
	SavingsPence       int64  `json:"savings_pence"`
	TargetRansomPence  int64  `json:"target_ransom_pence"`
	BecameLandowner    bool   `json:"became_landowner"`
}

func toColonialMarketResponse(m simulation.ColonialLabourMarket) colonialMarketResponse {
	resp := colonialMarketResponse{
		ID:                       string(m.ID),
		Colony:                   m.Colony,
		FreeLabourers:            m.FreeLabourers,
		AnnualWagePence:          int64(m.AnnualWagePence),
		LandAvailable:            m.LandAvailable,
		WakefieldSchemeApplied:   m.WakefieldSchemeApplied,
		IndependenceYears:        m.IndependenceYears,
		SurplusLabourExtractable: m.SurplusLabourExtractable,
	}
	if !m.CreatedAt.IsZero() {
		resp.CreatedAt = m.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")
	}
	return resp
}

// CreateColonialMarket handles POST /v1/colonial-markets. Persists a new
// colonial labour market in its unregulated state.
func (h *Handler) CreateColonialMarket(w http.ResponseWriter, r *http.Request) {
	var req colonialMarketRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	market := simulation.ColonialLabourMarket{
		Colony:          req.Colony,
		FreeLabourers:   req.FreeLabourers,
		AnnualWagePence: simulation.Pence(req.AnnualWagePence),
		LandAvailable:   req.LandAvailable,
	}
	stored, err := h.ColonialMarkets.CreateColonialLabourMarket(r.Context(), market)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrAlreadyExists):
			writeError(w, http.StatusConflict, "colony already exists")
		case errors.Is(err, simulation.ErrColonyRequired),
			errors.Is(err, simulation.ErrFreeLabourersNonNegative),
			errors.Is(err, simulation.ErrAnnualWagePositive):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/v1/colonial-markets/%s", string(stored.ID)))
	writeJSON(w, http.StatusCreated, toColonialMarketResponse(stored))
}

// ListColonialMarkets handles GET /v1/colonial-markets.
func (h *Handler) ListColonialMarkets(w http.ResponseWriter, r *http.Request) {
	markets, err := h.ColonialMarkets.ListColonialLabourMarkets(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]colonialMarketResponse, len(markets))
	for i, m := range markets {
		out[i] = toColonialMarketResponse(m)
	}
	writeJSON(w, http.StatusOK, out)
}

// GetColonialMarket handles GET /v1/colonial-markets/{id}.
func (h *Handler) GetColonialMarket(w http.ResponseWriter, r *http.Request) {
	id := simulation.ColonialLabourMarketID(r.PathValue("id"))
	if id.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	market, err := h.ColonialMarkets.GetColonialLabourMarket(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "colonial market not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toColonialMarketResponse(market))
}

// RegulateColonialMarket handles POST /v1/colonial-markets/{id}/regulate.
// Applies a SystematicColonisation scheme: computes the independence
// period at the colonial wage given the sufficient price, then persists
// the regulated state.
func (h *Handler) RegulateColonialMarket(w http.ResponseWriter, r *http.Request) {
	id := simulation.ColonialLabourMarketID(r.PathValue("id"))
	if id.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req regulateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	price := simulation.SufficientPrice{
		PricePerAcrePence: simulation.Pence(req.SufficientPrice.PricePerAcrePence),
		DesiredAcres:      req.SufficientPrice.DesiredAcres,
	}
	if err := price.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	market, err := h.ColonialMarkets.GetColonialLabourMarket(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "colonial market not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	scheme := simulation.SystematicColonisation{
		Colony:          market.Colony,
		SufficientPrice: price,
	}
	regulated := simulation.ColonialLabourRegulation(market, scheme)
	updated, err := h.ColonialMarkets.UpdateColonialLabourMarket(r.Context(), id, store.ColonialLabourMarketUpdate{
		WakefieldSchemeApplied:   &regulated.WakefieldSchemeApplied,
		IndependenceYears:        &regulated.IndependenceYears,
		SurplusLabourExtractable: &regulated.SurplusLabourExtractable,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "colonial market not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toColonialMarketResponse(updated))
}

// ComputeIndependence handles POST /v1/colonial-markets/{id}/independence.
// Projects how many wage-years a labourer must work at the colonial wage
// to save enough for the sufficient price; reports whether independence
// is reachable within the optional yearsLimit.
func (h *Handler) ComputeIndependence(w http.ResponseWriter, r *http.Request) {
	id := simulation.ColonialLabourMarketID(r.PathValue("id"))
	if id.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req independenceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.WorkerID == "" {
		writeError(w, http.StatusBadRequest, "worker_id is required")
		return
	}
	price := simulation.SufficientPrice{
		PricePerAcrePence: simulation.Pence(req.SufficientPrice.PricePerAcrePence),
		DesiredAcres:      req.SufficientPrice.DesiredAcres,
	}
	if err := price.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.YearsLimit < 0 {
		writeError(w, http.StatusBadRequest, "years_limit cannot be negative")
		return
	}
	market, err := h.ColonialMarkets.GetColonialLabourMarket(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "colonial market not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	indep := simulation.ComputeWageWorkerIndependence(req.WorkerID, market, price, req.YearsLimit)
	writeJSON(w, http.StatusOK, independenceResponse{
		WorkerID:           indep.WorkerID,
		AnnualSavingsPence: int64(indep.AnnualSavingsPence),
		YearsWorked:        indep.YearsWorked,
		SavingsPence:       int64(indep.SavingsPence),
		TargetRansomPence:  int64(indep.TargetRansomPence),
		BecameLandowner:    indep.BecameLandowner,
	})
}
