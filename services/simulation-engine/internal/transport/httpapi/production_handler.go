package httpapi

import (
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/production"
)

// workingDayDTO is the shared JSON shape for a WorkingDay.
type workingDayDTO struct {
	Total           int64 `json:"total"`
	NecessaryLabour int64 `json:"necessary_labour"`
	SurplusLabour   int64 `json:"surplus_labour"`
}

func dtoFromWorkingDay(wd production.WorkingDay) workingDayDTO {
	return workingDayDTO{
		Total:           int64(wd.Total),
		NecessaryLabour: int64(wd.NecessaryLabour),
		SurplusLabour:   int64(wd.SurplusLabour),
	}
}

// recordWorkingDayRequest is the POST /v1/production/working-day body.
type recordWorkingDayRequest struct {
	Total            int64 `json:"total"`
	LabourPowerValue int64 `json:"labour_power_value"`
}

// RecordWorkingDay handles POST /v1/production/working-day.
// Builds a WorkingDay from total and labour-power value (both in LabourMinutes).
func (h *Handler) RecordWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req recordWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Total <= 0 {
		writeError(w, http.StatusBadRequest, "total must be positive")
		return
	}
	if req.LabourPowerValue <= 0 {
		writeError(w, http.StatusBadRequest, "labour_power_value must be positive")
		return
	}
	if req.LabourPowerValue >= req.Total {
		writeError(w, http.StatusBadRequest, "labour_power_value must be less than total")
		return
	}
	wd := production.WorkingDay{
		Total:           production.LabourMinutes(req.Total),
		NecessaryLabour: production.NecessaryLabour(req.LabourPowerValue),
		SurplusLabour:   production.SurplusLabour(req.Total - req.LabourPowerValue),
	}
	writeJSON(w, http.StatusOK, dtoFromWorkingDay(wd))
}

// shortenWorkingDayRequest is the POST /v1/production/working-day/shorten body.
type shortenWorkingDayRequest struct {
	Total               int64 `json:"total"`
	NecessaryLabour     int64 `json:"necessary_labour"`
	SurplusLabour       int64 `json:"surplus_labour"`
	NewLabourPowerValue int64 `json:"new_labour_power_value"`
}

// shortenWorkingDayResponse is the POST /v1/production/working-day/shorten response.
type shortenWorkingDayResponse struct {
	WorkingDay           workingDayDTO `json:"working_day"`
	RelativeSurplusValue int64         `json:"relative_surplus_value"`
}

// ShortenWorkingDay handles POST /v1/production/working-day/shorten.
// Returns the updated WorkingDay and the RelativeSurplusValue delta (new SL − old SL).
func (h *Handler) ShortenWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req shortenWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Total <= 0 {
		writeError(w, http.StatusBadRequest, "total must be positive")
		return
	}
	if req.NecessaryLabour <= 0 || req.SurplusLabour < 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour must be positive; surplus_labour cannot be negative")
		return
	}
	if req.NecessaryLabour+req.SurplusLabour != req.Total {
		writeError(w, http.StatusBadRequest, "necessary_labour + surplus_labour must equal total")
		return
	}
	if req.NewLabourPowerValue <= 0 {
		writeError(w, http.StatusBadRequest, "new_labour_power_value must be positive")
		return
	}
	if req.NewLabourPowerValue >= req.Total {
		writeError(w, http.StatusBadRequest, "new_labour_power_value must be less than total")
		return
	}
	wd := production.WorkingDay{
		Total:           production.LabourMinutes(req.Total),
		NecessaryLabour: production.NecessaryLabour(req.NecessaryLabour),
		SurplusLabour:   production.SurplusLabour(req.SurplusLabour),
	}
	newWD := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(req.NewLabourPowerValue))
	delta := int64(newWD.SurplusLabour) - int64(wd.SurplusLabour)
	writeJSON(w, http.StatusOK, shortenWorkingDayResponse{
		WorkingDay:           dtoFromWorkingDay(newWD),
		RelativeSurplusValue: delta,
	})
}

// rateResponse is the GET /v1/production/rate-of-surplus-value response.
type rateResponse struct {
	Rate            float64 `json:"rate"`
	SurplusLabour   int64   `json:"surplus_labour"`
	NecessaryLabour int64   `json:"necessary_labour"`
}

// GetProductionRate handles GET /v1/production/rate-of-surplus-value?necessary=&surplus=.
func (h *Handler) GetProductionRate(w http.ResponseWriter, r *http.Request) {
	necessaryStr := r.URL.Query().Get("necessary")
	surplusStr := r.URL.Query().Get("surplus")
	if necessaryStr == "" || surplusStr == "" {
		writeError(w, http.StatusBadRequest, "necessary and surplus query params required")
		return
	}
	necessary, err1 := strconv.ParseInt(necessaryStr, 10, 64)
	surplus, err2 := strconv.ParseInt(surplusStr, 10, 64)
	if err1 != nil || err2 != nil || necessary <= 0 || surplus < 0 {
		writeError(w, http.StatusBadRequest, "necessary must be a positive integer; surplus must be a non-negative integer")
		return
	}
	rate := production.RateOfSurplusValue(
		production.SurplusLabour(surplus),
		production.NecessaryLabour(necessary),
	)
	writeJSON(w, http.StatusOK, rateResponse{
		Rate:            rate,
		SurplusLabour:   surplus,
		NecessaryLabour: necessary,
	})
}

// extraSurplusRequest is the POST /v1/production/extra-surplus-value body.
type extraSurplusRequest struct {
	IndividualValue int64 `json:"individual_value"`
	SocialValue     int64 `json:"social_value"`
	Quantity        int64 `json:"quantity"`
}

// extraSurplusResponse is the POST /v1/production/extra-surplus-value response.
type extraSurplusResponse struct {
	ExtraSurplusValue int64 `json:"extra_surplus_value"`
	PerUnit           int64 `json:"per_unit"`
	IndividualValue   int64 `json:"individual_value"`
	SocialValue       int64 `json:"social_value"`
	Quantity          int64 `json:"quantity"`
}

// ComputeExtraSurplusValue handles POST /v1/production/extra-surplus-value.
func (h *Handler) ComputeExtraSurplusValue(w http.ResponseWriter, r *http.Request) {
	var req extraSurplusRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndividualValue <= 0 {
		writeError(w, http.StatusBadRequest, "individual_value must be positive")
		return
	}
	if req.SocialValue <= 0 {
		writeError(w, http.StatusBadRequest, "social_value must be positive")
		return
	}
	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "quantity must be positive")
		return
	}
	total := production.ExtraSurplusValue(
		production.IndividualValue(req.IndividualValue),
		production.SocialValue(req.SocialValue),
		production.Quantity(req.Quantity),
	)
	perUnit := int64(0)
	if req.SocialValue > req.IndividualValue {
		perUnit = req.SocialValue - req.IndividualValue
	}
	writeJSON(w, http.StatusOK, extraSurplusResponse{
		ExtraSurplusValue: int64(total),
		PerUnit:           perUnit,
		IndividualValue:   req.IndividualValue,
		SocialValue:       req.SocialValue,
		Quantity:          req.Quantity,
	})
}
