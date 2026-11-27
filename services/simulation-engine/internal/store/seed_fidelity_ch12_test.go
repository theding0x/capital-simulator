package store

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// Adjust here if the migration number or id range changes (also edit the SQL ids).
const (
	ch12Migration = "migrations/00073_v2_ch12_seed_fidelity.sql"
	ch12IDFirst   = 1203 // 5eed00000000000000 + "1203" = 22 chars, same width as 00045
)

const (
	ch12Prefix  = "5eed00000000000000"
	ch12DayMin  = 600 // Marx: "working-days ... of ten hours each"
	ch12Created = "1871-01-02 00:00:00.000000"
	ch12Week    = 7 // the spinner's week
	ch12Year    = 365

	// The Bakewell shortening, 4 years to 1: Marx's sheep "not fit for the
	// butcher until four or five years old" against "fattened when only one year old".
	ch12SheepOld = 4 * ch12Year
	ch12SheepNew = ch12Year

	// 00045 ids this migration re-points or restores.
	ch12LocomotiveID       = "5eed000000000000001202"
	ch12CattleConstraint   = "5eed000000000000001211"
	ch12BakewellConstraint = "5eed000000000000001214"
	ch12BakewellShorten    = "5eed000000000000001230"
	ch12BuilderFinancing   = "5eed000000000000001240"

	// The reasons Up writes. Every UPDATE of Up and Down is compared with a statement built
	// from these constants and the domain values, character for character.
	ch12SheepReason  = `Bakewell sheep may be fattened "when only one year old" (Marx, Ch. 12): the natural floor after breeding is one year (365 days), against four or five years before Bakewell`
	ch12CattleReason = `livestock maturation - Marx: "several years to raise big cattle" (Good: four- and five-year-old oxen). 1095 days (three years) is an illustrative figure, not his`

	// The two 00045 reason strings Down restores. Down's UPDATEs are compared with them
	// character for character, and each must also occur verbatim in 00045.
	ch12OldCattleReason   = "livestock maturation — cattle require approximately three years to reach market weight (pre-Bakewell)"
	ch12OldBakewellReason = "Bakewell selective breeding reduces cattle maturation from 3 years to ~1.5 years"
)

func ch12ID(off int) string { return fmt.Sprintf("%s%04d", ch12Prefix, ch12IDFirst+off) }

// offsets from ch12IDFirst; days are working days (calendar-like: a week is 7).
const (
	ch12ManOfWar = iota // 365   D: "one year or more"
	ch12Road            // 84    I: "a few months"
	ch12Railway         // 730   I: "a job of years"
	ch12Carpet          // 7     M: "about a week"
	ch12Gobelin         // 1095  I: "takes years"
	ch12Sheep           // 1460  D: "four or five years old" (low end)
	ch12House           // 168   I: Marx gives no number
)

type ch12Row struct {
	off           int
	ic, commodity string
	days          int64
}

var ch12Rows = []ch12Row{
	{ch12ManOfWar, "armoured-ship-yard-ic", "man-of-war", 365},
	{ch12Road, "country-road-works-ic", "country-road", 84},
	{ch12Railway, "railway-contractor-ic", "railway", 730},
	{ch12Carpet, "carpet-works-ic", "carpet", 7},
	{ch12Gobelin, "tapestry-works-ic", "gobelin-tapestry", 1095},
	{ch12Sheep, "dishley-grange-ic", "sheep", 1460},
	{ch12House, "london-speculative-builder-ic", "house", 168},
}

var (
	ch12Ws      = regexp.MustCompile(`\s+`)
	ch12Comma   = regexp.MustCompile(`\s*,\s*`)
	ch12IDToken = regexp.MustCompile(`'(5eed[0-9a-z]+)'`)
	ch12ICToken = regexp.MustCompile(`'([a-z0-9-]+-ic)'`)
	ch12BareDel = regexp.MustCompile(`DELETE FROM\s+\w+\s*;`)
	ch12WPRow   = regexp.MustCompile(`\('5eed[0-9a-z]+',\s*'([a-z0-9-]+)',\s*'([a-z0-9-]+)',\s*(\d+),\s*'(?:discrete|connected)',\s*\d+,`)
)

