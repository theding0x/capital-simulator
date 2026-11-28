package commodity_test

import (
	"errors"
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// shilling is the seed scale used throughout the Ch.9 fixtures: 1s = 120
// labour-minutes (see minutesPerPound).
const shilling commodity.LabourMinutes = 120

const partsEpsilon = 1e-9

// Ch.9 §2: "30s. value of yarn = 24s. const. + 3s. var. + 3s. surpl." in 20
// lbs. of yarn. "8/10ths of this value, or the 24s. ... is contained in
// 8/10ths of the product or in 16 lbs. of yarn", and the remaining 4 lbs. hold
// the newly created value, half of it the variable part and half the surplus.
func TestApportionProduct_Ch9YarnSixteenTwoTwo(t *testing.T) {
	t.Parallel()
	account := commodity.ProductionAccount{
		Constant: 24 * shilling,
		Variable: 3 * shilling,
		Surplus:  commodity.SurplusValue(3 * shilling),
	}
	parts, err := commodity.ApportionProduct(account, 20)
	if err != nil {
		t.Fatalf("ApportionProduct: %v", err)
	}
	for _, c := range []struct {
		name      string
		got, want commodity.Quantity
	}{
		{"constant", parts.Constant, 16},
		{"variable", parts.Variable, 2},
		{"surplus", parts.Surplus, 2},
	} {
		if math.Abs(float64(c.got-c.want)) > partsEpsilon {
			t.Errorf("%s part = %v lbs, want %v", c.name, c.got, c.want)
		}
	}
}

// The ratio is scale-free: the same yarn counted in shillings or in labour
// minutes apportions identically.
func TestApportionProduct_ScaleFree(t *testing.T) {
	t.Parallel()
	inShillings := commodity.ProductionAccount{Constant: 24, Variable: 3, Surplus: 3}
	inMinutes := commodity.ProductionAccount{
		Constant: 24 * shilling, Variable: 3 * shilling, Surplus: commodity.SurplusValue(3 * shilling),
	}
	a, err := commodity.ApportionProduct(inShillings, 20)
	if err != nil {
		t.Fatalf("shillings: %v", err)
	}
	b, err := commodity.ApportionProduct(inMinutes, 20)
	if err != nil {
		t.Fatalf("minutes: %v", err)
	}
	if math.Abs(float64(a.Constant-b.Constant)) > partsEpsilon ||
		math.Abs(float64(a.Variable-b.Variable)) > partsEpsilon ||
		math.Abs(float64(a.Surplus-b.Surplus)) > partsEpsilon {
		t.Errorf("scale changed the split: %+v vs %+v", a, b)
	}
}

// Ch.9 §1, spinning mill 1871: 10,000 lbs. of yarn at £510 = £378 const. +
// £52 var. + £80 surpl. The parts of the product must be proportional to the
// parts of the value, and together make up the whole product.
func TestApportionProduct_SpinningMill1871_PartsProportionalAndExhaustive(t *testing.T) {
	t.Parallel()
	account := commodity.ProductionAccount{
		Constant: 378 * minutesPerPound,
		Variable: 52 * minutesPerPound,
		Surplus:  commodity.SurplusValue(80 * minutesPerPound),
	}
	const total commodity.Quantity = 10_000
	parts, err := commodity.ApportionProduct(account, total)
	if err != nil {
		t.Fatalf("ApportionProduct: %v", err)
	}
	if sum := parts.Constant + parts.Variable + parts.Surplus; math.Abs(float64(sum-total)) > 1e-6 {
		t.Errorf("c+v+s parts = %v lbs, want the whole product %v", sum, total)
	}
	whole := float64(account.ExpandedCapital())
	for _, c := range []struct {
		name  string
		part  commodity.Quantity
		value float64
	}{
		{"constant", parts.Constant, float64(account.Constant)},
		{"variable", parts.Variable, float64(account.Variable)},
		{"surplus", parts.Surplus, float64(account.Surplus)},
	} {
		if got, want := float64(c.part)/float64(total), c.value/whole; math.Abs(got-want) > partsEpsilon {
			t.Errorf("%s: share of product = %v, share of value = %v", c.name, got, want)
		}
	}
}

// Whatever the account, the three parts exhaust the product and none is
// negative.
func TestApportionProduct_PartitionInvariant(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		account commodity.ProductionAccount
		total   commodity.Quantity
	}{
		{"yarn", commodity.ProductionAccount{Constant: 24, Variable: 3, Surplus: 3}, 20},
		{"no constant capital", commodity.ProductionAccount{Variable: 3, Surplus: 3}, 7},
		{"no surplus", commodity.ProductionAccount{Constant: 24, Variable: 3}, 20},
		{"constant only", commodity.ProductionAccount{Constant: 5}, 9},
		{"fractional total", commodity.ProductionAccount{Constant: 410, Variable: 90, Surplus: 90}, 0.125},
		{"zero total", commodity.ProductionAccount{Constant: 24, Variable: 3, Surplus: 3}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parts, err := commodity.ApportionProduct(tc.account, tc.total)
			if err != nil {
				t.Fatalf("ApportionProduct: %v", err)
			}
			if parts.Constant < 0 || parts.Variable < 0 || parts.Surplus < 0 {
				t.Errorf("negative part: %+v", parts)
			}
			if sum := parts.Constant + parts.Variable + parts.Surplus; math.Abs(float64(sum-tc.total)) > 1e-9 {
				t.Errorf("parts sum to %v, want %v", sum, tc.total)
			}
		})
	}
}

