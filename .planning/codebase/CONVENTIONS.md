# Coding Conventions

**Analysis Date:** 2026-05-23

## Go Import Grouping

Three groups, separated by blank lines. Enforced in every file.

```go
import (
    // 1. stdlib
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"

    // 2. third-party
    pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"

    // 3. local (module-internal)
    "github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)
```

Never merge groups. The blank lines are mandatory, not optional.

## File Naming

- Go source: `snake_case.go` — one dominant concept per file.
  - Examples: `cooperation.go`, `labour_power.go`, `time_wage.go`
- Every load-bearing concept from a chapter gets its own sibling file; do not append to a neighbor file.
  - `cooperation.go` not "appended to agent.go"
  - `production_account.go` not "appended to commodity.go"
- Test files: `<source_file>_test.go`, co-located with the source.

## Package Naming

- Domain packages match their directory: `package agent`, `package commodity`, `package simulation`.
- Store packages: `package store`.
- Transport packages: `package httpapi`.
- No suffix like `_service` or `_pkg`.

## Naming Conventions

**Types:**
- Exported struct names use `PascalCase`: `Cooperation`, `WorkingDay`, `ProductionAccount`.
- Named type aliases carry semantic meaning: `type LabourMinutes int64`, `type CooperationSize int`, `type Pence int64`.

**Functions:**
- Domain functions: `PascalCase` verbs — `CollectiveWorkingDay`, `AverageSocialLabour`, `SplitSurplus`.
- Methods on structs: `PascalCase` — `(c Commodity) Value(qty Quantity) LabourMinutes`.
- Constructors: `New<TypeName>ID()` for ID types — `NewCooperationID()`, `NewAgentID()`, `NewProductionAccountID()`.

**Constants and package-level vars:**
- Constants: `PascalCase` — `CooperationMinSize`, `CapitalistClass`, `CircuitMCM`.
- Sentinel errors: `Err<Domain><Condition>` — `ErrCooperationNoMembers`, `ErrInsufficientFunds`, `ErrNotFound`.

**JSON tags:**
- All struct fields use `snake_case` JSON tags, no exceptions.
- Example: `CapitalistID AgentID \`json:"capitalist_id"\``, `WorkingDayMinutes LabourMinutes \`json:"working_day_minutes"\``.
- No `camelCase`, no omission of tags on exported fields that appear in wire responses.

## ID Generation Pattern

IDs are 96-bit hex strings (24 hex characters) generated via `crypto/rand`. Every domain concept that persists to the store gets its own named ID type with an `IsZero()` method and a `NewXxxID()` constructor.

```go
// In services/agent-service/internal/agent/cooperation.go
type CooperationID string

func (id CooperationID) IsZero() bool { return id == "" }

func NewCooperationID() CooperationID {
    b := make([]byte, 12)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        panic(err)
    }
    return CooperationID(hex.EncodeToString(b))
}
```

**Forbidden:**
- `google/uuid` — not imported anywhere.
- `math/rand` for IDs or tokens.
- Inline `fmt.Sprintf("%x", ...)` — use `hex.EncodeToString`.

## LabourMinutes — Canonical Value-Magnitude Unit

`LabourMinutes int64` is the single unit for all socially necessary labour-time. Every field that represents a duration of labour, a value expressed in labour-time, or a capital magnitude in time-form must use this type.

```go
// services/agent-service/internal/agent/labour_power.go
type LabourMinutes int64

// services/commodity-service/internal/commodity/commodity.go
SNLTPerUnit LabourMinutes `json:"snlt_per_unit"`

// services/agent-service/internal/agent/cooperation.go
WorkingDayMinutes LabourMinutes `json:"working_day_minutes"`
```

**Never use `float64` or `int` for labour-time fields on persistent structs.** `float64` appears only in derived/computed values like `CollectiveProductivePower` factors or `RateOfSurplusValue`, never in stored field types.

## Store Interface Pattern

Every service has exactly three store files:

