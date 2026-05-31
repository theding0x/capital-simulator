package rent

import (
	"errors"
	"testing"
)

// Fixtures drawn from Marx, Vol. III Ch. 37 (Moore-Aveling numbering):
// The canonical three-class picture: Duke of Sutherland as the Landowner,
// a Scottish tenant-farmer as the AgriculturalCapitalist, and the best
// Highland parcel (fertility grade A=1, regulating) vs. the better grades.
//
// Capital advanced = £10 money-capital ≈ 4800 LabourMinutes.
// Ground-rent paid by the farmer to the Duke = 2400 LabourMinutes (half the
// advance — a high differential rent on better soils, per Marx's schematic).

// ── RentForm.IsValid ──────────────────────────────────────────────────────────

func TestRentFormIsValid(t *testing.T) {
	t.Parallel()

	valid := []RentForm{RentFormDifferential, RentFormAbsolute, RentFormMonopoly}
	for _, f := range valid {
		if !f.IsValid() {
			t.Errorf("RentForm(%d).IsValid() = false, want true", f)
		}
	}

	invalid := []RentForm{RentForm(0), RentForm(4)}
	for _, f := range invalid {
		if f.IsValid() {
			t.Errorf("RentForm(%d).IsValid() = true, want false", f)
		}
	}
}

// ── GroundRent.Validate ───────────────────────────────────────────────────────

func TestGroundRentValidate(t *testing.T) {
	t.Parallel()

	ownerID := "duke-of-sutherland-0001"
	capitalistID := "highland-tenant-farmer-01"

	t.Run("happy", func(t *testing.T) {
		t.Parallel()
		g := GroundRent{
			Form:                RentFormDifferential,
			AmountLabourMinutes: 2400,
			LandOwnerID:         ownerID,
			CapitalistID:        capitalistID,
		}
		if err := g.Validate(); err != nil {
			t.Errorf("Validate() = %v, want nil", err)
		}
	})

	t.Run("negative_rent", func(t *testing.T) {
		t.Parallel()
		g := GroundRent{
			Form:                RentFormAbsolute,
			AmountLabourMinutes: -1,
			LandOwnerID:         ownerID,
			CapitalistID:        capitalistID,
		}
		if err := g.Validate(); !errors.Is(err, ErrNegativeRent) {
			t.Errorf("Validate() = %v, want ErrNegativeRent", err)
		}
	})

	t.Run("same_owner_operator", func(t *testing.T) {
		t.Parallel()
		g := GroundRent{
			Form:                RentFormDifferential,
			AmountLabourMinutes: 2400,
			LandOwnerID:         "owner-farmer-same-0001",
			CapitalistID:        "owner-farmer-same-0001",
		}
		if err := g.Validate(); !errors.Is(err, ErrSameOwnerOperator) {
			t.Errorf("Validate() = %v, want ErrSameOwnerOperator", err)
		}
	})

	t.Run("invalid_form_zero", func(t *testing.T) {
		t.Parallel()
		g := GroundRent{
			Form:                RentForm(0),
			AmountLabourMinutes: 2400,
			LandOwnerID:         ownerID,
			CapitalistID:        capitalistID,
		}
		if err := g.Validate(); !errors.Is(err, ErrInvalidRentForm) {
			t.Errorf("Validate() = %v, want ErrInvalidRentForm", err)
		}
	})

	t.Run("invalid_form_99", func(t *testing.T) {
		t.Parallel()
		g := GroundRent{
			Form:                RentForm(99),
			AmountLabourMinutes: 2400,
			LandOwnerID:         ownerID,
			CapitalistID:        capitalistID,
		}
		if err := g.Validate(); !errors.Is(err, ErrInvalidRentForm) {
			t.Errorf("Validate() = %v, want ErrInvalidRentForm", err)
		}
	})

	t.Run("boundary_zero_amount", func(t *testing.T) {
		// The regulating (worst) parcel pays zero differential rent — it sets the
		// floor from which better parcels earn surplus. Zero is legal.
		t.Parallel()
		g := GroundRent{
			Form:                RentFormDifferential,
			AmountLabourMinutes: 0,
			LandOwnerID:         ownerID,
			CapitalistID:        capitalistID,
		}
		if err := g.Validate(); err != nil {
			t.Errorf("Validate() = %v, want nil for zero-rent regulating parcel", err)
		}
	})
}

// ── New*ID constructors ───────────────────────────────────────────────────────

func TestNewIDsDistinct(t *testing.T) {
	t.Parallel()

	t.Run("LandParcelID", func(t *testing.T) {
		t.Parallel()
		a, b := NewLandParcelID(), NewLandParcelID()
		if len(string(a)) != 24 {
			t.Errorf("NewLandParcelID len = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two successive NewLandParcelID() calls must differ")
		}
	})

	t.Run("LandownerID", func(t *testing.T) {
		t.Parallel()
		a, b := NewLandownerID(), NewLandownerID()
		if len(string(a)) != 24 {
			t.Errorf("NewLandownerID len = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two successive NewLandownerID() calls must differ")
		}
	})

	t.Run("AgriculturalCapitalistID", func(t *testing.T) {
		t.Parallel()
		a, b := NewAgriculturalCapitalistID(), NewAgriculturalCapitalistID()
		if len(string(a)) != 24 {
			t.Errorf("NewAgriculturalCapitalistID len = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two successive NewAgriculturalCapitalistID() calls must differ")
		}
	})

	t.Run("GroundRentID", func(t *testing.T) {
		t.Parallel()
		a, b := NewGroundRentID(), NewGroundRentID()
		if len(string(a)) != 24 {
			t.Errorf("NewGroundRentID len = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two successive NewGroundRentID() calls must differ")
		}
	})
}
