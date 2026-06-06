# =====================================================================
# Multi-stage build for the Go backend.
# Produces a tiny static binary on top of alpine:3.20.
# =====================================================================

# ============== build stage ==============
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Cache go modules
COPY go.mod go.sum* ./
RUN go mod download

# Copy source and build
COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE}" \
    -o /out/crm-server ./cmd/server

# ============== runtime stage ==============
FROM alpine:3.20

# wget is included via busybox; add tzdata + ca-certs.
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/crm-server /app/crm-server
COPY migrations ./migrations

RUN mkdir -p /app/uploads/photos /app/uploads/files /app/uploads/chat /app/uploads/posts && \
    chown -R app:app /app

USER app

EXPOSE 8080

# A simple healthcheck for orchestrators that don't read compose-level checks.
HEALTHCHECK --interval=15s --timeout=5s --retries=4 --start-period=10s \
  CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/app/crm-server"]
