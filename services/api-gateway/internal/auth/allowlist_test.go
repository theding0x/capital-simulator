package auth

import "testing"

func TestIsComputePath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		path string
		want bool
	}{
		{"/v1/observatory/levers", true},
		{"/v1/observatory/levers/extra", true},
		{"/v1/commodities", false},
		{"/v1/owners", false},
		{"/v1/observatory/snapshot", false}, // GET-only path, not a write allowlist entry
	}
	for _, c := range cases {
		if got := IsComputePath(c.path); got != c.want {
			t.Errorf("IsComputePath(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}
