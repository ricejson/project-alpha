# ---- build stage ----
FROM golang:1.23-alpine AS build
WORKDIR /src

# Pre-fetch modules for better layer caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build a fully static binary (no cgo) so it runs on Alpine
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
        -o /out/server ./cmd/server

# ---- runtime stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=build /out/server /app/server
USER app
EXPOSE 8080
ENV HTTP_ADDR=:8080 \
    APP_ENV=production \
    GIN_MODE=release
ENTRYPOINT ["/app/server"]
