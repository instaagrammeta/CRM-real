# CRM Real Estate — Frontend (Nuxt 3 + Vue 3)

Frontend-и SPA барои CRM-и амлоки ғайриманқул. Бо backend-и Go (Gin/GORM/Postgres) аз тариқи REST API + WebSocket кор мекунад.

## Стек

- **Nuxt 3** + **Vue 3** + **TypeScript**
- **Pinia** (state management)
- **Tailwind CSS** (зерсохти Donezo green)
- **VueUse** (utility composables)
- **vuedraggable** (Kanban drag-and-drop)
- **WebSocket** мизоҷ (real-time chat + notifications)

## Сохтори лоиҳа

```
frontend/
├── app.vue                  # bootstrap (auth.restore + WS)
├── nuxt.config.ts           # Nuxt config (modules, runtime, head)
├── tailwind.config.ts       # Donezo green palette + helpers
├── tsconfig.json
│
├── assets/css/main.css      # Tailwind layers + utility classes (.btn-*, .card, .badge-*, ...)
├── public/                  # favicon, статикии ҷамъиятӣ
│
├── types/api.ts             # TS interfaces (mirror Go models)
│
├── composables/
│   ├── useApi.ts            # $fetch wrapper + JWT + 401 → /login
│   ├── useWebSocket.ts      # Singleton WS hub with auto-reconnect
│   ├── useToast.ts          # In-app toast service
│   └── useFormat.ts         # дата, валюта, initials, timeAgo
│
├── stores/
│   ├── auth.ts              # login/logout/restore/me
│   └── notifications.ts     # list, unread count, mark-read, WS push
│
├── middleware/
│   └── auth.global.ts       # JWT guard + admin guard + public routes
│
├── plugins/
│   └── websocket.client.ts  # WS event → store/toast routing
│
├── layouts/
│   ├── default.vue          # collapsible sidebar + topbar + notifications
│   ├── auth.vue             # login layout
│   └── public.vue           # ipoteka/rasrochka public pages
│
├── components/
│   ├── Modal.vue, ConfirmDialog.vue, Toaster.vue, NotificationsDropdown.vue
│   ├── Avatar.vue, Spinner.vue, EmptyState.vue, PageHeader.vue, NavItem.vue, StatCard.vue
│   └── Kanban/
│       ├── KanbanBoard.vue   # доскаи horizontal + drag-drop колоннаҳо
│       ├── KanbanColumn.vue  # колонна + drag-drop кардҳо
│       └── KanbanCard.vue    # кард барои lid ё request (auto-detect)
│
└── pages/
    ├── login.vue, index.vue, profile.vue, notifications.vue
    ├── lids.vue, zayavka.vue, zadacha.vue
    ├── chat.vue, posts.vue, houses.vue, baza.vue, karta.vue, obiekt.vue
    ├── sim-cards.vue
    ├── ipoteka.vue, ipoteka/[slug].vue
    ├── rasrochka.vue, rasrochka/[slug].vue
    └── admin/
        ├── users.vue
        ├── ipoteka.vue
        └── rasrochka.vue
```

## Шурӯъ кардан

### 1. Зависимостҳо

```bash
cd frontend
npm install
```

### 2. Env

```bash
cp .env.example .env
# .env
# NUXT_PUBLIC_API_BASE=http://localhost:8080
# NUXT_PUBLIC_WS_BASE=ws://localhost:8080
```

### 3. Dev

```bash
npm run dev
# → http://localhost:3000
```

Backend бояд дар port 8080 фаъол бошад (ниг. `../backend-go`).

### 4. Production build

```bash
npm run build
npm run preview
```

### 5. Docker

```bash
docker build -t crm-frontend .
docker run -p 3000:3000 \
  -e NUXT_PUBLIC_API_BASE=http://your-api:8080 \
  -e NUXT_PUBLIC_WS_BASE=ws://your-api:8080 \
  crm-frontend
```

## Хусусиятҳо

### 🔐 Auth
- JWT дар `localStorage` (`auth_token`) + http-only cookie
- Auto-attach `Authorization: Bearer ...` ба тамоми $fetch
- Хатои 401 → автоматӣ ба `/login` бо redirect

### ⚡ Real-time
- Singleton WebSocket бо exponential reconnect (1s → 30s)
- `useWebSocket().on('chat:message', ...)` барои подписка
- Pages худ ба events подписка мешаванд (chat, lids, zayavka)

### 🔔 Notifications
- `notifications` store + WS `notification` event
- Dropdown дар topbar бо unread count
- Toast pop-up барои оғоҳномаҳои нав
- Саҳифаи пурра `/notifications`

### 🤖 Telegram
- Профил → `Generate link token` → `/start <token>` дар бот → пайваст
- Рӯйхати pуйвасткунӣҳо + унлинк
- Тамоми оғоҳномаҳо ба Telegram низ мерасанд

### 🎨 UI/UX
- Donezo green theme (`#1f7a4d`)
- Tailwind utility classes + кастоми `.btn-primary`, `.card`, `.badge-brand`, ...
- Compact / collapsible sidebar (Cmd/Ctrl-friendly)
- Mobile responsive (sidebar slides over)
- Animation: `fade-in`, `slide-in-right`, `pulse-soft`

### 📋 Kanban
- Drag-drop колоннаҳо (handle) ва кардҳо (хам байни колоннаҳо)
- Optimistic UI updates + WS sync
- Boards CRUD + Columns CRUD + Cards CRUD
- Истифода дар `/lids` ва `/zayavka`

### 💬 Chat
- WS-based real-time
- File attachments (multipart upload)
- Message bubbles (own = brand-600, others = page)
- Edit / delete (танҳо барои худи муаллиф ё admin)

## Endpoints (мухтасар)

Тамоми API endpoints дар `composables/useApi.ts` тавассути JWT + base URL фарохонида мешаванд.
Барои тафсилот ниг. backend `README.md`.

## Сохтан-кушодан

```bash
# Backend (Go)
cd backend-go
docker compose up -d
# → :8080

# Frontend (Nuxt)
cd ../frontend
npm install
npm run dev
# → :3000
```

Логин: `admin` / `nav-xona@2026` (агар .env-и backend `ADMIN_PASSWORD`-ро тағйир надода бошад).

## Маҳдудиятҳо

- `node_modules` дар sandbox-и Kiro насб нашуд (network: INTEGRATIONS_ONLY); `npm install`-ро дар компютери худатон иҷро кунед — ҳама зависимостҳо дар `package.json` рӯйхат шудаанд.
- Картаи интерактивӣ (`/karta`) ҳозир ба Google Maps пайванд аст; тағйир ба Leaflet ё Mapbox дар фазаи навбатӣ ба нақша гирифта мешавад.
- Excel/PDF export (барои тестҳо мебуд) метавонад тавассути backend endpoints (бо `excelize`) илова шавад.
