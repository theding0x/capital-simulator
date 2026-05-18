package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ch. 33 — The Modern Theory of Colonisation.
//
// Marx closes Vol. I by reading Wakefield's colonial theory against the
// grain. In the colonies — virgin soil colonised by free immigrants —
// capital meets producers who are still owners of the conditions of their
// labour. Without prior expropriation, the wage relation cannot
// reproduce itself: Mr. Peel arrives at Swan River with £50,000 of
// capital and 300 working-class persons, and is "left without a servant
// to make his bed." Capital, Wakefield reluctantly admits, is not a
// thing but a social relation. He proposes a remedy — a "sufficient
// price" set on virgin land by the state, prohibitively high relative to
// colonial wages, to delay the labourer's independence long enough that
// the wage market reproduces itself. The chapter therefore confirms by
// negation the central argument of Part VIII: the historical
// precondition of capitalist production is the violent separation of the
// producer from the means of production.

// ColonialLabourMarketID is a 96-bit hex identifier for a persisted
// colonial labour market scenario.
type ColonialLabourMarketID string

// IsZero reports whether the ID is unset.
func (id ColonialLabourMarketID) IsZero() bool { return id == "" }

// NewColonialLabourMarketID generates a 96-bit hex identifier via
// crypto/rand.
func NewColonialLabourMarketID() ColonialLabourMarketID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return ColonialLabourMarketID(hex.EncodeToString(b))
}

// ColonialLabourMarket is the colonial context in which capital meets
// labour. Colony names the settlement (e.g. "Swan River", "Saint
// Domingo"). FreeLabourers is the population of immigrant workers
// nominally available for wage labour. AnnualWagePence is the prevailing
// colonial wage per labourer per year. LandAvailable records whether
// virgin public land is still effectively open to settlement —
// Wakefield's central anti-capitalist datum. WakefieldSchemeApplied
// records whether a SystematicColonisation scheme with a positive
// sufficient price has been applied; a scheme with zero sufficient
// price (no real ransom) is treated as no scheme and leaves this
// false. When true, ColonialLabourRegulation has set the dependence
// period based on the scheme's sufficient price. SurplusLabourExtractable is
// the binary read-out of the chapter's main thesis: without a scheme,
// labour walks off into independent production and capital cannot
// extract surplus; with a sufficient price, the wage relation persists.
type ColonialLabourMarket struct {
	ID                       ColonialLabourMarketID `json:"id"`
	Colony                   string                 `json:"colony"`
	FreeLabourers            int64                  `json:"free_labourers"`
	AnnualWagePence          Pence                  `json:"annual_wage_pence"`
	LandAvailable            bool                   `json:"land_available"`
	WakefieldSchemeApplied   bool                   `json:"wakefield_scheme_applied"`
	IndependenceYears        int64                  `json:"independence_years"`
	SurplusLabourExtractable bool                   `json:"surplus_labour_extractable"`
	CreatedAt                time.Time              `json:"created_at"`
}

// SufficientPrice is Wakefield's artificial price floor on virgin land.
// PricePerAcrePence is the state-imposed price set independently of
// supply and demand; DesiredAcres is the typical holding a labourer
// would seek if he became an independent settler. The product
// PricePerAcrePence * DesiredAcres is the savings target the wage-worker
// must accumulate before he can exit the labour market.
type SufficientPrice struct {
	PricePerAcrePence Pence `json:"price_per_acre_pence"`
	DesiredAcres      int64 `json:"desired_acres"`
}

// Ransom returns the total savings (in pence) the labourer must
// accumulate to buy DesiredAcres at the sufficient price.
func (p SufficientPrice) Ransom() Pence {
	if p.PricePerAcrePence <= 0 || p.DesiredAcres <= 0 {
		return 0
	}
	return Pence(int64(p.PricePerAcrePence) * p.DesiredAcres)
}

// WageWorkerIndependence is a per-worker projection of when the
// labourer has saved enough at the prevailing wage to buy land at the
// sufficient price. WorkerID is a free-form label (the immigrant's name
// or batch number). YearsWorked is the wage years required at
// AnnualSavingsPence to reach TargetRansomPence. BecameLandowner reports
// whether the projection terminates in independence — false when the
// sufficient price exceeds feasible accumulation under the supplied
// yearsLimit.
type WageWorkerIndependence struct {
	WorkerID           string `json:"worker_id"`
	AnnualSavingsPence Pence  `json:"annual_savings_pence"`
	YearsWorked        int64  `json:"years_worked"`
	SavingsPence       Pence  `json:"savings_pence"`
	TargetRansomPence  Pence  `json:"target_ransom_pence"`
	BecameLandowner    bool   `json:"became_landowner"`
}

// SystematicColonisation is Wakefield's full scheme: an artificial land
// price set by the colonial government, the resulting land-sale fund,
// and the count of replacement labourers imported with that fund.
//
// LandFundPence grows as `Ransom() * LandSalesCompleted` — every
// labourer who finally buys himself out of the labour market deposits
// his ransom into the fund, which is then used to ship fresh
// have-nothings from Europe so the wage market is replenished.
type SystematicColonisation struct {
	Colony             string          `json:"colony"`
	SufficientPrice    SufficientPrice `json:"sufficient_price"`
	LandSalesCompleted int64           `json:"land_sales_completed"`
	LandFundPence      Pence           `json:"land_fund_pence"`
	ImportedLabourers  int64           `json:"imported_labourers"`
}

