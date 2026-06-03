package store

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Memory implements Store, CircuitStore, and LabourPowerStore for tests and local dev.
type Memory struct {
	mu                  sync.RWMutex
	agents              map[agent.ID]agent.Agent
	circuits            map[agent.ID]agent.CapitalCircuit
	circuitLegs         map[agent.ID]agent.CircuitLeg
	labourWorkers       map[agent.AgentID]agent.Worker
	labourCapitalists   map[agent.AgentID]agent.Capitalist
	offerings           map[agent.AgentID]agent.LabourPowerOffering
	purchases           map[agent.PurchaseID]agent.LabourPowerPurchase
	labourProcesses     map[agent.LabourProcessID]agent.LabourProcess
	workingDays         map[agent.WorkingDayID]agent.WorkingDay
	relaySchedules      map[agent.RelayScheduleID]agent.RelaySchedule
	cooperations        map[agent.CooperationID]agent.Cooperation
	manufactures        map[agent.ManufactureID]agent.Manufacture
	wageForms           map[agent.AgentID]agent.WageForm
	workingSessions     map[agent.WorkingSessionID]agent.WorkingSession
	pieceWages          map[agent.AgentID]agent.PieceWage
	subContracts        map[agent.SubContractID]agent.SubContract
	piecePricePoints    map[agent.AgentID][]agent.PiecePricePoint
	nationalIntensities map[agent.CountryCode]agent.NationalIntensity
	dayWages            map[agent.CountryCode]agent.DayWage
	spindleRatios       map[agent.CountryCode]agent.SpindleRatio
	now                 func() time.Time
}

func NewMemory() *Memory {
	return &Memory{
		agents:              make(map[agent.ID]agent.Agent),
		circuits:            make(map[agent.ID]agent.CapitalCircuit),
		circuitLegs:         make(map[agent.ID]agent.CircuitLeg),
		labourWorkers:       make(map[agent.AgentID]agent.Worker),
		labourCapitalists:   make(map[agent.AgentID]agent.Capitalist),
		offerings:           make(map[agent.AgentID]agent.LabourPowerOffering),
		purchases:           make(map[agent.PurchaseID]agent.LabourPowerPurchase),
		labourProcesses:     make(map[agent.LabourProcessID]agent.LabourProcess),
		workingDays:         make(map[agent.WorkingDayID]agent.WorkingDay),
		relaySchedules:      make(map[agent.RelayScheduleID]agent.RelaySchedule),
		cooperations:        make(map[agent.CooperationID]agent.Cooperation),
		manufactures:        make(map[agent.ManufactureID]agent.Manufacture),
		wageForms:           make(map[agent.AgentID]agent.WageForm),
		workingSessions:     make(map[agent.WorkingSessionID]agent.WorkingSession),
		pieceWages:          make(map[agent.AgentID]agent.PieceWage),
		subContracts:        make(map[agent.SubContractID]agent.SubContract),
		piecePricePoints:    make(map[agent.AgentID][]agent.PiecePricePoint),
		nationalIntensities: make(map[agent.CountryCode]agent.NationalIntensity),
		dayWages:            make(map[agent.CountryCode]agent.DayWage),
		spindleRatios:       make(map[agent.CountryCode]agent.SpindleRatio),
		now:                 time.Now,
	}
}

