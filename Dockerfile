# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY --from=builder /app/api .
# Tablas NOM se leen desde filesystem (no están embebidas)
COPY data/ ./data/

EXPOSE 8080

CMD ["./api"]
