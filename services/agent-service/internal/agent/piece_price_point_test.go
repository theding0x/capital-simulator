package agent

import "testing"

func TestRepricePieceWage(t *testing.T) {
	t.Parallel()
	// Seed Spinning Mill fixture (Ch. 21): 1½d. = 6 farthings per piece at the
	// socially-normal 24 pieces a day -> implied daily wage 144 farthings (3s.).
	baseline := PieceWage{PricePence: 6, NormalOutput: 24}

	tests := []struct {
		name       string
		factor     float64
		wantPrice  int64
		wantOutput int64
	}{
		{"unchanged productivity is a no-op", 1.0, 6, 24},
		{"double productivity halves the price", 2.0, 3, 48},
		{"one-and-a-half productivity lowers price to two-thirds", 1.5, 4, 36},
		{"sixfold productivity drops the price to a single farthing", 6.0, 1, 144},
		{"non-positive factor returns the baseline", 0, 6, 24},
		{"negative factor returns the baseline", -2.0, 6, 24},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := RepricePieceWage(baseline, tt.factor)
			if got.PricePence != tt.wantPrice {
				t.Errorf("price: got %d, want %d", got.PricePence, tt.wantPrice)
			}
			if got.NormalOutput != tt.wantOutput {
				t.Errorf("output: got %d, want %d", got.NormalOutput, tt.wantOutput)
			}
		})
	}
}

// The wage form's central claim: the day-wage is unchanged as productivity rises;
// only the price per piece falls. price * output must stay at the baseline 144f.
func TestRepricePieceWageHoldsDailyWageConstant(t *testing.T) {
	t.Parallel()
	baseline := PieceWage{PricePence: 6, NormalOutput: 24}
	const wantDailyWage = 6 * 24 // 144 farthings
	for _, factor := range []float64{1.0, 1.5, 2.0, 3.0, 6.0} {
		got := RepricePieceWage(baseline, factor)
		if dw := got.PricePence * got.NormalOutput; dw != wantDailyWage {
			t.Errorf("factor %v: implied daily wage = %d, want %d", factor, dw, wantDailyWage)
		}
	}
}

func TestNewPiecePricePointID(t *testing.T) {
	t.Parallel()
	id := NewPiecePricePointID()
	if id.IsZero() {
		t.Fatal("expected non-zero id")
	}
	if len(string(id)) != 24 {
		t.Errorf("expected 24 hex chars, got %d", len(string(id)))
	}
	if id == NewPiecePricePointID() {
		t.Error("expected unique ids")
	}
}
