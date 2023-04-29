PROJECT_VERSION := 0.1.0
DOCKER_REPO     := synfinatic
PROJECT_NAME    := gorawsocket

DIST_DIR ?= dist/
GOOS ?= $(shell uname -s | tr "[:upper:]" "[:lower:]")
ARCH ?= $(shell uname -m)
ifeq ($(ARCH),x86_64)
GOARCH             := amd64
else
GOARCH             := $(ARCH)  # no idea if this works for other platforms....
endif

ifneq ($(BREW_INSTALL),1)
PROJECT_TAG               := $(shell git describe --tags 2>/dev/null $(git rev-list --tags --max-count=1))
PROJECT_COMMIT            := $(shell git rev-parse HEAD || echo "")
PROJECT_DELTA             := $(shell DELTA_LINES=$$(git diff | wc -l); if [ $${DELTA_LINES} -ne 0 ]; then echo $${DELTA_LINES} ; else echo "''" ; fi)
else
PROJECT_TAG               := Homebrew
endif

BUILDINFOSDET ?=
PROGRAM_ARGS ?=
ifeq ($(PROJECT_TAG),)
PROJECT_TAG               := NO-TAG
endif
ifeq ($(PROJECT_COMMIT),)
PROJECT_COMMIT            := NO-CommitID
endif
ifeq ($(PROJECT_DELTA),)
PROJECT_DELTA             :=
endif

VERSION_PKG               := $(shell echo $(PROJECT_VERSION) | sed 's/^v//g')
LICENSE                   := 3 Clause BSD
URL                       := https://github.com/$(DOCKER_REPO)/$(PROJECT_NAME)
DESCRIPTION               := gorawsocket
BUILDINFOS                ?= $(shell date +%FT%T%z)$(BUILDINFOSDET)
LDFLAGS                   := -X "main.Version=$(PROJECT_VERSION)" -X "main.Delta=$(PROJECT_DELTA)" -X "main.Buildinfos=$(BUILDINFOS)" -X "main.Tag=$(PROJECT_TAG)" -X "main.CommitID=$(PROJECT_COMMIT)"
OUTPUT_NAME               := $(DIST_DIR)$(PROJECT_NAME)-$(PROJECT_VERSION)  # default for current platform

ALL: $(DIST_DIR)$(PROJECT_NAME) ## Build binary for this platform

include help.mk  # place after ALL target and before all other targets

$(DIST_DIR)$(PROJECT_NAME):	$(wildcard */*.go) .prepare
	go build -ldflags='$(LDFLAGS)' -o $(DIST_DIR)$(PROJECT_NAME) ./cmd/rawsocktest//...
	@echo "Created: $(DIST_DIR)$(PROJECT_NAME)"

INSTALL_PREFIX ?= /usr/local

install: $(DIST_DIR)$(PROJECT_NAME)  ## install binary in $INSTALL_PREFIX
	install -d $(INSTALL_PREFIX)/bin
	install -c $(DIST_DIR)$(PROJECT_NAME) $(INSTALL_PREFIX)/bin

uninstall:  ## Uninstall binary from $INSTALL_PREFIX
	rm $(INSTALL_PREFIX)/bin/$(PROJECT_NAME)


.PHONY: shasum
shasum:
	@which shasum >/dev/null || (echo "Missing 'shasum' binary" ; exit 1)
	@echo "foo" | shasum -a 256 >/dev/null || (echo "'shasum' does not support: -a 256"; exit 1)


tags: $(wildcard cmd/rawsocktest/*.go) $(wildcard gorawsocket/*.go)  ## Create tags file for vim, etc
	@echo Make sure you have Universal Ctags installed: https://github.com/universal-ctags/ctags
	ctags --recurse=yes --exclude=.git --exclude=\*.sw?  --exclude=dist --exclude=docs


.validate-release: ALL
	@TAG=$$(./$(DIST_DIR)$(PROJECT_NAME) version 2>/dev/null | grep '(v$(PROJECT_VERSION))'); \
		if test -z "$$TAG"; then \
		echo "Build tag from does not match PROJECT_VERSION=v$(PROJECT_VERSION) in Makefile:" ; \
		./$(DIST_DIR)$(PROJECT_NAME) version 2>/dev/null | grep built ; \
		exit 1 ; \
	fi


clean-all: clean ## clean _everything_

clean: ## Remove all binaries in dist
	rm -rf dist/*

clean-go: ## Clean Go cache
	go clean -i -r -cache -modcache

go-get:  ## Get our go modules
	go get -v all

.prepare: $(DIST_DIR)

.PHONY: build-race
build-race: .prepare ## Build race detection binary
	go build -race -ldflags='$(LDFLAGS)' -o $(OUTPUT_NAME) ./cmd/rawsocktest/...

debug: .prepare ## Run debug in dlv
	dlv debug ./cmd/rawsocktest/

.PHONY: unittest
unittest: ## Run go unit tests
	go test -race -covermode=atomic -coverprofile=coverage.out  ./...

.PHONY: test-race
test-race: ## Run `go test -race` on the code
	@echo checking code for races...
	go test -race ./...

.PHONY: vet
vet: ## Run `go vet` on the code
	@echo checking code is vetted...
	for x in $(shell go list ./...); do echo $$x ; go vet $$x ; done

test: vet unittest lint test-homebrew ## Run important tests

precheck: test test-fmt test-tidy ## Run all tests that happen in a PR

# run everything but `lint` because that runs via it's own workflow
.build-tests: vet unittest test-tidy test-fmt

$(DIST_DIR):
	@if test ! -d $(DIST_DIR); then mkdir -p $(DIST_DIR) ; fi

.PHONY: fmt
fmt: ## Format Go code
	@gofmt -s -w */*.go */*/*.go

.PHONY: test-fmt
test-fmt: fmt ## Test to make sure code if formatted correctly
	@if test `git diff cmd/rawsocktest | wc -l` -gt 0; then \
	    echo "Code changes detected when running 'go fmt':" ; \
	    git diff -Xfiles ; \
	    exit -1 ; \
	fi

.PHONY: test-tidy
test-tidy:  ## Test to make sure go.mod is tidy
	@go mod tidy
	@if test `git diff go.mod | wc -l` -gt 0; then \
	    echo "Need to run 'go mod tidy' to clean up go.mod" ; \
	    exit -1 ; \
	fi

lint:  ## Run golangci-lint
	golangci-lint run


$(OUTPUT_NAME): $(wildcard */*.go) .prepare
	go build -ldflags='$(LDFLAGS)' -o $(OUTPUT_NAME) ./cmd/rawsocktest/...

.PHONY: loc
loc:  ## Print LOC stats
	wc -l $$(find . -name "*.go")