func (m *Memory) Create(_ context.Context, a agent.Agent) (agent.Agent, error) {
	if err := a.Validate(); err != nil {
		return agent.Agent{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if a.ID.IsZero() {
		a.ID = agent.NewID()
	}
	now := m.now()
	a.CreatedAt = now
	a.UpdatedAt = now
	m.agents[a.ID] = a
	return a, nil
}

func (m *Memory) Get(_ context.Context, id agent.ID) (agent.Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.agents[id]
	if !ok {
		return agent.Agent{}, ErrNotFound
	}
	return a, nil
}

func (m *Memory) List(_ context.Context) ([]agent.Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Agent, 0, len(m.agents))
	for _, a := range m.agents {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func (m *Memory) ListByClass(_ context.Context, class agent.Class) ([]agent.Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.Agent
	for _, a := range m.agents {
		if a.Class == class {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func (m *Memory) Update(_ context.Context, id agent.ID, u Update) (agent.Agent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.agents[id]
	if !ok {
		return agent.Agent{}, ErrNotFound
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return agent.Agent{}, err
	}
	next.UpdatedAt = m.now()
	m.agents[id] = next
	return next, nil
}

func (m *Memory) Delete(_ context.Context, id agent.ID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.agents[id]; !ok {
		return ErrNotFound
	}
	delete(m.agents, id)
	return nil
}

// CreateCircuit atomically inserts the circuit and updates the agent's
// money_balance by circuit.SurplusValue. Returns ErrNotFound if the agent
// doesn't exist; returns agent.ErrInsufficientFunds if balance < MAdvanced.
func (m *Memory) CreateCircuit(_ context.Context, c agent.CapitalCircuit) (agent.CapitalCircuit, error) {
	if err := c.Validate(); err != nil {
		return agent.CapitalCircuit{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.agents[c.AgentID]
	if !ok {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	if c.MAdvanced > a.MoneyBalance {
		return agent.CapitalCircuit{}, agent.ErrInsufficientFunds
	}
	if c.ID.IsZero() {
		c.ID = agent.NewID()
	}
	c.CreatedAt = m.now()
	m.circuits[c.ID] = c
	a.MoneyBalance += c.SurplusValue
	a.UpdatedAt = m.now()
	m.agents[a.ID] = a
	return c, nil
}

func (m *Memory) GetCircuit(_ context.Context, id agent.ID) (agent.CapitalCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.circuits[id]
	if !ok {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	return c, nil
}

func (m *Memory) ListCircuits(_ context.Context, agentID agent.ID) ([]agent.CapitalCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.CapitalCircuit
	for _, c := range m.circuits {
		if c.AgentID == agentID {
			out = append(out, c)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCircuitLeg records one closed leg of a coordinated multi-agent
// exchange. Unlike CreateCircuit it touches no balance and does not require the
// agent to exist — an owner may live only in market-service.
func (m *Memory) CreateCircuitLeg(_ context.Context, l agent.CircuitLeg) (agent.CircuitLeg, error) {
	if err := l.Validate(); err != nil {
		return agent.CircuitLeg{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if l.ID.IsZero() {
		l.ID = agent.NewID()
	}
	l.CreatedAt = m.now()
	m.circuitLegs[l.ID] = l
	return l, nil
}

func (m *Memory) ListCircuitLegs(_ context.Context, agentID agent.ID) ([]agent.CircuitLeg, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.CircuitLeg
	for _, l := range m.circuitLegs {
		if l.AgentID == agentID {
			out = append(out, l)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreateWorker(_ context.Context, w agent.Worker) (agent.Worker, error) {
	if err := w.Validate(); err != nil {
		return agent.Worker{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if w.ID.IsZero() {
		w.ID = agent.NewAgentID()
	}
	now := m.now()
	w.CreatedAt = now
	w.UpdatedAt = now
	w.Kind = agent.AgentKindWorker
	m.labourWorkers[w.ID] = w
	return w, nil
}

func (m *Memory) GetWorker(_ context.Context, id agent.AgentID) (agent.Worker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	w, ok := m.labourWorkers[id]
	if !ok {
		return agent.Worker{}, ErrNotFound
	}
	return w, nil
}

func (m *Memory) ListWorkers(_ context.Context) ([]agent.Worker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Worker, 0, len(m.labourWorkers))
	for _, w := range m.labourWorkers {
		out = append(out, w)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreateCapitalist(_ context.Context, c agent.Capitalist) (agent.Capitalist, error) {
	if err := c.Validate(); err != nil {
		return agent.Capitalist{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if c.ID.IsZero() {
		c.ID = agent.NewAgentID()
	}
	now := m.now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Kind = agent.AgentKindCapitalist
	m.labourCapitalists[c.ID] = c
	return c, nil
}

func (m *Memory) GetCapitalist(_ context.Context, id agent.AgentID) (agent.Capitalist, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.labourCapitalists[id]
	if !ok {
		return agent.Capitalist{}, ErrNotFound
	}
	return c, nil
}

func (m *Memory) ListCapitalists(_ context.Context) ([]agent.Capitalist, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Capitalist, 0, len(m.labourCapitalists))
	for _, c := range m.labourCapitalists {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreateOffering(_ context.Context, o agent.LabourPowerOffering) (agent.LabourPowerOffering, error) {
	if err := o.Validate(); err != nil {
		return agent.LabourPowerOffering{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if o.ID.IsZero() {
		o.ID = agent.NewAgentID()
	}
	o.CreatedAt = m.now()
	m.offerings[o.ID] = o
	return o, nil
}

func (m *Memory) ListOfferings(_ context.Context) ([]agent.LabourPowerOffering, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.LabourPowerOffering, 0, len(m.offerings))
	for _, o := range m.offerings {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreatePurchase(_ context.Context, p agent.LabourPowerPurchase) (agent.LabourPowerPurchase, error) {
	if err := p.Validate(); err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID.IsZero() {
		p.ID = agent.NewPurchaseID()
	}
	p.CreatedAt = m.now()
	m.purchases[p.ID] = p
	return p, nil
}

func (m *Memory) GetPurchase(_ context.Context, id agent.PurchaseID) (agent.LabourPowerPurchase, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.purchases[id]
	if !ok {
		return agent.LabourPowerPurchase{}, ErrNotFound
	}
	return p, nil
}

func (m *Memory) ListPurchases(_ context.Context) ([]agent.LabourPowerPurchase, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.LabourPowerPurchase, 0, len(m.purchases))
	for _, p := range m.purchases {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreateLabourProcess(_ context.Context, lp agent.LabourProcess) (agent.LabourProcess, error) {
	if err := lp.Validate(); err != nil {
		return agent.LabourProcess{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if lp.ID.IsZero() {
		lp.ID = agent.NewLabourProcessID()
	}
	lp.CreatedAt = m.now()
	m.labourProcesses[lp.ID] = lp
	return lp, nil
}

func (m *Memory) GetLabourProcess(_ context.Context, id agent.LabourProcessID) (agent.LabourProcess, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	lp, ok := m.labourProcesses[id]
	if !ok {
		return agent.LabourProcess{}, ErrNotFound
	}
	return lp, nil
}

func (m *Memory) CreateWorkingDay(_ context.Context, wd agent.WorkingDay) (agent.WorkingDay, error) {
	if err := wd.Validate(); err != nil {
		return agent.WorkingDay{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if wd.ID.IsZero() {
		wd.ID = agent.NewWorkingDayID()
	}
	wd.CreatedAt = m.now()
	m.workingDays[wd.ID] = wd
	return wd, nil
}

func (m *Memory) GetWorkingDay(_ context.Context, id agent.WorkingDayID) (agent.WorkingDay, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	wd, ok := m.workingDays[id]
	if !ok {
		return agent.WorkingDay{}, ErrNotFound
	}
	return wd, nil
}

func (m *Memory) CreateRelaySchedule(_ context.Context, rs agent.RelaySchedule) (agent.RelaySchedule, error) {
	if err := rs.Validate(); err != nil {
		return agent.RelaySchedule{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if rs.ID.IsZero() {
		rs.ID = agent.NewRelayScheduleID()
	}
	rs.CreatedAt = m.now()
	m.relaySchedules[rs.ID] = rs
	return rs, nil
}

func (m *Memory) GetRelaySchedule(_ context.Context, id agent.RelayScheduleID) (agent.RelaySchedule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rs, ok := m.relaySchedules[id]
	if !ok {
		return agent.RelaySchedule{}, ErrNotFound
	}
	return rs, nil
}

func (m *Memory) CreateCooperation(_ context.Context, c agent.Cooperation) (agent.Cooperation, error) {
	if err := c.Validate(); err != nil {
		return agent.Cooperation{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if c.ID.IsZero() {
		c.ID = agent.NewCooperationID()
	}
	c.CreatedAt = m.now()
	stored := agent.Cooperation{
		ID:           c.ID,
		Name:         c.Name,
		CapitalistID: c.CapitalistID,
		Members:      append([]agent.CooperationMember(nil), c.Members...),
		CreatedAt:    c.CreatedAt,
	}
	m.cooperations[c.ID] = stored
	return stored, nil
}

func (m *Memory) GetCooperation(_ context.Context, id agent.CooperationID) (agent.Cooperation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.cooperations[id]
	if !ok {
		return agent.Cooperation{}, ErrNotFound
	}
	out := c
	out.Members = append([]agent.CooperationMember(nil), c.Members...)
	return out, nil
}

func (m *Memory) ListCooperations(_ context.Context) ([]agent.Cooperation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Cooperation, 0, len(m.cooperations))
	for _, c := range m.cooperations {
		copyC := c
		copyC.Members = append([]agent.CooperationMember(nil), c.Members...)
		out = append(out, copyC)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) ListCooperationsByCapitalist(_ context.Context, capitalistID agent.AgentID) ([]agent.Cooperation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.Cooperation
	for _, c := range m.cooperations {
		if c.CapitalistID != capitalistID {
			continue
		}
		copyC := c
		copyC.Members = append([]agent.CooperationMember(nil), c.Members...)
		out = append(out, copyC)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func cloneManufacture(m agent.Manufacture) agent.Manufacture {
	out := m
	out.Roles = append([]agent.DetailRole(nil), m.Roles...)
	return out
}

func (m *Memory) CreateManufacture(_ context.Context, mf agent.Manufacture) (agent.Manufacture, error) {
	if err := mf.Validate(); err != nil {
		return agent.Manufacture{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if mf.ID.IsZero() {
		mf.ID = agent.NewManufactureID()
	}
	mf.CreatedAt = m.now()
	stored := cloneManufacture(mf)
	m.manufactures[mf.ID] = stored
	return cloneManufacture(stored), nil
}

func (m *Memory) GetManufacture(_ context.Context, id agent.ManufactureID) (agent.Manufacture, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mf, ok := m.manufactures[id]
	if !ok {
		return agent.Manufacture{}, ErrNotFound
	}
	return cloneManufacture(mf), nil
}

func (m *Memory) ListManufactures(_ context.Context) ([]agent.Manufacture, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Manufacture, 0, len(m.manufactures))
	for _, mf := range m.manufactures {
		out = append(out, cloneManufacture(mf))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) ListManufacturesByCapitalist(_ context.Context, capitalistID agent.AgentID) ([]agent.Manufacture, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.Manufacture
	for _, mf := range m.manufactures {
		if mf.CapitalistID != capitalistID {
			continue
		}
		out = append(out, cloneManufacture(mf))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) ListManufacturesByForm(_ context.Context, form agent.ManufactureForm) ([]agent.Manufacture, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.Manufacture
	for _, mf := range m.manufactures {
		if mf.Form != form {
			continue
		}
		out = append(out, cloneManufacture(mf))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreateWageForm(_ context.Context, wf agent.WageForm) (agent.WageForm, error) {
	if err := wf.Validate(); err != nil {
		return agent.WageForm{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if wf.ID.IsZero() {
		wf.ID = agent.NewWageFormID()
	}
	wf.CreatedAt = m.now()
	m.wageForms[wf.AgentID] = wf
	return wf, nil
}

func (m *Memory) GetWageForm(_ context.Context, agentID agent.AgentID) (agent.WageForm, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	wf, ok := m.wageForms[agentID]
	if !ok {
		return agent.WageForm{}, ErrNotFound
	}
	return wf, nil
}

func (m *Memory) ListWageForms(_ context.Context) ([]agent.WageForm, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.WageForm, 0, len(m.wageForms))
	for _, wf := range m.wageForms {
		out = append(out, wf)
	}
	return out, nil
}

func (m *Memory) CreateWorkingSession(_ context.Context, s agent.WorkingSession) (agent.WorkingSession, error) {
	if err := s.Validate(); err != nil {
		return agent.WorkingSession{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.ID.IsZero() {
		s.ID = agent.NewWorkingSessionID()
	}
	s.CreatedAt = m.now()
	m.workingSessions[s.ID] = s
	return s, nil
}

func (m *Memory) GetWorkingSession(_ context.Context, id agent.WorkingSessionID) (agent.WorkingSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.workingSessions[id]
	if !ok {
		return agent.WorkingSession{}, ErrNotFound
	}
	return s, nil
}

func (m *Memory) ListWorkingSessions(_ context.Context, agentID agent.AgentID) ([]agent.WorkingSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.WorkingSession
	for _, s := range m.workingSessions {
		if s.AgentID == agentID {
			out = append(out, s)
		}
	}
	if out == nil {
		out = []agent.WorkingSession{}
	}
	return out, nil
}

func (m *Memory) CreatePieceWage(_ context.Context, pw agent.PieceWage) (agent.PieceWage, error) {
	if err := pw.Validate(); err != nil {
		return agent.PieceWage{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.pieceWages[pw.AgentID]; exists {
		return agent.PieceWage{}, ErrAlreadyExists
	}
	if pw.ID.IsZero() {
		pw.ID = agent.NewPieceWageID()
	}
	pw.CreatedAt = m.now()
	m.pieceWages[pw.AgentID] = pw
	return pw, nil
}

func (m *Memory) GetPieceWage(_ context.Context, agentID agent.AgentID) (agent.PieceWage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pw, ok := m.pieceWages[agentID]
	if !ok {
		return agent.PieceWage{}, ErrNotFound
	}
	return pw, nil
}

func (m *Memory) ListPieceWages(_ context.Context) ([]agent.PieceWage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.PieceWage, 0, len(m.pieceWages))
	for _, pw := range m.pieceWages {
		out = append(out, pw)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (m *Memory) AppendPiecePricePoint(_ context.Context, p agent.PiecePricePoint) (agent.PiecePricePoint, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	existing := m.piecePricePoints[p.AgentID]
	if n := len(existing); n > 0 && existing[n-1].ProductivityFactor == p.ProductivityFactor {
		return existing[n-1], false, nil
	}
	if p.ID.IsZero() {
		p.ID = agent.NewPiecePricePointID()
	}
	p.Sequence = int64(len(existing))
	if p.OccurredAt.IsZero() {
		p.OccurredAt = m.now()
	}
	m.piecePricePoints[p.AgentID] = append(existing, p)
	return p, true, nil
}

func (m *Memory) ListPiecePricePoints(_ context.Context, agentID agent.AgentID) ([]agent.PiecePricePoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pts := m.piecePricePoints[agentID]
	out := make([]agent.PiecePricePoint, len(pts))
	copy(out, pts)
	return out, nil
}

func (m *Memory) CreateSubContract(_ context.Context, sc agent.SubContract) (agent.SubContract, error) {
	if err := sc.Validate(); err != nil {
		return agent.SubContract{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if sc.ID.IsZero() {
		sc.ID = agent.NewSubContractID()
	}
	sc.CreatedAt = m.now()
	m.subContracts[sc.ID] = sc
	return sc, nil
}

func (m *Memory) GetSubContract(_ context.Context, id agent.SubContractID) (agent.SubContract, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sc, ok := m.subContracts[id]
	if !ok {
		return agent.SubContract{}, ErrNotFound
	}
	return sc, nil
}

func (m *Memory) UpsertIntensity(_ context.Context, ni agent.NationalIntensity) (agent.NationalIntensity, error) {
	if err := ni.Validate(); err != nil {
		return agent.NationalIntensity{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if existing, ok := m.nationalIntensities[ni.CountryCode]; ok {
		ni.CreatedAt = existing.CreatedAt
	} else {
		ni.CreatedAt = m.now()
	}
	m.nationalIntensities[ni.CountryCode] = ni
	return ni, nil
}

func (m *Memory) ListIntensities(_ context.Context) ([]agent.NationalIntensity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.NationalIntensity, 0, len(m.nationalIntensities))
	for _, ni := range m.nationalIntensities {
		out = append(out, ni)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CountryCode < out[j].CountryCode })
	return out, nil
}

func (m *Memory) CreateDayWage(_ context.Context, w agent.DayWage) (agent.DayWage, error) {
	if err := w.Validate(); err != nil {
		return agent.DayWage{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.dayWages[w.CountryCode]; exists {
		return agent.DayWage{}, ErrAlreadyExists
	}
	w.CreatedAt = m.now()
	m.dayWages[w.CountryCode] = w
	return w, nil
}

func (m *Memory) GetDayWage(_ context.Context, country agent.CountryCode) (agent.DayWage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	w, ok := m.dayWages[country]
	if !ok {
		return agent.DayWage{}, ErrNotFound
	}
	return w, nil
}

func (m *Memory) ListDayWages(_ context.Context) ([]agent.DayWage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.DayWage, 0, len(m.dayWages))
	for _, w := range m.dayWages {
		out = append(out, w)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CountryCode < out[j].CountryCode })
	return out, nil
}

func (m *Memory) ListSpindleRatios(_ context.Context) ([]agent.SpindleRatio, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.SpindleRatio, 0, len(m.spindleRatios))
	for _, sr := range m.spindleRatios {
		out = append(out, sr)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CountryCode < out[j].CountryCode })
	return out, nil
}
