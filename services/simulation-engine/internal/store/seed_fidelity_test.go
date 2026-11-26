package store

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

const (
	ch15SeedPrefix = "5eed00000000000000150"
	ch15MillID     = "5eed00000000000000150f02"

	// ch15MinutesPerPenny is the plough-passage calibration (66 men x 60 min = 360d),
	// asserted in checkCh15Calibration and reused by the plough fixture.
	ch15MinutesPerPenny = 11
)

// ch15SetID captures every machine id a statement writes: machines.id and
// factory_machines.machine_id.
var ch15SetID = regexp.MustCompile("SET (?:id|machine_id) = '([^']+)'")

// ch15Old is the 00002 machine id (not hex); ch15New is its hex-valid successor.
func ch15Old(n int) string { return fmt.Sprintf("%sm%02d", ch15SeedPrefix, n) }
func ch15New(n int) string { return fmt.Sprintf("%sa%02d", ch15SeedPrefix, n) }

// ch15SQLOnly drops comment lines so that assertions see executable SQL only.
func ch15SQLOnly(s string) string {
	var kept []string
	for _, line := range strings.Split(s, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// 00072 restates the Vol. I Ch. 15 seed in Marx's figures (issue #450): the
// steam-plough at his three-pence an hour (its price and life are assumed), the
// old spinning machine relabelled as Baynes's 450-spindle mule, and the machine
// ids made hex-valid. Every number below is recomputed from the domain, not typed
// twice.
func TestMigration00072_V1Ch15SeedFidelity(t *testing.T) {
	t.Parallel()
	raw, err := migrationsFS.ReadFile("migrations/00072_v1_ch15_seed_fidelity.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	upRaw, downRaw, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}
	up, down := ch15SQLOnly(upRaw), ch15SQLOnly(downRaw)

	checkCh15Calibration(t)
	checkCh15Plough(t, up)
	checkCh15Mule(t, up)
	checkCh15NeedlesRenameOnly(t, up)
	checkCh15Ids(t, up, down)
	checkCh15IdPairing(t, up, down)
	checkCh15DownRestores(t, down)
	checkCh15Hygiene(t, up, down)
}

// Marx §2: 66 men in one hour "produce ... a value of 30 shillings"; their 15s
// wages are half an hour each. 3,960 min over 360d is exact: 11 min a penny.
func checkCh15Calibration(t *testing.T) {
	t.Helper()
	men, minutesEach := int64(66), int64(60)
	hand := men * minutesEach
	hourPence := int64(30 * 12)
	if hand%hourPence != 0 {
		t.Fatalf("%d min is not a whole number of pence at %dd", hand, hourPence)
	}
	perPenny := hand / hourPence
	if perPenny != ch15MinutesPerPenny {
		t.Fatalf("minutes per penny = %d, want %d", perPenny, ch15MinutesPerPenny)
	}
	wagePence := int64(15 * 12)
	if got, want := wagePence*perPenny, men*30; got != want {
		t.Fatalf("15s = %d min, want 66 men x 30 min = %d", got, want)
	}
}

// ch15Plough is the seeded steam-plough. Its price and life are assumed: Marx gives
// neither, and the 5,475-day life is borrowed from his §3B illustration of "a
// machine" worked 8 hours daily for 15 years.
func ch15Plough() machinery.Machine {
	const (
		men, minutesEach = int64(66), int64(60)
		costPence        = int64(3)
		minutesPerPenny  = int64(ch15MinutesPerPenny)
		fifteenYears     = machinery.LifespanDays(15 * 365)
	)
	hours := machinery.StandardWorkingDayHours
	return machinery.Machine{
		MachineValue:      machinery.MachineValue(costPence * minutesPerPenny * hours * int64(fifteenYears)),
		LifespanDays:      fifteenYears,
		ProductivePower:   machinery.ProductivePower(hours),
		HandLabourPerUnit: machinery.LabourMinutes(men * minutesEach),
	}
}

// m05: three-pence an hour (33 min) x 8 h x 5,475 days = 1,445,400.
func checkCh15Plough(t *testing.T, up string) {
	t.Helper()
	m := ch15Plough()
	hours := int64(m.ProductivePower)
	if got := machinery.DailyWearAndTear(m); got != 264 {
		t.Errorf("DailyWearAndTear = %d, want 264", got)
	}
	if got := machinery.ValueTransferredPerUnit(m, hours); got != 33 {
		t.Errorf("ValueTransferredPerUnit = %d, want 33 (3d)", got)
	}
	want := fmt.Sprintf("UPDATE machines SET id = '%s', machine_value = %d, lifespan_days = %d, productive_power = %d WHERE id = '%s';",
		ch15New(5), m.MachineValue, m.LifespanDays, m.ProductivePower, ch15Old(5))
	if !strings.Contains(up, want) {
		t.Errorf("Up missing %q", want)
	}
}

func ch15MuleMill() machinery.Factory {
	handPerSpindle := machinery.LabourMinutes(60 * 60) // "sixty hours" for 13 oz
	return machinery.Factory{
		Name:       "Baynes Mule Mill",
		PrimeMover: machinery.PrimeMover{Kind: machinery.PrimeMoverSteam, Horsepower: 1},
		Machines: []machinery.Machine{{
			MachineValue:      480_000,
			LifespanDays:      2000,
			ProductivePower:   450,
			HandLabourPerUnit: handPerSpindle,
		}},
		Attendants:                  5,
		LivingLabourPerAttendant:    300,
		NecessaryLabourPerAttendant: 150,
	}
}

// m06 / f02: 450 spindles x 3,600 min = 27,000 h of hand-spinning, one
// horse-power, 2.5 operatives for a ten-hour day = 1,500 min.
func checkCh15Mule(t *testing.T, up string) {
	t.Helper()
	mill := ch15MuleMill()
	mule := mill.Machines[0]
	r := mill.RunTick()
	if r.UnitsProduced != 450 {
		t.Errorf("UnitsProduced = %d, want 450", r.UnitsProduced)
	}
	if r.HandLabourSaved != machinery.LabourMinutes(27000*60) {
		t.Errorf("HandLabourSaved = %d, want 27,000 h = %d", r.HandLabourSaved, 27000*60)
	}
	if r.LivingLabour != 1500 {
		t.Errorf("LivingLabour = %d, want 1500", r.LivingLabour)
	}
	if got := machinery.LabourSavingRatio(1_620_000, 150*60); got != 180 {
		t.Errorf("LabourSavingRatio = %d, want 180", got)
	}
	if err := mill.Validate(); err != nil {
		t.Errorf("Validate = %v, want nil", err)
	}
	for _, want := range []string{
		fmt.Sprintf("UPDATE machines SET id = '%s', name = 'Self-acting Mule (450 spindles)', working_tool = 'spindle', productive_power = %d, hand_labour_per_unit = %d WHERE id = '%s';",
			ch15New(6), mule.ProductivePower, mule.HandLabourPerUnit, ch15Old(6)),
		fmt.Sprintf("UPDATE factories SET name = '%s', prime_mover_horsepower = %d, attendants = %d, living_labour_per_attendant = %d, necessary_labour_per_attendant = %d WHERE id = '%s';",
			mill.Name, mill.PrimeMover.Horsepower, mill.Attendants, mill.LivingLabourPerAttendant, mill.NecessaryLabourPerAttendant, ch15MillID),
	} {
		if !strings.Contains(up, want) {
			t.Errorf("Up missing %q", want)
		}
	}
}

// The needle-machines are Marx's 145,000 a day; only their ids change.
func checkCh15NeedlesRenameOnly(t *testing.T, up string) {
	t.Helper()
	for n := 1; n <= 4; n++ {
		want := fmt.Sprintf("UPDATE machines SET id = '%s' WHERE id = '%s';", ch15New(n), ch15Old(n))
		if !strings.Contains(up, want) {
			t.Errorf("Up missing rename-only %q", want)
		}
	}
}

func ch15IDCounts(t *testing.T, section string, wantLen int) map[string]int {
	t.Helper()
	matches := ch15SetID.FindAllStringSubmatch(section, -1)
	if len(matches) != wantLen {
		t.Fatalf("found %d machine id assignments, want %d", len(matches), wantLen)
	}
	counts := map[string]int{}
	for _, m := range matches {
		counts[m[1]]++
	}
	return counts
}

// Six machines are renamed in machines and re-pointed in factory_machines: each
// id appears twice per direction, and every new id is 96-bit hex.
func checkCh15Ids(t *testing.T, up, down string) {
	t.Helper()
	upCounts := ch15IDCounts(t, up, 12)
	for id := range upCounts {
		if len(id) != 24 || !strings.HasPrefix(id, ch15SeedPrefix) {
			t.Errorf("Up id %q: want 24 chars with prefix %s", id, ch15SeedPrefix)
		}
		if _, err := hex.DecodeString(id); err != nil {
			t.Errorf("Up id %q is not hex: %v", id, err)
		}
	}
	downCounts := ch15IDCounts(t, down, 12)
	for n := 1; n <= 6; n++ {
		if got := upCounts[ch15New(n)]; got != 2 {
			t.Errorf("Up: %s appears %d times, want 2", ch15New(n), got)
		}
		if got := downCounts[ch15Old(n)]; got != 2 {
			t.Errorf("Down: %s appears %d times, want 2", ch15Old(n), got)
		}
	}
	if len(upCounts) != 6 || len(downCounts) != 6 {
		t.Errorf("distinct ids: Up %d, Down %d, want 6 each", len(upCounts), len(downCounts))
	}
}

// The counts above would still pass if two WHERE clauses were swapped, so assert
// every old<->new pairing as a full statement: factory_machines in both directions
// for all six machines, and the rename-only machines rows back in Down (m05 and m06
// are checked in full by checkCh15Plough, checkCh15Mule and checkCh15DownRestores).
func checkCh15IdPairing(t *testing.T, up, down string) {
	t.Helper()
	for n := 1; n <= 6; n++ {
		oldID, newID := ch15Old(n), ch15New(n)
		wantUp := fmt.Sprintf("UPDATE factory_machines SET machine_id = '%s' WHERE machine_id = '%s';", newID, oldID)
		wantDown := fmt.Sprintf("UPDATE factory_machines SET machine_id = '%s' WHERE machine_id = '%s';", oldID, newID)
		if !strings.Contains(up, wantUp) {
			t.Errorf("Up missing %q", wantUp)
		}
		if !strings.Contains(down, wantDown) {
			t.Errorf("Down missing %q", wantDown)
		}
		if n > 4 {
			continue
		}
		wantMachines := fmt.Sprintf("UPDATE machines SET id = '%s' WHERE id = '%s';", oldID, newID)
		if !strings.Contains(down, wantMachines) {
			t.Errorf("Down missing %q", wantMachines)
		}
	}
}

// Down restores 00002 / 00070 exactly.
func checkCh15DownRestores(t *testing.T, down string) {
	t.Helper()
	for _, want := range []string{
		fmt.Sprintf("UPDATE machines SET id = '%s', machine_value = 720000, lifespan_days = 500, productive_power = 1 WHERE id = '%s';", ch15Old(5), ch15New(5)),
		fmt.Sprintf("UPDATE machines SET id = '%s', name = 'Power-loom', working_tool = 'shuttle', productive_power = 366, hand_labour_per_unit = 4400 WHERE id = '%s';", ch15Old(6), ch15New(6)),
		fmt.Sprintf("UPDATE factories SET name = 'Manchester Weaving Shed', prime_mover_horsepower = 60, attendants = 1, living_labour_per_attendant = 480, necessary_labour_per_attendant = 240 WHERE id = '%s';", ch15MillID),
	} {
		if !strings.Contains(down, want) {
			t.Errorf("Down missing %q", want)
		}
	}
}

// Up and Down are symmetric, Up carries none of 00002's invented loom figures,
// and neither direction creates or removes rows.
func checkCh15Hygiene(t *testing.T, up, down string) {
	t.Helper()
	if u, d := strings.Count(up, "UPDATE"), strings.Count(down, "UPDATE"); u != 13 || d != 13 {
		t.Errorf("UPDATE count: Up %d, Down %d, want 13 each", u, d)
	}
	for _, stale := range []string{"720000", "'Power-loom'", "4400", "= 366"} {
		if strings.Contains(up, stale) {
			t.Errorf("Up must not contain %q", stale)
		}
	}
	for _, verb := range []string{"DELETE", "INSERT"} {
		if strings.Contains(up, verb) || strings.Contains(down, verb) {
			t.Errorf("migration must not contain %s", verb)
		}
	}
}
