package simulation

import "errors"

// Ch. 27 — Expropriation of the Agricultural Population: the named structural
// invariants.
//
// The persisted record of one expropriation wave is EnclosureEvent (see
// enclosure_event.go), which the HTTP/store layer uses. These three value types
// express the chapter's named structural invariants that the single persisted
// record does not itself enforce:
//
//   - not every dispossessed peasant becomes a wage-labourer
//     (NewProletarians ≤ PeasantsDispossessed);
//   - enclosure conserves acreage — the commons become private land, acre for
//     acre (PrivateAcres == CommonAcres);
//   - the proletariat the transition creates is the sum of the proletarians
//     each episode produced (TotalProletarians == Σ NewProletarians).
//
// They are pure value types: invariant-bearing, side-effect free, exercised by
// the chapter's Marx fixtures in agrarian_transition_test.go.

var (
	// ErrNewProletariansExceedDispossessed — an episode cannot turn more
	// peasants into proletarians than it dispossessed.
	ErrNewProletariansExceedDispossessed = errors.New("simulation: new_proletarians cannot exceed peasants_dispossessed")

	// ErrDispossessedNegative — dispossession/proletarian counts cannot be negative.
	ErrDispossessedNegative = errors.New("simulation: dispossession counts cannot be negative")

	// ErrAcresNotConserved — enclosure conserves acreage: private_acres must
	// equal common_acres (the commons become private land, acre for acre).
	ErrAcresNotConserved = errors.New("simulation: private_acres must equal common_acres (acre conservation)")

	// ErrProletarianSumMismatch — the transition's total proletarians must
	// equal the sum of the proletarians its episodes produced.
	ErrProletarianSumMismatch = errors.New("simulation: total_proletarians must equal the sum over episodes")
)

// ExpropriationEpisode is one episode of agrarian dispossession: the peasants
// torn from the land and, of those, the number who become "free" wage-labourers
// (proletarians). Not all dispossessed become proletarians — some emigrate,
// vagabond, or perish under the bloody legislation of Ch. 28 — so the chapter's
// only monotone inequality holds: NewProletarians ≤ PeasantsDispossessed.
type ExpropriationEpisode struct {
	Period               string `json:"period"`
	PeasantsDispossessed int64  `json:"peasants_dispossessed"`
	NewProletarians      int64  `json:"new_proletarians"`
}

// Validate enforces the episode invariant: a period is named, counts are
// non-negative, and the proletarians created cannot exceed the peasants
// dispossessed.
func (e ExpropriationEpisode) Validate() error {
	if e.Period == "" {
		return ErrPeriodRequired
	}
	if e.PeasantsDispossessed < 0 || e.NewProletarians < 0 {
		return ErrDispossessedNegative
	}
	if e.NewProletarians > e.PeasantsDispossessed {
		return ErrNewProletariansExceedDispossessed
	}
	return nil
}

// Enclosure records common land converted to private property in one wave. The
// chapter's conservation invariant holds: the acres taken from the commons
// reappear, acre for acre, as private acres — PrivateAcres == CommonAcres.
type Enclosure struct {
	Period       string `json:"period"`
	CommonAcres  int64  `json:"common_acres"`
	PrivateAcres int64  `json:"private_acres"`
}

// Validate enforces acre conservation: a positive acreage is enclosed from the
// commons, and exactly that acreage is added to private holdings.
func (e Enclosure) Validate() error {
	if e.Period == "" {
		return ErrPeriodRequired
	}
	if e.CommonAcres <= 0 {
		return ErrAcresNonPositive
	}
	if e.PrivateAcres != e.CommonAcres {
		return ErrAcresNotConserved
	}
	return nil
}

// AgrarianTransition is the whole movement from a peasantry to a proletariat: a
// sequence of expropriation episodes and the total proletariat they produce.
// Its invariant ties the aggregate to its parts: TotalProletarians equals the
// sum of NewProletarians across the episodes.
type AgrarianTransition struct {
	Episodes          []ExpropriationEpisode `json:"episodes"`
	TotalProletarians int64                  `json:"total_proletarians"`
}

// SumNewProletarians returns the proletarians produced across all episodes.
func (t AgrarianTransition) SumNewProletarians() int64 {
	var sum int64
	for _, e := range t.Episodes {
		sum += e.NewProletarians
	}
	return sum
}

// Validate enforces each episode's own invariant and then the aggregation
// invariant: TotalProletarians == Σ NewProletarians.
func (t AgrarianTransition) Validate() error {
	for _, e := range t.Episodes {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	if t.TotalProletarians != t.SumNewProletarians() {
		return ErrProletarianSumMismatch
	}
	return nil
}
