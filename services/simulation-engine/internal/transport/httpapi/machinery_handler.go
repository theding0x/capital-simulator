package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// machineDTO is the JSON shape for a Machine over the wire.
type machineDTO struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	MotorMechanism           string `json:"motor_mechanism"`
	TransmittingMechanism    string `json:"transmitting_mechanism"`
	WorkingTool              string `json:"working_tool"`
	MachineValue             int64  `json:"machine_value"`
	LifespanDays             int64  `json:"lifespan_days"`
	ProductivePower          int64  `json:"productive_power"`
	HandLabourPerUnit        int64  `json:"hand_labour_per_unit"`
	AccumulatedWear          int64  `json:"accumulated_wear"`
	AccumulatedDepreciation  int64  `json:"accumulated_depreciation"`
	DailyWearAndTear         int64  `json:"daily_wear_and_tear"`
	ValueTransferredPerUnit  int64  `json:"value_transferred_per_unit"`
	DailyLabourDisplaced     int64  `json:"daily_labour_displaced"`
	CreatedAt                string `json:"created_at,omitempty"`
}

func dtoFromMachine(m machinery.Machine) machineDTO {
	d := machineDTO{
		ID:                      string(m.ID),
		Name:                    m.Name,
		MotorMechanism:          m.MotorMechanism,
		TransmittingMechanism:   m.TransmittingMechanism,
		WorkingTool:             m.WorkingTool,
		MachineValue:            int64(m.MachineValue),
		LifespanDays:            int64(m.LifespanDays),
		ProductivePower:         int64(m.ProductivePower),
		HandLabourPerUnit:       int64(m.HandLabourPerUnit),
		AccumulatedWear:         int64(m.AccumulatedWear.Value),
		AccumulatedDepreciation: int64(m.AccumulatedDepreciation.Value),
		DailyWearAndTear:        int64(machinery.DailyWearAndTear(m)),
		ValueTransferredPerUnit: int64(machinery.ValueTransferredPerUnit(m, int64(m.ProductivePower))),
		DailyLabourDisplaced:    int64(machinery.LabourDisplaced(m.HandLabourPerUnit, int64(m.ProductivePower))),
	}
	if !m.CreatedAt.IsZero() {
		d.CreatedAt = m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return d
}

// createMachineRequest accepts a Machine's intrinsic fields. The server
// assigns the ID, CreatedAt, and any accumulated counters.
type createMachineRequest struct {
	Name                  string `json:"name"`
	MotorMechanism        string `json:"motor_mechanism"`
	TransmittingMechanism string `json:"transmitting_mechanism"`
	WorkingTool           string `json:"working_tool"`
	MachineValue          int64  `json:"machine_value"`
	LifespanDays          int64  `json:"lifespan_days"`
	ProductivePower       int64  `json:"productive_power"`
	HandLabourPerUnit     int64  `json:"hand_labour_per_unit"`
}

// CreateMachine handles POST /v1/machines.
func (h *Handler) CreateMachine(w http.ResponseWriter, r *http.Request) {
	if h.Machines == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	var req createMachineRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	mc, err := h.Machines.CreateMachine(r.Context(), machinery.Machine{
		Name:                  req.Name,
		MotorMechanism:        req.MotorMechanism,
		TransmittingMechanism: req.TransmittingMechanism,
		WorkingTool:           req.WorkingTool,
		MachineValue:          machinery.MachineValue(req.MachineValue),
		LifespanDays:          machinery.LifespanDays(req.LifespanDays),
		ProductivePower:       machinery.ProductivePower(req.ProductivePower),
		HandLabourPerUnit:     machinery.LabourMinutes(req.HandLabourPerUnit),
	})
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, dtoFromMachine(mc))
}

// ListMachines handles GET /v1/machines.
func (h *Handler) ListMachines(w http.ResponseWriter, r *http.Request) {
	if h.Machines == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	list, err := h.Machines.ListMachines(r.Context())
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	out := make([]machineDTO, 0, len(list))
	for _, m := range list {
		out = append(out, dtoFromMachine(m))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// GetMachine handles GET /v1/machines/{id}.
func (h *Handler) GetMachine(w http.ResponseWriter, r *http.Request) {
	if h.Machines == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	id := r.PathValue("id")
	mc, err := h.Machines.GetMachine(r.Context(), machinery.MachineID(id))
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dtoFromMachine(mc))
}

// machineWearResponse is GET /v1/machines/{id}/wear.
type machineWearResponse struct {
	ID                      string `json:"id"`
	AccumulatedWear         int64  `json:"accumulated_wear"`
	AccumulatedDepreciation int64  `json:"accumulated_depreciation"`
	RemainingValue          int64  `json:"remaining_value"`
}

// GetMachineWear handles GET /v1/machines/{id}/wear.
func (h *Handler) GetMachineWear(w http.ResponseWriter, r *http.Request) {
	if h.Machines == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	id := r.PathValue("id")
	mc, err := h.Machines.GetMachine(r.Context(), machinery.MachineID(id))
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	remaining := int64(mc.MachineValue) - int64(mc.AccumulatedWear.Value) - int64(mc.AccumulatedDepreciation.Value)
	if remaining < 0 {
		remaining = 0
	}
	writeJSON(w, http.StatusOK, machineWearResponse{
		ID:                      string(mc.ID),
		AccumulatedWear:         int64(mc.AccumulatedWear.Value),
		AccumulatedDepreciation: int64(mc.AccumulatedDepreciation.Value),
		RemainingValue:          remaining,
	})
}

// factoryDTO is the JSON shape for a Factory over the wire.
type factoryDTO struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	PrimeMover           primeMoverDTO `json:"prime_mover"`
	Machines             []machineDTO `json:"machines"`
	TickCount            int64        `json:"tick_count"`
	TotalProductivePower int64        `json:"total_productive_power"`
	DailyValueTransfer   int64        `json:"daily_value_transfer"`
	CreatedAt            string       `json:"created_at,omitempty"`
}

