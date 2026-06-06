.PHONY: run build tidy fmt vet test docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/crm-server ./cmd/server

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

vet:
	go vet ./...

test:
	go test ./... -v

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down
