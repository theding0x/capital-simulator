package store

import (
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/pkg/surplusrate"
	val "github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// Every package-level identifier here is prefixed v2ch16 (or checkV2Ch16): the Ch. 15
// seed-fidelity file owns ch15*, the Ch. 12 one owns ch12*.

const (
	// v2ch16WeekMinutes is a week of seven 1,440-minute days: every seeded turnover
	// period is a whole number of these.
	v2ch16WeekMinutes int64 = 7 * 24 * 60

	// v2ch16YearMinutes is Ch. 16's year: Marx reckons "a year of 50 weeks" (sec. I, and
	// Ch. 15 likewise). Only the seeds apply it: n is a client-supplied input.
	v2ch16YearMinutes int64 = 50 * v2ch16WeekMinutes

	v2ch16Advances  = "variable_capital_advances"
	v2ch16Rates     = "per_circuit_surplus_rates"
	v2ch16Annuals   = "annual_surplus_rates"
	v2ch16Contrasts = "annual_surplus_rate_contrasts"
)

// v2ch16Cols names the nine columns 00074 restates, in the order of v2ch16Fixture.cols:
// the table, which of the fixture's three ids holds the row, and the column.
var v2ch16Cols = [9]struct {
	table string
	id    int
	col   string
}{
	{v2ch16Advances, 0, "pence"},
	{v2ch16Advances, 0, "turnover_period_minutes"},
	{v2ch16Advances, 0, "turnovers_per_year_basis_points"},
	{v2ch16Rates, 1, "surplus_pence"},
	{v2ch16Rates, 1, "variable_pence"},
	{v2ch16Rates, 1, "basis_points"},
	{v2ch16Annuals, 2, "n_basis_points"},
	{v2ch16Annuals, 2, "per_circuit_basis_points"},
	{v2ch16Annuals, 2, "annual_basis_points"},
}

// v2ch16Fixture is one seeded circuit. ids are its advance, per-circuit-rate and
// annual-rate rows. advance, surplus and weeks are the inputs; cols is the literal end
// state of the nine restated columns that derive from them.
type v2ch16Fixture struct {
	name             string
	ids              [3]string
	advance, surplus val.Pence
	weeks            int64
	cols             [9]int64
}

// v2ch16Fixtures returns the four seeded circuits. A and B are Marx's own (M, Ch. 16
// sec. I) and 00074 must not move them. The spinning mill is Vol. II Ch. 1's weekly 50
// pounds of wages and 78 pounds of surplus (M), its one-week turnover period assumed (I).
// The locomotive works is those weekly figures laid out for Ch. 12's twelve weeks (D).
func v2ch16Fixtures() (a, b, spinner, loco v2ch16Fixture) {
	wage, surplus := money.FromLSD(50, 0, 0), money.FromLSD(78, 0, 0)
	a = v2ch16Fixture{
		name: "A", ids: [3]string{"5eed0000000000001601", "5eed0000000000001607", "5eed0000000000001611"},
		advance: money.FromLSD(500, 0, 0), surplus: money.FromLSD(500, 0, 0), weeks: 5,
		cols: [9]int64{120000, 50400, 100000, 120000, 120000, 10000, 100000, 10000, 100000},
	}
	b = v2ch16Fixture{
		name: "B", ids: [3]string{"5eed0000000000001602", "5eed0000000000001608", "5eed0000000000001612"},
		advance: money.FromLSD(5000, 0, 0), surplus: money.FromLSD(5000, 0, 0), weeks: 50,
		cols: [9]int64{1200000, 504000, 10000, 1200000, 1200000, 10000, 10000, 10000, 10000},
	}
	spinner = v2ch16Fixture{
		name: "Spinning Mill", ids: [3]string{"5eed0000000000001603", "5eed0000000000001609", "5eed0000000000001613"},
		advance: wage, surplus: surplus, weeks: 1,
		cols: [9]int64{12000, 10080, 500000, 18720, 12000, 15600, 500000, 15600, 780000},
	}
	loco = v2ch16Fixture{
		name: "Locomotive Works", ids: [3]string{"5eed0000000000001605", "5eed0000000000001610", "5eed0000000000001614"},
		advance: 12 * wage, surplus: 12 * surplus, weeks: 12,
		// n = 50/12 rounds to 41,667 bp, so the annual rate is 65,001, 1 bp over 650%.
		cols: [9]int64{144000, 120960, 41667, 224640, 144000, 15600, 41667, 15600, 65001},
	}
	return a, b, spinner, loco
}

// derived recomputes cols with the domain functions: n from the 50-week year and the
// period (round half up), s/v from s and v, the annual rate from n and s/v.
func (f v2ch16Fixture) derived(t *testing.T) [9]int64 {
	t.Helper()
	period := f.weeks * v2ch16WeekMinutes
	n := surplusrate.RoundHalfUpDiv(10_000*v2ch16YearMinutes, period)
	rate, err := val.ComputePerCircuitRate("x", f.surplus, f.advance)
	if err != nil {
		t.Fatalf("%s: ComputePerCircuitRate: %v", f.name, err)
	}
	annual := val.ComputeAnnualSurplusRate("x", n, rate.BasisPoints, "1871").AnnualBasisPoints
	adv, sur := int64(f.advance), int64(f.surplus)
	return [9]int64{adv, period, n, sur, adv, rate.BasisPoints, n, rate.BasisPoints, annual}
}

// --- mini SQL harness ----------------------------------------------------------
//
// 00051 inserts the rows and 00074 updates them by seed id. Parsing the inserts and
// applying the updates proves that every id pairs with the right table and column (a
// typo'd id is a silent zero-row UPDATE in real SQL) and that Down restores 00051's rows
// exactly, which string matching alone cannot show.

// v2ch16DB maps table -> id -> column -> literal (quotes stripped).
type v2ch16DB map[string]map[string]map[string]string

var (
	v2ch16Insert  = regexp.MustCompile(`(?s)INSERT INTO (\w+)\s*\(([^)]*)\)\s*VALUES(.*?);`)
	v2ch16Tuple   = regexp.MustCompile(`\(([^()]*)\)`)
	v2ch16Update  = regexp.MustCompile(`^UPDATE (\w+) SET (.+) WHERE id = '([^']+)';$`)
	v2ch16WhereID = regexp.MustCompile(`WHERE id = '([^']+)'`)
)

// v2ch16SQLOnly drops whole-line comments (00051 has some inside VALUES).
func v2ch16SQLOnly(s string) string {
	var kept []string
	for _, line := range strings.Split(s, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

func v2ch16Split(s string) []string {
	parts := strings.Split(s, ",")
	for i, p := range parts {
		parts[i] = strings.Trim(strings.TrimSpace(p), "'")
	}
	return parts
}

func v2ch16Parse(t *testing.T, sql string) v2ch16DB {
	t.Helper()
	db := v2ch16DB{}
	for _, m := range v2ch16Insert.FindAllStringSubmatch(v2ch16SQLOnly(sql), -1) {
		table, cols := m[1], v2ch16Split(m[2])
		for _, tup := range v2ch16Tuple.FindAllStringSubmatch(m[3], -1) {
			vals := v2ch16Split(tup[1])
			if len(vals) != len(cols) {
				t.Fatalf("%s: %d values for %d columns in %q", table, len(vals), len(cols), tup[0])
			}
			row := map[string]string{}
			for i, c := range cols {
				row[c] = vals[i]
			}
			if db[table] == nil {
				db[table] = map[string]map[string]string{}
			}
			db[table][row["id"]] = row
		}
	}
	return db
}

func (db v2ch16DB) clone() v2ch16DB {
	out := v2ch16DB{}
	for table, rows := range db {
		out[table] = map[string]map[string]string{}
		for id, row := range rows {
			out[table][id] = map[string]string{}
			for col, lit := range row {
				out[table][id][col] = lit
			}
		}
	}
	return out
}

// v2ch16Apply runs a migration section against db. Only one-line
// "UPDATE t SET c = v, ... WHERE id = 'x';" statements are understood: anything else
// (an INSERT, a DELETE, a multi-line statement) fails the test as unsupported.
func v2ch16Apply(t *testing.T, db v2ch16DB, section string) {
	t.Helper()
	for _, line := range strings.Split(v2ch16SQLOnly(section), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		m := v2ch16Update.FindStringSubmatch(line)
		if m == nil {
			t.Fatalf("unsupported statement (want one-line UPDATE ... WHERE id): %q", line)
		}
		row, ok := db[m[1]][m[3]]
		if !ok {
			t.Fatalf("UPDATE matches no %s row with id %s: %q", m[1], m[3], line)
		}
		for _, set := range strings.Split(m[2], ",") {
			col, lit, ok := strings.Cut(strings.TrimSpace(set), " = ")
			if _, exists := row[col]; !ok || !exists {
				t.Fatalf("%s.%s does not exist: %q", m[1], col, line)
			}
			row[col] = strings.Trim(lit, "'")
		}
	}
}

func (db v2ch16DB) num(t *testing.T, table, id, col string) int64 {
	t.Helper()
	n, err := strconv.ParseInt(db[table][id][col], 10, 64)
	if err != nil {
		t.Fatalf("%s.%s on %s = %q: %v", table, col, id, db[table][id][col], err)
	}
	return n
}

// v2ch16Read returns a migration's Up and Down sections. It insists the file opens with
// the Up marker and holds exactly one Down marker, so a deleted marker cannot pass.
func v2ch16Read(t *testing.T, name string) (up, down string) {
	t.Helper()
	raw, err := migrationsFS.ReadFile("migrations/" + name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	text := string(raw)
	if !strings.HasPrefix(strings.TrimSpace(text), "-- +goose Up") || strings.Count(text, "-- +goose Down") != 1 {
		t.Fatalf(`%s must open with "-- +goose Up" and hold exactly one "-- +goose Down"`, name)
	}
	up, down, _ = strings.Cut(text, "-- +goose Down")
	return up, down
}

// --- the migration -----------------------------------------------------------

// 00074 restates two of 00051's four Vol. II Ch. 16 fixtures in place (issue #450): the
// spinning mill at Marx's Ch. 1 weekly 50 and 78 pounds, the locomotive works at those
// figures over Ch. 12's twelve weeks, both with n taken from Ch. 16's 50-week year.
// Capitals A and B were already Marx's own and must not move.
func TestMigration00074_V2Ch16SeedFidelity(t *testing.T) {
	t.Parallel()
	seedUp, _ := v2ch16Read(t, "00051_v2_ch16_seed.sql")
	up, down := v2ch16Read(t, "00074_v2_ch16_seed_fidelity.sql")
	a, b, spinner, loco := v2ch16Fixtures()

	pre := v2ch16Parse(t, seedUp)
	post := pre.clone()
	v2ch16Apply(t, post, up)

	for _, f := range []v2ch16Fixture{a, b, spinner, loco} {
		checkV2Ch16EndState(t, post, f)
	}
	checkV2Ch16Untouched(t, pre, post, a, b)
	checkV2Ch16Reproduction(t, post, a, b, spinner, loco)
	checkV2Ch16Down(t, pre, post, down)
	checkV2Ch16Ids(t, up, down, spinner, loco)
}

// After Up every restated column holds its literal, and the literals are the ones the
// domain recomputes on the 50-week year, so the year, the rates and the rounding are all
// pinned and a wrong formula here cannot agree with a wrong migration.
func checkV2Ch16EndState(t *testing.T, post v2ch16DB, f v2ch16Fixture) {
	t.Helper()
	if got := f.derived(t); got != f.cols {
		t.Errorf("%s: literals %v != recomputed %v on a 50-week year", f.name, f.cols, got)
	}
	for i, c := range v2ch16Cols {
		if got := post.num(t, c.table, f.ids[c.id], c.col); got != f.cols[i] {
			t.Errorf("%s: %s.%s = %d, want %d", f.name, c.table, c.col, got, f.cols[i])
		}
	}
}

// Capitals A and B and their A-vs-B contrast are Marx's own: Up leaves every cell of
// their rows (created_at and the capital id included) as 00051 seeded it.
func checkV2Ch16Untouched(t *testing.T, pre, post v2ch16DB, a, b v2ch16Fixture) {
	t.Helper()
	for _, f := range []v2ch16Fixture{a, b} {
		for _, c := range v2ch16Cols {
			id := f.ids[c.id]
			if !reflect.DeepEqual(pre[c.table][id], post[c.table][id]) {
				t.Errorf("%s: %s row %s changed: %v -> %v", f.name, c.table, id, pre[c.table][id], post[c.table][id])
			}
		}
	}
	if !reflect.DeepEqual(pre[v2ch16Contrasts], post[v2ch16Contrasts]) {
		t.Error("the A-vs-B contrast row must be untouched")
	}
}

// v2ch16Near reports |got-want| <= tol.
func v2ch16Near(got, want, tol int64) bool {
	d := got - want
	if d < 0 {
		d = -d
	}
	return d <= tol
}

// What the seeded rows reproduce over a year, read back through the domain. A and B each
// pay 5,000 pounds of wages and make 5,000 of surplus (Ch. 16 sec. I: only the rate
// differs). The spinner pays 50 x 50 = 2,500 and makes 50 x 78 = 3,900 pounds exactly.
// The locomotive works' n is rounded, so it lands near the exact 2,500 and 3,900: within
// half a basis point of n on its 144,000p advance, 7.2p for wages and 7.2 x 1.56 = 11.2p
// for surplus.
func checkV2Ch16Reproduction(t *testing.T, post v2ch16DB, a, b, spinner, loco v2ch16Fixture) {
	t.Helper()
	repro := func(f v2ch16Fixture) (wages, surplus int64) {
		r := val.ComputeVariableCapitalReproduction(val.VariableCapitalAdvance{
			Pence:                       val.Pence(post.num(t, v2ch16Advances, f.ids[0], "pence")),
			TurnoversPerYearBasisPoints: post.num(t, v2ch16Advances, f.ids[0], "turnovers_per_year_basis_points"),
		}, post.num(t, v2ch16Rates, f.ids[1], "basis_points"), v2ch16YearMinutes)
		return int64(r.TotalVariablePaidAnnualPence), int64(r.TotalSurplusAnnualPence)
	}
	pounds := func(n int64) int64 { return int64(money.FromLSD(n, 0, 0)) }
	for _, f := range []v2ch16Fixture{a, b} {
		if w, s := repro(f); w != pounds(5000) || s != pounds(5000) {
			t.Errorf("%s reproduces wages %d, surplus %d; want 5,000 pounds each", f.name, w, s)
		}
	}
	if w, s := repro(spinner); w != pounds(50*50) || s != pounds(50*78) {
		t.Errorf("spinner reproduces wages %d, surplus %d; want 2,500 and 3,900 pounds", w, s)
	}
	if w, s := repro(loco); !v2ch16Near(w, pounds(50*50), 8) || !v2ch16Near(s, pounds(50*78), 12) {
		t.Errorf("locomotive reproduces wages %d, surplus %d; want 2,500 and 3,900 pounds to within 8p and 12p", w, s)
	}
}

// Down restores 00051's rows exactly: applying it to the post-Up state gives back the
// parsed 00051 rows, and Up really did change something.
func checkV2Ch16Down(t *testing.T, pre, post v2ch16DB, down string) {
	t.Helper()
	if reflect.DeepEqual(post, pre) {
		t.Error("Up changed nothing")
	}
	restored := post.clone()
	v2ch16Apply(t, restored, down)
	if !reflect.DeepEqual(restored, pre) {
		t.Errorf("Down does not restore 00051:\n got  %v\n want %v", restored, pre)
	}
}

// Each direction is one UPDATE of each of the six restated rows and nothing else. It
// addresses them by seed id, never by a bare table (which would wipe user data on rollback).
func checkV2Ch16Ids(t *testing.T, up, down string, spinner, loco v2ch16Fixture) {
	t.Helper()
	want := append(append([]string{}, spinner.ids[:]...), loco.ids[:]...)
	sort.Strings(want)
	for name, section := range map[string]string{"Up": up, "Down": down} {
		var got []string
		for _, m := range v2ch16WhereID.FindAllStringSubmatch(v2ch16SQLOnly(section), -1) {
			got = append(got, m[1])
		}
		sort.Strings(got)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s ids = %v, want each of %v once", name, got, want)
		}
	}
}
