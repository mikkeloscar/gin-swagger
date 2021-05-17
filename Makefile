.PHONY: build test generate clean

BINARY       ?= gin-swagger
SOURCES      = $(shell find . -name '*.go')
GOPKGS       = $(shell go list ./...)
BUILD_FLAGS  ?= -v
GO           ?= go
LDFLAGS      ?= -X main.version=$(VERSION) -w -s

default: build

clean:
	@rm -f $(BINARY)

test:
	$(GO) test -v $(GOPKGS)

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 $(GO) build -o $(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)"
