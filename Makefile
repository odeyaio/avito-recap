.PHONY: install generate generate-go generate-web generate-check format lint lint-go lint-web test test-go test-web build build-go build-web check compose-up compose-down migrate-up migrate-down catalog-import

install:
	pnpm install --frozen-lockfile
	go -C apps/backend mod download

generate: generate-go generate-web

generate-go:
	go -C apps/backend generate ./...

generate-web:
	pnpm --filter @avito-recap/frontend generate

generate-check: generate
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then git diff --exit-code; else printf '%s\n' "Skipping generation diff check: not a git repository"; fi

format:
	go -C apps/backend fmt ./...
	pnpm format

lint: lint-go lint-web

lint-go:
	docker run --rm -v "$(CURDIR):/workspace" -w /workspace/apps/backend golangci/golangci-lint:v2.12.2 golangci-lint run --config /workspace/.golangci.yaml ./...

lint-web:
	pnpm lint
	pnpm format:check

test: test-go test-web

test-go:
	go -C apps/backend test -race ./...

test-web:
	pnpm test

build: build-go build-web

build-go:
	mkdir -p dist
	go -C apps/backend build -o ../../dist/avito-recap ./cmd/avito-recap

build-web:
	pnpm build

check: generate-check lint test build

compose-up:
	docker compose up --build

compose-down:
	docker compose down

migrate-up:
	docker compose run --rm migrate

migrate-down:
	docker compose run --rm migrate down 1

catalog-import:
	go -C apps/backend run ./cmd/catalog-import \
		-achievements-file ../../catalog/achievements.yaml \
		-behaviors-file ../../catalog/behaviors.yaml
