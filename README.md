# Real Estate CRM

CRM-и пурра барои агентиҳои амлоки ғайриманқул.

> **Stack:** Go (Gin + GORM) backend · Nuxt 3 + Vue 3 frontend · PostgreSQL 16 · Nginx · WebSocket real-time · Telegram bot · Docker compose · GitHub Actions CI/CD.

## 📐 Архитектура

```
                              ┌────────────┐
                  HTTP / WS   │   Nginx    │
   Browser ─────────────────► │  :80 / :443│ ──┬──► Frontend (Nuxt SSR :3000)
                              └────────────┘   │
                                               ├──► Backend  (Go + Gin :8080)
                                               │       │
                                               │       └──► PostgreSQL 16
                                               │
                                               └──► /uploads (shared volume)

                              Telegram Bot ──► Backend (long-polling /api)
```

## 📁 Сохтори монорепо

```
CRM-real/
├── backend-go/                  # Go backend (Gin + GORM)
│   ├── cmd/server/main.go
│   ├── internal/{auth, config, database, handlers, middleware, models, router,
│   │             services, ws, logger, uploads}/
│   ├── migrations/
│   ├── Dockerfile
│   ├── docker-compose.yml       # standalone backend + postgres
│   └── README.md
│
├── frontend/                    # Nuxt 3 + Vue 3 + Tailwind + Pinia
│   ├── pages/, components/, layouts/, composables/, stores/, plugins/, middleware/
│   ├── types/api.ts             # mirrors Go models
│   ├── Dockerfile
│   └── README.md
│
├── nginx/
│   └── nginx.conf               # reverse proxy + WebSocket + uploads
│
├── .github/
│   ├── dependabot.yml
│   └── workflows/
│       ├── backend-ci.yml             # gofmt + vet + golangci-lint + tests
│       ├── frontend-ci.yml            # nuxt build + typecheck
│       ├── docker.yml                 # build & push to ghcr.io
│       ├── docker-compose-smoke.yml   # spins up full stack and curls /healthz
│       └── codeql.yml                 # security scanning (Go + TS)
│
├── docker-compose.yml           # FULL STACK (production)
├── docker-compose.dev.yml       # DB-only (local dev)
├── Makefile                     # `make up`, `make dev-up`, `make psql`...
├── .env.example
├── app.py                       # legacy Flask (kept for reference)
└── templates/                   # legacy Jinja templates
```

## 🚀 Шурӯъ кардан (production stack)

```bash
git clone https://github.com/instaagrammeta/CRM-real.git
cd CRM-real

cp .env.example .env
# .env-ро вироиш кунед: JWT_SECRET, ADMIN_PASSWORD, DB_PASSWORD, TELEGRAM_BOT_TOKEN

make up
# ё:  docker compose up -d --build
```

Хидматҳои дастрас:

| URL | Хидмат |
|---|---|
| `http://localhost/` | Nuxt frontend (тавассути Nginx) |
| `http://localhost/api/` | Backend REST API |
| `ws://localhost/ws` | Real-time WebSocket |
| `http://localhost/uploads/...` | файлҳои бор кардашуда |
| `http://localhost/healthz` | Backend health check |

Логин: `admin` / **парол** аз `.env` (`ADMIN_PASSWORD`).

## 🛠 Шурӯъ кардан (development)

Барои hot-reload-и backend ва frontend:

```bash
# 1. Танҳо DB-ро дар Docker гардонед:
make dev-up
# Postgres → :5432, Adminer (UI) → http://localhost:8081

# 2. Backend локалӣ:
cd backend-go
cp .env.example .env
go mod tidy
make run        # ё `go run ./cmd/server`
# → http://localhost:8080

# 3. Frontend локалӣ:
cd frontend
cp .env.example .env
npm install
npm run dev
# → http://localhost:3000
```

## 🐳 Docker Compose командаҳои муҳим

| Команда | Тавсиф |
|---|---|
| `make up` | Full stack: postgres + backend + frontend + nginx |
| `make down` | Stop |
| `make logs` | Tail logs |
| `make restart` | Restart all services |
| `make build` | Rebuild images |
| `make clean` | Stop **ва** volumes-ро тоза кунад (DB пок мешавад!) |
| `make dev-up` / `make dev-down` | Танҳо Postgres + Adminer |
| `make psql` | Кушодани `psql` дар container-и postgres |

## 🔄 CI/CD (GitHub Actions)

| Workflow | Trigger | Чӣ кор мекунад |
|---|---|---|
| `backend-ci.yml` | push/PR ба `backend-go/` | gofmt · go vet · golangci-lint · go build · go test (бо Postgres service) |
| `frontend-ci.yml` | push/PR ба `frontend/` | npm ci · nuxt prepare · vue-tsc · nuxt build |
| `docker-compose-smoke.yml` | push/PR ба infra | `docker compose up` → curl `/healthz` → login → tear down |
| `docker.yml` | push ба `main`/tags | Multi-arch image-ҳоро ба `ghcr.io` push мекунад |
| `codeql.yml` | push/PR + ҳафтавор | Скани амнияти статикӣ барои Go ва TypeScript |
| `dependabot.yml` | автоматӣ | PR-и нав барои нав кардани go modules + npm + actions |

Ҳамаи workflow-ҳо дар `Pull request` тригер мешаванд — ҳар PR метавонад дар як ҷадвал тафтиш шавад.

## 🔐 Амнияти production

Барои деплой ба сервери воқеӣ:

1. **Парол ва секретҳоро тағйир диҳед** (`JWT_SECRET`, `DB_PASSWORD`, `ADMIN_PASSWORD`).
2. Дар `nginx/nginx.conf` блоки `server { listen 443 ssl; ... }` илова кунед бо сертификати Let's Encrypt (тавассути `certbot` ё `traefik`).
3. Postgres port-ро ба берун expose накунед (port-и `5432` дар `docker-compose.yml` пешфарз баста аст).
4. `CORS_ORIGINS=https://yourdomain.com` гузоред.
5. Backup-и автоматии Postgres гузоред (масалан тавассути `cron` + `pg_dump`).

## 📦 Docker images

Тавассути workflow-и `docker.yml` имиҷҳо ба GitHub Container Registry push мешаванд:

- `ghcr.io/instaagrammeta/crm-backend:latest`
- `ghcr.io/instaagrammeta/crm-frontend:latest`

## 🗺 Roadmap

Ниг. PR-и қаблӣ ва issues. Қадамҳои навбатии тавсияшаванда:
- [ ] Тестҳои backend (`go test ./...`) бо coverage > 60%
- [ ] HTTPS + Let's Encrypt дар `nginx/`
- [ ] PostgreSQL backup script (cron-based)
- [ ] Картаи интерактивӣ (Leaflet ё Mapbox) дар `/karta`
- [ ] Excel/PDF export-ҳои порт-кардашуда аз нусхаи Flask
- [ ] Multi-tenancy (барои фурӯши SaaS)

## 📜 Litsenzия

Хусусии internal — пас аз омодагӣ litsenz интихоб кунед (MIT / Apache-2.0).
