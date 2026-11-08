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

// 00033 corrects the 1847 Ten Hours' seed to 6h necessary + 4h surplus
// (issue #450): a 600-minute day, matching the Factory Act's limit.
func TestMigration00033_Ch10TenHoursSeed(t *testing.T) {
	t.Parallel()
	raw, err := migrationsFS.ReadFile("migrations/00033_v1_ch10_seed_fidelity.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	up, down, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}
	if !strings.Contains(up, "necessary_labour_minutes = 360, surplus_labour_minutes = 240") {
		t.Error("Up must set the 1847 day to 360 + 240 minutes")
	}
	if !strings.Contains(down, "necessary_labour_minutes = 300, surplus_labour_minutes = 300") {
		t.Error("Down must restore 300 + 300 minutes")
	}
}

// 00034 renames the Ch.13 twelve-man seed to Marx's twelve masons (issue #450)
// and must roll back to 00011's Arkwright spinning floor.
func TestMigration00034_Ch13TwelveMasonsSeed(t *testing.T) {
	t.Parallel()
	raw, err := migrationsFS.ReadFile("migrations/00034_v1_ch13_seed_fidelity.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	up, down, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}
	for _, want := range []string{"name = 'Twelve Masons'", "name = 'Mason 01'", "name = 'Mason 10'"} {
		if !strings.Contains(up, want) {
			t.Errorf("Up missing %q", want)
		}
	}
	for _, want := range []string{"name = 'Arkwright Spinning Floor'", "name = 'Spinner 01'", "name = 'Richard Arkwright'"} {
		if !strings.Contains(down, want) {
			t.Errorf("Down missing %q", want)
		}
	}
	if strings.Count(up, "UPDATE") != strings.Count(down, "UPDATE") {
		t.Errorf("Up has %d UPDATEs but Down has %d", strings.Count(up, "UPDATE"), strings.Count(down, "UPDATE"))
	}
}
