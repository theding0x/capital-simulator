package store

import (
	"encoding/hex"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// 00112 restates the Vol. III Ch. 5 economy seed in Ch. 5's own figures (issue
// #450). Tags: M = Marx's figure, D = derived from his, I = a figure the seed
// picks because Marx gives none. The migrations run through a small interpreter
// for the statement shapes they use, so a guard that matches no row, a DELETE
// without a WHERE or a duplicate id fails here rather than on MySQL.
const (
	ch05Migration = "migrations/00112_v3_ch05_seed_fidelity.sql"
	ch05Seed      = "migrations/00010_v3_ch05_seed.sql"
	ch05IDPrefix  = "5eed00000000000000" + "05" // 18-char seed prefix + chapter 05; two more chars follow
)

var (
	ch05InsertRe = regexp.MustCompile(`^INSERT INTO economy_analyses \(id, kind, constant_capital, variable_capital, surplus_value, saving, old_profit_rate_bp, new_profit_rate_bp, profit_rate_gain_bp, created_at\) VALUES (.+)$`)
	ch05RowRe    = regexp.MustCompile(`\(\s*'(5eed[0-9a-f]+)'\s*,\s*'([a-z_]+)'\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*'[^']*'\s*\)`)
	ch05UpdateRe = regexp.MustCompile(`^UPDATE economy_analyses SET (.+?) WHERE id = '(5eed[0-9a-f]+)'((?: AND [a-z_]+ = \d+)*)$`)
	ch05DeleteRe = regexp.MustCompile(`^DELETE FROM economy_analyses WHERE id (?:= '(5eed[0-9a-f]+)'|IN \((.+)\))$`)
	ch05AssignRe = regexp.MustCompile(`^([a-z_]+) = (\d+)$`)
	ch05GuardRe  = regexp.MustCompile(` AND ([a-z_]+) = (\d+)`)
	ch05IDRe     = regexp.MustCompile(`'(5eed[0-9a-z]*)'`) // any id-shaped literal, valid or not
)

// ch05Row is one economy_analyses row; ch05DB is the table plus what ran on it.
type ch05Row struct {
	id, kind                              string
	c, v, s, saving, oldBP, newBP, gainBP int64
}

type ch05DB struct {
	rows    map[string]ch05Row
	updated map[string]bool
	deleted []string
	counts  map[string]int // statements run, by leading keyword
}

func ch05NewDB() *ch05DB {
	return &ch05DB{rows: map[string]ch05Row{}, updated: map[string]bool{}, counts: map[string]int{}}
}

// clone copies the rows only: the copy starts a fresh history of statements.
func (db *ch05DB) clone() *ch05DB {
	c := ch05NewDB()
	maps.Copy(c.rows, db.rows)
	return c
}

// col returns the numeric field an UPDATE column names, or nil if it is not one.
func (r *ch05Row) col(name string) *int64 {
	return map[string]*int64{
		"constant_capital": &r.c, "variable_capital": &r.v, "surplus_value": &r.s, "saving": &r.saving,
		"old_profit_rate_bp": &r.oldBP, "new_profit_rate_bp": &r.newBP, "profit_rate_gain_bp": &r.gainBP,
	}[name]
}

func ch05ID(seq int) string { return fmt.Sprintf("%s%02d", ch05IDPrefix, seq) }

// ch05Statements drops full-line comments, splits on semicolons (as goose does)
// and collapses whitespace.
func ch05Statements(section string) []string {
	var code []string
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			code = append(code, line)
		}
	}
	var out []string
	for _, stmt := range strings.Split(strings.Join(code, "\n"), ";") {
		if stmt = strings.Join(strings.Fields(stmt), " "); stmt != "" {
			out = append(out, stmt)
		}
	}
	return out
}

func ch05Exec(db *ch05DB, stmt string) error {
	keyword, _, _ := strings.Cut(stmt, " ")
	db.counts[keyword]++
	switch {
	case ch05InsertRe.MatchString(stmt):
		return db.insert(stmt)
	case keyword == "UPDATE":
		return db.update(stmt)
	case keyword == "DELETE":
		return db.remove(stmt)
	}
	return fmt.Errorf("unexpected statement (DDL or an unknown shape): %.70s", stmt)
}

