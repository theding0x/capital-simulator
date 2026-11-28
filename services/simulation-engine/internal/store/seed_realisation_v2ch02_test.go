package store

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// Every package-level identifier here is prefixed v2ch02 so it cannot collide with
// the ch12*/ch15*/v2ch16* identifiers of the other seed-fidelity files.

// The two ids are the productive_circuits rows 00028 seeds. They are 20 characters
// ("5eed" + 12 zeros + chapter 02 + a/b), shorter than the 5eed00000000000000<CC>
// layout later chapters use; 00028 is append-only, so 00076 keeps them as they are.
const (
	v2ch02SeedPrefix = "5eed"
	v2ch02SimpleID   = "5eed000000000000020a" // Spinning Mill 1871, simple reproduction (00028)
	v2ch02ExtendedID = "5eed000000000000020b" // Spinning Mill 1871, extended reproduction (00028)
)

var (
	v2ch02Quoted = regexp.MustCompile(`'(5eed[0-9a-f]+)'`)
	v2ch02Insert = regexp.MustCompile(`(?s)INSERT INTO (\w+)\s*\(([^)]*)\)\s*VALUES\s*(.*?);`)
	v2ch02Tuple  = regexp.MustCompile(`\(\s*'(5eed[0-9a-f]+)'\s*,\s*(\d+)\s*,\s*(\d+|'[^']*')\s*\)`)
	v2ch02Delete = regexp.MustCompile(`DELETE FROM (\w+) WHERE productive_circuit_id IN \(([^)]*)\);`)
)

func v2ch02Read(t *testing.T, name string) (up, down string) {
	t.Helper()
	raw, err := migrationsFS.ReadFile("migrations/" + name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	upRaw, downRaw, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatalf("%s has no Down section", name)
	}
	return v2ch02SQLOnly(upRaw), v2ch02SQLOnly(downRaw)
}

