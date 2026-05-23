# Testing Patterns

**Analysis Date:** 2026-05-23

## Test Framework

**Runner:**
- Go's standard `testing` package — no external test runner.
- No testify, gomock, or assertion libraries. Assertions are hand-written with `t.Fatalf`/`t.Errorf`.

**Run Commands:**
```bash
make test             # go test ./... — all packages
make vet              # go vet ./... — runs before test in CI
make vet test build   # full local check
```

No coverage flag is enforced in the Makefile. Run manually with:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test File Organization

**Location:** Co-located with source in the same directory.
- `services/agent-service/internal/agent/cooperation.go` → `cooperation_test.go` in the same directory.
- `services/agent-service/internal/store/memory.go` → `memory_test.go` in the same directory.
- `services/agent-service/internal/transport/httpapi/cooperation_handler.go` → `cooperation_handler_test.go` in the same directory.

**Naming:**
- `<concept>_test.go` for domain tests.
- `<concept>_handler_test.go` for HTTP handler tests.
- `memory_test.go` for store memory tests (sometimes split: `labour_power_memory_test.go`, `working_day_memory_test.go`).

**Directory layout:**
```
services/<svc>/internal/
├── <domain>/
│   ├── cooperation.go
│   ├── cooperation_test.go   ← domain unit tests
│   └── ...
├── store/
│   ├── memory.go
│   ├── memory_test.go        ← store unit tests (Memory only)
│   └── ...
└── transport/httpapi/
    ├── cooperation_handler.go
    ├── cooperation_handler_test.go  ← HTTP integration tests
    └── ...
```

## Package Naming in Tests

Two styles are used — both valid. The choice is per-file:

**Black-box tests** (`package <pkg>_test`) — most common, especially for domain logic tests accessed via the public API:
```go
// services/agent-service/internal/agent/cooperation_test.go
package agent_test

import "github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
```

**White-box tests** (`package <pkg>`) — used when the test needs access to unexported identifiers, or when the test is for HTTP handler internals:
```go
// services/agent-service/internal/transport/httpapi/cooperation_handler_test.go
package httpapi   // same package — can access unexported response types

// services/commodity-service/internal/commodity/commodity_test.go
package commodity // helper functions like makeLinen() are package-private
```

## t.Parallel() — Mandatory

**Every test function and every subtest must call `t.Parallel()` as its first statement.** All 74 test files in the repo follow this rule (714 `t.Parallel()` calls total).

```go
func TestCollectiveWorkingDay_Para2_1200Workers(t *testing.T) {
    t.Parallel()
    got := agent.CollectiveWorkingDay(1200, agent.LabourMinutes(720))
    if got != agent.LabourMinutes(864_000) {
        t.Fatalf("CollectiveWorkingDay(1200, 720) = %d, want 864000", got)
    }
}
```

For table-driven tests, both the outer test and each subtest call `t.Parallel()`:
```go
func TestCommodity_Validate(t *testing.T) {
    t.Parallel()
    cases := []struct{ ... }{ ... }
    for _, tc := range cases {
        tc := tc  // capture loop variable (pre-Go 1.22)
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            ...
        })
    }
}
```

## Marx Textual Examples as Fixtures

Tests use names and magnitudes drawn directly from *Capital* — the same examples Marx uses in the text. Test comments cite the chapter, section, and paragraph/footnote.

**Canonical fixtures used across the codebase:**

| Fixture | Value | Source |
|---------|-------|--------|
| 20 yards linen = 1 coat | `SNLTPerUnit: 30` (linen), `SNLTPerUnit: 600` (coat) | Vol. I, Ch. 1, §1 |
| 12 men × 12 h = 144 h collective working-day | 12 workers × 720 min = 8640 min | Vol. I, Ch. 13, §3 |
| Burke's five-man platoon | `CooperationMinSize = 5` | Vol. I, Ch. 13, §3 fn. 1 |
| Spinning mill 100% surplus rate | `c=1440, v=360, s=360` | Vol. I, Ch. 7, §2 |
| Spinning mill 1871 (153 11/13%) | `v=124800, s=192000` | Vol. I, Ch. 9, §1 |
| Abraham sequence (full accumulation) | `c=8000, v=2000, s'=100%` | Vol. I, Ch. 24 |
| Working-day = 720 min (12 h) | `LabourMinutes(720)` | canonical throughout |