var (
	// ErrColonyRequired is returned when a colonial labour market lacks a
	// colony name.
	ErrColonyRequired = errors.New("simulation: colony name is required")
	// ErrFreeLabourersNonNegative is returned when free labourers is
	// negative.
	ErrFreeLabourersNonNegative = errors.New("simulation: free labourers cannot be negative")
	// ErrAnnualWagePositive is returned when the colonial wage is not
	// strictly positive.
	ErrAnnualWagePositive = errors.New("simulation: annual wage must be positive")
	// ErrSufficientPricePositive is returned when the sufficient price is
	// not strictly positive.
	ErrSufficientPricePositive = errors.New("simulation: sufficient price must be positive")
	// ErrDesiredAcresPositive is returned when desired acres is not
	// strictly positive.
	ErrDesiredAcresPositive = errors.New("simulation: desired acres must be positive")
)

// Validate checks the colonial market invariants.
func (m ColonialLabourMarket) Validate() error {
	if m.Colony == "" {
		return ErrColonyRequired
	}
	if m.FreeLabourers < 0 {
		return ErrFreeLabourersNonNegative
	}
	if m.AnnualWagePence <= 0 {
		return ErrAnnualWagePositive
	}
	return nil
}

// Validate checks the sufficient-price invariants.
func (p SufficientPrice) Validate() error {
	if p.PricePerAcrePence <= 0 {
		return ErrSufficientPricePositive
	}
	if p.DesiredAcres <= 0 {
		return ErrDesiredAcresPositive
	}
	return nil
}

// AnnualSavingsPence is the per-worker per-year savings the simulation
// assumes a colonial wage labourer can accumulate. Marx takes Wakefield
// at his word that colonial wages allow the worker "so great a share
// that he soon becomes a capitalist"; the simulation conservatively
// models savings at half the colonial wage, leaving the other half for
// subsistence.
func AnnualSavingsPence(market ColonialLabourMarket) Pence {
	if market.AnnualWagePence <= 0 {
		return 0
	}
	return market.AnnualWagePence / 2
}

// ColonialLabourRegulation applies Wakefield's systematic colonisation
// scheme to a colonial labour market and returns the regulated state.
//
// Invariants:
//
//   - FreeLabourers is preserved — the scheme does not change how many
//     immigrants are in the colony, only how long they remain
//     wage-workers before exiting.
//   - A scheme with zero sufficient price is no scheme: the market is
//     returned unchanged (WakefieldSchemeApplied stays false). This is
//     the Peel-at-Swan-River case, where capital meets producers who
//     can walk off into independent production unimpeded.
//   - Otherwise IndependenceYears is set to the integer ceiling of
//     ransom / annualSavings, using AnnualSavingsPence(market) as the
//     savings rate assumption (half the colonial wage),
//     WakefieldSchemeApplied is set to true, and
//     SurplusLabourExtractable becomes true when IndependenceYears > 0:
//     by extending the dependence period the scheme makes the wage
//     relation reproducible.
//
// The function is pure: no I/O, no state outside the returned struct.
func ColonialLabourRegulation(market ColonialLabourMarket, scheme SystematicColonisation) ColonialLabourMarket {
	annualSavings := AnnualSavingsPence(market)
	ransom := scheme.SufficientPrice.Ransom()
	if annualSavings <= 0 || ransom <= 0 {
		return market
	}
	out := market
	out.WakefieldSchemeApplied = true
	years := int64(ransom) / int64(annualSavings)
	if int64(ransom)%int64(annualSavings) != 0 {
		years++
	}
	out.IndependenceYears = years
	out.SurplusLabourExtractable = years > 0
	return out
}

// ComputeWageWorkerIndependence projects whether a particular labourer
// reaches independence at the colonial wage given the sufficient price
// in force. If yearsLimit <= 0 the function uses the natural ceiling of
// ransom / annualSavings.
//
// Returns BecameLandowner=true when SavingsPence >= TargetRansomPence at
// the end of YearsWorked. Annual savings default to AnnualSavingsPence
// of the market.
func ComputeWageWorkerIndependence(workerID string, market ColonialLabourMarket, price SufficientPrice, yearsLimit int64) WageWorkerIndependence {
	annual := AnnualSavingsPence(market)
	target := price.Ransom()
	out := WageWorkerIndependence{
		WorkerID:           workerID,
		AnnualSavingsPence: annual,
		TargetRansomPence:  target,
	}
	if annual <= 0 || target <= 0 {
		return out
	}
	required := int64(target) / int64(annual)
	if int64(target)%int64(annual) != 0 {
		required++
	}
	years := required
	if yearsLimit > 0 && yearsLimit < years {
		years = yearsLimit
	}
	out.YearsWorked = years
	out.SavingsPence = Pence(int64(annual) * years)
	out.BecameLandowner = out.SavingsPence >= target
	return out
}

// RecordLandSale appends one completed land sale to the scheme: the
// labourer's ransom is deposited into the LandFundPence, and one
// replacement labourer is treated as imported on that fund. The function
// is pure.
func RecordLandSale(scheme SystematicColonisation) SystematicColonisation {
	out := scheme
	out.LandSalesCompleted++
	out.LandFundPence += scheme.SufficientPrice.Ransom()
	out.ImportedLabourers++
	return out
}
