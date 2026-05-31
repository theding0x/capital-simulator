package credit

import (
	"testing"
)

// ── Fixtures from Capital Vol. III Ch. 28 ────────────────────────────────────
//
// Marx cites the Bank of England 1844 figures:
//   - Total notes outstanding: ~£20m (2_000_000_000 pence)
//   - Banking Department reserve: ~£8m  (800_000_000 pence)
// A plausible coin/capital split: £14m coin function, £4m capital-transfer,
// summing to £18m < £20m total (leaving headroom for reserve use).

func fixtureBoE1844(t *testing.T) CurrencyObservation {
	t.Helper()
	reserve := ReserveFund{
		Amount:      800_000_000,
		Description: "Bank of England Banking Department reserve 1844",
	}
	obs, err := NewCurrencyObservation(
		"Bank of England 1844",
		2_000_000_000,
		1_400_000_000,
		400_000_000,
		reserve,
		"Bank Charter Act 1844: £20m notes outstanding; £8m reserve; "+
			"coin function dominates (retail/wage payments); "+
			"capital transfers handled via bills and book entries",
	)
	if err != nil {
		t.Fatalf("fixtureBoE1844: %v", err)
	}
	return obs
}

// fixtureScottishBanking demonstrates the capital-transfer-dominant case:
// Scottish credit velocity meant relatively few notes circulated as coin.
func fixtureScottishBanking(t *testing.T) CurrencyObservation {
	t.Helper()
	reserve := ReserveFund{
		Amount:      5_000_000,
		Description: "Scottish bank metallic reserve, minimal due to credit velocity",
	}
	// £100 in Scotland circulates like £420 in England (Marx's observation).
	// Model: £10m outstanding, £1m coin, £7m capital-transfer.
	obs, err := NewCurrencyObservation(
		"Scottish Banking System — Credit Velocity",
		1_000_000_000,
		100_000_000,
		700_000_000,
		reserve,
		"Scottish credit system: high velocity reduces need for coin; "+
			"£100 Scottish ≈ £420 English in circulation effect; "+
			"capital-transfer function dominates",
	)
	if err != nil {
		t.Fatalf("fixtureScottishBanking: %v", err)
	}
	return obs
}

// ── NewCurrencyObservation happy paths ────────────────────────────────────────

func TestNewCurrencyObservation_BoE1844(t *testing.T) {
	t.Parallel()
	obs := fixtureBoE1844(t)
	if obs.TotalCurrencyOutstanding != 2_000_000_000 {
		t.Errorf("Total: got %d, want 2000000000", obs.TotalCurrencyOutstanding)
	}
	if obs.CoinFunctionAmount != 1_400_000_000 {
		t.Errorf("Coin: got %d, want 1400000000", obs.CoinFunctionAmount)
	}
	if obs.CapitalTransferAmount != 400_000_000 {
		t.Errorf("CapitalTransfer: got %d, want 400000000", obs.CapitalTransferAmount)
	}
	if obs.ReserveFund.Amount != 800_000_000 {
		t.Errorf("Reserve: got %d, want 800000000", obs.ReserveFund.Amount)
	}
	if obs.ID.IsZero() == false {
		// ID must remain zero until store assigns it.
		t.Error("ID must be zero before store persists")
	}
	if !obs.ID.IsZero() {
		t.Error("fresh observation must have zero ID")
	}
}

func TestNewCurrencyObservation_ScottishBanking(t *testing.T) {
	t.Parallel()
	obs := fixtureScottishBanking(t)
	if obs.TotalCurrencyOutstanding != 1_000_000_000 {
		t.Errorf("Total: got %d, want 1000000000", obs.TotalCurrencyOutstanding)
	}
	if obs.CoinFunctionAmount != 100_000_000 {
		t.Errorf("Coin: got %d, want 100000000", obs.CoinFunctionAmount)
	}
}

func TestNewCurrencyObservation_SplitEqualTotal(t *testing.T) {
	t.Parallel()
	// Boundary: coin + capital == total exactly — must succeed.
	reserve := ReserveFund{Amount: 0}
	obs, err := NewCurrencyObservation("exact", 1000, 600, 400, reserve, "")
	if err != nil {
		t.Fatalf("expected no error when split == total, got: %v", err)
	}
	if obs.CoinFunctionAmount+obs.CapitalTransferAmount != obs.TotalCurrencyOutstanding {
		t.Error("split must equal total")
	}
}

func TestNewCurrencyObservation_ZeroSplit(t *testing.T) {
	t.Parallel()
	// Boundary: both amounts zero — valid (all currency unclassified).
	reserve := ReserveFund{Amount: 0}
	_, err := NewCurrencyObservation("zero split", 1000, 0, 0, reserve, "")
	if err != nil {
		t.Fatalf("zero split must be valid: %v", err)
	}
}

// ── Invariant failures ────────────────────────────────────────────────────────

func TestNewCurrencyObservation_NonPositiveTotal(t *testing.T) {
	t.Parallel()
	reserve := ReserveFund{Amount: 0}
	_, err := NewCurrencyObservation("bad", 0, 0, 0, reserve, "")
	if err != ErrNonPositiveCurrency {
		t.Errorf("got %v, want ErrNonPositiveCurrency", err)
	}
}

