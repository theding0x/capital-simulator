package httpapi

import (
	"errors"
	"net/http"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 8 — Fixed Capital and Circulating Capital

// CreateComponent handles POST /v1/capital-components.
func (h *Handler) CreateComponent(w http.ResponseWriter, r *http.Request) {
	var c comp.CapitalComponent
	if err := decodeJSON(r, &c); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Composition.CreateComponent(r.Context(), c)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "component already exists")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Location", "/v1/capital-components/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// GetComponent handles GET /v1/capital-components/{id}.
func (h *Handler) GetComponent(w http.ResponseWriter, r *http.Request) {
	id := comp.CapitalComponentID(r.PathValue("id"))
	c, err := h.Composition.GetComponent(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "component not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// ListComponents handles GET /v1/capital-components.
// Optional query params: industrial_capital_id, kind, role.
func (h *Handler) ListComponents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	icID := q.Get("industrial_capital_id")
	kind := comp.CapitalKind(q.Get("kind"))
	role := comp.RoleInProcess(q.Get("role"))

	items, err := h.Composition.ListComponents(r.Context(), icID, kind, role)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// RegisterFixedItem handles POST /v1/fixed-items.
func (h *Handler) RegisterFixedItem(w http.ResponseWriter, r *http.Request) {
	var item comp.FixedCapitalItem
	if err := decodeJSON(r, &item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Composition.RegisterFixedItem(r.Context(), item)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "fixed item already exists")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Location", "/v1/fixed-items/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// GetFixedItem handles GET /v1/fixed-items/{id}.
func (h *Handler) GetFixedItem(w http.ResponseWriter, r *http.Request) {
	id := comp.FixedCapitalItemID(r.PathValue("id"))
	item, err := h.Composition.GetFixedItem(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "fixed item not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// RecordSubcomponent handles POST /v1/fixed-items/{id}/subcomponents.
func (h *Handler) RecordSubcomponent(w http.ResponseWriter, r *http.Request) {
	var sub comp.FixedCapitalSubcomponent
	if err := decodeJSON(r, &sub); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sub.FixedCapitalItemID = comp.FixedCapitalItemID(r.PathValue("id"))
	created, err := h.Composition.RecordSubcomponent(r.Context(), sub)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// RecordWear handles POST /v1/fixed-items/{id}/wear.
func (h *Handler) RecordWear(w http.ResponseWriter, r *http.Request) {
	var wear comp.WearAndTear
	if err := decodeJSON(r, &wear); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	wear.FixedCapitalItemID = comp.FixedCapitalItemID(r.PathValue("id"))
	recorded, err := h.Composition.RecordWear(r.Context(), wear)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, recorded)
}

// GetSinkingFund handles GET /v1/fixed-items/{id}/sinking-fund.
func (h *Handler) GetSinkingFund(w http.ResponseWriter, r *http.Request) {
	itemID := comp.FixedCapitalItemID(r.PathValue("id"))
	sf, err := h.Composition.GetSinkingFund(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "sinking fund not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sf)
}

// RecordRepair handles POST /v1/fixed-items/{id}/repairs.
func (h *Handler) RecordRepair(w http.ResponseWriter, r *http.Request) {
	var repair comp.Repair
	if err := decodeJSON(r, &repair); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	repair.FixedCapitalItemID = comp.FixedCapitalItemID(r.PathValue("id"))
	recorded, err := h.Composition.RecordRepair(r.Context(), repair)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, recorded)
}

// RecordCirculatingCycle handles POST /v1/circulating-cycles.
func (h *Handler) RecordCirculatingCycle(w http.ResponseWriter, r *http.Request) {
	var c comp.CirculatingCycle
	if err := decodeJSON(r, &c); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	recorded, err := h.Composition.RecordCirculatingCycle(r.Context(), c)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, recorded)
}

// RecordReinvestment handles POST /v1/fixed-items/{id}/reinvestments.
func (h *Handler) RecordReinvestment(w http.ResponseWriter, r *http.Request) {
	var reinvestment comp.SinkingFundReinvestment
	if err := decodeJSON(r, &reinvestment); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	reinvestment.FixedCapitalItemID = comp.FixedCapitalItemID(r.PathValue("id"))
	recorded, err := h.Composition.RecordReinvestment(r.Context(), reinvestment)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, recorded)
}
