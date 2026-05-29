---
name: store-scaffold
description: Generate the three-file persistence layer (store.go interface, memory.go in-memory impl, mongo.go MongoDB impl) for a new domain type in a capital-simulator service. Produces the sentinel-error pattern, partial-update Update struct, Memory store with case-insensitive uniqueness on the chosen key, and Mongo store with the matching unique index. Use whenever the user asks to scaffold a store, add persistence for a new type, set up a new internal/store/ package, or says "I need a Store for <T>". Especially relevant when a chapter introduces a new service (agent-service, market-service, simulation-engine) that needs its first domain type persisted.
---

# store-scaffold

The repo has one persistence pattern, repeated per domain type:

```
services/<svc>/internal/store/
├── store.go    Store interface + ErrNotFound/ErrAlreadyExists + Update struct
├── memory.go   in-memory impl (used by tests; fallback in dev)
└── mongo.go    MongoDB impl with unique index on the natural key
```

It's the same three files every time, with the same shape. The differences
are the domain type, the field set, and the natural-key column. This skill
turns "scaffold a Store for `agent.Worker`" into a correct first cut of all
three files.

## When the user invokes this

- "Scaffold a store for `Worker` in agent-service"
- "Add persistence for the new `Order` type"
- "I need a Store for X"
- A chapter introduces a new domain type and the user opens a service's
  `internal/store/` directory expecting it to be empty.

If the service has an existing `store/` package, append to it rather than
overwriting - re-running this for a second domain type in the same service
should add to the same `store.go`, not replace it.

## Inputs to gather

Before generating, get from the user (ask once if missing):

1. **Service name.** `agent-service`, `market-service`, etc. The Go path
   `services/<svc>/internal/store/` must exist.
2. **Domain package.** Usually mirrors the service - `agent` for
   agent-service. The package import path is
   `github.com/theding0x/capital-simulator/services/<svc>/internal/<pkg>`.
3. **Domain type name.** The Go type, exported, e.g. `Worker` or `Order`.
   Type must already be defined (with `Validate() error` and an `ID` type
   following `commodity.NewID()` style).
4. **Natural key field.** The field that gives "one canonical record per
   X" - `Name` for commodity, `Symbol` for some market type, etc. Must be
   a string field on the type. If there isn't one, the user should say so
   and we skip uniqueness handling (rare).
5. **Patchable fields.** Which fields the PATCH endpoint should be able
   to update. Typically every domain field except `ID`, `CreatedAt`,
   `UpdatedAt`. Each becomes a `*T` field on the `Update` struct.

## Read these first (templates lag)

The canonical reference is the existing commodity store. Read it before
generating - it always represents the current pattern, while the
templates in this skill might drift:

- `services/commodity-service/internal/store/store.go`
- `services/commodity-service/internal/store/memory.go`
- `services/commodity-service/internal/store/mongo.go`
- `services/commodity-service/internal/store/memory_test.go` (for the
  test pattern with Marx fixtures)

If anything in the templates here disagrees with the current commodity
store, prefer the commodity store and update the templates.

## Generation

The templates live in this skill at `references/`:

- `references/store.go.tmpl`
- `references/memory.go.tmpl`
- `references/mongo.go.tmpl`

Each uses `<<...>>` placeholders. Don't sed them blindly - the
substitutions are field-aware (e.g. the `Update.Apply` method needs one
`if u.X != nil { out.X = *u.X }` block per patchable field). Read the
template, fill it in for the specific type, write the file.

Required substitutions:

| Placeholder        | Example value           | Notes |
|--------------------|-------------------------|-------|
| `<<svc>>`          | `agent-service`         | Path component |
| `<<pkg>>`          | `agent`                 | Domain package |
| `<<Type>>`         | `Worker`                | Domain type name (PascalCase) |
| `<<type>>`         | `worker`                | Lowercase form for error strings |
| `<<types>>`        | `workers`               | Plural lowercase, for error strings + collection name |
| `<<key>>`          | `Name`                  | Natural-key Go field |
| `<<key_lower>>`    | `name`                  | bson/json tag of the key |
| `<<UpdateFields>>` | (see template)          | One `*T` field per patchable field |

For `<<UpdateFields>>` and the corresponding `Apply` / `IsEmpty` /
`Mongo.$set` blocks, generate them per the type's actual fields. Each
patchable field follows the pattern below; a nested struct field
(e.g. `UseValue.Description`) becomes a flattened `*string` named after
the dotted path: `UseValueDescription *string`.

```go
// Update struct field
<<FieldName>> *<<FieldType>>

// Apply branch
if u.<<FieldName>> != nil {
    out.<<dotted.path>> = *u.<<FieldName>>
}

// Mongo $set entry
if u.<<FieldName>> != nil {
    set["<<bson.path>>"] = next.<<dotted.path>>
}
```

## After generation

1. Add a `memory_test.go` next to `memory.go` that exercises Create /
   Get / List / Update / Delete and conflict on duplicate key. Use one
   of Marx's textual examples for the fixture (e.g. for a `Worker`,
   "weaver of 20 yards of linen"). Always `t.Parallel()`.
2. Run `go mod tidy` if the service is wiring in the mongo driver for
   the first time (pull `go.mongodb.org/mongo-driver/...`).
3. Ask the user to run `make vet test build` - the sandbox can't run
   the Go toolchain (CLAUDE.md sandbox limits).

## What this skill does NOT do

- **Define the domain type.** `Worker` / `Order` / etc. is the chapter's
  argument and must already exist. If it doesn't, ask the user to write
  it first.
- **Wire the store into main.go or HTTP routes.** That's a follow-up
  step (mirroring `services/commodity-service/cmd/commodity-service/main.go`).
- **Use an ORM.** Mongo driver direct, per CLAUDE.md anti-patterns.

## Anti-patterns

- Don't add a new ID type if `<<pkg>>.NewID()` already exists. The
  pattern is one ID type per domain package, hex-encoded 96 bits from
  crypto/rand.
- Don't import `github.com/google/uuid`. Repo convention is the local
  `NewID()`.
- Don't return raw mongo errors to callers. Map duplicate-key →
  `ErrAlreadyExists` and `ErrNoDocuments` → `ErrNotFound`.
- Don't add a transaction layer. Single-document operations only at
  this stage; the simulation isn't multi-doc consistent yet.
