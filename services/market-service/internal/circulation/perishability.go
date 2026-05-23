package circulation

// PerishabilityID is a 96-bit hex identifier.
type PerishabilityID string

func (id PerishabilityID) IsZero() bool { return id == "" }

// NewPerishabilityID generates a 96-bit hex identifier via crypto/rand.
func NewPerishabilityID() PerishabilityID { return PerishabilityID(newHexID()) }

// Perishability sets an absolute upper bound on the selling time for a
// commodity. A SellingPhase whose elapsed time exceeds WindowNanos forces
// OutcomeSpoiled and a zero-surplus value-loss event.
// WindowNanos == 0 means no perishability constraint within the sim horizon.
// §: "The more perishable a commodity is … the less is it suited to be an
// object of capitalist production."
type Perishability struct {
	ID          PerishabilityID `json:"id"`
	CommodityID CommodityID     `json:"commodity_id"`
	// WindowNanos is the maximum selling-phase duration before spoilage.
	// 0 = no constraint.
	WindowNanos int64 `json:"window_nanos"`
}

// IsConstrained reports whether this commodity has a spoilage window.
func (p Perishability) IsConstrained() bool { return p.WindowNanos > 0 }

// Spoiled reports whether elapsedNanos has exceeded the perishability window.
// Always returns false if WindowNanos == 0.
func (p Perishability) Spoiled(elapsedNanos int64) bool {
	return p.WindowNanos > 0 && elapsedNanos > p.WindowNanos
}

// MarketSeparationID is a 96-bit hex identifier.
type MarketSeparationID string

func (id MarketSeparationID) IsZero() bool { return id == "" }

// NewMarketSeparationID generates a 96-bit hex identifier via crypto/rand.
func NewMarketSeparationID() MarketSeparationID { return MarketSeparationID(newHexID()) }

// MarketSeparation records that the selling-market and buying-market for one
// industrial capital are spatially distinct. This does not alter the
// time-decomposition: CirculationTime.TotalNanos() is still SellingTimeNanos +
// BuyingTimeNanos regardless of spatial separation.
// §: "C—M and M—C may be separate not only in time but also in space."
type MarketSeparation struct {
	ID                  MarketSeparationID  `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	SellingMarketID     string              `json:"selling_market_id"`
	BuyingMarketID      string              `json:"buying_market_id"`
}

// IsSeparated reports whether the selling and buying markets differ.
func (ms MarketSeparation) IsSeparated() bool { return ms.SellingMarketID != ms.BuyingMarketID }

// CirculationCompressionTargetID is a 96-bit hex identifier.
type CirculationCompressionTargetID string

func (id CirculationCompressionTargetID) IsZero() bool { return id == "" }

// NewCirculationCompressionTargetID generates a 96-bit hex identifier.
func NewCirculationCompressionTargetID() CirculationCompressionTargetID {
	return CirculationCompressionTargetID(newHexID())
}

// CirculationCompressionTarget models the theoretical limit-case where
// circulation time approaches zero (payment-on-delivery in the buyer's own
// means of production). TargetCirculationNanos == 0 is the asymptote.
// §: "if this payment is made in his own means of production, the time of
// circulation approaches zero."
type CirculationCompressionTarget struct {
	ID                     CirculationCompressionTargetID `json:"id"`
	IndustrialCapitalID    IndustrialCapitalID            `json:"industrial_capital_id"`
	TargetCirculationNanos int64                          `json:"target_circulation_nanos"`
}

// IsPaymentOnDelivery reports whether this target represents the zero-limit case.
func (c CirculationCompressionTarget) IsPaymentOnDelivery() bool {
	return c.TargetCirculationNanos == 0
}
