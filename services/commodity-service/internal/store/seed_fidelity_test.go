package store

import (
	"strings"
	"testing"
)

// 00004 corrects the Ch.9 seeds (issue #450): £410 c = 984,000 min and the
// 1871 mill's c = £378 = 907,200 min. Both must roll back to 00003's values.
func TestMigration00004_Ch9SeedFidelity(t *testing.T) {
	t.Parallel()
	raw, err := migrationsFS.ReadFile("migrations/00004_v1_ch09_seed_fidelity.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	up, down, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}
	for _, want := range []string{"constant = 984000", "constant = 907200"} {
		if !strings.Contains(up, want) {
			t.Errorf("Up missing %q", want)
		}
	}
	for _, want := range []string{"constant = 9840000", "constant = 499200"} {
		if !strings.Contains(down, want) {
			t.Errorf("Down missing %q", want)
		}
	}
}