// v2ch02SQLOnly drops comment lines so assertions see executable SQL only.
func v2ch02SQLOnly(s string) string {
	var kept []string
	for _, line := range strings.Split(s, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// 00076 gives both Spinning Mill 1871 exemplars of 00028 a C' of £500 (42,200
// capital-value + 7,800 surplus-value, legacy pence) and a full sale of it.
func TestMigration00076_V2Ch02RealisationSeed(t *testing.T) {
	t.Parallel()
	up, down := v2ch02Read(t, "00076_v2_ch02_realisation_seed.sql")

	inserts := v2ch02Insert.FindAllStringSubmatch(up, -1)
	if len(inserts) != 2 {
		t.Fatalf("got %d INSERT statements, want 2 (commodities, realisations)", len(inserts))
	}
	byTable := map[string][][]string{}
	for _, in := range inserts {
		byTable[in[1]] = v2ch02Tuple.FindAllStringSubmatch(in[3], -1)
	}
	commodities := byTable["productive_circuit_commodities"]
	realisations := byTable["productive_circuit_realisations"]
	if len(commodities) != 2 || len(realisations) != 2 {
		t.Fatalf("rows: commodities=%d realisations=%d, want 2 and 2", len(commodities), len(realisations))
	}

	for i, id := range []string{v2ch02SimpleID, v2ch02ExtendedID} {
		c, r := commodities[i], realisations[i]
		if c[1] != id || r[1] != id {
			t.Fatalf("row %d ids = %s / %s, want %s", i, c[1], r[1], id)
		}
		if !strings.HasPrefix(id, v2ch02SeedPrefix) {
			t.Errorf("id %s lacks the %s seed prefix", id, v2ch02SeedPrefix)
		}
		v2ch02CheckSale(t, id, c, r)
	}
	v2ch02CheckParentsExist(t)
	v2ch02CheckDown(t, down)
}

// v2ch02CheckParentsExist requires both seeded ids to be productive_circuits rows
// inserted by 00028, so the foreign keys of 00075 hold on a fresh volume.
func v2ch02CheckParentsExist(t *testing.T) {
	t.Helper()
	up, _ := v2ch02Read(t, "00028_v2_ch02_seed.sql")
	var parents string
	for _, in := range v2ch02Insert.FindAllStringSubmatch(up, -1) {
		if in[1] == "productive_circuits" {
			parents += in[3]
		}
	}
	for _, id := range []string{v2ch02SimpleID, v2ch02ExtendedID} {
		if !strings.Contains(parents, "'"+id+"'") {
			t.Errorf("%s is not a productive_circuits id seeded by 00028", id)
		}
	}
}

// v2ch02CheckSale rebuilds C' from the seeded columns and requires the seeded
// sale to pass the domain guard: the sale is C' "at its value and in its
// entirety", while the capital-value part alone would be rejected.
func v2ch02CheckSale(t *testing.T, id string, commodity, sale []string) {
	t.Helper()
	original, _ := strconv.ParseInt(commodity[2], 10, 64)
	surplus, _ := strconv.ParseInt(commodity[3], 10, 64)
	realised, _ := strconv.ParseInt(sale[2], 10, 64)
	c := &circulation.ExitCommodityCapital{
		ProductiveCircuitID: circulation.ProductiveCircuitID(id),
		ValueOriginal:       circulation.Pence(original),
		ValueSurplus:        circulation.Pence(surplus),
	}
	if err := c.Validate(); err != nil {
		t.Errorf("%s: seeded C' invalid: %v", id, err)
	}
	if original != 42200 || surplus != 7800 || c.Total() != 50000 {
		t.Errorf("%s: C' = %d + %d = %d, want 42200 + 7800 = 50000 (£500)", id, original, surplus, c.Total())
	}
	if err := circulation.ValidateRealisation(c, circulation.Pence(realised)); err != nil {
		t.Errorf("%s: seeded sale of %d fails the realisation guard: %v", id, realised, err)
	}
	if realised != 50000 {
		t.Errorf("%s: realised = %d, want 50000", id, realised)
	}
	if err := circulation.ValidateRealisation(c, c.ValueOriginal); !errors.Is(err, circulation.ErrPartialRealisation) {
		t.Errorf("%s: capital-value yarn alone: got %v, want ErrPartialRealisation", id, err)
	}
}

// v2ch02CheckDown requires the Down to delete only the two seeded ids from both
// tables, realisations first, and never a bare DELETE FROM.
func v2ch02CheckDown(t *testing.T, down string) {
	t.Helper()
	deletes := v2ch02Delete.FindAllStringSubmatch(down, -1)
	if got := strings.Count(down, "DELETE FROM"); got != 2 || len(deletes) != 2 {
		t.Fatalf("Down has %d DELETE FROM, %d id-scoped; want 2 and 2", got, len(deletes))
	}
	if deletes[0][1] != "productive_circuit_realisations" || deletes[1][1] != "productive_circuit_commodities" {
		t.Errorf("Down order = %s, %s; want realisations then commodities", deletes[0][1], deletes[1][1])
	}
	for _, d := range deletes {
		var ids []string
		for _, m := range v2ch02Quoted.FindAllStringSubmatch(d[2], -1) {
			ids = append(ids, m[1])
		}
		if len(ids) != 2 || ids[0] != v2ch02SimpleID || ids[1] != v2ch02ExtendedID {
			t.Errorf("Down on %s deletes %v, want exactly [%s %s]", d[1], ids, v2ch02SimpleID, v2ch02ExtendedID)
		}
	}
}

// 00075 declares both tables 1:1 with productive_circuits (primary key = foreign
// key), and its Down drops them in dependency order.
func TestMigration00075_V2Ch02CircuitRealisationDDL(t *testing.T) {
	t.Parallel()
	up, down := v2ch02Read(t, "00075_v2_ch02_circuit_realisation.sql")
	for _, table := range []string{"productive_circuit_commodities", "productive_circuit_realisations"} {
		stmt := regexp.MustCompile(`(?s)CREATE TABLE ` + table + `\s*\((.*?)\) ENGINE`).FindStringSubmatch(up)
		if stmt == nil {
			t.Fatalf("no CREATE TABLE %s", table)
		}
		for _, want := range []string{
			"PRIMARY KEY (productive_circuit_id)",
			"FOREIGN KEY (productive_circuit_id)",
			"REFERENCES productive_circuits (id)",
		} {
			if !strings.Contains(stmt[1], want) {
				t.Errorf("%s lacks %q", table, want)
			}
		}
	}
	for _, col := range [][2]string{
		{"value_original", "BIGINT"}, {"value_surplus", "BIGINT"},
		{"realised_pence", "BIGINT"}, {"sold_at", "DATETIME(6)"},
	} {
		if !regexp.MustCompile(col[0] + `\s+` + regexp.QuoteMeta(col[1]) + `\s+NOT NULL`).MatchString(up) {
			t.Errorf("column %s %s is not declared NOT NULL", col[0], col[1])
		}
	}
	if strings.Index(down, "productive_circuit_realisations") > strings.Index(down, "productive_circuit_commodities") ||
		!strings.Contains(down, "DROP TABLE IF EXISTS productive_circuit_realisations") ||
		!strings.Contains(down, "DROP TABLE IF EXISTS productive_circuit_commodities") {
		t.Errorf("Down must drop both tables; got:\n%s", down)
	}
}
