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

// ProductivityScope says how far a productivity gain reaches. Marx, Ch. 12:
// one capitalist's gain is only temporary extra surplus-value (individual vs
// social value); the general rate of surplus-value is affected "only when the
// increase in the productiveness of labour has seized upon those branches of
// production that are connected with ... the necessary means of subsistence".
type ProductivityScope string

const (
	// ScopeWorkplace is a single firm's gain (the default). It yields extra
	// surplus-value and leaves the value of labour-power untouched.
	ScopeWorkplace ProductivityScope = "workplace"
	// ScopeWageGoodsBranch is a productivity change across the wage-goods
	// branches; only this scope lowers the value of labour-power.
	ScopeWageGoodsBranch ProductivityScope = "wage_goods_branch"
)

func (s ProductivityScope) Valid() bool {
	return s == ScopeWorkplace || s == ScopeWageGoodsBranch
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
	// Scope defaults to ScopeWorkplace. SocialValue (per article, in labour
	// minutes) and Quantity are required for ScopeWorkplace only.
	Scope       ProductivityScope `json:"scope,omitempty"`
	SocialValue int64             `json:"social_value,omitempty"`
	Quantity    int64             `json:"quantity,omitempty"`
}

// relativeSurplusResponse is the bridge result: the Ch. 12 working-day
// reshaped by the productivity gain, plus the relative-SV delta.
type relativeSurplusResponse struct {
	ProductivityFactor   float64           `json:"productivity_factor"`
	OldWorkingDay        workingDayDTO     `json:"old_working_day"`
	NewWorkingDay        workingDayDTO     `json:"new_working_day"`
	NewLabourPowerValue  int64             `json:"new_labour_power_value"`
	RelativeSurplusValue int64             `json:"relative_surplus_value"`
	Scope                ProductivityScope `json:"scope"`
	// ExtraSurplusValue is the temporary individual-vs-social gain; set for
	// ScopeWorkplace only.
	ExtraSurplusValue *int64                 `json:"extra_surplus_value,omitempty"`
	Source            *productivitySourceDTO `json:"source,omitempty"`
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
	Scope           ProductivityScope  `json:"scope,omitempty"`
	SocialValue     int64              `json:"social_value,omitempty"`
	Quantity        int64              `json:"quantity,omitempty"`
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
		Scope:              req.Scope,
		SocialValue:        req.SocialValue,
		Quantity:           req.Quantity,
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
	scope := req.Scope
	if scope == "" {
		scope = ScopeWorkplace
	}
	if !scope.Valid() {
		return relativeSurplusResponse{}, errors.New("scope must be 'workplace' or 'wage_goods_branch'")
	}
	old := production.WorkingDay{
		Total:           labour.LabourMinutes(req.WorkingDayTotal),
		NecessaryLabour: production.NecessaryLabour(req.CurrentLPV),
		SurplusLabour:   production.SurplusLabour(req.WorkingDayTotal - req.CurrentLPV),
	}
	if scope == ScopeWorkplace {
		return computeWorkplaceGain(req, old, source)
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
		Scope:                ScopeWageGoodsBranch,
		Source:               source,
	}, nil
}

// computeWorkplaceGain handles a single firm's productivity gain: the value of
// labour-power and the working-day split stay as they were, and the gain is the
// temporary extra surplus-value between the firm's individual value and the
// social value of its article [Ch. 12 §1].
func computeWorkplaceGain(req relativeSurplusRequest, old production.WorkingDay, source *productivitySourceDTO) (relativeSurplusResponse, error) {
	if req.SocialValue <= 0 {
		return relativeSurplusResponse{}, errors.New("social_value must be positive for scope 'workplace'")
	}
	if req.Quantity <= 0 {
		return relativeSurplusResponse{}, errors.New("quantity must be positive for scope 'workplace'")
	}
	iv := production.ApplyProductivityToSNLT(
		labour.LabourMinutes(req.SocialValue),
		labour.ProductivityFactor(req.ProductivityFactor),
	)
	extra := int64(production.ExtraSurplusValue(
		production.IndividualValue(iv),
		production.SocialValue(req.SocialValue),
		production.Quantity(req.Quantity),
	))
	return relativeSurplusResponse{
		ProductivityFactor:   req.ProductivityFactor,
		OldWorkingDay:        workingDayDTOFromDomain(old),
		NewWorkingDay:        workingDayDTOFromDomain(old),
		NewLabourPowerValue:  req.CurrentLPV,
		RelativeSurplusValue: 0,
		Scope:                ScopeWorkplace,
		ExtraSurplusValue:    &extra,
		Source:               source,
	}, nil
}

func workingDayDTOFromDomain(wd production.WorkingDay) workingDayDTO {
	return dtoFromWorkingDay(wd)
}
