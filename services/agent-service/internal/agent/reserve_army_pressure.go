package agent

import "errors"

// ReserveArmyPressure models the downward pull that the industrial reserve
// army (Ch. 25 — the general law of capitalist accumulation) exerts on the
// wage. Marx's argument, which Ch. 12 hands forward to the wages chapters in
// its competition footnotes: the relative surplus-population is "the
// background upon which the law of the demand and supply of labour does its
// work" — as the reserve grows relative to the active army, competition among
// workers bids the price of labour down below the value of labour-power.
//
// Magnitude is the head-count of the reserve army (the unemployed and
// half-employed); WorkforceTotal is the whole labour force it is measured
// against (active + reserve). Magnitude / WorkforceTotal is the same ratio the
// Ch. 25 type IndustrialReserveArmy carries as RelativeProportion, and
// Magnitude corresponds to that type's Size. The pressure is a pure, stateless
// input threaded into the wage-form scenarios (Ch. 17, 19, 20, 21); it is
// never persisted.
type ReserveArmyPressure struct {
	Magnitude      int64 `json:"magnitude"`
	WorkforceTotal int64 `json:"workforce_total"`
}

// ErrReserveArmyExceedsWorkforce is returned when the reserve army is larger
// than the total workforce it is measured against — an impossible partition.
var ErrReserveArmyExceedsWorkforce = errors.New("agent: reserve army magnitude cannot exceed workforce total")

// IsActive reports whether the pressure compresses wages at all. A zero
// magnitude (full employment) or an unset workforce leaves the wage untouched.
func (p ReserveArmyPressure) IsActive() bool {
	return p.Magnitude > 0 && p.WorkforceTotal > 0
}

// Validate enforces the partition 0 <= Magnitude <= WorkforceTotal. A pressure
// with a zero workforce is treated as "no pressure" and is valid (IsActive is
// false); a magnitude exceeding the workforce is rejected.
func (p ReserveArmyPressure) Validate() error {
	if p.Magnitude < 0 || p.WorkforceTotal < 0 {
		return errors.New("agent: reserve army magnitude and workforce total cannot be negative")
	}
	if p.Magnitude > p.WorkforceTotal {
		return ErrReserveArmyExceedsWorkforce
	}
	return nil
}

// PressureBP returns the compression intensity in basis points:
// Magnitude / WorkforceTotal expressed in ten-thousandths, clamped to
// [0, 10000]. 0 bp is full employment (no compression); 10000 bp is a reserve
// as large as the whole workforce (the wage is driven to its floor). Inactive
// pressure is 0.
func (p ReserveArmyPressure) PressureBP() int64 {
	if !p.IsActive() {
		return 0
	}
	bp := p.Magnitude * 10000 / p.WorkforceTotal
	if bp < 0 {
		return 0
	}
	if bp > 10000 {
		return 10000
	}
	return bp
}

// CompressPence returns the money-wage after reserve-army compression:
// nominal × (1 − PressureBP/10000), rounded half up and floored at zero. A
// non-positive nominal, or inactive pressure, returns the nominal unchanged.
func (p ReserveArmyPressure) CompressPence(nominal int64) int64 {
	return p.compress(nominal)
}

// CompressMinutes is the value-magnitude analogue of CompressPence: it
// compresses a wage expressed in LabourMinutes (e.g. the Ch. 17 value of
// labour-power) by the same proportion, showing the price of labour pushed
// below its value.
func (p ReserveArmyPressure) CompressMinutes(nominal LabourMinutes) LabourMinutes {
	return LabourMinutes(p.compress(int64(nominal)))
}

// compress applies the proportional reserve-army discount to a non-negative
// magnitude expressed in any unit (pence or labour-minutes).
func (p ReserveArmyPressure) compress(nominal int64) int64 {
	if nominal <= 0 || !p.IsActive() {
		return nominal
	}
	retained := 10000 - p.PressureBP()
	return roundHalfUpDiv(nominal*retained, 10000)
}

// roundHalfUpDiv divides num by den (den > 0, num >= 0), rounding halves up.
func roundHalfUpDiv(num, den int64) int64 {
	return (num + den/2) / den
}
