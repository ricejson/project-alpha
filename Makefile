SHELL := /bin/bash

BACKEND_DIR := backend
FRONTEND_DIR := frontend

BACKEND_ADDR ?= :8080
FRONTEND_HOST ?= 127.0.0.1
FRONTEND_PORT ?= 5173
DATABASE_URL ?= host=localhost user=project_alpha password=project_alpha dbname=project_alpha port=5432 sslmode=disable TimeZone=Asia/Shanghai
CORS_ALLOWED_ORIGINS ?= http://localhost:$(FRONTEND_PORT)
VITE_API_BASE_URL ?= http://localhost:8080/api/v1

.PHONY: help dev db backend frontend install fmt lint test build precommit clean deploy

help:
	@echo "project-alpha make targets:"
	@echo "  make dev       Start PostgreSQL, backend, and frontend"
	@echo "  make db        Start PostgreSQL only"
	@echo "  make backend   Start Gin backend on $(BACKEND_ADDR)"
	@echo "  make frontend  Start Vite frontend on http://$(FRONTEND_HOST):$(FRONTEND_PORT)"
	@echo "  make install   Install frontend dependencies"
	@echo "  make fmt       Format Go code"
	@echo "  make lint      Run Go vet and frontend production build"
	@echo "  make test      Run backend and frontend tests"
	@echo "  make build     Build frontend"
	@echo "  make precommit Run all pre-commit hooks"
	@echo "  make clean     Remove frontend build outputs"
	@echo "  make deploy    Deploy to Fly.io + Cloudflare Pages (one command)"

dev: db
	@trap 'kill 0' INT TERM EXIT; \
	$(MAKE) backend & \
	$(MAKE) frontend & \
	wait

db:
	docker compose up -d postgres

backend:
	cd $(BACKEND_DIR) && \
	APP_ENV=development \
	HTTP_ADDR="$(BACKEND_ADDR)" \
	DATABASE_URL="$(DATABASE_URL)" \
	CORS_ALLOWED_ORIGINS="$(CORS_ALLOWED_ORIGINS)" \
	go run ./cmd/server

frontend:
	cd $(FRONTEND_DIR) && \
	VITE_API_BASE_URL="$(VITE_API_BASE_URL)" \
	npm run dev -- --host $(FRONTEND_HOST) --port $(FRONTEND_PORT)

install:
	cd $(FRONTEND_DIR) && npm ci

fmt:
	cd $(BACKEND_DIR) && gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

lint:
	cd $(BACKEND_DIR) && go vet ./...
	cd $(FRONTEND_DIR) && npm run build

test:
	cd $(BACKEND_DIR) && go test ./...
	cd $(FRONTEND_DIR) && npm test

build:
	cd $(FRONTEND_DIR) && npm run build

precommit:
	pre-commit run --all-files

clean:
	rm -rf $(FRONTEND_DIR)/dist
	rm -f $(FRONTEND_DIR)/*.tsbuildinfo

deploy:
	./scripts/deploy.sh
