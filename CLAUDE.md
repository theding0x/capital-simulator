# CLAUDE.md

Project memory for Claude. Read this once per session; do not re-derive.

## What this is

Microservice simulation of an economy modeled chapter-by-chapter on Marx's
*Capital, Vol. I*. One branch + one PR + one HTML in `chapters/volume-1/` per chapter.

## Stack (fixed; do not propose alternatives)

Go 1.25 monorepo · React 18 + Vite + TS · MySQL 8 · Redis · Docker · k8s
(kustomize) · GitHub. Module path: `github.com/theding0x/capital-simulator`.

## Layout

```
pkg/{log,httpx,mysql,redis}/    shared Go libs
services/<svc>/cmd/<svc>/       main.go (one per service)
services/<svc>/internal/        domain, store, transport
services/<svc>/Dockerfile       multi-stage, distroless, build context = repo root
web/                            Vite+React+TS, nginx prod image, /api proxy
deploy/k8s/                     kustomize: namespace + infra + per-service
chapters/volume-1/NN-<slug>.html per-chapter summaries (large; do not Read in full)
docs/architecture.md            authoritative topology + chapter status
```

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
2. Open a **draft** PR against `main` populated from `.github/pull_request_template.md`, with the *planned* changes drawn from the chapter spec under `chapters/volume-X/Y-<slug>.spec.md` (Chapter, Summary, Services touched, Chapter HTML, planned tests, Notes for review).
3. Implement the planned chapter changes (domain types in the relevant service(s), tests using Marx's examples, seed migration `NNNNN_chNN_seed.sql` for every new domain type, chapter HTML drop, `docs/architecture.md` roadmap row flipped to Done).
4. Commit signed (`commit.gpgsign=true` is set per-repo). Use a multi-line conventional commit; the PR description fills from it.
5. Compare the committed changes against the draft PR and update the PR description so it reflects what actually landed.
6. Wait for GitHub Actions to finish running checks on the PR.
7. If checks fail, fetch the output (`gh run view --log-failed`), fix the failure, and commit again. Loop to step 5.
8. If all checks pass, mark the PR ready for review and notify the user that the PR is ready to merge into `main`.

> **Note:** If the user says "implement the next pending chapter", look up the next `Pending` row in `docs/architecture.md` (Roadmap table) — that's the authoritative source for chapter ordering and primary services.

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

- `chapters/volume-1/*.html` — Marx source text, 100KB+ per file. Use `python3` strip to get prose if needed.
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

## Anti-patterns (do not propose)

- Per-service `go.mod` / Go workspaces. Single module is intentional.
- An ORM. Use `database/sql` with `github.com/go-sql-driver/mysql` directly via the Store interface.
- A web router (react-router). One App.tsx until a chapter forces it.
- Editing existing migration files. Migrations are append-only; add a new numbered file instead.
- Putting schema DDL inline in Go code. All DDL belongs in `internal/store/migrations/`.
- HTTPS termination in services. The cluster handles that.
- `interface{}`/`any` in domain types. Use concrete types and named ID/value types.

## Done is when

`make vet test build` passes locally · `npm run lint && npm run build` passes
· chapter HTML present · `docs/architecture.md` updated · signed commit · PR open.
