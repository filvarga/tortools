# vim: set ts=2:
GO_CMD=go
GO_INSTALL=$(GO_CMD) install
GO_BUILD=$(GO_CMD) build

BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

GO_TOOLS := tortool

#CONTEXT = "https://github.com/filvarga/tortools.git\#$(BRANCH):docker"
#LDFLAGS=-ldflags "-X main.context=$(CONTEXT)"

all: install

tools:
	@$(foreach TOOL,$(GO_TOOLS),$(GO_BUILD) $(LDFLAGS) -o build/$(TOOL) ./cmd/$(TOOL);)

install:
	@$(foreach TOOL,$(GO_TOOLS),$(GO_INSTALL) $(LDFLAGS) ./cmd/$(TOOL);)
