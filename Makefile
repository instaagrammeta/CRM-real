.PHONY: help up down logs restart build clean dev-up dev-down backend-run frontend-dev test fmt psql

help:
	@echo "CRM Real Estate — make targets"
	@echo "  make up          – start full stack (postgres + backend + frontend + nginx)"
	@echo "  make down        – stop full stack"
	@echo "  make logs        – tail logs of all services"
	@echo "  make restart     – restart full stack"
	@echo "  make build       – rebuild Docker images"
	@echo "  make clean       – stop & remove volumes (DROPS DB!)"
	@echo "  make dev-up      – start dev DB only (Postgres + Adminer)"
	@echo "  make dev-down    – stop dev DB"
	@echo "  make backend-run – run Go backend locally"
	@echo "  make frontend-dev – run Nuxt frontend locally"
	@echo "  make test        – run all backend tests"
	@echo "  make fmt         – format Go code"
	@echo "  make psql        – open psql in running postgres container"

# ---------- prod stack ----------
up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f --tail=200

restart:
	docker compose restart

build:
	docker compose build

clean:
	docker compose down -v

# ---------- dev stack ----------
dev-up:
	docker compose -f docker-compose.dev.yml up -d

dev-down:
	docker compose -f docker-compose.dev.yml down

# ---------- local processes ----------
backend-run:
	$(MAKE) -C backend-go run

frontend-dev:
	cd frontend && npm run dev

test:
	$(MAKE) -C backend-go test

fmt:
	$(MAKE) -C backend-go fmt

# ---------- DB shell ----------
psql:
	docker compose exec postgres psql -U $${DB_USER:-crm} -d $${DB_NAME:-crm_real}
