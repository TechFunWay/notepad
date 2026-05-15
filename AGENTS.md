# AGENTS.md

Compact guidance for AI agents. Also see `CLAUDE.md` for fuller context.

## Commands

```bash
make dev-server    # Go backend on :8904 (kills existing process first)
make dev-web       # Vite dev server on :3000, proxies /api -> :8904
make build         # Frontend + Go binary -> release/<VERSION>/notepad
make cross-compile # linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
make clean         # Remove build artifacts
```

No tests. No linter. No typecheck. No CI. Skip all verification steps.

## Architecture gotchas

- **Two-package monorepo**: `server/` (Go, `go.mod` root) and `web/` (Vue, `package.json` root). They are independent — no workspace tooling links them.
- **Frontend embeds into Go**: `make build` runs `npm ci && npm run build` in `web/`, copies `web/dist` into `server/static/dist/`, then builds the Go binary. The binary embeds `server/static/dist/` via `//go:embed` in `server/static/embed.go`.
- **SQLite is pure Go** (`modernc.org/sqlite`): `CGO_ENABLED=0` everywhere. Never enable CGO.
- **Single DB connection**: `DB.SetMaxOpenConns(1)` in `server/database/database.go`. Concurrent writes will queue.
- **No ORM**: all raw SQL. `database.DB` is the global `*sql.DB`.

## Migrations

Custom upgrade system in `server/database/upgrade.go`, not a standard migration tool.

- Add new upgrades by appending to the `upgrades` slice: `{"X.Y.Z", upgrade_fn}`
- Version must match the `VERSION` file (without `v` prefix)
- Upgrade function receives `*sql.Tx` — all statements in one transaction
- `nil` fn = version-only record (no schema change)
- Anti-downgrade: DB version > binary version = fatal error at startup

## Auth & routing

- First registered user auto-becomes admin (role column, not a separate table)
- JWT via `Authorization: Bearer` header or `token` cookie
- SPA fallback: all non-`/api` routes serve `index.html`
- Frontend routes: `/login`, `/register`, `/forgot-password` are public; everything else requires auth; `/admin/*` requires admin role

## Environment variables

| Var | Default | Notes |
|-----|---------|-------|
| PORT | 8904 | Server port |
| DB_PATH | ./data/notepad.db | SQLite path |
| JWT_SECRET | random per start | Set a fixed value for persistence |
| DATA_DIR | ./data | Data directory |

## Build info injection

Version, build time, and git commit injected via ldflags at build time:
`-X main.Version=... -X main.BuildTime=... -X main.GitCommit=...`
Source of truth for version: `VERSION` file at repo root.

## CLI tools (run from binary)

```bash
./notepad find-admin      # Show admin username
./notepad recover-admin   # Reset admin password
./notepad list-users      # List all users
```

These dispatch from `server/main.go` when the first arg doesn't start with `-`.

## Skills

- **Location**: All skills are stored in the `.skills/` directory at the project root.
- **Usage**: Reference skills by their path relative to the project root (e.g., `.skills/skill-name.md`).
- **Organization**: Keep skill files organized and up-to-date. Remove obsolete skills.

## Deployment targets

- **Docker**: `Dockerfile` multi-stage build, data volume at `/app/data`, healthcheck at `/api/health`
- **Feiniu NAS (飞牛)**: FPK package built by `make fpk` → `scripts/build-fnpack.sh`
- **Direct binary**: just run `./notepad`, data in `./data/`