type primeMoverDTO struct {
	Kind       string `json:"kind"`
	Horsepower int64  `json:"horsepower"`
}

func dtoFromFactory(f machinery.Factory) factoryDTO {
	d := factoryDTO{
		ID:                   string(f.ID),
		Name:                 f.Name,
		PrimeMover:           primeMoverDTO{Kind: string(f.PrimeMover.Kind), Horsepower: f.PrimeMover.Horsepower},
		Machines:             make([]machineDTO, 0, len(f.Machines)),
		TickCount:            f.TickCount,
		TotalProductivePower: int64(f.TotalProductivePower()),
		DailyValueTransfer:   int64(f.DailyValueTransfer()),
	}
	for _, m := range f.Machines {
		d.Machines = append(d.Machines, dtoFromMachine(m))
	}
	if !f.CreatedAt.IsZero() {
		d.CreatedAt = f.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return d
}

// createFactoryRequest is POST /v1/factories.
// Each Machines entry may carry a bare `id` (referencing an existing record)
// or the full Machine fields (which the server then persists).
type createFactoryRequest struct {
	Name       string                 `json:"name"`
	PrimeMover primeMoverDTO          `json:"prime_mover"`
	Machines   []createMachineEntry   `json:"machines"`
}

type createMachineEntry struct {
	ID                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	MotorMechanism        string `json:"motor_mechanism,omitempty"`
	TransmittingMechanism string `json:"transmitting_mechanism,omitempty"`
	WorkingTool           string `json:"working_tool,omitempty"`
	MachineValue          int64  `json:"machine_value,omitempty"`
	LifespanDays          int64  `json:"lifespan_days,omitempty"`
	ProductivePower       int64  `json:"productive_power,omitempty"`
	HandLabourPerUnit     int64  `json:"hand_labour_per_unit,omitempty"`
}

// CreateFactory handles POST /v1/factories.
func (h *Handler) CreateFactory(w http.ResponseWriter, r *http.Request) {
	if h.Factories == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	var req createFactoryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.PrimeMover.Kind == "" {
		writeError(w, http.StatusBadRequest, "prime_mover.kind is required")
		return
	}
	if len(req.Machines) == 0 {
		writeError(w, http.StatusBadRequest, "at least one machine entry is required")
		return
	}
	machines := make([]machinery.Machine, 0, len(req.Machines))
	for _, e := range req.Machines {
		machines = append(machines, machinery.Machine{
			ID:                    machinery.MachineID(e.ID),
			Name:                  e.Name,
			MotorMechanism:        e.MotorMechanism,
			TransmittingMechanism: e.TransmittingMechanism,
			WorkingTool:           e.WorkingTool,
			MachineValue:          machinery.MachineValue(e.MachineValue),
			LifespanDays:          machinery.LifespanDays(e.LifespanDays),
			ProductivePower:       machinery.ProductivePower(e.ProductivePower),
			HandLabourPerUnit:     machinery.LabourMinutes(e.HandLabourPerUnit),
		})
	}
	f, err := h.Factories.CreateFactory(r.Context(), machinery.Factory{
		Name: req.Name,
		PrimeMover: machinery.PrimeMover{
			Kind:       machinery.PrimeMoverKind(req.PrimeMover.Kind),
			Horsepower: req.PrimeMover.Horsepower,
		},
		Machines: machines,
	})
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, dtoFromFactory(f))
}

// GetFactory handles GET /v1/factories/{id}.
func (h *Handler) GetFactory(w http.ResponseWriter, r *http.Request) {
	if h.Factories == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	id := r.PathValue("id")
	f, err := h.Factories.GetFactory(r.Context(), machinery.FactoryID(id))
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dtoFromFactory(f))
}

// ListFactories handles GET /v1/factories.
func (h *Handler) ListFactories(w http.ResponseWriter, r *http.Request) {
	if h.Factories == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	list, err := h.Factories.ListFactories(r.Context())
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	out := make([]factoryDTO, 0, len(list))
	for _, f := range list {
		out = append(out, dtoFromFactory(f))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// tickResponse is POST /v1/factories/{id}/tick.
type tickResponse struct {
	Factory          factoryDTO `json:"factory"`
	ValueTransferred int64      `json:"value_transferred"`
	UnitsProduced    int64      `json:"units_produced"`
	HandLabourSaved  int64      `json:"hand_labour_saved"`
}

// TickFactory handles POST /v1/factories/{id}/tick.
func (h *Handler) TickFactory(w http.ResponseWriter, r *http.Request) {
	if h.Factories == nil {
		writeError(w, http.StatusServiceUnavailable, "machinery store not configured")
		return
	}
	id := r.PathValue("id")
	f, result, err := h.Factories.AdvanceTick(r.Context(), machinery.FactoryID(id))
	if err != nil {
		writeMachineryError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tickResponse{
		Factory:          dtoFromFactory(f),
		ValueTransferred: int64(result.ValueTransferred),
		UnitsProduced:    result.UnitsProduced,
		HandLabourSaved:  int64(result.HandLabourSaved),
	})
}

func writeMachineryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, machinery.ErrMachineValue),
		errors.Is(err, machinery.ErrLifespanDays),
		errors.Is(err, machinery.ErrProductivePower),
		errors.Is(err, machinery.ErrPrimeMoverKind),
		errors.Is(err, machinery.ErrFactoryNoMachines),
		errors.Is(err, machinery.ErrFactoryNoMover):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
