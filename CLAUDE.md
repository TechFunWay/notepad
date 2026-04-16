# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Multi-user notepad web application targeting Feiniu NAS (飞牛NAS), Docker, and direct binary deployment. Single Go binary embeds the Vue 3 frontend via `//go:embed`. Version tracked in `VERSION` file.

## Build & Development Commands

```bash
# Local development (run in separate terminals)
make dev-server          # Go server on :8904 (default)
make dev-web             # Vite dev server, proxies /api -> :8904

# Production build
make build               # Frontend + Go binary -> release/<VERSION>/notepad

# Cross-platform builds
make cross-compile       # linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
make docker              # Multi-platform Docker image (wycto/notepad)
make fpk                 # Feiniu NAS FPK package
make all-build           # cross-compile + fpk
```

No test suite or linter is configured in this project.

## Architecture

### Backend (Go + Gin + SQLite)

- **Entry**: `server/main.go` — dispatches between CLI mode (`cmd/cli.go`) and server mode (`cmd/server.go`) based on CLI args
- **Config**: Environment variables only — `PORT` (8904), `DB_PATH`, `JWT_SECRET`, `DATA_DIR`. See `server/config/config.go`
- **Database**: `modernc.org/sqlite` (pure Go, no CGO). WAL mode, max 1 connection. Custom migration system in `server/database/` with `schema_migrations` table — add new migrations as numbered functions in `migrations.go`
- **Models**: Raw SQL, no ORM. `database.DB` is the global `*sql.DB` instance
- **Auth**: JWT via `Authorization: Bearer` header or `token` cookie. First registered user becomes admin automatically
- **Static**: Frontend built into `server/static/dist/`, embedded via Go embed. SPA fallback serves `index.html` for non-API routes

### Frontend (Vue 3 + Element Plus + TipTap)

- **Build**: Vite with `@vitejs/plugin-vue`. Dev proxy in `vite.config.js` forwards `/api` to `:8904`
- **State**: Pinia stores in `web/src/stores/` — `auth.js` (token/user in localStorage), `config.js` (site title, registration toggle)
- **API layer**: Axios instance in `web/src/api/request.js` — auto-attaches Bearer token, auto-redirects on 401
- **Rich text**: TipTap editor in `web/src/components/TiptapEditor.vue`
- **Routing**: HTML5 history mode. Auth guards in `web/src/router/index.js` — unauthenticated users redirected to `/login`

### API Routes (all under `/api`)

Public: `/auth/register`, `/auth/login`, `/auth/security-question`, `/auth/forgot-password`, `/public-config`, `/version`, `/health`

Authenticated: `/notes` (CRUD), `/notes/tags`, `/auth/logout`, `/auth/change-password`

Admin: `/users` (CRUD), `/configs` (read + update by key)

### Database Schema (migrations)

- `users` — id, username, password_hash, security_question, security_answer_hash, role, timestamps
- `notes` — id, user_id (FK), title, content, tags, timestamps
- `configs` — id, key, value, description, updated_at

## Key Patterns

- Version/build info injected via Go ldflags: `-X main.Version=... -X main.BuildTime=... -X main.GitCommit=...`
- `CGO_ENABLED=0` for all Go builds (pure Go SQLite)
- CLI tools: `recover-admin` (reset admin password), `find-admin` (show admin info), `list-users` (list all users)
- Deployment data volume at `/app/data` (Docker) or `./data` (local)
