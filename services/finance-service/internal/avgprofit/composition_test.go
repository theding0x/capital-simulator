package avgprofit

import "testing"

// --- OrganicComposition tests ---

// Sphere A composition: 600c + 100v → total 700.
// VariablePercent = 100*100/700 = 14 (integer).
// ConstantPercent = 600*100/700 = 85 (integer).
func TestOrganicComposition_SphereA(t *testing.T) {
	t.Parallel()
	oc := OrganicComposition{Constant: 600, Variable: 100}
	if got := oc.Total(); got != 700 {
		t.Errorf("Total() = %d, want 700", got)
	}
	if got := oc.VariablePercent(); got != 14 {
		t.Errorf("VariablePercent() = %d, want 14", got)
	}
	if got := oc.ConstantPercent(); got != 85 {
		t.Errorf("ConstantPercent() = %d, want 85", got)
	}
}

// Sphere B composition: 100c + 600v → total 700.
// VariablePercent = 600*100/700 = 85 (integer).
// ConstantPercent = 100*100/700 = 14 (integer).
func TestOrganicComposition_SphereB(t *testing.T) {
	t.Parallel()
	oc := OrganicComposition{Constant: 100, Variable: 600}
	if got := oc.VariablePercent(); got != 85 {
		t.Errorf("VariablePercent() = %d, want 85", got)
	}
	if got := oc.ConstantPercent(); got != 14 {
		t.Errorf("ConstantPercent() = %d, want 14", got)
	}
}

// European capital: 84c + 16v → total 100.
// VariablePercent = 16*100/100 = 16.
// ConstantPercent = 84*100/100 = 84.
func TestOrganicComposition_European(t *testing.T) {
	t.Parallel()
	oc := OrganicComposition{Constant: 84, Variable: 16}
	if got := oc.VariablePercent(); got != 16 {
		t.Errorf("VariablePercent() = %d, want 16", got)
	}
	if got := oc.ConstantPercent(); got != 84 {
		t.Errorf("ConstantPercent() = %d, want 84", got)
	}
}

// Asian capital: 16c + 84v → total 100.
func TestOrganicComposition_Asian(t *testing.T) {
	t.Parallel()
	oc := OrganicComposition{Constant: 16, Variable: 84}
	if got := oc.VariablePercent(); got != 84 {
		t.Errorf("VariablePercent() = %d, want 84", got)
	}
	if got := oc.ConstantPercent(); got != 16 {
		t.Errorf("ConstantPercent() = %d, want 16", got)
	}
}

// Invariant: VariablePercent == V*100/(c+v).
func TestOrganicComposition_VariablePercent_Invariant(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		c    ConstantCapital
		v    VariableCapital
	}{
		{"Sphere A", 600, 100},
		{"Sphere B", 100, 600},
		{"European", 84, 16},
		{"Asian", 16, 84},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			oc := OrganicComposition{Constant: tc.c, Variable: tc.v}
			total := int64(tc.c) + int64(tc.v)
			want := int64(tc.v) * 100 / total
			if got := oc.VariablePercent(); got != want {
				t.Errorf("VariablePercent() = %d, want %d", got, want)
			}
		})
	}
}

// Guard: zero total capital → VariablePercent and ConstantPercent both 0.
func TestOrganicComposition_ZeroTotal(t *testing.T) {
	t.Parallel()
	oc := OrganicComposition{Constant: 0, Variable: 0}
	if got := oc.Total(); got != 0 {
		t.Errorf("Total() = %d, want 0", got)
	}
	if got := oc.VariablePercent(); got != 0 {
		t.Errorf("VariablePercent() = %d, want 0 for zero total", got)
	}
	if got := oc.ConstantPercent(); got != 0 {
		t.Errorf("ConstantPercent() = %d, want 0 for zero total", got)
	}
}

// --- ClassifyComposition tests ---

// Below-average variable-percent → KindHigher (more constant capital).
// Sphere A var% = 14; if social average is 50 → A is KindHigher.
func TestClassifyComposition_KindHigher(t *testing.T) {
	t.Parallel()
	kind := ClassifyComposition(14, 50)
	if kind != KindHigher {
		t.Errorf("ClassifyComposition(14, 50) = %q, want %q", kind, KindHigher)
	}
}

// Above-average variable-percent → KindLower (less constant capital).
// Sphere B var% = 85; if social average is 50 → B is KindLower.
func TestClassifyComposition_KindLower(t *testing.T) {
	t.Parallel()
	kind := ClassifyComposition(85, 50)
	if kind != KindLower {
		t.Errorf("ClassifyComposition(85, 50) = %q, want %q", kind, KindLower)
	}
}

// Equal variable-percent → KindAverage.
func TestClassifyComposition_KindAverage(t *testing.T) {
	t.Parallel()
	kind := ClassifyComposition(50, 50)
	if kind != KindAverage {
		t.Errorf("ClassifyComposition(50, 50) = %q, want %q", kind, KindAverage)
	}
}

// European (var% 16) relative to Asian (var% 84): European is KindHigher.
func TestClassifyComposition_EuropeanVsAsian(t *testing.T) {
	t.Parallel()
	europeanVarPct := int64(16)
	asianVarPct := int64(84)
	// European has lower var% than Asian social average → KindHigher.
	if got := ClassifyComposition(europeanVarPct, asianVarPct); got != KindHigher {
		t.Errorf("European vs Asian: got %q, want KindHigher", got)
	}
	// Asian has higher var% than European average → KindLower.
	if got := ClassifyComposition(asianVarPct, europeanVarPct); got != KindLower {
		t.Errorf("Asian vs European: got %q, want KindLower", got)
	}
}

// --- TechnicalComposition tests ---

// Spinning mill example: 200 Mp, 4 workers → 50 Mp per worker.
func TestTechnicalComposition_MpPerWorker(t *testing.T) {
	t.Parallel()
	tc := TechnicalComposition{MeansOfProduction: 200, Workers: 4}
	if got := tc.MpPerWorker(); got != 50 {
		t.Errorf("MpPerWorker() = %d, want 50", got)
	}
}

// Guard: zero workers → 0 (no division by zero).
func TestTechnicalComposition_ZeroWorkers(t *testing.T) {
	t.Parallel()
	tc := TechnicalComposition{MeansOfProduction: 100, Workers: 0}
	if got := tc.MpPerWorker(); got != 0 {
		t.Errorf("MpPerWorker() = %d, want 0 when Workers=0", got)
	}
}

// Negative workers guard: treat as zero (Workers <= 0).
func TestTechnicalComposition_NegativeWorkers(t *testing.T) {
	t.Parallel()
	tc := TechnicalComposition{MeansOfProduction: 100, Workers: -5}
	if got := tc.MpPerWorker(); got != 0 {
		t.Errorf("MpPerWorker() = %d, want 0 when Workers<0", got)
	}
}

// Exact integer division: 300 Mp, 3 workers → 100 Mp per worker.
func TestTechnicalComposition_ExactDivision(t *testing.T) {
	t.Parallel()
	tc := TechnicalComposition{MeansOfProduction: 300, Workers: 3}
	if got := tc.MpPerWorker(); got != 100 {
		t.Errorf("MpPerWorker() = %d, want 100", got)
	}
}
