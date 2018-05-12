.PHONY: build test generate clean

BINARY       ?= gin-swagger
SOURCES      = $(shell find . -name '*.go')
GOPKGS       = $(shell go list ./...)
TEMPLATES    = $(shell find templates/ -type f -name '*.gotmpl')
BUILD_FLAGS  ?= -v
GO           ?= go
LDFLAGS      ?= -X main.version=$(VERSION) -w -s

default: build

clean:
	@rm $(BINARY)

test:
	$(GO) test -v $(GOPKGS)

generate: bindata.go

bindata.go: $(TEMPLATES)
	$(GO) generate .

build: $(BINARY)

$(BINARY): bindata.go $(SOURCES)
	CGO_ENABLED=0 $(GO) build -o $(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)"