func TestApportionProduct_Rejects(t *testing.T) {
	t.Parallel()
	valid := commodity.ProductionAccount{Constant: 24, Variable: 3, Surplus: 3}
	t.Run("empty account", func(t *testing.T) {
		t.Parallel()
		_, err := commodity.ApportionProduct(commodity.ProductionAccount{}, 20)
		if !errors.Is(err, commodity.ErrEmptyAccount) {
			t.Errorf("err = %v, want ErrEmptyAccount", err)
		}
	})
	t.Run("negative total product", func(t *testing.T) {
		t.Parallel()
		_, err := commodity.ApportionProduct(valid, -1)
		if !errors.Is(err, commodity.ErrNegativeProduct) {
			t.Errorf("err = %v, want ErrNegativeProduct", err)
		}
	})
	t.Run("not a number", func(t *testing.T) {
		t.Parallel()
		if _, err := commodity.ApportionProduct(valid, commodity.Quantity(math.NaN())); err == nil {
			t.Error("a NaN product was accepted")
		}
	})
	t.Run("malformed account", func(t *testing.T) {
		t.Parallel()
		bad := commodity.ProductionAccount{Constant: -1, Variable: 3, Surplus: 3}
		if _, err := commodity.ApportionProduct(bad, 20); err == nil {
			t.Error("an account with negative constant capital was accepted")
		}
	})
}

// c+v+s of three int64 magnitudes can exceed MaxInt64; the apportioning must
// not wrap around into negative or nonsensical parts.
func TestApportionProduct_HugeAccountDoesNotOverflow(t *testing.T) {
	t.Parallel()
	huge := commodity.ProductionAccount{Constant: math.MaxInt64, Variable: math.MaxInt64, Surplus: math.MaxInt64}
	parts, err := commodity.ApportionProduct(huge, 21)
	if err != nil {
		t.Fatalf("ApportionProduct: %v", err)
	}
	for name, got := range map[string]commodity.Quantity{
		"constant": parts.Constant, "variable": parts.Variable, "surplus": parts.Surplus,
	} {
		if math.Abs(float64(got)-7) > 1e-9 {
			t.Errorf("%s part = %v, want 7 (a third of the product)", name, got)
		}
	}
}

// Multiplying a huge total by a value before dividing overflows to +Inf, which
// json cannot encode. The share is at most 1, so total x share never does.
func TestApportionProduct_HugeTotalStaysFinite(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		account commodity.ProductionAccount
		total   commodity.Quantity
	}{
		{"max float total", commodity.ProductionAccount{Constant: 24, Variable: 3, Surplus: 3}, commodity.Quantity(math.MaxFloat64)},
		{"large total, huge constant", commodity.ProductionAccount{Constant: math.MaxInt64, Variable: 1, Surplus: 1}, 1e300},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parts, err := commodity.ApportionProduct(tc.account, tc.total)
			if err != nil {
				t.Fatalf("ApportionProduct: %v", err)
			}
			for name, got := range map[string]commodity.Quantity{
				"constant": parts.Constant, "variable": parts.Variable, "surplus": parts.Surplus,
			} {
				if math.IsInf(float64(got), 0) || math.IsNaN(float64(got)) {
					t.Errorf("%s part = %v, want a finite number", name, got)
				}
			}
		})
	}
}
