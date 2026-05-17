package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 31 — Genesis of the Industrial Capitalist.

type createCapitalOriginRequest struct {
	Source      string `json:"source"`
	AmountPence int64  `json:"amount_pence"`
	Period      string `json:"period"`
}

type capitalOriginResponse struct {
	ID                string `json:"id"`
	HistoricalStageID string `json:"historical_stage_id"`
	Source            string `json:"source"`
	AmountPence       int64  `json:"amount_pence"`
	Period            string `json:"period"`
	CreatedAt         string `json:"created_at"`
}

type createColonialTransferRequest struct {
	From       string `json:"from"`
	To         string `json:"to"`
	ValuePence int64  `json:"value_pence"`
	Method     string `json:"method"`
}

type colonialTransferResponse struct {
	ID                string `json:"id"`
	HistoricalStageID string `json:"historical_stage_id"`
	From              string `json:"from"`
	To                string `json:"to"`
	ValuePence        int64  `json:"value_pence"`
	Method            string `json:"method"`
	CreatedAt         string `json:"created_at"`
}

type createNationalDebtRequest struct {
	AmountPence     int64  `json:"amount_pence"`
	InterestRateBps int64  `json:"interest_rate_bps"`
	CreditorClass   string `json:"creditor_class"`
}

type nationalDebtResponse struct {
	ID                string `json:"id"`
	HistoricalStageID string `json:"historical_stage_id"`
	AmountPence       int64  `json:"amount_pence"`
	InterestRateBps   int64  `json:"interest_rate_bps"`
	CreditorClass     string `json:"creditor_class"`
	CreatedAt         string `json:"created_at"`
}

type protectionSystemResponse struct {
	ID                string `json:"id"`
	HistoricalStageID string `json:"historical_stage_id"`
	TariffRateBps     int64  `json:"tariff_rate_bps"`
	Beneficiary       string `json:"beneficiary"`
	PeriodStart       string `json:"period_start"`
	PeriodEnd         string `json:"period_end"`
	CreatedAt         string `json:"created_at"`
}

type industrialCapitalGenesisResponse struct {
	HistoricalStageID       string                     `json:"historical_stage_id"`
	Origins                 []capitalOriginResponse    `json:"origins"`
	ColonialTransfers       []colonialTransferResponse `json:"colonial_transfers"`
	NationalDebts           []nationalDebtResponse     `json:"national_debts"`
	ProtectionSystems       []protectionSystemResponse `json:"protection_systems"`
	TotalCapitalFormedPence int64                      `json:"total_capital_formed_pence"`
}

func toCapitalOriginResponse(c simulation.CapitalOrigin) capitalOriginResponse {
	return capitalOriginResponse{
		ID:                string(c.ID),
		HistoricalStageID: string(c.HistoricalStageID),
		Source:            c.Source,
		AmountPence:       int64(c.AmountPence),
		Period:            c.Period,
		CreatedAt:         c.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

func toColonialTransferResponse(t simulation.ColonialTransfer) colonialTransferResponse {
	return colonialTransferResponse{
		ID:                string(t.ID),
		HistoricalStageID: string(t.HistoricalStageID),
		From:              t.From,
		To:                t.To,
		ValuePence:        int64(t.ValuePence),
		Method:            t.Method,
		CreatedAt:         t.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

func toNationalDebtResponse(d simulation.NationalDebt) nationalDebtResponse {
	return nationalDebtResponse{
		ID:                string(d.ID),
		HistoricalStageID: string(d.HistoricalStageID),
		AmountPence:       int64(d.AmountPence),
		InterestRateBps:   d.InterestRateBps,
		CreditorClass:     d.CreditorClass,
		CreatedAt:         d.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

func toProtectionSystemResponse(s simulation.ProtectionSystem) protectionSystemResponse {
	return protectionSystemResponse{
		ID:                string(s.ID),
		HistoricalStageID: string(s.HistoricalStageID),
		TariffRateBps:     s.TariffRateBps,
		Beneficiary:       s.Beneficiary,
		PeriodStart:       s.PeriodStart,
		PeriodEnd:         s.PeriodEnd,
		CreatedAt:         s.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

// CreateCapitalOrigin handles POST /v1/historical-stages/{id}/capital-origins.
func (h *Handler) CreateCapitalOrigin(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))
	if stageID.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if _, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "historical stage not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req createCapitalOriginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	c := simulation.CapitalOrigin{
		HistoricalStageID: stageID,
		Source:            req.Source,
		AmountPence:       simulation.Pence(req.AmountPence),
		Period:            req.Period,
	}
	if err := c.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.CapitalOrigins.CreateCapitalOrigin(r.Context(), c)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "capital origin already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/historical-stages/%s/capital-origins", string(stageID)))
	writeJSON(w, http.StatusCreated, toCapitalOriginResponse(created))
}

// CreateColonialTransfer handles POST /v1/historical-stages/{id}/colonial-transfers.
func (h *Handler) CreateColonialTransfer(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))
	if stageID.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if _, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "historical stage not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req createColonialTransferRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	t := simulation.ColonialTransfer{
		HistoricalStageID: stageID,
		From:              req.From,
		To:                req.To,
		ValuePence:        simulation.Pence(req.ValuePence),
		Method:            req.Method,
	}
	if err := t.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.ColonialTransfers.CreateColonialTransfer(r.Context(), t)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "colonial transfer already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/historical-stages/%s/colonial-transfers", string(stageID)))
	writeJSON(w, http.StatusCreated, toColonialTransferResponse(created))
}

