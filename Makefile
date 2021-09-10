SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
TEST_TIMEOUT?=15m
TEST_PARALLEL?=2
DOCKER_BUILDKIT?=1
export DOCKER_BUILDKIT

export GO111MODULE := on

# Install all the build and lint dependencies
setup:
	go mod tidy
	git config core.hooksPath .githooks
.PHONY: setup

test:
	go test $(TEST_OPTIONS) -p $(TEST_PARALLEL) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=$(TEST_TIMEOUT)
.PHONY: test

cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

fmt:
	gofumpt -w backend
	gofumpt -w bin/crypt
	gofumpt -w config
	gofumpt -w encoding
.PHONY: fmt


ci: test
.PHONY: ci

build:
	go build -o crypt ./bin/crypt
.PHONY: build

install:
	go install ./bin/crypt

deps:
	go get -u github.com/bketelsen/crypt/bin/crypt
	go mod tidy
	go mod verify
	go mod vendor
.PHONY: deps

todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude-dir=bin \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

.DEFAULT_GOAL := build
