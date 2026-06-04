package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ch. 31 — Genesis of the Industrial Capitalist.
//
// The industrial capitalist did not emerge through thrift. Its genesis lay in
// colonial plunder, the slave trade, national debt, and the system of
// protection — each a form of primitive accumulation operating through the
// state and world market rather than the domestic labour process.

// CapitalOriginID is a 96-bit hex identifier for a capital origin record.
type CapitalOriginID string

func (id CapitalOriginID) IsZero() bool { return id == "" }

func NewCapitalOriginID() CapitalOriginID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return CapitalOriginID(hex.EncodeToString(b))
}

// CapitalOrigin records a single mechanism through which initial industrial
// capital was formed. Source is one of: "usury", "commerce", "colonial-plunder",
// "national-debt", "taxation", "guild-master-accumulation".
type CapitalOrigin struct {
	ID                CapitalOriginID   `json:"id"`
	HistoricalStageID HistoricalStageID `json:"historical_stage_id"`
	Source            string            `json:"source"`
	AmountPence       Pence             `json:"amount_pence"`
	Period            string            `json:"period"`
	CreatedAt         time.Time         `json:"created_at"`
}

var (
	ErrCapitalOriginSourceRequired = errors.New("simulation: capital origin source is required")
	ErrCapitalOriginAmountPositive = errors.New("simulation: amount_pence must be positive")
)

func (c CapitalOrigin) Validate() error {
	if c.Source == "" {
		return ErrCapitalOriginSourceRequired
	}
	if c.AmountPence <= 0 {
		return ErrCapitalOriginAmountPositive
	}
	return nil
}

// ColonialTransferID is a 96-bit hex identifier for a colonial transfer record.
type ColonialTransferID string

func (id ColonialTransferID) IsZero() bool { return id == "" }

func NewColonialTransferID() ColonialTransferID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return ColonialTransferID(hex.EncodeToString(b))
}

// ColonialTransfer records looted wealth flowing from periphery to metropole —
// the most direct expression of primitive accumulation through the world market.
type ColonialTransfer struct {
	ID                ColonialTransferID `json:"id"`
	HistoricalStageID HistoricalStageID  `json:"historical_stage_id"`
	From              string             `json:"from"`
	To                string             `json:"to"`
	ValuePence        Pence              `json:"value_pence"`
	Method            string             `json:"method"`
	CreatedAt         time.Time          `json:"created_at"`
}

var (
	ErrColonialTransferFromRequired   = errors.New("simulation: from is required")
	ErrColonialTransferToRequired     = errors.New("simulation: to is required")
	ErrColonialTransferValuePositive  = errors.New("simulation: value_pence must be positive")
	ErrColonialTransferMethodRequired = errors.New("simulation: method is required")
)

func (t ColonialTransfer) Validate() error {
	if t.From == "" {
		return ErrColonialTransferFromRequired
	}
	if t.To == "" {
		return ErrColonialTransferToRequired
	}
	if t.ValuePence <= 0 {
		return ErrColonialTransferValuePositive
	}
	if t.Method == "" {
		return ErrColonialTransferMethodRequired
	}
	return nil
}

// NationalDebtID is a 96-bit hex identifier for a national debt record.
type NationalDebtID string

func (id NationalDebtID) IsZero() bool { return id == "" }

func NewNationalDebtID() NationalDebtID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return NationalDebtID(hex.EncodeToString(b))
}

// NationalDebt records public debt as an instrument of primitive accumulation —
// the state borrows from the propertied class, committing future tax revenue
// (extorted from wage-labourers) to service the interest.
type NationalDebt struct {
	ID                NationalDebtID    `json:"id"`
	HistoricalStageID HistoricalStageID `json:"historical_stage_id"`
	AmountPence       Pence             `json:"amount_pence"`
	InterestRateBps   int64             `json:"interest_rate_bps"`
	CreditorClass     string            `json:"creditor_class"`
	CreatedAt         time.Time         `json:"created_at"`
}

