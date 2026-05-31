package httpapi

import (
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createDR1EntryInput is one soil row supplied by the client for a DR I table.
type createDR1EntryInput struct {
	Grade               rent.SoilGrade `json:"grade"`
	CapitalAdvancedBP   int64          `json:"capital_advanced_bp"`
	OutputQuarters      int64          `json:"output_quarters"`
	PriceOfProductionBP int64          `json:"price_of_production_bp"`
}

// createDR1TableRequest is the body for POST /v1/rent/dr1-tables.
type createDR1TableRequest struct {
	Description string                `json:"description"`
	Entries     []createDR1EntryInput `json:"entries"`
}

// CreateDifferentialRentITable handles POST /v1/rent/dr1-tables.
// SurplusProductQtrs and MoneyRentBP are derived server-side via ComputeDifferentialRentI.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateDifferentialRentITable(w http.ResponseWriter, r *http.Request) {
	var req createDR1TableRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Description == "" {
		writeError(w, http.StatusUnprocessableEntity, "description must not be empty")
		return
	}
	if len(req.Entries) < 1 {
		writeError(w, http.StatusUnprocessableEntity, "entries must contain at least one row")
		return
	}
	for i, e := range req.Entries {
		if !e.Grade.IsValid() {
			writeError(w, http.StatusUnprocessableEntity, "entries["+itoa(i)+"]: grade must be 1–4")
			return
		}
		if e.CapitalAdvancedBP <= 0 {
			writeError(w, http.StatusUnprocessableEntity, "entries["+itoa(i)+"]: capital_advanced_bp must be > 0")
			return
		}
		if e.OutputQuarters <= 0 {
			writeError(w, http.StatusUnprocessableEntity, "entries["+itoa(i)+"]: output_quarters must be > 0")
			return
		}
		if e.PriceOfProductionBP <= 0 {
			writeError(w, http.StatusUnprocessableEntity, "entries["+itoa(i)+"]: price_of_production_bp must be > 0")
			return
		}
	}

	raw := make([]rent.DifferentialRentIEntry, len(req.Entries))
	for i, e := range req.Entries {
		raw[i] = rent.DifferentialRentIEntry{
			ID:                  rent.NewDifferentialRentIEntryID(),
			Grade:               e.Grade,
			CapitalAdvancedBP:   e.CapitalAdvancedBP,
			OutputQuarters:      e.OutputQuarters,
			PriceOfProductionBP: e.PriceOfProductionBP,
			CreatedAt:           time.Now().UTC(),
		}
	}

	tbl, computed := rent.ComputeDifferentialRentI(raw)
	tbl.ID = rent.NewDifferentialRentITableID()
	tbl.Description = req.Description
	tbl.CreatedAt = time.Now().UTC()
	for i := range computed {
		computed[i].TableID = string(tbl.ID)
	}

	savedTbl, savedEntries, err := h.Store.CreateDifferentialRentITable(r.Context(), tbl, computed)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if savedEntries == nil {
		savedEntries = []rent.DifferentialRentIEntry{}
	}
	w.Header().Set("Location", "/v1/rent/dr1-tables/"+string(savedTbl.ID))
	writeJSON(w, http.StatusCreated, map[string]any{
		"table":   savedTbl,
		"entries": savedEntries,
	})
}

// ListDifferentialRentITables handles GET /v1/rent/dr1-tables.
// Returns a JSON object with an "items" array of table headers — never null.
func (h *Handler) ListDifferentialRentITables(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListDifferentialRentITables(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.DifferentialRentITable{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetDifferentialRentITable handles GET /v1/rent/dr1-tables/{id}.
// Returns {"table":..,"entries":..} or 404.
func (h *Handler) GetDifferentialRentITable(w http.ResponseWriter, r *http.Request) {
	id := rent.DifferentialRentITableID(r.PathValue("id"))
	tbl, entries, err := h.Store.GetDifferentialRentITable(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if entries == nil {
		entries = []rent.DifferentialRentIEntry{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"table":   tbl,
		"entries": entries,
	})
}

// ListDifferentialRentIEntries handles GET /v1/rent/dr1-tables/{id}/entries.
// Returns {"items":[...]} or 404 if the table is unknown.
func (h *Handler) ListDifferentialRentIEntries(w http.ResponseWriter, r *http.Request) {
	tableID := rent.DifferentialRentITableID(r.PathValue("id"))
	items, err := h.Store.ListDifferentialRentIEntries(r.Context(), tableID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.DifferentialRentIEntry{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createLocationRentFactorRequest is the body for POST /v1/rent/location-rent-factors.
type createLocationRentFactorRequest struct {
	ParcelID         string `json:"parcel_id"`
	TransportCostBP  int64  `json:"transport_cost_bp"`
	RentEquivalentBP int64  `json:"rent_equivalent_bp"`
}

// CreateLocationRentFactor handles POST /v1/rent/location-rent-factors.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateLocationRentFactor(w http.ResponseWriter, r *http.Request) {
	var req createLocationRentFactorRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.TransportCostBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "transport_cost_bp must be >= 0")
		return
	}
	if req.RentEquivalentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "rent_equivalent_bp must be >= 0")
		return
	}
	f := rent.LocationRentFactor{
		ParcelID:         req.ParcelID,
		TransportCostBP:  req.TransportCostBP,
		RentEquivalentBP: req.RentEquivalentBP,
	}
	saved, err := h.Store.CreateLocationRentFactor(r.Context(), f)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/location-rent-factors/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListLocationRentFactors handles GET /v1/rent/location-rent-factors.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListLocationRentFactors(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListLocationRentFactors(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.LocationRentFactor{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// itoa converts a small non-negative integer to its decimal string representation
// without importing strconv (avoiding an extra import for a single use-site).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
