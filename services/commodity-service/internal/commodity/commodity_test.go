package commodity

import (
	"strings"
	"testing"
)

func makeLinen() Commodity {
	return Commodity{
		Name: "linen",
		UseValue: UseValue{
			Description: "linen for clothing",
			Unit:        "yards",
		},
		ConcreteLabour: ConcreteLabour{Kind: "weaving"},
		SNLTPerUnit:    30, // 30 minutes per yard
	}
}

func makeCoat() Commodity {
	return Commodity{
		Name: "coat",
		UseValue: UseValue{
			Description: "a coat to wear",
			Unit:        "coats",
		},
		ConcreteLabour: ConcreteLabour{Kind: "tailoring"},
		SNLTPerUnit:    600, // 10 hours per coat
	}
}

func TestCommodity_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		mut     func(*Commodity)
		wantErr string
	}{
		{name: "ok", mut: func(*Commodity) {}},
		{
			name:    "empty name",
			mut:     func(c *Commodity) { c.Name = "" },
			wantErr: "name is required",
		},
		{
			name:    "empty use-value description",
			mut:     func(c *Commodity) { c.UseValue.Description = "" },
			wantErr: "use_value.description is required",
		},
		{
			name:    "empty unit",
			mut:     func(c *Commodity) { c.UseValue.Unit = "" },
			wantErr: "use_value.unit is required",
		},
		{
			name:    "negative snlt",
			mut:     func(c *Commodity) { c.SNLTPerUnit = -1 },
			wantErr: "snlt_per_unit must be non-negative",
		},
		{
			name:    "missing concrete labour kind",
			mut:     func(c *Commodity) { c.ConcreteLabour.Kind = "" },
			wantErr: "concrete_labour.kind is required",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := makeLinen()
			tc.mut(&c)
			err := c.Validate()
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Validate() = %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestNewID_Unique(t *testing.T) {
	t.Parallel()
	seen := make(map[ID]struct{}, 1000)
	for i := 0; i < 1000; i++ {
		id := NewID()
		if id.IsZero() {
			t.Fatalf("NewID returned zero ID")
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("NewID produced a duplicate %q", id)
		}
		seen[id] = struct{}{}
	}
}

func TestLabourMinutes_String(t *testing.T) {
	t.Parallel()
	cases := map[LabourMinutes]string{
		0:   "0m",
		45:  "45m",
		60:  "1h",
		90:  "1h 30m",
		600: "10h",
	}
	for in, want := range cases {
		if got := in.String(); got != want {
			t.Errorf("LabourMinutes(%d).String() = %q, want %q", in, got, want)
		}
	}
}

func TestLabourMinutes_Hours(t *testing.T) {
	t.Parallel()
	if got := LabourMinutes(90).Hours(); got != 1.5 {
		t.Errorf("Hours() = %v, want 1.5", got)
	}
}
