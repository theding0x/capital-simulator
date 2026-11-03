package surplus

// Capital Vol. I, Ch. 16 — the §1 supporting cast: who counts as productive,
// who the collective labourer is, and how fully labour is subsumed to capital.

// AgentID identifies a worker held by agent-service. simulation-engine
// cannot import agent-service's internal package, so the ID crosses the
// service boundary as a plain string.
type AgentID string

// ProductiveLabour records whether an agent's labour is productive in the
// capitalist sense. "That labourer alone is productive, who produces
// surplus-value for the capitalist" — usefulness is not the test.
type ProductiveLabour struct {
	AgentID              AgentID `json:"agent_id"`
	ProducesSurplusValue bool    `json:"produces_surplus_value"`
}

// CollectiveLabourer groups the agents whose combined output constitutes
// capitalist production. Here it is a grouping only; cooperation chapters
// compute its productive power.
type CollectiveLabourer []AgentID

// ProductiveMembers returns the members of the collective whose labour
// produces surplus-value, judged against the given records. Agents with no
// record are treated as unproductive.
func (c CollectiveLabourer) ProductiveMembers(records []ProductiveLabour) []AgentID {
	productive := make(map[AgentID]bool, len(records))
	for _, r := range records {
		if r.ProducesSurplusValue {
			productive[r.AgentID] = true
		}
	}
	out := []AgentID{}
	for _, id := range c {
		if productive[id] {
			out = append(out, id)
		}
	}
	return out
}

// SubjectionKind is the degree of labour's subsumption to capital.
type SubjectionKind string

const (
	// SubjectionFormal is control without technical transformation: capital
	// seizes the labour-process as it finds it. Yields absolute surplus-value.
	SubjectionFormal SubjectionKind = "formal"
	// SubjectionReal is the technical process revolutionised by capital, the
	// precondition for relative surplus-value as a general method.
	SubjectionReal SubjectionKind = "real"
)

// Valid reports whether k is one of the two recognised kinds.
func (k SubjectionKind) Valid() bool {
	return k == SubjectionFormal || k == SubjectionReal
}

// SurplusOrigin returns the Origin of surplus-value characteristic of this
// kind of subjection. Unknown kinds return the empty Origin.
func (k SubjectionKind) SurplusOrigin() Origin {
	switch k {
	case SubjectionFormal:
		return OriginAbsolute
	case SubjectionReal:
		return OriginRelative
	}
	return ""
}