func (db *ch05DB) insert(stmt string) error {
	values := ch05InsertRe.FindStringSubmatch(stmt)[1]
	if rest := strings.Trim(ch05RowRe.ReplaceAllString(values, ""), ", "); rest != "" {
		return fmt.Errorf("INSERT row not understood: %.70s", rest)
	}
	for _, m := range ch05RowRe.FindAllStringSubmatch(values, -1) {
		var n [7]int64
		for i := range n {
			n[i], _ = strconv.ParseInt(m[3+i], 10, 64) // \d+ by the pattern, small by construction
		}
		if _, dup := db.rows[m[1]]; dup {
			return fmt.Errorf("duplicate id %s on INSERT", m[1])
		}
		db.rows[m[1]] = ch05Row{id: m[1], kind: m[2], c: n[0], v: n[1], s: n[2], saving: n[3], oldBP: n[4], newBP: n[5], gainBP: n[6]}
	}
	return nil
}

func (db *ch05DB) update(stmt string) error {
	m := ch05UpdateRe.FindStringSubmatch(stmt)
	if m == nil {
		return fmt.Errorf("UPDATE without WHERE id, or not understood: %.70s", stmt)
	}
	row, ok := db.rows[m[2]]
	if !ok {
		return fmt.Errorf("UPDATE of absent id %s", m[2])
	}
	for _, g := range ch05GuardRe.FindAllStringSubmatch(m[3], -1) {
		want, _ := strconv.ParseInt(g[2], 10, 64)
		if p := row.col(g[1]); p == nil || *p != want {
			return fmt.Errorf("%s: guard %s = %s matches nothing, so the UPDATE would be a silent no-op", m[2], g[1], g[2])
		}
	}
	for _, a := range strings.Split(m[1], ", ") {
		g := ch05AssignRe.FindStringSubmatch(a)
		if g == nil || row.col(g[1]) == nil {
			return fmt.Errorf("%s: SET clause not understood: %q", m[2], a)
		}
		*row.col(g[1]), _ = strconv.ParseInt(g[2], 10, 64)
	}
	db.rows[m[2]], db.updated[m[2]] = row, true
	return nil
}

func (db *ch05DB) remove(stmt string) error {
	m := ch05DeleteRe.FindStringSubmatch(stmt)
	if m == nil {
		return fmt.Errorf("DELETE without WHERE id, or not understood: %.70s", stmt)
	}
	ids := []string{m[1]}
	if m[2] != "" {
		ids = ids[:0]
		for _, g := range ch05IDRe.FindAllStringSubmatch(m[2], -1) {
			ids = append(ids, g[1])
		}
	}
	for _, id := range ids {
		if _, ok := db.rows[id]; !ok {
			return fmt.Errorf("DELETE of absent id %s", id)
		}
		delete(db.rows, id)
		db.deleted = append(db.deleted, id)
	}
	return nil
}

func ch05Run(t *testing.T, db *ch05DB, section string) {
	t.Helper()
	for _, stmt := range ch05Statements(section) {
		if err := ch05Exec(db, stmt); err != nil {
			t.Errorf("%v", err)
		}
	}
}

// ch05State is 00010's Up (pre), 00112's Up on top of it (post), then 00112's
// Down on top of that (back).
type ch05State struct {
	raw, up, down   string
	pre, post, back *ch05DB
}

// ch05Read returns a migration's text and its two goose sections. A deleted or
// duplicated marker fails here: without one, goose would apply nothing or the
// Down as part of the Up.
func ch05Read(t *testing.T, path string) (raw, up, down string) {
	t.Helper()
	b, err := migrationsFS.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	raw = string(b)
	if !strings.HasPrefix(raw, "-- +goose Up") {
		t.Fatalf("%s must open with -- +goose Up", path)
	}
	if n := strings.Count(raw, "-- +goose Down"); n != 1 {
		t.Fatalf("%s holds %d -- +goose Down markers, want exactly 1", path, n)
	}
	up, down, _ = strings.Cut(raw, "-- +goose Down")
	return raw, up, down
}

func ch05Load(t *testing.T) ch05State {
	t.Helper()
	_, seedUp, _ := ch05Read(t, ch05Seed)
	s := ch05State{pre: ch05NewDB()}
	ch05Run(t, s.pre, seedUp)
	s.raw, s.up, s.down = ch05Read(t, ch05Migration)
	s.post = s.pre.clone()
	ch05Run(t, s.post, s.up)
	s.back = s.post.clone()
	ch05Run(t, s.back, s.down)
	return s
}

func ch05Economy(r ch05Row) profit.ConstantCapitalEconomy {
	return profit.ConstantCapitalEconomy{
		Kind:            profit.EconomyKind(r.kind),
		ConstantCapital: profit.ConstantCapital(r.c),
		VariableCapital: profit.VariableCapital(r.v),
		SurplusValue:    profit.SurplusValue(r.s),
		Saving:          profit.ConstantCapital(r.saving),
	}
}

