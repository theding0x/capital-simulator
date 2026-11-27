package store

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
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

// The seeded 1847 Ten Hours' day must satisfy the limit the FactoryAct type
// derives for young persons (Capital Vol. I Ch. 10 §6). The stored row cannot
// carry this check itself: the working_days schema holds no constraint column,
// so nothing in SQL ties the seeded 600-minute day (6 h necessary + 4 h
// surplus) to the Act. This test reads the seed's own numbers and ties them to
// the limit derived from the type: the day is admitted, and one more minute of
// surplus is not.
func TestMigration00033_TenHoursSeedSatisfiesFactoryAct1847(t *testing.T) {
	t.Parallel()
	raw, err := migrationsFS.ReadFile("migrations/00033_v1_ch10_seed_fidelity.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	up, _, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}
	// Pin the 1847 row's id so the test cannot pick up another UPDATE.
	re := regexp.MustCompile(`necessary_labour_minutes\s*=\s*(\d+),\s*surplus_labour_minutes\s*=\s*(\d+)\s*WHERE id = '5eed00000000000000100004'`)
	m := re.FindStringSubmatch(up)
	if m == nil {
		t.Fatalf("Up has no necessary/surplus assignment for the 1847 row 5eed00000000000000100004:\n%s", up)
	}
	necessary, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		t.Fatalf("parse necessary_labour_minutes %q: %v", m[1], err)
	}
	surplus, err := strconv.ParseInt(m[2], 10, 64)
	if err != nil {
		t.Fatalf("parse surplus_labour_minutes %q: %v", m[2], err)
	}
	day := agent.WorkingDay{
		NecessaryLabourMinutes: agent.NecessaryLabourMinutes(necessary),
		SurplusLabourMinutes:   agent.SurplusLabourMinutes(surplus),
	}
	act1847 := agent.FactoryAct{Year: 1847, ChildLimitMinutes: 390, YoungPersonLimitMinutes: 600}
	c, err := act1847.Constraint(agent.AgeBandYoungPerson)
	if err != nil {
		t.Fatalf("derive 1847 young-person constraint: %v", err)
	}
	if err := agent.ValidateConstraint(day, c); err != nil {
		t.Errorf("seeded day (%d + %d min) vs %q limit %d = %v, want nil",
			necessary, surplus, c.Label, c.StatutoryLimitMinutes, err)
	}
	day.SurplusLabourMinutes++
	if err := agent.ValidateConstraint(day, c); !errors.Is(err, agent.ErrWorkingDayExceedsStatutoryLimit) {
		t.Errorf("seeded day plus one minute of surplus = %v, want ErrWorkingDayExceedsStatutoryLimit", err)
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