| File | Purpose |
|------|---------|
| `internal/store/store.go` | Interface definition + sentinel errors + partial-update types |
| `internal/store/memory.go` | In-memory implementation — used by unit tests |
| `internal/store/mysql.go` | Production MySQL implementation |

**`store.go` owns the contract:**
```go
// Sentinel errors — callers branch with errors.Is, never string comparison.
var (
    ErrNotFound      = errors.New("commodity: not found")
    ErrAlreadyExists = errors.New("commodity: already exists")
)

// Store interface — all persistence calls go through this seam.
type Store interface {
    Create(ctx context.Context, c commodity.Commodity) (commodity.Commodity, error)
    Get(ctx context.Context, id commodity.ID) (commodity.Commodity, error)
    List(ctx context.Context) ([]commodity.Commodity, error)
    Update(ctx context.Context, id commodity.ID, u Update) (commodity.Commodity, error)
    Delete(ctx context.Context, id commodity.ID) error
}
```

**`mysql.go` embeds migrations and runs them on construction:**
```go
//go:embed migrations
var migrationsFS embed.FS

func NewMySQL(ctx context.Context, db *sql.DB) (*MySQL, error) {
    sub, _ := fs.Sub(migrationsFS, "migrations")
    if err := pkgmysql.Migrate(ctx, db, sub); err != nil {
        return nil, err
    }
    return &MySQL{db: db, now: time.Now}, nil
}
```

**No ORM.** Use `database/sql` with positional `?` placeholders and `ExecContext`/`QueryContext` directly.

## Error Handling

**Sentinel errors in store layer:**
- Define as `var Err... = errors.New(...)` in `store.go`.
- Return them from `memory.go` and `mysql.go` for not-found and duplicate cases.
- MySQL duplicate detection uses a helper `isDuplicate(err)` that checks for MySQL error 1062.

**HTTP layer maps via `errors.Is`:**
```go
func writeAppError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, store.ErrNotFound):
        writeError(w, http.StatusNotFound, err.Error())
    case errors.Is(err, store.ErrAlreadyExists):
        writeError(w, http.StatusConflict, err.Error())
    case errors.Is(err, agent.ErrInsufficientFunds),
         errors.Is(err, agent.ErrCooperationNoMembers), ...:
        writeError(w, http.StatusBadRequest, err.Error())
    default:
        writeError(w, http.StatusInternalServerError, err.Error())
    }
}
```

**Domain validation errors** are returned from `Validate()` on struct methods. The handler calls `Validate()` before the store call and returns 400 on failure.

## HTTP Route Conventions

Routes use Go 1.22+ `net/http` mux syntax with the method prefix. Registered via `s.HandleFunc("METHOD /path", handler)` through `pkg/httpx.Server`.

```go
// services/agent-service/internal/transport/httpapi/routes.go
func Register(s *httpx.Server, h *Handler) {
    s.HandleFunc("POST /v1/cooperations", h.CreateCooperation)
    s.HandleFunc("GET /v1/cooperations", h.ListCooperations)
    s.HandleFunc("GET /v1/cooperations/{id}", h.GetCooperation)
    s.HandleFunc("POST /v1/cooperations/{id}/collective-working-day", h.ComputeCollectiveWorkingDay)
}
```

**Paths always start with `/v1/`.** Path parameters use `r.PathValue("id")` (Go 1.22+). No gorilla/mux, no chi.

**Handler helpers** (defined in `handler.go` for each service):
- `decodeJSON(r, &dst)` — decodes body with `DisallowUnknownFields()`, returns 400 error string on failure.
- `writeJSON(w, status, body)` — sets `Content-Type: application/json` and encodes.
- `writeError(w, status, msg)` — writes `{"error": "..."}`.
- `writeAppError(w, err)` — maps sentinel errors to HTTP status codes.

**List endpoints** always return a non-nil slice. Initialize with `make([]T, 0, ...)` so JSON encodes `[]` not `null`.

## Migration File Naming

