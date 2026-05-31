// Package rent — unit tests for Vol. III Ch. 40 (Differential Rent II).
// Fixtures drawn from Marx, Vol. III Ch. 40: Soil D doubled-investment.
// Capital £2.50 (250p), output 4 qtr (400), surplus-profit £9 (900p) per tranche.
package rent

import (
	"strings"
	"testing"
)

// ch40Tranche returns a SuccessiveInvestment with the Marx Ch. 40 Soil-D values.
func ch40Tranche(parcelID string, tranche int) SuccessiveInvestment {
	return SuccessiveInvestment{
		ParcelID:        parcelID,
		TrancheNumber:   tranche,
		CapitalBP:       250,
		OutputQuarters:  400,
		SurplusProfitBP: 900,
	}
}

// ── ComputeDifferentialRentII ─────────────────────────────────────────────────

func TestComputeDifferentialRentII_DoubledInvestment(t *testing.T) {
	t.Parallel()
	invs := []SuccessiveInvestment{
		ch40Tranche("soil-D", 1),
		ch40Tranche("soil-D", 2),
	}
	tbl := ComputeDifferentialRentII("soil-D", invs, 1)

	if tbl.ParcelID != "soil-D" {
		t.Errorf("ParcelID = %q, want %q", tbl.ParcelID, "soil-D")
	}
	if tbl.TotalCapitalBP != 500 {
		t.Errorf("TotalCapitalBP = %d, want 500", tbl.TotalCapitalBP)
	}
	if tbl.TotalOutputQuarters != 800 {
		t.Errorf("TotalOutputQuarters = %d, want 800", tbl.TotalOutputQuarters)
	}
	if tbl.TotalRentBP != 1800 {
		t.Errorf("TotalRentBP = %d, want 1800", tbl.TotalRentBP)
	}
	if tbl.RentPerAcreBP != 1800 {
		t.Errorf("RentPerAcreBP = %d, want 1800 (1 acre)", tbl.RentPerAcreBP)
	}
}

func TestComputeDifferentialRentII_AcresDivision(t *testing.T) {
	t.Parallel()
	invs := []SuccessiveInvestment{
		ch40Tranche("soil-D2", 1),
		ch40Tranche("soil-D2", 2),
	}
	tbl := ComputeDifferentialRentII("soil-D2", invs, 2)

	if tbl.RentPerAcreBP != 900 {
		t.Errorf("RentPerAcreBP = %d, want 900 (1800 / 2 acres)", tbl.RentPerAcreBP)
	}
	if tbl.TotalRentBP != 1800 {
		t.Errorf("TotalRentBP = %d, want 1800 (unchanged by acres divisor)", tbl.TotalRentBP)
	}
}

func TestComputeDifferentialRentII_AcresGuard(t *testing.T) {
	t.Parallel()
	invs := []SuccessiveInvestment{
		ch40Tranche("soil-D3", 1),
		ch40Tranche("soil-D3", 2),
	}

	// acres=0 must be treated as 1.
	tbl0 := ComputeDifferentialRentII("soil-D3", invs, 0)
	if tbl0.RentPerAcreBP != tbl0.TotalRentBP {
		t.Errorf("acres=0: RentPerAcreBP=%d, TotalRentBP=%d; want equal (guard to 1)",
			tbl0.RentPerAcreBP, tbl0.TotalRentBP)
	}

	// acres=-3 must be treated as 1.
	tblNeg := ComputeDifferentialRentII("soil-D3", invs, -3)
	if tblNeg.RentPerAcreBP != tblNeg.TotalRentBP {
		t.Errorf("acres=-3: RentPerAcreBP=%d, TotalRentBP=%d; want equal (guard to 1)",
			tblNeg.RentPerAcreBP, tblNeg.TotalRentBP)
	}
}

func TestComputeDifferentialRentII_Empty(t *testing.T) {
	t.Parallel()

	// nil investments
	tblNil := ComputeDifferentialRentII("empty-parcel", nil, 1)
	if tblNil.TotalCapitalBP != 0 {
		t.Errorf("nil: TotalCapitalBP = %d, want 0", tblNil.TotalCapitalBP)
	}
	if tblNil.TotalOutputQuarters != 0 {
		t.Errorf("nil: TotalOutputQuarters = %d, want 0", tblNil.TotalOutputQuarters)
	}
	if tblNil.TotalRentBP != 0 {
		t.Errorf("nil: TotalRentBP = %d, want 0", tblNil.TotalRentBP)
	}
	if tblNil.RentPerAcreBP != 0 {
		t.Errorf("nil: RentPerAcreBP = %d, want 0", tblNil.RentPerAcreBP)
	}
	if tblNil.ParcelID != "empty-parcel" {
		t.Errorf("nil: ParcelID = %q, want %q", tblNil.ParcelID, "empty-parcel")
	}

	// empty (non-nil) investments
	tblEmpty := ComputeDifferentialRentII("empty-parcel", []SuccessiveInvestment{}, 1)
	if tblEmpty.TotalRentBP != 0 {
		t.Errorf("empty: TotalRentBP = %d, want 0", tblEmpty.TotalRentBP)
	}
	if tblEmpty.RentPerAcreBP != 0 {
		t.Errorf("empty: RentPerAcreBP = %d, want 0", tblEmpty.RentPerAcreBP)
	}
}

// ── ID constructors ───────────────────────────────────────────────────────────

func TestCh40NewIDsDistinct(t *testing.T) {
	t.Parallel()

	// Each ID must be 24 hex characters (96 bits = 12 bytes = 24 hex chars).
	const wantLen = 24

	siID := NewSuccessiveInvestmentID()
	dr2ID := NewDifferentialRentIITableID()
	ltID := NewLeaseTermID()

	ids := []struct {
		name string
		val  string
	}{
		{"SuccessiveInvestmentID", string(siID)},
		{"DifferentialRentIITableID", string(dr2ID)},
		{"LeaseTermID", string(ltID)},
	}

	for _, tc := range ids {
		if len(tc.val) != wantLen {
			t.Errorf("%s: len=%d, want %d", tc.name, len(tc.val), wantLen)
		}
		// Must be lowercase hex.
		if tc.val != strings.ToLower(tc.val) {
			t.Errorf("%s: %q is not lowercase hex", tc.name, tc.val)
		}
	}

	// All three must be distinct.
	if string(siID) == string(dr2ID) {
		t.Errorf("SuccessiveInvestmentID == DifferentialRentIITableID: %q", siID)
	}
	if string(siID) == string(ltID) {
		t.Errorf("SuccessiveInvestmentID == LeaseTermID: %q", siID)
	}
	if string(dr2ID) == string(ltID) {
		t.Errorf("DifferentialRentIITableID == LeaseTermID: %q", dr2ID)
	}

	// Two successive calls to the same constructor must also be distinct.
	siID2 := NewSuccessiveInvestmentID()
	if siID == siID2 {
		t.Errorf("NewSuccessiveInvestmentID returned duplicate: %q", siID)
	}
}
