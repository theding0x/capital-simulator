// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 33 ("The Medium of Circulation in the Credit System").
//
// In a developed credit system the bulk of commercial transactions are settled
// by the mere crossing-out of claims — bills cancel against bills in the
// clearing house, and only the net balance requires actual money. The same
// banknote circulates many times in a day, so the effective money supply is
// actual notes × velocity. Currency is needed in full only when credit
// collapses into crisis, at which point the payment system seizes.
//
// Key invariants:
//   - ClearingHouseSettlement.MoneyUsed <= TotalClaims (clearing always saves money).
//   - ClearingHouseSettlement.NetBalance == TotalClaims - MoneyUsed.
//   - CurrencyVelocity.Velocity >= 1 (notes circulate at least once).
//   - CreditContraction.Phase.IsValid().
//   - If Phase == PhaseCrisis, CausesPaymentCrisis is forced true.
//   - If Phase != PhaseCrisis, CausesPaymentCrisis must be false.
//   - EffectiveMoneySupply.Amount == ActualNotesAmount * Velocity.
//   - EffectiveMoneySupply.Velocity >= 1.
//
// IndustrialCyclePhase, PhaseCrisis, and ErrInvalidCyclePhase are reused from
// rate_of_interest.go (Ch. 22) — they are not redeclared here.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 33 sentinel errors.
var (
	// ErrMoneyUsedExceedsClaims is returned when money_used > total_claims.
	ErrMoneyUsedExceedsClaims = errors.New("credit: money_used must be <= total_claims")

	// ErrVelocityTooLow is returned when velocity < 1.
	ErrVelocityTooLow = errors.New("credit: velocity must be >= 1")

	// ErrCrisisMustCausePaymentCrisis is returned when a non-crisis phase is
	// submitted with causes_payment_crisis == true.
	ErrCrisisMustCausePaymentCrisis = errors.New("credit: phase crisis must cause payment crisis")
)

// ── ClearingHouseSettlementID ────────────────────────────────────────────────

// ClearingHouseSettlementID identifies a persisted ClearingHouseSettlement.
// 96-bit hex from crypto/rand, matching every other domain ID in the simulator.
type ClearingHouseSettlementID string

// NewClearingHouseSettlementID returns a fresh random ClearingHouseSettlementID.
func NewClearingHouseSettlementID() ClearingHouseSettlementID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return ClearingHouseSettlementID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty ClearingHouseSettlementID.
func (id ClearingHouseSettlementID) IsZero() bool { return id == "" }

// ── ClearingHouseSettlement ──────────────────────────────────────────────────