// CreateNationalDebt handles POST /v1/historical-stages/{id}/national-debts.
func (h *Handler) CreateNationalDebt(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))
	if stageID.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if _, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "historical stage not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req createNationalDebtRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	d := simulation.NationalDebt{
		HistoricalStageID: stageID,
		AmountPence:       simulation.Pence(req.AmountPence),
		InterestRateBps:   req.InterestRateBps,
		CreditorClass:     req.CreditorClass,
	}
	if err := d.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.NationalDebts.CreateNationalDebt(r.Context(), d)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "national debt already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/historical-stages/%s/national-debts", string(stageID)))
	writeJSON(w, http.StatusCreated, toNationalDebtResponse(created))
}

// GetIndustrialCapitalGenesis handles GET /v1/historical-stages/{id}/genesis.
// Returns the full IndustrialCapitalGenesis summary for the given stage.
func (h *Handler) GetIndustrialCapitalGenesis(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))
	if stageID.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	_, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "historical stage not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	origins, err := h.CapitalOrigins.ListCapitalOriginsByStage(r.Context(), stageID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	transfers, err := h.ColonialTransfers.ListColonialTransfersByStage(r.Context(), stageID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	debts, err := h.NationalDebts.ListNationalDebtsByStage(r.Context(), stageID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	systems, err := h.ProtectionSystems.ListProtectionSystemsByStage(r.Context(), stageID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	genesis := simulation.ComputeGenesis(stageID, origins, transfers, debts, systems)

	originResps := make([]capitalOriginResponse, len(genesis.Origins))
	for i, o := range genesis.Origins {
		originResps[i] = toCapitalOriginResponse(o)
	}
	transferResps := make([]colonialTransferResponse, len(genesis.ColonialTransfers))
	for i, t := range genesis.ColonialTransfers {
		transferResps[i] = toColonialTransferResponse(t)
	}
	debtResps := make([]nationalDebtResponse, len(genesis.NationalDebts))
	for i, d := range genesis.NationalDebts {
		debtResps[i] = toNationalDebtResponse(d)
	}
	systemResps := make([]protectionSystemResponse, len(genesis.ProtectionSystems))
	for i, s := range genesis.ProtectionSystems {
		systemResps[i] = toProtectionSystemResponse(s)
	}

	writeJSON(w, http.StatusOK, industrialCapitalGenesisResponse{
		HistoricalStageID:       string(genesis.HistoricalStageID),
		Origins:                 originResps,
		ColonialTransfers:       transferResps,
		NationalDebts:           debtResps,
		ProtectionSystems:       systemResps,
		TotalCapitalFormedPence: int64(genesis.TotalCapitalFormedPence),
	})
}
