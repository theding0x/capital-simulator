package store

import (
	"strings"
	"testing"
)

// 00032 corrects the Ch.7 spinning seed to the 12-hour means (issue #450):
// 20 lbs cotton, spindle wear 480, advanced capital 3,240 min (27s).
func TestMigration00032_Ch7SeedFidelity(t *testing.T) {
	t.Parallel()
	raw, err := migrationsFS.ReadFile("migrations/00032_v1_ch07_seed_fidelity.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	up, down, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}
	for _, want := range []string{`"quantity":20`, `"wear_per_run":480`, "product_quantity = 20", "money_capital = 3240"} {
		if !strings.Contains(up, want) {
			t.Errorf("Up missing %q", want)
		}
	}
	for _, want := range []string{`"quantity":10`, `"wear_per_run":240`, "product_quantity = 10", "money_capital = 1800"} {
		if !strings.Contains(down, want) {
			t.Errorf("Down missing %q", want)
		}
	}
}