// ClearingHouseSettlement is the persisted entity for a clearing-house
// settlement operation (Vol. III Ch. 33). Bills of exchange cancel against one
// another; only the net balance requires actual money. TotalClaims is the gross
// sum of mutual claims; MoneyUsed is the coin or notes actually transferred;
// NetBalance == TotalClaims - MoneyUsed.
//
// Invariants:
//
//	MoneyUsed <= TotalClaims
//	NetBalance == TotalClaims - MoneyUsed
type ClearingHouseSettlement struct {
	ID          ClearingHouseSettlementID `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	TotalClaims int64                     `json:"total_claims"` // pence; gross mutual claims
	MoneyUsed   int64                     `json:"money_used"`   // pence; actual money transferred
	NetBalance  int64                     `json:"net_balance"`  // pence; == TotalClaims - MoneyUsed
	CreatedAt   time.Time                 `json:"created_at"`
}

// NewClearingHouseSettlement constructs a ClearingHouseSettlement, enforcing
// invariants. ID and CreatedAt are left zero — the store assigns them.
//
//	moneyUsed > totalClaims → ErrMoneyUsedExceedsClaims
//
// Fixture (Marx, Vol. III Ch. 33): London Clearing House 1840s — bills worth
// £300 000 cancel against one another; only £1 000 in actual money crosses.
func NewClearingHouseSettlement(name, description string, totalClaims, moneyUsed int64) (ClearingHouseSettlement, error) {
	if moneyUsed > totalClaims {
		return ClearingHouseSettlement{}, ErrMoneyUsedExceedsClaims
	}
	return ClearingHouseSettlement{
		Name:        name,
		Description: description,
		TotalClaims: totalClaims,
		MoneyUsed:   moneyUsed,
		NetBalance:  totalClaims - moneyUsed,
	}, nil
}

// ── CurrencyVelocityID ───────────────────────────────────────────────────────

// CurrencyVelocityID identifies a CurrencyVelocity. It is a value type — not
// persisted — but carries an ID for referencing.
type CurrencyVelocityID string

// NewCurrencyVelocityID returns a fresh random CurrencyVelocityID.
func NewCurrencyVelocityID() CurrencyVelocityID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CurrencyVelocityID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CurrencyVelocityID.
func (id CurrencyVelocityID) IsZero() bool { return id == "" }

// ── CurrencyVelocity ─────────────────────────────────────────────────────────

// CurrencyVelocity is the value type for the number of times a banknote
// circulates in a given period. A value type — not persisted independently.
//
// Invariant: Velocity >= 1.
type CurrencyVelocity struct {
	ID          CurrencyVelocityID `json:"id"`
	Velocity    int64              `json:"velocity"` // times per period; >= 1
	Description string             `json:"description"`
}

// NewCurrencyVelocity builds and validates a CurrencyVelocity. The ID is left
// zero; callers may assign one with NewCurrencyVelocityID.
//
//	velocity < 1 → ErrVelocityTooLow
//
// Fixture (Marx, Vol. III Ch. 33): Scotland — a £1 note issued on Monday
// returns to its bank on Saturday having passed through 9 hands.
func NewCurrencyVelocity(velocity int64, description string) (CurrencyVelocity, error) {
	if velocity < 1 {
		return CurrencyVelocity{}, ErrVelocityTooLow
	}
	return CurrencyVelocity{
		Velocity:    velocity,
		Description: description,
	}, nil
}

// ── CreditContractionID ──────────────────────────────────────────────────────

// CreditContractionID identifies a CreditContraction. It is a value type — not
// persisted — but carries an ID for referencing.
type CreditContractionID string

// NewCreditContractionID returns a fresh random CreditContractionID.
func NewCreditContractionID() CreditContractionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CreditContractionID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CreditContractionID.
func (id CreditContractionID) IsZero() bool { return id == "" }

// ── CreditContraction ────────────────────────────────────────────────────────

// CreditContraction records a phase of the industrial cycle and whether that
// phase causes a payment crisis — a seizure of the payment system when credit
// collapses. In a crisis (PhaseCrisis) a payment crisis is inevitable and
// CausesPaymentCrisis is forced to true. A value type — not persisted.
//
// Invariants:
//
//	Phase.IsValid()
//	Phase == PhaseCrisis → CausesPaymentCrisis == true (forced)
//	Phase != PhaseCrisis → CausesPaymentCrisis must be false
type CreditContraction struct {
	ID                  CreditContractionID  `json:"id"`
	Phase               IndustrialCyclePhase `json:"phase"`
	CausesPaymentCrisis bool                 `json:"causes_payment_crisis"`
}

// NewCreditContraction builds and validates a CreditContraction. The ID is left
// zero; callers may assign one with NewCreditContractionID.
//
//	!phase.IsValid()                          → ErrInvalidCyclePhase
//	phase == PhaseCrisis                      → forces CausesPaymentCrisis = true
//	phase != PhaseCrisis && causesPaymentCrisis → ErrCrisisMustCausePaymentCrisis
//
// Fixture (Marx, Vol. III Ch. 33): crisis phase — the payment system seizes;
// bills of exchange are dishonoured; everyone demands actual coin.
func NewCreditContraction(phase IndustrialCyclePhase, causesPaymentCrisis bool) (CreditContraction, error) {
	if !phase.IsValid() {
		return CreditContraction{}, ErrInvalidCyclePhase
	}
	if phase == PhaseCrisis {
		causesPaymentCrisis = true
	} else if causesPaymentCrisis {
		return CreditContraction{}, ErrCrisisMustCausePaymentCrisis
	}
	return CreditContraction{
		Phase:               phase,
		CausesPaymentCrisis: causesPaymentCrisis,
	}, nil
}

// ── EffectiveMoneySupply ─────────────────────────────────────────────────────

// EffectiveMoneySupply is the compute-only (not persisted) value type for the
// effective money supply when velocity is taken into account.
//
//	Amount = ActualNotesAmount * Velocity
//
// Invariant: Velocity >= 1.
type EffectiveMoneySupply struct {
	ActualNotesAmount int64 `json:"actual_notes_amount"` // pence
	Velocity          int64 `json:"velocity"`            // >= 1
	Amount            int64 `json:"amount"`              // pence; == ActualNotesAmount * Velocity
}

// NewEffectiveMoneySupply computes the effective money supply. Not persisted —
// stateless compute only.
//
//	velocity < 1 → ErrVelocityTooLow
//
// Fixture (Marx, Vol. III Ch. 33): £500 in actual notes circulating 5 times
// per day → effective supply £2 500.
func NewEffectiveMoneySupply(actualNotesAmount, velocity int64) (EffectiveMoneySupply, error) {
	if velocity < 1 {
		return EffectiveMoneySupply{}, ErrVelocityTooLow
	}
	return EffectiveMoneySupply{
		ActualNotesAmount: actualNotesAmount,
		Velocity:          velocity,
		Amount:            actualNotesAmount * velocity,
	}, nil
}
