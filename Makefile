LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
GOOSE_BIN:=$(LOCAL_BIN)/goose



.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	go build -o bin/main cmd/main.go

.PHONY: test
test:
	go test -v ./...



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



.PHONY: install-goose
install-goose:
ifeq ($(wildcard $(GOOSE_BIN)),)
	$(info Downloading goose)
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
GOOSE_BIN:=$(LOCAL_BIN)/goose
endif

.PHONY: migrate-up
migrate-up: install-goose
	$(GOOSE_BIN) -dir migrations postgres "host=localhost user=postgres password=postgres dbname=news-bot-db sslmode=disable" up

.PHONY: migrate-check
migrate-check: install-goose
	$(GOOSE_BIN) -dir migrations postgres "host=localhost user=postgres password=postgres dbname=news-bot-db sslmode=disable" status



.PHONY: local-up
local-up:
	docker-compose -f docker-compose.yaml up -d

.PHONY: local-ps
local-ps:
	docker-compose -f docker-compose.yaml ps

.PHONY: local-down
local-down:
	docker-compose -f docker-compose.yaml down -v