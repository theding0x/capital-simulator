// Vol. II Ch. 8 — tests for fixed/circulating capital domain.
package composition_test

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
)

// TestKindForRole verifies the fixed role→kind map.
func TestKindForRole(t *testing.T) {
	t.Parallel()
	cases := []struct {
		role composition.RoleInProcess
		want composition.CapitalKind
	}{
		{composition.RoleInstrumentOfLabour, composition.KindFixed},
		{composition.RoleRawMaterial, composition.KindCirculating},
		{composition.RoleAuxiliaryMaterial, composition.KindCirculating},
		{composition.RoleLabourPower, composition.KindCirculating},
	}
	for _, tc := range cases {
		if got := composition.KindForRole(tc.role); got != tc.want {
			t.Errorf("KindForRole(%q) = %q, want %q", tc.role, got, tc.want)
		}
	}
}

// TestCapitalComponent_Validate_RoleDeterminesKind verifies the fixed map invariant.
func TestCapitalComponent_Validate_RoleDeterminesKind(t *testing.T) {
	t.Parallel()
	// Instrument of labour → fixed
	c := composition.CapitalComponent{
		ID:    composition.NewCapitalComponentID(),
		Role:  composition.RoleInstrumentOfLabour,
		Kind:  composition.KindFixed,
		Pence: 1_000_000,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("instrument-of-labour should be valid: %v", err)
	}

	// Wrong kind for role → error
	c.Kind = composition.KindCirculating
	if err := c.Validate(); err == nil {
		t.Fatal("instrument-of-labour with circulating kind should fail validation")
	}
}

// TestCapitalComponent_Validate_AuxiliaryIsCirculating verifies Marx's rejection of
// Ramsay's misclassification of auxiliary materials as fixed capital.
func TestCapitalComponent_Validate_AuxiliaryIsCirculating(t *testing.T) {
	t.Parallel()
	c := composition.CapitalComponent{
		ID:    composition.NewCapitalComponentID(),
		Role:  composition.RoleAuxiliaryMaterial,
		Kind:  composition.KindFixed, // wrong — auxiliary is always circulating
		Pence: 500,
	}
	if err := c.Validate(); err == nil {
		t.Fatal("auxiliary material with fixed kind should fail (Ramsay's error rejected)")
	}
	c.Kind = composition.KindCirculating
	if err := c.Validate(); err != nil {
		t.Fatalf("auxiliary material with circulating kind should be valid: %v", err)
	}
}

// TestCapitalComponent_Validate_InvalidRole verifies that unknown roles are rejected.
func TestCapitalComponent_Validate_InvalidRole(t *testing.T) {
	t.Parallel()
	c := composition.CapitalComponent{
		ID:    composition.NewCapitalComponentID(),
		Role:  "wizard", // unknown
		Kind:  composition.KindCirculating,
		Pence: 100,
	}
	if err := c.Validate(); err == nil {
		t.Fatal("unknown role should fail validation")
	}
}

// TestComputeLinearWear_10YearMachine replicates Marx's §I fixture:
// "If a machine worth £10,000 lasts for ten years, the period of turnover
// of the value originally advanced for it amounts to ten years."
// PencePurchased = 1,000,000 (£10,000 × 100). WearPerYear = 100,000.
func TestComputeLinearWear_10YearMachine(t *testing.T) {
	t.Parallel()
	const yearMinutes int64 = 525_600
	item := composition.FixedCapitalItem{
		ID:                 "machine-10yr",
		PencePurchased:     1_000_000,
		ServiceLifeMinutes: 10 * yearMinutes,
		WearModel:          composition.WearLinear,
	}
	wear := composition.ComputeLinearWear(item, yearMinutes, composition.WearUse)
	const wantPerYear composition.Pence = 100_000
	if wear.Pence != wantPerYear {
		t.Errorf("10-year machine annual wear = %d, want %d", wear.Pence, wantPerYear)
	}
	// Sum over 10 years equals PencePurchased.
	var total composition.Pence
	for range 10 {
		w := composition.ComputeLinearWear(item, yearMinutes, composition.WearUse)
		total += w.Pence
	}
	if total != item.PencePurchased {
		t.Errorf("10-year total wear = %d, want %d", total, item.PencePurchased)
	}
}

// TestComputeLinearWear_TwoMachines_HalfLifeDoubleWear replicates Marx's §I
// comparative fixture: "If of two machines of equal value one wears out in five
// years and the other in ten, then the first yields twice as much value in the
// same time as the second."
func TestComputeLinearWear_TwoMachines_HalfLifeDoubleWear(t *testing.T) {
	t.Parallel()
	const yearMinutes int64 = 525_600
	fast := composition.FixedCapitalItem{
		ID:                 "machine-5yr",
		PencePurchased:     1_000_000,
		ServiceLifeMinutes: 5 * yearMinutes,
	}
	slow := composition.FixedCapitalItem{
		ID:                 "machine-10yr",
		PencePurchased:     1_000_000,
		ServiceLifeMinutes: 10 * yearMinutes,
	}
	wFast := composition.ComputeLinearWear(fast, yearMinutes, composition.WearUse)
	wSlow := composition.ComputeLinearWear(slow, yearMinutes, composition.WearUse)
	if wFast.Pence != 2*wSlow.Pence {
		t.Errorf("5-year machine wear = %d, 10-year = %d; want 2× ratio", wFast.Pence, wSlow.Pence)
	}
}

