package httpapi

import (
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

const defaultReferenceDayMinutes = 600

type registerIntensityRequest struct {
	CountryCode string  `json:"country_code"`
	Factor      float64 `json:"factor"`
}

type registerDayWageRequest struct {
	CountryCode       string `json:"country_code"`
	NominalPence      int64  `json:"nominal_pence"`
	WorkingDayMinutes int64  `json:"working_day_minutes"`
}

// RegisterIntensity handles POST /v1/intensities.
func (h *Handler) RegisterIntensity(w http.ResponseWriter, r *http.Request) {
	var req registerIntensityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ni := agent.NationalIntensity{
		CountryCode: agent.CountryCode(req.CountryCode),
		Factor:      req.Factor,
	}
	if err := ni.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.NationalWageStore.UpsertIntensity(r.Context(), ni)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

// ListIntensities handles GET /v1/intensities.
func (h *Handler) ListIntensities(w http.ResponseWriter, r *http.Request) {
	items, err := h.NationalWageStore.ListIntensities(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if items == nil {
		items = []agent.NationalIntensity{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// RegisterDayWage handles POST /v1/wages.
func (h *Handler) RegisterDayWage(w http.ResponseWriter, r *http.Request) {
	var req registerDayWageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	dw := agent.DayWage{
		CountryCode:       agent.CountryCode(req.CountryCode),
		NominalPence:      req.NominalPence,
		WorkingDayMinutes: req.WorkingDayMinutes,
	}
	if err := dw.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.NationalWageStore.CreateDayWage(r.Context(), dw)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

// GetStandardisedWage handles GET /v1/wages/{country}/standardised.
func (h *Handler) GetStandardisedWage(w http.ResponseWriter, r *http.Request) {
	country := agent.CountryCode(r.PathValue("country"))
	refParam := r.URL.Query().Get("reference_day_minutes")
	refDay := int64(defaultReferenceDayMinutes)
	if refParam != "" {
		parsed, err := strconv.ParseInt(refParam, 10, 64)
		if err != nil || parsed <= 0 {
			writeError(w, http.StatusBadRequest, "reference_day_minutes must be a positive integer")
			return
		}
		refDay = parsed
	}
	dw, err := h.NationalWageStore.GetDayWage(r.Context(), country)
	if err != nil {
		writeAppError(w, err)
		return
	}
	sw := agent.StandardiseWage(dw, refDay)
	writeJSON(w, http.StatusOK, sw)
}

// GetWageComparison handles GET /v1/comparisons.
func (h *Handler) GetWageComparison(w http.ResponseWriter, r *http.Request) {
	refParam := r.URL.Query().Get("reference_day_minutes")
	refDay := int64(defaultReferenceDayMinutes)
	if refParam != "" {
		parsed, err := strconv.ParseInt(refParam, 10, 64)
		if err != nil || parsed <= 0 {
			writeError(w, http.StatusBadRequest, "reference_day_minutes must be a positive integer")
			return
		}
		refDay = parsed
	}

	ctx := r.Context()
	wages, err := h.NationalWageStore.ListDayWages(ctx)
	if err != nil {
		writeAppError(w, err)
		return
	}
	intensities, err := h.NationalWageStore.ListIntensities(ctx)
	if err != nil {
		writeAppError(w, err)
		return
	}
	spindleRatios, err := h.NationalWageStore.ListSpindleRatios(ctx)
	if err != nil {
		writeAppError(w, err)
		return
	}

	intensityByCountry := make(map[agent.CountryCode]agent.NationalIntensity, len(intensities))
	for _, ni := range intensities {
		intensityByCountry[ni.CountryCode] = ni
	}

	countries := make([]agent.CountryCode, 0, len(wages))
	standardised := make([]agent.StandardisedWage, 0, len(wages))
	relativePrices := make([]agent.RelativeLabourPrice, 0, len(wages))

	for _, dw := range wages {
		countries = append(countries, dw.CountryCode)
		standardised = append(standardised, agent.StandardiseWage(dw, refDay))
		if ni, ok := intensityByCountry[dw.CountryCode]; ok {
			relativePrices = append(relativePrices, agent.ComputeRelativePrice(dw, ni))
		}
	}

	cmp := agent.WageComparison{
		Countries:         countries,
		DayWages:          wages,
		StandardisedWages: standardised,
		RelativePrices:    relativePrices,
		SpindleRatios:     spindleRatios,
	}
	writeJSON(w, http.StatusOK, cmp)
}
