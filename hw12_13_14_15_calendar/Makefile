build: test build-calendar build-scheduler build-sender

build-calendar:
	go build -o ./bin/calendar ./cmd/calendar

build-scheduler:
	go build -o ./bin/calendar_scheduler ./cmd/scheduler

build-sender:
	go build -o ./bin/calendar_sender ./cmd/sender

test: lint
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build test lint

up:
	sh ./scripts/deploy.sh up

down:
	sh ./scripts/deploy.sh down

integration-tests:
	sh ./scripts/tests-run.sh

migrate-up:
	goose -dir migrations postgres $(CAL_STORAGE_DSN) up

migrate-down:
	goose -dir migrations postgres $(CAL_STORAGE_DSN) down

generate:
	go generate ./internal/server/grpcsrv/grpc_server.go
