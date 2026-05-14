// Package labour holds the canonical labour-time and productivity types
// shared by every chapter implemented in simulation-engine. It exists so
// that Ch. 11 (surplus), Ch. 12 (production / relative surplus-value),
// and Ch. 15 (machinery) can trade values without explicit casts.
//
// Marx, Capital Vol. I, Ch. 1 §1: "The labour-time socially necessary is
// that required to produce an article under the normal conditions of
// production, and with the average degree of skill and intensity prevalent
// at the time." LabourMinutes is the integer-minute representation of
// that magnitude. The unit is deliberately fine-grained (minute, not
// hour) so the working-day arithmetic of Ch. 10 lands on integers.
package labour

// LabourMinutes is the canonical value-magnitude unit. All chapters that
// reason about value, wear-and-tear, working-day length, surplus, or
// productivity should use this type rather than redefining their own.
type LabourMinutes int64

// ProductivityFactor is a multiplicative scalar on socially-necessary
// labour-time: value ∝ 1 / productivity (Ch. 12 §1). A factor of 2.0
// means each unit of product now embodies half the SNLT it previously
// did, freeing the rest of the working day for surplus labour. The
// factor is the lingua franca by which Ch. 13 (co-operation), Ch. 14
// (manufacture), and Ch. 15 (machinery) all feed Ch. 12's relative
// surplus-value calculator.
type ProductivityFactor float64
