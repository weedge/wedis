
.PHONY: help
help:
	@echo "help:"
	@echo "use [ make build] cmd to exec go build on local"
	@echo "use [ make wire-gen ] cmd to generate server by using injector"


# guides https://github.com/google/wire/blob/main/docs/guide.md
.PHONY: wire-gen
wire-gen:
	@go install github.com/google/wire/cmd/wire@latest
	@cd ./internal/srv && wire && cd -


# Module name.
NAME := wedis
# Project main package location.
CMD_DIR := ./cmd
# Project output directory.
OUTPUT_DIR := ./bin
# Build directory.
BUILD_DIR := ./build
# Current version of the project.
GOCOMMON     := $(shell if [ ! -f go.mod ]; then echo $(ROOT)/vendor/; fi)github.com/weedge/craftsman/cloudwego/payment
#VERSION      ?= latest
VERSION      ?= $(shell git describe --tags --always --dirty)
BRANCH       ?= $(shell git branch | grep \* | cut -d ' ' -f2)
GITCOMMIT    ?= $(shell git rev-parse HEAD)
GITTREESTATE ?= $(if $(shell git status --porcelain),dirty,clean)
BUILDDATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
appVersion   ?= $(VERSION)

.PHONY: build build-local
build: build-local
build-local:
	@go build -v -o $(OUTPUT_DIR)/$(NAME)                                  \
	  -ldflags "-s -w -X $(GOCOMMON)/version.module=$(NAME)                \
	    -X $(GOCOMMON)/version.version=$(VERSION)                          \
	    -X $(GOCOMMON)/version.branch=$(BRANCH)                            \
	    -X $(GOCOMMON)/version.gitCommit=$(GITCOMMIT)                      \
	    -X $(GOCOMMON)/version.gitTreeState=$(GITTREESTATE)                \
	    -X $(GOCOMMON)/version.buildDate=$(BUILDDATE)"                     \
	  $(CMD_DIR) && chmod +x $(OUTPUT_DIR)/$(NAME);

.DEFAULT_GOAL=help