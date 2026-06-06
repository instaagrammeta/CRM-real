# CRM Real Estate — Go Backend

Backend-и Go (Gin + GORM + PostgreSQL) барои CRM-и амлоки ғайриманқул.
Нусхаи аз нав навишташуда аз варианти аввалаи Flask (Python) ба Go.

## Хусусиятҳо

- 🔐 **JWT auth** + bcrypt (sha256-и legacy низ дастгирӣ мешавад)
- 🗄 **PostgreSQL** + GORM AutoMigrate + migrations SQL
- ⚡ **WebSocket** (gorilla/websocket) барои чат ва оғоҳномаҳои real-time
- 🔔 **Системаи оғоҳнома** дохилӣ + push тавассути **Telegram bot**
- 📦 Модулҳо: users, auth, tasks, lids, kanban, requests, houses, realty,
  objekt (шахматка), sim cards, posts, chat, folders, ипотека, рассрочка,
  dashboard, notifications
- 🐳 **Docker** + `docker-compose` (Postgres + backend дар як команда)
- 📁 **File uploads** (`/uploads/...`) аз 50 MB
- 📝 Logging бо **zerolog**
- 🧱 Сохтори тоза (Blueprints / clean architecture)

## Сохтори лоиҳа

```
backend-go/
├── cmd/server/         # main.go - entrypoint
├── internal/
│   ├── auth/           # JWT, bcrypt
│   ├── config/         # env config
│   ├── database/       # gorm + AutoMigrate
│   ├── handlers/       # HTTP handlers (CRUD)
│   ├── middleware/     # auth, admin, CORS, logger, recovery
│   ├── models/         # GORM models
│   ├── router/         # маршрутҳо
│   ├── services/       # notifier, telegram, seed
│   ├── uploads/        # file save helper
│   ├── ws/             # WebSocket hub
│   └── logger/         # zerolog setup
├── migrations/         # SQL migrations (snapshot)
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Шурӯъ кардан (локалӣ)

```bash
cp .env.example .env
# .env-ро вироиш кунед — JWT_SECRET, DB_PASSWORD, TELEGRAM_BOT_TOKEN

go mod tidy
go run ./cmd/server
```

Сервер ба `http://localhost:8080` шунавӣ мекунад.

### Тавассути Docker

```bash
docker compose up -d --build
```

Ин Postgres-ро дар port 5432 ва backend-ро дар 8080 мегардонад.

## Аввалин корбар

Ҳангоми оғози аввал, корбари admin аз `.env` сохта мешавад:
- логин: `admin`
- парол: `nav-xona@2026` (онро дар production иваз кунед!)

## Endpoints (асосӣ)

### Public
- `POST /api/login` — `{ login, password }` → `{ token, user, role }`
- `POST /api/logout`
- `GET  /api/banks` , `/api/banks/by-slug/:slug`
- `GET  /api/installment-objects` , `/api/installment-objects/by-slug/:slug`
- `GET  /healthz`
- `GET  /uploads/...` — статикии файлҳо

### Authenticated (JWT: `Authorization: Bearer <token>` ё cookie `auth_token`)
- `GET  /api/check-auth` , `GET /api/me`
- `GET  /api/users` , `/api/users/list`
- `GET  /api/tasks` , `POST /api/tasks` , `PUT /api/tasks/:id` , `DELETE /api/tasks/:id`
- `GET/POST/PUT/DELETE /api/lids` , `/api/lids/:id`
- Канбан лидҳо: `/api/kanban/boards` , `/api/kanban/columns` , `/api/kanban/leads` ,
  `POST /api/kanban/leads/move`
- Канбан заявкаҳо: `/api/requests-boards` , `/api/requests-columns` ,
  `/api/requests-new` , `POST /api/requests-move`
- `/api/houses` , `/api/realty/...` , `/api/objekt/...`
- `/api/company-phones` , `/api/sim-cards` , `/api/sim-tariffs` , `/api/sim-phones/stats`
- `/api/posts` , `/api/messages` , `/api/folders` , `/api/folders/:id/files`
- `GET /api/dashboard/stats`
- `GET /api/notifications` , `POST /api/notifications/:id/read` , `POST /api/notifications/read-all`

### Telegram pуйвастсозӣ
1. Корбар: `POST /api/me/telegram/link` → `{ token, link_url, command }`
2. Дар Telegram: pбот-ро кушоед ва `/start <token>` равон кунед
3. Bot чат-ро бо корбар pуйваст мекунад. Ҳама оғоҳномаҳои нав
   ҳамзамон ба ин чат фиристода мешаванд.

### WebSocket
- `GET /ws?token=<jwt>` — pуйвастан
- Сервер event-ҳои зеринро мефиристад:
  - `chat:message`, `chat:update`, `chat:delete`
  - `notification` (фақат ба худи корбар)
  - `lead:created`, `lead:moved`, `lead:updated`
  - `request:created`, `request:moved`
  - `task:assigned`, `tariff:expiring`

### Admin only
- `POST /api/users` , `PUT /api/users/:id` , `DELETE /api/users/:id`
- `POST/PUT/DELETE /api/banks` ,  `/api/conditions/:id`
- `POST/PUT/DELETE /api/installment-objects` , `/api/installment-conditions/:id`

## Telegram bot

Барои фаъол кардан:
1. Дар `@BotFather` бот созед
2. Token-ро ба `TELEGRAM_BOT_TOKEN` дар `.env` гузоред
3. Бот командаҳои зеринро дастгирӣ мекунад:
   - `/start <token>` — pуйвастан
   - `/me` — ҳолати ҳисоб
   - `/stop` — қатъ кардани оғоҳномаҳо
   - `/help`

## Migrations

`internal/database/database.go` ҳангоми оғоз AutoMigrate мекунад
(тамоми ҷадвалҳо аз модели Go).

Барои deploy-и production иловатан snapshot-и SQL дар `migrations/` мавҷуд аст,
ки бо `golang-migrate` ё `psql` иҷро карда мешавад.

## Frontend

Frontend (Nuxt + Vue) дар марҳилаи навбатӣ навишта мешавад.
