GRADLE ?= gradle

.PHONY: help fmt test test-unit vet build web-install web-build android-lint android-test android-build extension-check extension-test dev db-up db-down migrate backup clean

help:
	@echo "StockPilot development commands"
	@echo "  make fmt              Format Go code"
	@echo "  make test             Run backend tests"
	@echo "  make vet              Run Go vet"
	@echo "  make build            Build backend"
	@echo "  make web-build        Build frontend"
	@echo "  make android-lint     Run Android lint"
	@echo "  make android-test     Run Android unit tests"
	@echo "  make android-build    Build Android debug APK"
	@echo "  make extension-check  Validate extension sources"
	@echo "  make extension-test   Run extension unit tests"
	@echo "  make db-up            Start PostgreSQL with Docker Compose"
	@echo "  make migrate          Apply SQL migrations"

fmt:
	gofmt -w $$(find cmd internal -name '*.go')

test:
	go test ./...

test-unit:
	go test ./internal/domain ./internal/idgen ./internal/config ./internal/httpapi

vet:
	go vet ./...

build:
	CGO_ENABLED=0 go build -trimpath -o bin/stockpilot ./cmd/server

web-install:
	cd web && npm install --no-audit --no-fund

web-build:
	cd web && npm run build

android-lint:
	cd android && $(GRADLE) :app:lintDebug

android-test:
	cd android && $(GRADLE) :app:testDebugUnitTest

android-build:
	cd android && $(GRADLE) :app:assembleDebug

extension-check:
	cd extension && npm run check

extension-test:
	cd extension && npm test

dev:
	go run ./cmd/server

db-up:
	docker compose up -d db

db-down:
	docker compose down

migrate:
	@for f in migrations/*.up.sql; do \
		echo "Applying $$f"; \
		docker compose exec -T db psql -U stockpilot -d stockpilot -v ON_ERROR_STOP=1 < "$$f"; \
	done

backup:
	./scripts/backup.sh

clean:
	rm -rf bin web/dist coverage android/.gradle android/app/build
