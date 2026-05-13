package httpapi

import (
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

// Ch. 16 — Absolute and Relative Surplus-Value. All endpoints are
// stateless; they compute one period's partition given the inputs.

// workingDayPartitionDTO is the JSON shape of a Ch. 16 WorkingDay. We
// reuse `workingDayDTO` from production_handler.go for Ch. 12 results;
// the §1 partition uses bare LabourMinutes (no named field types) per
// the chapter spec.

type absoluteSurplusValueRequest struct {
	WorkingDayMinutes      int64 `json:"working_day_minutes"`
	NecessaryLabourMinutes int64 `json:"necessary_labour_minutes"`
	ExtensionMinutes       int64 `json:"extension_minutes"`
}

type absoluteSurplusValueResponse struct {
	OldWorkingDay        workingDayDTO `json:"old_working_day"`
	NewWorkingDay        workingDayDTO `json:"new_working_day"`
	AbsoluteSurplusValue surplusValueDTO `json:"absolute_surplus_value"`
}

type surplusValueDTO struct {
	LabourMinutes int64  `json:"labour_minutes"`
	Origin        string `json:"origin"`
}

// ComputeAbsoluteSurplusValue handles POST /v1/surplus-value/absolute.
// Prolongs the working day by extension_minutes, holding necessary labour
// fixed. The §1 mechanism for absolute surplus-value.
func (h *Handler) ComputeAbsoluteSurplusValue(w http.ResponseWriter, r *http.Request) {
	var req absoluteSurplusValueRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.WorkingDayMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "working_day_minutes must be positive")
		return
	}
	if req.NecessaryLabourMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour_minutes must be positive")
		return
	}
	if req.NecessaryLabourMinutes > req.WorkingDayMinutes {
		writeError(w, http.StatusBadRequest, "necessary_labour_minutes must be <= working_day_minutes")
		return
	}
	if req.ExtensionMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "extension_minutes must be positive")
		return
	}
	old := surplus.WorkingDay{
		Total:           labour.LabourMinutes(req.WorkingDayMinutes),
		NecessaryLabour: labour.LabourMinutes(req.NecessaryLabourMinutes),
		SurplusLabour:   labour.LabourMinutes(req.WorkingDayMinutes - req.NecessaryLabourMinutes),
	}
	newWD, asv := surplus.ProlongWorkingDay(old, labour.LabourMinutes(req.ExtensionMinutes))
	writeJSON(w, http.StatusOK, absoluteSurplusValueResponse{
		OldWorkingDay:        dtoFromSurplusWorkingDay(old),
		NewWorkingDay:        dtoFromSurplusWorkingDay(newWD),
		AbsoluteSurplusValue: dtoFromSurplusValue(asv.SurplusValue),
	})
}

type relativeSurplusValueRequest struct {
	WorkingDayMinutes      int64   `json:"working_day_minutes"`
	NecessaryLabourMinutes int64   `json:"necessary_labour_minutes"`
	ProductivityFactor     float64 `json:"productivity_factor"`
}

type relativeSurplusValueResponse struct {
	OldWorkingDay        workingDayDTO   `json:"old_working_day"`
	NewWorkingDay        workingDayDTO   `json:"new_working_day"`
	RelativeSurplusValue surplusValueDTO `json:"relative_surplus_value"`
}

