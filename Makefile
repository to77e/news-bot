include .env
export

SERVICE_NAME=news-fetching-bot
SERVICE_PATH=to77e/news-fetching-bot

LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
GOOSE_BIN:=$(LOCAL_BIN)/goose



.PHONY: run
run:
	go run cmd/news-fetching-bot/main.go

.PHONY: build
build:
	go build \
	-ldflags=" \
		-X 'github.com/$(SERVICE_PATH)/internal/config.version=`git tag --sort=-version:refname | head -n 1`\
	" \
	-o bin/news-fetching-bot cmd/news-fetching-bot/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: deps
deps:
	go mod tidy
	go mod download


# Commands for linting
.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

.PHONY: lint
lint: install-lint
	$(info Running lint...)
	$(GOLANGCI_BIN) run --config=.golangci.yaml ./...


# Commands for working with migrations
.PHONY: install-goose
install-goose:
ifeq ($(wildcard $(GOOSE_BIN)),)
	$(info Downloading goose)
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
GOOSE_BIN:=$(LOCAL_BIN)/goose
endif

.PHONY: migrate-up
migrate-up: install-goose
	$(GOOSE_BIN) -dir migrations postgres "host=postgres user=postgres password=postgres dbname=news-fetching-bot-db sslmode=disable" up

.PHONY: migrate-down
migrate-down: install-goose
	$(GOOSE_BIN) -dir migrations postgres "host=postgres user=postgres password=postgres dbname=news-fetching-bot-db sslmode=disable" down

.PHONY: migrate-status
migrate-status: install-goose
	$(GOOSE_BIN) -dir migrations postgres "host=postgres user=postgres password=postgres dbname=news-fetching-bot-db sslmode=disable" status


# Commands for local deployment
.PHONY: local-up
local-up:
	docker-compose -f docker-compose.yaml up -d

.PHONY: local-ps
local-ps:
	docker-compose -f docker-compose.yaml ps

.PHONY: local-down
local-down:
	docker-compose -f docker-compose.yaml down -v

# Commands for deploying to remote server
.PHONY: dev-up
dev-up:
	DOCKER_HOST=$(NEWS_FETCHING_BOT_REMOTE_HOST) docker-compose -f docker-compose.yaml up -d

.PHONY: dev-down
dev-down:
	DOCKER_HOST=$(NEWS_FETCHING_BOT_REMOTE_HOST) docker-compose -f docker-compose.yaml down -v -rmi $(SERVICE_NAME)