```
NNNNN_v<V>_ch<NN>_<slug>.sql
NNNNN_v<V>_ch<NN>_seed.sql
```

Where:
- `NNNNN` is a five-digit zero-padded sequence number (globally unique per service).
- `V` is the volume number: `1`, `2`, or `3`.
- `NN` is a two-digit chapter number.
- `<slug>` is a short descriptive name using underscores.

Examples from `services/agent-service/internal/store/migrations/`:
```
00001_v1_ch04_agents.sql
00010_v1_ch13_cooperations.sql
00011_v1_ch13_seed.sql
00022_v1_ch22_seed_day_wages.sql
```

**Rules:**
- Never edit an existing migration file. Append a new numbered file.
- DDL only in schema migrations; never inline in Go.
- Every chapter that adds a domain type ships both a schema migration and a seed migration.

## Seed ID Pattern

Seed row IDs must start with `5eed00000000000000` followed by a chapter-specific suffix so they are visually distinct from `NewXxxID()` output and never collide with random IDs.

Format: `5eed00000000000000<CC><XX>` where `CC` = two-digit chapter number (hex or decimal), `XX` = row-discriminator.

```sql
-- Ch. 13 seeds: 5eed00000000000000 13 NN
'5eed00000000000000130001'  -- Spinner 01
'5eed00000000000000130100'  -- Richard Arkwright (capitalist)
'5eed00000000000000130c01'  -- Arkwright Spinning Floor (cooperation)
```

**`-- +goose Down` must DELETE only by known seed IDs**, never `DELETE FROM table` (which would wipe user data on migration rollback).

## Transport Layer Structure

Each service's transport layer is at `services/<svc>/internal/transport/httpapi/`:

| File | Purpose |
|------|---------|
| `handler.go` | `Handler` struct, `New(...)`, shared helpers (`decodeJSON`, `writeJSON`, `writeError`, `writeAppError`) |
| `routes.go` | `Register(s, h)` — wires all `HandleFunc` calls |
| `<concept>_handler.go` | Handler methods for one domain concept |

## React / TypeScript Conventions

**`web/src/types.ts`** mirrors every Go response struct. Field names use `snake_case` to match Go JSON tags exactly.

```typescript
// mirrors services/commodity-service/internal/commodity/commodity.go
export interface Commodity {
  id: string;
  name: string;
  use_value: UseValue;
  concrete_labour: ConcreteLabour;
  snlt_per_unit: number; // labour-minutes
  created_at: string;
  updated_at: string;
}
```

**`web/src/api.ts`** contains one exported function per endpoint.

**`web/src/chapters/registry.ts`** registers each chapter component with `volume`, `circuitNode[]`, `part`, and `status: "done"`.

**Chapter components** live at `web/src/chapters/vol<V>/Ch<NN><Title>.tsx`. Registered and imported in `web/src/App.tsx`.

**No `react-router`.** One `App.tsx` composes all panels. No client-side routing.

## Anti-Patterns — Never Do These

| Anti-Pattern | Why Forbidden |
|---|---|
| `interface{}` or `any` in domain struct fields | Destroys type safety; use concrete named types |
| ORM (gorm, ent, etc.) | Not in go.mod; use `database/sql` directly |
| Per-service `go.mod` or `go.work` | Single module root is intentional |
| Editing an existing migration file | Migrations are append-only |
| DDL inline in Go (`CREATE TABLE` in `.go` files) | All DDL belongs in `.sql` files under `migrations/` |
| `google/uuid` for IDs | Use `crypto/rand` hex via `New<Type>ID()` |
| `float64` for stored labour-time fields | Use `LabourMinutes int64` |
| `react-router` | Not installed; one App.tsx |
| HTTPS termination in services | The cluster handles TLS at the edge |
| `math/rand` for IDs | Not cryptographically safe; `crypto/rand` only |
| Referencing migration files by number in Go | Use `//go:embed migrations` in the store package |

---

*Convention analysis: 2026-05-23*
