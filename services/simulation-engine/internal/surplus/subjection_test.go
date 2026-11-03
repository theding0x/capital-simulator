package surplus_test

import (
	"reflect"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

// Capital Vol. I, Ch. 16, §1: "That labourer alone is productive, who
// produces surplus-value for the capitalist." Marx's schoolmaster example:
// the schoolmaster is productive if he works for a school-owning
// entrepreneur who profits from him, though not when paid out of revenue.
func TestCollectiveLabourer_ProductiveMembers_Schoolmaster(t *testing.T) {
	t.Parallel()
	collective := surplus.CollectiveLabourer{"schoolmaster", "factory-hand", "unrecorded"}
	records := []surplus.ProductiveLabour{
		{AgentID: "schoolmaster", ProducesSurplusValue: true},
		{AgentID: "factory-hand", ProducesSurplusValue: false},
	}
	got := collective.ProductiveMembers(records)
	want := []surplus.AgentID{"schoolmaster"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ProductiveMembers = %v, want %v", got, want)
	}
}

func TestCollectiveLabourer_ProductiveMembers_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	got := surplus.CollectiveLabourer(nil).ProductiveMembers(nil)
	if got == nil || len(got) != 0 {
		t.Fatalf("ProductiveMembers = %#v, want empty non-nil slice", got)
	}
}

// Formal subjection (control without technical change) is the regime of
// absolute surplus-value; real subjection (revolutionised technique) is the
// regime of relative surplus-value.
func TestSubjectionKind_SurplusOrigin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		kind surplus.SubjectionKind
		want surplus.Origin
	}{
		{surplus.SubjectionFormal, surplus.OriginAbsolute},
		{surplus.SubjectionReal, surplus.OriginRelative},
		{surplus.SubjectionKind("other"), ""},
	}
	for _, c := range cases {
		if got := c.kind.SurplusOrigin(); got != c.want {
			t.Errorf("%q.SurplusOrigin() = %q, want %q", c.kind, got, c.want)
		}
	}
}

func TestSubjectionKind_Valid(t *testing.T) {
	t.Parallel()
	if !surplus.SubjectionFormal.Valid() || !surplus.SubjectionReal.Valid() {
		t.Fatal("formal and real must be valid")
	}
	if surplus.SubjectionKind("").Valid() || surplus.SubjectionKind("other").Valid() {
		t.Fatal("unknown kinds must be invalid")
	}
}
