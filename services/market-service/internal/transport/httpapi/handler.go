// Package httpapi exposes market-service over HTTP. Routes implement the
// surface for owners, offers, exchanges, the universal equivalent, the
// money-commodity, and prices — the domain of Capital Vol. I, Ch. 2.
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
)

// Handler bundles the dependencies for the HTTP layer.
type Handler struct {
	Store  store.Store
	Logger *slog.Logger
}

// New constructs a Handler.
func New(s store.Store, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{Store: s, Logger: logger}
}

// --- request / response types -----------------------------------------------

type createOwnerRequest struct {
	Name string `json:"name"`
}

type createOfferRequest struct {
	OwnerID          market.OwnerID     `json:"owner_id"`
	CommodityID      market.CommodityID `json:"commodity_id"`
	Quantity         float64            `json:"quantity"`
	SeeksKind        string             `json:"seeks_kind"`
	SeeksCommodityID market.CommodityID `json:"seeks_commodity_id,omitempty"`
}

type createExchangeRequest struct {
	GiverID             market.OwnerID       `json:"giver_id"`
	ReceiverID          market.OwnerID       `json:"receiver_id"`
	GiverCommodityID    market.CommodityID   `json:"giver_commodity_id"`
	GiverQty            float64              `json:"giver_qty"`
	ReceiverCommodityID market.CommodityID   `json:"receiver_commodity_id"`
	ReceiverQty         float64              `json:"receiver_qty"`
	RealisedValue       market.RealisedValue `json:"realised_value"`
}

type setUniversalEquivalentRequest struct {
	CommodityID market.CommodityID `json:"commodity_id"`
}

type setMoneyCommodityRequest struct {
	CommodityID market.CommodityID `json:"commodity_id"`
}

type computePriceRequest struct {
	CommodityID      market.CommodityID `json:"commodity_id"`
	MoneyCommodityID market.CommodityID `json:"money_commodity_id"`
	CommoditySNLT    int64              `json:"commodity_snlt"`
	MoneySNLT        int64              `json:"money_snlt"`
	UnitQty          int64              `json:"unit_qty"`
}

// --- Owner handlers ---------------------------------------------------------

