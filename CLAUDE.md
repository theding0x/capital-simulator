# CLAUDE.md

Project memory for Claude. Read this once per session; do not re-derive.

## What this is

Microservice simulation of an economy modeled chapter-by-chapter on Marx's
*Capital, Vol. I*. One branch + one PR per chapter. Chapter source text and
spec live in the **red-vault** Obsidian vault at
`marx-engels/1867/capital-volume-i/{texts,specs}/` and are read via the
`obsidian` MCP server, not from this repo.

## Stack (fixed; do not propose alternatives)

Go 1.25 monorepo · React 18 + Vite + TS · MySQL 8 · Redis · Docker · k8s
(kustomize) · GitHub. Module path: `github.com/theding0x/capital-simulator`.

## Layout

```
pkg/{log,httpx,mysql}/          shared Go libs
services/<svc>/cmd/<svc>/       main.go (one per service)
services/<svc>/internal/        domain, store, transport
services/<svc>/Dockerfile       multi-stage, distroless, build context = repo root
web/                            Vite+React+TS, nginx prod image, /api proxy
deploy/k8s/                     kustomize: namespace + infra + per-service
docs/architecture.md            authoritative topology + chapter status
```

Chapter source text and specs are **not in this repo** — they live in the
red-vault Obsidian vault, accessed through the `obsidian` MCP server.

## Services & ports

| svc                | port | role                                                              |
|--------------------|------|-------------------------------------------------------------------|
| api-gateway        | 8080 | external entrypoint; reverse-proxies to domain services           |
| commodity-service  | 8081 | Ch. 1 - commodity, value, value-forms, fetishism                  |
| agent-service      | 8082 | (placeholder until Ch. 4)                                         |
| market-service     | 8083 | (placeholder until Ch. 2-3)                                       |
| simulation-engine  | 8084 | (placeholder)                                                     |

## Conventions (don't re-explain to the user)

- **Go imports**: stdlib group, blank line, third-party, blank line, local.
- **HTTP**: pkg/httpx.Server. Routes use Go 1.22+ mux syntax (`POST /v1/...`).
- **Persistence**: every domain service gets `internal/store/{store.go,memory.go,mysql.go}`. Store interface, Memory for tests, MySQL for prod.
- **Migrations**: `github.com/pressly/goose/v3`. SQL files live at `internal/store/migrations/NNNNN_chNN_<slug>.sql`, embedded via `//go:embed` and applied by `pkg/mysql.Migrate` inside `NewMySQL`. Add a new numbered file per chapter; never edit existing ones.
- **Seeds**: every chapter that adds a domain type also ships a `NNNNN_chNN_seed.sql` migration that inserts Marx-faithful exemplars (named after his actual fixtures), with a `-- +goose Down` that DELETEs every seeded id. Seed IDs follow `5eed00000000000000<CC><…>` so they're recognisable on sight and never collide with `commodity.NewID()`. The dashboard must come up populated on a fresh MySQL volume — empty panels are a regression.
- **Errors**: sentinel `ErrNotFound`/`ErrAlreadyExists` in store; HTTP layer maps via `errors.Is`.
- **Time/labour**: `LabourMinutes int64` is the canonical value-magnitude unit.
- **IDs**: `commodity.NewID()` style — 96-bit hex from crypto/rand. No google/uuid.
- **Tests**: `t.Parallel()` by default. Use Marx's textual examples as test fixtures (20 yards linen = 1 coat, etc.).
- **Naming**: JSON tags use `snake_case`. Match the wire format the React types expect.
- **Web**: no router yet. One `App.tsx` composes feature panels. Shared `types.ts` mirrors Go structs.

## Commands

```bash
go mod tidy                 # after touching go.mod
make vet test build         # full Go check
make run-<service>          # run one service from source
docker compose up --build   # full stack locally
cd web && npm run lint      # tsc typecheck (no eslint)
cd web && npm run build     # vite production build
```

## Chapter workflow (do this; don't re-discuss)