func TestMigration00112_Ch5SeedFidelity(t *testing.T) {
	t.Parallel()

	t.Run("pre-image is 00010's three rows and post-image adds the fourth", func(t *testing.T) {
		t.Parallel()
		s := ch05Load(t)
		wantPre := []string{ch05ID(1), ch05ID(2), ch05ID(3)}
		if got := slices.Sorted(maps.Keys(s.pre.rows)); !slices.Equal(got, wantPre) {
			t.Errorf("00010 seeds %v, want %v", got, wantPre)
		}
		if got, want := slices.Sorted(maps.Keys(s.post.rows)), append(wantPre, ch05ID(4)); !slices.Equal(got, want) {
			t.Errorf("after 00112 the ids are %v, want %v", got, want)
		}
	})

	t.Run("stored rates are the domain's, before and after", func(t *testing.T) {
		t.Parallel()
		s := ch05Load(t)
		for label, db := range map[string]*ch05DB{"00010": s.pre, "00112": s.post} {
			for _, id := range slices.Sorted(maps.Keys(db.rows)) {
				r := db.rows[id]
				e := ch05Economy(r)
				if err := e.Validate(); err != nil {
					t.Errorf("%s %s: Validate() = %v", label, id, err)
				}
				a := profit.ComputeEconomyAnalysis(e)
				if a.OldProfitRate != r.oldBP || a.NewProfitRate != r.newBP || a.ProfitRateGain != r.gainBP {
					t.Errorf("%s %s: stored %d/%d/%d bp, domain gives %d/%d/%d", label, id,
						r.oldBP, r.newBP, r.gainBP, a.OldProfitRate, a.NewProfitRate, a.ProfitRateGain)
				}
			}
		}
	})

	t.Run("each restated row carries Marx's figure", func(t *testing.T) {
		t.Parallel()
		s := ch05Load(t)
		ch05CheckRags(t, s)
		ch05CheckShaft(t, s)
		ch05CheckPair(t, s)
		ch05CheckWorkday(t, s)
	})

	t.Run("Down restores 00010 exactly", func(t *testing.T) {
		t.Parallel()
		s := ch05Load(t)
		if !reflect.DeepEqual(s.back.rows, s.pre.rows) {
			t.Errorf("Down leaves\n %+v\nwant 00010's\n %+v", s.back.rows, s.pre.rows)
		}
	})

	t.Run("hygiene and provenance text", func(t *testing.T) {
		t.Parallel()
		ch05CheckHygiene(t, ch05Load(t))
	})
}

// ch05SameBase reports whether a restated row kept the base capital 00010 gave it.
func ch05SameBase(a, b ch05Row) bool { return a.c == b.c && a.v == b.v && a.s == b.s }

// Row 0501 is Marx's woollen rags: £44 a ton (M, Ch. 5 §IV), 10 tons reclaimed (I).
func ch05CheckRags(t *testing.T, s ch05State) {
	t.Helper()
	const tons, pricePerTon = 10, 44
	row := s.post.rows[ch05ID(1)]
	want := int64(profit.WasteReduction{
		Kind: profit.KindWasteReduction, OriginalWaste: tons, RecoveredWaste: tons, PricePerUnit: pricePerTon,
	}.ConstantCapitalSaving())
	if row.kind != string(profit.KindWasteReduction) || row.saving != want {
		t.Errorf("0501: kind %q saving %d, want waste_reduction and %d (%d tons at £%d)", row.kind, row.saving, want, tons, pricePerTon)
	}
	if !ch05SameBase(row, s.pre.rows[row.id]) {
		t.Errorf("0501: base capital changed from 00010's %+v", s.pre.rows[row.id])
	}
}

// Row 0502 is Marx's single shaft, which saves 60 per cent of the shafting (M, Ch. 5
// §III). The £500 other firms' shafting is I, the £200 of this mill's is D.
func ch05CheckShaft(t *testing.T, s ch05State) {
	t.Helper()
	const solo, shared, perCent = 500, 200, 60
	row := s.post.rows[ch05ID(2)]
	want := int64(profit.SharedCondition{
		Kind: profit.KindSharedConditions, SoloCost: solo, SharedCost: shared, ParticipantCount: 1,
	}.SavingPerParticipant())
	if row.kind != string(profit.KindSharedConditions) || row.saving != want {
		t.Errorf("0502: kind %q saving %d, want shared_conditions and %d", row.kind, row.saving, want)
	}
	if row.saving*100 != perCent*solo {
		t.Errorf("0502: saving %d of %d is not %d per cent", row.saving, int64(solo), int64(perCent))
	}
	if !ch05SameBase(row, s.pre.rows[row.id]) {
		t.Errorf("0502: base capital changed from 00010's %+v", s.pre.rows[row.id])
	}
}