// CreateOwner handles POST /v1/owners.
func (h *Handler) CreateOwner(w http.ResponseWriter, r *http.Request) {
	var req createOwnerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	o := market.Owner{Name: req.Name}
	created, err := h.Store.CreateOwner(r.Context(), o)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListOwners handles GET /v1/owners.
func (h *Handler) ListOwners(w http.ResponseWriter, r *http.Request) {
	out, err := h.Store.ListOwners(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []market.Owner{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// GetOwner handles GET /v1/owners/{id}.
func (h *Handler) GetOwner(w http.ResponseWriter, r *http.Request) {
	id := market.OwnerID(r.PathValue("id"))
	o, err := h.Store.GetOwner(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// --- Offer handlers ---------------------------------------------------------

// CreateOffer handles POST /v1/offers.
func (h *Handler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var req createOfferRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	o := market.Offer{
		OwnerID:          req.OwnerID,
		CommodityID:      req.CommodityID,
		Quantity:         req.Quantity,
		SeeksKind:        req.SeeksKind,
		SeeksCommodityID: req.SeeksCommodityID,
	}
	created, err := h.Store.CreateOffer(r.Context(), o)
	if err != nil {
		if errors.Is(err, market.ErrOfferInvalid) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListOffers handles GET /v1/offers.
func (h *Handler) ListOffers(w http.ResponseWriter, r *http.Request) {
	out, err := h.Store.ListOffers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []market.Offer{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// DeleteOffer handles DELETE /v1/offers/{id}.
func (h *Handler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	id := market.OfferID(r.PathValue("id"))
	if err := h.Store.DeleteOffer(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Exchange handlers ------------------------------------------------------

// CreateExchange handles POST /v1/exchanges.
func (h *Handler) CreateExchange(w http.ResponseWriter, r *http.Request) {
	var req createExchangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	e := market.Exchange{
		GiverID:             req.GiverID,
		ReceiverID:          req.ReceiverID,
		GiverCommodityID:    req.GiverCommodityID,
		GiverQty:            req.GiverQty,
		ReceiverCommodityID: req.ReceiverCommodityID,
		ReceiverQty:         req.ReceiverQty,
		RealisedValue:       req.RealisedValue,
	}
	created, err := h.Store.CreateExchange(r.Context(), e)
	if err != nil {
		if errors.Is(err, market.ErrSelfExchange) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListExchanges handles GET /v1/exchanges.
func (h *Handler) ListExchanges(w http.ResponseWriter, r *http.Request) {
	out, err := h.Store.ListExchanges(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []market.Exchange{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// GetExchange handles GET /v1/exchanges/{id}.
func (h *Handler) GetExchange(w http.ResponseWriter, r *http.Request) {
	id := market.ExchangeID(r.PathValue("id"))
	e, err := h.Store.GetExchange(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, e)
}

// --- Universal equivalent handler -------------------------------------------

// SetUniversalEquivalent handles POST /v1/universal-equivalent.
func (h *Handler) SetUniversalEquivalent(w http.ResponseWriter, r *http.Request) {
	var req setUniversalEquivalentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CommodityID.IsZero() {
		writeError(w, http.StatusBadRequest, "commodity_id is required")
		return
	}
	ue, err := h.Store.SetUniversalEquivalent(r.Context(), market.UniversalEquivalent{CommodityID: req.CommodityID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ue)
}

// GetUniversalEquivalent handles GET /v1/universal-equivalent.
func (h *Handler) GetUniversalEquivalent(w http.ResponseWriter, r *http.Request) {
	ue, err := h.Store.GetUniversalEquivalent(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ue)
}

// --- Money commodity handlers -----------------------------------------------

// SetMoneyCommodity handles POST /v1/money-commodity.
func (h *Handler) SetMoneyCommodity(w http.ResponseWriter, r *http.Request) {
	var req setMoneyCommodityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CommodityID.IsZero() {
		writeError(w, http.StatusBadRequest, "commodity_id is required")
		return
	}
	mc, err := h.Store.SetMoneyCommodity(r.Context(), market.MoneyCommodity{CommodityID: req.CommodityID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mc)
}

// GetMoneyCommodity handles GET /v1/money-commodity.
func (h *Handler) GetMoneyCommodity(w http.ResponseWriter, r *http.Request) {
	mc, err := h.Store.GetMoneyCommodity(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mc)
}

// --- Price handlers ---------------------------------------------------------

// ComputePrice handles POST /v1/prices — compute and store a price.
func (h *Handler) ComputePrice(w http.ResponseWriter, r *http.Request) {
	var req computePriceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CommodityID.IsZero() || req.MoneyCommodityID.IsZero() {
		writeError(w, http.StatusBadRequest, "commodity_id and money_commodity_id are required")
		return
	}
	if req.UnitQty == 0 {
		req.UnitQty = 1
	}
	amount, err := market.ComputePrice(req.CommoditySNLT, req.MoneySNLT, req.UnitQty)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, err := h.Store.SetPrice(r.Context(), market.Price{
		CommodityID:      req.CommodityID,
		MoneyCommodityID: req.MoneyCommodityID,
		Amount:           amount,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

// GetPrice handles GET /v1/prices/{commodityID}.
func (h *Handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	id := market.CommodityID(r.PathValue("commodityID"))
	p, err := h.Store.GetPrice(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// ListPrices handles GET /v1/prices.
func (h *Handler) ListPrices(w http.ResponseWriter, r *http.Request) {
	out, err := h.Store.ListPrices(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []market.Price{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// --- helpers ----------------------------------------------------------------

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return errors.New("invalid json: " + err.Error())
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
