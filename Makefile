# vim: set ts=2:
GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build

BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

GOTOOLS := tortool

CONTEXT = "https://github.com/filvarga/tortools.git\#$(BRANCH):docker"

LDFLAGS=-ldflags "-X main.context=$(CONTEXT)" 

all: install

tools:
	@$(foreach TOOL,$(GOTOOLS),$(GOBUILD) $(LDFLAGS) -o build/$(TOOL) ./cmd/$(TOOL);)

install:
	@$(foreach TOOL,$(GOTOOLS),$(GOINSTALL) $(LDFLAGS) ./cmd/$(TOOL);)
