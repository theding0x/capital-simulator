// ch03_money_handler.go — HTTP handlers for Capital Vol. I, Ch. 3 money-form
// records: hoards, payment obligations, world-money transfers, and the
// quantity-of-money formula M = ΣP / V.
package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
)

// --- request types ----------------------------------------------------------

type createHoardRequest struct {
	OwnerID market.OwnerID `json:"owner_id"`
	Amount  market.Pence   `json:"amount"`
}

type createPaymentObligationRequest struct {
	CreditorID market.OwnerID `json:"creditor_id"`
	DebtorID   market.OwnerID `json:"debtor_id"`
	Amount     market.Pence   `json:"amount"`
	DueAt      time.Time      `json:"due_at"`
}

type createWorldMoneyTransferRequest struct {
	SenderID   market.OwnerID `json:"sender_id"`
	ReceiverID market.OwnerID `json:"receiver_id"`
	GoldMg     int64          `json:"gold_mg"`
}

// --- Hoard handlers ---------------------------------------------------------

// CreateHoard handles POST /v1/hoards.
func (h *Handler) CreateHoard(w http.ResponseWriter, r *http.Request) {
	var req createHoardRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.CreateHoard(r.Context(), market.Hoard{
		OwnerID: req.OwnerID,
		Amount:  req.Amount,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListHoards handles GET /v1/hoards.
func (h *Handler) ListHoards(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListHoards(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []market.Hoard{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetHoard handles GET /v1/hoards/{id}.
func (h *Handler) GetHoard(w http.ResponseWriter, r *http.Request) {
	id := market.HoardID(r.PathValue("id"))
	hoard, err := h.Store.GetHoard(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, hoard)
}

// --- PaymentObligation handlers --------------------------------------------

// CreatePaymentObligation handles POST /v1/payment-obligations.
func (h *Handler) CreatePaymentObligation(w http.ResponseWriter, r *http.Request) {
	var req createPaymentObligationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.CreatePaymentObligation(r.Context(), market.PaymentObligation{
		CreditorID: req.CreditorID,
		DebtorID:   req.DebtorID,
		Amount:     req.Amount,
		DueAt:      req.DueAt,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListPaymentObligations handles GET /v1/payment-obligations.
func (h *Handler) ListPaymentObligations(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListPaymentObligations(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []market.PaymentObligation{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// SettlePaymentObligation handles POST /v1/payment-obligations/{id}/settle.
func (h *Handler) SettlePaymentObligation(w http.ResponseWriter, r *http.Request) {
	id := market.PaymentObligationID(r.PathValue("id"))
	settled, err := h.Store.SettlePaymentObligation(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, settled)
}

// --- WorldMoneyTransfer handlers -------------------------------------------

// CreateWorldMoneyTransfer handles POST /v1/world-money-transfers.
func (h *Handler) CreateWorldMoneyTransfer(w http.ResponseWriter, r *http.Request) {
	var req createWorldMoneyTransferRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.CreateWorldMoneyTransfer(r.Context(), market.WorldMoneyTransfer{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		GoldMg:     req.GoldMg,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListWorldMoneyTransfers handles GET /v1/world-money-transfers.
func (h *Handler) ListWorldMoneyTransfers(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListWorldMoneyTransfers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []market.WorldMoneyTransfer{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// --- Circulation / MoneyRequired -------------------------------------------

// MoneyRequired handles GET /v1/circulation/money-required?sum_of_prices=&velocity=.
// Implements Marx Vol. I Ch. 3 §2.b — M = ΣP / V.
func (h *Handler) MoneyRequired(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sumOfPrices, err := strconv.ParseInt(q.Get("sum_of_prices"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "sum_of_prices must be an integer")
		return
	}
	velocity, err := strconv.ParseInt(q.Get("velocity"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "velocity must be an integer")
		return
	}
	moneyRequired, err := market.MoneyRequired(sumOfPrices, velocity)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"money_required": moneyRequired})
}