Test function names embed the Marx reference:
```go
func TestCollectiveWorkingDay_Para2_1200Workers(t *testing.T)
func TestAverageSocialLabour_BurkeFootnote_FiveMen(t *testing.T)
func TestValue_MarxLinenAndCoat(t *testing.T)
func TestComputeRate_SpinningMill1871(t *testing.T)
func TestRunExtendedReproduction_AbrahamSequence(t *testing.T)
```

Test function comments cite the original text with the exact quotation when the test encodes a specific textual claim:
```go
// Capital Vol. I, Ch. 13, § (para 2):
// "If a working-day of 12 hours be embodied in six shillings,
//  1,200 such days will be embodied in 1,200 times 6 shillings."
func TestCollectiveWorkingDay_Para2_1200Workers(t *testing.T) {
```

## Helper / Fixture Functions

Package-level helper functions create reusable Marx fixtures. These live at the top of the `_test.go` file (not in a separate `testdata` file) and are declared in white-box style so they are accessible within the package.

```go
// services/commodity-service/internal/commodity/commodity_test.go (package commodity)
func makeLinen() Commodity {
    return Commodity{
        Name: "linen",
        UseValue: UseValue{Description: "linen for clothing", Unit: "yards"},
        ConcreteLabour: ConcreteLabour{Kind: "weaving"},
        SNLTPerUnit: 30, // 30 minutes per yard
    }
}

func makeCoat() Commodity {
    return Commodity{
        Name: "coat",
        UseValue: UseValue{Description: "a coat to wear", Unit: "coats"},
        ConcreteLabour: ConcreteLabour{Kind: "tailoring"},
        SNLTPerUnit: 600, // 10 hours per coat
    }
}
```

Handler tests use a `newTestServer` / `newAgentTestServer` helper that wires the in-memory store into a `httptest.NewServer`:
```go
func newAgentTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
    t.Helper()
    st := store.NewMemory()
    h := New(st, slog.Default())
    mux := http.NewServeMux()
    mux.HandleFunc("POST /v1/cooperations", h.CreateCooperation)
    // ... all routes
    ts := httptest.NewServer(mux)
    t.Cleanup(ts.Close)
    return ts, st
}
```

## Memory Store vs MySQL Store in Tests

**Only `Memory` is tested in unit tests.** `MySQL` is exercised by CI integration tests against a real database.

```go
// Correct: use store.NewMemory() in tests
func TestMemory_CreateGet(t *testing.T) {
    t.Parallel()
    m := store.NewMemory()
    ctx := context.Background()
    created, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000))
    ...
}
```

**Standard store memory test coverage:**
1. Create → Get round-trip: returned struct is identical to created struct.
2. Get missing ID → `errors.Is(err, store.ErrNotFound)`.
3. Create duplicate → `errors.Is(err, store.ErrAlreadyExists)` (where the interface supports it).
4. Update: fields change, unmodified fields unchanged.
5. Delete: subsequent Get returns `ErrNotFound`.
6. List: returns a non-nil slice (possibly empty).

```go
func TestMemory_Get_NotFound(t *testing.T) {
    t.Parallel()
    m := store.NewMemory()
    _, err := m.Get(context.Background(), agent.NewID())
    if !errors.Is(err, store.ErrNotFound) {
        t.Errorf("want ErrNotFound, got %v", err)
    }
}
```

## HTTP Handler Test Patterns

Two styles are used:

**Style 1 — Full httptest.Server (integration-style):** Used when the test exercises routing or needs realistic URL construction. The `newTestServer` helper spins up a real HTTP server.
```go
func TestCreateCooperation_BurkePlatoon(t *testing.T) {
    t.Parallel()
    ts, st := newAgentTestServer(t)
    cap := seedCapitalist(t, st, "Robert Owen")
    body := `{"name":"spinning floor", ...}`
    res, err := http.Post(ts.URL+"/v1/cooperations", "application/json", strings.NewReader(body))
    if err != nil {
        t.Fatalf("POST: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusCreated {
        t.Fatalf("status = %d, want 201", res.StatusCode)
    }
    var got cooperationResponse
    json.NewDecoder(res.Body).Decode(&got)
    // assert fields...
}
```

