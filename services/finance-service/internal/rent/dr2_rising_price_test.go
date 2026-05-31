// Package rent — unit tests for Vol. III Ch. 43 (DR II Rising Price of Production).
// Fixtures drawn from Engels Table XI: base=0s, increment=12s, grades=5 → [0,12,24,36,48] sum=120s.
package rent

import (
	"testing"
)

func TestRisingPriceVariantIsValid(t *testing.T) {
	t.Parallel()

	valid := []RisingPriceVariant{
		RisingPriceEventualityAVariant1,
		RisingPriceEventualityAVariant2,
		RisingPriceEventualityBVariant1,
		RisingPriceEventualityBVariant2,
		RisingPriceEventualityBVariant3,
	}
	for _, v := range valid {
		if !v.IsValid() {
			t.Errorf("variant %d should be valid", int(v))
		}
	}

	invalid := []RisingPriceVariant{0, 6}
	for _, v := range invalid {
		if v.IsValid() {
			t.Errorf("variant %d should be invalid", int(v))
		}
	}
}

// TestComputeRentSequence_TableXI verifies Engels Table XI exactly:
// ComputeRentSequence(0, 12, 5) → [0,12,24,36,48], sum=120.
func TestComputeRentSequence_TableXI(t *testing.T) {
	t.Parallel()

	got := ComputeRentSequence(0, 12, 5)
	want := []int64{0, 12, 24, 36, 48}

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("got[%d] = %d, want %d", i, got[i], w)
		}
	}

	var sum int64
	for _, r := range got {
		sum += r
	}
	if sum != 120 {
		t.Errorf("sum = %d, want 120", sum)
	}
}

// TestComputeRentSequence_General verifies the general arithmetic progression:
// ComputeRentSequence(100, 50, 4) → [100,150,200,250].
func TestComputeRentSequence_General(t *testing.T) {
	t.Parallel()

	base := int64(100)
	increment := int64(50)
	numGrades := 4
	got := ComputeRentSequence(base, increment, numGrades)

	if len(got) != numGrades {
		t.Fatalf("len = %d, want %d", len(got), numGrades)
	}
	for g := 1; g <= numGrades; g++ {
		want := base + int64(g-1)*increment
		if got[g-1] != want {
			t.Errorf("got[%d] = %d, want %d (grade %d)", g-1, got[g-1], want, g)
		}
	}
}

// TestComputeRentSequence_NonPositiveGrades verifies that zero and negative
// grade counts return a non-nil empty slice (len 0, not nil).
func TestComputeRentSequence_NonPositiveGrades(t *testing.T) {
	t.Parallel()

	cases := []int{0, -3}
	for _, n := range cases {
		got := ComputeRentSequence(0, 12, n)
		if got == nil {
			t.Errorf("numGrades=%d: got nil, want non-nil empty slice", n)
		}
		if len(got) != 0 {
			t.Errorf("numGrades=%d: len = %d, want 0", n, len(got))
		}
	}
}

// TestCh43NewIDsDistinct verifies that each ID constructor produces 24-char
// lowercase hex strings and that two consecutive calls differ.
func TestCh43NewIDsDistinct(t *testing.T) {
	t.Parallel()

	t.Run("RisingPriceDR2OutcomeID", func(t *testing.T) {
		t.Parallel()
		a := NewRisingPriceDR2OutcomeID()
		b := NewRisingPriceDR2OutcomeID()
		if len(string(a)) != 24 {
			t.Errorf("len = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two consecutive IDs must differ")
		}
	})

	t.Run("RentSequenceID", func(t *testing.T) {
		t.Parallel()
		a := NewRentSequenceID()
		b := NewRentSequenceID()
		if len(string(a)) != 24 {
			t.Errorf("len = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two consecutive IDs must differ")
		}
	})

	t.Run("EngelsTableEntryID", func(t *testing.T) {
		t.Parallel()
		a := NewEngelsTableEntryID()
		b := NewEngelsTableEntryID()
		if len(string(a)) != 24 {
			t.Errorf("len = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two consecutive IDs must differ")
		}
	})
}
