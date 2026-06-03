package agent

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"math"
	"time"
)

// PiecePricePointID is a 96-bit hex identifier for PiecePricePoint records.
type PiecePricePointID string

func (id PiecePricePointID) IsZero() bool { return id == "" }

// NewPiecePricePointID returns a random 96-bit hex PiecePricePointID.
func NewPiecePricePointID() PiecePricePointID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return PiecePricePointID(hex.EncodeToString(b))
}

// PiecePricePoint is one observation in a piece-wage's price-over-productivity
// time series (Ch. 21): "the price of labour ... falls in the same proportion
// as the productiveness of labour rises." As productivity rises the daily wage
// is unchanged and the socially-normal output grows, so the price per piece
// falls in inverse proportion. The simulation-engine tick scheduler appends one
// point each time the Ch. 15 factory productivity factor changes (#219).
type PiecePricePoint struct {
	ID                 PiecePricePointID `json:"id"`
	AgentID            AgentID           `json:"agent_id"`
	Sequence           int64             `json:"sequence"`
	ProductivityFactor float64           `json:"productivity_factor"`
	PricePence         int64             `json:"price_pence"`
	NormalOutput       int64             `json:"normal_output"`
	OccurredAt         time.Time         `json:"occurred_at"`
}

// RepricePieceWage applies a productivity rise to a piece-wage and returns the
// resulting price and normal output, holding the implied daily wage constant —
// Marx's Ch. 21 law of piece-wages. Given the baseline implied daily wage
// W = PricePence * NormalOutput, a productivity factor k raises the normal
// output to round(NormalOutput * k) and lowers the price per piece to
// W / newOutput. A factor <= 0 or == 1 is a no-op; a factor of 2 halves the
// price (seed Spinning Mill: 6f at 24 pieces -> 3f at 48 pieces); a factor of
// 1.5 lowers it to two-thirds (6f -> 4f).
//
// The baseline PieceWage is never mutated: callers reprice from the immutable
// baseline each pass, so successive productivity factors do not compound.
func RepricePieceWage(baseline PieceWage, productivityFactor float64) PieceWage {
	if productivityFactor <= 0 {
		return baseline
	}
	dailyWage := baseline.PricePence * baseline.NormalOutput
	newOutput := int64(math.Round(float64(baseline.NormalOutput) * productivityFactor))
	if newOutput < 1 {
		newOutput = 1
	}
	out := baseline
	out.NormalOutput = newOutput
	out.PricePence = dailyWage / newOutput
	return out
}