**Style 2 — httptest.NewRecorder (unit-style):** Used for handlers that don't require a full server, typically pure-computation endpoints.
```go
func TestCreateWageForm(t *testing.T) {
    t.Parallel()
    h := httpapi.New(store.NewMemory(), nil)
    b, _ := json.Marshal(map[string]any{"agent_id": "abc123", ...})
    req := httptest.NewRequest(http.MethodPost, "/v1/wage-forms", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    h.CreateWageForm(rr, req)
    if rr.Code != http.StatusCreated {
        t.Fatalf("status = %d, want %d; body: %s", rr.Code, http.StatusCreated, rr.Body.String())
    }
}
```

**What handler tests always verify:**
- `201 Created` on successful POST creation.
- `200 OK` on successful GET/computation.
- `400 Bad Request` on malformed JSON body or invalid domain input.
- `404 Not Found` when the store returns `ErrNotFound` (via `errors.Is`).
- `409 Conflict` when the store returns `ErrAlreadyExists`.
- Response body decodes without error into the expected response type.
- Specific numeric fields match Marx's examples (not just "non-zero").

## What Is NOT Tested

- **`MySQL` store implementation** — tested only by CI against a live database, not in unit tests.
- **`pkg/httpx.Server`** shutdown/signal handling — tested at integration level.
- **Migration SQL correctness** — verified by applying migrations against the real DB in CI.
- **React components** — checked with Playwright, not Go tests.
- **The `//go:embed` wiring** — exercised only when `NewMySQL` runs.

## Domain Logic Test Patterns

**Pure function tests** verify happy-path, boundary, and invariant cases:

```go
// Happy-path: known Marx fixture
func TestCollectiveWorkingDay_Para7_100Men(t *testing.T) {
    t.Parallel()
    got := agent.CollectiveWorkingDay(100, agent.LabourMinutes(720))
    if got != agent.LabourMinutes(72_000) {
        t.Fatalf("CollectiveWorkingDay(100, 720) = %d, want 72000", got)
    }
}

// Boundary: zero inputs
func TestCollectiveWorkingDay_Zero(t *testing.T) {
    t.Parallel()
    if got := agent.CollectiveWorkingDay(0, 720); got != 0 { ... }
}

// Invariant: mathematical law holds for all inputs in the set
func TestInvariant_AverageSocialLabourReversesCollective(t *testing.T) {
    t.Parallel()
    for _, n := range []agent.CooperationSize{1, 5, 12, 100, 1200} {
        d := agent.LabourMinutes(720)
        collective := agent.CollectiveWorkingDay(n, d)
        got := agent.AverageSocialLabour(collective, n)
        if got != d {
            t.Fatalf(...)
        }
    }
}
```

**Validate() tests** use table-driven cases, one case per validation rule:
```go
func TestCommodity_Validate(t *testing.T) {
    t.Parallel()
    cases := []struct {
        name    string
        mut     func(*Commodity)
        wantErr string
    }{
        {name: "ok", mut: func(*Commodity) {}},
        {name: "empty name", mut: func(c *Commodity) { c.Name = "" }, wantErr: "name is required"},
        {name: "negative snlt", mut: func(c *Commodity) { c.SNLTPerUnit = -1 }, wantErr: "snlt_per_unit must be non-negative"},
    }
    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            c := makeLinen()
            tc.mut(&c)
            err := c.Validate()
            if tc.wantErr == "" {
                if err != nil { t.Fatalf("Validate() = %v, want nil", err) }
                return
            }
            if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
                t.Fatalf("Validate() = %v, want error containing %q", err, tc.wantErr)
            }
        })
    }
}
```

## Seed Migration Tests

Verify that seed IDs follow the `5eed00000000000000<CC>*` prefix and are unique within their file. This is typically a review-time check rather than an automated test — the CI applies migrations against MySQL and catches duplicate-key violations.

---

*Testing analysis: 2026-05-23*