var (
	ErrNationalDebtAmountPositive       = errors.New("simulation: amount_pence must be positive")
	ErrNationalDebtInterestRatePositive = errors.New("simulation: interest_rate_bps must be positive")
	ErrNationalDebtCreditorRequired     = errors.New("simulation: creditor_class is required")
)

func (d NationalDebt) Validate() error {
	if d.AmountPence <= 0 {
		return ErrNationalDebtAmountPositive
	}
	if d.InterestRateBps <= 0 {
		return ErrNationalDebtInterestRatePositive
	}
	if d.CreditorClass == "" {
		return ErrNationalDebtCreditorRequired
	}
	return nil
}

// ProtectionSystemID is a 96-bit hex identifier for a protection system record.
type ProtectionSystemID string

func (id ProtectionSystemID) IsZero() bool { return id == "" }

func NewProtectionSystemID() ProtectionSystemID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return ProtectionSystemID(hex.EncodeToString(b))
}

// ProtectionSystem records state-backed capital formation through tariffs and
// trade restrictions — "an artificial means of manufacturing manufacturers."
type ProtectionSystem struct {
	ID                ProtectionSystemID `json:"id"`
	HistoricalStageID HistoricalStageID  `json:"historical_stage_id"`
	TariffRateBps     int64              `json:"tariff_rate_bps"`
	Beneficiary       string             `json:"beneficiary"`
	PeriodStart       string             `json:"period_start"`
	PeriodEnd         string             `json:"period_end"`
	CreatedAt         time.Time          `json:"created_at"`
}

var (
	ErrProtectionSystemBeneficiaryRequired = errors.New("simulation: protection system beneficiary is required")
	ErrProtectionSystemTariffNegative      = errors.New("simulation: tariff_rate_bps cannot be negative")
)

// Validate enforces invariants on a protection system record: a beneficiary must
// be named and the tariff rate cannot be negative.
func (s ProtectionSystem) Validate() error {
	if s.Beneficiary == "" {
		return ErrProtectionSystemBeneficiaryRequired
	}
	if s.TariffRateBps < 0 {
		return ErrProtectionSystemTariffNegative
	}
	return nil
}

// IndustrialCapitalGenesis aggregates all primitive-accumulation mechanisms
// that produced the industrial capitalist for a given historical stage.
//
// Invariant: TotalCapitalFormedPence == sum(Origins[*].AmountPence) +
// sum(ColonialTransfers[*].ValuePence).
type IndustrialCapitalGenesis struct {
	HistoricalStageID       HistoricalStageID  `json:"historical_stage_id"`
	Origins                 []CapitalOrigin    `json:"origins"`
	ColonialTransfers       []ColonialTransfer `json:"colonial_transfers"`
	NationalDebts           []NationalDebt     `json:"national_debts"`
	ProtectionSystems       []ProtectionSystem `json:"protection_systems"`
	TotalCapitalFormedPence Pence              `json:"total_capital_formed_pence"`
}

// ComputeGenesis assembles an IndustrialCapitalGenesis from its constituent
// parts. TotalCapitalFormedPence is the sum of capital origin amounts and
// colonial transfer values; national debts and protection systems are
// structural levers, not direct capital stocks.
func ComputeGenesis(stageID HistoricalStageID, origins []CapitalOrigin, transfers []ColonialTransfer, debts []NationalDebt, systems []ProtectionSystem) IndustrialCapitalGenesis {
	if origins == nil {
		origins = []CapitalOrigin{}
	}
	if transfers == nil {
		transfers = []ColonialTransfer{}
	}
	if debts == nil {
		debts = []NationalDebt{}
	}
	if systems == nil {
		systems = []ProtectionSystem{}
	}
	var total Pence
	for _, o := range origins {
		total += o.AmountPence
	}
	for _, t := range transfers {
		total += t.ValuePence
	}
	return IndustrialCapitalGenesis{
		HistoricalStageID:       stageID,
		Origins:                 origins,
		ColonialTransfers:       transfers,
		NationalDebts:           debts,
		ProtectionSystems:       systems,
		TotalCapitalFormedPence: total,
	}
}
