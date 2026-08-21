SHELL := /bin/bash

.PHONY: fmt migrate test vet build web-install web-lint web-test web-build verify up down

fmt:
	gofmt -w $$(find cmd internal -name '*.go')

migrate:
	go run ./cmd/migrate

test:
	go test ./...

vet:
	go vet ./...

build:
	go build ./cmd/api ./cmd/migrate

web-install:
	cd web && npm install --no-audit --no-fund

web-lint:
	cd web && npm run lint

web-test:
	cd web && npm test

web-build:
	cd web && npm run build

verify: fmt test vet build web-lint web-test web-build

up:
	docker compose up --build

down:
	docker compose down -v
