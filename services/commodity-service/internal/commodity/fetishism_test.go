package commodity

import (
	"strings"
	"testing"
)

func TestSocialRelationsOf(t *testing.T) {
	t.Parallel()
	linen := makeLinen()
	linen.ID = "linen-id"
	coat := makeCoat()
	coat.ID = "coat-id"
	iron := makeIron()

	rel := SocialRelationsOf(linen, []Commodity{linen, coat, iron})

	if rel.Subject.ID != "linen-id" {
		t.Errorf("Subject.ID = %v, want linen-id", rel.Subject.ID)
	}
	if rel.LabourPerUnit != 30 {
		t.Errorf("LabourPerUnit = %d, want 30", rel.LabourPerUnit)
	}
	if rel.ConcreteLabour.Kind != "weaving" {
		t.Errorf("ConcreteLabour.Kind = %v, want weaving", rel.ConcreteLabour.Kind)
	}
	if got := len(rel.LabourRelations); got != 2 {
		t.Fatalf("LabourRelations = %d, want 2", got)
	}
	// Every counterpart relation should equate to the subject's value (30m).
	for _, lr := range rel.LabourRelations {
		if lr.LabourTime != 30 {
			t.Errorf("LabourTime = %d, want 30 for counterpart %s",
				lr.LabourTime, lr.Counterpart.Name)
		}
	}
	if !strings.Contains(rel.Note, "Fetishism") {
		t.Errorf("Note should reference the fetishism critique; got %q", rel.Note)
	}
}

func TestSocialRelationsOf_SkipsSelfAndZeroValue(t *testing.T) {
	t.Parallel()
	linen := makeLinen()
	linen.ID = "linen-id"
	coat := makeCoat()
	coat.ID = "coat-id"
	free := makeIron()
	free.SNLTPerUnit = 0

	rel := SocialRelationsOf(linen, []Commodity{linen, coat, free})
	if len(rel.LabourRelations) != 1 {
		t.Errorf("expected only coat in relations, got %d entries",
			len(rel.LabourRelations))
	}
}
