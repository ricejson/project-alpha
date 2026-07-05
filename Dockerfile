# syntax=docker/dockerfile:1.7

# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Cache module layer
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build a static, stripped binary
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server && \
    go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

# ---------- Runtime stage ----------
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /out/server /app/server
COPY --from=builder /out/seed /app/seed
COPY --from=builder /src/seed.sql /app/seed.sql

ENV APP_ENV=production \
    HTTP_ADDR=:8080 \
    TZ=Asia/Shanghai

USER app
EXPOSE 8080

ENTRYPOINT ["/app/server"]