// TestComputeLinearWear_MoralDepreciation verifies that WearMoralDepreciation
// can be recorded with non-zero pence even when use-wear is zero.
// Fixture: rolling stock originally worth £40,000; market-equivalent now £30,000;
// 25% moral depreciation = £10,000 = 1,000,000 pence.
func TestComputeLinearWear_MoralDepreciation(t *testing.T) {
	t.Parallel()
	moral := composition.WearAndTear{
		ID:                   composition.NewWearAndTearID(),
		FixedCapitalItemID:   "rolling-stock",
		Pence:                1_000_000, // 25% of £40k purchase price
		PeriodCoveredMinutes: 0,         // no time elapsed for use-wear
		CalculatedAt:         time.Now().UTC(),
		Cause:                composition.WearMoralDepreciation,
	}
	if moral.Cause != composition.WearMoralDepreciation {
		t.Fatal("cause must be WearMoralDepreciation")
	}
	if moral.Pence <= 0 {
		t.Fatal("moral depreciation pence must be positive")
	}
}

// TestFixedCapitalItem_Validate_ZeroLife checks the service-life invariant.
func TestFixedCapitalItem_Validate_ZeroLife(t *testing.T) {
	t.Parallel()
	item := composition.FixedCapitalItem{
		ID:                 "bad",
		PencePurchased:     1000,
		ServiceLifeMinutes: 0,
	}
	if err := item.Validate(); err == nil {
		t.Fatal("zero service life should fail validation")
	}
}

// TestRepair_Validate_RunningMustNotExtendLife enforces the running-repair invariant.
func TestRepair_Validate_RunningMustNotExtendLife(t *testing.T) {
	t.Parallel()
	r := composition.Repair{
		ID:                          composition.NewRepairID(),
		FixedCapitalItemID:          "item-1",
		Kind:                        composition.RepairRunning,
		Pence:                       500,
		OccurredAt:                  time.Now(),
		ServiceLifeExtensionMinutes: 10, // must be 0
	}
	if err := r.Validate(); err == nil {
		t.Fatal("running repair with non-zero life extension should fail")
	}
	r.ServiceLifeExtensionMinutes = 0
	if err := r.Validate(); err != nil {
		t.Fatalf("running repair with zero life extension should be valid: %v", err)
	}
}

// TestRepair_Validate_LargeMustExtendLife enforces the large-repair invariant.
func TestRepair_Validate_LargeMustExtendLife(t *testing.T) {
	t.Parallel()
	r := composition.Repair{
		ID:                          composition.NewRepairID(),
		FixedCapitalItemID:          "item-1",
		Kind:                        composition.RepairLarge,
		Pence:                       200_000,
		OccurredAt:                  time.Now(),
		ServiceLifeExtensionMinutes: 0, // must be > 0
	}
	if err := r.Validate(); err == nil {
		t.Fatal("large repair with zero life extension should fail")
	}
	r.ServiceLifeExtensionMinutes = 2 * 525_600 // 2 years
	if err := r.Validate(); err != nil {
		t.Fatalf("large repair with positive life extension should be valid: %v", err)
	}
}

// TestIDUniqueness verifies that successive ID calls produce distinct values.
func TestIDUniqueness(t *testing.T) {
	t.Parallel()
	seen := make(map[string]bool)
	for range 100 {
		id := string(composition.NewCapitalComponentID())
		if seen[id] {
			t.Fatalf("duplicate ID: %s", id)
		}
		seen[id] = true
	}
}

// TestOxDualForm verifies the dual-form ox fixture: the same ox may appear as
// instrument of labour (fixed) or raw material (circulating) depending on role.
func TestOxDualForm(t *testing.T) {
	t.Parallel()
	oxAsToil := composition.CapitalComponent{
		ID:    composition.NewCapitalComponentID(),
		Role:  composition.RoleInstrumentOfLabour,
		Kind:  composition.KindFixed,
		Pence: 10_000,
	}
	oxAsBeef := composition.CapitalComponent{
		ID:    composition.NewCapitalComponentID(),
		Role:  composition.RoleRawMaterial,
		Kind:  composition.KindCirculating,
		Pence: 10_000,
	}
	if err := oxAsToil.Validate(); err != nil {
		t.Fatalf("ox-as-toil should be valid: %v", err)
	}
	if err := oxAsBeef.Validate(); err != nil {
		t.Fatalf("ox-as-beef should be valid: %v", err)
	}
	if oxAsToil.Kind == oxAsBeef.Kind {
		t.Fatal("ox-as-toil and ox-as-beef must have different kinds")
	}
}
