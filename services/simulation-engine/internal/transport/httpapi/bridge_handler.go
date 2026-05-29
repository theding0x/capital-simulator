package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/production"
)

// ProductivitySource enumerates where a productivity factor can be sourced
// from inside the Part IV bridge endpoint. Each value corresponds to one
// of Marx's three forms by which capital realises relative surplus-value
// (Ch. 13 simple co-operation, Ch. 14 manufacture, Ch. 15 machinery).
type ProductivitySource string

const (
	SourceCooperation ProductivitySource = "cooperation"
	SourceManufacture ProductivitySource = "manufacture"
	SourceFactory     ProductivitySource = "factory"
)

func (s ProductivitySource) Valid() bool {
	switch s {
	case SourceCooperation, SourceManufacture, SourceFactory:
		return true
	}
	return false
}

// ProductivityFetcher resolves a productivity factor for a domain record
// in one of the three Part IV chapters. The default implementation makes
// HTTP calls (gateway-mediated); tests inject a stub.
type ProductivityFetcher interface {
	Fetch(ctx context.Context, source ProductivitySource, id string) (float64, error)
}

// ErrUnknownSource is returned when a request names a source the fetcher
// does not know how to resolve.
var ErrUnknownSource = errors.New("productivity: unknown source")

// relativeSurplusRequest is POST /v1/production/relative-surplus.
// Pass the factor directly when the caller has already computed it.
type relativeSurplusRequest struct {
	WorkingDayTotal    int64   `json:"working_day_total"`
	CurrentLPV         int64   `json:"current_lpv"`
	ProductivityFactor float64 `json:"productivity_factor"`
}

// relativeSurplusResponse is the bridge result: the Ch. 12 working-day
// reshaped by the productivity gain, plus the relative-SV delta.
type relativeSurplusResponse struct {
	ProductivityFactor   float64                  `json:"productivity_factor"`
	OldWorkingDay        workingDayDTO  `json:"old_working_day"`
	NewWorkingDay        workingDayDTO  `json:"new_working_day"`
	NewLabourPowerValue  int64                    `json:"new_labour_power_value"`
	RelativeSurplusValue int64                    `json:"relative_surplus_value"`
	Source               *productivitySourceDTO   `json:"source,omitempty"`
}

type productivitySourceDTO struct {
	Source ProductivitySource `json:"source"`
	ID     string             `json:"id"`
}

// RelativeSurplus handles POST /v1/production/relative-surplus.
// Direct: caller supplies the productivity factor. This is the
// workhorse the bridge endpoint composes on top of.
func (h *Handler) RelativeSurplus(w http.ResponseWriter, r *http.Request) {
	var req relativeSurplusRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := computeRelativeSurplus(req, nil)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// relativeSurplusFromProductivityRequest is the bridge endpoint body.
type relativeSurplusFromProductivityRequest struct {
	WorkingDayTotal int64              `json:"working_day_total"`
	CurrentLPV      int64              `json:"current_lpv"`
	Source          ProductivitySource `json:"source"`
	SourceID        string             `json:"source_id"`
}

// RelativeSurplusFromProductivity handles POST
// /v1/production/relative-surplus-from-productivity. It fetches the
// productivity factor from one of the Part IV domains (Ch.13 cooperation,
// Ch.14 manufacture, Ch.15 factory) and applies it to Ch.12's relative-SV
// formula. This is the end-to-end path that makes Marx's Part IV arc
// visible in the simulation.
func (h *Handler) RelativeSurplusFromProductivity(w http.ResponseWriter, r *http.Request) {
	var req relativeSurplusFromProductivityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !req.Source.Valid() {
		writeError(w, http.StatusBadRequest, "source must be 'cooperation', 'manufacture', or 'factory'")
		return
	}
	if req.SourceID == "" {
		writeError(w, http.StatusBadRequest, "source_id is required")
		return
	}
	if h.Productivity == nil {
		writeError(w, http.StatusServiceUnavailable, "productivity fetcher not configured")
		return
	}
	factor, err := h.Productivity.Fetch(r.Context(), req.Source, req.SourceID)
	if err != nil {
		if errors.Is(err, ErrUnknownSource) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	inner := relativeSurplusRequest{
		WorkingDayTotal:    req.WorkingDayTotal,
		CurrentLPV:         req.CurrentLPV,
		ProductivityFactor: factor,
	}
	src := &productivitySourceDTO{Source: req.Source, ID: req.SourceID}
	resp, err := computeRelativeSurplus(inner, src)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// computeRelativeSurplus does the §1 arithmetic. Shared by both bridge
// endpoints. `source` is attached to the response when set.
func computeRelativeSurplus(req relativeSurplusRequest, source *productivitySourceDTO) (relativeSurplusResponse, error) {
	if req.WorkingDayTotal <= 0 {
		return relativeSurplusResponse{}, errors.New("working_day_total must be positive")
	}
	if req.CurrentLPV <= 0 {
		return relativeSurplusResponse{}, errors.New("current_lpv must be positive")
	}
	if req.CurrentLPV >= req.WorkingDayTotal {
		return relativeSurplusResponse{}, errors.New("current_lpv must be less than working_day_total")
	}
	if req.ProductivityFactor <= 0 {
		return relativeSurplusResponse{}, errors.New("productivity_factor must be positive")
	}
	old := production.WorkingDay{
		Total:           labour.LabourMinutes(req.WorkingDayTotal),
		NecessaryLabour: production.NecessaryLabour(req.CurrentLPV),
		SurplusLabour:   production.SurplusLabour(req.WorkingDayTotal - req.CurrentLPV),
	}
	newLPV := production.ApplyProductivityToSNLT(
		labour.LabourMinutes(req.CurrentLPV),
		labour.ProductivityFactor(req.ProductivityFactor),
	)
	// Guard: relative SV requires newLPV < total, else there is no working day to shorten.
	if int64(newLPV) >= req.WorkingDayTotal {
		return relativeSurplusResponse{}, errors.New("productivity gain is too small to reduce labour-power value below the working day")
	}
	newWD := production.ShortenNecessaryLabour(old, production.LabourPowerValue(newLPV))
	delta := int64(newWD.SurplusLabour) - int64(old.SurplusLabour)
	return relativeSurplusResponse{
		ProductivityFactor:   req.ProductivityFactor,
		OldWorkingDay:        workingDayDTOFromDomain(old),
		NewWorkingDay:        workingDayDTOFromDomain(newWD),
		NewLabourPowerValue:  int64(newLPV),
		RelativeSurplusValue: delta,
		Source:               source,
	}, nil
}

func workingDayDTOFromDomain(wd production.WorkingDay) workingDayDTO {
	return dtoFromWorkingDay(wd)
}
