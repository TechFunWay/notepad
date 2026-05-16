# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Multi-user notepad web application targeting Feiniu NAS (飞牛NAS), Docker, and direct binary deployment. Single Go binary embeds the Vue 3 frontend via `//go:embed`. Version tracked in `VERSION` file.

## Build & Development Commands

```bash
# Local development - run in separate terminals
make dev-server          # Go backend on :8904 (auto-kills existing process)
make dev-web             # Vite dev server on :3000, proxies /api -> :8904

# Full-stack local (builds frontend then starts backend)
make dev                 # npm ci -> npm run build -> copy dist -> go run on :8904

# Production (frontend + Go binary -> release/<VERSION>/notepad)
make build

# Cross-platform
make cross-compile       # linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
make docker              # Multi-platform Docker image (wycto/notepad)
make build-fpk           # Cross-compile + Feiniu NAS FPK package

make clean               # Remove release/, dist/, and web/dist
```

No test suite, linter, type checker, or CI is configured in this project.

## Architecture

Two independent packages in a monorepo: `server/` (Go, `go.mod` root) and `web/` (Vue 3, `package.json` root). No workspace tooling links them.

### Backend (Go + Gin + SQLite)

- **Entry**: `server/main.go` — dispatches between CLI mode (`cmd/cli.go`) and server mode (`cmd/server.go`) when first arg doesn't start with `-`
- **Config**: `config.Load()` reads env vars with sensible defaults. See `server/config/config.go`. Supports `PORT`, `DATA_DIR`, `DB_PATH`, `JWT_SECRET`, `WEB_DIR`, `UPLOAD_DIR`, `SHARE_DIRS` (colon-separated, served under `/uploads`)
- **Module path**: `notepad` (in `server/go.mod`)
- **Database**: `modernc.org/sqlite` (pure Go, no CGO). WAL mode, `SetMaxOpenConns(1)`. Global `database.DB *sql.DB`. See `server/database/database.go`
- **Migrations**: Custom upgrade system in `server/database/upgrade.go`. Uses `upgrade_records` table (not `schema_migrations`). Append to the `upgrades` slice as `{"version", func(tx)}`. Version in `VERSION` file (no `v` prefix). Anti-downgrade protection at startup
- **Auth**: JWT via `Authorization: Bearer` header or `token` cookie. First registered user becomes admin automatically
- **Middleware**: `CORS()` (wide open), `RequireAuth()` (parses JWT, sets `userID`/`username`/`role` in Gin context), `RequireAdmin()` (checks role == "admin")
- **File upload**: POST `/api/upload` (auth required). Accepts jpg/png/gif/webp only, saves to `UPLOAD_DIR` with UUID filename, served under `/uploads/`
- **Logger**: Simple file logger in `server/logger/logger.go`, writes to `DATA_DIR/logs/`
- **Static**: Frontend built into `server/static/dist/`, embedded via `//go:embed` in `server/static/embed.go`. SPA fallback serves `index.html` for non-API routes. Priority: external `webDir` path > external `dist` dir > embedded FS
- **API routes** (all under `/api`):
  - Public: `auth/register`, `auth/login`, `auth/security-question`, `auth/verify-answer`, `auth/forgot-password`, `public-config`, `version`, `health`
  - Authenticated: `auth/logout`, `auth/change-password`, `auth/security-question` (PUT), `upload`, `notes` (CRUD), `notes/tags`
  - Admin: `users` (CRUD), `configs` (list + update by key)
- **Database tables**: `users` (id, username, password_hash, security_question, security_answer_hash, role, timestamps), `notes` (id, user_id FK, title, content, tags, timestamps), `configs` (id, key, value, description, updated_at)
- **CLI tools**: `./notepad find-admin`, `./notepad recover-admin`, `./notepad list-users`

### Frontend (Vue 3 + Element Plus + TipTap + Pinia)

- **Build**: Vite with `@vitejs/plugin-vue`, `@` alias maps to `src/`. Dev proxy in `vite.config.js` forwards `/api` to `:8904`
- **Routing**: HTML5 history mode. Auth guards redirect unauthenticated to `/login`, non-admin away from `/admin/*`. Routes: `/login`, `/register`, `/forgot-password` (public), `/` + `/notes-list` + `/profile` (auth), `/admin/users` + `/admin/configs` (admin)
- **State**: Pinia stores in `web/src/stores/` — `auth.js` (token/user in localStorage), `config.js` (site title, registration toggle)
- **API layer**: Axios instance in `web/src/api/request.js` — auto-attaches Bearer token from localStorage, redirects to `/login` on 401
- **Rich text**: TipTap editor in `web/src/components/TiptapEditor.vue` with extensions for placeholder, color, highlight, underline, text-align, image

## Git 提交规范

- 提交备注**全部使用中文**，禁止使用英文
- 格式：`类型：描述`
- 常用类型：`新增`、`修复`、`优化`、`更新`、`清理`、`文档`、`重构`
- 尾部追加 `Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>`
- 示例：
  ```
  新增：手机端浮动新建按钮

  Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
  ```

## Key Patterns

- Version/build info injected via Go ldflags: `-X main.Version=$(VERSION) -X main.BuildTime=... -X main.GitCommit=...`
- `CGO_ENABLED=0` for all Go builds — `modernc.org/sqlite` is pure Go, never enable CGO
- Deployment data volume at `/app/data` (Docker) or `./data` (local)
- DB_PATH and DATA_DIR can differ — DB_PATH is `DATA_DIR/notepad.db` by default
- All frontend API modules (`web/src/api/`) return Axios promises with automatic error handling
- Raw SQL, no ORM. `database.DB` is accessed globally