// ch12Parsed is one working_periods row read back out of a seed migration.
type ch12Parsed struct {
	ic, commodity string
	days          int64
}

// ch12Seeds bundles the neighbouring seed migrations the new one is checked against.
type ch12Seeds struct {
	raw45, raw47           string
	rows45, rows47, rows73 []ch12Parsed
}

func ch12Norm(s string) string {
	return strings.TrimSpace(ch12Comma.ReplaceAllString(ch12Ws.ReplaceAllString(s, " "), ","))
}

func ch12SQLOnly(s string) string { // drop full-line comments
	var kept []string
	for _, line := range strings.Split(s, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// ch12Has is the whitespace-normalised containment check, used only for the multi-line
// INSERT tuples, the DELETE, and the numeric facts of the 00045 pre-image.
func ch12Has(t *testing.T, sql, want string) {
	t.Helper()
	if !strings.Contains(ch12Norm(sql), ch12Norm(want)) {
		t.Errorf("missing %q", want)
	}
}

// ch12Stmt returns the ;-separated statement of sql (comments already stripped) that ends
// WHERE id = 'id', trimmed of surrounding whitespace and of its ';', and otherwise untouched.
func ch12Stmt(t *testing.T, sql, id string) string {
	t.Helper()
	suffix := "WHERE id = '" + id + "'"
	for _, s := range strings.Split(sql, ";") {
		if s = strings.TrimSpace(s); strings.HasSuffix(s, suffix) {
			return s
		}
	}
	t.Fatalf("no statement ends %s", suffix)
	return ""
}

// ch12Exact asserts that the statement for id is, character for character and with no
// whitespace normalisation, want. This pins every column an UPDATE sets and every byte
// of every reason text.
func ch12Exact(t *testing.T, sql, id, want string) string {
	t.Helper()
	got := ch12Stmt(t, sql, id)
	if got != want {
		t.Errorf("statement for %s differs\n got:  %q\n want: %q", id, got, want)
	}
	return got
}

func ch12Period(r ch12Row) wp.WorkingPeriod {
	return wp.WorkingPeriod{
		ID: wp.WorkingPeriodID(ch12ID(r.off)), IndustrialCapitalID: r.ic, CommodityID: r.commodity,
		WorkingDayCount: r.days, Mode: wp.ModeConnected, WorkingDayLengthMinutes: ch12DayMin,
	}
}

func ch12Tuple(p wp.WorkingPeriod) string {
	return fmt.Sprintf("('%s', '%s', '%s', %d, '%s', %d, '%s')",
		p.ID, p.IndustrialCapitalID, p.CommodityID, p.WorkingDayCount, p.Mode, p.WorkingDayLengthMinutes, ch12Created)
}

func ch12Days(off int) int64 {
	for _, r := range ch12Rows {
		if r.off == off {
			return r.days
		}
	}
	panic(fmt.Sprintf("ch12: no new row at offset %d", off))
}

func ch12Read(t *testing.T, name string) string {
	t.Helper()
	b, err := migrationsFS.ReadFile(name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(b)
}

func ch12Atoi(t *testing.T, s string) int64 {
	t.Helper()
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		t.Fatalf("parse %q: %v", s, err)
	}
	return n
}

// ch12Values returns the VALUES body of the first INSERT INTO table in sql.
func ch12Values(t *testing.T, sql, table string) string {
	t.Helper()
	re := regexp.MustCompile(`(?s)INSERT INTO ` + regexp.QuoteMeta(table) + `\s*\([^)]*\)\s*VALUES(.*?);`)
	m := re.FindStringSubmatch(sql)
	if m == nil {
		t.Fatalf("no INSERT INTO %s", table)
	}
	return m[1]
}

func ch12ParseRows(t *testing.T, up string) []ch12Parsed {
	t.Helper()
	var out []ch12Parsed
	for _, m := range ch12WPRow.FindAllStringSubmatch(ch12Values(t, up, "working_periods"), -1) {
		out = append(out, ch12Parsed{ic: m[1], commodity: m[2], days: ch12Atoi(t, m[3])})
	}
	return out
}

func ch12UpOf(raw string) string {
	up, _, _ := strings.Cut(raw, "-- +goose Down")
	return ch12SQLOnly(up)
}

// 00073 restates the Vol. II Ch. 12 Bakewell and London-builder fixtures and adds the
// missing reference working periods (issue #450). Every figure below is recomputed from
// the domain or from Marx's own words, not typed twice: the Bakewell shortening is
// validated by the domain type, the financing is built with money.FromLSD, the whole set
// is replayed through the in-memory store with its natural floors live, and each UPDATE
// of Up and Down is compared whole with a statement built from those values.
func TestMigration00073_V2Ch12SeedFidelity(t *testing.T) {
	t.Parallel()
	upRaw, downRaw, ok := strings.Cut(ch12Read(t, ch12Migration), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}
	up, down := ch12SQLOnly(upRaw), ch12SQLOnly(downRaw)
	raw45 := ch12Read(t, "migrations/00045_v2_ch12_seed.sql")
	raw47 := ch12Read(t, "migrations/00047_v2_ch13_seed.sql")
	s := ch12Seeds{
		raw45: raw45, raw47: raw47,
		rows45: ch12ParseRows(t, ch12UpOf(raw45)),
		rows47: ch12ParseRows(t, ch12UpOf(raw47)),
		rows73: ch12ParseRows(t, up),
	}

	checkCh12Periods(t, up)
	checkCh12Ids(t, up, down, s.raw45)
	checkCh12Collisions(t, s)
	checkCh12MarxOrderings(t, s.rows45)
	checkCh12Bakewell(t, up)
	checkCh12Constraints(t, up)
	checkCh12MemoryRoundTrip(t)
	checkCh12Builder(t, up)
	checkCh12Down(t, down, s.raw45)
	checkCh12Hygiene(t, up, down)
}

func checkCh12Periods(t *testing.T, up string) {
	t.Helper()
	for _, r := range ch12Rows {
		p := ch12Period(r)
		if err := p.Validate(); err != nil {
			t.Errorf("%s: %v", p.CommodityID, err)
		}
		ch12Has(t, up, ch12Tuple(p))
	}
}

func checkCh12Ids(t *testing.T, up, down, seed45 string) {
	t.Helper()
	want := map[string]bool{}
	for _, r := range ch12Rows {
		want[ch12ID(r.off)] = true
	}
	if len(want) != len(ch12Rows) {
		t.Errorf("new ids are not unique: %d distinct for %d rows", len(want), len(ch12Rows))
	}
	in45 := map[string]bool{}
	for _, m := range ch12IDToken.FindAllStringSubmatch(seed45, -1) {
		in45[m[1]] = true
	}
	for _, section := range []string{up, down} {
		for _, m := range ch12IDToken.FindAllStringSubmatch(section, -1) {
			id := m[1]
			if len(id) != 22 || !strings.HasPrefix(id, ch12Prefix+"12") {
				t.Errorf("id %q: want 22 chars with prefix %s12", id, ch12Prefix)
			}
			if _, err := hex.DecodeString(id); err != nil {
				t.Errorf("id %q is not hex: %v", id, err)
			}
			if !in45[id] && !want[id] {
				t.Errorf("id %q is neither a 00045 id nor one of the seven new ids", id)
			}
		}
	}
	// A new period is named once by its INSERT; the sheep and house periods are also
	// named once more, as the new parent of the re-pointed shortening and financing.
	alsoParent := map[int]int{ch12Sheep: 1, ch12House: 1}
	for _, r := range ch12Rows {
		id := ch12ID(r.off)
		if in45[id] {
			t.Errorf("new id %s collides with 00045", id)
		}
		if got, wantUp := strings.Count(up, "'"+id+"'"), 1+alsoParent[r.off]; got != wantUp {
			t.Errorf("Up: %s appears %d times, want %d (INSERT plus parent references)", id, got, wantUp)
		}
		if got := strings.Count(down, "'"+id+"'"); got != 1 {
			t.Errorf("Down: %s appears %d times, want 1 (DELETE only)", id, got)
		}
	}
}

func checkCh12Collisions(t *testing.T, s ch12Seeds) {
	t.Helper()
	old := map[string]bool{}
	for _, src := range []string{s.raw45, s.raw47} {
		for _, m := range ch12ICToken.FindAllStringSubmatch(src, -1) {
			old[m[1]] = true
		}
	}
	for _, r := range s.rows73 {
		if old[r.ic] {
			t.Errorf("industrial_capital_id %q already used by 00045 or 00047 (the Ch. 13 panel filters by it)", r.ic)
		}
	}
	if len(s.rows45) != 2 || len(s.rows73) != 7 {
		t.Errorf("Ch. 12 working periods: %d in 00045 + %d in 00073, want 2 + 7 = 9", len(s.rows45), len(s.rows73))
	}
	seen := map[string]int{}
	for _, rows := range [][]ch12Parsed{s.rows45, s.rows47, s.rows73} {
		for _, r := range rows {
			seen[r.commodity]++
		}
	}
	// The nine products of the Ch. 12 spec. Corn, cattle and timber are covered by name only:
	// their Ch. 13 rows (00047) hold labour-days (10, 5 and 30 after 00071), not Ch. 12
	// working-period spans. The spec's second timber row (100 years) is not seeded.
	for _, c := range []string{"yarn", "locomotive", "corn", "cattle", "timber", "country-road", "railway", "carpet", "gobelin-tapestry"} {
		if seen[c] == 0 {
			t.Errorf("spec product %q has no working period in 00045, 00047 or 00073", c)
		}
	}
	for c, n := range seen {
		if n > 1 {
			t.Errorf("commodity %q has %d working periods; ?commodity_id= would return several", c, n)
		}
	}
}

func ch12DaysOf(t *testing.T, rows []ch12Parsed, commodity string) int64 {
	t.Helper()
	for _, r := range rows {
		if r.commodity == commodity {
			return r.days
		}
	}
	t.Fatalf("no 00045 working period for %q", commodity)
	return 0
}

// checkCh12MarxOrderings asserts what Marx's words fix: the spinner's week, the
// locomotive's three months (twelve of them), a man-of-war of a year or more, and railway
// and Gobelin of "years". He fixes no order among road, railway, Gobelin and man-of-war.
func checkCh12MarxOrderings(t *testing.T, rows45 []ch12Parsed) {
	t.Helper()
	spinner, loco := ch12DaysOf(t, rows45, "yarn"), ch12DaysOf(t, rows45, "locomotive")
	ship, road, rail := ch12Days(ch12ManOfWar), ch12Days(ch12Road), ch12Days(ch12Railway)
	carpet, gobelin := ch12Days(ch12Carpet), ch12Days(ch12Gobelin)
	for _, c := range []struct {
		ok  bool
		msg string
	}{
		{carpet == spinner && spinner == ch12Week, "carpet is the spinner's week (7 days): 'about a week'"},
		{loco == 12*ch12Week, "locomotive is twelve weeks (84 days): 'three months'"},
		{carpet < road && road < ship, "carpet < country road < man-of-war"},
		{loco < ship, "locomotive < man-of-war"},
		{ship >= ch12Year, "man-of-war is at least a year: 'one year or more'"},
		{rail >= 2*ch12Year && gobelin >= 2*ch12Year, "railway and Gobelin are 'years': at least two"},
		{carpet < gobelin, "carpet < Gobelin"},
	} {
		if !c.ok {
			t.Error(c.msg)
		}
	}
}

func checkCh12Bakewell(t *testing.T, up string) {
	t.Helper()
	s := wp.WorkingPeriodShortening{
		Kind: wp.ShorteningBakewellBreeding, OldWorkingDayCount: ch12SheepOld, NewWorkingDayCount: ch12SheepNew,
	}
	if err := s.Validate(); err != nil {
		t.Errorf("Bakewell shortening: %v", err)
	}
	if s.OldWorkingDayCount != 4*s.NewWorkingDayCount {
		t.Errorf("Bakewell: %d to %d is not four to one", s.OldWorkingDayCount, s.NewWorkingDayCount)
	}
	if s.OldWorkingDayCount < 4*ch12Year || s.OldWorkingDayCount > 5*ch12Year {
		t.Errorf("Bakewell old %d is outside 'four or five years old'", s.OldWorkingDayCount)
	}
	if s.NewWorkingDayCount != ch12Year {
		t.Errorf("Bakewell new %d is not 'only one year old'", s.NewWorkingDayCount)
	}
	if sheep := ch12Days(ch12Sheep); s.OldWorkingDayCount != sheep {
		t.Errorf("Bakewell old %d != the sheep period's %d days", s.OldWorkingDayCount, sheep)
	}
	ch12Exact(t, up, ch12BakewellShorten, fmt.Sprintf(
		"UPDATE working_period_shortenings SET working_period_id = '%s', old_working_day_count = %d, new_working_day_count = %d WHERE id = '%s'",
		ch12ID(ch12Sheep), s.OldWorkingDayCount, s.NewWorkingDayCount, ch12BakewellShorten))
}

func checkCh12Constraints(t *testing.T, up string) {
	t.Helper()
	nc := wp.NaturalWorkingPeriodConstraint{CommodityID: "sheep", MinWorkingDays: ch12Year, Reason: ch12SheepReason}
	if err := nc.Validate(); err != nil {
		t.Errorf("sheep constraint: %v", err)
	}
	// Bakewell reduces to, never below, the floor.
	if ch12SheepNew != nc.MinWorkingDays {
		t.Errorf("Bakewell new %d != the sheep floor %d", ch12SheepNew, nc.MinWorkingDays)
	}
	sheep := ch12Exact(t, up, ch12BakewellConstraint, fmt.Sprintf(
		"UPDATE natural_working_period_constraints SET commodity_id = '%s', min_working_days = %d, reason = '%s' WHERE id = '%s'",
		nc.CommodityID, nc.MinWorkingDays, nc.Reason, ch12BakewellConstraint))
	// What Marx said, as the readable statement of the reason.
	for _, words := range []string{"when only one year old", "four or five years"} {
		if !strings.Contains(sheep, words) {
			t.Errorf("sheep reason must carry Marx's words %q: %s", words, sheep)
		}
	}
	for _, gone := range []string{"bakewell-cattle", "547"} {
		if strings.Contains(up, gone) {
			t.Errorf("Up still mentions %q", gone)
		}
	}
	checkCh12CattleReason(t, up)
	for _, id := range []string{"1210", "1212", "1213", "1220", "1221", "1231", "1232"} {
		if strings.Contains(up, "'"+ch12Prefix+id+"'") {
			t.Errorf("Up touches %s%s, which this migration leaves alone", ch12Prefix, id)
		}
	}
}

// checkCh12CattleReason: the cattle row keeps its floor (1095 is the spec's figure, not
// Marx's) and only loses the "pre-Bakewell" framing.
func checkCh12CattleReason(t *testing.T, up string) {
	t.Helper()
	stmt := ch12Exact(t, up, ch12CattleConstraint, fmt.Sprintf(
		"UPDATE natural_working_period_constraints SET reason = '%s' WHERE id = '%s'", ch12CattleReason, ch12CattleConstraint))
	for _, words := range []string{"several years to raise big cattle", "illustrative"} {
		if !strings.Contains(stmt, words) {
			t.Errorf("cattle reason must carry %q: %s", words, stmt)
		}
	}
	if strings.Contains(strings.ToLower(stmt), "bakewell") {
		t.Errorf("cattle reason still mentions Bakewell: %s", stmt)
	}
	if strings.Contains(up, "min_working_days = 1095") {
		t.Error("the cattle floor (1095) must not be changed")
	}
}

func checkCh12MemoryRoundTrip(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	st := NewMemory()
	for _, c := range []struct {
		commodity string
		min       int64
	}{
		{"corn", 365}, {"cattle", 1095}, {"timber", 3650}, {"merino-wool", 365}, // 00045, unchanged
		{"sheep", 365}, // 00073: the Bakewell floor
	} {
		if _, err := st.CreateNaturalConstraint(ctx, wp.NaturalWorkingPeriodConstraint{
			CommodityID: c.commodity, MinWorkingDays: c.min, Reason: "fixture"}); err != nil {
			t.Fatalf("constraint %s: %v", c.commodity, err)
		}
	}
	periods := []wp.WorkingPeriod{
		{ID: "5eed000000000000001201", IndustrialCapitalID: "spinning-mill-1871-ic", CommodityID: "yarn",
			WorkingDayCount: ch12Week, Mode: wp.ModeDiscrete, WorkingDayLengthMinutes: ch12DayMin},
		{ID: wp.WorkingPeriodID(ch12LocomotiveID), IndustrialCapitalID: "locomotive-works-ic", CommodityID: "locomotive",
			WorkingDayCount: 12 * ch12Week, Mode: wp.ModeConnected, WorkingDayLengthMinutes: ch12DayMin},
	}
	for _, r := range ch12Rows {
		periods = append(periods, ch12Period(r))
	}
	for _, p := range periods { // none may be below its own natural floor
		if _, err := st.CreateWorkingPeriod(ctx, p); err != nil {
			t.Errorf("%s (%d days): %v", p.CommodityID, p.WorkingDayCount, err)
		}
	}
	// negative control: the floor is live
	_, err := st.CreateWorkingPeriod(ctx, wp.WorkingPeriod{ID: "ch12-negative-control", IndustrialCapitalID: "x",
		CommodityID: "sheep", WorkingDayCount: ch12Year - 1, Mode: wp.ModeConnected, WorkingDayLengthMinutes: ch12DayMin})
	if !errors.Is(err, wp.ErrPrematureDelivery) {
		t.Errorf("364-day sheep: err = %v, want ErrPrematureDelivery", err)
	}
	sheepID, houseID := wp.WorkingPeriodID(ch12ID(ch12Sheep)), wp.WorkingPeriodID(ch12ID(ch12House))
	if _, err := st.RecordShortening(ctx, sheepID, wp.WorkingPeriodShortening{
		Kind: wp.ShorteningBakewellBreeding, OldWorkingDayCount: ch12SheepOld, NewWorkingDayCount: ch12SheepNew}); err != nil {
		t.Errorf("record Bakewell shortening: %v", err)
	}
	own, borrowed := int64(money.FromLSD(100, 0, 0)), int64(money.FromLSD(2000, 0, 0))
	if _, err := st.RecordCreditFinancing(ctx, houseID, wp.WorkingPeriodCreditFinancing{
		OwnCapitalPence: own, BorrowedCapitalPence: borrowed}); err != nil {
		t.Errorf("record builder financing: %v", err)
	}
	if got, err := st.GetWorkingPeriod(ctx, houseID); err != nil || got.CreditFinancing == nil {
		t.Errorf("house period: err=%v financing=%v", err, got.CreditFinancing)
	}
}

// checkCh12Builder: Marx gives only the ratio ("an enterprise which exceeds his resources
// twenty to fifty times"); borrowed = 20 x own makes the enterprise 21 x own, inside the
// range on either reading. The pound sums are illustrative.
func checkCh12Builder(t *testing.T, up string) {
	t.Helper()
	own, borrowed := int64(money.FromLSD(100, 0, 0)), int64(money.FromLSD(2000, 0, 0))
	total := own + borrowed
	cf := wp.WorkingPeriodCreditFinancing{
		WorkingPeriodID: wp.WorkingPeriodID(ch12ID(ch12House)), OwnCapitalPence: own, BorrowedCapitalPence: borrowed,
	}
	if err := cf.Validate(); err != nil {
		t.Errorf("builder financing: %v", err)
	}
	if borrowed != 20*own || total != 21*own {
		t.Errorf("borrowed/own = %d/%d, total/own = %d/%d: want exactly 20 and 21", borrowed, own, total, own)
	}
	// Documentation of Marx's range: it cannot fail independently of the exact 20x/21x check above.
	if borrowed < 20*own || borrowed > 49*own || total < 20*own || total > 50*own {
		t.Errorf("financing is outside the window both readings share: 20 <= borrowed/own <= 49, 20 <= total/own <= 50")
	}
	ch12Exact(t, up, ch12BuilderFinancing, fmt.Sprintf(
		"UPDATE working_period_credit_financings SET working_period_id = '%s', own_capital_pence = %d, borrowed_capital_pence = %d WHERE id = '%s'",
		cf.WorkingPeriodID, own, borrowed, ch12BuilderFinancing))
	if strings.Contains(up, "mortgage_provider_id") {
		t.Error("Up must leave mortgage_provider_id as the Vol. III placeholder")
	}
}

func checkCh12Down(t *testing.T, down, seed45 string) {
	t.Helper()
	ch12Exact(t, down, ch12BuilderFinancing, fmt.Sprintf(
		"UPDATE working_period_credit_financings SET working_period_id = '%s', own_capital_pence = 2400, borrowed_capital_pence = 48000 WHERE id = '%s'",
		ch12LocomotiveID, ch12BuilderFinancing))
	ch12Exact(t, down, ch12CattleConstraint, fmt.Sprintf(
		"UPDATE natural_working_period_constraints SET reason = '%s' WHERE id = '%s'", ch12OldCattleReason, ch12CattleConstraint))
	ch12Exact(t, down, ch12BakewellConstraint, fmt.Sprintf(
		"UPDATE natural_working_period_constraints SET commodity_id = 'bakewell-cattle', min_working_days = 547, reason = '%s' WHERE id = '%s'",
		ch12OldBakewellReason, ch12BakewellConstraint))
	ch12Exact(t, down, ch12BakewellShorten, fmt.Sprintf(
		"UPDATE working_period_shortenings SET working_period_id = '%s', old_working_day_count = 1095, new_working_day_count = 547 WHERE id = '%s'",
		ch12LocomotiveID, ch12BakewellShorten))
	ids := make([]string, 0, len(ch12Rows))
	for _, r := range ch12Rows {
		ids = append(ids, "'"+ch12ID(r.off)+"'")
	}
	ch12Has(t, down, "DELETE FROM working_periods WHERE id IN ("+strings.Join(ids, ",")+")")
	if strings.Contains(down, "INSERT") {
		t.Error("Down must not INSERT")
	}
	if ch12BareDel.MatchString(down) {
		t.Error("Down has a DELETE without WHERE")
	}
	// Pre-image guard: every fact Down restores is what 00045 actually seeded. The two
	// reason strings must occur in 00045 verbatim, the numeric facts modulo layout.
	for _, reason := range []string{ch12OldBakewellReason, ch12OldCattleReason} {
		if !strings.Contains(seed45, "'"+reason+"'") {
			t.Errorf("00045 does not contain the reason %q verbatim", reason)
		}
	}
	for _, fact := range []string{
		fmt.Sprintf("'%s', 'bakewell-cattle', 547, '%s'", ch12BakewellConstraint, ch12OldBakewellReason),
		fmt.Sprintf("'%s', 'cattle', 1095, '%s'", ch12CattleConstraint, ch12OldCattleReason),
		fmt.Sprintf("'%s', '%s', 'bakewell_breeding', 1095, 547", ch12BakewellShorten, ch12LocomotiveID),
		fmt.Sprintf("'%s', '%s', 2400, 48000", ch12BuilderFinancing, ch12LocomotiveID),
	} {
		ch12Has(t, seed45, fact)
	}
}

func checkCh12Hygiene(t *testing.T, up, down string) {
	t.Helper()
	for _, c := range []struct {
		name, sql, verb string
		want            int
	}{
		{"Up", up, "INSERT INTO", 1}, {"Up", up, "UPDATE ", 4}, {"Up", up, "DELETE", 0},
		{"Down", down, "INSERT", 0}, {"Down", down, "UPDATE ", 4}, {"Down", down, "DELETE FROM", 1},
	} {
		if got := strings.Count(c.sql, c.verb); got != c.want {
			t.Errorf("%s: %q appears %d times, want %d", c.name, c.verb, got, c.want)
		}
	}
}