func TestNewCurrencyObservation_NegativeTotal(t *testing.T) {
	t.Parallel()
	reserve := ReserveFund{Amount: 0}
	_, err := NewCurrencyObservation("bad", -1, 0, 0, reserve, "")
	if err != ErrNonPositiveCurrency {
		t.Errorf("got %v, want ErrNonPositiveCurrency", err)
	}
}

func TestNewCurrencyObservation_SplitExceedsTotal(t *testing.T) {
	t.Parallel()
	reserve := ReserveFund{Amount: 0}
	_, err := NewCurrencyObservation("bad", 1000, 700, 400, reserve, "")
	if err != ErrFunctionSplitExceedsTotal {
		t.Errorf("got %v, want ErrFunctionSplitExceedsTotal", err)
	}
}

func TestNewCurrencyObservation_NegativeReserve(t *testing.T) {
	t.Parallel()
	reserve := ReserveFund{Amount: -1}
	_, err := NewCurrencyObservation("bad", 1000, 100, 200, reserve, "")
	if err != ErrNegativeReserve {
		t.Errorf("got %v, want ErrNegativeReserve", err)
	}
}

// ── CurrencyFunction.IsValid ──────────────────────────────────────────────────

func TestCurrencyFunction_IsValid(t *testing.T) {
	t.Parallel()
	cases := []struct {
		fn   CurrencyFunction
		want bool
	}{
		{FunctionCoin, true},
		{FunctionCapitalTransfer, true},
		{0, false},
		{3, false},
		{-1, false},
	}
	for _, tc := range cases {
		got := tc.fn.IsValid()
		if got != tc.want {
			t.Errorf("CurrencyFunction(%d).IsValid() = %v, want %v", tc.fn, got, tc.want)
		}
	}
}

// ── CurrencyObservation.CoinShareBP ──────────────────────────────────────────

func TestCoinShareBP_BoE1844(t *testing.T) {
	t.Parallel()
	obs := fixtureBoE1844(t)
	// 1_400_000_000 * 10000 / 2_000_000_000 = 7000 bp (70%)
	got := obs.CoinShareBP()
	if got != 7000 {
		t.Errorf("CoinShareBP: got %d, want 7000", got)
	}
}

func TestCoinShareBP_ScottishBanking(t *testing.T) {
	t.Parallel()
	obs := fixtureScottishBanking(t)
	// 100_000_000 * 10000 / 1_000_000_000 = 1000 bp (10%)
	got := obs.CoinShareBP()
	if got != 1000 {
		t.Errorf("CoinShareBP: got %d, want 1000", got)
	}
}

func TestCoinShareBP_ZeroTotal(t *testing.T) {
	t.Parallel()
	// Guard against division by zero on a manually constructed zero struct.
	var obs CurrencyObservation
	got := obs.CoinShareBP()
	if got != 0 {
		t.Errorf("CoinShareBP with zero total: got %d, want 0", got)
	}
}

// ── ClassifyCurrency ──────────────────────────────────────────────────────────

func TestClassifyCurrency_CoinDominant(t *testing.T) {
	t.Parallel()
	obs := fixtureBoE1844(t)
	got := ClassifyCurrency(obs)
	want := "predominantly coin function: currency serves mainly revenue expenditure (retail/consumer)"
	if got != want {
		t.Errorf("ClassifyCurrency BoE1844: got %q, want %q", got, want)
	}
}

func TestClassifyCurrency_CapitalTransferDominant(t *testing.T) {
	t.Parallel()
	obs := fixtureScottishBanking(t)
	got := ClassifyCurrency(obs)
	want := "predominantly capital-transfer function: currency moves productive capital between capitalists"
	if got != want {
		t.Errorf("ClassifyCurrency Scottish: got %q, want %q", got, want)
	}
}

func TestClassifyCurrency_Mixed(t *testing.T) {
	t.Parallel()
	reserve := ReserveFund{Amount: 0}
	// 50% coin = 5000 bp → mixed
	obs, err := NewCurrencyObservation("mixed", 1000, 500, 300, reserve, "")
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	got := ClassifyCurrency(obs)
	want := "mixed function: currency serves both revenue expenditure and capital transfer"
	if got != want {
		t.Errorf("ClassifyCurrency mixed: got %q, want %q", got, want)
	}
}

// ── NewTookeFullartonAnalysis ─────────────────────────────────────────────────

func TestNewTookeFullartonAnalysis(t *testing.T) {
	t.Parallel()
	a := NewTookeFullartonAnalysis()
	if len(a.ConflatedCategories) != 3 {
		t.Errorf("ConflatedCategories: got %d, want 3", len(a.ConflatedCategories))
	}
	if a.CorrectDistinction == "" {
		t.Error("CorrectDistinction must not be empty")
	}
}

// ── NewCurrencyObservationID ──────────────────────────────────────────────────

func TestNewCurrencyObservationID_Format(t *testing.T) {
	t.Parallel()
	id := NewCurrencyObservationID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
}

func TestNewCurrencyObservationID_Unique(t *testing.T) {
	t.Parallel()
	a := NewCurrencyObservationID()
	b := NewCurrencyObservationID()
	if a == b {
		t.Error("two IDs must not be equal")
	}
}

func TestCurrencyObservationID_IsZero(t *testing.T) {
	t.Parallel()
	var id CurrencyObservationID
	if !id.IsZero() {
		t.Error("zero value must be zero")
	}
	id = NewCurrencyObservationID()
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}