// Row 0503 is left alone: its magnitudes are Ch. 7's OrganicCompositionEffect pair
// (M, Ch. 7), borrowed. The economy of 2,000 is the difference between the two c (D).
func ch05CheckPair(t *testing.T, s ch05State) {
	t.Helper()
	pair := profit.OrganicCompositionEffect{SA: 1000, VA: 1000, CA: 9000, SB: 1000, VB: 1000, CB: 11000}
	row := s.post.rows[ch05ID(3)]
	if !reflect.DeepEqual(row, s.pre.rows[row.id]) {
		t.Errorf("0503 must stay as 00010 left it: %+v", row)
	}
	if row.c != int64(pair.CB) || row.v != int64(pair.VB) || row.s != int64(pair.SB) || row.c-row.saving != int64(pair.CA) {
		t.Errorf("0503 %+v is not Ch. 7's pair c=%d/%d v=%d s=%d", row, pair.CA, pair.CB, pair.VB, pair.SB)
	}
}

// Row 0504 is the working-day, 12 against 18 hours (M, Ch. 5 §I). At 18 the charge
// per hour is two thirds of that at 12. The 3,600 of fixed capital is I.
func ch05CheckWorkday(t *testing.T, s ch05State) {
	t.Helper()
	const original, extended = 12, 18
	row := s.post.rows[ch05ID(4)]
	ext := profit.WorkdayExtensionEffect{
		Kind: profit.KindWorkdayExtension, FixedCapital: profit.ConstantCapital(row.c),
		OriginalHours: original, ExtendedHours: extended,
	}
	c12, c18 := ext.FixedChargePerHour(original), ext.FixedChargePerHour(extended)
	if 3*c18 != 2*c12 {
		t.Errorf("0504: charge %d at %dh is not two thirds of %d at %dh", c18, extended, c12, original)
	}
	if row.kind != string(profit.KindWorkdayExtension) || row.saving != original*int64(c12-c18) || row.saving*3 != row.c {
		t.Errorf("0504: kind %q saving %d, want workday_extension and %d = a third of c=%d", row.kind, row.saving, original*int64(c12-c18), row.c)
	}
	if row.v != row.s {
		t.Errorf("0504: v=%d and s=%d should be equal", row.v, row.s)
	}
}

func ch05CheckHygiene(t *testing.T, s ch05State) {
	t.Helper()
	for _, section := range []string{s.up, s.down} {
		for _, m := range ch05IDRe.FindAllStringSubmatch(section, -1) {
			if _, err := hex.DecodeString(m[1]); err != nil || len(m[1]) != 22 || !strings.HasPrefix(m[1], ch05IDPrefix) {
				t.Errorf("id %q is not a 22-char hex seed id with prefix %s", m[1], ch05IDPrefix)
			}
		}
	}
	for _, line := range strings.Split(s.raw, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--") && strings.Contains(line, ";") {
			t.Errorf("a semicolon in a comment line, which goose would split on: %q", line)
		}
	}
	// The header keeps the provenance a reader needs: where the pair really comes
	// from, Marx's two rag prices, the shaft, the tags and the money scale.
	for _, want := range []string{"Ch. 7", "£4 4s.", "£44", "29,000", "M = Marx", "one stored unit = £1"} {
		if !strings.Contains(s.raw, want) {
			t.Errorf("header lost %q", want)
		}
	}
	if got, want := slices.Sorted(slices.Values(s.back.deleted)), []string{ch05ID(4)}; !slices.Equal(got, want) {
		t.Errorf("Down deletes ids %v, want only %v", got, want)
	}
	if got, want := slices.Sorted(maps.Keys(s.post.updated)), []string{ch05ID(1), ch05ID(2)}; !slices.Equal(got, want) {
		t.Errorf("Up updates ids %v, want %v (0503 is Ch. 7's and stays)", got, want)
	}
	for _, c := range []struct {
		db        *ch05DB
		kw        string
		want      int
		direction string
	}{
		{s.post, "UPDATE", 2, "Up"}, {s.back, "UPDATE", 2, "Down"},
		{s.post, "INSERT", 1, "Up"}, {s.back, "INSERT", 0, "Down"}, {s.post, "DELETE", 0, "Up"},
	} {
		if got := c.db.counts[c.kw]; got != c.want {
			t.Errorf("%s runs %d %s statements, want %d", c.direction, got, c.kw, c.want)
		}
	}
}
