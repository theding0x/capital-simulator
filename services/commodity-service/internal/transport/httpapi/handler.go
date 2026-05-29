// Package httpapi exposes commodity-service over HTTP. Routes implement the
// CRUD surface for commodities plus the value/exchange/value-form/social-
// relations endpoints required by Capital Vol. I, Ch. 1.
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/store"
)

// Handler bundles the dependencies for the HTTP layer.
type Handler struct {
	Store                  store.Store
	ProductionAccountStore store.ProductionAccountStore
	Logger                 *slog.Logger
}

// New constructs a Handler.
func New(s store.Store, pas store.ProductionAccountStore, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{Store: s, ProductionAccountStore: pas, Logger: logger}
}

// --- request/response types -------------------------------------------------

type createRequest struct {
	Name           string                  `json:"name"`
	UseValue       commodity.UseValue      `json:"use_value"`
	ConcreteLabour commodity.ConcreteLabour `json:"concrete_labour"`
	SNLTPerUnit    commodity.LabourMinutes `json:"snlt_per_unit"`
}

func (r createRequest) toCommodity() commodity.Commodity {
	return commodity.Commodity{
		Name:           strings.TrimSpace(r.Name),
		UseValue:       r.UseValue,
		ConcreteLabour: r.ConcreteLabour,
		SNLTPerUnit:    r.SNLTPerUnit,
	}
}

type updateRequest struct {
	Name                      *string                  `json:"name,omitempty"`
	UseValueDescription       *string                  `json:"use_value_description,omitempty"`
	UseValueUnit              *string                  `json:"use_value_unit,omitempty"`
	ConcreteLabourKind        *string                  `json:"concrete_labour_kind,omitempty"`
	ConcreteLabourDescription *string                  `json:"concrete_labour_description,omitempty"`
	SNLTPerUnit               *commodity.LabourMinutes `json:"snlt_per_unit,omitempty"`
}

func (r updateRequest) toUpdate() store.Update {
	return store.Update{
		Name:                      r.Name,
		UseValueDescription:       r.UseValueDescription,
		UseValueUnit:              r.UseValueUnit,
		ConcreteLabourKind:        r.ConcreteLabourKind,
		ConcreteLabourDescription: r.ConcreteLabourDescription,
		SNLTPerUnit:               r.SNLTPerUnit,
	}
}

type valueRequest struct {
	Quantity commodity.Quantity `json:"quantity"`
}

type valueResponse struct {
	Commodity commodity.Commodity     `json:"commodity"`
	Quantity  commodity.Quantity      `json:"quantity"`
	Value     commodity.LabourMinutes `json:"value"`
	ValueHours float64                `json:"value_hours"`
}

type exchangeRequest struct {
	BaseID  commodity.ID       `json:"base_id"`
	QuoteID commodity.ID       `json:"quote_id"`
	BaseQty commodity.Quantity `json:"base_qty"`
}

// --- handlers ---------------------------------------------------------------

// Create handles POST /v1/commodities.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	c := req.toCommodity()
	if err := c.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.Create(r.Context(), c)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// List handles GET /v1/commodities.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.Store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []commodity.Commodity{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// Get handles GET /v1/commodities/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := commodity.ID(r.PathValue("id"))
	c, err := h.Store.Get(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// Update handles PATCH /v1/commodities/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := commodity.ID(r.PathValue("id"))
	var req updateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	c, err := h.Store.Update(r.Context(), id, req.toUpdate())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// Delete handles DELETE /v1/commodities/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := commodity.ID(r.PathValue("id"))
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Value handles POST /v1/commodities/{id}/value - compute the value (in
// abstract human labour-minutes) of a given quantity of the commodity.
func (h *Handler) Value(w http.ResponseWriter, r *http.Request) {
	id := commodity.ID(r.PathValue("id"))
	var req valueRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	c, err := h.Store.Get(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	v := c.Value(req.Quantity)
	writeJSON(w, http.StatusOK, valueResponse{
		Commodity: c, Quantity: req.Quantity, Value: v, ValueHours: v.Hours(),
	})
}

// ExchangeRatio handles POST /v1/exchange-ratio.
func (h *Handler) ExchangeRatio(w http.ResponseWriter, r *http.Request) {
	var req exchangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.BaseQty <= 0 {
		req.BaseQty = 1
	}
	base, err := h.Store.Get(r.Context(), req.BaseID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	quote, err := h.Store.Get(r.Context(), req.QuoteID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	ratio, err := commodity.ExchangeRatioBetween(base, quote, req.BaseQty)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ratio)
}

// ValueForm handles GET /v1/commodities/{id}/value-form?kind=...&money_id=...
func (h *Handler) ValueForm(w http.ResponseWriter, r *http.Request) {
	id := commodity.ID(r.PathValue("id"))
	kind := commodity.ValueFormKind(strings.ToLower(r.URL.Query().Get("kind")))
	if kind == "" {
		kind = commodity.KindExpanded
	}
	if err := kind.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	subject, err := h.Store.Get(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	pop, err := h.Store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch kind {
	case commodity.KindSimple:
		quoteID := commodity.ID(r.URL.Query().Get("quote_id"))
		if quoteID.IsZero() {
			writeError(w, http.StatusBadRequest, "quote_id is required for simple form")
			return
		}
		quote, err := h.Store.Get(r.Context(), quoteID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		sf, err := commodity.SimpleFormOf(subject, quote)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"kind": kind, "form": sf})

	case commodity.KindExpanded:
		ef, err := commodity.ExpandedFormOf(subject, pop)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"kind": kind, "form": ef})

	case commodity.KindGeneral:
		gf, err := commodity.GeneralFormOf(subject, pop)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"kind": kind, "form": gf})

	case commodity.KindMoney:
		moneyID := commodity.ID(r.URL.Query().Get("money_id"))
		if moneyID.IsZero() {
			moneyID = subject.ID
		}
		money, err := h.Store.Get(r.Context(), moneyID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		mf, err := commodity.MoneyFormOf(money, pop)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"kind": kind, "form": mf})
	}
}

// SocialRelations handles GET /v1/commodities/{id}/social-relations - the
// un-mystified view that surfaces the labour relations behind exchange.
func (h *Handler) SocialRelations(w http.ResponseWriter, r *http.Request) {
	id := commodity.ID(r.PathValue("id"))
	subject, err := h.Store.Get(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	pop, err := h.Store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, commodity.SocialRelationsOf(subject, pop))
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
