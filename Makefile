PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=NewApp \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=appd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=appcli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 
BUILD_FLAGS := -ldflags '$(ldflags)'

export GO111MODULE := on

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iovnsd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iovnscli

build: go.sum
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o ./cmd/iovnsd -mod=readonly $(BUILD_FLAGS) ./cmd/iovnsd
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o ./cmd/iovnscli -mod=readonly $(BUILD_FLAGS) ./cmd/iovnscli

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