PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=iovns \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=iovnsd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=iovnscli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 
BUILD_FLAGS := -ldflags '$(ldflags)'
BINDIR ?= $(GOPATH)/bin

export GO111MODULE := on

include contrib/devtools/Makefile

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iovnsd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iovnscli

build: go.sum
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o ./build/iovnsd -mod=readonly $(BUILD_FLAGS) ./cmd/iovnsd
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o ./build/iovnscli -mod=readonly $(BUILD_FLAGS) ./cmd/iovnscli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

test:
	go test -mod=readonly -race ./...

update-swagger-docs: statik
	$(BINDIR)/statik -src=swagger-ui/swagger-ui -dest=swagger-ui -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
    	echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
.PHONY: update-swagger-docs
