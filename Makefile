WEB_DIR := web
ASSET_DIR := internal/server/assets/dist
BINARY := bin/craftlane

.PHONY: web-build web-test test test-race integration build fmt vet lint e2e

web-build:
	cd $(WEB_DIR) && pnpm build
	@test -f $(ASSET_DIR)/index.html

web-test:
	cd $(WEB_DIR) && pnpm test

test: web-build web-test
	go test ./...

test-race: web-build
	go test -race ./...

integration: web-build
	go test ./tests/integration

build: web-build
	go build -o $(BINARY) ./cmd/craftlane

fmt:
	go fmt ./...

vet: web-build
	go vet ./...

lint:
	cd $(WEB_DIR) && pnpm lint && pnpm typecheck

e2e: build
	cd $(WEB_DIR) && pnpm exec playwright test --config ../playwright.config.ts
