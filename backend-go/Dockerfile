# ============== build stage ==============
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/crm-server ./cmd/server

# ============== runtime stage ==============
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/crm-server /app/crm-server
COPY migrations ./migrations

RUN mkdir -p /app/uploads/photos /app/uploads/files /app/uploads/chat /app/uploads/posts && \
    chown -R app:app /app

USER app

EXPOSE 8080

ENTRYPOINT ["/app/crm-server"]
