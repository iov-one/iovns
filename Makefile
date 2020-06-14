PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

LEDGER_ENABLED ?= true

# process build tags
build_tags =
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=iovns \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=iovnsd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=iovnscli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)' -tags '$(build_tags)'
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

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

test:
	@# iovnscli binary is required for it tests
	go install -mod=readonly $(BUILD_FLAGS)  ./cmd/iovnscli

	go vet -mod=readonly  ./...
	go test -mod=readonly -race ./...

# Test fast
tf:
	go test -short ./...

update-swagger-docs: statik
	$(BINDIR)/statik -src=swagger-ui/swagger-ui -dest=swagger-ui -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
    	echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
.PHONY: update-swagger-docs
