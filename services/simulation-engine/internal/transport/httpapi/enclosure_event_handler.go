package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// Ch. 27 — Expropriation of the Agricultural Population.

type enclosureEventResponse struct {
	ID                  string `json:"id"`
	Period              string `json:"period"`
	AcresEnclosed       int64  `json:"acres_enclosed"`
	PopulationDisplaced int64  `json:"population_displaced"`
	Beneficiary         string `json:"beneficiary"`
	CreatedAt           string `json:"created_at"`
}

type createEnclosureEventRequest struct {
	Period              string `json:"period"`
	AcresEnclosed       int64  `json:"acres_enclosed"`
	PopulationDisplaced int64  `json:"population_displaced"`
	Beneficiary         string `json:"beneficiary"`
}

func toEnclosureEventResponse(e simulation.EnclosureEvent) enclosureEventResponse {
	return enclosureEventResponse{
		ID:                  string(e.ID),
		Period:              e.Period,
		AcresEnclosed:       e.AcresEnclosed,
		PopulationDisplaced: e.PopulationDisplaced,
		Beneficiary:         e.Beneficiary,
		CreatedAt:           e.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

// CreateEnclosureEvent handles POST /v1/enclosure-events.
func (h *Handler) CreateEnclosureEvent(w http.ResponseWriter, r *http.Request) {
	var req createEnclosureEventRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	e := simulation.EnclosureEvent{
		Period:              req.Period,
		AcresEnclosed:       req.AcresEnclosed,
		PopulationDisplaced: req.PopulationDisplaced,
		Beneficiary:         req.Beneficiary,
	}
	if err := e.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.EnclosureEvents.CreateEnclosureEvent(r.Context(), e)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", "/v1/enclosure-events/"+string(created.ID))
	writeJSON(w, http.StatusCreated, toEnclosureEventResponse(created))
}

// ListEnclosureEvents handles GET /v1/enclosure-events.
func (h *Handler) ListEnclosureEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.EnclosureEvents.ListEnclosureEvents(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]enclosureEventResponse, 0, len(events))
	for _, e := range events {
		out = append(out, toEnclosureEventResponse(e))
	}
	writeJSON(w, http.StatusOK, out)
}