1. Create a new branch off `main` named `volume-X/chapter-Y` (X = volume number, Y = chapter number — no slug suffix).
2. Open a **draft** PR against `main` populated from `.github/pull_request_template.md`, with the *planned* changes drawn from the chapter spec at `marx-engels/1867/capital-volume-i/specs/NN-<slug>.spec.md` in the red-vault Obsidian vault, fetched via the `obsidian` MCP server (Chapter, Summary, Services touched, planned tests, Notes for review).
3. Implement the planned chapter changes (domain types in the relevant service(s), tests using Marx's examples, seed migration `NNNNN_chNN_seed.sql` for every new domain type, `docs/architecture.md` roadmap row flipped to Done).
4. Commit signed (`commit.gpgsign=true` is set per-repo). Use a multi-line conventional commit; the PR description fills from it.
5. Compare the committed changes against the draft PR and update the PR description so it reflects what actually landed.
6. Wait for GitHub Actions to finish running checks on the PR.
7. If checks fail, fetch the output (`gh run view --log-failed`), fix the failure, and commit again. Loop to step 5.
8. If all checks pass, mark the PR ready for review and notify the user that the PR is ready to merge into `main`.

> **Note:** If the user says "implement the next pending chapter", look up the next `Pending` row in `docs/architecture.md` (Roadmap table) — that's the authoritative source for chapter ordering and primary services.

## Chapter workflow (swarm) — replaces steps 3–7 above when running under RuFlo

Use this section instead of manually executing steps 3–7 when a RuFlo swarm is available. Steps 1–2 (branch + draft PR) and step 8 (mark ready) remain manual.

### 0. Initialise the swarm (one-off per chapter)

```bash
npx @claude-flow/cli@latest swarm init \
  --topology hierarchical \
  --max-agents 8 \
  --strategy specialized
```

Then spawn the five agents below as background tasks via the Claude Code Task tool. Each agent description is its complete brief — paste it verbatim as the task prompt.

### Agent 1 — architect

**Role:** Read the chapter spec and existing code, then publish a design contract that all other agents consume. Runs first; no other agent starts until it stores its output.

**Steps:**

1. Fetch the chapter spec from the red-vault Obsidian vault:
   `marx-engels/1867/capital-volume-i/specs/NN-<slug>.spec.md`
   using `mcp__obsidian__obsidian_get_file_contents`. Also fetch the corresponding source text at
   `marx-engels/1867/capital-volume-i/texts/NN-<slug>.md` for fixture names.
2. Identify the primary service(s) listed in `docs/architecture.md` for this chapter.
3. For each new domain concept in the spec, decide:
   - Which package it belongs to (e.g. `services/agent-service/internal/agent/`).
   - Whether it is a pure function, a persistent entity, or both.
   - The Go type name, field names (concrete types only — no `interface{}` or `any`), and `snake_case` JSON tags.
   - The canonical `LabourMinutes int64` unit for any value-magnitude field.
   - The ID constructor name (e.g. `NewCooperationID()`) and pattern (96-bit hex via `crypto/rand`).
4. List every new HTTP endpoint: method, path (`/v1/...`), request/response shape, and which service handles it.
5. List every new migration file needed: `NNNNN_chNN_<slug>.sql` and `NNNNN_chNN_seed.sql`, with the seed ID prefix `5eed00000000000000<CC><XX>` where `CC` = two-digit chapter number.
6. Identify any new api-gateway proxy rules (`services/api-gateway/internal/routes/routes.go`).
7. Identify new React types needed in `web/src/types.ts` and new API functions in `web/src/api.ts`.
8. Store the full design contract:
   ```bash
   npx @claude-flow/cli@latest memory store \
     --key "design-chNN" \
     --value "<JSON contract>" \
     --namespace tasks
   ```
9. Post completion:
   ```bash
   npx @claude-flow/cli@latest hooks post-task --task-id "architect-chNN" --success true
   ```

### Agent 2 — coder

**Role:** Implement all Go and React files specified in the architect's design contract. Does not write tests — the tester agent owns test files.

**Steps:**

1. Retrieve the design contract:
   ```bash
   npx @claude-flow/cli@latest memory search --query "design-chNN" --namespace tasks
   ```
2. For each new domain concept, create `services/<svc>/internal/<pkg>/<concept>.go` containing:
   - Exported struct, ID type, `New<Concept>ID()`, `Validate()` if the spec defines invariants.
   - Pure domain functions (no DB calls, no HTTP).
   - No `interface{}` or `any`. No inline SQL or DDL.
   - Go import groups: stdlib / blank line / third-party / blank line / local.
3. Extend `services/<svc>/internal/store/store.go` with new Store interface methods.
4. Add implementations to `services/<svc>/internal/store/memory.go` (in-memory, for tests) and `services/<svc>/internal/store/mysql.go` (production). Use `database/sql` directly — no ORM.
5. Create migration files in `services/<svc>/internal/store/migrations/`:
   - Schema DDL: `NNNNN_chNN_<slug>.sql` — all DDL belongs here, never inline in Go.
   - Seed: `NNNNN_chNN_seed.sql` — Marx-faithful exemplars with seed IDs `5eed00000000000000<CC><XX>`. Include `-- +goose Down` that DELETEs every seeded row by ID.
6. Create `services/<svc>/internal/transport/httpapi/<concept>_handler.go` with handler functions. Wire into `services/<svc>/internal/transport/httpapi/routes.go` using Go 1.22+ mux syntax (`mux.HandleFunc("POST /v1/...", h.handleX)`). Register the handler in `handler.go`'s `New` function.
7. Add reverse-proxy rules to `services/api-gateway/internal/routes/routes.go` for every new path prefix.
8. Update `web/src/types.ts` with TypeScript mirror types (snake_case field names to match JSON tags).
9. Add API functions to `web/src/api.ts` for each new endpoint.
10. Create `web/src/chapters/ChNN<Title>.tsx` (and `.css` if needed) for the chapter panel. Register it in `web/src/chapters/registry.ts` with `status: "done"`. Import it in `web/src/App.tsx`.
11. Update the `docs/architecture.md` roadmap row to `Done` and add a `### Ch. NN — what was built` section.
12. Post completion:
    ```bash
    npx @claude-flow/cli@latest hooks post-task --task-id "coder-chNN" --success true
    ```

### Agent 3 — tester

**Role:** Write all `_test.go` files and verify that every invariant from the spec has test coverage. Runs after the coder finishes (depends on domain files existing).

**Steps:**

1. Retrieve the design contract from memory (same key as coder).
2. For each new domain file `<concept>.go`, create `<concept>_test.go` in the same package. Every test function must call `t.Parallel()` as its first statement.
3. Use Marx's textual examples as fixtures — names and magnitudes drawn from the chapter source text, not invented. Canonical examples: "20 yards linen = 1 coat", "Burke's five-man platoon", "Caslon type-foundry 4/2/1", "Spinning Mill 1871", etc. Read the chapter text from the vault if in doubt.
4. For each pure function, verify:
   - Happy-path: known Marx fixture produces the documented output.
   - Boundary: zero-value, minimum, and at-capacity inputs.
   - Invariants: e.g., partition equations hold (`Total == NecessaryLabour + SurplusLabour`), value conservation, `SurplusValue >= 0`.
5. For the store, test `Memory` only (MySQL is integration-tested via CI). Verify:
   - `ErrNotFound` returned for a missing ID.
   - `ErrAlreadyExists` returned on duplicate creation where applicable.
   - Round-trip: create → get returns identical struct.
6. For each seed migration, verify that seed IDs follow the pattern `5eed00000000000000<CC>*` and are unique within the file.
7. For the HTTP handler, write a handler test in `services/<svc>/internal/transport/httpapi/<concept>_handler_test.go`. Use `httptest.NewRecorder` and the in-memory store. Verify:
   - `201` on successful creation with correct `Location` header where applicable.
   - `404` from `errors.Is(err, store.ErrNotFound)` mapping.
   - `400` on malformed JSON body.
   - `200` list endpoint returns a non-nil (possibly empty) slice, never `null`.
8. Post completion:
   ```bash
   npx @claude-flow/cli@latest hooks post-task --task-id "tester-chNN" --success true
   ```

### Agent 4 — reviewer

**Role:** Read all files produced by the coder and tester and flag any violation of the project conventions in CLAUDE.md. Does not fix — reports findings as a structured list stored in memory.

**Steps:**

1. Read every new and modified file in the chapter branch diff.
2. Check Go conventions:
   - Import groups: stdlib / blank line / third-party / blank line / local. No merged groups.
   - JSON tags: all `snake_case`. No `camelCase` tags.
   - HTTP routes: Go 1.22+ mux syntax (`"METHOD /v1/path"`). No `http.HandleFunc("/path", ...)` without method prefix.
   - IDs: `<Concept>NewID()` using `crypto/rand`. No `google/uuid`, no `math/rand`.
   - `LabourMinutes int64` for all value-magnitude fields — not `float64`, not `int`.
   - Errors: store sentinel `ErrNotFound` / `ErrAlreadyExists` defined in `store.go`. HTTP handler maps via `errors.Is`.
3. Check persistence conventions:
   - Store interface defined in `store.go`. Memory and MySQL implementations in `memory.go` and `mysql.go`.
   - No new `go.mod` files. Single module root.
   - No edited existing migration files — only new numbered files added.
   - No DDL in Go source — all DDL in `.sql` files under `migrations/`.
4. Check React conventions:
   - `web/src/types.ts` updated with a mirror type for every new Go response struct.
   - `web/src/api.ts` updated with a function for every new endpoint.
   - Chapter component registered in `web/src/chapters/registry.ts` with correct `status: "done"`.
   - `App.tsx` imports and renders the new chapter component.
   - No `react-router` import anywhere.
5. Check `docs/architecture.md`:
   - Roadmap row status changed to `Done`.
   - `### Ch. NN — what was built` section present and accurate.
6. Store findings:
   ```bash
   npx @claude-flow/cli@latest memory store \
     --key "review-chNN" \
     --value "<findings JSON>" \
     --namespace tasks
   ```
7. Post completion:
   ```bash
   npx @claude-flow/cli@latest hooks post-task --task-id "reviewer-chNN" --success true
   ```

### Agent 5 — security-auditor

**Role:** Run a focused security and correctness audit against the anti-patterns listed in CLAUDE.md. Reports findings separately from the reviewer; does not overlap with convention checks.

**Steps:**

1. Read every new and modified `.go` file in the chapter diff.
2. Check for forbidden patterns:
   - `interface{}` or `any` in any domain struct field or function signature inside `internal/<pkg>/`.
   - Inline DDL: `CREATE TABLE`, `ALTER TABLE`, `DROP TABLE` in any `.go` file.
   - HTTPS termination: `tls.Listen`, `ListenAndServeTLS`, or TLS config in any service's `cmd/` or `internal/` package.
   - `google/uuid` import. All IDs must use `crypto/rand` hex strings.
   - `math/rand` for IDs or tokens (acceptable only for simulation tick randomness where non-cryptographic is documented).
   - `go.mod` or `go.work` files added outside the repo root.
   - Migrations referenced by number in Go code rather than via `//go:embed` in the store package.
3. Check MySQL store transactions:
   - Any operation that modifies two tables or reads-then-writes a balance must use `db.BeginTx` / `tx.Commit` / `tx.Rollback`. No multi-step updates outside a transaction.
   - `SELECT ... FOR UPDATE` must precede any balance mutation inside a transaction.
4. Check seed security:
   - Seed IDs must not collide with `NewID()` output space. Verify seed IDs start with `5eed00000000000000`.
   - `-- +goose Down` must DELETE only by known seed IDs, not `DELETE FROM <table>` (which would wipe user data on rollback).
5. Store findings:
   ```bash
   npx @claude-flow/cli@latest memory store \
     --key "security-chNN" \
     --value "<findings JSON>" \
     --namespace tasks
   ```
6. Post completion:
   ```bash
   npx @claude-flow/cli@latest hooks post-task --task-id "security-auditor-chNN" --success true
   ```

### Swarm coordination rules

- Agents 2, 3, 4, and 5 must not start until the architect (Agent 1) has posted its task completion and the design contract key `design-chNN` is present in memory.
- Agents 4 and 5 must not start until Agent 2 (coder) and Agent 3 (tester) have both posted completion.
- If the reviewer or security-auditor reports any finding with severity `error`, the coder re-opens and fixes; then the reviewer/auditor re-runs only the affected checks.
- The main session (not a swarm agent) handles: `git commit -S`, pushing the branch, and updating the draft PR description (step 5 of the manual workflow).
- After all five agents post success with no open `error`-severity findings, hand off to the manual workflow at step 5 (update PR description) and then step 6 (wait for CI).

## Sandbox limits (Cowork)

- No Go toolchain installable. Cannot run `go build`/`go test`. Always ask the user to run `make vet test build` after material Go changes.
- `.git/*.lock` files cannot be `unlink`ed from the sandbox; use `mv` to rename them out of the way before each git op. `git commit -S` will fail (no GPG key here) — leave commits to the user.
- Web build cannot be verified here either.

## Output style (Claude → user)

- No preambles ("I'll help you..."). No re-summaries of the user's prompt.
- Skip restating what was just done unless asked. The user reads the diff.
- For multi-step work: brief plan → execute → terse handoff.
- File paths use Windows form when shown to the user (`C:\Users\...`); bash uses the `/sessions/.../mnt/...` mount.
- Do not include emojis unless the user does first.
- For commit messages: conventional commit + multi-paragraph body referencing the chapter.

## Files NOT to read in full

- `web/dist/`, `web/node_modules/`, `bin/` — generated.
- `go.sum` — generated by `go mod tidy`; do not hand-edit.

## Where to look for X

- Domain types for current chapter: `services/commodity-service/internal/commodity/`
- Store contract: `services/commodity-service/internal/store/store.go`
- HTTP routes: `services/<svc>/internal/transport/httpapi/routes.go`
- Gateway proxy targets: `services/api-gateway/cmd/api-gateway/main.go`
- Env config patterns: `pkg/mysql/client.go` (`ConfigFromEnv`)
- Migration files: `services/<svc>/internal/store/migrations/`
- Migration runner: `pkg/mysql/migrate.go` (`Migrate`)
- React API client: `web/src/api.ts`
- Wire types: `web/src/types.ts` (mirror of Go structs — keep in sync)
- Chapter source text (Marx prose): vault `marx-engels/1867/capital-volume-i/texts/NN-<slug>.md` (read via `mcp__obsidian__obsidian_get_file_contents`)
- Chapter spec (code-facing view): vault `marx-engels/1867/capital-volume-i/specs/NN-<slug>.spec.md` (read via `mcp__obsidian__obsidian_get_file_contents`, or use the `chapter-spec` skill)

## Anti-patterns (do not propose)

- Per-service `go.mod` / Go workspaces. Single module is intentional.
- An ORM. Use `database/sql` with `github.com/go-sql-driver/mysql` directly via the Store interface.
- A web router (react-router). One App.tsx until a chapter forces it.
- Editing existing migration files. Migrations are append-only; add a new numbered file instead.
- Putting schema DDL inline in Go code. All DDL belongs in `internal/store/migrations/`.
- HTTPS termination in services. The cluster handles that.
- `interface{}`/`any` in domain types. Use concrete types and named ID/value types.
- Referencing `chapters/volume-1/`. That directory was retired; the vault is the source of truth.

## Done is when

`make vet test build` passes locally · `npm run lint && npm run build` passes
· `docs/architecture.md` updated · signed commit · PR open.
