package auth

import "strings"

// computePaths are write-method endpoints that perform NO database mutation —
// pure calculators and the per-session in-memory Atlas observatory. Non-owners
// (guests) may call these so they can play with the simulation.
//
// FAIL-CLOSED: anything not listed here is owner-only for write methods. The
// only failure mode is "a forgotten compute panel degrades to owner-only for
// guests", never "a forgotten write leaks to guests". Add a prefix here only
// after confirming the handler does not call a store mutation (see Task 8).
var computePaths = []string{
	"/v1/observatory/levers", // Atlas observatory — per-session in-memory, zero DB writes
}

// IsComputePath reports whether path is a non-persisting compute endpoint that
// non-owners may invoke with a write method.
func IsComputePath(path string) bool {
	for _, p := range computePaths {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}