// ComputeRelativeSurplusValue handles POST /v1/surplus-value/relative.
// Applies productivity_factor to necessary labour, holding the working
// day fixed. The §1 mechanism for relative surplus-value.
func (h *Handler) ComputeRelativeSurplusValue(w http.ResponseWriter, r *http.Request) {
	var req relativeSurplusValueRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.WorkingDayMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "working_day_minutes must be positive")
		return
	}
	if req.NecessaryLabourMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour_minutes must be positive")
		return
	}
	if req.NecessaryLabourMinutes > req.WorkingDayMinutes {
		writeError(w, http.StatusBadRequest, "necessary_labour_minutes must be <= working_day_minutes")
		return
	}
	if req.ProductivityFactor <= 0 {
		writeError(w, http.StatusBadRequest, "productivity_factor must be positive")
		return
	}
	old := surplus.WorkingDay{
		Total:           labour.LabourMinutes(req.WorkingDayMinutes),
		NecessaryLabour: labour.LabourMinutes(req.NecessaryLabourMinutes),
		SurplusLabour:   labour.LabourMinutes(req.WorkingDayMinutes - req.NecessaryLabourMinutes),
	}
	newWD, rsv := surplus.ReduceNecessaryLabour(old, surplus.ProductivityFactor(req.ProductivityFactor))
	writeJSON(w, http.StatusOK, relativeSurplusValueResponse{
		OldWorkingDay:        dtoFromSurplusWorkingDay(old),
		NewWorkingDay:        dtoFromSurplusWorkingDay(newWD),
		RelativeSurplusValue: dtoFromSurplusValue(rsv.SurplusValue),
	})
}

type surplusValueRateResponse struct {
	SurplusLabourMinutes   int64   `json:"surplus_labour_minutes"`
	NecessaryLabourMinutes int64   `json:"necessary_labour_minutes"`
	SurplusValueMinutes    int64   `json:"surplus_value_minutes"`
	TotalCapitalAdvanced   int64   `json:"total_capital_advanced"`
	RateOfSurplusValue     float64 `json:"rate_of_surplus_value"`
	RateOfProfit           float64 `json:"rate_of_profit"`
	MillCritiqueHolds      bool    `json:"mill_critique_holds"`
}

// GetSurplusValueRate handles GET /v1/surplus-value/rate. Surfaces Marx's
// correction of Mill: rate of profit is always lower than rate of
// surplus-value once constant capital is non-zero. Inputs are query
// parameters so the endpoint stays GET-able from a link.
//
//	?surplus_labour=…&necessary_labour=…&total_capital=…
//
// `surplus_value` defaults to `surplus_labour` when omitted (the
// magnitudes are identical in the §1 framing).
func (h *Handler) GetSurplusValueRate(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	surplusLabour, err := parsePositiveInt(q.Get("surplus_labour"), "surplus_labour")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	necessaryLabour, err := parsePositiveInt(q.Get("necessary_labour"), "necessary_labour")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	totalCapital, err := parsePositiveInt(q.Get("total_capital"), "total_capital")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	surplusValue := surplusLabour
	if s := q.Get("surplus_value"); s != "" {
		v, err := parsePositiveInt(s, "surplus_value")
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		surplusValue = v
	}
	rsv := surplus.RateSurplusValue(labour.LabourMinutes(surplusLabour), labour.LabourMinutes(necessaryLabour))
	rop := surplus.RateOfProfit(labour.LabourMinutes(surplusValue), labour.LabourMinutes(totalCapital))
	writeJSON(w, http.StatusOK, surplusValueRateResponse{
		SurplusLabourMinutes:   surplusLabour,
		NecessaryLabourMinutes: necessaryLabour,
		SurplusValueMinutes:    surplusValue,
		TotalCapitalAdvanced:   totalCapital,
		RateOfSurplusValue:     rsv,
		RateOfProfit:           rop,
		MillCritiqueHolds:      rop < rsv,
	})
}

func dtoFromSurplusWorkingDay(wd surplus.WorkingDay) workingDayDTO {
	return workingDayDTO{
		Total:           int64(wd.Total),
		NecessaryLabour: int64(wd.NecessaryLabour),
		SurplusLabour:   int64(wd.SurplusLabour),
	}
}

func dtoFromSurplusValue(sv surplus.SurplusValue) surplusValueDTO {
	return surplusValueDTO{
		LabourMinutes: int64(sv.LabourMinutes),
		Origin:        string(sv.Origin),
	}
}

func parsePositiveInt(s, field string) (int64, error) {
	if s == "" {
		return 0, errfmt(field + " is required")
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v <= 0 {
		return 0, errfmt(field + " must be a positive integer")
	}
	return v, nil
}

func errfmt(msg string) error { return &handlerError{msg: msg} }

type handlerError struct{ msg string }

func (e *handlerError) Error() string { return e.msg }
