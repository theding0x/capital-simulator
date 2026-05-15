package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// PieceWageID is a 96-bit hex identifier for PieceWage records.
type PieceWageID string

func (id PieceWageID) IsZero() bool { return id == "" }

func NewPieceWageID() PieceWageID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return PieceWageID(hex.EncodeToString(b))
}

// SubContractID is a 96-bit hex identifier for SubContract records.
type SubContractID string

func (id SubContractID) IsZero() bool { return id == "" }

func NewSubContractID() SubContractID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return SubContractID(hex.EncodeToString(b))
}

// QualityOutcome records whether a piece-work session was accepted or rejected.
type QualityOutcome string

const (
	QualityAccepted QualityOutcome = "accepted"
	QualityRejected QualityOutcome = "rejected"
)

// PieceWage is a piece-rate wage contract: the price per accepted piece and the
// number of pieces socially expected in one working day.
type PieceWage struct {
	ID           PieceWageID `json:"id"`
	AgentID      AgentID     `json:"agent_id"`
	WageFormID   WageFormID  `json:"wage_form_id,omitempty"`
	PricePence   int64       `json:"price_pence"`
	NormalOutput int64       `json:"normal_output"`
	CreatedAt    time.Time   `json:"created_at"`
}

// PieceSession holds the output of a piece-work session for earnings calculation.
type PieceSession struct {
	PiecesProduced int64          `json:"pieces_produced"`
	QualityOutcome QualityOutcome `json:"quality_outcome"`
}

// SubContract models the sweating system: a head labourer contracts with the
// capitalist at PieceRatePence and pays assistants AssistantRatePence, keeping
// the spread.
type SubContract struct {
	ID                 SubContractID `json:"id"`
	HeadLabourerID     AgentID       `json:"head_labourer_id"`
	AssistantIDs       []AgentID     `json:"assistant_ids"`
	PieceRatePence     int64         `json:"piece_rate_pence"`
	AssistantRatePence int64         `json:"assistant_rate_pence"`
	CreatedAt          time.Time     `json:"created_at"`
}

var (
	ErrPieceWageNormalOutput = errors.New("piece_wage: normal_output must be positive")
	ErrPieceWagePricePence   = errors.New("piece_wage: price_pence must be positive")
	ErrSubContractNoHead     = errors.New("sub_contract: head_labourer_id is required")
	ErrSubContractSpread     = errors.New("sub_contract: piece_rate_pence must exceed assistant_rate_pence")
)

func (pw PieceWage) Validate() error {
	if pw.AgentID.IsZero() {
		return errors.New("piece_wage: agent_id is required")
	}
	if pw.WageFormID.IsZero() && pw.PricePence <= 0 {
		return ErrPieceWagePricePence
	}
	if pw.NormalOutput <= 0 {
		return ErrPieceWageNormalOutput
	}
	return nil
}

func (sc SubContract) Validate() error {
	if sc.HeadLabourerID.IsZero() {
		return ErrSubContractNoHead
	}
	if sc.PieceRatePence <= sc.AssistantRatePence {
		return ErrSubContractSpread
	}
	return nil
}

// ComputePiecePrice returns the piece price as dailyWage / normalOutput.
// Use farthings (¼d.) as the monetary unit so the invariant
// ComputePiecePrice(144, 24) == 6 holds exactly in integer arithmetic.
func ComputePiecePrice(dailyWage, normalOutput int64) int64 {
	return dailyWage / normalOutput
}

// ComputePieceValue returns the piece value as dayValueProduct / normalOutput.
// Always exceeds ComputePiecePrice when dayValueProduct > dailyWage because
// the day's value product includes the unpaid surplus labour.
func ComputePieceValue(dayValueProduct, normalOutput int64) int64 {
	return dayValueProduct / normalOutput
}

// ComputePieceEarnings returns total earnings for a session: accepted pieces
// multiplied by the piece price. Rejected sessions earn nothing — quality
// control is enforced by the wage form itself.
func ComputePieceEarnings(s PieceSession, pw PieceWage) int64 {
	if s.QualityOutcome != QualityAccepted {
		return 0
	}
	return s.PiecesProduced * pw.PricePence
}

// SubContractSpread returns the earnings per piece the head labourer keeps
// after paying assistants — the "gain of middlemen" Marx identifies.
func SubContractSpread(sc SubContract) int64 {
	return sc.PieceRatePence - sc.AssistantRatePence
}
