package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ch. 29 — Genesis of the Capitalist Farmer.
//
// The capitalist farmer emerged through three successive tenant forms — bailiff,
// metayer, capitalist farmer — over the 15th–17th centuries. Two forces drove the
// transformation: long farm leases that fixed nominal rents in money, and the
// 16th-century depreciation of the precious metals which eroded the real burden
// of those fixed rents. The farmer kept the difference between rising corn prices
// and a flat nominal rent obligation, converting it into capital without any
// productive act on their part.

// TenantForm names a stage in the historical evolution of the agricultural tenant.
type TenantForm string

const (
	TenantFormBailiff          TenantForm = "bailiff"
	TenantFormMetayer          TenantForm = "metayer"
	TenantFormCapitalistFarmer TenantForm = "capitalist-farmer"
)

// Order returns the tenant form's position in the historical sequence
// bailiff (1) → metayer (2) → capitalist-farmer (3). The genesis of the
// capitalist farmer is a one-way progression through these forms (L-H2); Order
// makes that ordering comparable. An unrecognised form returns 0, sorting before
// all defined forms.
func (f TenantForm) Order() int {
	switch f {
	case TenantFormBailiff:
		return 1
	case TenantFormMetayer:
		return 2
	case TenantFormCapitalistFarmer:
		return 3
	default:
		return 0
	}
}

// HighestTenantForm returns the highest Order() among the given tenant forms, or
// 0 when there are none — the high-water mark of a stage's progression through
// bailiff → metayer → capitalist-farmer.
func HighestTenantForm(forms []TenantForm) int {
	highest := 0
	for _, f := range forms {
		if o := f.Order(); o > highest {
			highest = o
		}
	}
	return highest
}

// CheckTenantFormProgression enforces L-H2: the genesis of the capitalist farmer
// is a one-way progression, so a stage that has already reached a given tenant
// form cannot afterwards record a lower one. It returns ErrFarmTenureFormRegression
// when next.Order() is below the highest form already present; equal or higher
// forms are allowed (a stage may record several tenures at its current form and
// may advance).
func CheckTenantFormProgression(existing []TenantForm, next TenantForm) error {
	if next.Order() < HighestTenantForm(existing) {
		return ErrFarmTenureFormRegression
	}
	return nil
}

var (
	ErrFarmTenureFormRequired        = errors.New("simulation: farm tenure form is required")
	ErrFarmTenureFormInvalid         = errors.New("simulation: form must be bailiff, metayer, or capitalist-farmer")
	ErrFarmTenureLeasePeriodRequired = errors.New("simulation: lease_period_years must be positive")
	ErrFarmTenureRentNegative        = errors.New("simulation: rent_pence cannot be negative")
	ErrFarmTenureCapitalNegative     = errors.New("simulation: capital_advanced_pence cannot be negative")
	ErrDepreciationFactorInvalid     = errors.New("simulation: depreciation_factor must be in (0, 1]")
	ErrFarmTenureFormRegression      = errors.New("simulation: farm tenure form cannot regress below an existing form for the stage (L-H2 is one-way)")
)

// FarmTenureID is a 96-bit hex identifier for a persisted farm tenure record.
type FarmTenureID string

// IsZero reports whether the ID is unset.
func (id FarmTenureID) IsZero() bool { return id == "" }

// NewFarmTenureID generates a 96-bit hex identifier via crypto/rand.
func NewFarmTenureID() FarmTenureID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return FarmTenureID(hex.EncodeToString(b))
}

// FarmTenure records a single tenure arrangement between an agricultural tenant
// and a landlord. The combination of a long lease, a fixed nominal rent, and
// rising corn prices is the mechanism by which the tenant farmer accumulated
// capital — Marx's "genesis of the capitalist farmer".
type FarmTenure struct {
	ID                   FarmTenureID      `json:"id"`
	HistoricalStageID    HistoricalStageID `json:"historical_stage_id"`
	Form                 TenantForm        `json:"form"`
	LeasePeriodYears     int64             `json:"lease_period_years"`
	RentPence            Pence             `json:"rent_pence"`
	CapitalAdvancedPence Pence             `json:"capital_advanced_pence"`
	RevenuePence         Pence             `json:"revenue_pence"`
	WageCostsPence       Pence             `json:"wage_costs_pence"`
	CreatedAt            time.Time         `json:"created_at"`
}

// Validate enforces invariants on a farm tenure record.
func (f FarmTenure) Validate() error {
	if f.Form == "" {
		return ErrFarmTenureFormRequired
	}
	if f.Form != TenantFormBailiff && f.Form != TenantFormMetayer && f.Form != TenantFormCapitalistFarmer {
		return ErrFarmTenureFormInvalid
	}
	if f.LeasePeriodYears <= 0 {
		return ErrFarmTenureLeasePeriodRequired
	}
	if f.RentPence < 0 {
		return ErrFarmTenureRentNegative
	}
	if f.CapitalAdvancedPence < 0 {
		return ErrFarmTenureCapitalNegative
	}
	return nil
}

// MoneyDepreciation captures the effect of 16th-century silver depreciation on
// a nominally fixed rent. When currency depreciated, the farmer's real rent
// burden fell while nominal corn revenues rose — a windfall that swelled profit
// without any action on the farmer's part.
type MoneyDepreciation struct {
	NominalRentPence   Pence   `json:"nominal_rent_pence"`
	RealRentPence      Pence   `json:"real_rent_pence"`
	DepreciationFactor float64 `json:"depreciation_factor"`
}

// ComputeRealRent deflates nominalRent by md.DepreciationFactor.
// A factor of 0.5 means the currency halved in value over the lease term,
// so the real rent obligation is half the nominal figure.
func ComputeRealRent(nominalRent Pence, md MoneyDepreciation) Pence {
	return Pence(int64(float64(nominalRent) * md.DepreciationFactor))
}

// FarmingSurplus records the capitalist farmer's appropriated surplus over a
// lease period. Profit = Revenue − NominalRent − WageCosts; over a long lease
// with a fixed nominal rent and rising commodity prices, Profit swells without
// any additional capital investment — the primitive accumulation of the farmer.
type FarmingSurplus struct {
	Revenue     Pence `json:"revenue"`
	NominalRent Pence `json:"nominal_rent"`
	WageCosts   Pence `json:"wage_costs"`
	Profit      Pence `json:"profit"`
}

// ComputeFarmingSurplus derives the farmer's surplus from revenue, rent, and
// wage costs. Profit can be negative if costs exceed revenue.
func ComputeFarmingSurplus(revenue, nominalRent, wageCosts Pence) FarmingSurplus {
	return FarmingSurplus{
		Revenue:     revenue,
		NominalRent: nominalRent,
		WageCosts:   wageCosts,
		Profit:      revenue - nominalRent - wageCosts,
	}
}
