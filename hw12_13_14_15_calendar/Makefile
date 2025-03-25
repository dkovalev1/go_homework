BIN:=./bin/calendar
BIN_SENDER:=./bin/sender
BIN_SCHEDULER:=./bin/scheduler
GINKGO_FLAGS := -r --keep-going --randomize-suites --randomize-all --no-color

DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: BIN BIN_SENDER BIN_SCHEDULER
	echo "building..."

BIN:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

BIN_SENDER:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

BIN_SCHEDULER:	
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

compose-build:
	docker compose -f deployments/docker-compose.yaml build

compose-up:
	docker compose -f deployments/docker-compose.yaml up -d

compose-down:
	docker compose -f deployments/docker-compose.yaml down

compose-test-up:
	docker compose -f deployments/docker-compose-test.yaml up -d

compose-test-down:
	docker compose -f deployments/docker-compose-test.yaml down

up: compose-build compose-up

down: compose-down

test: compose-test-up
	go test -race ./internal/... ./pkg/...
	make compose-test-down

integration-tests: install-ginkgo compose-up 
	ginkgo $(GINKGO_FLAGS) --timeout=5m --poll-progress-after=0s integration 2>&1

end_to_end : install-ginkgo build compose-test-up 
	ginkgo $(GINKGO_FLAGS) --timeout=5m --poll-progress-after=0s end_to_end 2>&1

install-ginkgo:
	go install github.com/onsi/ginkgo/v2/ginkgo

install-lint-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
#	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

generate:
	cd api && protoc --go_out=. --go-grpc_out=. EventService.proto

.PHONY: build run build-img run-img version test lint end_to_end
